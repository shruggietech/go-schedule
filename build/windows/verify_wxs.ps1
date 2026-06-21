<#
.SYNOPSIS
  Sanity-checks build/windows/goschedule.wxs before the MSI is built in CI.

.DESCRIPTION
  Cheap guard against the WiX source drifting from reality:
    * the three expected binaries are referenced as File sources,
    * the Windows service Name is exactly "goschedd" (the name the CLI
      `gosched service ...` control layer expects),
    * the install folder is "go-schedule" and the package is per-machine.
  If a StageDir is provided, also asserts the referenced .exe files exist there.

.PARAMETER StageDir
  Optional path to the staged build output containing the three .exe files.

.EXAMPLE
  pwsh build/windows/verify_wxs.ps1 -StageDir $stage
#>
param(
  [string]$StageDir
)

$ErrorActionPreference = 'Stop'
$wxsPath = Join-Path $PSScriptRoot 'goschedule.wxs'
if (-not (Test-Path $wxsPath)) { throw "wxs not found at $wxsPath" }
$wxs = Get-Content $wxsPath -Raw

$expectedBinaries = @('goschedd.exe', 'gosched-gui.exe', 'gosched.exe')
$fail = @()

foreach ($bin in $expectedBinaries) {
  if ($wxs -notmatch [regex]::Escape("Source=`"`$(StageDir)\$bin`"")) {
    $fail += "wxs does not reference binary: $bin"
  }
}

if ($wxs -notmatch 'ServiceInstall[^>]*Name="goschedd"') {
  $fail += 'ServiceInstall Name must be "goschedd"'
}
if ($wxs -notmatch 'ServiceControl[^>]*Name="goschedd"') {
  $fail += 'ServiceControl Name must be "goschedd"'
}
if ($wxs -notmatch 'Start="auto"') {
  $fail += 'service Start must be "auto" (start on boot)'
}
if ($wxs -notmatch 'Directory Id="INSTALLFOLDER" Name="go-schedule"') {
  $fail += 'install folder must be "go-schedule"'
}
if ($wxs -notmatch 'Scope="perMachine"') {
  $fail += 'package Scope must be "perMachine" (requires elevation)'
}

if ($StageDir) {
  foreach ($bin in $expectedBinaries) {
    $p = Join-Path $StageDir $bin
    if (-not (Test-Path $p)) { $fail += "staged binary missing: $p" }
  }
}

if ($fail.Count -gt 0) {
  Write-Error ("WiX sanity check failed:`n - " + ($fail -join "`n - "))
  exit 1
}

Write-Host 'WiX sanity check passed.'
