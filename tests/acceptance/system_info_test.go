package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSystemInfo(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.criblio_system_info.my_systeminfo", "items.#"),
						resource.TestCheckResourceAttrSet("data.criblio_system_info.my_systeminfo", "items.0.build"),
						resource.TestCheckResourceAttrWith("data.criblio_system_info.my_systeminfo", "items.0.build", func(value string) error {
							if value == "" || value == "{}" {
								return fmt.Errorf("build attribute is empty, version cannot be extracted")
							}
							return nil
						}),
					),
				},
			},
		})
	})
}
