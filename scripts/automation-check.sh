#!/bin/sh
# Evaluate existing automation policy through the native source-checkout helper.
set -eu
if [ "$#" -gt 1 ]; then
  printf 'usage: %s [root]\n' "$0" >&2
  exit 2
fi
TARGET=${1:-.}
case "$TARGET" in
  /* | [A-Za-z]:*) ;;
  *) TARGET="$PWD/$TARGET" ;;
esac
SOURCE=$(CDPATH='' cd "$(dirname "$0")/.." && pwd)
cd "$SOURCE"
exec "${GO:-go}" run ./scripts/maintainer-check automation "$TARGET"
