package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSourceHTTPRawAuthTokensRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name string
		ext  string
		want string
	}{
		{"empty objects", `[{}, {"authType":null,"token":null,"tokenSecret":null,"description":null,"metadata":null}]`, ""},
		{"empty values", `[{"authType":"","token":"","metadata":[]}]`, ""},
		{"mixed", `[{}, {"authType":"manual","token":"inline-token","description":"manual token","metadata":[{"name":"env","value":"test"}]}, {}, {"authType":"secret","tokenSecret":"stored-token"}]`, `[{"authType":"manual","token":"inline-token","description":"manual token","metadata":[{"name":"env","value":"test"}]},{"authType":"secret","tokenSecret":"stored-token"}]`},
		{"legacy extended token", `[{"token":"legacy-extended-token"}]`, `[{"token":"legacy-extended-token"}]`},
		{"without extended tokens", "", ""},
		{"explicit empty list", `[]`, `[]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := map[string]any{"id": "http-auth-test", "type": "http_raw", "host": "0.0.0.0", "port": 10080, "authTokens": []string{"legacy-token"}}
			if test.ext != "" {
				var extended any
				if err := json.Unmarshal([]byte(test.ext), &extended); err != nil {
					t.Fatal(err)
				}
				fixture["authTokensExt"] = extended
			}
			payload, err := json.Marshal(fixture)
			if err != nil {
				t.Fatal(err)
			}
			for _, model := range []any{&SourceModel{}, &PackSourceModel{}} {
				if err := json.Unmarshal(payload, model); err != nil {
					t.Fatal(err)
				}
				data, err := json.Marshal(model)
				if err != nil {
					t.Fatal(err)
				}
				var got map[string]any
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got["authTokens"], []any{"legacy-token"}) {
					t.Fatalf("%T changed legacy authentication", model)
				}
				var want any
				if test.want != "" {
					if err := json.Unmarshal([]byte(test.want), &want); err != nil {
						t.Fatal(err)
					}
				}
				if !reflect.DeepEqual(got["authTokensExt"], want) {
					t.Errorf("%T: extended tokens were not preserved/normalized as expected", model)
				}
			}
		})
	}
}

func TestNormalizeHTTPRawExtendedAuthTokensPreservesOtherValues(t *testing.T) {
	for _, token := range []map[string]any{
		{"token": "active-token"},
		{"tokenSecret": "stored-secret"},
		{"authType": "manual"},
		{"enabled": false},
		{"futureField": "configured"},
	} {
		raw := map[string]any{"type": "http_raw", "authTokensExt": []any{map[string]any{}, token}}
		normalizeHTTPRawExtendedAuthTokens(raw)
		if !reflect.DeepEqual(raw["authTokensExt"], []any{token}) {
			t.Fatal("normalization removed a populated token entry")
		}
	}
	raw := map[string]any{"type": "http", "authTokensExt": []any{map[string]any{}}}
	normalizeHTTPRawExtendedAuthTokens(raw)
	if len(raw["authTokensExt"].([]any)) != 1 {
		t.Fatal("normalization changed another source type")
	}
}

func TestSourceHTTPRawReadRemovesPriorPlaceholders(t *testing.T) {
	for _, populated := range []bool{false, true} {
		t.Run(fmt.Sprintf("populated=%v", populated), func(t *testing.T) {
			fixture := `{"id":"http-auth-test","type":"http_raw","host":"0.0.0.0","port":10080,"authTokens":["legacy-token"],"authTokensExt":[{}]}`
			if populated {
				fixture = `{"id":"http-auth-test","type":"http_raw","host":"0.0.0.0","port":10080,"authTokens":["legacy-token"],"authTokensExt":[{}, {"token":"active-token"}]}`
			}
			var model SourceModel
			if err := json.Unmarshal([]byte(fixture), &model); err != nil {
				t.Fatal(err)
			}
			model.GroupID = types.StringValue("default")
			want := model.InputHttpRaw.AuthTokensExt
			var raw map[string]any
			if err := json.Unmarshal([]byte(fixture), &raw); err != nil {
				t.Fatal(err)
			}
			// Restore placeholders to simulate state written by an older provider.
			oldTokens, err := SourceAPIValueToTerraformValue(raw["authTokensExt"], types.ListType{ElemType: types.ObjectType{AttrTypes: InputHttpRawAuthTokensExtAttrTypes()}})
			if err != nil {
				t.Fatal(err)
			}
			model.InputHttpRaw.AuthTokensExt = oldTokens.(types.List)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("refresh used method %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := fmt.Fprintf(w, `{"items":[%s]}`, fixture); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			res := &SourceResource{api: newSourceAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"}))}
			ctx := context.Background()
			var schemaResponse resource.SchemaResponse
			res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
			state := tfsdk.State{Schema: schemaResponse.Schema}
			if diags := state.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			response := resource.ReadResponse{State: state}
			res.Read(ctx, resource.ReadRequest{State: state}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			if diags := response.State.Get(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			if !model.InputHttpRaw.AuthTokensExt.Equal(want) {
				t.Fatal("refresh did not remove only the prior placeholders")
			}
			legacy, err := SourceTerraformValueToJSON(model.InputHttpRaw.AuthTokens)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(legacy, []any{"legacy-token"}) {
				t.Fatal("refresh changed active legacy credentials")
			}
		})
	}
}

func TestNormalizeHTTPRawExtendedAuthTokenStatePreservesUnknowns(t *testing.T) {
	attrs := InputHttpRawAuthTokensExtAttrTypes()
	for _, tokens := range []types.List{
		types.ListNull(types.ObjectType{AttrTypes: attrs}),
		types.ListUnknown(types.ObjectType{AttrTypes: attrs}),
		types.ListValueMust(types.ObjectType{AttrTypes: attrs}, nil),
		types.ListValueMust(types.ObjectType{AttrTypes: attrs}, []attr.Value{types.ObjectUnknown(attrs)}),
	} {
		if !normalizeHTTPRawExtendedAuthTokenState(tokens).Equal(tokens) {
			t.Fatal("normalization changed null, empty, or unknown state")
		}
	}
	attributes := map[string]attr.Value{}
	for name, typ := range attrs {
		null, err := SourceTerraformNullValue(typ)
		if err != nil {
			t.Fatal(err)
		}
		attributes[name] = null
	}
	attributes["token"] = types.StringUnknown()
	tokens := types.ListValueMust(types.ObjectType{AttrTypes: attrs}, []attr.Value{types.ObjectValueMust(attrs, attributes)})
	if !normalizeHTTPRawExtendedAuthTokenState(tokens).Equal(tokens) {
		t.Fatal("normalization removed an unknown credential")
	}
}
