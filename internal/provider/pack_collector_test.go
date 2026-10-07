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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPackCollectorReusesCollectorSchema(t *testing.T) {
	ctx := context.Background()
	var collector, pack resource.SchemaResponse
	(&CollectorResource{}).Schema(ctx, resource.SchemaRequest{}, &collector)
	(&PackCollectorResource{}).Schema(ctx, resource.SchemaRequest{}, &pack)
	if len(pack.Schema.Attributes) != len(collector.Schema.Attributes)+1 {
		t.Fatal("pack collector must add only pack to the collector schema")
	}
	for name, want := range collector.Schema.Attributes {
		got, ok := pack.Schema.Attributes[name]
		if !ok || !got.GetType().Equal(want.GetType()) || got.IsRequired() != want.IsRequired() || got.IsOptional() != want.IsOptional() || got.IsComputed() != want.IsComputed() {
			t.Errorf("schema differs for %s", name)
		}
	}
	if !pack.Schema.Attributes["pack"].IsRequired() {
		t.Fatal("pack must be required")
	}
	// Reusing the actual variant structs keeps payload defaults and sensitive
	// field handling identical to worker-group collectors.
	baseType := reflect.TypeOf(CollectorModel{})
	packType := reflect.TypeOf(PackCollectorModel{})
	for i := 0; i < baseType.NumField(); i++ {
		field := baseType.Field(i)
		if field.Type.Kind() != reflect.Pointer {
			continue
		}
		got, ok := packType.FieldByName(field.Name)
		if !ok || got.Type != field.Type {
			t.Errorf("variant %s does not share the collector model", field.Name)
		}
	}
}

func TestPackCollectorCRUD(t *testing.T) {
	for _, collectorType := range []string{"rest", "cribl_lake"} {
		t.Run(collectorType, func(t *testing.T) {
			var saved map[string]any
			var methods []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/api/v1/m/default/packs" {
					if err := json.NewEncoder(w).Encode(map[string]any{"items": []map[string]string{{"id": "Actual-Pack"}}}); err != nil {
						t.Error(err)
					}
					return
				}
				wantPath := "/api/v1/m/default/p/Actual-Pack/lib/jobs"
				if r.Method != http.MethodPost {
					wantPath += "/events"
				}
				if r.URL.Path != wantPath {
					t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
					http.NotFound(w, r)
					return
				}
				methods = append(methods, r.Method)
				switch r.Method {
				case http.MethodPost, http.MethodPatch:
					if r.Method == http.MethodPost && r.URL.Query().Get("id") != "events" {
						t.Error("create must include collector ID query parameter")
					}
					if err := json.NewDecoder(r.Body).Decode(&saved); err != nil {
						t.Error(err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					if saved["id"] != "events" || saved["type"] != "collection" {
						t.Errorf("incorrect job identity: %#v", saved)
					}
					if saved["pack"] != nil || saved["groupId"] != nil {
						t.Error("path parameters leaked into body")
					}
					if input, ok := saved["input"].(map[string]any); !ok || input["type"] != "collection" {
						t.Errorf("collection input default missing: %#v", saved["input"])
					}
				case http.MethodGet:
					if saved == nil {
						http.NotFound(w, r)
						return
					}
				case http.MethodDelete:
					saved = nil
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					t.Errorf("unexpected method %s", r.Method)
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{saved}, "count": 1}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			api := newPackCollectorAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"}))
			var model PackCollectorModel
			if err := json.Unmarshal([]byte(fmt.Sprintf(`{"collector":{"type":%q},"ttl":"4h"}`, collectorType)), &model); err != nil {
				t.Fatal(err)
			}
			model.GroupID = types.StringValue("default")
			model.Pack = types.StringValue("actual-pack")
			model.ID = types.StringValue("events")
			ctx := context.Background()
			created, err := api.Create(ctx, oneOfRequestModelWithHoistedIdentity(model))
			if err != nil {
				t.Fatal(err)
			}
			applyPackCollectorAPIToState(created, &model, true, false)
			if model.Pack.ValueString() != "actual-pack" || model.GroupID.ValueString() != "default" {
				t.Fatal("API response lost configured scope")
			}
			if _, err := api.Read(ctx, model); err != nil {
				t.Fatal(err)
			}
			res := &PackCollectorResource{api: api}
			var schema resource.SchemaResponse
			res.Schema(ctx, resource.SchemaRequest{}, &schema)
			imported := resource.ImportStateResponse{State: tfsdk.State{
				Schema: schema.Schema,
				Raw:    tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil),
			}}
			res.ImportState(ctx, resource.ImportStateRequest{ID: `{"group_id":"default","id":"events","pack":"actual-pack"}`}, &imported)
			if imported.Diagnostics.HasError() {
				t.Fatalf("import: %v", imported.Diagnostics)
			}
			var importedModel PackCollectorModel
			if diags := imported.State.Get(ctx, &importedModel); diags.HasError() {
				t.Fatalf("decode import state: %v", diags)
			}
			if importedModel.Pack.ValueString() != "actual-pack" || importedModel.ID.ValueString() != "events" {
				t.Fatal("import lost the pack or collector identity")
			}
			if (collectorType == "rest" && importedModel.InputCollectorRest == nil) || (collectorType == "cribl_lake" && importedModel.InputCollectorCriblLake == nil) {
				t.Fatal("import did not populate the collector variant")
			}
			// Set the nested value, as users do in the collector block.
			if model.InputCollectorRest != nil {
				model.InputCollectorRest.Ttl = types.StringValue("8h")
			} else {
				model.InputCollectorCriblLake.Ttl = types.StringValue("8h")
			}
			updated, err := api.Update(ctx, oneOfRequestModelWithHoistedIdentity(model))
			if err != nil {
				t.Fatal(err)
			}
			if updated.Ttl.ValueString() != "8h" {
				t.Fatalf("updated ttl = %q", updated.Ttl.ValueString())
			}
			if err := api.Delete(ctx, model); err != nil {
				t.Fatal(err)
			}
			if _, err := api.Read(ctx, model); !restclient.IsNotFound(err) {
				t.Fatalf("read after delete: %v", err)
			}
			if !reflect.DeepEqual(methods, []string{"POST", "GET", "GET", "PATCH", "DELETE", "GET"}) {
				t.Fatalf("unexpected operations: %v", methods)
			}
		})
	}
}

func TestPackCollectorImportRequiresPack(t *testing.T) {
	var response resource.ImportStateResponse
	(&PackCollectorResource{}).ImportState(context.Background(), resource.ImportStateRequest{ID: `{"group_id":"default","id":"events"}`}, &response)
	if !response.Diagnostics.HasError() {
		t.Fatal("import without pack must fail before making an API request")
	}
}

func TestPackCollectorNormalizesAPIValues(t *testing.T) {
	var model PackCollectorModel
	err := json.Unmarshal([]byte(`{"id":"events","collector":{"type":"google_cloud_storage"},"input":{"throttleRatePerSec":42},"schedule":{"run":{"earliest":1700000000}},"savedState":[]}`), &model)
	if err != nil {
		t.Fatal(err)
	}
	if model.InputCollectorGCS == nil {
		t.Fatal("collector alias was not recognized")
	}
	input := model.InputCollectorGCS.Input.Attributes()
	if input["throttle_rate_per_sec"].(types.String).ValueString() != "42" {
		t.Fatal("numeric API throttle was not normalized")
	}
}
