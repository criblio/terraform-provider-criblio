package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestGlobalVar(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ExpectNonEmptyPlan: true,
					ConfigDirectory:    config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "id", "new"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "description", "test"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "lib", "cribl"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "tags", "cribl,sample"),
					),
				},
			},
		})
	})
}

