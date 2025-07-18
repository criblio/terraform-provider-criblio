// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetPipelineByPackAndIDRequest struct {
	// Unique ID to GET for pack
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// pack ID to GET
	Pack string `pathParam:"style=simple,explode=false,name=pack"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
}

func (o *GetPipelineByPackAndIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *GetPipelineByPackAndIDRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

func (o *GetPipelineByPackAndIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

// GetPipelineByPackAndIDResponseBody - a list of Pipeline objects
type GetPipelineByPackAndIDResponseBody struct {
	Items []shared.Pipeline `json:"items,omitempty"`
}

func (o *GetPipelineByPackAndIDResponseBody) GetItems() []shared.Pipeline {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetPipelineByPackAndIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Pipeline objects
	Object *GetPipelineByPackAndIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetPipelineByPackAndIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetPipelineByPackAndIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetPipelineByPackAndIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetPipelineByPackAndIDResponse) GetObject() *GetPipelineByPackAndIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetPipelineByPackAndIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
