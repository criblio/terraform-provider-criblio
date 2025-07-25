// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDRequest struct {
	// Group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// Project Id
	ProjectID string `pathParam:"style=simple,explode=false,name=projectId"`
	// Pipeline Id
	PipelineID string `pathParam:"style=simple,explode=false,name=pipelineId"`
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDRequest) GetProjectID() string {
	if o == nil {
		return ""
	}
	return o.ProjectID
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDRequest) GetPipelineID() string {
	if o == nil {
		return ""
	}
	return o.PipelineID
}

// GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponseBody - A list of Pipeline objects in specified Project
type GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponseBody struct {
	Items []shared.Pipeline `json:"items,omitempty"`
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponseBody) GetItems() []shared.Pipeline {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// A list of Pipeline objects in specified Project
	Object *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse) GetObject() *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetSystemProjectsPipelinesByGroupIDAndProjectIDAndPipelineIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
