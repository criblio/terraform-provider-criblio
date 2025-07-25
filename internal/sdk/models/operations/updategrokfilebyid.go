// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateGrokFileByIDRequest struct {
	// Unique ID to PATCH
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// GrokFile object to be updated
	GrokFile shared.GrokFile `request:"mediaType=application/json"`
}

func (o *UpdateGrokFileByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdateGrokFileByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *UpdateGrokFileByIDRequest) GetGrokFile() shared.GrokFile {
	if o == nil {
		return shared.GrokFile{}
	}
	return o.GrokFile
}

// UpdateGrokFileByIDResponseBody - a list of GrokFile objects
type UpdateGrokFileByIDResponseBody struct {
	Items []shared.GrokFile `json:"items,omitempty"`
}

func (o *UpdateGrokFileByIDResponseBody) GetItems() []shared.GrokFile {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateGrokFileByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of GrokFile objects
	Object *UpdateGrokFileByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateGrokFileByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateGrokFileByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateGrokFileByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateGrokFileByIDResponse) GetObject() *UpdateGrokFileByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateGrokFileByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
