#!/bin/bash

# Simple test script to download and inspect the diag bundle from /system/diag/download
# This helps understand the structure before building the full CLI tool

set -e

# Configuration - set these environment variables or edit below
SERVER_URL="${CRIBL_SERVER_URL:-https://main-org.cribl.cloud}"
BEARER_TOKEN="${CRIBL_BEARER_TOKEN:-}"
WORKSPACE_ID="${CRIBL_WORKSPACE_ID:-main}"
ORG_ID="${CRIBL_ORGANIZATION_ID:-org}"

# Build the full URL
if [[ "$SERVER_URL" == *"localhost"* ]] || [[ "$SERVER_URL" == *"127.0.0.1"* ]] || [[ "$SERVER_URL" == *":9000"* ]]; then
  # On-prem
  API_URL="${SERVER_URL}/api/v1/system/diag/download"
else
  # Cloud
  API_URL="https://${WORKSPACE_ID}-${ORG_ID}.${SERVER_URL#https://}/api/v1/system/diag/download"
fi

OUTPUT_DIR="./test-diag-output"
ARCHIVE_FILE="${OUTPUT_DIR}/diag-bundle.tar.gz"

echo "🔍 Testing /system/diag/download endpoint"
echo "📍 URL: $API_URL"
echo ""

# Check for authentication
if [ -z "$BEARER_TOKEN" ]; then
  echo "❌ Error: CRIBL_BEARER_TOKEN environment variable not set"
  echo ""
  echo "Set it with:"
  echo "  export CRIBL_BEARER_TOKEN='your-token-here'"
  echo ""
  echo "Or set other auth variables:"
  echo "  export CRIBL_SERVER_URL='https://main-org.cribl.cloud'"
  echo "  export CRIBL_WORKSPACE_ID='main'"
  echo "  export CRIBL_ORGANIZATION_ID='org'"
  exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Download the diag bundle
echo "📥 Downloading diag bundle..."
HTTP_CODE=$(curl -s -w "%{http_code}" -o "$ARCHIVE_FILE" \
  -H "Authorization: Bearer $BEARER_TOKEN" \
  -H "Accept: application/tar+gzip" \
  "$API_URL")

if [ "$HTTP_CODE" -ne 200 ]; then
  echo "❌ Error: HTTP $HTTP_CODE"
  echo "Response:"
  cat "$ARCHIVE_FILE"
  rm -f "$ARCHIVE_FILE"
  exit 1
fi

echo "✅ Downloaded: $ARCHIVE_FILE ($(du -h "$ARCHIVE_FILE" | cut -f1))"
echo ""

# Extract the archive
echo "📦 Extracting archive..."
EXTRACT_DIR="${OUTPUT_DIR}/extracted"
mkdir -p "$EXTRACT_DIR"
tar -xzf "$ARCHIVE_FILE" -C "$EXTRACT_DIR"

echo "✅ Extracted to: $EXTRACT_DIR"
echo ""

# List the structure
echo "📁 Archive structure:"
find "$EXTRACT_DIR" -type f -name "*.yaml" -o -name "*.yml" -o -name "*.json" | head -20
echo ""

# Find a single resource type (e.g., sources)
echo "🔍 Looking for source configurations..."
SOURCE_FILES=$(find "$EXTRACT_DIR" -type f \( -path "*/sources/*.yaml" -o -path "*/sources/*.yml" \) | head -5)

if [ -n "$SOURCE_FILES" ]; then
  echo "Found source files:"
  echo "$SOURCE_FILES" | while read -r file; do
    echo "  📄 $(basename "$file")"
  done
  echo ""
  
  # Show first source file as example
  FIRST_SOURCE=$(echo "$SOURCE_FILES" | head -1)
  if [ -n "$FIRST_SOURCE" ]; then
    echo "📖 Example source file: $FIRST_SOURCE"
    echo "─────────────────────────────────────────"
    head -30 "$FIRST_SOURCE"
    echo "─────────────────────────────────────────"
  fi
else
  echo "⚠️  No source files found in expected locations"
  echo ""
  echo "Full directory structure:"
  find "$EXTRACT_DIR" -type d | head -20
fi

echo ""
echo "✨ Done! Explore the extracted files in: $EXTRACT_DIR"
echo ""
echo "To find specific resource types:"
echo "  find $EXTRACT_DIR -path '*/sources/*' -type f    # Sources"
echo "  find $EXTRACT_DIR -path '*/destinations/*' -type f  # Destinations"
echo "  find $EXTRACT_DIR -path '*/pipelines/*' -type f   # Pipelines"
echo "  find $EXTRACT_DIR -path '*/routes/*' -type f      # Routes"
echo "  find $EXTRACT_DIR -path '*/packs/*' -type f       # Packs"

