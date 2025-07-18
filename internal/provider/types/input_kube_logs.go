// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputKubeLogs struct {
	BreakerRulesets     []types.String             `tfsdk:"breaker_rulesets"`
	Connections         []InputKubeLogsConnection  `tfsdk:"connections"`
	Description         types.String               `tfsdk:"description"`
	Disabled            types.Bool                 `tfsdk:"disabled"`
	EnableLoadBalancing types.Bool                 `tfsdk:"enable_load_balancing"`
	Environment         types.String               `tfsdk:"environment"`
	ID                  types.String               `tfsdk:"id"`
	Interval            types.Float64              `tfsdk:"interval"`
	Metadata            []InputKubeLogsMetadatum   `tfsdk:"metadata"`
	Persistence         *InputKubeLogsDiskSpooling `tfsdk:"persistence"`
	Pipeline            types.String               `tfsdk:"pipeline"`
	Pq                  *InputKubeLogsPq           `tfsdk:"pq"`
	PqEnabled           types.Bool                 `tfsdk:"pq_enabled"`
	Rules               []InputKubeLogsRule        `tfsdk:"rules"`
	SendToRoutes        types.Bool                 `tfsdk:"send_to_routes"`
	StaleChannelFlushMs types.Float64              `tfsdk:"stale_channel_flush_ms"`
	Streamtags          []types.String             `tfsdk:"streamtags"`
	Timestamps          types.Bool                 `tfsdk:"timestamps"`
	Type                types.String               `tfsdk:"type"`
}
