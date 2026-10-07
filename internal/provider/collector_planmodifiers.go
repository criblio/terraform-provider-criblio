package provider

import (
	"context"

	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var collectorBlockNames = []string{
	"input_collector_splunk",
	"input_collector_rest",
	"input_collector_s3",
	"input_collector_azure_blob",
	"input_collector_cribl_lake",
	"input_collector_database",
	"input_collector_gcs",
	"input_collector_health_check",
	"input_collector_script",
	"input_collector_filesystem",
}

func collectorEnvironmentPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		collectorHoistedString{field: "environment"},
	}
}

func collectorIgnoreGroupJobsLimitPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		collectorHoistedBool{field: "ignore_group_jobs_limit"},
	}
}

func collectorResumeOnBootPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		collectorHoistedBool{field: "resume_on_boot"},
	}
}

func collectorTTLPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		collectorHoistedString{field: "ttl"},
	}
}

func collectorWorkerAffinityPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		collectorHoistedBool{field: "worker_affinity"},
	}
}

func collectorPreferConfigOrStatePlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{
		custom_objectplanmodifier.PreferConfigOrState(),
	}
}

// Collector fields are configurable at both the root and within a variant. The
// generic hoisting modifiers prefer prior state for read-only fields, which would
// freeze these configurable root mirrors when the nested configuration changes.
type collectorHoistedString struct {
	field string
}

func (m collectorHoistedString) Description(context.Context) string {
	return "Uses the configured collector value, preserving state only when omitted."
}

func (m collectorHoistedString) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m collectorHoistedString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	value, diags := collectorHoistedValue(ctx, req.Config, req.State, m.field, req.StateValue, types.StringUnknown())
	resp.Diagnostics.Append(diags...)
	resp.PlanValue = value
}

type collectorHoistedBool struct {
	field string
}

func (m collectorHoistedBool) Description(context.Context) string {
	return "Uses the configured collector value, preserving state only when omitted."
}

func (m collectorHoistedBool) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m collectorHoistedBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	value, diags := collectorHoistedValue(ctx, req.Config, req.State, m.field, req.StateValue, types.BoolUnknown())
	resp.Diagnostics.Append(diags...)
	resp.PlanValue = value
}

func collectorHoistedValue[T attr.Value](ctx context.Context, config tfsdk.Config, state tfsdk.State, field string, prior, unknown T) (T, diag.Diagnostics) {
	for _, blockName := range collectorBlockNames {
		root := path.Root(blockName)
		var block types.Object
		if diags := config.GetAttribute(ctx, root, &block); diags.HasError() {
			return unknown, diags
		}
		if block.IsUnknown() {
			return unknown, nil
		}
		if block.IsNull() {
			continue
		}
		value, ok := block.Attributes()[field].(T)
		if !ok {
			return unknown, diag.Diagnostics{diag.NewErrorDiagnostic("Unexpected collector field type", "Cannot read "+blockName+"."+field+". Please report this issue to the provider developers.")}
		}
		// Unknown expressions must stay unknown; replacing them with prior state
		// would promise Terraform an old value that may change during apply.
		if !value.IsNull() {
			return value, nil
		}
		if state.Raw.IsNull() || prior.IsNull() || prior.IsUnknown() {
			return unknown, nil
		}
		var priorBlock types.Object
		if diags := state.GetAttribute(ctx, root, &priorBlock); diags.HasError() {
			return unknown, diags
		}
		// Only preserve defaults belonging to the same collector variant.
		if !priorBlock.IsNull() && !priorBlock.IsUnknown() {
			return prior, nil
		}
		return unknown, nil
	}
	return unknown, nil
}
