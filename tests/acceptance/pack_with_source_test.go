package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPackSource(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack test for on-prem: uses HelloPacks which causes API 500 errors")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.source_pack", "id", "pack-with-source"),
						resource.TestCheckResourceAttr("criblio_pack.source_pack", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_pack.source_pack", "description", "Pack with source"),
						resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "input_tcp.type", "tcp"),
						resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "group_id", "default"),
					),
				},
			},
		})
	})
}
