// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputKubeMetrics struct {
	Connections  []InputKubeMetricsConnection `tfsdk:"connections"`
	Description  types.String                 `tfsdk:"description"`
	Disabled     types.Bool                   `tfsdk:"disabled"`
	Environment  types.String                 `tfsdk:"environment"`
	ID           types.String                 `tfsdk:"id"`
	Interval     types.Float64                `tfsdk:"interval"`
	Metadata     []InputKubeMetricsMetadatum  `tfsdk:"metadata"`
	Persistence  *InputKubeMetricsPersistence `tfsdk:"persistence"`
	Pipeline     types.String                 `tfsdk:"pipeline"`
	Pq           *InputKubeMetricsPq          `tfsdk:"pq"`
	PqEnabled    types.Bool                   `tfsdk:"pq_enabled"`
	Rules        []InputKubeMetricsRule       `tfsdk:"rules"`
	SendToRoutes types.Bool                   `tfsdk:"send_to_routes"`
	Streamtags   []types.String               `tfsdk:"streamtags"`
	Type         types.String                 `tfsdk:"type"`
}
