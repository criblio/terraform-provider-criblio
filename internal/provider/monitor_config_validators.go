package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithValidateConfig = &MonitorResource{}

// ValidateConfig preserves the upstream oneOf constraint after query objects are flattened.
func (r *MonitorResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var queries types.Map
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("query"), &queries)...)
	if resp.Diagnostics.HasError() || queries.IsNull() || queries.IsUnknown() {
		return
	}
	for label, value := range queries.Elements() {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		query := value.(types.Object).Attributes()
		mode, promql, refs := query["mode"], query["promql"], query["config_refs"]
		if mode.IsUnknown() || promql.IsUnknown() || refs.IsUnknown() {
			continue
		}
		inline := !mode.IsNull() && !promql.IsNull()
		inherited := !refs.IsNull()
		if inline == inherited {
			resp.Diagnostics.AddAttributeError(path.Root("query").AtMapKey(label), "Invalid monitor query",
				"Configure either mode and promql for an inline query, or config_refs for an inherited query; exactly one form is required.")
		}
	}
}
