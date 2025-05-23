// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputSystemMetrics struct {
	Connections  []ConnectionSystemMetrics `tfsdk:"connections"`
	Container    *InputContainer           `tfsdk:"container"`
	Description  types.String              `tfsdk:"description"`
	Disabled     types.Bool                `tfsdk:"disabled"`
	Environment  types.String              `tfsdk:"environment"`
	Host         *HostSystemMetrics        `tfsdk:"host"`
	ID           types.String              `tfsdk:"id"`
	Interval     types.Float64             `tfsdk:"interval"`
	Metadata     []MetadatumSystemMetrics  `tfsdk:"metadata"`
	Persistence  *PersistenceSystemMetrics `tfsdk:"persistence"`
	Pipeline     types.String              `tfsdk:"pipeline"`
	Pq           *PqSystemMetrics          `tfsdk:"pq"`
	PqEnabled    types.Bool                `tfsdk:"pq_enabled"`
	Process      *ProcessSystemMetrics     `tfsdk:"process"`
	SendToRoutes types.Bool                `tfsdk:"send_to_routes"`
	Status       *TFStatus                 `tfsdk:"status"`
	Streamtags   []types.String            `tfsdk:"streamtags"`
	Type         types.String              `tfsdk:"type"`
}
