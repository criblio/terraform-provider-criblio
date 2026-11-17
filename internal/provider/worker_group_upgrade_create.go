package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// createWorkerGroupUpgrade submits the action and preserves its configured inputs.
// The returned job IDs acknowledge submission, not completion of the upgrade.
func createWorkerGroupUpgrade(ctx context.Context, client *restclient.Client, model WorkerGroupUpgradeModel) (*WorkerGroupUpgradeModel, error) {
	if model.GroupID.IsNull() || model.GroupID.IsUnknown() || model.GroupID.ValueString() == "" {
		return nil, fmt.Errorf("group_id must identify an existing self-managed Worker Group")
	}
	jobs, err := restclient.Post[WorkerGroupUpgradeModel, []string](ctx, client,
		"/system/distributed/upgrade/"+url.PathEscape(model.GroupID.ValueString()), model)
	if err != nil {
		var httpErr *restclient.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusBadRequest {
			var response struct {
				Message string `json:"message"`
			}
			if json.Unmarshal([]byte(httpErr.Body), &response) == nil {
				message := strings.ToLower(response.Message)
				if strings.Contains(message, "managed group") || strings.Contains(message, "managed worker group") {
					return nil, fmt.Errorf("worker group %q is managed by Cribl and cannot be upgraded with criblio_worker_group_upgrade. Set group_id to an existing self-managed Worker Group, including a hybrid group. Cribl manages upgrades for cloud-managed groups: %w", model.GroupID.ValueString(), err)
				}
			}
		}
		return nil, err
	}
	values := make([]attr.Value, 0)
	if jobs != nil {
		for _, job := range *jobs {
			values = append(values, types.StringValue(job))
		}
	}
	model.JobIDs = types.ListValueMust(types.StringType, values)
	return &model, nil
}
