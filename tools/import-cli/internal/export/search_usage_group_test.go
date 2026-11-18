package export

import (
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/require"
)

const searchUsageGroupWire = `{
	"id": "power_users",
	"description": "Members with elevated search limits.",
	"enabled": true,
	"rules": {
		"maxResultsPerSearch": {"type": "maxResultsPerSearch", "limit": 100000},
		"maxRunningTimePerSearch": {"type": "maxRunningTimePerSearch", "limit": 300, "description": "Five minutes per search."},
		"maxConcurrentAdhocSearchesPerUser": {"type": "maxConcurrentAdhocSearchesPerUser", "limit": "50%"},
		"maxConcurrentSearches": {"type": "maxConcurrentSearches", "limit": 10}
	},
	"users": {
		"user-1": {"id": "user-1", "email": "search.admin@example.com"}
	},
	"usersCount": 1
}`

func TestSearchUsageGroupExportKeepsRulesAndUsers(t *testing.T) {
	var model provider.SearchUsageGroupModel
	require.NoError(t, json.Unmarshal([]byte(searchUsageGroupWire), &model))

	attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_search_usage_group", registry.Entry{}))
	require.NoError(t, err)

	require.Equal(t, "power_users", attrs["id"].String)
	require.Equal(t, "Members with elevated search limits.", attrs["description"].String)
	require.True(t, attrs["enabled"].Bool)

	// usersCount is derived by the API, so it must not round-trip into HCL.
	require.NotContains(t, attrs, "users_count")

	rules := attrs["rules"].Map
	require.Equal(t, "maxResultsPerSearch", rules["max_results_per_search"].Map["type"].String)
	require.Equal(t, float64(100000), rules["max_results_per_search"].Map["limit"].Number)
	require.Equal(t, float64(300), rules["max_running_time_per_search"].Map["limit"].Number)
	// The API owns each limit type's description and metric and overwrites
	// whatever is sent, so neither is part of the exported configuration.
	require.NotContains(t, rules["max_running_time_per_search"].Map, "description")
	require.NotContains(t, rules["max_results_per_search"].Map, "metric")
	// Scheduling limits keep whichever NumberOrPercent form the API reports.
	require.Equal(t, hcl.KindString, rules["max_concurrent_adhoc_searches_per_user"].Map["limit"].Kind)
	require.Equal(t, "50%", rules["max_concurrent_adhoc_searches_per_user"].Map["limit"].String)
	require.Equal(t, hcl.KindNumber, rules["max_concurrent_searches"].Map["limit"].Kind)
	require.Equal(t, float64(10), rules["max_concurrent_searches"].Map["limit"].Number)
	// Rules the Usage Group does not set stay out of the exported configuration.
	require.True(t, rules["max_bytes_read_per_search"].IsNull())

	users := attrs["users"].Map
	require.Len(t, users, 1)
	require.Equal(t, "user-1", users["user-1"].Map["id"].String)
	require.Equal(t, "search.admin@example.com", users["user-1"].Map["email"].String)
}

// Built-in Usage Groups are skipped whether or not --exclude-defaults is set.
// A full export without the flag previously emitted the built-in `default`
// group, whose empty `rules` prunes away in generated HCL, leaving a config that
// failed `terraform plan` with `The argument "rules" is required`.
func TestSearchUsageGroupExportSkipsBuiltInsWithoutExcludeDefaults(t *testing.T) {
	require.True(t, isSearchUsageGroupType("criblio_search_usage_group"))
	require.False(t, isSearchUsageGroupType("criblio_search_dataset"))
	for _, id := range []string{"default", "default_api_credentials", "system"} {
		require.True(t, DefaultResource("criblio_search_usage_group", map[string]string{"id": id}, nil, IncludeOverride{}), id)
	}
}

// The built-in Usage Groups cannot be created or deleted through the API, so
// they must never be exported as managed resources. A live export of a
// Cribl.Cloud workspace reports exactly these three.
func TestSearchUsageGroupExportSkipsDefaultGroup(t *testing.T) {
	for _, tc := range []struct {
		id   string
		skip bool
	}{
		{"default", true},
		{"default_api_credentials", true},
		{"system", true},
		{"power_users", false},
	} {
		t.Run(tc.id, func(t *testing.T) {
			skipped := DefaultResource("criblio_search_usage_group", map[string]string{"id": tc.id}, nil, IncludeOverride{})
			require.Equal(t, tc.skip, skipped)
		})
	}
}
