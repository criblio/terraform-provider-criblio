package provider

import (
	"context"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/url"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GroupDataSource{}
var _ datasource.DataSourceWithConfigure = &GroupDataSource{}

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

// GroupDataSource is the data source implementation.
type GroupDataSource struct {
	client *restclient.Client
}

// GroupDataSourceModel describes the data model.
type GroupDataSourceModel struct {
	Cloud               *configGroupCloud `tfsdk:"cloud"`
	Description         types.String      `tfsdk:"description"`
	EstimatedIngestRate types.Float64     `tfsdk:"estimated_ingest_rate"`
	Fields              types.String      `queryParam:"style=form,explode=true,name=fields" tfsdk:"fields"`
	Git                 *groupGitModel    `tfsdk:"git"`
	ID                  types.String      `tfsdk:"id"`
	Inherits            types.String      `tfsdk:"inherits"`
	IsFleet             types.Bool        `tfsdk:"is_fleet"`
	MaxWorkerAge        types.String      `tfsdk:"max_worker_age"`
	Name                types.String      `tfsdk:"name"`
	OnPrem              types.Bool        `tfsdk:"on_prem"`
	Provisioned         types.Bool        `tfsdk:"provisioned"`
	Streamtags          types.List        `tfsdk:"streamtags"`
	Tags                types.String      `tfsdk:"tags"`
	Type                types.String      `tfsdk:"type"`
	WorkerRemoteAccess  types.Bool        `tfsdk:"worker_remote_access"`
}

type groupGitModel struct {
	Commit       types.String          `tfsdk:"commit"`
	LocalChanges types.Int64           `tfsdk:"local_changes"`
	Log          []groupGitCommitModel `tfsdk:"log"`
}

type groupGitCommitModel struct {
	AuthorEmail types.String `tfsdk:"author_email"`
	AuthorName  types.String `tfsdk:"author_name"`
	Date        types.String `tfsdk:"date"`
	Hash        types.String `tfsdk:"hash"`
	Message     types.String `tfsdk:"message"`
	Short       types.String `tfsdk:"short"`
}

// Metadata returns the data source type name.
func (r *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the data source.
func (r *GroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group DataSource",

		Attributes: map[string]schema.Attribute{
			"cloud": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"provider": schema.StringAttribute{
						Computed: true,
					},
					"region": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"estimated_ingest_rate": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: `Estimated ingest rate for the group. Supported values map to Max est ingest rate (MB/s):` + "\n" +
					`  - 1024 -> 12 MB/s` + "\n" +
					`  - 2048 -> 24 MB/s` + "\n" +
					`  - 3072 -> 36 MB/s` + "\n" +
					`  - 4096 -> 48 MB/s` + "\n" +
					`  - 5120 -> 60 MB/s` + "\n" +
					`  - 7168 -> 84 MB/s` + "\n" +
					`  - 10240 -> 120 MB/s` + "\n" +
					`  - 13312 -> 156 MB/s` + "\n" +
					`  - 15360 -> 180 MB/s`,
			},
			"fields": schema.StringAttribute{
				Optional:    true,
				Description: `fields to add to results: git.commit, git.localChanges, git.log`,
			},
			"git": schema.SingleNestedAttribute{
				Computed:    true,
				Description: `Git status fields requested through the fields argument.`,
				Attributes: map[string]schema.Attribute{
					"commit": schema.StringAttribute{
						Computed:    true,
						Description: `Commit hash of the currently committed configuration version.`,
					},
					"local_changes": schema.Int64Attribute{
						Computed:    true,
						Description: `Number of local configuration changes not yet committed.`,
					},
					"log": schema.ListNestedAttribute{
						Computed:    true,
						Description: `List of recent configuration commits.`,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"author_email": schema.StringAttribute{Computed: true},
								"author_name":  schema.StringAttribute{Computed: true},
								"date":         schema.StringAttribute{Computed: true},
								"hash":         schema.StringAttribute{Computed: true},
								"message":      schema.StringAttribute{Computed: true},
								"short":        schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: `Group id`,
			},
			"inherits": schema.StringAttribute{
				Computed: true,
			},
			"is_fleet": schema.BoolAttribute{
				Computed: true,
			},
			"max_worker_age": schema.StringAttribute{
				Computed:    true,
				Description: `This is only configurable for hybrid worker groups.`,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"on_prem": schema.BoolAttribute{
				Computed:    true,
				Description: `Whether this is an on-premises group. Cannot be true when cloud is set.`,
			},
			"provisioned": schema.BoolAttribute{
				Computed: true,
			},
			"streamtags": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"tags": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"worker_remote_access": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (r *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := fmt.Sprintf("/master/groups/%s", url.PathEscape(data.ID.ValueString()))
	if !data.Fields.IsNull() && !data.Fields.IsUnknown() {
		path += "?fields=" + url.QueryEscape(data.Fields.ValueString())
	}
	apiModel, err := restclient.Get[groupAPIModel](ctx, r.client, path)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	if apiModel == nil {
		resp.Diagnostics.AddError("unexpected response from API", "empty response body")
		return
	}
	data.applyGroupAPIModel(apiModel)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (data *GroupDataSourceModel) applyGroupAPIModel(api *groupAPIModel) {
	if api == nil {
		return
	}
	if api.Cloud == nil {
		data.Cloud = nil
	} else {
		data.Cloud = &configGroupCloud{
			Provider: types.StringPointerValue(api.Cloud.Provider),
			Region:   types.StringValue(api.Cloud.Region),
		}
	}
	data.Description = types.StringPointerValue(api.Description)
	data.EstimatedIngestRate = types.Float64PointerValue(api.EstimatedIngestRate)
	data.Git = groupGitModelFromAPI(api.Git)
	if api.ID != "" {
		data.ID = types.StringValue(api.ID)
	}
	data.Inherits = types.StringPointerValue(api.Inherits)
	data.IsFleet = types.BoolPointerValue(api.IsFleet)
	data.MaxWorkerAge = types.StringPointerValue(api.MaxWorkerAge)
	data.Name = types.StringPointerValue(api.Name)
	data.OnPrem = types.BoolPointerValue(api.OnPrem)
	data.Provisioned = types.BoolPointerValue(api.Provisioned)
	data.Streamtags = groupStringListFromSlice(api.Streamtags)
	data.Tags = types.StringPointerValue(api.Tags)
	data.Type = types.StringPointerValue(api.Type)
	data.WorkerRemoteAccess = types.BoolPointerValue(api.WorkerRemoteAccess)
}

func groupGitModelFromAPI(api *groupGitAPI) *groupGitModel {
	if api == nil {
		return nil
	}

	model := &groupGitModel{
		Commit:       types.StringPointerValue(api.Commit),
		LocalChanges: types.Int64PointerValue(api.LocalChanges),
	}
	if api.Log == nil {
		return model
	}

	model.Log = make([]groupGitCommitModel, 0, len(api.Log))
	for _, entry := range api.Log {
		model.Log = append(model.Log, groupGitCommitModel{
			AuthorEmail: types.StringPointerValue(entry.AuthorEmail),
			AuthorName:  types.StringPointerValue(entry.AuthorName),
			Date:        types.StringPointerValue(entry.Date),
			Hash:        types.StringPointerValue(entry.Hash),
			Message:     types.StringPointerValue(entry.Message),
			Short:       types.StringPointerValue(entry.Short),
		})
	}
	return model
}
