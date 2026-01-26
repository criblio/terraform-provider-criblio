package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestPackPipeline(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_pack.pipeline_pack", "id", "pack-with-pipeline"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "id", "my_id"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.0.id", "serde"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.0.filter", "true"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.0.conf", 
						  "{\"mode\":\"extract\",\"src_field\":\"_raw\",\"type\":\"json\"}"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.1.id", "eval"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.1.filter", "true"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.1.conf", 
						  "{\"add\":[{\"disabled\":false,\"name\":\"_value\",\"value\":\"1\"}],\"remove\":[\"host\",\"_raw\",\"source\",\"cribl_breaker\"]}"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.2.id", "publish_metrics"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.2.filter", "true"),
						resource.TestCheckResourceAttr("criblio_pack_pipeline.my_packpipeline", "conf.functions.2.conf", 
						  "{\"dimensions\":[\"!_*\",\"*\"],\"fields\":[{\"inFieldName\":\"_value\",\"metricType\":\"gauge\",\"outFieldExpr\":\"saas_env\"}],\"overwrite\":false,\"removeDimensions\":[\"cribl_pipe\"]}"),
					),
				},
				{
					ConfigDirectory: config.TestNameDirectory(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	})
}
