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
)

// Inject the local API client so these Terraform tests never use local credentials.
type notificationDeliveryTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *notificationDeliveryTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
}

func (p *notificationDeliveryTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewNotificationResource}
}

func TestNotificationDeliveryTerraformLifecycle(t *testing.T) {
	for _, configured := range []bool{true, false} {
		t.Run(fmt.Sprintf("configured=%v", configured), func(t *testing.T) {
			var mu sync.Mutex
			var notification map[string]any
			writes := map[string]int{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				wantPath := "/api/v1/m/default/notifications/delivery-test"
				if r.Method == http.MethodPost {
					wantPath = "/api/v1/m/default/notifications"
				}
				if r.URL.Path != wantPath {
					t.Errorf("API path = %s, want %s", r.URL.Path, wantPath)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				switch r.Method {
				case http.MethodPost, http.MethodPatch:
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode request: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					writes[r.Method]++
					if configured {
						wantMode, wantTemplate := "direct", "default-slack"
						if r.Method == http.MethodPatch {
							wantTemplate = "default-email"
						}
						if writes[http.MethodPatch] == 2 {
							wantMode = "policy"
						}
						wantPairs := []any{}
						if wantMode == "direct" {
							wantPairs = append(wantPairs, map[string]any{"targetId": "system_notifications", "templateId": wantTemplate})
						}
						if body["mode"] != wantMode || !reflect.DeepEqual(body["templateTargetPairs"], wantPairs) {
							t.Errorf("delivery payload = %v / %v, want %s / %v", body["mode"], body["templateTargetPairs"], wantMode, wantPairs)
						}
					} else if r.Method == http.MethodPost {
						for _, key := range []string{"mode", "templateTargetPairs"} {
							if _, ok := body[key]; ok {
								t.Errorf("omitted field %s was sent on create", key)
							}
						}
					}
					notification = body
					// Simulate server defaults for legacy configurations.
					if _, ok := notification["mode"]; !ok {
						notification["mode"] = "direct"
					}
					if _, ok := notification["templateTargetPairs"]; !ok {
						notification["templateTargetPairs"] = []any{}
					}
				case http.MethodGet:
					if notification == nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
				case http.MethodDelete:
					notification = nil
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					t.Errorf("unexpected method: %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{notification}}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()

			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
			const resourceName = "criblio_notification.test"
			config := func(disabled bool, mode, template string) string {
				delivery := ""
				if configured {
					pairs := "[]"
					if mode == "direct" {
						pairs = fmt.Sprintf(`[{ target_id = "system_notifications", template_id = %q }]`, template)
					}
					delivery = fmt.Sprintf("mode = %q\ntemplate_target_pairs = %s", mode, pairs)
				}
				return fmt.Sprintf(`resource "criblio_notification" "test" {
  id = "delivery-test"
  group = "default"
  condition = "high-volume"
  conf = { name = "in_http", time_window = "60s", data_volume = "1GB" }
  targets = []
  disabled = %t
  %s
}`, disabled, delivery)
			}
			checks := func(template string) resource.TestCheckFunc {
				check := []resource.TestCheckFunc{resource.TestCheckResourceAttr(resourceName, "mode", "direct")}
				if configured {
					check = append(check,
						resource.TestCheckResourceAttr(resourceName, "template_target_pairs.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "template_target_pairs.0.target_id", "system_notifications"),
						resource.TestCheckResourceAttr(resourceName, "template_target_pairs.0.template_id", template))
				} else {
					check = append(check, resource.TestCheckResourceAttr(resourceName, "template_target_pairs.#", "0"))
				}
				return resource.ComposeAggregateTestCheckFunc(check...)
			}
			createdConfig := config(false, "direct", "default-slack")
			updatedConfig := config(true, "direct", "default-email")
			steps := []resource.TestStep{
				{Config: createdConfig, Check: checks("default-slack")},
				{Config: createdConfig, PlanOnly: true},
				{Config: updatedConfig, Check: checks("default-email")},
				{Config: updatedConfig, PlanOnly: true},
			}
			// Retain the legacy ID and both supported group keys on import.
			for _, id := range []string{"delivery-test", `{"id":"delivery-test","group":"default"}`, `{"id":"delivery-test","group_id":"default"}`} {
				steps = append(steps,
					resource.TestStep{ResourceName: resourceName, ImportState: true, ImportStateId: id, ImportStateVerify: true},
					resource.TestStep{Config: updatedConfig, PlanOnly: true})
			}
			if configured {
				policyConfig := config(true, "policy", "")
				steps = append(steps,
					resource.TestStep{Config: policyConfig, Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "mode", "policy"),
						resource.TestCheckResourceAttr(resourceName, "template_target_pairs.#", "0"))},
					resource.TestStep{Config: policyConfig, PlanOnly: true})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&notificationDeliveryTestProvider{
						CriblioProvider: CriblioProvider{version: "test"}, client: client,
					}),
				},
				Steps: steps,
			})
			mu.Lock()
			defer mu.Unlock()
			wantUpdates := 1
			if configured {
				wantUpdates = 2
			}
			if writes[http.MethodPost] != 1 || writes[http.MethodPatch] != wantUpdates {
				t.Errorf("API writes = %v, want one create and %d updates", writes, wantUpdates)
			}
		})
	}
}
