package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchNotificationTarget(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is 'prohibited by current license'")
	}

	id := "tf-search-nt-" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	resourceName := "criblio_notification_target.my_notificationtarget"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: notificationTargetConfig(id, "search"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", id),
						resource.TestCheckResourceAttr(resourceName, "sns_target.id", id),
						resource.TestCheckResourceAttr(resourceName, "sns_target.type", "sns"),
					),
				},
				{
					Config:   notificationTargetConfig(id, "search"),
					PlanOnly: true,
				},
			},
		})
	})
}
