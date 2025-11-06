#!/bin/bash

# Ultra-simple one-liner version for quick testing
# Usage: BEARER_TOKEN=your-token ./test-download-simple.sh

BEARER_TOKEN="${CRIBL_BEARER_TOKEN:-${BEARER_TOKEN}}"
WORKSPACE_ID="${CRIBL_WORKSPACE_ID:-main}"
ORG_ID="${CRIBL_ORGANIZATION_ID:-org}"

if [ -z "$BEARER_TOKEN" ]; then
  echo "Usage: CRIBL_BEARER_TOKEN=token WORKSPACE_ID=main ORG_ID=org $0"
  exit 1
fi

# Download
curl -s -H "Authorization: Bearer $BEARER_TOKEN" \
  "https://${WORKSPACE_ID}-${ORG_ID}.cribl.cloud/api/v1/system/diag/download" \
  -o diag-bundle.tar.gz

# Extract and find first source
mkdir -p extracted
tar -xzf diag-bundle.tar.gz -C extracted 2>/dev/null

# Show structure
echo "Archive contents:"
find extracted -type f | head -20

# Find and show first source
FIRST_SOURCE=$(find extracted -path "*/sources/*.yaml" -o -path "*/sources/*.yml" | head -1)
if [ -n "$FIRST_SOURCE" ]; then
  echo ""
  echo "First source file: $FIRST_SOURCE"
  echo "─────────────────────────────"
  cat "$FIRST_SOURCE"
fi

