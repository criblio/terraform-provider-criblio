// Package export converts discovery results into generator ResourceItems.
package export

import (
	"context"
	"errors"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	ptypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// convertOneResource fetches a single resource via the converter, builds HCL attrs, and returns a ResourceItem or skip message.
func convertOneResource(ctx context.Context, client *sdk.CriblIo, r discovery.Result, e registry.Entry, idMap map[string]string) (item *generator.ResourceItem, skipMsg string) {
	if skipResourceByID(r.TypeName, idMap) {
		return nil, fmt.Sprintf("%s %v: skipped by config", r.TypeName, idMap)
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
	// criblio_routes and criblio_pack_routes: inject additional_properties from model into each route.
	// HCL skips it by default (readOnlyAttrs); API returns it, causing drift. Align both resources.
	injectRoutesAdditionalProperties := func(routes []ptypes.RoutesRoute, routesVal hcl.Value) hcl.Value {
		if routesVal.Kind != hcl.KindList || len(routes) > len(routesVal.List) {
			return routesVal
		}
		for i := range routes {
			if i < len(routesVal.List) {
				route := routes[i]
				if !route.AdditionalProperties.IsNull() && !route.AdditionalProperties.IsUnknown() {
					apStr := route.AdditionalProperties.ValueString()
					if apStr != "" && routesVal.List[i].Kind == hcl.KindMap && routesVal.List[i].Map != nil {
						routesVal.List[i].Map["additional_properties"] = hcl.Value{Kind: hcl.KindString, String: apStr}
					}
				}
			}
		}
		return routesVal
	}
	if r.TypeName == "criblio_routes" {
		if pm, ok := model.(*provider.RoutesResourceModel); ok && attrs["routes"].Kind != hcl.KindNull {
			attrs["routes"] = injectRoutesAdditionalProperties(pm.Routes, attrs["routes"])
		}
	}
	if r.TypeName == "criblio_pack_routes" {
		// items is skipped (Computed); populate routes from model.Items[0] since flatten had no items to merge.
		if pm, ok := model.(*provider.PackRoutesResourceModel); ok && len(pm.Items) > 0 {
			routesWrap := struct {
				Routes []ptypes.RoutesRoute `tfsdk:"routes"`
			}{Routes: pm.Items[0].Routes}
			routesAttrs, err := hcl.ModelToValue(&routesWrap, opts)
			if err == nil && routesAttrs["routes"].Kind != hcl.KindNull {
				routesVal := injectRoutesAdditionalProperties(pm.Items[0].Routes, routesAttrs["routes"])
				attrs["routes"] = routesVal
			}
		}
	}
	if r.TypeName == "criblio_appscope_config" {
		custom.ApplyAppscopeConfigDefaults(attrs)
	}
	if r.TypeName == "criblio_project" {
		custom.ApplyProjectDefaults(attrs)
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
	filterAttrsBySchema(attrs, e.ModelTypeName)
	if r.TypeName == "criblio_pack_destination" && idMap["pack"] != "" {
		attrs["pack"] = hcl.Value{Kind: hcl.KindString, String: idMap["pack"]}
	}
	buildIDMap := idMap
	if r.TypeName == "criblio_key" && idMap["id"] != "" {
		buildIDMap = make(map[string]string, len(idMap)+1)
		for k, v := range idMap {
			buildIDMap[k] = v
		}
		buildIDMap["key_id"] = idMap["id"]
	}
	importID, idErr := generator.BuildImportID(e.ImportIDFormat, buildIDMap)
	if idErr != nil {
		return nil, fmt.Sprintf("%s %v: import ID: %s", r.TypeName, idMap, sanitizeConvertError(idErr))
	}
	name := generator.StableResourceNameFromMap(e.TypeName, idMap)
	_ = hcl.ReplaceSecretValuesWithVariableRefs(attrs, name)
	if r.TypeName == "criblio_secret" {
		attrs["value"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.SecretValueVariableName(name)}
	}
	if r.TypeName == "criblio_certificate" {
		attrs["cert"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.CertificateCertVariableName(name)}
		attrs["priv_key"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.CertificatePrivKeyVariableName(name)}
	}
	it := generator.ResourceItem{
		TypeName:  e.TypeName,
		Name:      name,
		Attrs:     attrs,
		ImportID:  importID,
		GroupID:   groupIDForOutput(e.TypeName, groupIDFromIDMap(idMap)),
	}
	if r.TypeName == "criblio_appscope_config" {
		// Deeply nested Optional+Computed fields (cacertpath, buffer, allowbinary, headers, etc.)
		// are set to null/empty by the provider Read but absent from HCL, causing perpetual drift.
		it.LifecycleIgnoreChanges = []string{"config"}
	}
	if r.TypeName == "criblio_certificate" {
		// ca, in_use are computed; cert, priv_key, passphrase are sensitive/read-only from API.
		it.LifecycleIgnoreChanges = []string{"ca", "cert", "in_use", "passphrase", "priv_key"}
	}
	if r.TypeName == "criblio_secret" {
		// value is sensitive; description, tags may be computed.
		it.LifecycleIgnoreChanges = []string{"description", "tags", "value"}
	}
	return &it, ""
}

// appendResourceItemFromModel builds HCL attrs and appends a ResourceItem to out.Items (used for criblio_group and shared conversion path).
func appendResourceItemFromModel(out *ExportResult, typeName string, e registry.Entry, idMap map[string]string, model interface{}) error {
	opts := hclOptionsForType(typeName, e)
	attrs, attrErr := hcl.ModelToValue(model, opts)
	if attrErr != nil {
		return fmt.Errorf("model to value: %w", attrErr)
	}
	if flattenItemsToTopLevelTypes[typeName] {
		flattenItemsListToTopLevel(attrs)
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
	filterAttrsBySchema(attrs, e.ModelTypeName)
	if skipResourceWhenLibCribl(attrs) {
		return ErrSkipResourceLibCribl
	}
	importID, idErr := generator.BuildImportID(e.ImportIDFormat, idMap)
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
	_ = hcl.ReplaceSecretValuesWithVariableRefs(attrs, name)
	if typeName == "criblio_secret" {
		attrs["value"] = hcl.Value{Kind: hcl.KindVariableRef, VarName: hcl.SecretValueVariableName(name)}
	}
	it := generator.ResourceItem{
		TypeName:  e.TypeName,
		Name:      name,
		Attrs:     attrs,
		ImportID:  importID,
		GroupID:   groupIDForOutput(e.TypeName, groupIDFromIDMap(idMap)),
	}
	// criblio_group: API may not return product; state has null. Add ignore_changes to suppress drift.
	if typeName == "criblio_group" {
		it.LifecycleIgnoreChanges = []string{"product"}
	}
	out.Items = append(out.Items, it)
	return nil
}
