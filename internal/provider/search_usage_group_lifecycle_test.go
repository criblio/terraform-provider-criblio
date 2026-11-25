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
)

type usageGroupTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *usageGroupTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
	resp.DataSourceData = &ProviderClients{RC: p.client}
}

func (p *usageGroupTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSearchUsageGroupResource}
}

func (p *usageGroupTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewSearchUsageGroupDataSource}
}

// canonicalRuleText mirrors the live API: every limit type carries a fixed
// description, and some carry a fixed metric. The API fills them in when absent
// and replaces whatever the client sends.
var canonicalRuleText = map[string][2]string{
	"maxResultsPerSearch":               {"Max number of events that could be returned by a single job.", ""},
	"maxRunningTimePerSearch":           {"Max number of seconds a single job can run.", "running_time"},
	"maxConcurrentSearches":             {"Max number of concurrent searches", ""},
	"maxConcurrentAdhocSearchesPerUser": {"Max number of concurrent adhoc searches that could be run by a user.", ""},
}

// applyCanonicalRuleText rewrites a stored group's rules the way the API does.
func applyCanonicalRuleText(group map[string]any) {
	rules, ok := group["rules"].(map[string]any)
	if !ok {
		return
	}
	for key, raw := range rules {
		rule, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		text, known := canonicalRuleText[key]
		if !known {
			continue
		}
		rule["description"] = text[0]
		if text[1] != "" {
			rule["metric"] = text[1]
		}
		rules[key] = rule
	}
}

// TestSearchUsageGroupTerraformLifecycleAPIOwnedRuleFields pins the behaviour
// that drove excluding rule `description` and `metric` from the schema: the API
// returns them whether or not they were sent, so the plan has to stay clean
// anyway, including after an import. Runs against a fake API, so it needs no
// credentials.
func TestSearchUsageGroupTerraformLifecycleAPIOwnedRuleFields(t *testing.T) {
	var mu sync.Mutex
	groups := map[string]map[string]any{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		const base = "/api/v1/m/default_search/search/usage-groups"
		if !strings.HasPrefix(req.URL.Path, base) {
			t.Errorf("unexpected route: %s", req.URL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := strings.TrimPrefix(strings.TrimPrefix(req.URL.Path, base), "/")

		write := func(group map[string]any) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{group}, "count": 1})
		}

		switch {
		case req.Method == http.MethodPost && id == "":
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Errorf("decode create body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if body["users"] == nil {
				body["users"] = map[string]any{}
			}
			applyCanonicalRuleText(body)
			body["usersCount"] = len(body["users"].(map[string]any))
			groups[body["id"].(string)] = body
			write(body)
		case req.Method == http.MethodPatch && id != "":
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Errorf("decode update body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if body["users"] == nil {
				body["users"] = map[string]any{}
			}
			applyCanonicalRuleText(body)
			body["usersCount"] = len(body["users"].(map[string]any))
			groups[id] = body
			write(body)
		case req.Method == http.MethodGet && id != "":
			group, ok := groups[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			write(group)
		case req.Method == http.MethodDelete && id != "":
			group, ok := groups[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			delete(groups, id)
			write(group)
		default:
			t.Errorf("unhandled %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
	const name = "criblio_search_usage_group.test"

	config := func(limit int) string {
		return fmt.Sprintf(`resource "criblio_search_usage_group" "test" {
  id          = "power_users"
  description = "Members with elevated search limits."
  enabled     = true

  rules = {
    max_results_per_search = {
      type  = "maxResultsPerSearch"
      limit = %d
    }
    max_running_time_per_search = {
      type  = "maxRunningTimePerSearch"
      limit = 300
    }
    max_concurrent_adhoc_searches_per_user = {
      type  = "maxConcurrentAdhocSearchesPerUser"
      limit = "50%%"
    }
  }
}

data "criblio_search_usage_group" "test" {
  id         = criblio_search_usage_group.test.id
  depends_on = [criblio_search_usage_group.test]
}`, limit)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&usageGroupTestProvider{
				CriblioProvider: CriblioProvider{version: "test"},
				client:          client,
			}),
		},
		Steps: []resource.TestStep{
			{
				Config: config(100000),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", "power_users"),
					resource.TestCheckResourceAttr(name, "rules.max_results_per_search.limit", "100000"),
					resource.TestCheckResourceAttr(name, "rules.max_concurrent_adhoc_searches_per_user.limit", "50%"),
					// The API owns these; they are filled from its response.
					resource.TestCheckResourceAttr(name, "rules.max_results_per_search.description", "Max number of events that could be returned by a single job."),
					resource.TestCheckResourceAttr(name, "rules.max_running_time_per_search.metric", "running_time"),
				),
			},
			// The API returned canonical description and metric on create; the
			// plan must still be empty.
			{Config: config(100000), PlanOnly: true},
			{
				Config: config(75000),
				Check:  resource.TestCheckResourceAttr(name, "rules.max_results_per_search.limit", "75000"),
			},
			{Config: config(75000), PlanOnly: true},
			// Import reconstructs state from a response carrying the canonical
			// fields, and must still leave a clean plan.
			{ResourceName: name, ImportState: true, ImportStateId: "power_users", ImportStateVerify: true},
			{Config: config(75000), PlanOnly: true},
		},
	})
}
