// PreferStateIfJSONEqual returns a plan modifier that uses state value when plan and
// state are semantically equal JSON. Suppresses diff from JSON formatting differences.
package stringplanmodifier

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreferStateIfJSONEqual returns a plan modifier that prefers state when plan and state
// represent the same JSON (ignoring formatting).
func PreferStateIfJSONEqual() planmodifier.String {
	return preferStateIfJSONEqual{}
}

type preferStateIfJSONEqual struct{}

func (m preferStateIfJSONEqual) Description(_ context.Context) string {
	return "Use state value when plan and state are semantically equal JSON."
}

func (m preferStateIfJSONEqual) MarkdownDescription(_ context.Context) string {
	return "Use state value when plan and state are semantically equal JSON."
}

func (m preferStateIfJSONEqual) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}
	planStr := req.PlanValue.ValueString()
	stateStr := req.StateValue.ValueString()
	if jsonEqual(planStr, stateStr) {
		resp.PlanValue = req.StateValue
	}
}

func jsonEqual(a, b string) bool {
	var va, vb interface{}
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		return false
	}
	// Re-marshal both to canonical form for comparison (handles key order, whitespace)
	ca, err := json.Marshal(va)
	if err != nil {
		return false
	}
	cb, err := json.Marshal(vb)
	if err != nil {
		return false
	}
	return string(ca) == string(cb)
}
