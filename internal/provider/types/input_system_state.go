// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputSystemState struct {
	Collectors          *Collectors                  `tfsdk:"collectors"`
	Connections         []InputSystemStateConnection `tfsdk:"connections"`
	Description         types.String                 `tfsdk:"description"`
	Disabled            types.Bool                   `tfsdk:"disabled"`
	DisableNativeModule types.Bool                   `tfsdk:"disable_native_module"`
	Environment         types.String                 `tfsdk:"environment"`
	ID                  types.String                 `tfsdk:"id"`
	Interval            types.Float64                `tfsdk:"interval"`
	Metadata            []InputSystemStateMetadatum  `tfsdk:"metadata"`
	Persistence         *InputSystemStatePersistence `tfsdk:"persistence"`
	Pipeline            types.String                 `tfsdk:"pipeline"`
	Pq                  *InputSystemStatePq          `tfsdk:"pq"`
	PqEnabled           types.Bool                   `tfsdk:"pq_enabled"`
	SendToRoutes        types.Bool                   `tfsdk:"send_to_routes"`
	Streamtags          []types.String               `tfsdk:"streamtags"`
	Type                types.String                 `tfsdk:"type"`
}
