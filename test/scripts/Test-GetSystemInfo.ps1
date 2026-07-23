<#
.SYNOPSIS
    Record a snapshot of this machine into a local SQLite database.

.DESCRIPTION
    Captures host identity, logged-in user, running process count, uptime,
    network addresses, and listening ports, and writes them as one snapshot with
    attached detail rows.

    As a scheduler test this exercises a materially different execution profile
    from the heartbeat: it spawns subprocesses, calls platform-specific tooling,
    and writes multiple related rows. That is where cross-platform execution
    bugs actually surface. Run on a schedule it also becomes a rudimentary
    machine-state history worth diffing.

    Probes degrade rather than abort. A probe that cannot run on this host
    records NULL and warns; the snapshot is still written. NULL means "could not
    determine" and is never used for a legitimate zero, because a process count
    of zero and an unavailable process count support opposite conclusions.

    The Windows-only networking cmdlets (Get-NetIPAddress, Get-NetTCPConnection,
    Get-NetUDPEndpoint, Get-Uptime) do not exist in PowerShell on Linux or
    macOS, even though PowerShell itself runs there. This script therefore
    branches on $IsWindows and falls through to the same POSIX tools its shell
    twin uses. That branch is the likeliest source of a works-on-my-machine
    defect in this feature and is deliberately called out here.

    Exit codes: 0 the snapshot was recorded (even with degraded probes),
    1 the snapshot could not be written, 2 usage error or unmet prerequisite.

.PARAMETER Database
    Well-known name ('system') or an explicit path.
    Default: 'system'.
    Alias: d

.PARAMETER InvocationSource
    Free-form tag recorded on the snapshot, so snapshots produced by different
    scheduled tasks stay distinguishable.
    Default: 'manual'.
    Alias: i

.PARAMETER SkipPorts
    Skip the listening-port probe. It is the slowest, and on most platforms it
    is the one that wants elevation to attribute sockets to processes.
    Alias: s

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
    .\Test-GetSystemInfo.ps1
    Record one snapshot of this machine.

.EXAMPLE
    .\Test-GetSystemInfo.ps1 -InvocationSource hourly-audit -SkipPorts
    What a scheduled task might invoke when only host state matters.

.NOTES
    Requires PowerShell 7+ and sqlite3 3.33.0 or later.
    Full documentation: docs/test-scripts.md
#>
[CmdletBinding(SupportsShouldProcess=$false,ConfirmImpact='None',DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("d")]
    [string]$Database = 'system',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("i")]
    [string]$InvocationSource = 'manual',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("s")]
    [Switch]$SkipPorts,

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

    function Test-HasCommand {
        # Guard every external-tool call with this.
        #
        # `& ss -lntup 2>$null` does NOT degrade gracefully when ss is absent:
        # PowerShell raises CommandNotFoundException, which 2>$null does not
        # suppress, so the outer catch fires and the whole fallback chain is
        # skipped. On macOS -- no ss, but netstat present -- that turned an
        # answerable probe into 'unavailable'. Windows never reaches this code
        # and Linux has ss, so only macOS shows it.
        [CmdletBinding()]
        Param([Parameter(Mandatory=$true)][string]$Name)
        return [bool](Get-Command $Name -CommandType Application -ErrorAction SilentlyContinue)
    }

    function Get-ProcessCount {
        [CmdletBinding()]
        Param()
        try { return (Get-Process -ErrorAction Stop).Count }
        catch { Write-Log "process count unavailable" -Level Warn; return $null }
    }

    function Get-UptimeSeconds {
        [CmdletBinding()]
        Param()
        try {
            if ($IsWindows) {
                return [long](Get-Uptime -ErrorAction Stop).TotalSeconds
            }
            if (Test-Path '/proc/uptime') {
                return [long][double](((Get-Content '/proc/uptime' -Raw).Trim() -split '\s+')[0])
            }
            $boot = (& sysctl -n kern.boottime 2>$null)
            if ($boot -match 'sec\s*=\s*(\d+)') {
                return [long]([DateTimeOffset]::UtcNow.ToUnixTimeSeconds() - [long]$Matches[1])
            }
        } catch { }
        Write-Log "uptime unavailable" -Level Warn
        return $null
    }

    function Get-Addresses {
        # Fixed, documented fallback order. An ordering that varies by machine
        # produces data whose provenance differs between hosts while looking
        # identical in the database.
        [CmdletBinding()]
        Param()
        $rows = @()
        try {
            if ($IsWindows) {
                foreach ($a in (Get-NetIPAddress -ErrorAction Stop)) {
                    $rows += @{
                        family    = if ($a.AddressFamily -eq 'IPv4') { 'ipv4' } else { 'ipv6' }
                        address   = $a.IPAddress
                        interface = $a.InterfaceAlias
                        scope     = "$($a.PrefixOrigin)"
                    }
                }
                $script:AddressesProbe = 'ok'
                return $rows
            }
            $out = $null
            if (Test-HasCommand 'ip')       { $out = & ip -o addr 2>$null }
            if (-not $out -and (Test-HasCommand 'ifconfig')) { $out = & ifconfig 2>$null }
            if (-not $out) {
                Write-Log "address probe unavailable (no ip/ifconfig)" -Level Warn
                $script:AddressesProbe = 'unavailable'
                return @()
            }
            foreach ($line in $out) {
                if ($line -match 'inet6?\s+([0-9a-fA-F:.]+)') {
                    $addr = $Matches[1]
                    $iface = if ($line -match '^\d+:\s*(\S+)') { $Matches[1] } else { $null }
                    $rows += @{
                        family    = if ($addr -match ':') { 'ipv6' } else { 'ipv4' }
                        address   = $addr
                        interface = $iface
                        scope     = $null
                    }
                }
            }
        } catch {
            Write-Log "address probe unavailable" -Level Warn
            $script:AddressesProbe = 'unavailable'
            return @()
        }
        # 'ok' means the probe ran. Zero rows from a probe that ran is a real
        # answer; zero rows from one that could not run is not, and the two must
        # not collapse into the same record.
        if ($rows.Count -eq 0) { Write-Log "no addresses determined" -Level Warn }
        $script:AddressesProbe = 'ok'
        return $rows
    }

    function Get-ListeningPorts {
        [CmdletBinding()]
        Param()
        $rows = @()
        try {
            if ($IsWindows) {
                foreach ($c in (Get-NetTCPConnection -State Listen -ErrorAction Stop)) {
                    $rows += @{
                        protocol = 'tcp'; family = $null
                        address  = $c.LocalAddress; port = [int]$c.LocalPort
                        process  = $null
                    }
                }
                try {
                    foreach ($u in (Get-NetUDPEndpoint -ErrorAction Stop)) {
                        $rows += @{
                            protocol = 'udp'; family = $null
                            address  = $u.LocalAddress; port = [int]$u.LocalPort
                            process  = $null
                        }
                    }
                } catch { Write-Log "udp endpoint probe unavailable" -Level Warn }
                $script:PortsProbe = 'ok'
                return $rows
            }
            $out = $null
            if (Test-HasCommand 'ss')      { $out = & ss -lntup 2>$null }
            if (-not $out -and (Test-HasCommand 'netstat')) { $out = & netstat -an 2>$null }
            if (-not $out -and (Test-HasCommand 'lsof'))    { $out = & lsof -nP -iTCP -sTCP:LISTEN 2>$null }
            if (-not $out) {
                Write-Log "port probe unavailable (no ss/netstat/lsof)" -Level Warn
                $script:PortsProbe = 'unavailable'
                return @()
            }
            foreach ($line in $out) {
                if ($line -match '^(tcp|udp)\S*\s+.*?(\S+):(\d+)\s') {
                    $rows += @{
                        protocol = $Matches[1]; family = $null
                        address  = $Matches[2]; port = [int]$Matches[3]
                        process  = $null
                    }
                }
            }
        } catch {
            Write-Log "port probe unavailable" -Level Warn
            $script:PortsProbe = 'unavailable'
            return @()
        }
        $script:PortsProbe = 'ok'
        return $rows
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $ThisScriptPath = $MyInvocation.MyCommand.Path

    # Probe status vocabulary: 'ok' (ran; the result may legitimately be zero
    # rows), 'unavailable' (no tool on this host could answer), 'skipped' (the
    # caller asked not to). Recorded so a reader can tell an empty result from
    # an unanswerable one -- which is the whole point of NULL meaning "could not
    # determine" everywhere else in this schema.
    $script:AddressesProbe = 'unavailable'
    $script:PortsProbe     = 'skipped'

#_______________________________________________________________________________
## Execute Operations

    # Catch help text requests
    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }

    . (Join-Path $PSScriptRoot 'lib/Sqlite.ps1')
    if ($Quiet) { $script:LogQuiet = $true }

    Initialize-TestSqlite -SqliteExe $SqliteExe -InstallSqlite:$InstallSqlite

    $dbPath = Resolve-TestDatabase -Name $Database
    Write-Log ("database: {0}" -f $dbPath)
    Initialize-SystemSchema -Database $dbPath

    $now      = [DateTimeOffset]::Now
    $platform = if ($IsWindows) { 'windows' } elseif ($IsMacOS) { 'darwin' } else { 'linux' }

    $snapSql = @'
INSERT INTO snapshot (unixtime_ms, iso_local, iso_utc, tz_offset_minutes, hostname,
                      username, process_count, uptime_seconds, os_platform, os_release,
                      script_pid, script_flavor, invocation_source,
                      addresses_probe, ports_probe)
VALUES (CAST(:ms AS INTEGER), :isolocal, :isoutc, CAST(:tz AS INTEGER), :host,
        :user,
        CASE WHEN :procs = '' THEN NULL ELSE CAST(:procs AS INTEGER) END,
        CASE WHEN :uptime = '' THEN NULL ELSE CAST(:uptime AS INTEGER) END,
        :platform, :release,
        CAST(:pid AS INTEGER), 'powershell', :source, :aprobe, :pprobe);
SELECT last_insert_rowid();
'@
    $procs  = Get-ProcessCount
    $uptime = Get-UptimeSeconds
    # Probes run before the insert so their status is part of the snapshot row
    # rather than a second write that could fail independently.
    $addresses = Get-Addresses
    $ports = @()
    if ($SkipPorts) {
        $script:PortsProbe = 'skipped'
    } else {
        $ports = Get-ListeningPorts
    }
    $params = @{
        ms       = [long]$now.ToUnixTimeMilliseconds()
        isolocal = $now.ToString('yyyy-MM-ddTHH:mm:ss.fffzzz')
        isoutc   = $now.ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
        tz       = [int]$now.Offset.TotalMinutes
        host     = [System.Net.Dns]::GetHostName()
        user     = ($env:USERNAME ?? $env:USER)
        procs    = if ($null -eq $procs) { '' } else { $procs }
        uptime   = if ($null -eq $uptime) { '' } else { $uptime }
        platform = $platform
        release  = [System.Runtime.InteropServices.RuntimeInformation]::OSDescription
        pid      = $PID
        source   = $InvocationSource
        aprobe   = $script:AddressesProbe
        pprobe   = $script:PortsProbe
    }
    $out = Invoke-Sqlite -Database $dbPath -Sql $snapSql -Parameters $params -Mode 'list'
    $snapshotId = ($out | Where-Object { $_ -match '^\d+$' } | Select-Object -Last 1)
    if (-not $snapshotId) {
        Write-Log "snapshot insert did not return an id; nothing was recorded." -Level Error
        exit 1
    }

    # Child rows only ever written for a snapshot that was successfully
    # inserted, so a partial failure never orphans them.
    foreach ($a in $addresses) {
        Invoke-Sqlite -Database $dbPath -Mode 'list' -Sql @'
INSERT INTO snapshot_address (snapshot_id, family, address, interface, scope)
VALUES (CAST(:sid AS INTEGER), :family, :address,
        CASE WHEN :iface = '' THEN NULL ELSE :iface END,
        CASE WHEN :scope = '' THEN NULL ELSE :scope END);
'@ -Parameters @{
            sid = $snapshotId; family = $a.family; address = $a.address
            iface = ($a.interface ?? ''); scope = ($a.scope ?? '')
        } | Out-Null
    }

    $portCount = 0
    if (-not $SkipPorts) {
        foreach ($p in $ports) {
            if ($p.port -lt 1 -or $p.port -gt 65535) { continue }
            Invoke-Sqlite -Database $dbPath -Mode 'list' -Sql @'
INSERT INTO snapshot_port (snapshot_id, protocol, family, address, port, process_name)
VALUES (CAST(:sid AS INTEGER), :proto, NULL,
        CASE WHEN :address = '' THEN NULL ELSE :address END,
        CAST(:port AS INTEGER), NULL);
'@ -Parameters @{
                sid = $snapshotId; proto = $p.protocol
                address = ($p.address ?? ''); port = $p.port
            } | Out-Null
            $portCount++
        }
    }

    Write-Log ("snapshot {0} recorded: {1} address(es) [{2}], {3} listening port(s) [{4}]" -f `
        $snapshotId, $addresses.Count, $script:AddressesProbe, $portCount, $script:PortsProbe) -Level Success
    Write-Output $dbPath
    exit 0

#_______________________________________________________________________________
## End of script
