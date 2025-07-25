// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputGrafanaGrafana2 struct {
	ActivityLogSampleRate types.Float64                       `tfsdk:"activity_log_sample_rate"`
	CaptureHeaders        types.Bool                          `tfsdk:"capture_headers"`
	Connections           []InputGrafanaConnection2           `tfsdk:"connections"`
	Description           types.String                        `tfsdk:"description"`
	Disabled              types.Bool                          `tfsdk:"disabled"`
	EnableHealthCheck     types.Bool                          `tfsdk:"enable_health_check"`
	EnableProxyHeader     types.Bool                          `tfsdk:"enable_proxy_header"`
	Environment           types.String                        `tfsdk:"environment"`
	Host                  types.String                        `tfsdk:"host"`
	ID                    types.String                        `tfsdk:"id"`
	IPAllowlistRegex      types.String                        `tfsdk:"ip_allowlist_regex"`
	IPDenylistRegex       types.String                        `tfsdk:"ip_denylist_regex"`
	KeepAliveTimeout      types.Float64                       `tfsdk:"keep_alive_timeout"`
	LokiAPI               types.String                        `tfsdk:"loki_api"`
	LokiAuth              *InputGrafanaLokiAuth2              `tfsdk:"loki_auth"`
	MaxActiveReq          types.Float64                       `tfsdk:"max_active_req"`
	MaxRequestsPerSocket  types.Int64                         `tfsdk:"max_requests_per_socket"`
	Metadata              []InputGrafanaMetadatum2            `tfsdk:"metadata"`
	Pipeline              types.String                        `tfsdk:"pipeline"`
	Port                  types.Float64                       `tfsdk:"port"`
	Pq                    *InputGrafanaPq2                    `tfsdk:"pq"`
	PqEnabled             types.Bool                          `tfsdk:"pq_enabled"`
	PrometheusAPI         types.String                        `tfsdk:"prometheus_api"`
	PrometheusAuth        *InputGrafanaPrometheusAuth2        `tfsdk:"prometheus_auth"`
	RequestTimeout        types.Float64                       `tfsdk:"request_timeout"`
	SendToRoutes          types.Bool                          `tfsdk:"send_to_routes"`
	SocketTimeout         types.Float64                       `tfsdk:"socket_timeout"`
	Streamtags            []types.String                      `tfsdk:"streamtags"`
	TLS                   *InputGrafanaTLSSettingsServerSide2 `tfsdk:"tls"`
	Type                  types.String                        `tfsdk:"type"`
}
