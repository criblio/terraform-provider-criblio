package provider

import (
	"context"
	"encoding/json"
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

func TestNotificationTriggerRoundTrip(t *testing.T) {
	for _, triggerType := range []string{"custom", "resultsCount"} {
		t.Run(triggerType, func(t *testing.T) {
			ctx := context.Background()
			conf := map[string]any{
				"savedQueryId": "daily-volume", "triggerType": triggerType,
				"triggerComparator": ">", "triggerCount": float64(0),
			}
			if triggerType == "custom" {
				conf["trigger"] = "baselineComplete && (actual == 0 || actual < lower || actual > upper)"
			}
			fixture := map[string]any{
				"id": "daily-volume-alert", "condition": "search", "group": "default_search", "conf": conf,
			}
			var methods []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/api/v1/m/default_search/notifications/daily-volume-alert"
				if r.Method == http.MethodPost {
					wantPath = "/api/v1/m/default_search/notifications"
				}
				if r.URL.Path != wantPath {
					t.Errorf("API path = %s, want %s", r.URL.Path, wantPath)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				methods = append(methods, r.Method)
				if r.Method == http.MethodPost || r.Method == http.MethodPatch {
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode request: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if !reflect.DeepEqual(body["conf"], conf) {
						t.Errorf("request conf = %v, want %v", body["conf"], conf)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{fixture}}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()
			res := &NotificationResource{api: newNotificationAPI(restclient.New(restclient.Config{
				BaseURL: server.URL, BearerToken: "test-token",
			}))}
			var schemaResponse resource.SchemaResponse
			res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
			state := tfsdk.State{
				Schema: schemaResponse.Schema,
				Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
			}
			payload, err := json.Marshal(fixture)
			if err != nil {
				t.Fatal(err)
			}
			var model NotificationModel
			if err := json.Unmarshal(payload, &model); err != nil {
				t.Fatal(err)
			}
			plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: state.Raw}
			if diags := plan.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			created := resource.CreateResponse{State: state}
			res.Create(ctx, resource.CreateRequest{Plan: plan}, &created)
			if created.Diagnostics.HasError() {
				t.Fatal(created.Diagnostics)
			}
			assertNotificationConf(t, created.State, conf)
			read := resource.ReadResponse{State: created.State}
			res.Read(ctx, resource.ReadRequest{State: created.State}, &read)
			if read.Diagnostics.HasError() {
				t.Fatal(read.Diagnostics)
			}
			assertNotificationConf(t, read.State, conf)
			if !read.State.Raw.Equal(created.State.Raw) {
				t.Fatal("refresh changed notification state")
			}
			// Change the expression (or legacy threshold) and verify the API and
			// subsequent refresh retain the updated configuration.
			if triggerType == "custom" {
				conf["trigger"] = "baselineComplete && (actual == 0 || actual < 8000 || actual > 12000)"
			} else {
				conf["triggerCount"] = float64(10)
			}
			value, err := NotificationAPIValueToTerraformValue(conf, types.ObjectType{AttrTypes: NotificationConfAttrTypes()})
			if err != nil {
				t.Fatal(err)
			}
			model.Conf = value.(types.Object)
			if diags := plan.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			updated := resource.UpdateResponse{State: read.State}
			res.Update(ctx, resource.UpdateRequest{Plan: plan, State: read.State}, &updated)
			if updated.Diagnostics.HasError() {
				t.Fatal(updated.Diagnostics)
			}
			assertNotificationConf(t, updated.State, conf)
			read = resource.ReadResponse{State: updated.State}
			res.Read(ctx, resource.ReadRequest{State: updated.State}, &read)
			if read.Diagnostics.HasError() {
				t.Fatal(read.Diagnostics)
			}
			assertNotificationConf(t, read.State, conf)
			if !read.State.Raw.Equal(updated.State.Raw) {
				t.Fatal("refresh changed updated notification state")
			}
			imported := resource.ImportStateResponse{State: state}
			res.ImportState(ctx, resource.ImportStateRequest{ID: `{"group":"default_search","id":"daily-volume-alert"}`}, &imported)
			if imported.Diagnostics.HasError() {
				t.Fatal(imported.Diagnostics)
			}
			assertNotificationConf(t, imported.State, conf)
			if want := []string{"POST", "GET", "PATCH", "GET", "GET"}; !reflect.DeepEqual(methods, want) {
				t.Fatalf("API methods = %v, want %v", methods, want)
			}
		})
	}
}

func assertNotificationConf(t *testing.T, state tfsdk.State, want map[string]any) {
	t.Helper()
	var model NotificationModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatal(diags)
	}
	got, err := NotificationTerraformValueToJSON(model.Conf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("state conf = %v, want %v", got, want)
	}
}
