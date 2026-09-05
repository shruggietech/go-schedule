<#
.SYNOPSIS
  Sanity-checks build/windows/goschedule.wxs before the MSI is built in CI.

.DESCRIPTION
  Cheap guard against the WiX source drifting from reality:
    * the three installed executables and installer-private cleanup helper are referenced,
    * the canonical icon feeds installed-apps and both shortcut identities,
    * independently selectable shortcut features retain stable identities,
    * destructive cleanup remains an explicit, guarded commit action,
    * package-owned full UI has one success dialog and guarded finish actions,
    * the Explorer-visible Summary Subject is the approved project copy,
    * the Windows service Name is exactly "goschedd" (the name the CLI
      `gosched service ...` control layer expects),
    * the install folder is "go-schedule" and the package is per-machine.
  If a StageDir is provided, also asserts the referenced .exe files exist there.

.PARAMETER StageDir
  Optional path to the staged build output containing the four .exe files.

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

$expectedInstalledBinaries = @('goschedd.exe', 'gosched-gui.exe', 'gosched.exe')
$expectedStageBinaries = $expectedInstalledBinaries + @('gosched-cleanup.exe')
$fail = @()

foreach ($bin in $expectedInstalledBinaries) {
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

# The package, installed-apps entry, and both shortcuts deliberately share
# one Icon table row. Parse these relationships structurally so harmless XML
# formatting or attribute-order changes do not weaken the check.
try {
  [xml]$wxsXml = $wxs
  $ns = [System.Xml.XmlNamespaceManager]::new($wxsXml.NameTable)
  $ns.AddNamespace('w', 'http://wixtoolset.org/schemas/v4/wxs')
  $ns.AddNamespace('util', 'http://wixtoolset.org/schemas/v4/wxs/util')
  $icon = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Icon[@Id="GoSchedule.ico"]', $ns)
  $summary = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:SummaryInformation', $ns)
  $arpIcon = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Property[@Id="ARPPRODUCTICON"]', $ns)
  $arpNoRemove = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Property[@Id="ARPNOREMOVE"]', $ns)
  $arpNoModify = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Property[@Id="ARPNOMODIFY"]', $ns)
  $modifyPath = $wxsXml.SelectSingleNode('//w:RegistryValue[@Id="ApplicationManagementModifyPath"]', $ns)
  $maintenanceDialog = $wxsXml.SelectSingleNode('//w:Dialog[@Id="GoScheduleMaintenanceTypeDlg"]', $ns)
  $maintenanceRemove = $wxsXml.SelectSingleNode('//w:Dialog[@Id="GoScheduleMaintenanceTypeDlg"]/w:Control[@Id="RemoveButton"]', $ns)
  $shortcut = $wxsXml.SelectSingleNode('//w:Shortcut[@Id="GuiShortcut"]', $ns)
  $desktopShortcut = $wxsXml.SelectSingleNode('//w:Shortcut[@Id="DesktopShortcut"]', $ns)
  $mainFeature = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Feature[@Id="Main"]', $ns)
  $startMenuFeature = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Feature[@Id="Main"]/w:Feature[@Id="StartMenuShortcut"]', $ns)
  $desktopFeature = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Feature[@Id="Main"]/w:Feature[@Id="DesktopShortcut"]', $ns)
  $removeData = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Property[@Id="GOSCHEDULE_REMOVE_DATA"]', $ns)
  $cleanupBinary = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Binary[@Id="GoScheduleCleanup"]', $ns)
  $wipeAction = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:CustomAction[@Id="WipeApplicationData"]', $ns)
  $wipeSequence = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:InstallExecuteSequence/w:Custom[@Action="WipeApplicationData"]', $ns)
  $closeGui = $wxsXml.SelectSingleNode('//util:CloseApplication[@Id="CloseRunningGui"]', $ns)
  $ownedUi = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:UI[@Id="GoScheduleFeatureTreeUI"]', $ns)
  $adminGroup = $wxsXml.SelectSingleNode('//util:Group[@Id="GoScheduleAdminGroup"]', $ns)
  $installingUser = $wxsXml.SelectSingleNode('//util:User[@Id="InstallingUser"]', $ns)
  $adminGroupRef = $wxsXml.SelectSingleNode('//util:User[@Id="InstallingUser"]/util:GroupRef[@Id="GoScheduleAdminGroup"]', $ns)
  $adminComponentRef = $wxsXml.SelectSingleNode('//w:Feature[@Id="Main"]/w:ComponentRef[@Id="AdminAccessProvisioning"]', $ns)
  $managementComponentRef = $wxsXml.SelectSingleNode('//w:Feature[@Id="Main"]/w:ComponentRef[@Id="ApplicationManagementRegistration"]', $ns)

  $expectedSubject = 'go-schedule: cross-platform task scheduler'
  if (-not $summary) {
    $fail += 'SummaryInformation is missing'
  } elseif ($summary.Description -ne $expectedSubject) {
    $fail += "SummaryInformation Description must be `"$expectedSubject`""
  } elseif ($summary.Description.Contains([char]0x2014)) {
    $fail += 'SummaryInformation Description must not contain U+2014'
  }

  if (-not $icon) {
    $fail += 'canonical Icon Id="GoSchedule.ico" is missing'
  } elseif ($icon.SourceFile -ne 'brand/platform/windows/go-schedule.ico') {
    $fail += 'canonical Icon SourceFile must be "brand/platform/windows/go-schedule.ico"'
  } else {
    $repoRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
    $iconPath = Join-Path $repoRoot $icon.SourceFile
    if (-not (Test-Path -LiteralPath $iconPath -PathType Leaf)) {
      $fail += "canonical Icon source does not exist: $iconPath"
    }
  }

  if (-not $arpIcon) {
    $fail += 'ARPPRODUCTICON property is missing'
  } elseif ($arpIcon.Value -ne 'GoSchedule.ico') {
    $fail += 'ARPPRODUCTICON must reference "GoSchedule.ico"'
  }

  if (-not $arpNoRemove -or $arpNoRemove.Value -ne '1') {
    $fail += 'ARPNOREMOVE must be 1 so Windows Settings uses full maintenance before removal'
  }
  if ($arpNoModify) {
    $fail += 'ARPNOMODIFY must remain absent so Windows Settings exposes maintenance'
  }
  if (-not $modifyPath -or $modifyPath.Root -ne 'HKLM' -or
      $modifyPath.Key -ne 'Software\Microsoft\Windows\CurrentVersion\Uninstall\[ProductCode]' -or
      $modifyPath.Name -ne 'ModifyPath' -or $modifyPath.Type -ne 'expandable' -or
      $modifyPath.Value -ne 'MsiExec.exe /I[ProductCode]' -or $modifyPath.KeyPath -ne 'yes') {
    $fail += 'ApplicationManagementModifyPath must own the current product maintenance command'
  }
  if (-not $managementComponentRef) {
    $fail += 'Main feature must reference ApplicationManagementRegistration'
  }
  if (-not $maintenanceDialog -or -not $maintenanceRemove) {
    $fail += 'package-owned maintenance dialog and Remove control are required'
  } elseif ($maintenanceRemove.DisableCondition -match 'ARPNOREMOVE') {
    $fail += 'package-owned maintenance Remove must ignore ARPNOREMOVE'
  }

  if (-not $shortcut) {
    $fail += 'GuiShortcut is missing'
  } elseif ($shortcut.Icon -ne 'GoSchedule.ico') {
    $fail += 'GuiShortcut Icon must reference "GoSchedule.ico"'
  }
  if (-not $desktopShortcut) {
    $fail += 'DesktopShortcut is missing'
  } elseif ($desktopShortcut.Icon -ne 'GoSchedule.ico') {
    $fail += 'DesktopShortcut Icon must reference "GoSchedule.ico"'
  }

  if (-not $mainFeature) {
    $fail += 'Main feature is missing'
  } elseif ($mainFeature.AllowAbsent -ne 'no' -or $mainFeature.Display -ne 'expand') {
    $fail += 'Main feature must be required and initially expanded'
  }
  $shortcutFeatureContracts = @(
    @{ Node = $startMenuFeature; Id = 'StartMenuShortcut'; Level = '1'; Component = 'AppShortcut' },
    @{ Node = $desktopFeature; Id = 'DesktopShortcut'; Level = '2'; Component = 'DesktopShortcutComponent' }
  )
  foreach ($contract in $shortcutFeatureContracts) {
    $feature = $contract.Node
    if (-not $feature) {
      $fail += "$($contract.Id) feature is missing"
      continue
    }
    if ($feature.Level -ne $contract.Level -or $feature.AllowAdvertise -ne 'no' -or $feature.InstallDefault -ne 'local') {
      $fail += "$($contract.Id) feature has the wrong selection contract"
    }
    if (-not $feature.SelectSingleNode("w:ComponentRef[@Id='$($contract.Component)']", $ns)) {
      $fail += "$($contract.Id) feature does not own $($contract.Component)"
    }
  }
  if ($mainFeature -and $mainFeature.SelectSingleNode('w:ComponentRef[@Id="AppShortcut"]', $ns)) {
    $fail += 'AppShortcut must be owned only by StartMenuShortcut'
  }

  if (-not $removeData -or $removeData.Value -ne '0' -or $removeData.Secure -ne 'yes') {
    $fail += 'GOSCHEDULE_REMOVE_DATA must default to 0 and be Secure="yes"'
  }
  $removeDataLaunch = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Launch[contains(@Condition, "GOSCHEDULE_REMOVE_DATA")]', $ns)
  if (-not $removeDataLaunch -or $removeDataLaunch.Condition -notmatch '= "0"' -or $removeDataLaunch.Condition -notmatch '= "1"') {
    $fail += 'GOSCHEDULE_REMOVE_DATA must reject values other than 0 and 1'
  }
  if (-not $cleanupBinary -or $cleanupBinary.SourceFile -ne '$(StageDir)\gosched-cleanup.exe') {
    $fail += 'GoScheduleCleanup must embed $(StageDir)\gosched-cleanup.exe'
  }
  if (-not $wipeAction -or $wipeAction.BinaryRef -ne 'GoScheduleCleanup' -or $wipeAction.ExeCommand -ne 'wipe' -or
      $wipeAction.Execute -ne 'commit' -or $wipeAction.Impersonate -ne 'no' -or $wipeAction.Return -ne 'ignore') {
    $fail += 'WipeApplicationData must be an ignored non-impersonated commit cleanup action'
  }
  $expectedWipeCondition = 'Installed AND REMOVE~="ALL" AND NOT UPGRADINGPRODUCTCODE AND NOT REINSTALL AND GOSCHEDULE_REMOVE_DATA="1"'
  if (-not $wipeSequence -or $wipeSequence.Before -ne 'InstallFinalize' -or $wipeSequence.Condition -ne $expectedWipeCondition) {
    $fail += 'WipeApplicationData has the wrong schedule or condition'
  }
  if (-not $closeGui -or $closeGui.Target -ne 'gosched-gui.exe' -or $closeGui.TerminateProcess -ne '1') {
    $fail += 'CloseRunningGui must terminate gosched-gui.exe before file replacement/removal'
  }

  if (-not $ownedUi) {
    $fail += 'package-owned FeatureTree UI is missing'
  } else {
    foreach ($dialogId in @('GoScheduleMaintenanceTypeDlg', 'GoScheduleUninstallDlg', 'GoScheduleWipeConfirmDlg', 'GoScheduleExitDlg')) {
      if (-not $ownedUi.SelectSingleNode("w:Dialog[@Id='$dialogId']", $ns)) {
        $fail += "$dialogId is missing from the package-owned UI"
      }
    }
    $successRows = $ownedUi.SelectNodes('w:InstallUISequence/w:Show[@OnExit="success"]', $ns)
    if ($successRows.Count -ne 1 -or $successRows[0].Dialog -ne 'GoScheduleExitDlg') {
      $fail += 'package-owned UI must schedule exactly one GoScheduleExitDlg success row'
    }
    $adminSuccessRows = $ownedUi.SelectNodes('w:AdminUISequence/w:Show[@OnExit="success"]', $ns)
    if ($adminSuccessRows.Count -ne 1 -or $adminSuccessRows[0].Dialog -ne 'GoScheduleExitDlg') {
      $fail += 'package-owned UI must schedule exactly one administrative success row'
    }
    $directRemove = $ownedUi.SelectSingleNode('w:InstallUISequence/w:Show[@Dialog="GoScheduleUninstallDlg"]', $ns)
    if (-not $directRemove -or $directRemove.Before -ne 'ProgressDlg' -or $directRemove.Condition -notmatch 'Preselected' -or $directRemove.Condition -notmatch 'REMOVE~="ALL"') {
      $fail += 'direct full-UI uninstall must route through GoScheduleUninstallDlg'
    }
    $maintenanceRemove = $ownedUi.SelectSingleNode('w:Publish[@Dialog="GoScheduleMaintenanceTypeDlg" and @Control="RemoveButton" and @Event="NewDialog"]', $ns)
    if (-not $maintenanceRemove -or $maintenanceRemove.Value -ne 'GoScheduleUninstallDlg' -or $maintenanceRemove.Order -ne '4') {
      $fail += 'maintenance Remove must route through GoScheduleUninstallDlg'
    }
  }
  if ($wxsXml.SelectSingleNode('//*[local-name()="WixUI"]')) {
    $fail += 'stock WixUI composition would introduce competing dialog sequence rows'
  }
  foreach ($actionId in @('LaunchGui', 'OpenDocs')) {
    $action = $wxsXml.SelectSingleNode("/w:Wix/w:Package/w:CustomAction[@Id='$actionId']", $ns)
    if (-not $action -or $action.BinaryRef -ne 'Wix4UtilCA_$(sys.BUILDARCHSHORT)' -or
        $action.DllEntry -ne 'WixUnelevatedShellExec' -or $action.Execute -ne 'immediate' -or $action.Return -ne 'ignore') {
      $fail += "$actionId must use ignored immediate WixUnelevatedShellExec"
    }
  }

  if (-not $adminGroup) {
    $fail += 'administrative util:Group Id="GoScheduleAdminGroup" is missing'
  } else {
    $expectedGroupAttributes = @{
      Name = 'goschedadmin'; CreateGroup = 'yes'
      FailIfExists = 'no'; RemoveOnUninstall = 'no'; UpdateIfExists = 'yes'; Vital = 'yes'
    }
    foreach ($attribute in $expectedGroupAttributes.Keys) {
      if ($adminGroup.GetAttribute($attribute) -ne $expectedGroupAttributes[$attribute]) {
        $fail += "administrative group $attribute must be `"$($expectedGroupAttributes[$attribute])`""
      }
    }
    if ($adminGroup.GetAttribute('Domain')) {
      $fail += 'administrative group Domain must be empty for elevated local-group creation'
    }
  }
  if (-not $installingUser) {
    $fail += 'util:User Id="InstallingUser" is missing'
  } else {
    $expectedUserAttributes = @{
      Name = '[LogonUser]'; Domain = '[%USERDOMAIN]'; CreateUser = 'no'; FailIfExists = 'no'
      RemoveOnUninstall = 'no'; UpdateIfExists = 'yes'; Vital = 'yes'
    }
    foreach ($attribute in $expectedUserAttributes.Keys) {
      if ($installingUser.GetAttribute($attribute) -ne $expectedUserAttributes[$attribute]) {
        $fail += "installing user $attribute must be `"$($expectedUserAttributes[$attribute])`""
      }
    }
  }
  if (-not $adminGroupRef) { $fail += 'installing user is not a member of GoScheduleAdminGroup' }
  if (-not $adminComponentRef) { $fail += 'Main feature does not include AdminAccessProvisioning' }
} catch {
  $fail += "wxs is not valid XML: $($_.Exception.Message)"
}

# The install folder must land on the machine PATH, or every documented bare
# `gosched ...` command fails after a normal install (issue #5). Assert each
# attribute separately so a partial edit, a per-user entry, or one that
# survives uninstall, is reported for what it is rather than passing.
if ($wxs -notmatch '<Environment[^>]*Name="PATH"') {
  $fail += 'no <Environment> element adding INSTALLFOLDER to PATH'
} else {
  $envEl = [regex]::Match($wxs, '<Environment\b[^>]*Name="PATH"[^>]*>').Value
  if ($envEl -notmatch 'Value="\[INSTALLFOLDER\]"') {
    $fail += 'PATH <Environment> Value must be "[INSTALLFOLDER]"'
  }
  if ($envEl -notmatch 'System="yes"') {
    $fail += 'PATH <Environment> must be System="yes" (perMachine package)'
  }
  if ($envEl -notmatch 'Permanent="no"') {
    $fail += 'PATH <Environment> must be Permanent="no" (removed on uninstall)'
  }
  if ($envEl -notmatch 'Part="last"') {
    $fail += 'PATH <Environment> must be Part="last" (append, not replace)'
  }
}

if ($StageDir) {
  foreach ($bin in $expectedStageBinaries) {
    $p = Join-Path $StageDir $bin
    if (-not (Test-Path $p)) { $fail += "staged binary missing: $p" }
  }
}

if ($fail.Count -gt 0) {
  Write-Error ("WiX sanity check failed:`n - " + ($fail -join "`n - "))
  exit 1
}

Write-Output 'WiX sanity check passed.'
