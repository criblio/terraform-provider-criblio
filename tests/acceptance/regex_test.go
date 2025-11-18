package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRegex(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "description", "test_regex_2"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "lib", "custom"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "tags", "test"),
						resource.TestCheckResourceAttr("criblio_regex.my_regex", "id", "test_regex_2"),
					),
				},
			},
		})
	})
}
