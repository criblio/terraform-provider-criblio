package custom

import (
	"encoding/json"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

// ApplyAppscopeConfigDefaults sets required attribute defaults for criblio_appscope_config so
// generated HCL passes provider validation. The API may return null/empty for config.event.type
// (must be "ndjson") and config.metric.watch / config.event.watch (must be non-null list).
// After setting defaults, prunes null and empty-string values recursively so Optional+Computed
// attributes are omitted from HCL—writing them causes perpetual "(known after apply)" drift.
func ApplyAppscopeConfigDefaults(attrs map[string]hcl.Value) {
	config, ok := attrs["config"]
	if !ok || config.Kind != hcl.KindMap || config.Map == nil {
		return
	}
	emptyWatch := hcl.Value{Kind: hcl.KindList, List: []hcl.Value{}}
	if event, ok := config.Map["event"]; ok && event.Kind == hcl.KindMap && event.Map != nil {
		tv := event.Map["type"]
		if tv.Kind == hcl.KindNull || (tv.Kind == hcl.KindString && tv.String == "") {
			event.Map["type"] = hcl.Value{Kind: hcl.KindString, String: "ndjson"}
			config.Map["event"] = event
		}
		if w := event.Map["watch"]; w.Kind == hcl.KindNull {
			event.Map["watch"] = emptyWatch
			config.Map["event"] = event
		}
	}
	// config.metric.watch is required (NotNull validator); ensure metric block exists with watch list.
	// Use one minimal item - empty list can fail "value must be configured" in some provider versions.
	minimalWatch := hcl.Value{Kind: hcl.KindList, List: []hcl.Value{
		{Kind: hcl.KindMap, Map: map[string]hcl.Value{
			"enabled": {Kind: hcl.KindBool, Bool: false},
			"name":    {Kind: hcl.KindString, String: ".*"},
			"type":    {Kind: hcl.KindString, String: "metric"},
			"value":   {Kind: hcl.KindString, String: ".*"},
		}},
	}}
	metric, ok := config.Map["metric"]
	if !ok || metric.Kind != hcl.KindMap || metric.Map == nil {
		metric = hcl.Value{Kind: hcl.KindMap, Map: map[string]hcl.Value{"watch": minimalWatch}}
		config.Map["metric"] = metric
	} else {
		w := metric.Map["watch"]
		if w.Kind != hcl.KindList || len(w.List) == 0 {
			metric.Map["watch"] = minimalWatch
			config.Map["metric"] = metric
		}
	}
	// Prune null and empty-string values so Optional+Computed attrs are omitted (avoids drift).
	config = pruneAppscopeConfig(config)
	attrs["config"] = config
}

// pruneAppscopeConfig removes null, empty-string, and empty-list map entries recursively.
// Optional+Computed attributes written as null, "", or [] cause "(known after apply)" drift.
func pruneAppscopeConfig(v hcl.Value) hcl.Value {
	if v.Kind != hcl.KindMap && v.Kind != hcl.KindList {
		return v
	}
	if v.Kind == hcl.KindMap {
		out := make(map[string]hcl.Value, len(v.Map))
		for k, val := range v.Map {
			if val.Kind == hcl.KindNull {
				continue
			}
			if val.Kind == hcl.KindString && val.String == "" {
				continue
			}
			if val.Kind == hcl.KindList && len(val.List) == 0 {
				continue
			}
			out[k] = pruneAppscopeConfig(val)
		}
		return hcl.Value{Kind: hcl.KindMap, Map: out}
	}
	list := make([]hcl.Value, len(v.List))
	for i, el := range v.List {
		list[i] = pruneAppscopeConfig(el)
	}
	return hcl.Value{Kind: hcl.KindList, List: list}
}

// ApplyPackDefaults sets required attribute defaults for criblio_pack so generated HCL
// passes provider validation. Pack tags require "domain" key when tags is set.
func ApplyPackDefaults(attrs map[string]hcl.Value) {
	tags, ok := attrs["tags"]
	if !ok || tags.Kind != hcl.KindMap || tags.Map == nil {
		return
	}
	emptyList := hcl.Value{Kind: hcl.KindList, List: []hcl.Value{}}
	cur, hasDomain := tags.Map["domain"]
	if !hasDomain || cur.Kind == hcl.KindNull || (cur.Kind == hcl.KindList && cur.List == nil) {
		tags.Map["domain"] = emptyList
		attrs["tags"] = tags
	}
}

// ApplyPackVarsDefaults fixes pack_vars attributes: args must be list of object (not JSON string);
// id, description, lib, tags, type, value must not have extra JSON quotes (e.g. "\"test\"" -> "test").
func ApplyPackVarsDefaults(attrs map[string]hcl.Value) {
	// args: convert JSON string "[]" or "[{...}]" to list of objects
	if v, ok := attrs["args"]; ok && v.Kind == hcl.KindString {
		argsList := parseJSONListOfObjects(v.String)
		if argsList != nil {
			attrs["args"] = hcl.Value{Kind: hcl.KindList, List: argsList}
		}
	}
	// Strip extra quotes from string attrs that may come from JSON (e.g. "\"test\"" -> "test")
	for _, k := range []string{"id", "description", "lib", "tags", "type", "value"} {
		if v, ok := attrs[k]; ok && v.Kind == hcl.KindString && v.String != "" {
			s := v.String
			// Remove surrounding JSON quotes: "\"test\"" or \"test\"
			if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
				s = s[1 : len(s)-1]
			}
			if len(s) >= 2 && s[0] == '\\' && s[1] == '"' {
				s = s[2:]
				if len(s) >= 1 && s[len(s)-1] == '"' {
					s = s[:len(s)-1]
				}
			}
			if s != v.String {
				attrs[k] = hcl.Value{Kind: hcl.KindString, String: s}
			}
		}
	}
}

// parseJSONListOfObjects parses s as JSON array of objects; returns nil if invalid.
func parseJSONListOfObjects(s string) []hcl.Value {
	var raw []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return nil
	}
	out := make([]hcl.Value, 0, len(raw))
	for _, m := range raw {
		obj := make(map[string]hcl.Value)
		for k, v := range m {
			switch tv := v.(type) {
			case string:
				obj[k] = hcl.Value{Kind: hcl.KindString, String: tv}
			case float64:
				obj[k] = hcl.Value{Kind: hcl.KindNumber, Number: tv}
			case bool:
				obj[k] = hcl.Value{Kind: hcl.KindBool, Bool: tv}
			default:
				obj[k] = hcl.Value{Kind: hcl.KindNull}
			}
		}
		out = append(out, hcl.Value{Kind: hcl.KindMap, Map: obj})
	}
	return out
}

// ApplyGroupDefaults populates criblio_group attrs from the model when API returns values
// (description, is_fleet, on_prem, tags, type, etc.) so generated HCL matches state and avoids plan drift.
func ApplyGroupDefaults(attrs map[string]hcl.Value, model interface{}) {
	gm, ok := model.(*provider.GroupResourceModel)
	if !ok || gm == nil {
		return
	}
	setStr := func(attr, s string) {
		if s != "" {
			attrs[attr] = hcl.Value{Kind: hcl.KindString, String: s}
		}
	}
	setBool := func(attr string, b bool) {
		attrs[attr] = hcl.Value{Kind: hcl.KindBool, Bool: b}
	}
	setNum := func(attr string, n float64) {
		attrs[attr] = hcl.Value{Kind: hcl.KindNumber, Number: n}
	}
	if !gm.Description.IsNull() && !gm.Description.IsUnknown() {
		setStr("description", gm.Description.ValueString())
	}
	if !gm.IsFleet.IsNull() && !gm.IsFleet.IsUnknown() {
		setBool("is_fleet", gm.IsFleet.ValueBool())
	}
	if !gm.OnPrem.IsNull() && !gm.OnPrem.IsUnknown() {
		setBool("on_prem", gm.OnPrem.ValueBool())
	}
	if !gm.Tags.IsNull() && !gm.Tags.IsUnknown() {
		setStr("tags", gm.Tags.ValueString())
	}
	if !gm.Type.IsNull() && !gm.Type.IsUnknown() {
		setStr("type", gm.Type.ValueString())
	}
	if !gm.WorkerRemoteAccess.IsNull() && !gm.WorkerRemoteAccess.IsUnknown() {
		setBool("worker_remote_access", gm.WorkerRemoteAccess.ValueBool())
	}
	if !gm.MaxWorkerAge.IsNull() && !gm.MaxWorkerAge.IsUnknown() && gm.MaxWorkerAge.ValueString() != "" {
		setStr("max_worker_age", gm.MaxWorkerAge.ValueString())
	}
	if !gm.Name.IsNull() && !gm.Name.IsUnknown() && gm.Name.ValueString() != "" {
		setStr("name", gm.Name.ValueString())
	}
	if !gm.EstimatedIngestRate.IsNull() && !gm.EstimatedIngestRate.IsUnknown() {
		setNum("estimated_ingest_rate", gm.EstimatedIngestRate.ValueFloat64())
	}
	if !gm.Provisioned.IsNull() && !gm.Provisioned.IsUnknown() {
		setBool("provisioned", gm.Provisioned.ValueBool())
	}
	if gm.Cloud != nil && !gm.Cloud.Provider.IsNull() && !gm.Cloud.Provider.IsUnknown() {
		if cloud, ok := attrs["cloud"]; ok && cloud.Kind == hcl.KindMap && cloud.Map != nil {
			cloud.Map["provider"] = hcl.Value{Kind: hcl.KindString, String: gm.Cloud.Provider.ValueString()}
			if !gm.Cloud.Region.IsNull() && !gm.Cloud.Region.IsUnknown() {
				cloud.Map["region"] = hcl.Value{Kind: hcl.KindString, String: gm.Cloud.Region.ValueString()}
			}
			attrs["cloud"] = cloud
		} else if !gm.Cloud.Region.IsNull() && !gm.Cloud.Region.IsUnknown() {
			attrs["cloud"] = hcl.Value{Kind: hcl.KindMap, Map: map[string]hcl.Value{
				"provider": {Kind: hcl.KindString, String: gm.Cloud.Provider.ValueString()},
				"region":   {Kind: hcl.KindString, String: gm.Cloud.Region.ValueString()},
			}}
		}
	}
}

// ApplyProjectDefaults sets required attribute defaults for criblio_project so generated HCL
// passes provider validation. The provider requires "subscriptions" and "destinations".
func ApplyProjectDefaults(attrs map[string]hcl.Value) {
	emptyList := hcl.Value{Kind: hcl.KindList, List: []hcl.Value{}}
	if v, ok := attrs["subscriptions"]; !ok || v.Kind == hcl.KindNull {
		attrs["subscriptions"] = emptyList
	}
	if v, ok := attrs["destinations"]; !ok || v.Kind == hcl.KindNull {
		attrs["destinations"] = emptyList
	}
}

