#!/usr/bin/env bash
# Run canned inspection queries against the go-schedule test databases.
#
# POSIX twin of Test-ReadTestDB.ps1. Same queries, same reporting rules, same
# exit codes; every PowerShell -FooBar parameter is --foo-bar here.
#
# Two reporting rules are contract rather than presentation, because breaking
# them would produce confident numbers that mean less than they appear to:
#   - A query that excludes rows says how many it excluded and why. A percentile
#     over an unstated subset is a number drawn from unknown evidence.
#   - Drift is never pooled across expected-source values and never shown
#     without its source. A measured value and a derived one are different kinds
#     of number, and averaging them produces neither.
#
# Exit codes: 0 success, 1 runtime failure, 2 usage error / unmet prerequisite.
# Full documentation: docs/test-scripts.md

set -euo pipefail
# shellcheck source-path=SCRIPTDIR
# shellcheck source=lib/sqlite.sh
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/sqlite.sh"

DATABASE="heartbeat"
QUERY="summary"
DO_LIST=0
FORMAT="Table"
LIMIT=20
INTERVAL_SECONDS=0
EXPLICIT_SQLITE=""
DO_INSTALL=0

usage() {
    cat <<'EOF'
Test-ReadTestDB.sh -- canned inspection queries over the test databases.

  -d, --database NAME|PATH   'heartbeat', 'system', or a path (default: heartbeat)
  -k, --query NAME           which canned query (default: summary)
  -n, --list                 list the available queries and exit
  -f, --format FORM          Table | Json | Csv (default: Table)
  -m, --limit N              row cap for row-returning queries (default: 20)
  -i, --interval-seconds N   expected interval for gap and reliability checks;
                             inferred when omitted, and the output says so
      --sqlite-exe PATH      explicit sqlite3 path
      --install-sqlite       download+verify the pinned sqlite3
  -q, --quiet                suppress informational output
  -h, --help                 this text
EOF
}

list_queries() {
    cat <<'EOF'
Available queries (--query NAME):

  cadence    [heartbeat] What were the observed intervals?
  drift      [heartbeat] How far from expected did firings land, by expected-source?
  failures   [heartbeat] Which runs reported failure?
  gaps       [heartbeat] Which expected firings were missed or badly delayed?
  hosts      [both     ] Which hosts and users produced records?
  listeners  [system   ] What is listening now, and what changed since the previous snapshot?
  overlaps   [heartbeat] Which runs overlapped in time?
  recent     [both     ] What are the most recent records?
  restarts   [heartbeat] Where are the session boundaries?
  schema     [both     ] What is the stored structure?
  summary    [both     ] How many records, over what period, from how many sessions and hosts?

Databases: 'heartbeat', 'system', or an explicit path.
EOF
}

while [ $# -gt 0 ]; do
    case "$1" in
        -d|--database)         DATABASE="$2"; shift 2 ;;
        -k|--query)            QUERY="$2"; shift 2 ;;
        -n|--list)             DO_LIST=1; shift ;;
        -f|--format)           FORMAT="$2"; shift 2 ;;
        -m|--limit)            LIMIT="$2"; shift 2 ;;
        -i|--interval-seconds) INTERVAL_SECONDS="$2"; shift 2 ;;
        --sqlite-exe)          EXPLICIT_SQLITE="$2"; shift 2 ;;
        --install-sqlite)      DO_INSTALL=1; shift ;;
        -q|--quiet)            LOG_QUIET=1; shift ;;
        -h|--help)             usage; exit 0 ;;
        *) die_usage "Unknown option: $1 (try --help)" ;;
    esac
done

[ "$DO_LIST" -eq 1 ] && { list_queries; exit 0; }

case "$QUERY" in
    summary|recent|cadence|drift|gaps|overlaps|failures|restarts|hosts|listeners|schema) ;;
    *) die_usage "Unknown query '$QUERY'. Use --list to see the available queries." ;;
esac
case "$LIMIT" in ''|*[!0-9]*) die_usage "--limit must be a positive integer." ;; esac
[ "$LIMIT" -gt 0 ] || die_usage "--limit must be a positive integer."

init_test_sqlite "$EXPLICIT_SQLITE" "$DO_INSTALL"
DB_PATH="$(resolve_test_database "$DATABASE")"
[ -f "$DB_PATH" ] || die_usage "No database at $DB_PATH. Run a recording script first."
log INFO "database: $DB_PATH"

IS_HEARTBEAT=0
case "$QUERY" in cadence|drift|gaps|overlaps|failures|restarts) IS_HEARTBEAT=1 ;; esac
[ "$DATABASE" = "heartbeat" ] && IS_HEARTBEAT=1

case "$FORMAT" in
    Json|json) MODE=json ;;
    Csv|csv)   MODE=csv ;;
    Table|table) MODE=box ;;
    *) die_usage "--format must be Table, Json, or Csv." ;;
esac

# Interval resolution, and saying which it was. An inferred interval is an
# assumption, and assumptions that look like measurements are how wrong
# conclusions get confident.
INTERVAL_USED="$INTERVAL_SECONDS"
INTERVAL_INFERRED=0
if [ "$INTERVAL_USED" -le 0 ]; then
    case "$QUERY" in
        gaps|drift)
            INTERVAL_USED="$(sqlite_exec "$DB_PATH" list \
"SELECT CAST(ROUND(delta/1000.0) AS INTEGER) FROM
 (SELECT started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta FROM beat)
 WHERE delta IS NOT NULL AND delta > 0
 GROUP BY 1 ORDER BY COUNT(*) DESC, 1 ASC LIMIT 1;" | tail -n1)"
            case "$INTERVAL_USED" in ''|*[!0-9]*) INTERVAL_USED=0 ;; esac
            INTERVAL_INFERRED=1
            if [ "$INTERVAL_USED" -gt 0 ]; then
                log WARN "interval not supplied; inferred ${INTERVAL_USED}s from the most common observed interval"
            else
                log WARN "interval not supplied and could not be inferred (too few records)"
            fi
            ;;
    esac
fi

report_excluded() {
    # A query that drops rows must say so. Silence reads as "nothing was
    # dropped", which is a different and stronger claim.
    local count_sql="$1" reason="$2" n
    n="$(sqlite_exec "$DB_PATH" list "$count_sql" | tail -n1)"
    case "$n" in ''|*[!0-9]*) return 0 ;; esac
    [ "$n" -gt 0 ] && log WARN "excluded $n record(s): $reason"
    return 0
}

if [ "$QUERY" = "drift" ]; then
    report_excluded "SELECT COUNT(*) FROM beat WHERE drift_ms IS NULL;" \
        "no expected moment was available, so drift is not computable for them"
    if [ "$INTERVAL_USED" -gt 0 ]; then
        report_excluded \
            "SELECT COUNT(*) FROM beat WHERE expected_source='boundary' AND ABS(drift_ms) > $((INTERVAL_USED * 250));" \
            "boundary-derived drift exceeds a quarter of the interval and is NOT reliable at that magnitude"
    fi
fi
if [ "$QUERY" = "gaps" ] && [ "$INTERVAL_INFERRED" -eq 1 ]; then
    log WARN "the interval above was inferred, not supplied -- treat the result accordingly"
fi

case "$QUERY" in
    summary)
        if [ "$IS_HEARTBEAT" -eq 1 ]; then
            SQL="SELECT COUNT(*) AS beats, COUNT(DISTINCT session_id) AS sessions,
 COUNT(DISTINCT hostname) AS hosts, MIN(started_iso) AS first_seen,
 datetime(MAX(started_ms)/1000,'unixepoch') AS last_seen_utc,
 ROUND((MAX(started_ms)-MIN(started_ms))/60000.0,1) AS span_minutes,
 SUM(CASE WHEN outcome='failed' THEN 1 ELSE 0 END) AS failures FROM beat;"
        else
            SQL="SELECT COUNT(*) AS snapshots, COUNT(DISTINCT hostname) AS hosts,
 MIN(iso_local) AS first_seen, MAX(iso_local) AS last_seen,
 ROUND((MAX(unixtime_ms)-MIN(unixtime_ms))/60000.0,1) AS span_minutes FROM snapshot;"
        fi ;;
    recent)
        if [ "$IS_HEARTBEAT" -eq 1 ]; then
            SQL="SELECT id, started_iso, label, sequence, duration_ms, expected_source, drift_ms, outcome
 FROM beat ORDER BY started_ms DESC LIMIT $LIMIT;"
        else
            SQL="SELECT id, iso_local, hostname, username, process_count, uptime_seconds, invocation_source
 FROM snapshot ORDER BY unixtime_ms DESC LIMIT $LIMIT;"
        fi ;;
    cadence)
        SQL="WITH d AS (SELECT started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta FROM beat)
 SELECT COUNT(*) AS intervals, MIN(delta) AS min_ms, CAST(AVG(delta) AS INTEGER) AS mean_ms,
 MAX(delta) AS max_ms FROM d WHERE delta IS NOT NULL;" ;;
    drift)
        SQL="SELECT expected_source, COUNT(*) AS n, MIN(drift_ms) AS min_ms,
 CAST(AVG(drift_ms) AS INTEGER) AS mean_ms, MAX(ABS(drift_ms)) AS max_abs_ms
 FROM beat WHERE drift_ms IS NOT NULL GROUP BY expected_source ORDER BY expected_source;" ;;
    gaps)
        [ "$INTERVAL_USED" -gt 0 ] || die_usage "cannot detect gaps without an interval; pass --interval-seconds"
        SQL="WITH d AS (SELECT started_iso, started_ms,
 started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta FROM beat)
 SELECT started_iso AS resumed_at, delta AS gap_ms, ROUND(delta/1000.0,1) AS gap_seconds,
 CAST(ROUND(delta/1000.0/$INTERVAL_USED) AS INTEGER) AS intervals_missed
 FROM d WHERE delta IS NOT NULL AND delta > $((INTERVAL_USED * 2000))
 ORDER BY started_ms LIMIT $LIMIT;" ;;
    overlaps)
        SQL="SELECT a.id AS run_a, b.id AS run_b, a.started_iso AS a_started,
 b.started_iso AS b_started, a.duration_ms AS a_duration_ms, b.duration_ms AS b_duration_ms
 FROM beat a JOIN beat b ON a.id < b.id AND a.started_ms <= b.finished_ms
 AND b.started_ms <= a.finished_ms ORDER BY a.started_ms LIMIT $LIMIT;" ;;
    failures)
        SQL="SELECT id, started_iso, label, exit_code, outcome, duration_ms
 FROM beat WHERE exit_code <> 0 ORDER BY started_ms DESC LIMIT $LIMIT;" ;;
    restarts)
        SQL="SELECT session_id, pid, COUNT(*) AS beats, MIN(started_iso) AS first_beat,
 datetime(MAX(started_ms)/1000,'unixepoch') AS last_beat_utc
 FROM beat GROUP BY session_id, pid ORDER BY MIN(started_ms) DESC LIMIT $LIMIT;" ;;
    hosts)
        if [ "$IS_HEARTBEAT" -eq 1 ]; then
            SQL="SELECT hostname, username, COUNT(*) AS beats FROM beat
 GROUP BY hostname, username ORDER BY beats DESC LIMIT $LIMIT;"
        else
            SQL="SELECT hostname, username, COUNT(*) AS snapshots FROM snapshot
 GROUP BY hostname, username ORDER BY snapshots DESC LIMIT $LIMIT;"
        fi ;;
    listeners)
        SQL="WITH latest AS (SELECT id FROM snapshot ORDER BY unixtime_ms DESC LIMIT 1),
 prev AS (SELECT id FROM snapshot ORDER BY unixtime_ms DESC LIMIT 1 OFFSET 1)
 SELECT p.protocol, p.port, p.address, p.process_name,
 CASE WHEN (SELECT id FROM prev) IS NULL THEN 'no-previous-snapshot'
      WHEN EXISTS (SELECT 1 FROM snapshot_port q WHERE q.snapshot_id=(SELECT id FROM prev)
                   AND q.protocol=p.protocol AND q.port=p.port)
      THEN 'unchanged' ELSE 'NEW' END AS since_previous
 FROM snapshot_port p WHERE p.snapshot_id=(SELECT id FROM latest)
 ORDER BY p.protocol, p.port LIMIT $LIMIT;" ;;
    schema)
        SQL="SELECT sql FROM sqlite_master WHERE sql IS NOT NULL ORDER BY type DESC, name;" ;;
esac

sqlite_exec "$DB_PATH" "$MODE" "$SQL"
exit 0
