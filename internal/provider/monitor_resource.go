// Schema aligned with IAetosMonitorConf (packages/metrics-types/src/shared/types.ts).
// NOT code-generated — intentionally diverges from the upstream OpenAPI MonitorConf,
// which maps to System Insights, not the Customer Metrics interface.
package provider

import (
	"context"
	"fmt"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	custom_stringplanmodifier "github.com/criblio/terraform-provider-criblio/internal/tfplanmodifiers/stringplanmodifier"
	custom_validators "github.com/criblio/terraform-provider-criblio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ = jsontypes.NormalizedType{}
var _ = types.String{}

var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithConfigure = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

type MonitorResource struct {
	client *restclient.Client
	api    MonitorAPI
}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

func (r *MonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Cribl Customer Metrics monitor (`IAetosMonitorConf`). " +
			"When the Terraform provider provisions a monitor, the backend detects the " +
			"`speakeasy-sdk/terraform` User-Agent and stamps `managed_by = \"terraform\"` on the config.",
		Attributes: map[string]schema.Attribute{
			// ── Identity ─────────────────────────────────────────────────────
			"id": schema.StringAttribute{
				Required: true,
				Optional: false,
				Computed: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
				},
				Description: "Unique identifier for the monitor. Immutable — changing it forces a replacement.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Human-readable name for the monitor.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Whether the monitor is active and evaluated.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Monitor type. One of: threshold, change, anomaly, forecast, logs.",
			},
			"description": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Optional human-readable description.",
			},
			// ── Inheritable scalar fields (JSON) ─────────────────────────────
			"priority": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Monitor priority as JSON (Inheritable<IScalarValue>). ` +
					`Inline example: jsonencode({value = "P2"}).`,
			},
			"team": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Owning team as JSON (Inheritable<IScalarValue>). ` +
					`Inline example: jsonencode({value = "ops"}).`,
			},
			// ── Query (JSON) ──────────────────────────────────────────────────
			"query": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Monitor queries keyed by label as JSON (Record<QueryLabel, IAetosMonitorQuery>). ` +
					`Example: jsonencode({A = {mode = "promql", promql = "up"}}).`,
			},
			"expr": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
				Description: "Derived query expressions as a JSON array (IQueryExpression[]). Use jsonencode([]) when no expressions are needed.",
			},
			// ── Firing logic (JSON) ───────────────────────────────────────────
			"firing_condition": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Firing and recovery delays as JSON (Inheritable<IFiringCondition>). ` +
					`Inline example: jsonencode({fire_delay = 300, clear_delay = 60}).`,
			},
			"firing_rule": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Alert thresholds as JSON (Inheritable<IFiringRule>). ` +
					`Example: jsonencode({label = "down", threshold = [{severity = "critical", limit = 0, includedTags = [], excludedTags = []}]}).`,
			},
			// ── Notification / metadata / silence ─────────────────────────────
			"metadata": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
				Description: "Arbitrary key-value metadata as a JSON object. Use jsonencode({}) for an empty map.",
			},
			"notification": schema.StringAttribute{
				Required:   true,
				Optional:   false,
				Computed:   false,
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
				Description: `Notification configuration as JSON (Inheritable<INotificationConfig>). ` +
					`Example: jsonencode({enabled = false, type = "policy", config = []}).`,
			},
			"silence": schema.ListAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				ElementType: types.StringType,
				Description: "IDs of silence windows that suppress this monitor's alerts. Use [] for none.",
			},
			// ── Optional fields ───────────────────────────────────────────────
			"dataset_id": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Default dataset for query execution and alert history storage (e.g. \"metrics\"). Required for the monitor to evaluate queries in production; omit only in test environments where no datasets are configured.",
			},
			"detection_config": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
				Description: "Type-specific detection configuration as JSON. Present for change, anomaly, outlier, and forecast monitors.",
			},
			"unit": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Optional unit label applied to the monitor's values.",
			},
			// ── Computed (read-only, stamped by backend) ──────────────────────
			"managed_by": schema.StringAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Set to "terraform" when the monitor was provisioned via the Terraform provider. Read-only — the backend stamps this automatically based on the User-Agent.`,
			},
		},
	}
}

func (r *MonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.api = newMonitorAPI(r.client)
}

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model MonitorModel
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(plan.As(ctx, &model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiModel, err := r.api.Create(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	applyMonitorAPIToState(apiModel, &model, true, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model MonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiModel, err := r.api.Read(ctx, model)
	if restclient.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	applyMonitorAPIToState(apiModel, &model, true, isMonitorImportState(&model))
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model MonitorModel
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(plan.As(ctx, &model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiModel, err := r.api.Update(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		return
	}
	applyMonitorAPIToState(apiModel, &model, true, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model MonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.api.Delete(ctx, model); err != nil && !restclient.IsNotFound(err) {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// isMonitorImportState returns true when the state was produced by an import
// (required fields are null/unknown), signalling applyMonitorAPIToState to
// fill every field from the API response rather than only computed ones.
func isMonitorImportState(state *MonitorModel) bool {
	if state == nil {
		return false
	}
	return state.Enabled.IsNull() || state.Enabled.IsUnknown() ||
		state.Name.IsNull() || state.Name.IsUnknown() ||
		state.Type.IsNull() || state.Type.IsUnknown() ||
		state.Query.IsNull() || state.Query.IsUnknown() ||
		state.FiringCondition.IsNull() || state.FiringCondition.IsUnknown() ||
		state.FiringRule.IsNull() || state.FiringRule.IsUnknown()
}

// applyMonitorAPIToState copies API response fields into the Terraform state.
//
// preserveInputs=true: only update computed fields and (if fillMissingInputs)
// fields that are null/unknown in state (import case).
// preserveInputs=false: overwrite all fields (data source reads).
func applyMonitorAPIToState(api *MonitorModel, state *MonitorModel, preserveInputs bool, fillMissingInputs bool) {
	if api == nil || state == nil {
		return
	}

	copyStr := func(dst *types.String, src types.String) {
		if !preserveInputs || (fillMissingInputs && (dst.IsNull() || dst.IsUnknown())) {
			if !src.IsNull() && !src.IsUnknown() {
				*dst = src
			}
		}
	}
	copyBool := func(dst *types.Bool, src types.Bool) {
		if !preserveInputs || (fillMissingInputs && (dst.IsNull() || dst.IsUnknown())) {
			if !src.IsNull() && !src.IsUnknown() {
				*dst = src
			}
		}
	}
	copyJSON := func(dst *jsontypes.Normalized, src jsontypes.Normalized) {
		if !preserveInputs || (fillMissingInputs && (dst.IsNull() || dst.IsUnknown())) {
			if !src.IsNull() && !src.IsUnknown() {
				*dst = src
			}
		}
	}
	copyList := func(dst *types.List, src types.List) {
		if !preserveInputs || (fillMissingInputs && (dst.IsNull() || dst.IsUnknown())) {
			if !src.IsNull() && !src.IsUnknown() {
				*dst = src
			}
		}
	}

	copyStr(&state.ID, api.ID)
	copyStr(&state.Name, api.Name)
	copyBool(&state.Enabled, api.Enabled)
	copyStr(&state.Type, api.Type)
	copyStr(&state.Description, api.Description)
	copyJSON(&state.Priority, api.Priority)
	copyJSON(&state.Team, api.Team)
	copyStr(&state.DatasetId, api.DatasetId)
	copyJSON(&state.Query, api.Query)
	copyJSON(&state.Expr, api.Expr)
	copyJSON(&state.FiringCondition, api.FiringCondition)
	copyJSON(&state.FiringRule, api.FiringRule)
	copyJSON(&state.Metadata, api.Metadata)
	copyJSON(&state.Notification, api.Notification)
	copyList(&state.Silence, api.Silence)
	copyJSON(&state.DetectionConfig, api.DetectionConfig)
	copyStr(&state.Unit, api.Unit)

	// managed_by is always computed — copy unconditionally from API response.
	if !api.ManagedBy.IsNull() && !api.ManagedBy.IsUnknown() {
		state.ManagedBy = api.ManagedBy
	} else if state.ManagedBy.IsNull() || state.ManagedBy.IsUnknown() {
		state.ManagedBy = types.StringValue("")
	}

	// Ensure silence list always has a concrete element type even when empty.
	if elementType := state.Silence.ElementType(context.Background()); elementType == nil {
		state.Silence = types.ListNull(types.StringType)
	}
}

func MonitorDebug(value any) string {
	return fmt.Sprintf("%v", value)
}
