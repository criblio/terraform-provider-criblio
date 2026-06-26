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

func TestApplyAPIToStateTypesOmittedNestedObjects(t *testing.T) {
	api := &SubscriptionModel{
		Consumer: types.ObjectNull(SubscriptionConsumerAttrTypes()),
	}
	state := &SubscriptionModel{}

	applySubscriptionAPIToState(api, state, false, false)

	want := types.ObjectType{AttrTypes: SubscriptionConsumerAttrTypes()}
	if !reflect.DeepEqual(state.Consumer.Type(context.Background()), want) {
		t.Fatalf("consumer type = %v", state.Consumer.Type(context.Background()))
	}
}

func TestApplyAPIToStateTypesOmittedObjectMaps(t *testing.T) {
	api := &ProjectModel{
		Consumers: types.MapNull(types.ObjectType{AttrTypes: ProjectConsumersAttrTypes()}),
	}
	state := &ProjectModel{}

	applyProjectAPIToState(api, state, false, false)

	if !state.Consumers.IsNull() {
		t.Fatalf("consumers should remain null, got %#v", state.Consumers)
	}
	want := types.ObjectType{AttrTypes: ProjectConsumersAttrTypes()}
	if got := state.Consumers.ElementType(context.Background()); got == nil || !got.Equal(want) {
		t.Fatalf("consumers element type = %v, want %v", got, want)
	}
}
