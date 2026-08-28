package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNotificationTargetImportNormalizesSystemFieldsElementType(t *testing.T) {
	api := NotificationTargetModel{
		SnsTarget: &SnsTargetModel{SystemFields: types.ListNull(types.StringType)},
	}
	state := NotificationTargetModel{
		SnsTarget: &SnsTargetModel{SystemFields: types.List{}},
	}

	applyNotificationTargetAPIToState(&api, &state, true, true)

	elementType := state.SnsTarget.SystemFields.ElementType(context.Background())
	if elementType == nil || !elementType.Equal(types.StringType) {
		t.Fatalf("expected string list element type, got %v", elementType)
	}
}

func TestNotificationTargetCreateResolvesUnknownSystemFieldsFromAPI(t *testing.T) {
	api := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{
			SystemFields: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cribl_pipe")}),
		},
	}
	state := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{SystemFields: types.ListUnknown(types.StringType)},
	}

	applyNotificationTargetAPIToState(&api, &state, true, false)

	if state.SlackTarget.SystemFields.IsUnknown() {
		t.Fatal("expected system_fields to be known after create")
	}
	if !state.SlackTarget.SystemFields.Equal(api.SlackTarget.SystemFields) {
		t.Fatalf("system_fields = %v, want %v", state.SlackTarget.SystemFields, api.SlackTarget.SystemFields)
	}
}

func TestNotificationTargetCreateNormalizesUnknownSystemFieldsWhenAPIOmitsIt(t *testing.T) {
	api := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{SystemFields: types.ListNull(types.StringType)},
	}
	state := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{SystemFields: types.ListUnknown(types.StringType)},
	}

	applyNotificationTargetAPIToState(&api, &state, true, false)

	if !state.SlackTarget.SystemFields.IsNull() {
		t.Fatalf("system_fields = %v, want null", state.SlackTarget.SystemFields)
	}
}

func TestNotificationTargetCreatePreservesConfiguredSystemFields(t *testing.T) {
	configured := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("configured")})
	api := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{
			SystemFields: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("api")}),
		},
	}
	state := NotificationTargetModel{
		SlackTarget: &SlackTargetModel{SystemFields: configured},
	}

	applyNotificationTargetAPIToState(&api, &state, true, false)

	if !state.SlackTarget.SystemFields.Equal(configured) {
		t.Fatalf("system_fields = %v, want configured value %v", state.SlackTarget.SystemFields, configured)
	}
}
