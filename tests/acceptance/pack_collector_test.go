package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackCollector(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Uses the cloud HelloPacks installation fixture")
	}
	for _, kind := range []string{"rest", "cribl_lake"} {
		t.Run(kind, func(t *testing.T) {
			suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
			packID := "test-pack-collector-" + suffix
			id := "collector-" + suffix
			name := "criblio_pack_collector.test"
			block := "input_collector_" + kind
			base := packCollectorDependencies(packID, id, kind)
			config := func(ttl string) string {
				return base + packCollectorConfig(id, kind, ttl)
			}
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: providerFactory,
				Steps: []resource.TestStep{
					{
						Config: config("4h"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(name, "id", id),
							resource.TestCheckResourceAttr(name, "group_id", "default"),
							resource.TestCheckResourceAttr(name, "pack", packID),
							resource.TestCheckResourceAttr(name, block+".collector.type", kind),
							resource.TestCheckResourceAttr(name, block+".ttl", "4h"),
							resource.TestCheckResourceAttrPair("data.criblio_pack_collector.test", "id", name, "id"),
							resource.TestCheckResourceAttr("data.criblio_pack_collector.test", block+".collector.type", kind),
						),
					},
					{
						Config: config("8h"),
						Check:  resource.TestCheckResourceAttr(name, block+".ttl", "8h"),
					},
					{
						Config: config("8h"),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
						},
					},
					{
						ResourceName:      name,
						ImportState:       true,
						ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q,"pack":%q}`, id, packID),
						ImportStateVerify: true,
					},
					// Delete the collector while its pack still exists.
					{Config: base},
				},
			})
		})
	}
}

func packCollectorDependencies(packID, id, kind string) string {
	config := fmt.Sprintf(`
resource "criblio_pack" "test" {
  id           = %q
  group_id     = "default"
  disabled     = true
  display_name = "Pack collector acceptance test"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}
`, packID)
	if kind == "cribl_lake" {
		config += criblLakeDatasetConfig(id, "Pack collector acceptance test", "pack-collector")
	}
	return config
}

func packCollectorConfig(id, kind, ttl string) string {
	conf := `collect_url = "'https://example.com/events'"
        collect_method = "get"
        authentication = "none"`
	if kind == "cribl_lake" {
		conf = `dataset = criblio_cribl_lake_dataset.my_dataset.id`
	}
	return fmt.Sprintf(`
resource "criblio_pack_collector" "test" {
  group_id = "default"
  pack     = criblio_pack.test.id
  id       = %q
  input_collector_%s = {
    ttl = %q
    collector = {
      type = %q
      conf = {
        %s
      }
    }
    input = {
      type           = "collection"
      send_to_routes = true
    }
  }
}

data "criblio_pack_collector" "test" {
  group_id = criblio_pack_collector.test.group_id
  pack     = criblio_pack_collector.test.pack
  id       = criblio_pack_collector.test.id
}
`, id, kind, ttl, kind, conf)
}
