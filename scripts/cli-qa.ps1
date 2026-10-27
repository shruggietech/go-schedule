[CmdletBinding()]
param(
    [ValidateSet('capture', 'resume', 'finalize', 'validate', 'render', 'inventory', 'build')]
    [string]$Action = 'capture',
    [string]$SourceDirectory = (Split-Path -Parent $PSScriptRoot),
    [string]$OutputDirectory,
    [string]$BinaryDirectory,
    [string]$WSLDistribution,
    [string]$WSLUser,
    [string]$Reviews,
    [string]$Budget,
    [switch]$Incidents,
    [string]$GoExecutable
)

# Maintainer entry point. Every console child uses redirected, hidden I/O.
$ErrorActionPreference = 'Stop'
if ($PSVersionTable.PSVersion.Major -lt 7) { throw 'Use the existing PowerShell 7 installation.' }
$qaSource = [IO.Path]::GetFullPath($SourceDirectory)
$qaTemporary = Join-Path $qaSource '.test-temp'
[void][IO.Directory]::CreateDirectory($qaTemporary)

function Invoke-QaHidden {
    param([string]$Executable, [string[]]$Arguments, [int]$TimeoutMilliseconds = 8000000, [switch]$CaptureOutput, [switch]$Evaluation)
    $qaStart = [Diagnostics.ProcessStartInfo]::new()
    $qaStart.FileName = $Executable
    $qaStart.WorkingDirectory = $qaSource
    $qaStart.UseShellExecute = $false
    $qaStart.CreateNoWindow = $true
    $qaStart.RedirectStandardInput = $true
    $qaStart.RedirectStandardOutput = $true
    $qaStart.RedirectStandardError = $true
    $qaStart.StandardOutputEncoding = [Text.UTF8Encoding]::new($false)
    $qaStart.StandardErrorEncoding = [Text.UTF8Encoding]::new($false)
    foreach ($qaArgument in $Arguments) { $qaStart.ArgumentList.Add($qaArgument) }
    $qaStart.Environment['GOMAXPROCS'] = '2'
    $qaStart.Environment['GOFLAGS'] = '-p=2'
    $qaStart.Environment['CGO_ENABLED'] = '0'
    $qaStart.Environment['GOOS'] = 'windows'
    $qaStart.Environment['GOARCH'] = 'amd64'
    $qaStart.Environment['GOTOOLCHAIN'] = 'local'
    $qaStart.Environment['GOPROXY'] = 'off'
    $qaProcess = [Diagnostics.Process]::Start($qaStart)
    if (-not $qaProcess.HasExited) {
        try { $qaProcess.PriorityClass = [Diagnostics.ProcessPriorityClass]::BelowNormal }
        catch {
            if (-not $qaProcess.HasExited) { $qaProcess.Kill($true); $qaProcess.Dispose(); throw }
        }
    }
    $qaProcess.StandardInput.Close()
    $qaOutput = $qaProcess.StandardOutput.ReadToEndAsync()
    $qaDiagnostic = $qaProcess.StandardError.ReadToEndAsync()
    $qaTimedOut = -not $qaProcess.WaitForExit($TimeoutMilliseconds)
    if ($qaTimedOut) {
        $qaProcess.Kill($true)
        $qaProcess.WaitForExit()
    }
    $qaCapturedOutput = $qaOutput.GetAwaiter().GetResult()
    if (-not $CaptureOutput) { [Console]::Out.Write($qaCapturedOutput) }
    [Console]::Error.Write($qaDiagnostic.GetAwaiter().GetResult())
    $qaExit = $qaProcess.ExitCode
    $qaProcess.Dispose()
    if ($qaTimedOut) { throw "The owned QA command exceeded its execution, cleanup and reporting allowance ($TimeoutMilliseconds ms). Its process tree was stopped. Retain OutputDirectory and the original BinaryDirectory, budget, WSL distribution and user; run resume with those same inputs. The last durable case checkpoints and worker ownership records are retained for recovery." }
    if ($Evaluation -and $qaExit -in @(0, 2, 3)) { return $qaExit }
    if ($qaExit -ne 0) { throw "The owned QA command returned exit $qaExit; read the diagnostic above." }
    if ($CaptureOutput) { return $qaCapturedOutput }
}

$qaMustHaveInputs = $Action -in @('resume', 'finalize', 'validate', 'render')
if ($qaMustHaveInputs -and -not $BinaryDirectory) { throw 'Supply the original BinaryDirectory for this existing run.' }
if ($BinaryDirectory) { $BinaryDirectory = [IO.Path]::GetFullPath($BinaryDirectory) }
$qaExecutable = if ($BinaryDirectory) { Join-Path $BinaryDirectory 'windows/gosched-qa.exe' } else { '' }
if (-not $qaExecutable -or -not [IO.File]::Exists($qaExecutable)) {
    if ($qaMustHaveInputs) { throw 'The original coordinator executable is unavailable; do not substitute new binaries for a resume.' }
    $qaGo = $GoExecutable
    if (-not $qaGo) {
        $qaModule = [IO.File]::ReadAllText((Join-Path $qaSource 'go.mod'))
        if ($qaModule -notmatch '(?m)^go\s+(\d+\.\d+(?:\.\d+)?)\s*$') { throw 'Cannot identify the required Go version from go.mod.' }
        $qaRequiredGo = [version]$Matches[1]
        $qaCachedCandidates = @()
        $qaCacheRoots = @($env:GOMODCACHE, (Join-Path $env:USERPROFILE 'go/pkg/mod'))
        if ($env:GOPATH) { foreach ($qaGoPath in $env:GOPATH.Split(';')) { $qaCacheRoots += (Join-Path $qaGoPath 'pkg/mod') } }
        foreach ($qaCache in $qaCacheRoots) {
            if ($qaCache) { $qaCachedCandidates += (Join-Path $qaCache ('golang.org/toolchain@v0.0.1-go' + $qaRequiredGo.ToString() + '.windows-amd64/bin/go.exe')) }
        }
        $qaGoCandidates = $qaCachedCandidates + @((Join-Path $env:ProgramFiles 'Go/bin/go.exe')) + @((Get-Command go -CommandType Application -All -ErrorAction SilentlyContinue).Source)
        foreach ($qaCandidate in ($qaGoCandidates | Select-Object -Unique)) {
            if (-not $qaCandidate -or -not [IO.File]::Exists($qaCandidate)) { continue }
            $qaVersionText = Invoke-QaHidden -Executable $qaCandidate -Arguments @('version') -TimeoutMilliseconds 30000 -CaptureOutput
            if ($qaVersionText -match 'go version go(\d+\.\d+(?:\.\d+)?)' -and [version]$Matches[1] -ge $qaRequiredGo) { $qaGo = $qaCandidate; break }
        }
        if (-not $qaGo) { throw 'No suitable existing Go toolchain was found. Select its path with GoExecutable.' }
    }
    $qaExecutable = Join-Path $qaTemporary ('gosched-qa-bootstrap-' + [guid]::NewGuid().ToString() + '.exe')
    Invoke-QaHidden -Executable $qaGo -Arguments @('build', '-p=2', '-o', $qaExecutable, './cmd/gosched-qa') -TimeoutMilliseconds 900000
    if ($Action -eq 'inventory') { Invoke-QaHidden -Executable $qaExecutable -Arguments @('inventory'); return }
    if (-not $BinaryDirectory) { $BinaryDirectory = Join-Path $qaTemporary ('cli-qa-inputs-' + [guid]::NewGuid().ToString()) }
    $qaBuildArguments = @('build', '--source', $qaSource, '--bin-dir', $BinaryDirectory, '--go', $qaGo)
    if ($Budget) { $qaBuildArguments += @('--budget', [IO.Path]::GetFullPath($Budget)) }
    Invoke-QaHidden -Executable $qaExecutable -Arguments $qaBuildArguments
    $qaExecutable = Join-Path $BinaryDirectory 'windows/gosched-qa.exe'
}
if ($Action -eq 'build') { Write-Output "Evaluation inputs: $BinaryDirectory"; return }
if ($Action -eq 'inventory') { Invoke-QaHidden -Executable $qaExecutable -Arguments @('inventory'); return }
if (-not $OutputDirectory) {
    if ($Action -ne 'capture') { throw 'Supply OutputDirectory for an existing run.' }
    $OutputDirectory = Join-Path $qaTemporary ('cli-qa-run-' + [guid]::NewGuid().ToString())
}
$OutputDirectory = [IO.Path]::GetFullPath($OutputDirectory)
$qaArguments = @($Action, '--out', $OutputDirectory)
if ($Action -in @('capture', 'resume')) {
    if (-not $WSLDistribution -or -not $WSLUser) { throw 'Select an existing WSL2 distribution and unelevated WSL user.' }
    $qaArguments += @('--source', $qaSource, '--bin-dir', $BinaryDirectory, '--wsl-distro', $WSLDistribution, '--wsl-user', $WSLUser)
    if ($Budget) { $qaArguments += @('--budget', [IO.Path]::GetFullPath($Budget)) }
    if ($Incidents) { $qaArguments += '--incidents' }
}
if ($Action -eq 'finalize') {
    if (-not $Reviews) { throw 'Supply the completed agent Reviews JSON against this frozen corpus.' }
    $qaArguments += @('--reviews', [IO.Path]::GetFullPath($Reviews))
}
$qaTimeout = 8000000
if ($Action -in @('capture', 'resume')) {
    $qaPolicyArguments = @('timing-policy')
    if ($Budget) { $qaPolicyArguments += @('--budget', [IO.Path]::GetFullPath($Budget)) }
    if ($Action -eq 'resume') { $qaPolicyArguments += @('--resume-out', $OutputDirectory) }
    $qaPolicy = (Invoke-QaHidden -Executable $qaExecutable -Arguments $qaPolicyArguments -TimeoutMilliseconds 60000 -CaptureOutput) | ConvertFrom-Json
    $qaTimeout = [int]$qaPolicy.outer_timeout_ms
}
$qaEvaluationExit = Invoke-QaHidden -Executable $qaExecutable -Arguments $qaArguments -TimeoutMilliseconds $qaTimeout -Evaluation
Write-Output "Evaluation inputs: $BinaryDirectory"
Write-Output "Evaluation artifacts: $OutputDirectory"
exit $qaEvaluationExit
