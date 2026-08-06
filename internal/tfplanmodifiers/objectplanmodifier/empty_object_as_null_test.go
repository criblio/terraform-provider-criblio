package objectplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestEmptyObjectAsNull(t *testing.T) {
	objectType := map[string]attr.Type{"enabled": types.BoolType}
	emptyState := types.ObjectValueMust(objectType, map[string]attr.Value{"enabled": types.BoolNull()})
	nonEmptyState := types.ObjectValueMust(objectType, map[string]attr.Value{"enabled": types.BoolValue(true)})

	tests := []struct {
		name      string
		config    types.Object
		state     types.Object
		wantState bool
	}{
		{name: "null config and empty state", config: types.ObjectNull(objectType), state: emptyState, wantState: true},
		{name: "null config and populated state", config: types.ObjectNull(objectType), state: nonEmptyState},
		{name: "configured object", config: nonEmptyState, state: emptyState},
	}

	modifier := EmptyObjectAsNull()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := planmodifier.ObjectRequest{ConfigValue: tt.config, PlanValue: tt.config, StateValue: tt.state}
			response := &planmodifier.ObjectResponse{PlanValue: tt.config}
			modifier.PlanModifyObject(context.Background(), request, response)
			assert.Equal(t, tt.wantState, response.PlanValue.Equal(tt.state))
		})
	}
}

func TestPreferConfigOrStateWithNullConfig(t *testing.T) {
	objectType := map[string]attr.Type{"enabled": types.BoolType}
	nullObject := types.ObjectNull(objectType)
	refreshedObject := types.ObjectValueMust(objectType, map[string]attr.Value{"enabled": types.BoolValue(true)})

	modifier := PreferConfigOrState()
	for _, state := range []types.Object{nullObject, refreshedObject} {
		request := planmodifier.ObjectRequest{
			ConfigValue: nullObject,
			PlanValue:   types.ObjectUnknown(objectType),
			State:       tfsdk.State{Raw: tftypes.NewValue(tftypes.String, "existing")},
			StateValue:  state,
		}
		response := &planmodifier.ObjectResponse{PlanValue: request.PlanValue}
		modifier.PlanModifyObject(context.Background(), request, response)
		assert.True(t, response.PlanValue.Equal(state))
	}

	request := planmodifier.ObjectRequest{
		ConfigValue: nullObject,
		PlanValue:   types.ObjectUnknown(objectType),
		State:       tfsdk.State{Raw: tftypes.NewValue(tftypes.String, nil)},
		StateValue:  nullObject,
	}
	response := &planmodifier.ObjectResponse{PlanValue: request.PlanValue}
	modifier.PlanModifyObject(context.Background(), request, response)
	assert.True(t, response.PlanValue.IsUnknown(), "new resource should allow API defaults")
}
