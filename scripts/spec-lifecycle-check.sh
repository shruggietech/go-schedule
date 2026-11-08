#!/bin/sh
# Evaluate existing lifecycle policy through the native source-checkout helper.
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
# Git for Windows may leave glob characters unconverted in native argv.
# Resolve the directory literally with shell builtins before starting Go.
if [ "${OS:-}" = Windows_NT ]; then
  TARGET=$(CDPATH='' cd -- "$TARGET" && pwd -W)
fi
SOURCE=$(CDPATH='' cd "$(dirname "$0")/.." && pwd)
cd "$SOURCE"
exec "${GO:-go}" run ./scripts/maintainer-check lifecycle "$TARGET"
