// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputSplunkSearch struct {
	AuthHeaderExpr       types.String                   `tfsdk:"auth_header_expr"`
	AuthType             types.String                   `tfsdk:"auth_type"`
	BreakerRulesets      []types.String                 `tfsdk:"breaker_rulesets"`
	Connections          []InputSplunkSearchConnection  `tfsdk:"connections"`
	CredentialsSecret    types.String                   `tfsdk:"credentials_secret"`
	CronSchedule         types.String                   `tfsdk:"cron_schedule"`
	Description          types.String                   `tfsdk:"description"`
	Disabled             types.Bool                     `tfsdk:"disabled"`
	Earliest             types.String                   `tfsdk:"earliest"`
	Encoding             types.String                   `tfsdk:"encoding"`
	Endpoint             types.String                   `tfsdk:"endpoint"`
	EndpointHeaders      []EndpointHeader               `tfsdk:"endpoint_headers"`
	EndpointParams       []EndpointParam                `tfsdk:"endpoint_params"`
	Environment          types.String                   `tfsdk:"environment"`
	ID                   types.String                   `tfsdk:"id"`
	IgnoreGroupJobsLimit types.Bool                     `tfsdk:"ignore_group_jobs_limit"`
	JobTimeout           types.String                   `tfsdk:"job_timeout"`
	KeepAliveTime        types.Float64                  `tfsdk:"keep_alive_time"`
	Latest               types.String                   `tfsdk:"latest"`
	LoginURL             types.String                   `tfsdk:"login_url"`
	LogLevel             types.String                   `tfsdk:"log_level"`
	MaxMissedKeepAlives  types.Float64                  `tfsdk:"max_missed_keep_alives"`
	Metadata             []InputSplunkSearchMetadatum   `tfsdk:"metadata"`
	OauthHeaders         []InputSplunkSearchOauthHeader `tfsdk:"oauth_headers"`
	OauthParams          []InputSplunkSearchOauthParam  `tfsdk:"oauth_params"`
	OutputMode           types.String                   `tfsdk:"output_mode"`
	Password             types.String                   `tfsdk:"password"`
	Pipeline             types.String                   `tfsdk:"pipeline"`
	Pq                   *InputSplunkSearchPq           `tfsdk:"pq"`
	PqEnabled            types.Bool                     `tfsdk:"pq_enabled"`
	RejectUnauthorized   types.Bool                     `tfsdk:"reject_unauthorized"`
	RequestTimeout       types.Float64                  `tfsdk:"request_timeout"`
	RetryRules           *InputSplunkSearchRetryRules   `tfsdk:"retry_rules"`
	Search               types.String                   `tfsdk:"search"`
	SearchHead           types.String                   `tfsdk:"search_head"`
	Secret               types.String                   `tfsdk:"secret"`
	SecretParamName      types.String                   `tfsdk:"secret_param_name"`
	SendToRoutes         types.Bool                     `tfsdk:"send_to_routes"`
	StaleChannelFlushMs  types.Float64                  `tfsdk:"stale_channel_flush_ms"`
	Streamtags           []types.String                 `tfsdk:"streamtags"`
	TextSecret           types.String                   `tfsdk:"text_secret"`
	Token                types.String                   `tfsdk:"token"`
	TokenAttributeName   types.String                   `tfsdk:"token_attribute_name"`
	TokenTimeoutSecs     types.Float64                  `tfsdk:"token_timeout_secs"`
	TTL                  types.String                   `tfsdk:"ttl"`
	Type                 types.String                   `tfsdk:"type"`
	Username             types.String                   `tfsdk:"username"`
	UseRoundRobinDNS     types.Bool                     `tfsdk:"use_round_robin_dns"`
}
