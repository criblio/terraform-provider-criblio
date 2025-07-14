package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func TestParquetSchemas(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ConfigDirectory:         config.TestNameDirectory(),
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_parquet_schema.my_parquet_schema", "description", "test parquest"),
						resource.TestCheckResourceAttr("criblio_parquet_schema.my_parquet_schema", "id", "my_parquet_schema"),
						resource.TestCheckResourceAttr("criblio_parquet_schema.my_parquet_schema", "description", "test parquest"),
					),
				},
			},
		})
	})
}

