#!/bin/sh
# Offline consistency check for Spec-Kit lifecycle metadata.

set -eu

ROOT=${1:-.}
SPECS="$ROOT/specs"
INVENTORY="$SPECS/README.md"
FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT HUP INT TERM

report() {
  printf '%s\n' "$*" >> "$FAILURES"
}

if [ ! -f "$INVENTORY" ]; then
  printf 'spec-lifecycle-check: inventory not found: %s\n' "$INVENTORY" >&2
  exit 1
fi

found=0
for spec in "$SPECS"/[0-9][0-9][0-9]-*/spec.md; do
  [ -f "$spec" ] || continue
  found=$((found + 1))
  feature=$(basename "$(dirname "$spec")")
  status=$(sed -n 's/^\*\*Status\*\*: *//p' "$spec" | head -n 1)
  delivery=$(sed -n 's/^\*\*Delivery\*\*: *//p' "$spec" | head -n 1)
  tasks=$(dirname "$spec")/tasks.md

  case "$status" in
    Draft | Ready | 'In Progress' | Implemented | Deferred | Superseded | Abandoned) ;;
    '') report "$spec: missing status" ;;
    *) report "$spec: invalid status: $status" ;;
  esac

  if ! grep -Fq "| $feature | $status |" "$INVENTORY"; then
    report "$spec: missing inventory row for $feature with state $status"
  fi

  open=0
  checked=0
  if [ -f "$tasks" ]; then
    open=$(grep -Ec '^- \[ \] T[0-9]+' "$tasks" || true)
    checked=$(grep -Eic '^- \[[xX]\] T[0-9]+' "$tasks" || true)
  fi

  case "$status" in
    Draft | Ready)
      if [ "$open" -eq 0 ]; then
        report "$spec: $status has no actionable task"
      fi
      ;;
    'In Progress')
      if [ "$open" -eq 0 ] || [ "$checked" -eq 0 ]; then
        report "$spec: In Progress must have completed and actionable tasks"
      fi
      ;;
    Implemented)
      if [ -z "$delivery" ]; then
        report "$spec: missing delivery evidence"
      elif ! printf '%s\n' "$delivery" | grep -Eiq \
        '(commit[[:space:]]|releases?[[:space:]]|pull/[0-9]+|PR #[0-9]+|review branch[[:space:]])'; then
        report "$spec: invalid delivery evidence: $delivery"
      fi
      if [ "$open" -ne 0 ]; then
        report "$spec: Implemented has $open unresolved task(s)"
      fi
      ;;
    Deferred | Superseded | Abandoned)
      if ! grep -Eq '^\*\*Disposition\*\*: .+' "$spec"; then
        report "$spec: $status is missing a disposition"
      fi
      ;;
  esac
done

if [ "$found" -eq 0 ]; then
  report "$SPECS: no feature specifications found"
fi

inventory_count=$(grep -Ec '^\| [0-9][0-9][0-9]-[^|]+ \|' "$INVENTORY" || true)
if [ "$inventory_count" -ne "$found" ]; then
  report "$INVENTORY: inventory has $inventory_count feature row(s), discovered $found"
fi

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'spec-lifecycle-check: FAILED with %s issue(s)\n' \
    "$(wc -l < "$FAILURES" | tr -d ' ')" >&2
  exit 1
fi

printf 'spec-lifecycle-check: OK - %s specification(s) are lifecycle-consistent\n' "$found"
