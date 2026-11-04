package export

import (
	"context"
	"encoding/json"
	"regexp"
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

type groupSettingsExportTestProvider struct{ provider.CriblioProvider }

func (p *groupSettingsExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *groupSettingsExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewGroupSystemSettingsResource}
}

func TestGroupSystemSettingsLimitsExport(t *testing.T) {
	reg := buildTestRegistry(t)
	entry, ok := reg.ByTypeName("criblio_group_system_settings")
	if !ok {
		t.Fatal("missing group settings registry entry")
	}
	ids := map[string]string{"group_id": "test-fleet"}
	body := json.RawMessage(`{"limits":{"minFreeSpace":"2 GB","maxPQSize":"1 TB","metricsMaxDiskSpace":"1 GB","maxMetrics":10000,"metricsMaxCardinality":1000,"lookupMaxSize":"100 MB","lookupMaxTotalSize":"1 GB","samples":{"maxSize":"2 MB"},"enableMetricsPersistence":false,"metricsDropList":[],"metricsGCPeriod":"60s","kvStoreWALMaxBytes":1024}}`)
	model, err := converter.ConvertRawItemWithIdentifiers(entry, body, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var out ExportResult
	if err := appendResourceItemFromModel(&out, entry.TypeName, entry, ids, model, nil, []string{"test-fleet"}, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 {
		t.Fatalf("exported %d settings, want 1", len(out.Items))
	}
	item := out.Items[0]
	limits := item.Attrs["limits"].Map
	for key, want := range map[string]string{"min_free_space": "2 GB", "max_pq_size": "1 TB", "metrics_max_disk_space": "1 GB", "lookup_max_size": "100 MB", "lookup_max_total_size": "1 GB", "metrics_gc_period": "60s"} {
		if got := limits[key]; got.Kind != hcl.KindString || got.String != want {
			t.Fatalf("exported %s = %#v, want %q", key, got, want)
		}
	}
	for key, want := range map[string]float64{"max_metrics": 10000, "metrics_max_cardinality": 1000, "kv_store_wal_max_bytes": 1024} {
		if got := limits[key]; got.Kind != hcl.KindNumber || got.Number != want {
			t.Fatalf("exported %s = %#v, want %v", key, got, want)
		}
	}
	if got := limits["samples"].Map["max_size"]; got.Kind != hcl.KindString || got.String != "2 MB" {
		t.Fatalf("exported sample limit = %#v", got)
	}
	if got := limits["enable_metrics_persistence"]; got.Kind != hcl.KindBool || got.Bool {
		t.Fatalf("exported enable_metrics_persistence = %#v, want false", got)
	}
	if got := limits["metrics_drop_list"]; got.Kind != hcl.KindList || len(got.List) != 0 {
		t.Fatalf("exported metrics_drop_list = %#v, want empty list", got)
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&groupSettingsExportTestProvider{}),
		},
		Steps: []resource.TestStep{
			{Config: `resource "criblio_group_system_settings" "minimal" {
  group_id = "test-fleet"
  limits = { min_free_space = "2 GB" }
}`, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: string(file.Bytes()), PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: `resource "criblio_group_system_settings" "invalid_size" {
  group_id = "test-fleet"
  limits = { min_free_space = "invalid" }
}`, PlanOnly: true, ExpectError: regexp.MustCompile("Invalid Attribute Value Match")},
			{Config: `resource "criblio_group_system_settings" "invalid_range" {
  group_id = "test-fleet"
  limits = { azure_blob_state_flush_interval = 49 }
}`, PlanOnly: true, ExpectError: regexp.MustCompile("Invalid Attribute Value")},
		},
	})
}
