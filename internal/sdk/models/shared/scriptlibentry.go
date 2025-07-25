// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type ScriptLibEntry struct {
	ID string `json:"id"`
	// Command to execute for this script
	Command     string  `json:"command"`
	Description *string `json:"description,omitempty"`
	// Arguments to pass when executing this script
	Args []string `json:"args,omitempty"`
	// Extra environment variables to set when executing script
	Env                  map[string]string `json:"env,omitempty"`
	AdditionalProperties any               `additionalProperties:"true" json:"-"`
}

func (s ScriptLibEntry) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(s, "", false)
}

func (s *ScriptLibEntry) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &s, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ScriptLibEntry) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ScriptLibEntry) GetCommand() string {
	if o == nil {
		return ""
	}
	return o.Command
}

func (o *ScriptLibEntry) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *ScriptLibEntry) GetArgs() []string {
	if o == nil {
		return nil
	}
	return o.Args
}

func (o *ScriptLibEntry) GetEnv() map[string]string {
	if o == nil {
		return nil
	}
	return o.Env
}

func (o *ScriptLibEntry) GetAdditionalProperties() any {
	if o == nil {
		return nil
	}
	return o.AdditionalProperties
}
