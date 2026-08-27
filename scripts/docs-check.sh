#!/bin/sh
# docs-check.sh — documentation integrity gate for the go-schedule docs site.
#
# For every docs/*.md page it asserts:
#   1. a YAML front-matter block (--- … ---) with `title:` and `nav_order:`;
#   2. every non-http(s), non-fragment Markdown link resolves on disk; and
#   3. no link escapes the docs/ directory (content outside docs/ must be an
#      absolute https://github.com/… URL, which is skipped).
#   4. the custom syntax theme provides a complete, safe dark-palette contract.
#   5. every published fenced block declares an approved content category.
# It also checks that the pointer README(s) reference an existing docs/ page.
#
# Pure POSIX sh + coreutils: no network, no Ruby, no build. Anchors (#frag) are
# stripped, not validated. Contract:
# specs/010-docs-site-pages/contracts/docs-check.md
set -eu

DOCS_DIR="docs"
POINTERS="test/scripts/README.md"
THEME_SCSS="$DOCS_DIR/_sass/custom/custom.scss"

# Failures accumulate here, one `file: reason: detail` line each. Counting lines
# afterwards avoids the subshell-scoping trap of piping into a while-loop.
FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT

report() { printf '%s: %s: %s\n' "$1" "$2" "$3" >> "$FAILURES"; }

# require_style <extended-regex> <contract-description>
require_style() {
  if ! grep -Eq "$1" "$THEME_SCSS"; then
    report "$THEME_SCSS" "missing dark-theme contract" "$2"
  fi
}

check_theme_contract() {
  if [ ! -f "$THEME_SCSS" ]; then
    report "$THEME_SCSS" "missing stylesheet" "$THEME_SCSS"
    return 0
  fi

  require_style '\[class\][[:space:]]*\{[[:space:]]*color:[[:space:]]*inherit;' \
    "safe fallback for every classified token"
  require_style '\.n,[[:space:]]*\.nx[[:space:]]*\{[[:space:]]*color:[[:space:]]*#58a6ff;' \
    "name tokens use Anchor Blue"
  require_style '\.l,[[:space:]]*\.ld[[:space:]]*\{[[:space:]]*color:[[:space:]]*#f2b84b;' \
    "literal tokens use Hold Amber"
  require_style '\.nd,[[:space:]]*\.ne,[[:space:]]*\.ni,[[:space:]]*\.nl,[[:space:]]*\.py,[[:space:]]*\.gu[[:space:]]*\{[[:space:]]*color:[[:space:]]*#58a6ff;' \
    "decorator, exception, label, and heading tokens use Anchor Blue"
  require_style '\.hll[[:space:]]*\{[[:space:]]*background-color:[[:space:]]*#17262d;' \
    "highlighted lines use a dark surface"
  require_style '::selection[[:space:]]*\{[[:space:]]*color:[[:space:]]*#0d171c;[[:space:]]*background-color:[[:space:]]*#58a6ff;' \
    "selected code uses Panel ink on Anchor Blue"
  # These patterns intentionally match literal SCSS variables.
  # shellcheck disable=SC2016
  require_style 'padding:[[:space:]]*\$sp-3[[:space:]]+\$gutter-spacing-sm;' \
    "endorsement uses vertical spacing and the small navigation gutter"
  # shellcheck disable=SC2016
  require_style '@include[[:space:]]+mq\(md\)[[:space:]]*\{[[:space:]]*padding-right:[[:space:]]*\$gutter-spacing;[[:space:]]*padding-left:[[:space:]]*\$gutter-spacing;' \
    "endorsement uses the desktop navigation gutter"
}

# check_fences <file> — require one of the documented fence categories.
check_fences() {
  awk '
    /^```/ {
      if (!open) {
        language = substr($0, 4)
        if (language !~ /^(sh|bash|powershell|text)$/) {
          detail = language == "" ? "untagged opening fence" : "unsupported fence: " language
          print NR "|" detail
        }
        open = 1
        next
      }
      if ($0 != "```") {
        print NR "|closing fence must be plain triple backticks"
      }
      open = 0
    }
    END {
      if (open) print NR "|unclosed fenced block"
    }
  ' "$1" | while IFS='|' read -r line detail; do
    report "$1:$line" "invalid code fence" "$detail"
  done
}

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
check_theme_contract
for f in "$DOCS_DIR"/*.md; do
  [ -e "$f" ] || continue
  page_count=$((page_count + 1))
  check_frontmatter "$f"
  check_fences "$f"
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
printf 'docs-check: OK — %s pages, links, front matter, fences, and theme contract clean\n' "$page_count"
