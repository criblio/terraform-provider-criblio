# Authentication in CLI Tool

## Overview

**Authentication is handled automatically** by the `CriblTerraformHook` which is automatically registered when you initialize the SDK. You do NOT need to reimplement authentication logic.

## How It Works

The `CriblTerraformHook` (located in `internal/sdk/internal/hooks/cribl_terraform_hook.go`) is automatically registered in the SDK initialization via `internal/sdk/internal/hooks/registration.go`.

When you create an SDK client:
```go
import "github.com/criblio/terraform-provider-criblio/internal/sdk"

client := sdk.New()
```

The hook automatically:
1. Reads environment variables (`CRIBL_BEARER_TOKEN`, `CRIBL_CLIENT_ID`, etc.)
2. Reads credentials file (`~/.cribl/credentials`)
3. Handles OAuth token acquisition and caching
4. Handles on-prem authentication
5. Constructs correct URLs for workspace and gateway endpoints

## Supported Authentication Methods

### 1. Bearer Token (Simplest)
```bash
export CRIBL_BEARER_TOKEN="your-token"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_ORGANIZATION_ID="org"
```

### 2. OAuth (Client Credentials)
```bash
export CRIBL_CLIENT_ID="your-client-id"
export CRIBL_CLIENT_SECRET="your-client-secret"
export CRIBL_ORGANIZATION_ID="org"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_CLOUD_DOMAIN="cribl.cloud"  # optional
```

### 3. Credentials File
Create `~/.cribl/credentials`:
```ini
[default]
client_id = your-client-id
client_secret = your-client-secret
organization_id = org
workspace = main
cloud_domain = cribl.cloud  # optional
```

Then use:
```bash
export CRIBL_PROFILE="default"  # optional, uses "default" if not set
```

### 4. On-Prem Authentication
```bash
# Option 1: Bearer token
export CRIBL_ONPREM_SERVER_URL="http://localhost:9000"
export CRIBL_BEARER_TOKEN="your-token"

# Option 2: Username/password
export CRIBL_ONPREM_SERVER_URL="http://localhost:9000"
export CRIBL_ONPREM_USERNAME="admin"
export CRIBL_ONPREM_PASSWORD="admin"
```

Or via credentials file:
```ini
[onprem]
onprem_server_url = http://localhost:9000
onprem_username = admin
onprem_password = admin
```

## Implementation Example

```go
package main

import (
    "context"
    "github.com/criblio/terraform-provider-criblio/internal/sdk"
)

func main() {
    ctx := context.Background()
    
    // Create SDK client - authentication handled automatically by hook
    client := sdk.New()
    
    // Use the client - all requests are automatically authenticated
    response, err := client.System.Diag.GetDiagBundle(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use response...
}
```

## Precedence Order

The hook follows this precedence (highest to lowest):

1. **Environment Variables** - Highest priority
2. **Credentials File** (`~/.cribl/credentials`) - Lower priority
3. **Defaults** - Lowest priority (e.g., workspace="main", org="ian", domain="cribl.cloud")

## Token Caching

The hook automatically caches OAuth tokens:
- Tokens are cached per session key (client_id:client_secret)
- Tokens are refreshed when they expire (within 60 minutes of expiration)
- On-prem tokens are cached per server/username/password combination

## URL Construction

The hook automatically constructs URLs:
- **Cloud**: `https://{workspace}-{org}.{domain}`
- **On-Prem**: Uses the server URL directly
- **Gateway**: `https://gateway.{domain}` for management endpoints

## What You DON'T Need to Do

❌ **Don't** manually set Authorization headers  
❌ **Don't** implement OAuth token acquisition  
❌ **Don't** implement on-prem authentication  
❌ **Don't** manage token caching  
❌ **Don't** construct workspace URLs manually  

✅ **Just** create the SDK client and use it!

## References

- Hook implementation: `internal/sdk/internal/hooks/cribl_terraform_hook.go`
- Hook registration: `internal/sdk/internal/hooks/registration.go`
- Credentials file reader: `internal/sdk/internal/hooks/credentials.go`

