package provider

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// nullSearchUsageGroupAttrs builds an all-null attribute map for the given
// schema, so each test only has to set the attributes it cares about.
func nullSearchUsageGroupAttrs(t *testing.T, attrTypes map[string]attr.Type) map[string]attr.Value {
	t.Helper()
	values := make(map[string]attr.Value, len(attrTypes))
	for name, typ := range attrTypes {
		value, err := SearchUsageGroupTerraformNullValue(typ)
		if err != nil {
			t.Fatalf("null value for %s: %v", name, err)
		}
		values[name] = value
	}
	return values
}

// A scheduling limit written as `limit = 10` in HCL arrives as a number inside
// the dynamic attribute and has to reach the API as a JSON number, not "10".
func TestSearchUsageGroupMarshalJSONKeepsNumberAndPercentLimits(t *testing.T) {
	scheduling := nullSearchUsageGroupAttrs(t, SearchUsageGroupRulesMaxConcurrentSearchesAttrTypes())
	scheduling["type"] = types.StringValue("maxConcurrentSearches")
	scheduling["limit"] = types.DynamicValue(types.NumberValue(big.NewFloat(10)))

	perUser := nullSearchUsageGroupAttrs(t, SearchUsageGroupRulesMaxConcurrentAdhocSearchesPerUserAttrTypes())
	perUser["type"] = types.StringValue("maxConcurrentAdhocSearchesPerUser")
	perUser["limit"] = types.DynamicValue(types.StringValue("50%"))

	results := nullSearchUsageGroupAttrs(t, SearchUsageGroupRulesMaxResultsPerSearchAttrTypes())
	results["type"] = types.StringValue("maxResultsPerSearch")
	results["limit"] = types.DynamicValue(types.NumberValue(big.NewFloat(100000)))

	rules := nullSearchUsageGroupAttrs(t, SearchUsageGroupRulesAttrTypes())
	rules["max_concurrent_searches"] = types.ObjectValueMust(SearchUsageGroupRulesMaxConcurrentSearchesAttrTypes(), scheduling)
	rules["max_concurrent_adhoc_searches_per_user"] = types.ObjectValueMust(SearchUsageGroupRulesMaxConcurrentAdhocSearchesPerUserAttrTypes(), perUser)
	rules["max_results_per_search"] = types.ObjectValueMust(SearchUsageGroupRulesMaxResultsPerSearchAttrTypes(), results)

	users := map[string]attr.Value{
		"user-1": types.ObjectValueMust(SearchUsageGroupUsersAttrTypes(), map[string]attr.Value{
			"id":    types.StringValue("user-1"),
			"email": types.StringValue("search.admin@example.com"),
		}),
	}

	model := SearchUsageGroupModel{
		ID:         types.StringValue("power_users"),
		Enabled:    types.BoolValue(true),
		Rules:      types.ObjectValueMust(SearchUsageGroupRulesAttrTypes(), rules),
		Users:      types.MapValueMust(types.ObjectType{AttrTypes: SearchUsageGroupUsersAttrTypes()}, users),
		UsersCount: types.Int64Value(1),
	}

	raw, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("marshal usage group: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("decode request body: %v", err)
	}

	if _, ok := body["usersCount"]; ok {
		t.Error("usersCount is API-owned and must not be sent in the request body")
	}

	apiRules, ok := body["rules"].(map[string]any)
	if !ok {
		t.Fatalf("rules missing from request body: %s", raw)
	}
	for name, want := range map[string]any{
		"maxConcurrentSearches":             float64(10),
		"maxConcurrentAdhocSearchesPerUser": "50%",
		"maxResultsPerSearch":               float64(100000),
	} {
		rule, ok := apiRules[name].(map[string]any)
		if !ok {
			t.Fatalf("rule %s missing from request body: %s", name, raw)
		}
		if got := rule["limit"]; got != want {
			t.Errorf("rule %s limit = %#v, want %#v", name, got, want)
		}
	}

	apiUsers, ok := body["users"].(map[string]any)
	if !ok {
		t.Fatalf("users missing from request body: %s", raw)
	}
	user, ok := apiUsers["user-1"].(map[string]any)
	if !ok {
		t.Fatalf("user-1 missing from request body: %s", raw)
	}
	if user["email"] != "search.admin@example.com" {
		t.Errorf("user-1 email = %#v", user["email"])
	}
}

// Reads and imports decode both NumberOrPercent forms, which a string-typed
// `limit` could not do.
func TestSearchUsageGroupUnmarshalJSONReadsBothLimitForms(t *testing.T) {
	var model SearchUsageGroupModel
	wire := `{
		"id": "power_users",
		"enabled": true,
		"rules": {
			"maxResultsPerSearch": {"type": "maxResultsPerSearch", "limit": 100000},
			"maxConcurrentSearches": {"type": "maxConcurrentSearches", "limit": 10},
			"maxConcurrentAdhocSearchesPerUser": {"type": "maxConcurrentAdhocSearchesPerUser", "limit": "50%"}
		},
		"users": {"user-1": {"id": "user-1", "email": "search.admin@example.com"}},
		"usersCount": 1
	}`
	if err := json.Unmarshal([]byte(wire), &model); err != nil {
		t.Fatalf("unmarshal usage group: %v", err)
	}

	if model.ID.ValueString() != "power_users" {
		t.Errorf("id = %q", model.ID.ValueString())
	}
	if model.UsersCount.ValueInt64() != 1 {
		t.Errorf("users_count = %d", model.UsersCount.ValueInt64())
	}
	if got := len(model.Users.Elements()); got != 1 {
		t.Fatalf("users count = %d", got)
	}

	rules := model.Rules.Attributes()
	numeric := rules["max_results_per_search"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := numeric.UnderlyingValue().(types.Float64); !ok || value.ValueFloat64() != 100000 {
		t.Errorf("max_results_per_search limit = %#v", numeric.UnderlyingValue())
	}

	scheduling := rules["max_concurrent_searches"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := scheduling.UnderlyingValue().(types.Float64); !ok || value.ValueFloat64() != 10 {
		t.Errorf("max_concurrent_searches limit = %#v", scheduling.UnderlyingValue())
	}

	percent := rules["max_concurrent_adhoc_searches_per_user"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := percent.UnderlyingValue().(types.String); !ok || value.ValueString() != "50%" {
		t.Errorf("max_concurrent_adhoc_searches_per_user limit = %#v", percent.UnderlyingValue())
	}
}

// Editing a Usage Group in the Cribl UI rewrites limits as strings: it adds a
// `coordinatorHeapMemoryLimit` rule carrying the sentinel `"default"` and turns
// other limits into numeric strings such as `"8"`. Decoding has to survive that,
// or `terraform plan` fails outright against a group a user has touched. This is
// the payload a Cribl.Cloud workspace returned.
func TestSearchUsageGroupUnmarshalJSONReadsUIWrittenStringLimits(t *testing.T) {
	var model SearchUsageGroupModel
	wire := `{
		"id": "power_users",
		"enabled": false,
		"rules": {
			"coordinatorHeapMemoryLimit": {"limit": "default", "type": "coordinatorHeapMemoryLimit", "description": "Maximum heap memory that can be allocated to the coordinator process"},
			"maxConcurrentAdhocSearchesPerUser": {"limit": "8", "type": "maxConcurrentAdhocSearchesPerUser"},
			"maxResultsPerSearch": {"limit": 100000, "type": "maxResultsPerSearch"},
			"maxRunningTimePerSearch": {"limit": 300, "type": "maxRunningTimePerSearch", "metric": "running_time"}
		},
		"users": {},
		"usersCount": 0
	}`
	if err := json.Unmarshal([]byte(wire), &model); err != nil {
		t.Fatalf("unmarshal UI-written usage group: %v", err)
	}

	rules := model.Rules.Attributes()

	// A LimitRule holding a sentinel string rather than a number.
	heap := rules["coordinator_heap_memory_limit"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := heap.UnderlyingValue().(types.String); !ok || value.ValueString() != "default" {
		t.Errorf("coordinator_heap_memory_limit limit = %#v", heap.UnderlyingValue())
	}

	// A scheduling limit the UI rewrote from a percentage to a numeric string.
	adhoc := rules["max_concurrent_adhoc_searches_per_user"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := adhoc.UnderlyingValue().(types.String); !ok || value.ValueString() != "8" {
		t.Errorf("max_concurrent_adhoc_searches_per_user limit = %#v", adhoc.UnderlyingValue())
	}

	// Numeric limits in the same payload still decode as numbers.
	results := rules["max_results_per_search"].(types.Object).Attributes()["limit"].(types.Dynamic)
	if value, ok := results.UnderlyingValue().(types.Float64); !ok || value.ValueFloat64() != 100000 {
		t.Errorf("max_results_per_search limit = %#v", results.UnderlyingValue())
	}
}
