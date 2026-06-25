package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNotification(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	resourceName := "criblio_notification.my_notification"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: notificationConfig(false, "60s"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_notification"),
					resource.TestCheckResourceAttr(resourceName, "condition", "high-volume"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "group", "default"),
					resource.TestCheckResourceAttr(resourceName, "conf.time_window", "60s"),
				),
			},
			{
				Config: notificationConfig(true, "120s"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "conf.time_window", "120s"),
				),
			},
			{
				Config:   notificationConfig(true, "120s"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "test_notification",
				ImportStateVerify: true,
			},
		},
	})
}

func notificationConfig(disabled bool, timeWindow string) string {
	disabledValue := "false"
	if disabled {
		disabledValue = "true"
	}

	return `resource "criblio_notification" "my_notification" {
  condition = "high-volume"
  disabled  = ` + disabledValue + `
  group     = "default"
  id        = "test_notification"

  conf = {
    name                 = "cribl_http:test_search_source"
    time_window          = "` + timeWindow + `"
    data_volume          = "1GB"
    notify_on_resolution = true
  }

  targets = []
}
`
}
