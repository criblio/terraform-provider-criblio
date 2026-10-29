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

type notificationTemplateTestProvider struct {
	CriblioProvider
	client *restclient.Client
}

func (p *notificationTemplateTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &ProviderClients{RC: p.client}
	resp.DataSourceData = &ProviderClients{RC: p.client}
}

func (p *notificationTemplateTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewNotificationTemplateResource}
}

func (p *notificationTemplateTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewNotificationTemplateDataSource, NewNotificationTemplatesDataSource}
}

func TestNotificationTemplateTerraformLifecycle(t *testing.T) {
	for _, targetType := range []string{"webhook", "pager_duty", "slack", "sns", "smtp"} {
		t.Run(targetType, func(t *testing.T) {
			const original = "{\n  \"subject\": \"[{{status}}] {{monitor_name}}\",\n  \"body\": \"{{#each labels}}{{@key}}={{this}} {{/each}}\"\n}\n"
			var mu sync.Mutex
			creates, updates, deletes := 0, 0, 0
			templates := map[string]map[string]any{
				"default-email": {"id": "default-email", "type": "smtp", "engine": "handlebars", "content": `{"subject":"{{monitor_name}}"}`},
			}
			if targetType == "sns" {
				// Adopt an existing API object through declarative import, including
				// its response wrapper and API-owned metadata.
				templates["custom"] = map[string]any{"id": "custom", "type": "sns", "content": original, "description": "original", "engine": "handlebars", "createdBy": "test-user", "lastUpdated": int64(0)}
			}
			const base = "/api/v1/notification-templates"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				if req.URL.Query().Get("engine") != "handlebars" || (req.URL.Path != base && !strings.HasPrefix(req.URL.Path, base+"/")) {
					t.Errorf("incorrect template API route: %s", req.URL)
					w.WriteHeader(http.StatusBadRequest)
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
					for _, field := range []string{"engine", "createdBy", "lastUpdated"} {
						if _, ok := body[field]; ok {
							t.Errorf("request sent API-owned field %s", field)
						}
					}
					if req.Method == http.MethodPost {
						creates++
						var ok bool
						id, ok = body["id"].(string)
						if !ok || id == "" || body["type"] == nil {
							t.Error("create omitted identity or target type")
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						templates[id] = body
						templates[id]["createdBy"] = "test-user"
					} else {
						updates++
						for field := range body {
							if field != "content" && field != "description" {
								t.Errorf("PATCH sent immutable or unsupported field %s", field)
							}
						}
						if templates[id] == nil {
							t.Errorf("update addressed nonexistent template %s", id)
							w.WriteHeader(http.StatusNotFound)
							return
						}
						for key, value := range body {
							templates[id][key] = value
						}
					}
					templates[id]["engine"] = "handlebars"
					templates[id]["lastUpdated"] = int64(creates + updates)
				case http.MethodDelete:
					deletes++
					if templates[id] == nil || id == "default-email" {
						t.Errorf("deleted unexpected template %s", id)
					}
					delete(templates, id)
					w.WriteHeader(http.StatusNoContent)
					return
				case http.MethodGet:
				default:
					t.Errorf("unexpected method %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				items := []map[string]any{}
				if req.URL.Path == base && req.Method == http.MethodGet {
					for _, template := range templates {
						items = append(items, template)
					}
				} else if template := templates[id]; template != nil {
					items = append(items, template)
				} else {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"items": items, "count": len(items)}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
			const name = "criblio_notification_template.test"
			const updated = "Updated: {{#if description}}{{description}}{{else}}No description{{/if}}\nUnicode: ✓"
			config := func(id, kind, description, content string) string {
				return fmt.Sprintf(`resource "criblio_notification_template" "test" {
  id = %q
  type = %q
  description = %q
  content = %q
}
data "criblio_notification_template" "by_id" {
  id = criblio_notification_template.test.id
  depends_on = [criblio_notification_template.test]
}
data "criblio_notification_template" "builtin" {
  id = "default-email"
}
data "criblio_notification_templates" "all" {
  depends_on = [criblio_notification_template.test]
}`, id, kind, description, content)
			}
			steps := []resource.TestStep{
				{Config: config("custom", targetType, "original", original), Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "content", original),
					resource.TestCheckResourceAttr(name, "engine", "handlebars"),
					resource.TestCheckResourceAttr(name, "created_by", "test-user"),
					resource.TestCheckResourceAttr("data.criblio_notification_template.by_id", "content", original),
					resource.TestCheckResourceAttr("data.criblio_notification_template.builtin", "id", "default-email"),
					resource.TestCheckResourceAttr("data.criblio_notification_templates.all", "items.#", "2"),
				)},
				{Config: config("custom", targetType, "original", original), PlanOnly: true},
				{Config: config("custom", targetType, "updated", updated), Check: resource.TestCheckResourceAttr(name, "content", updated)},
				{Config: config("custom", targetType, "updated", updated), PlanOnly: true},
				{ResourceName: name, ImportState: true, ImportStateId: "custom", ImportStateVerify: true},
				{Config: config("custom", targetType, "updated", updated), PlanOnly: true},
			}
			if targetType == "smtp" {
				steps = append(steps,
					resource.TestStep{Config: config("custom", "sns", "updated", updated), ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionReplace)}}},
					resource.TestStep{Config: config("renamed", "sns", "updated", updated), ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionReplace)}}},
					resource.TestStep{Config: config("renamed", "sns", "updated", updated), PlanOnly: true})
			}
			if targetType == "webhook" {
				steps = append([]resource.TestStep{steps[0], {
					PreConfig: func() {
						mu.Lock()
						defer mu.Unlock()
						templates["custom"]["content"] = "Changed outside Terraform"
						delete(templates["custom"], "description")
					},
					Config:           config("custom", targetType, "original", original),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, plancheck.ResourceActionUpdate)}},
					Check:            resource.TestCheckResourceAttr(name, "content", original),
				}}, steps[1:]...)
			}
			if targetType == "sns" {
				steps[0].Config += `
import {
  to = criblio_notification_template.test
  id = "custom"
}`
				withoutDescription := func(content string) string {
					return strings.Replace(config("custom", targetType, "updated", content), "  description = \"updated\"\n", "", 1)
				}
				steps = append(steps,
					resource.TestStep{Config: withoutDescription(updated), PlanOnly: true},
					resource.TestStep{Config: withoutDescription("New content"), Check: resource.TestCheckResourceAttr(name, "description", "updated")},
					resource.TestStep{Config: config("custom", targetType, "", "New content"), Check: resource.TestCheckResourceAttr(name, "description", "")},
					resource.TestStep{Config: config("custom", targetType, "", "New content"), PlanOnly: true})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&notificationTemplateTestProvider{CriblioProvider: CriblioProvider{version: "test"}, client: client}),
				}, Steps: steps,
			})
			mu.Lock()
			defer mu.Unlock()
			wantCreates := 1
			wantDeletes := 1
			wantUpdates := 1
			if targetType == "webhook" {
				wantUpdates = 2
			}
			if targetType == "smtp" {
				wantCreates = 3
				wantDeletes = 3
			}
			if targetType == "sns" {
				wantCreates = 0
				wantUpdates = 3
			}
			if creates != wantCreates || updates != wantUpdates || deletes != wantDeletes || len(templates) != 1 {
				t.Errorf("API calls create=%d update=%d delete=%d remaining=%d; want %d, %d, %d, 1", creates, updates, deletes, len(templates), wantCreates, wantUpdates, wantDeletes)
			}
		})
	}
}
