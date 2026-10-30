package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMonitorUpstreamRequiredConfiguration(t *testing.T) {
	ctx := context.Background()
	var schema resource.SchemaResponse
	NewMonitorResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	base := collectorPlanTestObject(typ, map[string]any{
		"id": "monitor-validation-test", "name": "Monitor validation test", "dataset_id": "metrics",
		"enabled": true, "type": "threshold", "managed_by": "external-provisioner",
		"expr": []tftypes.Value{}, "silence": []tftypes.Value{},
		"firing_condition":                   map[string]any{"fire_delay": float64(300), "clear_delay": float64(60)},
		"firing_rule":                        map[string]any{"label": "A", "threshold": []tftypes.Value{}},
		"aetos_metadata":                     map[string]tftypes.Value{},
		"monitor_policy_notification_config": map[string]any{"config": []tftypes.Value{}, "enabled": false, "type": "policy"},
		"priority_scalar_value":              map[string]any{"value": "P3"}, "team_scalar_value": map[string]any{"value": ""},
	})
	var values map[string]tftypes.Value
	if err := base.As(&values); err != nil {
		t.Fatal(err)
	}
	queryType := typ.AttributeTypes["query"].(tftypes.Map)
	values["query"] = tftypes.NewValue(queryType, map[string]tftypes.Value{
		"A": collectorPlanTestObject(queryType.ElementType.(tftypes.Object), map[string]any{"mode": "promql", "promql": "cpu_usage_percent"}),
	})
	for _, missing := range []string{"", "query", "expr", "silence", "firing_condition", "firing_rule", "aetos_metadata", "monitor_policy_notification_config", "priority_scalar_value", "team_scalar_value"} {
		name := missing
		if name == "" {
			name = "valid"
		}
		t.Run(name, func(t *testing.T) {
			configured := make(map[string]tftypes.Value, len(values))
			for key, value := range values {
				configured[key] = value
			}
			if missing != "" {
				configured[missing] = tftypes.NewValue(typ.AttributeTypes[missing], nil)
			}
			server := providerserver.NewProtocol6(New("test")())()
			response, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
				TypeName: "criblio_monitor", Config: collectorPlanTestDynamic(t, tftypes.NewValue(typ, configured)),
			})
			if err != nil {
				t.Fatal(err)
			}
			hasError := false
			for _, diagnostic := range response.Diagnostics {
				if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
					hasError = true
					if missing == "" {
						t.Errorf("valid configuration: %s: %s", diagnostic.Summary, diagnostic.Detail)
					}
				}
			}
			if missing != "" && !hasError {
				t.Fatalf("missing required %s was accepted", missing)
			}
		})
	}
}

func TestMonitorUpstreamWireFields(t *testing.T) {
	const input = `{
		"id":"wire-test", "name":"Wire test", "datasetId":"metrics", "enabled":true, "type":"threshold",
		"managedBy":"external-provisioner", "expr":[], "query":{"A":{"mode":"promql","promql":"cpu_usage_percent"}}, "silence":[],
		"$profileRef":{"id":"profile", "fingerprint":"profile-fingerprint"},
		"firingCondition":{"$configRefs":[{"id":"condition", "fingerprint":"condition-fingerprint"}]},
		"firingRule":{"label":"A", "threshold":[]}, "metadata":{},
		"notification":{"type":"policy", "enabled":false, "config":[]}, "priority":{"value":"P3"}, "team":{"value":""},
		"templateParams":{"A":{"$configRefs":[{"id":"parameters", "fingerprint":"parameters-fingerprint"}]}}
	}`
	var model MonitorModel
	if err := json.Unmarshal([]byte(input), &model); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	var want, got map[string]any
	if err := json.Unmarshal([]byte(input), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"$profileRef", "firingCondition", "templateParams", "managedBy", "silence", "expr", "notification"} {
		wantJSON, err := json.Marshal(want[key])
		if err != nil {
			t.Fatal(err)
		}
		gotJSON, err := json.Marshal(got[key])
		if err != nil {
			t.Fatal(err)
		}
		if string(gotJSON) != string(wantJSON) {
			t.Errorf("%s: got %s, want %s", key, gotJSON, wantJSON)
		}
	}
}

// Existing API-shaped configurations must plan without query or expression synthesis.
func TestMonitorPlanAcceptsUpstreamQueryAndExpr(t *testing.T) {
	ctx := context.Background()
	var schema resource.SchemaResponse
	NewMonitorResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	for _, name := range []string{"query", "expr"} {
		attribute := schema.Schema.Attributes[name]
		if attribute == nil || !attribute.IsRequired() || attribute.IsComputed() {
			t.Fatalf("%s must remain a required, user-configurable upstream attribute", name)
		}
	}
	for _, name := range []string{"queries", "alert_expression"} {
		if _, exists := schema.Schema.Attributes[name]; exists {
			t.Fatalf("%s is not part of the upstream monitor schema", name)
		}
	}

	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	queryType := typ.AttributeTypes["query"].(tftypes.Map)
	exprType := typ.AttributeTypes["expr"].(tftypes.List)
	query := tftypes.NewValue(queryType, map[string]tftypes.Value{
		"A": collectorPlanTestObject(queryType.ElementType.(tftypes.Object), map[string]any{
			"mode": "promql", "dataset_id": "metrics", "promql": "cpu_usage_percent",
		}),
	})
	for _, monitorType := range []string{"threshold", "forecast", "outlier"} {
		for _, withExpression := range []bool{false, true} {
			name := "empty-expr"
			expressions := []tftypes.Value{}
			if withExpression {
				name = "configured-expr"
				expressions = append(expressions, collectorPlanTestObject(exprType.ElementType.(tftypes.Object), map[string]any{
					"label": "Alert Query", "text": "A", "query_labels": []tftypes.Value{tftypes.NewValue(tftypes.String, "A")},
				}))
			}
			t.Run(monitorType+"/"+name, func(t *testing.T) {
				expr := tftypes.NewValue(exprType, expressions)
				base := collectorPlanTestObject(typ, map[string]any{
					"id": "monitor-plan-test", "name": "Monitor plan test", "dataset_id": "metrics",
					"enabled": true, "type": monitorType,
				})
				var configured map[string]tftypes.Value
				if err := base.As(&configured); err != nil {
					t.Fatal(err)
				}
				configured["query"] = query
				configured["expr"] = expr
				config := tftypes.NewValue(typ, configured)
				server := providerserver.NewProtocol6(New("test")())()
				response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
					TypeName:   "criblio_monitor",
					PriorState: collectorPlanTestDynamic(t, tftypes.NewValue(typ, nil)),
					Config:     collectorPlanTestDynamic(t, config), ProposedNewState: collectorPlanTestDynamic(t, config),
				})
				if err != nil {
					t.Fatal(err)
				}
				for _, diagnostic := range response.Diagnostics {
					if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
						t.Fatalf("plan: %s: %s", diagnostic.Summary, diagnostic.Detail)
					}
				}
				planned, err := response.PlannedState.Unmarshal(typ)
				if err != nil {
					t.Fatal(err)
				}
				var fields map[string]tftypes.Value
				if err := planned.As(&fields); err != nil {
					t.Fatal(err)
				}
				if !fields["query"].Equal(query) || !fields["expr"].Equal(expr) {
					t.Fatal("plan changed the configured query or expr")
				}
			})
		}
	}
}
