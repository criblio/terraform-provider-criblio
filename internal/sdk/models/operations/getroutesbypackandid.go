// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetRoutesByPackAndIDRequest struct {
	// Unique ID to GET for pack
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// pack ID to GET
	Pack string `pathParam:"style=simple,explode=false,name=pack"`
}

func (o *GetRoutesByPackAndIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *GetRoutesByPackAndIDRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

// GetRoutesByPackAndIDResponseBody - a list of Routes objects
type GetRoutesByPackAndIDResponseBody struct {
	Items []shared.Routes `json:"items,omitempty"`
}

func (o *GetRoutesByPackAndIDResponseBody) GetItems() []shared.Routes {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetRoutesByPackAndIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Routes objects
	Object *GetRoutesByPackAndIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetRoutesByPackAndIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetRoutesByPackAndIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetRoutesByPackAndIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetRoutesByPackAndIDResponse) GetObject() *GetRoutesByPackAndIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetRoutesByPackAndIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
