package export

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type eventBreakerExportTestProvider struct{ provider.CriblioProvider }

func (p *eventBreakerExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *eventBreakerExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewEventBreakerRulesetResource, provider.NewPackBreakersResource}
}

func TestEventBreakerJSONOptionsExport(t *testing.T) {
	reg := buildTestRegistry(t)
	for _, typeName := range []string{"criblio_event_breaker_ruleset", "criblio_pack_breakers"} {
		t.Run(typeName, func(t *testing.T) {
			entry, ok := reg.ByTypeName(typeName)
			if !ok {
				t.Fatal("missing breaker registry entry")
			}
			ids := map[string]string{"group_id": "default", "id": "json-test"}
			if typeName == "criblio_pack_breakers" {
				ids["pack"] = "test-pack"
			}
			var steps []resource.TestStep
			for _, parentFields := range [][]string{{"source", "host*"}, {}} {
				fixture := map[string]any{
					"id": "json-test", "lib": "custom",
					"rules": []any{map[string]any{
						"name": "api-events", "condition": "true", "type": "json_array",
						"timestampAnchorRegex": "/^/", "timestamp": map[string]any{"type": "auto"},
						"jsonArrayField": "items.events", "parentFieldsToCopy": parentFields,
						"jsonExtractAll": false,
					}},
				}
				body, err := json.Marshal(fixture)
				if err != nil {
					t.Fatal(err)
				}
				model, err := converter.ConvertRawItemWithIdentifiers(entry, json.RawMessage(body), toRequestParams(ids))
				if err != nil {
					t.Fatal(err)
				}
				var out ExportResult
				if err := appendResourceItemFromModel(&out, typeName, entry, ids, model, nil, []string{"default"}, false, IncludeOverride{}); err != nil {
					t.Fatal(err)
				}
				if len(out.Items) != 1 || len(out.Items[0].Attrs["rules"].List) != 1 {
					t.Fatal("export lost ruleset or rule")
				}
				item := out.Items[0]
				rule := item.Attrs["rules"].List[0].Map
				if got := rule["json_array_field"]; got.Kind != hcl.KindString || got.String != "items.events" {
					t.Fatalf("exported JSON array field = %#v", got)
				}
				if got := rule["json_extract_all"]; got.Kind != hcl.KindBool || got.Bool {
					t.Fatalf("exported JSON extract all = %#v, want false", got)
				}
				fields := rule["parent_fields_to_copy"]
				if fields.Kind != hcl.KindList || len(fields.List) != len(parentFields) {
					t.Fatalf("exported parent fields = %#v, want %v", fields, parentFields)
				}
				for i, field := range parentFields {
					if got := fields.List[i]; got.Kind != hcl.KindString || got.String != field {
						t.Fatalf("exported parent field %d = %#v, want %q", i, got, field)
					}
				}
				file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
				if err != nil {
					t.Fatal(err)
				}
				steps = append(steps, resource.TestStep{Config: string(file.Bytes()), PlanOnly: true, ExpectNonEmptyPlan: true})
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"criblio": providerserver.NewProtocol6WithError(&eventBreakerExportTestProvider{}),
				},
				Steps: steps,
			})
		})
	}
}
