package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Guard the fields that the 4.20.3 upstream refresh would otherwise remove or
// newly require. Both saved state and existing configurations depend on them.
func TestSpecRefreshPreservesLegacyAttributes(t *testing.T) {
	response, err := providerserver.NewProtocol6(New("test")())().GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("%s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
	for _, tt := range []struct {
		name       string
		path       string
		dataSource bool
	}{
		{"criblio_app", "cribl.hidden", false},
		{"criblio_app", "cribl.hidden", true},
		{"criblio_search_datatype", "schema_map", false},
		{"criblio_search_datatype", "schema_map", true},
		{"criblio_system_info", "items.is_highside", true},
		{"criblio_search_datatype", "additional_extractions.name", false},
		{"criblio_search_datatype", "additional_extractions.source_field", false},
		{"criblio_collector", "input_collector_database.schedule.run.state_tracking.tracking_column", false},
		{"criblio_pack_collector", "input_collector_database.schedule.run.state_tracking.tracking_column", false},
	} {
		kind := "resource"
		if tt.dataSource {
			kind = "data"
		}
		t.Run(tt.name+"/"+tt.path+"/"+kind, func(t *testing.T) {
			schemas := response.ResourceSchemas
			if tt.dataSource {
				schemas = response.DataSourceSchemas
			}
			schema, ok := schemas[tt.name]
			if !ok {
				t.Fatalf("missing schema %s", tt.name)
			}
			attributes := schema.Block.Attributes
			var found *tfprotov6.SchemaAttribute
			for _, name := range strings.Split(tt.path, ".") {
				found = nil
				for _, attribute := range attributes {
					if attribute.Name == name {
						found = attribute
						break
					}
				}
				if found == nil {
					t.Fatalf("missing %s.%s", tt.name, tt.path)
				}
				if found.NestedType != nil {
					attributes = found.NestedType.Attributes
				} else {
					attributes = nil
				}
			}
			if found.Required {
				t.Fatalf("%s.%s must not require a new configuration value", tt.name, tt.path)
			}
		})
	}
}
