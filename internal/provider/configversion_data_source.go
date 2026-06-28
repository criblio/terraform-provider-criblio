package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ConfigVersionDataSource{}
var _ datasource.DataSourceWithConfigure = &ConfigVersionDataSource{}

func NewConfigVersionDataSource() datasource.DataSource {
	return &ConfigVersionDataSource{}
}

type ConfigVersionDataSource struct {
	client *restclient.Client
}

type ConfigVersionDataSourceModel struct {
	ID    types.String   `tfsdk:"id"`
	Items []types.String `tfsdk:"items"`
}

func (d *ConfigVersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_version"
}

func (d *ConfigVersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ConfigVersion DataSource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: `Group ID`,
			},
			"items": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *ConfigVersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = clients.RC
}

func (d *ConfigVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model ConfigVersionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	items, err := restclient.Get[[]string](ctx, d.client, fmt.Sprintf("/master/groups/%s/configVersion", url.PathEscape(model.ID.ValueString())))
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	model.Items = nil
	if items != nil {
		model.Items = make([]types.String, 0, len(*items))
		for _, item := range *items {
			model.Items = append(model.Items, types.StringValue(item))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
