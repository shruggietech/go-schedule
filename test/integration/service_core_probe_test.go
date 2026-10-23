package integration

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestServiceProbeSelectsTerminalRunsAndRetainsFailureEvidence(t *testing.T) {
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		t.Skip("PowerShell 7 is required for the executable service-probe fixture")
	}
	script := filepath.Join("..", "windows", "Invoke-ServiceCoreCI.ps1")
	// Load only function definitions. Never execute the service-installing body.
	fixture := `
param($Source, $Evidence)
$ErrorActionPreference = 'Stop'
$tokens=$null; $errors=$null
$ast=[System.Management.Automation.Language.Parser]::ParseFile($Source,[ref]$tokens,[ref]$errors)
if ($errors.Count) { throw 'Probe has parse errors' }
$ast.FindAll({param($n) $n -is [System.Management.Automation.Language.FunctionDefinitionAst]},$false) | ForEach-Object { Invoke-Expression $_.Extent.Text }
function Invoke-ProbeCli {
    @(
        [pscustomobject]@{id='queued';trigger='manual';ended_at=$null;outcome='queued'},
        [pscustomobject]@{id='other';trigger='schedule';ended_at='2026-10-22T00:00:00Z';outcome='success'},
        [pscustomobject]@{id='completed';trigger='manual';ended_at='2026-10-22T00:00:00Z';outcome='failure'}
    )
}
function Start-Sleep {}
$run=Wait-ProbeRun -TaskId fixture -Trigger manual -TimeoutSeconds 0
if ($run.id -ne 'completed' -or $run.outcome -ne 'failure') { throw 'Terminal manual failure was hidden' }
function Invoke-ProbeCli { @([pscustomobject]@{id='skipped';trigger='manual';ended_at=$null;outcome='skipped'}) }
$skipped=Wait-ProbeRun -TaskId fixture -Trigger manual -TimeoutSeconds 0
if ($skipped.id -ne 'skipped') { throw 'Skipped terminal result hidden' }
function Invoke-ProbeCli { @([pscustomobject]@{id='queued';trigger='manual';ended_at=$null;outcome='queued'}) }
try { Wait-ProbeRun -TaskId fixture -Trigger manual -TimeoutSeconds 0; throw 'Unexpected terminal run' }
catch { if ($_ -notmatch 'Timed out waiting for manual run') { throw } }
$script:evidence=$Evidence
$script:probeResult=[ordered]@{status='failed';primary_error='original failure';diagnostics=@{runs=@('fixture-only')};cleanup=@()}
Write-ProbeEvidence
$script:createdTaskIds=@('fixture')
$script:serviceName='fixture-service'
function Invoke-ProbeCli { throw 'diagnostic unavailable' }
function Get-CimInstance { [pscustomobject]@{Name='fixture-service';State='Stopped';StartName='LocalSystem';ExitCode=0} }
$env:ProgramData=Join-Path ([IO.Path]::GetDirectoryName($Evidence)) 'fixture-data'
$logDirectory=Join-Path $env:ProgramData 'goschedule/logs'
[void][IO.Directory]::CreateDirectory($logDirectory)
[IO.File]::WriteAllText((Join-Path $logDirectory 'goschedule.log'), (('x'*5000+[Environment]::NewLine)*110))
Save-ProbeFailure
Invoke-ProbeCleanup -Name controlled -Operation { throw 'cleanup failure' }
Write-ProbeEvidence
$saved=Get-Content -Raw -LiteralPath $Evidence | ConvertFrom-Json
if ($saved.primary_error -ne 'original failure' -or $saved.cleanup.Count -ne 1 -or $saved.cleanup[0].error -notmatch 'cleanup failure' -or $saved.diagnostics.tasks[0].id -ne 'fixture' -or $saved.diagnostics.errors.Count -ne 3 -or $saved.diagnostics.log_tail.Count -ne 100 -or $saved.diagnostics.log_tail[0].Length -ne 4096) { throw 'Failure or cleanup evidence lost' }
$script:probeResult.status='proven'
$script:probeResult.primary_error=$null
$script:probeResult.cleanup=@()
Invoke-ProbeCleanup -Name cleanup-only -Operation { throw 'cleanup-only failure' }
Write-ProbeEvidence
$saved=Get-Content -Raw -LiteralPath $Evidence | ConvertFrom-Json
if ($saved.status -ne 'failed' -or $saved.primary_error -ne $null -or $saved.cleanup[0].error -notmatch 'cleanup-only failure') { throw 'Cleanup failure reported as proven' }
$bytes=[IO.File]::ReadAllBytes($Evidence)
if ($bytes.Length -ge 3 -and $bytes[0] -eq 239 -and $bytes[1] -eq 187 -and $bytes[2] -eq 191) { throw 'Evidence has a BOM' }
`
	command := scriptCommand(pwsh, "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", "& {"+fixture+"}", script, filepath.Join(t.TempDir(), "failure.json"))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("isolated service-probe fixture: %v\n%s", err, output)
	}
}
