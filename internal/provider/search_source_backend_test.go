package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSearchSourceBackendFullReplacement(t *testing.T) {
	for _, test := range []struct {
		name, protocol, field, token string
		blocked                      bool
	}{
		{"encrypted inline", "http_raw", "authTokensExt", `{"token":"#42:iv:ciphertext#","metadata":[{"name":"env","value":"prod"}]}`, false},
		{"raw reference", "http_raw", "authTokensExt", `{"authType":"secret","tokenSecret":"folder/secret","token":"********"}`, false},
		{"HEC reference", "splunk_hec", "authTokens", `{"authType":"secret","tokenSecret":"folder/secret","token":"********","enabled":false,"allowedIndexesAtToken":["main"]}`, false},
		{"masked reference", "http_raw", "authTokensExt", `{"authType":"secret","tokenSecret":"****/******"}`, true},
		{"masked inline", "http_raw", "authTokensExt", `{"token":"********"}`, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var token map[string]any
			if err := json.Unmarshal([]byte(test.token), &token); err != nil {
				t.Fatal(err)
			}
			existing := map[string]any{"id": "source", "type": test.protocol, test.field: []any{token}}
			writes := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path != "/api/v1/m/default_search/search/local_search/sources/source" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				if r.Method == http.MethodPatch {
					writes++
					// v1 stores the request as the complete source representation.
					existing = nil
					if err := json.NewDecoder(r.Body).Decode(&existing); err != nil {
						t.Error(err)
					}
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{existing}}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			api := newSearchSourceAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"}))
			_, err := api.updateWithPreservedAuth(context.Background(), SearchSourceModel{ID: types.StringValue("source"), Type: types.StringValue(test.protocol), Description: types.StringValue("changed"), AuthTokens: types.ListNull(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()})})
			if test.blocked {
				if err == nil || writes != 0 {
					t.Fatalf("masked credentials must block before PATCH: writes=%d, err=%v", writes, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if writes != 1 {
				t.Fatalf("writes=%d, want 1", writes)
			}
			entries, ok := existing[test.field].([]any)
			if !ok || len(entries) != 1 {
				t.Fatal("replacement lost authentication")
			}
			entry := entries[0].(map[string]any)
			if existing["description"] != "changed" {
				t.Fatal("update did not apply the configured change")
			}
			if token["tokenSecret"] == nil && entry["token"] != token["token"] {
				t.Fatal("changed the preserved inline credential")
			}
			if entry["tokenSecret"] != nil && entry["token"] != nil {
				t.Fatal("secret reference retained a resolved token")
			}
			for _, field := range []string{"metadata", "allowedIndexesAtToken", "enabled"} {
				if token[field] != nil && entry[field] == nil {
					t.Fatalf("lost token property %s", field)
				}
			}
		})
	}
}

func TestSearchSourceBackendMaskedReferenceState(t *testing.T) {
	for _, protocol := range []string{"http_raw", "splunk_hec", "cribl_http"} {
		t.Run(protocol, func(t *testing.T) {
			field := "token_secret"
			if protocol == "cribl_http" {
				field = "token"
			}
			makeTokens := func(value string) types.List {
				attributes := map[string]attr.Value{
					"token": types.StringNull(), "token_secret": types.StringNull(),
					"description": types.StringNull(), "enabled": types.BoolNull(),
				}
				attributes[field] = types.StringValue(value)
				object := types.ObjectValueMust(SearchSourceAuthTokensAttrTypes(), attributes)
				return types.ListValueMust(types.ObjectType{AttrTypes: SearchSourceAuthTokensAttrTypes()}, []attr.Value{object})
			}
			prior := makeTokens("folder/secret")
			if got := mergeSearchSourceAuthTokenState(makeTokens("****/******"), prior, protocol); !got.Equal(prior) {
				t.Fatal("masked reference replaced a known ID")
			}
		})
	}
}

func TestSearchSourceBackendHECSecretAndManualSelection(t *testing.T) {
	for _, protocol := range []string{"http_raw", "splunk_hec"} {
		field := "authTokens"
		if protocol == "http_raw" {
			field = "authTokensExt"
		}
		var model SearchSourceModel
		body := `{"type":"` + protocol + `","` + field + `":[{"authType":"manual","token":"inline","tokenSecret":"inactive"},{"authType":"secret","tokenSecret":"folder/secret","token":"resolved"}]}`
		if err := json.Unmarshal([]byte(body), &model); err != nil {
			t.Fatal(err)
		}
		if !model.AuthTokens.Elements()[0].(types.Object).Attributes()["token_secret"].IsNull() {
			t.Fatal("manual selection used an inactive reference")
		}
		encoded, err := json.Marshal(model)
		if err != nil {
			t.Fatal(err)
		}
		var wire map[string]any
		if err := json.Unmarshal(encoded, &wire); err != nil {
			t.Fatal(err)
		}
		entry := wire[field].([]any)[1].(map[string]any)
		if entry["authType"] != "secret" || entry["tokenSecret"] != "folder/secret" || entry["token"] != nil {
			t.Fatal("HEC/Raw reference contract mismatch")
		}
	}
}
