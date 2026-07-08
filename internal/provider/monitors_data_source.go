// Hand-written: do not regenerate (listed in .codegen-ignore).
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

// MonitorsDataSource implements the criblio_monitors data source (list).
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
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all Aetos monitors from `GET /products/aetos/monitors`.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
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
						"silence": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of silence window IDs.",
						},
					},
				},
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
	// The API returns CountedMonitorConf {count, items[]}. restclient.Get unwraps
	// items automatically via decodeEnvelope when T is a slice.
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
			obj := map[string]attr.Value{
				"id":               item.ID,
				"name":             item.Name,
				"enabled":          item.Enabled,
				"type":             item.Type,
				"description":      item.Description,
				"dataset_id":       item.DatasetID,
				"detection_config": item.DetectionConfig,
				"unit":             item.Unit,
				"managed_by":       item.ManagedBy,
				"priority":         item.Priority,
				"team":             item.Team,
				"query":            item.Query,
				"expr":             item.Expr,
				"firing_condition": item.FiringCondition,
				"firing_rule":      item.FiringRule,
				"metadata":         item.Metadata,
				"notification":     item.Notification,
				"silence":          item.Silence,
			}
			values = append(values, types.ObjectValueMust(attrTypes, obj))
		}
	}
	model.Items = types.ListValueMust(types.ObjectType{AttrTypes: MonitorsItemAttrTypes()}, values)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// MonitorsItemAttrTypes returns the Terraform attribute type map for a single
// monitor item in the list data source.
func MonitorsItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"enabled":          types.BoolType,
		"type":             types.StringType,
		"description":      types.StringType,
		"dataset_id":       types.StringType,
		"detection_config": types.StringType,
		"unit":             types.StringType,
		"managed_by":       types.StringType,
		"priority":         jsontypes.NormalizedType{},
		"team":             jsontypes.NormalizedType{},
		"query":            jsontypes.NormalizedType{},
		"expr":             jsontypes.NormalizedType{},
		"firing_condition": jsontypes.NormalizedType{},
		"firing_rule":      jsontypes.NormalizedType{},
		"metadata":         jsontypes.NormalizedType{},
		"notification":     jsontypes.NormalizedType{},
		"silence":          types.ListType{ElemType: types.StringType},
	}
}
