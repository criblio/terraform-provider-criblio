// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TLSSettingsClientSideAzureEventhub struct {
	Disabled           types.Bool `tfsdk:"disabled"`
	RejectUnauthorized types.Bool `tfsdk:"reject_unauthorized"`
}
