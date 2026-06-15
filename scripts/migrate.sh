#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RESOURCE="${RESOURCE:-}"
BUCKET="${BUCKET:-}"
DRY_RUN="${DRY_RUN:-false}"
MODE="${1:-migrate}"

cd "$ROOT"

usage() {
  cat <<'EOF'
Usage:
  make migrate RESOURCE=<name> [DRY_RUN=true]
  make migrate-batch BUCKET=A [DRY_RUN=true]
EOF
}

is_true() {
  case "${1:-}" in
    1|true|TRUE|yes|YES) return 0 ;;
    *) return 1 ;;
  esac
}

confirm() {
  local prompt="$1"
  if is_true "$DRY_RUN"; then
    echo "DRY_RUN=true: would pause for approval: $prompt"
    return 0
  fi
  printf "%s Type yes to continue: " "$prompt"
  local answer
  read -r answer
  [[ "$answer" == "yes" ]]
}

genignored_resource() {
  local resource="$1"
  grep -qx "internal/provider/${resource}_resource.go" .genignore 2>/dev/null
}

companion_files() {
  local resource="$1"
  awk -v resource="$resource" '
    $0 ~ "^internal/provider/" resource "_" && $0 !~ "_resource.go$" && $0 !~ "_resource_sdk.go$" { print }
    $0 ~ "^internal/provider/" resource "[^/]*\\.go$" && $0 !~ "_resource.go$" && $0 !~ "_resource_sdk.go$" { print }
  ' .genignore 2>/dev/null | sort -u
}

contract_path() {
  local resource="$1"
  for path in \
    "internal/provider/${resource}_behavioral_contract.json" \
    "tools/behavioral-extractor/testdata/golden/${resource}.json" \
    "tools/agentic-codegen/testdata/${resource}_behavioral_contract.json"
  do
    [[ -f "$path" ]] && { echo "$path"; return 0; }
  done
  return 1
}

contract_decision() {
  local resource="$1"
  local path
  path="$(contract_path "$resource" 2>/dev/null || true)"
  [[ -z "$path" ]] && return 1
  sed -n 's/.*"codegen_decision"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$path" | head -1
}

has_companion_sdk_refs() {
  local resource="$1"
  local files
  files="$(companion_files "$resource")"
  [[ -z "$files" ]] && return 1
  while IFS= read -r file; do
    [[ -f "$file" ]] || continue
    if grep -Eq '(^|[^A-Za-z0-9_])sdk\.|internal/sdk' "$file"; then
      return 0
    fi
  done <<< "$files"
  return 1
}

has_plan_modifier_definitions() {
  local resource="$1"
  local file="internal/provider/${resource}_resource.go"
  [[ -f "$file" ]] || return 1
  grep -Eq 'type .*planmodifier|func .*PlanModify|SuppressWhitespaceDiff|UseHoistedValue|PreferConfigOrState' "$file"
}

classify_resource() {
  local resource="$1"
  local decision
  decision="$(contract_decision "$resource" || true)"
  case "$decision" in
    codegen_direct) echo "A"; return 0 ;;
    overlay_augmented) echo "B"; return 0 ;;
    overlay_plus_companion|hand_written) echo "C"; return 0 ;;
  esac

  if ! genignored_resource "$resource"; then
    echo "A"
    return 0
  fi

  if has_companion_sdk_refs "$resource" || has_plan_modifier_definitions "$resource"; then
    echo "C"
    return 0
  fi

  case "$resource" in
    lookupfile|destination) echo "B" ;;
    *) echo "A" ;;
  esac
}

extract_contract() {
  local resource="$1"
  if [[ -d tools/behavioral-extractor ]]; then
    go run ./tools/behavioral-extractor --resource "$resource" --dry-run
    return 0
  fi
  local path
  path="$(contract_path "$resource" 2>/dev/null || true)"
  [[ -n "$path" ]] && { cat "$path"; return 0; }
  echo "No behavioral extractor or recorded contract found for ${resource}" >&2
  return 1
}

show_overlay_additions() {
  local resource="$1"
  local path
  path="$(contract_path "$resource" 2>/dev/null || true)"
  echo "--- overlay additions for ${resource} ---"
  case "$resource" in
    lookupfile)
      cat <<'EOF'
- target: "$.components.schemas.LookupFile.properties.content"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.description"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.tags"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.version"
  update:
    x-terraform-prefer-state: true
EOF
      ;;
    *)
      if [[ -n "$path" ]]; then
        echo "# Review ${path}; no deterministic overlay additions are recorded for ${resource}."
      else
        echo "# No deterministic overlay additions are recorded for ${resource}."
      fi
      ;;
  esac
}

append_lookupfile_overlay() {
  if grep -q 'LookupFile.properties.content' terraform-overlay.yml 2>/dev/null; then
    echo "overlay already contains lookupfile prefer-state entries"
    return 0
  fi
  cat >> terraform-overlay.yml <<'EOF'

# Added by make migrate RESOURCE=lookupfile.
- target: "$.components.schemas.LookupFile.properties.content"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.description"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.tags"
  update:
    x-terraform-prefer-state: true
- target: "$.components.schemas.LookupFile.properties.version"
  update:
    x-terraform-prefer-state: true
EOF
}

run_generate() {
  if is_true "$DRY_RUN"; then
    echo "DRY_RUN=true: would run make generate"
    return 0
  fi
  make generate
}

run_build() {
  if is_true "$DRY_RUN"; then
    echo "DRY_RUN=true: would run go build ./..."
    return 0
  fi
  if go build ./...; then
    echo "go build ./... passed"
  else
    echo "go build ./... failed" >&2
    return 1
  fi
}

run_acceptance() {
  local resource="$1"
  local test="./tests/acceptance"
  if [[ -f "tests/acceptance/${resource}_test.go" ]]; then
    test="./tests/acceptance -run ${resource}"
  fi
  if is_true "$DRY_RUN"; then
    echo "DRY_RUN=true: would run acceptance test: go test -v -timeout 20m ${test}"
    return 0
  fi
  # shellcheck disable=SC2086
  if go test -v -timeout 20m ${test}; then
    echo "acceptance test passed"
  else
    echo "acceptance test failed" >&2
    return 1
  fi
}

migrate_bucket_a() {
  local resource="$1"
  echo "Bucket A path: codegen-direct"
  run_generate
  run_build
  run_acceptance "$resource"
}

migrate_bucket_b() {
  local resource="$1"
  echo "Bucket B path: extractor -> overlay approval -> make generate"
  echo "--- behavioral contract ---"
  extract_contract "$resource"
  show_overlay_additions "$resource"
  confirm "Apply overlay additions for ${resource}?" || { echo "overlay approval rejected"; return 1; }
  if ! is_true "$DRY_RUN" && [[ "$resource" == "lookupfile" ]]; then
    append_lookupfile_overlay
  fi
  run_generate
  run_build
  run_acceptance "$resource"
}

migrate_bucket_c() {
  local resource="$1"
  echo "Bucket C path: extractor checkpoint -> make generate -> synthesizer checkpoint -> write"
  echo "Checkpoint 1: behavioral contract"
  extract_contract "$resource"
  confirm "Approve behavioral contract for ${resource}?" || { echo "contract approval rejected"; return 1; }
  run_generate
  echo "Checkpoint 2: synthesized diff"
  if [[ -d tools/agentic-codegen ]]; then
    if is_true "$DRY_RUN"; then
      go run ./tools/agentic-codegen --resource "$resource" --dry-run
    else
      go run ./tools/agentic-codegen --resource "$resource"
    fi
  else
    echo "tools/agentic-codegen is not available; would run synthesizer for ${resource}"
  fi
  run_build
  run_acceptance "$resource"
}

migrate_one() {
  local resource="$1"
  [[ -n "$resource" ]] || { usage; return 2; }
  local bucket
  bucket="$(classify_resource "$resource")"
  echo "Resource: ${resource}"
  echo "Classifier signals:"
  if genignored_resource "$resource"; then
    echo "- .genignore: yes"
  else
    echo "- .genignore: no"
  fi
  if has_companion_sdk_refs "$resource"; then
    echo "- companion sdk refs: yes"
  else
    echo "- companion sdk refs: no"
  fi
  if has_plan_modifier_definitions "$resource"; then
    echo "- plan modifier definitions: yes"
  else
    echo "- plan modifier definitions: no"
  fi
  echo "Bucket: ${bucket}"
  case "$bucket" in
    A) migrate_bucket_a "$resource" ;;
    B) migrate_bucket_b "$resource" ;;
    C) migrate_bucket_c "$resource" ;;
    *) echo "unknown bucket ${bucket}" >&2; return 1 ;;
  esac
}

batch_bucket_a() {
  local requested="${BUCKET:-}"
  [[ "$requested" == "A" ]] || { echo "make migrate-batch currently supports BUCKET=A only" >&2; return 2; }
  echo "Batch Bucket A path: unattended codegen-direct migration"
  local resources=()
  while IFS= read -r file; do
    local base resource
    base="$(basename "$file")"
    resource="${base%_resource.go}"
    [[ "$(classify_resource "$resource")" == "A" ]] && resources+=("$resource")
  done < <(find internal/provider -maxdepth 1 -name '*_resource.go' | sort)
  printf 'Bucket A resources (%d): %s\n' "${#resources[@]}" "${resources[*]}"
  if is_true "$DRY_RUN"; then
    echo "DRY_RUN=true: would run make generate once, then go build ./..."
    return 0
  fi
  make generate
  go build ./...
}

case "$MODE" in
  migrate) migrate_one "$RESOURCE" ;;
  batch) batch_bucket_a ;;
  *) usage; exit 2 ;;
esac
