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
// marking and modifier ordering, for both resources that share collector fields.
func TestCollectorHoistedFieldsPlan(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	before := map[string]any{
		"ttl": "4h", "environment": "dev", "ignore_group_jobs_limit": false,
		"resume_on_boot": false, "worker_affinity": false,
	}
	after := map[string]any{
		"ttl": "8h", "environment": "prod", "ignore_group_jobs_limit": true,
		"resume_on_boot": true, "worker_affinity": true,
	}
	unknown := map[string]any{}
	for field := range before {
		unknown[field] = tftypes.UnknownValue
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
				name   string
				nested map[string]any
				root   map[string]any
				want   map[string]any
				create bool
			}{
				{name: "nested_update", nested: after, want: after},
				{name: "nested_unchanged", nested: before, want: before},
				{name: "nested_unknown", nested: unknown, want: unknown},
				{name: "nested_omitted", nested: map[string]any{}, want: before},
				{name: "root_update", nested: map[string]any{}, root: after, want: after},
				{name: "create", nested: after, want: after, create: true},
			} {
				t.Run(metadata.TypeName+"/"+kind+"/"+tc.name, func(t *testing.T) {
					stateFields := map[string]any{"id": "events", "group_id": "default", block: before}
					configFields := map[string]any{"id": "events", "group_id": "default", block: tc.nested}
					planFields := map[string]any{"id": "events", "group_id": "default", block: tc.nested}
					if metadata.TypeName == "criblio_pack_collector" {
						stateFields["pack"], configFields["pack"], planFields["pack"] = "my-pack", "my-pack", "my-pack"
					}
					for field, value := range before {
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
