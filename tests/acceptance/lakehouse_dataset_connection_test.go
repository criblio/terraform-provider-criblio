package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLakehouseDatasetConnection(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}
	lakehouseID := os.Getenv("CRIBL_TEST_LAKEHOUSE_ID")
	if lakehouseID == "" {
		t.Skip("Set CRIBL_TEST_LAKEHOUSE_ID to a ready lakehouse ID to run this test")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	datasetID := "tf_lake_dataset_conn_" + suffix
	resourceName := "criblio_lakehouse_dataset_connection.my_connection"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: lakehouseDatasetConnectionConfig(lakehouseID, datasetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "lakehouse_id", lakehouseID),
					resource.TestCheckResourceAttr(resourceName, "lake_dataset_id", datasetID),
				),
			},
			{
				Config:   lakehouseDatasetConnectionConfig(lakehouseID, datasetID),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"lakehouse_id":"` + lakehouseID + `","lake_dataset_id":"` + datasetID + `"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func lakehouseDatasetConnectionConfig(lakehouseID, datasetID string) string {
	return `resource "criblio_cribl_lake_dataset" "my_dataset" {
  description              = "Terraform lake dataset connection parent"
  format                   = "json"
  id                       = "` + datasetID + `"
  lake_id                  = "default"
  retention_period_in_days = 30
}

resource "criblio_lakehouse_dataset_connection" "my_connection" {
  lake_dataset_id = criblio_cribl_lake_dataset.my_dataset.id
  lakehouse_id    = "` + lakehouseID + `"
}
`
}
