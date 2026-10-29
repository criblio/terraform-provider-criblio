package provider

import (
	"context"

	custom_boolplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/boolplanmodifier"
	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/objectplanmodifier"
	custom_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/stringplanmodifier"
	"github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
}

func collectorEnvironmentPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		custom_stringplanmodifier.PreferState(),
		custom_stringplanmodifier.UseHoistedValue(collectorHoistedSources("environment")),
		stringplanmodifier.UseStateForUnknown(),
	}
}

func collectorIgnoreGroupJobsLimitPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		custom_boolplanmodifier.UseHoistedValue(collectorHoistedSources("ignore_group_jobs_limit")),
	}
}

func collectorResumeOnBootPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		custom_boolplanmodifier.UseHoistedValue(collectorHoistedSources("resume_on_boot")),
	}
}

func collectorTTLPlanModifiers() []planmodifier.String {
	return []planmodifier.String{
		custom_stringplanmodifier.PreferState(),
		custom_stringplanmodifier.UseHoistedValue(collectorHoistedSources("ttl")),
		stringplanmodifier.UseStateForUnknown(),
		collectorConfiguredTTL{},
	}
}

func collectorWorkerAffinityPlanModifiers() []planmodifier.Bool {
	return []planmodifier.Bool{
		custom_boolplanmodifier.UseHoistedValue(collectorHoistedSources("worker_affinity")),
	}
}

func collectorPreferConfigOrStatePlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{
		custom_objectplanmodifier.PreferConfigOrState(),
	}
}

func collectorHoistedSources(fieldName string) []utils.HoistedSource {
	sources := make([]utils.HoistedSource, 0, len(collectorBlockNames))
	for _, blockName := range collectorBlockNames {
		root := path.Root(blockName)
		sources = append(sources, utils.HoistedSource{
			AssociatedTypePath: root,
			FieldPath:          root.AtName(fieldName),
		})
	}
	return sources
}

// collectorConfiguredTTL overrides the legacy state preference only when TTL is
// explicitly configured. Omitted TTL retains the existing collector behavior.
type collectorConfiguredTTL struct{}

func (collectorConfiguredTTL) Description(context.Context) string {
	return "Uses an explicitly configured TTL instead of its previous state value."
}

func (m collectorConfiguredTTL) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (collectorConfiguredTTL) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	if !req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue
		return
	}
	// Include filesystem for TTL without changing the other legacy field hooks.
	for _, blockName := range append([]string{"input_collector_filesystem"}, collectorBlockNames...) {
		root := path.Root(blockName)
		var block types.Object
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, root, &block)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if block.IsNull() {
			continue
		}
		if block.IsUnknown() {
			resp.PlanValue = types.StringUnknown()
			return
		}
		var ttl types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, root.AtName("ttl"), &ttl)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !ttl.IsNull() {
			resp.PlanValue = ttl
		}
		return
	}
}
