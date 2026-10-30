package export

import (
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestMonitorExportKeepsUpstreamConfiguration(t *testing.T) {
	var model provider.MonitorModel
	require.NoError(t, json.Unmarshal([]byte(`{
		"id":"monitor-export", "managedBy":"terraform", "silence":[], "expr":[],
		"$profileRef":{"id":"profile", "fingerprint":"profile-fingerprint"},
		"firingCondition":{"$configRefs":[{"id":"condition", "fingerprint":"condition-fingerprint"}]},
		"notification":{"type":"policy", "config":[], "enabled":false}
		,"query":{"A":{"$configRefs":[{"id":"query", "fingerprint":"query-fingerprint"}]}},
		"metadata":{"owner":"example"}
	}`), &model))
	attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_monitor", registry.Entry{}))
	require.NoError(t, err)
	require.Equal(t, "terraform", attrs["managed_by"].String)
	require.Equal(t, hcl.KindList, attrs["silence"].Kind)
	require.Empty(t, attrs["silence"].List)
	require.Equal(t, "profile", attrs["profile_ref"].Map["id"].String)
	refs := attrs["firing_condition_inheritable"].Map["config_refs"]
	require.Len(t, refs.List, 1)
	require.Equal(t, "condition", refs.List[0].Map["id"].String)
	require.Contains(t, attrs, "monitor_policy_notification_config")
	require.Equal(t, "query", attrs["query"].Map["A"].Map["config_refs"].List[0].Map["id"].String)
	require.Equal(t, "example", attrs["aetos_metadata"].Map["owner"].String)
}

func TestMonitoringResourceExport(t *testing.T) {
	for _, tc := range []struct {
		name, wire string
		model      any
		computed   []string
	}{
		{"criblio_monitor_shared_config", `{"id":"shared","config":{"priority":{"value":"P3"}},"fingerprint":"hash"}`, &provider.MonitorSharedConfigModel{}, []string{"fingerprint"}},
		{"criblio_monitor_profile", `{"id":"profile","name":"Platform","configIds":["shared"]}`, &provider.MonitorProfileModel{}, nil},
		{"criblio_monitor_silence", `{"id":"silence","monitor_id":"cpu","begins_at":1893456000000,"expires_at":1893542400000,"tags":{"host":["a","b"]},"created_at":123456,"created_by":"api-user"}`, &provider.MonitorSilenceModel{}, []string{"created_at", "created_by"}},
		{"criblio_notification_policy", `{"id":"policy","order":10,"conditions":[[{"key":"count","operator":"=","value":3}]],"templateTargetPairs":[{"targetId":"target","templateId":"template"}]}`, &provider.NotificationPolicyModel{}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(tc.wire), tc.model))
			attrs, err := hcl.ModelToValue(tc.model, hclOptionsForType(tc.name, registry.Entry{}))
			require.NoError(t, err)
			require.Contains(t, attrs, "id")
			for _, name := range tc.computed {
				require.NotContains(t, attrs, name)
			}
			switch tc.name {
			case "criblio_monitor_shared_config":
				require.Contains(t, attrs["config"].String, "P3")
			case "criblio_monitor_profile":
				require.Equal(t, "shared", attrs["config_ids"].List[0].String)
			case "criblio_monitor_silence":
				require.Contains(t, attrs["tags"].String, `"b"`)
			case "criblio_notification_policy":
				require.Contains(t, attrs["conditions"].String, `"value":3`)
			}
		})
	}
}
