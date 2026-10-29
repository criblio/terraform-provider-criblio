package export

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type notificationTemplateExportTestProvider struct{ provider.CriblioProvider }

func (p *notificationTemplateExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *notificationTemplateExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewNotificationTemplateResource}
}

type notificationTemplateTransport struct{ body string }

func (tr notificationTemplateTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet || req.URL.Path != "/api/v1/notification-templates/custom-email" || req.URL.Query().Get("engine") != "handlebars" {
		return nil, fmt.Errorf("unexpected template export request: %s %s", req.Method, req.URL)
	}
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tr.body)), Request: req}, nil
}

func TestNotificationTemplateExportConfiguration(t *testing.T) {
	ctx := context.Background()
	reg := buildTestRegistry(t)
	entry, ok := reg.ByTypeName("criblio_notification_template")
	if !ok || entry.ImportIDFormat != "id" {
		t.Fatalf("template registry entry missing or incorrect: %+v", entry)
	}
	const content = "{\n  \"subject\":\"{{monitor_name}}\",\n  \"body\":\"{{#each labels}}{{@key}}={{this}}{{/each}}\"\n}"
	raw, err := json.Marshal(map[string]any{"id": "custom-email", "type": "smtp", "description": "custom", "content": content, "engine": "handlebars", "createdBy": "test-user", "lastUpdated": int64(42)})
	if err != nil {
		t.Fatal(err)
	}
	client := &importclient.Client{REST: restclient.New(restclient.Config{
		BaseURL: "https://example.invalid", BearerToken: "test-token", HTTPClient: &http.Client{Transport: notificationTemplateTransport{body: `{"items":[` + string(raw) + `],"count":1}`}},
	})}
	ids := map[string]string{"id": "custom-email"}
	fromAPI, err := converter.Convert(ctx, client, entry, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	fromRaw, err := converter.ConvertRawItemWithIdentifiers(entry, raw, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var steps []resource.TestStep
	for _, model := range []any{fromAPI, fromRaw} {
		var result ExportResult
		if err := appendResourceItemFromModel(&result, entry.TypeName, entry, ids, model, nil, nil, false, IncludeOverride{}); err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 1 {
			t.Fatalf("export returned %d resources, want 1", len(result.Items))
		}
		item := result.Items[0]
		if item.Attrs["content"].String != content || item.Attrs["type"].String != "smtp" {
			t.Fatal("export changed template content or target type")
		}
		for _, field := range []string{"engine", "created_by", "last_updated"} {
			if _, ok := item.Attrs[field]; ok {
				t.Errorf("export configured read-only field %s", field)
			}
		}
		file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
		if err != nil {
			t.Fatal(err)
		}
		steps = append(steps, resource.TestStep{Config: string(file.Bytes()), PlanOnly: true, ExpectNonEmptyPlan: true})
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(&notificationTemplateExportTestProvider{}),
	}, Steps: steps})
	for _, id := range []string{"default-email", "default-slack", "default-sns", "default-pagerduty"} {
		if !skipResourceByID(entry.TypeName, map[string]string{"id": id}) {
			t.Errorf("export did not exclude immutable template %s", id)
		}
	}
	if skipResourceByID(entry.TypeName, ids) {
		t.Fatal("export excluded a custom template")
	}
}
