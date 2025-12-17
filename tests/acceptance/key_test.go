package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestKey(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_key.my_key", "id", "key-001"),
						resource.TestCheckResourceAttr("criblio_key.my_key", "description", "My Key Metadata"),
						resource.TestCheckResourceAttr("criblio_key.my_key", "algorithm", "aes-256-cbc"),
						resource.TestCheckResourceAttr("criblio_key.my_key", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_key.my_key", "kms", "local"),
						resource.TestCheckResourceAttr("criblio_key.my_key", "use_iv", "true"),
					),
				},
			},
		})
	})
}
