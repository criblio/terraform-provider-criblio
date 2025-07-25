// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type TypeValues string

const (
	TypeValuesValues TypeValues = "values"
)

func (e TypeValues) ToPointer() *TypeValues {
	return &e
}
func (e *TypeValues) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "values":
		*e = TypeValues(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeValues: %v", v)
	}
}

type SearchQueryValues struct {
	Type   TypeValues `json:"type"`
	Values []string   `json:"values"`
}

func (o *SearchQueryValues) GetType() TypeValues {
	if o == nil {
		return TypeValues("")
	}
	return o.Type
}

func (o *SearchQueryValues) GetValues() []string {
	if o == nil {
		return []string{}
	}
	return o.Values
}

type SearchQueryEarliestType string

const (
	SearchQueryEarliestTypeStr    SearchQueryEarliestType = "str"
	SearchQueryEarliestTypeNumber SearchQueryEarliestType = "number"
)

type SearchQueryEarliest struct {
	Str    *string  `queryParam:"inline"`
	Number *float64 `queryParam:"inline"`

	Type SearchQueryEarliestType
}

func CreateSearchQueryEarliestStr(str string) SearchQueryEarliest {
	typ := SearchQueryEarliestTypeStr

	return SearchQueryEarliest{
		Str:  &str,
		Type: typ,
	}
}

func CreateSearchQueryEarliestNumber(number float64) SearchQueryEarliest {
	typ := SearchQueryEarliestTypeNumber

	return SearchQueryEarliest{
		Number: &number,
		Type:   typ,
	}
}

func (u *SearchQueryEarliest) UnmarshalJSON(data []byte) error {

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = SearchQueryEarliestTypeStr
		return nil
	}

	var number float64 = float64(0)
	if err := utils.UnmarshalJSON(data, &number, "", true, true); err == nil {
		u.Number = &number
		u.Type = SearchQueryEarliestTypeNumber
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for SearchQueryEarliest", string(data))
}

func (u SearchQueryEarliest) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.Number != nil {
		return utils.MarshalJSON(u.Number, "", true)
	}

	return nil, errors.New("could not marshal union type SearchQueryEarliest: all fields are null")
}

type SearchQueryLatestType string

const (
	SearchQueryLatestTypeStr    SearchQueryLatestType = "str"
	SearchQueryLatestTypeNumber SearchQueryLatestType = "number"
)

type SearchQueryLatest struct {
	Str    *string  `queryParam:"inline"`
	Number *float64 `queryParam:"inline"`

	Type SearchQueryLatestType
}

func CreateSearchQueryLatestStr(str string) SearchQueryLatest {
	typ := SearchQueryLatestTypeStr

	return SearchQueryLatest{
		Str:  &str,
		Type: typ,
	}
}

func CreateSearchQueryLatestNumber(number float64) SearchQueryLatest {
	typ := SearchQueryLatestTypeNumber

	return SearchQueryLatest{
		Number: &number,
		Type:   typ,
	}
}

func (u *SearchQueryLatest) UnmarshalJSON(data []byte) error {

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = SearchQueryLatestTypeStr
		return nil
	}

	var number float64 = float64(0)
	if err := utils.UnmarshalJSON(data, &number, "", true, true); err == nil {
		u.Number = &number
		u.Type = SearchQueryLatestTypeNumber
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for SearchQueryLatest", string(data))
}

func (u SearchQueryLatest) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.Number != nil {
		return utils.MarshalJSON(u.Number, "", true)
	}

	return nil, errors.New("could not marshal union type SearchQueryLatest: all fields are null")
}

type TypeInline string

const (
	TypeInlineInline TypeInline = "inline"
)

func (e TypeInline) ToPointer() *TypeInline {
	return &e
}
func (e *TypeInline) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "inline":
		*e = TypeInline(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeInline: %v", v)
	}
}

type SearchQueryInline struct {
	Earliest       *SearchQueryEarliest `json:"earliest"`
	Latest         *SearchQueryLatest   `json:"latest"`
	ParentSearchID *string              `json:"parentSearchId,omitempty"`
	Query          *string              `json:"query"`
	SampleRate     *float64             `json:"sampleRate,omitempty"`
	Timezone       *string              `json:"timezone,omitempty"`
	Type           TypeInline           `json:"type"`
}

func (o *SearchQueryInline) GetEarliest() *SearchQueryEarliest {
	if o == nil {
		return nil
	}
	return o.Earliest
}

func (o *SearchQueryInline) GetLatest() *SearchQueryLatest {
	if o == nil {
		return nil
	}
	return o.Latest
}

func (o *SearchQueryInline) GetParentSearchID() *string {
	if o == nil {
		return nil
	}
	return o.ParentSearchID
}

func (o *SearchQueryInline) GetQuery() *string {
	if o == nil {
		return nil
	}
	return o.Query
}

func (o *SearchQueryInline) GetSampleRate() *float64 {
	if o == nil {
		return nil
	}
	return o.SampleRate
}

func (o *SearchQueryInline) GetTimezone() *string {
	if o == nil {
		return nil
	}
	return o.Timezone
}

func (o *SearchQueryInline) GetType() TypeInline {
	if o == nil {
		return TypeInline("")
	}
	return o.Type
}

type TypeSaved string

const (
	TypeSavedSaved TypeSaved = "saved"
)

func (e TypeSaved) ToPointer() *TypeSaved {
	return &e
}
func (e *TypeSaved) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "saved":
		*e = TypeSaved(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TypeSaved: %v", v)
	}
}

type SearchQuerySaved struct {
	Query   *string             `json:"query,omitempty"`
	QueryID string              `json:"queryId"`
	RunMode *SavesSearchRunMode `json:"runMode,omitempty"`
	Type    TypeSaved           `json:"type"`
}

func (o *SearchQuerySaved) GetQuery() *string {
	if o == nil {
		return nil
	}
	return o.Query
}

func (o *SearchQuerySaved) GetQueryID() string {
	if o == nil {
		return ""
	}
	return o.QueryID
}

func (o *SearchQuerySaved) GetRunMode() *SavesSearchRunMode {
	if o == nil {
		return nil
	}
	return o.RunMode
}

func (o *SearchQuerySaved) GetType() TypeSaved {
	if o == nil {
		return TypeSaved("")
	}
	return o.Type
}

type SearchQueryType string

const (
	SearchQueryTypeSearchQuerySaved  SearchQueryType = "SearchQuery_Saved"
	SearchQueryTypeSearchQueryInline SearchQueryType = "SearchQuery_Inline"
	SearchQueryTypeSearchQueryValues SearchQueryType = "SearchQuery_Values"
)

type SearchQuery struct {
	SearchQuerySaved  *SearchQuerySaved  `queryParam:"inline"`
	SearchQueryInline *SearchQueryInline `queryParam:"inline"`
	SearchQueryValues *SearchQueryValues `queryParam:"inline"`

	Type SearchQueryType
}

func CreateSearchQuerySearchQuerySaved(searchQuerySaved SearchQuerySaved) SearchQuery {
	typ := SearchQueryTypeSearchQuerySaved

	return SearchQuery{
		SearchQuerySaved: &searchQuerySaved,
		Type:             typ,
	}
}

func CreateSearchQuerySearchQueryInline(searchQueryInline SearchQueryInline) SearchQuery {
	typ := SearchQueryTypeSearchQueryInline

	return SearchQuery{
		SearchQueryInline: &searchQueryInline,
		Type:              typ,
	}
}

func CreateSearchQuerySearchQueryValues(searchQueryValues SearchQueryValues) SearchQuery {
	typ := SearchQueryTypeSearchQueryValues

	return SearchQuery{
		SearchQueryValues: &searchQueryValues,
		Type:              typ,
	}
}

func (u *SearchQuery) UnmarshalJSON(data []byte) error {

	var searchQueryValues SearchQueryValues = SearchQueryValues{}
	if err := utils.UnmarshalJSON(data, &searchQueryValues, "", true, true); err == nil {
		u.SearchQueryValues = &searchQueryValues
		u.Type = SearchQueryTypeSearchQueryValues
		return nil
	}

	var searchQuerySaved SearchQuerySaved = SearchQuerySaved{}
	if err := utils.UnmarshalJSON(data, &searchQuerySaved, "", true, true); err == nil {
		u.SearchQuerySaved = &searchQuerySaved
		u.Type = SearchQueryTypeSearchQuerySaved
		return nil
	}

	var searchQueryInline SearchQueryInline = SearchQueryInline{}
	if err := utils.UnmarshalJSON(data, &searchQueryInline, "", true, true); err == nil {
		u.SearchQueryInline = &searchQueryInline
		u.Type = SearchQueryTypeSearchQueryInline
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for SearchQuery", string(data))
}

func (u SearchQuery) MarshalJSON() ([]byte, error) {
	if u.SearchQuerySaved != nil {
		return utils.MarshalJSON(u.SearchQuerySaved, "", true)
	}

	if u.SearchQueryInline != nil {
		return utils.MarshalJSON(u.SearchQueryInline, "", true)
	}

	if u.SearchQueryValues != nil {
		return utils.MarshalJSON(u.SearchQueryValues, "", true)
	}

	return nil, errors.New("could not marshal union type SearchQuery: all fields are null")
}
