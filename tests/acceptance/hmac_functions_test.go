package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHmacFunctions(t *testing.T) {
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
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "id", "my_hmacfunction_test"),
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "description", "test hmac function"),
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_hmac_function.my_hmacfunction", "header_name", "signature"),
					),
				},
			},
		})
	})
}
