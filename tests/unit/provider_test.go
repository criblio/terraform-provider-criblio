package tests

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	//this should be criblio instead of speakeasy for the pathing
	"github.com/speakeasy/terraform-provider-criblio/internal/provider"
)

var (
	providerFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(provider.New("999.99.9")()),
	}

	providerConfigUs = `provider "criblio" {
  server_url = "https://app.cribl-playground.cloud/organizations/beautiful-nguyen-y8y4azd/workspaces/tfprovider/app/api/v1"
}
`
)
