package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLakeDestination(t *testing.T) {
	t.Skip("Skipping lake destination due to RPC timeouts")
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "id", "cribl-lake-11"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "group_id", "default"),
					),
				},
			},
		})
	})
}
