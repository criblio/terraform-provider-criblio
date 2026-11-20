package export

import (
	"encoding/json"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
	"github.com/stretchr/testify/require"
)

// A Guard Rule that belongs to no ruleset: the API omits `rulesets` entirely.
const guardRuleWireWithoutRulesets = `{
	"id": "employee_badge_id",
	"groupId": "default",
	"description": "Employee badge identifiers.",
	"regex": "/\\bEMP-\\d{6}\\b/"
}`

// The same rule as the API reports it when the field is present but empty.
const guardRuleWireEmptyRulesets = `{
	"id": "employee_badge_id",
	"groupId": "default",
	"regex": "/\\bEMP-\\d{6}\\b/",
	"rulesets": []
}`

// `rulesets` is a required attribute, so generated HCL must always carry it.
// A full export previously emitted a rule with no `rulesets` at all, and
// `terraform plan` failed with `The argument "rulesets" is required, but no
// definition was found`, because the generator drops null and empty-list
// attributes.
func TestGuardRuleExportAlwaysEmitsRulesets(t *testing.T) {
	for name, wire := range map[string]string{
		"absent": guardRuleWireWithoutRulesets,
		"empty":  guardRuleWireEmptyRulesets,
	} {
		t.Run(name, func(t *testing.T) {
			var model provider.GuardRuleModel
			require.NoError(t, json.Unmarshal([]byte(wire), &model))

			attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_guard_rule", registry.Entry{}))
			require.NoError(t, err)
			ensureGuardRuleRulesets(attrs)

			require.Equal(t, hcl.KindList, attrs["rulesets"].Kind)
			require.Empty(t, attrs["rulesets"].List)

			// Rendering is where the attribute was lost, so assert on the HCL.
			out, err := hcl.ResourceBlockBytes("criblio_guard_rule", "guard_rule_default_employee_badge_id", attrs, hcl.DefaultResourceBlockOptions())
			require.NoError(t, err)
			require.Regexp(t, `rulesets\s*= \[\]`, string(out))
		})
	}
}

// A rule that does belong to rulesets keeps them verbatim.
func TestGuardRuleExportKeepsRulesets(t *testing.T) {
	var model provider.GuardRuleModel
	require.NoError(t, json.Unmarshal([]byte(`{
		"id": "credit_card",
		"groupId": "default",
		"regex": "/\\b(?:\\d[ -]*?){13,16}\\b/",
		"rulesets": ["pci-compliance"]
	}`), &model))

	attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_guard_rule", registry.Entry{}))
	require.NoError(t, err)
	ensureGuardRuleRulesets(attrs)

	require.Len(t, attrs["rulesets"].List, 1)
	require.Equal(t, "pci-compliance", attrs["rulesets"].List[0].String)
}
