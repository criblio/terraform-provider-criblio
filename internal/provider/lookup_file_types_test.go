package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestApplyLookupFileAPIToStatePreservesOmittedFields(t *testing.T) {
	api := LookupFileModel{
		GroupID: types.StringValue("default_search"),
		ID:      types.StringValue("test_lookup.csv"),
		Mode:    types.StringValue("memory"),
	}
	state := LookupFileModel{
		Content:     types.StringValue("a,b\n1,2"),
		Description: types.StringValue("test lookup"),
		GroupID:     types.StringValue("default_search"),
		ID:          types.StringValue("test_lookup.csv"),
		Mode:        types.StringValue("memory"),
		PendingTask: types.ObjectNull(LookupFilePendingTaskAttrTypes()),
		Tags:        types.StringValue("test"),
		Version:     types.StringValue("v1"),
	}

	applyLookupFileAPIToState(&api, &state, true, false)

	if got := state.Content.ValueString(); got != "a,b\n1,2" {
		t.Fatalf("content = %q", got)
	}
	if got := state.Description.ValueString(); got != "test lookup" {
		t.Fatalf("description = %q", got)
	}
	if got := state.Tags.ValueString(); got != "test" {
		t.Fatalf("tags = %q", got)
	}
	if got := state.Version.ValueString(); got != "v1" {
		t.Fatalf("version = %q", got)
	}
}
