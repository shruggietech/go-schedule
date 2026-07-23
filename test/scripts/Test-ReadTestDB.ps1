<#
.SYNOPSIS
    Run canned inspection queries against the go-schedule test databases.

.DESCRIPTION
    The recording scripts write; this one reads. It carries a fixed set of named
    queries chosen to answer the questions that actually come up after a test
    window, so a maintainer does not have to remember the schema or write SQL at
    the moment they are trying to diagnose something.

    Two reporting rules are contract rather than presentation, because breaking
    them would produce confident numbers that mean less than they appear to:

      - A query that excludes rows says how many it excluded and why. A
        percentile over an unstated subset is a number drawn from unknown
        evidence.
      - Drift is never pooled across expected-source values, and never shown
        without its source. A measured value and a derived one are different
        kinds of number, and averaging them produces neither.

    Boundary-derived drift larger than a quarter of the interval is flagged
    unreliable: past that magnitude a late firing and an early next one are not
    distinguishable, so reporting the figure as fact would be a guess wearing a
    decimal point.

    Exit codes: 0 success, 1 runtime failure, 2 usage error or unmet
    prerequisite.

.PARAMETER Database
    'heartbeat', 'system', or an explicit path.
    Alias: d

.PARAMETER Query
    Which canned query to run ("kind"). Use -List to see them all.
    Default: 'summary'.
    Alias: k (not 'q' -- PowerShell aliases are case-insensitive, so -Q would
    collide with -Quiet)

.PARAMETER List
    List the available queries with the question each one answers, then exit.
    Alias: n

.PARAMETER Format
    Output form: Table, Json, or Csv.
    Default: 'Table'.
    Alias: f

.PARAMETER Limit
    Maximum rows for queries that return rows.
    Default: 20.
    Alias: m

.PARAMETER IntervalSeconds
    Expected schedule interval, used by the gap and reliability checks. When
    omitted it is inferred from the most common observed interval, and the
    output says that it was inferred.
    Alias: i

.PARAMETER AnchorIso
    One real firing time of the schedule, RFC 3339, from `gosched task show`.
    With -IntervalSeconds this makes -Query drift compute true dispatch latency
    at read time, from the raw start timestamps.

    Read time is the right place for this. The anchor cannot be known before the
    task exists, because this scheduler derives it from the task's creation
    moment -- so passing it to the recorder is a chicken-and-egg problem. Drift
    is a derived quantity; deriving it when the anchor is actually knowable
    costs nothing and works on beats already recorded.
    Alias: a

.PARAMETER SqliteExe
    Explicit path to sqlite3. Highest precedence in the search order.

.PARAMETER InstallSqlite
    Download the pinned sqlite3 build, verify its checksum, and install it.

.PARAMETER Quiet
    Suppress informational output. Warnings and errors are never suppressed.
    Alias: q

.PARAMETER Help
    Print this help text to the terminal.
    Alias: h

.EXAMPLE
    .\Test-ReadTestDB.ps1 -Database heartbeat -Query summary
    How many beats, over what period, from how many sessions.

.EXAMPLE
    .\Test-ReadTestDB.ps1 -Database heartbeat -Query drift -IntervalSeconds 60
    The dispatch-latency question, broken down by how each figure was derived.

.EXAMPLE
    .\Test-ReadTestDB.ps1 -Database heartbeat -Query gaps -IntervalSeconds 60
    Which expected firings were missed or badly delayed.

.EXAMPLE
    .\Test-ReadTestDB.ps1 -Database system -Query listeners -Format Json
    Listening ports from the latest snapshot, machine-readable.

.NOTES
    Requires PowerShell 7+ and sqlite3 3.33.0 or later.
    Full documentation: docs/test-scripts.md
#>
[CmdletBinding(SupportsShouldProcess=$false,ConfirmImpact='None',DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("d")]
    [string]$Database = 'heartbeat',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("k")]
    [string]$Query = 'summary',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("n")]
    [Switch]$List,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("f")]
    [ValidateSet('Table','Json','Csv')]
    [string]$Format = 'Table',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("m")]
    [int]$Limit = 20,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("i")]
    [int]$IntervalSeconds = 0,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("a")]
    [string]$AnchorIso = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [string]$SqliteExe = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Switch]$InstallSqlite,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("q")]
    [Switch]$Quiet,

    [Parameter(Mandatory=$true,ParameterSetName='HelpText')]
    [Alias("h")]
    [Switch]$Help
)
#_______________________________________________________________________________
## Declare Functions

    function Show-QueryList {
        [CmdletBinding()]
        Param()
        Write-Output "Available queries (-Query <name>):"
        Write-Output ""
        foreach ($k in ($script:Queries.Keys | Sort-Object)) {
            Write-Output ("  {0,-10} [{1,-9}] {2}" -f $k, $script:Queries[$k].Db, $script:Queries[$k].Question)
        }
        Write-Output ""
        Write-Output "Databases: 'heartbeat', 'system', or an explicit path."
    }

    function Get-InferredInterval {
        # The modal observed inter-beat interval, in seconds. Used only when the
        # caller did not declare one; the output always says which happened,
        # because an inferred interval is an assumption and assumptions that
        # look like measurements are how wrong conclusions get confident.
        [CmdletBinding()]
        Param([Parameter(Mandatory=$true)][string]$DbPath)
        $sql = @'
SELECT CAST(ROUND(delta / 1000.0) AS INTEGER) AS secs, COUNT(*) AS n
FROM (SELECT started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta FROM beat)
WHERE delta IS NOT NULL AND delta > 0
GROUP BY secs ORDER BY n DESC, secs ASC LIMIT 1;
'@
        $out = Invoke-Sqlite -Database $DbPath -Sql $sql -Mode 'list'
        $line = ($out | Where-Object { $_ -match '^\d+\|' } | Select-Object -First 1)
        if ($line) { return [int]($line -split '\|')[0] }
        return 0
    }

    function Get-ExcludedNote {
        # FR-013a: a query that drops rows must say so. Silence here reads as
        # "nothing was dropped", which is a different and stronger claim.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)][string]$DbPath,
            [Parameter(Mandatory=$true)][string]$CountSql,
            [Parameter(Mandatory=$true)][string]$Reason
        )
        $out = Invoke-Sqlite -Database $DbPath -Sql $CountSql -Mode 'list'
        $n = ($out | Where-Object { $_ -match '^\d+$' } | Select-Object -First 1)
        if ($n -and [int]$n -gt 0) {
            Write-Log ("excluded {0} record(s): {1}" -f $n, $Reason) -Level Warn
        }
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $ThisScriptPath = $MyInvocation.MyCommand.Path

    $script:AnchorMs = $null

    $script:Queries = @{
        'summary' = @{ Db = 'both'; Question = 'How many records, over what period, from how many sessions and hosts?' }
        'recent'  = @{ Db = 'both'; Question = 'What are the most recent records?' }
        'cadence' = @{ Db = 'heartbeat'; Question = 'What were the observed intervals (min/p50/p95/p99/max)?' }
        'drift'   = @{ Db = 'heartbeat'; Question = 'True dispatch latency, by expected-source (needs an anchor).' }
        'jitter'  = @{ Db = 'heartbeat'; Question = 'Variation around the observed firing phase (no anchor needed).' }
        'gaps'    = @{ Db = 'heartbeat'; Question = 'Which expected firings were missed or badly delayed?' }
        'overlaps'= @{ Db = 'heartbeat'; Question = 'Which runs overlapped in time?' }
        'failures'= @{ Db = 'heartbeat'; Question = 'Which runs reported failure?' }
        'restarts'= @{ Db = 'heartbeat'; Question = 'Where are the session boundaries?' }
        'hosts'   = @{ Db = 'both'; Question = 'Which hosts and users produced records?' }
        'listeners'= @{ Db = 'system'; Question = 'What is listening now, and what changed since the previous snapshot?' }
        'schema'  = @{ Db = 'both'; Question = 'What is the stored structure?' }
    }

#_______________________________________________________________________________
## Execute Operations

    # Catch help text requests
    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }

    . (Join-Path $PSScriptRoot 'lib/Sqlite.ps1')
    if ($Quiet) { $script:LogQuiet = $true }

    if ($List) { Show-QueryList; exit 0 }

    if (-not $script:Queries.ContainsKey($Query)) {
        Write-Log ("Unknown query '{0}'. Use -List to see the available queries." -f $Query) -Level Error
        exit 2
    }
    if ($Limit -le 0) { Write-Log "-Limit must be positive." -Level Error; exit 2 }

    Initialize-TestSqlite -SqliteExe $SqliteExe -InstallSqlite:$InstallSqlite

    $dbPath = Resolve-TestDatabase -Name $Database
    if (-not (Test-Path -LiteralPath $dbPath -PathType Leaf)) {
        Write-Log ("No database at {0}. Run a recording script first." -f $dbPath) -Level Error
        exit 2
    }
    Write-Log ("database: {0}" -f $dbPath)

    $isHeartbeat = ($Query -in @('cadence','drift','jitter','gaps','overlaps','failures','restarts')) -or
                   ($Database -eq 'heartbeat')
    $mode = switch ($Format) { 'Json' { 'json' } 'Csv' { 'csv' } default { 'box' } }

    # Interval resolution, and saying which it was.
    $intervalUsed = $IntervalSeconds
    $intervalInferred = $false
    if ($intervalUsed -le 0 -and $Query -in @('gaps','drift','jitter')) {
        $intervalUsed = Get-InferredInterval -DbPath $dbPath
        $intervalInferred = $true
        if ($intervalUsed -gt 0) {
            Write-Log ("interval not supplied; inferred {0}s from the most common observed interval" -f $intervalUsed) -Level Warn
        } else {
            Write-Log "interval not supplied and could not be inferred (too few records)" -Level Warn
        }
    }

    if ($AnchorIso) {
        if ($intervalUsed -le 0) {
            Write-Log "-AnchorIso needs -IntervalSeconds to reconstruct the firing grid." -Level Error
            exit 2
        }
        try {
            $script:AnchorMs = [long]([DateTimeOffset]::Parse($AnchorIso)).ToUnixTimeMilliseconds()
        } catch {
            Write-Log ("-AnchorIso '{0}' is not a parseable timestamp." -f $AnchorIso) -Level Error
            exit 2
        }
        Write-Log ("drift derived at read time from anchor {0} and interval {1}s" -f $AnchorIso, $intervalUsed)
    }

    $sql = switch ($Query) {
        'summary' {
            if ($isHeartbeat) {
@'
SELECT COUNT(*) AS beats, COUNT(DISTINCT session_id) AS sessions,
       COUNT(DISTINCT hostname) AS hosts,
       MIN(started_iso) AS first_seen,
       datetime(MAX(started_ms)/1000,'unixepoch') AS last_seen_utc,
       ROUND((MAX(started_ms)-MIN(started_ms))/60000.0,1) AS span_minutes,
       SUM(CASE WHEN outcome='failed' THEN 1 ELSE 0 END) AS failures
FROM beat;
'@
            } else {
@'
SELECT COUNT(*) AS snapshots, COUNT(DISTINCT hostname) AS hosts,
       MIN(iso_local) AS first_seen, MAX(iso_local) AS last_seen,
       ROUND((MAX(unixtime_ms)-MIN(unixtime_ms))/60000.0,1) AS span_minutes
FROM snapshot;
'@
            }
        }
        'recent' {
            if ($isHeartbeat) {
                "SELECT id, started_iso, label, sequence, duration_ms, expected_source, drift_ms, outcome FROM beat ORDER BY started_ms DESC LIMIT $Limit;"
            } else {
                "SELECT id, iso_local, hostname, username, process_count, uptime_seconds, invocation_source FROM snapshot ORDER BY unixtime_ms DESC LIMIT $Limit;"
            }
        }
        'cadence' {
@'
WITH d AS (SELECT started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta FROM beat)
SELECT COUNT(*) AS intervals,
       MIN(delta) AS min_ms,
       CAST(AVG(delta) AS INTEGER) AS mean_ms,
       MAX(delta) AS max_ms
FROM d WHERE delta IS NOT NULL;
'@
        }
        'drift' {
            if ($script:AnchorMs -ne $null) {
                # Read-time derivation from raw start timestamps: reconstruct the
                # anchor + k*interval grid and measure each start against it.
                $im = [long]$intervalUsed * 1000
                $an = $script:AnchorMs
@"
WITH g AS (
  SELECT started_ms,
         started_ms - (CAST(ROUND((started_ms - $an) * 1.0 / $im) AS INTEGER) * $im + $an) AS d
  FROM beat)
SELECT 'anchor (read-time)' AS source,
       COUNT(*) AS n,
       MIN(d) AS min_ms,
       CAST(AVG(d) AS INTEGER) AS mean_ms,
       MAX(ABS(d)) AS max_abs_ms
FROM g;
"@
            } else {
@'
SELECT expected_source,
       COUNT(*) AS n,
       MIN(drift_ms) AS min_ms,
       CAST(AVG(drift_ms) AS INTEGER) AS mean_ms,
       MAX(ABS(drift_ms)) AS max_abs_ms
FROM beat
WHERE drift_ms IS NOT NULL
GROUP BY expected_source
ORDER BY expected_source;
'@
            }
        }
        'jitter' {
            if ($intervalUsed -le 0) {
                Write-Log "jitter needs an interval; pass -IntervalSeconds" -Level Error
                exit 2
            }
            $im = [long]$intervalUsed * 1000
@"
WITH o AS (
  SELECT started_ms % $im AS off FROM beat WHERE interval_seconds IS NOT NULL),
     a AS (SELECT AVG(off) AS phase FROM o)
SELECT COUNT(*) AS n,
       CAST((SELECT phase FROM a) AS INTEGER) AS phase_ms,
       CAST(MIN(off) - (SELECT phase FROM a) AS INTEGER) AS min_ms,
       CAST(MAX(off) - (SELECT phase FROM a) AS INTEGER) AS max_ms,
       CAST(MAX(off) - MIN(off) AS INTEGER) AS spread_ms
FROM o;
"@
        }
        'gaps' {
            $threshold = if ($intervalUsed -gt 0) { $intervalUsed * 2000 } else { 0 }
            if ($threshold -le 0) {
                Write-Log "cannot detect gaps without an interval; pass -IntervalSeconds" -Level Error
                exit 2
            }
@"
WITH d AS (
  SELECT started_iso, started_ms,
         started_ms - LAG(started_ms) OVER (ORDER BY started_ms) AS delta
  FROM beat)
SELECT started_iso AS resumed_at,
       delta AS gap_ms,
       ROUND(delta/1000.0,1) AS gap_seconds,
       CAST(ROUND(delta/1000.0/$intervalUsed) AS INTEGER) AS intervals_missed
FROM d WHERE delta IS NOT NULL AND delta > $threshold
ORDER BY started_ms LIMIT $Limit;
"@
        }
        'overlaps' {
@"
SELECT a.id AS run_a, b.id AS run_b,
       a.started_iso AS a_started, b.started_iso AS b_started,
       a.duration_ms AS a_duration_ms, b.duration_ms AS b_duration_ms
FROM beat a JOIN beat b
  ON a.id < b.id
 AND a.started_ms <= b.finished_ms
 AND b.started_ms <= a.finished_ms
ORDER BY a.started_ms LIMIT $Limit;
"@
        }
        'failures' {
            "SELECT id, started_iso, label, exit_code, outcome, duration_ms FROM beat WHERE exit_code <> 0 ORDER BY started_ms DESC LIMIT $Limit;"
        }
        'restarts' {
@"
SELECT session_id, pid, COUNT(*) AS beats,
       MIN(started_iso) AS first_beat,
       datetime(MAX(started_ms)/1000,'unixepoch') AS last_beat_utc
FROM beat GROUP BY session_id, pid ORDER BY MIN(started_ms) DESC LIMIT $Limit;
"@
        }
        'hosts' {
            if ($isHeartbeat) {
                "SELECT hostname, username, COUNT(*) AS beats FROM beat GROUP BY hostname, username ORDER BY beats DESC LIMIT $Limit;"
            } else {
                "SELECT hostname, username, COUNT(*) AS snapshots FROM snapshot GROUP BY hostname, username ORDER BY snapshots DESC LIMIT $Limit;"
            }
        }
        'listeners' {
            # Read the most recent snapshot whose port probe actually RAN.
            #
            # Reading "the newest snapshot" unconditionally was wrong: when the
            # newest one came from a host or twin where no port tool was
            # available, the query returned nothing -- indistinguishable from
            # "this machine is listening on no ports". That is the exact
            # conflation of "could not determine" with "determined zero" that
            # the rest of this schema goes out of its way to avoid.
            #
            # Legacy snapshots (schema < 3) have no probe column. They are
            # treated as usable but flagged, because their provenance is
            # genuinely unknown rather than known-good.
            $usable = Invoke-Sqlite -Database $dbPath -Mode 'list' -Sql @'
SELECT id, iso_local, COALESCE(ports_probe,'unknown')
FROM snapshot
WHERE ports_probe = 'ok' OR ports_probe IS NULL
ORDER BY unixtime_ms DESC LIMIT 1;
'@
            $row = ($usable | Where-Object { $_ -match '^\d+\|' } | Select-Object -First 1)
            if (-not $row) {
                Write-Log "No snapshot has usable port data." -Level Error
                Write-Log "Every snapshot recorded ports_probe = 'unavailable' or 'skipped'." -Level Error
                Write-Log "Re-run Test-GetSystemInfo without -SkipPorts on a host with ss, netstat, or Get-NetTCPConnection." -Level Error
                Write-Output "no snapshot with usable port data"
                exit 0
            }
            $parts  = $row -split '\|'
            $sid    = $parts[0]
            $siso   = $parts[1]
            $sprobe = $parts[2]

            $newest = Invoke-Sqlite -Database $dbPath -Mode 'list' -Sql @'
SELECT id, iso_local, COALESCE(ports_probe,'unknown')
FROM snapshot ORDER BY unixtime_ms DESC LIMIT 1;
'@
            $nrow = ($newest | Where-Object { $_ -match '^\d+\|' } | Select-Object -First 1)
            $nparts = $nrow -split '\|'
            Write-Log ("showing snapshot {0} ({1}), ports_probe={2}" -f $sid, $siso, $sprobe)
            if ($nparts[0] -ne $sid) {
                Write-Log ("the NEWEST snapshot is {0} ({1}) with ports_probe={2} -- not shown, because its port probe did not run" -f `
                    $nparts[0], $nparts[1], $nparts[2]) -Level Warn
            }
            if ($sprobe -eq 'unknown') {
                Write-Log "this snapshot predates probe-status recording (schema < 3); an empty result may mean the probe never ran" -Level Warn
            }

            # The comparison baseline must itself have usable port data,
            # otherwise every port reads as NEW against a snapshot that simply
            # never looked.
@"
WITH prev AS (
  SELECT id FROM snapshot
  WHERE (ports_probe = 'ok' OR ports_probe IS NULL) AND id <> $sid
    AND unixtime_ms <= (SELECT unixtime_ms FROM snapshot WHERE id = $sid)
  ORDER BY unixtime_ms DESC LIMIT 1)
SELECT p.protocol, p.port, p.address, p.process_name,
       CASE WHEN (SELECT id FROM prev) IS NULL THEN 'no-comparable-snapshot'
            WHEN EXISTS (SELECT 1 FROM snapshot_port q
                         WHERE q.snapshot_id=(SELECT id FROM prev)
                           AND q.protocol=p.protocol AND q.port=p.port)
            THEN 'unchanged' ELSE 'NEW' END AS since_previous
FROM snapshot_port p
WHERE p.snapshot_id = $sid
ORDER BY p.protocol, p.port LIMIT $Limit;
"@
        }
        'schema' { "SELECT sql FROM sqlite_master WHERE sql IS NOT NULL ORDER BY type DESC, name;" }
    }

    # Exclusion disclosure, before the results so it is not lost below them.
    if ($Query -eq 'drift') {
        Get-ExcludedNote -DbPath $dbPath `
            -CountSql "SELECT COUNT(*) FROM beat WHERE drift_ms IS NULL;" `
            -Reason "no expected moment was available (record with -AnchorIso to measure latency); try -Query jitter"
        # Legacy rows written before v0.5.1 used epoch-boundary snapping, which
        # is only correct for a schedule that happens to sit on that grid. This
        # scheduler anchors interval schedules to task creation time, so those
        # figures are dominated by a constant phase offset and are not latency.
        Get-ExcludedNote -DbPath $dbPath `
            -CountSql "SELECT COUNT(*) FROM beat WHERE expected_source='boundary';" `
            -Reason "LEGACY epoch-boundary rows from before v0.5.1 -- these are phase offset, NOT latency; re-record with -AnchorIso"
    }
    if ($Query -eq 'jitter') {
        Write-Log "jitter measures variation around the schedule's OWN observed phase." -Level Warn
        Write-Log "It cannot detect uniform lateness: a scheduler consistently late by a" -Level Warn
        Write-Log "fixed amount has zero jitter. For absolute latency, record with -AnchorIso." -Level Warn
    }
    if ($Query -eq 'gaps' -and $intervalInferred) {
        Write-Log "the interval above was inferred, not supplied -- treat the result accordingly" -Level Warn
    }

    $result = Invoke-Sqlite -Database $dbPath -Sql $sql -Mode $mode
    $result | ForEach-Object { Write-Output $_ }
    exit 0

#_______________________________________________________________________________
## End of script
