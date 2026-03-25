package provider

import (
	"encoding/json"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func tfElementConfigFromShared(cfg *shared.ElementConfigType) *tfTypes.ElementConfigType {
	if cfg == nil {
		return nil
	}
	out := &tfTypes.ElementConfigType{}
	if cfg.AdditionalProperties == nil {
		out.AdditionalProperties = jsontypes.NewNormalizedNull()
		return out
	}
	raw, err := json.Marshal(cfg.AdditionalProperties)
	if err != nil {
		out.AdditionalProperties = jsontypes.NewNormalizedNull()
		return out
	}
	out.AdditionalProperties = jsontypes.NewNormalizedValue(string(raw))

	m := elementConfigAsMap(cfg.AdditionalProperties)
	if m == nil {
		return out
	}
	out.XAxis = types.StringPointerValue(stringFromMap(m, "xAxis"))
	out.YAxis = types.StringPointerValue(stringFromMap(m, "yAxis"))
	out.Columns = types.StringPointerValue(stringFromMap(m, "columns"))
	out.MaxRows = types.StringPointerValue(stringFromMap(m, "maxRows"))
	out.Series = types.StringPointerValue(stringFromMap(m, "series"))
	out.GroupBy = types.StringPointerValue(stringFromMap(m, "groupBy"))
	return out
}

func sharedElementConfigFromTF(tf *tfTypes.ElementConfigType) *shared.ElementConfigType {
	if tf == nil {
		return nil
	}
	m := map[string]any{}
	if !tf.AdditionalProperties.IsUnknown() && !tf.AdditionalProperties.IsNull() {
		_ = json.Unmarshal([]byte(tf.AdditionalProperties.ValueString()), &m)
	}
	setString := func(key string, attr types.String) {
		if attr.IsUnknown() || attr.IsNull() {
			return
		}
		m[key] = attr.ValueString()
	}
	setString("xAxis", tf.XAxis)
	setString("yAxis", tf.YAxis)
	setString("columns", tf.Columns)
	setString("maxRows", tf.MaxRows)
	setString("series", tf.Series)
	setString("groupBy", tf.GroupBy)
	if len(m) == 0 {
		return &shared.ElementConfigType{AdditionalProperties: nil}
	}
	return &shared.ElementConfigType{AdditionalProperties: m}
}

func elementConfigAsMap(v any) map[string]any {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

func stringFromMap(m map[string]any, key string) *string {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return &s
}
