// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Cloud struct {
	Provider types.String `tfsdk:"provider"`
	Region   types.String `tfsdk:"region"`
}
