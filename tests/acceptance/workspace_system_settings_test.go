package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkspaceSystemSettings(t *testing.T) {
	workspace := os.Getenv("CRIBL_TEST_WORKSPACE_SETTINGS_WORKSPACE")
	if workspace == "" {
		t.Skip("Set CRIBL_TEST_WORKSPACE_SETTINGS_WORKSPACE to a disposable workspace to test workspace settings")
	}
	if workspace != os.Getenv("CRIBL_WORKSPACE_ID") {
		t.Fatal("CRIBL_WORKSPACE_ID must match CRIBL_TEST_WORKSPACE_SETTINGS_WORKSPACE")
	}
	name := "criblio_workspace_system_settings.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{Config: workspaceSettingsAcceptanceConfig(true), Check: resource.TestCheckResourceAttr(name, "upgrade_settings.disable_automatic_upgrade", "true")},
			{Config: workspaceSettingsAcceptanceConfig(false), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "upgrade_settings.disable_automatic_upgrade", "false"),
				resource.TestCheckResourceAttr(name, "upgrade_settings.automatic_upgrade_check_period", "24h"),
				resource.TestCheckResourceAttr("data.criblio_workspace_system_settings.test", "upgrade_settings.disable_automatic_upgrade", "false"))},
			{Config: workspaceSettingsAcceptanceConfig(false), PlanOnly: true},
			{ResourceName: name, ImportState: true, ImportStateId: "conf", ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"api", "apps", "backups", "custom_logo", "pii", "proxy", "rollback", "shutdown", "sni", "sockets", "support", "system", "tls", "upgrade_group_settings", "workers", "managed_fields", "upgrade_settings.enable_legacy_edge_upgrade", "upgrade_settings.package_urls", "upgrade_settings.reject_unauthorized", "upgrade_settings.upgrade_source"}},
		},
	})
}

func workspaceSettingsAcceptanceConfig(disabled bool) string {
	return fmt.Sprintf(`resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = {
    disable_automatic_upgrade      = %t
    automatic_upgrade_check_period = "24h"
  }
}
data "criblio_workspace_system_settings" "test" {
  depends_on = [criblio_workspace_system_settings.test]
}`, disabled)
}
