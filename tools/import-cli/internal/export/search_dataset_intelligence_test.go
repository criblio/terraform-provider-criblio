package export

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/converter"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

type datasetIntelligenceExportTestProvider struct{ provider.CriblioProvider }

func (p *datasetIntelligenceExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func TestDefaultExportWithoutDatasetIntelligencePrerequisites(t *testing.T) {
	for _, message := range []string{"AI is disabled", "AI consent not accepted"} {
		t.Run(message, func(t *testing.T) {
			var probes atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				body := `{"items":[]}`
				switch req.URL.Path {
				case "/api/v1/products/stream/groups", "/api/v1/products/edge/groups":
					body = `{"items":[{"id":"default","name":"default"}]}`
				case "/api/v1/m/default_search/search/datasets":
					body = `{"items":[{"id":"logs"},{"id":"other"}]}`
				case "/api/v1/ai/settings/dataset-intelligence/logs", "/api/v1/ai/settings/dataset-intelligence/other":
					probes.Add(1)
					w.WriteHeader(http.StatusForbidden)
					body = fmt.Sprintf(`{"error":%q}`, message)
				}
				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})}
			ctx := context.Background()
			reg := buildTestRegistry(t)
			// Use the same unfiltered discovery and export stages as the default CLI.
			results, err := discovery.Discover(ctx, client, reg, nil, nil, nil, false)
			require.NoError(t, err)
			foundIntelligence := false
			for _, result := range results {
				require.NoError(t, result.Err, result.TypeName)
				if result.TypeName == "criblio_search_dataset_intelligence" {
					foundIntelligence = true
					require.Zero(t, result.Count)
					require.Empty(t, result.Identifiers)
					require.True(t, result.InventoryComplete)
				}
			}
			require.True(t, foundIntelligence)
			require.EqualValues(t, 1, probes.Load(), "workspace-wide unavailability should stop further probes")
			result, err := ToResourceItems(ctx, client, reg, results, []string{"default", "default_search"}, nil, 1, false, IncludeOverride{}, nil)
			require.NoError(t, err)
			require.NotEmpty(t, result.Items, "unrelated resources should still be exported")
			foundGroup := false
			for _, item := range result.Items {
				require.NotEqual(t, "criblio_search_dataset_intelligence", item.TypeName)
				if item.TypeName == "criblio_group" {
					foundGroup = true
				}
			}
			require.True(t, foundGroup)
			outputDir := t.TempDir()
			err = generator.WriteAllModuleDirectories(outputDir, result.Items, hcl.DefaultResourceBlockOptions(), nil)
			require.NoError(t, err, "default export should reach Terraform file generation")
		})
	}
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
