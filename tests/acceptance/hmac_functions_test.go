package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHmacFunctions(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config:             hmacConfig,
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "id", "my_hmacfunction"),
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "description", "test hmac function"),
						//resource.TestCheckResourceAttr("data.criblio_hmac_function.my_hmacfunction", "id", "my_hmacfunction"),
					),
				},
			},
		})
	})
}

const hmacConfig = `
                  resource "criblio_hmac_function" "my_hmacfunction" {
		  description       = "test hmac function"
		  group_id          = "default"
		  header_expression = "'hmac sha256 ' + C.Crypto.createHmac('test', C.Secret('yourSecret','text').value, 'sha256','hex')"
		  header_name       = "signature"
		  id                = "my_hmacfunction"
		  lib               = "cribl"
		  string_builders = [
		    "true"
		  ]
		  string_delim = "true"
		}

		output "hmac_function" {
		  value = criblio_hmac_function.my_hmacfunction
		}

		data "criblio_hmac_function" "my_hmacfunction" {
		  group_id = "default"
		}

		provider "criblio" {
		  server_url = "https://app.cribl-playground.cloud"
		  organization_id = "beautiful-nguyen-y8y4azd"
		  workspace_id = "tfprovider"
		  version = "999.99.9"
		}
		`
