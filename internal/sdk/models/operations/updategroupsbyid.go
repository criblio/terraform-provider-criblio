// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateGroupsByIDRequest struct {
	// Group id
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

func (o *UpdateGroupsByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// UpdateGroupsByIDResponseBody - a list of ConfigGroup objects
type UpdateGroupsByIDResponseBody struct {
	Items []shared.Group `json:"items,omitempty"`
}

func (o *UpdateGroupsByIDResponseBody) GetItems() []shared.Group {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateGroupsByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of ConfigGroup objects
	Object *UpdateGroupsByIDResponseBody
}

func (o *UpdateGroupsByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateGroupsByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateGroupsByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateGroupsByIDResponse) GetObject() *UpdateGroupsByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}
