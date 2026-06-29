package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bannerMessage struct {
	Created         types.Float64  `tfsdk:"created"`
	CustomThemes    []types.String `tfsdk:"custom_themes"`
	Enabled         types.Bool     `tfsdk:"enabled"`
	ID              types.String   `tfsdk:"id"`
	InvertFontColor types.Bool     `tfsdk:"invert_font_color"`
	Link            types.String   `tfsdk:"link"`
	LinkDisplay     types.String   `tfsdk:"link_display"`
	Message         types.String   `tfsdk:"message"`
	Theme           types.String   `tfsdk:"theme"`
	Type            types.String   `tfsdk:"type"`
}

type configGroupCloud struct {
	Provider types.String `tfsdk:"provider"`
	Region   types.String `tfsdk:"region"`
}

type packInstallInfo struct {
	Author              types.String                    `tfsdk:"author"`
	Description         types.String                    `tfsdk:"description"`
	DisplayName         types.String                    `tfsdk:"display_name"`
	Exports             []types.String                  `tfsdk:"exports"`
	ID                  types.String                    `tfsdk:"id"`
	Inputs              types.Float64                   `tfsdk:"inputs"`
	MinLogStreamVersion types.String                    `tfsdk:"min_log_stream_version"`
	Outputs             types.Float64                   `tfsdk:"outputs"`
	Settings            map[string]jsontypes.Normalized `tfsdk:"settings"`
	Source              types.String                    `tfsdk:"source"`
	Spec                types.String                    `tfsdk:"spec"`
	Tags                *packInstallInfoTags            `tfsdk:"tags"`
	Version             types.String                    `tfsdk:"version"`
	Warnings            jsontypes.Normalized            `tfsdk:"warnings"`
}

type packInstallInfoTags struct {
	DataType   []types.String `tfsdk:"data_type"`
	Domain     []types.String `tfsdk:"domain"`
	Streamtags []types.String `tfsdk:"streamtags"`
	Technology []types.String `tfsdk:"technology"`
}

type packRequestBodyTags struct {
	DataType   []types.String `tfsdk:"data_type"`
	Domain     []types.String `tfsdk:"domain"`
	Streamtags []types.String `tfsdk:"streamtags"`
	Technology []types.String `tfsdk:"technology"`
}
