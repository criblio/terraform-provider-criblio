#!/usr/bin/env bash
set -euo pipefail

fail=0

check_no_match() {
  local description="$1"
  shift
  if "$@"; then
    echo "FAIL: ${description}" >&2
    fail=1
  else
    echo "ok: ${description}"
  fi
}

check_missing() {
  local path="$1"
  if [[ -e "$path" ]]; then
    echo "FAIL: ${path} still exists" >&2
    fail=1
  else
    echo "ok: ${path} is absent"
  fi
}

check_no_match 'clients.Legacy under internal/provider' grep -R "clients.Legacy" internal/provider
check_no_match 'sdk.CriblIo under internal outside deleted SDK' grep -R "sdk.CriblIo" internal --exclude-dir=sdk
check_no_match 'internal/sdk imports outside deleted SDK' grep -R "internal/sdk" internal tools --exclude-dir=sdk
check_no_match 'provider reflect merge helpers under internal/provider' grep -R -E "tfReflect|provider/reflect|merge\\(ctx" internal/provider
check_no_match 'provider imports of internal plan modifiers' grep -R "internal/planmodifiers" internal/provider
check_no_match 'provider imports of Speakeasy model types' grep -R "internal/provider/types" internal/provider
check_no_match 'Speakeasy mentions under internal and tools' grep -R -i "speakeasy" internal tools
check_no_match 'Speakeasy annotations or text in openapi.yml' grep -i "speakeasy" openapi.yml

check_missing internal/sdk
check_missing internal/provider/types
check_missing internal/planmodifiers
check_missing internal/provider/reflect
check_missing .speakeasy
check_missing .genignore

if [[ $fail -ne 0 ]]; then
  exit 1
fi
