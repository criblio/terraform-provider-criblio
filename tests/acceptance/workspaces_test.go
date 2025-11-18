package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkspaces(t *testing.T) {
	t.Skip("Skipping workspaces test due to gateway endpoint inconsistency and schema changes")
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.criblio_workspaces.my_workspaces", "items.#"),
					),
				},
			},
		})
	})
}
