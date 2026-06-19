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
