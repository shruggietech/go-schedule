#!/bin/sh
# Contract regressions for the maintainer automation baseline.

set -eu

ROOT=$(CDPATH='' cd "$(dirname "$0")/../.." && pwd)
CHECK="$ROOT/scripts/automation-check.sh"
VERIFY="$ROOT/scripts/verify.sh"
MODE=${1:-all}

fail() {
  printf 'automation-check-test: FAIL: %s\n' "$*" >&2
  exit 1
}

assert_contains() {
  haystack=$1
  needle=$2
  case "$haystack" in
    *"$needle"*) : ;;
    *) fail "expected output to contain '$needle'; got: $haystack" ;;
  esac
}

run_expect_pass() {
  label=$1
  shift
  if ! output=$("$@" 2>&1); then
    fail "$label should pass; got: $output"
  fi
}

run_expect_fail() {
  label=$1
  needle=$2
  shift 2
  if output=$("$@" 2>&1); then
    fail "$label should fail"
  fi
  assert_contains "$output" "$needle"
}

make_fixture() {
  fixture=$1
  manifest=$2
  mkdir -p "$fixture/.github/workflows" "$fixture/scripts"
  cat > "$fixture/.github/workflows/ci.yml" <<'EOF'
name: fixture
jobs:
  test:
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
      - uses: actions/upload-artifact@v7
      - uses: softprops/action-gh-release@v3
EOF
  cat > "$fixture/scripts/verify.sh" <<EOF
#!/bin/sh
if [ "\${1:-}" = list ]; then
  printf '%s\\n' $manifest
  exit 0
fi
exit 2
EOF
}

run_automation_cases() {
  [ -f "$CHECK" ] || fail "missing $CHECK"

  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT HUP INT TERM
  good="$tmp/good"
  make_fixture "$good" 'format vet lint race gui coverage docs automation'

  run_expect_pass approved sh "$CHECK" "$good"

  old="$tmp/old"
  cp -R "$good" "$old"
  sed 's#actions/checkout@v7#actions/checkout@v4#' \
    "$good/.github/workflows/ci.yml" > "$old/.github/workflows/ci.yml"
  run_expect_fail obsolete 'actions/checkout@v4' sh "$CHECK" "$old"

  unknown="$tmp/unknown"
  cp -R "$good" "$unknown"
  printf '      - uses: example/unknown@v1\n' >> \
    "$unknown/.github/workflows/ci.yml"
  run_expect_fail unknown 'example/unknown@v1' sh "$CHECK" "$unknown"

  missing="$tmp/missing"
  make_fixture "$missing" 'format vet lint race gui coverage docs'
  run_expect_fail missing 'gate manifest differs' sh "$CHECK" "$missing"

  duplicate="$tmp/duplicate"
  make_fixture "$duplicate" \
    'format vet lint race gui coverage docs automation automation'
  run_expect_fail duplicate 'gate manifest differs' sh "$CHECK" "$duplicate"

  extra="$tmp/extra"
  make_fixture "$extra" \
    'format vet lint race gui coverage docs automation build'
  run_expect_fail extra 'gate manifest differs' sh "$CHECK" "$extra"

  rm -rf "$tmp"
  trap - EXIT HUP INT TERM
}

run_verify_cases() {
  [ -f "$VERIFY" ] || fail "missing $VERIFY"

  run_expect_fail unknown-mode 'usage:' sh "$VERIFY" unknown

  if output=$(GOFMT=false sh "$VERIFY" all 2>&1); then
    fail 'controlled aggregate child failure should fail'
  fi
  assert_contains "$output" '[format]'
  case "$output" in
    *'[vet]'*) fail 'aggregate continued after the controlled format failure' ;;
    *) : ;;
  esac
}

case "$MODE" in
  all)
    run_automation_cases
    run_verify_cases
    ;;
  automation)
    run_automation_cases
    ;;
  failure-propagation)
    run_verify_cases
    ;;
  *)
    fail "unknown test mode: $MODE"
    ;;
esac

printf 'automation-check-test: OK (%s)\n' "$MODE"
