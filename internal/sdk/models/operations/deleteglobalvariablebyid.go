// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type DeleteGlobalVariableByIDRequest struct {
	// Unique ID to DELETE
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
}

func (o *DeleteGlobalVariableByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *DeleteGlobalVariableByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

// DeleteGlobalVariableByIDResponseBody - a list of Global Variable objects
type DeleteGlobalVariableByIDResponseBody struct {
	Items []shared.GlobalVar `json:"items,omitempty"`
}

func (o *DeleteGlobalVariableByIDResponseBody) GetItems() []shared.GlobalVar {
	if o == nil {
		return nil
	}
	return o.Items
}

type DeleteGlobalVariableByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Global Variable objects
	Object *DeleteGlobalVariableByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *DeleteGlobalVariableByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteGlobalVariableByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteGlobalVariableByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeleteGlobalVariableByIDResponse) GetObject() *DeleteGlobalVariableByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *DeleteGlobalVariableByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
