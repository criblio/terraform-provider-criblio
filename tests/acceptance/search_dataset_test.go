package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDataset(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: searchDatasetConfig("tf_acc_s3_dataset", "test", "true"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "id", "tf_acc_s3_dataset"),
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "type", "s3"),
					),
				},
				{
					Config: searchDatasetConfig("tf_acc_s3_dataset", "test updated", "true"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "description", "test updated"),
					),
				},
				{Config: searchDatasetConfig("tf_acc_s3_dataset", "test updated", "true"), PlanOnly: true},
				{
					ResourceName:      "criblio_search_dataset.my_s3_dataset",
					ImportState:       true,
					ImportStateId:     "tf_acc_s3_dataset",
					ImportStateVerify: true,
				},
			},
		})
	})
}

func searchDatasetConfig(id, description, filter string) string {
	return `
resource "criblio_search_dataset" "my_s3_dataset" {
  s3_dataset = {
    auto_detect_region = false
    bucket             = "test_bucket"
    description        = "` + description + `"
    extra_paths = [
      {
        auto_detect_region = false
        bucket             = "test_bucket"
        filter             = "` + filter + `"
        path               = "logs/*.log"
        region             = "us-east-1"
      }
    ]
    filter = "` + filter + `"
    id     = "` + id + `"
    metadata = {
      enable_acceleration = false
    }
    path        = "logs/*.log"
    provider_id = "S3"
    region      = "us-east-1"
    storage_classes = [
      "STANDARD"
    ]
    type = "s3"
  }
}
`
}
