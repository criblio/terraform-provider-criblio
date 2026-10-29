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
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type searchDatasetIDTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *searchDatasetIDTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
}

func (p *searchDatasetIDTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSearchDatasetResource}
}

func TestSearchDatasetIDTerraformLifecycle(t *testing.T) {
	var mu sync.Mutex
	datasets := make(map[string]map[string]any)
	creates, updates, deletes := 0, 0, 0
	const base = "/api/v1/m/default_search/search/datasets"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if req.URL.Path != base && !strings.HasPrefix(req.URL.Path, base+"/") {
			t.Errorf("unexpected API path %s", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id := strings.TrimPrefix(req.URL.Path, base+"/")
		switch req.Method {
		case http.MethodPost, http.MethodPatch:
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			bodyID, ok := body["id"].(string)
			if !ok || bodyID == "" {
				t.Error("request missing configured ID")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.Method == http.MethodPost {
				id = bodyID
				creates++
			} else {
				updates++
				if id != bodyID || datasets[id] == nil {
					t.Errorf("update addressed %q with body ID %q", id, bodyID)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			datasets[id] = body
		case http.MethodGet:
			if datasets[id] == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		case http.MethodDelete:
			if datasets[id] == nil {
				t.Errorf("delete addressed missing dataset %q", id)
			}
			delete(datasets, id)
			deletes++
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			t.Errorf("unexpected method %s", req.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{datasets[id]}}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
	config := func(id, description string) string {
		return fmt.Sprintf(`resource "criblio_search_dataset" "test" {
  s3_dataset = {
    id = %q
    description = %q
    provider_id = "S3"
    bucket = "logs"
  }
}
output "dataset_id" {
  value = criblio_search_dataset.test.id
}`, id, description)
	}
	const address = "criblio_search_dataset.test"
	checks := func(id string, action plancheck.ResourceActionType) resource.ConfigPlanChecks {
		return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
			plancheck.ExpectKnownValue(address, tfjsonpath.New("id"), knownvalue.StringExact(id)),
			plancheck.ExpectKnownOutputValue("dataset_id", knownvalue.StringExact(id)),
			plancheck.ExpectResourceAction(address, action),
		}}
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&searchDatasetIDTestProvider{
				CriblioProvider: CriblioProvider{version: "test"}, client: client,
			}),
		},
		Steps: []resource.TestStep{
			{Config: config("logs", "original"), ConfigPlanChecks: checks("logs", plancheck.ResourceActionCreate),
				Check: resource.TestCheckResourceAttr(address, "id", "logs")},
			{Config: config("logs", "original"), PlanOnly: true},
			{Config: config("logs", "updated"), ConfigPlanChecks: checks("logs", plancheck.ResourceActionUpdate),
				Check: resource.TestCheckResourceAttr(address, "s3_dataset.description", "updated")},
			{Config: config("logs", "updated"), PlanOnly: true},
			{ResourceName: address, ImportState: true, ImportStateId: "logs", ImportStateVerify: true},
			{Config: config("logs", "updated"), PlanOnly: true},
			{Config: config("new-logs", "updated"), ConfigPlanChecks: checks("new-logs", plancheck.ResourceActionReplace),
				Check: resource.TestCheckResourceAttr(address, "id", "new-logs")},
			{Config: config("new-logs", "updated"), PlanOnly: true},
			{ResourceName: address, ImportState: true, ImportStateId: "new-logs", ImportStateVerify: true},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	if creates != 2 || updates != 1 || deletes != 2 || len(datasets) != 0 {
		t.Errorf("API calls: creates=%d updates=%d deletes=%d remaining=%d; want 2, 1, 2, 0", creates, updates, deletes, len(datasets))
	}
}
