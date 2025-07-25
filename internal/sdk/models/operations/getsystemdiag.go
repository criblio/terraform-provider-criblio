// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// GetSystemDiagResponseBody - a list of Diag objects
type GetSystemDiagResponseBody struct {
	Items []shared.Diag `json:"items,omitempty"`
}

func (o *GetSystemDiagResponseBody) GetItems() []shared.Diag {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetSystemDiagResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Diag objects
	Object *GetSystemDiagResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetSystemDiagResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetSystemDiagResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetSystemDiagResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetSystemDiagResponse) GetObject() *GetSystemDiagResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetSystemDiagResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
