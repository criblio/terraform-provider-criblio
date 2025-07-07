package tests

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/speakeasy/terraform-provider-criblio/internal/provider"
)

var (
	providerFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(provider.New("999.99.9")()),
	}
)
