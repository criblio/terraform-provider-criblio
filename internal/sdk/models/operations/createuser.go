// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// CreateUserResponseBody - a list of User objects
type CreateUserResponseBody struct {
	Items []shared.User `json:"items,omitempty"`
}

func (o *CreateUserResponseBody) GetItems() []shared.User {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateUserResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of User objects
	Object *CreateUserResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateUserResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateUserResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateUserResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateUserResponse) GetObject() *CreateUserResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateUserResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
