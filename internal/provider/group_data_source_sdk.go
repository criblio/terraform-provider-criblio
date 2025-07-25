// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *GroupDataSourceModel) RefreshFromOperationsGetGroupsByIDResponseBody(ctx context.Context, resp *operations.GetGroupsByIDResponseBody) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp != nil {
		r.Items = []tfTypes.Group{}
		if len(r.Items) > len(resp.Items) {
			r.Items = r.Items[:len(resp.Items)]
		}
		for itemsCount, itemsItem := range resp.Items {
			var items tfTypes.Group
			if itemsItem.Cloud == nil {
				items.Cloud = nil
			} else {
				items.Cloud = &tfTypes.Cloud{}
				items.Cloud.Provider = types.StringValue(string(itemsItem.Cloud.Provider))
				items.Cloud.Region = types.StringValue(itemsItem.Cloud.Region)
			}
			items.EstimatedIngestRate = types.Float64PointerValue(itemsItem.EstimatedIngestRate)
			items.ID = types.StringValue(itemsItem.ID)
			items.IsFleet = types.BoolPointerValue(itemsItem.IsFleet)
			items.Name = types.StringPointerValue(itemsItem.Name)
			items.OnPrem = types.BoolPointerValue(itemsItem.OnPrem)
			items.Provisioned = types.BoolValue(itemsItem.Provisioned)
			items.Streamtags = make([]types.String, 0, len(itemsItem.Streamtags))
			for _, v := range itemsItem.Streamtags {
				items.Streamtags = append(items.Streamtags, types.StringValue(v))
			}
			items.WorkerRemoteAccess = types.BoolPointerValue(itemsItem.WorkerRemoteAccess)
			if itemsCount+1 > len(r.Items) {
				r.Items = append(r.Items, items)
			} else {
				r.Items[itemsCount].Cloud = items.Cloud
				r.Items[itemsCount].EstimatedIngestRate = items.EstimatedIngestRate
				r.Items[itemsCount].ID = items.ID
				r.Items[itemsCount].IsFleet = items.IsFleet
				r.Items[itemsCount].Name = items.Name
				r.Items[itemsCount].OnPrem = items.OnPrem
				r.Items[itemsCount].Provisioned = items.Provisioned
				r.Items[itemsCount].Streamtags = items.Streamtags
				r.Items[itemsCount].WorkerRemoteAccess = items.WorkerRemoteAccess
			}
		}
	}

	return diags
}

func (r *GroupDataSourceModel) ToOperationsGetGroupsByIDRequest(ctx context.Context) (*operations.GetGroupsByIDRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var id string
	id = r.ID.ValueString()

	fields := new(string)
	if !r.Fields.IsUnknown() && !r.Fields.IsNull() {
		*fields = r.Fields.ValueString()
	} else {
		fields = nil
	}
	out := operations.GetGroupsByIDRequest{
		ID:     id,
		Fields: fields,
	}

	return &out, diags
}
