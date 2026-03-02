// Package export converts discovery results into generator ResourceItems.
package export

// readOnlyAttrsByType lists attribute names (tfsdk) that are read-only (Computed only) per type.
// For oneOf resources, the OneOf.ReadOnlyAttr is added at runtime so this only needs extra entries when not using OneOf.
var readOnlyAttrsByType = map[string][]string{
	"criblio_certificate":             {"ca", "in_use", "passphrase"}, // provider computes these; API may not return cert/priv_key on GET
	"criblio_global_var":              {"items"},                      // provider marks items as Computed only; config comes from GetByID and we flatten Items[0]
	"criblio_group_system_settings":   {"items"},                      // provider marks items as read-only; configurable attrs are top-level (api, backups, etc.)
	"criblio_pack_vars":               {"items"},                       // items is Computed; config comes from flattenFirstItemToAttrs (description, lib, tags, type, value)
	"criblio_pack_breakers":           {"items"},                       // items is read-only (Computed); id, group_id, pack are configurable
	"criblio_pack_pipeline":           {"items"},                      // GET returns Routes (items); conf filled from lib pipeline or minimal default
	"criblio_routes":                  {"id", "additional_properties"}, // id read-only; additional_properties in route items is optional/empty, omit from HCL
	"criblio_pack_routes":             {"id", "additional_properties", "items"}, // items is Computed; routes come from items[0], same as criblio_routes
	"criblio_search_dataset":          {"id", "description", "provider_id", "type"}, // all computed/read-only; config comes from type-specific blocks only
	"criblio_search_dataset_provider": {"id", "description", "provider_id", "type"},
}

// flattenItemsToAttrsTypes are resource types whose API returns payload in a list (Items) and schema
// has top-level attributes; we fetch via GetByID then flatten Items[0] into attrs (snake_case, no nested block).
var flattenItemsToAttrsTypes = map[string]bool{
	"criblio_global_var": true,
	"criblio_pack_vars":  true, // same as global_var: Items[0] -> description, lib, tags, type, value
}

// flattenItemsToTopLevelTypes are resource types where attrs["items"] (after ModelToValue) is a list
// whose first element should be merged to top-level and "items" removed (e.g. criblio_pack, criblio_group_system_settings).
var flattenItemsToTopLevelTypes = map[string]bool{
	"criblio_pack":                  true,
	"criblio_pack_breakers":         true, // first item (EventBreakerRuleset: description, lib, rules, etc.) becomes top-level; items is computed
	"criblio_pack_lookups":          true, // first item (LookupFile: content, description, mode, tags) becomes top-level; items is computed
	"criblio_pack_pipeline":         true, // first item (Pipeline: id, conf) becomes top-level; items is computed
	"criblio_pack_routes":           true, // same as criblio_routes: first item (id, routes, comments, groups) becomes top-level
	"criblio_pack_vars":             true, // first item (args, description, id, lib, tags, type, value) becomes top-level
	"criblio_group_system_settings": true,
	"criblio_search_usage_group":    true, // API returns list; schema has required top-level rules, items is computed
}
