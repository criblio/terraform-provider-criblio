package export

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type datasetIntelligenceExportTestProvider struct{ provider.CriblioProvider }

func (p *datasetIntelligenceExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *datasetIntelligenceExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewSearchDatasetIntelligenceResource}
}

type datasetIntelligenceTransport struct{ body string }

func (tr datasetIntelligenceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet || req.URL.Path != "/api/v1/ai/settings/dataset-intelligence/logs" || req.URL.RawQuery != "" {
		return nil, fmt.Errorf("unexpected intelligence export request: %s %s", req.Method, req.URL)
	}
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tr.body)), Request: req}, nil
}

func TestSearchDatasetIntelligenceExportConfiguration(t *testing.T) {
	ctx := context.Background()
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_search_dataset_intelligence")
	if !ok {
		t.Fatal("missing intelligence registry entry")
	}
	const raw = `{"datasetId":"logs","status":"partial","content":"# Log analysis","generatedAt":100,"updatedAt":200,"failureReason":"One tool failed"}`
	client := &importclient.Client{REST: restclient.New(restclient.Config{
		BaseURL: "https://example.invalid", BearerToken: "test-token", HTTPClient: &http.Client{Transport: datasetIntelligenceTransport{body: raw}},
	})}
	ids := map[string]string{"dataset_id": "logs"}
	importID, err := generator.BuildImportID(entry.ImportIDFormat, ids)
	if err != nil || importID != "logs" {
		t.Fatalf("invalid intelligence import ID %q: %v", importID, err)
	}
	model, err := converter.Convert(ctx, client, entry, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var result ExportResult
	if err := appendResourceItemFromModel(&result, entry.TypeName, entry, ids, model, nil, nil, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("export returned %d items, want 1", len(result.Items))
	}
	item := result.Items[0]
	if len(item.Attrs) != 1 || item.Attrs["dataset_id"].String != "logs" {
		t.Fatalf("export should configure only dataset_id: %+v", item.Attrs)
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(&datasetIntelligenceExportTestProvider{}),
	}, Steps: []resource.TestStep{{Config: string(file.Bytes()), PlanOnly: true, ExpectNonEmptyPlan: true}}})
}
