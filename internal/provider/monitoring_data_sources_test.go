package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMonitoringListDataSources(t *testing.T) {
	for _, tc := range []struct {
		name, path, wire string
		newDataSource    func() datasource.DataSource
		config           map[string]any
	}{
		{"shared", "/products/lakehouse_engine_metrics/shared-configs", `{"id":"shared","config":{"priority":{"value":"P3"}},"fingerprint":"hash"}`, NewMonitorSharedConfigsDataSource, nil},
		{"profile", "/products/lakehouse_engine_metrics/config-profiles", `{"id":"profile","name":"Platform","configIds":["shared"]}`, NewMonitorProfilesDataSource, nil},
		{"silence", "/a/__cribl_metrics/products/lakehouse_engine_metrics/silences", `{"id":"silence","monitor_id":"cpu","begins_at":1893456000000,"expires_at":1893542400000,"created_at":123,"tags":{"host":["a","b"]}}`, NewMonitorSilencesDataSource, nil},
		{"policy", "/notification-policies", `{"id":"policy","order":10,"conditions":[[{"key":"team","operator":"=","value":"platform"}]],"templateTargetPairs":[{"targetId":"target","templateId":"template"}]}`, NewNotificationPoliciesDataSource, nil},
		{"current-state", "/products/lakehouse_engine_metrics/engines/engine/datasets/metrics/alert_activity/current", `{"monitor_id":"cpu","fingerprint":"instance","status":"warning","pending_status":"critical"}`, NewMonitorActiveAlertsDataSource, map[string]any{"engine_id": "engine", "dataset_id": "metrics"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/api/v1"+tc.path {
					t.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
					http.NotFound(w, req)
					return
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []json.RawMessage{json.RawMessage(tc.wire)}, "count": 1}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			d := tc.newDataSource()
			var configured datasource.ConfigureResponse
			d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
			if configured.Diagnostics.HasError() {
				t.Fatal(configured.Diagnostics)
			}
			var schema datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schema)
			raw := collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), tc.config)
			response := datasource.ReadResponse{State: tfsdk.State{Schema: schema.Schema}}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: raw}}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			var items types.List
			if diagnostics := response.State.GetAttribute(ctx, path.Root("items"), &items); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if len(items.Elements()) != 1 {
				t.Fatalf("items: %s", items.String())
			}
			if tc.name == "current-state" {
				attrs := items.Elements()[0].(types.Object).Attributes()
				if attrs["status"].(types.String).ValueString() != "warning" || attrs["pending_status"].(types.String).ValueString() != "critical" {
					t.Fatalf("current states were lost: %v", attrs)
				}
			}
		})
	}
}
