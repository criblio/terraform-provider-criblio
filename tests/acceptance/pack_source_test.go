package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackSourceGenerated(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack source test for on-prem: uses source pack installation")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	packID := "test-pack-source-" + suffix
	sourceID := "test_packsource_" + suffix

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: packSourceResourceConfig(packID, sourceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_pack.source_pack", "id", packID),
					resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "id", sourceID),
					resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "group_id", "default"),
					resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "pack", packID),
					resource.TestCheckResourceAttr("criblio_pack_source.my_packsource", "input_tcp.type", "tcp"),
				),
			},
			{
				Config: packSourceResourceConfig(packID, sourceID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "criblio_pack_source.my_packsource",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q,"pack":%q}`, sourceID, packID),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"items",
				},
			},
		},
	})
}

func packSourceResourceConfig(packID, sourceID string) string {
	return fmt.Sprintf(`resource "criblio_pack" "source_pack" {
  id           = %q
  group_id     = "default"
  description  = "Pack with source"
  disabled     = true
  display_name = "Pack with source"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_source" "my_packsource" {
  group_id = "default"
  pack     = criblio_pack.source_pack.id
  id       = %q

  input_tcp = {
    auth_type           = "manual"
    description         = "my_description"
    disabled            = true
    enable_header       = false
    enable_proxy_header = false
    host                = "0.0.0.0"
    id                  = %q
    ip_whitelist_regex  = "/.*/"
    max_active_cxn      = 1000
    pipeline            = "main"
    port                = 55140
    pq_enabled          = false
    send_to_routes      = false
    socket_ending_max_wait = 30
    socket_idle_timeout    = 0
    socket_max_lifespan    = 0
    stale_channel_flush_ms = 1500
    streamtags      = ["tcp"]
    type            = "tcp"
  }
}
`, packID, sourceID, sourceID)
}
