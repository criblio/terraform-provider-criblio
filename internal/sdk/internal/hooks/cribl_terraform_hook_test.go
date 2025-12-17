package hooks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
)

func TestNewCriblTerraformHook(t *testing.T) {
	var testData CriblTerraformHook

	output := NewCriblTerraformHook()
	if fmt.Sprintf("%+v", &testData) != fmt.Sprintf("%+v", output) {
		t.Errorf("NewCriblTerraformHook returned incorrect struct, expected '%+v' got '%+v'", &testData, output)
	}
}

func TestTerraformSDKInit(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")

	myUrl := "foobar"
	var myClient HTTPClient

	test := NewCriblTerraformHook()
	url, _ := test.SDKInit(myUrl, myClient)

	// With environment variables set, the hook should construct the URL from them
	expectedURL := "https://punk'n-biz.cribl.cloud"
	if url != expectedURL {
		t.Errorf("creds hook init returned %s, expected %s", url, expectedURL)
	}
	if test.client != myClient {
		t.Errorf("creds hook init test.client %+v, expected %+v", test.client, myClient)
	}
	if test.baseURL != expectedURL {
		t.Errorf("creds hook init test.baseURL %s, expected %s", test.baseURL, expectedURL)
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook.orgID returned %s, expected %s", test.orgID, "biz")
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook.workspaceID returned %s, expected %s", test.orgID, "biz")
	}
}

func TestTerraformSDKInitGetCredentialsFailure(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("HOME", "/root")

	var myClient HTTPClient

	myUrl := "foobar"
	test := NewCriblTerraformHook()
	urlOut, _ := test.SDKInit(myUrl, myClient)

	if urlOut != myUrl {
		t.Errorf("SDKInit returned incorrect url - expected %s, got %s", myUrl, urlOut)
	}
}

func TestTerraformSDKInitWithCloudDomain(t *testing.T) {
	// Set all environment variables including cloud domain
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-playground.cloud")

	myUrl := "should-be-overridden"
	var myClient HTTPClient

	test := NewCriblTerraformHook()
	url, _ := test.SDKInit(myUrl, myClient)

	// Should construct URL from all environment variables
	expectedURL := "https://test-workspace-test-org.cribl-playground.cloud"
	if url != expectedURL {
		t.Errorf("creds hook init returned %s, expected %s", url, expectedURL)
	}
	if test.baseURL != expectedURL {
		t.Errorf("creds hook init test.baseURL %s, expected %s", test.baseURL, expectedURL)
	}
	if test.orgID != "test-org" {
		t.Errorf("*CriblTerraformHook.orgID returned %s, expected %s", test.orgID, "test-org")
	}
	if test.workspaceID != "test-workspace" {
		t.Errorf("*CriblTerraformHook.workspaceID returned %s, expected %s", test.workspaceID, "test-workspace")
	}

	// Clean up environment variables
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
}

func TestTerraformBeforeRequestWithCloudDomain(t *testing.T) {
	// Set all environment variables including cloud domain
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "staging-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "staging-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-staging.cloud")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-bearer-token")

	myUrl := "should-be-overridden"
	var myClient HTTPClient

	test := NewCriblTerraformHook()
	test.SDKInit(myUrl, myClient)

	// Create a test request
	myReq, _ := http.NewRequest("GET", "/", nil)
	var myCtx BeforeRequestContext

	// Call BeforeRequest
	finalCtx, _ := test.BeforeRequest(myCtx, myReq)

	// Should construct URL with cloud domain from environment
	expectedUrlString := "https://staging-workspace-staging-org.cribl-staging.cloud/api/v1/"
	if finalCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook with cloud domain returned %s, expected %s", finalCtx.URL.String(), expectedUrlString)
	}

	// Should have bearer token from environment
	expectedHeaderString := "map[Authorization:[Bearer test-bearer-token]]"
	if fmt.Sprintf("%+v", finalCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook finalCtx.Header returned %+v, expected %+v", fmt.Sprintf("%+v", finalCtx.Header), expectedHeaderString)
	}

	// Clean up environment variables
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestProviderConfigTakesPrecedenceOverEnvironment(t *testing.T) {
	// Set environment variables that should be overridden by provider config
	os.Setenv("CRIBL_CLIENT_ID", "env-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "env-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "env-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "env-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl.cloud")
	os.Setenv("CRIBL_BEARER_TOKEN", "provider-bearer-token")

	// Create hook and initialize
	test := NewCriblTerraformHook()
	test.SDKInit("initial-url", nil)

	// Simulate provider configuration that should override environment
	providerSecurity := shared.Security{
		OrganizationID: StringPtr("provider-org"),
		WorkspaceID:    StringPtr("provider-workspace"),
		CloudDomain:    StringPtr("cribl-playground.cloud"),
	}

	// Create security source that returns provider config
	securitySource := func(ctx context.Context) (interface{}, error) {
		return providerSecurity, nil
	}

	// Create test request context with security source
	myCtx := BeforeRequestContext{
		HookContext: HookContext{
			Context:        context.Background(),
			SecuritySource: securitySource,
		},
	}

	// Create a test request
	myReq, _ := http.NewRequest("GET", "/", nil)

	// Call BeforeRequest - should use provider config over environment
	finalCtx, _ := test.BeforeRequest(myCtx, myReq)

	// Should construct URL using provider config, NOT environment variables
	expectedUrlString := "https://provider-workspace-provider-org.cribl-playground.cloud/api/v1/"
	if finalCtx.URL.String() != expectedUrlString {
		t.Errorf("Provider config precedence failed. Got %s, expected %s", finalCtx.URL.String(), expectedUrlString)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

// Helper function to create string pointers
func StringPtr(s string) *string {
	return &s
}

func TestTerraformBeforeRequest(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "Paradise City")

	var myClient HTTPClient
	var ctx BeforeRequestContext
	myUrl := "foobar"
	myReq, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	returnedCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	expectedHeaderString := "map[Authorization:[Bearer Paradise City]]"
	expectedUrlString := "https://punk'n-biz.cribl.cloud/api/v1/"

	if returnedCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.Method, "GET")
	}
	if returnedCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", returnedCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", returnedCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.orgID, "biz")
	}
}

func TestTerraformBeforeRequestMultiUse(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "Paradise City")

	var myClient HTTPClient
	var ctx BeforeRequestContext
	myUrl := "foobar"
	myReq, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	returnedCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	expectedHeaderString := "map[Authorization:[Bearer Paradise City]]"
	expectedUrlString := "https://punk'n-biz.cribl.cloud/api/v1/"

	if returnedCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.Method, "GET")
	}
	if returnedCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook returnedCtx.Url returned %s, expected %s", returnedCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", returnedCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", returnedCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.orgID, "biz")
	}

	finalCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	if finalCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", finalCtx.Method, "GET")
	}
	if finalCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook after retry finalCtx.URL returned %s, expected %s", finalCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", finalCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", finalCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", test.orgID, "biz")
	}
}

func TestTerraformBeforeRequestWithSecuritySource(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "Paradise City")

	var myClient HTTPClient
	var ctx BeforeRequestContext
	myUrl := "foobar"
	myReq, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx.HookContext.SecuritySource = func(context.Context) (interface{}, error) { var foo interface{}; return foo, nil }

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	returnedCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	expectedHeaderString := "map[Authorization:[Bearer Paradise City]]"
	expectedUrlString := "https://punk'n-biz.cribl.cloud/api/v1/"

	if returnedCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.Method, "GET")
	}
	if returnedCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", returnedCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", returnedCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.orgID, "biz")
	}
}

func TestTerraformBeforeRequestWithSecuritySourceMultiUse(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "Paradise City")

	var myClient HTTPClient
	var ctx BeforeRequestContext
	myUrl := "foobar"
	myReq, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx.HookContext.SecuritySource = func(context.Context) (interface{}, error) { var foo interface{}; return foo, nil }

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	returnedCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	expectedHeaderString := "map[Authorization:[Bearer Paradise City]]"
	expectedUrlString := "https://punk'n-biz.cribl.cloud/api/v1/"

	if returnedCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.Method, "GET")
	}
	if returnedCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", returnedCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", returnedCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", returnedCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook returnedCtx.Method returned %s, expected %s", test.orgID, "biz")
	}

	finalCtx, err := test.BeforeRequest(ctx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	if finalCtx.Method != "GET" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", finalCtx.Method, "GET")
	}
	if finalCtx.URL.String() != expectedUrlString {
		t.Errorf("*CriblTerraformHook after retry finalCtx.URL returned %s, expected %s", finalCtx.URL.String(), expectedUrlString)
	}
	if fmt.Sprintf("%+v", finalCtx.Header) != expectedHeaderString {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %+v, expected %+v", fmt.Sprintf("%+v", finalCtx.Header), expectedHeaderString)
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", test.workspaceID, "punk'n")
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook finalCtx.Method returned %s, expected %s", test.orgID, "biz")
	}
}

func TestTerraformBeforeRequestWithoutBearerToken(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth/token" {
			res.Write([]byte(`{"access_token":"my-access-token","expires_in":300}`))
			res.WriteHeader(200)
		}
	}))
	defer func() { testServer.Close() }()

	var beforeCtx BeforeRequestContext
	myClient := MockHTTPClient{}
	myUrl := testServer.URL
	myReq, err := http.NewRequest("POST", myUrl, nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx, cancel := context.WithTimeout(myReq.Context(), 10*time.Millisecond)
	defer cancel()
	beforeCtx.Context = ctx

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	outReq, err := test.BeforeRequest(beforeCtx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	if fmt.Sprintf("%s", outReq.Header) != "map[Authorization:[Bearer my-access-token]]" {
		t.Errorf("*getBearerToken output.Token returned %s, expected %s", fmt.Sprintf("%s", outReq.Header), "map[Authorization:[Bearer my-access-token]]")
	}
}

func TestTerraformBeforeRequestWithoutBearerTokenMultiUse(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth/token" {
			res.Write([]byte(`{"access_token":"my-access-token","expires_in":300}`))
			res.WriteHeader(200)
		}
	}))
	defer func() { testServer.Close() }()

	var beforeCtx BeforeRequestContext
	myClient := MockHTTPClient{}
	myUrl := testServer.URL
	myReq, err := http.NewRequest("POST", myUrl, nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx, cancel := context.WithTimeout(myReq.Context(), 10*time.Millisecond)
	defer cancel()
	beforeCtx.Context = ctx

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)
	outReq, err := test.BeforeRequest(beforeCtx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	if fmt.Sprintf("%s", outReq.Header) != "map[Authorization:[Bearer my-access-token]]" {
		t.Errorf("*getBearerToken output.Token returned %s, expected %s", fmt.Sprintf("%s", outReq.Header), "map[Authorization:[Bearer my-access-token]]")
	}

	finalReq, err := test.BeforeRequest(beforeCtx, myReq)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	if fmt.Sprintf("%s", finalReq.Header) != "map[Authorization:[Bearer my-access-token]]" {
		t.Errorf("*getBearerToken output.Token returned %s, expected %s", fmt.Sprintf("%s", finalReq.Header), "map[Authorization:[Bearer my-access-token]]")
	}
}

func TestTerraformGetBearerToken(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")
	os.Setenv("CRIBL_BEARER_TOKEN", "Paradise City")

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth/token" {
			res.Write([]byte(`{"access_token":"my-access-token","expires_in":300}`))
			res.WriteHeader(200)
		}
	}))
	defer func() { testServer.Close() }()

	myClient := MockHTTPClient{}
	myUrl := testServer.URL
	myReq, err := http.NewRequest("POST", myUrl, nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx, cancel := context.WithTimeout(myReq.Context(), 10*time.Millisecond)
	defer cancel()

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)

	output, err := test.getBearerToken(ctx, "my-client-id", "my-client-string", "my-default-audience")
	if err != nil {
		t.Errorf("Error generating bearer token for testing: %s", err)
	}

	futureTime := time.Now().Add(300 * time.Second)
	if output.Token != "my-access-token" {
		t.Errorf("*getBearerToken output.Token returned %s, expected %s", output.Token, "my-access-token")
	}
	if output.ExpiresAt.Unix() != futureTime.Unix() {
		t.Errorf("*getBearerToken output.ExpiresAt returned %s, expected %s", output.ExpiresAt, futureTime)
	}
}

func TestTerraformAfterError(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth/token" {
			res.Write([]byte(`{"access_token":"my-access-token","expires_in":300}`))
			res.WriteHeader(200)
		}
	}))
	defer func() { testServer.Close() }()

	var hookCtx HookContext
	localCtx := context.TODO()
	ctx, cancel := context.WithTimeout(localCtx, 10*time.Millisecond)
	defer cancel()

	hookCtx.BaseURL = testServer.URL
	hookCtx.Context = ctx

	afterCtx := AfterErrorContext{HookContext: hookCtx}
	myClient := MockHTTPClient{}
	myUrl := testServer.URL
	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)

	myResp := http.Response{StatusCode: http.StatusUnauthorized}
	_, err := test.AfterError(afterCtx, &myResp, errors.New("Testing error for AfterError"))
	if err == nil {
		t.Errorf("Expected error was not returned from test.AfterError")
	}

	expectedErrorString := "authentication handled by custom hook"
	if err.Error() != expectedErrorString {
		t.Errorf("test.AfterError returned unexpected error, got '%s', expected '%s'", err.Error(), expectedErrorString)
	}
}

func TestTerraformAfterErrorMultiUse(t *testing.T) {
	os.Setenv("CRIBL_CLIENT_ID", "foo")
	os.Setenv("CRIBL_CLIENT_SECRET", "bar")
	os.Setenv("CRIBL_ORGANIZATION_ID", "biz")
	os.Setenv("CRIBL_WORKSPACE_ID", "punk'n")

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth/token" {
			res.Write([]byte(`{"access_token":"my-access-token","expires_in":300}`))
			res.WriteHeader(200)
		}
	}))
	defer func() { testServer.Close() }()

	var hookCtx HookContext
	localCtx := context.TODO()
	ctx, cancel := context.WithTimeout(localCtx, 10*time.Millisecond)
	defer cancel()

	hookCtx.BaseURL = testServer.URL
	hookCtx.Context = ctx

	afterCtx := AfterErrorContext{HookContext: hookCtx}
	myClient := MockHTTPClient{}
	myUrl := testServer.URL
	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(myUrl, myClient)

	myResp := http.Response{StatusCode: http.StatusUnauthorized}
	_, err := test.AfterError(afterCtx, &myResp, errors.New("Testing error for AfterError"))
	if err == nil {
		t.Errorf("Expected error was not returned from test.AfterError")
	}

	expectedErrorString := "authentication handled by custom hook"
	if err.Error() != expectedErrorString {
		t.Errorf("test.AfterError returned unexpected error, got '%s', expected '%s'", err.Error(), expectedErrorString)
	}
}

func TestGovDomainOAuth2(t *testing.T) {
	// Test that gov domains use Okta OAuth2 with form-encoded data
	receivedContentType := ""
	receivedBody := ""
	requestedURL := ""

	// Create custom HTTP client that intercepts OAuth2 requests
	customClient := &MockHTTPClient{}
	customClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		// Capture OAuth2 request details
		if strings.Contains(req.URL.Path, "/oauth2/") && strings.Contains(req.URL.Path, "/v1/token") {
			requestedURL = req.URL.String()
			receivedContentType = req.Header.Get("Content-Type")
			bodyBytes, _ := io.ReadAll(req.Body)
			receivedBody = string(bodyBytes)

			// Return mock token response for OAuth2 requests
			response := &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"gov-access-token","expires_in":3600,"token_type":"Bearer"}`)),
			}
			response.Header.Set("Content-Type", "application/json")
			return response, nil
		}

		// For other requests, return empty response
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
	}

	os.Setenv("CRIBL_CLIENT_ID", "test-client-id")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-client-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-gov-staging.cloud")
	os.Setenv("CRIBL_OKTA_DOMAIN", "criblgov-stg.okta.com")
	os.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "ausfridm9cpg2Y5nW4h7")
	os.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "ausfridm9cpg2Y5nW4h7")
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	var beforeCtx BeforeRequestContext

	// Use gov domain workspace URL
	govWorkspaceURL := "https://test-workspace-test-org.cribl-gov-staging.cloud"

	myReq, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	beforeCtx.Context = ctx

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(govWorkspaceURL, customClient)

	outReq, err := test.BeforeRequest(beforeCtx, myReq)
	if err != nil {
		t.Errorf("Error in BeforeRequest for gov domain: %s", err)
		return
	}

	// Verify OAuth2 token was obtained
	expectedAuthHeader := "Bearer gov-access-token"
	actualAuthHeader := outReq.Header.Get("Authorization")
	if actualAuthHeader != expectedAuthHeader {
		t.Errorf("Expected Authorization header '%s', got '%s'", expectedAuthHeader, actualAuthHeader)
	}

	// Verify Okta URL was used
	expectedURLPattern := "criblgov-stg.okta.com/oauth2/ausfridm9cpg2Y5nW4h7/v1/token"
	if !strings.Contains(requestedURL, expectedURLPattern) {
		t.Errorf("Expected OAuth2 URL to contain '%s', got '%s'", expectedURLPattern, requestedURL)
	}

	// Verify form-encoded content type was used (not JSON)
	expectedContentType := "application/x-www-form-urlencoded"
	if receivedContentType != expectedContentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", expectedContentType, receivedContentType)
	}

	// Verify form data contains required fields
	if !strings.Contains(receivedBody, "grant_type=client_credentials") {
		t.Errorf("Form data should contain grant_type=client_credentials, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "client_id=test-client-id") {
		t.Errorf("Form data should contain client_id, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "client_secret=test-client-secret") {
		t.Errorf("Form data should contain client_secret, got: %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "audience=") {
		t.Errorf("Form data should contain audience parameter, got: %s", receivedBody)
	}

	// Clean up environment variables
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_OKTA_DOMAIN", "")
	os.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "")
	os.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestGovDomainOAuth2MissingAuthServerID(t *testing.T) {
	// Test that missing CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID returns an error
	os.Setenv("CRIBL_CLIENT_ID", "test-client-id")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-client-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-gov-staging.cloud")
	os.Setenv("CRIBL_OKTA_DOMAIN", "criblgov-stg.okta.com")
	os.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "")         // Not provided
	os.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "") // Not provided - should cause error
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	// Create custom HTTP client (shouldn't be called due to error)
	customClient := &MockHTTPClient{}
	customClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		t.Errorf("HTTP client should not be called when auth server ID is missing")
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
	}

	var beforeCtx BeforeRequestContext
	govWorkspaceURL := "https://test-workspace-test-org.cribl-gov-staging.cloud"

	myReq, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Errorf("Error generating http request for testing: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	beforeCtx.Context = ctx

	test := NewCriblTerraformHook()
	_, _ = test.SDKInit(govWorkspaceURL, customClient)

	// This should return an error due to missing auth server ID
	_, err = test.BeforeRequest(beforeCtx, myReq)
	if err == nil {
		t.Errorf("Expected error when CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID is missing, but got none")
		return
	}

	expectedError := "CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID environment variable is required for gov domains"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}

	// Clean up environment variables
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_OKTA_DOMAIN", "")
	os.Setenv("CRIBL_OKTA_AUTH_SERVER_ID", "")
	os.Setenv("CRIBL_OKTA_DEFAULT_AUTH_SERVER_ID", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestGatewayRoutingWithProviderConfig(t *testing.T) {
	// Set environment variables that should be overridden by provider config
	os.Setenv("CRIBL_CLIENT_ID", "env-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "env-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "env-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "env-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl.cloud")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-bearer-token")

	hook := NewCriblTerraformHook()
	hook.SDKInit("initial-url", nil)

	// Simulate provider configuration
	providerSecurity := shared.Security{
		OrganizationID: StringPtr("provider-org"),
		WorkspaceID:    StringPtr("provider-workspace"),
		CloudDomain:    StringPtr("cribl-playground.cloud"),
	}

	// Create security source that returns provider config
	securitySource := func(ctx context.Context) (interface{}, error) {
		return providerSecurity, nil
	}

	// Create test request context with security source
	myCtx := BeforeRequestContext{
		HookContext: HookContext{
			Context:        context.Background(),
			SecuritySource: securitySource,
		},
	}

	// Test gateway path with provider config
	gatewayReq, err := http.NewRequest("POST", "/v1/organizations/provider-org/workspaces", nil)
	if err != nil {
		t.Fatalf("Error creating test request: %v", err)
	}

	finalReq, err := hook.BeforeRequest(myCtx, gatewayReq)
	if err != nil {
		t.Fatalf("BeforeRequest failed: %v", err)
	}

	// Should route to gateway URL using provider cloud domain (no /api prefix)
	expectedURL := "https://gateway.cribl-playground.cloud/v1/organizations/provider-org/workspaces"
	if finalReq.URL.String() != expectedURL {
		t.Errorf("Gateway routing with provider config failed. Got %q, expected %q", finalReq.URL.String(), expectedURL)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestGatewayHostOverride(t *testing.T) {
	// Test that hardcoded gateway.cribl.cloud gets overridden with correct domain
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "cribl-playground.cloud")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-bearer-token")

	hook := NewCriblTerraformHook()
	hook.SDKInit("initial-url", nil)

	// Simulate a request that the SDK would send to hardcoded gateway.cribl.cloud
	hardcodedReq, err := http.NewRequest("POST", "https://gateway.cribl.cloud/v1/organizations/test-org/workspaces", nil)
	if err != nil {
		t.Fatalf("Error creating test request: %v", err)
	}

	var ctx BeforeRequestContext
	finalReq, err := hook.BeforeRequest(ctx, hardcodedReq)
	if err != nil {
		t.Fatalf("BeforeRequest failed: %v", err)
	}

	// Should override the host to use the correct domain
	expectedHost := "gateway.cribl-playground.cloud"
	if finalReq.URL.Host != expectedHost {
		t.Errorf("Gateway URL.Host override failed. Got %q, expected %q", finalReq.URL.Host, expectedHost)
	}

	// Should also set the explicit Host field
	if finalReq.Host != expectedHost {
		t.Errorf("Gateway Host field override failed. Got %q, expected %q", finalReq.Host, expectedHost)
	}

	// Path should remain the same (no /api prefix needed for gateway)
	expectedPath := "/v1/organizations/test-org/workspaces"
	if finalReq.URL.Path != expectedPath {
		t.Errorf("Path should remain unchanged for gateway. Got %q, expected %q", finalReq.URL.Path, expectedPath)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestOnPremSDKInit(t *testing.T) {
	// Test that SDKInit handles on-prem configuration
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")

	hook := NewCriblTerraformHook()
	url, _ := hook.SDKInit("should-be-overridden", nil)

	expectedURL := "http://localhost:9000"
	if url != expectedURL {
		t.Errorf("SDKInit for on-prem returned %q, expected %q", url, expectedURL)
	}
	if hook.baseURL != expectedURL {
		t.Errorf("Hook baseURL %q, expected %q", hook.baseURL, expectedURL)
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
}

func TestOnPremAuthenticationWithBearerToken(t *testing.T) {
	// Test on-prem authentication using direct bearer token
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_BEARER_TOKEN", "my-onprem-token")

	hook := NewCriblTerraformHook()
	hook.SDKInit("http://localhost:9000", nil)

	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	ctx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	resultReq, err := hook.BeforeRequest(ctx, req)
	if err != nil {
		t.Fatalf("BeforeRequest failed: %v", err)
	}

	// Verify token was set
	expectedAuth := "Bearer my-onprem-token"
	if resultReq.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization %q, got %q", expectedAuth, resultReq.Header.Get("Authorization"))
	}

	// Verify URL was set correctly
	expectedURL := "http://localhost:9000/api/v1/test"
	if resultReq.URL.String() != expectedURL {
		t.Errorf("Expected URL %q, got %q", expectedURL, resultReq.URL.String())
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestOnPremAuthenticationWithUsernamePassword(t *testing.T) {
	// Test on-prem authentication using username/password
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_ONPREM_USERNAME", "admin")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "admin")
	os.Setenv("CRIBL_BEARER_TOKEN", "") // Not using direct token

	// Create a test server that returns a successful auth response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" && r.Method == "POST" {
			// Verify request body
			bodyBytes, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(bodyBytes), "username") || !strings.Contains(string(bodyBytes), "password") {
				t.Errorf("Request body should contain username and password")
			}
			// Return successful auth response
			w.Write([]byte(`{"token": "test-token-from-login", "forcePasswordChange": false}`))
		}
	}))
	defer testServer.Close()

	hook := NewCriblTerraformHook()
	// Override the server URL to use test server
	os.Setenv("CRIBL_ONPREM_SERVER_URL", testServer.URL)

	// Use http.DefaultClient to make actual HTTP requests to the test server
	hook.SDKInit(testServer.URL, &http.Client{})

	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	ctx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	resultReq, err := hook.BeforeRequest(ctx, req)
	if err != nil {
		t.Fatalf("BeforeRequest failed: %v", err)
	}

	// Verify token was set
	expectedAuth := "Bearer test-token-from-login"
	if resultReq.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization %q, got %q", expectedAuth, resultReq.Header.Get("Authorization"))
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_ONPREM_USERNAME", "")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestOnPremTokenCaching(t *testing.T) {
	// Test that tokens are cached and reused
	requestCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" && r.Method == "POST" {
			requestCount++
			w.Write([]byte(`{"token": "cached-token", "forcePasswordChange": false}`))
		}
	}))
	defer testServer.Close()

	// Set environment variables with the test server URL BEFORE creating hook
	os.Setenv("CRIBL_ONPREM_SERVER_URL", testServer.URL)
	os.Setenv("CRIBL_ONPREM_USERNAME", "admin")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "admin")
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	// Create hook and initialize with test server URL
	hook := NewCriblTerraformHook()
	hook.SDKInit(testServer.URL, &http.Client{})

	req1, _ := http.NewRequest("GET", "/api/v1/test1", nil)
	req2, _ := http.NewRequest("GET", "/api/v1/test2", nil)
	ctx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	// First request - should fetch token
	_, err := hook.BeforeRequest(ctx, req1)
	if err != nil {
		t.Fatalf("First BeforeRequest failed: %v", err)
	}

	// Second request - should use cached token
	_, err = hook.BeforeRequest(ctx, req2)
	if err != nil {
		t.Fatalf("Second BeforeRequest failed: %v", err)
	}

	// Should only have called the auth endpoint once
	if requestCount != 1 {
		t.Errorf("Expected auth endpoint to be called 1 time, got %d", requestCount)
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_ONPREM_USERNAME", "")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "")
}

func TestOnPremAuthenticationErrors(t *testing.T) {
	// Test error handling when on-prem auth fails
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_ONPREM_USERNAME", "admin")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "wrong-password")
	os.Setenv("CRIBL_BEARER_TOKEN", "")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" && r.Method == "POST" {
			// Return 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid credentials"}`))
		}
	}))
	defer testServer.Close()

	os.Setenv("CRIBL_ONPREM_SERVER_URL", testServer.URL)

	hook := NewCriblTerraformHook()
	hook.SDKInit(testServer.URL, &http.Client{})

	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	ctx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	_, err := hook.BeforeRequest(ctx, req)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "failed to get on-prem bearer token") {
		t.Errorf("Expected error about failing to get token, got: %v", err)
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_ONPREM_USERNAME", "")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "")
}

func TestOnPremURLRouting(t *testing.T) {
	// Test that on-prem requests are routed correctly
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	hook := NewCriblTerraformHook()
	hook.SDKInit("http://localhost:9000", nil)

	testCases := []struct {
		path     string
		expected string
	}{
		{"/api/v1/sources", "http://localhost:9000/api/v1/sources"},
		{"/api/v1/destinations", "http://localhost:9000/api/v1/destinations"},
		{"/api/v1/system/health", "http://localhost:9000/api/v1/system/health"},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", tc.path, nil)
		ctx := BeforeRequestContext{
			HookContext: HookContext{
				Context: context.Background(),
			},
		}

		resultReq, err := hook.BeforeRequest(ctx, req)
		if err != nil {
			t.Errorf("BeforeRequest failed for %s: %v", tc.path, err)
			continue
		}

		if resultReq.URL.String() != tc.expected {
			t.Errorf("For path %q, expected URL %q, got %q", tc.path, tc.expected, resultReq.URL.String())
		}
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestOnPremMissingCredentials(t *testing.T) {
	// Test error when on-prem URL is set but no credentials provided
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
	os.Setenv("CRIBL_ONPREM_USERNAME", "")
	os.Setenv("CRIBL_ONPREM_PASSWORD", "")

	hook := NewCriblTerraformHook()
	hook.SDKInit("http://localhost:9000", nil)

	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	ctx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	_, err := hook.BeforeRequest(ctx, req)
	if err == nil {
		t.Fatalf("Expected error when no credentials provided")
	}

	if !strings.Contains(err.Error(), "authentication requires either environment variables or a cribl config file. Please refer to the provider docs") {
		t.Errorf("Expected error about missing credentials, got: %v", err)
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
}

func TestOnPremAllowedEndpoints(t *testing.T) {
	// Test that allowed endpoints work fine for on-prem deployments
	allowedPaths := []struct {
		path        string
		description string
	}{
		{"/sources", "Sources endpoint"},
		{"/destinations", "Destinations endpoint"},
		{"/routes", "Routes endpoint"},
		{"/pipelines", "Pipelines endpoint"},
		{"/packs", "Packs endpoint"},
		{"/system/inputs", "System inputs"},
		{"/system/outputs", "System outputs"},
		{"/groups/default", "Worker groups"},
	}

	// Set up on-prem configuration
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "http://localhost:9000")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-token")

	hook := NewCriblTerraformHook()
	hook.SDKInit("http://localhost:9000", nil)

	for _, tc := range allowedPaths {
		req, _ := http.NewRequest("GET", tc.path, nil)
		ctx := BeforeRequestContext{
			HookContext: HookContext{
				Context: context.Background(),
			},
		}

		resultReq, err := hook.BeforeRequest(ctx, req)
		if err != nil {
			t.Errorf("Unexpected error for allowed endpoint %s (%s): %v", tc.path, tc.description, err)
			continue
		}

		// Verify that the URL was set correctly
		if !strings.HasPrefix(resultReq.URL.String(), "http://localhost:9000") {
			t.Errorf("Expected URL to start with http://localhost:9000, got %s", resultReq.URL.String())
		}
	}

	// Clean up
	os.Setenv("CRIBL_ONPREM_SERVER_URL", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestOnPremWithProviderConfig(t *testing.T) {
	// Test that on-prem configuration works when server_url is provided via provider config

	// Set environment variables
	os.Setenv("CRIBL_BEARER_TOKEN", "provider-provided-token")

	// Simulate provider configuration that provides on-prem server URL
	// The provider config's server_url gets passed as baseURL to SDKInit
	hook := NewCriblTerraformHook()
	hook.SDKInit("http://localhost:9000", &MockHTTPClient{})

	// Create test request context
	myCtx := BeforeRequestContext{
		HookContext: HookContext{
			Context: context.Background(),
		},
	}

	req, _ := http.NewRequest("GET", "/api/v1/sources", nil)
	resultReq, err := hook.BeforeRequest(myCtx, req)

	if err != nil {
		t.Fatalf("BeforeRequest failed: %v", err)
	}

	// Verify token was set
	expectedAuth := "Bearer provider-provided-token"
	if resultReq.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization %q, got %q", expectedAuth, resultReq.Header.Get("Authorization"))
	}

	// Clean up
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestProviderConfigWithNilCredentialsFile(t *testing.T) {
	// Test scenario: Provider config provides client_id and client_secret,
	// but no credentials file exists (config is nil)
	// This tests the fix for the nil pointer dereference panic

	// Clear environment variables to ensure we're testing provider config only
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-bearer-token") // Use bearer token to avoid OAuth complexity
	os.Setenv("HOME", "/nonexistent")                    // Ensure no credentials file is found

	// Create hook and initialize with a template URL (not concrete)
	// This allows provider variables to override it
	hook := NewCriblTerraformHook()
	hook.SDKInit("https://{workspaceId}-{organizationId}.{cloudDomain}", &http.Client{})

	// Simulate provider configuration with client_id and client_secret
	providerSecurity := shared.Security{
		ClientOauth: &shared.SchemeClientOauth{
			ClientID:     "provider-client-id",
			ClientSecret: "provider-client-secret",
		},
		OrganizationID: StringPtr("provider-org"),
		WorkspaceID:    StringPtr("provider-workspace"),
		CloudDomain:    StringPtr("cribl-playground.cloud"),
	}

	// Create security source that returns provider config
	securitySource := func(ctx context.Context) (interface{}, error) {
		return providerSecurity, nil
	}

	// Create test request context with security source
	myCtx := BeforeRequestContext{
		HookContext: HookContext{
			Context:        context.Background(),
			SecuritySource: securitySource,
		},
	}

	// Create a test request
	myReq, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Error creating test request: %v", err)
	}

	// Call BeforeRequest - should not panic even though config is nil
	// It should create a minimal config object with provider credentials
	finalReq, err := hook.BeforeRequest(myCtx, myReq)
	if err != nil {
		t.Fatalf("BeforeRequest should not fail with provider config: %v", err)
	}

	// Verify bearer token was used
	expectedAuth := "Bearer test-bearer-token"
	if finalReq.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization %q, got %q", expectedAuth, finalReq.Header.Get("Authorization"))
	}

	// The main goal: verify that BeforeRequest didn't panic when config is nil
	// and that it successfully extracted provider variables
	// Verify that hook stored the provider org and workspace IDs
	if hook.orgID != "provider-org" {
		t.Errorf("Expected orgID 'provider-org', got %q", hook.orgID)
	}
	if hook.workspaceID != "provider-workspace" {
		t.Errorf("Expected workspaceID 'provider-workspace', got %q", hook.workspaceID)
	}

	// Verify that the request was processed successfully (no panic occurred)
	if finalReq == nil {
		t.Error("finalReq should not be nil - indicates a panic occurred")
	}

	// Clean up
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestProviderConfigWithPointerType(t *testing.T) {
	// Test that provider config works when Security is passed as a pointer
	// This tests the type assertion fix for pointer types

	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
	os.Setenv("CRIBL_CLOUD_DOMAIN", "")
	os.Setenv("CRIBL_BEARER_TOKEN", "test-bearer-token")
	os.Setenv("HOME", "/nonexistent")

	hook := NewCriblTerraformHook()
	hook.SDKInit("initial-url", nil)

	// Simulate provider configuration as a pointer (as it might be passed)
	providerSecurity := &shared.Security{
		ClientOauth: &shared.SchemeClientOauth{
			ClientID:     "pointer-client-id",
			ClientSecret: "pointer-client-secret",
		},
		OrganizationID: StringPtr("pointer-org"),
		WorkspaceID:    StringPtr("pointer-workspace"),
		CloudDomain:    StringPtr("cribl-staging.cloud"),
	}

	// Create security source that returns pointer type
	securitySource := func(ctx context.Context) (interface{}, error) {
		return providerSecurity, nil
	}

	myCtx := BeforeRequestContext{
		HookContext: HookContext{
			Context:        context.Background(),
			SecuritySource: securitySource,
		},
	}

	myReq, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Error creating test request: %v", err)
	}

	// Should not panic and should extract credentials correctly
	finalReq, err := hook.BeforeRequest(myCtx, myReq)
	if err != nil {
		t.Fatalf("BeforeRequest should not fail with pointer type: %v", err)
	}

	// Verify bearer token was used (since we set it)
	expectedAuth := "Bearer test-bearer-token"
	if finalReq.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization %q, got %q", expectedAuth, finalReq.Header.Get("Authorization"))
	}

	// Verify org and workspace IDs were extracted from pointer
	if hook.orgID != "pointer-org" {
		t.Errorf("Expected orgID 'pointer-org', got %q", hook.orgID)
	}
	if hook.workspaceID != "pointer-workspace" {
		t.Errorf("Expected workspaceID 'pointer-workspace', got %q", hook.workspaceID)
	}

	// Clean up
	os.Setenv("CRIBL_BEARER_TOKEN", "")
}

func TestDoRequestWithRetrySuccess(t *testing.T) {
	// Test successful request without retries
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"success-token","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success but got error: %v", err)
	}

	if string(body) != `{"access_token":"success-token","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"success-token","expires_in":300}`, string(body))
	}

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt, got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryNetworkError(t *testing.T) {
	// Test retry on network errors
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			if attemptCount < 2 {
				return nil, errors.New("network error")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"retry-success","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success after retry but got error: %v", err)
	}

	if string(body) != `{"access_token":"retry-success","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"retry-success","expires_in":300}`, string(body))
	}

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryUnexpectedEOF(t *testing.T) {
	// Test retry on unexpected EOF (body read error)
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			// First attempt: return response with body that will fail to read
			if attemptCount == 1 {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&failingReader{err: io.ErrUnexpectedEOF}),
				}, nil
			}
			// Second attempt: return successful response
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"eof-retry-success","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success after EOF retry but got error: %v", err)
	}

	if string(body) != `{"access_token":"eof-retry-success","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"eof-retry-success","expires_in":300}`, string(body))
	}

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetry429(t *testing.T) {
	// Test retry on 429 Too Many Requests
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			if attemptCount < 3 {
				return &http.Response{
					StatusCode: http.StatusTooManyRequests,
					Body:       io.NopCloser(strings.NewReader(`{"error":"rate limited"}`)),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"429-retry-success","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success after 429 retry but got error: %v", err)
	}

	if string(body) != `{"access_token":"429-retry-success","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"429-retry-success","expires_in":300}`, string(body))
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryMaxRetriesExhausted(t *testing.T) {
	// Test that max retries are exhausted and error is returned
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			return nil, errors.New("persistent network error")
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	_, err := hook.doRequestWithRetry(req)
	if err == nil {
		t.Fatalf("Expected error after max retries but got success")
	}

	if !strings.Contains(err.Error(), "failed to make request with retry") {
		t.Errorf("Expected error about failed request, got: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts (max retries), got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryNon200Status(t *testing.T) {
	// Test that non-200 status codes (except 429) don't retry
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"error":"bad request"}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	_, err := hook.doRequestWithRetry(req)
	if err == nil {
		t.Fatalf("Expected error for non-200 status but got success")
	}

	if !strings.Contains(err.Error(), "status=400") {
		t.Errorf("Expected error with status 400, got: %v", err)
	}

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt (no retry for 400), got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryContextCancellation(t *testing.T) {
	// Test that context cancellation is respected
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			if attemptCount == 1 {
				return nil, errors.New("network error")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"token","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(`{"grant_type":"client_credentials"}`))

	// Cancel context during retry backoff
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err := hook.doRequestWithRetry(req)
	if err == nil {
		t.Fatalf("Expected error due to context cancellation but got success")
	}

	if !strings.Contains(err.Error(), "request cancelled") {
		t.Errorf("Expected cancellation error, got: %v", err)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryRequestBodyReuse(t *testing.T) {
	// Test that request body is properly reused on retries
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	requestBodies := []string{}
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			// Read request body to verify it's sent on each retry
			if req.Body != nil {
				bodyBytes, _ := io.ReadAll(req.Body)
				requestBodies = append(requestBodies, string(bodyBytes))
			}

			if attemptCount < 2 {
				return nil, errors.New("network error")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"body-reuse-success","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	requestBody := `{"grant_type":"client_credentials","client_id":"test"}`
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://login.cribl.cloud/oauth/token", strings.NewReader(requestBody))

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success after retry but got error: %v", err)
	}

	if string(body) != `{"access_token":"body-reuse-success","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"body-reuse-success","expires_in":300}`, string(body))
	}

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}

	// Verify request body was sent on both attempts
	if len(requestBodies) != 2 {
		t.Errorf("Expected 2 request bodies, got %d", len(requestBodies))
	}
	for i, body := range requestBodies {
		if body != requestBody {
			t.Errorf("Attempt %d: Expected request body %q, got %q", i+1, requestBody, body)
		}
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

func TestDoRequestWithRetryNoBody(t *testing.T) {
	// Test request with no body
	os.Setenv("CRIBL_CLIENT_ID", "test-client")
	os.Setenv("CRIBL_CLIENT_SECRET", "test-secret")
	os.Setenv("CRIBL_ORGANIZATION_ID", "test-org")
	os.Setenv("CRIBL_WORKSPACE_ID", "test-workspace")

	attemptCount := 0
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"no-body-success","expires_in":300}`)),
			}, nil
		},
	}

	hook := NewCriblTerraformHook()
	hook.client = mockClient
	hook.SDKInit("https://test-workspace-test-org.cribl.cloud", mockClient)

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://login.cribl.cloud/oauth/token", nil)

	body, err := hook.doRequestWithRetry(req)
	if err != nil {
		t.Fatalf("Expected success but got error: %v", err)
	}

	if string(body) != `{"access_token":"no-body-success","expires_in":300}` {
		t.Errorf("Expected body %q, got %q", `{"access_token":"no-body-success","expires_in":300}`, string(body))
	}

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt, got %d", attemptCount)
	}

	// Clean up
	os.Setenv("CRIBL_CLIENT_ID", "")
	os.Setenv("CRIBL_CLIENT_SECRET", "")
	os.Setenv("CRIBL_ORGANIZATION_ID", "")
	os.Setenv("CRIBL_WORKSPACE_ID", "")
}

// failingReader is a reader that always returns an error
type failingReader struct {
	err error
}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func (r *failingReader) Close() error {
	return nil
}
