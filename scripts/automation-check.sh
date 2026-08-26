#!/bin/sh
# Offline contract check for workflow action majors and local verification gates.

set -eu

ROOT=${1:-.}

if [ ! -d "$ROOT/.github/workflows" ]; then
  printf 'automation-check: workflow directory not found: %s\n' \
    "$ROOT/.github/workflows" >&2
  exit 1
fi

FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT HUP INT TERM

report() {
  printf '%s\n' "$*" >> "$FAILURES"
}

approved_for() {
  case "$1" in
    actions/checkout) printf '%s' 'actions/checkout@v7' ;;
    actions/setup-go) printf '%s' 'actions/setup-go@v7' ;;
    actions/upload-artifact) printf '%s' 'actions/upload-artifact@v7' ;;
    softprops/action-gh-release) printf '%s' 'softprops/action-gh-release@v3' ;;
    *) return 1 ;;
  esac
}

check_action() {
  file=$1
  line=$2
  ref=$3
  family=${ref%@*}

  if approved=$(approved_for "$family"); then
    if [ "$ref" != "$approved" ]; then
      report "$file:$line: action reference $ref is not approved; expected $approved"
    fi
  else
    report "$file:$line: unaudited action reference $ref"
  fi
}

found_workflow=0
for file in "$ROOT"/.github/workflows/*.yml "$ROOT"/.github/workflows/*.yaml; do
  [ -f "$file" ] || continue
  found_workflow=1
  awk '
    match($0, /uses:[[:space:]]*[^[:space:]#]+/) {
      ref = substr($0, RSTART, RLENGTH)
      sub(/^uses:[[:space:]]*/, "", ref)
      gsub(/["'"'"']/, "", ref)
      print FNR "|" ref
    }
  ' "$file" | while IFS='|' read -r line ref; do
    check_action "$file" "$line" "$ref"
  done
done

if [ "$found_workflow" -eq 0 ]; then
  report "$ROOT/.github/workflows: no .yml or .yaml workflow files found"
fi

EXPECTED_GATES='format
vet
lint
race
gui
coverage
docs
automation'

VERIFY="$ROOT/scripts/verify.sh"
if [ ! -f "$VERIFY" ]; then
  report "$VERIFY: verification driver not found"
elif ! observed=$(sh "$VERIFY" list 2>&1); then
  report "$VERIFY: could not read gate manifest: $observed"
elif [ "$observed" != "$EXPECTED_GATES" ]; then
  report "gate manifest differs"
  report "expected:"
  report "$EXPECTED_GATES"
  report "observed:"
  report "$observed"
fi

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'automation-check: FAILED with %s issue(s)\n' \
    "$(grep -c '^[^[:space:]]' "$FAILURES")" >&2
  exit 1
fi

printf 'automation-check: OK - approved actions and 8-gate manifest\n'
