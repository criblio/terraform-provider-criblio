package export

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/generator"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

type searchMonitorTestProvider struct {
	provider.CriblioProvider
	client *restclient.Client
}

func (p *searchMonitorTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	clients := &provider.ProviderClients{RC: p.client}
	resp.ResourceData = clients
	resp.DataSourceData = clients
}

// Discover and import both query representations, then plan the exported HCL
// against the same API. Any attempted write during import indicates drift.
func TestSearchMonitorExportImport(t *testing.T) {
	const collection = "/api/v1/m/default_search/alert/monitors"
	const simple = `{"id":"simple","name":"Simple","product":"stream","enabled":false,"firing_after":0,"ok_after":0,"schedule_interval_seconds":60,"query":"up","rules":[],"params":{}}`
	const detailed = `{"id":"detailed","name":"Detailed","description":"Structured query","product":"stream","enabled":true,"firing_after":10,"ok_after":5,"schedule_interval_seconds":60,"isDefault":true,
"query":{"metricName":"events","labelFilters":[],"timeRange":"1m","operation":{"operation":"sum"}},"params":{"threshold":1.5,"active":false,"name":"worker"},
"rules":[{"name":"empty","includedTags":{},"excludedTags":{},"showOnChart":false,"conditions":[]},{"name":"range","includedTags":{"env":["prod","qa"]},"excludedTags":{},"showOnChart":true,"conditions":[{"enabled":true,"labels":{},"condition":{"type":"within_range","range":{"min":0,"max":10}}}]}],
"notification_policies":[],"notifications_enabled":false,"silences":[],"sqlOverride":{"instant":"SELECT 1","range":"SELECT 2"}}`
	var renamed atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		items := ""
		switch req.URL.Path {
		case "/api/v1/products/stream/groups", "/api/v1/products/edge/groups":
		case collection:
			switch req.URL.Query().Get("offset") {
			case "0":
				fmt.Fprintf(w, `{"items":[%s],"offset":0,"totalCount":2}`, simple)
			case "1":
				fmt.Fprintf(w, `{"items":[%s],"offset":1,"totalCount":2}`, detailed)
			default:
				t.Errorf("unexpected list request %s", req.URL)
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		case collection + "/simple":
			items = simple
		case collection + "/detailed":
			if req.Method == http.MethodPatch {
				var update struct {
					Name      string
					IsDefault bool `json:"isDefault"`
				}
				if err := json.NewDecoder(req.Body).Decode(&update); err != nil || !update.IsDefault || update.Name != "Renamed" {
					t.Errorf("built-in monitor update = %+v, err = %v", update, err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				renamed.Store(true)
			}
			items = detailed
			if renamed.Load() {
				items = strings.Replace(detailed, `"Detailed"`, `"Renamed"`, 1)
			}
		default:
			t.Errorf("unexpected path %s", req.URL)
			http.NotFound(w, req)
			return
		}
		if req.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if req.Method != http.MethodGet && req.Method != http.MethodPatch {
			t.Errorf("unexpected %s %s during import", req.Method, req.URL)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, `{"items":[%s]}`, items)
	}))
	defer server.Close()
	ctx := t.Context()
	client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}
	reg := buildTestRegistry(t)
	results, err := discovery.Discover(ctx, client, reg, []string{"criblio_search_monitor"}, nil, nil, false)
	require.NoError(t, err)
	out, err := ToResourceItems(ctx, client, reg, results, nil, nil, 1, false, IncludeOverride{}, nil)
	require.NoError(t, err)
	require.Len(t, out.Items, 2)
	custom, err := ToResourceItems(ctx, client, reg, results, nil, nil, 1, true, IncludeOverride{}, nil)
	require.NoError(t, err)
	require.Len(t, custom.Items, 1)
	require.Equal(t, "simple", custom.Items[0].ImportID)
	config := string(generator.ImportBlocksBytes(out.Items))
	for _, item := range out.Items {
		block, err := hcl.ResourceBlockBytes(item.TypeName, item.Name, item.Attrs, hcl.DefaultResourceBlockOptions())
		require.NoError(t, err)
		require.NotContains(t, string(block), "is_default")
		config += string(block)
	}
	config += `
data "criblio_search_monitor" "one" { id = "detailed" }
data "criblio_search_monitors" "all" {}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&searchMonitorTestProvider{client: client.REST}),
		},
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.criblio_search_monitor.one", "name", "Detailed"),
				resource.TestCheckResourceAttr("data.criblio_search_monitors.all", "items.#", "2"),
				resource.TestCheckResourceAttr("data.criblio_search_monitors.all", "items.1.product", "stream"),
				resource.TestCheckResourceAttr("data.criblio_search_monitor.one", "is_default", "true"),
			)},
			{Config: config, PlanOnly: true},
			{Config: strings.Replace(config, `"Detailed"`, `"Renamed"`, 1), Check: resource.TestCheckResourceAttr(
				"data.criblio_search_monitor.one", "is_default", "true",
			)},
		},
	})
	require.True(t, renamed.Load(), "expected an update to the imported built-in monitor")
}
