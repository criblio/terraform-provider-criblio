// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AuthTokensExtHTTPRaw struct {
	Description types.String                    `tfsdk:"description"`
	Metadata    []AuthTokensExtMetadatumHTTPRaw `tfsdk:"metadata"`
	Token       types.String                    `tfsdk:"token"`
}
