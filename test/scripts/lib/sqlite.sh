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
    [ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]
}

test_sqlite_candidate() {
    # Echoes the version when usable and new enough; otherwise nothing.
    # A candidate that exists but is broken or too old must be treated as
    # not-found so the search continues -- a stale sqlite3 early in the order
    # must never shadow a good one later in it.
    local candidate="$1" raw token
    raw="$("$candidate" --version 2>/dev/null)" || return 0
    [ -n "$raw" ] || return 0
    token="$(printf '%s' "$raw" | awk '{print $1}')"
    if version_ge "$token" "$SQLITE_MIN_VERSION"; then
        printf '%s\n' "$token"
    else
        log DEBUG "Ignoring sqlite3 at '$candidate': version $token is below the required $SQLITE_MIN_VERSION."
    fi
}

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
    local script attempt out rc pair key value esc

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
            esc="${value//\'/\'\'}"
            script="$script.param set :$key '$esc'
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
    log ERROR "Database still contended after 3 attempts. The record was NOT written."
    die_runtime "This is a test-harness failure, not a scheduler defect."
}

#_______________________________________________________________________________
# Schemas -- kept byte-identical in intent to lib/Sqlite.ps1

init_heartbeat_schema() {
    sqlite_exec "$1" list "CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT);
INSERT OR IGNORE INTO meta(key,value) VALUES('schema_version','1');
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
  expected_source TEXT NOT NULL CHECK (expected_source IN ('env','boundary','none')),
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
  invocation_source TEXT NOT NULL
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

iso_local() { date '+%Y-%m-%dT%H:%M:%S%z'; }
iso_utc()   { date -u '+%Y-%m-%dT%H:%M:%SZ'; }
tz_offset_minutes() {
    local off sign hh mm
    off="$(date '+%z')"
    sign="${off:0:1}"; hh="${off:1:2}"; mm="${off:3:2}"
    hh="${hh#0}"; mm="${mm#0}"
    printf '%s%s\n' "$([ "$sign" = "-" ] && printf -- '-')" "$(( ${hh:-0} * 60 + ${mm:-0} ))"
}
