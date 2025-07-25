// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetCollectorByIDRequest struct {
	// Unique ID to GET
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

func (o *GetCollectorByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// GetCollectorByIDResponseBody - a list of Collector objects
type GetCollectorByIDResponseBody struct {
	Items []shared.Collector `json:"items,omitempty"`
}

func (o *GetCollectorByIDResponseBody) GetItems() []shared.Collector {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetCollectorByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Collector objects
	Object *GetCollectorByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetCollectorByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetCollectorByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetCollectorByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetCollectorByIDResponse) GetObject() *GetCollectorByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetCollectorByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
