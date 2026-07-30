package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSourceRequiredFieldsRemainRequired(t *testing.T) {
	var response resource.SchemaResponse
	(&SourceResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("source schema diagnostics: %v", response.Diagnostics)
	}

	tests := []struct {
		variant string
		field   string
	}{
		{variant: "input_elastic", field: "elastic_api"},
		{variant: "input_splunk_hec", field: "splunk_hec_api"},
		{variant: "input_prometheus_rw", field: "prometheus_api"},
		{variant: "input_sqs", field: "queue_type"},
	}
	for _, test := range tests {
		t.Run(test.variant+"."+test.field, func(t *testing.T) {
			variant := response.Schema.Attributes[test.variant].(schema.SingleNestedAttribute)
			attribute := variant.Attributes[test.field].(schema.StringAttribute)
			if !attribute.Required || attribute.Optional || attribute.Computed {
				t.Fatalf("attribute must be required: %#v", attribute)
			}
		})
	}
}

func TestSourceMetadataUsesLowercaseAPIName(t *testing.T) {
	metadata := types.ObjectValueMust(InputHttpMetadataAttrTypes(), map[string]attr.Value{
		"name":  types.StringValue("source"),
		"value": types.StringValue("http"),
	})
	value, err := SourceTerraformValueToJSON(metadata)
	if err != nil {
		t.Fatalf("convert metadata: %v", err)
	}
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal metadata: %v", err)
	}
	if string(raw) != `{"name":"source","value":"http"}` {
		t.Fatalf("unexpected metadata JSON: %s", raw)
	}
}

func TestPackSourceTCPMetadataPayloadUsesLowercaseAPIName(t *testing.T) {
	metadata := types.ObjectValueMust(InputTcpMetadataAttrTypes(), map[string]attr.Value{
		"name":  types.StringValue("my_name"),
		"value": types.StringValue(`"my_value"`),
	})
	model := InputTcpModel{
		Metadata: types.ListValueMust(
			types.ObjectType{AttrTypes: InputTcpMetadataAttrTypes()},
			[]attr.Value{metadata},
		),
	}
	payload, err := model.terraformPayload()
	if err != nil {
		t.Fatalf("build TCP payload: %v", err)
	}
	raw, err := json.Marshal(payload["metadata"])
	if err != nil {
		t.Fatalf("marshal TCP metadata: %v", err)
	}
	if string(raw) != `[{"name":"my_name","value":"\"my_value\""}]` {
		t.Fatalf("unexpected TCP metadata JSON: %s", raw)
	}
}
