// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateJobsKeepByIDRequest struct {
	// Job Instance id
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// Group ID
	GroupID *string `queryParam:"style=form,explode=true,name=groupId"`
}

func (o *UpdateJobsKeepByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdateJobsKeepByIDRequest) GetGroupID() *string {
	if o == nil {
		return nil
	}
	return o.GroupID
}

// UpdateJobsKeepByIDResponseBody - a list of JobInfo objects
type UpdateJobsKeepByIDResponseBody struct {
	Items []shared.JobInfo `json:"items,omitempty"`
}

func (o *UpdateJobsKeepByIDResponseBody) GetItems() []shared.JobInfo {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateJobsKeepByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of JobInfo objects
	Object *UpdateJobsKeepByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateJobsKeepByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateJobsKeepByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateJobsKeepByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateJobsKeepByIDResponse) GetObject() *UpdateJobsKeepByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateJobsKeepByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
