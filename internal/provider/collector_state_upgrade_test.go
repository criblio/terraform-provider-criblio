package provider

import (
	"encoding/json"
	"testing"
)

func TestUpgradeCollectorStateJSON(t *testing.T) {
	upgraded, err := upgradeCollectorStateJSON([]byte(`{
		"id":"collector",
		"input_collector_rest":{"schedule":{"run":{"earliest":1700000000,"latest":-10.5}}},
		"input_collector_s3":null
	}`))
	if err != nil {
		t.Fatalf("upgradeCollectorStateJSON returned error: %v", err)
	}
	var state map[string]any
	if err := json.Unmarshal(upgraded, &state); err != nil {
		t.Fatalf("decode upgraded state: %v", err)
	}
	rest := state["input_collector_rest"].(map[string]any)
	run := rest["schedule"].(map[string]any)["run"].(map[string]any)
	if got := run["earliest"]; got != "1700000000" {
		t.Fatalf("earliest = %#v, want 1700000000", got)
	}
	if got := run["latest"]; got != "-10.5" {
		t.Fatalf("latest = %#v, want -10.5", got)
	}
}
