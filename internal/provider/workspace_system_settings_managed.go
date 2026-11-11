package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const workspaceSettingsEndpoint = "/system/settings/conf"
const workspaceSettingsManagedFieldsKey = "workspace_settings_managed_fields"

// workspaceSystemSettingsManagedResource preserves unconfigured fields across
// section-replacement PATCHes. Its schema and data source remain generated.
type workspaceSystemSettingsManagedResource struct {
	WorkspaceSystemSettingsResource
}

var _ resource.ResourceWithModifyPlan = &workspaceSystemSettingsManagedResource{}

func newWorkspaceSystemSettingsManagedResource() resource.Resource {
	return &workspaceSystemSettingsManagedResource{}
}

func (r *workspaceSystemSettingsManagedResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	if !req.Config.Raw.IsFullyKnown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("managed_fields"), types.ListUnknown(types.StringType))...)
		return
	}
	var model WorkspaceSystemSettingsModel
	// The effective plan respects lifecycle.ignore_changes.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := workspaceSettingsPayload(model)
	if err != nil {
		resp.Diagnostics.AddError("Unable to record workspace settings ownership", err.Error())
		return
	}
	fields, diags := settingsManagedFields(ctx, settingsMask(payload))
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("managed_fields"), fields)...)
}

func (r *workspaceSystemSettingsManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var private settingsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	r.write(ctx, req.Plan, &resp.State, private, &resp.Diagnostics)
}

func (r *workspaceSystemSettingsManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var private settingsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	r.write(ctx, req.Plan, &resp.State, private, &resp.Diagnostics)
}

func (r *workspaceSystemSettingsManagedResource) write(ctx context.Context, plan tfsdk.Plan, state *tfsdk.State, private settingsPrivate, diags *diag.Diagnostics) {
	var model WorkspaceSystemSettingsModel
	diags.Append(plan.Get(ctx, &model)...)
	if diags.HasError() {
		return
	}
	configured, err := workspaceSettingsPayload(model)
	if err != nil {
		diags.AddError("Unable to encode workspace settings", err.Error())
		return
	}
	current, err := restclient.Get[map[string]json.RawMessage](ctx, r.client, workspaceSettingsEndpoint)
	if err != nil {
		diags.AddError("Unable to read existing workspace settings", err.Error())
		return
	}
	if current == nil || len(*current) == 0 {
		diags.AddError("Unable to read existing workspace settings", "The API returned an empty configuration; no settings were changed.")
		return
	}
	// Sending only touched sections also avoids reapplying operator-owned sections.
	merged := make(map[string]json.RawMessage, len(configured))
	for name, value := range configured {
		existing := map[string]json.RawMessage{}
		if raw, ok := (*current)[name]; ok {
			existing[name] = raw
		}
		section, err := mergeSettingsFields(existing, map[string]json.RawMessage{name: value})
		if err != nil {
			diags.AddError("Unable to merge workspace settings", err.Error())
			return
		}
		merged[name] = section[name]
	}
	mask := settingsMask(configured)
	refreshed, err := workspaceSettingsState(ctx, *current, mask, model)
	if err != nil {
		diags.AddError("Unable to read current workspace settings", err.Error())
		return
	}
	remote, err := workspaceSettingsPayload(refreshed)
	if err != nil {
		diags.AddError("Unable to encode current workspace settings", err.Error())
		return
	}
	// A redacted secret cannot establish that a new configured value already
	// exists remotely. Always send it when Terraform invokes a write.
	hasSecret := false
	if !model.API.IsNull() && !model.API.IsUnknown() {
		if ssl, ok := model.API.Attributes()["ssl"].(types.Object); ok && !ssl.IsNull() && !ssl.IsUnknown() {
			if secret, ok := ssl.Attributes()["passphrase"].(types.String); ok {
				hasSecret = !secret.IsNull() && !secret.IsUnknown()
			}
		}
	}
	if len(merged) != 0 && (hasSecret || !reflect.DeepEqual(remote, configured)) {
		if err := validateWorkspaceUpgradeSettings(merged); err != nil {
			diags.AddAttributeError(path.Root("upgrade_settings"), "Invalid automatic upgrade settings", err.Error())
			return
		}
		if err := r.patch(ctx, merged, *current); err != nil {
			diags.AddError("Unable to update workspace settings", err.Error())
			return
		}
		current, err = restclient.Get[map[string]json.RawMessage](ctx, r.client, workspaceSettingsEndpoint)
		if err != nil {
			diags.AddError("Unable to read updated workspace settings", err.Error())
			return
		}
	}
	if current == nil || len(*current) == 0 {
		diags.AddError("Unable to read updated workspace settings", "The API returned an empty configuration.")
		return
	}
	// Keep omitted or redacted write-only values from the effective plan.
	model, err = workspaceSettingsState(ctx, *current, settingsMask(configured), model)
	if err != nil {
		diags.AddError("Unable to refresh workspace settings", err.Error())
		return
	}
	setWorkspaceSettingsState(ctx, state, private, model, settingsMask(configured), diags)
}

func (r *workspaceSystemSettingsManagedResource) patch(ctx context.Context, payload, current map[string]json.RawMessage) error {
	err := restclient.PatchNoResponse(ctx, r.client, workspaceSettingsEndpoint, payload)
	if err == nil {
		return nil
	}
	// Some Cloud backends compare the incoming API address with the saved
	// address even when the entire API section was omitted from the PATCH.
	// Retry only that specific rejection, and only when API settings are unmanaged.
	var httpErr *restclient.HTTPError
	if _, touchesAPI := payload["api"]; touchesAPI || !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusForbidden {
		return err
	}
	var response struct {
		Message string `json:"message"`
	}
	if json.Unmarshal([]byte(httpErr.Body), &response) != nil || response.Message != "Updating API host/port is disabled in this deployment" {
		return err
	}
	raw, ok := current["api"]
	if !ok || string(raw) == "null" {
		return err
	}
	var api struct {
		SSL struct {
			Passphrase string `json:"passphrase"`
		} `json:"ssl"`
	}
	if decodeErr := json.Unmarshal(raw, &api); decodeErr != nil {
		return fmt.Errorf("read API settings for Cloud compatibility retry: %w", decodeErr)
	}
	if isLikelyMaskedSecret(api.SSL.Passphrase) {
		return fmt.Errorf("%w; the API section contains a masked SSL passphrase and cannot safely be resent unchanged", err)
	}
	// Preserve the complete section, including fields unknown to the provider:
	// supplying only host/port would replace and discard its other settings.
	withAPI := make(map[string]json.RawMessage, len(payload)+1)
	for name, value := range payload {
		withAPI[name] = value
	}
	withAPI["api"] = raw
	return restclient.PatchNoResponse(ctx, r.client, workspaceSettingsEndpoint, withAPI)
}

// Validate the merged section so an existing, unmanaged check period can satisfy
// the API dependency without requiring Terraform to take ownership of it.
func validateWorkspaceUpgradeSettings(payload map[string]json.RawMessage) error {
	raw, ok := payload["upgradeSettings"]
	if !ok {
		return nil
	}
	var settings struct {
		Disabled *bool   `json:"disableAutomaticUpgrade"`
		Period   *string `json:"automaticUpgradeCheckPeriod"`
	}
	if err := json.Unmarshal(raw, &settings); err != nil {
		return fmt.Errorf("decode upgrade settings: %w", err)
	}
	if settings.Disabled != nil && !*settings.Disabled && (settings.Period == nil || *settings.Period == "") {
		return fmt.Errorf("enabling automatic upgrades requires automatic_upgrade_check_period; set it to a duration such as \"24h\" because the existing workspace settings do not contain a check period to preserve; no settings were changed")
	}
	return nil
}

func (r *workspaceSystemSettingsManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var prior WorkspaceSystemSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	current, err := restclient.Get[map[string]json.RawMessage](ctx, r.client, workspaceSettingsEndpoint)
	if restclient.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read workspace settings", err.Error())
		return
	}
	if current == nil || len(*current) == 0 {
		resp.Diagnostics.AddError("Unable to read workspace settings", "The API returned an empty configuration.")
		return
	}
	var mask settingsFieldMask
	if req.Private != nil {
		data, diags := req.Private.GetKey(ctx, workspaceSettingsManagedFieldsKey)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(data) != 0 {
			if err := json.Unmarshal(data, &mask); err != nil {
				resp.Diagnostics.AddError("Unable to read workspace settings ownership", err.Error())
				return
			}
		}
	}
	model, err := workspaceSettingsState(ctx, *current, mask, prior)
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh workspace settings", err.Error())
		return
	}
	if mask == nil {
		// An ID-only import hydrates supported settings. The next apply records
		// the configured subset, including an empty ownership mask.
		payload, err := workspaceSettingsPayload(model)
		if err != nil {
			resp.Diagnostics.AddError("Unable to record workspace settings ownership", err.Error())
			return
		}
		mask = settingsMask(payload)
	}
	var private settingsPrivate
	if resp.Private != nil {
		private = resp.Private
	}
	setWorkspaceSettingsState(ctx, &resp.State, private, model, mask, &resp.Diagnostics)
}

func (r *workspaceSystemSettingsManagedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "conf" {
		resp.Diagnostics.AddError("Invalid workspace settings import ID", "Use the fixed import ID conf. The provider configuration selects the workspace.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func workspaceSettingsPayload(model WorkspaceSystemSettingsModel) (map[string]json.RawMessage, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	var payload map[string]json.RawMessage
	err = json.Unmarshal(data, &payload)
	return payload, err
}

func workspaceSettingsState(ctx context.Context, current map[string]json.RawMessage, mask settingsFieldMask, prior WorkspaceSystemSettingsModel) (WorkspaceSystemSettingsModel, error) {
	selected := current
	var err error
	if mask != nil {
		selected, err = selectSettingsFields(current, mask)
		if err != nil {
			return WorkspaceSystemSettingsModel{}, err
		}
	}
	data, err := json.Marshal(selected)
	if err != nil {
		return WorkspaceSystemSettingsModel{}, err
	}
	var model WorkspaceSystemSettingsModel
	if err := json.Unmarshal(data, &model); err != nil {
		return model, err
	}
	model.ID = types.StringValue("conf")
	// API SSL passphrases are write-only or redacted in some deployments.
	if !model.API.IsNull() && !prior.API.IsNull() && !prior.API.IsUnknown() {
		attrs := model.API.Attributes()
		oldSSL, ok := prior.API.Attributes()["ssl"].(types.Object)
		ssl, currentOK := attrs["ssl"].(types.Object)
		if ok && currentOK && !oldSSL.IsNull() && !oldSSL.IsUnknown() && !ssl.IsNull() {
			values := ssl.Attributes()
			old, oldOK := oldSSL.Attributes()["passphrase"].(types.String)
			value, valueOK := values["passphrase"].(types.String)
			if oldOK && valueOK && !old.IsNull() && !old.IsUnknown() && (value.IsNull() || isLikelyMaskedSecret(value.ValueString())) {
				values["passphrase"] = old
				updated, diags := types.ObjectValue(ssl.AttributeTypes(ctx), values)
				if diags.HasError() {
					return model, fmt.Errorf("restore SSL passphrase: %v", diags)
				}
				attrs["ssl"] = updated
				model.API, diags = types.ObjectValue(model.API.AttributeTypes(ctx), attrs)
				if diags.HasError() {
					return model, fmt.Errorf("restore API settings: %v", diags)
				}
			}
		}
	}
	return model, nil
}

func setWorkspaceSettingsState(ctx context.Context, state *tfsdk.State, private settingsPrivate, model WorkspaceSystemSettingsModel, mask settingsFieldMask, diags *diag.Diagnostics) {
	fields, fieldDiags := settingsManagedFields(ctx, mask)
	diags.Append(fieldDiags...)
	if diags.HasError() {
		return
	}
	model.ManagedFields = fields
	diags.Append(state.Set(ctx, &model)...)
	if private != nil {
		data, err := json.Marshal(mask)
		if err != nil {
			diags.AddError("Unable to record workspace settings ownership", err.Error())
			return
		}
		diags.Append(private.SetKey(ctx, workspaceSettingsManagedFieldsKey, data)...)
	}
}
