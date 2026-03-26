// Package export converts discovery results into generator ResourceItems.
package export

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	ptypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/custom"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/hcl"
	"github.com/criblio/terraform-provider-criblio/tools/import-cli/internal/registry"
)

// addOneOfBlockFromFirstItem extracts the first element from the model's read-only list (e.g. Items)
// and converts it to the type-specific block using the oneOf config, merging into attrs.
// When SupportedBlockNames is set, resolves API type to a supported block (or returns ErrUnsupportedOneOfType).
// When UnsupportedDiscriminatorValues is set, those types are rejected before further resolution.
// When NestedDiscriminatorField is set and the top-level discriminator doesn't resolve,
// the exporter parses the nested object to find the real discriminator (e.g. collector.type).
func addOneOfBlockFromFirstItem(model interface{}, attrs map[string]hcl.Value, oneOf *registry.OneOfConfig) error {
	itemMap := firstItemMapFromModel(model, oneOf.ReadOnlyAttr)
	if len(itemMap) == 0 {
		// Source and pack_source use []InputUnion1 (struct slices); firstItemMapFromModel only handles
		// []map[string]jsontypes.Normalized. Build the input_* block from the non-nil union branch.
		return addOneOfFromInputUnionStruct(model, attrs, oneOf)
	}
	raw := itemMap[oneOf.DiscriminatorField]
	var discStr string
	if err := json.Unmarshal([]byte(raw), &discStr); err != nil {
		discStr = strings.Trim(raw, `"`)
	}
	for _, unsup := range oneOf.UnsupportedDiscriminatorValues {
		if discStr == unsup {
			return ErrUnsupportedOneOfType
		}
	}
	var alias map[string]string
	if len(oneOf.SupportedBlockNames) > 0 {
		suffix, ok := hcl.ResolveOneOfBlockNameRaw(raw, oneOf.SupportedBlockNames, oneOf.BlockNamePrefix)
		if !ok && oneOf.NestedDiscriminatorField != "" {
			nestedRaw := resolveNestedDiscriminator(itemMap, oneOf.NestedDiscriminatorField)
			if nestedRaw != "" {
				suffix, ok = hcl.ResolveOneOfBlockNameRaw(nestedRaw, oneOf.SupportedBlockNames, oneOf.BlockNamePrefix)
			}
		}
		if !ok {
			return ErrUnsupportedOneOfType
		}
		alias = map[string]string{discStr: suffix}
	} else {
		alias = oneOf.DiscriminatorAlias
	}
	blockName, blockValue, err := hcl.ItemMapToBlock(itemMap, oneOf.DiscriminatorField, oneOf.BlockNamePrefix, oneOf.BlockNameSuffix, oneOf.KeysToSkip, alias)
	if err != nil {
		return err
	}
	if blockName != "" && !blockValue.IsNull() {
		attrs[blockName] = blockValue
	}
	return nil
}

func attrsHasOutputBlock(attrs map[string]hcl.Value) bool {
	for k := range attrs {
		if strings.HasPrefix(k, "output_") {
			return true
		}
	}
	return false
}

// addPackDestinationOneOfFromStoredItem adds the oneOf block (e.g. output_cribl_lake) to attrs from the
// raw API response stored by the converter when the model has no Items field.
func addPackDestinationOneOfFromStoredItem(idMap map[string]string, attrs map[string]hcl.Value, oneOf *registry.OneOfConfig) {
	if oneOf == nil {
		return
	}
	groupID := idMap["group_id"]
	pack := idMap["pack"]
	id := idMap["id"]
	if groupID == "" || pack == "" || id == "" {
		return
	}
	itemJSON, ok := custom.GetAndClearPackOutputFirstItem(groupID, pack, id)
	if !ok || len(itemJSON) == 0 {
		return
	}
	itemMap := rawJSONToItemMap(itemJSON)
	if len(itemMap) == 0 {
		return
	}
	blockName, blockValue, err := hcl.ItemMapToBlock(itemMap, oneOf.DiscriminatorField, oneOf.BlockNamePrefix, oneOf.BlockNameSuffix, oneOf.KeysToSkip, oneOf.DiscriminatorAlias)
	if err != nil || blockName == "" || blockValue.IsNull() {
		return
	}
	attrs[blockName] = blockValue
}

// resolveNestedDiscriminator parses a dot-separated field path (e.g. "collector.type")
// from the item map and returns the raw JSON value of the inner field.
func resolveNestedDiscriminator(itemMap map[string]string, path string) string {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	parentRaw, ok := itemMap[parts[0]]
	if !ok || parentRaw == "" {
		return ""
	}
	var parentObj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(parentRaw), &parentObj); err != nil {
		return ""
	}
	inner, ok := parentObj[parts[1]]
	if !ok {
		return ""
	}
	return string(inner)
}

// addOneOfFromInputUnionStruct emits input_<type> from Items[0] when the model uses InputUnion1 structs
// (criblio_source, criblio_pack_source). The map-based path in firstItemMapFromModel does not apply.
func addOneOfFromInputUnionStruct(model interface{}, attrs map[string]hcl.Value, oneOf *registry.OneOfConfig) error {
	if oneOf == nil {
		return nil
	}
	var items []ptypes.InputUnion1
	switch m := model.(type) {
	case *provider.SourceResourceModel:
		items = m.Items
	case *provider.PackSourceResourceModel:
		items = m.Items
	default:
		return nil
	}
	if len(items) == 0 {
		return nil
	}
	u := items[0]
	val := reflect.ValueOf(u)
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.Kind() != reflect.Ptr || f.IsNil() {
			continue
		}
		innerPtr := f.Interface()
		blockMap, err := hcl.ModelToValue(innerPtr, nil)
		if err != nil {
			return err
		}
		itemMap, err := hcl.TFBlockModelToAPIItemMap(blockMap, oneOf.KeysToSkip)
		if err != nil {
			return err
		}
		if len(itemMap) == 0 {
			return nil
		}
		raw := itemMap[oneOf.DiscriminatorField]
		if raw == "" {
			return fmt.Errorf("input union branch missing discriminator %q", oneOf.DiscriminatorField)
		}
		var discStr string
		if err := json.Unmarshal([]byte(raw), &discStr); err != nil {
			discStr = strings.Trim(raw, `"`)
		}
		for _, unsup := range oneOf.UnsupportedDiscriminatorValues {
			if discStr == unsup {
				return ErrUnsupportedOneOfType
			}
		}
		var alias map[string]string
		if len(oneOf.SupportedBlockNames) > 0 {
			suffix, ok := hcl.ResolveOneOfBlockNameRaw(raw, oneOf.SupportedBlockNames, oneOf.BlockNamePrefix)
			if !ok {
				return ErrUnsupportedOneOfType
			}
			alias = map[string]string{discStr: suffix}
		} else {
			alias = oneOf.DiscriminatorAlias
		}
		blockName, blockValue, err := hcl.ItemMapToBlock(itemMap, oneOf.DiscriminatorField, oneOf.BlockNamePrefix, oneOf.BlockNameSuffix, oneOf.KeysToSkip, alias)
		if err != nil {
			return err
		}
		if blockName != "" && !blockValue.IsNull() {
			attrs[blockName] = blockValue
		}
		return nil
	}
	return nil
}
