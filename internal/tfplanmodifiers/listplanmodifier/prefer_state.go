// PreferState returns a plan modifier that sets the planned value to the state value
// when state is known. Used to suppress drift when config and state differ.
package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreferState returns a plan modifier that prefers state over plan when state is known.
func PreferState() planmodifier.List {
	return preferState{}
}

type preferState struct{}

func (m preferState) Description(_ context.Context) string {
	return "When state has a value, use it to suppress config-vs-state drift."
}

func (m preferState) MarkdownDescription(_ context.Context) string {
	return "When state has a value, use it to suppress config-vs-state drift."
}

func (m preferState) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}
