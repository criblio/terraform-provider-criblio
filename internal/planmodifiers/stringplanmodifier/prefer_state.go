// PreferState returns a plan modifier that sets the planned value to the state value
// when state is known. Used to suppress drift when config and state differ (e.g. API
// returns different defaults). Unlike SuppressDiff, this runs even when plan is known.
package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreferState returns a plan modifier that prefers state over plan when state is known.
func PreferState() planmodifier.String {
	return preferState{}
}

type preferState struct{}

func (m preferState) Description(_ context.Context) string {
	return "When state has a value, use it to suppress config-vs-state drift."
}

func (m preferState) MarkdownDescription(_ context.Context) string {
	return "When state has a value, use it to suppress config-vs-state drift."
}

func (m preferState) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Only prefer state when it has a non-null value. When state is null/unknown, pass through
	// so config wins (avoids "planned value cty.NullVal does not match config value").
	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}
