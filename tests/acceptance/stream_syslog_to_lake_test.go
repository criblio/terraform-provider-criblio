package tests

import (
	"os"
	"time"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestStreamSyslogToLake(t *testing.T) {
        if os.Getenv("DEPLOYMENT") == "onprem" {
                time.Sleep(1 * time.Second)
        }
	t.Run("plan-diff", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories:  providerFactory,
			PreventPostDestroyRefresh: true,
			Steps: []resource.TestStep{
				{
					ImportState:     true,
					ImportStateId:   "criblio_group.syslog_worker_group",
					ResourceName:    "criblio_group.syslog_worker_group",
					ConfigDirectory: config.TestNameDirectory(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("criblio_group.syslog_worker_group", "id", "syslog-workers"),
						resource.TestCheckResourceAttr("criblio_group.syslog_worker_group", "name", "syslog-workers"),
						resource.TestCheckResourceAttr("criblio_group.syslog_worker_group", "product", "stream"),
						resource.TestCheckResourceAttr("criblio_source.syslog_source", "id", "syslog-input"),
						resource.TestCheckResourceAttr("criblio_source.syslog_source", "group_id", "syslog-workers"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "id", "cribl-lake-2"),
						resource.TestCheckResourceAttr("criblio_destination.cribl_lake", "group_id", "syslog-workers"),
					),
				},
			},
		})
	})
}
