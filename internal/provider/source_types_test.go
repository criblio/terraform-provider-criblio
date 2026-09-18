package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSourceMarshalJSONInfersS3Type(t *testing.T) {
	payload, err := json.Marshal(SourceModel{
		ID: types.StringValue("wiz_runtime_s3"),
		InputS3: &InputS3Model{
			Type:      types.StringNull(),
			QueueName: types.StringValue("test-queue"),
		},
	})
	if err != nil {
		t.Fatalf("marshal S3 source: %v", err)
	}

	var output map[string]any
	if err := json.Unmarshal(payload, &output); err != nil {
		t.Fatalf("unmarshal S3 source payload: %v", err)
	}
	if got := output["type"]; got != "s3" {
		t.Fatalf("type = %#v, want s3; payload=%s", got, payload)
	}
}

func TestApplySourceAPIToStatePreservesImportedGroupID(t *testing.T) {
	state := SourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("diag_kafka"),
	}
	api := SourceModel{
		GroupID: types.StringValue("diag"),
		ID:      types.StringValue("diag_kafka"),
		InputKafka: &InputKafkaModel{
			GroupID: types.StringValue("diag"),
		},
	}

	applySourceAPIToState(&api, &state, true, true)

	if got := state.GroupID.ValueString(); got != "default" {
		t.Fatalf("root group_id = %q, want default", got)
	}
	if state.InputKafka == nil {
		t.Fatal("input_kafka was not hydrated")
	}
	if got := state.InputKafka.GroupID.ValueString(); got != "diag" {
		t.Fatalf("input_kafka.group_id = %q, want diag", got)
	}
}
