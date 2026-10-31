#!/bin/sh
# Exercise the canonical coverage driver's artifact boundary and cleanup.

set -eu

ROOT=$(CDPATH='' cd -- "$(dirname "$0")/../.." && pwd)
mkdir -p "$ROOT/.test-temp"
fixture=$(mktemp -d "$ROOT/.test-temp/coverage-profile-test.XXXXXX")
fixture=$(CDPATH='' cd -- "$fixture" && pwd)
case "$fixture" in
    "$ROOT"/.test-temp/coverage-profile-test.*) ;;
    *) printf 'coverage-profile-test: fixture escaped its owned parent\n' >&2; exit 1 ;;
esac
trap 'rm -rf -- "$fixture"' EXIT HUP INT TERM

cat > "$fixture/go" <<'EOF'
#!/bin/sh
set -eu
[ "$1" = test ] || exit 2
printf '%s\n' "$COVERAGE_PROFILE" > "$PROFILE_RECORD"
if ! git check-ignore --quiet -- "$COVERAGE_PROFILE"; then
    printf 'coverage profile is an unignored mutable source input: %s\n' "$COVERAGE_PROFILE" >&2
    exit 1
fi
printf 'mode: set\n' > "$COVERAGE_PROFILE"
if [ "$PROFILE_TEST_MODE" = failure ]; then
    printf 'controlled coverage child failure\n' >&2
    exit 7
fi
for package in engine schedule timezone store catchup logbus; do
    printf 'github.com/shruggietech/go-schedule/v2/internal/%s/fixture.go:1.1,1.2 1 1\n' "$package" >> "$COVERAGE_PROFILE"
done
EOF
chmod +x "$fixture/go"

for mode in success failure; do
    record="$fixture/$mode.profile"
    status=0
    PROFILE_RECORD="$record" PROFILE_TEST_MODE="$mode" GO="$fixture/go" \
        sh "$ROOT/scripts/verify.sh" coverage > "$fixture/$mode.output" 2>&1 || status=$?
    expected=0
    if [ "$mode" = failure ]; then expected=7; fi
    if [ "$status" -ne "$expected" ]; then
        cat "$fixture/$mode.output" >&2
        printf 'coverage-profile-test: %s returned %s, expected %s\n' "$mode" "$status" "$expected" >&2
        exit 1
    fi
    profile=$(cat "$record")
    if [ -e "$ROOT/$profile" ]; then
        printf 'coverage-profile-test: %s left its mutable profile: %s\n' "$mode" "$profile" >&2
        exit 1
    fi
    if [ "$mode" = failure ] && ! grep -Fq 'controlled coverage child failure' "$fixture/$mode.output"; then
        printf 'coverage-profile-test: original child failure output was lost\n' >&2
        exit 1
    fi
done

printf 'coverage-profile-test: OK (ignored source boundary and success/failure cleanup)\n'
