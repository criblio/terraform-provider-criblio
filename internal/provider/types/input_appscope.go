// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputAppscope struct {
	AuthToken           types.String                        `tfsdk:"auth_token"`
	AuthType            types.String                        `tfsdk:"auth_type"`
	BreakerRulesets     []types.String                      `tfsdk:"breaker_rulesets"`
	Connections         []InputAppscopeConnection           `tfsdk:"connections"`
	Description         types.String                        `tfsdk:"description"`
	Disabled            types.Bool                          `tfsdk:"disabled"`
	EnableProxyHeader   types.Bool                          `tfsdk:"enable_proxy_header"`
	EnableUnixPath      types.Bool                          `tfsdk:"enable_unix_path"`
	Environment         types.String                        `tfsdk:"environment"`
	Filter              *InputAppscopeFilter                `tfsdk:"filter"`
	Host                types.String                        `tfsdk:"host"`
	ID                  types.String                        `tfsdk:"id"`
	IPWhitelistRegex    types.String                        `tfsdk:"ip_whitelist_regex"`
	MaxActiveCxn        types.Float64                       `tfsdk:"max_active_cxn"`
	Metadata            []InputAppscopeMetadatum            `tfsdk:"metadata"`
	Persistence         *InputAppscopePersistence           `tfsdk:"persistence"`
	Pipeline            types.String                        `tfsdk:"pipeline"`
	Port                types.Float64                       `tfsdk:"port"`
	Pq                  *InputAppscopePq                    `tfsdk:"pq"`
	PqEnabled           types.Bool                          `tfsdk:"pq_enabled"`
	SendToRoutes        types.Bool                          `tfsdk:"send_to_routes"`
	SocketEndingMaxWait types.Float64                       `tfsdk:"socket_ending_max_wait"`
	SocketIdleTimeout   types.Float64                       `tfsdk:"socket_idle_timeout"`
	SocketMaxLifespan   types.Float64                       `tfsdk:"socket_max_lifespan"`
	StaleChannelFlushMs types.Float64                       `tfsdk:"stale_channel_flush_ms"`
	Streamtags          []types.String                      `tfsdk:"streamtags"`
	TextSecret          types.String                        `tfsdk:"text_secret"`
	TLS                 *InputAppscopeTLSSettingsServerSide `tfsdk:"tls"`
	Type                types.String                        `tfsdk:"type"`
	UnixSocketPath      types.String                        `tfsdk:"unix_socket_path"`
	UnixSocketPerms     types.String                        `tfsdk:"unix_socket_perms"`
}
