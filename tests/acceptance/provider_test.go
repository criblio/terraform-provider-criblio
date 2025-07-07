package tests

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/speakeasy/terraform-provider-criblio/internal/provider"
)

var (
	providerFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(provider.New("999.99.9")()),
	}

	providerConfigUs = `provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "main"
  version = "999.99.9"
}
`
)

func TestFoo (t *testing.T) {
	fmt.Println("client_id = ", os.Getenv("CRIBL_CLIENT_ID"))
	fmt.Println("client_secret = ", os.Getenv("CRIBL_CLIENT_SECRET"))
	fmt.Println("bearer_token = ", os.Getenv("CRIBL_BEARER_TOKEN"))
}

