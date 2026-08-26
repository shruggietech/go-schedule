#!/bin/sh
# Canonical, non-mutating pre-push verification driver.

set -eu

ROOT=$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)
GO=${GO:-go}
GOFMT=${GOFMT:-gofmt}
SH=${SH:-sh}
GATES='format vet lint race gui coverage docs automation'

usage() {
  printf 'usage: %s {list|all|format|vet|lint|race|gui|coverage|docs|automation}\n' "$0" >&2
}

run_gate() {
  gate=$1
  printf '[%s]\n' "$gate"

  case "$gate" in
    format)
      unformatted=$("$GOFMT" -l internal cmd test)
      if [ -n "$unformatted" ]; then
        printf '%s\n' "$unformatted" >&2
        printf 'format: unformatted Go files found\n' >&2
        return 1
      fi
      ;;
    vet)
      CGO_ENABLED=0 "$GO" vet ./...
      ;;
    lint)
      if [ -z "${GOTOOLCHAIN:-}" ]; then
        module_go=$(awk '$1 == "go" { print $2; exit }' go.mod)
        if [ -z "$module_go" ]; then
          printf 'lint: go.mod has no go directive\n' >&2
          return 1
        fi
        GOTOOLCHAIN="go$module_go"
        export GOTOOLCHAIN
      fi
      CGO_ENABLED=0 "$GO" run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
      ;;
    race)
      # Enumerate with cgo off so excluded desktop packages cannot make package
      # discovery depend on OpenGL headers; the selected tests still run with cgo.
      all_packages=$(CGO_ENABLED=0 "$GO" list ./...)
      packages=$(printf '%s\n' "$all_packages" | grep -vE '/cmd/gosched-gui$|/gui$' || true)
      if [ -z "$packages" ]; then
        printf 'race: package selection is empty\n' >&2
        return 1
      fi
      # Deliberate word splitting passes the newline-delimited package list as arguments.
      # shellcheck disable=SC2086
      CGO_ENABLED=1 "$GO" test -race $packages
      ;;
    gui)
      "$GO" test ./gui/...
      ;;
    coverage)
      coverage_profile=$(mktemp)
      if CGO_ENABLED=0 COVERAGE_PROFILE="$coverage_profile" \
        "$SH" scripts/coverage-gate.sh; then
        rm -f "$coverage_profile"
      else
        status=$?
        rm -f "$coverage_profile"
        return "$status"
      fi
      ;;
    docs)
      "$SH" scripts/docs-check.sh
      ;;
    automation)
      "$SH" scripts/automation-check.sh "$ROOT"
      ;;
    *)
      printf 'verify: unknown gate: %s\n' "$gate" >&2
      return 2
      ;;
  esac
}

cd "$ROOT"

case "${1:-}" in
  list)
    for gate in $GATES; do
      printf '%s\n' "$gate"
    done
    ;;
  all)
    for gate in $GATES; do
      run_gate "$gate"
    done
    ;;
  format | vet | lint | race | gui | coverage | docs | automation)
    run_gate "$1"
    ;;
  *)
    usage
    exit 2
    ;;
esac
