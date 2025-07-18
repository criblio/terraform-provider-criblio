// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputTCP struct {
	AuthType            types.String                   `tfsdk:"auth_type"`
	BreakerRulesets     []types.String                 `tfsdk:"breaker_rulesets"`
	Connections         []InputTCPConnection           `tfsdk:"connections"`
	Description         types.String                   `tfsdk:"description"`
	Disabled            types.Bool                     `tfsdk:"disabled"`
	EnableHeader        types.Bool                     `tfsdk:"enable_header"`
	EnableProxyHeader   types.Bool                     `tfsdk:"enable_proxy_header"`
	Environment         types.String                   `tfsdk:"environment"`
	Host                types.String                   `tfsdk:"host"`
	ID                  types.String                   `tfsdk:"id"`
	IPWhitelistRegex    types.String                   `tfsdk:"ip_whitelist_regex"`
	MaxActiveCxn        types.Float64                  `tfsdk:"max_active_cxn"`
	Metadata            []InputTCPMetadatum            `tfsdk:"metadata"`
	Pipeline            types.String                   `tfsdk:"pipeline"`
	Port                types.Float64                  `tfsdk:"port"`
	Pq                  *InputTCPPq                    `tfsdk:"pq"`
	PqEnabled           types.Bool                     `tfsdk:"pq_enabled"`
	Preprocess          *InputTCPPreprocess            `tfsdk:"preprocess"`
	SendToRoutes        types.Bool                     `tfsdk:"send_to_routes"`
	SocketEndingMaxWait types.Float64                  `tfsdk:"socket_ending_max_wait"`
	SocketIdleTimeout   types.Float64                  `tfsdk:"socket_idle_timeout"`
	SocketMaxLifespan   types.Float64                  `tfsdk:"socket_max_lifespan"`
	StaleChannelFlushMs types.Float64                  `tfsdk:"stale_channel_flush_ms"`
	Streamtags          []types.String                 `tfsdk:"streamtags"`
	TLS                 *InputTCPTLSSettingsServerSide `tfsdk:"tls"`
	Type                types.String                   `tfsdk:"type"`
}
