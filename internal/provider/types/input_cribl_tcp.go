// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputCriblTCP struct {
	Connections         []ConnectionCriblTCP           `tfsdk:"connections"`
	Description         types.String                   `tfsdk:"description"`
	Disabled            types.Bool                     `tfsdk:"disabled"`
	EnableLoadBalancing types.Bool                     `tfsdk:"enable_load_balancing"`
	EnableProxyHeader   types.Bool                     `tfsdk:"enable_proxy_header"`
	Environment         types.String                   `tfsdk:"environment"`
	Host                types.String                   `tfsdk:"host"`
	ID                  types.String                   `tfsdk:"id"`
	MaxActiveCxn        types.Float64                  `tfsdk:"max_active_cxn"`
	Metadata            []MetadatumCriblTCP            `tfsdk:"metadata"`
	Pipeline            types.String                   `tfsdk:"pipeline"`
	Port                types.Float64                  `tfsdk:"port"`
	Pq                  *PqCriblTCP                    `tfsdk:"pq"`
	PqEnabled           types.Bool                     `tfsdk:"pq_enabled"`
	SendToRoutes        types.Bool                     `tfsdk:"send_to_routes"`
	SocketEndingMaxWait types.Float64                  `tfsdk:"socket_ending_max_wait"`
	SocketIdleTimeout   types.Float64                  `tfsdk:"socket_idle_timeout"`
	SocketMaxLifespan   types.Float64                  `tfsdk:"socket_max_lifespan"`
	Status              *TFStatus                      `tfsdk:"status"`
	Streamtags          []types.String                 `tfsdk:"streamtags"`
	TLS                 *TLSSettingsServerSideCriblTCP `tfsdk:"tls"`
	Type                types.String                   `tfsdk:"type"`
}
