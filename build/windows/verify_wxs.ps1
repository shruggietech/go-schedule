<#
.SYNOPSIS
  Sanity-checks build/windows/goschedule.wxs before the MSI is built in CI.

.DESCRIPTION
  Cheap guard against the WiX source drifting from reality:
    * the three expected binaries are referenced as File sources,
    * the canonical icon feeds both installed-apps and Start Menu identity,
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

# The package, installed-apps entry, and Start Menu shortcut deliberately share
# one Icon table row. Parse these relationships structurally so harmless XML
# formatting or attribute-order changes do not weaken the check.
try {
  [xml]$wxsXml = $wxs
  $ns = [System.Xml.XmlNamespaceManager]::new($wxsXml.NameTable)
  $ns.AddNamespace('w', 'http://wixtoolset.org/schemas/v4/wxs')
  $ns.AddNamespace('util', 'http://wixtoolset.org/schemas/v4/wxs/util')
  $icon = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Icon[@Id="GoSchedule.ico"]', $ns)
  $arpIcon = $wxsXml.SelectSingleNode('/w:Wix/w:Package/w:Property[@Id="ARPPRODUCTICON"]', $ns)
  $shortcut = $wxsXml.SelectSingleNode('//w:Shortcut[@Id="GuiShortcut"]', $ns)
  $adminGroup = $wxsXml.SelectSingleNode('//util:Group[@Id="GoScheduleAdminGroup"]', $ns)
  $installingUser = $wxsXml.SelectSingleNode('//util:User[@Id="InstallingUser"]', $ns)
  $adminGroupRef = $wxsXml.SelectSingleNode('//util:User[@Id="InstallingUser"]/util:GroupRef[@Id="GoScheduleAdminGroup"]', $ns)
  $adminComponentRef = $wxsXml.SelectSingleNode('//w:Feature[@Id="Main"]/w:ComponentRef[@Id="AdminAccessProvisioning"]', $ns)

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

  if (-not $shortcut) {
    $fail += 'GuiShortcut is missing'
  } elseif ($shortcut.Icon -ne 'GoSchedule.ico') {
    $fail += 'GuiShortcut Icon must reference "GoSchedule.ico"'
  }

  if (-not $adminGroup) {
    $fail += 'administrative util:Group Id="GoScheduleAdminGroup" is missing'
  } else {
    $expectedGroupAttributes = @{
      Name = 'goschedadmin'; Domain = '[ComputerName]'; CreateGroup = 'yes'
      FailIfExists = 'no'; RemoveOnUninstall = 'no'; UpdateIfExists = 'yes'; Vital = 'yes'
    }
    foreach ($attribute in $expectedGroupAttributes.Keys) {
      if ($adminGroup.GetAttribute($attribute) -ne $expectedGroupAttributes[$attribute]) {
        $fail += "administrative group $attribute must be `"$($expectedGroupAttributes[$attribute])`""
      }
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
# attribute separately so a partial edit — a per-user entry, or one that
# survives uninstall — is reported for what it is rather than passing.
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
  foreach ($bin in $expectedBinaries) {
    $p = Join-Path $StageDir $bin
    if (-not (Test-Path $p)) { $fail += "staged binary missing: $p" }
  }
}

if ($fail.Count -gt 0) {
  Write-Error ("WiX sanity check failed:`n - " + ($fail -join "`n - "))
  exit 1
}

Write-Output 'WiX sanity check passed.'
