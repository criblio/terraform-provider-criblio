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

const guardRuleWireWithRulesets = `{
	"id": "credit_card",
	"groupId": "default",
	"regex": "/\\b(?:\\d[ -]*?){13,16}\\b/",
	"rulesets": ["pci-compliance"]
}`

func guardRuleAttrs(t *testing.T, wire string) map[string]hcl.Value {
	t.Helper()
	var model provider.GuardRuleModel
	require.NoError(t, json.Unmarshal([]byte(wire), &model))
	attrs, err := hcl.ModelToValue(&model, hclOptionsForType("criblio_guard_rule", registry.Entry{}))
	require.NoError(t, err)
	return attrs
}

// A Guard Rule in no ruleset cannot round-trip and must not be exported.
// `rulesets` is required by the provider schema and the API enforces at least
// one entry, so such a rule has no valid representation: omitting the attribute
// failed `terraform plan` with `The argument "rulesets" is required`, and
// emitting `rulesets = []` failed apply with `HTTP 400: At least one ruleset is
// required`. Both broke TestIntegration_FullExport_Cloud.
func TestGuardRuleExportSkipsRuleWithoutRulesets(t *testing.T) {
	for name, wire := range map[string]string{
		"absent": guardRuleWireWithoutRulesets,
		"empty":  guardRuleWireEmptyRulesets,
	} {
		t.Run(name, func(t *testing.T) {
			require.True(t, skipGuardRuleWithoutRulesets("criblio_guard_rule", guardRuleAttrs(t, wire)))
		})
	}
}

// A rule that does belong to rulesets exports normally, keeping them verbatim.
func TestGuardRuleExportKeepsRulesets(t *testing.T) {
	attrs := guardRuleAttrs(t, guardRuleWireWithRulesets)

	require.False(t, skipGuardRuleWithoutRulesets("criblio_guard_rule", attrs))
	require.Len(t, attrs["rulesets"].List, 1)
	require.Equal(t, "pci-compliance", attrs["rulesets"].List[0].String)

	out, err := hcl.ResourceBlockBytes("criblio_guard_rule", "guard_rule_default_credit_card", attrs, hcl.DefaultResourceBlockOptions())
	require.NoError(t, err)
	require.Regexp(t, `rulesets\s*= \["pci-compliance"\]`, string(out))
}

// The skip is specific to Guard Rules; other types with an empty list attribute
// are untouched.
func TestGuardRuleSkipDoesNotApplyToOtherTypes(t *testing.T) {
	attrs := map[string]hcl.Value{"rulesets": {Kind: hcl.KindList, List: []hcl.Value{}}}
	require.False(t, skipGuardRuleWithoutRulesets("criblio_event_breaker_ruleset", attrs))
}
