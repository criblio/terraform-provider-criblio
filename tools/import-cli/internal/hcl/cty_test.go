package hcl

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// Regression: search dashboard elements are a list of one-of maps (e.g. dashboard_element_input
// vs dashboard_element_visualization). normalizeListOfMaps adds null for absent branches;
// an extra PruneNulls on each element before ValueToCty removed those nulls and caused
// cty.ListVal to panic on inconsistent element types.
func TestValueToCty_listHeterogeneousUnionMapBranches(t *testing.T) {
	v := Value{
		Kind: KindList,
		List: []Value{
			{
				Kind: KindMap,
				Map: map[string]Value{
					"branch_a": {Kind: KindString, String: "only-a"},
				},
			},
			{
				Kind: KindMap,
				Map: map[string]Value{
					"branch_b": {Kind: KindString, String: "only-b"},
				},
			},
		},
	}
	_, err := ValueToCty(v)
	require.NoError(t, err)
}

// Nested one-of objects (e.g. dashboard_element_input vs dashboard_element_visualization) with
// different inner config shapes must not panic: alternating branches yield incompatible cty object
// types until column-wise homogenization.
func TestValueToCty_listNestedUnionObjectBranches(t *testing.T) {
	innerA := Value{Kind: KindMap, Map: map[string]Value{
		"config": {Kind: KindMap, Map: map[string]Value{
			"input_shape": {Kind: KindString, String: "a"},
		}},
	}}
	innerB := Value{Kind: KindMap, Map: map[string]Value{
		"config": {Kind: KindMap, Map: map[string]Value{
			"viz_shape": {Kind: KindNumber, Number: 1},
		}},
	}}
	v := Value{
		Kind: KindList,
		List: []Value{
			{Kind: KindMap, Map: map[string]Value{"branch_a": innerA}},
			{Kind: KindMap, Map: map[string]Value{"branch_b": innerB}},
		},
	}
	_, err := ValueToCty(v)
	require.NoError(t, err)
}

func TestUnifySiblingValues_dynamicNullPlaceholderDoesNotWipeObjects(t *testing.T) {
	obj := cty.ObjectVal(map[string]cty.Value{"id": cty.StringVal("tile-1")})
	col := []cty.Value{obj, cty.NullVal(cty.DynamicPseudoType)}
	out, err := unifySiblingValues(col)
	require.NoError(t, err)
	require.True(t, out[0].IsKnown())
	require.True(t, out[0].GetAttr("id").IsKnown())
	require.Equal(t, "tile-1", out[0].GetAttr("id").AsString())
	require.True(t, out[1].IsNull())
	require.True(t, out[1].Type().IsObjectType(), "placeholder should become typed null matching column")
}

func TestPruneSearchDashboardElementsCty_omitsExclusiveNullBranches(t *testing.T) {
	viz := cty.ObjectVal(map[string]cty.Value{"id": cty.StringVal("tile-1")})
	inpNull := cty.NullVal(cty.Object(map[string]cty.Type{
		"id": cty.String,
	}))
	row := cty.ObjectVal(map[string]cty.Value{
		"dashboard_element_visualization": viz,
		"dashboard_element_input":         inpNull,
	})
	out, err := PruneSearchDashboardElementsCty(cty.TupleVal([]cty.Value{row}))
	require.NoError(t, err)
	r0 := out.Index(cty.NumberIntVal(0))
	require.True(t, r0.Type().HasAttribute("dashboard_element_visualization"))
	require.False(t, r0.Type().HasAttribute("dashboard_element_input"))
}

// Regression: list homogenization can yield a non-null but empty dashboard_element_input object
// alongside a real visualization; fixed priority wrongly preferred input over visualization.
func TestPruneSearchDashboardElementsCty_prefersVisualizationOverEmptyInputShell(t *testing.T) {
	viz := cty.ObjectVal(map[string]cty.Value{
		"id":   cty.StringVal("chart-1"),
		"type": cty.StringVal("chart.bar"),
	})
	inpShell := cty.ObjectVal(map[string]cty.Value{
		"id":   cty.NullVal(cty.String),
		"type": cty.NullVal(cty.String),
	})
	row := cty.ObjectVal(map[string]cty.Value{
		"dashboard_element_visualization": viz,
		"dashboard_element_input":         inpShell,
	})
	out, err := PruneSearchDashboardElementsCty(cty.TupleVal([]cty.Value{row}))
	require.NoError(t, err)
	r0 := out.Index(cty.NumberIntVal(0))
	require.True(t, r0.Type().HasAttribute("dashboard_element_visualization"))
	require.False(t, r0.Type().HasAttribute("dashboard_element_input"))
}

func TestCtyReplaceUnknownWithNullForHCL_hclwriteSafe(t *testing.T) {
	unknownObj := cty.UnknownVal(cty.Object(map[string]cty.Type{
		"x": cty.String,
	}))
	v := cty.ObjectVal(map[string]cty.Value{
		"elements": cty.ListVal([]cty.Value{unknownObj}),
	})
	out, err := ctyReplaceUnknownWithNullForHCL(v)
	require.NoError(t, err)
	require.True(t, out.IsWhollyKnown())
	require.NotPanics(t, func() { hclwrite.TokensForValue(out) })
}
