// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *CommitResourceModel) RefreshFromOperationsCreateVersionCommitResponseBody(ctx context.Context, resp *operations.CreateVersionCommitResponseBody) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp != nil {
		r.Items = []tfTypes.GitCommitSummary{}
		if len(r.Items) > len(resp.Items) {
			r.Items = r.Items[:len(resp.Items)]
		}
		for itemsCount, itemsItem := range resp.Items {
			var items tfTypes.GitCommitSummary
			items.Author.Email = types.StringValue(itemsItem.Author.Email)
			items.Author.Name = types.StringValue(itemsItem.Author.Name)
			items.Branch = types.StringValue(itemsItem.Branch)
			items.Commit = types.StringValue(itemsItem.Commit)
			items.Files.Created = make([]types.String, 0, len(itemsItem.Files.Created))
			for _, v := range itemsItem.Files.Created {
				items.Files.Created = append(items.Files.Created, types.StringValue(v))
			}
			items.Files.Deleted = make([]types.String, 0, len(itemsItem.Files.Deleted))
			for _, v := range itemsItem.Files.Deleted {
				items.Files.Deleted = append(items.Files.Deleted, types.StringValue(v))
			}
			items.Files.Modified = make([]types.String, 0, len(itemsItem.Files.Modified))
			for _, v := range itemsItem.Files.Modified {
				items.Files.Modified = append(items.Files.Modified, types.StringValue(v))
			}
			items.Files.Renamed = make([]types.String, 0, len(itemsItem.Files.Renamed))
			for _, v := range itemsItem.Files.Renamed {
				items.Files.Renamed = append(items.Files.Renamed, types.StringValue(v))
			}
			items.Summary.Changes = types.Float64Value(itemsItem.Summary.Changes)
			items.Summary.Deletions = types.Float64Value(itemsItem.Summary.Deletions)
			items.Summary.Insertions = types.Float64Value(itemsItem.Summary.Insertions)
			if itemsCount+1 > len(r.Items) {
				r.Items = append(r.Items, items)
			} else {
				r.Items[itemsCount].Author = items.Author
				r.Items[itemsCount].Branch = items.Branch
				r.Items[itemsCount].Commit = items.Commit
				r.Items[itemsCount].Files = items.Files
				r.Items[itemsCount].Summary = items.Summary
			}
		}
	}

	return diags
}

func (r *CommitResourceModel) ToSharedGitCommitParams(ctx context.Context) (*shared.GitCommitParams, diag.Diagnostics) {
	var diags diag.Diagnostics

	effective := new(bool)
	if !r.Effective.IsUnknown() && !r.Effective.IsNull() {
		*effective = r.Effective.ValueBool()
	} else {
		effective = nil
	}
	files := make([]string, 0, len(r.Files))
	for _, filesItem := range r.Files {
		files = append(files, filesItem.ValueString())
	}
	group := new(string)
	if !r.Group.IsUnknown() && !r.Group.IsNull() {
		*group = r.Group.ValueString()
	} else {
		group = nil
	}
	var message string
	message = r.Message.ValueString()

	out := shared.GitCommitParams{
		Effective: effective,
		Files:     files,
		Group:     group,
		Message:   message,
	}

	return &out, diags
}
