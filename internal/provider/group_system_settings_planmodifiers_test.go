package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupSystemSettingsEmptyObjectsAreOptionalComputed(t *testing.T) {
	var response resource.SchemaResponse
	NewGroupSystemSettingsResource().Schema(context.Background(), resource.SchemaRequest{}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)

	for _, name := range []string{"backups", "custom_logo", "pii", "rollback", "sni", "sockets", "upgrade_group_settings"} {
		attribute, ok := response.Schema.Attributes[name].(schema.SingleNestedAttribute)
		require.True(t, ok, "%s should be a single nested attribute", name)
		assert.True(t, attribute.Optional, "%s should be optional", name)
		assert.True(t, attribute.Computed, "%s should be computed", name)
	}
}
