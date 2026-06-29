package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestCustomBannerCreateUsesPostThenGet(t *testing.T) {
	var postCount, patchCount, getCount int
	server := newCustomBannerTestServer(t, http.StatusOK, &postCount, &patchCount, &getCount)
	defer server.Close()

	resource := &CustomBannerResource{client: restclient.New(restclient.Config{
		BaseURL:     server.URL,
		BearerToken: "test",
	})}
	model := customBannerTestModel()

	if err := resource.createCustomBanner(context.Background(), &model); err != nil {
		t.Fatalf("createCustomBanner returned error: %v", err)
	}
	if postCount != 1 || patchCount != 0 || getCount != 1 {
		t.Fatalf("POST/PATCH/GET counts = %d/%d/%d, want 1/0/1", postCount, patchCount, getCount)
	}
	if model.ID.ValueString() != customBannerID {
		t.Fatalf("model ID = %q, want %q", model.ID.ValueString(), customBannerID)
	}
}

func TestCustomBannerCreateFallsBackToPatch(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusMethodNotAllowed} {
		t.Run(fmt.Sprintf("post_%d", status), func(t *testing.T) {
			var postCount, patchCount, getCount int
			server := newCustomBannerTestServer(t, status, &postCount, &patchCount, &getCount)
			defer server.Close()

			resource := &CustomBannerResource{client: restclient.New(restclient.Config{
				BaseURL:     server.URL,
				BearerToken: "test",
			})}
			model := customBannerTestModel()

			if err := resource.createCustomBanner(context.Background(), &model); err != nil {
				t.Fatalf("createCustomBanner returned error: %v", err)
			}
			if postCount != 1 || patchCount != 1 || getCount != 1 {
				t.Fatalf("POST/PATCH/GET counts = %d/%d/%d, want 1/1/1", postCount, patchCount, getCount)
			}
		})
	}
}

func TestCustomBannerPlanDecodeAllowsUnknownComputedLists(t *testing.T) {
	stringListType := types.ListType{ElemType: types.StringType}
	itemType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"created":           types.Float64Type,
		"custom_themes":     stringListType,
		"enabled":           types.BoolType,
		"id":                types.StringType,
		"invert_font_color": types.BoolType,
		"link":              types.StringType,
		"link_display":      types.StringType,
		"message":           types.StringType,
		"theme":             types.StringType,
		"type":              types.StringType,
	}}
	bannerType := map[string]attr.Type{
		"created":           types.Float64Type,
		"custom_themes":     stringListType,
		"enabled":           types.BoolType,
		"id":                types.StringType,
		"invert_font_color": types.BoolType,
		"items":             types.ListType{ElemType: itemType},
		"link":              types.StringType,
		"link_display":      types.StringType,
		"message":           types.StringType,
		"theme":             types.StringType,
		"type":              types.StringType,
	}
	plan := types.ObjectValueMust(bannerType, map[string]attr.Value{
		"created":           types.Float64Unknown(),
		"custom_themes":     types.ListUnknown(types.StringType),
		"enabled":           types.BoolValue(true),
		"id":                types.StringUnknown(),
		"invert_font_color": types.BoolUnknown(),
		"items":             types.ListUnknown(itemType),
		"link":              types.StringValue("https://status.example.com"),
		"link_display":      types.StringValue("View status page"),
		"message":           types.StringValue("Scheduled maintenance"),
		"theme":             types.StringValue("purple"),
		"type":              types.StringValue("custom"),
	})

	var model CustomBannerResourceModel
	diags := plan.As(context.Background(), &model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})

	if diags.HasError() {
		t.Fatalf("expected unknown computed lists to decode without diagnostics, got: %v", diags)
	}
	if model.CustomThemes != nil {
		t.Fatalf("expected unknown custom_themes to decode as nil, got %#v", model.CustomThemes)
	}
	if model.Items != nil {
		t.Fatalf("expected unknown items to decode as nil, got %#v", model.Items)
	}
	if got := model.Message.ValueString(); got != "Scheduled maintenance" {
		t.Fatalf("expected message to decode, got %q", got)
	}
}

func newCustomBannerTestServer(t *testing.T, postStatus int, postCount, patchCount, getCount *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/system/banners":
			*postCount++
			w.WriteHeader(postStatus)
			if postStatus >= 200 && postStatus < 300 {
				_, _ = w.Write([]byte(customBannerTestEnvelope()))
			} else {
				_, _ = w.Write([]byte(`{"status":"error"}`))
			}
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v1/system/banners/custom-banner":
			*patchCount++
			_, _ = w.Write([]byte(customBannerTestEnvelope()))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/system/banners/custom-banner":
			*getCount++
			_, _ = w.Write([]byte(customBannerTestEnvelope()))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
}

func customBannerTestModel() CustomBannerResourceModel {
	return CustomBannerResourceModel{
		Enabled: types.BoolValue(true),
		Message: types.StringValue("hello"),
		Theme:   types.StringValue("purple"),
		Type:    types.StringValue("custom"),
	}
}

func customBannerTestEnvelope() string {
	return `{"items":[{"id":"custom-banner","enabled":true,"message":"hello","theme":"purple","type":"custom"}]}`
}
