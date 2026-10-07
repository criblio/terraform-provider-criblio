package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Exercise the framework's complete planning pipeline, including its unknown
// marking and modifier ordering, for TTL in both collector resources.
func TestCollectorTTLPlan(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	before := map[string]any{"ttl": "4h"}
	after := map[string]any{"ttl": "8h"}
	unknown := map[string]any{}
	nulls := map[string]any{}
	for field := range before {
		unknown[field] = tftypes.UnknownValue
		nulls[field] = nil
	}
	// Discovery fixtures configure these two fields but omit TTL and the other
	// job flags; the S3 fixture omits all five. Those omissions must remain stable.
	discoveryDefaults := map[string]any{
		"environment": "demo", "ignore_group_jobs_limit": false,
		"ttl": nil, "resume_on_boot": nil, "worker_affinity": nil,
	}
	for _, res := range []resource.Resource{NewCollectorResource(), NewPackCollectorResource()} {
		var metadata resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "criblio"}, &metadata)
		var schema resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schema)
		objectType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
		for _, kind := range []string{"splunk", "rest", "s3", "azure_blob", "cribl_lake", "database", "gcs", "health_check", "script", "filesystem"} {
			block := "input_collector_" + kind
			for _, tc := range []struct {
				name      string
				nested    map[string]any
				root      map[string]any
				prior     map[string]any
				proposed  map[string]any
				want      map[string]any
				create    bool
				unchanged bool
			}{
				{name: "nested_update", nested: after, want: after},
				{name: "nested_unchanged", nested: before, want: before, unchanged: true},
				{name: "nested_unknown", nested: unknown, want: unknown},
				{name: "nested_omitted", nested: map[string]any{}, proposed: before, want: before, unchanged: true},
				{name: "root_update", nested: map[string]any{}, root: after, want: after},
				{name: "create", nested: after, want: after, create: true},
				{name: "nulls_unchanged", nested: nulls, prior: nulls, want: nulls, unchanged: true},
				{name: "nulls_update", nested: after, prior: nulls, want: after},
				{name: "nulls_unknown", nested: unknown, prior: nulls, want: unknown},
				{name: "nulls_create", nested: nulls, prior: nulls, want: unknown, create: true},
				{name: "discovery_omitted_settings", nested: discoveryDefaults, prior: discoveryDefaults, want: discoveryDefaults, unchanged: true},
			} {
				t.Run(metadata.TypeName+"/"+kind+"/"+tc.name, func(t *testing.T) {
					prior := tc.prior
					if prior == nil {
						prior = before
					}
					stateFields := map[string]any{"id": "events", "group_id": "default", block: prior}
					configFields := map[string]any{"id": "events", "group_id": "default", block: tc.nested}
					planFields := map[string]any{"id": "events", "group_id": "default", block: tc.nested}
					if tc.proposed != nil {
						planFields[block] = tc.proposed
					}
					if metadata.TypeName == "criblio_pack_collector" {
						stateFields["pack"], configFields["pack"], planFields["pack"] = "my-pack", "my-pack", "my-pack"
					}
					for field, value := range prior {
						stateFields[field], planFields[field] = value, value
					}
					for field, value := range tc.root {
						configFields[field], planFields[field] = value, value
					}
					state := collectorPlanTestObject(objectType, stateFields)
					if tc.create {
						state = tftypes.NewValue(objectType, nil)
					}
					config := collectorPlanTestObject(objectType, configFields)
					plan := collectorPlanTestObject(objectType, planFields)
					response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
						TypeName: metadata.TypeName, PriorState: collectorPlanTestDynamic(t, state),
						Config: collectorPlanTestDynamic(t, config), ProposedNewState: collectorPlanTestDynamic(t, plan),
					})
					if err != nil {
						t.Fatal(err)
					}
					for _, diagnostic := range response.Diagnostics {
						if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
							t.Fatalf("plan: %s: %s", diagnostic.Summary, diagnostic.Detail)
						}
					}
					planned, err := response.PlannedState.Unmarshal(objectType)
					if err != nil {
						t.Fatal(err)
					}
					if tc.unchanged && !planned.Equal(state) {
						t.Error("unchanged configuration produced a non-empty plan")
					}
					var fields map[string]tftypes.Value
					if err := planned.As(&fields); err != nil {
						t.Fatal(err)
					}
					for field, value := range tc.want {
						want := tftypes.NewValue(objectType.AttributeTypes[field], value)
						if !fields[field].Equal(want) {
							t.Errorf("%s planned %s, want %s", field, fields[field], want)
						}
					}
				})
			}
		}
	}
}

func collectorPlanTestObject(typ tftypes.Object, values map[string]any) tftypes.Value {
	fields := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for name, fieldType := range typ.AttributeTypes {
		if nested, ok := values[name].(map[string]any); ok {
			fields[name] = collectorPlanTestObject(fieldType.(tftypes.Object), nested)
			continue
		}
		fields[name] = tftypes.NewValue(fieldType, values[name])
	}
	return tftypes.NewValue(typ, fields)
}

func collectorPlanTestDynamic(t *testing.T, value tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dynamic, err := tfprotov6.NewDynamicValue(value.Type(), value)
	if err != nil {
		t.Fatal(err)
	}
	return &dynamic
}
