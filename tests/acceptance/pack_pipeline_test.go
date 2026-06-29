package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackPipelineGenerated(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping pack pipeline test for on-prem: uses source pack installation")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	packID := "test-pack-pipeline-" + suffix
	pipelineID := "test_packpipeline_" + suffix

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config: packPipelineConfig(packID, pipelineID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "id", packID),
					resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "id", pipelineID),
					resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "group_id", "default"),
					resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "pack", packID),
					resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.#", "1"),
				),
			},
			{
				Config: packPipelineConfig(packID, pipelineID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "criblio_pack_pipeline.my_packpipeline",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf(`{"group_id":"default","id":%q,"pack":%q}`, pipelineID, packID),
				ImportStateVerify: true,
			},
		},
	})
}

func packPipelineConfig(packID, pipelineID string) string {
	return fmt.Sprintf(`resource "criblio_pack" "pipeline_pack" {
  id           = %q
  group_id     = "default"
  description  = "Pack with pipeline"
  disabled     = true
  display_name = "Pack with pipeline"
  source       = "file:/opt/cribl_data/failover/groups/default/default/HelloPacks"
  version      = "1.0.0"
}

resource "criblio_pack_pipeline" "my_packpipeline" {
  group_id = "default"
  id       = %q
  pack     = criblio_pack.pipeline_pack.id
  conf = {
    async_func_timeout = 9066
    description        = "my_description"
    functions = [
      {
        id     = "eval"
        filter = "true"
        conf = jsonencode({
          add = [
            {
              disabled = false
              name     = "_value"
              value    = "1"
            }
          ]
        })
      },
    ]
    groups = {
      default = {
        name = "default"
      }
    }
    streamtags = [
      "tags"
    ]
  }
}
`, packID, pipelineID)
}
