package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestCustomBanner(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping custom banner test for on-prem: global UI banner behavior differs by deployment")
	}

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: customBannerConfig("Scheduled maintenance window: Saturday 2am-4am UTC", "purple"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "enabled", "true"),
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "id", "custom-banner"),
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "message", "Scheduled maintenance window: Saturday 2am-4am UTC"),
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "theme", "purple"),
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "type", "custom"),
					),
				},
				{
					Config: customBannerConfig("Maintenance complete. Systems are operating normally.", "green"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "message", "Maintenance complete. Systems are operating normally."),
						resource.TestCheckResourceAttr("criblio_custom_banner.my_custom_banner", "theme", "green"),
					),
				},
				{
					Config: customBannerConfig("Maintenance complete. Systems are operating normally.", "green"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
				{
					ResourceName:      "criblio_custom_banner.my_custom_banner",
					ImportState:       true,
					ImportStateId:     "custom-banner",
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"created",
						"custom_themes",
						"invert_font_color",
					},
				},
			},
		})
	})
}

func customBannerConfig(message, theme string) string {
	return `resource "criblio_custom_banner" "my_custom_banner" {
  enabled = true
  message = "` + message + `"
  theme   = "` + theme + `"
  type    = "custom"

  link         = "https://status.example.com"
  link_display = "View status page"
}
`
}
