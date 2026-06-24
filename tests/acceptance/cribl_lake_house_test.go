package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCriblLakeHouse(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	id := "tf_lakehouse_" + suffix
	resourceName := "criblio_cribl_lake_house.my_lakehouse"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: criblLakeHouseConfig(id, "Terraform lakehouse", "small"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform lakehouse"),
					resource.TestCheckResourceAttr(resourceName, "tier_size", "small"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: criblLakeHouseConfig(id, "Terraform lakehouse updated", "small"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "Terraform lakehouse updated"),
			},
			{
				Config:   criblLakeHouseConfig(id, "Terraform lakehouse updated", "small"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           id,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func criblLakeHouseConfig(id, description, tierSize string) string {
	return `resource "criblio_cribl_lake_house" "my_lakehouse" {
  description = "` + description + `"
  id          = "` + id + `"
  tier_size   = "` + tierSize + `"
}
`
}
