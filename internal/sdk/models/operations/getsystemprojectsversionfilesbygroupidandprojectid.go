// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetSystemProjectsVersionFilesByGroupIDAndProjectIDRequest struct {
	// Group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// Project Id
	ProjectID string `pathParam:"style=simple,explode=false,name=projectId"`
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDRequest) GetProjectID() string {
	if o == nil {
		return ""
	}
	return o.ProjectID
}

// GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponseBody - A list of GitFilesResponse objects
type GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponseBody struct {
	Items []shared.GitFilesResponse `json:"items,omitempty"`
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponseBody) GetItems() []shared.GitFilesResponse {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// A list of GitFilesResponse objects
	Object *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse) GetObject() *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetSystemProjectsVersionFilesByGroupIDAndProjectIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
