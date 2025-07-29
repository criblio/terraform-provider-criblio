package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDatasetProvider(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:    config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_searchdatasetprovider", "id", "test_http"),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_elastic_provider", "id", "test_elastic"),
						resource.TestCheckResourceAttr("criblio_search_dataset_provider.my_s3_provider", "id", "test_s3"),
					),
				},
			},
		})
	})
}
