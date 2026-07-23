<#
.SYNOPSIS
    Record a heartbeat into a local SQLite database, so a scheduled task's
    firing behavior can be measured after the fact.

.DESCRIPTION
    Writes one beat per invocation and exits. That default is deliberate: the
    point of this script is to verify that goschedd supplies the cadence, and a
    script that sleeps in a loop tests its own Start-Sleep rather than the
    scheduler.

    Each beat records when the run started and finished, the host and process
    that produced it, a per-invocation session identity, the outcome, and -- the
    reason the script exists -- how far the firing landed from when it was
    expected. The expected moment comes from one of three sources, in
    precedence order, and the source is recorded alongside every drift figure:

      env       GOSCHED_SCHEDULED_TIME from the environment. The executor does
                not currently set this; the tier exists so that a future release
                which does is consumed with no change here.
      anchor    -AnchorIso plus -IntervalSeconds. One real firing time from
                the schedule reconstructs the whole grid (anchor + k*interval),
                so the difference is genuine absolute dispatch latency. Get the
                anchor from `gosched task show <id>`.
      none      Neither available. No expected moment and no drift are recorded,
                and the reader excludes the beat from drift statistics. A drift
                of zero for an unmeasurable run would be worse than nothing.

    There is deliberately no epoch-boundary tier. Snapping to the nearest
    multiple of the interval from the Unix epoch is correct only if the schedule
    happens to sit on that grid -- and this scheduler anchors an interval
    schedule to the task's creation time, so "every 1 minute" created at :06
    fires at :06 forever. Epoch snapping then reports a constant offset as
    though it were lateness, indistinguishable from the real thing. An earlier
    version did exactly that and reported ~6.4s of "drift" for a scheduler that
    was in fact firing within a quarter second.

    The beat is written once, at the end of the run. Two writes per beat would
    double database contention to buy a mid-flight-crash record that a
    maintainer cannot act on differently from a missed firing anyway.

    Exit codes: 0 success, 1 runtime failure, 2 usage error or unmet
    prerequisite (no usable sqlite3).

.PARAMETER Database
    Well-known name ('heartbeat') or an explicit path. A well-known name
    resolves inside the per-user test directory, which is reported on stderr.
    Default: 'heartbeat'.
    Alias: d

.PARAMETER IntervalSeconds
    The schedule interval you registered this task with. Enables the reader's
    missed-firing detection, and -- together with -AnchorIso -- drift. On its
    own it is NOT enough to measure drift; see -AnchorIso.
    Alias: i

.PARAMETER Label
    Free-form tag recorded on each beat, so several scheduled tasks can share
    one database and stay distinguishable.
    Alias: l

.PARAMETER AnchorIso
    Any one firing time from the schedule, RFC 3339, taken from
    `gosched task show <id>` next-runs. Combined with -IntervalSeconds this
    reconstructs the full firing grid, so drift becomes true dispatch latency.
    Without it no drift is recorded at all -- see the note in the description.
    Alias: a

.PARAMETER Loop
    Opt-in continuous mode ("repeat"). Bounded by -MaxBeats, -DurationSeconds,
    or the 3600-second default. There is no unbounded form.
    Alias: r (not 'l' -- PowerShell aliases are case-insensitive, so it would
    collide with -Label)

.PARAMETER MaxBeats
    Stop after this many beats. Whichever bound is reached first wins.
    Alias: m

.PARAMETER DurationSeconds
    Stop after this many seconds. Checked between beats, so one deliberately
    slow run may overrun it -- interrupting a run mid-write would corrupt the
    record the bound exists to protect.
    Default: 3600 when -Loop is set and no other bound is given.
    Alias: t

.PARAMETER SleepSeconds
    Deliberately extend the run by this many seconds. This is how overlap
    policies get exercised: make the run longer than its interval and watch what
    queue_one, skip, and allow_concurrent each do.
    Alias: s

.PARAMETER FailWith
    Exit with this code after recording the beat, to exercise failure handling
    and alerting. Rejects 0 and 2, which are reserved for success and for unmet
    prerequisites -- an induced failure must never be mistakable for either.
    Alias: f

.PARAMETER SqliteExe
    Explicit path to sqlite3. Highest precedence in the search order.

.PARAMETER InstallSqlite
    Download the pinned sqlite3 build, verify its checksum, and install it into
    test/scripts/.bin/. The only option here that touches the network.

.PARAMETER Quiet
    Suppress informational output. Warnings and errors are never suppressed.
    Alias: q

.PARAMETER Help
    Print this help text to the terminal.
    Alias: h

.EXAMPLE
    .\Test-Heartbeat.ps1
    Record a single beat. The cheapest possible confidence check.

.EXAMPLE
    .\Test-Heartbeat.ps1 -IntervalSeconds 60 -AnchorIso 2026-07-23T12:08:06Z
    What a task registered with "every 1 minute" should invoke. The anchor comes
    from `gosched task show`; together they yield true dispatch latency.

.EXAMPLE
    .\Test-Heartbeat.ps1 -SleepSeconds 90 -IntervalSeconds 60
    A run deliberately longer than its interval, for overlap-policy testing.

.EXAMPLE
    .\Test-Heartbeat.ps1 -FailWith 1
    Record the beat, then report failure so the daemon's alerting reacts.

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
    [Alias("i")]
    [int]$IntervalSeconds = 0,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("l")]
    [string]$Label = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("a")]
    [string]$AnchorIso = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("r")]
    [Switch]$Loop,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("m")]
    [int]$MaxBeats = 0,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("t")]
    [int]$DurationSeconds = 0,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("s")]
    [int]$SleepSeconds = 0,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("f")]
    [int]$FailWith = 0,

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

    function Get-EpochMs {
        [CmdletBinding()]
        Param([Parameter(Mandatory=$false)][datetime]$When = (Get-Date))
        return [long]([DateTimeOffset]$When).ToUnixTimeMilliseconds()
    }

    function Resolve-ExpectedMoment {
        # Returns a hashtable: Source ('env'|'anchor'|'none'), ExpectedMs, DriftMs.
        # Precedence is strict, so a future release that starts exporting the
        # scheduled time is picked up here without any other change.
        #
        # There is deliberately no epoch-boundary tier. Snapping to the nearest
        # multiple of the interval from the Unix epoch is only correct when the
        # schedule happens to be aligned to that grid, and this scheduler
        # anchors an interval schedule to the task's creation time -- so
        # "every 1 minute" created at :06 fires at :06 forever. Epoch snapping
        # then reports a constant ~6s "drift" that is pure phase offset, not
        # lateness, and the script has no way to tell the two apart. A
        # measurement that is right only when you got lucky is worse than no
        # measurement, because it is reported with the same confidence either
        # way. Supply -AnchorIso to get a real one.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)][long]$StartedMs,
            [Parameter(Mandatory=$true)][int]$Interval,
            [Parameter(Mandatory=$false)][string]$Anchor
        )
        if ($env:GOSCHED_SCHEDULED_TIME) {
            try {
                $parsed = [DateTimeOffset]::Parse($env:GOSCHED_SCHEDULED_TIME)
                $ms = [long]$parsed.ToUnixTimeMilliseconds()
                return @{ Source = 'env'; ExpectedMs = $ms; DriftMs = ($StartedMs - $ms) }
            } catch {
                Write-Log "GOSCHED_SCHEDULED_TIME set but unparseable; ignoring it." -Level Warn
            }
        }
        if ($Anchor -and $Interval -gt 0) {
            try {
                $anchorMs = [long]([DateTimeOffset]::Parse($Anchor)).ToUnixTimeMilliseconds()
            } catch {
                Write-Log ("-AnchorIso '{0}' is not a parseable timestamp." -f $Anchor) -Level Error
                exit 2
            }
            # Any firing time from the schedule works as the anchor: they all
            # sit on the same anchor + k*interval grid, k signed.
            $intervalMs = [long]$Interval * 1000
            $k = [long][math]::Round([double]($StartedMs - $anchorMs) / $intervalMs)
            $expected = $anchorMs + ($k * $intervalMs)
            return @{ Source = 'anchor'; ExpectedMs = $expected; DriftMs = ($StartedMs - $expected) }
        }
        return @{ Source = 'none'; ExpectedMs = $null; DriftMs = $null }
    }

    function Get-SchedulerEnv {
        [CmdletBinding()]
        Param()
        $found = @{}
        foreach ($e in [System.Environment]::GetEnvironmentVariables().GetEnumerator()) {
            if ($e.Key -like 'GOSCHED_*') { $found[$e.Key] = [string]$e.Value }
        }
        return ($found | ConvertTo-Json -Compress)
    }

    function Write-Beat {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)][string]$DbPath,
            [Parameter(Mandatory=$true)][string]$SessionId,
            [Parameter(Mandatory=$true)][int]$Sequence,
            [Parameter(Mandatory=$true)][long]$StartedMs,
            [Parameter(Mandatory=$true)][long]$FinishedMs,
            [Parameter(Mandatory=$true)][int]$ExitCode
        )
        $expected = Resolve-ExpectedMoment -StartedMs $StartedMs -Interval $IntervalSeconds -Anchor $AnchorIso
        $startedIso = ([DateTimeOffset]::FromUnixTimeMilliseconds($StartedMs)).ToLocalTime().ToString('yyyy-MM-ddTHH:mm:ss.fffzzz')

        $sql = @'
INSERT INTO beat (session_id, sequence, label, hostname, username, pid,
                  started_ms, started_iso, finished_ms, duration_ms,
                  expected_ms, expected_source, drift_ms, interval_seconds,
                  exit_code, outcome, sched_env)
VALUES (:session, CAST(:seq AS INTEGER), :label, :host, :user, CAST(:pid AS INTEGER),
        CAST(:started AS INTEGER), :startediso, CAST(:finished AS INTEGER), CAST(:duration AS INTEGER),
        CASE WHEN :expected IS NULL THEN NULL ELSE CAST(:expected AS INTEGER) END,
        :source,
        CASE WHEN :drift IS NULL THEN NULL ELSE CAST(:drift AS INTEGER) END,
        CASE WHEN :interval = '0' THEN NULL ELSE CAST(:interval AS INTEGER) END,
        CAST(:exitcode AS INTEGER), :outcome, :schedenv);
'@
        $params = @{
            session    = $SessionId
            seq        = $Sequence
            label      = if ($Label) { $Label } else { $null }
            host       = [System.Net.Dns]::GetHostName()
            user       = $env:USERNAME ?? $env:USER
            pid        = $PID
            started    = $StartedMs
            startediso = $startedIso
            finished   = $FinishedMs
            duration   = ($FinishedMs - $StartedMs)
            expected   = $expected.ExpectedMs
            source     = $expected.Source
            drift      = $expected.DriftMs
            interval   = $IntervalSeconds
            exitcode   = $ExitCode
            outcome    = if ($ExitCode -eq 0) { 'ok' } else { 'failed' }
            schedenv   = (Get-SchedulerEnv)
        }
        Invoke-Sqlite -Database $DbPath -Sql $sql -Parameters $params | Out-Null

        $driftText = if ($null -eq $expected.DriftMs) {
            'drift n/a (no expected moment available)'
        } else {
            ("drift {0}ms (source: {1})" -f $expected.DriftMs, $expected.Source)
        }
        Write-Log ("beat {0} seq {1} recorded; {2}" -f $SessionId.Substring(0,8), $Sequence, $driftText) -Level Success
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $ThisScriptPath = $MyInvocation.MyCommand.Path
    $LoopDefaultDurationSeconds = 3600

#_______________________________________________________________________________
## Execute Operations

    # Catch help text requests
    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }

    . (Join-Path $PSScriptRoot 'lib/Sqlite.ps1')
    if ($Quiet) { $script:LogQuiet = $true }

    # Reserved codes. An induced failure that could return 0 or 2 would be
    # indistinguishable from a success or from a missing sqlite3.
    if ($FailWith -eq 2 -or ($PSBoundParameters.ContainsKey('FailWith') -and $FailWith -eq 0)) {
        Write-Log "-FailWith must not be 0 or 2; those codes are reserved for success and for unmet prerequisites." -Level Error
        exit 2
    }
    if ($SleepSeconds -lt 0 -or $IntervalSeconds -lt 0 -or $MaxBeats -lt 0 -or $DurationSeconds -lt 0) {
        Write-Log "Numeric options must not be negative." -Level Error
        exit 2
    }

    Initialize-TestSqlite -SqliteExe $SqliteExe -InstallSqlite:$InstallSqlite

    $dbPath = Resolve-TestDatabase -Name $Database
    Write-Log ("database: {0}" -f $dbPath)
    Initialize-HeartbeatSchema -Database $dbPath

    $sessionId = [guid]::NewGuid().ToString('N')
    $exitCode  = $FailWith

    if (-not $Loop) {
        # Default path: exactly one beat. The scheduler owns the cadence.
        $startedMs = Get-EpochMs
        if ($SleepSeconds -gt 0) { Start-Sleep -Seconds $SleepSeconds }
        Write-Beat -DbPath $dbPath -SessionId $sessionId -Sequence 1 `
                   -StartedMs $startedMs -FinishedMs (Get-EpochMs) -ExitCode $exitCode
        Write-Output $dbPath
        exit $exitCode
    }

    # Opt-in continuous mode. Bounded always: an unbounded loop launched under a
    # scheduler is a resource incident, not a test.
    $effectiveDuration = if ($DurationSeconds -gt 0) { $DurationSeconds }
                         elseif ($MaxBeats -gt 0) { 0 }
                         else { $LoopDefaultDurationSeconds }
    $loopStart = Get-EpochMs
    $sequence  = 0
    $cadence   = if ($IntervalSeconds -gt 0) { $IntervalSeconds } else { 1 }

    Write-Log ("loop mode: max beats {0}, duration {1}s, cadence {2}s" -f `
        $(if ($MaxBeats -gt 0) { $MaxBeats } else { 'unset' }),
        $(if ($effectiveDuration -gt 0) { $effectiveDuration } else { 'unset' }),
        $cadence)

    while ($true) {
        # Both bounds checked between beats; first to trip ends the loop.
        if ($MaxBeats -gt 0 -and $sequence -ge $MaxBeats) { break }
        if ($effectiveDuration -gt 0 -and ((Get-EpochMs) - $loopStart) -ge ([long]$effectiveDuration * 1000)) { break }

        $sequence++
        $startedMs = Get-EpochMs
        if ($SleepSeconds -gt 0) { Start-Sleep -Seconds $SleepSeconds }
        Write-Beat -DbPath $dbPath -SessionId $sessionId -Sequence $sequence `
                   -StartedMs $startedMs -FinishedMs (Get-EpochMs) -ExitCode 0

        if ($MaxBeats -gt 0 -and $sequence -ge $MaxBeats) { break }
        Start-Sleep -Seconds $cadence
    }

    Write-Log ("loop finished after {0} beat(s)" -f $sequence) -Level Success
    Write-Output $dbPath
    exit $exitCode

#_______________________________________________________________________________
## End of script
