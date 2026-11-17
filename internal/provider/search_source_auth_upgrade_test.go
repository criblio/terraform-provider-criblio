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
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The previous resource computed auth_tokens and did not record ownership.
type legacySearchSourceAuthResource struct{ SearchSourceResource }

func (r *legacySearchSourceAuthResource) Schema(ctx context.Context, req frameworkresource.SchemaRequest, resp *frameworkresource.SchemaResponse) {
	r.SearchSourceResource.Schema(ctx, req, resp)
	attribute := resp.Schema.Attributes["auth_tokens"].(schema.ListNestedAttribute)
	attribute.PlanModifiers = []planmodifier.List{listplanmodifier.UseStateForUnknown()}
	resp.Schema.Attributes["auth_tokens"] = attribute
}

func (r *legacySearchSourceAuthResource) ModifyPlan(context.Context, frameworkresource.ModifyPlanRequest, *frameworkresource.ModifyPlanResponse) {
}
func (r *legacySearchSourceAuthResource) Create(ctx context.Context, req frameworkresource.CreateRequest, resp *frameworkresource.CreateResponse) {
	r.SearchSourceResource.Create(ctx, req, resp)
	if resp.Private != nil {
		resp.Diagnostics.Append(resp.Private.SetKey(ctx, searchSourceAuthManagedKey, nil)...)
	}
}

type legacySearchSourceAuthProvider struct{ searchSourceAuthTestProvider }

func (p *legacySearchSourceAuthProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{func() frameworkresource.Resource { return &legacySearchSourceAuthResource{} }}
}

func TestSearchSourceLegacyAuthUpgradePreservesTokens(t *testing.T) {
	for _, token := range []string{"legacy-token", "new-token"} {
		t.Run(token, func(t *testing.T) { testSearchSourceLegacyAuthUpgradePreservesTokens(t, token, false) })
	}
}

func TestSearchSourceStoredSecretPreservation(t *testing.T) {
	for _, secret := range []string{"legacy-token", "new-token"} {
		t.Run(secret, func(t *testing.T) { testSearchSourceLegacyAuthUpgradePreservesTokens(t, secret, true) })
	}
}

func testSearchSourceLegacyAuthUpgradePreservesTokens(t *testing.T, adoptedToken string, storedSecret bool) {
	tokenField := "token"
	if storedSecret {
		tokenField = "token_secret"
	}
	var mu sync.Mutex
	var source map[string]any
	updates := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
				t.Error(err)
				return
			}
			source["authTokens"] = []any{"legacy-token"}
			source["authTokensExt"] = []any{map[string]any{"token": "legacy-token"}}
			if storedSecret {
				source["authTokensExt"] = []any{map[string]any{"authType": "secret", "tokenSecret": "legacy-token", "token": "resolved-value"}}
			}
			source["disabled"] = false
		case http.MethodPatch:
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
				return
			}
			updates++
			if updates == 1 {
				if entries, ok := body["authTokensExt"].([]any); !ok || len(entries) != 1 {
					t.Error("full-replacement update did not preserve legacy credentials")
				}
			}
			if storedSecret {
				expectedSecret := adoptedToken
				if updates == 1 {
					expectedSecret = "legacy-token"
				}
				if entries, ok := body["authTokensExt"].([]any); ok {
					for _, item := range entries {
						entry := item.(map[string]any)
						if entry["authType"] != "secret" || entry["tokenSecret"] != expectedSecret || entry["token"] != nil {
							t.Error("update converted a stored-secret reference into an inline token")
						}
					}
				}
			}
			// The v1 backend replaces the complete stored source.
			source = body
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{source}}); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
	newProvider := func() *searchSourceAuthTestProvider { return &searchSourceAuthTestProvider{client: client} }
	legacy := map[string]func() (tfprotov6.ProviderServer, error){"criblio": providerserver.NewProtocol6WithError(&legacySearchSourceAuthProvider{*newProvider()})}
	current := map[string]func() (tfprotov6.ProviderServer, error){"criblio": providerserver.NewProtocol6WithError(newProvider())}
	config := func(description, tokens string) string {
		return fmt.Sprintf(`resource "criblio_search_source" "test" {
 id = "upgrade-auth"
 type = "http_raw"
 host = "0.0.0.0"
 port = 20001
 description = %q
 tls = {disabled = false, cert_path = "$CRIBL_CLOUD_CRT", priv_key_path = "$CRIBL_CLOUD_KEY", min_version = "TLSv1.2"}
 %s
}`, description, tokens)
	}
	const name = "criblio_search_source.test"
	steps := []resource.TestStep{
		{ProtoV6ProviderFactories: legacy, Config: config("old", ""), Check: resource.TestCheckResourceAttr(name, "auth_tokens.0."+tokenField, "legacy-token")},
		{ProtoV6ProviderFactories: current, Config: config("old", ""), PlanOnly: true},
		{ProtoV6ProviderFactories: current, Config: config("changed", ""), Check: resource.TestCheckResourceAttr(name, "auth_tokens.0."+tokenField, "legacy-token")},
		{ProtoV6ProviderFactories: current, Config: config("changed", ""), PlanOnly: true},
		{ProtoV6ProviderFactories: current, Config: config("changed", fmt.Sprintf(`auth_tokens = [{%s = %q}]`, tokenField, adoptedToken)), Check: resource.TestCheckResourceAttr(name, "auth_tokens.0."+tokenField, adoptedToken)},
	}
	if storedSecret {
		steps = append([]resource.TestStep{
			steps[0],
			{ProtoV6ProviderFactories: current, ResourceName: name, ImportState: true, ImportStateId: "upgrade-auth", ImportStateVerify: true},
		}, steps[1:]...)
	}
	if adoptedToken == "legacy-token" {
		// Terraform does not persist private ownership metadata on a no-op.
		// Preserve legacy credentials until a write or explicit clear occurs.
		steps = append(steps,
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", ""), Check: resource.TestCheckResourceAttr(name, "auth_tokens.0."+tokenField, "legacy-token")},
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", ""), PlanOnly: true},
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", `auth_tokens = []`), Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", `auth_tokens = []`), PlanOnly: true},
		)
	} else {
		steps = append(steps,
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", ""), Check: resource.TestCheckResourceAttr(name, "auth_tokens.#", "0")},
			resource.TestStep{ProtoV6ProviderFactories: current, Config: config("changed", ""), PlanOnly: true},
		)
	}
	resource.UnitTest(t, resource.TestCase{Steps: steps})
	mu.Lock()
	defer mu.Unlock()
	wantUpdates := 3
	if adoptedToken == "legacy-token" {
		wantUpdates = 2
	}
	if updates != wantUpdates {
		t.Fatalf("updates = %d, want %d", updates, wantUpdates)
	}
	if tokens, ok := source["authTokensExt"].([]any); !ok || len(tokens) != 0 {
		t.Fatal("owned credentials were not cleared")
	}
}
