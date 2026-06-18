package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCertificate(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_certificate.my_certificate"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "id", "my-demo-cert-001"),
						resource.TestCheckResourceAttr(resourceName, "description", "Demo x509 certificate for Cribl configuration"),
						resource.TestCheckResourceAttr(resourceName, "in_use.#", "0"),
					),
				},
				{
					ConfigDirectory: config.TestNameDirectory(),
				},
				{
					ConfigDirectory: config.TestNameDirectory(),
					PlanOnly:        true,
				},
				{
					ResourceName:            resourceName,
					ImportState:             true,
					ImportStateId:           `{"group_id":"default","id":"my-demo-cert-001"}`,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"priv_key", "passphrase"},
				},
			},
		})
	})
}
