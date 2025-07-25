// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateSchemaRequest struct {
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// New Schema object
	SchemaLibEntry shared.SchemaLibEntry `request:"mediaType=application/json"`
}

func (o *CreateSchemaRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreateSchemaRequest) GetSchemaLibEntry() shared.SchemaLibEntry {
	if o == nil {
		return shared.SchemaLibEntry{}
	}
	return o.SchemaLibEntry
}

// CreateSchemaResponseBody - a list of Schema objects
type CreateSchemaResponseBody struct {
	Items []shared.SchemaLibEntry `json:"items,omitempty"`
}

func (o *CreateSchemaResponseBody) GetItems() []shared.SchemaLibEntry {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateSchemaResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Schema objects
	Object *CreateSchemaResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateSchemaResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateSchemaResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateSchemaResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateSchemaResponse) GetObject() *CreateSchemaResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateSchemaResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
