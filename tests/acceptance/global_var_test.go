package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGlobalVar(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_global_var.my_globalvar"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "sample_globalvar"),
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "description", "test"),
						resource.TestCheckResourceAttr(resourceName, "lib", "test"),
						resource.TestCheckResourceAttr(resourceName, "tags", "test"),
						resource.TestCheckResourceAttr(resourceName, "args.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "args.0.name", "val"),
						resource.TestCheckResourceAttr(resourceName, "args.0.type", "number"),
					),
				},
				{
					Config: globalVarUpdatedConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "test updated"),
						resource.TestCheckResourceAttr(resourceName, "value", "200"),
						resource.TestCheckResourceAttr(resourceName, "args.#", "1"),
					),
				},
				{
					Config:   globalVarUpdatedConfig,
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default","id":"sample_globalvar"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

const globalVarUpdatedConfig = `resource "criblio_global_var" "my_globalvar" {
  args = [
    {
      name = "val"
      type = "number"
    }
  ]
  description = "test updated"
  group_id    = "default"
  id          = "sample_globalvar"
  lib         = "test"
  tags        = "test"
  type        = "number"
  value       = "200"
}
`
