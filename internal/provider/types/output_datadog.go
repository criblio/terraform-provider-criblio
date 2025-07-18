// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OutputDatadog struct {
	AllowAPIKeyFromEvents         types.Bool                          `tfsdk:"allow_api_key_from_events"`
	APIKey                        types.String                        `tfsdk:"api_key"`
	AuthType                      types.String                        `tfsdk:"auth_type"`
	BatchByTags                   types.Bool                          `tfsdk:"batch_by_tags"`
	Compress                      types.Bool                          `tfsdk:"compress"`
	Concurrency                   types.Float64                       `tfsdk:"concurrency"`
	ContentType                   types.String                        `tfsdk:"content_type"`
	CustomURL                     types.String                        `tfsdk:"custom_url"`
	Description                   types.String                        `tfsdk:"description"`
	Environment                   types.String                        `tfsdk:"environment"`
	ExtraHTTPHeaders              []OutputDatadogExtraHTTPHeader      `tfsdk:"extra_http_headers"`
	FailedRequestLoggingMode      types.String                        `tfsdk:"failed_request_logging_mode"`
	FlushPeriodSec                types.Float64                       `tfsdk:"flush_period_sec"`
	Host                          types.String                        `tfsdk:"host"`
	ID                            types.String                        `tfsdk:"id"`
	MaxPayloadEvents              types.Float64                       `tfsdk:"max_payload_events"`
	MaxPayloadSizeKB              types.Float64                       `tfsdk:"max_payload_size_kb"`
	Message                       types.String                        `tfsdk:"message"`
	OnBackpressure                types.String                        `tfsdk:"on_backpressure"`
	Pipeline                      types.String                        `tfsdk:"pipeline"`
	PqCompress                    types.String                        `tfsdk:"pq_compress"`
	PqControls                    *OutputDatadogPqControls            `tfsdk:"pq_controls"`
	PqMaxFileSize                 types.String                        `tfsdk:"pq_max_file_size"`
	PqMaxSize                     types.String                        `tfsdk:"pq_max_size"`
	PqMode                        types.String                        `tfsdk:"pq_mode"`
	PqOnBackpressure              types.String                        `tfsdk:"pq_on_backpressure"`
	PqPath                        types.String                        `tfsdk:"pq_path"`
	RejectUnauthorized            types.Bool                          `tfsdk:"reject_unauthorized"`
	ResponseHonorRetryAfterHeader types.Bool                          `tfsdk:"response_honor_retry_after_header"`
	ResponseRetrySettings         []OutputDatadogResponseRetrySetting `tfsdk:"response_retry_settings"`
	SafeHeaders                   []types.String                      `tfsdk:"safe_headers"`
	SendCountersAsCount           types.Bool                          `tfsdk:"send_counters_as_count"`
	Service                       types.String                        `tfsdk:"service"`
	Severity                      types.String                        `tfsdk:"severity"`
	Site                          types.String                        `tfsdk:"site"`
	Source                        types.String                        `tfsdk:"source"`
	Streamtags                    []types.String                      `tfsdk:"streamtags"`
	SystemFields                  []types.String                      `tfsdk:"system_fields"`
	Tags                          []types.String                      `tfsdk:"tags"`
	TextSecret                    types.String                        `tfsdk:"text_secret"`
	TimeoutRetrySettings          *OutputDatadogTimeoutRetrySettings  `tfsdk:"timeout_retry_settings"`
	TimeoutSec                    types.Float64                       `tfsdk:"timeout_sec"`
	TotalMemoryLimitKB            types.Float64                       `tfsdk:"total_memory_limit_kb"`
	Type                          types.String                        `tfsdk:"type"`
	UseRoundRobinDNS              types.Bool                          `tfsdk:"use_round_robin_dns"`
}
