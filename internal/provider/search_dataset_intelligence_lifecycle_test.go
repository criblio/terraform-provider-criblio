package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

type datasetIntelligenceTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *datasetIntelligenceTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
	resp.DataSourceData = &ProviderClients{RC: p.client}
}

func (p *datasetIntelligenceTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSearchDatasetIntelligenceResource}
}

func (p *datasetIntelligenceTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewSearchDatasetIntelligenceDataSource}
}

func TestSearchDatasetIntelligenceTerraformLifecycle(t *testing.T) {
	for _, adopt := range []bool{false, true} {
		t.Run(fmt.Sprintf("import_%t", adopt), func(t *testing.T) {
			var mu sync.Mutex
			creates, deletes := 0, 0
			files := map[string]map[string]any{}
			if adopt {
				files["logs"] = map[string]any{"datasetId": "logs", "generatedAt": int64(100), "status": "processing"}
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				const base = "/api/v1/ai/settings/dataset-intelligence/"
				if !strings.HasPrefix(req.URL.Path, base) || req.URL.RawQuery != "" {
					t.Errorf("incorrect intelligence route: %s", req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				id := strings.TrimPrefix(req.URL.Path, base)
				switch req.Method {
				case http.MethodPost:
					var body map[string]any
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil || len(body) != 0 {
						t.Errorf("generation must send no configurable envelope fields: body=%v, err=%v", body, err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					creates++
					files[id] = map[string]any{"datasetId": id, "status": "processing", "generatedAt": int64(100)}
					w.WriteHeader(http.StatusCreated)
				case http.MethodGet:
					if files[id] == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
				case http.MethodDelete:
					deletes++
					delete(files, id)
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					t.Errorf("unexpected intelligence mutation: %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				if err := json.NewEncoder(w).Encode(files[id]); err != nil {
					t.Errorf("encode envelope: %v", err)
				}
			}))
			defer server.Close()
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
			const name = "criblio_search_dataset_intelligence.test"
			config := func(id string) string {
				return fmt.Sprintf(`resource "criblio_search_dataset_intelligence" "test" {
  dataset_id = %q
}
data "criblio_search_dataset_intelligence" "test" {
  dataset_id = criblio_search_dataset_intelligence.test.dataset_id
  depends_on = [criblio_search_dataset_intelligence.test]
}`, id)
			}
			refresh := func(status, content, failure string) resource.TestStep {
				return resource.TestStep{
					PreConfig: func() {
						mu.Lock()
						defer mu.Unlock()
						files["logs"] = map[string]any{"datasetId": "logs", "generatedAt": int64(100), "status": status}
						if content != "" {
							files["logs"]["content"] = content
							files["logs"]["updatedAt"] = int64(200)
						}
						if failure != "" {
							files["logs"]["failureReason"] = failure
						}
					},
					Config:           config("logs"),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionNoop)}},
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(name, "status", status),
						resource.TestCheckResourceAttr("data.criblio_search_dataset_intelligence.test", "status", status),
						func() resource.TestCheckFunc {
							if content == "" {
								return resource.ComposeAggregateTestCheckFunc(resource.TestCheckNoResourceAttr(name, "content"), resource.TestCheckNoResourceAttr(name, "updated_at"))
							}
							return resource.TestCheckResourceAttr(name, "content", content)
						}(),
						func() resource.TestCheckFunc {
							if failure == "" {
								return resource.TestCheckNoResourceAttr(name, "failure_reason")
							}
							return resource.TestCheckResourceAttr(name, "failure_reason", failure)
						}(),
					),
				}
			}
			steps := []resource.TestStep{
				{Config: config("logs"), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr(name, "status", "processing"), resource.TestCheckNoResourceAttr(name, "content"))},
				{Config: config("logs"), PlanOnly: true},
				refresh("complete", "# Logs\n\nFields: `status`, `message`. Unicode: ✓", ""),
				refresh("partial", "# Partial analysis", "One tool failed"),
				refresh("failed", "", "Analysis failed"),
				refresh("processing", "", ""),
				{ResourceName: name, ImportState: true, ImportStateId: "logs", ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "dataset_id"},
				{Config: config("other"), ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionReplace)}}},
				{Config: config("other"), PlanOnly: true},
				{PreConfig: func() {
					mu.Lock()
					defer mu.Unlock()
					delete(files, "other")
				}, Config: config("other"), ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionCreate)}}},
				// Removing the resource disables intelligence, leaving the dataset alone.
				{Config: "terraform {}"},
			}
			if adopt {
				steps[0].Config += `
import {
  to = criblio_search_dataset_intelligence.test
  id = "logs"
}`
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
				"criblio": providerserver.NewProtocol6WithError(&datasetIntelligenceTestProvider{CriblioProvider: CriblioProvider{version: "test"}, client: client}),
			}, Steps: steps})
			mu.Lock()
			defer mu.Unlock()
			wantCreates := 3
			if adopt {
				wantCreates = 2
			}
			if creates != wantCreates || deletes != 2 || len(files) != 0 {
				t.Errorf("creates=%d deletes=%d remaining=%d; want %d, 2, 0", creates, deletes, len(files), wantCreates)
			}
		})
	}
}
