package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackBreakersGenerated(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack breakers test for on-prem: uses source pack installation")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	packID := "test-pack-breakers-" + suffix
	breakerID := "test_packbreakers_" + suffix

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: packBreakersConfig(packID, breakerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_pack.breakers_pack", "id", packID),
					resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "id", breakerID),
					resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "group_id", "default"),
					resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "pack", packID),
					resource.TestCheckResourceAttr("criblio_pack_breakers.my_packbreakers", "rules.#", "1"),
				),
			},
			{
				Config: packBreakersConfig(packID, breakerID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "criblio_pack_breakers.my_packbreakers",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q,"pack":%q}`, breakerID, packID),
				ImportStateVerify: true,
			},
		},
	})
}

func packBreakersConfig(packID, breakerID string) string {
	return fmt.Sprintf(`resource "criblio_pack" "breakers_pack" {
  id           = %q
  group_id     = "default"
  description  = "Pack breakers"
  disabled     = true
  display_name = "Pack breakers"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_breakers" "my_packbreakers" {
  description    = "test"
  group_id       = "default"
  id             = %q
  lib            = "custom"
  min_raw_length = 256
  pack           = criblio_pack.breakers_pack.id
  tags           = "test"
  rules = [
    {
      condition           = "PASS_THROUGH_SOURCE_TYPE"
      disabled            = false
      event_breaker_regex = "/[\\n\\r]+(?!\\s)/"
      fields              = []
      max_event_bytes     = 51200
      name                = "test"
      parser_enabled      = false
      should_use_data_raw = false
      timestamp = {
        length = 150
        type   = "auto"
      }
      timestamp_anchor_regex = "/^/"
      timestamp_earliest     = "-420weeks"
      timestamp_latest       = "+1week"
      timestamp_timezone     = "local"
      type                   = "regex"
    },
  ]
}
`, packID, breakerID)
}
