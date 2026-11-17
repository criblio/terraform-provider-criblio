package provider

import (
	"context"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

// resolvePackIDForRestAPI returns the pack ID the API expects. Cribl 4.17.0+ normalizes
// pack IDs to lowercase; pre-4.17.0 packs may retain mixed case. The API paths are
// case-sensitive, so pack-scoped resources do a case-insensitive list scan and use the
// server's actual ID when present. Successful listings are cached per client and
// group for one minute and shared by concurrent lookups.
func resolvePackIDForRestAPI(ctx context.Context, client *restclient.Client, groupID, configPackID string) string {
	if client == nil {
		return strings.ToLower(configPackID)
	}
	return client.ResolvePackID(ctx, groupID, configPackID)
}
