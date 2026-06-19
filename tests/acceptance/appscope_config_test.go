package tests

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppscopeConfig(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}
	t.Run("plan-diff", func(t *testing.T) {
		resourceName := "criblio_appscope_config.my_appscopeconfig"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: appscopeConfig("A sample AppScope configuration", "cribl, test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "A sample AppScope configuration"),
						resource.TestCheckResourceAttr(resourceName, "group_id", "default"),
						resource.TestCheckResourceAttr(resourceName, "lib", "cribl"),
						resource.TestCheckResourceAttr(resourceName, "id", "sample_appscope_config"),
						resource.TestCheckResourceAttr(resourceName, "tags", "cribl, test"),
						resource.TestCheckResourceAttr(resourceName, "config.metric.watch.0.type", "statsd"),
						resource.TestCheckResourceAttr(resourceName, "config.event.transport.tls.enable", "false"),
					),
				},
				{
					Config: appscopeConfig("Updated AppScope configuration", "cribl, updated"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "description", "Updated AppScope configuration"),
						resource.TestCheckResourceAttr(resourceName, "tags", "cribl, updated"),
					),
				},
				{
					Config:   appscopeConfig("Updated AppScope configuration", "cribl, updated"),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     `{"group_id":"default","id":"sample_appscope_config"}`,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func appscopeConfig(description string, tags string) string {
	return `resource "criblio_appscope_config" "my_appscopeconfig" {
  config = {
    cribl = {
      use_scope_source_transport = true
    }
    event = {
      enable = true
      type   = "ndjson"
      format = {
        enhancefs      = true
        maxeventpersec = 10000
      }
      watch = [
        {
          allowbinary = true
          enabled     = true
          name        = "(stdout)|(stderr)"
          type        = "console"
          value       = ".*"
        }
      ]
      transport = {
        host = "127.0.0.1"
        port = 9109
        tls = {
          enable = false
        }
        type = "tcp"
      }
    }
    libscope = {
      commanddir    = "/tmp"
      configevent   = true
      summaryperiod = 10
      log = {
        level = "warning"
        transport = {
          buffer = "line"
          type   = "file"
        }
      }
    }
    metric = {
      enable = true
      format = {
        type      = "ndjson"
        verbosity = 4
      }
      transport = {
        host = "127.0.0.1"
        port = 8125
        type = "udp"
      }
      watch = [
        {
          type = "statsd"
        }
      ]
    }
    payload = {
      dir    = "/tmp"
      enable = false
    }
  }
  description = "` + description + `"
  group_id    = "default"
  id          = "sample_appscope_config"
  lib         = "cribl"
  tags        = "` + tags + `"
}
`
}
