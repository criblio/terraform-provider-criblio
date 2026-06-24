#!/usr/bin/env bash
set -euo pipefail

testdata_root="tests/acceptance/testdata"

if [[ ! -d "$testdata_root" ]]; then
  echo "acceptance testdata root not found: $testdata_root" >&2
  exit 1
fi

testdata_dirs="$(
  find "$testdata_root" -mindepth 1 -maxdepth 1 -type d |
    sed -E "s#^${testdata_root}/##" |
    sort -u
)"

referenced_dirs="$(
  {
    for test_file in tests/acceptance/*_test.go; do
      [[ -f "$test_file" ]] || continue
      if grep -q 'config\.TestNameDirectory()' "$test_file"; then
        sed -nE 's/^func (Test[A-Za-z0-9_]+)\(t \*testing\.T\).*/\1/p' "$test_file"
      fi
    done

    grep -RhoE 'testdata/Test[A-Za-z0-9_]+' tests/acceptance --include='*_test.go' 2>/dev/null |
      sed -E 's#^testdata/##' ||
      true

    for test_file in tests/acceptance/*_test.go; do
      [[ -f "$test_file" ]] || continue
      if grep -q '"testdata"' "$test_file"; then
        sed -nE 's/.*"(Test[A-Za-z0-9_]+)".*/\1/p' "$test_file"
      fi
    done
  } | sort -u
)"

unused="$(
  comm -23 \
    <(printf '%s\n' "$testdata_dirs" | sed '/^$/d') \
    <(printf '%s\n' "$referenced_dirs" | sed '/^$/d')
)"

if [[ -z "$unused" ]]; then
  echo "No unused acceptance testdata directories found."
  exit 0
fi

echo "Unused acceptance testdata directories:"
printf '%s\n' "$unused" | sed "s#^#${testdata_root}/#"
exit 1
