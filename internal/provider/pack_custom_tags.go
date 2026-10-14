package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func builtInPackTagKey(key string) bool {
	switch key {
	case "dataType", "data_type", "domain", "streamtags", "technology":
		return true
	}
	return false
}

type packTagsValidator struct{}

func (packTagsValidator) Description(context.Context) string {
	return "Pack tags must be an object; built-in fields must be string lists."
}

func (value packTagsValidator) MarkdownDescription(ctx context.Context) string {
	return value.Description(ctx)
}

func (packTagsValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || packTagHasUnknown(req.ConfigValue) {
		return
	}
	if _, err := tagsAPIMap(ctx, req.ConfigValue); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid pack tags", err.Error())
	}
}

func packTagHasUnknown(value attr.Value) bool {
	if value.IsUnknown() {
		return true
	}
	if value.IsNull() {
		return false
	}
	switch value := value.(type) {
	case types.Dynamic:
		return packTagHasUnknown(value.UnderlyingValue())
	case types.Object:
		for _, child := range value.Attributes() {
			if packTagHasUnknown(child) {
				return true
			}
		}
	case types.Map:
		for _, child := range value.Elements() {
			if packTagHasUnknown(child) {
				return true
			}
		}
	case types.List:
		for _, child := range value.Elements() {
			if packTagHasUnknown(child) {
				return true
			}
		}
	case types.Tuple:
		for _, child := range value.Elements() {
			if packTagHasUnknown(child) {
				return true
			}
		}
	case types.Set:
		for _, child := range value.Elements() {
			if packTagHasUnknown(child) {
				return true
			}
		}
	}
	return false
}

func (tags *packTagsAPI) UnmarshalJSON(raw []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	var decoded packTagsAPI
	for key, target := range map[string]*[]string{
		"dataType": &decoded.DataType, "domain": &decoded.Domain,
		"streamtags": &decoded.Streamtags, "technology": &decoded.Technology,
	} {
		if field, exists := fields[key]; exists {
			if err := json.Unmarshal(field, target); err != nil {
				return fmt.Errorf("decode pack tag %q: %w", key, err)
			}
		}
	}
	values, err := decodePackTagValues(raw)
	if err != nil {
		return err
	}
	decoded.Raw = values
	*tags = decoded
	return nil
}

func packTagAttributes(tags types.Dynamic) (map[string]attr.Value, error) {
	if tags.IsNull() {
		return nil, nil
	}
	if tags.IsUnknown() || tags.IsUnderlyingValueUnknown() {
		return nil, fmt.Errorf("pack tags must be known before updating settings")
	}
	switch value := tags.UnderlyingValue().(type) {
	case types.Object:
		return value.Attributes(), nil
	case types.Map:
		return value.Elements(), nil
	default:
		return nil, fmt.Errorf("pack tags must be an object or map")
	}
}

func tagsAPIMap(ctx context.Context, tags types.Dynamic) (map[string]any, error) {
	attributes, err := packTagAttributes(tags)
	if err != nil || attributes == nil {
		return nil, err
	}
	out := make(map[string]any, len(attributes))
	for key, value := range attributes {
		if strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("pack tag keys must not be empty")
		}
		if key == "dataType" {
			return nil, fmt.Errorf("use data_type instead of dataType in Terraform pack tags")
		}
		if value.IsNull() {
			if builtInPackTagKey(key) {
				continue
			}
			return nil, fmt.Errorf("pack tag %q must not be null; remove the key to delete it", key)
		}
		converted, err := packTagJSONValue(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("pack tag %q: %w", key, err)
		}
		if builtInPackTagKey(key) {
			elements, ok := converted.([]any)
			if !ok {
				return nil, fmt.Errorf("pack tag %q must be a list of strings", key)
			}
			for _, element := range elements {
				if _, ok := element.(string); !ok {
					return nil, fmt.Errorf("pack tag %q must contain only strings", key)
				}
			}
		}
		if key == "data_type" {
			key = "dataType"
		}
		out[key] = converted
	}
	return out, nil
}

func packTagJSONValue(ctx context.Context, value attr.Value) (any, error) {
	if value.IsUnknown() {
		return nil, fmt.Errorf("value must be known before updating settings")
	}
	if value.IsNull() {
		return nil, nil
	}
	switch value := value.(type) {
	case types.Dynamic:
		return packTagJSONValue(ctx, value.UnderlyingValue())
	case types.String:
		return value.ValueString(), nil
	case types.Bool:
		return value.ValueBool(), nil
	case types.Number:
		return json.Number(value.ValueBigFloat().Text('f', -1)), nil
	case types.Float64:
		return value.ValueFloat64(), nil
	case types.Int64:
		return value.ValueInt64(), nil
	case types.Object:
		return packTagJSONObject(ctx, value.Attributes())
	case types.Map:
		return packTagJSONObject(ctx, value.Elements())
	case types.List:
		return packTagJSONArray(ctx, value.Elements())
	case types.Tuple:
		return packTagJSONArray(ctx, value.Elements())
	case types.Set:
		return packTagJSONArray(ctx, value.Elements())
	default:
		return nil, fmt.Errorf("unsupported tag value type %T", value)
	}
}

func packTagJSONObject(ctx context.Context, values map[string]attr.Value) (map[string]any, error) {
	out := make(map[string]any, len(values))
	for key, value := range values {
		converted, err := packTagJSONValue(ctx, value)
		if err != nil {
			return nil, err
		}
		out[key] = converted
	}
	return out, nil
}

func packTagJSONArray(ctx context.Context, values []attr.Value) ([]any, error) {
	out := make([]any, 0, len(values))
	for _, value := range values {
		converted, err := packTagJSONValue(ctx, value)
		if err != nil {
			return nil, err
		}
		out = append(out, converted)
	}
	return out, nil
}

func readPackSettingsTags(ctx context.Context, client *restclient.Client, groupID, packID string) (map[string]any, bool, error) {
	settings, err := restclient.Get[struct {
		Package *struct {
			Tags json.RawMessage `json:"tags"`
		} `json:"package"`
	}](ctx, client, fmt.Sprintf("/m/%s/p/%s/pack/settings", url.PathEscape(groupID), url.PathEscape(packID)))
	if restclient.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read pack settings: %w", err)
	}
	if settings == nil || settings.Package == nil {
		return nil, false, fmt.Errorf("pack settings response contained no package")
	}
	if len(settings.Package.Tags) == 0 {
		return nil, true, nil
	}
	values, err := decodePackTagValues(settings.Package.Tags)
	if err != nil {
		return nil, false, fmt.Errorf("decode pack settings tags: %w", err)
	}
	return values, true, nil
}

func decodePackTags(values map[string]any) (*packTagsAPI, error) {
	raw, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("encode pack tags: %w", err)
	}
	var tags packTagsAPI
	if err := json.Unmarshal(raw, &tags); err != nil {
		return nil, fmt.Errorf("decode pack tags: %w", err)
	}
	return &tags, nil
}

func mergePackSettingsTags(remote, configured map[string]any) map[string]any {
	merged := make(map[string]any, len(remote)+len(configured))
	for key, value := range remote {
		if configured != nil && !builtInPackTagKey(key) {
			continue
		}
		merged[key] = value
	}
	for key, value := range configured {
		merged[key] = value
	}
	return merged
}

func packRequestBodyTagsObjectFromAPI(tags *packTagsAPI) types.Dynamic {
	if tags == nil {
		return types.DynamicNull()
	}
	values := tags.Raw
	if values == nil {
		values = map[string]any{
			"dataType":   nonNilStringSlice(tags.DataType),
			"domain":     nonNilStringSlice(tags.Domain),
			"streamtags": nonNilStringSlice(tags.Streamtags),
			"technology": nonNilStringSlice(tags.Technology),
		}
	}
	value, err := packTagsValue(values)
	if err != nil {
		return types.DynamicNull()
	}
	return value
}

func packTagsValue(values map[string]any) (types.Dynamic, error) {
	attributes := make(map[string]attr.Value, len(values))
	attributeTypes := make(map[string]attr.Type, len(values))
	for key, value := range values {
		if key == "dataType" {
			key = "data_type"
		}
		converted, err := packTagAttrValue(value)
		if err != nil {
			return types.Dynamic{}, err
		}
		attributes[key] = converted
		attributeTypes[key] = converted.Type(context.Background())
	}
	return types.DynamicValue(types.ObjectValueMust(attributeTypes, attributes)), nil
}

func packTagAttrValue(value any) (attr.Value, error) {
	switch value := value.(type) {
	case nil:
		return types.DynamicNull(), nil
	case string:
		return types.StringValue(value), nil
	case bool:
		return types.BoolValue(value), nil
	case float64:
		return types.Float64Value(value), nil
	case json.Number:
		number, _, err := big.ParseFloat(value.String(), 10, 512, big.ToNearestEven)
		if err != nil {
			return nil, err
		}
		return types.NumberValue(number), nil
	case []string:
		return types.ListValueMust(types.StringType, stringAttrValues(value)), nil
	case []any:
		elements := make([]attr.Value, 0, len(value))
		elementTypes := make([]attr.Type, 0, len(value))
		stringsOnly := true
		for _, element := range value {
			converted, err := packTagAttrValue(element)
			if err != nil {
				return nil, err
			}
			elements = append(elements, converted)
			elementTypes = append(elementTypes, converted.Type(context.Background()))
			if _, ok := converted.(types.String); !ok {
				stringsOnly = false
			}
		}
		if stringsOnly {
			return types.ListValueMust(types.StringType, elements), nil
		}
		return types.TupleValueMust(elementTypes, elements), nil
	case map[string]any:
		attributes := make(map[string]attr.Value, len(value))
		attributeTypes := make(map[string]attr.Type, len(value))
		for key, element := range value {
			converted, err := packTagAttrValue(element)
			if err != nil {
				return nil, err
			}
			attributes[key] = converted
			attributeTypes[key] = converted.Type(context.Background())
		}
		return types.ObjectValueMust(attributeTypes, attributes), nil
	default:
		return nil, fmt.Errorf("unsupported API tag value %T", value)
	}
}

func packTagsForRead(ctx context.Context, apiTags, reference types.Dynamic) types.Dynamic {
	actual, err := packTagAttributes(apiTags)
	if err != nil {
		return apiTags
	}
	expected, err := packTagAttributes(reference)
	if err != nil || expected == nil {
		return apiTags
	}
	added := false
	for key, value := range actual {
		if _, exists := expected[key]; !exists && !builtInPackTagKey(key) {
			expected[key] = value
			added = true
		}
	}
	if added {
		attributeTypes := make(map[string]attr.Type, len(expected))
		for key, value := range expected {
			attributeTypes[key] = value.Type(ctx)
		}
		reference = types.DynamicValue(types.ObjectValueMust(attributeTypes, expected))
	}
	return preservePackTags(ctx, apiTags, reference)
}

func preservePackTags(ctx context.Context, apiTags, reference types.Dynamic) types.Dynamic {
	if reference.IsNull() || reference.IsUnknown() || reference.IsUnderlyingValueUnknown() || apiTags.IsNull() {
		return apiTags
	}
	actual, err := packTagAttributes(apiTags)
	if err != nil {
		return apiTags
	}
	expected, err := packTagAttributes(reference)
	if err != nil {
		return apiTags
	}
	values := make(map[string]attr.Value, len(expected))
	for key, planned := range expected {
		if planned.IsNull() && builtInPackTagKey(key) {
			values[key] = planned
		} else if value, exists := actual[key]; exists {
			values[key] = value
		} else if builtInPackTagKey(key) {
			values[key] = types.ListValueMust(types.StringType, nil)
		} else {
			value, err := planned.Type(ctx).ValueFromTerraform(ctx, tftypes.NewValue(planned.Type(ctx).TerraformType(ctx), nil))
			if err != nil {
				return apiTags
			}
			values[key] = value
		}
	}
	converted, err := packTagJSONObject(ctx, values)
	if err != nil {
		return apiTags
	}
	raw, err := json.Marshal(converted)
	if err != nil {
		return apiTags
	}
	typ := reference.UnderlyingValue().Type(ctx)
	terraformValue, err := tftypes.ValueFromJSON(raw, typ.TerraformType(ctx))
	if err != nil {
		return apiTags
	}
	underlying, err := typ.ValueFromTerraform(ctx, terraformValue)
	if err != nil {
		return apiTags
	}
	preserved, err := packTagJSONValue(ctx, underlying)
	if err != nil {
		return apiTags
	}
	preservedJSON, err := json.Marshal(preserved)
	if err != nil || !bytes.Equal(raw, preservedJSON) {
		return apiTags
	}
	return types.DynamicValue(underlying)
}

// PackTagsFromAPIJSON converts API tags to Terraform tags without renaming custom keys.
func PackTagsFromAPIJSON(raw []byte) (types.Dynamic, error) {
	values, err := decodePackTagValues(raw)
	if err != nil {
		return types.Dynamic{}, fmt.Errorf("decode pack tags: %w", err)
	}
	if values == nil {
		return types.DynamicNull(), nil
	}
	return packTagsValue(values)
}

func decodePackTagValues(raw []byte) (map[string]any, error) {
	var values map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&values); err != nil {
		return nil, err
	}
	return values, nil
}
