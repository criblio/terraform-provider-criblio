package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEventBreakerJSONOptionsRoundTrip(t *testing.T) {
	for _, pack := range []bool{false, true} {
		name := "group"
		if pack {
			name = "pack"
		}
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			rule := map[string]any{
				"name": "api-events", "condition": "true", "type": "json_array",
				"timestampAnchorRegex": "/^/", "timestamp": map[string]any{"type": "auto"},
				"jsonArrayField": "items.events", "parentFieldsToCopy": []any{"source", "host*"},
				"jsonExtractAll": false,
			}
			fixture := map[string]any{"id": "json-test", "rules": []any{rule}}
			basePath := "/api/v1/m/default/lib/breakers"
			ids := map[string]string{"group_id": "default", "id": "json-test"}
			if pack {
				basePath = "/api/v1/m/default/p/test-pack/lib/breakers"
				ids["pack"] = "test-pack"
			}
			var methods []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if pack && r.URL.Path == "/api/v1/m/default/packs" {
					if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{map[string]any{"id": "test-pack"}}}); err != nil {
						t.Errorf("encode packs: %v", err)
					}
					return
				}
				wantPath := basePath + "/json-test"
				if r.Method == http.MethodPost {
					wantPath = basePath
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
					if !reflect.DeepEqual(body["rules"], fixture["rules"]) {
						t.Errorf("request rules = %#v, want %#v", body["rules"], fixture["rules"])
					}
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": []any{fixture}}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test-token"})
			var res resource.ResourceWithImportState = &EventBreakerRulesetResource{api: newEventBreakerRulesetAPI(client)}
			if pack {
				res = &PackBreakersResource{api: newPackBreakersAPI(client)}
			}
			var schemaResponse resource.SchemaResponse
			res.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
			state := tfsdk.State{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil)}
			plan := tfsdk.Plan{Schema: schemaResponse.Schema, Raw: state.Raw}
			setPlan := func() {
				t.Helper()
				payload, err := json.Marshal(fixture)
				if err != nil {
					t.Fatal(err)
				}
				var model any
				if pack {
					m := &PackBreakersModel{}
					if err := json.Unmarshal(payload, m); err != nil {
						t.Fatal(err)
					}
					m.GroupID, m.Pack = types.StringValue("default"), types.StringValue("test-pack")
					model = m
				} else {
					m := &EventBreakerRulesetModel{}
					if err := json.Unmarshal(payload, m); err != nil {
						t.Fatal(err)
					}
					m.GroupID = types.StringValue("default")
					model = m
				}
				if diags := plan.Set(ctx, model); diags.HasError() {
					t.Fatal(diags)
				}
			}
			assertRules := func(state tfsdk.State) {
				t.Helper()
				var rules types.List
				if diags := state.GetAttribute(ctx, path.Root("rules"), &rules); diags.HasError() {
					t.Fatal(diags)
				}
				got, err := EventBreakerRulesetTerraformValueToJSON(rules)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, fixture["rules"]) {
					t.Fatalf("state rules = %#v, want %#v", got, fixture["rules"])
				}
			}
			setPlan()
			created := resource.CreateResponse{State: state}
			res.Create(ctx, resource.CreateRequest{Plan: plan}, &created)
			if created.Diagnostics.HasError() {
				t.Fatal(created.Diagnostics)
			}
			assertRules(created.State)
			read := resource.ReadResponse{State: created.State}
			res.Read(ctx, resource.ReadRequest{State: created.State}, &read)
			if read.Diagnostics.HasError() {
				t.Fatal(read.Diagnostics)
			}
			assertRules(read.State)
			// Clearing the array path selects a top-level array; clearing parent
			// fields must send [] rather than omitting the requested update.
			rule["jsonArrayField"], rule["parentFieldsToCopy"], rule["jsonExtractAll"] = "", []any{}, true
			setPlan()
			updated := resource.UpdateResponse{State: read.State}
			res.Update(ctx, resource.UpdateRequest{Plan: plan, State: read.State}, &updated)
			if updated.Diagnostics.HasError() {
				t.Fatal(updated.Diagnostics)
			}
			assertRules(updated.State)
			read = resource.ReadResponse{State: updated.State}
			res.Read(ctx, resource.ReadRequest{State: updated.State}, &read)
			if read.Diagnostics.HasError() {
				t.Fatal(read.Diagnostics)
			}
			assertRules(read.State)
			importID, err := json.Marshal(ids)
			if err != nil {
				t.Fatal(err)
			}
			imported := resource.ImportStateResponse{State: state}
			res.ImportState(ctx, resource.ImportStateRequest{ID: string(importID)}, &imported)
			if imported.Diagnostics.HasError() {
				t.Fatal(imported.Diagnostics)
			}
			assertRules(imported.State)
			if want := []string{"POST", "GET", "GET", "PATCH", "GET", "GET", "GET"}; !reflect.DeepEqual(methods, want) {
				t.Fatalf("API methods = %v, want %v", methods, want)
			}
		})
	}
}
