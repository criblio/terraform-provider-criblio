package provider

import (
	"context"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

type packIDListResponse struct {
	ID string `json:"id"`
}

// resolvePackIDForRestAPI returns the pack ID the API expects. Cribl 4.17.0+ normalizes
// pack IDs to lowercase; pre-4.17.0 packs may retain mixed case. The API paths are
// case-sensitive, so pack-scoped resources do a case-insensitive list scan and use the
// server's actual ID when present.
func resolvePackIDForRestAPI(ctx context.Context, client *restclient.Client, groupID, configPackID string) string {
	if client == nil {
		return strings.ToLower(configPackID)
	}
	items, err := restclient.Get[[]packIDListResponse](ctx, client, "/m/"+groupID+"/packs")
	if err != nil || items == nil {
		return strings.ToLower(configPackID)
	}
	for _, item := range *items {
		if strings.EqualFold(item.ID, configPackID) {
			return item.ID
		}
	}
	return strings.ToLower(configPackID)
}
