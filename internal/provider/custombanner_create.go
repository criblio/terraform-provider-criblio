package provider

import (
	"context"
	"errors"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

const (
	customBannerID   = "custom-banner"
	customBannerPath = "/system/banners/custom-banner"
)

func (r *CustomBannerResource) createCustomBanner(ctx context.Context, data *CustomBannerResourceModel) error {
	apiModel := data.toCustomBannerAPI()

	// POST /system/banners is not in the upstream spec. It is an undocumented
	// bootstrap endpoint that seeds the banner entity on fresh workspaces.
	// PATCH /system/banners/custom-banner returns 400 until this entity exists.
	// If POST returns 400/404/405, fall back to PATCH (entity already exists or
	// this Cribl version has no POST endpoint).
	if _, err := restclient.Post[customBannerAPI, customBannerAPI](ctx, r.client, "/system/banners", apiModel); err != nil {
		if !shouldPatchCustomBannerAfterPost(err) {
			return err
		}
		if _, patchErr := restclient.Patch[customBannerAPI, customBannerAPI](ctx, r.client, customBannerPath, apiModel); patchErr != nil {
			return patchErr
		}
	}

	return r.refreshCustomBannerState(ctx, data)
}

func shouldPatchCustomBannerAfterPost(err error) bool {
	if restclient.IsNotFound(err) {
		return true
	}
	var httpErr *restclient.HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode {
		case 400, 404, 405:
			return true
		}
	}
	return false
}
