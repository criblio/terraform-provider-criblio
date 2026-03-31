package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// searchDashboardConfigMapUseStateWhenNullOnly runs after UseStateForUnknown. The proposed plan
// can still contain a known map whose values are all TF-null or JSON null (e.g. "color": null
// placeholders) while prior state has no config after Read — that mismatch causes perpetual
// updates. When config is omitted in HCL, prefer prior state; otherwise strip null-like entries.
type searchDashboardConfigMapPlanModifier struct{}

func searchDashboardConfigMapUseStateWhenNullOnly() planmodifier.Map {
	return searchDashboardConfigMapPlanModifier{}
}

func (searchDashboardConfigMapPlanModifier) Description(_ context.Context) string {
	return "Aligns dashboard element config maps with prior state when the plan only contains null values and config was not set in Terraform."
}

func (m searchDashboardConfigMapPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m searchDashboardConfigMapPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}
	elems := req.PlanValue.Elements()
	if len(elems) == 0 {
		return
	}
	elemType := req.PlanValue.ElementType(ctx)
	kept := make(map[string]attr.Value, len(elems))
	for k, v := range elems {
		nv, ok := v.(jsontypes.Normalized)
		if !ok || !isJSONNullNormalized(nv) {
			kept[k] = v
		}
	}
	if len(kept) > 0 {
		newMap, diags := types.MapValueFrom(ctx, elemType, kept)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.PlanValue = newMap
		return
	}
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}
	// HCL set this map (e.g. config = { color = null }). The plan must match config; an empty map
	// is invalid and triggers "planned value does not match config value".
	resp.PlanValue = req.ConfigValue
}
