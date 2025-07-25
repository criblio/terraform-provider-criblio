// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LookupFileDataSource{}
var _ datasource.DataSourceWithConfigure = &LookupFileDataSource{}

func NewLookupFileDataSource() datasource.DataSource {
	return &LookupFileDataSource{}
}

// LookupFileDataSource is the data source implementation.
type LookupFileDataSource struct {
	// Provider configured SDK client.
	client *sdk.CriblIo
}

// LookupFileDataSourceModel describes the data model.
type LookupFileDataSourceModel struct {
	GroupID types.String              `tfsdk:"group_id"`
	Items   []tfTypes.LookupFileUnion `tfsdk:"items"`
}

// Metadata returns the data source type name.
func (r *LookupFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lookup_file"
}

// Schema defines the schema for the data source.
func (r *LookupFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "LookupFile DataSource",

		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: `The consumer group to which this instance belongs. Defaults to 'Cribl'.`,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"lookup_file1": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed: true,
								},
								"file_info": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"filename": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								"id": schema.StringAttribute{
									Computed: true,
								},
								"mode": schema.StringAttribute{
									Computed: true,
								},
								"pending_task": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"error": schema.StringAttribute{
											Computed:    true,
											Description: `Error message if task has failed`,
										},
										"id": schema.StringAttribute{
											Computed:    true,
											Description: `Task ID (generated).`,
										},
										"type": schema.StringAttribute{
											Computed:    true,
											Description: `Task type`,
										},
									},
								},
								"size": schema.Float64Attribute{
									Computed:    true,
									Description: `File size. Optional.`,
								},
								"tags": schema.StringAttribute{
									Computed:    true,
									Description: `One or more tags related to this lookup. Optional.`,
								},
								"version": schema.StringAttribute{
									Computed:    true,
									Description: `Unique string generated for each modification of this lookup`,
								},
							},
						},
						"lookup_file2": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"content": schema.StringAttribute{
									Computed:    true,
									Description: `File content.`,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"id": schema.StringAttribute{
									Computed: true,
								},
								"mode": schema.StringAttribute{
									Computed: true,
								},
								"pending_task": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"error": schema.StringAttribute{
											Computed:    true,
											Description: `Error message if task has failed`,
										},
										"id": schema.StringAttribute{
											Computed:    true,
											Description: `Task ID (generated).`,
										},
										"type": schema.StringAttribute{
											Computed:    true,
											Description: `Task type`,
										},
									},
								},
								"size": schema.Float64Attribute{
									Computed:    true,
									Description: `File size. Optional.`,
								},
								"tags": schema.StringAttribute{
									Computed:    true,
									Description: `One or more tags related to this lookup. Optional.`,
								},
								"version": schema.StringAttribute{
									Computed:    true,
									Description: `Unique string generated for each modification of this lookup`,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *LookupFileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *LookupFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *LookupFileDataSourceModel
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

	request, requestDiags := data.ToOperationsListLookupFileRequest(ctx)
	resp.Diagnostics.Append(requestDiags...)

	if resp.Diagnostics.HasError() {
		return
	}
	res, err := r.client.Lookups.ListLookupFile(ctx, *request)
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
	if !(res.Object != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	resp.Diagnostics.Append(data.RefreshFromOperationsListLookupFileResponseBody(ctx, res.Object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
