package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHmacFunction(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_hmac_function.my_hmac_function"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hmacFunctionConfig("test hmac function"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "my_hmac_function_test"),
					resource.TestCheckResourceAttr(resourceName, "description", "test hmac function"),
					resource.TestCheckResourceAttr(resourceName, "header_name", "signature"),
				),
			},
			{
				Config: hmacFunctionConfig("updated hmac function"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "updated hmac function"),
			},
			{
				Config:   hmacFunctionConfig("updated hmac function"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"my_hmac_function_test"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func hmacFunctionConfig(description string) string {
	return `resource "criblio_hmac_function" "my_hmac_function" {
  description       = "` + description + `"
  group_id          = "default"
  header_expression = "'hmac sha256 ' + C.Crypto.createHmac('test', C.Secret('yourSecret','text').value, 'sha256','hex')"
  header_name       = "signature"
  id                = "my_hmac_function_test"
  lib               = "custom"
  string_builders   = ["true"]
  string_delim      = "true"
}
`
}
