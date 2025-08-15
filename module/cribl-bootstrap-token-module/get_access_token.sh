#!/bin/bash

# Get Bootstrap Token Script
# This script retrieves the auth bootstrap token from Cribl Cloud workspace

# Function to display usage
usage() {
    echo "Usage: $0 -i CLIENT_ID -s CLIENT_SECRET -o ORG [-w WORKSPACE] [-u WORKSPACE_URL] [-e ENVIRONMENT] [-g GROUP]"
    echo
    echo "Options:"
    echo "  -i CLIENT_ID       Cribl Cloud Client ID (required)"
    echo "  -s CLIENT_SECRET   Cribl Cloud Client Secret (required)"
    echo "  -o ORG            Organization name (required)"
    echo "  -w WORKSPACE      Workspace name (required if -u not provided)"
    echo "  -u WORKSPACE_URL  Full workspace URL (e.g., https://workspace-org.cribl.cloud)"
    echo "  -e ENVIRONMENT    Environment: production or staging (default: production)"
    echo "  -g GROUP          Group parameter (default: defaultHybrid)"
    echo "  -h                Display this help message"
    echo
    echo "Examples:"
    echo "  $0 -i myid -s mysecret -o myorg -w myworkspace"
    echo "  $0 -i myid -s mysecret -o myorg -u https://myworkspace-myorg.cribl.cloud"
    exit 1
}

# Default values
ENVIRONMENT="production"
GROUP="defaultHybrid"

# Parse command line arguments
while getopts "i:s:o:w:u:e:g:h" opt; do
    case $opt in
        i) CLIENT_ID="$OPTARG";;
        s) CLIENT_SECRET="$OPTARG";;
        o) ORG="$OPTARG";;
        w) WORKSPACE="$OPTARG";;
        u) WORKSPACE_URL="$OPTARG";;
        e) ENVIRONMENT="$OPTARG";;
        g) GROUP="$OPTARG";;
        h) usage;;
        \?) echo "Invalid option: -$OPTARG" >&2; usage;;
    esac
done

# Validate required parameters
if [ -z "$CLIENT_ID" ] || [ -z "$CLIENT_SECRET" ] || [ -z "$ORG" ]; then
    echo "Error: CLIENT_ID, CLIENT_SECRET, and ORG are required" >&2
    usage
fi

if [ -z "$WORKSPACE_URL" ] && [ -z "$WORKSPACE" ]; then
    echo "Error: Either WORKSPACE or WORKSPACE_URL must be provided" >&2
    usage
fi

# Set domain based on environment
if [ "$ENVIRONMENT" = "staging" ]; then
    DOMAIN="cribl-staging.cloud"
    TOKEN_URL="https://login.cribl-staging.cloud/oauth/token"
    AUDIENCE="https://api.cribl-staging.cloud"
else
    DOMAIN="cribl.cloud"
    TOKEN_URL="https://login.cribl.cloud/oauth/token"
    AUDIENCE="https://api.cribl.cloud"
fi

# Step 1: Get OAuth token
echo "Retrieving OAuth token..." >&2
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_URL" \
    -H "Content-Type: application/json" \
    -d "{
        \"grant_type\": \"client_credentials\",
        \"client_id\": \"$CLIENT_ID\",
        \"client_secret\": \"$CLIENT_SECRET\",
        \"audience\": \"$AUDIENCE\"
    }")

# Extract access token
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Error: Failed to get access token" >&2
    echo "Response: $TOKEN_RESPONSE" >&2
    exit 1
fi

echo "Successfully obtained access token" >&2

# Step 2: Construct workspace URL if not provided
if [ -z "$WORKSPACE_URL" ]; then
    WORKSPACE_URL="https://${WORKSPACE}-${ORG}.${DOMAIN}"
    echo "Constructed workspace URL: $WORKSPACE_URL" >&2
else
    echo "Using provided workspace URL: $WORKSPACE_URL" >&2
fi

# Step 3: Get bootstrap token
BOOTSTRAP_URL="${WORKSPACE_URL}/api/v1/system/instance/distributed"
echo "Fetching bootstrap token from: $BOOTSTRAP_URL" >&2

BOOTSTRAP_RESPONSE=$(curl -s -X GET "$BOOTSTRAP_URL" \
    -H "Accept: application/json, text/plain, */*" \
    -H "Accept-Language: en-US,en;q=0.9" \
    -H "Accept-Encoding: gzip, deflate, br" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Cache-Control: no-cache" \
    -H "Connection: keep-alive" \
    -H "Pragma: no-cache" \
    -H "Referer: ${WORKSPACE_URL}/?group=${GROUP}" \
    -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36" \
    -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
    -H "sec-ch-ua-mobile: ?0" \
    -H 'sec-ch-ua-platform: "macOS"' \
    -H "sec-fetch-dest: empty" \
    -H "sec-fetch-mode: cors" \
    -H "sec-fetch-site: same-origin")

# Extract bootstrap token
BOOTSTRAP_TOKEN=$(echo "$BOOTSTRAP_RESPONSE" | grep -o '"authToken":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$BOOTSTRAP_TOKEN" ]; then
    echo "Error: Failed to get bootstrap token" >&2
    echo "Response: $BOOTSTRAP_RESPONSE" >&2
    exit 1
fi

echo "Successfully retrieved bootstrap token" >&2
echo "$BOOTSTRAP_TOKEN"