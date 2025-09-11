package hooks

import (
	"context"
	"errors"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/models/shared"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
