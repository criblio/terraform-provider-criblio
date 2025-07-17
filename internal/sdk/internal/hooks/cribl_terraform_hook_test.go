package hooks

import (
	"context"
	"errors"
	"fmt"
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
		t.Errorf(fmt.Sprintf("NewCriblTerraformHook returned incorrect struct, expected '%+v' got '%+v'", &testData, output))
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

	if url != myUrl {
		t.Errorf("creds hook init returned %s, expected %s", url, myUrl)
	}
	if test.client != myClient {
		t.Errorf("creds hook init test.client %+v, expected %+v", test.client, myClient)
	}
	if test.baseURL != myUrl {
		t.Errorf("creds hook init test.baseURLd %s, expected %s", url, myUrl)
	}
	if test.orgID != "biz" {
		t.Errorf("*CriblTerraformHook.orgID returned %s, expected %s", test.orgID, "biz")
	}
	if test.workspaceID != "punk'n" {
		t.Errorf("*CriblTerraformHook.workspaceID returned %s, expected %s", test.orgID, "biz")
	}
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
	expectedUrlString := "foobar/organizations/biz/workspaces/punk'n/app/api/v1/"

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
	expectedUrlString := "foobar/organizations/biz/workspaces/punk'n/app/api/v1/"

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
