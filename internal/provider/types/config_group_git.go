// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConfigGroupGit struct {
	Commit       types.String  `tfsdk:"commit"`
	LocalChanges types.Float64 `tfsdk:"local_changes"`
	Log          []Commit      `tfsdk:"log"`
}
