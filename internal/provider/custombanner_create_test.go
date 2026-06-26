package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
