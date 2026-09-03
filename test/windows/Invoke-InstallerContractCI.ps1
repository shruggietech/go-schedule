<#
.SYNOPSIS
  Exercises the S039 MSI's silent installer contract on a disposable CI host.

.DESCRIPTION
  This probe is deliberately restricted to GitHub-hosted Windows automation.
  It verifies compiled feature behavior, silent non-launch, maintenance,
  upgrade migration, preserve-by-default removal, explicit wipe, invalid-mode
  rejection, sentinels, security-state preservation, and cleanup-result
  evidence. It is not the clean Windows 11 attended-desktop gate owned by #94.
#>
[CmdletBinding(SupportsShouldProcess=$true,ConfirmImpact='High')]
param(
  [Parameter(Mandatory)]
  [string]$BaselineMsiPath,

  [Parameter(Mandatory)]
  [string]$MsiPath,

  [Parameter(Mandatory)]
  [string]$EvidencePath,

  [Parameter(Mandatory)]
  [ValidateNotNullOrEmpty()]
  [string]$ArtifactOrigin
)

$ErrorActionPreference = 'Stop'
$adminGroupName = 'goschedadmin'
$installDirectory = 'C:\Program Files\go-schedule'
$machineDataRoot = Join-Path $env:ProgramData 'goschedule'
$profileDataRoot = Join-Path $env:APPDATA 'fyne\tech.shruggie.goschedule'
$machineSentinel = Join-Path $env:ProgramData 'goschedule-s039-sentinel.txt'
$profileSentinel = Join-Path $env:APPDATA 'fyne\tech.shruggie.goschedule-s039-sentinel.txt'
$operations = [System.Collections.Generic.List[object]]::new()
$observations = [System.Collections.Generic.List[object]]::new()
$problems = [System.Collections.Generic.List[string]]::new()

function Get-MsiProperty {
  param(
    [Parameter(Mandatory)] [string]$PackagePath,
    [Parameter(Mandatory)] [string]$PropertyName
  )

  $installer = New-Object -ComObject WindowsInstaller.Installer
  $database = $installer.OpenDatabase($PackagePath, 0)
  $view = $database.OpenView(
    "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='$PropertyName'"
  )
  $view.Execute() | Out-Null
  $record = $view.Fetch()
  if (-not $record) { return '' }
  $record.StringData(1)
}

function Test-ProductInstalled {
  param([Parameter(Mandatory)] [string]$ProductCode)

  $key = "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\$ProductCode"
  Test-Path -LiteralPath $key
}

function Assert-ProductState {
  param(
    [Parameter(Mandatory)] [string]$Label,
    [Parameter(Mandatory)] [string]$ProductCode,
    [Parameter(Mandatory)] [bool]$Installed
  )

  $actual = Test-ProductInstalled -ProductCode $ProductCode
  if ($actual -ne $Installed) {
    throw "$Label product registration is $actual; expected $Installed"
  }
  $observations.Add([pscustomobject]@{
    kind = 'product-state'
    label = $Label
    product_code = $ProductCode
    installed = $actual
  })
}

function Get-ShortcutPaths {
  $commonPrograms = [Environment]::GetFolderPath(
    [Environment+SpecialFolder]::CommonPrograms
  )
  $commonDesktop = [Environment]::GetFolderPath(
    [Environment+SpecialFolder]::CommonDesktopDirectory
  )
  [pscustomobject]@{
    StartMenu = Join-Path $commonPrograms 'go-schedule\go-schedule.lnk'
    Desktop = Join-Path $commonDesktop 'go-schedule.lnk'
  }
}

function Test-EquivalentWindowsPath {
  param(
    [Parameter(Mandatory)] [string]$Actual,
    [Parameter(Mandatory)] [string]$Expected
  )

  $actualFull = [IO.Path]::GetFullPath($Actual).TrimEnd('\', '/')
  $expectedFull = [IO.Path]::GetFullPath($Expected).TrimEnd('\', '/')
  return $actualFull.Equals($expectedFull, [StringComparison]::OrdinalIgnoreCase)
}

function Assert-ShortcutState {
  param(
    [Parameter(Mandatory)] [string]$Label,
    [Parameter(Mandatory)] [bool]$StartMenu,
    [Parameter(Mandatory)] [bool]$Desktop
  )

  $paths = Get-ShortcutPaths
  foreach ($contract in @(
    @{ Kind = 'StartMenu'; Path = $paths.StartMenu; Expected = $StartMenu },
    @{ Kind = 'Desktop'; Path = $paths.Desktop; Expected = $Desktop }
  )) {
    $present = Test-Path -LiteralPath $contract.Path -PathType Leaf
    if ($present -ne $contract.Expected) {
      throw "$Label $($contract.Kind) shortcut presence is $present; expected $($contract.Expected)"
    }
    if ($present) {
      $shell = New-Object -ComObject WScript.Shell
      $shortcut = $shell.CreateShortcut($contract.Path)
      $expectedTarget = Join-Path $installDirectory 'gosched-gui.exe'
      if (-not (Test-EquivalentWindowsPath -Actual $shortcut.TargetPath -Expected $expectedTarget)) {
        throw "$Label $($contract.Kind) shortcut target '$($shortcut.TargetPath)' is not '$expectedTarget'"
      }
      if (-not (Test-EquivalentWindowsPath -Actual $shortcut.WorkingDirectory -Expected $installDirectory)) {
        throw "$Label $($contract.Kind) shortcut working directory '$($shortcut.WorkingDirectory)' is not '$installDirectory'"
      }
      if ($shortcut.Description -ne 'Open the go-schedule desktop app') {
        throw "$Label $($contract.Kind) shortcut description '$($shortcut.Description)' is not canonical"
      }
    }
  }
  $observations.Add([pscustomobject]@{
    kind = 'shortcut-state'
    label = $Label
    start_menu = $StartMenu
    desktop = $Desktop
    start_menu_path = $paths.StartMenu
    desktop_path = $paths.Desktop
  })
}

function Assert-NoCompletionLaunch {
  param(
    [Parameter(Mandatory)] [string]$Label,
    [Parameter(Mandatory)] [string]$LogPath,
    [Parameter(Mandatory)]
    [AllowEmptyCollection()]
    [int[]]$GuiProcessIdsBefore
  )

  $newGui = @(Get-Process -Name 'gosched-gui' -ErrorAction SilentlyContinue |
    Where-Object { $_.Id -notin $GuiProcessIdsBefore })
  if ($newGui.Count) {
    throw "$Label launched gosched-gui.exe during a silent installer operation"
  }
  $unexpectedActions = @(Select-String -LiteralPath $LogPath `
    -Pattern 'Action start.*(LaunchGui|OpenDocs)' -ErrorAction SilentlyContinue)
  if ($unexpectedActions.Count) {
    throw "$Label executed a completion action during a silent installer operation"
  }
}

function Invoke-MsiOperation {
  param(
    [Parameter(Mandatory)] [string]$Label,
    [Parameter(Mandatory)] [string]$Operation,
    [Parameter(Mandatory)] [string]$Package,
    [string[]]$Properties = @(),
    [switch]$ExpectFailure
  )

  $logPath = Join-Path $script:LogDirectory "$Label.log"
  $arguments = @(
    $Operation,
    ('"{0}"' -f $Package)
  ) + @($Properties) + @(
    '/qn',
    '/norestart',
    '/L*v',
    ('"{0}"' -f $logPath)
  )
  $guiBefore = @(Get-Process -Name 'gosched-gui' -ErrorAction SilentlyContinue |
    Select-Object -ExpandProperty Id)
  $process = Start-Process -FilePath 'msiexec.exe' -ArgumentList $arguments `
    -Wait -PassThru -WindowStyle Hidden
  if (-not (Test-Path -LiteralPath $logPath -PathType Leaf)) {
    throw "$Label did not produce a verbose MSI log"
  }
  Assert-NoCompletionLaunch -Label $Label -LogPath $logPath `
    -GuiProcessIdsBefore $guiBefore
  $accepted = $process.ExitCode -in 0, 3010
  if ($ExpectFailure -and $accepted) {
    throw "$Label unexpectedly succeeded"
  }
  if (-not $ExpectFailure -and -not $accepted) {
    throw "$Label failed with exit code $($process.ExitCode)"
  }
  $operations.Add([pscustomobject]@{
    label = $Label
    operation = $Operation
    properties = @($Properties)
    exit_code = $process.ExitCode
    expected_failure = [bool]$ExpectFailure
    log_path = $logPath
    log_sha256 = (Get-FileHash -LiteralPath $logPath -Algorithm SHA256).Hash.ToLowerInvariant()
    completion_actions_observed = $false
  })
}

function Write-SeedData {
  param([Parameter(Mandatory)] [string]$Label)

  New-Item -ItemType Directory -Path $machineDataRoot -Force | Out-Null
  New-Item -ItemType Directory -Path $profileDataRoot -Force | Out-Null
  New-Item -ItemType Directory -Path (Split-Path $profileSentinel -Parent) `
    -Force | Out-Null
  [IO.File]::WriteAllText(
    (Join-Path $machineDataRoot 's039-machine.txt'),
    "$Label-machine",
    [Text.UTF8Encoding]::new($false)
  )
  [IO.File]::WriteAllText(
    (Join-Path $profileDataRoot 'preferences.json'),
    "{`"s039`":`"$Label-profile`"}",
    [Text.UTF8Encoding]::new($false)
  )
  [IO.File]::WriteAllText(
    $machineSentinel,
    "$Label-machine-sentinel",
    [Text.UTF8Encoding]::new($false)
  )
  [IO.File]::WriteAllText(
    $profileSentinel,
    "$Label-profile-sentinel",
    [Text.UTF8Encoding]::new($false)
  )
  $seed = [pscustomobject]@{
    machine_file = Join-Path $machineDataRoot 's039-machine.txt'
    machine_sha256 = (Get-FileHash -LiteralPath (Join-Path $machineDataRoot 's039-machine.txt') -Algorithm SHA256).Hash.ToLowerInvariant()
    profile_file = Join-Path $profileDataRoot 'preferences.json'
    profile_sha256 = (Get-FileHash -LiteralPath (Join-Path $profileDataRoot 'preferences.json') -Algorithm SHA256).Hash.ToLowerInvariant()
    machine_sentinel_sha256 = (Get-FileHash -LiteralPath $machineSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
    profile_sentinel_sha256 = (Get-FileHash -LiteralPath $profileSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
  }
  $observations.Add([pscustomobject]@{
    kind = 'populated-data-inventory'
    label = $Label
    machine_root = $machineDataRoot
    profile_root = $profileDataRoot
    inventory = $seed
  })
  $seed
}

function Assert-SeedPreserved {
  param(
    [Parameter(Mandatory)] $Seed,
    [Parameter(Mandatory)] [string]$Label
  )

  foreach ($contract in @(
    @{ Path = $Seed.machine_file; Hash = $Seed.machine_sha256 },
    @{ Path = $Seed.profile_file; Hash = $Seed.profile_sha256 },
    @{ Path = $machineSentinel; Hash = $Seed.machine_sentinel_sha256 },
    @{ Path = $profileSentinel; Hash = $Seed.profile_sentinel_sha256 }
  )) {
    if (-not (Test-Path -LiteralPath $contract.Path -PathType Leaf)) {
      throw "$Label did not preserve $($contract.Path)"
    }
    $hash = (Get-FileHash -LiteralPath $contract.Path -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($hash -ne $contract.Hash) {
      throw "$Label changed $($contract.Path)"
    }
  }
  $observations.Add([pscustomobject]@{
    kind = 'preserved-data'
    label = $Label
    machine_file_sha256 = (Get-FileHash -LiteralPath $Seed.machine_file -Algorithm SHA256).Hash.ToLowerInvariant()
    profile_file_sha256 = (Get-FileHash -LiteralPath $Seed.profile_file -Algorithm SHA256).Hash.ToLowerInvariant()
    machine_sentinel_sha256 = (Get-FileHash -LiteralPath $machineSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
    profile_sentinel_sha256 = (Get-FileHash -LiteralPath $profileSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
  })
}

function Assert-WipeState {
  param(
    [Parameter(Mandatory)] $Seed,
    [Parameter(Mandatory)] [string]$Label,
    [Parameter(Mandatory)] [bool]$ExpectComplete
  )

  $ownedRemain = (Test-Path -LiteralPath $machineDataRoot) -or
    (Test-Path -LiteralPath $profileDataRoot)
  if ($ExpectComplete -and $ownedRemain) {
    throw "$Label left a declared owned root after a complete wipe"
  }
  if (-not $ExpectComplete -and -not $ownedRemain) {
    throw "$Label did not retain the expected locked residual"
  }
  foreach ($contract in @(
    @{ Path = $machineSentinel; Hash = $Seed.machine_sentinel_sha256 },
    @{ Path = $profileSentinel; Hash = $Seed.profile_sentinel_sha256 }
  )) {
    if (-not (Test-Path -LiteralPath $contract.Path -PathType Leaf)) {
      throw "$Label removed out-of-scope sentinel $($contract.Path)"
    }
    $hash = (Get-FileHash -LiteralPath $contract.Path -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($hash -ne $contract.Hash) {
      throw "$Label changed out-of-scope sentinel $($contract.Path)"
    }
  }
  $observations.Add([pscustomobject]@{
    kind = 'wipe-state'
    label = $Label
    expected_complete = $ExpectComplete
    owned_roots_remain = $ownedRemain
    machine_sentinel_sha256 = (Get-FileHash -LiteralPath $machineSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
    profile_sentinel_sha256 = (Get-FileHash -LiteralPath $profileSentinel -Algorithm SHA256).Hash.ToLowerInvariant()
  })
}

function Get-CleanupResultEvidence {
  $key = 'HKLM:\Software\ShruggieTech\go-schedule-uninstall'
  $defaultReport = Join-Path $env:ProgramData (
    'ShruggieTech\go-schedule-uninstall\' +
    'b6f3c2e1-7a4d-4c9e-9b2a-1f6d8e5a0c34\cleanup-result.json'
  )
  $keyPresent = Test-Path -LiteralPath $key
  $values = if ($keyPresent) { Get-ItemProperty -LiteralPath $key } else { $null }
  $reportPath = if ($values -and $values.ReportPath) {
    [string]$values.ReportPath
  } else { $defaultReport }
  $reportPresent = Test-Path -LiteralPath $reportPath -PathType Leaf
  $report = if ($reportPresent) {
    Get-Content -LiteralPath $reportPath -Raw | ConvertFrom-Json
  } else { $null }
  [pscustomobject]@{
    registry_key = $key
    registry_present = $keyPresent
    report_path = $reportPath
    report_present = $reportPresent
    state = if ($values) { [string]$values.State } elseif ($report) { [string]$report.state } else { '' }
    remaining_count = if ($values) { [int]$values.RemainingCount } elseif ($report) { [int]$report.remaining } else { 0 }
    report = $report
  }
}

function Assert-GroupPreserved {
  param(
    [Parameter(Mandatory)] [string]$ExpectedSid,
    [Parameter(Mandatory)] [string]$Label
  )

  $group = Get-LocalGroup -Name $adminGroupName -ErrorAction SilentlyContinue
  if (-not $group -or $group.SID.Value -ne $ExpectedSid) {
    throw "$Label did not preserve the $adminGroupName group identity"
  }
}

function Remove-ProbePath {
  param([Parameter(Mandatory)] [string]$Path)

  if (-not (Test-Path -LiteralPath $Path)) { return }
  $full = [IO.Path]::GetFullPath($Path)
  $allowed = @(
    [IO.Path]::GetFullPath($machineDataRoot),
    [IO.Path]::GetFullPath($profileDataRoot),
    [IO.Path]::GetFullPath($machineSentinel),
    [IO.Path]::GetFullPath($profileSentinel)
  )
  if ($full -notin $allowed) {
    throw "Refusing cleanup outside exact S039 probe paths: $full"
  }
  Remove-Item -LiteralPath $full -Recurse -Force
}

if (-not $IsWindows) { throw 'S039 installer contract requires Windows.' }
if ($env:GITHUB_ACTIONS -ne 'true' -or
    $env:RUNNER_ENVIRONMENT -ne 'github-hosted' -or
    $env:RUNNER_OS -ne 'Windows' -or
    -not $env:RUNNER_TEMP) {
  throw 'Refusing destructive installer contract outside a GitHub-hosted disposable Windows runner.'
}
$principal = [Security.Principal.WindowsPrincipal]::new(
  [Security.Principal.WindowsIdentity]::GetCurrent()
)
if (-not $principal.IsInRole(
    [Security.Principal.WindowsBuiltInRole]::Administrator
  )) {
  throw 'S039 installer contract requires an elevated disposable runner.'
}
if (-not $PSCmdlet.ShouldProcess(
    $env:RUNNER_NAME,
    'Exercise machine-wide S039 MSI install, maintenance, upgrade, and removal'
  )) {
  Write-Output 'installer-contract-ci: SKIPPED - lifecycle action was declined'
  return
}

$resolvedBaseline = (Resolve-Path -LiteralPath $BaselineMsiPath).Path
$resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
$resolvedEvidence = [IO.Path]::GetFullPath($EvidencePath)
$runnerRoot = [IO.Path]::GetFullPath($env:RUNNER_TEMP).TrimEnd('\') + '\'
if (-not $resolvedEvidence.StartsWith(
    $runnerRoot,
    [StringComparison]::OrdinalIgnoreCase
  )) {
  throw 'Evidence path must remain under RUNNER_TEMP.'
}
$evidenceDirectory = Split-Path $resolvedEvidence -Parent
$script:LogDirectory = Join-Path $evidenceDirectory 'logs'
New-Item -ItemType Directory -Path $script:LogDirectory -Force | Out-Null

$candidateHash = (Get-FileHash -LiteralPath $resolvedMsi -Algorithm SHA256).Hash.ToLowerInvariant()
$candidateVersion = Get-MsiProperty -PackagePath $resolvedMsi -PropertyName 'ProductVersion'
$candidateProductCode = Get-MsiProperty -PackagePath $resolvedMsi -PropertyName 'ProductCode'
$baselineHash = (Get-FileHash -LiteralPath $resolvedBaseline -Algorithm SHA256).Hash.ToLowerInvariant()
$baselineVersion = Get-MsiProperty -PackagePath $resolvedBaseline -PropertyName 'ProductVersion'
$baselineProductCode = Get-MsiProperty -PackagePath $resolvedBaseline -PropertyName 'ProductCode'
$operatingSystem = Get-CimInstance Win32_OperatingSystem
$identity = [Security.Principal.WindowsIdentity]::GetCurrent().Name
$groupSid = ''
$lockedStream = $null

try {
  if ((Test-ProductInstalled -ProductCode $candidateProductCode) -or
      (Test-ProductInstalled -ProductCode $baselineProductCode)) {
    throw 'go-schedule test product is already installed; runner is not clean.'
  }
  if (Get-LocalGroup -Name $adminGroupName -ErrorAction SilentlyContinue) {
    throw "$adminGroupName already exists; runner is not clean."
  }
  foreach ($path in @(
    $machineDataRoot,
    $profileDataRoot,
    $machineSentinel,
    $profileSentinel
  )) {
    if (Test-Path -LiteralPath $path) {
      throw "S039 probe path already exists; runner is not clean: $path"
    }
  }

  # Fresh-install defaults.
  Invoke-MsiOperation -Label 'default-install' -Operation '/i' -Package $resolvedMsi
  Assert-ProductState -Label 'default install' -ProductCode $candidateProductCode `
    -Installed $true
  Assert-ShortcutState -Label 'default install' -StartMenu $true -Desktop $false
  $group = Get-LocalGroup -Name $adminGroupName
  $groupSid = $group.SID.Value
  Invoke-MsiOperation -Label 'default-preserve-uninstall' -Operation '/x' -Package $resolvedMsi
  Assert-ProductState -Label 'default uninstall' -ProductCode $candidateProductCode `
    -Installed $false
  Assert-GroupPreserved -ExpectedSid $groupSid -Label 'default uninstall'

  # Both shortcuts, followed by maintenance removal and restoration.
  Invoke-MsiOperation -Label 'both-install' -Operation '/i' -Package $resolvedMsi `
    -Properties @('ADDLOCAL=ALL')
  Assert-ShortcutState -Label 'both install' -StartMenu $true -Desktop $true
  $coreHash = (Get-FileHash -LiteralPath (Join-Path $installDirectory 'gosched.exe') -Algorithm SHA256).Hash
  Invoke-MsiOperation -Label 'maintenance-remove-start-menu' -Operation '/i' `
    -Package $resolvedMsi -Properties @('REMOVE=StartMenuShortcut')
  Assert-ShortcutState -Label 'maintenance removed Start Menu' -StartMenu $false -Desktop $true
  Invoke-MsiOperation -Label 'maintenance-add-start-menu' -Operation '/i' `
    -Package $resolvedMsi -Properties @('ADDLOCAL=StartMenuShortcut')
  Assert-ShortcutState -Label 'maintenance restored Start Menu' -StartMenu $true -Desktop $true
  if ((Get-FileHash -LiteralPath (Join-Path $installDirectory 'gosched.exe') -Algorithm SHA256).Hash -ne $coreHash) {
    throw 'Shortcut maintenance changed the core CLI binary.'
  }
  Invoke-MsiOperation -Label 'both-preserve-uninstall' -Operation '/x' -Package $resolvedMsi
  Assert-ProductState -Label 'both uninstall' -ProductCode $candidateProductCode `
    -Installed $false

  # Invalid public removal value fails before execution.
  Invoke-MsiOperation -Label 'invalid-remove-data' -Operation '/i' -Package $resolvedMsi `
    -Properties @('GOSCHEDULE_REMOVE_DATA=2') -ExpectFailure
  if (Test-ProductInstalled -ProductCode $candidateProductCode) {
    throw 'Invalid GOSCHEDULE_REMOVE_DATA installed the product.'
  }

  # Repair never consumes an otherwise valid wipe opt-in.
  Invoke-MsiOperation -Label 'repair-control-install' -Operation '/i' -Package $resolvedMsi
  $repairSeed = Write-SeedData -Label 'repair-control'
  Invoke-MsiOperation -Label 'repair-with-wipe-property' -Operation '/fa' `
    -Package $resolvedMsi -Properties @('GOSCHEDULE_REMOVE_DATA=1')
  Assert-SeedPreserved -Seed $repairSeed -Label 'repair control'
  Invoke-MsiOperation -Label 'repair-control-uninstall' -Operation '/x' -Package $resolvedMsi
  Assert-ProductState -Label 'repair control uninstall' `
    -ProductCode $candidateProductCode -Installed $false
  Assert-SeedPreserved -Seed $repairSeed -Label 'repair control preserve uninstall'

  # A same-authoring baseline proves stable optional-feature migration. The
  # wipe property is deliberately present to prove upgrade exclusion.
  Invoke-MsiOperation -Label 'upgrade-baseline-install' -Operation '/i' `
    -Package $resolvedBaseline `
    -Properties @('ADDLOCAL=Main,DesktopShortcut','REMOVE=StartMenuShortcut')
  Assert-ShortcutState -Label 'upgrade baseline' -StartMenu $false -Desktop $true
  $upgradeSeed = Write-SeedData -Label 'upgrade-control'
  Invoke-MsiOperation -Label 'candidate-upgrade' -Operation '/i' -Package $resolvedMsi `
    -Properties @('GOSCHEDULE_REMOVE_DATA=1')
  Assert-ShortcutState -Label 'candidate upgrade' -StartMenu $false -Desktop $true
  Assert-SeedPreserved -Seed $upgradeSeed -Label 'upgrade control'
  Invoke-MsiOperation -Label 'upgrade-control-uninstall' -Operation '/x' -Package $resolvedMsi
  Assert-ProductState -Label 'upgrade control uninstall' `
    -ProductCode $candidateProductCode -Installed $false
  Assert-SeedPreserved -Seed $upgradeSeed -Label 'upgrade control preserve uninstall'

  # Explicit wipe with neither shortcut.
  Invoke-MsiOperation -Label 'neither-install' -Operation '/i' -Package $resolvedMsi `
    -Properties @('ADDLOCAL=Main','REMOVE=StartMenuShortcut,DesktopShortcut')
  Assert-ShortcutState -Label 'neither install' -StartMenu $false -Desktop $false
  $wipeSeed = Write-SeedData -Label 'wipe-complete'
  Invoke-MsiOperation -Label 'explicit-wipe-uninstall' -Operation '/x' -Package $resolvedMsi `
    -Properties @('GOSCHEDULE_REMOVE_DATA=1')
  Assert-ProductState -Label 'explicit wipe uninstall' `
    -ProductCode $candidateProductCode -Installed $false
  Assert-WipeState -Seed $wipeSeed -Label 'explicit wipe' -ExpectComplete $true
  Assert-GroupPreserved -ExpectedSid $groupSid -Label 'explicit wipe'
  $completeResult = Get-CleanupResultEvidence
  if ($completeResult.registry_present -or $completeResult.report_present) {
    throw 'Complete wipe retained stale cleanup result evidence.'
  }
  $observations.Add([pscustomobject]@{
    kind = 'cleanup-result'
    label = 'explicit wipe'
    result = $completeResult
  })

  # Desktop-only explicit wipe with a locked owned file must retain truthful
  # partial-result evidence while software removal itself completes.
  Invoke-MsiOperation -Label 'desktop-only-install' -Operation '/i' -Package $resolvedMsi `
    -Properties @('ADDLOCAL=Main,DesktopShortcut','REMOVE=StartMenuShortcut')
  Assert-ShortcutState -Label 'desktop-only install' -StartMenu $false -Desktop $true
  $lockedSeed = Write-SeedData -Label 'wipe-locked'
  $lockedStream = [IO.File]::Open(
    $lockedSeed.machine_file,
    [IO.FileMode]::Open,
    [IO.FileAccess]::Read,
    [IO.FileShare]::None
  )
  Invoke-MsiOperation -Label 'locked-wipe-uninstall' -Operation '/x' -Package $resolvedMsi `
    -Properties @('GOSCHEDULE_REMOVE_DATA=1')
  Assert-ProductState -Label 'locked wipe uninstall' `
    -ProductCode $candidateProductCode -Installed $false
  Assert-WipeState -Seed $lockedSeed -Label 'locked wipe' -ExpectComplete $false
  $lockedResult = Get-CleanupResultEvidence
  $observations.Add([pscustomobject]@{
    kind = 'cleanup-result'
    label = 'locked wipe'
    result = $lockedResult
  })
  if (-not $lockedResult.registry_present -or
      -not $lockedResult.report_present -or
      $lockedResult.state -notin @('partial','internal-error') -or
      $lockedResult.remaining_count -lt 1) {
    throw 'Locked wipe did not retain a non-success cleanup result.'
  }
} catch {
  $problems.Add($_.Exception.Message)
} finally {
  if ($lockedStream) { $lockedStream.Dispose() }
  foreach ($productCode in @($candidateProductCode, $baselineProductCode)) {
    if ($productCode -and (Test-ProductInstalled -ProductCode $productCode)) {
      try {
        Invoke-MsiOperation -Label "cleanup-$($productCode.Trim('{}'))" `
          -Operation '/x' -Package $productCode
      } catch {
        $problems.Add("Cleanup uninstall failed for ${productCode}: $($_.Exception.Message)")
      }
    }
  }
  foreach ($path in @(
    $machineDataRoot,
    $profileDataRoot,
    $machineSentinel,
    $profileSentinel
  )) {
    try { Remove-ProbePath -Path $path } catch { $problems.Add($_.Exception.Message) }
  }
  if (Get-LocalGroup -Name $adminGroupName -ErrorAction SilentlyContinue) {
    try {
      Remove-LocalGroup -Name $adminGroupName
    } catch {
      $problems.Add("Local-group cleanup failed: $($_.Exception.Message)")
    }
  }
}

$evidence = [ordered]@{
  schema_version = 1
  evidence_class = 'hosted Windows Server silent installer contract'
  attended_windows_11_gate = 'not claimed; tracked by #94'
  artifact_origin = $ArtifactOrigin
  candidate_path = $resolvedMsi
  candidate_sha256 = $candidateHash
  candidate_version = $candidateVersion
  candidate_product_code = $candidateProductCode
  baseline_path = $resolvedBaseline
  baseline_sha256 = $baselineHash
  baseline_version = $baselineVersion
  baseline_product_code = $baselineProductCode
  windows_caption = $operatingSystem.Caption
  windows_version = $operatingSystem.Version
  windows_product_type = $operatingSystem.ProductType
  identity = $identity
  elevated = $true
  group_sid = $groupSid
  operations = @($operations)
  observations = @($observations)
  status = if ($problems.Count) { 'failed' } else { 'proven' }
  problems = @($problems)
}
$evidence | ConvertTo-Json -Depth 8 |
  Set-Content -LiteralPath $resolvedEvidence -Encoding utf8NoBOM

if ($problems.Count) {
  [Console]::Error.WriteLine(
    "installer-contract-ci: FAILED`n - " + ($problems -join "`n - ")
  )
  exit 1
}
Write-Output "installer-contract-ci: OK - evidence written to $resolvedEvidence"
