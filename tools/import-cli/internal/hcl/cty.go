// Conversion from our Value type to cty.Value for hclwrite SetAttributeValue.
package hcl

import (
	"fmt"
	"sort"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// VarRefPlaceholderPrefix is the prefix for a placeholder string that will be replaced
// with jsonencode(var.<name>) in the generated HCL. Used for nested secret refs.
const VarRefPlaceholderPrefix = "__VAR_REF__"
const VarRefPlaceholderSuffix = "__"

// PlainVarRefPlaceholderPrefix/Suffix are replaced with var.<name> (no jsonencode) for plain string secrets (e.g. criblio_secret.value).
const PlainVarRefPlaceholderPrefix = "__VAR_REF_PLAIN__"
const PlainVarRefPlaceholderSuffix = "__"

// ValueToCty converts a Value to cty.Value so it can be passed to hclwrite.Body.SetAttributeValue.
// Nested maps and lists are converted recursively; null and structure are preserved.
// KindSensitive with a variable name is emitted as a placeholder string; FileWithResources
// replaces it with jsonencode(var.xxx) so secrets are always JSON-encoded.
func ValueToCty(v Value) (cty.Value, error) {
	out, err := valueToCtyInner(v)
	if err != nil {
		return cty.NilVal, err
	}
	return ctyReplaceUnknownWithNullForHCL(out)
}

func valueToCtyInner(v Value) (cty.Value, error) {
	switch v.Kind {
	case KindNull:
		return cty.NullVal(cty.DynamicPseudoType), nil
	case KindSensitive:
		if isVariableName(v.Sensitive) {
			return cty.StringVal(VarRefPlaceholderPrefix + v.Sensitive + VarRefPlaceholderSuffix), nil
		}
		return cty.StringVal(v.Sensitive), nil
	case KindVariableRef:
		if isVariableName(v.VarName) {
			return cty.StringVal(PlainVarRefPlaceholderPrefix + v.VarName + PlainVarRefPlaceholderSuffix), nil
		}
		return cty.NullVal(cty.DynamicPseudoType), nil
	case KindString:
		return cty.StringVal(v.String), nil
	case KindBool:
		return cty.BoolVal(v.Bool), nil
	case KindNumber:
		return cty.NumberFloatVal(v.Number), nil
	case KindList:
		v = normalizeAndPruneList(v)
		elems, err := listToCtyElems(v.List)
		if err != nil {
			return cty.NilVal, err
		}
		if len(elems) == 0 {
			return cty.ListValEmpty(cty.DynamicPseudoType), nil
		}
		elems, err = homogenizeListElementObjectsUntilStable(elems)
		if err != nil {
			return cty.NilVal, err
		}
		elems, err = unifyListElementTypes(elems)
		if err != nil {
			return cty.NilVal, err
		}
		if cty.CanListVal(elems) {
			return cty.ListVal(elems), nil
		}
		// Provider schema expects a list; hclwrite accepts a tuple literal for mixed
		// element shapes when structural homogenization is not enough.
		return cty.TupleVal(elems), nil
	case KindMap:
		m := make(map[string]cty.Value, len(v.Map))
		keys := make([]string, 0, len(v.Map))
		for k := range v.Map {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			c, err := ValueToCty(v.Map[k])
			if err != nil {
				return cty.NilVal, fmt.Errorf("map.%s: %w", k, err)
			}
			m[k] = c
		}
		if len(m) == 0 {
			return cty.MapValEmpty(cty.DynamicPseudoType), nil
		}
		return cty.ObjectVal(m), nil
	default:
		return cty.NullVal(cty.DynamicPseudoType), nil
	}
}

// ctyReplaceUnknownWithNullForHCL turns unknown cty values into null. convert.UnifyUnsafe can
// produce unknowns when unifying dynamic nulls with concrete types; hclwrite cannot serialize
// unknowns (TokensForValue panics).
func ctyReplaceUnknownWithNullForHCL(v cty.Value) (cty.Value, error) {
	return cty.Transform(v, func(_ cty.Path, val cty.Value) (cty.Value, error) {
		if val.IsNull() || val.IsKnown() {
			return val, nil
		}
		ty := val.Type().WithoutOptionalAttributesDeep()
		return cty.NullVal(ty), nil
	})
}

// search dashboard elements nested object uses objectvalidator.ConflictsWith across these;
// explicit null still counts as "specified" in Terraform, so emitted HCL must omit absent branches.
var searchDashboardElementOneOfAttrs = []string{
	"dashboard_element",
	"dashboard_element_input",
	"dashboard_element_visualization",
}

// PruneSearchDashboardElementsCty keeps at most one of the exclusive element branches per list item.
func PruneSearchDashboardElementsCty(val cty.Value) (cty.Value, error) {
	if val.IsNull() || !val.IsKnown() {
		return val, nil
	}
	ty := val.Type()
	if !ty.IsListType() && !ty.IsTupleType() {
		return val, nil
	}
	l := val.LengthInt()
	out := make([]cty.Value, l)
	for i := 0; i < l; i++ {
		ev := val.Index(cty.NumberIntVal(int64(i)))
		var err error
		out[i], err = ctyPruneDashboardElementUnionObject(ev)
		if err != nil {
			return cty.NilVal, fmt.Errorf("elements[%d]: %w", i, err)
		}
	}
	if ty.IsListType() && cty.CanListVal(out) {
		return cty.ListVal(out), nil
	}
	return cty.TupleVal(out), nil
}

// StripSearchDashboardElementsConfigNullKeysCty removes null-valued entries from
// dashboard_element{,_input,_visualization}.config maps on each element.
//
// List homogenization (normalizeListOfMaps / homogenizeListElementObjectsUntilStable)
// unions keys across elements and inserts explicit null for missing keys so cty list
// types line up. For open-ended JSON config maps that is undesirable: Terraform
// then treats those as real nulls and the provider sees perpetual drift. Omitting
// keys is the correct representation for "not set".
func StripSearchDashboardElementsConfigNullKeysCty(val cty.Value) (cty.Value, error) {
	if val.IsNull() || !val.IsKnown() {
		return val, nil
	}
	ty := val.Type()
	if !ty.IsListType() && !ty.IsTupleType() {
		return val, nil
	}
	l := val.LengthInt()
	out := make([]cty.Value, l)
	for i := 0; i < l; i++ {
		ev := val.Index(cty.NumberIntVal(int64(i)))
		var err error
		out[i], err = stripSearchDashboardElementConfigNullKeys(ev)
		if err != nil {
			return cty.NilVal, fmt.Errorf("elements[%d]: %w", i, err)
		}
	}
	if ty.IsListType() && cty.CanListVal(out) {
		return cty.ListVal(out), nil
	}
	return cty.TupleVal(out), nil
}

func stripSearchDashboardElementConfigNullKeys(elem cty.Value) (cty.Value, error) {
	if elem.IsNull() || !elem.IsKnown() || !elem.Type().IsObjectType() {
		return elem, nil
	}
	out := elem
	for _, branch := range searchDashboardElementOneOfAttrs {
		if !out.Type().HasAttribute(branch) {
			continue
		}
		b := out.GetAttr(branch)
		if b.IsNull() || !b.IsKnown() || !b.Type().IsObjectType() || !b.Type().HasAttribute("config") {
			continue
		}
		cfg := b.GetAttr("config")
		newCfg := ctyStripNullKeyedEntries(cfg)
		if newCfg.RawEquals(cfg) {
			continue
		}
		b2 := ctyObjectSetAttr(b, "config", newCfg)
		out = ctyObjectSetAttr(out, branch, b2)
	}
	return out, nil
}

// ctyStripNullKeyedEntries drops map keys / object attributes whose value is null.
func ctyStripNullKeyedEntries(v cty.Value) cty.Value {
	if v.IsNull() || !v.IsKnown() {
		return v
	}
	t := v.Type()
	switch {
	case t.IsMapType():
		if v.LengthInt() == 0 {
			return v
		}
		em := v.AsValueMap()
		out := make(map[string]cty.Value, len(em))
		for k, ev := range em {
			if ev.IsNull() {
				continue
			}
			out[k] = ev
		}
		if len(out) == 0 {
			return cty.MapValEmpty(t.ElementType())
		}
		return cty.MapVal(out)
	case t.IsObjectType():
		atys := t.AttributeTypes()
		out := make(map[string]cty.Value, len(atys))
		for k := range atys {
			ev := v.GetAttr(k)
			if ev.IsNull() {
				continue
			}
			out[k] = ev
		}
		if len(out) == 0 {
			return cty.ObjectVal(map[string]cty.Value{})
		}
		return cty.ObjectVal(out)
	default:
		return v
	}
}

func ctyPruneDashboardElementUnionObject(elem cty.Value) (cty.Value, error) {
	if elem.IsNull() || !elem.IsKnown() {
		return elem, nil
	}
	if !elem.Type().IsObjectType() {
		return elem, nil
	}
	atys := elem.Type().AttributeTypes()
	// Homogenization can leave multiple exclusive branches non-null (e.g. a real
	// dashboard_element_visualization plus an empty dashboard_element_input shell).
	// Pick the branch with substantive content, not a fixed order that favors input.
	topTie := []string{
		"dashboard_element",
		"dashboard_element_visualization",
		"dashboard_element_input",
	}
	chosen, err := ctyPickExclusiveOneOf(elem, searchDashboardElementOneOfAttrs, searchDashboardBranchSubstanceScore, topTie, true)
	if err != nil {
		return cty.NilVal, err
	}
	newAttrs := make(map[string]cty.Value, len(atys))
	for name := range atys {
		v := elem.GetAttr(name)
		if !searchDashboardElementOneOfAttr(name) {
			newAttrs[name] = v
			continue
		}
		if chosen != "" && name == chosen {
			cleaned, cerr := pruneSearchDashboardBranchContents(name, v)
			if cerr != nil {
				return cty.NilVal, cerr
			}
			newAttrs[name] = cleaned
		}
	}
	return cty.ObjectVal(newAttrs), nil
}

func searchDashboardElementOneOfAttr(name string) bool {
	for _, o := range searchDashboardElementOneOfAttrs {
		if o == name {
			return true
		}
	}
	return false
}

// ctyValueSubstanceScore counts known, non-null leaf values (for picking real union branches).
func ctyValueSubstanceScore(v cty.Value) int {
	if v.IsNull() || !v.IsKnown() {
		return 0
	}
	t := v.Type()
	switch {
	case t.IsPrimitiveType():
		return 1
	case t.IsObjectType():
		sum := 0
		for name := range t.AttributeTypes() {
			sum += ctyValueSubstanceScore(v.GetAttr(name))
		}
		return sum
	case t.IsTupleType() || t.IsListType():
		if v.LengthInt() == 0 {
			return 0
		}
		sum := 0
		for i := 0; i < v.LengthInt(); i++ {
			sum += ctyValueSubstanceScore(v.Index(cty.NumberIntVal(int64(i))))
		}
		return sum
	default:
		return 0
	}
}

func searchDashboardBranchSubstanceScore(v cty.Value) int {
	s := ctyValueSubstanceScore(v)
	if v.Type().IsObjectType() && v.Type().HasAttribute("id") {
		id := v.GetAttr("id")
		if id.IsKnown() && !id.IsNull() && id.Type().Equals(cty.String) && id.AsString() != "" {
			s += 10000
		}
	}
	return s
}

// ctyPickExclusiveOneOf picks at most one of attrOrder on obj. If allowZeroScoreSingleton is
// false and there is only one candidate with substance score 0, returns "" (caller may null
// the parent). tieBreak resolves equal top scores among multiple candidates.
func ctyPickExclusiveOneOf(obj cty.Value, attrOrder []string, score func(cty.Value) int, tieBreak []string, allowZeroScoreSingleton bool) (string, error) {
	if !obj.Type().IsObjectType() {
		return "", nil
	}
	var candidates []string
	for _, name := range attrOrder {
		if !obj.Type().HasAttribute(name) {
			continue
		}
		v := obj.GetAttr(name)
		if v.IsNull() || !v.IsKnown() {
			continue
		}
		candidates = append(candidates, name)
	}
	if len(candidates) == 0 {
		return "", nil
	}
	if len(candidates) == 1 {
		if !allowZeroScoreSingleton && score(obj.GetAttr(candidates[0])) == 0 {
			return "", nil
		}
		return candidates[0], nil
	}
	bestScore := -1
	var best []string
	for _, c := range candidates {
		s := score(obj.GetAttr(c))
		if s > bestScore {
			bestScore = s
			best = []string{c}
		} else if s == bestScore {
			best = append(best, c)
		}
	}
	if bestScore == 0 {
		return "", nil
	}
	order := tieBreak
	if len(order) == 0 {
		order = candidates
	}
	for _, name := range order {
		for _, b := range best {
			if b == name {
				return b, nil
			}
		}
	}
	return best[0], nil
}

func ctyObjectKeepOneExclusive(obj cty.Value, exclusive []string, keep string) (cty.Value, error) {
	if !obj.Type().IsObjectType() {
		return obj, nil
	}
	exSet := make(map[string]bool, len(exclusive))
	for _, e := range exclusive {
		exSet[e] = true
	}
	atys := obj.Type().AttributeTypes()
	newAttrs := make(map[string]cty.Value, len(atys))
	for name := range atys {
		v := obj.GetAttr(name)
		if exSet[name] {
			if keep != "" && name == keep {
				newAttrs[name] = v
			}
			continue
		}
		newAttrs[name] = v
	}
	return cty.ObjectVal(newAttrs), nil
}

func ctyObjectSetAttr(obj cty.Value, name string, val cty.Value) cty.Value {
	atys := obj.Type().AttributeTypes()
	m := make(map[string]cty.Value, len(atys))
	for n := range atys {
		if n == name {
			m[n] = val
		} else {
			m[n] = obj.GetAttr(n)
		}
	}
	return cty.ObjectVal(m)
}

func pruneSearchDashboardBranchContents(branch string, v cty.Value) (cty.Value, error) {
	switch branch {
	case "dashboard_element_input", "dashboard_element_visualization":
		return pruneDashElemInputOrVizCty(v)
	default:
		return v, nil
	}
}

func pruneDashElemInputOrVizCty(v cty.Value) (cty.Value, error) {
	if !v.Type().IsObjectType() {
		return v, nil
	}
	var err error
	v, err = pruneNestedAttrSearch(v)
	if err != nil {
		return cty.NilVal, err
	}
	if !v.Type().HasAttribute("config") {
		return v, nil
	}
	cfg := v.GetAttr("config")
	if cfg.IsNull() || !cfg.IsKnown() || !cfg.Type().IsObjectType() || !cfg.Type().HasAttribute("default_value") {
		return v, nil
	}
	newCfg, err := pruneInputConfigDefaultValueUnion(cfg)
	if err != nil {
		return cty.NilVal, err
	}
	return ctyObjectSetAttr(v, "config", newCfg), nil
}

func pruneNestedAttrSearch(v cty.Value) (cty.Value, error) {
	if !v.Type().HasAttribute("search") {
		return v, nil
	}
	s := v.GetAttr("search")
	if s.IsNull() || !s.IsKnown() {
		return v, nil
	}
	newS, err := pruneSearchQueryObjectCty(s)
	if err != nil {
		return cty.NilVal, err
	}
	if newS.IsNull() {
		return ctyObjectSetAttr(v, "search", cty.NullVal(s.Type())), nil
	}
	return ctyObjectSetAttr(v, "search", newS), nil
}

func pruneSearchQueryObjectCty(s cty.Value) (cty.Value, error) {
	exclusive := []string{"search_query_inline", "search_query_saved", "search_query_values"}
	tie := []string{"search_query_inline", "search_query_saved", "search_query_values"}
	chosen, err := ctyPickExclusiveOneOf(s, exclusive, ctyValueSubstanceScore, tie, false)
	if err != nil {
		return cty.NilVal, err
	}
	if chosen == "" {
		return cty.NullVal(s.Type()), nil
	}
	out, err := ctyObjectKeepOneExclusive(s, exclusive, chosen)
	if err != nil {
		return cty.NilVal, err
	}
	if chosen != "search_query_inline" {
		return out, nil
	}
	inline := out.GetAttr("search_query_inline")
	inline2, err := pruneSearchQueryInlineNested(inline)
	if err != nil {
		return cty.NilVal, err
	}
	return ctyObjectSetAttr(out, "search_query_inline", inline2), nil
}

func pruneSearchQueryInlineNested(v cty.Value) (cty.Value, error) {
	if !v.Type().IsObjectType() {
		return v, nil
	}
	out := v
	strNumTie := []string{"str", "number"}
	for _, path := range []string{"earliest", "latest"} {
		if !out.Type().HasAttribute(path) {
			continue
		}
		sub := out.GetAttr(path)
		if sub.IsNull() || !sub.IsKnown() {
			continue
		}
		newSub, err := ctyPruneExclusiveStrNumber(sub, strNumTie)
		if err != nil {
			return cty.NilVal, err
		}
		out = ctyObjectSetAttr(out, path, newSub)
	}
	return out, nil
}

func ctyPruneExclusiveStrNumber(obj cty.Value, tie []string) (cty.Value, error) {
	exclusive := []string{"str", "number"}
	chosen, err := ctyPickExclusiveOneOf(obj, exclusive, ctyValueSubstanceScore, tie, false)
	if err != nil {
		return cty.NilVal, err
	}
	if chosen == "" {
		return cty.NullVal(obj.Type()), nil
	}
	return ctyObjectKeepOneExclusive(obj, exclusive, chosen)
}

func pruneInputConfigDefaultValueUnion(cfg cty.Value) (cty.Value, error) {
	dv := cfg.GetAttr("default_value")
	if dv.IsNull() || !dv.IsKnown() {
		return cfg, nil
	}
	exclusive := []string{"str", "number", "array_of_str", "default_value"}
	tie := []string{"default_value", "str", "number", "array_of_str"}
	chosen, err := ctyPickExclusiveOneOf(dv, exclusive, ctyValueSubstanceScore, tie, false)
	if err != nil {
		return cty.NilVal, err
	}
	if chosen == "" {
		return ctyObjectSetAttr(cfg, "default_value", cty.NullVal(dv.Type())), nil
	}
	newDV, err := ctyObjectKeepOneExclusive(dv, exclusive, chosen)
	if err != nil {
		return cty.NilVal, err
	}
	if chosen == "default_value" {
		inner := newDV.GetAttr("default_value")
		inner2, err := pruneDefaultValueTimerangeEarliestLatest(inner)
		if err != nil {
			return cty.NilVal, err
		}
		newDV = ctyObjectSetAttr(newDV, "default_value", inner2)
	}
	return ctyObjectSetAttr(cfg, "default_value", newDV), nil
}

func pruneDefaultValueTimerangeEarliestLatest(inner cty.Value) (cty.Value, error) {
	if inner.IsNull() || !inner.IsKnown() || !inner.Type().IsObjectType() {
		return inner, nil
	}
	out := inner
	strNumTie := []string{"str", "number"}
	for _, path := range []string{"earliest", "latest"} {
		if !out.Type().HasAttribute(path) {
			continue
		}
		sub := out.GetAttr(path)
		if sub.IsNull() || !sub.IsKnown() {
			continue
		}
		newSub, err := ctyPruneExclusiveStrNumber(sub, strNumTie)
		if err != nil {
			return cty.NilVal, err
		}
		out = ctyObjectSetAttr(out, path, newSub)
	}
	return out, nil
}

// normalizeAndPruneList normalizes list-of-maps (union keys), prunes nulls, and re-normalizes.
// cty requires all list elements to have the same type; optional/omitted attributes cause object types to differ.
func normalizeAndPruneList(v Value) Value {
	v = normalizeListOfMaps(v)
	v = PruneNulls(v)
	v = normalizeListOfMaps(v)
	return v
}

// listToCtyElems converts each list element to cty.Value. Do not call PruneNulls here:
// normalizeAndPruneList already aligned keys across elements with explicit nulls; pruning
// again drops union branches and yields inconsistent object types per element (cty.ListVal panics).
func listToCtyElems(list []Value) ([]cty.Value, error) {
	elems := make([]cty.Value, 0, len(list))
	for i, el := range list {
		c, err := ValueToCty(el)
		if err != nil {
			return nil, fmt.Errorf("list[%d]: %w", i, err)
		}
		elems = append(elems, c)
	}
	return elems, nil
}

// unifyListElementTypes ensures all elements have the same cty type so cty.ListVal succeeds.
// Same keys can yield different attribute types (e.g. null -> DynamicPseudoType vs string -> String).
func unifyListElementTypes(elems []cty.Value) ([]cty.Value, error) {
	if len(elems) <= 1 {
		return elems, nil
	}
	types := make([]cty.Type, len(elems))
	for i := range elems {
		types[i] = elems[i].Type()
	}
	unified, convs := convert.UnifyUnsafe(types)
	if unified != cty.NilType && convs != nil {
		for j := range elems {
			if convs[j] != nil {
				converted, err := convs[j](elems[j])
				if err != nil {
					return nil, fmt.Errorf("list[%d]: %w", j, err)
				}
				elems[j] = converted
			}
		}
	}
	if unified != cty.NilType && allSameType(elems) {
		return elems, nil
	}
	if cty.CanListVal(elems) {
		return elems, nil
	}
	// convert.Convert(_, DynamicPseudoType) is a no-op in go-cty; leave elems for
	// ValueToCty to emit TupleVal when CanListVal is still false.
	return elems, nil
}

// homogenizeListElementObjectsUntilStable walks list elements that are objects (or null)
// and recursively aligns attribute types across elements (union of keys + convert.Unify
// per attribute column). This fixes one-of / union maps where alternating branches yield
// object types like {a: O, b: dynamic null} vs {a: dynamic null, b: O} — incompatible for
// cty.ListVal because Dynamic null does not coalesce list element types (see ListVal).
func homogenizeListElementObjectsUntilStable(elems []cty.Value) ([]cty.Value, error) {
	if len(elems) <= 1 {
		return elems, nil
	}
	const maxPasses = 64
	prevSig := ""
	for pass := 0; pass < maxPasses; pass++ {
		if cty.CanListVal(elems) {
			return elems, nil
		}
		sig := valueTypesSignature(elems)
		if sig == prevSig {
			break
		}
		prevSig = sig
		var err error
		elems, err = homogenizeObjectListElementsOnePass(elems)
		if err != nil {
			return nil, err
		}
	}
	return elems, nil
}

func valueTypesSignature(elems []cty.Value) string {
	var b []byte
	for _, e := range elems {
		b = append(b, e.Type().FriendlyName()...)
		b = append(b, '|')
	}
	return string(b)
}

func homogenizeObjectListElementsOnePass(elems []cty.Value) ([]cty.Value, error) {
	if len(elems) == 0 {
		return elems, nil
	}
	keySet := make(map[string]struct{})
	for _, e := range elems {
		if e.IsNull() {
			continue
		}
		if !e.Type().IsObjectType() {
			return elems, nil
		}
		for k := range e.Type().AttributeTypes() {
			keySet[k] = struct{}{}
		}
	}
	if len(keySet) == 0 {
		return elems, nil
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	columns := make(map[string][]cty.Value, len(keys))
	for _, k := range keys {
		col := make([]cty.Value, len(elems))
		for i, e := range elems {
			if !e.IsNull() && e.Type().IsObjectType() && e.Type().HasAttribute(k) {
				col[i] = e.GetAttr(k)
			} else {
				col[i] = cty.NullVal(cty.DynamicPseudoType)
			}
		}
		var err error
		col, err = unifySiblingValues(col)
		if err != nil {
			return nil, err
		}
		col, err = homogenizeColumnValuesDeep(col)
		if err != nil {
			return nil, err
		}
		columns[k] = col
	}

	out := make([]cty.Value, len(elems))
	for i := range elems {
		m := make(map[string]cty.Value, len(keys))
		for _, k := range keys {
			m[k] = columns[k][i]
		}
		out[i] = cty.ObjectVal(m)
	}
	return out, nil
}

// unifySiblingValues aligns types for one attribute across list elements (a "column").
//
// Homogenization fills missing union branches with NullVal(DynamicPseudoType). If we pass
// those placeholders into convert.UnifyUnsafe alongside concrete object types, go-cty sets
// hasDynamic and uses unifyAllAsDynamic — which turns known objects into unknown values and
// breaks export (unknown → null in HCL). Skip dynamic-null placeholders when unifying;
// after unifying the concrete cells, rewrite placeholders to NullVal(unifiedType).
func unifySiblingValues(col []cty.Value) ([]cty.Value, error) {
	if len(col) == 0 {
		return col, nil
	}
	concreteTypes := make([]cty.Type, 0, len(col))
	concreteIdx := make([]int, 0, len(col))
	for i, v := range col {
		if v.IsNull() && v.Type() == cty.DynamicPseudoType {
			continue
		}
		if !v.IsKnown() {
			continue
		}
		concreteTypes = append(concreteTypes, v.Type())
		concreteIdx = append(concreteIdx, i)
	}
	if len(concreteTypes) == 0 {
		return col, nil
	}
	unified, convs := convert.UnifyUnsafe(concreteTypes)
	if unified == cty.NilType || convs == nil {
		return col, nil
	}
	out := make([]cty.Value, len(col))
	copy(out, col)
	for j, idx := range concreteIdx {
		if convs[j] != nil {
			cv, err := convs[j](col[idx])
			if err != nil {
				return nil, fmt.Errorf("column unify [%d]: %w", idx, err)
			}
			out[idx] = cv
		} else {
			out[idx] = col[idx]
		}
	}
	nullFill := cty.NullVal(unified.WithoutOptionalAttributesDeep())
	for i := range out {
		if out[i].IsNull() && out[i].Type() == cty.DynamicPseudoType {
			out[i] = nullFill
		}
	}
	return out, nil
}

func homogenizeColumnValuesDeep(col []cty.Value) ([]cty.Value, error) {
	if !columnAllNullOrObject(col) {
		return homogenizeListColumnSameLength(col)
	}
	return homogenizeObjectListElementsOnePass(col)
}

func columnAllNullOrObject(col []cty.Value) bool {
	seenObject := false
	for _, v := range col {
		if v.IsNull() {
			continue
		}
		if !v.Type().IsObjectType() {
			return false
		}
		seenObject = true
	}
	return seenObject
}

func homogenizeListColumnSameLength(col []cty.Value) ([]cty.Value, error) {
	n := -1
	for _, v := range col {
		if v.IsNull() {
			continue
		}
		if !v.Type().IsListType() {
			return col, nil
		}
		ln := v.LengthInt()
		if n < 0 {
			n = ln
		} else if n != ln {
			return col, nil
		}
	}
	if n <= 0 {
		return col, nil
	}

	rebuilt := make([][]cty.Value, n)
	for j := 0; j < n; j++ {
		slice := make([]cty.Value, len(col))
		for i, v := range col {
			if v.IsNull() {
				slice[i] = cty.NullVal(cty.DynamicPseudoType)
			} else {
				slice[i] = v.Index(cty.NumberIntVal(int64(j)))
			}
		}
		var err error
		slice, err = homogenizeListElementObjectsUntilStable(slice)
		if err != nil {
			return nil, err
		}
		rebuilt[j] = slice
	}

	out := make([]cty.Value, len(col))
	for i := range col {
		vals := make([]cty.Value, n)
		for j := 0; j < n; j++ {
			vals[j] = rebuilt[j][i]
		}
		if cty.CanListVal(vals) {
			out[i] = cty.ListVal(vals)
		} else {
			out[i] = cty.TupleVal(vals)
		}
	}
	return out, nil
}

// normalizeListOfMaps ensures every element in a list of maps has the same set of keys (union).
// Missing keys are set to null. This allows cty.ListVal to succeed (all elements same object type).
// Nested lists of maps are normalized recursively.
func normalizeListOfMaps(v Value) Value {
	if v.Kind != KindList || len(v.List) == 0 {
		return v
	}
	allMaps := true
	for _, el := range v.List {
		if el.Kind != KindMap {
			allMaps = false
			break
		}
	}
	if !allMaps {
		// Recursively normalize nested lists (e.g. list of maps containing lists).
		out := make([]Value, len(v.List))
		for i, el := range v.List {
			if el.Kind == KindList {
				out[i] = normalizeListOfMaps(el)
			} else if el.Kind == KindMap && el.Map != nil {
				m := make(map[string]Value, len(el.Map))
				for k, val := range el.Map {
					m[k] = normalizeListOfMaps(val)
				}
				out[i] = Value{Kind: KindMap, Map: m}
			} else {
				out[i] = el
			}
		}
		return Value{Kind: KindList, List: out}
	}
	// Collect union of keys.
	keySet := make(map[string]struct{})
	for _, el := range v.List {
		for k := range el.Map {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nullVal := Value{Kind: KindNull}
	out := make([]Value, 0, len(v.List))
	for _, el := range v.List {
		m := make(map[string]Value, len(keys))
		for _, k := range keys {
			if val, ok := el.Map[k]; ok {
				m[k] = normalizeListOfMaps(val)
			} else {
				m[k] = nullVal
			}
		}
		out = append(out, Value{Kind: KindMap, Map: m})
	}
	return Value{Kind: KindList, List: out}
}

func allSameType(elems []cty.Value) bool {
	if len(elems) <= 1 {
		return true
	}
	t := elems[0].Type()
	for i := 1; i < len(elems); i++ {
		if !elems[i].Type().Equals(t) {
			return false
		}
	}
	return true
}
