package restclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/useragent"
)

type testItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestMethods(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Authorization = %q, expected Bearer test-token", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != useragent.TerraformProvider {
			t.Errorf("User-Agent = %q, expected %q", r.Header.Get("User-Agent"), useragent.TerraformProvider)
		}

		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("GET path = %q", r.URL.Path)
			}
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-get"})
		case http.MethodPost:
			switch r.URL.Path {
			case "/api/v1/system/certificates":
				assertJSONBody(t, r, "from-post")
				writeJSON(t, w, testItem{ID: "cert-2", Name: "from-post"})
			case "/api/v1/system/no-response":
				assertJSONBody(t, r, "from-post-no-response")
				writeJSON(t, w, map[string]any{"items": []testItem{}})
			case "/api/v1/system/files":
				mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil {
					t.Errorf("failed to parse content type: %v", err)
				}
				if mediaType != "multipart/form-data" {
					t.Errorf("upload content type = %q, expected multipart/form-data", mediaType)
				}
				if err := r.ParseMultipartForm(1024); err != nil {
					t.Errorf("failed to parse multipart form: %v", err)
				}
				file, header, err := r.FormFile("file")
				if err != nil {
					t.Errorf("missing upload file: %v", err)
					return
				}
				defer file.Close()
				content, err := io.ReadAll(file)
				if err != nil {
					t.Errorf("failed to read upload file: %v", err)
				}
				if header.Filename != "lookup.csv" {
					t.Errorf("upload filename = %q, expected lookup.csv", header.Filename)
				}
				if string(content) != "a,b\n" {
					t.Errorf("upload content = %q, expected a,b\\n", string(content))
				}
				w.WriteHeader(http.StatusNoContent)
			default:
				t.Errorf("unexpected POST path %q", r.URL.Path)
			}
		case http.MethodPatch:
			if r.URL.Path == "/api/v1/system/no-response" {
				assertJSONBody(t, r, "from-patch-no-response")
				writeJSON(t, w, map[string]any{"items": []testItem{}})
				return
			}
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("PATCH path = %q", r.URL.Path)
			}
			assertJSONBody(t, r, "from-patch")
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-patch"})
		case http.MethodPut:
			if r.URL.Path == "/api/v1/system/files" {
				if r.Header.Get("Content-Type") != "text/csv" {
					t.Errorf("raw upload content type = %q, expected text/csv", r.Header.Get("Content-Type"))
				}
				content, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("failed to read raw upload body: %v", err)
				}
				if string(content) != "raw-a,raw-b\n" {
					t.Errorf("raw upload content = %q, expected raw-a,raw-b\\n", string(content))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("PUT path = %q", r.URL.Path)
			}
			assertJSONBody(t, r, "from-put")
			writeJSON(t, w, testItem{ID: "cert-1", Name: "from-put"})
		case http.MethodDelete:
			if r.URL.Path != "/api/v1/system/certificates/cert-1" {
				t.Errorf("DELETE path = %q", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "from-get" {
		t.Fatalf("Get name = %q, expected from-get", got.Name)
	}

	created, err := Post[testItem, testItem](context.Background(), client, "/system/certificates", testItem{Name: "from-post"})
	if err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if created.ID != "cert-2" {
		t.Fatalf("Post ID = %q, expected cert-2", created.ID)
	}

	if err := PostNoResponse(context.Background(), client, "/system/no-response", testItem{Name: "from-post-no-response"}); err != nil {
		t.Fatalf("PostNoResponse returned error: %v", err)
	}

	patched, err := Patch[testItem, testItem](context.Background(), client, "/system/certificates/cert-1", testItem{Name: "from-patch"})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if patched.Name != "from-patch" {
		t.Fatalf("Patch name = %q, expected from-patch", patched.Name)
	}

	if err := PatchNoResponse(context.Background(), client, "/system/no-response", testItem{Name: "from-patch-no-response"}); err != nil {
		t.Fatalf("PatchNoResponse returned error: %v", err)
	}

	put, err := Put[testItem, testItem](context.Background(), client, "/system/certificates/cert-1", testItem{Name: "from-put"})
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	if put.Name != "from-put" {
		t.Fatalf("Put name = %q, expected from-put", put.Name)
	}

	if err := Delete(context.Background(), client, "/system/certificates/cert-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if err := Upload(context.Background(), client, "/system/files", "lookup.csv", []byte("a,b\n")); err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}

	if err := PutRawNoResponse(context.Background(), client, "/system/files", "text/csv", []byte("raw-a,raw-b\n")); err != nil {
		t.Fatalf("PutRawNoResponse returned error: %v", err)
	}
}

func TestDecodeEnvelopeSingleAndSlice(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"count": 2,
			"items": []testItem{
				{ID: "one", Name: "first"},
				{ID: "two", Name: "second"},
			},
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	single, err := Get[testItem](context.Background(), client, "/system/certificates/one")
	if err != nil {
		t.Fatalf("Get single returned error: %v", err)
	}
	if single.ID != "one" {
		t.Fatalf("single ID = %q, expected one", single.ID)
	}

	list, err := Get[[]testItem](context.Background(), client, "/system/certificates")
	if err != nil {
		t.Fatalf("Get slice returned error: %v", err)
	}
	if len(*list) != 2 {
		t.Fatalf("slice length = %d, expected 2", len(*list))
	}
	if (*list)[1].ID != "two" {
		t.Fatalf("second ID = %q, expected two", (*list)[1].ID)
	}
}

func TestDecodeEnvelopeEmptySingleIsNotFound(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"count": 0,
			"items": []testItem{},
		})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	_, err := Get[testItem](context.Background(), client, "/system/certificates/missing")
	if err == nil {
		t.Fatal("Get empty envelope returned nil error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound = false, expected true for empty single-resource envelope: %v", err)
	}

	list, err := Get[[]testItem](context.Background(), client, "/system/certificates")
	if err != nil {
		t.Fatalf("Get empty list returned error: %v", err)
	}
	if len(*list) != 0 {
		t.Fatalf("empty list length = %d, expected 0", len(*list))
	}
}

func TestDecodePlainJSONAndNoContent(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/empty") {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeJSON(t, w, testItem{ID: "plain", Name: "plain-json"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	plain, err := Get[testItem](context.Background(), client, "/system/plain")
	if err != nil {
		t.Fatalf("Get plain returned error: %v", err)
	}
	if plain.ID != "plain" {
		t.Fatalf("plain ID = %q, expected plain", plain.ID)
	}

	empty, err := Patch[testItem, testItem](context.Background(), client, "/system/plain/empty", testItem{Name: "empty"})
	if err != nil {
		t.Fatalf("Patch empty returned error: %v", err)
	}
	if empty != nil {
		t.Fatalf("Patch empty = %#v, expected nil", empty)
	}
}

func TestErrors(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/missing":
			http.Error(w, "missing resource", http.StatusNotFound)
		case "/api/v1/bad":
			http.Error(w, "bad request", http.StatusBadRequest)
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	_, err := Get[testItem](context.Background(), client, "/missing")
	if err == nil {
		t.Fatal("Get missing returned nil error")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound = false, expected true for 404")
	}

	_, err = Get[testItem](context.Background(), client, "/bad")
	if err == nil {
		t.Fatal("Get bad returned nil error")
	}
	if IsNotFound(err) {
		t.Fatalf("IsNotFound = true, expected false for non-404")
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error type = %T, expected HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, expected 400", httpErr.StatusCode)
	}
	if !strings.Contains(httpErr.Body, "bad request") {
		t.Fatalf("body = %q, expected bad request", httpErr.Body)
	}

	if IsNotFound(fmt.Errorf("ordinary error")) {
		t.Fatalf("IsNotFound = true, expected false for ordinary error")
	}
}

func TestGetRetriesTransientECONNRESETResponse(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"status":"error","message":"Read ECONNRESET","error":{"errno":-104,"code":"ECONNRESET","syscall":"read"}}`))
			return
		}
		writeJSON(t, w, testItem{ID: "pipeline-1", Name: "pipeline"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	got, err := Get[testItem](context.Background(), client, "/pipelines/pipeline-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "pipeline" {
		t.Fatalf("name = %q, expected pipeline", got.Name)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestPostDoesNotRetryTransientECONNRESETResponse(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":"error","message":"Read ECONNRESET","error":{"errno":-104,"code":"ECONNRESET","syscall":"read"}}`))
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
	})

	_, err := Post[testItem, testItem](context.Background(), client, "/pipelines", testItem{Name: "pipeline"})
	if err == nil {
		t.Fatal("Post returned nil error, expected HTTP 500")
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected 1", requestCount)
	}
}

func TestReplaySafeRequestsRetryTooManyRequests(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "generic GET", method: http.MethodGet, path: "/any/endpoint"},
		{name: "generic PUT", method: http.MethodPut, path: "/any/endpoint"},
		{name: "generic DELETE", method: http.MethodDelete, path: "/any/endpoint"},
		{name: "product group POST", method: http.MethodPost, path: "/products/stream/groups"},
		{name: "legacy group POST", method: http.MethodPost, path: "/master/groups"},
		{name: "fleet PATCH", method: http.MethodPatch, path: "/m/fleet-id/pipelines/id"},
		{name: "app fleet POST", method: http.MethodPost, path: "/a/app-id/m/fleet-id/sources"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				if r.Method != test.method {
					t.Errorf("method = %q, expected %q", r.Method, test.method)
				}
				if test.method != http.MethodGet && test.method != http.MethodDelete {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("failed to read body: %v", err)
					}
					if string(body) != `{"name":"request"}` {
						t.Errorf("body = %q, expected request body", body)
					}
				}
				if requestCount == 1 {
					w.Header().Set("Retry-After", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
				writeJSON(t, w, testItem{ID: "retried", Name: "success"})
			}))
			defer server.Close()

			client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
			if _, err := do(context.Background(), client, test.method, test.path, "application/json", []byte(`{"name":"request"}`)); err != nil {
				t.Fatalf("%s returned error: %v", test.method, err)
			}
			if requestCount != 2 {
				t.Fatalf("request count = %d, expected 2", requestCount)
			}
		})
	}
}

func TestUnsafePostDoesNotRetryTooManyRequests(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	_, err := Post[testItem, testItem](context.Background(), client, "/unrelated", testItem{Name: "request"})
	if err == nil {
		t.Fatal("Post returned nil error, expected HTTP 429")
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected 1", requestCount)
	}
}

func TestFleetUploadRetriesTooManyRequestsWithBody(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("missing upload file: %v", err)
			return
		}
		defer file.Close()
		body, err := io.ReadAll(file)
		if err != nil {
			t.Errorf("failed to read upload body: %v", err)
			return
		}
		if header.Filename != "lookup.csv" || string(body) != "a,b\n" {
			t.Errorf("upload = (%q, %q), expected lookup.csv content", header.Filename, body)
		}
		if requestCount == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	if err := Upload(context.Background(), client, "/m/fleet-id/system/lookups", "lookup.csv", []byte("a,b\n")); err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, expected 2", requestCount)
	}
}

func TestRetryAfterDelay(t *testing.T) {
	now := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		value string
		want  time.Duration
		ok    bool
	}{
		{name: "seconds", value: "15", want: 15 * time.Second, ok: true},
		{name: "HTTP date", value: now.Add(45 * time.Second).Format(http.TimeFormat), want: 45 * time.Second, ok: true},
		{name: "past HTTP date", value: now.Add(-time.Second).Format(http.TimeFormat), want: 0, ok: true},
		{name: "missing", value: "", ok: false},
		{name: "invalid", value: "later", ok: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := retryAfterDelay(test.value, now)
			if ok != test.ok || got != test.want {
				t.Fatalf("retryAfterDelay(%q) = (%s, %t), expected (%s, %t)", test.value, got, ok, test.want, test.ok)
			}
		})
	}
}

func TestRetryAfterWaitHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := waitBeforeAPIRetry(ctx, nil, 0, "60"); !errors.Is(err, context.Canceled) {
		t.Fatalf("waitBeforeAPIRetry error = %v, expected context cancellation", err)
	}
}

func TestRetryAfterOverLimitDoesNotWait(t *testing.T) {
	values := []string{
		"61",
		"9999999999",
		"18446744073709551616",
	}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			err := waitBeforeAPIRetry(context.Background(), nil, 0, value)
			if !errors.Is(err, errRetryAfterExceedsLimit) {
				t.Fatalf("waitBeforeAPIRetry error = %v, expected maximum delay error", err)
			}
		})
	}
}

func TestRetryWaitBudgetReturnsOriginalResponse(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "1")
		http.Error(w, "admission throttled", http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:         server.URL,
		BearerToken:     "test-token",
		RetryWaitBudget: time.Nanosecond,
	})
	_, err := Get[testItem](context.Background(), client, "/any/endpoint")
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("Get error = %v, expected original HTTP 429", err)
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected no retry", requestCount)
	}
}

func TestSuccessfulRequestResetsRetryWaitBudget(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCounts := make(map[string]int)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounts[r.URL.Path]++
		if requestCounts[r.URL.Path] == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeJSON(t, w, testItem{ID: "retried", Name: "success"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:         server.URL,
		BearerToken:     "test-token",
		RetryWaitBudget: time.Millisecond,
	})
	if _, err := Get[testItem](context.Background(), client, "/first"); err != nil {
		t.Fatalf("first Get returned error: %v", err)
	}
	if _, err := Get[testItem](context.Background(), client, "/second"); err != nil {
		t.Fatalf("second Get returned error after budget reset: %v", err)
	}
	if requestCounts["/api/v1/first"] != 2 {
		t.Fatalf("first request count = %d, expected 2", requestCounts["/api/v1/first"])
	}
	if requestCounts["/api/v1/second"] != 2 {
		t.Fatalf("second request count = %d, expected 2", requestCounts["/api/v1/second"])
	}
}

func TestConcurrentRetryWaitsCountWallClockOnce(t *testing.T) {
	client := New(Config{RetryWaitBudget: 1100 * time.Millisecond})
	start := make(chan struct{})
	results := make(chan error, 5)

	for range 5 {
		go func() {
			<-start
			results <- waitBeforeAPIRetry(context.Background(), client, 0, "1")
		}()
	}
	close(start)
	for range 5 {
		if err := <-results; err != nil {
			t.Fatalf("concurrent retry wait returned error: %v", err)
		}
	}
	if err := waitBeforeAPIRetry(context.Background(), client, 0, "1"); !errors.Is(err, errRetryWaitBudgetExhausted) {
		t.Fatalf("subsequent retry wait error = %v, expected exhausted wall-clock budget", err)
	}
}

func TestServiceUnavailableIgnoresRetryAfter(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeJSON(t, w, testItem{ID: "retried", Name: "success"})
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	started := time.Now()
	if _, err := Get[testItem](context.Background(), client, "/any/endpoint"); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("503 retry took %s, expected exponential fallback rather than Retry-After", elapsed)
	}
}

func TestOverflowingRetryAfterReturnsOriginalResponse(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "18446744073709551616")
		http.Error(w, "admission throttled", http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	_, err := Get[testItem](context.Background(), client, "/any/endpoint")
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("Get error = %v, expected original HTTP 429", err)
	}
	if !strings.Contains(httpErr.Body, "admission throttled") {
		t.Fatalf("HTTP error body = %q, expected original response body", httpErr.Body)
	}
	if requestCount != 1 {
		t.Fatalf("request count = %d, expected no retry", requestCount)
	}
}

func TestRetryWaitsAtLeastRetryAfter(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	requestTimes := make([]time.Time, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTimes = append(requestTimes, time.Now())
		if len(requestTimes) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeJSON(t, w, testItem{ID: "retried", Name: "success"})
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	if _, err := Post[testItem, testItem](context.Background(), client, "/products/stream/groups", testItem{Name: "request"}); err != nil {
		t.Fatalf("Post returned error: %v", err)
	}
	if len(requestTimes) != 2 {
		t.Fatalf("request count = %d, expected 2", len(requestTimes))
	}
	if elapsed := requestTimes[1].Sub(requestTimes[0]); elapsed < time.Second {
		t.Fatalf("retry occurred after %s, expected at least 1s", elapsed)
	}
}

func TestTooManyRequestsStopsAfterRetryLimit(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token"})
	_, err := Post[testItem, testItem](context.Background(), client, "/products/stream/groups", testItem{Name: "request"})
	if err == nil {
		t.Fatal("Post returned nil error, expected HTTP 429")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("Post error = %v, expected HTTP 429", err)
	}
	if delay, ok := RetryAfter(err); !ok || delay != 0 {
		t.Fatalf("RetryAfter(error) = (%s, %t), expected (0, true)", delay, ok)
	}
	if requestCount != defaultAPIRetryMax+1 {
		t.Fatalf("request count = %d, expected %d", requestCount, defaultAPIRetryMax+1)
	}
}

func TestConfiguredRetryLimit(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{BaseURL: server.URL, BearerToken: "test-token", RetryMax: 7})
	_, err := Get[testItem](context.Background(), client, "/any/endpoint")
	if err == nil {
		t.Fatal("Get returned nil error, expected HTTP 429")
	}
	if requestCount != 8 {
		t.Fatalf("request count = %d, expected 8", requestCount)
	}
}

func TestConfigHelperAdmissionRetriesBeyondOrdinaryLimit(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount <= defaultAPIRetryMax+2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeJSON(t, w, testItem{ID: "retried", Name: "success"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:                  server.URL,
		BearerToken:              "test-token",
		ConfigHelperRetryTimeout: time.Second,
	})
	if _, err := Get[testItem](context.Background(), client, "/m/fleet-id/pipelines"); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if requestCount != defaultAPIRetryMax+3 {
		t.Fatalf("request count = %d, expected %d", requestCount, defaultAPIRetryMax+3)
	}
}

func TestConfigHelperAdmissionTimeoutReturnsHTTPError(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	fastAPIRetry(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:                  server.URL,
		BearerToken:              "test-token",
		ConfigHelperRetryTimeout: 10 * time.Millisecond,
	})
	_, err := Get[testItem](context.Background(), client, "/m/fleet-id/pipelines")
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("Get error = %v, expected original HTTP 429", err)
	}
}

func TestGatewayRouting(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	tests := []struct {
		name string
		path string
	}{
		{
			name: "v1 path",
			path: "/v1/organizations/org-id/workspaces",
		},
		{
			name: "api v1 path",
			path: "/api/v1/organizations/org-id/workspaces",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "gateway.cribl.cloud" {
			t.Errorf("Host = %q, expected gateway.cribl.cloud", r.Host)
		}
		if r.URL.Path != "/v1/organizations/org-id/workspaces" {
			t.Errorf("path = %q, expected /v1/organizations/org-id/workspaces", r.URL.Path)
		}
		writeJSON(t, w, testItem{ID: "workspace", Name: "gateway"})
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
		HTTPClient: &http.Client{
			Transport: rewriteTransport{target: targetURL},
		},
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Get[testItem](context.Background(), client, test.path)
			if err != nil {
				t.Fatalf("Get gateway returned error: %v", err)
			}
			if got.Name != "gateway" {
				t.Fatalf("gateway name = %q, expected gateway", got.Name)
			}
		})
	}
}

func TestProviderWorkspaceRouting(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")
	t.Setenv("CRIBL_ORGANIZATION_ID", "")
	t.Setenv("CRIBL_WORKSPACE_ID", "")
	t.Setenv("CRIBL_CLOUD_DOMAIN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "provider-workspace-provider-org.cribl-playground.cloud" {
			t.Errorf("Host = %q, expected provider-workspace-provider-org.cribl-playground.cloud", r.Host)
		}
		if r.URL.Path != "/api/v1/system/certificates/cert-1" {
			t.Errorf("path = %q, expected /api/v1/system/certificates/cert-1", r.URL.Path)
		}
		writeJSON(t, w, testItem{ID: "cert-1", Name: "provider-config"})
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}

	client := New(Config{
		ProviderOrgID:       "provider-org",
		ProviderWorkspaceID: "provider-workspace",
		ProviderCloudDomain: "cribl-playground.cloud",
		BearerToken:         "test-token",
		HTTPClient: &http.Client{
			Transport: rewriteTransport{target: targetURL},
		},
	})

	got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "provider-config" {
		t.Fatalf("name = %q, expected provider-config", got.Name)
	}
}

func TestAuthenticationRequired(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	client := New(Config{BaseURL: "http://127.0.0.1:1"})
	_, err := Get[testItem](context.Background(), client, "/system/certificates")
	if err == nil {
		t.Fatal("Get returned nil error, expected authentication error")
	}
	if !strings.Contains(err.Error(), "authentication requires bearer token or credentials") {
		t.Fatalf("error = %q, expected authentication message", err.Error())
	}
}

func TestUnauthorizedInvalidatesTokenAndRetriesRequest(t *testing.T) {
	auth.ClearTokenCache()
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	loginCount := 0
	apiCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/login":
			loginCount++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"token":"token-%d","forcePasswordChange":false}`, loginCount)))
		case "/api/v1/system/certificates/cert-1":
			apiCount++
			switch r.Header.Get("Authorization") {
			case "Bearer token-1":
				http.Error(w, "stale token", http.StatusUnauthorized)
			case "Bearer token-2":
				writeJSON(t, w, testItem{ID: "cert-1", Name: "fresh-token"})
			default:
				t.Errorf("unexpected Authorization = %q", r.Header.Get("Authorization"))
			}
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(Config{
		Credentials: &auth.CriblConfig{
			OnpremServerURL: server.URL,
			OnpremUsername:  "admin",
			OnpremPassword:  "secret",
		},
	})

	got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "fresh-token" {
		t.Fatalf("name = %q, expected fresh-token", got.Name)
	}
	if loginCount != 2 {
		t.Fatalf("login count = %d, expected 2", loginCount)
	}
	if apiCount != 2 {
		t.Fatalf("api count = %d, expected 2", apiCount)
	}
}

func TestConcurrentUnauthorizedRefreshesTokenOnce(t *testing.T) {
	auth.ClearTokenCache()
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	var mu sync.Mutex
	loginCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/auth/login":
			mu.Lock()
			loginCount++
			token := "stale-token"
			if loginCount > 1 {
				token = "fresh-token"
			}
			mu.Unlock()
			if token == "fresh-token" {
				time.Sleep(50 * time.Millisecond)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"token":"%s","forcePasswordChange":false}`, token)))
		case "/api/v1/system/certificates/cert-1":
			switch r.Header.Get("Authorization") {
			case "Bearer stale-token":
				http.Error(w, "stale token", http.StatusUnauthorized)
			case "Bearer fresh-token":
				writeJSON(t, w, testItem{ID: "cert-1", Name: "fresh-token"})
			default:
				t.Errorf("unexpected Authorization = %q", r.Header.Get("Authorization"))
			}
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	credentials := &auth.CriblConfig{
		OnpremServerURL: server.URL,
		OnpremUsername:  "admin",
		OnpremPassword:  "secret",
	}
	if token, err := auth.GetToken(context.Background(), credentials); err != nil {
		t.Fatalf("preload GetToken returned error: %v", err)
	} else if token != "stale-token" {
		t.Fatalf("preload token = %q, expected stale-token", token)
	}

	client := New(Config{Credentials: credentials})

	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := Get[testItem](context.Background(), client, "/system/certificates/cert-1")
			if err != nil {
				errs <- err
				return
			}
			if got.Name != "fresh-token" {
				errs <- fmt.Errorf("name = %q, expected fresh-token", got.Name)
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent Get returned error: %v", err)
		}
	}

	mu.Lock()
	gotLoginCount := loginCount
	mu.Unlock()
	if gotLoginCount != 2 {
		t.Fatalf("login count = %d, expected 2", gotLoginCount)
	}
}

func TestUserAgentOverride(t *testing.T) {
	t.Setenv("CRIBL_BEARER_TOKEN", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "custom-agent" {
			t.Errorf("User-Agent = %q, expected custom-agent", r.Header.Get("User-Agent"))
		}
		writeJSON(t, w, testItem{ID: "agent", Name: "custom"})
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:     server.URL,
		BearerToken: "test-token",
		UserAgent:   "custom-agent",
	})

	if _, err := Get[testItem](context.Background(), client, "/system/certificates"); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
}

type rewriteTransport struct {
	target *url.URL
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	out := req.Clone(req.Context())
	out.URL.Scheme = t.target.Scheme
	out.URL.Host = t.target.Host
	out.Host = req.URL.Host
	return http.DefaultTransport.RoundTrip(out)
}

func assertJSONBody(t *testing.T, r *http.Request, expectedName string) {
	t.Helper()

	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, expected application/json", r.Header.Get("Content-Type"))
	}
	var body testItem
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Errorf("failed to decode JSON body: %v", err)
		return
	}
	if body.Name != expectedName {
		t.Errorf("body name = %q, expected %q", body.Name, expectedName)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("failed to write JSON response: %v", err)
	}
}

func fastAPIRetry(t *testing.T) {
	t.Helper()

	oldMin := apiRetryWaitMin
	oldMax := apiRetryWaitMax
	apiRetryWaitMin = time.Millisecond
	apiRetryWaitMax = time.Millisecond
	t.Cleanup(func() {
		apiRetryWaitMin = oldMin
		apiRetryWaitMax = oldMax
	})
}
