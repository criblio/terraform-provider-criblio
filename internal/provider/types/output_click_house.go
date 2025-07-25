// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OutputClickHouse struct {
	AsyncInserts                  types.Bool                             `tfsdk:"async_inserts"`
	AuthHeaderExpr                types.String                           `tfsdk:"auth_header_expr"`
	AuthType                      types.String                           `tfsdk:"auth_type"`
	ColumnMappings                []OutputClickHouseColumnMapping        `tfsdk:"column_mappings"`
	Compress                      types.Bool                             `tfsdk:"compress"`
	Concurrency                   types.Float64                          `tfsdk:"concurrency"`
	CredentialsSecret             types.String                           `tfsdk:"credentials_secret"`
	Database                      types.String                           `tfsdk:"database"`
	DescribeTable                 types.String                           `tfsdk:"describe_table"`
	Description                   types.String                           `tfsdk:"description"`
	DumpFormatErrorsToDisk        types.Bool                             `tfsdk:"dump_format_errors_to_disk"`
	Environment                   types.String                           `tfsdk:"environment"`
	ExcludeMappingFields          []types.String                         `tfsdk:"exclude_mapping_fields"`
	ExtraHTTPHeaders              []OutputClickHouseExtraHTTPHeader      `tfsdk:"extra_http_headers"`
	FailedRequestLoggingMode      types.String                           `tfsdk:"failed_request_logging_mode"`
	FlushPeriodSec                types.Float64                          `tfsdk:"flush_period_sec"`
	Format                        types.String                           `tfsdk:"format"`
	ID                            types.String                           `tfsdk:"id"`
	LoginURL                      types.String                           `tfsdk:"login_url"`
	MappingType                   types.String                           `tfsdk:"mapping_type"`
	MaxPayloadEvents              types.Float64                          `tfsdk:"max_payload_events"`
	MaxPayloadSizeKB              types.Float64                          `tfsdk:"max_payload_size_kb"`
	OauthHeaders                  []OutputClickHouseOauthHeader          `tfsdk:"oauth_headers"`
	OauthParams                   []OutputClickHouseOauthParam           `tfsdk:"oauth_params"`
	OnBackpressure                types.String                           `tfsdk:"on_backpressure"`
	Password                      types.String                           `tfsdk:"password"`
	Pipeline                      types.String                           `tfsdk:"pipeline"`
	PqCompress                    types.String                           `tfsdk:"pq_compress"`
	PqControls                    *OutputClickHousePqControls            `tfsdk:"pq_controls"`
	PqMaxFileSize                 types.String                           `tfsdk:"pq_max_file_size"`
	PqMaxSize                     types.String                           `tfsdk:"pq_max_size"`
	PqMode                        types.String                           `tfsdk:"pq_mode"`
	PqOnBackpressure              types.String                           `tfsdk:"pq_on_backpressure"`
	PqPath                        types.String                           `tfsdk:"pq_path"`
	RejectUnauthorized            types.Bool                             `tfsdk:"reject_unauthorized"`
	ResponseHonorRetryAfterHeader types.Bool                             `tfsdk:"response_honor_retry_after_header"`
	ResponseRetrySettings         []OutputClickHouseResponseRetrySetting `tfsdk:"response_retry_settings"`
	SafeHeaders                   []types.String                         `tfsdk:"safe_headers"`
	Secret                        types.String                           `tfsdk:"secret"`
	SecretParamName               types.String                           `tfsdk:"secret_param_name"`
	SQLUsername                   types.String                           `tfsdk:"sql_username"`
	Streamtags                    []types.String                         `tfsdk:"streamtags"`
	SystemFields                  []types.String                         `tfsdk:"system_fields"`
	TableName                     types.String                           `tfsdk:"table_name"`
	TextSecret                    types.String                           `tfsdk:"text_secret"`
	TimeoutRetrySettings          *OutputClickHouseTimeoutRetrySettings  `tfsdk:"timeout_retry_settings"`
	TimeoutSec                    types.Float64                          `tfsdk:"timeout_sec"`
	TLS                           *OutputClickHouseTLSSettingsClientSide `tfsdk:"tls"`
	Token                         types.String                           `tfsdk:"token"`
	TokenAttributeName            types.String                           `tfsdk:"token_attribute_name"`
	TokenTimeoutSecs              types.Float64                          `tfsdk:"token_timeout_secs"`
	Type                          types.String                           `tfsdk:"type"`
	URL                           types.String                           `tfsdk:"url"`
	Username                      types.String                           `tfsdk:"username"`
	UseRoundRobinDNS              types.Bool                             `tfsdk:"use_round_robin_dns"`
	WaitForAsyncInserts           types.Bool                             `tfsdk:"wait_for_async_inserts"`
}
