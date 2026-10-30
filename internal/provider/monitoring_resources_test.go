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
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMonitoringResourceLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name, path, body, driftAPI, driftTF string
		newResource                         func() resource.Resource
		newModel                            func() any
	}{
		{"shared", "/products/lakehouse_engine_metrics/shared-configs", `{"id":"test","name":"Initial","config":{"priority":{"value":"P3"},"metadata":{"owner":"platform"}}}`, "name", "name", NewMonitorSharedConfigResource, func() any { return &MonitorSharedConfigModel{} }},
		{"profile", "/products/lakehouse_engine_metrics/config-profiles", `{"id":"test","name":"Initial","configIds":[]}`, "name", "name", NewMonitorProfileResource, func() any { return &MonitorProfileModel{} }},
		{"silence", "/a/__cribl_metrics/products/lakehouse_engine_metrics/silences", `{"id":"test","name":"Initial","monitor_id":"cpu","begins_at":1893456000000,"expires_at":1893542400000,"tags":{"host":["a","b"]}}`, "name", "name", NewMonitorSilenceResource, func() any { return &MonitorSilenceModel{} }},
		{"policy", "/notification-policies", `{"id":"test","order":10,"conditions":[[{"key":"team","operator":"=","value":"platform"},{"key":"count","operator":"=","value":3},{"key":"enabled","operator":"=","value":false},{"key":"unset","operator":"=","value":null}]],"templateTargetPairs":[{"templateId":"template","targetId":"target"}]}`, "order", "order", NewNotificationPolicyResource, func() any { return &NotificationPolicyModel{} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var stored map[string]any
			var lastPayload map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				requestPath := strings.TrimPrefix(req.URL.Path, "/api/v1")
				if requestPath != tc.path && requestPath != tc.path+"/test" {
					t.Errorf("unexpected path %s", req.URL.Path)
					http.NotFound(w, req)
					return
				}
				switch req.Method {
				case http.MethodPost, http.MethodPatch:
					var body map[string]any
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					lastPayload = body
					if req.Method == http.MethodPost {
						for _, field := range []string{"fingerprint", "created_at", "created_by"} {
							if _, ok := body[field]; ok {
								t.Errorf("create set server-owned %s", field)
							}
						}
					}
					stored = make(map[string]any, len(body))
					for key, value := range body {
						stored[key] = value
					}
					if tc.name == "shared" {
						stored["fingerprint"] = "server-hash"
					}
					if tc.name == "policy" {
						if _, ok := stored["disabled"]; !ok {
							stored["disabled"] = false
						}
						if _, ok := stored["final"]; !ok {
							stored["final"] = false
						}
						if _, ok := stored["waitToGroup"]; !ok {
							stored["waitToGroup"] = float64(0)
						}
					}
					if tc.name == "silence" {
						stored["created_at"] = float64(123456)
						stored["created_by"] = "api-user"
					}
				case http.MethodDelete:
					stored = nil
				case http.MethodGet:
					if stored == nil {
						http.NotFound(w, req)
						return
					}
				default:
					t.Errorf("unexpected method %s", req.Method)
				}
				items := []any{}
				if stored != nil {
					items = append(items, stored)
				}
				if err := json.NewEncoder(w).Encode(map[string]any{"items": items, "count": len(items)}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			r := tc.newResource()
			var configured resource.ConfigureResponse
			r.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{ProviderData: &ProviderClients{RC: restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})}}, &configured)
			if configured.Diagnostics.HasError() {
				t.Fatal(configured.Diagnostics)
			}
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			model := tc.newModel()
			if err := json.Unmarshal([]byte(tc.body), model); err != nil {
				t.Fatal(err)
			}
			plan := tfsdk.State{Schema: schema.Schema}
			if diagnostics := plan.Set(ctx, model); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			create := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema}}
			r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: schema.Schema, Raw: plan.Raw}}, &create)
			if create.Diagnostics.HasError() {
				t.Fatal(create.Diagnostics)
			}
			var expected map[string]any
			if err := json.Unmarshal([]byte(tc.body), &expected); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(lastPayload, expected) {
				t.Errorf("payload: got %#v, want %#v", lastPayload, expected)
			}
			if tc.name != "profile" {
				updateModel := tc.newModel()
				if diagnostics := create.State.Get(ctx, updateModel); diagnostics.HasError() {
					t.Fatal(diagnostics)
				}
				updatePlan := tfsdk.State{Schema: schema.Schema}
				if diagnostics := updatePlan.Set(ctx, updateModel); diagnostics.HasError() {
					t.Fatal(diagnostics)
				}
				update := resource.UpdateResponse{State: create.State}
				r.Update(ctx, resource.UpdateRequest{State: create.State, Plan: tfsdk.Plan{Schema: schema.Schema, Raw: updatePlan.Raw}}, &update)
				if update.Diagnostics.HasError() {
					t.Fatal(update.Diagnostics)
				}
				if tc.name == "silence" {
					if lastPayload["created_at"] != float64(123456) || lastPayload["created_by"] != "api-user" {
						t.Errorf("update lost audit values: %#v", lastPayload)
					}
				}
				create.State = update.State
			}
			if tc.name == "policy" {
				stored[tc.driftAPI] = float64(20)
			} else {
				stored[tc.driftAPI] = "Changed outside Terraform"
			}
			read := resource.ReadResponse{State: create.State}
			r.Read(ctx, resource.ReadRequest{State: create.State}, &read)
			if read.Diagnostics.HasError() {
				t.Fatal(read.Diagnostics)
			}
			refreshed := tc.newModel()
			if diagnostics := read.State.Get(ctx, refreshed); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			wire, err := json.Marshal(refreshed)
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]any
			if err := json.Unmarshal(wire, &actual); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual[tc.driftAPI], stored[tc.driftAPI]) {
				t.Fatalf("refresh did not detect drift: %s", wire)
			}
			imported := resource.ImportStateResponse{State: tfsdk.State{Schema: schema.Schema, Raw: tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil)}}
			r.(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: "test"}, &imported)
			if imported.Diagnostics.HasError() {
				t.Fatal(imported.Diagnostics)
			}
			importRead := resource.ReadResponse{State: imported.State}
			r.Read(ctx, resource.ReadRequest{State: imported.State}, &importRead)
			if importRead.Diagnostics.HasError() {
				t.Fatal(importRead.Diagnostics)
			}
			if !importRead.State.Raw.Equal(read.State.Raw) {
				t.Errorf("import did not hydrate the same API state")
			}
			deleted := resource.DeleteResponse{State: read.State}
			r.Delete(ctx, resource.DeleteRequest{State: read.State}, &deleted)
			if deleted.Diagnostics.HasError() {
				t.Fatal(deleted.Diagnostics)
			}
			missing := resource.ReadResponse{State: read.State}
			r.Read(ctx, resource.ReadRequest{State: read.State}, &missing)
			if missing.Diagnostics.HasError() || !missing.State.Raw.IsNull() {
				t.Fatalf("404 did not remove state: %v", missing.Diagnostics)
			}
		})
	}
}

func TestMonitorProfileChangesRequireReplacement(t *testing.T) {
	ctx := context.Background()
	var schema resource.SchemaResponse
	NewMonitorProfileResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	prior := collectorPlanTestObject(typ, map[string]any{"id": "profile", "name": "Platform", "config_ids": []tftypes.Value{}})
	for _, field := range []string{"name", "description", "config_ids", "parent_id"} {
		t.Run(field, func(t *testing.T) {
			var attrs map[string]tftypes.Value
			if err := prior.As(&attrs); err != nil {
				t.Fatal(err)
			}
			configured := make(map[string]tftypes.Value, len(attrs))
			for key, value := range attrs {
				configured[key] = value
			}
			attrs = configured
			value := any("changed")
			if field == "config_ids" {
				value = []tftypes.Value{tftypes.NewValue(tftypes.String, "shared")}
			}
			attrs[field] = tftypes.NewValue(typ.AttributeTypes[field], value)
			config := collectorPlanTestDynamic(t, tftypes.NewValue(typ, attrs))
			response, err := providerserver.NewProtocol6(New("test")())().PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName: "criblio_monitor_profile", PriorState: collectorPlanTestDynamic(t, prior), Config: config, ProposedNewState: config,
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, diagnostic := range response.Diagnostics {
				if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
					t.Fatalf("plan: %s: %s", diagnostic.Summary, diagnostic.Detail)
				}
			}
			if len(response.RequiresReplace) == 0 {
				t.Fatal("profile change must replace to preserve backend graph invariants")
			}
		})
	}
}
