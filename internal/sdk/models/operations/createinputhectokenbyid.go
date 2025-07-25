// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateInputHecTokenByIDRequest struct {
	// hec input id
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// AddHecTokenRequest object
	AddHecTokenRequest shared.AddHecTokenRequest `request:"mediaType=application/json"`
}

func (o *CreateInputHecTokenByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *CreateInputHecTokenByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreateInputHecTokenByIDRequest) GetAddHecTokenRequest() shared.AddHecTokenRequest {
	if o == nil {
		return shared.AddHecTokenRequest{}
	}
	return o.AddHecTokenRequest
}

// CreateInputHecTokenByIDResponseBody - a list of any objects
type CreateInputHecTokenByIDResponseBody struct {
	Items []map[string]any `json:"items,omitempty"`
}

func (o *CreateInputHecTokenByIDResponseBody) GetItems() []map[string]any {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateInputHecTokenByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of any objects
	Object *CreateInputHecTokenByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateInputHecTokenByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateInputHecTokenByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateInputHecTokenByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateInputHecTokenByIDResponse) GetObject() *CreateInputHecTokenByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateInputHecTokenByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
