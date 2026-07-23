#!/bin/sh
# Core-package statement coverage gate.
#
# This is the single implementation of the gate: CI runs it, and it is one of
# the CI-parity commands in CLAUDE.md that must be green before a push. Keeping
# one implementation is deliberate — when the local check and the CI check were
# two different measurements, a push went out that CI then rejected.
#
# Written in POSIX sh + awk rather than python so it runs unchanged in Git Bash
# on Windows, in WSL, and on the CI runner.
#
# Usage: scripts/coverage-gate.sh [threshold]      (default 80)

set -eu

THRESHOLD="${1:-80}"
CORE="engine schedule timezone store catchup logbus"
PROFILE="${COVERAGE_PROFILE:-cover.out}"

COVERPKG="./internal/engine,./internal/schedule,./internal/timezone,./internal/store,./internal/catchup,./internal/logbus"

# -count=1 defeats the test cache. Under -coverpkg every test binary is
# instrumented for all six target packages, so a cached result replays a
# coverage profile enumerating the file set as it was when that result was
# cached. A package whose own sources did not change is served from cache and
# drags stale blocks — including blocks for files that have since been deleted —
# into the merged profile, inflating the denominator and failing the gate for
# code that no longer exists.
go test -count=1 -coverpkg="$COVERPKG" ./... -coverprofile="$PROFILE" >/dev/null

# A block key can appear once per test binary in the merged profile, so blocks
# are deduplicated by key and counted as covered if any binary reached them.
# Packages are keyed by the first path segment after "internal/".
awk -v core="$CORE" -v threshold="$THRESHOLD" '
NR > 1 && NF >= 3 {
    key = $1; n = $2 + 0; c = $3 + 0
    if (!(key in seen)) { seen[key] = 1; stmts[key] = n; hits[key] = 0 }
    if (c > hits[key]) hits[key] = c
}
END {
    for (k in seen) {
        split(k, a, "internal/")
        if (a[2] == "") continue
        split(a[2], b, "/")
        p = b[1]
        total[p] += stmts[k]
        if (hits[k] > 0) covered[p] += stmts[k]
    }
    failed = 0
    n = split(core, list, " ")
    for (i = 1; i <= n; i++) {
        p = list[i]
        pct = total[p] ? 100 * covered[p] / total[p] : 0
        status = (pct >= threshold) ? "OK" : "FAIL"
        if (pct < threshold) failed = 1
        printf "%-10s%6.1f%%  %s\n", p, pct, status
    }
    exit failed
}
' "$PROFILE"
