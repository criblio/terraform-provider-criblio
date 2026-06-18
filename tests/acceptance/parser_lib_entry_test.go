package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestParserLibEntry(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping resource for On-Prem deployments as it is not supported")
	}

	resourceName := "criblio_parser_lib_entry.my_parser"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: parserLibEntryConfig("Parser for delimited logs", "delim"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_id", "default_search"),
					resource.TestCheckResourceAttr(resourceName, "id", "test_parser_lib_entry_phase2"),
					resource.TestCheckResourceAttr(resourceName, "description", "Parser for delimited logs"),
					resource.TestCheckResourceAttr(resourceName, "type", "delim"),
				),
			},
			{
				Config: parserLibEntryConfig("Parser for csv logs", "csv"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Parser for csv logs"),
					resource.TestCheckResourceAttr(resourceName, "type", "csv"),
				),
			},
			{Config: parserLibEntryConfig("Parser for csv logs", "csv"), PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     `{"group_id":"default_search","id":"test_parser_lib_entry_phase2"}`,
				ImportStateVerify: true,
			},
		},
	})
}

func parserLibEntryConfig(description, parserType string) string {
	return `resource "criblio_parser_lib_entry" "my_parser" {
  description = "` + description + `"
  group_id    = "default_search"
  id          = "test_parser_lib_entry_phase2"
  lib         = "custom"
  tags        = "phase2"
  type        = "` + parserType + `"
}
`
}
