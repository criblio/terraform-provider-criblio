// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OutputTcpjson struct {
	AuthToken                 types.String                        `tfsdk:"auth_token"`
	AuthType                  types.String                        `tfsdk:"auth_type"`
	Compression               types.String                        `tfsdk:"compression"`
	ConnectionTimeout         types.Float64                       `tfsdk:"connection_timeout"`
	Description               types.String                        `tfsdk:"description"`
	DNSResolvePeriodSec       types.Float64                       `tfsdk:"dns_resolve_period_sec"`
	Environment               types.String                        `tfsdk:"environment"`
	ExcludeSelf               types.Bool                          `tfsdk:"exclude_self"`
	Host                      types.String                        `tfsdk:"host"`
	Hosts                     []OutputTcpjsonHost                 `tfsdk:"hosts"`
	ID                        types.String                        `tfsdk:"id"`
	LoadBalanced              types.Bool                          `tfsdk:"load_balanced"`
	LoadBalanceStatsPeriodSec types.Float64                       `tfsdk:"load_balance_stats_period_sec"`
	LogFailedRequests         types.Bool                          `tfsdk:"log_failed_requests"`
	MaxConcurrentSenders      types.Float64                       `tfsdk:"max_concurrent_senders"`
	OnBackpressure            types.String                        `tfsdk:"on_backpressure"`
	Pipeline                  types.String                        `tfsdk:"pipeline"`
	Port                      types.Float64                       `tfsdk:"port"`
	PqCompress                types.String                        `tfsdk:"pq_compress"`
	PqControls                *OutputTcpjsonPqControls            `tfsdk:"pq_controls"`
	PqMaxFileSize             types.String                        `tfsdk:"pq_max_file_size"`
	PqMaxSize                 types.String                        `tfsdk:"pq_max_size"`
	PqMode                    types.String                        `tfsdk:"pq_mode"`
	PqOnBackpressure          types.String                        `tfsdk:"pq_on_backpressure"`
	PqPath                    types.String                        `tfsdk:"pq_path"`
	SendHeader                types.Bool                          `tfsdk:"send_header"`
	Streamtags                []types.String                      `tfsdk:"streamtags"`
	SystemFields              []types.String                      `tfsdk:"system_fields"`
	TextSecret                types.String                        `tfsdk:"text_secret"`
	ThrottleRatePerSec        types.String                        `tfsdk:"throttle_rate_per_sec"`
	TLS                       *OutputTcpjsonTLSSettingsClientSide `tfsdk:"tls"`
	TokenTTLMinutes           types.Float64                       `tfsdk:"token_ttl_minutes"`
	Type                      types.String                        `tfsdk:"type"`
	WriteTimeout              types.Float64                       `tfsdk:"write_timeout"`
}
