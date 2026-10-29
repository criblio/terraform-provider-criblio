package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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

type collectorSavedStateTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *collectorSavedStateTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
	resp.DataSourceData = &ProviderClients{RC: p.client}
}

func (p *collectorSavedStateTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewCollectorResource, NewPackCollectorResource}
}

func (p *collectorSavedStateTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewCollectorDataSource, NewCollectorsDataSource, NewPackCollectorDataSource}
}

func TestCollectorSavedStateTerraformLifecycle(t *testing.T) {
	for _, scope := range []string{"group", "pack"} {
		for _, configured := range []struct {
			name string
			hcl  string
			json string
		}{
			{"computed", "", "null"},
			{"legacy map", `tomap({cursor = "configured"})`, `{"cursor":"configured"}`},
			{"legacy literal", `{cursor = "configured"}`, `{"cursor":"configured"}`},
			{"nested object", `{seed = {value = 1, enabled = true, nullable = null, files = ["one", null, 2]}}`, `{"seed":{"value":1,"enabled":true,"nullable":null,"files":["one",null,2]}}`},
		} {
			for _, removal := range []string{"omitted", "null"} {
				t.Run(scope+"/"+configured.name+"/"+removal, func(t *testing.T) {
					var mu sync.Mutex
					var saved map[string]any
					var checkpoint any
					var expected any
					var unconfigured bool
					if err := json.Unmarshal([]byte(configured.json), &expected); err != nil {
						t.Fatal(err)
					}
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						mu.Lock()
						defer mu.Unlock()
						w.Header().Set("Content-Type", "application/json")
						if req.URL.Path == "/api/v1/m/default/packs" {
							if _, err := fmt.Fprint(w, `{"items":[{"id":"example"}]}`); err != nil {
								t.Error(err)
							}
							return
						}
						switch req.Method {
						case http.MethodPost, http.MethodPatch:
							saved = nil
							if err := json.NewDecoder(req.Body).Decode(&saved); err != nil {
								t.Error(err)
								w.WriteHeader(http.StatusBadRequest)
								return
							}
							want := expected
							if unconfigured {
								want = nil
							}
							if !reflect.DeepEqual(saved["savedState"], want) {
								t.Errorf("%s savedState = %#v, want %#v", req.Method, saved["savedState"], want)
							}
							if want == nil {
								if _, exists := saved["savedState"]; exists {
									t.Error("computed checkpoint must be omitted from the request")
								}
							}
						case http.MethodDelete:
							saved = nil
							w.WriteHeader(http.StatusNoContent)
							return
						}
						if saved == nil {
							http.NotFound(w, req)
							return
						}
						response := make(map[string]any, len(saved))
						for name, value := range saved {
							response[name] = value
						}
						if checkpoint != nil {
							response["savedState"] = checkpoint
						}
						if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{response}, "count": 1}); err != nil {
							t.Error(err)
						}
					}))
					defer server.Close()
					typeName := "criblio_collector"
					pack := ""
					importID := `{"group_id":"default","id":"events"}`
					if scope == "pack" {
						typeName = "criblio_pack_collector"
						pack = `pack = "example"`
						importID = `{"group_id":"default","id":"events","pack":"example"}`
					}
					address := typeName + ".test"
					config := func(description, state string) string {
						if state != "" {
							state = "saved_state = " + state
						}
						result := fmt.Sprintf(`resource %q "test" {
  group_id = "default"
  id = "events"
  %s
  input_collector_database = {
    description = %q
    collector = {type = "database"}
    %s
  }
}`, typeName, pack, description, state)
						if scope == "group" {
							result += fmt.Sprintf(`
data "criblio_collector" "test" {
  group_id = "default"
  id = %s.id
}
data "criblio_collectors" "test" {
  group_id = "default"
  depends_on = [%s]
}`, address, address)
						}
						return result
					}
					setCheckpoint := func(value int) func() {
						return func() {
							mu.Lock()
							defer mu.Unlock()
							checkpoint = map[string]any{"databaseCollectorState": map[string]any{"data": map[string]any{"id": map[string]any{"value": value}}}}
						}
					}
					checks := []resource.TestCheckFunc{}
					if configured.hcl == "" {
						checks = append(checks, resource.TestCheckResourceAttr(address, "input_collector_database.saved_state.databaseCollectorState.data.id.value", "3"))
					} else if strings.HasPrefix(configured.name, "legacy") {
						checks = append(checks, resource.TestCheckResourceAttr(address, "input_collector_database.saved_state.cursor", "configured"))
					} else {
						checks = append(checks, resource.TestCheckResourceAttr(address, "input_collector_database.saved_state.seed.value", "1"))
					}
					if scope == "group" {
						checks = append(checks,
							resource.TestCheckResourceAttr("data.criblio_collector.test", "input_collector_database.saved_state.databaseCollectorState.data.id.value", "3"),
							resource.TestCheckResourceAttr("data.criblio_collectors.test", "items.0.input_collector_database.saved_state.databaseCollectorState.data.id.value", "3"))
					}
					removedState := ""
					if removal == "null" {
						removedState = "null"
					}
					removalAction := plancheck.ResourceActionUpdate
					if configured.hcl == "" {
						removalAction = plancheck.ResourceActionNoop
					}
					client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
							"criblio": providerserver.NewProtocol6WithError(&collectorSavedStateTestProvider{CriblioProvider: CriblioProvider{version: "test"}, client: client}),
						},
						Steps: []resource.TestStep{
							{Config: config("initial", configured.hcl)},
							{PreConfig: setCheckpoint(3), Config: config("initial", configured.hcl),
								ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
								Check:            resource.ComposeAggregateTestCheckFunc(checks...)},
							{PreConfig: setCheckpoint(4), Config: config("updated", configured.hcl)},
							{Config: config("updated", configured.hcl), PlanOnly: true},
							{ResourceName: address, ImportState: true, ImportStateId: importID, ImportStateVerify: true,
								ImportStateVerifyIgnore: []string{"input_collector_database.saved_state"}},
							{PreConfig: func() {
								mu.Lock()
								defer mu.Unlock()
								unconfigured = true
							}, Config: config("updated", removedState),
								ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, removalAction)}},
								Check:            resource.TestCheckResourceAttr(address, "input_collector_database.saved_state.databaseCollectorState.data.id.value", "4")},
							{PreConfig: setCheckpoint(5), Config: config("updated", removedState),
								ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
								Check:            resource.TestCheckResourceAttr(address, "input_collector_database.saved_state.databaseCollectorState.data.id.value", "5")},
							{PreConfig: func() {
								mu.Lock()
								defer mu.Unlock()
								checkpoint = nil
							}, Config: config("updated", removedState),
								ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
								Check:            resource.TestCheckNoResourceAttr(address, "input_collector_database.saved_state")},
						},
					})
				})
			}
		}
	}
}
