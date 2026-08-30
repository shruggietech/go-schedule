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
  mkdir -p "$fixture/.github/workflows" "$fixture/scripts" "$fixture/specs"
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
  cat > "$fixture/.github/dependabot.yml" <<'EOF'
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    groups:
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: monthly
    open-pull-requests-limit: 5
    groups:
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
EOF
  cat > "$fixture/scripts/spec-lifecycle-check.sh" <<'EOF'
#!/bin/sh
printf 'fixture lifecycle OK\n'
EOF
  cat > "$fixture/.github/workflows/codeql.yml" <<'EOF'
name: CodeQL
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '17 4 * * 1'
permissions:
  contents: read
  security-events: write
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
      - uses: github/codeql-action/init@v4
        with:
          languages: go
          build-mode: manual
      - name: Build
        run: CGO_ENABLED=0 go build ./...
      - uses: github/codeql-action/analyze@v4
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

  old_codeql="$tmp/old-codeql"
  cp -R "$good" "$old_codeql"
  sed 's#github/codeql-action/init@v4#github/codeql-action/init@v3#' \
    "$good/.github/workflows/codeql.yml" > \
    "$old_codeql/.github/workflows/codeql.yml"
  run_expect_fail obsolete-codeql 'github/codeql-action/init@v3' \
    sh "$CHECK" "$old_codeql"

  unknown="$tmp/unknown"
  cp -R "$good" "$unknown"
  printf '      - uses: example/unknown@v1\n' >> \
    "$unknown/.github/workflows/ci.yml"
  run_expect_fail unknown 'example/unknown@v1' sh "$CHECK" "$unknown"

  missing_trigger="$tmp/missing-trigger"
  cp -R "$good" "$missing_trigger"
  sed '/  schedule:/,/cron:/d' "$good/.github/workflows/codeql.yml" > \
    "$missing_trigger/.github/workflows/codeql.yml"
  run_expect_fail missing-trigger 'weekly schedule trigger' \
    sh "$CHECK" "$missing_trigger"

  insufficient_permission="$tmp/insufficient-permission"
  cp -R "$good" "$insufficient_permission"
  sed '/  security-events: write/d' "$good/.github/workflows/codeql.yml" > \
    "$insufficient_permission/.github/workflows/codeql.yml"
  run_expect_fail insufficient-permission 'security-events: write' \
    sh "$CHECK" "$insufficient_permission"

  job_permission="$tmp/job-permission"
  cp -R "$good" "$job_permission"
  sed '/    runs-on: ubuntu-latest/a\    permissions: write-all' \
    "$good/.github/workflows/codeql.yml" > \
    "$job_permission/.github/workflows/codeql.yml"
  run_expect_fail job-permission 'job-level permissions override' \
    sh "$CHECK" "$job_permission"

  invalid_schedule="$tmp/invalid-schedule"
  cp -R "$good" "$invalid_schedule"
  sed "s/17 4 \* \* 1/0 0 31 2 */" \
    "$good/.github/workflows/codeql.yml" > \
    "$invalid_schedule/.github/workflows/codeql.yml"
  run_expect_fail invalid-schedule 'weekly schedule trigger' \
    sh "$CHECK" "$invalid_schedule"

  missing_analysis="$tmp/missing-analysis"
  cp -R "$good" "$missing_analysis"
  sed '/github\/codeql-action\/analyze@v4/d' \
    "$good/.github/workflows/codeql.yml" > \
    "$missing_analysis/.github/workflows/codeql.yml"
  run_expect_fail missing-analysis 'CodeQL analyze step' \
    sh "$CHECK" "$missing_analysis"

  missing_gomod="$tmp/missing-gomod"
  cp -R "$good" "$missing_gomod"
  sed '/  - package-ecosystem: gomod/,/  - package-ecosystem: github-actions/{
    /  - package-ecosystem: github-actions/!d
  }' "$good/.github/dependabot.yml" > \
    "$missing_gomod/.github/dependabot.yml"
  run_expect_fail missing-gomod 'gomod entry' sh "$CHECK" "$missing_gomod"

  noisy_cadence="$tmp/noisy-cadence"
  cp -R "$good" "$noisy_cadence"
  sed '0,/interval: weekly/s//interval: daily/' \
    "$good/.github/dependabot.yml" > \
    "$noisy_cadence/.github/dependabot.yml"
  run_expect_fail noisy-cadence 'weekly cadence' sh "$CHECK" "$noisy_cadence"

  missing_label="$tmp/missing-dependency-label"
  cp -R "$good" "$missing_label"
  sed 's/      - dependencies/      - maintenance/' \
    "$good/.github/dependabot.yml" > \
    "$missing_label/.github/dependabot.yml"
  run_expect_fail missing-dependency-label 'dependencies label' \
    sh "$CHECK" "$missing_label"

  bracket_secret="$tmp/bracket-secret"
  cp -R "$good" "$bracket_secret"
  sed "/    steps:/a\      - run: echo \"\${{ secrets['NAME'] }}\"" \
    "$good/.github/workflows/codeql.yml" > \
    "$bracket_secret/.github/workflows/codeql.yml"
  run_expect_fail bracket-secret 'must not consume' \
    sh "$CHECK" "$bracket_secret"

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
