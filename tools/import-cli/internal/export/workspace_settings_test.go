package export

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

type workspaceSettingsExportTestProvider struct {
	systemLimitsExportTestProvider
	baseURL string
}

func (p *workspaceSettingsExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: p.baseURL, BearerToken: "test"})}
}

func (p *workspaceSettingsExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewWorkspaceSystemSettingsResource}
}

func TestWorkspaceSettingsExportAndValidation(t *testing.T) {
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_workspace_system_settings")
	if !ok {
		t.Fatal("missing workspace settings registry entry")
	}
	ids := map[string]string{"id": "conf"}
	fixture := json.RawMessage(`{
  "upgradeSettings":{"disableAutomaticUpgrade":false,"packageUrls":[]},
  "workers":{"count":0}, "api":{"headers":{}},
  "sni":{}, "backups":{}, "tls":{"minVersion":"TLSv1.2"},
  "managedFields":["workers.count"]
}`)
	model, err := converter.ConvertRawItemWithIdentifiers(entry, fixture, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var out ExportResult
	if err := appendResourceItemFromModel(&out, entry.TypeName, entry, ids, model, nil, []string{"default", "default_search"}, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 || out.Items[0].ImportID != "conf" {
		t.Fatalf("workspace singleton export: %#v", out.Items)
	}
	item := out.Items[0]
	if got := item.Attrs["upgrade_settings"].Map["disable_automatic_upgrade"]; got.Kind != hcl.KindBool || got.Bool {
		t.Fatalf("lost explicit false: %#v", got)
	}
	if got := item.Attrs["workers"].Map["count"]; got.Kind != hcl.KindNumber || got.Number != 0 {
		t.Fatalf("lost explicit zero: %#v", got)
	}
	if got := item.Attrs["upgrade_settings"].Map["package_urls"]; got.Kind != hcl.KindList || len(got.List) != 0 {
		t.Fatalf("lost empty list: %#v", got)
	}
	for _, field := range []string{"managed_fields", "group_id", "workspace_id"} {
		if _, ok := item.Attrs[field]; ok {
			t.Fatalf("exported provider-owned metadata: %s", field)
		}
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/api/v1/system/settings/conf" {
			t.Errorf("export/import caused an unnecessary write or used the wrong scope: %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"count": 1, "items": []json.RawMessage{fixture}}); err != nil {
			t.Errorf("encode settings: %v", err)
		}
	}))
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&workspaceSettingsExportTestProvider{baseURL: server.URL}),
		},
		Steps: []resource.TestStep{
			{Config: string(file.Bytes()), PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: string(file.Bytes()), ResourceName: item.TypeName + "." + item.Name, ImportState: true, ImportStateId: "conf", ImportStatePersist: true},
			{Config: string(file.Bytes()), PlanOnly: true},
		},
	})
}
