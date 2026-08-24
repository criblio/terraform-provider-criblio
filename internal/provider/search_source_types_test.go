package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSearchSourceMarshalJSONUsesTokenSecret(t *testing.T) {
	authToken := types.ObjectValueMust(SearchSourceAuthTokensAttrTypes(), map[string]attr.Value{
		"description": types.StringNull(),
		"enabled":     types.BoolValue(true),
		"token":       types.StringValue("search-source-token"),
	})
	payload, err := json.Marshal(SearchSourceModel{
		AuthTokens: types.ListValueMust(authToken.Type(context.Background()), []attr.Value{authToken}),
	})
	if err != nil {
		t.Fatalf("marshal search source: %v", err)
	}

	var output map[string]any
	if err := json.Unmarshal(payload, &output); err != nil {
		t.Fatalf("unmarshal search source payload: %v", err)
	}
	authTokens, ok := output["authTokens"].([]any)
	if !ok || len(authTokens) != 1 {
		t.Fatalf("authTokens = %#v, want one token; payload=%s", output["authTokens"], payload)
	}
	token, ok := authTokens[0].(map[string]any)
	if !ok {
		t.Fatalf("authTokens[0] = %#v, want object; payload=%s", authTokens[0], payload)
	}
	if got := token["tokenSecret"]; got != "search-source-token" {
		t.Fatalf("tokenSecret = %#v, want search-source-token; payload=%s", got, payload)
	}
	if _, exists := token["token"]; exists {
		t.Fatalf("payload contains obsolete token key: %s", payload)
	}
}

func TestSearchSourceTLSFieldsAcceptAPIDefaults(t *testing.T) {
	var resp resource.SchemaResponse
	NewSearchSourceResource().Schema(context.Background(), resource.SchemaRequest{}, &resp)

	tlsAttribute, ok := resp.Schema.Attributes["tls"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("tls schema type = %T, want schema.SingleNestedAttribute", resp.Schema.Attributes["tls"])
	}
	for _, name := range []string{"cert_path", "min_version", "priv_key_path"} {
		attribute, ok := tlsAttribute.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s schema type = %T, want schema.StringAttribute", name, tlsAttribute.Attributes[name])
		}
		if !attribute.Optional || !attribute.Computed {
			t.Fatalf("%s Optional=%t Computed=%t, want both true", name, attribute.Optional, attribute.Computed)
		}
		if len(attribute.PlanModifiers) != 1 {
			t.Fatalf("%s has %d plan modifiers, want UseStateForUnknown", name, len(attribute.PlanModifiers))
		}
	}
}
