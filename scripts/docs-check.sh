#!/bin/sh
# docs-check.sh — documentation integrity gate for the go-schedule docs site.
#
# For every docs/*.md page it asserts:
#   1. a YAML front-matter block (--- … ---) with `title:` and `nav_order:`;
#   2. every non-http(s), non-fragment Markdown link resolves on disk; and
#   3. no link escapes the docs/ directory (content outside docs/ must be an
#      absolute https://github.com/… URL, which is skipped).
# It also checks that the pointer README(s) reference an existing docs/ page.
#
# Pure POSIX sh + coreutils: no network, no Ruby, no build. Anchors (#frag) are
# stripped, not validated. Contract:
# specs/010-docs-site-pages/contracts/docs-check.md
set -eu

DOCS_DIR="docs"
POINTERS="test/scripts/README.md"

# Failures accumulate here, one `file: reason: detail` line each. Counting lines
# afterwards avoids the subshell-scoping trap of piping into a while-loop.
FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT

report() { printf '%s: %s: %s\n' "$1" "$2" "$3" >> "$FAILURES"; }

# normalize <path> — collapse . and .. segments; print the cleaned path.
normalize() {
  oldIFS=$IFS
  IFS=/
  # shellcheck disable=SC2086
  set -- $1
  IFS=$oldIFS
  out=""
  for part in "$@"; do
    case "$part" in
      "" | .) : ;;
      ..) out=${out%/*} ;;
      *) out="$out/$part" ;;
    esac
  done
  printf '%s' "${out#/}"
}

# links_in <file> — print each Markdown link target on its own line.
links_in() {
  grep -oE '\]\([^)]+\)' "$1" 2>/dev/null | sed 's/^](//; s/)$//' || true
}

# check_link <file> <target> <enforce_no_escape 0|1>
check_link() {
  target=${2%%#*}      # strip anchor
  target=${target%% *} # strip optional Jekyll link title
  case "$2" in
    http://* | https://* | mailto:* | "#"*) return 0 ;;
  esac
  [ -n "$target" ] || return 0

  resolved=$(normalize "$(dirname "$1")/$target")
  if [ "$3" = "1" ]; then
    case "$resolved" in
      "$DOCS_DIR"/*) : ;;
      *) report "$1" "link escapes $DOCS_DIR/" "$2"; return 0 ;;
    esac
  fi
  [ -e "$resolved" ] || report "$1" "broken link" "$2"
}

# check_frontmatter <file>
check_frontmatter() {
  if [ "$(sed -n '1p' "$1")" != "---" ]; then
    report "$1" "missing front matter" "no opening --- on line 1"
    return 0
  fi
  close=$(awk 'NR>1 && $0=="---"{print NR; exit}' "$1")
  if [ -z "$close" ]; then
    report "$1" "missing front matter" "no closing ---"
    return 0
  fi
  fm=$(sed -n "2,$((close - 1))p" "$1")
  printf '%s\n' "$fm" | grep -q '^title:' || report "$1" "front matter missing key" "title"
  printf '%s\n' "$fm" | grep -q '^nav_order:' || report "$1" "front matter missing key" "nav_order"
}

page_count=0
for f in "$DOCS_DIR"/*.md; do
  [ -e "$f" ] || continue
  page_count=$((page_count + 1))
  check_frontmatter "$f"
  # No pipe here: the for-loop keeps report()'s writes in this shell.
  for target in $(links_in "$f"); do
    check_link "$f" "$target" 1
  done
done

for p in $POINTERS; do
  if [ ! -e "$p" ]; then
    report "$p" "missing pointer README" "$p"
    continue
  fi
  for target in $(links_in "$p"); do
    check_link "$p" "$target" 0
  done
done

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'docs-check: FAILED with %s issue(s) across %s page(s)\n' \
    "$(wc -l < "$FAILURES" | tr -d ' ')" "$page_count" >&2
  exit 1
fi
printf 'docs-check: OK — %s pages, links and front matter clean\n' "$page_count"
