<#
.SYNOPSIS
  Inspects a compiled go-schedule MSI without installing it.

.DESCRIPTION
  Reads Windows Installer tables through the built-in COM API and verifies the
  Summary Information Subject, canonical identity, icon, machine PATH, and
  local administrative-group relationships authored by the WiX source. This is
  compiled candidate/published artifact evidence, not native Explorer or
  lifecycle observation.
#>
param(
  [Parameter(Mandatory)]
  [string]$MsiPath,

  [string]$EvidencePath,

  [Parameter(Mandatory)]
  [ValidateSet('candidate', 'published')]
  [string]$ArtifactClass,

  [Parameter(Mandatory)]
  [ValidateNotNullOrEmpty()]
  [string]$ArtifactOrigin
)

$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'MSI inspection requires Windows.' }

$ArtifactOrigin = $ArtifactOrigin.Trim()
if (-not $ArtifactOrigin) { throw 'Artifact origin must not be blank.' }
if ($ArtifactClass -eq 'published') {
  [Uri]$originUri = $null
  if (-not [Uri]::TryCreate($ArtifactOrigin, [UriKind]::Absolute, [ref]$originUri) -or
      $originUri.Scheme -ne 'https' -or
      $originUri.Host -ne 'github.com' -or
      $originUri.AbsolutePath -notmatch '^/shruggietech/go-schedule/releases/download/[^/]+/[^/]+\.msi$') {
    throw "Published artifact origin must be this repository's absolute HTTPS release-asset URL."
  }
}

$resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
$installer = New-Object -ComObject WindowsInstaller.Installer
$database = $installer.OpenDatabase($resolvedMsi, 0)
$summary = $database.SummaryInformation(0)

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
$canonicalSubject = 'go-schedule: cross-platform task scheduler'

$summaryTitle = $summary.Property(2)
$summarySubject = $summary.Property(3)
$summaryAuthor = $summary.Property(4)
$summaryComments = $summary.Property(6)
if ($summaryTitle -ne 'Installation Database') {
  $fail.Add("Summary Title is '$summaryTitle'; expected 'Installation Database'")
}
if ($summarySubject -ne $canonicalSubject) {
  $fail.Add("Summary Subject is '$summarySubject'; expected '$canonicalSubject'")
}
if ($summarySubject.Contains([char]0x2014)) {
  $fail.Add('Summary Subject contains U+2014')
}
if ($summaryAuthor -ne 'ShruggieTech') {
  $fail.Add("Summary Author is '$summaryAuthor'; expected 'ShruggieTech'")
}
$canonicalComments = 'This installer database contains the logic and data required to install go-schedule.'
if ($summaryComments -ne $canonicalComments) {
  $fail.Add("Summary Comments is '$summaryComments'; expected '$canonicalComments'")
}

# Windows Explorer reads the same Shell property system. Query System.Subject
# through that native handler as direct evidence of the value Explorer exposes.
$shell = New-Object -ComObject Shell.Application
$shellFolder = $shell.NameSpace([System.IO.Path]::GetDirectoryName($resolvedMsi))
$shellItem = $shellFolder.ParseName([System.IO.Path]::GetFileName($resolvedMsi))
$explorerSubject = $shellItem.ExtendedProperty('System.Subject')
if ($explorerSubject -ne $canonicalSubject) {
  $fail.Add("Explorer System.Subject is '$explorerSubject'; expected '$canonicalSubject'")
}

$version = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='ProductVersion'"
if (-not $version) { $version = '<missing>' }

$productName = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='ProductName'"
if ($productName -ne 'go-schedule') {
  $fail.Add("Property.ProductName is '$productName'; expected 'go-schedule'")
}

$manufacturer = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='Manufacturer'"
if ($manufacturer -ne 'ShruggieTech') {
  $fail.Add("Property.Manufacturer is '$manufacturer'; expected 'ShruggieTech'")
}

$upgradeCode = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='UpgradeCode'"
$canonicalUpgradeCode = '{B6F3C2E1-7A4D-4C9E-9B2A-1F6D8E5A0C34}'
if ($upgradeCode.ToUpperInvariant() -ne $canonicalUpgradeCode) {
  $fail.Add("Property.UpgradeCode is '$upgradeCode'; expected '$canonicalUpgradeCode'")
}

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

$groupQuery = "SELECT ``Name``, ``Domain`` FROM ``Wix4Group`` WHERE ``Group``='GoScheduleAdminGroup'"
$adminGroupName = Get-MsiString -Database $database -Query $groupQuery -Field 1
if (-not $adminGroupName) {
  $adminGroupName = ''
  $adminGroupDomain = ''
  $fail.Add('Wix4Group.GoScheduleAdminGroup row is missing')
} else {
  $adminGroupDomain = Get-MsiString -Database $database -Query $groupQuery -Field 2
  if ($adminGroupName -ne 'goschedadmin') {
    $fail.Add("Wix4Group.GoScheduleAdminGroup.Name is '$adminGroupName'; expected 'goschedadmin'")
  }
  if ($adminGroupDomain) {
    $fail.Add(
      "Wix4Group.GoScheduleAdminGroup.Domain is '$adminGroupDomain'; " +
      'expected empty for elevated local-group creation'
    )
  }
}

$hash = (Get-FileHash -LiteralPath $resolvedMsi -Algorithm SHA256).Hash.ToLowerInvariant()
$status = if ($fail.Count -eq 0) { 'proven' } else { 'failed' }
$evidence = @(
  '# Windows MSI Artifact Evidence'
  ''
  "- Date: $(Get-Date -Format 'yyyy-MM-dd')"
  "- Artifact: ``$resolvedMsi``"
  "- Evidence class: **$ArtifactClass artifact**"
  "- Artifact origin: $ArtifactOrigin"
  "- SHA-256: ``$hash``"
  "- Product version: ``$version``"
  "- ProductName: ``$productName``"
  "- Manufacturer: ``$manufacturer``"
  "- UpgradeCode: ``$upgradeCode``"
  "- $ArtifactClass artifact status: **$status**"
  "- Summary Title (PID 2): ``$summaryTitle``"
  "- Summary Subject (PID 3): ``$summarySubject``"
  "- Summary Author (PID 4): ``$summaryAuthor``"
  "- Summary Comments (PID 6): ``$summaryComments``"
  "- Explorer property-system Subject: ``$explorerSubject``"
  "- Icon row: ``$iconName``"
  "- ARPPRODUCTICON: ``$arpIcon``"
  "- GuiShortcut Icon_: ``$shortcutIcon``"
  "- PATH row: ``$environmentName`` | ``$environmentValue`` | ``$environmentComponent``"
  "- Administrative group row: ``GoScheduleAdminGroup`` | ``$adminGroupName`` | domain ``$adminGroupDomain``"
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

Write-Output (
  "installer-inspect: OK - SHA-256 $hash, Subject '$summarySubject', canonical identity, icon, PATH, and local-group rows proven"
)
