package provider

import (
	"context"
	"fmt"
	"regexp"

	custom_boolplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/boolplanmodifier"
	custom_float64planmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/float64planmodifier"
	custom_listplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/listplanmodifier"
	custom_objectplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/objectplanmodifier"
	custom_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/stringplanmodifier"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &CustomBannerResource{}
var _ resource.ResourceWithConfigure = &CustomBannerResource{}
var _ resource.ResourceWithImportState = &CustomBannerResource{}

func NewCustomBannerResource() resource.Resource {
	return &CustomBannerResource{}
}

type CustomBannerResource struct {
	client *restclient.Client
}

type CustomBannerResourceModel struct {
	Created         types.Float64           `tfsdk:"created"`
	CustomThemes    []types.String          `tfsdk:"custom_themes"`
	Enabled         types.Bool              `tfsdk:"enabled"`
	ID              types.String            `tfsdk:"id"`
	InvertFontColor types.Bool              `tfsdk:"invert_font_color"`
	Items           []tfTypes.BannerMessage `tfsdk:"items"`
	Link            types.String            `tfsdk:"link"`
	LinkDisplay     types.String            `tfsdk:"link_display"`
	Message         types.String            `tfsdk:"message"`
	Theme           types.String            `tfsdk:"theme"`
	Type            types.String            `tfsdk:"type"`
}

type customBannerAPI struct {
	Created         *float64 `json:"created,omitempty"`
	CustomThemes    []string `json:"customThemes,omitempty"`
	Enabled         bool     `json:"enabled"`
	ID              *string  `json:"id,omitempty"`
	InvertFontColor *bool    `json:"invertFontColor,omitempty"`
	Link            *string  `json:"link,omitempty"`
	LinkDisplay     *string  `json:"linkDisplay,omitempty"`
	Message         string   `json:"message"`
	Theme           string   `json:"theme"`
	Type            string   `json:"type"`
}

func (r *CustomBannerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_banner"
}

func (r *CustomBannerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "CustomBanner Resource",
		Attributes: map[string]schema.Attribute{
			"created": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: `Time created`,
			},
			"custom_themes": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: `Show a banner on top of all pages`,
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
				},
			},
			"invert_font_color": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					custom_listplanmodifier.SuppressDiff(custom_listplanmodifier.ExplicitSuppress),
				},
				NestedObject: schema.NestedAttributeObject{
					PlanModifiers: []planmodifier.Object{
						custom_objectplanmodifier.SuppressDiff(custom_objectplanmodifier.ExplicitSuppress),
					},
					Attributes: map[string]schema.Attribute{
						"created": schema.Float64Attribute{
							Computed: true,
							PlanModifiers: []planmodifier.Float64{
								custom_float64planmodifier.SuppressDiff(custom_float64planmodifier.ExplicitSuppress),
							},
							Description: `Time created`,
						},
						"custom_themes": schema.ListAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.List{
								custom_listplanmodifier.SuppressDiff(custom_listplanmodifier.ExplicitSuppress),
							},
							ElementType: types.StringType,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								custom_boolplanmodifier.SuppressDiff(custom_boolplanmodifier.ExplicitSuppress),
							},
							Description: `Show a banner on top of all pages`,
						},
						"id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
						},
						"invert_font_color": schema.BoolAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								custom_boolplanmodifier.SuppressDiff(custom_boolplanmodifier.ExplicitSuppress),
							},
						},
						"link": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
							Description: `Optionally, provide a URL to append to the message`,
						},
						"link_display": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
							Description: `Optionally, display your link with a short text label instead of the raw URL (100-character limit)`,
						},
						"message": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
							Description: `Enter a message to display to all your Organization's users, across all Cribl products. Limited to one line and 100 characters; will be truncated as needed.`,
						},
						"theme": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
						},
						"type": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
							},
						},
					},
				},
			},
			"link": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: `Optionally, provide a URL to append to the message`,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^https?://`), "must match pattern "+regexp.MustCompile(`^https?://`).String()),
				},
			},
			"link_display": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: `Optionally, display your link with a short text label instead of the raw URL (100-character limit)`,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(100),
				},
			},
			"message": schema.StringAttribute{
				Required:    true,
				Description: `Enter a message to display to all your Organization's users, across all Cribl products. Limited to one line and 100 characters; will be truncated as needed.`,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(100),
				},
			},
			"theme": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^((#?[0-9a-fA-F]{6})|(orange)|(yellow)|(green)|(blue)|(purple)|(magenta)|(red)){1}$`), "must match pattern "+regexp.MustCompile(`^((#?[0-9a-fA-F]{6})|(orange)|(yellow)|(green)|(blue)|(purple)|(magenta)|(red)){1}$`).String()),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: `must be one of ["custom", "system"]`,
				Validators: []validator.String{
					stringvalidator.OneOf("custom", "system"),
				},
			},
		},
	}
}

func (r *CustomBannerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clients, ok := req.ProviderData.(*ProviderClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = clients.RC
}

func (r *CustomBannerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CustomBannerResourceModel
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(plan.As(ctx, &data, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.createCustomBanner(ctx, data); err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	resp.Diagnostics.Append(refreshPlan(ctx, plan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomBannerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CustomBannerResourceModel
	var item types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
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

	if err := r.refreshCustomBannerState(ctx, data); err != nil {
		if restclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomBannerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CustomBannerResourceModel
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	merge(ctx, req, resp, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := restclient.Patch[customBannerAPI, customBannerAPI](ctx, r.client, customBannerPath, data.toCustomBannerAPI()); err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	if err := r.refreshCustomBannerState(ctx, data); err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	resp.Diagnostics.Append(refreshPlan(ctx, plan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomBannerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CustomBannerResourceModel
	var item types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
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

	apiModel := data.toCustomBannerAPI()
	apiModel.Enabled = false
	if _, err := restclient.Patch[customBannerAPI, customBannerAPI](ctx, r.client, customBannerPath, apiModel); err != nil && !restclient.IsNotFound(err) {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
}

func (r *CustomBannerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" && req.ID != customBannerID {
		resp.Diagnostics.AddError("Invalid import ID", `The custom banner import ID must be "custom-banner".`)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), customBannerID)...)
}

func (r *CustomBannerResource) refreshCustomBannerState(ctx context.Context, data *CustomBannerResourceModel) error {
	apiModel, err := restclient.Get[customBannerAPI](ctx, r.client, customBannerPath)
	if err != nil {
		return err
	}
	if apiModel == nil {
		return fmt.Errorf("response envelope contained no items")
	}
	data.applyCustomBannerAPI(apiModel)
	return nil
}

func (data *CustomBannerResourceModel) toCustomBannerAPI() customBannerAPI {
	id := customBannerID
	return customBannerAPI{
		Created:         float64PointerFromValue(data.Created),
		CustomThemes:    stringSliceFromValues(data.CustomThemes),
		Enabled:         data.Enabled.ValueBool(),
		ID:              &id,
		InvertFontColor: boolPointerFromValue(data.InvertFontColor),
		Link:            stringPointerFromValue(data.Link),
		LinkDisplay:     stringPointerFromValue(data.LinkDisplay),
		Message:         data.Message.ValueString(),
		Theme:           data.Theme.ValueString(),
		Type:            data.Type.ValueString(),
	}
}

func (data *CustomBannerResourceModel) applyCustomBannerAPI(apiModel *customBannerAPI) {
	if apiModel == nil {
		return
	}
	item := bannerMessageFromAPI(apiModel)
	data.Items = []tfTypes.BannerMessage{item}
	data.Created = item.Created
	data.CustomThemes = item.CustomThemes
	data.Enabled = item.Enabled
	data.ID = item.ID
	data.InvertFontColor = item.InvertFontColor
	data.Link = item.Link
	data.LinkDisplay = item.LinkDisplay
	data.Message = item.Message
	data.Theme = item.Theme
	data.Type = item.Type
}

func bannerMessageFromAPI(apiModel *customBannerAPI) tfTypes.BannerMessage {
	return tfTypes.BannerMessage{
		Created:         types.Float64PointerValue(apiModel.Created),
		CustomThemes:    stringValuesFromSlice(apiModel.CustomThemes),
		Enabled:         types.BoolValue(apiModel.Enabled),
		ID:              types.StringPointerValue(apiModel.ID),
		InvertFontColor: types.BoolPointerValue(apiModel.InvertFontColor),
		Link:            types.StringPointerValue(apiModel.Link),
		LinkDisplay:     types.StringPointerValue(apiModel.LinkDisplay),
		Message:         types.StringValue(apiModel.Message),
		Theme:           types.StringValue(apiModel.Theme),
		Type:            types.StringValue(apiModel.Type),
	}
}
