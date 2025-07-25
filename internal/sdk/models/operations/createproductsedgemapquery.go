// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// CreateProductsEdgeMapQueryResponseBody - a list of EdgeMapQueryResult objects
type CreateProductsEdgeMapQueryResponseBody struct {
	Items []shared.EdgeMapQueryResult `json:"items,omitempty"`
}

func (o *CreateProductsEdgeMapQueryResponseBody) GetItems() []shared.EdgeMapQueryResult {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateProductsEdgeMapQueryResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of EdgeMapQueryResult objects
	Object *CreateProductsEdgeMapQueryResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateProductsEdgeMapQueryResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateProductsEdgeMapQueryResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateProductsEdgeMapQueryResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateProductsEdgeMapQueryResponse) GetObject() *CreateProductsEdgeMapQueryResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateProductsEdgeMapQueryResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
