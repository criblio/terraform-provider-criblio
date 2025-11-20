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
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:           config.TestNameDirectory(),
					ExpectNonEmptyPlan:        true,
					PreventPostDestroyRefresh: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_certificate.my_certificate", "id", "my-demo-cert-001"),
						resource.TestCheckResourceAttr("criblio_certificate.my_certificate", "description", "Demo x509 certificate for Cribl configuration"),
						resource.TestCheckResourceAttr("criblio_certificate.my_certificate", "in_use.0", "wef-prod"),
					),
				},
			},
		})
	})
}
