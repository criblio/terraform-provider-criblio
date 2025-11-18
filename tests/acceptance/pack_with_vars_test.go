package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPackVars(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(2 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.vars_pack", "id", "pack-with-vars"),
						resource.TestCheckResourceAttr("criblio_pack.vars_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack_vars.my_packvars", "id", "my_id"),
						resource.TestCheckResourceAttr("criblio_pack_vars.my_packvars", "value", "my_value"),
					),
				},
			},
		})
	})
}
