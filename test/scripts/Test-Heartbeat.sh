#!/usr/bin/env bash
# Record a heartbeat into a local SQLite database, so a scheduled task's firing
# behavior can be measured after the fact.
#
# POSIX twin of Test-Heartbeat.ps1. Same behavior, same recorded fields, same
# exit codes; every PowerShell -FooBar parameter is --foo-bar here.
#
# Writes one beat per invocation and exits. That default is deliberate: the
# point is to verify that goschedd supplies the cadence, and a script that
# sleeps in a loop tests its own sleep rather than the scheduler.
#
# The expected moment comes from one of three sources, in precedence order, and
# the source is recorded alongside every drift figure:
#   env       GOSCHED_SCHEDULED_TIME from the environment (not set by the
#             current executor; the tier exists so a future release is consumed
#             with no change here)
#   anchor    --anchor-iso plus --interval-seconds. One real firing time from
#             the schedule reconstructs the whole grid (anchor + k*interval), so
#             the difference is genuine dispatch latency. Get it from
#             `gosched task show <id>`.
#   none      neither available; no drift is recorded, and the reader excludes
#             the beat rather than reporting a meaningless zero
#
# There is deliberately no epoch-boundary tier. Snapping to the nearest multiple
# of the interval from the Unix epoch is correct only if the schedule sits on
# that grid, and this scheduler anchors interval schedules to task creation
# time -- so "every 1 minute" created at :06 fires at :06 forever, and epoch
# snapping reports a constant offset as though it were lateness.
#
# Exit codes: 0 success, 1 runtime failure, 2 usage error / unmet prerequisite.
# Full documentation: docs/test-scripts.md

set -euo pipefail
# shellcheck source-path=SCRIPTDIR
# shellcheck source=lib/sqlite.sh
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/sqlite.sh"

DATABASE="heartbeat"
INTERVAL_SECONDS=0
LABEL=""
ANCHOR_ISO=""
LOOP=0
MAX_BEATS=0
DURATION_SECONDS=0
SLEEP_SECONDS=0
FAIL_WITH=0
FAIL_WITH_SET=0
EXPLICIT_SQLITE=""
DO_INSTALL=0
LOOP_DEFAULT_DURATION_SECONDS=3600

usage() {
    sed -n '2,25p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
    cat <<'EOF'

Usage: Test-Heartbeat.sh [options]

  -d, --database NAME|PATH    'heartbeat' or an explicit path (default: heartbeat)
  -i, --interval-seconds N    declared schedule interval (gap detection; with
                              --anchor-iso also enables drift)
  -a, --anchor-iso TS         one firing time from the schedule, RFC 3339, from
                              `gosched task show`; with -i yields true drift
  -l, --label TEXT            tag recorded on each beat
  -r, --loop                  opt-in bounded continuous mode
  -m, --max-beats N           stop after N beats
  -t, --duration-seconds N    stop after N seconds (default 3600 under --loop)
  -s, --sleep-seconds N       deliberately extend the run (overlap testing)
  -f, --fail-with N           exit N after recording (rejects 0 and 2)
      --sqlite-exe PATH       explicit sqlite3 path
      --install-sqlite        download+verify the pinned sqlite3
  -q, --quiet                 suppress informational output
  -h, --help                  this text
EOF
}

while [ $# -gt 0 ]; do
    case "$1" in
        -d|--database)         DATABASE="$2"; shift 2 ;;
        -i|--interval-seconds) INTERVAL_SECONDS="$2"; shift 2 ;;
        -l|--label)            LABEL="$2"; shift 2 ;;
        -a|--anchor-iso)       ANCHOR_ISO="$2"; shift 2 ;;
        -r|--loop)             LOOP=1; shift ;;
        -m|--max-beats)        MAX_BEATS="$2"; shift 2 ;;
        -t|--duration-seconds) DURATION_SECONDS="$2"; shift 2 ;;
        -s|--sleep-seconds)    SLEEP_SECONDS="$2"; shift 2 ;;
        -f|--fail-with)        FAIL_WITH="$2"; FAIL_WITH_SET=1; shift 2 ;;
        --sqlite-exe)          EXPLICIT_SQLITE="$2"; shift 2 ;;
        --install-sqlite)      DO_INSTALL=1; shift ;;
        -q|--quiet)            LOG_QUIET=1; shift ;;
        -h|--help)             usage; exit 0 ;;
        *) die_usage "Unknown option: $1 (try --help)" ;;
    esac
done

# Reserved codes. An induced failure that could return 0 or 2 would be
# indistinguishable from a success or from a missing sqlite3.
if [ "$FAIL_WITH_SET" -eq 1 ] && { [ "$FAIL_WITH" -eq 0 ] || [ "$FAIL_WITH" -eq 2 ]; }; then
    die_usage "--fail-with must not be 0 or 2; those codes are reserved for success and for unmet prerequisites."
fi
for n in "$INTERVAL_SECONDS" "$MAX_BEATS" "$DURATION_SECONDS" "$SLEEP_SECONDS"; do
    case "$n" in ''|*[!0-9]*) die_usage "Numeric options must be non-negative integers." ;; esac
done

init_test_sqlite "$EXPLICIT_SQLITE" "$DO_INSTALL"
DB_PATH="$(resolve_test_database "$DATABASE")"
log INFO "database: $DB_PATH"
init_heartbeat_schema "$DB_PATH"

SESSION_ID="$(od -An -N16 -tx1 /dev/urandom 2>/dev/null | tr -d ' \n' || date +%s%N)"

resolve_expected() {
    # Echoes "SOURCE EXPECTED_MS DRIFT_MS"; EXPECTED/DRIFT are __NULL__ when
    # no source is available.
    local started="$1" expected drift interval_ms anchor_ms delta k
    if [ -n "${GOSCHED_SCHEDULED_TIME:-}" ]; then
        expected="$(date -d "$GOSCHED_SCHEDULED_TIME" +%s000 2>/dev/null || true)"
        if [ -n "$expected" ]; then
            printf 'env %s %s\n' "$expected" "$((started - expected))"
            return 0
        fi
        log WARN "GOSCHED_SCHEDULED_TIME set but unparseable; ignoring it."
    fi
    if [ -n "$ANCHOR_ISO" ] && [ "$INTERVAL_SECONDS" -gt 0 ]; then
        anchor_ms="$(date -d "$ANCHOR_ISO" +%s 2>/dev/null || true)"
        if [ -z "$anchor_ms" ]; then
            log ERROR "--anchor-iso '$ANCHOR_ISO' is not a parseable timestamp."
            exit 2
        fi
        anchor_ms=$((anchor_ms * 1000))
        interval_ms=$((INTERVAL_SECONDS * 1000))
        # Any firing time works as the anchor: they all sit on the same
        # anchor + k*interval grid, k signed. Round half away from zero.
        delta=$((started - anchor_ms))
        if [ "$delta" -ge 0 ]; then
            k=$(( (delta + interval_ms / 2) / interval_ms ))
        else
            k=$(( (delta - interval_ms / 2) / interval_ms ))
        fi
        expected=$((anchor_ms + k * interval_ms))
        printf 'anchor %s %s\n' "$expected" "$((started - expected))"
        return 0
    fi
    printf 'none __NULL__ __NULL__\n'
}

write_beat() {
    local seq="$1" started="$2" finished="$3" code="$4"
    local outcome source expected drift started_iso
    read -r source expected drift <<<"$(resolve_expected "$started")"
    outcome=ok; [ "$code" -eq 0 ] || outcome=failed
    started_iso="$(iso_local)"

    sqlite_exec "$DB_PATH" list "INSERT INTO beat (session_id, sequence, label, hostname, username, pid,
       started_ms, started_iso, finished_ms, duration_ms, expected_ms, expected_source,
       drift_ms, interval_seconds, exit_code, outcome, sched_env)
VALUES (:session, CAST(:seq AS INTEGER),
        CASE WHEN :label = '' THEN NULL ELSE :label END,
        :host, :user, CAST(:pid AS INTEGER),
        CAST(:started AS INTEGER), :startediso, CAST(:finished AS INTEGER),
        CAST(:duration AS INTEGER),
        CASE WHEN :expected IS NULL THEN NULL ELSE CAST(:expected AS INTEGER) END,
        :source,
        CASE WHEN :drift IS NULL THEN NULL ELSE CAST(:drift AS INTEGER) END,
        CASE WHEN :interval = '0' THEN NULL ELSE CAST(:interval AS INTEGER) END,
        CAST(:code AS INTEGER), :outcome, '{}');" \
        "session=$SESSION_ID" "seq=$seq" "label=$LABEL" \
        "host=$(hostname)" "user=$(id -un 2>/dev/null || printf '%s' "${USER:-unknown}")" \
        "pid=$$" "started=$started" "startediso=$started_iso" \
        "finished=$finished" "duration=$((finished - started))" \
        "expected=$expected" "source=$source" "drift=$drift" \
        "interval=$INTERVAL_SECONDS" "code=$code" "outcome=$outcome" >/dev/null

    if [ "$drift" = "__NULL__" ]; then
        log SUCCESS "beat ${SESSION_ID:0:8} seq $seq recorded; drift n/a (no expected moment available)"
    else
        log SUCCESS "beat ${SESSION_ID:0:8} seq $seq recorded; drift ${drift}ms (source: $source)"
    fi
}

if [ "$LOOP" -eq 0 ]; then
    # Default path: exactly one beat. The scheduler owns the cadence.
    STARTED="$(now_ms)"
    [ "$SLEEP_SECONDS" -gt 0 ] && sleep "$SLEEP_SECONDS"
    write_beat 1 "$STARTED" "$(now_ms)" "$FAIL_WITH"
    printf '%s\n' "$DB_PATH"
    exit "$FAIL_WITH"
fi

# Opt-in continuous mode. Bounded always: an unbounded loop launched under a
# scheduler is a resource incident, not a test.
EFFECTIVE_DURATION=0
if [ "$DURATION_SECONDS" -gt 0 ]; then
    EFFECTIVE_DURATION="$DURATION_SECONDS"
elif [ "$MAX_BEATS" -eq 0 ]; then
    EFFECTIVE_DURATION="$LOOP_DEFAULT_DURATION_SECONDS"
fi
CADENCE=1
[ "$INTERVAL_SECONDS" -gt 0 ] && CADENCE="$INTERVAL_SECONDS"
LOOP_START="$(now_ms)"
SEQUENCE=0
log INFO "loop mode: max beats ${MAX_BEATS:-unset}, duration ${EFFECTIVE_DURATION}s, cadence ${CADENCE}s"

while true; do
    # Both bounds checked between beats; first to trip ends the loop.
    [ "$MAX_BEATS" -gt 0 ] && [ "$SEQUENCE" -ge "$MAX_BEATS" ] && break
    if [ "$EFFECTIVE_DURATION" -gt 0 ] && \
       [ $(( $(now_ms) - LOOP_START )) -ge $((EFFECTIVE_DURATION * 1000)) ]; then
        break
    fi
    SEQUENCE=$((SEQUENCE + 1))
    STARTED="$(now_ms)"
    [ "$SLEEP_SECONDS" -gt 0 ] && sleep "$SLEEP_SECONDS"
    write_beat "$SEQUENCE" "$STARTED" "$(now_ms)" 0
    [ "$MAX_BEATS" -gt 0 ] && [ "$SEQUENCE" -ge "$MAX_BEATS" ] && break
    sleep "$CADENCE"
done

log SUCCESS "loop finished after $SEQUENCE beat(s)"
printf '%s\n' "$DB_PATH"
exit "$FAIL_WITH"
