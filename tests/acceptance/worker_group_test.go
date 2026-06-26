package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestWorkerGroup(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping cloud worker group test for on-prem deployments")
	}

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))
	groupID := "tf-stream-group-" + suffix
	resourceName := "criblio_group.azure_eastus_stream_group"

	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: workerGroupConfig(groupID, "created", "created"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", groupID),
						resource.TestCheckResourceAttr(resourceName, "name", groupID),
						resource.TestCheckResourceAttr(resourceName, "description", "Worker group created"),
						resource.TestCheckResourceAttr(resourceName, "product", "stream"),
						resource.TestCheckResourceAttr(resourceName, "type", "stream"),
					),
				},
				{
					Config: workerGroupConfig(groupID, "updated", "updated"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "Worker group updated"),
						resource.TestCheckResourceAttr(resourceName, "tags", "environment=updated"),
					),
				},
				{
					Config: workerGroupConfig(groupID, "updated", "updated"),
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
						"cloud",
						"estimated_ingest_rate",
						"provisioned",
						"worker_remote_access",
					},
				},
			},
		})
	})
}

func TestWorkerGroupRejectsEdgeCloudConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config:      edgeCloudGroupConfig("tf-edge-cloud-invalid"),
				ExpectError: regexp.MustCompile("cloud configuration is only allowed for stream product, not edge"),
			},
		},
	})
}

func workerGroupConfig(id, descriptionSuffix, tagSuffix string) string {
	return fmt.Sprintf(`resource "criblio_group" "azure_eastus_stream_group" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
  description           = "Worker group %s"
  estimated_ingest_rate = 1024
  id                    = %q
  is_fleet              = false
  name                  = %q
  on_prem               = false
  product               = "stream"
  provisioned           = false
  streamtags            = ["terraform", %q]
  tags                  = "environment=%s"
  type                  = "stream"
  worker_remote_access  = false
}
`, descriptionSuffix, id, id, tagSuffix, tagSuffix)
}

func edgeCloudGroupConfig(id string) string {
	return fmt.Sprintf(`resource "criblio_group" "edge_cloud_invalid" {
  cloud = {
    provider = "aws"
    region   = "us-east-1"
  }
  id       = %q
  name     = %q
  on_prem  = false
  product  = "edge"
  type     = "edge"
}
`, id, id)
}
