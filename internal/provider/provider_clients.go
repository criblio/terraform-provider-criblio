package provider

import (
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

// ProviderClients carries API clients for provider resources and data sources.
type ProviderClients struct {
	RC *restclient.Client
}
