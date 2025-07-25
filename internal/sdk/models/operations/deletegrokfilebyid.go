// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type DeleteGrokFileByIDRequest struct {
	// Unique ID to DELETE
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
}

func (o *DeleteGrokFileByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *DeleteGrokFileByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

// DeleteGrokFileByIDResponseBody - a list of GrokFile objects
type DeleteGrokFileByIDResponseBody struct {
	Items []shared.GrokFile `json:"items,omitempty"`
}

func (o *DeleteGrokFileByIDResponseBody) GetItems() []shared.GrokFile {
	if o == nil {
		return nil
	}
	return o.Items
}

type DeleteGrokFileByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of GrokFile objects
	Object *DeleteGrokFileByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *DeleteGrokFileByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteGrokFileByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteGrokFileByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeleteGrokFileByIDResponse) GetObject() *DeleteGrokFileByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *DeleteGrokFileByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
