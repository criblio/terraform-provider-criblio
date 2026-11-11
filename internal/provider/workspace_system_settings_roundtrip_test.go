package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type workspaceSettingsTestProvider struct{ systemLimitsTestProvider }

func (p *workspaceSettingsTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewWorkspaceSystemSettingsResource}
}

func (p *workspaceSettingsTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewWorkspaceSystemSettingsDataSource}
}

type workspaceSettingsBackend struct {
	mu         sync.Mutex
	stored     map[string]json.RawMessage
	patches    int
	failGET    bool
	redacted   bool
	omitSecret bool
	cloudGuard bool
	forbidden  string
	rejected   int
	t          *testing.T
}

func (b *workspaceSettingsBackend) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if req.URL.Path != "/api/v1/system/settings/conf" {
		b.t.Errorf("unexpected settings scope: %s", req.URL.Path)
		http.NotFound(w, req)
		return
	}
	switch req.Method {
	case http.MethodGet:
		if b.failGET {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	case http.MethodPatch:
		b.patches++
		var body map[string]json.RawMessage
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			b.t.Errorf("decode settings PATCH: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if b.forbidden != "" {
			b.rejected++
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": b.forbidden}); err != nil {
				b.t.Errorf("encode rejection: %v", err)
			}
			return
		}
		if b.cloudGuard {
			var saved, incoming map[string]json.RawMessage
			if err := json.Unmarshal(b.stored["api"], &saved); err != nil {
				b.t.Fatalf("decode saved API settings: %v", err)
			}
			if raw, ok := body["api"]; ok {
				if err := json.Unmarshal(raw, &incoming); err != nil {
					b.t.Fatalf("decode incoming API settings: %v", err)
				}
			}
			// Reproduce the backend bug: omitted host/port count as changes.
			if string(saved["host"]) != string(incoming["host"]) || string(saved["port"]) != string(incoming["port"]) {
				b.rejected++
				w.WriteHeader(http.StatusForbidden)
				if err := json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Updating API host/port is disabled in this deployment"}); err != nil {
					b.t.Errorf("encode rejection: %v", err)
				}
				return
			}
		}
		if raw, ok := body["upgradeSettings"]; ok {
			var upgrades map[string]json.RawMessage
			if err := json.Unmarshal(raw, &upgrades); err != nil {
				b.t.Errorf("decode upgrade settings: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if string(upgrades["disableAutomaticUpgrade"]) == "false" && len(upgrades["automaticUpgradeCheckPeriod"]) == 0 {
				http.Error(w, "automaticUpgradeCheckPeriod is required", http.StatusBadRequest)
				return
			}
		}
		for name, value := range body {
			if name == "id" || name == "managedFields" || name == "groupId" {
				b.t.Errorf("local identity or ownership leaked into PATCH: %s", name)
			}
			// Match the upstream API: replace supplied sections in full.
			b.stored[name] = value
		}
	default:
		b.t.Errorf("unexpected settings method: %s", req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	response := make(map[string]json.RawMessage, len(b.stored))
	for name, value := range b.stored {
		response[name] = value
	}
	if b.redacted {
		response["api"] = json.RawMessage(`{"ssl":{"certPath":"/cert.pem","passphrase":"********"}}`)
	}
	if b.omitSecret {
		response["api"] = json.RawMessage(`{"ssl":{"certPath":"/cert.pem"}}`)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"count": 1, "items": []any{response}}); err != nil {
		b.t.Errorf("encode settings: %v", err)
	}
}

func workspaceSettingsFactories(baseURL string) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(&workspaceSettingsTestProvider{
			systemLimitsTestProvider: systemLimitsTestProvider{baseURL: baseURL},
		}),
	}
}

func TestWorkspaceSettingsManagedSections(t *testing.T) {
	backend := &workspaceSettingsBackend{t: t, stored: map[string]json.RawMessage{
		"api":             json.RawMessage(`{"host":"0.0.0.0","port":9000}`),
		"workers":         json.RawMessage(`{"count":3,"futureWorkerField":9007199254740993}`),
		"upgradeSettings": json.RawMessage(`{"disableAutomaticUpgrade":false,"automaticUpgradeCheckPeriod":"24h","futureUpgradeField":9007199254740993,"packageUrls":[{"packageUrl":"https://example.com/package"}]}`),
	}}
	server := httptest.NewServer(backend)
	defer server.Close()
	config := func(disabled bool, extra string) string {
		return fmt.Sprintf(`resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = {
    disable_automatic_upgrade = %t
    %s
  }
}`, disabled, extra)
	}
	name := "criblio_workspace_system_settings.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: config(true, ""), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "id", "conf"),
				resource.TestCheckResourceAttr(name, "upgrade_settings.disable_automatic_upgrade", "true"),
				resource.TestCheckNoResourceAttr(name, "upgrade_settings.automatic_upgrade_check_period"),
				resource.TestCheckNoResourceAttr(name, "api.host"),
			)},
			{Config: config(true, ""), PlanOnly: true},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.stored["workers"] = json.RawMessage(`{"count":8}`)
			}, Config: config(true, ""), PlanOnly: true},
			{Config: config(false, "package_urls = []") + `
data "criblio_workspace_system_settings" "test" {}`, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "upgrade_settings.disable_automatic_upgrade", "false"),
				resource.TestCheckResourceAttr(name, "upgrade_settings.package_urls.#", "0"),
				resource.TestCheckResourceAttr("data.criblio_workspace_system_settings.test", "workers.count", "8"),
				resource.TestCheckResourceAttr("data.criblio_workspace_system_settings.test", "id", "conf"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "conf", ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"api", "workers", "upgrade_settings.automatic_upgrade_check_period", "managed_fields"}},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.stored["upgradeSettings"] = json.RawMessage(`{"disableAutomaticUpgrade":true,"automaticUpgradeCheckPeriod":"24h","futureUpgradeField":9007199254740993,"packageUrls":[]}`)
			}, Config: config(false, "package_urls = []"), PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: config(false, "package_urls = []")},
			{Config: `resource "criblio_workspace_system_settings" "test" {}`, Check: resource.TestCheckResourceAttr(name, "managed_fields.#", "0")},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.stored["upgradeSettings"] = json.RawMessage(`{"disableAutomaticUpgrade":true,"futureUpgradeField":9007199254740993}`)
			}, Config: `resource "criblio_workspace_system_settings" "test" {}`, PlanOnly: true},
			{Config: config(false, `automatic_upgrade_check_period = "24h"`)},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.stored["upgradeSettings"] = json.RawMessage(`{"futureUpgradeField":9007199254740993}`)
			}, Config: config(false, `automatic_upgrade_check_period = "24h"`)},
			{Config: config(false, `automatic_upgrade_check_period = "24h"`), PlanOnly: true},
		},
	})
	backend.mu.Lock()
	defer backend.mu.Unlock()
	var upgrades map[string]json.RawMessage
	if err := json.Unmarshal(backend.stored["upgradeSettings"], &upgrades); err != nil {
		t.Fatal(err)
	}
	if string(upgrades["futureUpgradeField"]) != "9007199254740993" || string(upgrades["disableAutomaticUpgrade"]) != "false" {
		t.Fatalf("settings were lost during reconciliation or destroy: %s", backend.stored["upgradeSettings"])
	}
	if string(backend.stored["api"]) != `{"host":"0.0.0.0","port":9000}` {
		t.Fatal("updating upgrades changed the API section")
	}
}

func TestWorkspaceSettingsEnableUpgradesRequiresPeriod(t *testing.T) {
	backend := &workspaceSettingsBackend{t: t, stored: map[string]json.RawMessage{
		"upgradeSettings": json.RawMessage(`{"disableAutomaticUpgrade":true}`),
	}}
	server := httptest.NewServer(backend)
	defer server.Close()
	config := func(disabled bool, period string) string {
		return fmt.Sprintf(`resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = {
    disable_automatic_upgrade = %t
    %s
  }
}`, disabled, period)
	}
	assertNoPatch := func() {
		backend.mu.Lock()
		defer backend.mu.Unlock()
		if backend.patches != 0 {
			t.Fatalf("invalid upgrade settings sent %d PATCHes", backend.patches)
		}
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: config(false, ""), ExpectError: regexp.MustCompile("Enabling automatic upgrades requires automatic_upgrade_check_period")},
			{PreConfig: assertNoPatch, Config: config(true, "")},
			{Config: config(false, ""), ExpectError: regexp.MustCompile("Enabling automatic upgrades requires automatic_upgrade_check_period")},
			{PreConfig: assertNoPatch, Config: config(false, `automatic_upgrade_check_period = "24h"`), Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "upgrade_settings.disable_automatic_upgrade", "false"),
				resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "upgrade_settings.automatic_upgrade_check_period", "24h"),
			)},
			{Config: config(false, `automatic_upgrade_check_period = "24h"`), PlanOnly: true},
			// An existing period can remain unmanaged when only toggling upgrades.
			{Config: config(true, "")},
			{Config: config(false, ""), Check: resource.TestCheckNoResourceAttr("criblio_workspace_system_settings.test", "upgrade_settings.automatic_upgrade_check_period")},
			{Config: config(false, ""), PlanOnly: true},
		},
	})
}

func TestWorkspaceSettingsCloudAPIGuard(t *testing.T) {
	for _, test := range []struct {
		name          string
		api           string
		configAPI     string
		forbidden     string
		expectedError string
		patches       int
	}{
		{name: "preserve unmanaged API", api: `{"host":"0.0.0.0","port":9000,"headers":{"X-Test":"keep"},"futureField":9007199254740993}`, patches: 2},
		{name: "other permissions", api: `{"host":"0.0.0.0","port":9000}`, forbidden: "Not authorized", expectedError: "Not authorized", patches: 1},
		{name: "explicit address change", api: `{"host":"0.0.0.0","port":9000}`, configAPI: `api = { port = 9001 }`, expectedError: "Updating API host/port is disabled", patches: 1},
		{name: "masked secret", api: `{"host":"0.0.0.0","port":9000,"ssl":{"passphrase":"********"}}`, expectedError: "masked SSL passphrase", patches: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			backend := &workspaceSettingsBackend{t: t, cloudGuard: true, forbidden: test.forbidden, stored: map[string]json.RawMessage{
				"api":             json.RawMessage(test.api),
				"workers":         json.RawMessage(`{"count":7}`),
				"upgradeSettings": json.RawMessage(`{"disableAutomaticUpgrade":true}`),
			}}
			server := httptest.NewServer(backend)
			defer server.Close()
			config := fmt.Sprintf(`resource "criblio_workspace_system_settings" "test" {
  %s
  upgrade_settings = {
    disable_automatic_upgrade = false
    automatic_upgrade_check_period = "24h"
  }
}`, test.configAPI)
			step := resource.TestStep{Config: config}
			steps := []resource.TestStep{step}
			if test.expectedError != "" {
				steps[0].ExpectError = regexp.MustCompile(test.expectedError)
			} else {
				steps[0].Check = resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "upgrade_settings.disable_automatic_upgrade", "false"),
					resource.TestCheckNoResourceAttr("criblio_workspace_system_settings.test", "api.host"),
					resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "managed_fields.#", "2"),
				)
				steps = append(steps, resource.TestStep{Config: config, PlanOnly: true})
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL), Steps: steps})
			backend.mu.Lock()
			defer backend.mu.Unlock()
			if backend.patches != test.patches {
				t.Errorf("PATCHes = %d, want %d", backend.patches, test.patches)
			}
			if string(backend.stored["api"]) != test.api || string(backend.stored["workers"]) != `{"count":7}` {
				t.Fatal("compatibility retry changed unmanaged settings")
			}
		})
	}
}

func TestWorkspaceSettingsReadFailurePreventsWrite(t *testing.T) {
	backend := &workspaceSettingsBackend{t: t, failGET: true}
	server := httptest.NewServer(backend)
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL),
		Steps: []resource.TestStep{{Config: `resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = { disable_automatic_upgrade = true }
}`, ExpectError: regexp.MustCompile("Unable to read existing workspace settings")}},
	})
	if backend.patches != 0 {
		t.Fatal("failed GET caused a replacement PATCH")
	}
}

func TestWorkspaceSettingsMaskedPassphrase(t *testing.T) {
	for _, mode := range []string{"masked", "omitted"} {
		t.Run(mode, func(t *testing.T) {
			backend := &workspaceSettingsBackend{t: t, redacted: mode == "masked", omitSecret: mode == "omitted", stored: map[string]json.RawMessage{
				"api": json.RawMessage(`{"ssl":{"certPath":"/cert.pem","passphrase":"old-password"}}`),
			}}
			server := httptest.NewServer(backend)
			defer server.Close()
			config := `resource "criblio_workspace_system_settings" "test" {
  api = { ssl = { cert_path = "/cert.pem", passphrase = "new-password" } }
}`
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL),
				Steps: []resource.TestStep{
					{Config: config, Check: resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "api.ssl.passphrase", "new-password")},
					{Config: config, PlanOnly: true},
				},
			})
			backend.mu.Lock()
			defer backend.mu.Unlock()
			var api struct{ SSL struct{ Passphrase string } }
			if err := json.Unmarshal(backend.stored["api"], &api); err != nil {
				t.Fatal(err)
			}
			if api.SSL.Passphrase != "new-password" || backend.patches != 1 {
				t.Fatalf("secret was not updated exactly once: %d PATCHes", backend.patches)
			}
		})
	}
}

func TestWorkspaceSettingsUnknownValuesAndIgnoredChanges(t *testing.T) {
	backend := &workspaceSettingsBackend{t: t, stored: map[string]json.RawMessage{
		"upgradeSettings": json.RawMessage(`{"disableAutomaticUpgrade":false,"automaticUpgradeCheckPeriod":"24h"}`),
	}}
	server := httptest.NewServer(backend)
	defer server.Close()
	config := `resource "terraform_data" "upgrade" { input = true }
resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = { disable_automatic_upgrade = terraform_data.upgrade.output }
}`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: workspaceSettingsFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: config},
			{Config: config, PlanOnly: true},
			{Config: `resource "criblio_workspace_system_settings" "test" {
  upgrade_settings = { disable_automatic_upgrade = false }
  lifecycle { ignore_changes = [upgrade_settings.disable_automatic_upgrade] }
}`, Check: resource.TestCheckResourceAttr("criblio_workspace_system_settings.test", "upgrade_settings.disable_automatic_upgrade", "true")},
		},
	})
}
