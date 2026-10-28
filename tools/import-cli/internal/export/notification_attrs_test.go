package export

import (
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestNotificationExportKeepsDeliveryConfiguration(t *testing.T) {
	var model provider.NotificationModel
	require.NoError(t, json.Unmarshal([]byte(`{
		"id":"direct-notification", "condition":"high-volume", "group":"default",
		"mode":"direct", "pack":"api-owned-pack",
		"templateTargetPairs":[{"targetId":"system_notifications", "templateId":"default-slack"}]
	}`), &model))
	attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_notification", registry.Entry{}))
	require.NoError(t, err)
	require.Equal(t, "direct", attrs["mode"].String)
	pairs := attrs["template_target_pairs"]
	require.Equal(t, hcl.KindList, pairs.Kind)
	require.Len(t, pairs.List, 1)
	require.Equal(t, "system_notifications", pairs.List[0].Map["target_id"].String)
	require.Equal(t, "default-slack", pairs.List[0].Map["template_id"].String)
	require.NotContains(t, attrs, "pack")
}
