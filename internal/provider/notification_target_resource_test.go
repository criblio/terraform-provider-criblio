package provider

import (
	"context"
	"testing"

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
