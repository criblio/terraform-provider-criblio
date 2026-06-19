package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestApplyAPIToStateTypesOmittedNestedObjectLists(t *testing.T) {
	api := &EventBreakerRulesetModel{
		Rules: types.ListNull(types.ObjectType{AttrTypes: EventBreakerRulesetRulesAttrTypes()}),
	}
	state := &EventBreakerRulesetModel{}

	applyEventBreakerRulesetAPIToState(api, state, false, false)

	want := types.ObjectType{AttrTypes: EventBreakerRulesetRulesAttrTypes()}
	if !reflect.DeepEqual(state.Rules.ElementType(context.Background()), want) {
		t.Fatalf("rules element type = %v", state.Rules.ElementType(context.Background()))
	}
}
