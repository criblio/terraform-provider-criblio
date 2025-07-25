// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// UpdateSystemSettingsAuthResponseBody - a list of AuthConfig objects
type UpdateSystemSettingsAuthResponseBody struct {
	Items []shared.AuthConfig `json:"items,omitempty"`
}

func (o *UpdateSystemSettingsAuthResponseBody) GetItems() []shared.AuthConfig {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateSystemSettingsAuthResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of AuthConfig objects
	Object *UpdateSystemSettingsAuthResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateSystemSettingsAuthResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateSystemSettingsAuthResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateSystemSettingsAuthResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateSystemSettingsAuthResponse) GetObject() *UpdateSystemSettingsAuthResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateSystemSettingsAuthResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
