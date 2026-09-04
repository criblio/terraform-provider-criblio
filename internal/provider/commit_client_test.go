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

func TestCommitCreateReturnsCreatedCommit(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/version/commit" {
			t.Fatalf("request = %s %s, want POST /api/v1/version/commit", r.Method, r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request: %v", err)
		}
		if strings.Contains(string(body), `"items"`) {
			t.Fatalf("computed items sent in request: %s", body)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{
  "count": 1,
  "items": [{
    "branch": "main",
    "commit": "a1b2c3d4e5f6",
    "summary": {"changes": 1, "deletions": 0, "insertions": 1}
  }]
}`)),
			Request: r,
		}, nil
	})}

	api := newCommitAPI(restclient.New(restclient.Config{
		BaseURL:     "https://example.test",
		BearerToken: "test",
		HTTPClient:  httpClient,
	}))
	commit, err := api.Create(context.Background(), CommitModel{
		Group:   types.StringValue("default"),
		Message: types.StringValue("test"),
	})
	if err != nil {
		t.Fatalf("create commit: %v", err)
	}
	items := commit.Items.Elements()
	if len(items) != 1 {
		t.Fatalf("items length = %d, want 1", len(items))
	}
	item := items[0].(types.Object).Attributes()
	if got := item["commit"].(types.String).ValueString(); got != "a1b2c3d4e5f6" {
		t.Fatalf("commit = %q, want a1b2c3d4e5f6", got)
	}
}
