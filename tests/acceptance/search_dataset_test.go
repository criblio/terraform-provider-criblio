package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDatasetV1(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	const (
		id           = "tf_acc_s3_dataset_v1"
		resourceName = "criblio_search_dataset.s3_v1"
	)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: searchDatasetV1Config(id, "test", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.search_version", "v1"),
				),
			},
			{
				Config: searchDatasetV1Config(id, "test updated", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "test updated"),
				),
			},
			{Config: searchDatasetV1Config(id, "test updated", "true"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     id,
				ImportStateVerify: true,
			},
		},
	})
}

func TestSearchDatasetV2(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	const (
		id           = "search_dataset_v2"
		resourceName = "criblio_search_dataset.v2"
	)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("../../examples/search-dataset"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "description", "Search v2 dataset"),
					resource.TestCheckResourceAttr(resourceName, "type", "s3"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.search_version", "v2"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.paths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.paths.0.bucket", "lake-main-beautiful-nguyen-y8y4azd"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.paths.0.filters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.paths.0.filters.0.data_type_id", "generic_ndjson"),
					resource.TestCheckResourceAttr(resourceName, "s3_dataset.paths.0.filters.0.filter", "**"),
				),
			},
			{ConfigDirectory: config.StaticDirectory("../../examples/search-dataset"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     id,
				ImportStateVerify: true,
			},
		},
	})
}

func searchDatasetV1Config(id, description, filter string) string {
	return `
resource "criblio_search_dataset" "s3_v1" {
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
    search_version = "v1"
    storage_classes = [
      "STANDARD"
    ]
    type = "s3"
  }
}
`
}
