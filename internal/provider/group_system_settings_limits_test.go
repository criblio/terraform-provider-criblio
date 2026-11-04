package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestGroupSystemSettingsLimitsRoundTrip(t *testing.T) {
	ctx := context.Background()
	configured := map[string]any{"minFreeSpace": "2 GB"}
	defaults := map[string]any{"metricsDirectory": "/var/cribl/metrics", "cpuProfileTTL": "10m"}
	responseLimits := func() map[string]any {
		limits := make(map[string]any, len(configured)+len(defaults))
		for key, value := range defaults {
			limits[key] = value
		}
		for key, value := range configured {
			limits[key] = value
		}
		return limits
	}
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/m/test-fleet/system/settings/conf" {
			t.Errorf("unexpected API path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		methods = append(methods, r.Method)
		if r.Method == http.MethodPatch {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			want := map[string]any{"limits": configured}
			if !reflect.DeepEqual(body, want) {
				t.Errorf("PATCH body = %#v, want %#v", body, want)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{map[string]any{"limits": responseLimits()}}}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
	res := &GroupSystemSettingsResource{api: newGroupSystemSettingsAPI(client)}
	var schemaResponse resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	state := tfsdk.State{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil)}
	plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: state.Raw}
	setPlan := func() {
		t.Helper()
		payload, err := json.Marshal(map[string]any{"limits": configured})
		if err != nil {
			t.Fatal(err)
		}
		var model GroupSystemSettingsModel
		if err := json.Unmarshal(payload, &model); err != nil {
			t.Fatal(err)
		}
		model.GroupID = types.StringValue("test-fleet")
		if diags := plan.Set(ctx, &model); diags.HasError() {
			t.Fatal(diags)
		}
	}
	assertLimits := func(state tfsdk.State, want map[string]any) {
		t.Helper()
		var model GroupSystemSettingsModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatal(diags)
		}
		got, err := GroupSystemSettingsTerraformValueToJSON(model.Limits)
		if err != nil {
			t.Fatal(err)
		}
		payload, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		var normalized map[string]any
		if err := json.Unmarshal(payload, &normalized); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(normalized, want) {
			t.Fatalf("state limits = %#v, want %#v", normalized, want)
		}
	}
	setPlan()
	created := resource.CreateResponse{State: state}
	res.Create(ctx, resource.CreateRequest{Plan: plan}, &created)
	if created.Diagnostics.HasError() {
		t.Fatal(created.Diagnostics)
	}
	assertLimits(created.State, configured)
	read := resource.ReadResponse{State: created.State}
	res.Read(ctx, resource.ReadRequest{State: created.State}, &read)
	if read.Diagnostics.HasError() {
		t.Fatal(read.Diagnostics)
	}
	assertLimits(read.State, configured)
	if !read.State.Raw.Equal(created.State.Raw) {
		t.Fatal("refresh changed partially configured limits")
	}
	configured = map[string]any{
		"minFreeSpace": "3 GB", "maxPQSize": "1 TB", "metricsMaxDiskSpace": "1 GB",
		"maxMetrics": float64(10000), "maxEPs": float64(0), "metricsMaxCardinality": float64(1000),
		"enableMetricsPersistence": false, "metricsDropList": []any{},
		"lookupAllowedPaths": []any{"/mnt/lookups"}, "samples": map[string]any{"maxSize": "2 MB"},
		"kvStoreWALMaxBytes": float64(1024), "metricsGCPeriod": "60s",
	}
	setPlan()
	updated := resource.UpdateResponse{State: read.State}
	res.Update(ctx, resource.UpdateRequest{Plan: plan, State: read.State}, &updated)
	if updated.Diagnostics.HasError() {
		t.Fatal(updated.Diagnostics)
	}
	assertLimits(updated.State, configured)
	read = resource.ReadResponse{State: updated.State}
	res.Read(ctx, resource.ReadRequest{State: updated.State}, &read)
	if read.Diagnostics.HasError() {
		t.Fatal(read.Diagnostics)
	}
	assertLimits(read.State, configured)
	if !read.State.Raw.Equal(updated.State.Raw) {
		t.Fatal("refresh changed updated limits")
	}
	imported := resource.ImportStateResponse{State: state}
	res.ImportState(ctx, resource.ImportStateRequest{ID: "test-fleet"}, &imported)
	if imported.Diagnostics.HasError() {
		t.Fatal(imported.Diagnostics)
	}
	read = resource.ReadResponse{State: imported.State}
	res.Read(ctx, resource.ReadRequest{State: imported.State}, &read)
	if read.Diagnostics.HasError() {
		t.Fatal(read.Diagnostics)
	}
	assertLimits(read.State, responseLimits())
	deleted := resource.DeleteResponse{State: updated.State}
	res.Delete(ctx, resource.DeleteRequest{State: updated.State}, &deleted)
	if deleted.Diagnostics.HasError() {
		t.Fatal(deleted.Diagnostics)
	}

	// The listing data source must also expose limits returned by settings/conf.
	ds := &GroupSystemSettingsDataSource{client: client}
	var dsSchema datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &dsSchema)
	dsState := tfsdk.State{Schema: dsSchema.Schema, Raw: tftypes.NewValue(dsSchema.Schema.Type().TerraformType(ctx), nil)}
	if diags := dsState.Set(ctx, &GroupSystemSettingsListDataSourceModel{
		GroupID: types.StringValue("test-fleet"),
		Items:   types.ListNull(types.ObjectType{AttrTypes: GroupSystemSettingsItemAttrTypes()}),
	}); diags.HasError() {
		t.Fatal(diags)
	}
	config := tfsdk.Config{Schema: dsSchema.Schema, Raw: dsState.Raw}
	dsRead := datasource.ReadResponse{State: dsState}
	ds.Read(ctx, datasource.ReadRequest{Config: config}, &dsRead)
	if dsRead.Diagnostics.HasError() {
		t.Fatal(dsRead.Diagnostics)
	}
	var limits types.Object
	if diags := dsRead.State.GetAttribute(ctx, path.Root("items").AtListIndex(0).AtName("limits"), &limits); diags.HasError() {
		t.Fatal(diags)
	}
	got, err := GroupSystemSettingsTerraformValueToJSON(limits)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var normalized map[string]any
	if err := json.Unmarshal(payload, &normalized); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(normalized, responseLimits()) {
		t.Fatalf("data source limits = %#v, want %#v", got, responseLimits())
	}
	if want := []string{"PATCH", "GET", "PATCH", "GET", "GET", "GET"}; !reflect.DeepEqual(methods, want) {
		t.Fatalf("API methods = %v, want %v", methods, want)
	}
}
