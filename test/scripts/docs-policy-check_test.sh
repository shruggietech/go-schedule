#!/bin/sh
set -eu

ROOT=${1:-.}
CHECK="$ROOT/scripts/docs-policy-check.sh"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

copy_fixture() {
  dest=$1
  for file in \
    README.md docs/README.md docs/cli.md docs/cron.md docs/gui-fields.md \
    internal/cli/cli.go internal/cli/cron_test.go \
    specs/001-task-scheduler/spec.md \
    specs/001-task-scheduler/data-model.md \
    specs/001-task-scheduler/quickstart.md \
    specs/001-task-scheduler/contracts/cli.md \
    specs/001-task-scheduler/contracts/local-api.md; do
    mkdir -p "$dest/$(dirname "$file")"
    cp "$ROOT/$file" "$dest/$file"
  done
}

GOOD="$TMP/good"
copy_fixture "$GOOD"
sh "$CHECK" "$GOOD"

STALE="$TMP/stale"
copy_fixture "$STALE"
printf '\nThe desktop editor remains plain-language only.\n' >> "$STALE/docs/cli.md"
if sh "$CHECK" "$STALE" >/dev/null 2>&1; then
  printf '%s\n' 'docs-policy fixture: stale categorical claim passed unexpectedly' >&2
  exit 1
fi

MISSING="$TMP/missing"
copy_fixture "$MISSING"
printf '%s\n' 'current product copy omitted' > "$MISSING/README.md"
if sh "$CHECK" "$MISSING" >/dev/null 2>&1; then
  printf '%s\n' 'docs-policy fixture: missing dual-syntax posture passed unexpectedly' >&2
  exit 1
fi

STALE_BREADTH="$TMP/stale-breadth"
copy_fixture "$STALE_BREADTH"
sed 's/lists, ranges, and field-local steps/lists and ranges/' \
  "$STALE_BREADTH/docs/cron.md" > "$STALE_BREADTH/docs/cron.md.tmp"
mv "$STALE_BREADTH/docs/cron.md.tmp" "$STALE_BREADTH/docs/cron.md"
if sh "$CHECK" "$STALE_BREADTH" >/dev/null 2>&1; then
  printf '%s\n' 'docs-policy fixture: missing composite breadth passed unexpectedly' >&2
  exit 1
fi

printf '%s\n' 'docs-policy fixtures: OK (aligned copy accepted; stale and missing copy rejected)'
