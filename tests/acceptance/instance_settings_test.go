package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestInstanceSettings(t *testing.T) {
	if os.Getenv("DEPLOYMENT") != "onprem" {
		t.Skip("Skipping data source for Cloud deployments as it is not supported")
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.criblio_instance_settings.my_instancesettings", "items.#"),
					),
				},
			},
		})
	})
}
