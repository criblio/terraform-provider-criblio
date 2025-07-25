// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputConfluentCloud struct {
	AuthenticationTimeout     types.Float64                                         `tfsdk:"authentication_timeout"`
	AutoCommitInterval        types.Float64                                         `tfsdk:"auto_commit_interval"`
	AutoCommitThreshold       types.Float64                                         `tfsdk:"auto_commit_threshold"`
	BackoffRate               types.Float64                                         `tfsdk:"backoff_rate"`
	Brokers                   []types.String                                        `tfsdk:"brokers"`
	Connections               []InputConfluentCloudConnection                       `tfsdk:"connections"`
	ConnectionTimeout         types.Float64                                         `tfsdk:"connection_timeout"`
	Description               types.String                                          `tfsdk:"description"`
	Disabled                  types.Bool                                            `tfsdk:"disabled"`
	Environment               types.String                                          `tfsdk:"environment"`
	FromBeginning             types.Bool                                            `tfsdk:"from_beginning"`
	GroupID                   types.String                                          `tfsdk:"group_id"`
	HeartbeatInterval         types.Float64                                         `tfsdk:"heartbeat_interval"`
	ID                        types.String                                          `tfsdk:"id"`
	InitialBackoff            types.Float64                                         `tfsdk:"initial_backoff"`
	KafkaSchemaRegistry       *InputConfluentCloudKafkaSchemaRegistryAuthentication `tfsdk:"kafka_schema_registry"`
	MaxBackOff                types.Float64                                         `tfsdk:"max_back_off"`
	MaxBytes                  types.Float64                                         `tfsdk:"max_bytes"`
	MaxBytesPerPartition      types.Float64                                         `tfsdk:"max_bytes_per_partition"`
	MaxRetries                types.Float64                                         `tfsdk:"max_retries"`
	MaxSocketErrors           types.Float64                                         `tfsdk:"max_socket_errors"`
	Metadata                  []InputConfluentCloudMetadatum                        `tfsdk:"metadata"`
	Pipeline                  types.String                                          `tfsdk:"pipeline"`
	Pq                        *InputConfluentCloudPq                                `tfsdk:"pq"`
	PqEnabled                 types.Bool                                            `tfsdk:"pq_enabled"`
	ReauthenticationThreshold types.Float64                                         `tfsdk:"reauthentication_threshold"`
	RebalanceTimeout          types.Float64                                         `tfsdk:"rebalance_timeout"`
	RequestTimeout            types.Float64                                         `tfsdk:"request_timeout"`
	Sasl                      *InputConfluentCloudAuthentication                    `tfsdk:"sasl"`
	SendToRoutes              types.Bool                                            `tfsdk:"send_to_routes"`
	SessionTimeout            types.Float64                                         `tfsdk:"session_timeout"`
	Streamtags                []types.String                                        `tfsdk:"streamtags"`
	TLS                       *InputConfluentCloudTLSSettingsClientSide             `tfsdk:"tls"`
	Topics                    []types.String                                        `tfsdk:"topics"`
	Type                      types.String                                          `tfsdk:"type"`
}
