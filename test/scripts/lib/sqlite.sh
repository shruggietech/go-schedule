#!/usr/bin/env bash
# Shared sqlite3 resolution, installation, and access for the go-schedule
# maintainer test scripts. Sourced by the three Test-*.sh twins.
#
# This is the POSIX counterpart of lib/Sqlite.ps1 and mirrors it function for
# function. There is exactly one implementation per twin on purpose: three
# copies of the resolution order would be three chances for them to disagree,
# and a disagreement here presents as an intermittent, platform-specific test
# failure rather than an obvious bug.
#
# Exit-code contract, shared by every script that sources this file:
#   0  success
#   1  runtime failure (probe failed, write failed, contention exhausted)
#   2  usage error or unmet prerequisite (bad arguments, no usable sqlite3,
#      unsupported platform for the installer)
#
# Results go to stdout; diagnostics and warnings go to stderr.
#
# Shell conventions for these scripts are documented in docs/test-scripts.md;
# there is no separate shell house standard, so that section is the contract.

set -euo pipefail

LOG_QUIET=0
SQLITE_EXE=""
SQLITE_MIN_VERSION="3.33.0"
LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_SCRIPTS_ROOT="$(dirname "$LIB_DIR")"

#_______________________________________________________________________________
# Logging

log() {
    # log LEVEL MESSAGE
    local level="$1"; shift
    if [ "$LOG_QUIET" -eq 1 ]; then
        case "$level" in INFO|SUCCESS|DEBUG) return 0 ;; esac
    fi
    printf '%s %-7s %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$level" "$*" >&2
}

die_runtime() { log ERROR "$*"; exit 1; }
die_usage()   { log ERROR "$*"; exit 2; }

#_______________________________________________________________________________
# Platform

test_platform() {
    local os arch
    case "$(uname -s)" in
        Darwin*) os=darwin ;;
        Linux*)  os=linux ;;
        MINGW*|MSYS*|CYGWIN*) os=windows ;;
        *) os=linux ;;
    esac
    case "$(uname -m)" in
        arm64|aarch64) arch=arm64 ;;
        *) arch=x64 ;;
    esac
    printf '%s-%s\n' "$os" "$arch"
}

package_manager_hint() {
    case "$(test_platform)" in
        darwin-*)  printf 'brew install sqlite\n' ;;
        windows-*) printf 'winget install SQLite.SQLite\n' ;;
        *)         printf 'sudo apt install sqlite3\n' ;;
    esac
}

#_______________________________________________________________________________
# sqlite3 resolution

version_ge() {
    # version_ge FOUND MINIMUM -- true when FOUND >= MINIMUM
    awk -v found="$1" -v minimum="$2" 'BEGIN {
        split(found,a,"."); split(minimum,b,".")
        for(i=1;i<=3;i++) { if(a[i]+0<b[i]+0) exit 1; if(a[i]+0>b[i]+0) exit 0 }
        exit 0
    }'
}

test_sqlite_candidate() (
    # A private, bounded probe preserves stderr and the original native exit.
    local candidate="$1" raw token rc probe_dir probe_pid watcher_pid
    if ! command -v "$candidate" >/dev/null 2>&1; then
        log WARN "sqlite3 candidate '$candidate' could not be launched: executable not found."
        return 0
    fi
    umask 077
    probe_dir="$(mktemp -d)" || return 0
    local timeout_exe=""
    if [ -x /usr/bin/timeout ]; then timeout_exe=/usr/bin/timeout
    elif command -v gtimeout >/dev/null 2>&1; then timeout_exe="$(command -v gtimeout)"; fi
    if [ -n "$timeout_exe" ]; then
        # Coreutils owns the whole probe group. Both channels stop at 4097
        # bytes, which lets us detect truncation without unbounded temp files.
        local started=$SECONDS
        if "$timeout_exe" --kill-after=1 10 "$BASH" -c '
            "$1" --version </dev/null 2> >(head -c 4097 >"$2") | head -c 4097 >"$3"
            result=${PIPESTATUS[0]}
            wait
            exit "$result"
        ' _ "$candidate" "$probe_dir/err" "$probe_dir/out" 2>"$probe_dir/launcher"; then rc=0; else rc=$?; fi
        if [ "$rc" -eq 124 ] || { [ "$rc" -eq 137 ] && [ "$((SECONDS-started))" -ge 10 ]; }; then : >"$probe_dir/timeout"; fi
    else
    (ulimit -f 8 2>/dev/null || true; exec "$candidate" --version </dev/null) >"$probe_dir/out" 2>"$probe_dir/err" &
    probe_pid=$!
    # The parent receives the sleep PID before waiting for the candidate, so it
    # can stop both watchdog processes immediately after a fast probe.
    local sleeper
    read -r sleeper < <(
        sleep 10 </dev/null >/dev/null 2>&1 &
        local timer=$!
        printf '%s\n' "$timer"
        if wait "$timer" 2>/dev/null && [ ! -f "$probe_dir/done" ]; then
            printf 'timeout' >"$probe_dir/timeout"
            kill -KILL "$probe_pid" 2>/dev/null || true
        fi
    )
    watcher_pid=$!
    if wait "$probe_pid" 2>/dev/null; then rc=0; else rc=$?; fi
    : >"$probe_dir/done"
    kill -KILL "$sleeper" 2>/dev/null || true
    wait "$watcher_pid" 2>/dev/null || true
    fi
    raw="$(head -c 4096 "$probe_dir/out")"
    local detail
    detail="$(head -c 4096 "$probe_dir/err" | tr '\000-\037\177' ' ')"
    if [ -f "$probe_dir/timeout" ]; then
        log WARN "sqlite3 candidate '$candidate' timed out after 10000 ms. stderr: $detail"
        rm -rf "$probe_dir"
        return 0
    fi
    local output_bytes error_bytes
    output_bytes="$(wc -c <"$probe_dir/out")"
    error_bytes="$(wc -c <"$probe_dir/err")"
    rm -rf "$probe_dir"
    if [ "$output_bytes" -gt 4096 ] || [ "$error_bytes" -gt 4096 ]; then
        log WARN "sqlite3 candidate '$candidate' exceeded the 4096-byte probe output limit (exit $rc)."
        return 0
    fi
    if [ "$rc" -ne 0 ]; then
        log WARN "sqlite3 candidate '$candidate' exited $rc. stderr: $detail"
        return 0
    fi
    token="$(printf '%s' "$raw" | awk 'NR==1 {print $1}')"
    if [[ ! "$token" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log WARN "sqlite3 candidate '$candidate' returned an empty or malformed version."
        return 0
    fi
    if version_ge "$token" "$SQLITE_MIN_VERSION"; then
        printf '%s\n' "$token"
    else
        log DEBUG "Ignoring sqlite3 at '$candidate': version $token is below the required $SQLITE_MIN_VERSION."
    fi
)

resolve_sqlite() {
    # Strict precedence: explicit path, then repo-local .bin/, then PATH.
    #
    # The fall-through-on-unusable rule applies only to the IMPLICIT
    # candidates, where continuing past a stale .bin/ to a good tool on PATH is
    # the helpful thing to do. An EXPLICIT path is different: the maintainer
    # named one specific tool, and quietly running a different one instead is
    # how you debug the wrong binary for an hour. A bad explicit path is a hard
    # usage error.
    local explicit="${1:-}" candidate ver
    if [ -n "$explicit" ]; then
        if [ ! -f "$explicit" ]; then
            log ERROR "--sqlite-exe '$explicit' does not exist."
            log ERROR "Refusing to silently fall back to another sqlite3 when you named a specific one."
            exit 2
        fi
        ver="$(test_sqlite_candidate "$explicit")"
        if [ -z "$ver" ]; then
            log ERROR "--sqlite-exe '$explicit' is not runnable or is older than $SQLITE_MIN_VERSION."
            exit 2
        fi
        log DEBUG "Using sqlite3 $ver at $explicit"
        printf '%s\n' "$explicit"
        return 0
    fi
    for candidate in \
        "$TEST_SCRIPTS_ROOT/.bin/sqlite3" \
        "$(command -v sqlite3 2>/dev/null || true)"
    do
        [ -n "$candidate" ] || continue
        [ -x "$candidate" ] || continue
        ver="$(test_sqlite_candidate "$candidate")"
        if [ -n "$ver" ]; then
            log DEBUG "Using sqlite3 $ver at $candidate"
            printf '%s\n' "$candidate"
            return 0
        fi
    done
    return 1
}

exit_no_sqlite() {
    # Unmet prerequisite is exit 2, not 1. The message names both remedies
    # because a maintainer who hits this needs to act, not to investigate.
    log ERROR "No usable sqlite3 found (need $SQLITE_MIN_VERSION or later)."
    log ERROR "Fix it either way:"
    log ERROR "  - rerun this script with --install-sqlite (downloads the pinned build, verifies its checksum)"
    log ERROR "  - or install it yourself: $(package_manager_hint)"
    exit 2
}

#_______________________________________________________________________________
# Installer -- verification precedes unpacking

manifest_field() {
    # manifest_field PLATFORM FIELD -- minimal JSON field read, no jq dependency
    local plat="$1" field="$2"
    python3 - "$LIB_DIR/sqlite-manifest.json" "$plat" "$field" <<'PY' 2>/dev/null || true
import json,sys
try:
    m=json.load(open(sys.argv[1]))
    e=m.get("platforms",{}).get(sys.argv[2])
    print(e.get(sys.argv[3],"") if e else "")
except Exception:
    print("")
PY
}

install_sqlite() {
    local plat url want tmp actual stage found dest bin
    plat="$(test_platform)"
    url="$(manifest_field "$plat" url)"
    want="$(manifest_field "$plat" sha256)"
    bin="$(manifest_field "$plat" binary)"

    if [ -z "$url" ]; then
        log ERROR "No prebuilt sqlite3 is published for $plat."
        log ERROR "Install it from your package manager instead; this installer does not build from source."
        exit 2
    fi
    if [ -z "$want" ]; then
        log ERROR "Manifest has no checksum for this platform; refusing to install."
        die_runtime "An unverified binary is worse than no binary. Use your package manager."
    fi

    mkdir -p "$TEST_SCRIPTS_ROOT/.bin"
    tmp="$(mktemp)"
    log INFO "Downloading $url"
    if ! curl -fsSL --max-time 300 -o "$tmp" "$url"; then
        rm -f "$tmp"
        die_runtime "Download failed."
    fi

    actual="$(sha256sum "$tmp" 2>/dev/null | awk '{print $1}')"
    [ -n "$actual" ] || actual="$(shasum -a 256 "$tmp" | awk '{print $1}')"
    if [ "$actual" != "$want" ]; then
        rm -f "$tmp"
        log ERROR "CHECKSUM MISMATCH -- download discarded and not installed."
        log ERROR "  expected $want"
        log ERROR "  actual   $actual"
        exit 1
    fi
    log SUCCESS "Checksum verified."

    stage="$(mktemp -d)"
    unzip -qo "$tmp" -d "$stage" || { rm -rf "$tmp" "$stage"; die_runtime "Could not unpack archive."; }
    rm -f "$tmp"
    found="$(find "$stage" -name "$bin" -type f | head -n1)"
    if [ -z "$found" ]; then
        rm -rf "$stage"
        die_runtime "Archive did not contain $bin."
    fi
    dest="$TEST_SCRIPTS_ROOT/.bin/$bin"
    cp "$found" "$dest"
    chmod +x "$dest"
    rm -rf "$stage"

    if [ -z "$(test_sqlite_candidate "$dest")" ]; then
        die_runtime "Installed binary does not run or is too old."
    fi
    log SUCCESS "Installed sqlite3 to $dest"
    printf '%s\n' "$dest"
}

#_______________________________________________________________________________
# Data directory and database resolution

test_data_dir() {
    # A user-writable directory, never the daemon's system-wide data dir: test
    # payloads must run unelevated, and disposable test output does not belong
    # in the directory holding live scheduler state.
    local base
    if [ -n "${GOSCHEDULE_TEST_DIR:-}" ]; then
        base="$GOSCHEDULE_TEST_DIR"
    elif [ "$(uname -s)" = "Darwin" ]; then
        base="$HOME/Library/Application Support/goschedule-test"
    elif [ -n "${XDG_DATA_HOME:-}" ]; then
        base="$XDG_DATA_HOME/goschedule-test"
    else
        base="$HOME/.local/share/goschedule-test"
    fi
    mkdir -p "$base"
    printf '%s\n' "$base"
}

resolve_test_database() {
    local name="$1"
    case "$name" in
        */*|*\\*|*.db) printf '%s\n' "$name" ;;
        *) printf '%s/%s.db\n' "$(test_data_dir)" "$name" ;;
    esac
}

#_______________________________________________________________________________
# Query execution -- bound parameters, WAL, bounded retry

sqlite_exec() {
    # sqlite_exec DATABASE MODE SQL [name=value ...]
    #
    # Bound parameters via '.param set'. The values written include the
    # hostname, the username, and network interface names -- all capable of
    # containing a quote, and all influenceable by whoever administers the
    # machine. String-interpolated SQL here would be both an injection vector
    # and an ordinary bug for any user named O'Brien.
    local db="$1" mode="$2" sql="$3"; shift 3
    local script attempt out rc pair key value encoded

    # PRAGMA statements return rows; left unredirected they land in the
    # caller's result and get mistaken for data. Silence them at the source.
    local null_device=/dev/null
    case "$(uname -s)" in MINGW*|MSYS*|CYGWIN*) null_device=nul ;; esac
    script=".bail on
.output $null_device
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout=5000;
.output stdout
.mode $mode
"
    for pair in "$@"; do
        key="${pair%%=*}"
        value="${pair#*=}"
        if [ "$value" = "__NULL__" ]; then
            script="$script.param set :$key NULL
"
        else
            # The CLI removes one quote layer before evaluating a parameter as
            # an SQL expression. Literal-looking IDs such as 1e999... can become
            # Inf, and a newline can become another dot command. Encode bytes
            # into one quoted expression whose evaluated value is exact text.
            encoded="$(printf '%s' "$value" | od -An -v -tx1 | tr -d ' \n')"
            script="$script.param set :$key \"CAST(X'$encoded' AS TEXT)\"
"
        fi
    done
    script="$script$sql"

    for attempt in 1 2 3; do
        set +e
        out="$(printf '%s\n' "$script" | "$SQLITE_EXE" "$db" 2>&1)"
        rc=$?
        set -e
        if [ $rc -eq 0 ]; then
            printf '%s\n' "$out"
            return 0
        fi
        case "$out" in
            *locked*|*busy*)
                log WARN "Database contention on attempt $attempt/3; retrying."
                sleep "0.$((25 * attempt))"
                ;;
            *)
                die_runtime "sqlite3 failed: $out"
                ;;
        esac
    done
    log ERROR "Database still contended after 3 attempts. The record was NOT written. Last sqlite3 error: $out"
    die_runtime "This is a test-harness failure, not a scheduler defect."
}

#_______________________________________________________________________________
# Schemas -- kept byte-identical in intent to lib/Sqlite.ps1

init_heartbeat_schema() {
    sqlite_exec "$1" list "CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT);
INSERT OR IGNORE INTO meta(key,value) VALUES('schema_version','2');
CREATE TABLE IF NOT EXISTS beat (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id TEXT NOT NULL,
  sequence INTEGER NOT NULL,
  label TEXT,
  hostname TEXT NOT NULL,
  username TEXT,
  pid INTEGER NOT NULL,
  started_ms INTEGER NOT NULL,
  started_iso TEXT NOT NULL,
  finished_ms INTEGER NOT NULL,
  duration_ms INTEGER NOT NULL,
  expected_ms INTEGER,
  expected_source TEXT NOT NULL CHECK (expected_source IN ('env','anchor','boundary','none')),
  drift_ms INTEGER,
  interval_seconds INTEGER,
  exit_code INTEGER NOT NULL,
  outcome TEXT NOT NULL CHECK (outcome IN ('ok','failed')),
  sched_env TEXT,
  CHECK (finished_ms >= started_ms),
  CHECK ((drift_ms IS NULL) = (expected_ms IS NULL)),
  CHECK ((expected_source = 'none') = (expected_ms IS NULL)),
  CHECK ((outcome = 'ok') = (exit_code = 0)),
  UNIQUE (session_id, sequence)
);
CREATE INDEX IF NOT EXISTS beat_started ON beat(started_ms);" >/dev/null
}

init_system_schema() {
    sqlite_exec "$1" list "PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT);
INSERT OR IGNORE INTO meta(key,value) VALUES('schema_version','1');
CREATE TABLE IF NOT EXISTS snapshot (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  unixtime_ms INTEGER NOT NULL,
  iso_local TEXT NOT NULL,
  iso_utc TEXT NOT NULL,
  tz_offset_minutes INTEGER NOT NULL,
  hostname TEXT NOT NULL,
  username TEXT,
  process_count INTEGER,
  uptime_seconds INTEGER,
  os_platform TEXT NOT NULL,
  os_release TEXT,
  script_pid INTEGER NOT NULL,
  script_flavor TEXT NOT NULL CHECK (script_flavor IN ('powershell','posix')),
  invocation_source TEXT NOT NULL,
  addresses_probe TEXT,
  ports_probe TEXT
);
CREATE INDEX IF NOT EXISTS snapshot_time ON snapshot(unixtime_ms);
CREATE TABLE IF NOT EXISTS snapshot_address (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  snapshot_id INTEGER NOT NULL REFERENCES snapshot(id) ON DELETE CASCADE,
  family TEXT NOT NULL CHECK (family IN ('ipv4','ipv6')),
  address TEXT NOT NULL,
  interface TEXT,
  scope TEXT
);
CREATE INDEX IF NOT EXISTS address_snapshot ON snapshot_address(snapshot_id);
CREATE TABLE IF NOT EXISTS snapshot_port (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  snapshot_id INTEGER NOT NULL REFERENCES snapshot(id) ON DELETE CASCADE,
  protocol TEXT NOT NULL CHECK (protocol IN ('tcp','udp')),
  family TEXT,
  address TEXT,
  port INTEGER NOT NULL CHECK (port BETWEEN 1 AND 65535),
  process_name TEXT
);
CREATE INDEX IF NOT EXISTS port_snapshot ON snapshot_port(snapshot_id, protocol, port);" >/dev/null

    # Forward-only, non-destructive migration for databases created before
    # schema 3. CREATE TABLE IF NOT EXISTS does nothing to an existing table, so
    # the probe-status columns have to be added explicitly. ALTER TABLE ADD
    # COLUMN cannot carry a CHECK in SQLite, so these stay plain TEXT in both
    # fresh and migrated databases -- a fresh database with a stricter schema
    # than a migrated one is a difference nobody would think to look for. The
    # writer enforces the vocabulary instead.
    local col has
    for col in addresses_probe ports_probe; do
        has="$(sqlite_exec "$1" list "SELECT COUNT(*) FROM pragma_table_info('snapshot') WHERE name='$col';" | tail -n1)"
        if [ "$has" = "0" ]; then
            log DEBUG "migrating system.db: adding $col"
            sqlite_exec "$1" list "ALTER TABLE snapshot ADD COLUMN $col TEXT;" >/dev/null
        fi
    done
    sqlite_exec "$1" list "INSERT OR REPLACE INTO meta(key,value) VALUES('schema_version','3');" >/dev/null
}

#_______________________________________________________________________________
# Initialization entry point

init_test_sqlite() {
    # init_test_sqlite EXPLICIT_PATH INSTALL_FLAG
    local explicit="${1:-}" install="${2:-0}"
    if [ "$install" -eq 1 ]; then
        SQLITE_EXE="$(install_sqlite | tail -n1)"
        return 0
    fi
    if ! SQLITE_EXE="$(resolve_sqlite "$explicit")"; then
        exit_no_sqlite
    fi
}

now_ms() {
    # Millisecond epoch. GNU date does %3N; BSD/macOS date does not, so fall
    # back to whole seconds rather than emitting the literal string "%3N" and
    # writing a nonsense timestamp into the database.
    local n
    n="$(date +%s%3N 2>/dev/null || true)"
    case "$n" in
        *N*|"") printf '%s000\n' "$(date +%s)" ;;
        *) printf '%s\n' "$n" ;;
    esac
}

parse_iso_epoch() {
    # Echo Unix seconds for an RFC 3339 timestamp; echo nothing if unparseable.
    #
    # GNU date and BSD date have incompatible parsing interfaces and macOS ships
    # BSD, where `date -d` is not a parse flag at all. Trying GNU first and
    # falling back keeps one code path for both. This is the portability trap
    # that shipped broken in v0.5.1: it cannot reproduce on a GNU-date host, so
    # only a real macOS runner catches it.
    local ts="$1" out norm off
    out="$(date -d "$ts" +%s 2>/dev/null || true)"
    if [ -n "$out" ]; then printf '%s\n' "$out"; return 0; fi

    # BSD: the format must be spelled out, fractional seconds are not accepted,
    # and a numeric offset must be +0500 rather than +05:00.
    norm="${ts}"
    case "$norm" in
        *Z)
            norm="${norm%Z}"
            norm="${norm%%.*}"
            out="$(TZ=UTC date -j -f '%Y-%m-%dT%H:%M:%S' "$norm" +%s 2>/dev/null || true)"
            [ -n "$out" ] && { printf '%s\n' "$out"; return 0; }
            ;;
        *[+-][0-9][0-9]:[0-9][0-9])
            off="${norm##*[+-]}"
            off="${off%%:*}${off##*:}"
            norm="${norm%[+-]*}"
            case "$ts" in *+*) off="+$off" ;; *) off="-$off" ;; esac
            norm="${norm%%.*}"
            out="$(date -j -f '%Y-%m-%dT%H:%M:%S%z' "${norm}${off}" +%s 2>/dev/null || true)"
            [ -n "$out" ] && { printf '%s\n' "$out"; return 0; }
            ;;
    esac
    norm="${ts%%.*}"
    out="$(date -j -f '%Y-%m-%dT%H:%M:%S' "$norm" +%s 2>/dev/null || true)"
    [ -n "$out" ] && printf '%s\n' "$out"
    return 0
}

iso_local() { date '+%Y-%m-%dT%H:%M:%S%z'; }
iso_utc()   { date -u '+%Y-%m-%dT%H:%M:%SZ'; }
tz_offset_minutes() {
    local off sign hh mm
    off="$(date '+%z')"
    sign="${off:0:1}"; hh="${off:1:2}"; mm="${off:3:2}"
    hh="${hh#0}"; mm="${mm#0}"
    printf '%s%s\n' "$([ "$sign" = "-" ] && printf -- '-')" "$(( ${hh:-0} * 60 + ${mm:-0} ))"
}
