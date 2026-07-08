// Hand-written: do not regenerate (listed in .codegen-ignore).
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

// MonitorResource implements the criblio_monitor Terraform resource.
// It manages Aetos monitors via GET/POST/PATCH/DELETE /products/aetos/monitors[/{id}].
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
		MarkdownDescription: "Manages an Aetos monitor (`IAetosMonitorConf`). " +
			"Create, read, update, and delete monitors via `GET/POST/PATCH/DELETE /products/aetos/monitors[/{id}]`.",
		Attributes: map[string]schema.Attribute{
			// ── Identifying ──────────────────────────────────────────────────────
			"id": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Unique identifier for the monitor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					custom_stringplanmodifier.SuppressDiff(custom_stringplanmodifier.ExplicitSuppress),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Human-readable name for the monitor.",
			},

			// ── Lifecycle ─────────────────────────────────────────────────────────
			"enabled": schema.BoolAttribute{
				Required:    true,
				Optional:    false,
				Computed:    false,
				Description: "Whether the monitor is actively evaluated.",
			},
			"type": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Monitor type (e.g. `threshold`, `anomaly`).",
			},

			// ── JSON blob fields ──────────────────────────────────────────────────
			"priority": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: `Monitor priority as a JSON object, e.g. ` + "`" + `jsonencode({ value = "P2" })` + "`" + `.`,
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"team": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: `Owning team as a JSON object, e.g. ` + "`" + `jsonencode({ value = "ops" })` + "`" + `.`,
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"query": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: `Map of named query expressions as a JSON object (` + "`" + `Record<string, QueryExpr>` + "`" + `). ` +
					"Use `jsonencode({...})` in HCL.",
				CustomType: jsontypes.NormalizedType{},
				Validators: []validator.String{custom_validators.IsValidJSON()},
			},
			"expr": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Array of derived expression objects as a JSON array. Use `jsonencode([...])` in HCL.",
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"firing_condition": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Firing / clear delay settings as a JSON object. Use `jsonencode({...})` in HCL.",
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"firing_rule": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Threshold / firing rule as a JSON object. Use `jsonencode({...})` in HCL.",
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"metadata": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Arbitrary metadata key/value map as a JSON object. Use `jsonencode({...})` in HCL.",
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},
			"notification": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Notification policy configuration as a JSON object. Use `jsonencode({...})` in HCL.",
				CustomType:  jsontypes.NormalizedType{},
				Validators:  []validator.String{custom_validators.IsValidJSON()},
			},

			// ── Scalar optional ───────────────────────────────────────────────────
			"description": schema.StringAttribute{
				Required: false,
				Optional: true,
				Computed: false,
			},
			"dataset_id": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Dataset the monitor queries against.",
			},
			"detection_config": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Detection configuration identifier.",
			},
			"unit": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				Description: "Unit for the metric values evaluated by this monitor.",
			},

			// ── List field ────────────────────────────────────────────────────────
			"silence": schema.ListAttribute{
				Required:    false,
				Optional:    true,
				Computed:    false,
				ElementType: types.StringType,
				Description: "List of silence window IDs that suppress alerts from this monitor.",
			},

			// ── Computed-only ─────────────────────────────────────────────────────
			"managed_by": schema.StringAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "Stamped to `terraform` by the backend when the Terraform provider User-Agent is detected. Read-only.",
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

// isMonitorImportState returns true when required fields are null/unknown,
// indicating this is an import (Terraform only provided the ID).
func isMonitorImportState(state *MonitorModel) bool {
	if state == nil {
		return false
	}
	if state.Name.IsNull() || state.Name.IsUnknown() {
		return true
	}
	if state.Enabled.IsNull() || state.Enabled.IsUnknown() {
		return true
	}
	return false
}

// applyMonitorAPIToState copies fields from the API response into the Terraform state.
//
//   - preserveInputs=true: only overwrite fields that are not already set in state
//     (used after create/update so user-provided values survive round-trips).
//   - fillMissingInputs=true: additionally fill null/unknown fields from the API
//     (used during import when state only contains the ID).
func applyMonitorAPIToState(api *MonitorModel, state *MonitorModel, preserveInputs bool, fillMissingInputs bool) {
	if api == nil || state == nil {
		return
	}

	// Helper: should we overwrite the state field with the API value?
	shouldSet := func(stateVal interface{ IsNull() bool; IsUnknown() bool }) bool {
		if !preserveInputs {
			return true
		}
		if fillMissingInputs && (stateVal.IsNull() || stateVal.IsUnknown()) {
			return true
		}
		return false
	}

	// ID — always propagate so the resource gets its server-assigned ID on create.
	if !api.ID.IsNull() && !api.ID.IsUnknown() {
		if shouldSet(&state.ID) {
			state.ID = api.ID
		} else if state.ID.IsNull() || state.ID.IsUnknown() {
			state.ID = api.ID
		}
	}

	if shouldSet(&state.Name) {
		if !api.Name.IsNull() && !api.Name.IsUnknown() {
			state.Name = api.Name
		}
	}
	if shouldSet(&state.Enabled) {
		if !api.Enabled.IsNull() && !api.Enabled.IsUnknown() {
			state.Enabled = api.Enabled
		}
	}
	if shouldSet(&state.Type) {
		if !api.Type.IsNull() && !api.Type.IsUnknown() {
			state.Type = api.Type
		}
	}
	if shouldSet(&state.Description) {
		if !api.Description.IsNull() && !api.Description.IsUnknown() {
			state.Description = api.Description
		}
	}
	if shouldSet(&state.DatasetID) {
		if !api.DatasetID.IsNull() && !api.DatasetID.IsUnknown() {
			state.DatasetID = api.DatasetID
		}
	}
	if shouldSet(&state.DetectionConfig) {
		if !api.DetectionConfig.IsNull() && !api.DetectionConfig.IsUnknown() {
			state.DetectionConfig = api.DetectionConfig
		}
	}
	if shouldSet(&state.Unit) {
		if !api.Unit.IsNull() && !api.Unit.IsUnknown() {
			state.Unit = api.Unit
		}
	}

	// JSON blob fields
	if shouldSet(&state.Priority) {
		if !api.Priority.IsNull() && !api.Priority.IsUnknown() {
			state.Priority = api.Priority
		}
	}
	if shouldSet(&state.Team) {
		if !api.Team.IsNull() && !api.Team.IsUnknown() {
			state.Team = api.Team
		}
	}
	if shouldSet(&state.Query) {
		if !api.Query.IsNull() && !api.Query.IsUnknown() {
			state.Query = api.Query
		}
	}
	if shouldSet(&state.Expr) {
		if !api.Expr.IsNull() && !api.Expr.IsUnknown() {
			state.Expr = api.Expr
		}
	}
	if shouldSet(&state.FiringCondition) {
		if !api.FiringCondition.IsNull() && !api.FiringCondition.IsUnknown() {
			state.FiringCondition = api.FiringCondition
		}
	}
	if shouldSet(&state.FiringRule) {
		if !api.FiringRule.IsNull() && !api.FiringRule.IsUnknown() {
			state.FiringRule = api.FiringRule
		}
	}
	if shouldSet(&state.Metadata) {
		if !api.Metadata.IsNull() && !api.Metadata.IsUnknown() {
			state.Metadata = api.Metadata
		}
	}
	if shouldSet(&state.Notification) {
		if !api.Notification.IsNull() && !api.Notification.IsUnknown() {
			state.Notification = api.Notification
		}
	}

	// silence list
	if shouldSet(&state.Silence) {
		if !api.Silence.IsNull() && !api.Silence.IsUnknown() {
			state.Silence = api.Silence
		}
	}
	if elemType := state.Silence.ElementType(context.Background()); elemType == nil {
		state.Silence = types.ListNull(types.StringType)
	}

	// managed_by — always propagate (computed-only, always set from API).
	if !api.ManagedBy.IsNull() && !api.ManagedBy.IsUnknown() {
		state.ManagedBy = api.ManagedBy
	} else if state.ManagedBy.IsNull() || state.ManagedBy.IsUnknown() {
		state.ManagedBy = types.StringNull()
	}
}
