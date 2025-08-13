#!/bin/bash

# Script to check if a lakehouse is ready
# Usage: ./check_lakehouse_status.sh <lakehouse_id>

set -e

LAKEHOUSE_ID="$1"
ORGANIZATION_ID="${CRIBL_ORG_ID:-beautiful-nguyen-y8y4azd}"
WORKSPACE_ID="${CRIBL_WORKSPACE_ID:-main}"
SERVER_URL="${CRIBL_SERVER_URL:-https://app.cribl-playground.cloud}"
TOKEN="${CRIBL_TOKEN}"

if [ -z "$LAKEHOUSE_ID" ]; then
    echo "Usage: $0 <lakehouse_id>"
    exit 1
fi

if [ -z "$TOKEN" ]; then
    echo "Error: CRIBL_TOKEN environment variable not set"
    exit 1
fi

# Check if lakehouse exists and get its status
RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
    "$SERVER_URL/api/v1/orgs/$ORGANIZATION_ID/workspaces/$WORKSPACE_ID/lakehouses/$LAKEHOUSE_ID")

if [ $? -ne 0 ]; then
    echo '{"ready": "false", "status": "error", "message": "Failed to connect to Cribl API"}'
    exit 1
fi

# Extract status from response (adjust the jq path based on actual API response structure)
STATUS=$(echo "$RESPONSE" | jq -r '.status // .state // "unknown"')

# Check if lakehouse is ready (adjust these status values based on your API)
case "$STATUS" in
    "ready"|"active"|"running")
        echo "{\"ready\": \"true\", \"status\": \"$STATUS\"}"
        ;;
    "creating"|"starting"|"provisioning")
        echo "{\"ready\": \"false\", \"status\": \"$STATUS\"}"
        ;;
    "error"|"failed")
        echo "{\"ready\": \"false\", \"status\": \"$STATUS\", \"error\": true}"
        ;;
    *)
        echo "{\"ready\": \"false\", \"status\": \"$STATUS\", \"unknown\": true}"
        ;;
esac 