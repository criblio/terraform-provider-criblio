package tests

import (
	"time"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGlobalVar(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(2 * time.Second)
        }

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "id", "sample_globalvar"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "description", "test"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "lib", "test"),
						resource.TestCheckResourceAttr("criblio_global_var.my_globalvar", "tags", "test"),
					),
				},
			},
		})
	})
}
