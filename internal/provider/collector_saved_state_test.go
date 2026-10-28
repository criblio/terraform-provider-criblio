package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestCollectorSavedStateCompatibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state string
		want  any
	}{
		{name: "database checkpoint", state: `{"databaseCollectorState":{"data":{"id":{"value":3}}}}`},
		{name: "object checkpoint", state: `{"cursor":{"bucket":"events","objects":["one","two"]}}`},
		{name: "mixed checkpoint", state: `{"legacy":"value","checkpoint":{"value":3}}`},
		{name: "numeric checkpoint", state: `{"offset":3}`},
		{name: "boolean checkpoint", state: `{"completed":true}`},
		{name: "array checkpoint", state: `{"files":["one","two"]}`},
		{name: "legacy strings", state: `{"cursor":"abc","offset":"3"}`, want: map[string]any{"cursor": "abc", "offset": "3"}},
		{name: "nullable strings", state: `{"cursor":null}`, want: map[string]any{}},
		{name: "empty object", state: `{}`, want: map[string]any{}},
		{name: "empty array", state: `[]`, want: map[string]any{}},
		{name: "null", state: `null`},
	}
	for _, variant := range []string{"splunk", "rest", "s3", "azure_blob", "cribl_lake", "database", "google_cloud_storage", "health_check", "script", "filesystem"} {
		for _, tt := range tests {
			t.Run(variant+"/"+tt.name, func(t *testing.T) {
				payload := []byte(fmt.Sprintf(`{"id":"events","collector":{"type":%q},"savedState":%s}`, variant, tt.state))
				var model CollectorModel
				if err := json.Unmarshal(payload, &model); err != nil {
					t.Fatalf("decode collector checkpoint: %v", err)
				}
				encoded, err := json.Marshal(model)
				if err != nil {
					t.Fatalf("encode collector update: %v", err)
				}
				var update map[string]any
				if err := json.Unmarshal(encoded, &update); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(update["savedState"], tt.want) {
					t.Errorf("savedState in update = %#v, want %#v", update["savedState"], tt.want)
				}
				if _, ok := update["savedState"]; tt.want == nil && ok {
					t.Error("runtime checkpoint must be omitted from updates")
				}
				if update["id"] != "events" {
					t.Errorf("collector identity = %v, want events", update["id"])
				}
			})
		}
	}
}

func TestCollectorImportAndRefreshWithRuntimeCheckpoint(t *testing.T) {
	for _, scope := range []struct {
		name        string
		importID    string
		newResource func(*restclient.Client) resource.ResourceWithImportState
	}{
		{
			name:     "group",
			importID: `{"group_id":"default","id":"events"}`,
			newResource: func(client *restclient.Client) resource.ResourceWithImportState {
				return &CollectorResource{api: newCollectorAPI(client)}
			},
		},
		{
			name:     "pack",
			importID: `{"group_id":"default","id":"events","pack":"example"}`,
			newResource: func(client *restclient.Client) resource.ResourceWithImportState {
				return &PackCollectorResource{api: newPackCollectorAPI(client)}
			},
		},
	} {
		for _, variant := range []struct {
			apiType    string
			block      string
			checkpoint string
		}{
			{"database", "input_collector_database", `{"databaseCollectorState":{"data":{"id":{"value":3}}}}`},
			{"google_cloud_storage", "input_collector_gcs", `{"cursor":{"bucket":"events","objects":["one","two"]}}`},
		} {
			t.Run(scope.name+"/"+variant.apiType, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if r.Method != http.MethodGet {
						t.Errorf("unexpected method: %s", r.Method)
					}
					if r.URL.Path == "/api/v1/m/default/packs" {
						if _, err := fmt.Fprint(w, `{"items":[{"id":"example"}]}`); err != nil {
							t.Error(err)
						}
						return
					}
					wantPath := "/api/v1/m/default/lib/jobs/events"
					if scope.name == "pack" {
						wantPath = "/api/v1/m/default/p/example/lib/jobs/events"
					}
					if r.URL.Path != wantPath {
						t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
						http.NotFound(w, r)
						return
					}
					if _, err := fmt.Fprintf(w, `{"items":[{"id":"events","collector":{"type":%q},"savedState":%s}],"count":1}`, variant.apiType, variant.checkpoint); err != nil {
						t.Error(err)
					}
				}))
				defer server.Close()
				ctx := context.Background()
				res := scope.newResource(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"}))
				var schema resource.SchemaResponse
				res.Schema(ctx, resource.SchemaRequest{}, &schema)
				imported := resource.ImportStateResponse{State: tfsdk.State{
					Schema: schema.Schema,
					Raw:    tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil),
				}}
				res.ImportState(ctx, resource.ImportStateRequest{ID: scope.importID}, &imported)
				if imported.Diagnostics.HasError() {
					t.Fatalf("import collector with checkpoint: %v", imported.Diagnostics)
				}
				savedStatePath := path.Root(variant.block).AtName("saved_state")
				var savedState types.Map
				if diags := imported.State.GetAttribute(ctx, savedStatePath, &savedState); diags.HasError() {
					t.Fatal(diags)
				}
				if !savedState.IsNull() {
					t.Errorf("imported runtime checkpoint = %v, want null", savedState)
				}
				// Existing configured string maps must survive refresh even when
				// the API enriches the collector with an incompatible checkpoint.
				legacy := types.MapValueMust(types.StringType, map[string]attr.Value{"cursor": types.StringValue("configured")})
				if diags := imported.State.SetAttribute(ctx, savedStatePath, legacy); diags.HasError() {
					t.Fatal(diags)
				}
				refreshed := resource.ReadResponse{State: imported.State}
				res.Read(ctx, resource.ReadRequest{State: imported.State}, &refreshed)
				if refreshed.Diagnostics.HasError() {
					t.Fatalf("refresh collector with checkpoint: %v", refreshed.Diagnostics)
				}
				if diags := refreshed.State.GetAttribute(ctx, savedStatePath, &savedState); diags.HasError() {
					t.Fatal(diags)
				}
				if !savedState.Equal(legacy) {
					t.Errorf("legacy saved_state after refresh = %v, want %v", savedState, legacy)
				}
			})
		}
	}
}
