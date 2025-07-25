// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
)

// CreateSystemSettingsUpgradeFromPackageResponseBody - a list of string objects
type CreateSystemSettingsUpgradeFromPackageResponseBody struct {
	Items []string `json:"items,omitempty"`
}

func (o *CreateSystemSettingsUpgradeFromPackageResponseBody) GetItems() []string {
	if o == nil {
		return nil
	}
	return o.Items
}

type CreateSystemSettingsUpgradeFromPackageResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// a list of string objects
	Object *CreateSystemSettingsUpgradeFromPackageResponseBody
	// Unexpected error
	Error *shared.Error
}

func (o *CreateSystemSettingsUpgradeFromPackageResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *CreateSystemSettingsUpgradeFromPackageResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *CreateSystemSettingsUpgradeFromPackageResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *CreateSystemSettingsUpgradeFromPackageResponse) GetObject() *CreateSystemSettingsUpgradeFromPackageResponseBody {
	if o == nil {
		return nil
	}
	return o.Object
}

func (o *CreateSystemSettingsUpgradeFromPackageResponse) GetError() *shared.Error {
	if o == nil {
		return nil
	}
	return o.Error
}
