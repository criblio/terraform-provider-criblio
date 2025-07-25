// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputKinesis struct {
	AssumeRoleArn           types.String             `tfsdk:"assume_role_arn"`
	AssumeRoleExternalID    types.String             `tfsdk:"assume_role_external_id"`
	AvoidDuplicates         types.Bool               `tfsdk:"avoid_duplicates"`
	AwsAPIKey               types.String             `tfsdk:"aws_api_key"`
	AwsAuthenticationMethod types.String             `tfsdk:"aws_authentication_method"`
	AwsSecret               types.String             `tfsdk:"aws_secret"`
	AwsSecretKey            types.String             `tfsdk:"aws_secret_key"`
	Connections             []InputKinesisConnection `tfsdk:"connections"`
	Description             types.String             `tfsdk:"description"`
	Disabled                types.Bool               `tfsdk:"disabled"`
	DurationSeconds         types.Float64            `tfsdk:"duration_seconds"`
	EnableAssumeRole        types.Bool               `tfsdk:"enable_assume_role"`
	Endpoint                types.String             `tfsdk:"endpoint"`
	Environment             types.String             `tfsdk:"environment"`
	GetRecordsLimit         types.Float64            `tfsdk:"get_records_limit"`
	GetRecordsLimitTotal    types.Float64            `tfsdk:"get_records_limit_total"`
	ID                      types.String             `tfsdk:"id"`
	LoadBalancingAlgorithm  types.String             `tfsdk:"load_balancing_algorithm"`
	Metadata                []InputKinesisMetadatum  `tfsdk:"metadata"`
	PayloadFormat           types.String             `tfsdk:"payload_format"`
	Pipeline                types.String             `tfsdk:"pipeline"`
	Pq                      *InputKinesisPq          `tfsdk:"pq"`
	PqEnabled               types.Bool               `tfsdk:"pq_enabled"`
	Region                  types.String             `tfsdk:"region"`
	RejectUnauthorized      types.Bool               `tfsdk:"reject_unauthorized"`
	ReuseConnections        types.Bool               `tfsdk:"reuse_connections"`
	SendToRoutes            types.Bool               `tfsdk:"send_to_routes"`
	ServiceInterval         types.Float64            `tfsdk:"service_interval"`
	ShardExpr               types.String             `tfsdk:"shard_expr"`
	ShardIteratorType       types.String             `tfsdk:"shard_iterator_type"`
	SignatureVersion        types.String             `tfsdk:"signature_version"`
	StreamName              types.String             `tfsdk:"stream_name"`
	Streamtags              []types.String           `tfsdk:"streamtags"`
	Type                    types.String             `tfsdk:"type"`
	VerifyKPLCheckSums      types.Bool               `tfsdk:"verify_kpl_check_sums"`
}
