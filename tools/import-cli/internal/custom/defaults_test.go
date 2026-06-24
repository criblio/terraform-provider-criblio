package custom

import (
	"testing"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
)

func TestApplySubscriptionDefaultsRemovesEmptyConsumer(t *testing.T) {
	attrs := map[string]hcl.Value{
		"consumer": {
			Kind: hcl.KindMap,
			Map: map[string]hcl.Value{
				"connections": {Kind: hcl.KindList, List: []hcl.Value{}},
				"disabled":    {Kind: hcl.KindNull},
				"type":        {Kind: hcl.KindNull},
			},
		},
	}

	ApplySubscriptionDefaults(attrs)

	if _, ok := attrs["consumer"]; ok {
		t.Fatalf("empty consumer was not removed: %#v", attrs["consumer"])
	}
}

func TestApplyProjectDefaultsRemovesEmptyConsumers(t *testing.T) {
	attrs := map[string]hcl.Value{
		"consumers": {Kind: hcl.KindMap, Map: map[string]hcl.Value{}},
	}

	ApplyProjectDefaults(attrs)

	if _, ok := attrs["consumers"]; ok {
		t.Fatalf("empty consumers was not removed: %#v", attrs["consumers"])
	}
	if subscriptions, ok := attrs["subscriptions"]; !ok || subscriptions.Kind != hcl.KindList {
		t.Fatalf("subscriptions default was not added: %#v", attrs["subscriptions"])
	}
	if destinations, ok := attrs["destinations"]; !ok || destinations.Kind != hcl.KindList {
		t.Fatalf("destinations default was not added: %#v", attrs["destinations"])
	}
}

func TestApplyProjectDefaultsKeepsConfiguredConsumers(t *testing.T) {
	attrs := map[string]hcl.Value{
		"consumers": {
			Kind: hcl.KindMap,
			Map: map[string]hcl.Value{
				"subscription_a": {Kind: hcl.KindString, String: "destination_a"},
			},
		},
	}

	ApplyProjectDefaults(attrs)

	if _, ok := attrs["consumers"]; !ok {
		t.Fatalf("configured consumers was removed")
	}
}

func TestApplySubscriptionDefaultsKeepsConfiguredConsumer(t *testing.T) {
	attrs := map[string]hcl.Value{
		"consumer": {
			Kind: hcl.KindMap,
			Map: map[string]hcl.Value{
				"type": {Kind: hcl.KindString, String: "subscription"},
			},
		},
	}

	ApplySubscriptionDefaults(attrs)

	if _, ok := attrs["consumer"]; !ok {
		t.Fatalf("configured consumer was removed")
	}
}
