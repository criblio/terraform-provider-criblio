package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSearchDatatype(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	id := "tf_search_datatype_" + acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	resourceName := "criblio_search_datatype.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: searchDatatypeConfig(id, 65536),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "lib", "custom"),
					resource.TestCheckResourceAttr(resourceName, "data_format", "ndjson"),
					resource.TestCheckResourceAttr(resourceName, "max_event_bytes", "65536"),
					resource.TestCheckResourceAttr(resourceName, "min_raw_length", "0"),
					resource.TestCheckResourceAttr(resourceName, "search_version", "v2"),
					resource.TestCheckResourceAttr(resourceName, "automatic_extraction.extraction_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "timestamp_extraction.type", "auto"),
					resource.TestCheckResourceAttr(resourceName, "timestamp_extraction.scan_depth", "150"),
				),
			},
			{
				Config: searchDatatypeConfig(id, 131072),
				Check:  resource.TestCheckResourceAttr(resourceName, "max_event_bytes", "131072"),
			},
			{
				Config:   searchDatatypeConfig(id, 131072),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     id,
				ImportStateVerify: true,
			},
		},
	})
}

func searchDatatypeConfig(id string, maxEventBytes int) string {
	return fmt.Sprintf(`
resource "criblio_search_datatype" "example" {
  id              = %q
  lib             = "custom"
  data_format     = "ndjson"
  max_event_bytes = %d
  min_raw_length  = 0
  search_version  = "v2"
  tags            = ""

  automatic_extraction = {
    extraction_type = "json"
  }

  timestamp_extraction = {
    type         = "auto"
    anchor_regex = "/^/"
    timezone     = "UTC"
    earliest     = "0"
    latest       = "+10years"
    scan_depth   = 150
  }
}
`, id, maxEventBytes)
}
