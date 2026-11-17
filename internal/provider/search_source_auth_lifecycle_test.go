package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

type searchSourceAuthTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *searchSourceAuthTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
}

func (p *searchSourceAuthTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSearchSourceResource}
}

func TestSearchSourceHTTPRawAuthTerraformLifecycle(t *testing.T) {
	var mu sync.Mutex
	var source map[string]any
	updates := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path != "/api/v1/m/default_search/search/local_search/sources" && r.URL.Path != "/api/v1/m/default_search/search/local_search/sources/auth-test" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost, http.MethodPatch:
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Enforce the actual Raw HTTP input contract independently of the
			// provider's shared Search schema. Objects in authTokens reproduce
			// the backend failure reported by the user.
			if legacy, ok := body["authTokens"].([]any); ok {
				for _, token := range legacy {
					if _, ok := token.(string); !ok {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(w, `{"status":"error","message":"s.startsWith is not a function"}`)
						return
					}
				}
			}
			if r.Method == http.MethodPost {
				source = body
				source["authTokens"] = []any{}
				source["authTokensExt"] = []any{}
			} else {
				updates++
				extended, ok := body["authTokensExt"].([]any)
				if !ok {
					t.Error("token update omitted authTokensExt")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				legacy, ok := body["authTokens"].([]any)
				if !ok || len(legacy) != 0 {
					t.Error("update did not clear legacy tokens")
				}
				for _, item := range extended {
					token, ok := item.(map[string]any)
					if !ok {
						t.Error("extended token is not an object")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					value, ok := token["token"].(string)
					if !ok || value == "" || token["tokenSecret"] != nil || token["enabled"] != nil {
						t.Error("update did not send a valid inline Raw HTTP token")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					wantToken := map[int]string{1: "alpha", 2: "beta", 4: "gamma", 6: "gamma"}[updates]
					if value != wantToken {
						t.Errorf("update %d did not send the configured token", updates)
					}
					// Simulate responses that cannot reveal plaintext credentials.
					if updates == 2 {
						token["token"] = "#42:iv:ciphertext#"
					} else {
						token["token"] = "********"
					}
				}
				source = body
				// BaseHttpIn.init backports the authoritative extended list into
				// authTokens before InputMgr returns the initialized config.
				mirrored := make([]any, 0, len(extended))
				for _, item := range extended {
					mirrored = append(mirrored, item.(map[string]any)["token"])
				}
				source["authTokens"] = mirrored
			}
			source["disabled"] = false
		case http.MethodGet:
			if source == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case http.MethodDelete:
			source = nil
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			t.Errorf("unexpected method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{source}}); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	config := func(tokens string) string {
		return fmt.Sprintf(`resource "criblio_search_source" "test" {
  id = "auth-test"
  type = "http_raw"
  host = "0.0.0.0"
  port = 20001
  description = "auth test"
  tls = {
    disabled = false
    cert_path = "$CRIBL_CLOUD_CRT"
    priv_key_path = "$CRIBL_CLOUD_KEY"
    min_version = "TLSv1.2"
  }
  %s
}`, tokens)
	}
	checks := func(count int, token string, wantUpdates int) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("criblio_search_source.test", "auth_tokens.#", fmt.Sprint(count)),
			func(state *terraform.State) error {
				mu.Lock()
				defer mu.Unlock()
				if updates != wantUpdates {
					return fmt.Errorf("updates = %d, want %d", updates, wantUpdates)
				}
				if count != 0 {
					return resource.TestCheckResourceAttr("criblio_search_source.test", "auth_tokens.0.token", token)(state)
				}
				return nil
			},
		)
	}
	alpha := config(`auth_tokens = [{enabled = true, description = "test token", token = "alpha"}]`)
	beta := config(`auth_tokens = [{token = "beta"}]`)
	empty := config(`auth_tokens = []`)
	gamma := config(`auth_tokens = [{enabled = true, token = "gamma"}]`)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&searchSourceAuthTestProvider{
				CriblioProvider: CriblioProvider{version: "test"},
				client:          restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"}),
			}),
		},
		Steps: []resource.TestStep{
			{Config: config(`auth_tokens = null`), Check: checks(0, "", 0)},
			{Config: alpha, Check: checks(1, "alpha", 1)},
			{Config: alpha, PlanOnly: true},
			{Config: beta, Check: checks(1, "beta", 2)},
			{Config: beta, PlanOnly: true},
			{ResourceName: "criblio_search_source.test", ImportState: true, ImportStateId: "auth-test", ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_tokens"}},
			{Config: empty, Check: checks(0, "", 3)},
			{Config: empty, PlanOnly: true},
			{Config: gamma, Check: checks(1, "gamma", 4)},
			{Config: config(`auth_tokens = null`), Check: checks(0, "", 5)},
			{Config: config(`auth_tokens = null`), PlanOnly: true},
			{Config: gamma, Check: checks(1, "gamma", 6)},
			{Config: config(""), Check: checks(0, "", 7)},
			{Config: config(""), PlanOnly: true},
		},
	})
}
