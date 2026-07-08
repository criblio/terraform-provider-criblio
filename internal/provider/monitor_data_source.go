// Hand-written: do not regenerate (listed in .codegen-ignore).
package provider

import (
	"context"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ = jsontypes.NormalizedType{}
var _ = types.String{}

var _ datasource.DataSource = &MonitorDataSource{}
var _ datasource.DataSourceWithConfigure = &MonitorDataSource{}

// MonitorDataSource implements the criblio_monitor data source.
type MonitorDataSource struct {
	client *restclient.Client
	api    MonitorAPI
}

func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

func (d *MonitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an Aetos monitor (`IAetosMonitorConf`) by ID from `GET /products/aetos/monitors/{id}`.",
		Attributes: map[string]schema.Attribute{
			// ── Lookup key ───────────────────────────────────────────────────────
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the monitor to look up.",
			},

			// ── Scalar fields ─────────────────────────────────────────────────────
			"name": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"dataset_id": schema.StringAttribute{
				Computed: true,
			},
			"detection_config": schema.StringAttribute{
				Computed: true,
			},
			"unit": schema.StringAttribute{
				Computed: true,
			},
			"managed_by": schema.StringAttribute{
				Computed: true,
			},

			// ── JSON blob fields ──────────────────────────────────────────────────
			"priority": schema.StringAttribute{
				Computed:    true,
				Description: "Monitor priority as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"team": schema.StringAttribute{
				Computed:    true,
				Description: "Owning team as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"query": schema.StringAttribute{
				Computed:    true,
				Description: "Map of named query expressions as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"expr": schema.StringAttribute{
				Computed:    true,
				Description: "Array of derived expression objects as a JSON array.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"firing_condition": schema.StringAttribute{
				Computed:    true,
				Description: "Firing / clear delay settings as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"firing_rule": schema.StringAttribute{
				Computed:    true,
				Description: "Threshold / firing rule as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"metadata": schema.StringAttribute{
				Computed:    true,
				Description: "Arbitrary metadata key/value map as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},
			"notification": schema.StringAttribute{
				Computed:    true,
				Description: "Notification policy configuration as a JSON object.",
				CustomType:  jsontypes.NormalizedType{},
			},

			// ── List field ────────────────────────────────────────────────────────
			"silence": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of silence window IDs that suppress alerts from this monitor.",
			},
		},
	}
}

func (d *MonitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clients, ok := req.ProviderData.(*ProviderClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = clients.RC
	d.api = newMonitorAPI(d.client)
}

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model MonitorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiModel, err := d.api.Read(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	applyMonitorAPIToState(apiModel, &model, false, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func MonitorDataSourceDebug(value any) string {
	return fmt.Sprintf("%v", value)
}
