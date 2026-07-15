// Package restclient provides the generic HTTP client used by migrated resources.
package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/auth"
	"github.com/criblio/terraform-provider-criblio/internal/useragent"
)

const apiRetryMax = 3

var (
	apiRetryWaitMin = 500 * time.Millisecond
	apiRetryWaitMax = 2 * time.Second
)

// Config holds REST client settings.
type Config struct {
	BaseURL             string
	ProviderOrgID       string
	ProviderWorkspaceID string
	ProviderCloudDomain string
	Credentials         *auth.CriblConfig
	BearerToken         string
	HTTPClient          *http.Client
	UserAgent           string
}

// Client sends authenticated requests to Cribl APIs.
type Client struct {
	baseURL             string
	providerCloudDomain string
	credentials         *auth.CriblConfig
	bearerToken         string
	httpClient          *http.Client
	userAgent           string
}

// HTTPError is returned for non-2xx responses other than 404.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("cribl API returned HTTP %d: %s", e.StatusCode, e.Body)
}

// NotFoundError is returned for 404 responses.
type NotFoundError struct {
	Path string
	Body string
}

func (e *NotFoundError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("cribl API path %q was not found", e.Path)
	}
	return fmt.Sprintf("cribl API path %q was not found: %s", e.Path, e.Body)
}

// New creates a REST client.
func New(config Config) *Client {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	agent := config.UserAgent
	if agent == "" {
		agent = useragent.TerraformProvider
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = auth.ConstructBaseURL(auth.ConstructBaseURLInput{
			ProviderOrgID:       config.ProviderOrgID,
			ProviderWorkspaceID: config.ProviderWorkspaceID,
			ProviderCloudDomain: config.ProviderCloudDomain,
		}, config.Credentials)
	}

	return &Client{
		baseURL:             strings.TrimRight(baseURL, "/"),
		providerCloudDomain: config.ProviderCloudDomain,
		credentials:         config.Credentials,
		bearerToken:         config.BearerToken,
		httpClient:          httpClient,
		userAgent:           agent,
	}
}

// Get sends a GET request and decodes the response.
func Get[T any](ctx context.Context, c *Client, path string) (*T, error) {
	body, err := do(ctx, c, http.MethodGet, path, "", nil)
	if err != nil {
		return nil, err
	}
	return decodeResponse[T](path, body)
}

// Post sends a POST request with a JSON body and decodes the response.
func Post[Req, Resp any](ctx context.Context, c *Client, path string, body Req) (*Resp, error) {
	responseBody, err := doJSON(ctx, c, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return decodeResponse[Resp](path, responseBody)
}

// PostNoResponse sends a POST request with a JSON body and ignores the response body.
func PostNoResponse[Req any](ctx context.Context, c *Client, path string, body Req) error {
	_, err := doJSON(ctx, c, http.MethodPost, path, body)
	return err
}

// Patch sends a PATCH request with a JSON body and decodes the response.
func Patch[Req, Resp any](ctx context.Context, c *Client, path string, body Req) (*Resp, error) {
	responseBody, err := doJSON(ctx, c, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return decodeResponse[Resp](path, responseBody)
}

// PatchNoResponse sends a PATCH request with a JSON body and ignores the response body.
func PatchNoResponse[Req any](ctx context.Context, c *Client, path string, body Req) error {
	_, err := doJSON(ctx, c, http.MethodPatch, path, body)
	return err
}

// Put sends a PUT request with a JSON body and decodes the response.
func Put[Req, Resp any](ctx context.Context, c *Client, path string, body Req) (*Resp, error) {
	responseBody, err := doJSON(ctx, c, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return decodeResponse[Resp](path, responseBody)
}

// PutNoResponse sends a PUT request with a JSON body and ignores the response body.
func PutNoResponse[Req any](ctx context.Context, c *Client, path string, body Req) error {
	_, err := doJSON(ctx, c, http.MethodPut, path, body)
	return err
}

// PutRawNoResponse sends a PUT request with the provided content type and raw body.
func PutRawNoResponse(ctx context.Context, c *Client, path, contentType string, body []byte) error {
	_, err := PutRaw(ctx, c, path, contentType, body)
	return err
}

// PutRaw sends a PUT request with the provided content type and raw body and returns the response body.
func PutRaw(ctx context.Context, c *Client, path, contentType string, body []byte) ([]byte, error) {
	return do(ctx, c, http.MethodPut, path, contentType, body)
}

// PatchRaw sends a PATCH request with the provided content type and raw body and returns the response body.
func PatchRaw(ctx context.Context, c *Client, path, contentType string, body []byte) ([]byte, error) {
	return do(ctx, c, http.MethodPatch, path, contentType, body)
}

// Delete sends a DELETE request.
func Delete(ctx context.Context, c *Client, path string) error {
	_, err := do(ctx, c, http.MethodDelete, path, "", nil)
	return err
}

// Upload sends multipart file content to path.
func Upload(ctx context.Context, c *Client, path, filename string, content []byte) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create multipart file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		return fmt.Errorf("failed to write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	_, err = do(ctx, c, http.MethodPost, path, writer.FormDataContentType(), body.Bytes())
	return err
}

// IsNotFound reports whether err is a NotFoundError.
func IsNotFound(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}

func doJSON[Req any](ctx context.Context, c *Client, method, path string, body Req) ([]byte, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}
	return do(ctx, c, method, path, "application/json", requestBody)
}

func do(ctx context.Context, c *Client, method, path, contentType string, body []byte) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("restclient client is required")
	}

	for attempt := 0; ; attempt++ {
		responseBody, statusCode, token, err := c.send(ctx, method, path, contentType, body)
		if err != nil {
			if shouldRetryAPIRequest(method, 0, nil, err, attempt) {
				if err := waitBeforeAPIRetry(ctx, attempt); err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}

		if statusCode == http.StatusUnauthorized && c.bearerToken == "" && os.Getenv("CRIBL_BEARER_TOKEN") == "" && c.credentials != nil {
			auth.InvalidateTokenValue(c.credentials, token)
			responseBody, statusCode, _, err = c.send(ctx, method, path, contentType, body)
			if err != nil {
				if shouldRetryAPIRequest(method, 0, nil, err, attempt) {
					if err := waitBeforeAPIRetry(ctx, attempt); err != nil {
						return nil, err
					}
					continue
				}
				return nil, err
			}
		}

		err = responseError(path, statusCode, responseBody)
		if err == nil {
			return responseBody, nil
		}
		if shouldRetryAPIRequest(method, statusCode, responseBody, nil, attempt) {
			if err := waitBeforeAPIRetry(ctx, attempt); err != nil {
				return nil, err
			}
			continue
		}
		return responseBody, err
	}
}

func (c *Client) send(ctx context.Context, method, path, contentType string, body []byte) ([]byte, int, string, error) {
	requestURL, err := c.requestURL(path)
	if err != nil {
		return nil, 0, "", err
	}

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to create request: %v", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Accept", "application/json")
	}

	token, err := c.token(ctx)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to authenticate for %s %s: %v", method, path, err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to read response body: %v", err)
	}

	return responseBody, resp.StatusCode, token, nil
}

func (c *Client) requestURL(path string) (string, error) {
	trimmedPath := auth.TrimPath(path)

	if auth.IsOnPrem(c.credentials) {
		if auth.IsRestrictedOnPremEndpoint(trimmedPath) {
			return "", fmt.Errorf("endpoint %q is not supported for on-prem deployments", trimmedPath)
		}
		return joinBasePath(c.credentials.OnpremServerURL, "api/v1", trimmedPath), nil
	}

	if auth.IsGatewayPath(path) {
		return joinBasePath(auth.ConstructGatewayURL(c.providerCloudDomain, c.credentials), "", gatewayRequestPath(path)), nil
	}

	if c.baseURL == "" {
		return "", fmt.Errorf("base URL is required")
	}
	return joinBasePath(c.baseURL, "api/v1", trimmedPath), nil
}

func (c *Client) token(ctx context.Context) (string, error) {
	switch {
	case c.bearerToken != "":
		return c.bearerToken, nil
	case os.Getenv("CRIBL_BEARER_TOKEN") != "":
		return os.Getenv("CRIBL_BEARER_TOKEN"), nil
	case c.credentials != nil:
		return auth.GetToken(ctx, c.credentials)
	default:
		return "", fmt.Errorf("authentication requires bearer token or credentials")
	}
}

func responseError(path string, statusCode int, body []byte) error {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return nil
	case statusCode == http.StatusNotFound:
		return &NotFoundError{
			Path: path,
			Body: string(body),
		}
	default:
		return &HTTPError{
			StatusCode: statusCode,
			Body:       string(body),
		}
	}
}

func shouldRetryAPIRequest(method string, statusCode int, body []byte, err error, attempt int) bool {
	if attempt >= apiRetryMax || !isRetryableAPIMethod(method) {
		return false
	}
	if err != nil {
		message := strings.ToLower(err.Error())
		return strings.Contains(message, "failed to send request") ||
			strings.Contains(message, "failed to read response body")
	}
	switch statusCode {
	case http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	case http.StatusInternalServerError:
		return responseBodyHasTransientError(body)
	default:
		return false
	}
}

func isRetryableAPIMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func responseBodyHasTransientError(body []byte) bool {
	message := strings.ToLower(string(body))
	return strings.Contains(message, "econnreset") ||
		strings.Contains(message, "connection reset") ||
		strings.Contains(message, "socket hang up")
}

func waitBeforeAPIRetry(ctx context.Context, attempt int) error {
	wait := apiRetryWaitMin << attempt
	if wait > apiRetryWaitMax {
		wait = apiRetryWaitMax
	}
	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func decodeResponse[T any](path string, body []byte) (*T, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, nil
	}

	var envelope struct {
		Items json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Items) > 0 && string(envelope.Items) != "null" {
		return decodeEnvelope[T](path, envelope.Items)
	}

	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("[")) {
		return decodeEnvelope[T](path, body)
	}

	var output T
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &output, nil
}

func decodeEnvelope[T any](path string, items json.RawMessage) (*T, error) {
	var output T
	outputValue := reflect.ValueOf(&output).Elem()
	if outputValue.Kind() == reflect.Slice {
		if err := json.Unmarshal(items, &output); err != nil {
			return nil, fmt.Errorf("failed to decode response envelope: %v", err)
		}
		return &output, nil
	}

	itemsValue := reflect.New(reflect.SliceOf(outputValue.Type()))
	if err := json.Unmarshal(items, itemsValue.Interface()); err != nil {
		return nil, fmt.Errorf("failed to decode response envelope: %v", err)
	}

	itemsSlice := itemsValue.Elem()
	if itemsSlice.Len() == 0 {
		return nil, &NotFoundError{Path: path, Body: `{"items":[]}`}
	}

	outputValue.Set(itemsSlice.Index(0))
	return &output, nil
}

func joinBasePath(baseURL, prefix, path string) string {
	parts := []string{strings.TrimRight(baseURL, "/")}
	if prefix != "" {
		parts = append(parts, strings.Trim(prefix, "/"))
	}
	if path != "" {
		parts = append(parts, strings.TrimLeft(path, "/"))
	}
	return strings.Join(parts, "/")
}

func gatewayRequestPath(path string) string {
	cleanPath := strings.TrimLeft(path, "/")
	return strings.TrimPrefix(cleanPath, "api/")
}
