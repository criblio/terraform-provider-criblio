// ElementConfigType is hand-maintained (see openapi ElementConfigType additionalProperties-only).

package shared

import (
	"encoding/json"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

// ElementConfigType - Chart/visualization-specific config. The Search API returns arbitrary
// JSON (strings or nested objects for axes, thresholds, etc.); only additionalProperties is used for wire JSON.
type ElementConfigType struct {
	AdditionalProperties any `additionalProperties:"true" json:"-"`
}

func (e ElementConfigType) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(e, "", false)
}

func (e *ElementConfigType) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &e, "", false, nil); err != nil {
		return err
	}
	return nil
}

func elementConfigPropsMap(ap any) map[string]any {
	if ap == nil {
		return nil
	}
	if m, ok := ap.(map[string]any); ok {
		return m
	}
	b, err := json.Marshal(ap)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

func stringProp(m map[string]any, key string) *string {
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

func (e *ElementConfigType) GetXAxis() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "xAxis")
}

func (e *ElementConfigType) GetYAxis() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "yAxis")
}

func (e *ElementConfigType) GetColumns() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "columns")
}

func (e *ElementConfigType) GetMaxRows() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "maxRows")
}

func (e *ElementConfigType) GetSeries() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "series")
}

func (e *ElementConfigType) GetGroupBy() *string {
	if e == nil {
		return nil
	}
	return stringProp(elementConfigPropsMap(e.AdditionalProperties), "groupBy")
}

func (e *ElementConfigType) GetAdditionalProperties() any {
	if e == nil {
		return nil
	}
	return e.AdditionalProperties
}
