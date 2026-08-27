package hcl

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestPruneNullCollectionPlaceholdersPreservesEmptyCollectionsAndTypes(t *testing.T) {
	input := cty.ObjectVal(map[string]cty.Value{
		"empty_list": cty.ListValEmpty(cty.String),
		"empty_map":  cty.MapValEmpty(cty.String),
		"null_value": cty.NullVal(cty.String),
		"set": cty.SetVal([]cty.Value{
			cty.StringVal("kept"),
			cty.NullVal(cty.String),
		}),
	})

	got := pruneNullCollectionPlaceholders(input, false)

	require.True(t, got.Type().HasAttribute("empty_list"))
	require.True(t, got.GetAttr("empty_list").Type().IsListType())
	require.True(t, got.Type().HasAttribute("empty_map"))
	require.True(t, got.GetAttr("empty_map").Type().IsMapType())
	require.False(t, got.Type().HasAttribute("null_value"))
	require.True(t, got.GetAttr("set").Type().IsSetType())
	require.Equal(t, 1, got.GetAttr("set").LengthInt())
}

func TestPruneNullCollectionPlaceholdersHandlesHeterogeneousMapAndSet(t *testing.T) {
	left := cty.ObjectVal(map[string]cty.Value{
		"left":  cty.StringVal("value"),
		"right": cty.NullVal(cty.String),
	})
	right := cty.ObjectVal(map[string]cty.Value{
		"left":  cty.NullVal(cty.String),
		"right": cty.StringVal("value"),
	})
	input := cty.ObjectVal(map[string]cty.Value{
		"map": cty.MapVal(map[string]cty.Value{"a": left, "b": right}),
		"set": cty.SetVal([]cty.Value{left, right}),
	})

	var got cty.Value
	require.NotPanics(t, func() {
		got = pruneNullCollectionPlaceholders(input, false)
	})
	require.True(t, got.GetAttr("map").Type().IsObjectType())
	require.True(t, got.GetAttr("set").Type().IsTupleType())
}

func TestValueToCtyPreservesNullListElements(t *testing.T) {
	got, err := ValueToCty(Value{Kind: KindList, List: []Value{
		{Kind: KindString, String: "before"},
		{Kind: KindNull},
		{Kind: KindString, String: "after"},
	}})
	require.NoError(t, err)
	require.Equal(t, 3, got.LengthInt())
	require.True(t, got.Index(cty.NumberIntVal(1)).IsNull())

	objects, err := ValueToCty(Value{Kind: KindList, List: []Value{{
		Kind: KindMap,
		Map:  map[string]Value{"optional": {Kind: KindNull}},
	}}})
	require.NoError(t, err)
	require.True(t, objects.Index(cty.NumberIntVal(0)).Type().HasAttribute("optional"))
	require.True(t, objects.Index(cty.NumberIntVal(0)).GetAttr("optional").IsNull())
}

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
func TestStripSearchDashboardElementsConfigNullKeysCty_dropsHomogenizedNullConfigKeys(t *testing.T) {
	// Simulates list homogenization adding "color" = null on one element while another has a real value.
	viz := cty.ObjectVal(map[string]cty.Value{
		"id": cty.StringVal("chart-1"),
		"config": cty.MapVal(map[string]cty.Value{
			"color": cty.NullVal(cty.String),
		}),
	})
	row := cty.ObjectVal(map[string]cty.Value{
		"dashboard_element_visualization": viz,
	})
	tup := cty.TupleVal([]cty.Value{row})
	out, err := StripSearchDashboardElementsConfigNullKeysCty(tup)
	require.NoError(t, err)
	cfg := out.Index(cty.NumberIntVal(0)).GetAttr("dashboard_element_visualization").GetAttr("config")
	require.True(t, cfg.IsKnown())
	require.Equal(t, 0, cfg.LengthInt(), "null-only config keys should be omitted")
}

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

// Regression: homogenization can set both search_query_inline (sparse shell from another row)
// and search_query_metric; schema forbids both — Terraform reports ConflictsWith / missing inline fields.
func TestPruneSearchQueryObjectCty_keepsExactlyOneIncludingMetric(t *testing.T) {
	sparseInline := cty.ObjectVal(map[string]cty.Value{
		"earliest": cty.ObjectVal(map[string]cty.Value{
			"number": cty.NumberFloatVal(-300),
		}),
	})
	metric := cty.ObjectVal(map[string]cty.Value{
		"type": cty.StringVal("metric"),
		"queries": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"local_id": cty.StringVal("q1"),
				"query":    cty.StringVal(`dataset="logs"`),
			}),
		}),
	})
	s := cty.ObjectVal(map[string]cty.Value{
		"search_query_inline": sparseInline,
		"search_query_metric": metric,
	})
	out, err := pruneSearchQueryObjectCty(s)
	require.NoError(t, err)
	require.True(t, out.Type().HasAttribute("search_query_metric"))
	require.False(t, out.Type().HasAttribute("search_query_inline"))
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
