package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProject(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_project.my_project.0", "id", "my_project"),
						resource.TestCheckResourceAttr("criblio_project.my_project.0", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_project.my_project.0", "description", "test project"),
					),
				},
			},
		})
	})
}
