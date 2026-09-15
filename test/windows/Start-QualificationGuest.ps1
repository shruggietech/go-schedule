<#
.SYNOPSIS
    Stage and diagnose an attended candidate installation inside Windows Sandbox.

.DESCRIPTION
    Verifies packaged inputs, copies them to guest-local storage, prepares offline
    prerequisites, and opens intentional attended installer UI. Console children
    run hidden with redirected I/O. Networking is not required. Logs and phase
    results export as they are produced. No attended observation is marked passed.

    Refuses outside the Sandbox account and virtual-machine environment. A timeout
    leaves Installer state intact and stops the session; reset before retrying.
    Windows PowerShell 5.1 is supported because it ships in a clean Sandbox.

.PARAMETER Scenario
    Fresh or upgrade installation session. Alias: s

.PARAMETER Fixture
    Run nondestructive process fixtures only, without installing or staging.
    Alias: f

.PARAMETER FixtureOutput
    Existing empty directory for fixture diagnostic output. Alias: o

.PARAMETER Help
    Print detailed help. Alias: h

.EXAMPLE
    .\Start-QualificationGuest.ps1 -Scenario fresh
    Invoked automatically by the prepared fresh.wsb configuration.

.EXAMPLE
    .\Start-QualificationGuest.ps1 -Fixture -FixtureOutput C:\fixture-output
    Exercises only hidden child success, failure, and deadline behavior.
#>
[CmdletBinding(SupportsShouldProcess=$true,ConfirmImpact='Medium',DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('s')]
    [ValidateSet('fresh','upgrade')]
    [string]$Scenario = 'fresh',
    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('f')]
    [switch]$Fixture,
    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('o')]
    [string]$FixtureOutput,
    [Parameter(Mandatory=$true,ParameterSetName='HelpText')]
    [Alias('h')]
    [switch]$Help
)
#_______________________________________________________________________________
## Declare Functions

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
        $color = switch ($Level) {
            'Info'    { 'Gray' }
            'Success' { 'Green' }
            'Warn'    { 'Yellow' }
            'Error'   { 'Red' }
            'Debug'   { 'DarkGray' }
        }
        Write-Host ("{0} {1}{2} {3}" -f $stamp, $tag, $label, $Message) -ForegroundColor $color
    }

    function Save-Record {
        [CmdletBinding()]
        Param([string]$Name, [object]$Value)
        $path = Join-Path $script:Exports $Name
        $stream = [IO.File]::Open($path, 'CreateNew', 'Write', 'None')
        $writer = [IO.StreamWriter]::new($stream, [Text.UTF8Encoding]::new($false))
        try { $writer.Write(($Value | ConvertTo-Json -Depth 12) + "`n") }
        finally { $writer.Dispose(); $stream.Dispose() }
    }

    function Invoke-DiagnosticProcess {
        [CmdletBinding()]
        Param([string]$Name, [string]$FilePath, [string]$Arguments,
            [int]$TimeoutSeconds = 1200, [switch]$Installer)
        $info = [Diagnostics.ProcessStartInfo]::new()
        $info.FileName = $FilePath
        $info.Arguments = $Arguments
        $info.UseShellExecute = $false
        $info.CreateNoWindow = $true
        $info.WindowStyle = 'Hidden'
        $info.RedirectStandardOutput = $true
        $info.RedirectStandardError = $true
        $info.RedirectStandardInput = $true
        $started = [datetime]::UtcNow
        $watch = [Diagnostics.Stopwatch]::StartNew()
        $process = [Diagnostics.Process]::Start($info)
        $process.StandardInput.Close()
        # Stream bytes directly to create-only export files, including before EOF.
        # A one-byte FileStream buffer avoids losing captured data at timeout.
        $stdoutPath = Join-Path $script:Exports "$Name-stdout.log"
        $stderrPath = Join-Path $script:Exports "$Name-stderr.log"
        $stdoutFile = [IO.FileStream]::new($stdoutPath, 'CreateNew', 'Write', 'Read', 1)
        $stderrFile = [IO.FileStream]::new($stderrPath, 'CreateNew', 'Write', 'Read', 1)
        $stdout = $process.StandardOutput.BaseStream.CopyToAsync($stdoutFile)
        $stderr = $process.StandardError.BaseStream.CopyToAsync($stderrFile)
        $nextProgress = 0
        $status = 'timed-out'
        $code = $null
        while ($watch.Elapsed.TotalSeconds -lt $TimeoutSeconds) {
            if ($watch.Elapsed.TotalSeconds -ge $nextProgress) {
                $line = '{0:o} {1}: elapsed {2:n0}s; logs in {3}' -f (
                    [datetime]::UtcNow, $Name, $watch.Elapsed.TotalSeconds,
                    $script:Exports)
                [IO.File]::AppendAllText((Join-Path $script:Exports 'progress.log'),
                    "$line`n", [Text.UTF8Encoding]::new($false))
                Write-Log $line -Source $Name
                $nextProgress += 10
            }
            if ($process.WaitForExit(1000)) {
                $code = $process.ExitCode
                $status = if ($code -eq 0) { 'completed' } else { 'failed' }
                if ($Installer -and $code -eq 3010) { $status = 'completed-reboot-required' }
                break
            }
        }
        # Never wait indefinitely for inherited pipe handles after child exit.
        if ($null -ne $code) {
            if (-not $stdout.Wait(3000) -or -not $stderr.Wait(3000)) {
                $status = 'failed'
            }
        }
        $captureComplete = $stdout.IsCompleted -and $stderr.IsCompleted -and
            -not $stdout.IsFaulted -and -not $stderr.IsFaulted
        if ($stdout.IsCompleted) { $stdoutFile.Dispose() }
        if ($stderr.IsCompleted) { $stderrFile.Dispose() }
        Save-Record "$Name-output.json" @{
            stdout_log = $stdoutPath; stderr_log = $stderrPath
            capture_complete = $captureComplete
        }
        Save-Record "$Name-phase.json" @{
            started_at = $started.ToString('o')
            completed_at = [datetime]::UtcNow.ToString('o')
            elapsed_seconds = $watch.Elapsed.TotalSeconds
            process_id = $process.Id; status = $status; exit_code = $code
            logs = $script:Exports
            reboot_required = ($Installer -and $code -eq 3010)
            attended_status = 'unavailable'
        }
        if ($status -ne 'completed') {
            throw "$Name $status. Do not start another installer; preserve logs and reset the guest."
        }
    }

    function Show-Instructions {
        [CmdletBinding()]
        Param([string]$Text)
        $answer = [Windows.Forms.MessageBox]::Show($Text,
            'go-schedule attended qualification', 'OKCancel', 'Information')
        if ($answer -ne 'OK') { throw 'Operator canceled qualification.' }
    }

    function Invoke-AttendedMsi {
        [CmdletBinding()]
        Param([string]$Name, [string]$Path)
        # An idle service-side msiexec is not an active install. Operations here
        # are serialized; native Installer rejects conflicting transactions (1618).
        $log = Join-Path $script:Exports "$Name-msi.log"
        Invoke-DiagnosticProcess -Name $Name -FilePath "$env:WINDIR\System32\msiexec.exe" `
            -Arguments "/i `"$Path`" /norestart /L*vx! `"$log`"" -Installer
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $ThisScriptPath = $MyInvocation.MyCommand.Path
    $ErrorActionPreference = 'Stop'
    $script:LogQuiet = $false
    $script:LogSilent = $false
    $script:Exports = 'C:\qualification-exports'

#_______________________________________________________________________________
## Execute Operations

    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }
    if ($Fixture) {
        if (-not [IO.Path]::IsPathRooted($FixtureOutput) -or
            -not (Test-Path -LiteralPath $FixtureOutput -PathType Container) -or
            @(Get-ChildItem -LiteralPath $FixtureOutput -Force).Count -ne 0) {
            throw 'FixtureOutput must be an existing empty absolute directory.'
        }
        if (-not $PSCmdlet.ShouldProcess($FixtureOutput, 'Write nondestructive process fixtures')) { exit 0 }
        $script:Exports = $FixtureOutput
        $engine = (Get-Process -Id $PID).Path
        Invoke-DiagnosticProcess 'success' $engine '-NoProfile -NonInteractive -Command "exit 0"' 10
        foreach ($case in @('failure','timeout')) {
            $fixtureArguments = if ($case -eq 'failure') {
                '-NoProfile -NonInteractive -Command "exit 7"'
            } else { '-NoProfile -NonInteractive -Command "[Console]::Out.WriteLine(''timeout-stdout''); [Console]::Error.WriteLine(''timeout-stderr''); Start-Sleep -Seconds 6"' }
            $refused = $false
            $fixtureDeadline = if ($case -eq 'failure') { 10 } else { 3 }
            try { Invoke-DiagnosticProcess $case $engine $fixtureArguments $fixtureDeadline }
            catch { $refused = $true }
            if (-not $refused) { throw "Fixture $case did not refuse success." }
        }
        $refused = $false
        try { Invoke-DiagnosticProcess 'reboot' $engine '-NoProfile -NonInteractive -Command "exit 3010"' 10 -Installer }
        catch { $refused = $true }
        if (-not $refused) { throw 'Reboot-required fixture permitted subsequent work.' }
        exit 0
    }
    # Accidental host execution is refused before any staging or installation.
    if ($env:USERNAME -ne 'WDAGUtilityAccount' -or
        (Get-CimInstance Win32_ComputerSystem).Model -ne 'Virtual Machine' -or
        $PSScriptRoot -ne 'C:\qualification-inputs') {
        throw 'This installer bootstrap runs only in the prepared Windows Sandbox.'
    }
    if (-not (Test-Path -LiteralPath $script:Exports -PathType Container) -or
        @(Get-ChildItem -LiteralPath $script:Exports -Force).Count -ne 0 -or
        (Test-Path -LiteralPath 'C:\qualification')) {
        throw 'Occupied or unavailable session destination; prepare a new package.'
    }
    if (-not $PSCmdlet.ShouldProcess('Windows Sandbox guest', 'Stage prerequisites and open attended installer')) { exit 0 }
    try {
        $manifest = Get-Content -LiteralPath "$PSScriptRoot\manifest.json" -Raw | ConvertFrom-Json
        $names = @{
            candidate = 'go-schedule_v1.4.0_windows_amd64.msi'; baseline = 'go-schedule_v1.1.1_windows_amd64.msi'
            powershell = 'powershell.zip'; webview2 = 'webview2.exe'
            bootstrap = 'Start-QualificationGuest.ps1'
            collector = 'Invoke-ReleaseCandidateAttended.ps1'
        }
        if ($manifest.schema_version -ne 1 -or $manifest.repository -ne 'shruggietech/go-schedule' -or
            $manifest.tag -ne 'v1.4.0' -or $manifest.commit -cnotmatch '^[a-f0-9]{40}$' -or
            $manifest.run_id -le 0 -or $manifest.run_attempt -le 0 -or
            @($manifest.inputs).Count -ne 6) { throw 'Invalid packaged provenance.' }
        [IO.Directory]::CreateDirectory('C:\qualification') | Out-Null
        $seen = @{}
        foreach ($inputFile in $manifest.inputs) {
            if (-not $names.ContainsKey($inputFile.role) -or $seen.ContainsKey($inputFile.role) -or
                $inputFile.path -cne $names[$inputFile.role]) { throw 'Unsafe packaged role/path.' }
            $seen[$inputFile.role] = $true
            $source = Join-Path $PSScriptRoot $inputFile.path
            $target = Join-Path 'C:\qualification' $inputFile.path
            Copy-Item -LiteralPath $source -Destination $target
            if ((Get-Item -LiteralPath $target).Length -ne $inputFile.bytes -or
                (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash.ToLowerInvariant() -cne $inputFile.sha256) {
                throw "Packaged input changed: $($inputFile.role)"
            }
        }
        Save-Record 'session.json' @{ scenario = $Scenario; candidate = $manifest; attended_status = 'unavailable' }
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        [IO.Compression.ZipFile]::ExtractToDirectory('C:\qualification\powershell.zip', 'C:\qualification\pwsh')
        if (-not (Test-Path -LiteralPath 'C:\qualification\pwsh\pwsh.exe')) { throw 'Portable PowerShell executable missing.' }
        Invoke-DiagnosticProcess 'powershell-check' 'C:\qualification\pwsh\pwsh.exe' `
            '-NoProfile -NonInteractive -Command "if ($PSVersionTable.PSVersion.Major -lt 7) { exit 1 }"' 30
        Invoke-DiagnosticProcess 'collector-initialize' 'C:\qualification\pwsh\pwsh.exe' `
            -Arguments ('-NoProfile -NonInteractive -File C:\qualification\Invoke-ReleaseCandidateAttended.ps1 -Action Initialize -MsiPath C:\qualification\go-schedule_v1.4.0_windows_amd64.msi -WorkspacePath C:\qualification\attended-workspace -Tag v1.4.0 -Commit {0} -RunId {1} -RunAttempt {2}' -f $manifest.commit, $manifest.run_id, $manifest.run_attempt) `
            -TimeoutSeconds 120
        Invoke-DiagnosticProcess 'webview2' 'C:\qualification\webview2.exe' '/silent /install' 600
        $runtime = @(
            Get-ItemProperty -LiteralPath 'HKLM:\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}' -ErrorAction SilentlyContinue
            Get-ItemProperty -LiteralPath 'HKCU:\SOFTWARE\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}' -ErrorAction SilentlyContinue
        ) | Where-Object { $_.pv -and $_.pv -ne '0.0.0.0' }
        if (@($runtime).Count -eq 0) { throw 'WebView2 runtime remains unavailable.' }
        Add-Type -AssemblyName System.Windows.Forms
        if ($Scenario -eq 'upgrade') {
            Show-Instructions 'Install the public v1.1.1 baseline. Do not close this Sandbox. Logs export automatically.'
            Invoke-AttendedMsi 'baseline' 'C:\qualification\go-schedule_v1.1.1_windows_amd64.msi'
            Show-Instructions 'Before clicking OK, create representative tasks, history, and appearance settings in v1.1.1, then close the application. This requires real observation, not an automatic pass.'
        }
        Show-Instructions 'Install the corrected candidate. Inputs are guest-local and offline prerequisites are prepared. Installer timing/logs export automatically. Do not launch another installer if this times out.'
        Invoke-AttendedMsi 'candidate' 'C:\qualification\go-schedule_v1.4.0_windows_amd64.msi'
        Show-Instructions 'Mechanical installation finished, NOT release qualification. Before clicking OK, review the native matrix and #229-#233 regressions using C:\qualification\Invoke-ReleaseCandidateAttended.ps1 and C:\qualification\attended-workspace. OK exports the workspace without marking observations passed. Normal-user and high/mixed-DPI checks need separate environments. Finalize through the repository collector on the host with its Go context after merging genuine evidence.'
        Copy-Item -LiteralPath 'C:\qualification\attended-workspace' -Destination $script:Exports -Recurse
        exit 0
    } catch {
        Save-Record 'session-failure.json' @{ error = $_.Exception.Message; attended_status = 'unavailable' }
        Write-Log $_.Exception.Message -Level Error
        exit 1
    }

#_______________________________________________________________________________
## End of script
