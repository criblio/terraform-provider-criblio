package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// EmptyObjectAsNull returns a plan modifier that treats an API-owned empty
// object in state as equivalent to an omitted object in configuration.
func EmptyObjectAsNull() planmodifier.Object {
	return emptyObjectAsNull{}
}

type emptyObjectAsNull struct{}

func (emptyObjectAsNull) Description(_ context.Context) string {
	return "Treat an empty object in state as equivalent to null configuration."
}

func (emptyObjectAsNull) MarkdownDescription(_ context.Context) string {
	return "Treat an empty object in state as equivalent to null configuration."
}

func (emptyObjectAsNull) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.ConfigValue.IsNull() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	for _, value := range req.StateValue.Attributes() {
		if !value.IsNull() {
			return
		}
	}
	resp.PlanValue = req.StateValue
}
