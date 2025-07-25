// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateSystemInputsByPackRequest struct {
	// pack ID to PATCH
	Pack     string `pathParam:"style=simple,explode=false,name=pack"`
	Disabled *bool  `queryParam:"style=form,explode=true,name=disabled"`
	// group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// Unique ID to PATCH for pack source
	ID    string       `pathParam:"style=simple,explode=false,name=id"`
	Input shared.Input `request:"mediaType=application/json"`
}

func (o *UpdateSystemInputsByPackRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

func (o *UpdateSystemInputsByPackRequest) GetDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.Disabled
}

func (o *UpdateSystemInputsByPackRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *UpdateSystemInputsByPackRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdateSystemInputsByPackRequest) GetInput() shared.Input {
	if o == nil {
		return shared.Input{}
	}
	return o.Input
}

// UpdateSystemInputsByPackResponseBody - a list of Pipeline objects
type UpdateSystemInputsByPackResponseBody struct {
	Items []shared.Pipeline `json:"items,omitempty"`
}

func (o *UpdateSystemInputsByPackResponseBody) GetItems() []shared.Pipeline {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateSystemInputsByPackResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Pipeline objects
	Object *UpdateSystemInputsByPackResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateSystemInputsByPackResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateSystemInputsByPackResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateSystemInputsByPackResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateSystemInputsByPackResponse) GetObject() *UpdateSystemInputsByPackResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateSystemInputsByPackResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
