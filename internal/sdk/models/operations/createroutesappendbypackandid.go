// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateRoutesAppendByPackAndIDRequest struct {
	// the route table to be appended to - currently default is the only supported value for pack
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// pack ID to POST
	Pack string `pathParam:"style=simple,explode=false,name=pack"`
	// RouteDefinitions object
	RequestBody []shared.RouteConf `request:"mediaType=application/json"`
}

func (o *CreateRoutesAppendByPackAndIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *CreateRoutesAppendByPackAndIDRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

func (o *CreateRoutesAppendByPackAndIDRequest) GetRequestBody() []shared.RouteConf {
	if o == nil {
		return []shared.RouteConf{}
	}
	return o.RequestBody
}

// CreateRoutesAppendByPackAndIDResponseBody - a list of any objects
type CreateRoutesAppendByPackAndIDResponseBody struct {
	Items []map[string]any `json:"items,omitempty"`
}

func (o *CreateRoutesAppendByPackAndIDResponseBody) GetItems() []map[string]any {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateRoutesAppendByPackAndIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of any objects
	Object *CreateRoutesAppendByPackAndIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateRoutesAppendByPackAndIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateRoutesAppendByPackAndIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateRoutesAppendByPackAndIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateRoutesAppendByPackAndIDResponse) GetObject() *CreateRoutesAppendByPackAndIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateRoutesAppendByPackAndIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
