package export

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	importclient "github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/client"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/discovery"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

type guardRuleExportTestProvider struct {
	provider.CriblioProvider
	baseURL string
}

func (p *guardRuleExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: p.baseURL, BearerToken: "test-token"})}
}

func (p *guardRuleExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewGuardRuleResource}
}

// TestGuardRuleExportEndToEnd runs the real export pipeline over two Guard
// Rules: one that belongs to a ruleset and one that does not. The second is the
// shape that broke TestIntegration_FullExport_Cloud twice — first as HCL missing
// the required `rulesets` argument, then as `rulesets = []`, which the API
// rejects with "At least one ruleset is required". It must be reported as
// skipped rather than exported, while the first exports and plans clean after
// import.
func TestGuardRuleExportEndToEnd(t *testing.T) {
	const prefix = "/api/v1/m/default/lib/sds-rules"
	rules := map[string]string{
		"credit_card":       `{"id":"credit_card","groupId":"default","regex":"/\\b(?:\\d[ -]*?){13,16}\\b/","description":"Card numbers.","rulesets":["pci-compliance"]}`,
		"employee_badge_id": `{"id":"employee_badge_id","groupId":"default","regex":"/\\bEMP-\\d{6}\\b/","description":"Employee badge identifiers."}`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, prefix) {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		id := strings.TrimPrefix(strings.TrimPrefix(req.URL.Path, prefix), "/")
		// The post-test destroy removes the imported rule.
		if req.Method == http.MethodDelete {
			rule, ok := rules[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			delete(rules, id)
			_, _ = w.Write([]byte(`{"items":[` + rule + `],"count":1}`))
			return
		}
		if req.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if id == "" {
			_, _ = w.Write([]byte(`{"items":[` + rules["credit_card"] + `,` + rules["employee_badge_id"] + `],"count":2}`))
			return
		}
		rule, ok := rules[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(`{"items":[` + rule + `],"count":1}`))
	}))
	defer server.Close()

	client := &importclient.Client{REST: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})}
	results := []discovery.Result{{
		TypeName: "criblio_guard_rule",
		Count:    2,
		Identifiers: []map[string]string{
			{"group_id": "default", "id": "credit_card"},
			{"group_id": "default", "id": "employee_badge_id"},
		},
		InventoryComplete: true,
	}}

	out, err := ToResourceItems(context.Background(), client, buildTestRegistry(t), results, []string{"default"}, nil, 1, false, IncludeOverride{}, nil)
	require.NoError(t, err)

	// Only the manageable rule is exported.
	require.Len(t, out.Items, 1)
	item := out.Items[0]
	require.Equal(t, "credit_card", item.Attrs["id"].String)
	require.Equal(t, "pci-compliance", item.Attrs["rulesets"].List[0].String)

	// The other is reported, not silently dropped.
	require.Contains(t, strings.Join(out.ConvertSkipped, "\n"), "employee_badge_id")
	require.Contains(t, strings.Join(out.ConvertSkipped, "\n"), "belongs to no ruleset")

	// The exported configuration round-trips: import, then a clean plan.
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	require.NoError(t, err)
	config := string(file.Bytes())
	require.Regexp(t, `rulesets\s*= \["pci-compliance"\]`, config)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&guardRuleExportTestProvider{baseURL: server.URL}),
		},
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       item.TypeName + "." + item.Name,
				ImportState:        true,
				ImportStateId:      item.ImportID,
				ImportStatePersist: true,
			},
			{Config: config, PlanOnly: true},
		},
	})
}
