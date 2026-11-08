#!/bin/sh
# Fixed adapter matrix; the policy corpus continues to reuse its compiled binary.
set -eu
ROOT=$(CDPATH='' cd -- "$(dirname "$0")/../.." && pwd)
. "$ROOT/test/scripts/lib/automation-fixture.sh"
fail() { printf 'public-entrypoints-test: FAIL: %s\n' "$*" >&2; exit 1; }
tmp=$(mktemp -d "$ROOT/.public-entrypoints-test.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
cases=0
invoke() {
  case "$selection" in
    absolute) sh "$ROOT/scripts/$command.sh" "$target" ;;
    relative) (cd "$tmp" && sh "$ROOT/scripts/$command.sh" "./$leaf") ;;
    default) (cd "$target" && sh "$ROOT/scripts/$command.sh") ;;
  esac
}
for leaf in ordinary 'root with spaces' 'root [1]' 'root [unclosed'; do
  target="$tmp/$leaf"
  make_fixture "$target" 'format vet lint race gui coverage docs automation'
  for command in automation-check spec-lifecycle-check; do
    for selection in absolute relative default; do
      if ! output=$(invoke 2>&1); then fail "$command $selection $leaf failed: $output"; fi
      case "$command:$output" in
        automation-check:*'automation-check: OK - '*) : ;;
        spec-lifecycle-check:*'spec-lifecycle-check: OK - 1 specification(s)'*) : ;;
        *) fail "$command $selection selected the wrong root: $output" ;;
      esac
      cases=$((cases + 1))
    done
  done
done
# A failing target distinguishes target selection from accidentally checking the
# source checkout, which also passes automation and has a generic success line.
leaf=wrong-lifecycle
target="$tmp/$leaf"
make_fixture "$target" 'format vet lint race gui coverage docs automation'
sed 's/Implemented/Finished/g' "$target/specs/001-fixture/spec.md" > "$tmp/invalid-spec"
cp "$tmp/invalid-spec" "$target/specs/001-fixture/spec.md"
for command in automation-check spec-lifecycle-check; do
  for selection in absolute relative default; do
    if output=$(invoke 2>&1); then fail "$command $selection ignored invalid selected target"; fi
    case "$output" in *'invalid status'*Finished*) : ;; *) fail "$command $selection lost target diagnostic: $output" ;; esac
    cases=$((cases + 1))
  done
done
printf 'public-entrypoints-test: OK (%s adapter cases)\n' "$cases"
