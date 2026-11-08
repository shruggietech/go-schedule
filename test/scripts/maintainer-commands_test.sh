#!/bin/sh
# Shared focused corpus for canonical automation and the existing native matrix.
set -eu
ROOT=$(CDPATH='' cd -- "$(dirname "$0")/../.." && pwd)
sh "$ROOT/test/scripts/spec-lifecycle-check_test.sh"
sh "$ROOT/test/scripts/public-entrypoints_test.sh"
sh "$ROOT/test/scripts/verify-driver_test.sh"
printf 'maintainer-commands-test: OK\n'
