package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRestoreSourcePlainAuthTokensIfAPIShrank(t *testing.T) {
	model := SourceModel{
		InputHttp: &InputHttpModel{
			AuthTokens: sourceTestStringList("old-token"),
		},
	}
	prior := priorSourcePlainAuthTokens{
		http: []types.String{types.StringValue("old-token"), types.StringValue("new-token")},
	}

	restoreSourcePlainAuthTokensIfAPIShrank(&model, prior)

	got := sourceStringListValues(model.InputHttp.AuthTokens)
	if len(got) != 2 {
		t.Fatalf("expected prior auth tokens to be restored, got %d", len(got))
	}
	if got[0].ValueString() != "old-token" || got[1].ValueString() != "new-token" {
		t.Fatalf("unexpected restored auth tokens: %#v", got)
	}
}

func TestSourceRequestModelWithHoistedIdentity(t *testing.T) {
	model := SourceModel{
		ID:        types.StringValue("source-id"),
		InputHttp: &InputHttpModel{Type: types.StringValue("http")},
	}

	request := sourceRequestModelWithHoistedIdentity(model)

	if request.InputHttp.ID.ValueString() != "source-id" {
		t.Fatalf("expected active input id to be hoisted, got %q", request.InputHttp.ID.ValueString())
	}
	if !model.InputHttp.ID.IsNull() {
		t.Fatalf("expected planned model to remain unchanged, got %q", model.InputHttp.ID.ValueString())
	}
}

func TestDestinationRequestModelWithHoistedIdentity(t *testing.T) {
	model := DestinationModel{
		ID:       types.StringValue("destination-id"),
		OutputS3: &OutputS3Model{Type: types.StringValue("s3")},
	}

	request := oneOfRequestModelWithHoistedIdentity(model)

	if request.OutputS3.ID.ValueString() != "destination-id" {
		t.Fatalf("expected active output id to be hoisted, got %q", request.OutputS3.ID.ValueString())
	}
	if !model.OutputS3.ID.IsNull() {
		t.Fatalf("expected planned model to remain unchanged, got %q", model.OutputS3.ID.ValueString())
	}
}

func TestCollectorRequestModelWithHoistedIdentity(t *testing.T) {
	model := CollectorModel{
		ID:                   types.StringValue("collector-id"),
		InputCollectorSplunk: &InputCollectorSplunkModel{},
	}

	request := oneOfRequestModelWithHoistedIdentity(model)

	if request.InputCollectorSplunk.ID.ValueString() != "collector-id" {
		t.Fatalf("expected active collector id to be hoisted, got %q", request.InputCollectorSplunk.ID.ValueString())
	}
}

func TestOneOfRequestModelHoistsIdentityIntoEveryActiveVariant(t *testing.T) {
	model := DestinationModel{
		ID:            types.StringValue("destination-id"),
		OutputDefault: &OutputDefaultModel{ID: types.StringValue("destination-id")},
		OutputS3:      &OutputS3Model{},
	}

	request := oneOfRequestModelWithHoistedIdentity(model)

	if request.OutputS3.ID.ValueString() != "destination-id" {
		t.Fatalf("expected later active output id to be hoisted, got %q", request.OutputS3.ID.ValueString())
	}
}

func TestSyncSourceLikeActiveInputClearsZeroValueInactiveVariant(t *testing.T) {
	api := SourceModel{InputCloudflareHec: &InputCloudflareHecModel{}}
	state := SourceModel{
		InputCloudflareHec: &InputCloudflareHecModel{},
		InputSplunk:        &InputSplunkModel{},
	}

	syncSourceLikeActiveInput(&api, &state)

	if state.InputSplunk != nil {
		t.Fatal("expected inactive input_splunk variant to be cleared")
	}
	if state.InputCloudflareHec == nil {
		t.Fatal("expected active input_cloudflare_hec variant to be preserved")
	}
}

func sourceTestStringList(values ...string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	return types.ListValueMust(types.StringType, elements)
}
