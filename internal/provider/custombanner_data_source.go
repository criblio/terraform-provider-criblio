package provider

import (
	"context"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CustomBannerDataSource{}
var _ datasource.DataSourceWithConfigure = &CustomBannerDataSource{}

func NewCustomBannerDataSource() datasource.DataSource {
	return &CustomBannerDataSource{}
}

type CustomBannerDataSource struct {
	client *restclient.Client
}

type CustomBannerDataSourceModel struct {
	Items []bannerMessage `tfsdk:"items"`
}

func (r *CustomBannerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_banner"
}

func (r *CustomBannerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "CustomBanner DataSource",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created": schema.Float64Attribute{
							Computed:    true,
							Description: `Time created (Unix epoch seconds)`,
						},
						"custom_themes": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: `Show a banner on top of all pages`,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"invert_font_color": schema.BoolAttribute{
							Computed: true,
						},
						"link": schema.StringAttribute{
							Computed:    true,
							Description: `Optionally, provide a URL to append to the message`,
						},
						"link_display": schema.StringAttribute{
							Computed:    true,
							Description: `Optionally, display your link with a short text label instead of the raw URL (100-character limit)`,
						},
						"message": schema.StringAttribute{
							Computed:    true,
							Description: `Enter a message to display to all your Organization's users, across all Cribl products. Limited to one line and 100 characters; will be truncated as needed.`,
						},
						"theme": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *CustomBannerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CustomBannerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CustomBannerDataSourceModel
	items, err := restclient.Get[[]customBannerAPI](ctx, r.client, customBannerPath)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	if items != nil {
		data.Items = make([]bannerMessage, 0, len(*items))
		for index := range *items {
			data.Items = append(data.Items, bannerMessageFromAPI(&(*items)[index]))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
