// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateEdgeFileIngestRequest struct {
	// Absolute path to file to ingest.
	FilePath *string `queryParam:"style=form,explode=true,name=filePath"`
	// Id of the pipeline to use.
	PipelineID *string `queryParam:"style=form,explode=true,name=pipelineId"`
	// Destination to send events to.
	OutputID *string `queryParam:"style=form,explode=true,name=outputId"`
	// Id to the pre-processing pipeline to use for routes.
	PreProcessingPipelineID *string `queryParam:"style=form,explode=true,name=preProcessingPipelineId"`
	// boolean condition required on whether to send events to routes.
	SendToRoutes *string `queryParam:"style=form,explode=true,name=sendToRoutes"`
	// Breaker rules to use on the file.
	BreakerRuleSet *string `queryParam:"style=form,explode=true,name=breakerRuleSet"`
}

func (o *CreateEdgeFileIngestRequest) GetFilePath() *string {
	if o == nil {
		return nil
	}
	return o.FilePath
}

func (o *CreateEdgeFileIngestRequest) GetPipelineID() *string {
	if o == nil {
		return nil
	}
	return o.PipelineID
}

func (o *CreateEdgeFileIngestRequest) GetOutputID() *string {
	if o == nil {
		return nil
	}
	return o.OutputID
}

func (o *CreateEdgeFileIngestRequest) GetPreProcessingPipelineID() *string {
	if o == nil {
		return nil
	}
	return o.PreProcessingPipelineID
}

func (o *CreateEdgeFileIngestRequest) GetSendToRoutes() *string {
	if o == nil {
		return nil
	}
	return o.SendToRoutes
}

func (o *CreateEdgeFileIngestRequest) GetBreakerRuleSet() *string {
	if o == nil {
		return nil
	}
	return o.BreakerRuleSet
}

// CreateEdgeFileIngestResponseBody - a list of any objects
type CreateEdgeFileIngestResponseBody struct {
	Items []map[string]any `json:"items,omitempty"`
}

func (o *CreateEdgeFileIngestResponseBody) GetItems() []map[string]any {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateEdgeFileIngestResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of any objects
	Object *CreateEdgeFileIngestResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateEdgeFileIngestResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateEdgeFileIngestResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateEdgeFileIngestResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateEdgeFileIngestResponse) GetObject() *CreateEdgeFileIngestResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateEdgeFileIngestResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
