<#
.SYNOPSIS
    Prove the Windows daemon's real LocalSystem IPC and execution boundary.

.DESCRIPTION
    Creates the repository's default local authorization group, installs the
    staged daemon as a Windows service, and proves manual, scheduled,
    nonzero-exit, and process-start-failure task runs through LocalSystem. The
    service and group are removed in a finally block. Run only on a disposable
    elevated Windows CI host. Ordinary-token IPC is covered by the separate
    native named-pipe regression because the CI runner is elevated.
#>
[CmdletBinding()]
Param(
    [Parameter(Mandatory=$true)]
    [string]$StageDirectory,

    [Parameter(Mandatory=$true)]
    [string]$EvidencePath
)

$ErrorActionPreference = 'Stop'
$groupName = 'goschedadmin'
$serviceName = 'goschedd'
$stage = [System.IO.Path]::GetFullPath($StageDirectory)
$evidence = [System.IO.Path]::GetFullPath($EvidencePath)
$cli = [System.IO.Path]::Combine($stage, 'gosched.exe')
$daemon = [System.IO.Path]::Combine($stage, 'goschedd.exe')
$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = [Security.Principal.WindowsPrincipal]::new($identity)
$createdTaskIds = [System.Collections.Generic.List[string]]::new()
$groupCreated = $false
$serviceInstalled = $false
$primaryFailure = $null
$marker = $null
$probeResult = [ordered]@{ status = 'running'; primary_error = $null; diagnostics = @{}; cleanup = @() }

function Invoke-ProbeCli {
    [CmdletBinding()]
    Param(
        [Parameter(Mandatory=$true)]
        [string[]]$Arguments,

        [Parameter(Mandatory=$false)]
        [Switch]$ExpectJson
    )
    $start = [Diagnostics.ProcessStartInfo]::new()
    $start.FileName = $script:cli
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardInput = $true
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    foreach ($argument in $Arguments) { $start.ArgumentList.Add($argument) }
    $process = [Diagnostics.Process]::Start($start)
    try {
        $process.StandardInput.Close()
        $stdout = $process.StandardOutput.ReadToEndAsync()
        $stderr = $process.StandardError.ReadToEndAsync()
        if (-not $process.WaitForExit(120000)) {
            $process.Kill($true)
            $process.WaitForExit()
            throw "gosched $($Arguments -join ' ') exceeded its 120-second command deadline"
        }
        $exitCode = $process.ExitCode
        $text = $stdout.GetAwaiter().GetResult().Trim()
        $diagnostic = $stderr.GetAwaiter().GetResult().Trim()
        if ($exitCode -ne 0) { $text = "$text $diagnostic".Trim() }
    } finally { $process.Dispose() }
    if ($exitCode -ne 0) {
        throw "gosched $($Arguments -join ' ') failed with exit $exitCode`: $text"
    }
    if ($ExpectJson) {
        return $text | ConvertFrom-Json -Depth 20
    }
    $text
}

function Wait-ProbeRun {
    [CmdletBinding()]
    Param(
        [Parameter(Mandatory=$true)]
        [string]$TaskId,

        [Parameter(Mandatory=$true)]
        [ValidateSet('manual','schedule')]
        [string]$Trigger,

        [Parameter(Mandatory=$false)]
        [int]$TimeoutSeconds = 180
    )
    $deadline = [DateTime]::UtcNow.AddSeconds($TimeoutSeconds)
    do {
        $runs = Invoke-ProbeCli -ExpectJson -Arguments @(
            '--json', 'runs', '--task', $TaskId, '--limit', '20'
        )
        $run = @($runs) | Where-Object { $_.trigger -eq $Trigger -and ($_.ended_at -or $_.outcome -eq 'skipped') } |
            Select-Object -First 1
        if ($run) { return $run }
        Start-Sleep -Milliseconds 1000
    } while ([DateTime]::UtcNow -lt $deadline)
    throw "Timed out waiting for $Trigger run of task $TaskId"
}


function Write-ProbeEvidence {
    [CmdletBinding()]
    Param()
    $directory = [IO.Path]::GetDirectoryName($script:evidence)
    [void][IO.Directory]::CreateDirectory($directory)
    [IO.File]::WriteAllText($script:evidence, ($script:probeResult | ConvertTo-Json -Depth 20), [Text.UTF8Encoding]::new($false))
}

function Invoke-ProbeCleanup {
    [CmdletBinding()]
    Param([string]$Name, [scriptblock]$Operation)
    $entry = [ordered]@{ operation = $Name; error = $null }
    try { & $Operation | Out-Null } catch {
        $entry.error = $_.Exception.Message
        $script:probeResult.status = 'failed'
    }
    $script:probeResult.cleanup += $entry
}

function Save-ProbeFailure {
    [CmdletBinding()]
    Param()
    $diagnostics = [ordered]@{ tasks = @(); health = $null; service = $null; errors = @() }
    foreach ($taskId in $script:createdTaskIds) {
        $task = [ordered]@{ id = $taskId; runs = @(); active = @() }
        try { $task.runs = Invoke-ProbeCli -ExpectJson -Arguments @('--json', 'runs', '--task', $taskId, '--limit', '20') }
        catch { $diagnostics.errors += "runs $taskId : $($_.Exception.Message)" }
        try { $task.active = Invoke-ProbeCli -ExpectJson -Arguments @('--json', 'runs', 'active') | Where-Object task_id -eq $taskId }
        catch { $diagnostics.errors += "active $taskId : $($_.Exception.Message)" }
        $diagnostics.tasks += $task
    }
    try { $diagnostics.health = Invoke-ProbeCli -Arguments @('health') }
    catch { $diagnostics.errors += "health: $($_.Exception.Message)" }
    try { $diagnostics.service = Get-CimInstance -ClassName Win32_Service -Filter "Name='$script:serviceName'" | Select-Object Name, State, StartName, ExitCode }
    catch { $diagnostics.errors += "service: $($_.Exception.Message)" }
    # The disposable service owns this log. Bound the retained tail and never
    # enumerate the runner environment or unrelated process/credential state.
    $logPath = [IO.Path]::Combine($env:ProgramData, 'goschedule', 'logs', 'goschedule.log')
    try {
        if (Test-Path -LiteralPath $logPath -PathType Leaf) {
            $diagnostics.log_tail = @([IO.File]::ReadLines($logPath) | Select-Object -Last 100 | ForEach-Object { if ($_.Length -gt 4096) { $_.Substring(0,4096) } else { $_ } })
        }
    } catch { $diagnostics.errors += "log: $($_.Exception.Message)" }
    $script:probeResult.diagnostics = $diagnostics
    Write-ProbeEvidence
}

if (-not $IsWindows) { throw 'Service-core CI probe requires Windows.' }
if (-not $principal.IsInRole(
    [Security.Principal.WindowsBuiltInRole]::Administrator
)) {
    throw 'Service-core CI probe requires an elevated disposable runner.'
}
foreach ($path in @($cli, $daemon)) {
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
        throw "Staged binary is missing: $path"
    }
}
if (Get-LocalGroup -Name $groupName -ErrorAction SilentlyContinue) {
    throw "Disposable runner is not clean: local group $groupName exists."
}
if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
    throw "Disposable runner is not clean: service $serviceName exists."
}
if (Test-Path -LiteralPath $evidence) {
    throw "Evidence path already exists: $evidence"
}

try {
    [void](New-LocalGroup -Name $groupName `
        -Description 'go-schedule authorized operators')
    $groupCreated = $true
    Add-LocalGroupMember -Group $groupName -Member $identity.Name
    $group = Get-LocalGroup -Name $groupName
    $tokenSids = @($identity.User.Value) + @(
        $identity.Groups | ForEach-Object { $_.Value }
    )
    if ($tokenSids -contains $group.SID.Value) {
        throw 'Probe token unexpectedly contains the newly created group SID.'
    }

    [void](Invoke-ProbeCli -Arguments @('service', 'install'))
    $serviceInstalled = $true
    [void](Invoke-ProbeCli -Arguments @('service', 'start'))
    $healthDeadline = [DateTime]::UtcNow.AddSeconds(180)
    $health = ''
    do {
        try {
            $health = Invoke-ProbeCli -Arguments @('health')
            break
        } catch {
            if ([DateTime]::UtcNow -ge $healthDeadline) { throw }
            Start-Sleep -Milliseconds 1000
        }
    } while ($true)

    $markerDirectory = [System.IO.Path]::GetDirectoryName($evidence)
    [void][System.IO.Directory]::CreateDirectory($markerDirectory)
    $serviceDataDirectory = [System.IO.Path]::Combine(
        $env:ProgramData,
        'goschedule'
    )
    $marker = [System.IO.Path]::Combine(
        $serviceDataDirectory,
        's038-service-marker.txt'
    )
    $markerEvidence = [System.IO.Path]::Combine(
        $markerDirectory,
        's038-service-marker.txt'
    )
    if (Test-Path -LiteralPath $marker) {
        throw "Marker already exists: $marker"
    }
    $systemCommand = [System.IO.Path]::Combine(
        $env:SystemRoot,
        'System32',
        'WindowsPowerShell',
        'v1.0',
        'powershell.exe'
    )
    $suffix = Get-Date -Format 'yyyyMMddHHmmssfff'
    $escapedMarker = $marker.Replace("'", "''")
    $successCommand = "[IO.File]::AppendAllText('$escapedMarker', 'S038-marker' + [Environment]::NewLine); [Console]::Out.Write('S038-output')"
    $success = Invoke-ProbeCli -ExpectJson -Arguments @(
        '--json', 'task', 'add', "s038-ci-success-$suffix",
        '--command', $systemCommand,
        '--schedule', '*/5 * * * * *', '--tz', 'UTC',
        '--arg', '-NoLogo', '--arg', '-NoProfile',
        '--arg', '-NonInteractive', '--arg', '-Command',
        '--arg', $successCommand
    )
    $successTask = $success.task
    $createdTaskIds.Add([string]$successTask.id)
    [void](Invoke-ProbeCli -Arguments @(
        'task', 'run-now', [string]$successTask.id
    ))
    $manualRun = Wait-ProbeRun -TaskId $successTask.id -Trigger manual
    $scheduledRun = Wait-ProbeRun -TaskId $successTask.id -Trigger schedule
    if ($manualRun.outcome -ne 'success' -or $manualRun.exit_code -ne 0 -or
        $manualRun.output -notmatch 'S038-output') {
        throw 'Manual LocalSystem execution contract failed.'
    }
    if ($scheduledRun.outcome -ne 'success' -or
        $scheduledRun.exit_code -ne 0 -or
        $scheduledRun.output -notmatch 'S038-output') {
        throw 'Scheduled LocalSystem execution contract failed.'
    }
    $markerLines = @(Get-Content -LiteralPath $marker)
    if (@($markerLines | Where-Object { $_ -eq 'S038-marker' }).Count -lt 2) {
        throw 'Manual and scheduled marker effects were not both observed.'
    }
    Copy-Item -LiteralPath $marker -Destination $markerEvidence

    $exitFailure = Invoke-ProbeCli -ExpectJson -Arguments @(
        '--json', 'task', 'add', "s038-ci-exit-$suffix",
        '--command', $systemCommand,
        '--at', '2099-01-01T00:00:00Z', '--tz', 'UTC',
        '--arg', '-NoLogo', '--arg', '-NoProfile',
        '--arg', '-NonInteractive', '--arg', '-Command',
        '--arg', "[Console]::Out.Write('S038-controlled-failure'); exit 7"
    )
    $createdTaskIds.Add([string]$exitFailure.task.id)
    [void](Invoke-ProbeCli -Arguments @(
        'task', 'run-now', [string]$exitFailure.task.id
    ))
    $exitRun = Wait-ProbeRun -TaskId $exitFailure.task.id -Trigger manual
    if ($exitRun.outcome -ne 'failure' -or $exitRun.exit_code -ne 7 -or
        $exitRun.output -notmatch 'S038-controlled-failure') {
        throw 'Nonzero child-exit contract failed.'
    }

    $missingCommand = [System.IO.Path]::Combine(
        $markerDirectory,
        'missing-s038-executable.exe'
    )
    $startFailure = Invoke-ProbeCli -ExpectJson -Arguments @(
        '--json', 'task', 'add', "s038-ci-start-$suffix",
        '--command', $missingCommand,
        '--at', '2099-01-02T00:00:00Z', '--tz', 'UTC'
    )
    $createdTaskIds.Add([string]$startFailure.task.id)
    [void](Invoke-ProbeCli -Arguments @(
        'task', 'run-now', [string]$startFailure.task.id
    ))
    $startRun = Wait-ProbeRun -TaskId $startFailure.task.id -Trigger manual
    if ($startRun.outcome -ne 'failure' -or
        $null -ne $startRun.exit_code -or
        $startRun.output -notmatch '^process start failed for') {
        throw 'Process-start failure contract failed.'
    }

    $service = Get-CimInstance -ClassName Win32_Service `
        -Filter "Name='$serviceName'"
    if ($service.StartName -ne 'LocalSystem' -or $service.State -ne 'Running') {
        throw 'Service did not execute the probe as a running LocalSystem service.'
    }
    $result = [ordered]@{
        status = 'proven'
        identity = $identity.Name
        user_sid = $identity.User.Value
        group_sid = $group.SID.Value
        token_contains_group_sid = $false
        service_account = $service.StartName
        service_state = $service.State
        health = $health
        command = $systemCommand
        environment_keys = @()
        manual_run = $manualRun
        scheduled_run = $scheduledRun
        exit_failure_run = $exitRun
        start_failure_run = $startRun
        marker_path = $marker
        marker_evidence_path = $markerEvidence
        marker_sha256 = (Get-FileHash -LiteralPath $marker `
            -Algorithm SHA256).Hash.ToLowerInvariant()
        marker_lines = $markerLines.Count
    }
    $probeResult = $result
    $probeResult.primary_error = $null
    $probeResult.diagnostics = @{}
    $probeResult.cleanup = @()
    Write-ProbeEvidence
    Write-Output "Service-core probe proven: $evidence"
} catch {
    $primaryFailure = $_
    $probeResult.status = 'failed'
    $probeResult.primary_error = $_.Exception.Message
    try { Save-ProbeFailure } catch { Write-Warning "Failure evidence checkpoint failed: $($_.Exception.Message)" }
} finally {
    foreach ($taskId in $createdTaskIds) {
        Invoke-ProbeCleanup -Name "remove task $taskId" -Operation { [void](Invoke-ProbeCli -Arguments @('task', 'rm', $taskId)) }
    }
    if ($serviceInstalled) {
        Invoke-ProbeCleanup -Name 'stop service' -Operation { [void](Invoke-ProbeCli -Arguments @('service', 'stop')) }
        Invoke-ProbeCleanup -Name 'uninstall service' -Operation { [void](Invoke-ProbeCli -Arguments @('service', 'uninstall')) }
    }
    if ($groupCreated) {
        Invoke-ProbeCleanup -Name 'remove group' -Operation {
            if (Get-LocalGroup -Name $groupName -ErrorAction SilentlyContinue) { Remove-LocalGroup -Name $groupName }
        }
    }
    if ($marker) {
        Invoke-ProbeCleanup -Name 'remove marker' -Operation {
            if (Test-Path -LiteralPath $marker -PathType Leaf) { Remove-Item -LiteralPath $marker -Force }
        }
    }
    try { Write-ProbeEvidence } catch {
        if (-not $primaryFailure) { $primaryFailure = $_ }
        else { Write-Warning "Final evidence write failed: $($_.Exception.Message)" }
    }
}
if ($primaryFailure) { throw $primaryFailure }
if (@($probeResult.cleanup | Where-Object error).Count) { throw 'Service-core probe cleanup failed; see retained evidence.' }
