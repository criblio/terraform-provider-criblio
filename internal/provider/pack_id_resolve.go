package provider

import (
	"context"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/sdk"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/operations"
)

// resolvePackIDForAPI returns the pack ID the API expects. Cribl 4.17.0+ normalizes pack IDs to
// lowercase; pre-4.17.0 packs may retain mixed-case. Does case-insensitive lookup against the packs
// list and returns the actual server ID, or lowercase of configPackID if not found. Used by
// pack-scoped resources (pack_pipeline, pack_source, pack_vars, etc.) to support both behaviors.
func resolvePackIDForAPI(ctx context.Context, client *sdk.CriblIo, groupID, configPackID string) string {
	listReq := operations.GetPacksByGroupRequest{GroupID: groupID}
	listRes, err := client.Packs.GetPacksByGroup(ctx, listReq)
	if err != nil || listRes == nil || listRes.Object == nil {
		return strings.ToLower(configPackID)
	}
	for _, item := range listRes.Object.Items {
		if strings.EqualFold(item.ID, configPackID) {
			return item.ID
		}
	}
	return strings.ToLower(configPackID)
}
