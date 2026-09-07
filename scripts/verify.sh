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
      unformatted=$("$GOFMT" -l internal cmd test desktop)
      if [ -n "$unformatted" ]; then
        printf '%s\n' "$unformatted" >&2
        printf 'format: unformatted Go files found\n' >&2
        return 1
      fi
      "$GO" run ./scripts/github-format
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
      CGO_ENABLED=0 "$GO" run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.0 run ./...
      ;;
    race)
      CGO_ENABLED=1 "$GO" test -race ./...
      ;;
    gui)
      (
        cd desktop
        "$GO" test -race ./...
        "$GO" run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean
      )
      (
        cd desktop/frontend
        npm test
        npm run build
      )
      ;;
    coverage)
      # Keep the profile relative to the repository so a Windows Go toolchain
      # invoked from WSL Bash receives a path it can open. A WSL /tmp path is not
      # meaningful to a native Windows process.
      coverage_profile=$(mktemp ./goschedule-cover.XXXXXX)
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
      "$SH" test/scripts/automation-check_test.sh automation
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
