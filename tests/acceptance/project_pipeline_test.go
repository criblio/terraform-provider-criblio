package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestProjectPipeline(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Skipping project pipeline test for on-prem: projects are prohibited by the current license")
	}

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	projectID := "test-project-pipeline-" + suffix
	pipelineID := "test_project_pipeline_" + suffix
	resourceName := "criblio_project_pipeline.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories:  providerFactory,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: projectPipelineConfig(projectID, pipelineID, "created"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_project.test", "id", projectID),
					resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "id", pipelineID),
					resource.TestCheckResourceAttr(resourceName, "conf.description", "created"),
					resource.TestCheckResourceAttr(resourceName, "conf.output", "default"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conf.functions.0.id", "code"),
					resource.TestCheckResourceAttrPair("data.criblio_project_pipeline.by_id", "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair("data.criblio_project_pipeline.by_id", "project_id", resourceName, "project_id"),
					testCheckListDataSourceHasItems("data.criblio_project_pipelines.all"),
				),
			},
			{
				Config: projectPipelineConfig(projectID, pipelineID, "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "conf.description", "updated"),
				),
			},
			{
				Config: projectPipelineConfig(projectID, pipelineID, "updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf(`{"group_id":"default","project_id":%q,"id":%q}`, projectID, pipelineID),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"conf.functions.0.conf"},
			},
		},
	})
}

func projectPipelineConfig(projectID, pipelineID, description string) string {
	return fmt.Sprintf(`resource "criblio_project" "test" {
  group_id     = "default"
  id           = %[1]q
  description  = "Project pipeline acceptance test"
  destinations = []
  subscriptions = []
}

resource "criblio_project_pipeline" "test" {
  group_id   = "default"
  project_id = criblio_project.test.id
  id         = %[2]q
  conf = {
    async_func_timeout = 60
    description        = %[3]q
    output             = "default"
    streamtags         = []
    functions = [
      {
        id       = "code"
        filter   = "true"
        disabled = false
        final    = true
        conf = jsonencode({
          code = "__e.project_pipeline_test = true"
        })
      }
    ]
  }
}

data "criblio_project_pipeline" "by_id" {
  group_id   = criblio_project_pipeline.test.group_id
  project_id = criblio_project_pipeline.test.project_id
  id         = criblio_project_pipeline.test.id
  depends_on = [criblio_project_pipeline.test]
}

data "criblio_project_pipelines" "all" {
  group_id   = criblio_project_pipeline.test.group_id
  project_id = criblio_project_pipeline.test.project_id
  depends_on = [criblio_project_pipeline.test]
}
`, projectID, pipelineID, description)
}
