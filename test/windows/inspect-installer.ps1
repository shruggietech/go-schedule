<#
.SYNOPSIS
  Inspects a compiled go-schedule MSI without installing it.

.DESCRIPTION
  Reads Windows Installer tables through the built-in COM API and verifies the
  canonical icon and machine PATH relationships authored by the WiX source.
  This is candidate/published artifact evidence, not native shell observation.
#>
param(
  [Parameter(Mandatory)]
  [string]$MsiPath,

  [string]$EvidencePath
)

$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'MSI inspection requires Windows.' }

$resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
$installer = New-Object -ComObject WindowsInstaller.Installer
$database = $installer.OpenDatabase($resolvedMsi, 0)

function Get-MsiString {
  param(
    [Parameter(Mandatory)] $Database,
    [Parameter(Mandatory)] [string]$Query,
    [int]$Field = 1
  )

  $view = $Database.OpenView($Query)
  $view.Execute() | Out-Null
  $record = $view.Fetch()
  if (-not $record) { return '' }
  $record.StringData($Field)
}

$fail = [System.Collections.Generic.List[string]]::new()
$canonicalIcon = 'GoSchedule.ico'

$version = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='ProductVersion'"
if (-not $version) { $version = '<missing>' }

$arpIcon = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='ARPPRODUCTICON'"
if ($arpIcon -ne $canonicalIcon) {
  $fail.Add("Property.ARPPRODUCTICON is '$arpIcon'; expected '$canonicalIcon'")
}

try {
  $iconName = Get-MsiString -Database $database `
    -Query "SELECT ``Name`` FROM ``Icon`` WHERE ``Name``='$canonicalIcon'"
} catch {
  $iconName = ''
}
if ($iconName -ne $canonicalIcon) {
  $fail.Add("Icon table has no '$canonicalIcon' row")
}

$shortcutIcon = Get-MsiString -Database $database `
  -Query "SELECT ``Icon_`` FROM ``Shortcut`` WHERE ``Shortcut``='GuiShortcut'"
if ($shortcutIcon -ne $canonicalIcon) {
  $fail.Add("Shortcut.GuiShortcut.Icon_ is '$shortcutIcon'; expected '$canonicalIcon'")
}

$environmentQuery = "SELECT ``Name``, ``Value``, ``Component_`` FROM ``Environment`` WHERE ``Environment``='PathEnv'"
$environmentName = Get-MsiString -Database $database -Query $environmentQuery -Field 1
if (-not $environmentName) {
  $environmentName = ''
  $environmentValue = ''
  $environmentComponent = ''
  $fail.Add('Environment.PathEnv row is missing')
} else {
  $environmentValue = Get-MsiString -Database $database -Query $environmentQuery -Field 2
  $environmentComponent = Get-MsiString -Database $database -Query $environmentQuery -Field 3
  if ($environmentName -ne '=-*PATH') {
    $fail.Add("Environment.PathEnv.Name is '$environmentName'; expected '=-*PATH'")
  }
  if ($environmentValue -ne '[~];[INSTALLFOLDER]') {
    $fail.Add("Environment.PathEnv.Value is '$environmentValue'; expected '[~];[INSTALLFOLDER]'")
  }
  if ($environmentComponent -ne 'Gosched') {
    $fail.Add("Environment.PathEnv.Component_ is '$environmentComponent'; expected 'Gosched'")
  }
}

$hash = (Get-FileHash -LiteralPath $resolvedMsi -Algorithm SHA256).Hash.ToLowerInvariant()
$status = if ($fail.Count -eq 0) { 'proven' } else { 'failed' }
$evidence = @(
  '# Windows MSI Artifact Evidence'
  ''
  "- Date: $(Get-Date -Format 'yyyy-MM-dd')"
  "- Artifact: ``$resolvedMsi``"
  "- SHA-256: ``$hash``"
  "- Product version: ``$version``"
  "- Candidate/published artifact status: **$status**"
  "- Icon row: ``$iconName``"
  "- ARPPRODUCTICON: ``$arpIcon``"
  "- GuiShortcut Icon_: ``$shortcutIcon``"
  "- PATH row: ``$environmentName`` | ``$environmentValue`` | ``$environmentComponent``"
)
if ($fail.Count -gt 0) {
  $evidence += ''
  $evidence += '## Failures'
  $evidence += $fail | ForEach-Object { "- $_" }
}

if ($EvidencePath) {
  $evidenceFile = [System.IO.Path]::GetFullPath($EvidencePath)
  $evidence | Set-Content -LiteralPath $evidenceFile -Encoding utf8NoBOM
  Write-Output "installer-inspect: evidence written to $evidenceFile"
}

if ($fail.Count -gt 0) {
  [Console]::Error.WriteLine(
    "installer-inspect: FAILED`n - " + ($fail -join "`n - ")
  )
  exit 1
}

Write-Output "installer-inspect: OK - version $version, canonical icon and PATH rows proven"
