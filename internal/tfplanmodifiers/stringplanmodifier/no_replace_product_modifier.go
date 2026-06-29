package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = StringNoReplaceProductModifierPlanModifier{}

type StringNoReplaceProductModifierPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringNoReplaceProductModifierPlanModifier) Description(_ context.Context) string {
	return "Allows product changes without forcing resource replacement for default groups (default, defaultHybrid, default_fleet)."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringNoReplaceProductModifierPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Allows product changes without forcing resource replacement for default groups (default, defaultHybrid, default_fleet)."
}

// PlanModifyString performs the plan modification.
// This modifier prevents resource replacement when:
// - product is "stream" AND group id is "default" or "defaultHybrid"
// - product is "edge" AND group id is "default_fleet"
// For other cases, it allows normal behavior (which may still trigger replacement from other modifiers).
func (v StringNoReplaceProductModifierPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Get the group id from state or plan
	var groupID types.String

	// Try to get id from state first (for existing resources)
	stateDiags := req.State.GetAttribute(ctx, path.Root("id"), &groupID)
	if stateDiags.HasError() || groupID.IsNull() || groupID.IsUnknown() {
		// If not in state, try plan (for new resources)
		planDiags := req.Plan.GetAttribute(ctx, path.Root("id"), &groupID)
		if planDiags.HasError() || groupID.IsNull() || groupID.IsUnknown() {
			// If we can't get the id, allow normal behavior
			return
		}
	}

	groupIDValue := groupID.ValueString()

	// Get the product value from the plan/config (the new value)
	var productValue string
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		productValue = req.ConfigValue.ValueString()
	} else if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		productValue = req.PlanValue.ValueString()
	} else {
		// If product is unknown, allow normal behavior
		return
	}

	// Check if the value actually changed
	stateValueChanged := false
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		stateProductValue := req.StateValue.ValueString()
		if stateProductValue != productValue {
			stateValueChanged = true
		}
	} else {
		// If there's no state value, this is a new resource, so no change to prevent
		return
	}

	// Only prevent replacement if the value actually changed
	if !stateValueChanged {
		return
	}

	// Check if we should prevent replacement based on the conditions
	shouldPreventReplacement := false

	switch productValue {
	case "stream":
		// For stream product, allow changes if group is default or defaultHybrid
		if groupIDValue == "default" || groupIDValue == "defaultHybrid" {
			shouldPreventReplacement = true
		}
	case "edge":
		// For edge product, allow changes if group is default_fleet
		if groupIDValue == "default_fleet" {
			shouldPreventReplacement = true
		}
	}

	// If we should prevent replacement, explicitly set RequiresReplace to false
	// This ensures that replacement is not forced for these specific cases
	if shouldPreventReplacement {
		// Explicitly prevent replacement by not setting RequiresReplace
		// The default behavior is to not require replacement
		resp.RequiresReplace = false
		return
	}

	// For other cases, force replacement (mimic RequiresReplaceIfConfigured behavior)
	// This ensures that non-default groups still require replacement when product changes
	resp.RequiresReplace = true
}

func NoReplaceProductModifier() planmodifier.String {
	return StringNoReplaceProductModifierPlanModifier{}
}
