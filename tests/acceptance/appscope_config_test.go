package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppscopeConfig(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config:             appConfig,
					ExpectNonEmptyPlan: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "description", "A sample AppScope configuration"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_appscope_config.my_appscopeconfig", "lib", "cribl"),
					),
				},
			},
		})
	})
}

const appConfig = `resource "criblio_appscope_config" "my_appscopeconfig" {
  config = {
    metric = {
      enable = true
      format = {
        type      = "ndjson"
        verbosity = 4
      }
      watch = [
        {
          type = "statsd"
        },
        {
          type = "fs"
        },
        {
          type = "net"
        },
        {
          type = "http"
        },
        {
          type = "dns"
        },
        {
          type = "process"
        }
      ]
      transport = {
        type = "udp"
        host = "127.0.0.1"
        port = 8125
      }
    }
    event = {
      enable = true
      type   = "ndjson"
      format = {
        maxeventpersec = 10000
        enhancefs      = true
      }
      watch = [
        {
          type       = "console"
          enabled    = true
          name       = "(stdout)|(stderr)"
          value      = ".*"
          allowbinary = true
        },
        {
          type    = "file"
          enabled = true
          name    = "(\\/logs?\\/)|(\\.log$)|(\\.log[.\\d])"
          value   = ".*"
        },
        {
          type    = "fs"
          enabled = true
          name    = ".*"
          field   = ".*"
          value   = ".*"
        },
        {
          type    = "net"
          enabled = true
          name    = ".*"
          field   = ".*"
          value   = ".*"
        },
        {
          type    = "dns"
          enabled = true
          name    = ".*"
          field   = ".*"
          value   = ".*"
        },
        {
          type    = "http"
          enabled = true
          name    = ".*"
          field   = ".*"
          value   = ".*"
          headers = ["host","user-agent","x-forwarded-for","x-appscope"]
        },
        {
          type    = "metric"
          enabled = false
          name    = ".*"
          field   = ".*"
          value   = ".*"
        }
      ]
      transport = {
        type = "tcp"
        host = "127.0.0.1"
        port = 9109
        tls = {
          enable = false
        }
      }
    }
    libscope = {
      configevent    = true
      summaryperiod  = 10
      commanddir     = "/tmp"
      log = {
        level = "warning"
        transport = {
          type   = "file"
          buffer = "line"
        }
      }
    }
    cribl = {
      use_scope_source_transport = true
    }
    payload = {
      enable = false
      dir    = "/tmp"
    }
  }
  description = "A sample AppScope configuration"
  group_id    = "default"
  id          = "new_appscope_config"
  lib         = "cribl"
  tags        = "cribl, test"
}

output "appscope_config" {
  value = criblio_appscope_config.my_appscopeconfig
}

data "criblio_appscope_config" "my_appscopeconfig" {
  group_id = "default"
}

provider "criblio" {
  server_url = "https://app.cribl-playground.cloud"
  organization_id = "beautiful-nguyen-y8y4azd"
  workspace_id = "tfprovider"
  version = "999.99.9"
}
`
