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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type workerGroupUpgradeTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *workerGroupUpgradeTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
}

func (p *workerGroupUpgradeTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewWorkerGroupUpgradeResource}
}

func (p *workerGroupUpgradeTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func TestWorkerGroupUpgradeTerraformLifecycle(t *testing.T) {
	var mu sync.Mutex
	var bodies []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if req.Method != http.MethodPost || req.URL.Path != "/api/v1/system/distributed/upgrade/hybrid-group" || req.URL.RawQuery != "" {
			t.Errorf("unexpected upgrade request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Errorf("decode upgrade request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, field := range []string{"trigger", "group", "group_id", "jobIds", "versionTo"} {
			if _, ok := body[field]; ok {
				t.Errorf("Terraform-only or operator field %s sent in request: %v", field, body)
			}
		}
		bodies = append(bodies, body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"count":2,"items":["job-%d","other-job-%d"]}`, len(bodies), len(bodies))
	}))
	defer server.Close()
	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
	const name = "criblio_worker_group_upgrade.test"
	config := func(trigger, settings string) string {
		return fmt.Sprintf(`resource "criblio_worker_group_upgrade" "test" {
  group_id = "hybrid-group"
  %s
  %s
}`, trigger, settings)
	}
	check := func(count int, mode string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(name, "job_ids.#", "2"),
			resource.TestCheckResourceAttr(name, "job_ids.0", fmt.Sprintf("job-%d", count)),
			resource.TestCheckResourceAttr(name, "job_ids.1", fmt.Sprintf("other-job-%d", count)),
			resource.TestCheckResourceAttr(name, "upgrade_mode", mode),
			func(_ *terraform.State) error {
				mu.Lock()
				defer mu.Unlock()
				if len(bodies) != count {
					return fmt.Errorf("upgrade submissions = %d; want %d", len(bodies), count)
				}
				return nil
			},
		)
	}
	step := func(cfg string, count int, mode string, action plancheck.ResourceActionType) resource.TestStep {
		return resource.TestStep{Config: cfg, Check: check(count, mode),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, action)}},
		}
	}
	first := config(`trigger = "upgrade-1"`, "")
	second := config(`trigger = "upgrade-2"`, "")
	settings := "upgrade_mode = \"batch\"\nupgrade_percentage = 0.5\nworker_retries = 5\nworker_retry_delay = 1000"
	custom := `package_urls = [{
  package_url = "https://packages.example/cribl-4.20.0-abcdef12-linux-x64.tgz"
  package_hash_url = "https://packages.example/cribl-4.20.0-abcdef12-linux-x64.tgz.sha256"
}]`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&workerGroupUpgradeTestProvider{CriblioProvider: CriblioProvider{version: "test"}, client: client}),
		},
		Steps: []resource.TestStep{
			step(first, 1, "rolling", plancheck.ResourceActionCreate),
			step(first, 1, "rolling", plancheck.ResourceActionNoop),
			{RefreshState: true, Check: check(1, "rolling")},
			step(second, 2, "rolling", plancheck.ResourceActionReplace),
			step(config(`trigger = "upgrade-2"`, settings), 3, "batch", plancheck.ResourceActionReplace),
			step(second, 4, "rolling", plancheck.ResourceActionReplace),
			step(config(`trigger = "upgrade-2"`, custom), 5, "rolling", plancheck.ResourceActionReplace),
			step(second, 6, "rolling", plancheck.ResourceActionReplace),
			step(config("", ""), 7, "rolling", plancheck.ResourceActionReplace),
			{Config: "terraform {}"},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != 7 {
		t.Fatalf("upgrade submissions including destroy = %d; want 7", len(bodies))
	}
	if len(bodies[0]) != 1 || bodies[0]["upgradeMode"] != "rolling" {
		t.Errorf("default request = %v; want only upgradeMode=rolling", bodies[0])
	}
	if bodies[2]["upgradeMode"] != "batch" || bodies[2]["upgradePercentage"] != 0.5 || bodies[2]["workerRetries"] != float64(5) || bodies[2]["workerRetryDelay"] != float64(1000) {
		t.Errorf("configured upgrade settings = %v", bodies[2])
	}
	packages, ok := bodies[4]["packageUrls"].([]any)
	if !ok || len(packages) != 1 {
		t.Fatalf("custom package request = %v", bodies[4])
	}
	if pkg, ok := packages[0].(map[string]any); !ok || pkg["packageUrl"] != "https://packages.example/cribl-4.20.0-abcdef12-linux-x64.tgz" || pkg["packageHashUrl"] != "https://packages.example/cribl-4.20.0-abcdef12-linux-x64.tgz.sha256" {
		t.Errorf("custom package fields = %v", packages[0])
	}
}
