package provider

import (
	custom_boolplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/boolplanmodifier"
	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/objectplanmodifier"
	custom_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/stringplanmodifier"
	"github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
