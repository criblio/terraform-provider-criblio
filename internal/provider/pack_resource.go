package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	cribl_listplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/listplanmodifier"
	cribl_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/planmodifiers/stringplanmodifier"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PackResource{}
var _ resource.ResourceWithImportState = &PackResource{}

func NewPackResource() resource.Resource {
	return &PackResource{}
}

// PackResource defines the resource implementation.
type PackResource struct {
	client *restclient.Client
}

// PackResourceModel describes the resource data model.
type PackResourceModel struct {
	AllowCustomFunctions types.Bool                   `tfsdk:"allow_custom_functions"`
	Author               types.String                 `tfsdk:"author"`
	Description          types.String                 `tfsdk:"description"`
	Disabled             types.Bool                   `queryParam:"style=form,explode=true,name=disabled" tfsdk:"disabled"`
	DisplayName          types.String                 `tfsdk:"display_name"`
	Exports              []types.String               `tfsdk:"exports"`
	Filename             types.String                 `queryParam:"style=form,explode=true,name=filename" tfsdk:"filename"`
	Force                types.Bool                   `tfsdk:"force"`
	GroupID              types.String                 `tfsdk:"group_id"`
	ID                   types.String                 `tfsdk:"id"`
	Inputs               types.Float64                `tfsdk:"inputs"`
	Items                []tfTypes.PackInstallInfo    `tfsdk:"items"`
	MinLogStreamVersion  types.String                 `tfsdk:"min_log_stream_version"`
	Outputs              types.Float64                `tfsdk:"outputs"`
	Source               types.String                 `tfsdk:"source"`
	Spec                 types.String                 `tfsdk:"spec"`
	Tags                 *tfTypes.PackRequestBodyTags `tfsdk:"tags"`
	Version              types.String                 `tfsdk:"version"`
}

func (r *PackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pack"
}

func (r *PackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pack Resource",
		Attributes: map[string]schema.Attribute{
			"allow_custom_functions": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Requires replacement if changed.`,
			},
			"author": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: `Pack author (from pack metadata). Config changes are applied via pack/settings PATCH.`,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: `Pack description (from pack metadata). Config changes are applied via pack/settings PATCH.`,
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: `Pack display name (from pack metadata). Config changes are applied via pack/settings PATCH.`,
			},
			"exports": schema.ListAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
				},
				ElementType: types.StringType,
				Description: `Requires replacement if changed.`,
			},
			"filename": schema.StringAttribute{
				Optional:    true,
				Description: `Local .crbl file path to upload. File is uploaded (PUT) then the pack is installed or updated in place (PATCH); changing filename updates the existing pack rather than replacing it. When set, description and display_name come from the pack file; omit them from config to avoid drift.`,
			},
			"force": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Requires replacement if changed.`,
			},
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `The consumer group to which this instance belongs. Defaults to 'Cribl'. Requires replacement if changed.`,
			},
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Pack name. Requires replacement if changed.`,
			},
			"inputs": schema.Float64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Requires replacement if changed.`,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					cribl_listplanmodifier.PreferState(),
				},
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
			"min_log_stream_version": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					cribl_stringplanmodifier.PreferState(),
				},
				Description: `Min LogStream version (from pack metadata). Preserved from state when not configured.`,
			},
			"outputs": schema.Float64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplaceIfConfigured(),
				},
				Description: `Requires replacement if changed.`,
			},
			"source": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					cribl_stringplanmodifier.PreferState(),
				},
				Description: `Pack source path (from pack metadata). Preserved from state when not configured.`,
			},
			"spec": schema.StringAttribute{
				Optional:    true,
				Description: `body string optional Specify a branch, tag or a semver spec`,
			},
			"tags": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"data_type": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: `Pack data_type tags (from pack metadata).`,
					},
					"domain": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: `Pack domain tags (from pack metadata).`,
					},
					"streamtags": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: `Pack streamtags (from pack metadata).`,
					},
					"technology": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: `Pack technology tags (from pack metadata).`,
					},
				},
				Description: `Pack tags (from pack metadata). Changes are reflected in state from the API; no replacement.`,
			},
			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: `Pack version (from pack metadata). Changes are reflected in state from the API; no replacement.`,
			},
		},
	}
}

func (r *PackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

type packTagsAPI struct {
	DataType   []string `json:"dataType,omitempty"`
	Domain     []string `json:"domain,omitempty"`
	Streamtags []string `json:"streamtags,omitempty"`
	Technology []string `json:"technology,omitempty"`
}

type packAPIModel struct {
	Author              *string         `json:"author,omitempty"`
	Description         *string         `json:"description,omitempty"`
	DisplayName         *string         `json:"displayName,omitempty"`
	Exports             []string        `json:"exports,omitempty"`
	ID                  string          `json:"id,omitempty"`
	Inputs              *float64        `json:"inputs,omitempty"`
	MinLogStreamVersion *string         `json:"minLogStreamVersion,omitempty"`
	Outputs             *float64        `json:"outputs,omitempty"`
	Settings            map[string]any  `json:"settings,omitempty"`
	Source              *string         `json:"source,omitempty"`
	Spec                *string         `json:"spec,omitempty"`
	Tags                *packTagsAPI    `json:"tags,omitempty"`
	Version             *string         `json:"version,omitempty"`
	Warnings            json.RawMessage `json:"warnings,omitempty"`
}

type packRequestBody struct {
	AllowCustomFunctions *bool        `json:"allowCustomFunctions,omitempty"`
	Author               *string      `json:"author,omitempty"`
	Description          *string      `json:"description,omitempty"`
	DisplayName          *string      `json:"displayName,omitempty"`
	Exports              []string     `json:"exports,omitempty"`
	Force                *bool        `json:"force,omitempty"`
	ID                   string       `json:"id,omitempty"`
	Inputs               *float64     `json:"inputs,omitempty"`
	MinLogStreamVersion  *string      `json:"minLogStreamVersion,omitempty"`
	Outputs              *float64     `json:"outputs,omitempty"`
	Source               *string      `json:"source,omitempty"`
	Spec                 *string      `json:"spec,omitempty"`
	Tags                 *packTagsAPI `json:"tags,omitempty"`
	Version              *string      `json:"version,omitempty"`
}

type packInstallItemsRequest struct {
	Items []packInstallItemRequest `json:"items"`
	Count int                      `json:"count"`
}

type packInstallItemRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Version     string   `json:"version,omitempty"`
	Warnings    []string `json:"warnings,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
}

func (r *PackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PackResourceModel
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

	var cfg PackResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	mergePackCreateConfigIntoModel(data, &cfg)

	uploadedSource, err := r.uploadPackFileFromModel(ctx, data, true)
	if err != nil {
		resp.Diagnostics.AddError("Failed to upload pack file", err.Error())
		return
	}

	packID := effectivePackIDForAPI(data)
	if uploadedSource != "" && r.packExists(ctx, data.GroupID.ValueString(), packID) {
		if err := r.installUploadedPack(ctx, data, uploadedSource); err != nil {
			resp.Diagnostics.AddError("failure to invoke API", err.Error())
			return
		}
	} else {
		apiModel, err := r.createPack(ctx, data)
		if err != nil {
			resp.Diagnostics.AddError("failure to invoke API", err.Error())
			return
		}
		if apiModel != nil {
			data.applyPackAPIModel(apiModel)
		}
	}

	if err := r.patchPackSettings(ctx, data.GroupID.ValueString(), r.effectivePackID(ctx, data), data); err != nil {
		resp.Diagnostics.AddError("pack settings sync failed", fmt.Sprintf("Could not update pack metadata via pack/settings: %v", err))
		return
	}
	if err := r.refreshPackState(ctx, data); err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	preservePackMetadataFromConfig(ctx, data, plan)
	resp.Diagnostics.Append(refreshPlan(ctx, plan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PackResourceModel
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

	savedDesc, savedDisplay, savedVersion := data.Description, data.DisplayName, data.Version
	if err := r.refreshPackState(ctx, data); err != nil {
		if restclient.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	if !savedDesc.IsNull() && !savedDesc.IsUnknown() {
		data.Description = savedDesc
	}
	if !savedDisplay.IsNull() && !savedDisplay.IsUnknown() {
		data.DisplayName = savedDisplay
	}
	if !savedVersion.IsNull() && !savedVersion.IsUnknown() {
		data.Version = savedVersion
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PackResourceModel
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	merge(ctx, req, resp, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	var cfg PackResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uploadedSource, err := r.uploadPackFileFromModel(ctx, data, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to upload pack file", err.Error())
		return
	}
	switch {
	case uploadedSource != "":
		if err := r.installUploadedPack(ctx, data, uploadedSource); err != nil {
			resp.Diagnostics.AddError("failure to invoke API", err.Error())
			return
		}
	case configuredString(cfg.Source) != "":
		var stateData PackResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if configuredString(cfg.Source) != configuredString(stateData.Source) {
			if _, err := r.patchPackByIDWithSource(ctx, data.GroupID.ValueString(), r.effectivePackID(ctx, data), shortNameForPatch(configuredString(cfg.Source)), boolPointerFromValue(data.Disabled)); err != nil {
				resp.Diagnostics.AddError("failure to invoke API", err.Error())
				return
			}
		}
	}

	if err := r.patchPackSettings(ctx, data.GroupID.ValueString(), r.effectivePackID(ctx, data), data); err != nil {
		resp.Diagnostics.AddError("pack settings sync failed", fmt.Sprintf("Could not update pack metadata via pack/settings: %v", err))
		return
	}
	if err := r.refreshPackState(ctx, data); err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	preservePackMetadataFromConfig(ctx, data, plan)
	resp.Diagnostics.Append(refreshPlan(ctx, plan, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PackResourceModel
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

	err := restclient.Delete(ctx, r.client, fmt.Sprintf("/m/%s/packs/%s", url.PathEscape(data.GroupID.ValueString()), url.PathEscape(r.effectivePackID(ctx, data))))
	if err == nil || restclient.IsNotFound(err) {
		return
	}
	var httpErr *restclient.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == 500 && (strings.Contains(httpErr.Body, "referenced by") || strings.Contains(httpErr.Body, "Cannot uninstall")) {
		resp.Diagnostics.AddError(
			"Cannot delete pack: it is in use",
			"The pack is referenced by a route or other resource. Remove or update the route (e.g. criblio_routes) so it no longer references this pack, then try again.",
		)
		return
	}
	resp.Diagnostics.AddError("failure to invoke API", err.Error())
}

func (r *PackResource) createPack(ctx context.Context, data *PackResourceModel) (*packAPIModel, error) {
	query := url.Values{}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		query.Set("disabled", fmt.Sprintf("%t", data.Disabled.ValueBool()))
	}
	path := fmt.Sprintf("/m/%s/packs", url.PathEscape(data.GroupID.ValueString()))
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	return restclient.Post[packRequestBody, packAPIModel](ctx, r.client, path, packRequestFromModel(data))
}

func (r *PackResource) refreshPackState(ctx context.Context, data *PackResourceModel) error {
	apiModel, err := restclient.Get[packAPIModel](ctx, r.client, fmt.Sprintf("/m/%s/packs/%s", url.PathEscape(data.GroupID.ValueString()), url.PathEscape(r.effectivePackID(ctx, data))))
	if err != nil {
		return err
	}
	if apiModel == nil {
		return fmt.Errorf("response envelope contained no items")
	}
	data.applyPackAPIModel(apiModel)
	return nil
}

func (r *PackResource) packExists(ctx context.Context, groupID, packID string) bool {
	if strings.TrimSpace(packID) == "" {
		return false
	}
	packID = resolvePackIDForRestAPI(ctx, r.client, groupID, packID)
	_, err := restclient.Get[packAPIModel](ctx, r.client, fmt.Sprintf("/m/%s/packs/%s", url.PathEscape(groupID), url.PathEscape(packID)))
	return err == nil
}

func (r *PackResource) uploadPackFileFromModel(ctx context.Context, data *PackResourceModel, requireExistingFile bool) (string, error) {
	if data.Filename.IsNull() || data.Filename.IsUnknown() || strings.TrimSpace(data.Filename.ValueString()) == "" {
		return "", nil
	}
	filename := data.Filename.ValueString()
	filePath, err := resolvePackFilePath(filename)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filePath); err != nil {
		if requireExistingFile || data.Source.IsNull() || data.Source.IsUnknown() {
			wd, _ := os.Getwd()
			return "", fmt.Errorf("file does not exist: %s\n\nWorking directory: %s\n\nPlease use an absolute path or Terraform's path functions:\n  filename = \"${path.module}/%s\"\n  or\n  filename = \"${path.root}/%s\"", filePath, wd, filepath.Base(filename), filepath.Base(filename))
		}
		return "", nil
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to read file %s: %v", filePath, err)
	}
	baseFilename := filepath.Base(filePath)
	fullSource, shortName, err := r.uploadPackFile(ctx, data.GroupID.ValueString(), fileContent, baseFilename)
	if err != nil {
		return "", err
	}
	if shortName == "" {
		shortName = baseFilename
	}
	data.Filename = types.StringValue(shortName)
	data.Source = types.StringValue(shortName)
	return fullSource, nil
}

func resolvePackFilePath(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filepath.Clean(filename), nil
	}
	filePath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path %s: %v", filename, err)
	}
	return filepath.Clean(filePath), nil
}

func (r *PackResource) uploadPackFile(ctx context.Context, groupID string, fileContent []byte, filename string) (storedSource string, storedShortName string, err error) {
	query := url.Values{}
	query.Set("filename", filename)
	query.Set("size", fmt.Sprintf("%d", len(fileContent)))
	body, err := restclient.PutRaw(ctx, r.client, fmt.Sprintf("/m/%s/packs?%s", url.PathEscape(groupID), query.Encode()), "application/octet-stream", fileContent)
	if err != nil {
		return "", "", err
	}
	var responseData struct {
		Source string           `json:"source"`
		Items  []map[string]any `json:"items"`
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &responseData); err != nil {
			return "", "", fmt.Errorf("failed to parse upload response: %v", err)
		}
	}
	source := responseData.Source
	if source == "" && len(responseData.Items) > 0 {
		if itemSource, ok := responseData.Items[0]["source"].(string); ok {
			source = itemSource
		}
	}
	if source == "" {
		source = "file:" + filename
	}
	shortName := source
	if strings.HasPrefix(source, "file:") {
		shortName = filepath.Base(strings.TrimPrefix(source, "file:"))
	}
	if shortName == "" {
		shortName = filename
	}
	return source, shortName, nil
}

func (r *PackResource) installUploadedPack(ctx context.Context, data *PackResourceModel, uploadedSource string) error {
	fullPath := fullSourcePath(uploadedSource)
	version := ""
	if len(data.Items) > 0 && !data.Items[0].Version.IsNull() && !data.Items[0].Version.IsUnknown() {
		version = data.Items[0].Version.ValueString()
	}
	if _, err := r.postPacksInstallWithItems(ctx, data.GroupID.ValueString(), r.effectivePackID(ctx, data), fullPath, version, packInstallDisplayName(data)); err == nil {
		return nil
	}
	_, err := r.patchPackByIDWithSource(ctx, data.GroupID.ValueString(), r.effectivePackID(ctx, data), shortNameForPatch(uploadedSource), boolPointerFromValue(data.Disabled))
	return err
}

func (r *PackResource) postPacksInstallWithItems(ctx context.Context, groupID, packID, sourceFullPath, version, displayName string) (*packAPIModel, error) {
	if strings.TrimSpace(packID) == "" {
		return nil, fmt.Errorf("pack id is empty")
	}
	if strings.TrimSpace(sourceFullPath) == "" {
		return nil, fmt.Errorf("pack source path is empty")
	}
	if displayName == "" {
		displayName = packID
	}
	body := packInstallItemsRequest{
		Items: []packInstallItemRequest{{
			ID:          packID,
			Name:        displayName,
			Source:      sourceFullPath,
			Version:     version,
			Warnings:    []string{},
			DisplayName: displayName,
		}},
		Count: 1,
	}
	return restclient.Post[packInstallItemsRequest, packAPIModel](ctx, r.client, fmt.Sprintf("/m/%s/packs", url.PathEscape(groupID)), body)
}

func (r *PackResource) patchPackByIDWithSource(ctx context.Context, groupID, packID, source string, disabled *bool) (*packAPIModel, error) {
	form := url.Values{}
	form.Set("source", source)
	if disabled != nil {
		form.Set("disabled", fmt.Sprintf("%t", *disabled))
	}
	path := fmt.Sprintf("/m/%s/packs/%s?%s", url.PathEscape(groupID), url.PathEscape(packID), form.Encode())
	body, err := restclient.PatchRaw(ctx, r.client, path, "application/x-www-form-urlencoded", []byte(form.Encode()))
	if err != nil {
		var httpErr *restclient.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == 500 && strings.Contains(httpErr.Body, "up to date") {
			return nil, nil
		}
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	var envelope struct {
		Items []packAPIModel `json:"items"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Items) > 0 {
		return &envelope.Items[0], nil
	}
	var apiModel packAPIModel
	if err := json.Unmarshal(body, &apiModel); err != nil {
		return nil, fmt.Errorf("failed to parse PATCH response: %v", err)
	}
	return &apiModel, nil
}

func (r *PackResource) patchPackSettings(ctx context.Context, groupID, packID string, data *PackResourceModel) error {
	displayName := data.ID.ValueString()
	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() && data.DisplayName.ValueString() != "" {
		displayName = data.DisplayName.ValueString()
	}
	version := "0.0.0"
	if !data.Version.IsNull() && !data.Version.IsUnknown() && data.Version.ValueString() != "" {
		version = data.Version.ValueString()
	} else if len(data.Items) > 0 && !data.Items[0].Version.IsNull() && !data.Items[0].Version.IsUnknown() {
		version = data.Items[0].Version.ValueString()
	}
	packageObj := map[string]any{
		"displayName": displayName,
		"tags":        tagsAPIMap(data.Tags),
		"version":     version,
	}
	if !data.Author.IsNull() && !data.Author.IsUnknown() {
		packageObj["author"] = data.Author.ValueString()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		packageObj["description"] = data.Description.ValueString()
	}
	if !data.MinLogStreamVersion.IsNull() && !data.MinLogStreamVersion.IsUnknown() {
		packageObj["minLogStreamVersion"] = data.MinLogStreamVersion.ValueString()
	}
	return restclient.PatchNoResponse(ctx, r.client, fmt.Sprintf("/m/%s/p/%s/pack/settings", url.PathEscape(groupID), url.PathEscape(packID)), map[string]any{"package": packageObj})
}

func packRequestFromModel(data *PackResourceModel) packRequestBody {
	return packRequestBody{
		AllowCustomFunctions: boolPointerFromValue(data.AllowCustomFunctions),
		Author:               stringPointerFromValue(data.Author),
		Description:          stringPointerFromValue(data.Description),
		DisplayName:          stringPointerFromValue(data.DisplayName),
		Exports:              stringSliceFromValues(data.Exports),
		Force:                boolPointerFromValue(data.Force),
		ID:                   data.ID.ValueString(),
		Inputs:               float64PointerFromValue(data.Inputs),
		MinLogStreamVersion:  stringPointerFromValue(data.MinLogStreamVersion),
		Outputs:              float64PointerFromValue(data.Outputs),
		Source:               stringPointerFromValue(data.Source),
		Spec:                 stringPointerFromValue(data.Spec),
		Tags:                 packTagsAPIFromTF(data.Tags),
		Version:              stringPointerFromValue(data.Version),
	}
}

func (data *PackResourceModel) applyPackAPIModel(api *packAPIModel) {
	if api == nil {
		return
	}
	item := packInstallInfoFromAPI(api)
	data.Items = []tfTypes.PackInstallInfo{item}
	data.Author = item.Author
	data.Description = item.Description
	data.DisplayName = item.DisplayName
	data.MinLogStreamVersion = item.MinLogStreamVersion
	data.Source = item.Source
	data.Version = item.Version
	if item.Tags != nil {
		data.Tags = &tfTypes.PackRequestBodyTags{
			DataType:   item.Tags.DataType,
			Domain:     item.Tags.Domain,
			Streamtags: item.Tags.Streamtags,
			Technology: item.Tags.Technology,
		}
	} else {
		data.Tags = nil
	}
}

func packInstallInfoFromAPI(api *packAPIModel) tfTypes.PackInstallInfo {
	item := tfTypes.PackInstallInfo{
		Author:              types.StringPointerValue(api.Author),
		Description:         types.StringPointerValue(api.Description),
		DisplayName:         types.StringPointerValue(api.DisplayName),
		Exports:             stringValuesFromSlice(api.Exports),
		ID:                  types.StringValue(api.ID),
		Inputs:              types.Float64PointerValue(api.Inputs),
		MinLogStreamVersion: types.StringPointerValue(api.MinLogStreamVersion),
		Outputs:             types.Float64PointerValue(api.Outputs),
		Settings:            normalizedSettings(api.Settings),
		Source:              types.StringPointerValue(api.Source),
		Spec:                types.StringPointerValue(api.Spec),
		Tags:                packInstallTagsFromAPI(api.Tags),
		Version:             types.StringPointerValue(api.Version),
		Warnings:            jsontypes.NewNormalizedValue("null"),
	}
	if len(api.Warnings) > 0 {
		item.Warnings = jsontypes.NewNormalizedValue(string(api.Warnings))
	}
	return item
}

func normalizedSettings(settings map[string]any) map[string]jsontypes.Normalized {
	if len(settings) == 0 {
		return nil
	}
	result := make(map[string]jsontypes.Normalized, len(settings))
	for key, value := range settings {
		data, _ := json.Marshal(value)
		result[key] = jsontypes.NewNormalizedValue(string(data))
	}
	return result
}

func packInstallTagsFromAPI(tags *packTagsAPI) *tfTypes.PackInstallInfoTags {
	if tags == nil {
		return nil
	}
	return &tfTypes.PackInstallInfoTags{
		DataType:   stringValuesFromSlice(tags.DataType),
		Domain:     stringValuesFromSlice(tags.Domain),
		Streamtags: stringValuesFromSlice(tags.Streamtags),
		Technology: stringValuesFromSlice(tags.Technology),
	}
}

func packTagsAPIFromTF(tags *tfTypes.PackRequestBodyTags) *packTagsAPI {
	if tags == nil {
		return nil
	}
	return &packTagsAPI{
		DataType:   stringSliceFromValues(tags.DataType),
		Domain:     stringSliceFromValues(tags.Domain),
		Streamtags: stringSliceFromValues(tags.Streamtags),
		Technology: stringSliceFromValues(tags.Technology),
	}
}

func tagsAPIMap(tags *tfTypes.PackRequestBodyTags) map[string]any {
	out := map[string]any{}
	if tags == nil {
		return out
	}
	out["dataType"] = stringSliceFromValues(tags.DataType)
	out["domain"] = stringSliceFromValues(tags.Domain)
	out["streamtags"] = stringSliceFromValues(tags.Streamtags)
	out["technology"] = stringSliceFromValues(tags.Technology)
	return out
}

func stringValuesFromSlice(values []string) []types.String {
	out := make([]types.String, 0, len(values))
	for _, value := range values {
		out = append(out, types.StringValue(value))
	}
	return out
}

func stringSliceFromValues(values []types.String) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if !value.IsNull() && !value.IsUnknown() {
			out = append(out, value.ValueString())
		}
	}
	return out
}

func stringPointerFromValue(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := value.ValueString()
	return &result
}

func boolPointerFromValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := value.ValueBool()
	return &result
}

func float64PointerFromValue(value types.Float64) *float64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := value.ValueFloat64()
	return &result
}

func configuredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

func fullSourcePath(source string) string {
	if source == "" || strings.HasPrefix(source, "file:") {
		return source
	}
	return "file:/opt/cribl_config/state/packs/" + source
}

func mergePackCreateConfigIntoModel(data *PackResourceModel, cfg *PackResourceModel) {
	if data == nil || cfg == nil {
		return
	}
	if strings.TrimSpace(configuredString(cfg.ID)) != "" {
		data.ID = cfg.ID
	}
	if strings.TrimSpace(configuredString(cfg.GroupID)) != "" {
		data.GroupID = cfg.GroupID
	}
	if strings.TrimSpace(configuredString(cfg.DisplayName)) != "" {
		data.DisplayName = cfg.DisplayName
	}
	if strings.TrimSpace(configuredString(cfg.Filename)) != "" {
		data.Filename = cfg.Filename
	}
}

func effectivePackIDForAPI(data *PackResourceModel) string {
	if len(data.Items) > 0 && strings.TrimSpace(configuredString(data.Items[0].ID)) != "" {
		return strings.TrimSpace(data.Items[0].ID.ValueString())
	}
	return strings.TrimSpace(data.ID.ValueString())
}

func (r *PackResource) effectivePackID(ctx context.Context, data *PackResourceModel) string {
	if len(data.Items) > 0 && strings.TrimSpace(configuredString(data.Items[0].ID)) != "" {
		return strings.TrimSpace(data.Items[0].ID.ValueString())
	}
	return resolvePackIDForRestAPI(ctx, r.client, data.GroupID.ValueString(), data.ID.ValueString())
}

func packInstallDisplayName(data *PackResourceModel) string {
	if strings.TrimSpace(configuredString(data.DisplayName)) != "" {
		return data.DisplayName.ValueString()
	}
	if len(data.Items) > 0 && strings.TrimSpace(configuredString(data.Items[0].DisplayName)) != "" {
		return data.Items[0].DisplayName.ValueString()
	}
	return data.ID.ValueString()
}

func shortNameForPatch(source string) string {
	if source == "" {
		return source
	}
	if strings.HasPrefix(source, "file:") {
		return filepath.Base(strings.TrimPrefix(source, "file:"))
	}
	return filepath.Base(source)
}

func preservePackMetadataFromConfig(ctx context.Context, data *PackResourceModel, plan types.Object) {
	var planData PackResourceModel
	if diags := plan.As(ctx, &planData, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		return
	}
	if !planData.Description.IsNull() && !planData.Description.IsUnknown() {
		data.Description = planData.Description
	}
	if !planData.DisplayName.IsNull() && !planData.DisplayName.IsUnknown() {
		data.DisplayName = planData.DisplayName
	}
	if !planData.Version.IsNull() && !planData.Version.IsUnknown() {
		data.Version = planData.Version
	}
}

func (r *PackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	dec := json.NewDecoder(strings.NewReader(req.ID))
	dec.DisallowUnknownFields()
	var data struct {
		GroupID string `json:"group_id"`
		ID      string `json:"id"`
	}
	if err := dec.Decode(&data); err != nil {
		resp.Diagnostics.AddError("Invalid ID", `The import ID is not valid. It is expected to be a JSON object string with the format: '{"group_id": "Cribl", "id": "observability-pack"}': `+err.Error())
		return
	}
	if data.GroupID == "" {
		resp.Diagnostics.AddError("Missing required field", `The field group_id is required but was not found in the json encoded ID. It's expected to be a value alike '"Cribl"'`)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), data.GroupID)...)
	if data.ID == "" {
		resp.Diagnostics.AddError("Missing required field", `The field id is required but was not found in the json encoded ID. It's expected to be a value alike '"observability-pack"'`)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.ID)...)
}
