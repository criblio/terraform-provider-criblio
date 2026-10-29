package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Exercise framework unknown marking and resource plan modification for every
// dataset variant, discovering variants from the actual resource schema.
func TestSearchDatasetIDPlan(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	var schema resource.SchemaResponse
	NewSearchDatasetResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	objectType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	variants := 0
	for block, typ := range objectType.AttributeTypes {
		if _, ok := typ.(tftypes.Object); !ok {
			continue
		}
		variants++
		for _, tc := range []struct {
			name         string
			id           any
			priorID      string
			want         any
			replace      bool
			unknownBlock bool
		}{
			{name: "known", id: "logs", want: "logs"},
			{name: "unknown", id: tftypes.UnknownValue, want: tftypes.UnknownValue},
			{name: "omitted", want: tftypes.UnknownValue},
			{name: "unknown block", unknownBlock: true, want: tftypes.UnknownValue},
			{name: "imported", id: "logs", priorID: "logs", want: "logs"},
			{name: "changed", id: "new-logs", priorID: "logs", want: "new-logs", replace: true},
			{name: "unknown update", id: tftypes.UnknownValue, priorID: "logs", want: tftypes.UnknownValue, replace: true},
			{name: "unknown block update", unknownBlock: true, priorID: "logs", want: tftypes.UnknownValue, replace: true},
			{name: "omitted update", priorID: "logs", want: tftypes.UnknownValue},
		} {
			t.Run(block+"/"+tc.name, func(t *testing.T) {
				var nested any = map[string]any{"id": tc.id}
				if tc.unknownBlock {
					nested = tftypes.UnknownValue
				}
				config := collectorPlanTestObject(objectType, map[string]any{block: nested})
				plan := collectorPlanTestObject(objectType, map[string]any{block: nested, "id": tftypes.UnknownValue})
				state := tftypes.NewValue(objectType, nil)
				if tc.priorID != "" {
					state = collectorPlanTestObject(objectType, map[string]any{block: map[string]any{"id": tc.priorID}, "id": tc.priorID})
				}
				response, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
					TypeName: "criblio_search_dataset", PriorState: collectorPlanTestDynamic(t, state),
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
				if want := tftypes.NewValue(tftypes.String, tc.want); !fields["id"].Equal(want) {
					t.Errorf("planned ID = %s, want %s", fields["id"], want)
				}
				if got := len(response.RequiresReplace) > 0; got != tc.replace {
					t.Errorf("requires replacement = %t, want %t", got, tc.replace)
				}
			})
		}
	}
	if variants == 0 {
		t.Fatal("no dataset variants tested")
	}
}

func TestSearchDatasetIDDestroyPlan(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	var schema resource.SchemaResponse
	NewSearchDatasetResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	null := collectorPlanTestDynamic(t, tftypes.NewValue(typ, nil))
	state := collectorPlanTestDynamic(t, collectorPlanTestObject(typ, map[string]any{
		"id": "logs", "s3_dataset": map[string]any{"id": "logs"},
	}))
	resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName: "criblio_search_dataset", Config: null, ProposedNewState: null, PriorState: state,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range resp.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("destroy plan: %s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
	plan, err := resp.PlannedState.Unmarshal(typ)
	if err != nil {
		t.Fatal(err)
	}
	if !plan.IsNull() || len(resp.RequiresReplace) != 0 {
		t.Fatal("destroy plan was modified")
	}
}
