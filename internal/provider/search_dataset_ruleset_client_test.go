package provider

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSearchDatasetRulesetMetricsUsesMetricsEndpoint(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got, want := r.URL.Path, "/api/v1/m/default_search/search/local_search/dataset-rulesets/metrics"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %q, want PATCH", r.Method)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"id":"metrics","rules":[]}`)),
			Request:    r,
		}, nil
	})}

	api := newSearchDatasetRulesetAPI(restclient.New(restclient.Config{
		BaseURL:     "https://example.test",
		BearerToken: "test-token",
		HTTPClient:  httpClient,
	}))
	result, err := api.Create(context.Background(), SearchDatasetRulesetModel{
		ID:    types.StringValue("metrics"),
		Rules: types.ListNull(types.ObjectType{AttrTypes: SearchDatasetRulesetRulesAttrTypes()}),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got := result.ID.ValueString(); got != "metrics" {
		t.Fatalf("ID = %q, want metrics", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
