package provider

import (
	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
)

// ProviderClients carries both API clients during the Speakeasy migration.
type ProviderClients struct {
	Legacy *sdk.CriblIo
	RC     *restclient.Client
}
