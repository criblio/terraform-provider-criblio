package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackDestinationGenerated(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack destination test for on-prem: uses source pack installation")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	packID := "test-pack-dest-" + suffix
	destinationID := "test_packdest_" + suffix

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: packDestinationConfig(packID, destinationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_pack.destination_pack", "id", packID),
					resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "id", destinationID),
					resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "group_id", "default"),
					resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "pack", packID),
					resource.TestCheckResourceAttr("criblio_pack_destination.my_packdest", "output_devnull.type", "devnull"),
				),
			},
			{
				Config: packDestinationConfig(packID, destinationID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "criblio_pack_destination.my_packdest",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q,"pack":%q}`, destinationID, packID),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"items",
				},
			},
		},
	})
}

func packDestinationConfig(packID, destinationID string) string {
	return fmt.Sprintf(`resource "criblio_pack" "destination_pack" {
  id           = %q
  group_id     = "default"
  description  = "Pack with destination"
  disabled     = true
  display_name = "Pack with destination"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_destination" "my_packdest" {
  pack     = criblio_pack.destination_pack.id
  group_id = "default"
  id       = %q
  output_devnull = {
    id   = %q
    type = "devnull"
  }
}
`, packID, destinationID, destinationID)
}
