package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLakeDestination(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "id", "cribl-lake-11"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.id", "cribl-lake-11"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.type", "cribl_lake"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.description", "Cribl Lake destination"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.dest_path", "default_logs"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.format", "json"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.compress", "gzip"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.base_file_name", "CriblOut"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.max_file_size_mb", "32"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.max_open_files", "100"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.write_high_water_mark", "64"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.on_backpressure", "block"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.max_file_open_time_sec", "300"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.max_file_idle_time_sec", "30"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "output_cribl_lake.0.max_retry_num", "20"),
					),
				},
			},
		})
	})
}
