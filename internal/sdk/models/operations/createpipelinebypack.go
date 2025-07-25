// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreatePipelineByPackRequest struct {
	// pack ID to POST
	Pack string `pathParam:"style=simple,explode=false,name=pack"`
	// group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// New Pipeline object
	Pipeline shared.Pipeline `request:"mediaType=application/json"`
}

func (o *CreatePipelineByPackRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

func (o *CreatePipelineByPackRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreatePipelineByPackRequest) GetPipeline() shared.Pipeline {
	if o == nil {
		return shared.Pipeline{}
	}
	return o.Pipeline
}

// CreatePipelineByPackResponseBody - a list of Routes objects
type CreatePipelineByPackResponseBody struct {
	Items []shared.Routes `json:"items,omitempty"`
}

func (o *CreatePipelineByPackResponseBody) GetItems() []shared.Routes {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreatePipelineByPackResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Routes objects
	Object *CreatePipelineByPackResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreatePipelineByPackResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreatePipelineByPackResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreatePipelineByPackResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreatePipelineByPackResponse) GetObject() *CreatePipelineByPackResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreatePipelineByPackResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
