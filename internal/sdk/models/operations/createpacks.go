// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreatePacksRequest struct {
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// the file to upload
	Filename *string `queryParam:"style=form,explode=true,name=filename"`
	Disabled *bool   `queryParam:"style=form,explode=true,name=disabled"`
	// PackRequestBody object
	PackRequestBody shared.PackRequestBody `request:"mediaType=application/json"`
}

func (o *CreatePacksRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreatePacksRequest) GetFilename() *string {
	if o == nil {
		return nil
	}
	return o.Filename
}

func (o *CreatePacksRequest) GetDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.Disabled
}

func (o *CreatePacksRequest) GetPackRequestBody() shared.PackRequestBody {
	if o == nil {
		return shared.PackRequestBody{}
	}
	return o.PackRequestBody
}

// CreatePacksResponseBody - a list of PackInstallInfo objects
type CreatePacksResponseBody struct {
	Items []shared.PackInstallInfo `json:"items,omitempty"`
}

func (o *CreatePacksResponseBody) GetItems() []shared.PackInstallInfo {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreatePacksResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of PackInstallInfo objects
	Object *CreatePacksResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreatePacksResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreatePacksResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreatePacksResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreatePacksResponse) GetObject() *CreatePacksResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreatePacksResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
