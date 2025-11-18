package tests

import (
	"time"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGrok(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_grok.my_grok", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_grok.my_grok", "id", "test_grok"),
					),
				},
			},
		})
	})
}
