// Schema aligned with IAetosMonitorConf — NOT code-generated.
package provider

import (
	"context"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ = jsontypes.NormalizedType{}
var _ = types.String{}

var _ datasource.DataSource = &MonitorsDataSource{}
var _ datasource.DataSourceWithConfigure = &MonitorsDataSource{}

type MonitorsDataSource struct {
	client *restclient.Client
}

type MonitorsListDataSourceModel struct {
	Items types.List `tfsdk:"items"`
}

func NewMonitorsDataSource() datasource.DataSource {
	return &MonitorsDataSource{}
}

func (d *MonitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

func (d *MonitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	itemAttrs := map[string]schema.Attribute{
		"id":          schema.StringAttribute{Computed: true},
		"name":        schema.StringAttribute{Computed: true},
		"enabled":     schema.BoolAttribute{Computed: true},
		"type":        schema.StringAttribute{Computed: true},
		"description": schema.StringAttribute{Computed: true},
		"priority":    schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"team":        schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"dataset_id":  schema.StringAttribute{Computed: true},
		"query":       schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"expr":        schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"firing_condition": schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"firing_rule":      schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"metadata":         schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"notification":     schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"silence":          schema.ListAttribute{Computed: true, ElementType: types.StringType},
		"detection_config": schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}},
		"unit":             schema.StringAttribute{Computed: true},
		"managed_by":       schema.StringAttribute{Computed: true},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all Cribl Customer Metrics monitors.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: itemAttrs},
			},
		},
	}
}

func (d *MonitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model MonitorsListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	items, err := restclient.Get[[]MonitorModel](ctx, d.client, "/products/aetos/monitors")
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	values := []attr.Value{}
	if items != nil {
		values = make([]attr.Value, 0, len(*items))
		attrTypes := MonitorsItemAttrTypes()
		for _, item := range *items {
			values = append(values, types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"id":               item.ID,
				"name":             item.Name,
				"enabled":          item.Enabled,
				"type":             item.Type,
				"description":      item.Description,
				"priority":         item.Priority,
				"team":             item.Team,
				"dataset_id":       item.DatasetId,
				"query":            item.Query,
				"expr":             item.Expr,
				"firing_condition": item.FiringCondition,
				"firing_rule":      item.FiringRule,
				"metadata":         item.Metadata,
				"notification":     item.Notification,
				"silence":          item.Silence,
				"detection_config": item.DetectionConfig,
				"unit":             item.Unit,
				"managed_by":       item.ManagedBy,
			}))
		}
	}
	model.Items = types.ListValueMust(types.ObjectType{AttrTypes: MonitorsItemAttrTypes()}, values)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// MonitorsItemAttrTypes returns the attr.Type map for a single monitor item in the list.
func MonitorsItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"enabled":          types.BoolType,
		"type":             types.StringType,
		"description":      types.StringType,
		"priority":         jsontypes.NormalizedType{},
		"team":             jsontypes.NormalizedType{},
		"dataset_id":       types.StringType,
		"query":            jsontypes.NormalizedType{},
		"expr":             jsontypes.NormalizedType{},
		"firing_condition": jsontypes.NormalizedType{},
		"firing_rule":      jsontypes.NormalizedType{},
		"metadata":         jsontypes.NormalizedType{},
		"notification":     jsontypes.NormalizedType{},
		"silence":          types.ListType{ElemType: types.StringType},
		"detection_config": jsontypes.NormalizedType{},
		"unit":             types.StringType,
		"managed_by":       types.StringType,
	}
}
