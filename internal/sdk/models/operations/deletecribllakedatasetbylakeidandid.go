// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type DeleteCriblLakeDatasetByLakeIDAndIDRequest struct {
	// lake id that contains the Datasets
	LakeID string `pathParam:"style=simple,explode=false,name=lakeId"`
	// dataset id to delete
	ID string `pathParam:"style=simple,explode=false,name=id"`
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDRequest) GetLakeID() string {
	if o == nil {
		return ""
	}
	return o.LakeID
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

// DeleteCriblLakeDatasetByLakeIDAndIDResponseBody - a list of CriblLakeDataset objects
type DeleteCriblLakeDatasetByLakeIDAndIDResponseBody struct {
	Items []shared.CriblLakeDataset `json:"items,omitempty"`
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponseBody) GetItems() []shared.CriblLakeDataset {
	if o == nil {
		return nil
	}
	return o.Items
}

type DeleteCriblLakeDatasetByLakeIDAndIDResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of CriblLakeDataset objects
	Object *DeleteCriblLakeDatasetByLakeIDAndIDResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponse) GetObject() *DeleteCriblLakeDatasetByLakeIDAndIDResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *DeleteCriblLakeDatasetByLakeIDAndIDResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
