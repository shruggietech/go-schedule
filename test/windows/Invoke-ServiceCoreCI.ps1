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

function Invoke-ProbeCli {
    [CmdletBinding()]
    Param(
        [Parameter(Mandatory=$true)]
        [string[]]$Arguments,

        [Parameter(Mandatory=$false)]
        [Switch]$ExpectJson
    )
    $output = & $script:cli @Arguments 2>&1
    $exitCode = $LASTEXITCODE
    $text = (($output | Out-String).Trim())
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
        [int]$TimeoutSeconds = 45
    )
    $deadline = [DateTime]::UtcNow.AddSeconds($TimeoutSeconds)
    do {
        $runs = Invoke-ProbeCli -ExpectJson -Arguments @(
            '--json', 'runs', '--task', $TaskId, '--limit', '20'
        )
        $run = @($runs) | Where-Object trigger -eq $Trigger |
            Select-Object -First 1
        if ($run -and $run.ended_at) { return $run }
        Start-Sleep -Milliseconds 1000
    } while ([DateTime]::UtcNow -lt $deadline)
    throw "Timed out waiting for $Trigger run of task $TaskId"
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
    $healthDeadline = [DateTime]::UtcNow.AddSeconds(30)
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
    $result | ConvertTo-Json -Depth 20 |
        Set-Content -LiteralPath $evidence -Encoding utf8
    Write-Output "Service-core probe proven: $evidence"
} finally {
    foreach ($taskId in $createdTaskIds) {
        & $cli task rm $taskId 2>&1 | Out-Null
    }
    if ($serviceInstalled) {
        & $cli service stop 2>&1 | Out-Null
        & $cli service uninstall 2>&1 | Out-Null
    }
    if ($groupCreated -and
        (Get-LocalGroup -Name $groupName -ErrorAction SilentlyContinue)) {
        Remove-LocalGroup -Name $groupName
    }
    if ($marker -and (Test-Path -LiteralPath $marker -PathType Leaf)) {
        Remove-Item -LiteralPath $marker -Force
    }
}
