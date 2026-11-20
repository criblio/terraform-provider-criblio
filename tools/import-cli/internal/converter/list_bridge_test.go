package converter

import (
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// The generator emits two shapes for the same array field: the generated model
// holds types.List, the resource model a slice of framework values. Neither is
// assignable or convertible to the other, so copyMatchingFields used to skip the
// field and leave it null, losing every element the API returned.
// criblio_guard_rule's required `rulesets` was dropped this way, which produced
// configurations that failed `terraform plan`.
func TestCopyMatchingFieldsCarriesListValuesIntoTypedSlices(t *testing.T) {
	source := provider.GuardRuleModel{
		ID:       types.StringValue("credit_card"),
		GroupID:  types.StringValue("default"),
		Regex:    types.StringValue("/x/"),
		Rulesets: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("pci-compliance"), types.StringValue("pii")}),
	}
	var target provider.GuardRuleResourceModel
	copyMatchingFields(reflect.ValueOf(&target).Elem(), reflect.ValueOf(source))

	require.Equal(t, "credit_card", target.ID.ValueString())
	require.Len(t, target.Rulesets, 2)
	require.Equal(t, "pci-compliance", target.Rulesets[0].ValueString())
	require.Equal(t, "pii", target.Rulesets[1].ValueString())
}

// A null list leaves the slice nil rather than producing an empty one, so a
// field the API did not report stays absent from the exported configuration.
func TestCopyMatchingFieldsLeavesNullListsAlone(t *testing.T) {
	source := provider.GuardRuleModel{
		ID:       types.StringValue("employee_badge_id"),
		Rulesets: types.ListNull(types.StringType),
	}
	var target provider.GuardRuleResourceModel
	copyMatchingFields(reflect.ValueOf(&target).Elem(), reflect.ValueOf(source))

	require.Equal(t, "employee_badge_id", target.ID.ValueString())
	require.Nil(t, target.Rulesets)
}
