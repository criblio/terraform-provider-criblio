// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type GetPacksExportByIDRequest struct {
	// Group Id
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
}

func (o *GetPacksExportByIDRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

// GetPacksExportByIDResponseBody - a list of any objects
type GetPacksExportByIDResponseBody struct {
	// number of items present in the items array
	Count *int64           `json:"count,omitempty"`
	Items []map[string]any `json:"items,omitempty"`
}

func (o *GetPacksExportByIDResponseBody) GetCount() *int64 {
	if o == nil {
		return nil
	}
	return o.Count
}

func (o *GetPacksExportByIDResponseBody) GetItems() []map[string]any {
	if o == nil {
		return nil
	}
	return o.Items
}

type GetPacksExportByIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of any objects
	Object *GetPacksExportByIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *GetPacksExportByIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetPacksExportByIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetPacksExportByIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetPacksExportByIDResponse) GetObject() *GetPacksExportByIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *GetPacksExportByIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
