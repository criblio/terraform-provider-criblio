package provider

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestSearchSourceAuthTokensWireFormat(t *testing.T) {
	for _, test := range []struct {
		name, protocol, want string
	}{
		{"raw HTTP", "http_raw", `{"authTokens":[],"authTokensExt":[{"description":"test token","token":"search-source-token"}],"type":"http_raw"}`},
		{"Cribl HTTP", "cribl_http", `{"authTokens":[{"description":"test token","enabled":true,"tokenSecret":"search-source-token"}],"type":"cribl_http"}`},
		{"Splunk HEC", "splunk_hec", `{"authTokens":[{"description":"test token","enabled":true,"token":"search-source-token"}],"type":"splunk_hec"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			model := SearchSourceModel{Type: types.StringValue(test.protocol), AuthTokens: searchSourceTestTokens("search-source-token", types.BoolValue(true), types.StringValue("test token"))}
			payload, err := json.Marshal(model)
			if err != nil {
				t.Fatal(err)
			}
			var got, want any
			if err := json.Unmarshal(payload, &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(test.want), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("payload = %s, want %s", payload, test.want)
			}
		})
	}
}

func searchSourceTestTokens(token string, enabled types.Bool, description types.String) types.List {
	value := types.ObjectValueMust(SearchSourceAuthTokensAttrTypes(), map[string]attr.Value{
		"description": description, "enabled": enabled, "token": types.StringValue(token), "token_secret": types.StringNull(),
	})
	return types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, []attr.Value{value})
}

func TestSearchSourceHTTPRawAuthTokensEmptyAndOmitted(t *testing.T) {
	for _, tokens := range []types.List{
		types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}),
		types.ListUnknown(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}),
		types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, nil),
	} {
		payload, err := json.Marshal(SearchSourceModel{Type: types.StringValue("http_raw"), AuthTokens: tokens})
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]any
		if err := json.Unmarshal(payload, &body); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"authTokens", "authTokensExt"} {
			if tokens.IsNull() || tokens.IsUnknown() {
				if _, ok := body[key]; ok {
					t.Fatalf("unmanaged %s was sent", key)
				}
				continue
			}
			if value, ok := body[key].([]any); !ok || len(value) != 0 {
				t.Fatalf("%s = %#v, want []", key, body[key])
			}
		}
	}
}

func TestSearchSourceHTTPRawRejectsDisabledToken(t *testing.T) {
	_, err := json.Marshal(SearchSourceModel{Type: types.StringValue("http_raw"), AuthTokens: searchSourceTestTokens("token", types.BoolValue(false), types.StringNull())})
	if err == nil || !strings.Contains(err.Error(), "remove the token entry") {
		t.Fatalf("error = %v, want disabled-token diagnostic", err)
	}
}

func TestSearchSourceRejectsNullTokenObject(t *testing.T) {
	tokens := types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, []attr.Value{types.ObjectNull(SearchSourceAuthTokensAttrTypes())})
	_, err := json.Marshal(SearchSourceModel{Type: types.StringValue("http_raw"), AuthTokens: tokens})
	if err == nil || !strings.Contains(err.Error(), "non-null token object") {
		t.Fatalf("error = %v, want null-object diagnostic", err)
	}
}

func TestSearchSourceDecodeAuthTokens(t *testing.T) {
	for _, test := range []struct{ name, body, token string }{
		{"legacy raw", `{"type":"http_raw","authTokens":["raw-token"]}`, "raw-token"},
		{"extended raw", `{"type":"http_raw","authTokens":[],"authTokensExt":[{}, {"token":"raw-token"}]}`, "raw-token"},
		{"mirrored raw", `{"type":"http_raw","authTokens":["raw-token"],"authTokensExt":[{"token":"raw-token"}]}`, "raw-token"},
		{"authoritative extended raw", `{"type":"http_raw","authTokens":["stale-legacy-token"],"authTokensExt":[{"token":"raw-token"}]}`, "raw-token"},
		{"null extended raw falls back to legacy", `{"type":"http_raw","authTokens":["raw-token"],"authTokensExt":null}`, "raw-token"},
		{"Cribl HTTP", `{"type":"cribl_http","authTokens":[{"tokenSecret":"stored-secret"}]}`, "stored-secret"},
		{"Cribl HTTP resolved secret", `{"type":"cribl_http","authTokens":[{"tokenSecret":"stored-secret","token":"resolved-token"}]}`, "stored-secret"},
		{"Splunk HEC", `{"type":"splunk_hec","authTokens":[{"token":"hec-token"}]}`, "hec-token"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var model SearchSourceModel
			if err := json.Unmarshal([]byte(test.body), &model); err != nil {
				t.Fatal(err)
			}
			if want := searchSourceTestTokens(test.token, types.BoolNull(), types.StringNull()); !model.AuthTokens.Equal(want) {
				t.Fatalf("decoded tokens = %s, want %s", model.AuthTokens, want)
			}
		})
	}
}

func TestSearchSourceEmptyExtendedTokensOverrideLegacy(t *testing.T) {
	var model SearchSourceModel
	if err := json.Unmarshal([]byte(`{"type":"http_raw","authTokens":["stale-token"],"authTokensExt":[]}`), &model); err != nil {
		t.Fatal(err)
	}
	want := types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, nil)
	if !model.AuthTokens.Equal(want) {
		t.Fatal("empty authoritative list resurrected legacy tokens")
	}
}

func TestSearchSourceRawStoredSecretRoundTrip(t *testing.T) {
	for _, body := range []string{
		`{"type":"http_raw","authTokensExt":[{"authType":"secret","tokenSecret":"http-secret"}]}`,
		`{"type":"http_raw","authTokensExt":[{"tokenSecret":"http-secret"}]}`,
		`{"type":"http_raw","authTokens":[{"authType":"secret","tokenSecret":"http-secret","token":"resolved-value"}]}`,
		`{"type":"http_raw","authTokens":["resolved-value"],"authTokensExt":[{"authType":"secret","tokenSecret":"http-secret","token":"resolved-value"}]}`,
	} {
		var model SearchSourceModel
		if err := json.Unmarshal([]byte(body), &model); err != nil {
			t.Fatal(err)
		}
		attributes := model.AuthTokens.Elements()[0].(types.Object).Attributes()
		if !attributes["token"].IsNull() || attributes["token_secret"].(types.String).ValueString() != "http-secret" {
			t.Fatal("read must retain the stored-secret ID without a resolved token")
		}
		payload, err := json.Marshal(model)
		if err != nil {
			t.Fatal(err)
		}
		var wire map[string]any
		if err := json.Unmarshal(payload, &wire); err != nil {
			t.Fatal(err)
		}
		token := wire["authTokensExt"].([]any)[0].(map[string]any)
		if token["authType"] != "secret" || token["tokenSecret"] != "http-secret" || token["token"] != nil {
			t.Fatal("write must send the stored-secret reference, not an inline token")
		}
		if len(wire["authTokens"].([]any)) != 0 {
			t.Fatal("legacy token list must be cleared when writing extended tokens")
		}
	}
}

func TestSearchSourceRawMixedAuthentication(t *testing.T) {
	var model SearchSourceModel
	if err := json.Unmarshal([]byte(`{"type":"http_raw","authTokensExt":[{"authType":"manual","tokenSecret":"","token":"inline-token"},{"authType":"secret","tokenSecret":"http-secret","token":"resolved-value"}]}`), &model); err != nil {
		t.Fatal(err)
	}
	if len(model.AuthTokens.Elements()) != 2 {
		t.Fatal("mixed inline and stored-secret entries must both be retained")
	}
	inline := model.AuthTokens.Elements()[0].(types.Object).Attributes()
	if inline["token"].(types.String).ValueString() != "inline-token" || !inline["token_secret"].IsNull() {
		t.Fatal("empty manual secret selector must not become a stored-secret reference")
	}
	if _, err := json.Marshal(model); err != nil {
		t.Fatal(err)
	}
	prior := model.AuthTokens
	if err := json.Unmarshal([]byte(`{"type":"http_raw","authTokensExt":[{"token":"********"},{"authType":"secret","tokenSecret":"changed-secret","token":"********"}]}`), &model); err != nil {
		t.Fatal(err)
	}
	merged := mergeSearchSourceAuthTokenState(model.AuthTokens, prior, "http_raw")
	if merged.Elements()[0].(types.Object).Attributes()["token"].(types.String).ValueString() != "inline-token" || merged.Elements()[1].(types.Object).Attributes()["token_secret"].(types.String).ValueString() != "changed-secret" {
		t.Fatal("masked inline tokens must preserve plaintext while reference drift remains visible")
	}
}

func TestSearchSourceAuthTokenChoices(t *testing.T) {
	for _, test := range []struct {
		name, protocol string
		token          map[string]any
	}{
		{"both", "http_raw", map[string]any{"token": "inline", "tokenSecret": "secret"}},
		{"neither", "http_raw", map[string]any{}},
		{"empty reference", "http_raw", map[string]any{"tokenSecret": ""}},
		{"unsupported protocol", "cribl_http", map[string]any{"tokenSecret": "secret"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := searchSourceAuthTokensToAPI(map[string]any{"type": test.protocol, "authTokens": []any{test.token}}); err == nil {
				t.Fatal("invalid authentication token was accepted")
			}
		})
	}
}

func TestSearchSourceAuthConfigValidation(t *testing.T) {
	ctx := context.Background()
	r := &SearchSourceResource{}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	objectType := schemaResponse.Schema.Type().TerraformType(ctx).(tftypes.Object)
	tokenType := objectType.AttributeTypes["auth_tokens"].(tftypes.List).ElementType.(tftypes.Object)
	for _, test := range []struct {
		name    string
		values  map[string]any
		invalid bool
	}{
		{"inline", map[string]any{"token": "inline"}, false},
		{"reference", map[string]any{"token_secret": "http-secret"}, false},
		{"both", map[string]any{"token": "inline", "token_secret": "http-secret"}, true},
		{"neither", map[string]any{}, true},
		{"unknown", map[string]any{"token_secret": tftypes.UnknownValue}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := collectorPlanTestObject(objectType, map[string]any{
				"type": "http_raw", "auth_tokens": []tftypes.Value{collectorPlanTestObject(tokenType, test.values)},
			})
			var response resource.ValidateConfigResponse
			r.ValidateConfig(ctx, resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: raw}}, &response)
			if response.Diagnostics.HasError() != test.invalid {
				t.Fatalf("validation diagnostics: %v", response.Diagnostics)
			}
		})
	}
}

func TestSearchSourceAuthTokenStatePreservesSecretsAndDefaults(t *testing.T) {
	for _, protocol := range []string{"http_raw", "cribl_http", "splunk_hec"} {
		for _, secret := range []string{"********", "#42:iv:ciphertext#"} {
			if protocol == "cribl_http" && secret != "********" {
				continue
			}
			prior := searchSourceTestTokens("new-plaintext", types.BoolNull(), types.StringNull())
			api := searchSourceTestTokens(secret, types.BoolValue(true), types.StringValue("API default"))
			if got := mergeSearchSourceAuthTokenState(api, prior, protocol); !got.Equal(prior) {
				t.Fatalf("%s state drifted after masked response", protocol)
			}
		}
	}
	prior := searchSourceTestTokens("old-token", types.BoolValue(true), types.StringNull())
	api := searchSourceTestTokens("new-token", types.BoolNull(), types.StringNull())
	want := searchSourceTestTokens("new-token", types.BoolValue(true), types.StringNull())
	if got := mergeSearchSourceAuthTokenState(api, prior, "http_raw"); !got.Equal(want) {
		t.Fatal("read lost enabled=true or hid plaintext token drift")
	}
	empty := types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, nil)
	if got := mergeSearchSourceAuthTokenState(empty, prior, "http_raw"); !got.Equal(empty) {
		t.Fatal("read preserved removed tokens")
	}
	if got := mergeSearchSourceAuthTokenState(api, types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}), "http_raw"); !got.Equal(api) {
		t.Fatal("import lost tokens")
	}
}

func TestSearchSourceAuthTokenStatePreservesNullAndEmpty(t *testing.T) {
	null := types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()})
	empty := types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, nil)
	for _, protocol := range []string{"http_raw", "cribl_http", "splunk_hec"} {
		for _, prior := range []types.List{null, empty} {
			if got := mergeSearchSourceAuthTokenState(empty, prior, protocol); !got.Equal(prior) {
				t.Fatalf("%s: cleared tokens = %s, want %s", protocol, got, prior)
			}
		}
	}
}

func TestSearchSourceAuthTokensAreTerraformManaged(t *testing.T) {
	var resp resource.SchemaResponse
	NewSearchSourceResource().Schema(context.Background(), resource.SchemaRequest{}, &resp)
	attribute, ok := resp.Schema.Attributes["auth_tokens"].(schema.ListNestedAttribute)
	if !ok || !attribute.Optional || !attribute.Computed {
		t.Fatalf("auth_tokens must remain optional and computed for legacy state compatibility: %#v", resp.Schema.Attributes["auth_tokens"])
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
