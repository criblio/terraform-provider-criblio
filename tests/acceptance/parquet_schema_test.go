package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestParquetSchema(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	resourceName := "criblio_parquet_schema.my_parquet_schema"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: parquetSchemaConfig("Phase 2 parquet schema", "STRING"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "id", "phase2_parquet_schema"),
					resource.TestCheckResourceAttr(resourceName, "description", "Phase 2 parquet schema"),
				),
			},
			{
				Config: parquetSchemaConfig("Phase 2 parquet schema updated", "BYTE_ARRAY"),
				Check:  resource.TestCheckResourceAttr(resourceName, "description", "Phase 2 parquet schema updated"),
			},
			{Config: parquetSchemaConfig("Phase 2 parquet schema updated", "BYTE_ARRAY"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default","id":"phase2_parquet_schema"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func parquetSchemaConfig(description, fieldType string) string {
	return `resource "criblio_parquet_schema" "my_parquet_schema" {
  description = "` + description + `"
  group_id    = "default"
  id          = "phase2_parquet_schema"
  schema      = jsonencode({ message = { type = "` + fieldType + `" } })
}
`
}
