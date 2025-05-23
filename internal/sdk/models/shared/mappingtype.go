// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type MappingType string

const (
	MappingTypeAutomatic MappingType = "automatic"
	MappingTypeCustom    MappingType = "custom"
)

func (e MappingType) ToPointer() *MappingType {
	return &e
}
func (e *MappingType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "automatic":
		fallthrough
	case "custom":
		*e = MappingType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for MappingType: %v", v)
	}
}
