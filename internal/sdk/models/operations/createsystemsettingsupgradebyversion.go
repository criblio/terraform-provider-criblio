// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateSystemSettingsUpgradeByVersionRequest struct {
	// Version number
	Version string `pathParam:"style=simple,explode=false,name=version"`
}

func (o *CreateSystemSettingsUpgradeByVersionRequest) GetVersion() string {
	if o == nil {
		return ""
	}
	return o.Version
}

// CreateSystemSettingsUpgradeByVersionResponseBody - a list of string objects
type CreateSystemSettingsUpgradeByVersionResponseBody struct {
	Items []string `json:"items,omitempty"`
}

func (o *CreateSystemSettingsUpgradeByVersionResponseBody) GetItems() []string {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateSystemSettingsUpgradeByVersionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of string objects
	Object *CreateSystemSettingsUpgradeByVersionResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateSystemSettingsUpgradeByVersionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateSystemSettingsUpgradeByVersionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateSystemSettingsUpgradeByVersionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateSystemSettingsUpgradeByVersionResponse) GetObject() *CreateSystemSettingsUpgradeByVersionResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateSystemSettingsUpgradeByVersionResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
