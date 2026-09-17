// Package export converts discovery results into generator ResourceItems.
package export

import (
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// flattenFirstItemToAttrs takes the first element from the model's list (e.g. Items from GetByID
// response) and merges it as top-level attributes (description, type, value, etc.). Used for
// resources like global_var that use GetByID and have a read-only items list.
func flattenFirstItemToAttrs(model interface{}, attrs map[string]hcl.Value, readOnlyAttrTfsdk string) error {
	itemMap := firstItemMapFromModel(model, readOnlyAttrTfsdk)
	if len(itemMap) == 0 {
		return nil
	}
	flat, err := hcl.ItemMapToFlatValues(itemMap, []string{"status"})
	if err != nil {
		return err
	}
	for k, v := range flat {
		attrs[k] = v
	}
	return nil
}

// flattenItemsListToTopLevel merges attrs["items"].List[0].Map into attrs and removes "items".
// Used for resources like criblio_pack where ModelToValue produces items = [ {...} ] and schema expects those fields at top level.
func flattenItemsListToTopLevel(attrs map[string]hcl.Value) {
	iv, ok := attrs["items"]
	if !ok {
		return
	}
	if iv.Kind == hcl.KindList && len(iv.List) > 0 {
		first := iv.List[0]
		if first.Kind == hcl.KindMap && first.Map != nil {
			for k, v := range first.Map {
				attrs[k] = v
			}
		}
	}
	delete(attrs, "items")
}

// filterAttrsBySchema removes from attrs any key not in the provider model's schema (tfsdk attribute names).
// Ensures generated HCL only includes supported attributes; unsupported fields are ignored.
func filterAttrsBySchema(attrs map[string]hcl.Value, modelTypeName string) {
	allowed := converter.AllAttributeNamesFromModel(modelTypeName)
	if len(allowed) == 0 {
		return
	}
	allowedSet := make(map[string]bool, len(allowed))
	for _, k := range allowed {
		allowedSet[k] = true
	}
	for k := range attrs {
		if !allowedSet[k] {
			delete(attrs, k)
		}
	}
}

// pruneNotificationTargetConfigs removes empty optional SMTP configuration
// objects returned for non-SMTP notification targets. Emitting conf = {} is
// invalid because the provider schema requires email_recipient when conf is set.
func pruneNotificationTargetConfigs(attrs map[string]hcl.Value) {
	targetConfigs, ok := attrs["target_configs"]
	if !ok || targetConfigs.Kind != hcl.KindList {
		return
	}
	for index := range targetConfigs.List {
		item := &targetConfigs.List[index]
		if item.Kind != hcl.KindMap {
			continue
		}
		conf, ok := item.Map["conf"]
		if ok && conf.Kind == hcl.KindMap && !hasNonNullValue(conf) {
			delete(item.Map, "conf")
		}
	}
	attrs["target_configs"] = targetConfigs
}

func hasNonNullValue(value hcl.Value) bool {
	switch value.Kind {
	case hcl.KindNull:
		return false
	case hcl.KindMap:
		for _, child := range value.Map {
			if hasNonNullValue(child) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// hclOptionsForType returns HCL conversion options for the given resource type,
// including skipping read-only attributes (and oneOf list attr when present) so generated config is valid.
func hclOptionsForType(typeName string, e registry.Entry) *hcl.Options {
	skip := readOnlyAttrsByType[typeName]
	if e.OneOf != nil {
		// oneOf resources store payload in a read-only list; skip that attr in HCL.
		skip = append(skip, e.OneOf.ReadOnlyAttr)
	}
	if len(skip) == 0 {
		return nil
	}
	topLevel := make(map[string]bool)
	nested := make(map[string]bool)
	for _, s := range skip {
		topLevel[s] = true
		// additional_properties appears nested (e.g. in route items); must be skipped at any level.
		if s == "additional_properties" {
			nested[s] = true
		}
		// id in route items (criblio_routes, criblio_pack_routes) is Computed; skip at nested level.
		if s == "id" && (typeName == "criblio_routes" || typeName == "criblio_pack_routes") {
			nested[s] = true
		}
	}
	opts := &hcl.Options{SkipAttributes: topLevel}
	if len(nested) > 0 {
		opts.SkipAttributesNested = nested
	}
	return opts
}
