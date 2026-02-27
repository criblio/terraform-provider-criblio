// Terraform resource block generation using hclwrite.
package hcl

import (
	"bytes"
	"fmt"
	"sort"

	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ResourceBlockOptions configures resource block generation.
type ResourceBlockOptions struct {
	// SkipNullAttributes when true omits attributes that are null from the block.
	SkipNullAttributes bool
	// SkipEmptyListAttributes when true omits attributes that are empty lists, avoiding
	// schema validation errors (e.g. listvalidator.SizeAtLeast(1)).
	SkipEmptyListAttributes bool
	// SkipPruneEmptyListsFor: when SkipEmptyListAttributes is true, do not prune empty lists
	// for these resource type + attribute pairs. E.g. criblio_pack tags must keep domain=[].
	// Map key = resource type name, value = attribute names.
	SkipPruneEmptyListsFor map[string][]string
	// AlwaysEmitEmptyListsFor: when SkipEmptyListAttributes is true, always emit these
	// attributes even when empty (empty list is valid). E.g. criblio_project subscriptions, destinations.
	// Map key = resource type name, value = attribute names.
	AlwaysEmitEmptyListsFor map[string][]string
}

// DefaultResourceBlockOptions returns options suitable for import CLI output:
// skip nulls, skip empty lists, with resource-specific rules for criblio_pack and criblio_project.
func DefaultResourceBlockOptions() *ResourceBlockOptions {
	return &ResourceBlockOptions{
		SkipNullAttributes:       true,
		SkipEmptyListAttributes:  true,
		SkipPruneEmptyListsFor:   map[string][]string{"criblio_pack": {"tags"}},
		AlwaysEmitEmptyListsFor:  map[string][]string{"criblio_project": {"subscriptions", "destinations"}},
	}
}

func attrInList(typeName, attrName string, m map[string][]string) bool {
	if m == nil {
		return false
	}
	for _, a := range m[typeName] {
		if a == attrName {
			return true
		}
	}
	return false
}

// ResourceBlock builds an hclwrite Block for a Terraform resource:
//   resource "typeName" "name" { ... }
// Attributes are taken from attrs (e.g. from ModelToValue). Nested attributes
// are rendered as nested blocks/objects. The returned block can be appended to
// an hclwrite.File body.
// When lifecycleIgnoreChanges is non-nil and non-empty, a lifecycle block with
// ignore_changes = lifecycleIgnoreChanges is appended (e.g. for cloud default group).
func ResourceBlock(typeName, name string, attrs map[string]Value, opts *ResourceBlockOptions, lifecycleIgnoreChanges []string) (*hclwrite.Block, error) {
	if typeName == "" || name == "" {
		return nil, fmt.Errorf("resource type and name are required")
	}
	block := hclwrite.NewBlock("resource", []string{typeName, name})
	body := block.Body()
	skipNull := opts != nil && opts.SkipNullAttributes

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := attrs[k]
		if opts != nil && opts.SkipEmptyListAttributes && !attrInList(typeName, k, opts.SkipPruneEmptyListsFor) {
			v = PruneEmptyLists(v)
		}
		if skipNull {
			v = PruneNulls(v)
		}
		if skipNull && v.Kind == KindNull {
			continue
		}
		if skipNull && v.Kind == KindMap && len(v.Map) == 0 {
			continue
		}
		if opts != nil && opts.SkipEmptyListAttributes && v.Kind == KindList && len(v.List) == 0 {
			if !attrInList(typeName, k, opts.AlwaysEmitEmptyListsFor) {
				continue
			}
		}
		ctyVal, err := ValueToCty(v)
		if err != nil {
			return nil, fmt.Errorf("attribute %q: %w", k, err)
		}
		body.SetAttributeValue(k, ctyVal)
	}
	if len(lifecycleIgnoreChanges) > 0 {
		lifecycleBlock := body.AppendNewBlock("lifecycle", nil)
		lifecycleBody := lifecycleBlock.Body()
		// Use unquoted references [api, items] not ["api", "items"] to avoid deprecation warning.
		elems := make([]hclwrite.Tokens, len(lifecycleIgnoreChanges))
		for i, attr := range lifecycleIgnoreChanges {
			elems[i] = hclwrite.TokensForIdentifier(attr)
		}
		lifecycleBody.SetAttributeRaw("ignore_changes", hclwrite.TokensForTuple(elems))
	}
	return block, nil
}

// ResourceInput is the input for building a single resource block.
type ResourceInput struct {
	TypeName                string
	Name                    string
	Attrs                   map[string]Value
	LifecycleIgnoreChanges  []string
}

// ResourceBlockBytes returns the HCL source for a single resource block (with newline).
// Useful for writing to a file or parsing in tests.
func ResourceBlockBytes(typeName, name string, attrs map[string]Value, opts *ResourceBlockOptions) ([]byte, error) {
	block, err := ResourceBlock(typeName, name, attrs, opts, nil)
	if err != nil {
		return nil, err
	}
	f := hclwrite.NewEmptyFile()
	f.Body().AppendBlock(block)
	return f.Bytes(), nil
}

// AppendResourceBlock appends a resource block to the given file body.
func AppendResourceBlock(body *hclwrite.Body, r ResourceInput, opts *ResourceBlockOptions) error {
	block, err := ResourceBlock(r.TypeName, r.Name, r.Attrs, opts, r.LifecycleIgnoreChanges)
	if err != nil {
		return err
	}
	body.AppendBlock(block)
	return nil
}

// FileWithResources builds an hclwrite.File containing the given resource blocks.
// Resources are sorted by TypeName then Name so output is deterministic.
// Secret variable placeholders (__VAR_REF__name__) in the file are replaced with jsonencode(var.name).
func FileWithResources(resources []ResourceInput, opts *ResourceBlockOptions) (*hclwrite.File, error) {
	// Sort for deterministic output (same input → same output)
	sorted := make([]ResourceInput, len(resources))
	copy(sorted, resources)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].TypeName != sorted[j].TypeName {
			return sorted[i].TypeName < sorted[j].TypeName
		}
		return sorted[i].Name < sorted[j].Name
	})
	f := hclwrite.NewEmptyFile()
	for _, r := range sorted {
		if err := AppendResourceBlock(f.Body(), r, opts); err != nil {
			return nil, fmt.Errorf("resource %s.%s: %w", r.TypeName, r.Name, err)
		}
	}
	// Replace plain variable refs (e.g. criblio_secret.value) with var.xxx
	plainNames := make([]string, 0)
	plainSeen := make(map[string]bool)
	for _, r := range sorted {
		for _, n := range CollectPlainVariableRefNames(r.Attrs) {
			if !plainSeen[n] {
				plainSeen[n] = true
				plainNames = append(plainNames, n)
			}
		}
	}
	b := f.Bytes()
	for _, n := range plainNames {
		quoted := `"` + PlainVarRefPlaceholderPrefix + n + PlainVarRefPlaceholderSuffix + `"`
		b = bytes.ReplaceAll(b, []byte(quoted), []byte("var."+n))
	}
	// Replace secret variable placeholders with jsonencode(var.xxx) so secrets are always JSON-encoded
	names := make([]string, 0)
	seen := make(map[string]bool)
	for _, r := range sorted {
		for _, n := range CollectSecretVariableNames(r.Attrs) {
			if !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
	}
	for _, n := range names {
		quoted := `"` + VarRefPlaceholderPrefix + n + VarRefPlaceholderSuffix + `"`
		b = bytes.ReplaceAll(b, []byte(quoted), []byte("jsonencode(var."+n+")"))
	}
	// Re-parse so the file reflects the replacement (for callers that use the file object)
	parsed, _ := hclwrite.ParseConfig(b, "", hclv2.Pos{Line: 1, Column: 1})
	if parsed != nil {
		return parsed, nil
	}
	return f, nil
}
