package tests

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMonitor(t *testing.T) {
	t.Run("plan-diff", func(t *testing.T) {
		resourceName := "criblio_monitor.example_threshold"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: monitorThreshold(80),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "id", "acc_test_threshold_monitor"),
						resource.TestCheckResourceAttr(resourceName, "type", "threshold"),
						resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
						resource.TestCheckResourceAttr(resourceName, "managed_by", "terraform"),
						resource.TestCheckResourceAttr(resourceName, "firing_rule.threshold.0.limit", "80"),
						resource.TestCheckResourceAttr(resourceName, "priority_scalar_value.value", "P3"),
					),
				},
				{
					Config: monitorThreshold(90),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "firing_rule.threshold.0.limit", "90"),
					),
				},
				{
					Config:   monitorThreshold(90),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     "acc_test_threshold_monitor",
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("detection-config-dispatch", func(t *testing.T) {
		resourceName := "criblio_monitor.example_forecast"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: monitorForecast(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "type", "forecast"),
						resource.TestCheckResourceAttr(resourceName, "forecast_config.algorithm", "linear"),
						resource.TestCheckResourceAttr(resourceName, "forecast_config.seasonality", "hourly"),
						resource.TestCheckNoResourceAttr(resourceName, "anomaly_config"),
					),
				},
				{
					Config:   monitorForecast(),
					PlanOnly: true,
				},
			},
		})
	})
}

func monitorThreshold(limit int) string {
	return `resource "criblio_monitor" "example_threshold" {
  id         = "acc_test_threshold_monitor"
  name       = "Acceptance Test - Threshold Monitor"
  enabled    = true
  type       = "threshold"
  dataset_id = "metrics"

  priority_scalar_value = { value = "P3" }
  team_scalar_value     = { value = "" }

  query = {
    A = {
      mode       = "promql"
      dataset_id = "metrics"
      promql     = "cpu_usage_percent"
    }
  }

  expr = []

  firing_condition = {
    fire_delay  = 300
    clear_delay = 60
  }

  firing_rule = {
    label = ""
    threshold = [{
      severity        = "warning"
      limit           = ` + strconv.Itoa(limit) + `
      operator        = "gt"
      included_tags   = []
      excluded_tags   = []
      times_triggered = 1
    }]
  }

  aetos_metadata = {}
}
`
}

func monitorForecast() string {
	return `resource "criblio_monitor" "example_forecast" {
  id         = "acc_test_forecast_monitor"
  name       = "Acceptance Test - Forecast Monitor"
  enabled    = true
  type       = "forecast"
  dataset_id = "metrics"

  forecast_config = {
    algorithm   = "linear"
    window      = "1h"
    deviations  = 2
    seasonality = "hourly"
  }

  priority_scalar_value = { value = "P3" }
  team_scalar_value     = { value = "" }

  query = {
    A = {
      mode       = "promql"
      dataset_id = "metrics"
      promql     = "cpu_usage_percent"
    }
  }

  expr = []

  firing_condition = {
    fire_delay  = 300
    clear_delay = 60
  }

  firing_rule = {
    label = ""
    threshold = [{
      severity        = "warning"
      limit           = 90
      operator        = "gt"
      included_tags   = []
      excluded_tags   = []
      times_triggered = 1
    }]
  }

  aetos_metadata = {}
}
`
}
