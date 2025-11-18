package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSource(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(2 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "id", "http-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.id", "http-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.type", "http"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.description", "HTTP listener for webhook events"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.host", "0.0.0.0"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.port", "10089"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.environment", "main"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.pipeline", "default"),
					),
				},
			},
		})
	})
}
