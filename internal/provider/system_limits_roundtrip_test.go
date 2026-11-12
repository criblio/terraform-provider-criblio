package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sync"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type systemLimitsTestProvider struct {
	CriblioProvider
	baseURL string
}

func (p *systemLimitsTestProvider) Configure(_ context.Context, _ frameworkprovider.ConfigureRequest, resp *frameworkprovider.ConfigureResponse) {
	clients := &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: p.baseURL, BearerToken: "test"})}
	resp.ResourceData = clients
	resp.DataSourceData = clients
}
func (p *systemLimitsTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSystemLimitsResource}
}

func (p *systemLimitsTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewSystemLimitsDataSource}
}

func TestSystemLimitsManagedLimits(t *testing.T) {
	var mu sync.Mutex
	var stored map[string]json.RawMessage
	if err := json.Unmarshal([]byte(`{
 "cpuProfileTTL":"10m", "enableMetricsPersistence":false,
 "metricsDirectory":"/var/cribl/metrics", "metricsFieldsBlacklist":[],
 "metricsGCPeriod":"60s", "metricsNeverDropList":[],
 "metricsWorkerIdBlacklist":[], "minFreeSpace":"2 GB",
 "samples":{"maxSize":"2 MB","futureSampleOption":9007199254740993},
 "lookupMaxSize":"50 MB", "maxEPs":0, "futureLimit":9007199254740993,
 "metricsDropList":["old.metric"],
 "lookupAllowedPaths":["/readonly"], "lookupMaxTotalSize":"5 GB",
 "enableWorkerInMemory":false
 }`), &stored); err != nil {
		t.Fatal(err)
	}
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path != "/api/v1/m/test-fleet/system/limits" {
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		methods = append(methods, r.Method)
		if r.Method == http.MethodPatch {
			var body map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode PATCH: %v", err)
				w.WriteHeader(400)
				return
			}
			// Emulate upstream replacement, not a partial merge.
			for key, value := range stored {
				if key != "metricsDropList" && key != "minFreeSpace" && key != "samples" && !reflect.DeepEqual(body[key], value) {
					t.Errorf("replacement PATCH changed unconfigured %s", key)
				}
			}
			for _, key := range []string{"cpuProfileTTL", "enableMetricsPersistence", "metricsDirectory", "metricsFieldsBlacklist", "metricsGCPeriod", "metricsNeverDropList", "metricsWorkerIdBlacklist", "minFreeSpace", "samples", "lookupMaxSize", "maxEPs", "futureLimit"} {
				if _, ok := body[key]; !ok {
					t.Errorf("replacement PATCH dropped %s", key)
				}
			}
			if string(body["futureLimit"]) != "9007199254740993" {
				t.Errorf("unknown integer corrupted: %s", body["futureLimit"])
			}
			if string(body["enableMetricsPersistence"]) != "false" || string(body["maxEPs"]) != "0" || string(body["metricsFieldsBlacklist"]) != "[]" {
				t.Error("PATCH lost false, zero or empty list")
			}
			var samples map[string]json.RawMessage
			if err := json.Unmarshal(body["samples"], &samples); err != nil {
				t.Errorf("decode samples: %v", err)
			}
			if string(samples["futureSampleOption"]) != "9007199254740993" {
				t.Error("PATCH dropped unknown nested sample field")
			}
			stored = body
		} else if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s", r.Method)
			w.WriteHeader(405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{stored}}); err != nil {
			t.Errorf("encode limits: %v", err)
		}
	}))
	defer server.Close()
	config := func(metrics, extra string) string {
		return fmt.Sprintf(`resource "criblio_system_limits" "test" {
 group_id = "test-fleet"
 metrics_drop_list = %s
 %s
}`, metrics, extra)
	}
	name := "criblio_system_limits.test"
	var sch frameworkresource.SchemaResponse
	NewSystemLimitsResource().Schema(context.Background(), frameworkresource.SchemaRequest{}, &sch)
	var ignored []string
	for key := range sch.Schema.Attributes {
		if key != "group_id" && key != "metrics_drop_list" {
			ignored = append(ignored, key)
		}
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"criblio": providerserver.NewProtocol6WithError(&systemLimitsTestProvider{baseURL: server.URL})},
		Steps: []resource.TestStep{
			{Config: config(`["cribl.cpu"]`, ""), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "metrics_drop_list.0", "cribl.cpu"),
				resource.TestCheckNoResourceAttr(name, "min_free_space"),
				resource.TestCheckResourceAttr(name, "lookup_allowed_paths.0", "/readonly"),
				resource.TestCheckResourceAttr(name, "lookup_max_size", "50 MB"),
				resource.TestCheckResourceAttr(name, "lookup_max_total_size", "5 GB"),
				resource.TestCheckResourceAttr(name, "enable_worker_in_memory", "false"),
				resource.TestCheckResourceAttr(name, "metrics_directory", "/var/cribl/metrics"),
				resource.TestCheckResourceAttr(name, "managed_fields.#", "1"),
				resource.TestCheckResourceAttr(name, "managed_fields.0", "metricsDropList"),
			)},
			{Config: config(`["cribl.cpu"]`, ""), PlanOnly: true},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				stored["lookupAllowedPaths"] = json.RawMessage(`["/updated"]`)
				stored["lookupMaxSize"] = json.RawMessage(`"100 MB"`)
				stored["lookupMaxTotalSize"] = json.RawMessage(`"10 GB"`)
				stored["enableWorkerInMemory"] = json.RawMessage(`true`)
				stored["metricsDirectory"] = json.RawMessage(`"/var/updated/metrics"`)
			}, Config: config(`["cribl.cpu"]`, ""), PlanOnly: true},
			{Config: config(`["cribl.cpu"]`, ""), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "lookup_allowed_paths.0", "/updated"),
				resource.TestCheckResourceAttr(name, "lookup_max_size", "100 MB"),
				resource.TestCheckResourceAttr(name, "lookup_max_total_size", "10 GB"),
				resource.TestCheckResourceAttr(name, "enable_worker_in_memory", "true"),
				resource.TestCheckResourceAttr(name, "metrics_directory", "/var/updated/metrics"),
			)},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				stored["minFreeSpace"] = json.RawMessage(`"500 MB"`)
			}, Config: config(`["cribl.cpu"]`, ""), PlanOnly: true},
			{Config: config(`[]`, ""), Check: resource.TestCheckResourceAttr(name, "metrics_drop_list.#", "0")},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				stored["metricsDropList"] = json.RawMessage(`["remote.metric"]`)
				stored["minFreeSpace"] = json.RawMessage(`"500 MB"`)
			}, Config: config(`[]`, ""), PlanOnly: true, ExpectNonEmptyPlan: true},
			{PreConfig: func() { mu.Lock(); defer mu.Unlock(); delete(stored, "metricsDropList") }, Config: config(`[]`, ""), PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config(`[]`, ""), PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config(`["cribl.cpu","cribl.mem"]`, "min_free_space = \"3 GB\"\nsamples = { max_size = \"1 MB\" }"), Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr(name, "metrics_drop_list.#", "2"), resource.TestCheckResourceAttr(name, "samples.max_size", "1 MB"))},
			{Config: config(`["cribl.cpu"]`, ""), Check: resource.TestCheckNoResourceAttr(name, "min_free_space")},
			{Config: config(`["cribl.cpu"]`, ""), PlanOnly: true},
			{Config: config(`["cribl.cpu"]`, "") + `
data "criblio_system_limits" "test" {
  group_id = criblio_system_limits.test.group_id
}`, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.criblio_system_limits.test", "min_free_space", "3 GB"),
				resource.TestCheckResourceAttr("data.criblio_system_limits.test", "metrics_drop_list.0", "cribl.cpu"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "test-fleet", ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "group_id", ImportStateVerifyIgnore: ignored},
			{Config: `resource "criblio_system_limits" "test" { group_id = "test-fleet" }`, Check: resource.TestCheckNoResourceAttr(name, "metrics_drop_list.#")},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				stored["metricsDropList"] = json.RawMessage(`["external.metric"]`)
			}, Config: `resource "criblio_system_limits" "test" { group_id = "test-fleet" }`, PlanOnly: true},
			{Config: config(`["cribl.cpu"]`, ""), Check: resource.TestCheckResourceAttr(name, "metrics_drop_list.0", "cribl.cpu")},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				stored["lookupAllowedPaths"] = json.RawMessage(`[]`)
			}, Config: config(`["cribl.cpu"]`, ""), Check: resource.TestCheckResourceAttr(name, "lookup_allowed_paths.#", "0")},
			{PreConfig: func() {
				mu.Lock()
				defer mu.Unlock()
				delete(stored, "lookupAllowedPaths")
			}, Config: config(`["cribl.cpu"]`, ""), Check: resource.TestCheckNoResourceAttr(name, "lookup_allowed_paths.#")},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	if string(stored["metricsDropList"]) != `["cribl.cpu"]` {
		t.Fatalf("final metrics: %s", stored["metricsDropList"])
	}
	if string(stored["minFreeSpace"]) != `"3 GB"` {
		t.Fatalf("removing configuration erased server value: %s", stored["minFreeSpace"])
	}
	if !reflect.DeepEqual(methods[:2], []string{"GET", "PATCH"}) {
		t.Fatalf("create methods %v", methods)
	}
}

func TestSystemLimitsOwnershipAfterRemoteRemoval(t *testing.T) {
	for _, tc := range []struct {
		name, field, config, initial, external, prefix string
	}{
		{"metrics list", "metricsDropList", `metrics_drop_list = ["cribl.cpu"]`, `["cribl.cpu"]`, `["external.metric"]`, ""},
		{"zero limit", "maxEPs", `max_eps = 0`, `0`, `10`, ""},
		{"nested object", "samples", `samples = { max_size = "1 MB" }`, `{"maxSize":"1 MB"}`, `{"maxSize":"2 MB"}`, ""},
		{"unknown metric", "metricsDropList", `metrics_drop_list = [terraform_data.metric.output]`, `["cribl.cpu"]`, `["external.metric"]`, `resource "terraform_data" "metric" { input = "cribl.cpu" }`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			stored := map[string]json.RawMessage{
				"minFreeSpace":         json.RawMessage(`"2 GB"`),
				"lookupAllowedPaths":   json.RawMessage(`["/readonly"]`),
				"lookupMaxSize":        json.RawMessage(`"500 MB"`),
				"lookupMaxTotalSize":   json.RawMessage(`"2 GB"`),
				"enableWorkerInMemory": json.RawMessage(`false`),
				"metricsDirectory":     json.RawMessage(`"/metrics"`),
			}
			patches := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				if r.URL.Path != "/api/v1/m/test-fleet/system/limits" {
					t.Errorf("unexpected path %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}
				switch r.Method {
				case http.MethodPatch:
					patches++
					var body map[string]json.RawMessage
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode PATCH: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if _, ok := body["managedFields"]; ok {
						t.Error("PATCH included Terraform ownership metadata")
					}
					if string(body[tc.field]) != tc.initial {
						t.Errorf("configured field = %s, want %s", body[tc.field], tc.initial)
					}
					stored = body
				case http.MethodGet:
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{stored}}); err != nil {
					t.Errorf("encode limits: %v", err)
				}
			}))
			defer server.Close()
			groupOnly := tc.prefix + "\n" + `resource "criblio_system_limits" "test" { group_id = "test-fleet" }`
			configured := tc.prefix + "\n" + fmt.Sprintf(`resource "criblio_system_limits" "test" {
                group_id = "test-fleet"
                %s
            }`, tc.config)
			steps := []resource.TestStep{{Config: configured}}
			if tc.name == "metrics list" {
				steps = append(steps, resource.TestStep{Config: `resource "criblio_system_limits" "test" {
                    group_id = "test-fleet"
                    lifecycle { ignore_changes = [metrics_drop_list] }
                }`, PlanOnly: true})
			}
			steps = append(steps, []resource.TestStep{
				// Plan-only refresh must retain ownership while the field is configured.
				{PreConfig: func() { mu.Lock(); defer mu.Unlock(); delete(stored, tc.field) }, Config: configured, PlanOnly: true, ExpectNonEmptyPlan: true},
				{Config: configured, PlanOnly: true, ExpectNonEmptyPlan: true},
				// Applying removal must persist ownership even when both values are null.
				{Config: groupOnly, Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_system_limits.test", "managed_fields.#", "0"),
					resource.TestCheckResourceAttr("criblio_system_limits.test", "lookup_allowed_paths.0", "/readonly"),
				)},
				{PreConfig: func() { mu.Lock(); defer mu.Unlock(); stored[tc.field] = json.RawMessage(tc.external) }, Config: groupOnly, PlanOnly: true},
				{Config: groupOnly, PlanOnly: true},
				{Config: configured, PlanOnly: true, ExpectNonEmptyPlan: true},
			}...)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"criblio": providerserver.NewProtocol6WithError(&systemLimitsTestProvider{baseURL: server.URL})},
				Steps:                    steps,
			})
			mu.Lock()
			defer mu.Unlock()
			if patches != 1 {
				t.Fatalf("ownership-only changes caused API writes: got %d PATCHes, want 1", patches)
			}
			if string(stored[tc.field]) != tc.external {
				t.Fatalf("external value changed: %s", stored[tc.field])
			}
		})
	}
}

func TestSystemLimitsRejectReadOnlyConfiguration(t *testing.T) {
	for _, tc := range []struct{ field, value string }{
		{"lookup_allowed_paths", `["/tmp"]`},
		{"lookup_max_size", `"500 MB"`},
		{"lookup_max_total_size", `"2 GB"`},
		{"enable_worker_in_memory", `false`},
		{"metrics_directory", `"/metrics"`},
	} {
		t.Run(tc.field, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){"criblio": providerserver.NewProtocol6WithError(&systemLimitsTestProvider{baseURL: "https://example.invalid"})},
				Steps: []resource.TestStep{{
					Config: fmt.Sprintf(`resource "criblio_system_limits" "test" {
  group_id = "default"
  %s = %s
}`, tc.field, tc.value),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile("Invalid Configuration for Read-Only Attribute"),
				}},
			})
		})
	}
}

func TestSystemLimitsFailedReadDoesNotPatch(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"forbidden", http.StatusForbidden, `{"message":"forbidden"}`},
		{"invalid JSON", http.StatusOK, `{"items":[`},
		{"empty item", http.StatusOK, `{"items":[{}]}`},
		{"no items", http.StatusOK, `{"items":[]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			patches := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					patches++
					t.Errorf("unexpected mutation %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				if _, err := w.Write([]byte(tc.body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			rc := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
			res := &systemLimitsManagedResource{SystemLimitsResource: SystemLimitsResource{client: rc}}
			ctx := context.Background()
			var sch frameworkresource.SchemaResponse
			res.Schema(ctx, frameworkresource.SchemaRequest{}, &sch)
			state := tfsdk.State{Schema: sch.Schema, Raw: tftypes.NewValue(sch.Schema.Type().TerraformType(ctx), nil)}
			plan := tfsdk.Plan{Schema: sch.Schema, Raw: state.Raw}
			var model SystemLimitsModel
			if err := json.Unmarshal([]byte(`{"metricsDropList":["cribl.cpu"]}`), &model); err != nil {
				t.Fatal(err)
			}
			model.GroupID = types.StringValue("test-fleet")
			if diags := plan.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			resp := frameworkresource.CreateResponse{State: state}
			res.Create(ctx, frameworkresource.CreateRequest{Plan: plan}, &resp)
			if !resp.Diagnostics.HasError() {
				t.Fatal("invalid GET did not fail create")
			}
			if patches != 0 {
				t.Fatal("failed GET caused a replacement PATCH")
			}
		})
	}
}
