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

func sourceTestStringList(values ...string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	return types.ListValueMust(types.StringType, elements)
}
