package hcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestinationItemToOutputBlock(t *testing.T) {
	t.Run("prometheus type", func(t *testing.T) {
		item := map[string]string{
			"type":    `"prometheus"`,
			"id":      `"elastic-prometheus"`,
			"url":     `"http://prometheus.example:9201"`,
			"timeoutSec": "30",
		}
		blockName, value, err := DestinationItemToOutputBlock(item)
		require.NoError(t, err)
		assert.Equal(t, "output_prometheus", blockName)
		assert.Equal(t, KindMap, value.Kind)
		assert.Contains(t, value.Map, "type")
		assert.Contains(t, value.Map, "id")
		assert.Contains(t, value.Map, "url")
		assert.Contains(t, value.Map, "timeout_sec")
		assert.Equal(t, "prometheus", value.Map["type"].String)
		assert.Equal(t, float64(30), value.Map["timeout_sec"].Number)
	})

	t.Run("open_telemetry type", func(t *testing.T) {
		item := map[string]string{
			"type": `"open_telemetry"`,
			"id":   `"elastic-otel"`,
		}
		blockName, value, err := DestinationItemToOutputBlock(item)
		require.NoError(t, err)
		assert.Equal(t, "output_open_telemetry", blockName)
		assert.Equal(t, KindMap, value.Kind)
		assert.Equal(t, "open_telemetry", value.Map["type"].String)
	})

	t.Run("skips status", func(t *testing.T) {
		item := map[string]string{
			"type":   `"devnull"`,
			"id":     `"devnull"`,
			"status": `{"health":"Green"}`,
		}
		_, value, err := DestinationItemToOutputBlock(item)
		require.NoError(t, err)
		assert.NotContains(t, value.Map, "status")
	})

	t.Run("empty item", func(t *testing.T) {
		blockName, value, err := DestinationItemToOutputBlock(nil)
		require.NoError(t, err)
		assert.Equal(t, "", blockName)
		assert.True(t, value.IsNull())
	})

	t.Run("missing type", func(t *testing.T) {
		item := map[string]string{"id": `"x"`}
		_, _, err := DestinationItemToOutputBlock(item)
		require.Error(t, err)
	})
}

func TestItemMapToBlock_generic(t *testing.T) {
	// Generic oneOf: same item shape, different prefix (e.g. input_collector_ for collectors).
	item := map[string]string{
		"type": `"rest"`,
		"id":   `"my-collector"`,
	}
	blockName, value, err := ItemMapToBlock(item, "type", "input_collector_", "", []string{"status"}, nil)
	require.NoError(t, err)
	assert.Equal(t, "input_collector_rest", blockName)
	assert.Equal(t, KindMap, value.Kind)
	assert.Equal(t, "rest", value.Map["type"].String)
}

func TestItemMapToBlock_collector_collection_alias(t *testing.T) {
	// API returns type "collection" for REST-style collectors; provider expects input_collector_rest.
	item := map[string]string{
		"type": `"collection"`,
		"id":   `"crowdstrike_ngsiem_api"`,
	}
	alias := map[string]string{"collection": "rest"}
	blockName, value, err := ItemMapToBlock(item, "type", "input_collector_", "", []string{"status"}, alias)
	require.NoError(t, err)
	assert.Equal(t, "input_collector_rest", blockName)
	assert.Equal(t, KindMap, value.Kind)
	assert.Equal(t, "collection", value.Map["type"].String, "item type is unchanged; only block name is aliased")
}

func TestItemMapToBlock_omits_empty_lists(t *testing.T) {
	// output_cribl_http and output_webhook have urls with SizeAtLeast(1); we must not emit urls = [].
	item := map[string]string{
		"type": `"cribl_http"`,
		"id":   `"cribl_http_route"`,
		"url":  `"https://example.com:10200"`,
		"urls": `[]`,
	}
	blockName, value, err := ItemMapToBlock(item, "type", "output_", "", []string{"status"}, nil)
	require.NoError(t, err)
	assert.Equal(t, "output_cribl_http", blockName)
	assert.Contains(t, value.Map, "url")
	assert.NotContains(t, value.Map, "urls", "empty urls list must be omitted to satisfy SizeAtLeast(1)")
}

func TestItemMapToBlock_notification_target_suffix(t *testing.T) {
	t.Run("smtp type gets _target suffix", func(t *testing.T) {
		item := map[string]string{"type": `"smtp"`, "id": `"system_email"`}
		blockName, _, err := ItemMapToBlock(item, "type", "", "_target", []string{"status"}, nil)
		require.NoError(t, err)
		assert.Equal(t, "smtp_target", blockName)
	})
	t.Run("WebhookTarget normalizes to webhook_target; suffix not duplicated", func(t *testing.T) {
		item := map[string]string{"type": `"WebhookTarget"`, "id": `"my-webhook"`}
		blockName, _, err := ItemMapToBlock(item, "type", "", "_target", []string{"status"}, nil)
		require.NoError(t, err)
		assert.Equal(t, "webhook_target", blockName)
	})
}

func TestResolveOneOfBlockNameRaw(t *testing.T) {
	supported := []string{"smtp_target", "slack_target", "pager_duty_target", "sns_target", "webhook_target"}
	t.Run("unsupported type returns false", func(t *testing.T) {
		suffix, ok := ResolveOneOfBlockNameRaw(`"bulletin_message"`, supported, "")
		assert.False(t, ok)
		assert.Empty(t, suffix)
	})
	t.Run("smtp resolves to smtp_target", func(t *testing.T) {
		suffix, ok := ResolveOneOfBlockNameRaw(`"smtp"`, supported, "")
		assert.True(t, ok)
		assert.Equal(t, "smtp_target", suffix)
	})
	t.Run("PascalCase WebhookTarget resolves to webhook_target", func(t *testing.T) {
		suffix, ok := ResolveOneOfBlockNameRaw(`"WebhookTarget"`, supported, "")
		assert.True(t, ok)
		assert.Equal(t, "webhook_target", suffix)
	})
	t.Run("with prefix matches suffix", func(t *testing.T) {
		outputSupported := []string{"output_prometheus", "output_webhook"}
		suffix, ok := ResolveOneOfBlockNameRaw(`"prometheus"`, outputSupported, "output_")
		assert.True(t, ok)
		assert.Equal(t, "prometheus", suffix)
	})
}
