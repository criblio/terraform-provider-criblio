package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestSearchDatasetMetadataValidation(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	var schema resource.SchemaResponse
	NewSearchDatasetResource().Schema(ctx, resource.SchemaRequest{}, &schema)
	typ := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	variants := 0
	for block, fieldType := range typ.AttributeTypes {
		variant, ok := fieldType.(tftypes.Object)
		if !ok {
			continue
		}
		metadata, ok := variant.AttributeTypes["metadata"].(tftypes.Object)
		if !ok {
			continue
		}
		enabledType, hasEnabled := metadata.AttributeTypes["enable_acceleration"]
		if !hasEnabled || !enabledType.Equal(tftypes.Bool) {
			continue
		}
		variants++
		for _, test := range []struct {
			name      string
			metadata  any
			wantError bool
		}{
			{"enabled-missing", map[string]any{"enable_acceleration": true}, true},
			{"enabled-explicit-null", map[string]any{"enable_acceleration": true, "earliest": nil}, true},
			{"enabled-configured", map[string]any{"enable_acceleration": true, "earliest": "-30d"}, false},
			{"disabled-missing", map[string]any{"enable_acceleration": false}, false},
			{"omitted-trigger", map[string]any{}, false},
			{"omitted-metadata", nil, false},
			{"unknown-trigger", map[string]any{"enable_acceleration": tftypes.UnknownValue}, false},
			{"unknown-earliest", map[string]any{"enable_acceleration": true, "earliest": tftypes.UnknownValue}, false},
			{"unknown-metadata", tftypes.UnknownValue, false},
		} {
			t.Run(block+"/"+test.name, func(t *testing.T) {
				config := collectorPlanTestObject(typ, map[string]any{block: map[string]any{"metadata": test.metadata}})
				response, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
					TypeName: "criblio_search_dataset", Config: collectorPlanTestDynamic(t, config),
				})
				if err != nil {
					t.Fatal(err)
				}
				errors := 0
				for _, diagnostic := range response.Diagnostics {
					if diagnostic.Severity != tfprotov6.DiagnosticSeverityError {
						continue
					}
					errors++
					if !test.wantError {
						t.Fatalf("unexpected validation error: %s: %s", diagnostic.Summary, diagnostic.Detail)
					}
					wantPath := tftypes.NewAttributePath().WithAttributeName(block).WithAttributeName("metadata").WithAttributeName("earliest")
					if diagnostic.Attribute == nil || !diagnostic.Attribute.Equal(wantPath) {
						t.Fatalf("diagnostic path = %v, want %v", diagnostic.Attribute, wantPath)
					}
					if !strings.Contains(diagnostic.Detail, "metadata.earliest is required when "+block+".metadata.enable_acceleration is true") {
						t.Fatalf("unclear diagnostic: %s", diagnostic.Detail)
					}
				}
				if (errors == 1) != test.wantError || errors > 1 {
					t.Fatalf("errors = %d, want error %t", errors, test.wantError)
				}
			})
		}
	}
	if variants == 0 {
		t.Fatal("no dataset metadata variants tested")
	}
	t.Logf("validated %d dataset metadata variants", variants)
}
