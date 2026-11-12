package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/internal/searchacl"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The generated schema describes a declarative projection; the backend exposes
// separate user/team lists and add/remove actions rather than CRUD on that model.
type searchDatasetAclManagedResource struct{ SearchDatasetACLResource }
type searchDatasetAclManagedDataSource struct{ SearchDatasetACLDataSource }

var _ resource.ResourceWithValidateConfig = &searchDatasetAclManagedResource{}

func newSearchDatasetAclManagedResource() resource.Resource {
	return &searchDatasetAclManagedResource{}
}

func newSearchDatasetAclManagedDataSource() datasource.DataSource {
	return &searchDatasetAclManagedDataSource{}
}

func (r *searchDatasetAclManagedResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model SearchDatasetACLModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !model.DatasetID.IsUnknown() && strings.TrimSpace(model.DatasetID.ValueString()) == "" {
		resp.Diagnostics.AddAttributeError(path.Root("dataset_id"), "Missing dataset ID", "Specify the ID of an existing Search dataset.")
	}
	for name, grants := range map[string]types.Map{"users": model.Users, "teams": model.Teams} {
		if grants.IsNull() || grants.IsUnknown() {
			continue
		}
		for principal, value := range grants.Elements() {
			if strings.TrimSpace(principal) == "" {
				resp.Diagnostics.AddAttributeError(path.Root(name), "Invalid ACL principal", "Principal IDs cannot be empty.")
			}
			if value.IsUnknown() {
				continue
			}
			policy := value.(types.String)
			if policy.IsNull() || !searchacl.ValidPolicy(policy.ValueString()) {
				resp.Diagnostics.AddAttributeError(path.Root(name).AtMapKey(principal), "Invalid dataset policy", "Use DatasetRead or DatasetMaintain; null policies are not allowed.")
			}
		}
	}
}

func aclGrants(ctx context.Context, model SearchDatasetACLModel, diags *diag.Diagnostics) searchacl.Grants {
	var grants searchacl.Grants
	if !model.Users.IsNull() {
		diags.Append(model.Users.ElementsAs(ctx, &grants.Users, false)...)
	}
	if !model.Teams.IsNull() {
		diags.Append(model.Teams.ElementsAs(ctx, &grants.Teams, false)...)
	}
	return grants
}

func setACLGrants(ctx context.Context, model *SearchDatasetACLModel, grants searchacl.Grants, diags *diag.Diagnostics) {
	model.Users = types.MapNull(types.StringType)
	model.Teams = types.MapNull(types.StringType)
	if grants.Users != nil {
		value, valueDiags := types.MapValueFrom(ctx, types.StringType, grants.Users)
		diags.Append(valueDiags...)
		model.Users = value
	}
	if grants.Teams != nil {
		value, valueDiags := types.MapValueFrom(ctx, types.StringType, grants.Teams)
		diags.Append(valueDiags...)
		model.Teams = value
	}
}

func (r *searchDatasetAclManagedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.write(ctx, req.Plan, &resp.State, &resp.Diagnostics)
}

func (r *searchDatasetAclManagedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.write(ctx, req.Plan, &resp.State, &resp.Diagnostics)
}

func (r *searchDatasetAclManagedResource) write(ctx context.Context, plan tfsdk.Plan, state *tfsdk.State, diags *diag.Diagnostics) {
	var model SearchDatasetACLModel
	diags.Append(plan.Get(ctx, &model)...)
	if diags.HasError() {
		return
	}
	desired := aclGrants(ctx, model, diags)
	if diags.HasError() {
		return
	}
	client := searchacl.New(r.client, model.DatasetID.ValueString())
	current, err := client.Read(ctx, desired.Users != nil, desired.Teams != nil)
	if err != nil {
		diags.AddError("Unable to read dataset ACLs", err.Error())
		return
	}
	writeErr := client.Apply(ctx, desired, current)
	refreshed, readErr := client.Read(ctx, desired.Users != nil, desired.Teams != nil)
	if readErr == nil {
		setACLGrants(ctx, &model, refreshed, diags)
	}
	// Retain the identity and managed scope even after a partial failure, so a
	// retry or destroy can reconcile successful changes from the first action.
	diags.Append(state.Set(ctx, &model)...)
	if writeErr != nil {
		diags.AddError("Unable to reconcile dataset ACLs", fmt.Sprintf("%s. User and team ACLs are separate API operations; some changes may have succeeded. Retry to reconcile the remaining grants.", writeErr))
	}
	if readErr != nil {
		diags.AddError("Unable to read updated dataset ACLs", readErr.Error())
	}
}

func (r *searchDatasetAclManagedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model SearchDatasetACLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	grants, err := searchacl.New(r.client, model.DatasetID.ValueString()).Read(ctx, !model.Users.IsNull(), !model.Teams.IsNull())
	if restclient.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read dataset ACLs", err.Error())
		return
	}
	setACLGrants(ctx, &model, grants, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *searchDatasetAclManagedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model SearchDatasetACLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := searchacl.New(r.client, model.DatasetID.ValueString())
	current, err := client.Read(ctx, !model.Users.IsNull(), !model.Teams.IsNull())
	if restclient.IsNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read dataset ACLs before removal", err.Error())
		return
	}
	var empty searchacl.Grants
	if !model.Users.IsNull() {
		empty.Users = map[string]string{}
	}
	if !model.Teams.IsNull() {
		empty.Teams = map[string]string{}
	}
	if err := client.Apply(ctx, empty, current); err != nil && !restclient.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to remove dataset ACLs", err.Error())
	}
}

func (r *searchDatasetAclManagedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" {
		resp.Diagnostics.AddError("Invalid dataset ACL import ID", "Use the existing dataset ID as the import ID.")
		return
	}
	// Non-null maps select both explicit ACL lists on the first import read.
	empty := types.MapValueMust(types.StringType, map[string]attr.Value{})
	model := SearchDatasetACLModel{DatasetID: types.StringValue(req.ID), Users: empty, Teams: empty}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *searchDatasetAclManagedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model SearchDatasetACLModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	grants, err := searchacl.New(d.client, model.DatasetID.ValueString()).Read(ctx, true, true)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read dataset ACLs", err.Error())
		return
	}
	setACLGrants(ctx, &model, grants, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
