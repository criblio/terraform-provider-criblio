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
		time.Sleep(1 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						// HTTP source
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "id", "http-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.id", "http-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.type", "http"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.description", "HTTP listener for webhook events"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.host", "0.0.0.0"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.port", "10089"),
						resource.TestCheckResourceAttr("criblio_source.my_http_source", "input_http.pipeline", "default"),
						// Cloudflare HEC source
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "id", "cloudflare-hec-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.id", "cloudflare-hec-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.type", "cloudflare_hec"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.description", "Cloudflare HTTP Event Collector listener"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.host", "0.0.0.0"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.port", "10093"),
						resource.TestCheckResourceAttr("criblio_source.my_cloudflare_hec_source", "input_cloudflare_hec.hec_api", "/services/collector/event"),
						// Wiz Webhook source
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "id", "wiz-webhook-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.id", "wiz-webhook-listener"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.type", "wiz_webhook"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.description", "Wiz webhook listener"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.disabled", "false"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.host", "0.0.0.0"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.port", "10092"),
						resource.TestCheckResourceAttr("criblio_source.my_wiz_webhook_source", "input_wiz_webhook.pipeline", "default"),
					),
				},
			},
		})
	})
}
