// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type UpdateAppscopeLibEntryByIDRequest struct {
	// Unique ID to PATCH
	ID string `pathParam:"style=simple,explode=false,name=id"`
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// AppscopeLibEntry object to be updated
	AppscopeLibEntry shared.AppscopeLibEntry `request:"mediaType=application/json"`
}

func (o *UpdateAppscopeLibEntryByIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdateAppscopeLibEntryByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *UpdateAppscopeLibEntryByIDRequest) GetAppscopeLibEntry() shared.AppscopeLibEntry {
	if o == nil {
		return shared.AppscopeLibEntry{}
	}
	return o.AppscopeLibEntry
}

// UpdateAppscopeLibEntryByIDResponseBody - a list of AppscopeLibEntry objects
type UpdateAppscopeLibEntryByIDResponseBody struct {
	Items []shared.AppscopeLibEntry `json:"items,omitempty"`
}

func (o *UpdateAppscopeLibEntryByIDResponseBody) GetItems() []shared.AppscopeLibEntry {
	if o == nil {
		return nil
	}
	return o.Items
}

type UpdateAppscopeLibEntryByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of AppscopeLibEntry objects
	Object *UpdateAppscopeLibEntryByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *UpdateAppscopeLibEntryByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateAppscopeLibEntryByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateAppscopeLibEntryByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateAppscopeLibEntryByIDResponse) GetObject() *UpdateAppscopeLibEntryByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *UpdateAppscopeLibEntryByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
