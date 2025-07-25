// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// GetEdgeProcessesResponseBody - a list of Process objects
type GetEdgeProcessesResponseBody struct {
	Items []shared.Process `json:"items,omitempty"`
}

func (o *GetEdgeProcessesResponseBody) GetItems() []shared.Process {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetEdgeProcessesResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Process objects
	Object *GetEdgeProcessesResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetEdgeProcessesResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetEdgeProcessesResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetEdgeProcessesResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetEdgeProcessesResponse) GetObject() *GetEdgeProcessesResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetEdgeProcessesResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
