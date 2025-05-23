// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// CreateHmacFunctionResponseBody - a list of HmacFunction objects
type CreateHmacFunctionResponseBody struct {
	// number of items present in the items array
	Count *int64                `json:"count,omitempty"`
	Items []shared.HmacFunction `json:"items,omitempty"`
}

func (o *CreateHmacFunctionResponseBody) GetCount() *int64 {
	if o == nil {
		return nil
	}
	return o.Count
}

func (o *CreateHmacFunctionResponseBody) GetItems() []shared.HmacFunction {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateHmacFunctionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of HmacFunction objects
	Object *CreateHmacFunctionResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateHmacFunctionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateHmacFunctionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateHmacFunctionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateHmacFunctionResponse) GetObject() *CreateHmacFunctionResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateHmacFunctionResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
