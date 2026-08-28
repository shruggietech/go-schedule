<#
.SYNOPSIS
  Exercises a go-schedule MSI on a clean disposable Windows host.

.DESCRIPTION
  Refuses contaminated baselines, then records install, fresh-environment CLI
  probes, reinstall, PATH cardinality, and uninstall. The package is removed in
  a finally block after installation begins. Run only on a disposable machine.
#>
param(
  [Parameter(Mandatory)]
  [string]$MsiPath,

  [Parameter(Mandatory)]
  [string]$EvidencePath,

  [Parameter(Mandatory)]
  [ValidateSet('candidate', 'published')]
  [string]$ArtifactClass,

  [Parameter(Mandatory)]
  [ValidateNotNullOrEmpty()]
  [string]$ArtifactOrigin,

  [switch]$PauseForNativeObservation
)

$ErrorActionPreference = 'Stop'
if (-not $IsWindows) { throw 'Installer lifecycle verification requires Windows.' }

$ArtifactOrigin = $ArtifactOrigin.Trim()
if (-not $ArtifactOrigin) { throw 'Artifact origin must not be blank.' }
$resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
$resolvedEvidence = [System.IO.Path]::GetFullPath($EvidencePath)
$installDir = 'C:\Program Files\go-schedule'
$observations = [System.Collections.Generic.List[string]]::new()
$failures = [System.Collections.Generic.List[string]]::new()
$installed = $false
$installAttempted = $false
$lifecycleStatus = 'failed'
$lifecycleError = $null
$nativeObservations = [ordered]@{
  'Start Menu' = [pscustomobject]@{ Status = 'unavailable'; Detail = 'not recorded by the operator' }
  'Installed apps' = [pscustomobject]@{ Status = 'unavailable'; Detail = 'not recorded by the operator' }
  'GUI title area' = [pscustomobject]@{ Status = 'unavailable'; Detail = 'not recorded by the operator' }
  'Taskbar' = [pscustomobject]@{ Status = 'unavailable'; Detail = 'not recorded by the operator' }
}

if ($ArtifactClass -eq 'published') {
  [Uri]$originUri = $null
  if (-not [Uri]::TryCreate($ArtifactOrigin, [UriKind]::Absolute, [ref]$originUri) -or
      $originUri.Scheme -ne 'https' -or
      $originUri.Host -ne 'github.com' -or
      $originUri.AbsolutePath -notmatch '^/shruggietech/go-schedule/releases/download/[^/]+/[^/]+\.msi$') {
    throw "Published artifact origin must be this repository's absolute HTTPS release-asset URL."
  }
}

function Get-MachinePath {
  [Environment]::GetEnvironmentVariable('Path', 'Machine')
}

function Get-InstallPathEntry {
  param([string]$PathValue)

  $canonical = $installDir.TrimEnd('\').ToLowerInvariant()
  @($PathValue -split ';' | Where-Object {
    $_.Trim().TrimEnd('\').ToLowerInvariant() -eq $canonical
  })
}

function Get-InstalledProduct {
  Get-ItemProperty `
    'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*' `
    -ErrorAction SilentlyContinue |
    Where-Object DisplayName -eq 'go-schedule' |
    Select-Object -First 1
}

function Read-NativeObservation {
  param([Parameter(Mandatory)] [string]$Surface)

  do {
    $status = (Read-Host "$Surface status (proven/failed/unavailable)").Trim().ToLowerInvariant()
  } while ($status -notin 'proven', 'failed', 'unavailable')

  $detail = (Read-Host "$Surface note or screenshot reference (optional)").Trim()
  if (-not $detail) { $detail = 'no note supplied' }
  [pscustomobject]@{ Status = $status; Detail = $detail }
}

function Write-Evidence {
  param(
    [Parameter(Mandatory)] [string]$Status,
    [string[]]$Problems = @()
  )

  $hash = (Get-FileHash -LiteralPath $resolvedMsi -Algorithm SHA256).Hash.ToLowerInvariant()
  $content = @(
    '# Windows Installer Lifecycle Evidence'
    ''
    "- Date: $(Get-Date -Format 'yyyy-MM-dd')"
    "- Windows: $([Environment]::OSVersion.VersionString)"
    "- Artifact: ``$resolvedMsi``"
    "- Evidence class: **$ArtifactClass artifact**"
    "- Artifact origin: $ArtifactOrigin"
    "- SHA-256: ``$hash``"
    "- Clean lifecycle status: **$Status**"
    ''
    '## Observations'
  )
  if ($observations.Count -eq 0) {
    $content += '- None; prerequisites stopped the run before mutation.'
  } else {
    $content += $observations | ForEach-Object { "- $_" }
  }
  if ($Problems.Count -gt 0) {
    $content += ''
    $content += '## Failures or unavailable prerequisites'
    $content += $Problems | ForEach-Object { "- $_" }
  }
  $content += ''
  $content += '## Native observations'
  foreach ($surface in $nativeObservations.Keys) {
    $entry = $nativeObservations[$surface]
    $content += "- ${surface}: **$($entry.Status)**; $($entry.Detail)."
  }

  $content | Set-Content -LiteralPath $resolvedEvidence -Encoding utf8NoBOM
}

function Invoke-Msi {
  param(
    [Parameter(Mandatory)] [ValidateSet('/i', '/x')] [string]$Operation,
    [Parameter(Mandatory)] [string]$LogPath
  )

  $arguments = "$Operation `"$resolvedMsi`" /qn /norestart /L*v `"$LogPath`""
  $process = Start-Process -FilePath msiexec.exe -ArgumentList $arguments -Wait -PassThru
  if ($process.ExitCode -notin 0, 3010) {
    throw "msiexec $Operation failed with exit code $($process.ExitCode); see $LogPath"
  }
}

function Invoke-FreshPathProbe {
  param([Parameter(Mandatory)] [string[]]$Arguments)

  $savedPath = $env:Path
  try {
    $machinePath = Get-MachinePath
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    $env:Path = @($machinePath, $userPath) -join ';'
    $output = & gosched @Arguments 2>&1
    if ($LASTEXITCODE -ne 0) {
      throw "gosched $($Arguments -join ' ') failed: $($output -join [Environment]::NewLine)"
    }
    ($output | Out-String).Trim()
  } finally {
    $env:Path = $savedPath
  }
}

$principal = [Security.Principal.WindowsPrincipal]::new(
  [Security.Principal.WindowsIdentity]::GetCurrent()
)
if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
  $failures.Add('PowerShell is not elevated; per-machine MSI verification requires Administrator.')
}
if (Get-Command gosched -ErrorAction SilentlyContinue) {
  $failures.Add('gosched already resolves on PATH; baseline is not clean.')
}
if (Test-Path -LiteralPath $installDir) {
  $failures.Add("install directory already exists: $installDir")
}
if (Get-InstalledProduct) {
  $failures.Add('go-schedule already has an installed-product registry entry.')
}
$beforePath = Get-MachinePath
if ((Get-InstallPathEntry $beforePath).Count -ne 0) {
  $failures.Add('machine PATH already contains the go-schedule install directory.')
}

if ($failures.Count -gt 0) {
  Write-Evidence -Status 'unavailable' -Problems $failures
  [Console]::Error.WriteLine(
    "installer-lifecycle: REFUSED before mutation`n - " + ($failures -join "`n - ")
  )
  exit 2
}

$logDir = Join-Path ([System.IO.Path]::GetTempPath()) 'go-schedule-installer-lifecycle'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null

try {
  $installAttempted = $true
  Invoke-Msi -Operation /i -LogPath (Join-Path $logDir 'install.log')
  $installed = $true
  $observations.Add('Install completed successfully.')

  $installedPath = Get-MachinePath
  $pathMatches = Get-InstallPathEntry $installedPath
  if ($pathMatches.Count -ne 1) {
    throw "machine PATH has $($pathMatches.Count) install-directory entries after install; expected 1"
  }
  $observations.Add('Machine PATH contains exactly one install-directory entry after install.')

  $version = Invoke-FreshPathProbe @('--version')
  $observations.Add("Fresh PATH probe: gosched --version succeeded ($version).")
  Invoke-FreshPathProbe @('task', 'list') | Out-Null
  $observations.Add('Fresh PATH probe: gosched task list succeeded.')

  if ($PauseForNativeObservation) {
    Write-Output 'Inspect Start Menu, Installed apps, GUI title area, and taskbar now.'
    foreach ($surface in @($nativeObservations.Keys)) {
      $nativeObservations[$surface] = Read-NativeObservation -Surface $surface
    }
  }

  Invoke-Msi -Operation /i -LogPath (Join-Path $logDir 'reinstall.log')
  $reinstallMatches = Get-InstallPathEntry (Get-MachinePath)
  if ($reinstallMatches.Count -ne 1) {
    throw "machine PATH has $($reinstallMatches.Count) entries after reinstall; expected 1"
  }
  $observations.Add('Reinstall completed with exactly one machine PATH entry.')

  Invoke-Msi -Operation /x -LogPath (Join-Path $logDir 'uninstall.log')
  $installed = $false
  $afterPath = Get-MachinePath
  if ((Get-InstallPathEntry $afterPath).Count -ne 0) {
    throw 'machine PATH still contains the install directory after uninstall'
  }
  if ($afterPath -ne $beforePath) {
    throw 'machine PATH differs from its exact pre-install value after uninstall'
  }
  if ($afterPath -match ';;') {
    throw 'machine PATH contains an empty segment after uninstall'
  }
  $observations.Add('Uninstall restored the exact pre-install machine PATH with no empty segment.')
  $lifecycleStatus = 'proven'
} catch {
  $lifecycleError = $_.Exception.Message
  $failures.Add($lifecycleError)
} finally {
  if ($installAttempted -and
      ($installed -or ($lifecycleStatus -eq 'failed' -and [bool](Get-InstalledProduct)))) {
    try {
      Invoke-Msi -Operation /x -LogPath (Join-Path $logDir 'cleanup-uninstall.log')
      $installed = $false
      $observations.Add('Cleanup uninstall completed after an interrupted lifecycle run.')
    } catch {
      $failures.Add("Cleanup uninstall failed; the package may remain installed: $($_.Exception.Message)")
    }
  }

  if ($installAttempted) {
    $finalProductPresent = [bool](Get-InstalledProduct)
    $finalInstallDirPresent = Test-Path -LiteralPath $installDir
    $finalPath = Get-MachinePath
    $finalPathEntries = (Get-InstallPathEntry $finalPath).Count
    $finalPathRestored = $finalPath -eq $beforePath
    $observations.Add(
      "Final machine state: product registered=$finalProductPresent; " +
      "install directory present=$finalInstallDirPresent; PATH entries=$finalPathEntries; " +
      "original machine PATH restored=$finalPathRestored."
    )
    if ($finalProductPresent) {
      $failures.Add('Final machine state still has an installed-product registry entry.')
    }
    if ($finalInstallDirPresent) {
      $failures.Add('Final machine state still has the installation directory.')
    }
    if ($finalPathEntries -ne 0) {
      $failures.Add("Final machine state has $finalPathEntries install-directory PATH entries; expected 0.")
    }
    if (-not $finalPathRestored) {
      $failures.Add('Final machine PATH does not match its exact pre-install value.')
    }
  }
}

if ($failures.Count -gt 0) {
  $lifecycleStatus = 'failed'
  if (-not $lifecycleError) { $lifecycleError = $failures[0] }
}
Write-Evidence -Status $lifecycleStatus -Problems $failures
if ($lifecycleStatus -eq 'failed') {
  [Console]::Error.WriteLine("installer-lifecycle: FAILED - $lifecycleError")
  exit 1
}
Write-Output "installer-lifecycle: OK - evidence written to $resolvedEvidence"
