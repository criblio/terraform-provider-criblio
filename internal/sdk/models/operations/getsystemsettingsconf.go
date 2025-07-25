// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// GetSystemSettingsConfResponseBody - a list of SystemSettingsConf objects
type GetSystemSettingsConfResponseBody struct {
	Items []shared.SystemSettingsConf `json:"items,omitempty"`
}

func (o *GetSystemSettingsConfResponseBody) GetItems() []shared.SystemSettingsConf {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetSystemSettingsConfResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of SystemSettingsConf objects
	Object *GetSystemSettingsConfResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetSystemSettingsConfResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetSystemSettingsConfResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetSystemSettingsConfResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetSystemSettingsConfResponse) GetObject() *GetSystemSettingsConfResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetSystemSettingsConfResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
