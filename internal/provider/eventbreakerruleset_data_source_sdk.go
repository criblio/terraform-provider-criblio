// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (r *EventBreakerRulesetDataSourceModel) RefreshFromOperationsListEventBreakerRulesetResponseBody(ctx context.Context, resp *operations.ListEventBreakerRulesetResponseBody) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp != nil {
	}

	return diags
}

func (r *EventBreakerRulesetDataSourceModel) ToOperationsListEventBreakerRulesetRequest(ctx context.Context) (*operations.ListEventBreakerRulesetRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupID string
	groupID = r.GroupID.ValueString()

	out := operations.ListEventBreakerRulesetRequest{
		GroupID: groupID,
	}

	return &out, diags
}
