package export

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
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

type systemLimitsExportTestProvider struct{ provider.CriblioProvider }

func (p *systemLimitsExportTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	resp.ResourceData = &provider.ProviderClients{RC: restclient.New(restclient.Config{BaseURL: "https://example.invalid", BearerToken: "test-token"})}
}

func (p *systemLimitsExportTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{provider.NewSystemLimitsResource, provider.NewGroupSystemSettingsResource}
}

func TestSystemLimitsExportAndValidation(t *testing.T) {
	entry, ok := buildTestRegistry(t).ByTypeName("criblio_system_limits")
	if !ok {
		t.Fatal("missing system limits registry entry")
	}
	if entry.ModelTypeName != "SystemLimitsResourceModel" {
		t.Fatalf("wrapper model mapping = %q", entry.ModelTypeName)
	}
	ids := map[string]string{"group_id": "test-fleet", "id": "test-fleet"}
	fixture := json.RawMessage(`{
		"cpuProfileTTL":"10m", "enableMetricsPersistence":false,
		"metricsDirectory":"/var/cribl/metrics", "metricsFieldsBlacklist":[],
		"metricsGCPeriod":"60s", "metricsNeverDropList":[],
		"metricsWorkerIdBlacklist":[], "minFreeSpace":"2 GB",
		"samples":{"maxSize":"2 MB"}, "maxPQSize":"1 TB", "maxEPs":0,
        "lookupMaxSize":"500 MB", "lookupMaxTotalSize":"2 GB", "enableWorkerInMemory":false,
		"eventsMetadataSources":[], "lookupAllowedPaths":[], "metricsDropList":[],
		"managedFields":["lookupAllowedPaths"]
	}`)
	model, err := converter.ConvertRawItemWithIdentifiers(entry, fixture, toRequestParams(ids))
	if err != nil {
		t.Fatal(err)
	}
	var out ExportResult
	if err := appendResourceItemFromModel(&out, entry.TypeName, entry, ids, model, nil, []string{"test-fleet"}, false, IncludeOverride{}); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 1 {
		t.Fatal("export filtered out the group singleton")
	}
	item := out.Items[0]
	if item.ImportID != "test-fleet" {
		t.Fatalf("import ID = %q", item.ImportID)
	}
	if got := item.Attrs["enable_metrics_persistence"]; got.Kind != hcl.KindBool || got.Bool {
		t.Fatalf("export lost explicit false: %#v", got)
	}
	if got := item.Attrs["max_eps"]; got.Kind != hcl.KindNumber || got.Number != 0 {
		t.Fatalf("export lost explicit zero: %#v", got)
	}
	for _, name := range []string{"metrics_fields_blacklist", "metrics_never_drop_list", "metrics_worker_id_blacklist", "events_metadata_sources", "metrics_drop_list"} {
		if got := item.Attrs[name]; got.Kind != hcl.KindList || len(got.List) != 0 {
			t.Fatalf("export lost empty required list %s: %#v", name, got)
		}
	}
	if got := item.Attrs["samples"].Map["max_size"]; got.Kind != hcl.KindString || got.String != "2 MB" {
		t.Fatalf("export lost nested samples: %#v", got)
	}
	if _, ok := item.Attrs["limits"]; ok {
		t.Fatal("limits must be flat resource attributes")
	}
	for _, field := range []string{"managed_fields", "lookup_allowed_paths", "lookup_max_size", "lookup_max_total_size", "enable_worker_in_memory", "metrics_directory"} {
		if _, ok := item.Attrs[field]; ok {
			t.Fatalf("read-only attribute %s must not become configuration", field)
		}
	}
	file, err := hcl.FileWithResources([]hcl.ResourceInput{{TypeName: item.TypeName, Name: item.Name, Attrs: item.Attrs}}, hcl.DefaultResourceBlockOptions())
	if err != nil {
		t.Fatal(err)
	}
	config := string(file.Bytes())
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"criblio": providerserver.NewProtocol6WithError(&systemLimitsExportTestProvider{}),
		},
		Steps: []resource.TestStep{
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: `resource "criblio_system_limits" "invalid" {
  group_id = "test-fleet"
  min_free_space = "2 GB"
}`, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: `resource "criblio_system_limits" "lookup_only" {
  group_id = "test-fleet"
  metrics_drop_list = ["custom.metric"]
}`, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: strings.ReplaceAll(config, "2 GB", "invalid"), PlanOnly: true, ExpectError: regexp.MustCompile("must match pattern")},
			{Config: strings.Replace(config, "{", "{\n  default_heartbeat_seconds = 0\n", 1), PlanOnly: true, ExpectError: regexp.MustCompile("at least 1")},
			{Config: `resource "criblio_group_system_settings" "invalid" {
  group_id = "test-fleet"
  limits = { min_free_space = "2 GB" }
}`, PlanOnly: true, ExpectError: regexp.MustCompile("Unsupported argument")},
		},
	})
}
