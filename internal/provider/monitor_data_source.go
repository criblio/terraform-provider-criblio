// Schema aligned with IAetosMonitorConf — NOT code-generated.
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
		MarkdownDescription: "Reads a single Cribl Customer Metrics monitor by ID.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Required: true},
			"name":        schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"type":        schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"priority": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Monitor priority as JSON (Inheritable<IScalarValue>).",
			},
			"team": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Owning team as JSON (Inheritable<IScalarValue>).",
			},
			"dataset_id": schema.StringAttribute{Computed: true},
			"query": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Monitor queries keyed by label as JSON.",
			},
			"expr": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Derived query expressions as JSON array.",
			},
			"firing_condition": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Firing and recovery delays as JSON.",
			},
			"firing_rule": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Alert thresholds as JSON.",
			},
			"metadata": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Key-value metadata as JSON.",
			},
			"notification": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Notification configuration as JSON.",
			},
			"silence": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "IDs of active silence windows.",
			},
			"detection_config": schema.StringAttribute{
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
				Description: "Type-specific detection configuration as JSON.",
			},
			"unit":       schema.StringAttribute{Computed: true},
			"managed_by": schema.StringAttribute{Computed: true},
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
