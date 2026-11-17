package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestWorkerGroupUpgradeSubmissionResponses(t *testing.T) {
	for _, tc := range []struct {
		name, body, wantError, groupID string
		status                         int
		wantJobs                       int
	}{
		{name: "job-envelope", status: 200, body: `{"count":2,"items":["job-1","job-2"]}`, wantJobs: 2},
		{name: "no-jobs", status: 200, body: `{"count":0,"items":[]}`},
		{name: "no-body", status: 200},
		{name: "managed-group-rejected", status: 400, body: `{"message":"Managed Worker Groups cannot be upgraded"}`, wantError: "Managed Worker Groups"},
		{name: "reported-managed-group", groupID: "default", status: 400, body: `{"status":"error","message":"Group 'default' is a managed group."}`, wantError: "Set group_id to an existing self-managed Worker Group"},
		{name: "self-managed-default", groupID: "default", status: 200, body: `{"items":["job-1"]}`, wantJobs: 1},
		{name: "invalid-settings", status: 400, body: `{"message":"Invalid upgrade percentage"}`, wantError: "Invalid upgrade percentage"},
		{name: "not-authorized", status: 403, body: `{"message":"Forbidden"}`, wantError: "Forbidden"},
		{name: "missing-group", status: 404, body: `{"message":"Worker Group not found"}`, wantError: "Worker Group not found"},
		{name: "invalid-response", status: 200, body: `{"items":[{}]}`, wantError: "decode"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			groupID := tc.groupID
			if groupID == "" {
				groupID = "group/with space"
			}
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls++
				if req.Method != http.MethodPost || req.URL.EscapedPath() != "/api/v1/system/distributed/upgrade/"+url.PathEscape(groupID) {
					t.Errorf("incorrect escaped route: %s %s", req.Method, req.URL)
				}
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}))
			defer server.Close()
			api := newWorkerGroupUpgradeAPI(restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"}))
			model := WorkerGroupUpgradeModel{GroupID: types.StringValue(groupID), Trigger: types.StringValue("stable"), UpgradeMode: types.StringValue("rolling")}
			result, err := api.Create(context.Background(), model)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) || result != nil {
					t.Fatalf("result=%v error=%v; want error containing %q and no state", result, err, tc.wantError)
				}
				if tc.status == http.StatusBadRequest {
					var httpErr *restclient.HTTPError
					if !errors.As(err, &httpErr) || httpErr.StatusCode != tc.status {
						t.Fatalf("original HTTP error was lost: %v", err)
					}
				}
			} else {
				if err != nil || result == nil {
					t.Fatalf("submission failed: %v", err)
				}
				if result.JobIDs.IsNull() || result.JobIDs.IsUnknown() || len(result.JobIDs.Elements()) != tc.wantJobs || !result.Trigger.Equal(model.Trigger) || !result.GroupID.Equal(model.GroupID) {
					t.Fatalf("incorrect submission state: %+v", result)
				}
			}
			if calls != 1 {
				t.Fatalf("POST calls = %d; want 1", calls)
			}
		})
	}
}

func TestWorkerGroupUpgradeConfigValidation(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	var schema resource.SchemaResponse
	NewWorkerGroupUpgradeResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	for _, tc := range []struct {
		name, field string
		value       any
		wantError   bool
	}{
		{"minimum-group", "group_id", "", true},
		{"rolling", "upgrade_mode", "rolling", false},
		{"regular", "upgrade_mode", "regular", false},
		{"batch", "upgrade_mode", "batch", false},
		{"invalid-mode", "upgrade_mode", "parallel", true},
		{"percentage-zero", "upgrade_percentage", float64(0), true},
		{"percentage-too-low", "upgrade_percentage", 0.000001, true},
		{"percentage-minimum", "upgrade_percentage", 0.00001, false},
		{"percentage-full", "upgrade_percentage", float64(1), false},
		{"percentage-not-percent", "upgrade_percentage", float64(50), true},
		{"retries-zero", "worker_retries", float64(0), true},
		{"retries-minimum", "worker_retries", float64(1), false},
		{"delay-too-low", "worker_retry_delay", float64(99), true},
		{"delay-minimum", "worker_retry_delay", float64(100), false},
		{"unknown-trigger", "trigger", tftypes.UnknownValue, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := map[string]any{"group_id": "default"}
			fields[tc.field] = tc.value
			config := collectorPlanTestDynamic(t, collectorPlanTestObject(typ, fields))
			response, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{TypeName: "criblio_worker_group_upgrade", Config: config})
			if err != nil {
				t.Fatal(err)
			}
			hasError := false
			for _, diagnostic := range response.Diagnostics {
				if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
					hasError = true
					if !tc.wantError {
						t.Errorf("unexpected diagnostic: %s: %s", diagnostic.Summary, diagnostic.Detail)
					}
				}
			}
			if hasError != tc.wantError {
				t.Fatalf("validation errors = %t; want %t", hasError, tc.wantError)
			}
		})
	}
	var imported resource.ImportStateResponse
	NewWorkerGroupUpgradeResource().(resource.ResourceWithImportState).ImportState(ctx, resource.ImportStateRequest{ID: "default"}, &imported)
	if !imported.Diagnostics.HasError() {
		t.Fatal("asynchronous actions must reject import")
	}
}
