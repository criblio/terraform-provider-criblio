package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type sourceHTTPAuthTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *sourceHTTPAuthTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
}

func (p *sourceHTTPAuthTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSourceResource}
}

func TestSourceHTTPRawAuthTerraformLifecycle(t *testing.T) {
	for _, test := range []struct {
		name string
		ext  []any
		hcl  string
	}{
		{"placeholders", []any{map[string]any{}, map[string]any{"authType": nil, "token": nil, "tokenSecret": nil, "description": nil, "metadata": nil}}, ""},
		{"mixed", []any{map[string]any{}, map[string]any{"authType": "secret", "tokenSecret": "stored-token"}, map[string]any{}}, `auth_tokens_ext = [{ auth_type = "secret", token_secret = "stored-token" }]`},
		{"no extended tokens", nil, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var mu sync.Mutex
			writes := 0
			source := map[string]any{
				"id": "http-auth-test", "type": "http_raw", "host": "0.0.0.0", "port": float64(10080),
				"description": "original", "authTokens": []any{"legacy-token"},
			}
			if test.ext != nil {
				source["authTokensExt"] = test.ext
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				if r.URL.Path != "/api/v1/m/default/system/inputs/http-auth-test" {
					t.Errorf("unexpected API path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				switch r.Method {
				case http.MethodGet:
				case http.MethodPatch:
					writes++
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode update: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if !reflect.DeepEqual(body["authTokens"], []any{"legacy-token"}) {
						t.Error("update changed legacy credentials")
					}
					if test.hcl == "" {
						if _, ok := body["authTokensExt"]; ok {
							t.Error("update sent empty extended-token placeholders")
						}
					} else if !reflect.DeepEqual(body["authTokensExt"], []any{map[string]any{"authType": "secret", "tokenSecret": "stored-token"}}) {
						t.Error("update changed meaningful extended credentials")
					}
					source["description"] = body["description"]
				case http.MethodDelete:
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					t.Errorf("unexpected API method: %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{source}}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()
			config := func(description string) string {
				return fmt.Sprintf(`resource "criblio_source" "test" {
  id = "http-auth-test"
  group_id = "default"
  input_http_raw = {
    id = "http-auth-test"
    host = "0.0.0.0"
    port = 10080
    description = %q
    auth_tokens = ["legacy-token"]
    %s
  }
}`, description, test.hcl)
			}
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&sourceHTTPAuthTestProvider{
						CriblioProvider: CriblioProvider{version: "test"}, client: client,
					}),
				},
				Steps: []resource.TestStep{
					{
						Config: config("original") + `
import {
  to = criblio_source.test
  id = "{\"group_id\":\"default\",\"id\":\"http-auth-test\"}"
}`,
						Check: func(_ *terraform.State) error {
							mu.Lock()
							defer mu.Unlock()
							if writes != 0 {
								return fmt.Errorf("importing an unchanged source performed %d updates", writes)
							}
							return nil
						},
					},
					{Config: config("original"), PlanOnly: true},
					{Config: config("updated"), Check: resource.TestCheckResourceAttr("criblio_source.test", "input_http_raw.description", "updated")},
					{Config: config("updated"), PlanOnly: true},
					{ResourceName: "criblio_source.test", ImportState: true, ImportStateId: `{"group_id":"default","id":"http-auth-test"}`, ImportStateVerify: true},
					{Config: config("updated"), PlanOnly: true},
				},
			})
			mu.Lock()
			defer mu.Unlock()
			if writes != 1 {
				t.Errorf("API updates = %d, want only the description update", writes)
			}
		})
	}
}
