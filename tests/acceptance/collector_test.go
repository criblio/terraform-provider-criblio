package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCollector(t *testing.T) {
	if os.Getenv("DEPLOYMENT") == "onprem" {
		time.Sleep(1 * time.Second)
	}

	t.Run("plan-diff", func(t *testing.T) {
		suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
		splunkID := "splunk-demo-collector-" + suffix
		restID := "rest-api-demo-collector-" + suffix
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: providerFactory,
			Steps: []resource.TestStep{
				{
					Config: collectorConfig(t, suffix),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector", "id", splunkID),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector",
							"input_collector_splunk.environment", "demo"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector",
							"input_collector_splunk.collector.type", "splunk"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector",
							"input_collector_splunk.collector.conf.authentication", "token"),
						resource.TestCheckResourceAttr("criblio_collector.splunk_access_log_collector",
							"input_collector_splunk.collector.conf.disable_time_filter", "false"),

						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "id", restID),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "group_id", "default"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "input_collector_rest.environment", "demo"),
						resource.TestCheckResourceAttr("criblio_collector.rest_api_collector", "input_collector_rest.collector.type", "rest"),

						/*
							resource.TestCheckResourceAttr("criblio_collector.script_collector", "group_id", "default"),
							resource.TestCheckResourceAttr("criblio_collector.script_collector",
								"input_collector_script.collector.type", "script"),
							resource.TestCheckResourceAttr("criblio_collector.script_collector",
								"input_collector_script.collector.conf.shell", "/bin/bash"),
							resource.TestCheckResourceAttr("criblio_collector.script_collector",
								"input_collector_script.collector.conf.discover_script", "echo 1"),
						*/
					),
				},
			},
		})
	})
}

func collectorConfig(t *testing.T, suffix string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("testdata", "TestCollector", "plan-diff", "main.tf"))
	if err != nil {
		t.Fatalf("read collector test config: %v", err)
	}
	config := string(content)
	config = strings.ReplaceAll(config, "splunk-demo-collector", "splunk-demo-collector-"+suffix)
	config = strings.ReplaceAll(config, "rest-api-demo-collector", "rest-api-demo-collector-"+suffix)
	return config
}
