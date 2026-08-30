#!/bin/sh
# Contract fixtures for specification lifecycle consistency.

set -eu

ROOT=$(CDPATH='' cd "$(dirname "$0")/../.." && pwd)
CHECK="$ROOT/scripts/spec-lifecycle-check.sh"

fail() {
  printf 'spec-lifecycle-check-test: FAIL: %s\n' "$*" >&2
  exit 1
}

write_spec() {
  root=$1
  status=$2
  delivery=$3
  mkdir -p "$root/specs/001-example"
  {
    printf '# Feature Specification: Example\n\n'
    printf '**Status**: %s\n\n' "$status"
    if [ -n "$delivery" ]; then
      printf '**Delivery**: %s\n' "$delivery"
    fi
  } > "$root/specs/001-example/spec.md"
}

write_inventory() {
  root=$1
  status=$2
  mkdir -p "$root/specs"
  {
    printf '# Specification Lifecycle\n\n'
    printf '| Feature | State | Delivery |\n'
    printf '| --- | --- | --- |\n'
    printf '| 001-example | %s | commit abc1234 |\n' "$status"
  } > "$root/specs/README.md"
}

write_tasks() {
  root=$1
  marker=$2
  printf -- '- [%s] T001 Example task in example.txt.\n' "$marker" > \
    "$root/specs/001-example/tasks.md"
}

expect_pass() {
  label=$1
  root=$2
  if ! output=$(sh "$CHECK" "$root" 2>&1); then
    fail "$label should pass; got: $output"
  fi
}

expect_fail() {
  label=$1
  needle=$2
  root=$3
  if output=$(sh "$CHECK" "$root" 2>&1); then
    fail "$label should fail"
  fi
  case "$output" in
    *"$needle"*) : ;;
    *) fail "$label expected '$needle'; got: $output" ;;
  esac
}

[ -f "$CHECK" ] || fail "missing $CHECK"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM

good="$tmp/good"
write_spec "$good" Implemented 'commit abc1234'
write_tasks "$good" x
write_inventory "$good" Implemented
expect_pass good "$good"

invalid="$tmp/invalid"
write_spec "$invalid" Finished 'commit abc1234'
write_tasks "$invalid" x
write_inventory "$invalid" Finished
expect_fail invalid 'invalid status' "$invalid"

missing_delivery="$tmp/missing-delivery"
write_spec "$missing_delivery" Implemented ''
write_tasks "$missing_delivery" x
write_inventory "$missing_delivery" Implemented
expect_fail missing-delivery 'missing delivery evidence' "$missing_delivery"

placeholder_delivery="$tmp/placeholder-delivery"
write_spec "$placeholder_delivery" Implemented 'Pending implementation'
write_tasks "$placeholder_delivery" x
write_inventory "$placeholder_delivery" Implemented
expect_fail placeholder-delivery 'invalid delivery evidence' "$placeholder_delivery"

arbitrary_delivery="$tmp/arbitrary-delivery"
write_spec "$arbitrary_delivery" Implemented 'done somewhere'
write_tasks "$arbitrary_delivery" x
write_inventory "$arbitrary_delivery" Implemented
expect_fail arbitrary-delivery 'invalid delivery evidence' "$arbitrary_delivery"

draft_complete="$tmp/draft-complete"
write_spec "$draft_complete" Draft ''
write_tasks "$draft_complete" x
write_inventory "$draft_complete" Draft
expect_fail draft-complete 'has no actionable task' "$draft_complete"

implemented_open="$tmp/implemented-open"
write_spec "$implemented_open" Implemented 'commit abc1234'
write_tasks "$implemented_open" ' '
write_inventory "$implemented_open" Implemented
expect_fail implemented-open 'unresolved task' "$implemented_open"

missing_inventory="$tmp/missing-inventory"
write_spec "$missing_inventory" Implemented 'commit abc1234'
write_tasks "$missing_inventory" x
mkdir -p "$missing_inventory/specs"
printf '# Specification Lifecycle\n' > "$missing_inventory/specs/README.md"
expect_fail missing-inventory 'missing inventory row' "$missing_inventory"

printf 'spec-lifecycle-check-test: OK\n'
