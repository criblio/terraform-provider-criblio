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

func TestApplyLookupFileAPIToStateRefreshesDownloadedContent(t *testing.T) {
	api := LookupFileModel{
		Content: types.StringValue("a,b\nremote,value"),
		GroupID: types.StringValue("default_search"),
		ID:      types.StringValue("test_lookup.csv"),
	}
	state := LookupFileModel{
		Content:     types.StringValue("a,b\nold,value"),
		Description: types.StringValue("preserved"),
		GroupID:     types.StringValue("default_search"),
		ID:          types.StringValue("test_lookup.csv"),
		PendingTask: types.ObjectNull(LookupFilePendingTaskAttrTypes()),
	}

	applyLookupFileAPIToState(&api, &state, true, false)

	if got := state.Content.ValueString(); got != "a,b\nremote,value" {
		t.Fatalf("content = %q", got)
	}
	if got := state.Description.ValueString(); got != "preserved" {
		t.Fatalf("description = %q", got)
	}
}

func TestApplyPackLookupsAPIToStateRefreshesDownloadedContent(t *testing.T) {
	api := PackLookupsModel{
		Content: types.StringValue("a,b\nremote,value"),
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("pack_lookup.csv"),
		Pack:    types.StringValue("knowledge"),
	}
	state := PackLookupsModel{
		Content:     types.StringValue("a,b\nold,value"),
		Description: types.StringValue("preserved"),
		GroupID:     types.StringValue("default"),
		ID:          types.StringValue("pack_lookup.csv"),
		Pack:        types.StringValue("knowledge"),
		PendingTask: types.ObjectNull(PackLookupsPendingTaskAttrTypes()),
	}

	applyPackLookupsAPIToState(&api, &state, true, false)

	if got := state.Content.ValueString(); got != "a,b\nremote,value" {
		t.Fatalf("content = %q", got)
	}
	if got := state.Description.ValueString(); got != "preserved" {
		t.Fatalf("description = %q", got)
	}
}
