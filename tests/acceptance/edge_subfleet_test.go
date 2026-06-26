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

func TestEdgeSubFleet(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping edge subfleet test for on-prem deployments")
	}

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf-edge-subfleet-" + suffix
	resourceName := "criblio_group.my_edge_subfleet"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: edgeSubfleetConfig(groupID, "created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", groupID),
						resource.TestCheckResourceAttr(resourceName, "name", groupID),
						resource.TestCheckResourceAttr(resourceName, "description", "Edge subfleet created"),
						resource.TestCheckResourceAttr(resourceName, "inherits", "default_fleet"),
						resource.TestCheckResourceAttr(resourceName, "is_fleet", "false"),
						resource.TestCheckResourceAttr(resourceName, "product", "edge"),
					),
				},
				{
					Config: edgeSubfleetConfig(groupID, "updated"),
					Check:  resource.TestCheckResourceAttr(resourceName, "description", "Edge subfleet updated"),
				},
				{
					Config: edgeSubfleetConfig(groupID, "updated"),
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

func edgeSubfleetConfig(id, descriptionSuffix string) string {
	return fmt.Sprintf(`resource "criblio_group" "my_edge_subfleet" {
  description          = "Edge subfleet %s"
  id                   = %q
  inherits             = "default_fleet"
  is_fleet             = false
  name                 = %q
  on_prem              = true
  product              = "edge"
  provisioned          = false
  streamtags           = ["terraform", "edge"]
  tags                 = "environment=%s"
  type                 = "edge"
  worker_remote_access = false
}
`, descriptionSuffix, id, id, descriptionSuffix)
}
