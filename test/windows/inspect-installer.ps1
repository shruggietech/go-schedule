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
  [string]$ArtifactOrigin,

  [string]$CandidateManifestPath,

  [string]$Repository,

  [string]$Tag,

  [string]$Commit,

  [string]$Workflow,

  [long]$RunId,

  [int]$RunAttempt
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

function Get-MsiRows {
  param(
    [Parameter(Mandatory)] $Database,
    [Parameter(Mandatory)] [string]$Query
  )

  $view = $Database.OpenView($Query)
  $view.Execute() | Out-Null
  $columns = $view.ColumnInfo(0)
  $rows = [System.Collections.Generic.List[object]]::new()
  while ($record = $view.Fetch()) {
    $values = [ordered]@{}
    for ($field = 1; $field -le $columns.FieldCount(); $field++) {
      $values[$columns.StringData($field)] = $record.StringData($field)
    }
    $rows.Add([pscustomobject]$values)
  }
  @($rows)
}

function Get-MsiRowsOrEmpty {
  param(
    [Parameter(Mandatory)] $Database,
    [Parameter(Mandatory)] [string]$Query
  )

  try {
    @(Get-MsiRows -Database $Database -Query $Query)
  } catch {
    @()
  }
}

function Require-SingleMsiRow {
  param(
    [Parameter(Mandatory)] $Database,
    [Parameter(Mandatory)] [string]$Query,
    [Parameter(Mandatory)] [string]$Description,
    [Parameter(Mandatory)]
    [AllowEmptyCollection()]
    [System.Collections.Generic.List[string]]$Failures
  )

  $rows = @(Get-MsiRowsOrEmpty -Database $Database -Query $Query)
  if ($rows.Count -ne 1) {
    $Failures.Add("$Description row count is $($rows.Count); expected 1")
    return $null
  }
  $rows[0]
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

$productCode = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='ProductCode'"
if (-not $productCode) {
  $fail.Add('Property.ProductCode is missing')
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

# S039 installer lifecycle relationships. These assertions intentionally read
# the compiled database rather than treating source XML as artifact evidence.
$featureEvidence = [System.Collections.Generic.List[string]]::new()
foreach ($featureContract in @(
  @{ Id = 'StartMenuShortcut'; Parent = 'Main'; Level = '1'; Shortcut = 'GuiShortcut'; Directory = 'AppMenuFolder' },
  @{ Id = 'DesktopShortcut'; Parent = 'Main'; Level = '2'; Shortcut = 'DesktopShortcut'; Directory = 'DesktopFolder' }
)) {
  $feature = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Feature_Parent``, ``Level`` FROM ``Feature`` WHERE ``Feature``='$($featureContract.Id)'" `
    -Description "Feature.$($featureContract.Id)" -Failures $fail
  if ($feature) {
    if ($feature.Feature_Parent -ne $featureContract.Parent) {
      $fail.Add("Feature.$($featureContract.Id).Feature_Parent is '$($feature.Feature_Parent)'; expected '$($featureContract.Parent)'")
    }
    if ($feature.Level -ne $featureContract.Level) {
      $fail.Add("Feature.$($featureContract.Id).Level is '$($feature.Level)'; expected '$($featureContract.Level)'")
    }
  }

  $mapping = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Component_`` FROM ``FeatureComponents`` WHERE ``Feature_``='$($featureContract.Id)'" `
    -Description "FeatureComponents.$($featureContract.Id)" -Failures $fail
  $shortcut = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Directory_``, ``Component_``, ``Target``, ``Description``, ``Icon_``, ``WkDir`` FROM ``Shortcut`` WHERE ``Shortcut``='$($featureContract.Shortcut)'" `
    -Description "Shortcut.$($featureContract.Shortcut)" -Failures $fail
  if ($shortcut) {
    if ($shortcut.Directory_ -ne $featureContract.Directory) {
      $fail.Add("Shortcut.$($featureContract.Shortcut).Directory_ is '$($shortcut.Directory_)'; expected '$($featureContract.Directory)'")
    }
    if ($mapping -and $shortcut.Component_ -ne $mapping.Component_) {
      $fail.Add("Shortcut.$($featureContract.Shortcut) is not owned by the component mapped to Feature.$($featureContract.Id)")
    }
    if ($shortcut.Target -ne '[INSTALLFOLDER]gosched-gui.exe') {
      $fail.Add("Shortcut.$($featureContract.Shortcut).Target is '$($shortcut.Target)'; expected '[INSTALLFOLDER]gosched-gui.exe'")
    }
    if ($shortcut.Description -ne 'Open the go-schedule desktop app') {
      $fail.Add("Shortcut.$($featureContract.Shortcut).Description is not canonical")
    }
    if ($shortcut.Icon_ -ne $canonicalIcon) {
      $fail.Add("Shortcut.$($featureContract.Shortcut).Icon_ is '$($shortcut.Icon_)'; expected '$canonicalIcon'")
    }
    if ($shortcut.WkDir -ne 'INSTALLFOLDER') {
      $fail.Add("Shortcut.$($featureContract.Shortcut).WkDir is '$($shortcut.WkDir)'; expected 'INSTALLFOLDER'")
    }
  }
  $mappedComponent = if ($mapping) { $mapping.Component_ } else { '<missing>' }
  $featureEvidence.Add("$($featureContract.Id): level=$($feature.Level); component=$mappedComponent; shortcut=$($featureContract.Shortcut)")
}

$dialogEvidence = [System.Collections.Generic.List[string]]::new()
foreach ($dialogId in @('GoScheduleUninstallDlg','GoScheduleWipeConfirmDlg','GoScheduleExitDlg')) {
  $dialog = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Dialog`` FROM ``Dialog`` WHERE ``Dialog``='$dialogId'" `
    -Description "Dialog.$dialogId" -Failures $fail
  if ($dialog) { $dialogEvidence.Add($dialogId) }
}

$launchControl = Require-SingleMsiRow -Database $database `
  -Query "SELECT ``Dialog_``, ``Control``, ``Type``, ``Property`` FROM ``Control`` WHERE ``Dialog_``='GoScheduleExitDlg' AND ``Property``='LAUNCH_GOSCHEDULE'" `
  -Description 'GoScheduleExitDlg.LAUNCH_GOSCHEDULE control' -Failures $fail
$docsControl = Require-SingleMsiRow -Database $database `
  -Query "SELECT ``Dialog_``, ``Control``, ``Type``, ``Property`` FROM ``Control`` WHERE ``Dialog_``='GoScheduleExitDlg' AND ``Property``='OPEN_GOSCHEDULE_DOCS'" `
  -Description 'GoScheduleExitDlg.OPEN_GOSCHEDULE_DOCS control' -Failures $fail
foreach ($controlContract in @(
  @{ Row = $launchControl; Description = 'LAUNCH_GOSCHEDULE' },
  @{ Row = $docsControl; Description = 'OPEN_GOSCHEDULE_DOCS' }
)) {
  if ($controlContract.Row -and $controlContract.Row.Type -ne 'CheckBox') {
    $fail.Add("GoScheduleExitDlg $($controlContract.Description) control type is '$($controlContract.Row.Type)'; expected 'CheckBox'")
  }
}

$removeMode = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='GOSCHEDULE_REMOVE_DATA'"
if ($removeMode -ne '0') {
  $fail.Add("Property.GOSCHEDULE_REMOVE_DATA is '$removeMode'; expected preserve default '0'")
}
$launchDefault = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='LAUNCH_GOSCHEDULE'"
if ($launchDefault -ne '1') {
  $fail.Add("Property.LAUNCH_GOSCHEDULE is '$launchDefault'; expected '1'")
}
$docsDefault = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='OPEN_GOSCHEDULE_DOCS'"
if ($docsDefault) {
  $fail.Add("Property.OPEN_GOSCHEDULE_DOCS is '$docsDefault'; expected unselected/absent")
}
$wipeConfirmation = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='GOSCHEDULE_WIPE_CONFIRMED'"
if ($wipeConfirmation) {
  $fail.Add('GOSCHEDULE_WIPE_CONFIRMED must be session-only and unselected by default')
}
$secureProperties = Get-MsiString -Database $database `
  -Query "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='SecureCustomProperties'"
$securePropertySet = @($secureProperties -split ';' | Where-Object { $_ })
if ($securePropertySet -notcontains 'GOSCHEDULE_REMOVE_DATA') {
  $fail.Add('SecureCustomProperties does not contain GOSCHEDULE_REMOVE_DATA')
}

$launchActionEvidence = [System.Collections.Generic.List[string]]::new()
foreach ($actionId in @('LaunchGui','OpenDocs')) {
  $action = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Type``, ``Source``, ``Target`` FROM ``CustomAction`` WHERE ``Action``='$actionId'" `
    -Description "CustomAction.$actionId" -Failures $fail
  if ($action) {
    if ($action.Target -ne 'WixUnelevatedShellExec') {
      $fail.Add("CustomAction.$actionId.Target is '$($action.Target)'; expected WixUnelevatedShellExec")
    }
    if ($action.Source -notmatch '^Wix4UtilCA_') {
      $fail.Add("CustomAction.$actionId.Source is '$($action.Source)'; expected a WiX Util custom-action binary")
    }
    $launchActionEvidence.Add("${actionId}: source=$($action.Source); target=$($action.Target)")
  }
  $executeRows = @(Get-MsiRowsOrEmpty -Database $database `
    -Query "SELECT ``Action`` FROM ``InstallExecuteSequence`` WHERE ``Action``='$actionId'")
  if ($executeRows.Count -ne 0) {
    $fail.Add("Completion action $actionId must not appear in InstallExecuteSequence")
  }
  $uiRows = @(Get-MsiRowsOrEmpty -Database $database `
    -Query "SELECT ``Action`` FROM ``InstallUISequence`` WHERE ``Action``='$actionId'")
  if ($uiRows.Count -ne 0) {
    $fail.Add("Completion action $actionId must be driven only by GoScheduleExitDlg controls")
  }
}

$finishEvents = @(Get-MsiRowsOrEmpty -Database $database `
  -Query "SELECT ``Event``, ``Argument``, ``Condition``, ``Ordering`` FROM ``ControlEvent`` WHERE ``Dialog_``='GoScheduleExitDlg' AND ``Control_``='Finish'")
foreach ($finishContract in @(
  @{ Action = 'LaunchGui'; Property = 'LAUNCH_GOSCHEDULE' },
  @{ Action = 'OpenDocs'; Property = 'OPEN_GOSCHEDULE_DOCS' }
)) {
  $matches = @($finishEvents | Where-Object {
    $_.Event -eq 'DoAction' -and $_.Argument -eq $finishContract.Action -and
    $_.Condition -match [regex]::Escape($finishContract.Property)
  })
  if ($matches.Count -ne 1) {
    $fail.Add("GoScheduleExitDlg.Finish has $($matches.Count) conditional $($finishContract.Action) events; expected 1")
  }
}

$cleanupBinary = Require-SingleMsiRow -Database $database `
  -Query "SELECT ``Name`` FROM ``Binary`` WHERE ``Name``='GoScheduleCleanup'" `
  -Description 'Binary.GoScheduleCleanup' -Failures $fail
$wipeAction = Require-SingleMsiRow -Database $database `
  -Query "SELECT ``Type``, ``Source``, ``Target`` FROM ``CustomAction`` WHERE ``Action``='WipeApplicationData'" `
  -Description 'CustomAction.WipeApplicationData' -Failures $fail
$wipeActionType = 0
if ($wipeAction) {
  $wipeActionType = [int]$wipeAction.Type
  if ($wipeAction.Source -ne 'GoScheduleCleanup') {
    $fail.Add("CustomAction.WipeApplicationData.Source is '$($wipeAction.Source)'; expected GoScheduleCleanup")
  }
  if ($wipeAction.Target -ne 'wipe') {
    $fail.Add("CustomAction.WipeApplicationData.Target is '$($wipeAction.Target)'; expected fixed verb 'wipe'")
  }
  # Type 2 executable + continue/ignore + in-script commit + no impersonation.
  if (($wipeActionType -band 0x3F) -ne 2 -or
      ($wipeActionType -band 0x40) -eq 0 -or
      ($wipeActionType -band 0x400) -eq 0 -or
      ($wipeActionType -band 0x200) -eq 0 -or
      ($wipeActionType -band 0x800) -eq 0) {
    $fail.Add("CustomAction.WipeApplicationData.Type is $wipeActionType; expected embedded EXE, ignored return, commit execution, and no impersonation")
  }
}

$wipeSequence = Require-SingleMsiRow -Database $database `
  -Query "SELECT ``Condition``, ``Sequence`` FROM ``InstallExecuteSequence`` WHERE ``Action``='WipeApplicationData'" `
  -Description 'InstallExecuteSequence.WipeApplicationData' -Failures $fail
$wipeCondition = if ($wipeSequence) { [string]$wipeSequence.Condition } else { '' }
$normalizedWipeCondition = ($wipeCondition -replace '[\s()]', '').ToUpperInvariant()
foreach ($conditionFragment in @(
  'INSTALLED',
  'REMOVE~="ALL"',
  'NOTUPGRADINGPRODUCTCODE',
  'NOTREINSTALL',
  'GOSCHEDULE_REMOVE_DATA="1"'
)) {
  if (-not $normalizedWipeCondition.Contains($conditionFragment)) {
    $fail.Add("WipeApplicationData condition '$wipeCondition' is missing $conditionFragment")
  }
}

$invalidActionRows = @(Get-MsiRowsOrEmpty -Database $database `
  -Query "SELECT ``Type`` FROM ``CustomAction`` WHERE ``Action``='RejectInvalidRemoveData'")
$invalidCondition = ''
if ($invalidActionRows.Count -eq 1) {
  if (([int]$invalidActionRows[0].Type -band 0x3F) -ne 19) {
    $fail.Add("CustomAction.RejectInvalidRemoveData.Type is $($invalidActionRows[0].Type); expected Type 19 error")
  }
  $invalidSequence = Require-SingleMsiRow -Database $database `
    -Query "SELECT ``Condition``, ``Sequence`` FROM ``InstallExecuteSequence`` WHERE ``Action``='RejectInvalidRemoveData'" `
    -Description 'InstallExecuteSequence.RejectInvalidRemoveData' -Failures $fail
  if ($invalidSequence) { $invalidCondition = [string]$invalidSequence.Condition }
} elseif ($invalidActionRows.Count -eq 0) {
  $launchGuards = @(Get-MsiRowsOrEmpty -Database $database `
    -Query "SELECT ``Condition``, ``Description`` FROM ``LaunchCondition``")
  $guard = @($launchGuards | Where-Object {
    $_.Condition -match 'GOSCHEDULE_REMOVE_DATA'
  })
  if ($guard.Count -ne 1) {
    $fail.Add("LaunchCondition removal-value guard row count is $($guard.Count); expected 1")
  } else {
    $invalidCondition = [string]$guard[0].Condition
  }
} else {
  $fail.Add("CustomAction.RejectInvalidRemoveData row count is $($invalidActionRows.Count); expected at most 1")
}
if ($invalidCondition -notmatch 'GOSCHEDULE_REMOVE_DATA' -or
    $invalidCondition -notmatch '0' -or $invalidCondition -notmatch '1') {
  $fail.Add("RejectInvalidRemoveData condition '$invalidCondition' does not constrain GOSCHEDULE_REMOVE_DATA to 0 or 1")
}

$closeEvidence = [System.Collections.Generic.List[string]]::new()
foreach ($closeContract in @(
  @{ Id = 'CloseStaleDaemon'; Target = 'goschedd.exe' },
  @{ Id = 'CloseRunningGui'; Target = 'gosched-gui.exe' }
)) {
  $closeRow = Require-SingleMsiRow -Database $database `
    -Query "SELECT * FROM ``Wix4CloseApplication`` WHERE ``CloseApplication``='$($closeContract.Id)'" `
    -Description "Wix4CloseApplication.$($closeContract.Id)" -Failures $fail
  if ($closeRow) {
    if ($closeRow.Target -ne $closeContract.Target) {
      $fail.Add("Wix4CloseApplication.$($closeContract.Id).Target is '$($closeRow.Target)'; expected '$($closeContract.Target)'")
    }
    $attributes = [int]$closeRow.Attributes
    if ($attributes -ne 32) {
      $fail.Add("Wix4CloseApplication.$($closeContract.Id).Attributes is $attributes; expected TerminateProcess=1 and RebootPrompt=no (32)")
    }
    $closeEvidence.Add("$($closeContract.Id): target=$($closeRow.Target); attributes=$attributes")
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
  "- Shortcut features: ``$($featureEvidence -join '; ')``"
  "- Lifecycle dialogs: ``$($dialogEvidence -join ', ')``"
  "- Completion actions: ``$($launchActionEvidence -join '; ')``"
  "- Removal property: ``GOSCHEDULE_REMOVE_DATA=$removeMode``; secure=$($securePropertySet -contains 'GOSCHEDULE_REMOVE_DATA')"
  "- Cleanup action: binary=``$($cleanupBinary.Name)``; type=$wipeActionType; condition=``$wipeCondition``"
  "- Invalid removal guard: condition=``$invalidCondition``"
  "- Close-application rows: ``$($closeEvidence -join '; ')``"
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

if ($CandidateManifestPath) {
  if ($ArtifactClass -ne 'candidate') {
    throw 'Candidate manifest output requires ArtifactClass candidate.'
  }
  if (-not $Repository -or -not $Tag -or -not $Commit -or -not $Workflow -or
      $RunId -le 0 -or $RunAttempt -le 0) {
    throw 'Candidate manifest output requires complete workflow identity.'
  }
  if ($Tag -notmatch '^v[0-9]+\.[0-9]+\.[0-9]+$' -or
      $Commit -notmatch '^[0-9a-f]{40}$') {
    throw 'Candidate manifest tag or commit has invalid syntax.'
  }
  if ($version -ne $Tag.TrimStart('v')) {
    throw "MSI ProductVersion '$version' does not match tag '$Tag'."
  }
  $candidateManifest = [ordered]@{
    repository = $Repository
    tag = $Tag
    commit = $Commit
    workflow = $Workflow
    run_id = $RunId
    run_attempt = $RunAttempt
    filename = [System.IO.Path]::GetFileName($resolvedMsi)
    bytes = (Get-Item -LiteralPath $resolvedMsi).Length
    sha256 = $hash
    product_version = $version
    product_code = $productCode.ToUpperInvariant()
  }
  $candidateManifestFile = [System.IO.Path]::GetFullPath(
    $CandidateManifestPath
  )
  $candidateManifestJson = $candidateManifest |
    ConvertTo-Json -Depth 5
  $candidateStream = [System.IO.File]::Open(
    $candidateManifestFile,
    [System.IO.FileMode]::CreateNew,
    [System.IO.FileAccess]::Write,
    [System.IO.FileShare]::None
  )
  $candidateComplete = $false
  try {
    $candidateWriter = [System.IO.StreamWriter]::new(
      $candidateStream,
      [System.Text.UTF8Encoding]::new($false)
    )
    try {
      $candidateWriter.Write("$candidateManifestJson`n")
      $candidateWriter.Flush()
      $candidateComplete = $true
    } finally {
      $candidateWriter.Dispose()
    }
  } finally {
    $candidateStream.Dispose()
    if (-not $candidateComplete -and
        (Test-Path -LiteralPath $candidateManifestFile)) {
      [System.IO.File]::Delete($candidateManifestFile)
    }
  }
  Write-Output (
    "installer-inspect: candidate manifest written to $candidateManifestFile"
  )
}

Write-Output (
  "installer-inspect: OK - SHA-256 $hash, Subject '$summarySubject', canonical identity, icon, PATH, and local-group rows proven"
)
