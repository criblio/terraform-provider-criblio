package client

// Example: How to use the SDK with automatic authentication
//
// This demonstrates how simple it is to use the SDK - authentication
// is handled automatically by CriblTerraformHook.

/*
import (
	"context"
	"fmt"
	"os"
	
	"github.com/criblio/terraform-provider-criblio/internal/sdk"
)

func ExampleUsage() {
	ctx := context.Background()
	
	// That's it! Authentication is handled automatically by the hook
	// The hook reads from:
	// 1. Environment variables (CRIBL_BEARER_TOKEN, CRIBL_CLIENT_ID, etc.)
	// 2. Credentials file (~/.cribl/credentials)
	// 3. Defaults (workspace="main", org="ian")
	client := sdk.New()
	
	// Download diag bundle - authentication is automatic
	response, err := client.System.Diag.GetDiagBundle(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	// Use the response (tar.gz bytes)
	fmt.Printf("Downloaded %d bytes\n", len(response.Bytes))
	
	// All other SDK calls are also automatically authenticated
	// No need to manually set headers or manage tokens!
}
*/

// Authentication Methods Supported Automatically:
//
// 1. Bearer Token:
//    export CRIBL_BEARER_TOKEN="token"
//    export CRIBL_WORKSPACE_ID="main"
//    export CRIBL_ORGANIZATION_ID="org"
//
// 2. OAuth:
//    export CRIBL_CLIENT_ID="id"
//    export CRIBL_CLIENT_SECRET="secret"
//    export CRIBL_ORGANIZATION_ID="org"
//    export CRIBL_WORKSPACE_ID="main"
//
// 3. Credentials File (~/.cribl/credentials):
//    [default]
//    client_id = id
//    client_secret = secret
//    organization_id = org
//    workspace = main
//
// 4. On-Prem:
//    export CRIBL_ONPREM_SERVER_URL="http://localhost:9000"
//    export CRIBL_ONPREM_USERNAME="admin"
//    export CRIBL_ONPREM_PASSWORD="admin"
//
// All of these are handled automatically - no code changes needed!

