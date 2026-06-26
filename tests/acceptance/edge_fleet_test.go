package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestEdgeFleet(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping edge fleet test for on-prem deployments")
	}

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf-edge-fleet-" + suffix
	resourceName := "criblio_group.my_edge_fleet"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: edgeFleetConfig(groupID, "created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", groupID),
						resource.TestCheckResourceAttr(resourceName, "name", groupID),
						resource.TestCheckResourceAttr(resourceName, "description", "Edge fleet created"),
						resource.TestCheckResourceAttr(resourceName, "is_fleet", "true"),
						resource.TestCheckResourceAttr(resourceName, "on_prem", "false"),
						resource.TestCheckResourceAttr(resourceName, "product", "edge"),
					),
				},
				{
					Config: edgeFleetConfig(groupID, "updated"),
					Check:  resource.TestCheckResourceAttr(resourceName, "description", "Edge fleet updated"),
				},
				{
					Config: edgeFleetConfig(groupID, "updated"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     groupID,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"worker_remote_access",
					},
				},
			},
		})
	})
}

func edgeFleetConfig(id, descriptionSuffix string) string {
	return fmt.Sprintf(`resource "criblio_group" "my_edge_fleet" {
  description          = "Edge fleet %s"
  id                   = %q
  is_fleet             = true
  name                 = %q
  on_prem              = false
  product              = "edge"
  provisioned          = false
  streamtags           = ["terraform", "edge"]
  tags                 = "environment=%s"
  type                 = "edge"
  worker_remote_access = false
}
`, descriptionSuffix, id, id, descriptionSuffix)
}
