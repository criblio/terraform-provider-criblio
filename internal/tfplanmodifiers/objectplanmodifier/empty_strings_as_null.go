package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmptyStringsAsNull treats configured empty strings as equivalent to null
// values returned by an API for an existing resource.
func EmptyStringsAsNull() planmodifier.Object {
	return emptyStringsAsNull{}
}

type emptyStringsAsNull struct{}

func (emptyStringsAsNull) Description(_ context.Context) string {
	return "Treat configured empty strings as equivalent to API null values."
}

func (emptyStringsAsNull) MarkdownDescription(_ context.Context) string {
	return "Treat configured empty strings as equivalent to API null values."
}

func (emptyStringsAsNull) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.State.Raw.IsNull() || req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	configAttributes := req.ConfigValue.Attributes()
	stateAttributes := req.StateValue.Attributes()
	planAttributes := resp.PlanValue.Attributes()
	merged := make(map[string]attr.Value, len(planAttributes))
	for name, value := range planAttributes {
		merged[name] = value
	}

	changed := false
	for name, configValue := range configAttributes {
		configString, configIsString := configValue.(types.String)
		stateString, stateIsString := stateAttributes[name].(types.String)
		if !configIsString || !stateIsString || configString.IsNull() || configString.IsUnknown() || configString.ValueString() != "" || !stateString.IsNull() {
			continue
		}
		merged[name] = stateString
		changed = true
	}
	if !changed {
		return
	}

	value, diagnostics := types.ObjectValue(resp.PlanValue.AttributeTypes(ctx), merged)
	resp.Diagnostics.Append(diagnostics...)
	if !diagnostics.HasError() {
		resp.PlanValue = value
	}
}
