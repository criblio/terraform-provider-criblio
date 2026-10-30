package tests

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMonitor(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		t.Skip("Monitor acceptance tests are disabled for on-prem deployments")
	}

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
						resource.TestCheckResourceAttr("data.criblio_monitor.example_threshold", "id", "acc_test_threshold_monitor"),
						testCheckListDataSourceHasItems("data.criblio_monitors.all"),
					),
				},
				{
					Config: monitorThreshold(90),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "firing_rule.threshold.0.limit", "90"),
						resource.TestCheckResourceAttr("data.criblio_monitor.example_threshold", "firing_rule.threshold.0.limit", "90"),
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
					Config: monitorForecast(2),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "type", "forecast"),
						resource.TestCheckResourceAttr(resourceName, "forecast_config.algorithm", "linear"),
						resource.TestCheckResourceAttr(resourceName, "forecast_config.seasonality", "hourly"),
						resource.TestCheckResourceAttr(resourceName, "forecast_config.deviations", "2"),
						resource.TestCheckNoResourceAttr(resourceName, "anomaly_config"),
					),
				},
				{
					Config: monitorForecast(3),
					Check:  resource.TestCheckResourceAttr(resourceName, "forecast_config.deviations", "3"),
				},
				{
					Config:   monitorForecast(3),
					PlanOnly: true,
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateId:     "acc_test_forecast_monitor",
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("outlier-dispatch", func(t *testing.T) {
		t.Skip("Outlier monitor acceptance test is disabled: the backend rejects outlier as an unsupported detection type")

		resourceName := "criblio_monitor.example_outlier"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					Config: monitorOutlier(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "type", "outlier"),
						resource.TestCheckResourceAttr(resourceName, "outlier_config.algorithm", "mad"),
						resource.TestCheckNoResourceAttr(resourceName, "anomaly_config"),
						resource.TestCheckNoResourceAttr(resourceName, "forecast_config"),
					),
				},
				{
					Config:   monitorOutlier(),
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
  managed_by = "terraform"
  silence    = []

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
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
    label = "A"
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

data "criblio_monitor" "example_threshold" {
  id = criblio_monitor.example_threshold.id
}

data "criblio_monitors" "all" {
  depends_on = [criblio_monitor.example_threshold]
}
`
}

func monitorForecast(deviations int) string {
	return `resource "criblio_monitor" "example_forecast" {
  id         = "acc_test_forecast_monitor"
  name       = "Acceptance Test - Forecast Monitor"
  enabled    = true
  type       = "forecast"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
  }

  forecast_config = {
    algorithm   = "linear"
    window      = "1h"
    deviations  = ` + strconv.Itoa(deviations) + `
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
    label = "A"
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

func monitorOutlier() string {
	return `resource "criblio_monitor" "example_outlier" {
  id         = "acc_test_outlier_monitor"
  name       = "Acceptance Test - Outlier Monitor"
  enabled    = true
  type       = "outlier"
  dataset_id = "metrics"
  managed_by = "terraform"
  silence    = []

  monitor_policy_notification_config = {
    type    = "policy"
    enabled = false
    config  = []
  }

  outlier_config = {
    algorithm = "mad"
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
    label = "A"
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
