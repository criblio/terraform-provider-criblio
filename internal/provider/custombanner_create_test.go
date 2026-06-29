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

func TestPreserveCustomBannerPlanKeepsPlannedItems(t *testing.T) {
	plan := customBannerPlanObjectWithItems(t, bannerMessage{
		Created: types.Float64Null(),
		Enabled: types.BoolValue(true),
		ID:      types.StringValue(customBannerID),
		Message: types.StringValue("Scheduled maintenance window: Saturday 2am-4am UTC"),
		Theme:   types.StringValue("purple"),
		Type:    types.StringValue("custom"),
	})
	data := &CustomBannerResourceModel{
		Enabled: types.BoolValue(true),
		Items: []bannerMessage{{
			Enabled: types.BoolValue(true),
			ID:      types.StringValue(customBannerID),
			Message: types.StringValue("Maintenance complete. Systems are operating normally."),
			Theme:   types.StringValue("green"),
			Type:    types.StringValue("custom"),
		}},
		Message: types.StringValue("Maintenance complete. Systems are operating normally."),
		Theme:   types.StringValue("green"),
		Type:    types.StringValue("custom"),
	}

	preserveCustomBannerPlan(context.Background(), data, plan)

	if len(data.Items) != 1 {
		t.Fatalf("expected one planned item, got %d", len(data.Items))
	}
	if got := data.Items[0].Theme.ValueString(); got != "purple" {
		t.Fatalf("expected planned item theme to be preserved, got %q", got)
	}
	if got := data.Items[0].Message.ValueString(); got != "Scheduled maintenance window: Saturday 2am-4am UTC" {
		t.Fatalf("expected planned item message to be preserved, got %q", got)
	}
	if got := data.Theme.ValueString(); got != "purple" {
		t.Fatalf("expected planned top-level theme to be preserved, got %q", got)
	}
}

func customBannerPlanObjectWithItems(t *testing.T, item bannerMessage) types.Object {
	t.Helper()

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
	itemValue := types.ObjectValueMust(itemType.AttrTypes, map[string]attr.Value{
		"created":           item.Created,
		"custom_themes":     types.ListNull(types.StringType),
		"enabled":           item.Enabled,
		"id":                item.ID,
		"invert_font_color": types.BoolNull(),
		"link":              types.StringNull(),
		"link_display":      types.StringNull(),
		"message":           item.Message,
		"theme":             item.Theme,
		"type":              item.Type,
	})

	return types.ObjectValueMust(bannerType, map[string]attr.Value{
		"created":           types.Float64Null(),
		"custom_themes":     types.ListNull(types.StringType),
		"enabled":           types.BoolValue(true),
		"id":                types.StringValue(customBannerID),
		"invert_font_color": types.BoolNull(),
		"items":             types.ListValueMust(itemType, []attr.Value{itemValue}),
		"link":              types.StringValue("https://status.example.com"),
		"link_display":      types.StringValue("View status page"),
		"message":           item.Message,
		"theme":             item.Theme,
		"type":              item.Type,
	})
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
