package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PackDataSource{}
var _ datasource.DataSourceWithConfigure = &PackDataSource{}

func NewPackDataSource() datasource.DataSource {
	return &PackDataSource{}
}

// PackDataSource is the data source implementation.
type PackDataSource struct {
	client *restclient.Client
}

// PackDataSourceModel describes the data model.
type PackDataSourceModel struct {
	Disabled types.Bool                `queryParam:"style=form,explode=true,name=disabled" tfsdk:"disabled"`
	GroupID  types.String              `tfsdk:"group_id"`
	ID       types.String              `tfsdk:"id"`
	Items    []tfTypes.PackInstallInfo `tfsdk:"items"`
	With     types.String              `queryParam:"style=form,explode=true,name=with" tfsdk:"with"`
}

// Metadata returns the data source type name.
func (r *PackDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pack"
}

// Schema defines the schema for the data source.
func (r *PackDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pack DataSource",

		Attributes: map[string]schema.Attribute{
			"disabled": schema.BoolAttribute{
				Optional: true,
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: `The consumer group to which this instance belongs. Defaults to 'Cribl'.`,
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: `Pack name`,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"author": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"exports": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"inputs": schema.Float64Attribute{
							Computed: true,
						},
						"min_log_stream_version": schema.StringAttribute{
							Computed: true,
						},
						"outputs": schema.Float64Attribute{
							Computed: true,
						},
						"settings": schema.MapAttribute{
							Computed:    true,
							ElementType: jsontypes.NormalizedType{},
						},
						"source": schema.StringAttribute{
							Computed: true,
						},
						"spec": schema.StringAttribute{
							Computed: true,
						},
						"tags": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"data_type": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"domain": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"streamtags": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"technology": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"version": schema.StringAttribute{
							Computed: true,
						},
						"warnings": schema.StringAttribute{
							CustomType:  jsontypes.NormalizedType{},
							Computed:    true,
							Description: `Parsed as JSON.`,
						},
					},
				},
			},
			"with": schema.StringAttribute{
				Optional:    true,
				Description: `Comma separated list of entities, "outputs", "inputs"`,
			},
		},
	}
}

func (r *PackDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(*ProviderClients)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *ProviderClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = clients.RC
}

func (r *PackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PackDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	packID := resolvePackIDForRestAPI(ctx, r.client, data.GroupID.ValueString(), data.ID.ValueString())
	path := fmt.Sprintf("/m/%s/packs/%s", url.PathEscape(data.GroupID.ValueString()), url.PathEscape(packID))
	query := url.Values{}
	if !data.With.IsNull() && !data.With.IsUnknown() {
		query.Set("with", data.With.ValueString())
	}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		query.Set("disabled", strings.ToLower(fmt.Sprintf("%t", data.Disabled.ValueBool())))
	}
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	apiModel, err := restclient.Get[packAPIModel](ctx, r.client, path)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	if apiModel == nil {
		resp.Diagnostics.AddError("unexpected response from API", "empty response body")
		return
	}
	data.Items = []tfTypes.PackInstallInfo{packInstallInfoFromAPI(apiModel)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
