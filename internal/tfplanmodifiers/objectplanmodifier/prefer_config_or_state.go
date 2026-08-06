package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// PreferConfigOrState returns a plan modifier that deeply merges config and state:
// - For each attribute, if config has a non-null value, use config
// - If config has null, use state value (preserving API-computed defaults)
// This allows user changes while preserving computed defaults for unspecified fields.
func PreferConfigOrState() planmodifier.Object {
	return preferConfigOrState{}
}

type preferConfigOrState struct{}

func (m preferConfigOrState) Description(_ context.Context) string {
	return "Deep merges config and state: uses config values when set, state values for unspecified fields."
}

func (m preferConfigOrState) MarkdownDescription(_ context.Context) string {
	return "Deep merges config and state: uses config values when set, state values for unspecified fields."
}

func (m preferConfigOrState) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// For new resources (no prior state), let normal flow proceed
	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}

	// If config is entirely null/unknown, use state (original PreferState behavior)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	// Deep merge recursively
	merged := deepMergeObjects(ctx, req.ConfigValue, req.StateValue)
	if merged != nil {
		resp.PlanValue = *merged
	}
}

// deepMergeObjects recursively merges config and state objects
func deepMergeObjects(ctx context.Context, config, state basetypes.ObjectValue) *basetypes.ObjectValue {
	configAttrs := config.Attributes()
	stateAttrs := state.Attributes()

	if len(configAttrs) == 0 || len(stateAttrs) == 0 {
		return &state
	}

	mergedAttrs := make(map[string]attr.Value)

	for key, stateVal := range stateAttrs {
		configVal, hasConfig := configAttrs[key]

		// Use state value if config doesn't have this key or config value is null/unknown
		if !hasConfig || configVal == nil || configVal.IsNull() || configVal.IsUnknown() {
			mergedAttrs[key] = stateVal
			continue
		}

		// Check if both are objects - if so, merge recursively
		configObj, configIsObj := configVal.(basetypes.ObjectValue)
		stateObj, stateIsObj := stateVal.(basetypes.ObjectValue)

		if configIsObj && stateIsObj && !configObj.IsNull() && !stateObj.IsNull() {
			// Recursively merge nested objects
			merged := deepMergeObjects(ctx, configObj, stateObj)
			if merged != nil {
				mergedAttrs[key] = *merged
			} else {
				mergedAttrs[key] = stateVal
			}
			continue
		}

		// Check if both are lists - handle empty list vs null
		configList, configIsList := configVal.(types.List)
		stateList, stateIsList := stateVal.(types.List)

		if configIsList && stateIsList {
			// If config list is empty and state list is also empty (or vice versa), use state
			// This handles the [] vs null equivalence
			if configList.IsNull() && len(stateList.Elements()) == 0 {
				mergedAttrs[key] = stateVal
				continue
			}
		}

		// Config has a real value - use it
		mergedAttrs[key] = configVal
	}

	// Include any config attrs that weren't in state (shouldn't happen normally)
	for key, configVal := range configAttrs {
		if _, exists := mergedAttrs[key]; !exists && configVal != nil {
			mergedAttrs[key] = configVal
		}
	}

	// Build the merged object using state's attribute types
	attrTypes := state.AttributeTypes(ctx)
	mergedObj, diags := basetypes.NewObjectValue(attrTypes, mergedAttrs)
	if diags.HasError() {
		return &state
	}

	return &mergedObj
}
