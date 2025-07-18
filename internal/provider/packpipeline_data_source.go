// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tfTypes "github.com/speakeasy/terraform-provider-criblio/internal/provider/types"
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PackPipelineDataSource{}
var _ datasource.DataSourceWithConfigure = &PackPipelineDataSource{}

func NewPackPipelineDataSource() datasource.DataSource {
	return &PackPipelineDataSource{}
}

// PackPipelineDataSource is the data source implementation.
type PackPipelineDataSource struct {
	// Provider configured SDK client.
	client *sdk.CriblIo
}

// PackPipelineDataSourceModel describes the data model.
type PackPipelineDataSourceModel struct {
	Conf    tfTypes.PipelineConf `tfsdk:"conf"`
	GroupID types.String         `tfsdk:"group_id"`
	ID      types.String         `tfsdk:"id"`
	Pack    types.String         `tfsdk:"pack"`
}

// Metadata returns the data source type name.
func (r *PackPipelineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pack_pipeline"
}

// Schema defines the schema for the data source.
func (r *PackPipelineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PackPipeline DataSource",

		Attributes: map[string]schema.Attribute{
			"conf": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"async_func_timeout": schema.Int64Attribute{
						Computed:    true,
						Description: `Time (in ms) to wait for an async function to complete processing of a data item`,
					},
					"description": schema.StringAttribute{
						Computed: true,
					},
					"functions": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"conf": schema.SingleNestedAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed:    true,
									Description: `Simple description of this step`,
								},
								"disabled": schema.BoolAttribute{
									Computed:    true,
									Description: `If true, data will not be pushed through this function`,
								},
								"filter": schema.StringAttribute{
									Computed:    true,
									Description: `Filter that selects data to be fed through this Function`,
								},
								"final": schema.BoolAttribute{
									Computed:    true,
									Description: `If enabled, stops the results of this Function from being passed to the downstream Functions`,
								},
								"group_id": schema.StringAttribute{
									Computed:    true,
									Description: `Group ID`,
								},
								"id": schema.StringAttribute{
									Computed:    true,
									Description: `Function ID`,
								},
							},
						},
						Description: `List of Functions to pass data through`,
					},
					"groups": schema.MapNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed:    true,
									Description: `Short description of this group`,
								},
								"disabled": schema.BoolAttribute{
									Computed:    true,
									Description: `Whether this group is disabled`,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"output": schema.StringAttribute{
						Computed:    true,
						Description: `The output destination for events processed by this Pipeline`,
					},
					"streamtags": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: `Tags for filtering and grouping in @{product}`,
					},
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: `The consumer group to which this instance belongs. Defaults to 'Cribl'.`,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pack": schema.StringAttribute{
				Required:    true,
				Description: `pack ID to GET`,
			},
		},
	}
}

func (r *PackPipelineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.CriblIo)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *sdk.CriblIo, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PackPipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PackPipelineDataSourceModel
	var item types.Object

	resp.Diagnostics.Append(req.Config.Get(ctx, &item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(item.As(ctx, &data, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)

	if resp.Diagnostics.HasError() {
		return
	}

	request, requestDiags := data.ToOperationsGetPipelineByPackRequest(ctx)
	resp.Diagnostics.Append(requestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}
	res, err := r.client.Pipelines.GetPipelineByPack(ctx, *request)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		if res != nil && res.RawResponse != nil {
			resp.Diagnostics.AddError("unexpected http request/response", debugResponse(res.RawResponse))
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("unexpected response from API", fmt.Sprintf("%v", res))
		return
	}
	if res.StatusCode != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode), debugResponse(res.RawResponse))
		return
	}
	if !(res.Object != nil && res.Object.Items != nil && len(res.Object.Items) > 0) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	resp.Diagnostics.Append(data.RefreshFromSharedPipeline(ctx, &res.Object.Items[0])...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
