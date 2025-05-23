// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// CreateGrokFileResponseBody - a list of GrokFile objects
type CreateGrokFileResponseBody struct {
	// number of items present in the items array
	Count *int64            `json:"count,omitempty"`
	Items []shared.GrokFile `json:"items,omitempty"`
}

func (o *CreateGrokFileResponseBody) GetCount() *int64 {
	if o == nil {
		return nil
	}
	return o.Count
}

func (o *CreateGrokFileResponseBody) GetItems() []shared.GrokFile {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateGrokFileResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of GrokFile objects
	Object *CreateGrokFileResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateGrokFileResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateGrokFileResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateGrokFileResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateGrokFileResponse) GetObject() *CreateGrokFileResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateGrokFileResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
