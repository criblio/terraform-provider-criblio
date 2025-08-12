package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLakeDatasetWithLakehouse(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					ImportState:     true,
					ImportStateId:   `{"id": "my_lake_dataset_with_lakehouse_10", "lake_id": "default"}`,
					ResourceName:    "criblio_cribl_lake_dataset.my_cribllakedataset",
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "id", "my_lake_dataset_with_lakehouse_11"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "bucket_name", "lake-tfprovider-beautiful-nguyen-y8y4azd"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "description", "my_description"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "format", "json"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "lake_id", "default"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "retention_period_in_days", "30"),
						resource.TestCheckResourceAttr("criblio_cribl_lake_dataset.my_cribllakedataset", "search_config.0.metadata.0.tags.0", "test_tag"),
					),
				},
			},
		})
	})
}
