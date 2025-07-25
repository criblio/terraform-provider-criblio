// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

type CreateSubscriptionRequest struct {
	// The consumer group to which this instance belongs. Defaults to 'Cribl'.
	GroupID string `pathParam:"style=simple,explode=false,name=groupId"`
	// Project Id
	Disabled *bool `queryParam:"style=form,explode=true,name=disabled"`
	// Project description
	Description *string `queryParam:"style=form,explode=true,name=description"`
	// filter
	Filter *string `queryParam:"style=form,explode=true,name=filter"`
	// pipeline to be used
	Pipeline string `queryParam:"style=form,explode=true,name=pipeline"`
	// pipeline to be used
	ID string `queryParam:"style=form,explode=true,name=id"`
	// Subscription object
	Subscription shared.Subscription `request:"mediaType=application/json"`
}

func (o *CreateSubscriptionRequest) GetGroupID() string {
	if o == nil {
		return ""
	}
	return o.GroupID
}

func (o *CreateSubscriptionRequest) GetDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.Disabled
}

func (o *CreateSubscriptionRequest) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *CreateSubscriptionRequest) GetFilter() *string {
	if o == nil {
		return nil
	}
	return o.Filter
}

func (o *CreateSubscriptionRequest) GetPipeline() string {
	if o == nil {
		return ""
	}
	return o.Pipeline
}

func (o *CreateSubscriptionRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *CreateSubscriptionRequest) GetSubscription() shared.Subscription {
	if o == nil {
		return shared.Subscription{}
	}
	return o.Subscription
}

// CreateSubscriptionResponseBody - a list of Subscription objects
type CreateSubscriptionResponseBody struct {
	Items []shared.Subscription `json:"items,omitempty"`
}

func (o *CreateSubscriptionResponseBody) GetItems() []shared.Subscription {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateSubscriptionResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of Subscription objects
	Object *CreateSubscriptionResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateSubscriptionResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateSubscriptionResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateSubscriptionResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateSubscriptionResponse) GetObject() *CreateSubscriptionResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateSubscriptionResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
