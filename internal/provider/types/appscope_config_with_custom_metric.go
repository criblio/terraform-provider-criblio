// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AppscopeConfigWithCustomMetric struct {
	Enable    types.Bool                            `tfsdk:"enable"`
	Format    AppscopeConfigWithCustomMetricFormat  `tfsdk:"format"`
	Transport AppscopeTransport                     `tfsdk:"transport"`
	Watch     []AppscopeConfigWithCustomMetricWatch `tfsdk:"watch"`
}
