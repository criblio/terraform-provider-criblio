// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateSystemLookupsByPackAndIDRequest struct {
	// Unique ID to PATCH for pack
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// pack ID to PATCH
	Pack     string `pathParam:"style=simple,explode=false,name=pack"`
	Disabled *bool  `queryParam:"style=form,explode=true,name=disabled"`
	// group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// Pipeline object to be updated in specified Project
	LookupFile shared.LookupFileInputUnion `request:"mediaType=application/json"`
}

func (o *CreateSystemLookupsByPackAndIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *CreateSystemLookupsByPackAndIDRequest) GetPack() string {
	if o == nil {
		return ""
	}
	return o.Pack
}

func (o *CreateSystemLookupsByPackAndIDRequest) GetDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.Disabled
}

func (o *CreateSystemLookupsByPackAndIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreateSystemLookupsByPackAndIDRequest) GetLookupFile() shared.LookupFileInputUnion {
	if o == nil {
		return shared.LookupFileInputUnion{}
	}
	return o.LookupFile
}

// CreateSystemLookupsByPackAndIDResponseBody - a list of Pipeline objects
type CreateSystemLookupsByPackAndIDResponseBody struct {
	Items []shared.Pipeline `json:"items,omitempty"`
}

func (o *CreateSystemLookupsByPackAndIDResponseBody) GetItems() []shared.Pipeline {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateSystemLookupsByPackAndIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Pipeline objects
	Object *CreateSystemLookupsByPackAndIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateSystemLookupsByPackAndIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateSystemLookupsByPackAndIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateSystemLookupsByPackAndIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateSystemLookupsByPackAndIDResponse) GetObject() *CreateSystemLookupsByPackAndIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateSystemLookupsByPackAndIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
