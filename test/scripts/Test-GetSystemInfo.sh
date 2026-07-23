#!/usr/bin/env bash
# Record a snapshot of this machine into a local SQLite database.
#
# POSIX twin of Test-GetSystemInfo.ps1. Same recorded fields, same exit codes;
# every PowerShell -FooBar parameter is --foo-bar here.
#
# Captures host identity, logged-in user, running process count, uptime, network
# addresses, and listening ports as one snapshot with attached detail rows. As a
# scheduler test this exercises subprocess spawning, platform tooling, and
# multi-row writes -- where cross-platform execution bugs actually surface.
#
# Probes degrade rather than abort: one that cannot run records NULL and warns,
# and the snapshot is still written. NULL means "could not determine" and is
# never used for a legitimate zero, because a process count of zero and an
# unavailable process count support opposite conclusions.
#
# Exit codes: 0 the snapshot was recorded (even with degraded probes), 1 the
# snapshot could not be written, 2 usage error / unmet prerequisite.
# Full documentation: docs/test-scripts.md

set -euo pipefail
# shellcheck source-path=SCRIPTDIR
# shellcheck source=lib/sqlite.sh
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/sqlite.sh"

DATABASE="system"
INVOCATION_SOURCE="manual"
SKIP_PORTS=0
EXPLICIT_SQLITE=""
DO_INSTALL=0

usage() {
    cat <<'EOF'
Test-GetSystemInfo.sh -- record a snapshot of this machine.

  -d, --database NAME|PATH       'system' or an explicit path (default: system)
  -i, --invocation-source TEXT   tag recorded on the snapshot (default: manual)
  -s, --skip-ports               skip the listening-port probe
      --sqlite-exe PATH          explicit sqlite3 path
      --install-sqlite           download+verify the pinned sqlite3
  -q, --quiet                    suppress informational output
  -h, --help                     this text
EOF
}

while [ $# -gt 0 ]; do
    case "$1" in
        -d|--database)          DATABASE="$2"; shift 2 ;;
        -i|--invocation-source) INVOCATION_SOURCE="$2"; shift 2 ;;
        -s|--skip-ports)        SKIP_PORTS=1; shift ;;
        --sqlite-exe)           EXPLICIT_SQLITE="$2"; shift 2 ;;
        --install-sqlite)       DO_INSTALL=1; shift ;;
        -q|--quiet)             LOG_QUIET=1; shift ;;
        -h|--help)              usage; exit 0 ;;
        *) die_usage "Unknown option: $1 (try --help)" ;;
    esac
done

init_test_sqlite "$EXPLICIT_SQLITE" "$DO_INSTALL"
DB_PATH="$(resolve_test_database "$DATABASE")"
log INFO "database: $DB_PATH"
init_system_schema "$DB_PATH"

#_______________________________________________________________________________
# Probes -- fixed, documented fallback order. An ordering that varies by machine
# produces data whose provenance differs between hosts while looking identical
# in the database.

probe_process_count() {
    local n
    n="$(ps -e 2>/dev/null | wc -l | tr -d ' ')" || true
    if [ -n "$n" ] && [ "$n" -gt 0 ] 2>/dev/null; then
        printf '%s\n' "$((n - 1))"
    else
        log WARN "process count unavailable"
        printf '\n'
    fi
}

probe_uptime_seconds() {
    local raw boot
    if [ -r /proc/uptime ]; then
        raw="$(cut -d' ' -f1 /proc/uptime)"
        printf '%s\n' "${raw%%.*}"
        return 0
    fi
    boot="$(sysctl -n kern.boottime 2>/dev/null | sed -n 's/.*sec = \([0-9]*\).*/\1/p')"
    if [ -n "$boot" ]; then
        printf '%s\n' "$(( $(date +%s) - boot ))"
        return 0
    fi
    log WARN "uptime unavailable"
    printf '\n'
}

probe_addresses() {
    # Order: ip -o addr, then ifconfig. Emits "family|address|interface".
    local out
    out="$(ip -o addr 2>/dev/null || true)"
    if [ -n "$out" ]; then
        printf '%s\n' "$out" | awk '{
            for (i=1;i<=NF;i++) {
                if ($i=="inet")  { split($(i+1),a,"/"); print "ipv4|" a[1] "|" $2 }
                if ($i=="inet6") { split($(i+1),a,"/"); print "ipv6|" a[1] "|" $2 }
            }}'
        return 0
    fi
    out="$(ifconfig 2>/dev/null || true)"
    if [ -n "$out" ]; then
        printf '%s\n' "$out" | awk '
            /^[a-zA-Z0-9]/ { iface=$1; sub(":","",iface) }
            /inet /  { print "ipv4|" $2 "|" iface }
            /inet6 / { print "ipv6|" $2 "|" iface }'
        return 0
    fi
    log WARN "address probe unavailable"
}

probe_ports() {
    # Order: ss, then netstat, then lsof. Emits "protocol|address|port".
    local out
    out="$(ss -lntu 2>/dev/null || true)"
    if [ -n "$out" ]; then
        printf '%s\n' "$out" | awk 'NR>1 {
            proto=tolower($1); split($5,a,":"); port=a[length(a)];
            if (port ~ /^[0-9]+$/ && (proto=="tcp"||proto=="udp")) print proto "||" port }'
        return 0
    fi
    out="$(netstat -an 2>/dev/null || true)"
    if [ -n "$out" ]; then
        printf '%s\n' "$out" | awk '/LISTEN|udp/ {
            proto=tolower($1); sub(/[46]$/,"",proto);
            split($4,a,":"); port=a[length(a)];
            if (port ~ /^[0-9]+$/ && (proto=="tcp"||proto=="udp")) print proto "||" port }'
        return 0
    fi
    log WARN "port probe unavailable"
}

#_______________________________________________________________________________
# Record

PROCS="$(probe_process_count)"
UPTIME="$(probe_uptime_seconds)"
NOW_MS="$(now_ms)"

SNAPSHOT_ID="$(sqlite_exec "$DB_PATH" list \
"INSERT INTO snapshot (unixtime_ms, iso_local, iso_utc, tz_offset_minutes, hostname,
   username, process_count, uptime_seconds, os_platform, os_release, script_pid,
   script_flavor, invocation_source)
VALUES (CAST(:ms AS INTEGER), :isolocal, :isoutc, CAST(:tz AS INTEGER), :host, :user,
        CASE WHEN :procs = '' THEN NULL ELSE CAST(:procs AS INTEGER) END,
        CASE WHEN :uptime = '' THEN NULL ELSE CAST(:uptime AS INTEGER) END,
        :platform, :release, CAST(:pid AS INTEGER), 'posix', :source);
SELECT last_insert_rowid();" \
    "ms=$NOW_MS" "isolocal=$(iso_local)" "isoutc=$(iso_utc)" \
    "tz=$(tz_offset_minutes)" "host=$(hostname)" \
    "user=$(id -un 2>/dev/null || printf '%s' "${USER:-unknown}")" \
    "procs=$PROCS" "uptime=$UPTIME" \
    "platform=$(uname -s | tr '[:upper:]' '[:lower:]')" \
    "release=$(uname -sr)" "pid=$$" "source=$INVOCATION_SOURCE" | tail -n1)"

case "$SNAPSHOT_ID" in
    ''|*[!0-9]*) die_runtime "snapshot insert did not return an id; nothing was recorded." ;;
esac

# Child rows only ever written for a snapshot that was successfully inserted,
# so a partial failure never orphans them.
ADDR_COUNT=0
while IFS='|' read -r family address iface; do
    [ -n "${address:-}" ] || continue
    sqlite_exec "$DB_PATH" list \
"INSERT INTO snapshot_address (snapshot_id, family, address, interface, scope)
VALUES (CAST(:sid AS INTEGER), :family, :address,
        CASE WHEN :iface = '' THEN NULL ELSE :iface END, NULL);" \
        "sid=$SNAPSHOT_ID" "family=$family" "address=$address" "iface=${iface:-}" >/dev/null
    ADDR_COUNT=$((ADDR_COUNT + 1))
done < <(probe_addresses)

PORT_COUNT=0
if [ "$SKIP_PORTS" -eq 0 ]; then
    while IFS='|' read -r proto address port; do
        case "${port:-}" in ''|*[!0-9]*) continue ;; esac
        [ "$port" -ge 1 ] && [ "$port" -le 65535 ] || continue
        sqlite_exec "$DB_PATH" list \
"INSERT INTO snapshot_port (snapshot_id, protocol, family, address, port, process_name)
VALUES (CAST(:sid AS INTEGER), :proto, NULL,
        CASE WHEN :address = '' THEN NULL ELSE :address END,
        CAST(:port AS INTEGER), NULL);" \
            "sid=$SNAPSHOT_ID" "proto=$proto" "address=${address:-}" "port=$port" >/dev/null
        PORT_COUNT=$((PORT_COUNT + 1))
    done < <(probe_ports)
fi

log SUCCESS "snapshot $SNAPSHOT_ID recorded: $ADDR_COUNT address(es), $PORT_COUNT listening port(s)"
printf '%s\n' "$DB_PATH"
exit 0
