package provider

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// systemLimitsManagedResource retains unconfigured limits across replacement PATCHes.
// The generated resource supplies the schema, configuration, import and no-op delete.
type systemLimitsManagedResource struct{ SystemLimitsResource }

var _ resource.Resource = &systemLimitsManagedResource{}
var _ resource.ResourceWithConfigure = &systemLimitsManagedResource{}
var _ resource.ResourceWithImportState = &systemLimitsManagedResource{}
var _ resource.ResourceWithModifyPlan = &systemLimitsManagedResource{}

func newSystemLimitsManagedResource() resource.Resource { return &systemLimitsManagedResource{} }

const systemLimitsManagedFieldsKey = "system_limits_managed_fields"

// This field is feature-gated by the backend. It remains in replacement
// PATCHes unchanged and is always refreshed into state.
var systemLimitsReadOnlyFields = []string{"metricsDirectory"}

func (r *systemLimitsManagedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	// Terraform does not persist private-only changes on a no-op apply. This
	// computed field ensures an ownership change has a durable state difference.
	if !req.Config.Raw.IsFullyKnown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("managed_fields"), types.ListUnknown(types.StringType))...)
		return
	}
	var model SystemLimitsModel
	// Use the effective plan so lifecycle.ignore_changes and the payload agree
	// about which fields remain managed.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configured, err := systemLimitsPayload(model)
	if err != nil {
		resp.Diagnostics.AddError("Unable to record system limits ownership", err.Error())
		return
	}
	fields, diags := settingsManagedFields(ctx, settingsMask(configured))
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("managed_fields"), fields)...)
	}
}

func (r *systemLimitsManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var private settingsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	r.write(ctx, req.Plan, nil, &resp.State, private, &resp.Diagnostics)
}

func (r *systemLimitsManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var private settingsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	var prior SystemLimitsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.write(ctx, req.Plan, &prior, &resp.State, private, &resp.Diagnostics)
}

func (r *systemLimitsManagedResource) write(ctx context.Context, plan tfsdk.Plan, prior *SystemLimitsModel, state *tfsdk.State, private settingsPrivate, diags *diag.Diagnostics) {
	var model SystemLimitsModel
	diags.Append(plan.Get(ctx, &model)...)
	if diags.HasError() {
		return
	}
	configured, err := systemLimitsPayload(model)
	if err != nil {
		diags.AddError("Unable to encode system limits", err.Error())
		return
	}
	if prior != nil {
		previous, err := systemLimitsPayload(*prior)
		if err != nil {
			diags.AddError("Unable to encode prior system limits", err.Error())
			return
		}
		if reflect.DeepEqual(previous, configured) {
			// Only ownership changed; the refreshed values already match the plan.
			copySystemLimitsReadOnly(*prior, &model)
			setSystemLimitsState(ctx, state, private, model, configured, diags)
			return
		}
	}
	endpoint := systemLimitsEndpoint(model)
	current, err := restclient.Get[map[string]json.RawMessage](ctx, r.client, endpoint)
	if err != nil {
		diags.AddError("Unable to read existing system limits", err.Error())
		return
	}
	if current == nil || len(*current) == 0 {
		diags.AddError("Unable to read existing system limits", "The API returned an empty configuration; no limits were changed.")
		return
	}
	readOnly, err := systemLimitsReadOnlyModel(*current)
	if err != nil {
		diags.AddError("Unable to read protected system limits", err.Error())
		return
	}
	copySystemLimitsReadOnly(readOnly, &model)
	merged, err := mergeSettingsFields(*current, configured)
	if err != nil {
		diags.AddError("Unable to merge system limits", err.Error())
		return
	}
	if err := restclient.PatchNoResponse(ctx, r.client, endpoint, merged); err != nil {
		diags.AddError("Unable to update system limits", err.Error())
		return
	}
	setSystemLimitsState(ctx, state, private, model, configured, diags)
}

func setSystemLimitsState(ctx context.Context, state *tfsdk.State, private settingsPrivate, model SystemLimitsModel, configured map[string]json.RawMessage, diags *diag.Diagnostics) {
	mask := settingsMask(configured)
	fields, fieldDiags := settingsManagedFields(ctx, mask)
	diags.Append(fieldDiags...)
	if diags.HasError() {
		return
	}
	model.ManagedFields = fields
	// Track configured writable fields and current read-only values without
	// expanding ownership to unrelated remote settings.
	diags.Append(state.Set(ctx, &model)...)
	if private != nil {
		diags.Append(writeSystemLimitsMask(ctx, private, mask)...)
	}
}

func (r *systemLimitsManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var prior SystemLimitsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	current, err := restclient.Get[map[string]json.RawMessage](ctx, r.client, systemLimitsEndpoint(prior))
	if restclient.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read system limits", err.Error())
		return
	}
	if current == nil || len(*current) == 0 {
		resp.Diagnostics.AddError("Unable to read system limits", "The API returned an empty configuration.")
		return
	}
	var mask settingsFieldMask
	if req.Private != nil {
		data, diags := req.Private.GetKey(ctx, systemLimitsManagedFieldsKey)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(data) != 0 {
			if err := json.Unmarshal(data, &mask); err != nil {
				resp.Diagnostics.AddError("Unable to read system limits ownership", err.Error())
				return
			}
		}
	}
	if mask == nil {
		// Legacy state records its populated fields; a group-only import hydrates
		// all supported fields. The next apply records the selected configuration.
		configured, err := systemLimitsPayload(prior)
		if err != nil {
			resp.Diagnostics.AddError("Unable to read system limits ownership", err.Error())
			return
		}
		if len(configured) != 0 {
			mask = settingsMask(configured)
		}
	}
	// Older state may have owned fields that are now read-only.
	for _, name := range systemLimitsReadOnlyFields {
		delete(mask, name)
	}
	selected := *current
	if mask != nil {
		selected, err = selectSettingsFields(*current, mask)
		if err != nil {
			resp.Diagnostics.AddError("Unable to refresh managed system limits", err.Error())
			return
		}
	}
	for _, name := range systemLimitsReadOnlyFields {
		if value, ok := (*current)[name]; ok {
			selected[name] = value
		}
	}
	data, err := json.Marshal(selected)
	if err != nil {
		resp.Diagnostics.AddError("Unable to decode system limits", err.Error())
		return
	}
	var model SystemLimitsModel
	if err := json.Unmarshal(data, &model); err != nil {
		resp.Diagnostics.AddError("Unable to decode system limits", err.Error())
		return
	}
	model.GroupID = prior.GroupID
	if mask == nil {
		payload, err := systemLimitsPayload(model)
		if err != nil {
			resp.Diagnostics.AddError("Unable to record system limits ownership", err.Error())
			return
		}
		mask = settingsMask(payload)
	}
	fields, fieldDiags := settingsManagedFields(ctx, mask)
	resp.Diagnostics.Append(fieldDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.ManagedFields = fields
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	if resp.Private != nil {
		resp.Diagnostics.Append(writeSystemLimitsMask(ctx, resp.Private, mask)...)
	}
}

func systemLimitsReadOnlyModel(current map[string]json.RawMessage) (SystemLimitsModel, error) {
	selected := make(map[string]json.RawMessage, len(systemLimitsReadOnlyFields))
	for _, name := range systemLimitsReadOnlyFields {
		if value, ok := current[name]; ok {
			selected[name] = value
		}
	}
	data, err := json.Marshal(selected)
	if err != nil {
		return SystemLimitsModel{}, err
	}
	var model SystemLimitsModel
	err = json.Unmarshal(data, &model)
	return model, err
}

func copySystemLimitsReadOnly(source SystemLimitsModel, target *SystemLimitsModel) {
	target.MetricsDirectory = source.MetricsDirectory
}

func systemLimitsEndpoint(model SystemLimitsModel) string {
	return "/m/" + url.PathEscape(model.GroupID.ValueString()) + "/system/limits"
}

func systemLimitsPayload(model SystemLimitsModel) (map[string]json.RawMessage, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeSystemLimitsMask(ctx context.Context, private settingsPrivate, mask settingsFieldMask) diag.Diagnostics {
	data, err := json.Marshal(mask)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Unable to record system limits ownership", err.Error())}
	}
	return private.SetKey(ctx, systemLimitsManagedFieldsKey, data)
}
