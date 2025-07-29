package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDataset(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "id", "S3"),
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "description", "test"),
						resource.TestCheckResourceAttr("criblio_search_dataset.my_s3_dataset", "type", "s3"),
					),
				},
			},
		})
	})
}
