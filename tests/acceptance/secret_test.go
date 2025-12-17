package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSecret(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_secret.my_secret", "id", "sec-test-001"),
						resource.TestCheckResourceAttr("criblio_secret.my_secret", "description", "API key for ingestion service"),
						resource.TestCheckResourceAttr("criblio_secret.my_secret", "secret_type", "text"),
						resource.TestCheckResourceAttr("criblio_secret.my_secret", "group_id", "default"),
					),
				},
			},
		})
	})
}
