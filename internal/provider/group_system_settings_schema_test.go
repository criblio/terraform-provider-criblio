package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupSystemSettingsOmittedTLSFieldsAreOptional(t *testing.T) {
	var response resource.SchemaResponse
	NewGroupSystemSettingsResource().Schema(context.Background(), resource.SchemaRequest{}, &response)
	require.False(t, response.Diagnostics.HasError())

	api := response.Schema.Attributes["api"].(schema.SingleNestedAttribute)
	ssl := api.Attributes["ssl"].(schema.SingleNestedAttribute)
	for _, name := range []string{"ca_path", "cert_path", "disabled", "passphrase", "priv_key_path"} {
		attribute := ssl.Attributes[name]
		assert.True(t, attribute.IsOptional(), "api.ssl.%s should be optional", name)
		assert.False(t, attribute.IsComputed(), "api.ssl.%s should not remain unknown after apply", name)
	}

	tls := response.Schema.Attributes["tls"].(schema.SingleNestedAttribute)
	for _, name := range []string{"default_cipher_list", "default_ecdh_curve", "max_version", "min_version", "reject_unauthorized"} {
		attribute := tls.Attributes[name]
		assert.True(t, attribute.IsOptional(), "tls.%s should be optional", name)
		assert.False(t, attribute.IsComputed(), "tls.%s should not remain unknown after apply", name)
	}
}
