// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateRoleByIDRequest struct {
	// Unique ID to PATCH
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// Role object to be updated
	Role shared.Role `request:"mediaType=application/json"`
}

func (o *UpdateRoleByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdateRoleByIDRequest) GetRole() shared.Role {
	if o == nil {
		return shared.Role{}
	}
	return o.Role
}

// UpdateRoleByIDResponseBody - a list of Role objects
type UpdateRoleByIDResponseBody struct {
	Items []shared.Role `json:"items,omitempty"`
}

func (o *UpdateRoleByIDResponseBody) GetItems() []shared.Role {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateRoleByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Role objects
	Object *UpdateRoleByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateRoleByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateRoleByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateRoleByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateRoleByIDResponse) GetObject() *UpdateRoleByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateRoleByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
