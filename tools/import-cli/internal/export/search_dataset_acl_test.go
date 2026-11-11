package export

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

type datasetACLExportTestProvider struct {
	systemLimitsExportTestProvider
	baseURL string
}

func (p *datasetACLExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: p.baseURL, BearerToken: "test"})}
}

func (p *datasetACLExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewSearchDatasetACLResource}
}

func TestDatasetACLExportImportCleanPlan(t *testing.T) {
	var mu sync.Mutex
	userPresent := true
	posts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		body := `{"items":[],"count":0}`
		switch req.Method {
		case http.MethodGet:
			if req.URL.Path == "/api/v1/m/default_search/search/datasets/logs/acl" {
				if userPresent {
					body = `{"items":[{"user":"client-id@clients","perms":[{"type":"datasets","gid":"default_search","id":"logs","policy":"DatasetRead"}]}],"count":1}`
				}
			} else if req.URL.Path != "/api/v1/m/default_search/search/datasets/logs/acl/teams" {
				t.Errorf("unexpected export/import path: %s", req.URL.Path)
			}
		case http.MethodPost:
			posts++
			if req.URL.Path != "/api/v1/m/default_search/search/datasets/logs/acl/apply" {
				t.Errorf("unexpected cleanup action: %s", req.URL.Path)
			}
			userPresent = false
		default:
			t.Errorf("unexpected method: %s", req.Method)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("write fixture: %v", err)
		}
	}))
	defer server.Close()
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_search_dataset_acl")
	if !ok {
		t.Fatal("missing ACL registry entry")
	}
	ids := map[string]string{"dataset_id": "logs"}
	client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}
	model, err := converter.Convert(context.Background(), client, entry, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var out ExportResult
	if err := appendResourceItemFromModel(&out, entry.TypeName, entry, ids, model, nil, []string{"default_search"}, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 || out.Items[0].ImportID != "logs" {
		t.Fatalf("unexpected ACL export: %#v", out.Items)
	}
	item := out.Items[0]
	if got := item.Attrs["teams"]; got.Kind != hcl.KindMap || len(got.Map) != 0 {
		t.Fatalf("export dropped empty teams map: %#v", got)
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	config := string(file.Bytes())
	if !strings.Contains(config, "client-id@clients") || strings.Contains(config, "var.") {
		t.Fatalf("principal ID was treated as a secret: %s", config)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&datasetACLExportTestProvider{baseURL: server.URL}),
		},
		Steps: []resource.TestStep{
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config, ResourceName: item.TypeName + "." + item.Name, ImportState: true, ImportStateId: "logs", ImportStatePersist: true},
			{Config: config, PlanOnly: true},
			{Config: config},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	if posts != 1 {
		t.Fatalf("export/import/unchanged apply caused writes beyond the single destroy action: %d", posts)
	}
}
