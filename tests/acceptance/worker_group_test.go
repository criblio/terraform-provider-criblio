package tests

import (
	"os"
	"time"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkerGroup(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(2 * time.Second)
        }
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ImportState:     true,
					ImportStateId:   "criblio_group.azure_eastus_stream_group",
					ResourceName:    "criblio_group.azure_eastus_stream_group",
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.azure_eastus_stream_group", "id", "my-group"),
						resource.TestCheckResourceAttr("criblio_group.azure_eastus_stream_group", "name", "my-group"),
						resource.TestCheckResourceAttr("criblio_group.azure_eastus_stream_group", "product", "stream"),
					),
				},
			},
		})
	})
}
