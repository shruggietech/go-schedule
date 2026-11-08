#!/bin/sh
# Exercise the exact driver with inert children, including a broken-stop mutant.
set -eu
ROOT=$(CDPATH='' cd -- "$(dirname "$0")/../.." && pwd)
REAL_SH=$(command -v sh)
fail() { printf 'verify-driver-test: FAIL: %s\n' "$*" >&2; exit 1; }
tmp=$(mktemp -d "$ROOT/.verify-driver-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
fixture="$tmp/checkout"
mkdir -p "$fixture/scripts" "$fixture/desktop/frontend" "$tmp/bin"
cp "$ROOT/scripts/verify.sh" "$fixture/scripts/verify.sh"
printf 'module example.invalid/driver\ngo 1.26.0\n' > "$fixture/go.mod"
cat > "$tmp/bin/gofmt-stub" <<'EOF'
#!/bin/sh
printf 'GOFMT %s\n' "$*" >> "$DRIVER_LOG"
if [ "$DRIVER_FAIL" = yes ]; then
  printf 'driver-sentinel: format child failed\n' >&2
  exit 37
fi
EOF
cat > "$tmp/bin/go-stub" <<'EOF'
#!/bin/sh
printf 'GO %s\n' "$*" >> "$DRIVER_LOG"
EOF
cat > "$tmp/bin/sh-stub" <<'EOF'
#!/bin/sh
# Never execute an argument, even if the driver erroneously reaches automation.
printf 'SH %s\n' "$*" >> "$DRIVER_LOG"
EOF
cat > "$tmp/bin/npm" <<'EOF'
#!/bin/sh
printf 'NPM %s\n' "$*" >> "$DRIVER_LOG"
EOF
chmod +x "$tmp/bin/"*
run_driver() {
  DRIVER_LOG="$tmp/children.log" DRIVER_FAIL="$1" \
    GOFMT="$tmp/bin/gofmt-stub" GO="$tmp/bin/go-stub" SH="$tmp/bin/sh-stub" \
    GOTOOLCHAIN=go1.26.0 PATH="$tmp/bin:$PATH" \
    "$REAL_SH" "$2" "$3"
}
assert_stop() {
  [ "$status" -eq 37 ] || return 1
  [ "$output" = '[format]
driver-sentinel: format child failed' ] || return 1
  [ "$(cat "$tmp/children.log")" = 'GOFMT -l internal cmd test desktop' ]
}
: > "$tmp/children.log"
if output=$(run_driver yes "$fixture/scripts/verify.sh" all 2>&1); then status=0; else status=$?; fi
assert_stop || fail "early child status/diagnostic/stop changed (status=$status): $output"
: > "$tmp/children.log"
if ! output=$(run_driver no "$fixture/scripts/verify.sh" all 2>&1); then fail "inert success aggregate failed: $output"; fi
expected='[format]
[vet]
[lint]
[race]
[gui]
[coverage]
[docs]
[automation]'
[ "$output" = "$expected" ] || fail "gate order changed: $output"
# This assertion makes removal from the canonical automation branch observable.
[ "$(grep -c '^SH test/scripts/maintainer-commands_test.sh$' "$tmp/children.log")" -eq 1 ] || fail 'canonical automation must invoke the shared command corpus once'
if ! output=$(run_driver no "$fixture/scripts/verify.sh" list 2>&1); then fail "list failed: $output"; fi
[ "$output" = 'format
vet
lint
race
gui
coverage
docs
automation' ] || fail "gate manifest changed: $output"
if output=$(run_driver no "$fixture/scripts/verify.sh" unknown 2>&1); then status=0; else status=$?; fi
[ "$status" -eq 2 ] || fail "unknown mode status changed: $status"
case "$output" in *usage:*) : ;; *) fail "unknown mode lost usage: $output" ;; esac
# Deliberately suppress aggregate stopping. All expensive roles stay inert, so
# this mutant finishes even after reaching the nested automation invocation.
sed 's/run_gate "$gate"$/run_gate "$gate" || :/' "$fixture/scripts/verify.sh" > "$fixture/scripts/broken.sh"
cmp -s "$fixture/scripts/verify.sh" "$fixture/scripts/broken.sh" && fail 'broken-stop mutation did not apply'
: > "$tmp/children.log"
if output=$(run_driver yes "$fixture/scripts/broken.sh" all 2>&1); then status=0; else status=$?; fi
if assert_stop; then fail 'guard accepted broken stopping'; fi
grep -q '^SH test/scripts/maintainer-commands_test.sh$' "$tmp/children.log" || fail 'broken-stop mutation did not reach the inert automation child'
printf 'verify-driver-test: OK (status 37, eight gates, bounded broken-stop rejection)\n'
