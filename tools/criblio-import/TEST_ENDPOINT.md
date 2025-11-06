# Testing the /system/diag/download Endpoint

Simple commands to test the endpoint and inspect the downloaded configuration bundle.

## Quick Test (Single Command)

### For Cribl.Cloud:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "https://main-org.cribl.cloud/api/v1/system/diag/download" \
  -o diag-bundle.tar.gz && \
  tar -tzf diag-bundle.tar.gz | head -20
```

### For On-Prem:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:9000/api/v1/system/diag/download" \
  -o diag-bundle.tar.gz && \
  tar -tzf diag-bundle.tar.gz | head -20
```

## Extract and Inspect

### Extract the archive:
```bash
mkdir -p diag-extracted
tar -xzf diag-bundle.tar.gz -C diag-extracted
```

### Find specific resource types:
```bash
# Sources
find diag-extracted -path "*/sources/*.yaml" -o -path "*/sources/*.yml"

# Destinations
find diag-extracted -path "*/destinations/*.yaml" -o -path "*/destinations/*.yml"

# Pipelines
find diag-extracted -path "*/pipelines/*.yaml" -o -path "*/pipelines/*.yml"

# Routes
find diag-extracted -path "*/routes/*.yaml" -o -path "*/routes/*.yml"

# Packs
find diag-extracted -path "*/packs/*.yaml" -o -path "*/packs/*.yml"
```

### View a single resource file:
```bash
# Get first source file
FIRST_SOURCE=$(find diag-extracted -path "*/sources/*.yaml" | head -1)
cat "$FIRST_SOURCE"

# Or with yq (if installed) for formatted output
yq eval . "$FIRST_SOURCE"
```

## Using the Test Scripts

### Full test script (recommended):
```bash
export CRIBL_BEARER_TOKEN="your-token"
export CRIBL_WORKSPACE_ID="main"
export CRIBL_ORGANIZATION_ID="org"

./test-download.sh
```

### Simple one-liner script:
```bash
CRIBL_BEARER_TOKEN="your-token" \
CRIBL_WORKSPACE_ID="main" \
CRIBL_ORGANIZATION_ID="org" \
  ./test-download-simple.sh
```

## Using OAuth Instead of Bearer Token

First get a bearer token via OAuth, then use it:

```bash
# Get token (example - adjust URL and credentials)
TOKEN=$(curl -s -X POST \
  "https://api.cribl.cloud/oauth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET&audience=https://api.cribl.cloud" \
  | jq -r '.access_token')

# Use token to download
curl -H "Authorization: Bearer $TOKEN" \
  "https://main-org.cribl.cloud/api/v1/system/diag/download" \
  -o diag-bundle.tar.gz
```

## Expected Archive Structure

The archive typically contains:
```
configs/
├── sources/
│   ├── source1.yaml
│   ├── source2.yaml
│   └── ...
├── destinations/
│   ├── dest1.yaml
│   └── ...
├── pipelines/
│   └── ...
├── routes/
│   └── ...
├── packs/
│   └── ...
├── groups/
│   └── ...
└── ... (other config types)
```

## Example YAML Structure

A typical source YAML file might look like:
```yaml
id: my-http-source
type: http
group: default
port: 8080
host: 0.0.0.0
maxBufferSize: 1024
tls:
  certPath: /path/to/cert
  rejectUnauthorized: true
```

## Troubleshooting

### 401 Unauthorized
- Verify your bearer token is valid
- Check token expiration
- Ensure workspace/organization IDs are correct

### Empty or missing files
- Check if you have any resources configured in Cribl
- Verify you're accessing the correct workspace

### Archive structure differs
- Different Cribl versions may have different structures
- Use `find diag-extracted -type f` to explore actual structure

