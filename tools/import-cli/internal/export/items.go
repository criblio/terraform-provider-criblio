// Package export converts discovery results into generator ResourceItems.
package export

import (
	"context"
	"errors"
	"fmt"
	"strings"

	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// convertOneResource fetches a single resource via the converter, builds HCL attrs, and returns a ResourceItem or skip message.
func convertOneResource(ctx context.Context, client *importclient.Client, r discovery.Result, e registry.Entry, idMap map[string]string, groupFilter []string, groupIDs []string, excludeDefaults bool, includeOverride IncludeOverride, out *ExportResult) (item *generator.ResourceItem, skipMsg string) {
	if skipExportForGroupFilter(r.TypeName, idMap, groupFilter, groupIDs) {
		return nil, ""
	}
	if skipResourceByID(r.TypeName, idMap) {
		return nil, fmt.Sprintf("%s %v: skipped by config", r.TypeName, idMap)
	}
	if defaultLookupFileByID(r.TypeName, idMap, includeOverride) {
		return nil, fmt.Sprintf("%s %v: built-in default lookup (skip export)", r.TypeName, idMap)
	}
	requestParams := toRequestParams(idMap)
	model, convErr := converter.Convert(ctx, client, e, requestParams)
	if convErr != nil {
		return nil, fmt.Sprintf("%s %v: %s", r.TypeName, idMap, sanitizeConvertError(convErr))
	}
	opts := hclOptionsForType(r.TypeName, e)
	attrs, attrErr := hcl.ModelToValue(model, opts)
	if attrErr != nil {
		return nil, fmt.Sprintf("%s %v: model to value: %s", r.TypeName, idMap, sanitizeConvertError(attrErr))
	}
	if flattenItemsToTopLevelTypes[r.TypeName] {
		flattenItemsListToTopLevel(attrs)
	}
	// criblio_pack: exports is an install-time parameter (RequiresReplaceIfConfigured). The API does
	// not return it on GET so state has exports=null after import, while items[0].exports (flattened
	// above) can be ["*"], causing destroy+recreate on every apply. Drop it from HCL config.
	if r.TypeName == "criblio_pack" {
		delete(attrs, "exports")
	}
	if (r.TypeName == "criblio_routes" || r.TypeName == "criblio_pack_routes") && attrs["routes"].Kind != hcl.KindNull {
		attrs["routes"] = normalizeRoutesForExport(attrs["routes"])
	}
	if r.TypeName == "criblio_appscope_config" {
		custom.ApplyAppscopeConfigDefaults(attrs)
	}
	if r.TypeName == "criblio_project" {
		custom.ApplyProjectDefaults(attrs)
	}
	if r.TypeName == "criblio_subscription" {
		custom.ApplySubscriptionDefaults(attrs)
	}
	if r.TypeName == "criblio_pack" {
		custom.ApplyPackDefaults(attrs)
	}
	if r.TypeName == "criblio_pack_vars" {
		custom.ApplyPackVarsDefaults(attrs)
	}
	if e.OneOf != nil {
		if oneOfErr := addOneOfBlockFromFirstItem(model, attrs, e.OneOf); oneOfErr != nil {
			if errors.Is(oneOfErr, ErrUnsupportedOneOfType) {
				return nil, fmt.Sprintf("%s %v: oneOf type unsupported by provider", r.TypeName, idMap)
			}
			return nil, fmt.Sprintf("%s %v: oneOf: %s", r.TypeName, idMap, sanitizeConvertError(oneOfErr))
		}
	}
	// criblio_pack_destination: model (DestinationResourceModel) has no Items field, so addOneOfBlockFromFirstItem
	// does nothing. Emit the oneOf block from the API response when present (stored by converter).
	if r.TypeName == "criblio_pack_destination" && !attrsHasOutputBlock(attrs) {
		addPackDestinationOneOfFromStoredItem(idMap, attrs, e.OneOf)
	}
	if flattenItemsToAttrsTypes[r.TypeName] {
		if flatErr := flattenFirstItemToAttrs(model, attrs, "items"); flatErr != nil {
			return nil, fmt.Sprintf("%s %v: flatten: %s", r.TypeName, idMap, sanitizeConvertError(flatErr))
		}
	}
	if skipResourceWhenLibCribl(attrs) {
		return nil, fmt.Sprintf("%s %v: lib is cribl (built-in, skip export)", r.TypeName, idMap)
	}
	if isLookupFileType(r.TypeName) && criblDefaultTag(attrs) {
		return nil, fmt.Sprintf("%s %v: cribl:default tag (built-in, skip export)", r.TypeName, idMap)
	}
	if isLookupFileType(r.TypeName) && DefaultResource(r.TypeName, idMap, attrs, includeOverride) {
		return nil, fmt.Sprintf("%s %v: built-in default lookup (skip export)", r.TypeName, idMap)
	}
	if excludeDefaults && DefaultResource(r.TypeName, idMap, attrs, includeOverride) {
		out.DefaultsSkipped++
		return nil, fmt.Sprintf("%s %v: built-in default (--exclude-defaults)", r.TypeName, idMap)
	}
	if r.TypeName == "criblio_search_datatype_ruleset" || r.TypeName == "criblio_search_dataset_ruleset" {
		if injErr := injectSearchRulesetRulesForExport(r.TypeName, model, attrs, e); injErr != nil {
			return nil, fmt.Sprintf("%s %v: search ruleset rules: %s", r.TypeName, idMap, sanitizeConvertError(injErr))
		}
	}
	filterAttrsBySchema(attrs, e.ModelTypeName)
	if r.TypeName == "criblio_pack_destination" && idMap["pack"] != "" {
		attrs["pack"] = hcl.Value{Kind: hcl.KindString, String: idMap["pack"]}
	}
	buildIDMap := importIDMapForType(r.TypeName, idMap)
	importID, idErr := generator.BuildImportID(e.ImportIDFormat, buildIDMap)
	if idErr != nil {
		return nil, fmt.Sprintf("%s %v: import ID: %s", r.TypeName, idMap, sanitizeConvertError(idErr))
	}
	name := generator.StableResourceNameFromMap(e.TypeName, idMap)
	ensureNotificationTargetSecretPlaceholders(r.TypeName, attrs, name)
	_ = hcl.ReplaceSecretValuesWithVariableRefs(attrs, name)
	if r.TypeName == "criblio_secret" {
		attrs["value"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.SecretValueVariableName(name)}
	}
	if r.TypeName == "criblio_certificate" {
		attrs["cert"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.CertificateCertVariableName(name)}
		attrs["priv_key"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.CertificatePrivKeyVariableName(name)}
	}
	files, hasLookupContent, contentErr := lookupContentAsset(ctx, client, r.TypeName, attrs, name, idMap)
	if contentErr != nil {
		return nil, fmt.Sprintf("%s %v: lookup content: %s", r.TypeName, idMap, sanitizeConvertError(contentErr))
	}
	if !hasLookupContent {
		return nil, fmt.Sprintf("%s %v: lookup content missing from API response", r.TypeName, idMap)
	}
	it := generator.ResourceItem{
		TypeName: e.TypeName,
		Name:     name,
		Attrs:    attrs,
		ImportID: importID,
		GroupID:  groupIDForOutput(e.TypeName, groupIDFromIDMap(idMap)),
		Files:    files,
	}
	it.LifecycleIgnoreChanges = lifecycleIgnoreChangesForConvertedResource(r.TypeName, attrs)
	return &it, ""
}

// normalizeRoutesForExport removes route values that are valid in state but invalid in config.
func normalizeRoutesForExport(routesVal hcl.Value) hcl.Value {
	if routesVal.Kind != hcl.KindList {
		return routesVal
	}
	for i := range routesVal.List {
		if routesVal.List[i].Kind != hcl.KindMap || routesVal.List[i].Map == nil {
			continue
		}
		route := routesVal.List[i].Map
		if v, ok := route["description"]; ok && v.Kind == hcl.KindString && v.String == "" {
			route["description"] = hcl.Value{Kind: hcl.KindNull}
		}
		if v, ok := route["clones"]; ok && v.Kind == hcl.KindList {
			route["clones"] = removeNullListItems(v)
			if len(route["clones"].List) == 0 {
				delete(route, "clones")
			}
		}
	}
	return routesVal
}

func removeNullListItems(v hcl.Value) hcl.Value {
	if v.Kind != hcl.KindList {
		return v
	}
	list := make([]hcl.Value, 0, len(v.List))
	for _, item := range v.List {
		if item.Kind == hcl.KindNull {
			continue
		}
		list = append(list, item)
	}
	v.List = list
	return v
}

func lifecycleIgnoreChangesForConvertedResource(typeName string, attrs map[string]hcl.Value) []string {
	if typeName == "criblio_appscope_config" {
		// Deeply nested Optional+Computed fields (cacertpath, buffer, allowbinary, headers, etc.)
		// are set to null/empty by the provider Read but absent from HCL, causing perpetual drift.
		return []string{"config"}
	}
	if typeName == "criblio_secret" {
		// value is sensitive; description, tags may be computed.
		return []string{"description", "tags", "value"}
	}
	if typeName == "criblio_pack_destination" {
		// OneOf block structure may differ slightly from provider Read (e.g. optional nulls).
		// Ignore the output block we emit to suppress drift.
		for k := range attrs {
			if strings.HasPrefix(k, "output_") {
				return []string{k}
			}
		}
	}
	return nil
}

// appendResourceItemFromModel builds HCL attrs and appends a ResourceItem to out.Items (used for criblio_group and shared conversion path).
func appendResourceItemFromModel(out *ExportResult, typeName string, e registry.Entry, idMap map[string]string, model interface{}, groupFilter []string, groupIDs []string, excludeDefaults bool, includeOverride IncludeOverride) error {
	if skipResourceByID(typeName, idMap) {
		return nil
	}
	if skipExportForGroupFilter(typeName, idMap, groupFilter, groupIDs) {
		return nil
	}
	opts := hclOptionsForType(typeName, e)
	attrs, attrErr := hcl.ModelToValue(model, opts)
	if attrErr != nil {
		return fmt.Errorf("model to value: %w", attrErr)
	}
	if flattenItemsToTopLevelTypes[typeName] {
		flattenItemsListToTopLevel(attrs)
	}
	if typeName == "criblio_pack" {
		delete(attrs, "exports")
	}
	// criblio_group requires product (stream|edge); provider model may not set it from API. Use product from idMap (we set it in ListGroupIdentifiersAndItems).
	// Emit streamtags = [] explicitly when empty so config matches state (API returns []).
	// ApplyGroupDefaults populates description, is_fleet, on_prem, tags, type, etc. from model to avoid plan drift.
	if typeName == "criblio_group" {
		if idMap["product"] != "" {
			attrs["product"] = hcl.Value{Kind: hcl.KindString, String: idMap["product"]}
		}
		if _, has := attrs["streamtags"]; !has || attrs["streamtags"].IsNull() {
			attrs["streamtags"] = hcl.Value{Kind: hcl.KindList, List: []hcl.Value{}}
		}
		custom.ApplyGroupDefaults(attrs, model)
	}
	if typeName == "criblio_appscope_config" {
		custom.ApplyAppscopeConfigDefaults(attrs)
	}
	if typeName == "criblio_project" {
		custom.ApplyProjectDefaults(attrs)
	}
	if typeName == "criblio_subscription" {
		custom.ApplySubscriptionDefaults(attrs)
	}
	if typeName == "criblio_pack" {
		custom.ApplyPackDefaults(attrs)
	}
	if typeName == "criblio_pack_vars" {
		custom.ApplyPackVarsDefaults(attrs)
	}
	if e.OneOf != nil {
		if err := addOneOfBlockFromFirstItem(model, attrs, e.OneOf); err != nil {
			return fmt.Errorf("oneOf: %w", err)
		}
	}
	if flattenItemsToAttrsTypes[typeName] {
		if err := flattenFirstItemToAttrs(model, attrs, "items"); err != nil {
			return fmt.Errorf("flatten: %w", err)
		}
	}
	if typeName == "criblio_search_datatype_ruleset" || typeName == "criblio_search_dataset_ruleset" {
		if err := injectSearchRulesetRulesForExport(typeName, model, attrs, e); err != nil {
			return err
		}
	}
	filterAttrsBySchema(attrs, e.ModelTypeName)
	if skipResourceWhenLibCribl(attrs) {
		return ErrSkipResourceLibCribl
	}
	if excludeDefaults && DefaultResource(typeName, idMap, attrs, includeOverride) {
		out.DefaultsSkipped++
		return nil
	}
	buildIDMap := importIDMapForType(typeName, idMap)
	importID, idErr := generator.BuildImportID(e.ImportIDFormat, buildIDMap)
	if idErr != nil {
		return fmt.Errorf("import ID: %w", idErr)
	}
	// For criblio_group, name from group_id only so resource name stays "group_default" not "group_default_stream".
	nameMap := idMap
	if typeName == "criblio_group" {
		nameMap = make(map[string]string, len(idMap))
		for k, v := range idMap {
			if k != "product" {
				nameMap[k] = v
			}
		}
	}
	name := generator.StableResourceNameFromMap(e.TypeName, nameMap)
	ensureNotificationTargetSecretPlaceholders(typeName, attrs, name)
	_ = hcl.ReplaceSecretValuesWithVariableRefs(attrs, name)
	if typeName == "criblio_secret" {
		attrs["value"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.SecretValueVariableName(name)}
	}
	it := generator.ResourceItem{
		TypeName: e.TypeName,
		Name:     name,
		Attrs:    attrs,
		ImportID: importID,
		GroupID:  groupIDForOutput(e.TypeName, groupIDFromIDMap(idMap)),
	}
	// criblio_group: API may not return product; state has null. Add ignore_changes to suppress drift.
	if typeName == "criblio_group" {
		it.LifecycleIgnoreChanges = []string{"product"}
	}
	out.Items = append(out.Items, it)
	return nil
}

func ensureNotificationTargetSecretPlaceholders(typeName string, attrs map[string]hcl.Value, resourceName string) {
	if typeName != "criblio_notification_target" {
		return
	}
	block, ok := attrs["slack_target"]
	if !ok || block.Kind != hcl.KindMap || block.Map == nil {
		return
	}
	block.Map["url"] = hcl.Value{
		Kind:    hcl.KindVariableRef,
		VarName: hcl.SensitiveVariableName(resourceName, "slack_target_url"),
	}
	attrs["slack_target"] = block
}

func importIDMapForType(typeName string, idMap map[string]string) map[string]string {
	if typeName == "criblio_key" && idMap["id"] != "" {
		out := copyStringMap(idMap, 1)
		out["key_id"] = idMap["id"]
		return out
	}
	if typeName == "criblio_notification" && idMap["group"] == "" && idMap["group_id"] != "" {
		out := copyStringMap(idMap, 1)
		out["group"] = idMap["group_id"]
		return out
	}
	return idMap
}

func copyStringMap(in map[string]string, extra int) map[string]string {
	out := make(map[string]string, len(in)+extra)
	for k, v := range in {
		out[k] = v
	}
	return out
}
