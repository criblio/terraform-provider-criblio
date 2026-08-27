package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestCollectorSchemaKeepsInputOptionalForBackwardCompatibility(t *testing.T) {
	t.Parallel()

	var response resource.SchemaResponse
	(&CollectorResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)

	for _, variantName := range []string{
		"input_collector_splunk",
		"input_collector_rest",
		"input_collector_s3",
		"input_collector_azure_blob",
		"input_collector_cribl_lake",
		"input_collector_database",
		"input_collector_gcs",
		"input_collector_health_check",
		"input_collector_script",
		"input_collector_filesystem",
	} {
		t.Run(variantName, func(t *testing.T) {
			variant, ok := response.Schema.Attributes[variantName].(resourceschema.SingleNestedAttribute)
			if !ok {
				t.Fatalf("%s is not a single nested attribute", variantName)
			}
			input, ok := variant.Attributes["input"].(resourceschema.SingleNestedAttribute)
			if !ok {
				t.Fatalf("%s.input is not a single nested attribute", variantName)
			}
			if input.Required || !input.Optional || !input.Computed {
				t.Errorf("%s.input flags = required:%v optional:%v computed:%v", variantName, input.Required, input.Optional, input.Computed)
			}
		})
	}
}
