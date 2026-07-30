package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGroupSystemSettingsDetectsIDOnlyImportState(t *testing.T) {
	state := GroupSystemSettingsModel{GroupID: types.StringValue("my-hybrid-group")}
	if !isGroupSystemSettingsImportState(&state) {
		t.Fatal("expected ID-only group system settings state to be detected as import")
	}

	state.API = types.ObjectNull(GroupSystemSettingsAPIFieldAttrTypes())
	if !isGroupSystemSettingsImportState(&state) {
		t.Fatal("expected typed-null settings to remain import state")
	}
}
