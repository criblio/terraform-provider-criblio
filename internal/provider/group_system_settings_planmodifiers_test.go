package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func TestGroupSystemSettingsObjectFromAPIOrPrior(t *testing.T) {
	sslType := map[string]attr.Type{
		"disabled":   types.BoolType,
		"passphrase": types.StringType,
	}
	objectType := map[string]attr.Type{
		"base_url": types.StringType,
		"host":     types.StringType,
		"ssl":      types.ObjectType{AttrTypes: sslType},
	}
	prior := types.ObjectValueMust(objectType, map[string]attr.Value{
		"base_url": types.StringValue("https://old.example.com"),
		"host":     types.StringValue("old.example.com"),
		"ssl": types.ObjectValueMust(sslType, map[string]attr.Value{
			"disabled":   types.BoolValue(false),
			"passphrase": types.StringValue("configured-secret"),
		}),
	})
	api := types.ObjectValueMust(objectType, map[string]attr.Value{
		"base_url": types.StringValue("https://new.example.com"),
		"host":     types.StringNull(),
		"ssl": types.ObjectValueMust(sslType, map[string]attr.Value{
			"disabled":   types.BoolValue(true),
			"passphrase": types.StringNull(),
		}),
	})

	merged := groupSystemSettingsObjectFromAPIOrPrior(api, prior)
	attributes := merged.Attributes()
	assert.Equal(t, "https://new.example.com", attributes["base_url"].(types.String).ValueString())
	assert.True(t, attributes["host"].IsNull(), "ordinary API null should clear prior state")
	ssl := attributes["ssl"].(types.Object).Attributes()
	assert.True(t, ssl["disabled"].(types.Bool).ValueBool())
	assert.Equal(t, "configured-secret", ssl["passphrase"].(types.String).ValueString())

	removed := groupSystemSettingsObjectFromAPIOrPrior(types.ObjectNull(objectType), prior)
	assert.True(t, removed.IsNull(), "removed API object should clear prior state")
}
