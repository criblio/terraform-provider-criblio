package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const monitorReviewJSON = `{
  "id":"review-monitor", "name":"Review monitor", "datasetId":"metrics", "enabled":true, "type":"threshold",
  "managedBy":"terraform", "expr":[{"label":"Alert Query", "text":"A", "queryLabels":["A"]}],
  "query":{"A":{"mode":"promql", "promql":"old_metric"}}, "silence":[],
  "firingCondition":{"fire_delay":300,"clear_delay":60}, "firingRule":{"label":"A","threshold":[]},
  "metadata":{"owner":"example"}, "notification":{"type":"policy","enabled":false,"config":[]},
  "priority":{"value":"P3"}, "team":{"value":""}
}`

func monitorReviewObject(t *testing.T) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal([]byte(monitorReviewJSON), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func monitorReviewModel(t *testing.T, object map[string]any) MonitorModel {
	t.Helper()
	encoded, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	var model MonitorModel
	if err := json.Unmarshal(encoded, &model); err != nil {
		t.Fatal(err)
	}
	return model
}

func monitorReviewCheckWire(t *testing.T, model MonitorModel, object map[string]any) {
	t.Helper()
	encoded, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	var actual map[string]any
	if err := json.Unmarshal(encoded, &actual); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"query", "expr", "metadata", "notification"} {
		if !reflect.DeepEqual(actual[field], object[field]) {
			t.Errorf("%s: got %#v, want %#v", field, actual[field], object[field])
		}
	}
}

func TestMonitorInheritanceAndMetadataRoundTrip(t *testing.T) {
	for _, inherited := range []bool{false, true} {
		name := "inline"
		if inherited {
			name = "inherited"
		}
		t.Run(name, func(t *testing.T) {
			object := monitorReviewObject(t)
			if inherited {
				refs := map[string]any{"$configRefs": []any{map[string]any{"id": "shared", "fingerprint": "hash"}}}
				object["notification"] = refs
				object["metadata"] = refs
				object["query"] = map[string]any{"A": refs}
			}
			model := monitorReviewModel(t, object)
			monitorReviewCheckWire(t, model, object)
			var schema resource.SchemaResponse
			NewMonitorResource().Schema(context.Background(), resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			if diagnostics := state.Set(context.Background(), &model); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			var restored MonitorModel
			if diagnostics := state.Get(context.Background(), &restored); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			monitorReviewCheckWire(t, restored, object)
		})
	}
}

func TestMonitorReadRefreshesConfiguration(t *testing.T) {
	for _, name := range []string{"drift", "inheritance-and-cleared-expression", "import"} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			object := monitorReviewObject(t)
			object["query"].(map[string]any)["A"].(map[string]any)["promql"] = "updated_metric"
			object["expr"].([]any)[0].(map[string]any)["text"] = "A * 2"
			object["metadata"].(map[string]any)["owner"] = "updated"
			if name != "drift" {
				refs := map[string]any{"$configRefs": []any{map[string]any{"id": "shared", "fingerprint": "updated"}}}
				object["query"] = map[string]any{"A": refs}
				object["notification"] = refs
				object["expr"] = []any{}
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || !strings.HasSuffix(req.URL.Path, "/monitors/review-monitor") {
					t.Errorf("unexpected monitor request: %s %s", req.Method, req.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{object}, "count": 1}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
			r := NewMonitorResource().(*MonitorResource)
			r.api = newMonitorAPI(client)
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			prior := monitorReviewModel(t, monitorReviewObject(t))
			if name == "import" {
				state.Raw = collectorPlanTestObject(schema.Schema.Type().TerraformType(ctx).(tftypes.Object), map[string]any{"id": "review-monitor"})
			} else if diagnostics := state.Set(ctx, &prior); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			response := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &response)
			if response.Diagnostics.HasError() {
				t.Fatal(response.Diagnostics)
			}
			var refreshed MonitorModel
			if diagnostics := response.State.Get(ctx, &refreshed); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			monitorReviewCheckWire(t, refreshed, object)
		})
	}
}

func TestMonitorQueryUnionValidation(t *testing.T) {
	for _, tc := range []struct {
		name    string
		query   map[string]any
		invalid bool
	}{
		{"inline", map[string]any{"mode": "promql", "promql": "metric"}, false},
		{"inherited", map[string]any{"$configRefs": []any{}}, false},
		{"missing", map[string]any{}, true},
		{"partial-inline", map[string]any{"mode": "promql"}, true},
		{"both", map[string]any{"mode": "promql", "promql": "metric", "$configRefs": []any{}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			object := monitorReviewObject(t)
			object["query"] = map[string]any{"A": tc.query}
			model := monitorReviewModel(t, object)
			r := NewMonitorResource().(*MonitorResource)
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			if diagnostics := state.Set(ctx, &model); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			var response resource.ValidateConfigResponse
			r.ValidateConfig(ctx, resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: state.Raw}}, &response)
			if response.Diagnostics.HasError() != tc.invalid {
				t.Fatalf("query validation: got %v, expected invalid=%v", response.Diagnostics, tc.invalid)
			}
		})
	}
}
