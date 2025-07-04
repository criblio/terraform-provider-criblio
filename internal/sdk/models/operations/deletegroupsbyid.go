// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type DeleteGroupsByIDRequest struct {
	// Group id
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

func (o *DeleteGroupsByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// DeleteGroupsByIDResponseBody - a list of ConfigGroup objects
type DeleteGroupsByIDResponseBody struct {
	Items []shared.Group `json:"items,omitempty"`
}

func (o *DeleteGroupsByIDResponseBody) GetItems() []shared.Group {
	if o == nil {
		return nil
	}
	return o.Items
}

type DeleteGroupsByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of ConfigGroup objects
	Object *DeleteGroupsByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *DeleteGroupsByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteGroupsByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteGroupsByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeleteGroupsByIDResponse) GetObject() *DeleteGroupsByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *DeleteGroupsByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
