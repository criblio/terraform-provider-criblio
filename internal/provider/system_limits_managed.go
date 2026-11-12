package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"

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

type systemLimitsPrivate interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
	SetKey(context.Context, string, []byte) diag.Diagnostics
}

// Leaves are nil, objects contain the configured child fields. Persisting this
// mask keeps ownership when the API removes a configured value during refresh.
type systemLimitsMask map[string]systemLimitsMask

const systemLimitsManagedFieldsKey = "system_limits_managed_fields"

// These fields are operator-controlled or feature-gated by the backend. They
// remain in replacement PATCHes unchanged and are always refreshed into state.
var systemLimitsReadOnlyFields = []string{
	"lookupAllowedPaths", "lookupMaxSize", "lookupMaxTotalSize",
	"enableWorkerInMemory", "metricsDirectory",
}

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
	fields, diags := systemLimitsManagedFields(ctx, systemLimitsFieldMask(configured))
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("managed_fields"), fields)...)
	}
}

func (r *systemLimitsManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var private systemLimitsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	r.write(ctx, req.Plan, nil, &resp.State, private, &resp.Diagnostics)
}

func (r *systemLimitsManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var private systemLimitsPrivate
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

func (r *systemLimitsManagedResource) write(ctx context.Context, plan tfsdk.Plan, prior *SystemLimitsModel, state *tfsdk.State, private systemLimitsPrivate, diags *diag.Diagnostics) {
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
	merged, err := mergeSystemLimits(*current, configured)
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

func setSystemLimitsState(ctx context.Context, state *tfsdk.State, private systemLimitsPrivate, model SystemLimitsModel, configured map[string]json.RawMessage, diags *diag.Diagnostics) {
	mask := systemLimitsFieldMask(configured)
	fields, fieldDiags := systemLimitsManagedFields(ctx, mask)
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
	var mask systemLimitsMask
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
			mask = systemLimitsFieldMask(configured)
		}
	}
	// Older state may have owned fields that are now read-only.
	for _, name := range systemLimitsReadOnlyFields {
		delete(mask, name)
	}
	selected := *current
	if mask != nil {
		selected, err = selectSystemLimits(*current, mask)
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
		mask = systemLimitsFieldMask(payload)
	}
	fields, fieldDiags := systemLimitsManagedFields(ctx, mask)
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
	target.LookupAllowedPaths = source.LookupAllowedPaths
	target.LookupMaxSize = source.LookupMaxSize
	target.LookupMaxTotalSize = source.LookupMaxTotalSize
	target.EnableWorkerInMemory = source.EnableWorkerInMemory
	target.MetricsDirectory = source.MetricsDirectory
}

func systemLimitsManagedFields(ctx context.Context, mask systemLimitsMask) (types.List, diag.Diagnostics) {
	fields := make([]string, 0)
	var collect func(systemLimitsMask, string)
	collect = func(mask systemLimitsMask, prefix string) {
		for name, nested := range mask {
			field := prefix + name
			if len(nested) == 0 {
				fields = append(fields, field)
				continue
			}
			collect(nested, field+".")
		}
	}
	collect(mask, "")
	sort.Strings(fields)
	return types.ListValueFrom(ctx, types.StringType, fields)
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

func systemLimitsFieldMask(payload map[string]json.RawMessage) systemLimitsMask {
	mask := make(systemLimitsMask, len(payload))
	for name, raw := range payload {
		var nested map[string]json.RawMessage
		if json.Unmarshal(raw, &nested) == nil && nested != nil {
			mask[name] = systemLimitsFieldMask(nested)
		} else {
			mask[name] = nil
		}
	}
	return mask
}

func writeSystemLimitsMask(ctx context.Context, private systemLimitsPrivate, mask systemLimitsMask) diag.Diagnostics {
	data, err := json.Marshal(mask)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Unable to record system limits ownership", err.Error())}
	}
	return private.SetKey(ctx, systemLimitsManagedFieldsKey, data)
}

// RawMessage preserves unknown fields and exact integer values from newer servers.
func mergeSystemLimits(current, configured map[string]json.RawMessage) (map[string]json.RawMessage, error) {
	merged := make(map[string]json.RawMessage, len(current)+len(configured))
	for name, value := range current {
		merged[name] = value
	}
	for name, value := range configured {
		var requested, existing map[string]json.RawMessage
		if json.Unmarshal(value, &requested) == nil && requested != nil {
			if raw, ok := current[name]; ok && string(raw) != "null" {
				if err := json.Unmarshal(raw, &existing); err != nil {
					return nil, fmt.Errorf("existing %s is not an object: %w", name, err)
				}
			}
			nested, err := mergeSystemLimits(existing, requested)
			if err != nil {
				return nil, err
			}
			value, err = json.Marshal(nested)
			if err != nil {
				return nil, err
			}
		}
		merged[name] = value
	}
	return merged, nil
}

func selectSystemLimits(current map[string]json.RawMessage, mask systemLimitsMask) (map[string]json.RawMessage, error) {
	selected := make(map[string]json.RawMessage, len(mask))
	for name, nestedMask := range mask {
		value, ok := current[name]
		if !ok {
			continue
		}
		if nestedMask != nil && string(value) != "null" {
			var nested map[string]json.RawMessage
			if err := json.Unmarshal(value, &nested); err != nil {
				return nil, fmt.Errorf("managed %s is not an object: %w", name, err)
			}
			filtered, err := selectSystemLimits(nested, nestedMask)
			if err != nil {
				return nil, err
			}
			value, err = json.Marshal(filtered)
			if err != nil {
				return nil, err
			}
		}
		selected[name] = value
	}
	return selected, nil
}
