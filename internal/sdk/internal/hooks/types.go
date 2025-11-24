package hooks

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/config"
	"golang.org/x/sync/singleflight"
)

type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// CriblTerraformHook implements both authentication and URL routing for Cribl Terraform API
type CriblTerraformHook struct {
	client      HTTPClient
	sessions    sync.Map
	baseURL     string
	orgID       string
	workspaceID string
}

type session struct {
	Credentials *credentials
	Token       string
	ExpiresAt   *int64
	Scopes      []string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   *int64 `json:"expires_in"`
}

type credentials struct {
	ClientID             string
	ClientSecret         string
	TokenURL             string
	Scopes               []string
	AdditionalProperties map[string]string
}

type clientCredentialsHook struct {
	client   HTTPClient
	sessions sync.Map

	// sessionsGroup prevents concurrent token refreshes.
	sessionsGroup *singleflight.Group
}

type CriblConfig struct {
	ClientID       string `json:"client_id" ini:"client_id"`
	ClientSecret   string `json:"client_secret" ini:"client_secret"`
	OrganizationID string `json:"organization_id" ini:"organization_id"`
	Workspace      string `json:"workspace" ini:"workspace"`
	CloudDomain    string `json:"cloud_domain" ini:"cloud_domain"`

	// On-prem configuration fields
	OnpremServerURL string `json:"onprem_server_url" ini:"onprem_server_url"`
	OnpremUsername  string `json:"onprem_username" ini:"onprem_username"`
	OnpremPassword  string `json:"onprem_password" ini:"onprem_password"`
}

type CriblConfigFile struct {
	Profiles map[string]CriblConfig `json:"profiles"`
}

type FailEarly struct {
	Cause error
}

// HTTPClient provides an interface for supplying the SDK with a custom HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HookContext struct {
	SDK              any
	SDKConfiguration config.SDKConfiguration
	BaseURL          string
	Context          context.Context
	OperationID      string
	OAuth2Scopes     []string
	SecuritySource   func(context.Context) (interface{}, error)
}

type BeforeRequestContext struct {
	HookContext
}

type AfterSuccessContext struct {
	HookContext
}

type AfterErrorContext struct {
	HookContext
}

// sdkInitHook is called when the SDK is initializing. The hook can modify and return a new baseURL and HTTP client to be used by the SDK.
type sdkInitHook interface {
	SDKInit(baseURL string, client HTTPClient) (string, HTTPClient)
}

// beforeRequestHook is called before the SDK sends a request. The hook can modify the request before it is sent or return an error to stop the request from being sent.
type beforeRequestHook interface {
	BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error)
}

// afterSuccessHook is called after the SDK receives a response. The hook can modify the response before it is handled or return an error to stop the response from being handled.
type afterSuccessHook interface {
	AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error)
}

// afterErrorHook is called after the SDK encounters an error, or a non-successful response. The hook can modify the response if available otherwise modify the error.
// All afterErrorHook hooks are called and returning an error won't stop the other hooks from being called. But if you want to stop the other hooks from being called, you can return a FailEarly error wrapping your error.
type afterErrorHook interface {
	AfterError(hookCtx AfterErrorContext, res *http.Response, err error) (*http.Response, error)
}

type Hooks struct {
	sdkInitHooks      []sdkInitHook
	beforeRequestHook []beforeRequestHook
	afterSuccessHook  []afterSuccessHook
	afterErrorHook    []afterErrorHook
}
