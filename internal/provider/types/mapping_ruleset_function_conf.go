package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MappingRulesetFunctionConf struct {
	Description types.String                        `tfsdk:"description"`
	Disabled    types.Bool                          `tfsdk:"disabled"`
	Filter      types.String                        `tfsdk:"filter"`
	Final       types.Bool                          `tfsdk:"final"`
	GroupID     types.String                        `tfsdk:"group_id"`
	ID          types.String                        `tfsdk:"id"`
	Conf        *MappingRulesetFunctionSpecificConf `tfsdk:"conf"`
}

type MappingRulesetFunctionSpecificConf struct {
	Add []MappingRulesetAddField `tfsdk:"add"`
}

type MappingRulesetAddField struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
