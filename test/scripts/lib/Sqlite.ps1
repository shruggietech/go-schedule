<#
.SYNOPSIS
    Shared sqlite3 resolution, installation, and access for the go-schedule
    maintainer test scripts.

.DESCRIPTION
    Dot-sourced by Test-Heartbeat.ps1, Test-GetSystemInfo.ps1, and
    Test-ReadTestDB.ps1. This is the single implementation of everything those
    three share: locating a usable sqlite3, optionally installing a pinned one,
    resolving the per-user test data directory, opening databases in WAL mode,
    executing statements with bound parameters, creating the schemas, and
    logging.

    There is exactly one implementation on purpose. Three copies of the
    resolution order would be three chances for them to disagree, and a
    disagreement here does not present as an obvious bug -- it presents as an
    intermittent, platform-specific test failure that costs a day to find.

    Exit-code contract, shared by every script that dot-sources this file:
      0  success
      1  runtime failure (probe failed, write failed, contention exhausted)
      2  usage error or unmet prerequisite (bad arguments, no usable sqlite3,
         unsupported platform for the installer)

    Results go to stdout; diagnostics and warnings go to stderr. This deviates
    from the ShruggieTech Write-Log fixture, which uses Write-Host: the
    go-schedule constitution (principle III) requires diagnostics on stderr, and
    where the local project standard and the house standard disagree, the local
    one wins. The level names, ordering, and suppression semantics are otherwise
    unchanged.

.PARAMETER Help
    Print this help text to the terminal.
    Alias: h

.EXAMPLE
    . "$PSScriptRoot/lib/Sqlite.ps1"
    The only intended use: dot-sourced from a sibling test script.

.NOTES
    Requires PowerShell 7+. Requires sqlite3 3.33.0 or later at runtime; see
    lib/sqlite-manifest.json for why that specific floor.
#>
[CmdletBinding(SupportsShouldProcess=$false,ConfirmImpact='None',DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$false,ParameterSetName='HelpText')]
    [Alias("h")]
    [Switch]$Help
)
#_______________________________________________________________________________
## Declare Functions

    function Assert-PSVersion {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$false)]
            [version]$Minimum = '7.0'
        )
        $current = $PSVersionTable.PSVersion
        if ($current -lt $Minimum) {
            [Console]::Error.WriteLine(
                ("ALERT: PowerShell {0}+ required; running {1}. Relaunch with 'pwsh'." -f $Minimum, $current))
            exit 2
        }
    }

    function Write-Log {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true,Position=0)]
            [string]$Message,

            [Parameter(Mandatory=$false)]
            [ValidateSet('Info','Success','Warn','Error','Debug')]
            [string]$Level = 'Info',

            [Parameter(Mandatory=$false)]
            [string]$Source = $null
        )
        if ($script:LogSilent -and $Level -ne 'Error') { return }
        if ($script:LogQuiet -and (@('Info','Success','Debug') -contains $Level)) { return }

        $stamp = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss.fff')
        $tag   = if ($Source) { "[$Source] " } else { '' }
        $label = $Level.ToUpper().PadRight(7)
        [Console]::Error.WriteLine(("{0} {1}{2} {3}" -f $stamp, $tag, $label, $Message))
    }

    function Get-TestPlatform {
        [CmdletBinding()]
        Param()
        if ($IsWindows) { $os = 'windows' }
        elseif ($IsMacOS) { $os = 'darwin' }
        else { $os = 'linux' }
        $arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
            'Arm64' { 'arm64' }
            default { 'x64' }
        }
        return @{ OS = $os; Arch = $arch; Key = "$os-$arch" }
    }

    function Test-SqliteCandidate {
        # Returns the version string when the candidate is usable and new enough,
        # otherwise $null. A candidate that exists but is broken or too old must
        # be treated as not-found so the search continues -- a stale sqlite3
        # early in the order must never shadow a good one later in it.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Candidate,

            [Parameter(Mandatory=$false)]
            [version]$Minimum = '3.33.0'
        )
        $probe = [Diagnostics.ProcessStartInfo]::new()
        $probe.FileName = $Candidate
        $probe.ArgumentList.Add('--version')
        $probe.UseShellExecute = $false
        $probe.CreateNoWindow = $true
        $probe.RedirectStandardInput = $true
        $probe.RedirectStandardOutput = $true
        $probe.RedirectStandardError = $true
        $probe.Environment['GIT_TERMINAL_PROMPT'] = '0'
        $probe.Environment['GCM_INTERACTIVE'] = 'Never'
        $process = $null
        $clock = [Diagnostics.Stopwatch]::StartNew()
        try {
            $process = [Diagnostics.Process]::Start($probe)
            $process.StandardInput.Close()
            # Fixed buffers prevent a broken candidate from filling memory.
            $outputBuffer = [char[]]::new(4097)
            $errorBuffer = [char[]]::new(4097)
            # Incremental reads retain completed chunks if another read is still
            # pending at the shared deadline (including inherited pipe handles).
            $channels = @(
                @{ Reader=$process.StandardOutput; Buffer=$outputBuffer; Filled=0; Done=$false; Pending=$null },
                @{ Reader=$process.StandardError; Buffer=$errorBuffer; Filled=0; Done=$false; Pending=$null }
            )
            foreach ($channel in $channels) {
                $channel.Pending = $channel.Reader.ReadAsync($channel.Buffer, 0, $channel.Buffer.Length)
            }
            $timedOut = $false
            while (@($channels | Where-Object { -not $_.Done }).Count -gt 0) {
                foreach ($channel in $channels) {
                    if (-not $channel.Done -and $channel.Pending.IsCompleted) {
                        $read = $channel.Pending.GetAwaiter().GetResult()
                        $channel.Pending = $null
                        $channel.Filled += $read
                        if ($read -eq 0 -or $channel.Filled -eq $channel.Buffer.Length) {
                            $channel.Done = $true
                        } else {
                            $channel.Pending = $channel.Reader.ReadAsync($channel.Buffer, $channel.Filled, $channel.Buffer.Length - $channel.Filled)
                        }
                    }
                }
                $pending = @($channels | Where-Object { -not $_.Done })
                if ($pending.Count -eq 0) { break }
                $remaining = [Math]::Max(0, 10000 - [int]$clock.ElapsedMilliseconds)
                if ($remaining -eq 0 -or [Threading.Tasks.Task]::WaitAny([Threading.Tasks.Task[]]@($pending | ForEach-Object { $_.Pending }), $remaining) -eq -1) {
                    $timedOut = $true
                    break
                }
            }
            if (-not $timedOut) {
                $remaining = [Math]::Max(0, 10000 - [int]$clock.ElapsedMilliseconds)
                $timedOut = -not $process.WaitForExit($remaining)
            }
            if ($timedOut -and -not $process.HasExited) { $process.Kill($true); $process.WaitForExit() }
            foreach ($channel in $channels) {
                if ($channel.Pending -and $channel.Pending.IsCompletedSuccessfully) {
                    $channel.Filled += $channel.Pending.GetAwaiter().GetResult()
                }
            }
            $outputCount = $channels[0].Filled
            $errorCount = $channels[1].Filled
            $raw = [string]::new($outputBuffer, 0, [Math]::Min($outputCount, 4096))
            $detail = [string]::new($errorBuffer, 0, [Math]::Min($errorCount, 4096))
            $detail = [regex]::Replace($detail, '[\x00-\x1f\x7f]', ' ').Trim()
            if ($timedOut) {
                Write-Log ("sqlite3 candidate '{0}' timed out after 10000 ms. stderr: {1}" -f $Candidate, $detail) -Level Warn
                return $null
            }
            if ($process.ExitCode -ne 0) {
                Write-Log ("sqlite3 candidate '{0}' exited {1} (0x{1:X8}). stderr: {2}" -f $Candidate, $process.ExitCode, $detail) -Level Warn
                return $null
            }
            if ($outputCount -gt 4096 -or $errorCount -gt 4096) {
                Write-Log ("sqlite3 candidate '{0}' exceeded the 4096-character probe output limit." -f $Candidate) -Level Warn
                return $null
            }
        } catch {
            Write-Log ("sqlite3 candidate '{0}' could not be launched or read: {1}" -f $Candidate, $_.Exception.Message) -Level Warn
            return $null
        } finally {
            $clock.Stop()
            if ($process) { if (-not $process.HasExited) { $process.Kill($true); $process.WaitForExit() }; $process.Dispose() }
        }
        $token = $raw.Trim().Split(' ')[0]
        if ($token -notmatch '^\d+\.\d+\.\d+$') {
            Write-Log ("sqlite3 candidate '{0}' returned an empty or malformed version." -f $Candidate) -Level Warn
            return $null
        }
        try { $found = [version]$token } catch { return $null }
        if ($found -lt $Minimum) {
            Write-Log ("Ignoring sqlite3 at '{0}': version {1} is below the required {2}." -f $Candidate, $found, $Minimum) -Level Debug
            return $null
        }
        return $token
    }

    function Resolve-Sqlite {
        # Strict precedence: explicit path, then repo-local .bin/, then PATH.
        #
        # The fall-through-on-unusable rule (FR-016a) applies only to the
        # IMPLICIT candidates, where continuing past a stale .bin/ to a good
        # tool on PATH is the helpful thing to do. An EXPLICIT path is
        # different: the maintainer named one specific tool, and quietly
        # running a different one instead is how you debug the wrong binary for
        # an hour. A bad explicit path is a hard usage error.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$false)]
            [string]$Explicit
        )
        $exeName = if ($IsWindows) { 'sqlite3.exe' } else { 'sqlite3' }

        if ($Explicit) {
            if (-not (Test-Path -LiteralPath $Explicit -PathType Leaf)) {
                Write-Log ("-SqliteExe '{0}' does not exist." -f $Explicit) -Level Error
                Write-Log "Refusing to silently fall back to another sqlite3 when you named a specific one." -Level Error
                exit 2
            }
            $ver = Test-SqliteCandidate -Candidate $Explicit
            if (-not $ver) {
                Write-Log ("-SqliteExe '{0}' is not runnable or is older than 3.33.0." -f $Explicit) -Level Error
                exit 2
            }
            Write-Log ("Using sqlite3 {0} at {1}" -f $ver, $Explicit) -Level Debug
            return $Explicit
        }

        $candidates = @()
        $candidates += (Join-Path $script:TestScriptsRoot (Join-Path '.bin' $exeName))
        $onPath = Get-Command 'sqlite3' -ErrorAction SilentlyContinue
        if ($onPath) { $candidates += $onPath.Source }

        foreach ($c in $candidates) {
            if (-not $c) { continue }
            if ($c -ne 'sqlite3' -and -not (Test-Path -LiteralPath $c -PathType Leaf)) { continue }
            $ver = Test-SqliteCandidate -Candidate $c
            if ($ver) {
                Write-Log ("Using sqlite3 {0} at {1}" -f $ver, $c) -Level Debug
                return $c
            }
        }
        return $null
    }

    function Exit-NoSqlite {
        # Unmet prerequisite is exit 2, not 1. The message names both remedies
        # because a maintainer who hits this needs to act, not to investigate.
        [CmdletBinding()]
        Param()
        $p = Get-TestPlatform
        $hint = switch ($p.OS) {
            'windows' { 'winget install SQLite.SQLite' }
            'darwin'  { 'brew install sqlite' }
            default   { 'sudo apt install sqlite3' }
        }
        Write-Log "No usable sqlite3 found (need 3.33.0 or later)." -Level Error
        Write-Log "Fix it either way:" -Level Error
        Write-Log ("  - rerun this script with -InstallSqlite (downloads the pinned build, verifies its checksum)") -Level Error
        Write-Log ("  - or install it yourself: {0}" -f $hint) -Level Error
        exit 2
    }

    function Install-Sqlite {
        # Verification precedes unpacking. Nothing unverified is ever placed
        # where the resolution order would later find it.
        [CmdletBinding()]
        Param()
        $manifestPath = Join-Path $PSScriptRoot 'sqlite-manifest.json'
        if (-not (Test-Path -LiteralPath $manifestPath -PathType Leaf)) {
            Write-Log "Installer manifest not found: $manifestPath" -Level Error
            exit 1
        }
        $manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json
        $p = Get-TestPlatform
        $entry = $manifest.platforms.($p.Key)
        if (-not $entry) {
            Write-Log ("No prebuilt sqlite3 is published for {0}." -f $p.Key) -Level Error
            Write-Log ("Install it from your package manager instead; this installer does not build from source.") -Level Error
            exit 2
        }
        if (-not $entry.sha256) {
            Write-Log "Manifest has no checksum for this platform; refusing to install." -Level Error
            Write-Log "An unverified binary is worse than no binary. Use your package manager." -Level Error
            exit 1
        }

        $binDir = Join-Path $script:TestScriptsRoot '.bin'
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
        $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("gosched-sqlite-{0}.zip" -f ([guid]::NewGuid().ToString('N')))

        Write-Log ("Downloading {0}" -f $entry.url)
        try {
            Invoke-WebRequest -Uri $entry.url -OutFile $tmp -UseBasicParsing
        } catch {
            Write-Log ("Download failed: {0}" -f $_.Exception.Message) -Level Error
            if (Test-Path -LiteralPath $tmp) { Remove-Item -LiteralPath $tmp -Force }
            exit 1
        }

        $actual = (Get-FileHash -LiteralPath $tmp -Algorithm SHA256).Hash.ToLower()
        if ($actual -ne $entry.sha256.ToLower()) {
            Remove-Item -LiteralPath $tmp -Force
            Write-Log "CHECKSUM MISMATCH -- download discarded and not installed." -Level Error
            Write-Log ("  expected {0}" -f $entry.sha256) -Level Error
            Write-Log ("  actual   {0}" -f $actual) -Level Error
            exit 1
        }
        Write-Log "Checksum verified." -Level Success

        $stage = Join-Path ([System.IO.Path]::GetTempPath()) ("gosched-sqlite-x-{0}" -f ([guid]::NewGuid().ToString('N')))
        Expand-Archive -LiteralPath $tmp -DestinationPath $stage -Force
        Remove-Item -LiteralPath $tmp -Force

        $found = Get-ChildItem -Path $stage -Recurse -Filter $entry.binary | Select-Object -First 1
        if (-not $found) {
            Remove-Item -LiteralPath $stage -Recurse -Force
            Write-Log ("Archive did not contain {0}." -f $entry.binary) -Level Error
            exit 1
        }
        $dest = Join-Path $binDir $entry.binary
        Copy-Item -LiteralPath $found.FullName -Destination $dest -Force
        Remove-Item -LiteralPath $stage -Recurse -Force
        if (-not $IsWindows) { & chmod +x $dest 2>$null }

        $ver = Test-SqliteCandidate -Candidate $dest
        if (-not $ver) {
            Write-Log "Installed binary does not run or is too old." -Level Error
            exit 1
        }
        Write-Log ("Installed sqlite3 {0} to {1}" -f $ver, $dest) -Level Success
        return $dest
    }

    function Get-TestDataDir {
        # A user-writable directory, never the daemon's system-wide data dir:
        # test payloads must run unelevated, and disposable test output does not
        # belong in the directory holding live scheduler state.
        [CmdletBinding()]
        Param()
        if ($env:GOSCHEDULE_TEST_DIR) { $base = $env:GOSCHEDULE_TEST_DIR }
        elseif ($IsWindows) { $base = Join-Path $env:LOCALAPPDATA 'goschedule-test' }
        elseif ($IsMacOS) { $base = Join-Path $HOME 'Library/Application Support/goschedule-test' }
        elseif ($env:XDG_DATA_HOME) { $base = Join-Path $env:XDG_DATA_HOME 'goschedule-test' }
        else { $base = Join-Path $HOME '.local/share/goschedule-test' }
        New-Item -ItemType Directory -Path $base -Force | Out-Null
        return $base
    }

    function Resolve-TestDatabase {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Name
        )
        if ($Name -match '[\\/]' -or $Name -match '\.db$') { return $Name }
        return (Join-Path (Get-TestDataDir) ("{0}.db" -f $Name))
    }

    function Invoke-Sqlite {
        # Bound parameters via '.param set'. The values written include the
        # hostname, the username, and network interface names -- all capable of
        # containing a quote, and all influenceable by whoever administers the
        # machine. String-interpolated SQL here would be both an injection
        # vector and an ordinary bug for any user named O'Brien.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Database,

            [Parameter(Mandatory=$true)]
            [string]$Sql,

            [Parameter(Mandatory=$false)]
            [hashtable]$Parameters = @{},

            [Parameter(Mandatory=$false)]
            [string]$Mode = 'list',

            [Parameter(Mandatory=$false)]
            [int]$Attempts = 3,

            [Parameter(Mandatory=$false)]
            [Switch]$ReadOnly
        )
        # PRAGMA statements return rows. Left unredirected they land in the
        # caller's result -- and 'busy_timeout=5000' emitting a bare 5000 was
        # being parsed downstream as a record count, producing a confident
        # "excluded 5000 records" warning out of thin air. Silence them at the
        # source rather than filtering them out afterwards.
        $nullDevice = if ($IsWindows) { 'nul' } else { '/dev/null' }
        $lines = @(
            '.bail on'
            ('.output {0}' -f $nullDevice)
            'PRAGMA busy_timeout=5000;'
        )
        if (-not $ReadOnly) {
            $lines += @('PRAGMA foreign_keys=ON;', 'PRAGMA journal_mode=WAL;')
        }
        $lines += @(
            '.output stdout'
            ('.mode {0}' -f $Mode)
        )
        foreach ($k in $Parameters.Keys) {
            $v = $Parameters[$k]
            if ($null -eq $v) {
                $lines += ('.param set :{0} NULL' -f $k)
            } else {
                # Keep one CLI-quoted expression. SQLite otherwise evaluates
                # numeric-looking text (including session IDs) or dot commands.
                $encoded = [BitConverter]::ToString([Text.Encoding]::UTF8.GetBytes([string]$v)).Replace('-', '')
                $lines += ('.param set :{0} "CAST(X''{1}'' AS TEXT)"' -f $k, $encoded)
            }
        }
        $lines += $Sql
        $script = ($lines -join "`n")

        $databaseArguments = @($Database)
        if ($ReadOnly) { $databaseArguments = @('-readonly', $Database) }
        for ($i = 1; $i -le $Attempts; $i++) {
            $out = $script | & $script:ResolvedSqliteExe @databaseArguments 2>&1
            if ($LASTEXITCODE -eq 0) { return $out }
            $text = ($out | Out-String)
            if ($text -match 'locked|busy') {
                Write-Log ("Database contention on attempt {0}/{1}; retrying. sqlite3: {2}" -f $i, $Attempts, $text.Trim()) -Level Warn
                Start-Sleep -Milliseconds (250 * $i)
                continue
            }
            if ($ReadOnly -and $text -match 'not a database|database disk image is malformed') {
                Write-Log ("sqlite3 read-only inspection failed: {0}. Use an empty path or explicitly reset the disposable metrics database. No conversion was attempted." -f $text.Trim()) -Level Error
                exit 1
            }
            Write-Log ("sqlite3 failed: {0}" -f $text.Trim()) -Level Error
            exit 1
        }
        Write-Log ("Database still contended after {0} attempts. The record was NOT written. Last sqlite3 error: {1}" -f $Attempts, $text.Trim()) -Level Error
        Write-Log "This is a test-harness failure, not a scheduler defect." -Level Error
        exit 1
    }

    function Get-MetricsShapeSql {
        [CmdletBinding()]
        Param()
        return @'
WITH expected(name,definition) AS (VALUES ('meta','createtablemeta(keytextprimarykey,valuetext)'),('beat','createtablebeat(idintegerprimarykeyautoincrement,session_idtextnotnull,sequenceintegernotnull,labeltext,hostnametextnotnull,usernametext,pidintegernotnull,started_msintegernotnull,started_isotextnotnull,finished_msintegernotnull,duration_msintegernotnull,expected_msinteger,expected_sourcetextnotnullcheck(expected_sourcein(''env'',''anchor'',''boundary'',''none'')),drift_msinteger,interval_secondsinteger,exit_codeintegernotnull,outcometextnotnullcheck(outcomein(''ok'',''failed'')),sched_envtext,check(finished_ms>=started_ms),check((drift_msisnull)=(expected_msisnull)),check((expected_source=''none'')=(expected_msisnull)),check((outcome=''ok'')=(exit_code=0)),unique(session_id,sequence))'),('beat_started','createindexbeat_startedonbeat(started_ms)'),('snapshot','createtablesnapshot(idintegerprimarykeyautoincrement,unixtime_msintegernotnull,iso_localtextnotnull,iso_utctextnotnull,tz_offset_minutesintegernotnull,hostnametextnotnull,usernametext,process_countinteger,uptime_secondsinteger,os_platformtextnotnull,os_releasetext,script_pidintegernotnull,script_flavortextnotnullcheck(script_flavorin(''powershell'',''posix'')),invocation_sourcetextnotnull,addresses_probetext,ports_probetext)'),('snapshot_time','createindexsnapshot_timeonsnapshot(unixtime_ms)'),('snapshot_address','createtablesnapshot_address(idintegerprimarykeyautoincrement,snapshot_idintegernotnullreferencessnapshot(id)ondeletecascade,familytextnotnullcheck(familyin(''ipv4'',''ipv6'')),addresstextnotnull,interfacetext,scopetext)'),('address_snapshot','createindexaddress_snapshotonsnapshot_address(snapshot_id)'),('snapshot_port','createtablesnapshot_port(idintegerprimarykeyautoincrement,snapshot_idintegernotnullreferencessnapshot(id)ondeletecascade,protocoltextnotnullcheck(protocolin(''tcp'',''udp'')),familytext,addresstext,portintegernotnullcheck(portbetween1and65535),process_nametext)'),('port_snapshot','createindexport_snapshotonsnapshot_port(snapshot_id,protocol,port)')) SELECT CASE WHEN (SELECT COUNT(*) FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%')=9 AND NOT EXISTS(SELECT 1 FROM expected e LEFT JOIN sqlite_schema s ON s.name=e.name WHERE s.sql IS NULL OR lower(replace(replace(replace(replace(s.sql,' ',''),char(10),''),char(13),''),char(9),''))<>e.definition) THEN 1 ELSE 0 END;
'@
    }

    function Assert-MetricsSchema {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)][string]$Database,
            [Parameter(Mandatory=$false)][bool]$AllowEmpty = $true
        )
        # Inspect through a read-only connection before WAL or persistent writes.
        if (-not [IO.File]::Exists($Database) -or ([IO.FileInfo]::new($Database)).Length -eq 0) {
            if ($AllowEmpty) { return }
        } else {
            $objects = Invoke-Sqlite -Database $Database -ReadOnly -Sql "SELECT COUNT(*) FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%';"
            if ([string]$objects -eq '0' -and $AllowEmpty) { return }
            $shape = Invoke-Sqlite -Database $Database -ReadOnly -Sql (Get-MetricsShapeSql)
            if ([string]$shape -eq '1') {
                $version = Invoke-Sqlite -Database $Database -ReadOnly -Sql "SELECT value FROM meta WHERE key='schema_version';"
                if ([string]$version -eq '4') { return }
            }
        }
        Write-Log 'Unsupported metrics database; use an empty path or explicitly reset the disposable metrics database. No conversion was attempted.' -Level Error
        exit 1
    }

    function Initialize-MetricsSchema {
        [CmdletBinding()]
        Param([Parameter(Mandatory=$true)][string]$Database)
        Assert-MetricsSchema -Database $Database
        # Serialize fresh creators and retain every already committed record.
        $sql = @'
BEGIN IMMEDIATE;
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT);
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
CREATE INDEX IF NOT EXISTS beat_started ON beat(started_ms);
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
CREATE INDEX IF NOT EXISTS port_snapshot ON snapshot_port(snapshot_id, protocol, port);
INSERT OR IGNORE INTO meta(key,value) VALUES('schema_version','4');
COMMIT;
'@
        Invoke-Sqlite -Database $Database -Sql $sql | Out-Null
    }

    function Initialize-HeartbeatSchema {
        [CmdletBinding()]
        Param([Parameter(Mandatory=$true)][string]$Database)
        Initialize-MetricsSchema -Database $Database
    }

    function Initialize-SystemSchema {
        [CmdletBinding()]
        Param([Parameter(Mandatory=$true)][string]$Database)
        Initialize-MetricsSchema -Database $Database
    }

    function Initialize-TestSqlite {
        # Called by every script after argument parsing. Establishes
        # $script:ResolvedSqliteExe or exits 2.
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$false)][string]$SqliteExe,
            [Parameter(Mandatory=$false)][switch]$InstallSqlite
        )
        Assert-PSVersion
        if ($InstallSqlite) {
            $script:ResolvedSqliteExe = Install-Sqlite
            return
        }
        $found = Resolve-Sqlite -Explicit $SqliteExe
        if (-not $found) { Exit-NoSqlite }
        $script:ResolvedSqliteExe = $found
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $script:LogQuiet        = $false
    $script:LogSilent       = $false
    $script:TestScriptsRoot = Split-Path -Parent $PSScriptRoot

    # Deliberately NOT named $script:SqliteExe. This file is dot-sourced, so
    # '$script:' resolves to the CALLER's scope -- and every caller has a
    # -SqliteExe parameter. Declaring $script:SqliteExe here silently
    # overwrote that parameter with $null before it was ever read, so an
    # explicitly requested sqlite3 was ignored and the one on PATH used
    # instead. Namespacing the internal name is the fix. The collision is
    # invisible otherwise: nothing errors, it just quietly does the wrong thing.
    $script:ResolvedSqliteExe = $null

#_______________________________________________________________________________
## Execute Operations

    # Catch help text requests
    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $PSCommandPath -Detailed
        exit 0
    }

    # This file is a library. Dot-sourcing it defines the functions above and
    # does nothing else; there is deliberately no work performed at load time.

#_______________________________________________________________________________
## End of script
