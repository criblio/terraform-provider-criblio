package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCriblLakeDataset(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	id := "tf_lake_dataset_" + suffix
	resourceName := "criblio_cribl_lake_dataset.my_dataset"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: criblLakeDatasetConfig(id, "Terraform lake dataset", "tf_tag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "lake_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform lake dataset"),
					resource.TestCheckResourceAttr(resourceName, "format", "json"),
					resource.TestCheckResourceAttr(resourceName, "retention_period_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "search_config.metadata.tags.0", "tf_tag"),
				),
			},
			{
				Config: criblLakeDatasetConfig(id, "Terraform lake dataset updated", "tf_tag_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform lake dataset updated"),
					resource.TestCheckResourceAttr(resourceName, "search_config.metadata.tags.0", "tf_tag_updated"),
				),
			},
			{
				Config:   criblLakeDatasetConfig(id, "Terraform lake dataset updated", "tf_tag_updated"),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"lake_id":"default","id":"` + id + `"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func criblLakeDatasetConfig(id, description, tag string) string {
	return `resource "criblio_cribl_lake_dataset" "my_dataset" {
  description              = "` + description + `"
  format                   = "json"
  id                       = "` + id + `"
  lake_id                  = "default"
  retention_period_in_days = 30
  search_config = {
    metadata = {
      tags = ["` + tag + `"]
    }
  }
}
`
}
