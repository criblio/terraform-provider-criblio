// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
)

func (r *SubscriptionDataSourceModel) ToOperationsListSubscriptionRequest(ctx context.Context) (*operations.ListSubscriptionRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupID string
	groupID = r.GroupID.ValueString()

	out := operations.ListSubscriptionRequest{
		GroupID: groupID,
	}

	return &out, diags
}

func (r *SubscriptionDataSourceModel) RefreshFromSharedSubscription(ctx context.Context, resp *shared.Subscription) diag.Diagnostics {
	var diags diag.Diagnostics

	r.Description = types.StringPointerValue(resp.Description)
	r.Disabled = types.BoolPointerValue(resp.Disabled)
	r.Filter = types.StringPointerValue(resp.Filter)
	r.ID = types.StringValue(resp.ID)
	r.Pipeline = types.StringValue(resp.Pipeline)

	return diags
}
