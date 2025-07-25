// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type WarmPoolSizeEnum string

const (
	WarmPoolSizeEnumAuto WarmPoolSizeEnum = "auto"
)

func (e WarmPoolSizeEnum) ToPointer() *WarmPoolSizeEnum {
	return &e
}
func (e *WarmPoolSizeEnum) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "auto":
		*e = WarmPoolSizeEnum(v)
		return nil
	default:
		return fmt.Errorf("invalid value for WarmPoolSizeEnum: %v", v)
	}
}

type WarmPoolSizeType string

const (
	WarmPoolSizeTypeNumber           WarmPoolSizeType = "number"
	WarmPoolSizeTypeWarmPoolSizeEnum WarmPoolSizeType = "warmPoolSize_enum"
)

type WarmPoolSize struct {
	Number           *float64          `queryParam:"inline"`
	WarmPoolSizeEnum *WarmPoolSizeEnum `queryParam:"inline"`

	Type WarmPoolSizeType
}

func CreateWarmPoolSizeNumber(number float64) WarmPoolSize {
	typ := WarmPoolSizeTypeNumber

	return WarmPoolSize{
		Number: &number,
		Type:   typ,
	}
}

func CreateWarmPoolSizeWarmPoolSizeEnum(warmPoolSizeEnum WarmPoolSizeEnum) WarmPoolSize {
	typ := WarmPoolSizeTypeWarmPoolSizeEnum

	return WarmPoolSize{
		WarmPoolSizeEnum: &warmPoolSizeEnum,
		Type:             typ,
	}
}

func (u *WarmPoolSize) UnmarshalJSON(data []byte) error {

	var number float64 = float64(0)
	if err := utils.UnmarshalJSON(data, &number, "", true, true); err == nil {
		u.Number = &number
		u.Type = WarmPoolSizeTypeNumber
		return nil
	}

	var warmPoolSizeEnum WarmPoolSizeEnum = WarmPoolSizeEnum("")
	if err := utils.UnmarshalJSON(data, &warmPoolSizeEnum, "", true, true); err == nil {
		u.WarmPoolSizeEnum = &warmPoolSizeEnum
		u.Type = WarmPoolSizeTypeWarmPoolSizeEnum
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for WarmPoolSize", string(data))
}

func (u WarmPoolSize) MarshalJSON() ([]byte, error) {
	if u.Number != nil {
		return utils.MarshalJSON(u.Number, "", true)
	}

	if u.WarmPoolSizeEnum != nil {
		return utils.MarshalJSON(u.WarmPoolSizeEnum, "", true)
	}

	return nil, errors.New("could not marshal union type WarmPoolSize: all fields are null")
}

type SearchSettings struct {
	CompressObjectCacheArtifacts bool         `json:"compressObjectCacheArtifacts"`
	FieldSummaryMaxFields        float64      `json:"fieldSummaryMaxFields"`
	FieldSummaryMaxNestedDepth   float64      `json:"fieldSummaryMaxNestedDepth"`
	MaxConcurrentSearches        float64      `json:"maxConcurrentSearches"`
	MaxExecutorsPerSearch        float64      `json:"maxExecutorsPerSearch"`
	MaxResultsPerSearch          float64      `json:"maxResultsPerSearch"`
	SearchHistoryMaxJobs         float64      `json:"searchHistoryMaxJobs"`
	SearchQueueLength            float64      `json:"searchQueueLength"`
	WarmPoolSize                 WarmPoolSize `json:"warmPoolSize"`
	WriteOnlyProviderSecrets     bool         `json:"writeOnlyProviderSecrets"`
}

func (o *SearchSettings) GetCompressObjectCacheArtifacts() bool {
	if o == nil {
		return false
	}
	return o.CompressObjectCacheArtifacts
}

func (o *SearchSettings) GetFieldSummaryMaxFields() float64 {
	if o == nil {
		return 0.0
	}
	return o.FieldSummaryMaxFields
}

func (o *SearchSettings) GetFieldSummaryMaxNestedDepth() float64 {
	if o == nil {
		return 0.0
	}
	return o.FieldSummaryMaxNestedDepth
}

func (o *SearchSettings) GetMaxConcurrentSearches() float64 {
	if o == nil {
		return 0.0
	}
	return o.MaxConcurrentSearches
}

func (o *SearchSettings) GetMaxExecutorsPerSearch() float64 {
	if o == nil {
		return 0.0
	}
	return o.MaxExecutorsPerSearch
}

func (o *SearchSettings) GetMaxResultsPerSearch() float64 {
	if o == nil {
		return 0.0
	}
	return o.MaxResultsPerSearch
}

func (o *SearchSettings) GetSearchHistoryMaxJobs() float64 {
	if o == nil {
		return 0.0
	}
	return o.SearchHistoryMaxJobs
}

func (o *SearchSettings) GetSearchQueueLength() float64 {
	if o == nil {
		return 0.0
	}
	return o.SearchQueueLength
}

func (o *SearchSettings) GetWarmPoolSize() WarmPoolSize {
	if o == nil {
		return WarmPoolSize{}
	}
	return o.WarmPoolSize
}

func (o *SearchSettings) GetWriteOnlyProviderSecrets() bool {
	if o == nil {
		return false
	}
	return o.WriteOnlyProviderSecrets
}
