<#
.SYNOPSIS
    Exercise and record the security-sensitive lifecycle of a go-schedule MSI.

.DESCRIPTION
    Runs one of three disposable-host scenarios: a clean candidate lifecycle,
    an upgrade from the published v0.9.0 MSI, or a non-elevated access probe
    after login-session refresh. Fresh and upgrade scenarios install, inspect,
    and uninstall software and therefore require an elevated PowerShell 7
    session on a disposable Windows 11 host.

    Every Windows Installer process is noninteractive, hidden, and writes a
    verbose log beside the evidence file. Evidence includes process exit codes,
    log hashes and diagnostics, local-group SID and members, service state,
    product state, install-directory state, and machine PATH cardinality.

    Exit codes: 0 proven, 1 an assertion failed, 2 a prerequisite prevented the
    scenario from starting. A declined confirmation or -WhatIf stops at the
    first skipped phase with exit code 0 and no final evidence or cleanup.
    Fresh and upgrade runs intentionally preserve the goschedadmin group and
    membership after uninstall.

.PARAMETER Scenario
    Scenario to run: fresh, upgrade, or access-probe.
    Alias: s

.PARAMETER MsiPath
    Candidate MSI path. For access-probe, this identifies the installed
    candidate artifact for evidence provenance.
    Alias: m

.PARAMETER EvidencePath
    New Markdown evidence file to write. Existing files are refused.
    Alias: e

.PARAMETER ArtifactClass
    Candidate artifact classification: candidate or published.
    Alias: c

.PARAMETER ArtifactOrigin
    Human-readable candidate origin, or this repository's HTTPS release URL for
    a published artifact.
    Alias: o

.PARAMETER PriorMsiPath
    Published v0.9.0 MSI path. Required only for the upgrade scenario.
    Alias: p

.PARAMETER PriorArtifactOrigin
    Absolute v0.9.0 GitHub release-asset URL. Required for upgrade.
    Alias: r

.PARAMETER Quiet
    Suppress informational progress while retaining warnings and errors.
    Alias: q

.PARAMETER Silent
    Suppress all progress output. Genuine errors still reach the error stream.

.PARAMETER Help
    Print detailed help.
    Alias: h

.EXAMPLE
    .\Invoke-InstallerLifecycle.ps1 -Scenario fresh -MsiPath C:\verify\candidate.msi -EvidencePath C:\verify\fresh.md -ArtifactClass candidate -ArtifactOrigin 'local build from commit abc' -Confirm:$false
    Run install, repair, reinstall, and uninstall on a clean disposable host.

.EXAMPLE
    .\Invoke-InstallerLifecycle.ps1 -Scenario upgrade -MsiPath C:\verify\candidate.msi -PriorMsiPath C:\verify\v0.9.0.msi -EvidencePath C:\verify\upgrade.md -ArtifactClass candidate -ArtifactOrigin 'local build from commit abc' -PriorArtifactOrigin 'https://github.com/shruggietech/go-schedule/releases/download/v0.9.0/go-schedule_v0.9.0_windows_amd64.msi' -Confirm:$false
    Preprovision the durable group, install v0.9.0, and upgrade the candidate.

.EXAMPLE
    .\Invoke-InstallerLifecycle.ps1 -Scenario access-probe -MsiPath C:\verify\candidate.msi -EvidencePath C:\verify\access.md -ArtifactClass candidate -ArtifactOrigin 'local build from commit abc'
    After sign-out/sign-in, prove a non-elevated token reaches the daemon.
#>
[CmdletBinding(SupportsShouldProcess=$true,ConfirmImpact='High',DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias("s")]
    [ValidateSet('fresh','upgrade','access-probe')]
    [string]$Scenario,

    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias("m")]
    [string]$MsiPath,

    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias("e")]
    [string]$EvidencePath,

    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias("c")]
    [ValidateSet('candidate', 'published')]
    [string]$ArtifactClass,

    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias("o")]
    [ValidateNotNullOrEmpty()]
    [string]$ArtifactOrigin,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("p")]
    [string]$PriorMsiPath = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("r")]
    [string]$PriorArtifactOrigin = '',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias("q")]
    [Switch]$Quiet,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Switch]$Silent,

    [Parameter(Mandatory=$true,ParameterSetName='HelpText')]
    [Alias("h")]
    [Switch]$Help
)
#_______________________________________________________________________________
## Declare Functions

    function Assert-PSVersion {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$false)]
            [version]$Minimum = '7.0'
        )
        $current = $PSVersionTable.PSVersion
        if ($current -lt $Minimum) {
            Write-Host ("ALERT: PowerShell {0}+ required; running {1}. Relaunch with 'pwsh'." -f $Minimum, $current) -ForegroundColor Red
            exit 2
        }
    }

    function Write-Log {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true,Position=0)]
            [string]$Message,

            [Parameter(Mandatory=$false)]
            [ValidateSet('Info','Success','Warn','Error','Debug')]
            [string]$Level = 'Info',

            [Parameter(Mandatory=$false)]
            [string]$Source = $null
        )
        if ($script:LogSilent -and $Level -ne 'Error') { return }
        if ($script:LogQuiet -and (@('Info','Success','Debug') -contains $Level)) { return }

        $stamp = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss.fff')
        $tag   = if ($Source) { "[$Source] " } else { '' }
        $label = $Level.ToUpper().PadRight(7)
        $color = switch ($Level) {
            'Info'    { 'Gray' }
            'Success' { 'Green' }
            'Warn'    { 'Yellow' }
            'Error'   { 'Red' }
            'Debug'   { 'DarkGray' }
        }
        Write-Host ("{0} {1}{2} {3}" -f $stamp, $tag, $label, $Message) -ForegroundColor $color
    }

    function Get-MachinePath {
        [CmdletBinding()]
        Param()
        [Environment]::GetEnvironmentVariable('Path', 'Machine')
    }

    function Get-InstallPathCount {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [AllowEmptyString()]
            [string]$PathValue
        )
        $canonical = $script:InstallDir.TrimEnd('\').ToLowerInvariant()
        @($PathValue -split ';' | Where-Object {
            $_.Trim().TrimEnd('\').ToLowerInvariant() -eq $canonical
        }).Count
    }

    function Get-InstalledProduct {
        [CmdletBinding()]
        Param()
        # Wildcard expansion is intentional across uninstall registry subkeys.
        Get-ItemProperty -Path $script:UninstallRegistryPath `
            -ErrorAction SilentlyContinue |
            Where-Object DisplayName -eq 'go-schedule' |
            Select-Object -First 1
    }

    function Get-MsiProperty {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$PackagePath,

            [Parameter(Mandatory=$true)]
            [string]$PropertyName
        )
        $installer = New-Object -ComObject WindowsInstaller.Installer
        $database = $installer.OpenDatabase($PackagePath, 0)
        $query = "SELECT ``Value`` FROM ``Property`` WHERE ``Property``='$PropertyName'"
        $view = $database.OpenView($query)
        $view.Execute() | Out-Null
        $record = $view.Fetch()
        if (-not $record) { return '' }
        $record.StringData(1)
    }

    function Test-ReleaseOrigin {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Origin,

            [Parameter(Mandatory=$false)]
            [string]$ExactUrl = ''
        )
        [Uri]$uri = $null
        if (-not [Uri]::TryCreate(
            $Origin,
            [UriKind]::Absolute,
            [ref]$uri
        )) { return $false }
        if ($uri.Scheme -ne 'https' -or $uri.Host -ne 'github.com' -or `
            $uri.AbsolutePath -notmatch `
                '^/shruggietech/go-schedule/releases/download/[^/]+/[^/]+\.msi$') {
            return $false
        }
        -not $ExactUrl -or $uri.AbsoluteUri -eq $ExactUrl
    }

    function Get-SecurityState {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Label
        )
        $group = Get-LocalGroup -Name $script:AdminGroupName `
            -ErrorAction SilentlyContinue
        $members = @()
        if ($group) {
            $members = @(Get-LocalGroupMember -Group $script:AdminGroupName |
                ForEach-Object { $_.Name } | Sort-Object -Unique)
        }
        $service = Get-CimInstance -ClassName Win32_Service `
            -Filter "Name='goschedd'" -ErrorAction SilentlyContinue
        [pscustomobject]@{
            Label = $Label
            ProductPresent = [bool](Get-InstalledProduct)
            InstallDirectoryPresent = Test-Path -LiteralPath $script:InstallDir
            PathEntries = Get-InstallPathCount -PathValue (Get-MachinePath)
            GroupPresent = [bool]$group
            GroupSid = if ($group) { $group.SID.Value } else { '' }
            Members = $members
            IntendedMemberPresent = $members -contains $script:CurrentIdentity
            ServicePresent = [bool]$service
            ServiceState = if ($service) { $service.State } else { '' }
            ServiceStartMode = if ($service) { $service.StartMode } else { '' }
            ServiceAccount = if ($service) { $service.StartName } else { '' }
            ServicePath = if ($service) { $service.PathName } else { '' }
        }
    }

    function Get-LogDiagnostics {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$LogPath,

            [Parameter(Mandatory=$true)]
            [bool]$RequireGroupOrder
        )
        $text = Get-Content -LiteralPath $LogPath -Raw
        $groupMatches = [regex]::Matches(
            $text,
            '(?im)^.*(?:Wix6CreateGroup|CreateGroup).*$'
        )
        $serviceMatches = [regex]::Matches(
            $text,
            '(?im)^.*Action start.*StartServices.*$'
        )
        $membershipMatches = [regex]::Matches(
            $text,
            '(?im)^.*(?:Wix(?:4|6)ConfigureUsers|ConfigureUsers|AddUserToGroup).*$'
        )
        $groupIndex = if ($groupMatches.Count) {
            $groupMatches[$groupMatches.Count - 1].Index
        } else { -1 }
        $serviceIndex = if ($serviceMatches.Count) {
            $serviceMatches[0].Index
        } else { -1 }
        $membershipIndex = if ($membershipMatches.Count) {
            $membershipMatches[$membershipMatches.Count - 1].Index
        } else { -1 }
        [pscustomobject]@{
            AccessDenied = $text -match '(?i)0x80070005|E_ACCESSDENIED'
            Error26421 = $text -match '(?m)\b26421\b'
            Rollback = $text -match '(?im)^.*Return value 3.*$'
            GroupActionFound = $groupIndex -ge 0
            StartServicesFound = $serviceIndex -ge 0
            MembershipActionFound = $membershipIndex -ge 0
            GroupBeforeService = if ($RequireGroupOrder) {
                $groupIndex -ge 0 -and $serviceIndex -ge 0 -and `
                    $groupIndex -lt $serviceIndex
            } else { $true }
            MembershipBeforeService = if ($RequireGroupOrder) {
                $membershipIndex -ge 0 -and $serviceIndex -ge 0 -and `
                    $membershipIndex -lt $serviceIndex
            } else { $true }
        }
    }

    function Invoke-MsiOperation {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Label,

            [Parameter(Mandatory=$true)]
            [string]$PackagePath,

            [Parameter(Mandatory=$true)]
            [string[]]$OperationArguments,

            [Parameter(Mandatory=$true)]
            [bool]$RequireGroupOrder
        )
        Initialize-EvidenceDirectory
        $logPath = [System.IO.Path]::Combine(
            $script:LogDirectory,
            "$Label.log"
        )
        if ($OperationArguments.Count -lt 1) {
            throw "Windows Installer '$Label' has no operation switch."
        }
        $arguments = @(
            $OperationArguments[0],
            ('"{0}"' -f $PackagePath)
        )
        if ($OperationArguments.Count -gt 1) {
            $arguments += @($OperationArguments | Select-Object -Skip 1)
        }
        $arguments += @(
            '/qn',
            '/norestart',
            '/L*v',
            ('"{0}"' -f $logPath)
        )
        Write-Log "Running Windows Installer operation '$Label'." `
            -Level Info -Source 'msiexec'
        $process = Start-Process -FilePath 'msiexec.exe' `
            -ArgumentList $arguments -Wait -PassThru -WindowStyle Hidden
        if (-not (Test-Path -LiteralPath $logPath -PathType Leaf)) {
            throw "Windows Installer did not create verbose log: $logPath"
        }
        $diagnostics = Get-LogDiagnostics -LogPath $logPath `
            -RequireGroupOrder $RequireGroupOrder
        $result = [pscustomobject]@{
            Label = $Label
            ExitCode = $process.ExitCode
            LogPath = $logPath
            LogSha256 = (Get-FileHash -LiteralPath $logPath `
                -Algorithm SHA256).Hash.ToLowerInvariant()
            Diagnostics = $diagnostics
        }
        $script:Operations.Add($result)
        if ($process.ExitCode -notin 0, 3010) {
            throw "Windows Installer '$Label' failed with exit code $($process.ExitCode)"
        }
        if ($diagnostics.AccessDenied -or $diagnostics.Error26421 -or `
            $diagnostics.Rollback -or `
            -not $diagnostics.GroupBeforeService -or `
            -not $diagnostics.MembershipBeforeService) {
            throw "Windows Installer '$Label' verbose-log contract failed"
        }
        $result
    }

    function Initialize-EvidenceDirectory {
        [CmdletBinding()]
        Param()
        $evidenceDirectory = [System.IO.Path]::GetDirectoryName(
            $script:ResolvedEvidence
        )
        if (-not [System.IO.Directory]::Exists($evidenceDirectory)) {
            [void][System.IO.Directory]::CreateDirectory($evidenceDirectory)
        }
        if (-not [System.IO.Directory]::Exists($script:LogDirectory)) {
            [void][System.IO.Directory]::CreateDirectory($script:LogDirectory)
        }
    }

    function Assert-InstalledState {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [pscustomobject]$State
        )
        $expectedServicePath = [System.IO.Path]::Combine(
            $script:InstallDir,
            'goschedd.exe'
        )
        $actualServicePath = $State.ServicePath.Trim().Trim('"')
        if (-not $State.ProductPresent -or `
            -not $State.InstallDirectoryPresent -or `
            $State.PathEntries -ne 1 -or `
            -not $State.GroupPresent -or `
            -not $State.IntendedMemberPresent -or `
            -not $State.ServicePresent -or `
            $State.ServiceState -ne 'Running' -or `
            $State.ServiceStartMode -ne 'Auto' -or `
            $State.ServiceAccount -ne 'LocalSystem' -or `
            -not $actualServicePath.Equals(
                $expectedServicePath,
                [StringComparison]::OrdinalIgnoreCase
            )) {
            throw "Installed-state contract failed after '$($State.Label)'"
        }
    }

    function Write-Evidence {
        [CmdletBinding()]
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Status,

            [Parameter(Mandatory=$false)]
            [string[]]$Problems = @()
        )
        Initialize-EvidenceDirectory
        $lines = @(
            '# Windows Installer Lifecycle Evidence'
            ''
            "- Date: $(Get-Date -Format 'yyyy-MM-dd')"
            "- Windows: $([Environment]::OSVersion.VersionString)"
            "- Scenario: **$Scenario**"
            "- Candidate: ``$script:ResolvedMsi``"
            "- Evidence class: **$ArtifactClass artifact**"
            "- Artifact origin: $ArtifactOrigin"
            "- Candidate SHA-256: ``$script:CandidateHash``"
            "- Candidate product version: ``$script:CandidateVersion``"
            "- Candidate product code: ``$script:CandidateProductCode``"
            "- Intended identity: ``$script:CurrentIdentity``"
            "- Lifecycle status: **$Status**"
        )
        if ($script:ResolvedPriorMsi) {
            $lines += "- Prior artifact: ``$script:ResolvedPriorMsi``"
            $lines += "- Prior artifact origin: $PriorArtifactOrigin"
            $lines += "- Prior SHA-256: ``$script:PriorHash``"
        }
        if ($Scenario -eq 'access-probe') {
            $lines += @(
                ''
                '## Refreshed non-elevated access probe'
                "- Token contains goschedadmin SID: $script:AccessTokenContainsGroup"
                "- Command: ``$script:AccessProbeCommand``"
                "- Exit code: $script:AccessProbeExitCode"
                "- Output: $script:AccessProbeOutput"
            )
        }
        $lines += @('', '## MSI operations')
        if ($script:Operations.Count -eq 0) {
            $lines += '- None.'
        } else {
            foreach ($operation in $script:Operations) {
                $diag = $operation.Diagnostics
                $lines += (
                    "- $($operation.Label): exit=$($operation.ExitCode); " +
                    "log=``$($operation.LogPath)``; " +
                    "sha256=``$($operation.LogSha256)``; " +
                    "access-denied=$($diag.AccessDenied); " +
                    "error-26421=$($diag.Error26421); rollback=$($diag.Rollback); " +
                    "group-action=$($diag.GroupActionFound); " +
                    "membership-action=$($diag.MembershipActionFound); " +
                    "start-services=$($diag.StartServicesFound); " +
                    "group-before-service=$($diag.GroupBeforeService); " +
                    "membership-before-service=$($diag.MembershipBeforeService)."
                )
            }
        }
        $lines += @('', '## Phase observations')
        if ($script:States.Count -eq 0) {
            $lines += '- None.'
        } else {
            foreach ($state in $script:States) {
                $memberList = if ($state.Members.Count) {
                    $state.Members -join ', '
                } else { '<none>' }
                $lines += (
                    "- $($state.Label): product=$($state.ProductPresent); " +
                    "directory=$($state.InstallDirectoryPresent); " +
                    "PATH entries=$($state.PathEntries); " +
                    "group=$($state.GroupPresent); SID=``$($state.GroupSid)``; " +
                    "intended member=$($state.IntendedMemberPresent); " +
                    "members=``$memberList``; service=$($state.ServicePresent); " +
                    "state=``$($state.ServiceState)``; " +
                    "start=``$($state.ServiceStartMode)``; " +
                    "account=``$($state.ServiceAccount)``; " +
                    "path=``$($state.ServicePath)``."
                )
            }
        }
        if ($Problems.Count) {
            $lines += @('', '## Failures or unavailable prerequisites')
            $lines += $Problems | ForEach-Object { "- $_" }
        }
        $lines | Set-Content -LiteralPath $script:ResolvedEvidence `
            -Encoding utf8
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $script:LogQuiet  = $false
    $script:LogSilent = $false
    $ThisScriptPath = $MyInvocation.MyCommand.Path
    $script:AdminGroupName = 'goschedadmin'
    $script:InstallDir = 'C:\Program Files\go-schedule'
    $script:UninstallRegistryPath = `
        'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*'
    $script:Operations = [System.Collections.Generic.List[object]]::new()
    $script:States = [System.Collections.Generic.List[object]]::new()
    $script:ResolvedMsi = ''
    $script:ResolvedPriorMsi = ''
    $script:ResolvedEvidence = ''
    $script:LogDirectory = ''
    $script:CandidateHash = ''
    $script:PriorHash = ''
    $script:CleanupPackage = ''
    $script:CandidateVersion = ''
    $script:CandidateProductCode = ''
    $script:ExpectedPriorUrl = `
        'https://github.com/shruggietech/go-schedule/releases/download/v0.9.0/go-schedule_v0.9.0_windows_amd64.msi'
    $script:ExpectedPriorHash = `
        '0365f8ea592321100ffb2875d4c66649d2ab53407faddba5a1168a4ae1e1fb1c'
    $script:AccessTokenContainsGroup = $false
    $script:AccessProbeCommand = ''
    $script:AccessProbeExitCode = -1
    $script:AccessProbeOutput = '<not run>'
    $script:CurrentIdentity = `
        [Security.Principal.WindowsIdentity]::GetCurrent().Name

#_______________________________________________________________________________
## Execute Operations

    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }

    Assert-PSVersion
    if ($Quiet)  { $script:LogQuiet  = $true }
    if ($Silent) { $script:LogSilent = $true }
    $ErrorActionPreference = 'Stop'

    if (-not $IsWindows) {
        Write-Host 'FAIL: installer lifecycle verification requires Windows.' `
            -ForegroundColor Red
        exit 2
    }
    $script:ResolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
    $script:ResolvedEvidence = [System.IO.Path]::GetFullPath($EvidencePath)
    if (Test-Path -LiteralPath $script:ResolvedEvidence) {
        Write-Host "FAIL: evidence file already exists: $script:ResolvedEvidence" `
            -ForegroundColor Red
        exit 2
    }
    $script:CandidateHash = (Get-FileHash -LiteralPath $script:ResolvedMsi `
        -Algorithm SHA256).Hash.ToLowerInvariant()
    $script:CandidateVersion = Get-MsiProperty `
        -PackagePath $script:ResolvedMsi -PropertyName 'ProductVersion'
    $script:CandidateProductCode = Get-MsiProperty `
        -PackagePath $script:ResolvedMsi -PropertyName 'ProductCode'
    $evidenceDirectory = [System.IO.Path]::GetDirectoryName(
        $script:ResolvedEvidence
    )
    $evidenceStem = [System.IO.Path]::GetFileNameWithoutExtension(
        $script:ResolvedEvidence
    )
    $script:LogDirectory = [System.IO.Path]::Combine(
        $evidenceDirectory,
        "$evidenceStem.logs"
    )
    $principal = [Security.Principal.WindowsPrincipal]::new(
        [Security.Principal.WindowsIdentity]::GetCurrent()
    )
    $isAdministrator = $principal.IsInRole(
        [Security.Principal.WindowsBuiltInRole]::Administrator
    )

    $hostProblems = [System.Collections.Generic.List[string]]::new()
    $operatingSystem = Get-CimInstance -ClassName Win32_OperatingSystem
    if ($operatingSystem.ProductType -ne 1 -or `
        $operatingSystem.Caption -notmatch 'Windows 11') {
        $hostProblems.Add(
            "Unsupported host: $($operatingSystem.Caption); Windows 11 client required."
        )
    }
    if ($ArtifactClass -eq 'published' -and `
        -not (Test-ReleaseOrigin -Origin $ArtifactOrigin)) {
        $hostProblems.Add(
            'Published candidate origin is not this repository release asset URL.'
        )
    }
    if ($hostProblems.Count) {
        Write-Evidence -Status 'unavailable' -Problems $hostProblems
        Write-Host ("FAIL: lifecycle host or provenance unavailable.`n - " +
            ($hostProblems -join "`n - ")) -ForegroundColor Red
        exit 2
    }

    if ($Scenario -eq 'access-probe') {
        $problems = [System.Collections.Generic.List[string]]::new()
        if ($isAdministrator) {
            $problems.Add('Access probe must run from a non-elevated session.')
        }
        $state = Get-SecurityState -Label 'post-refresh access probe'
        $script:States.Add($state)
        try {
            Assert-InstalledState -State $state
        } catch {
            $problems.Add($_.Exception.Message)
        }
        $product = Get-InstalledProduct
        if (-not $product -or `
            [string]$product.DisplayVersion -ne $script:CandidateVersion -or `
            [string]$product.PSChildName -ne $script:CandidateProductCode) {
            $problems.Add(
                'Installed product identity does not match the candidate MSI.'
            )
        }
        $group = Get-LocalGroup -Name $script:AdminGroupName `
            -ErrorAction SilentlyContinue
        $tokenSids = @([Security.Principal.WindowsIdentity]::GetCurrent().Groups |
            ForEach-Object { $_.Value })
        $script:AccessTokenContainsGroup = [bool]$group -and `
            $tokenSids -contains $group.SID.Value
        if (-not $script:AccessTokenContainsGroup) {
            $problems.Add('Current token does not contain the goschedadmin SID.')
        }
        $configPath = [System.IO.Path]::Combine(
            $env:ProgramData,
            'goschedule',
            'config.json'
        )
        $configuredGroup = $script:AdminGroupName
        if (Test-Path -LiteralPath $configPath -PathType Leaf) {
            try {
                $config = Get-Content -LiteralPath $configPath -Raw |
                    ConvertFrom-Json
                if ($null -ne $config.admin_group) {
                    $configuredGroup = [string]$config.admin_group
                }
            } catch {
                $problems.Add("Could not read daemon config: $($_.Exception.Message)")
            }
        }
        if ($configuredGroup -ne $script:AdminGroupName) {
            $problems.Add(
                "Daemon admin_group is '$configuredGroup'; expected goschedadmin."
            )
        }
        $installedCli = [System.IO.Path]::Combine(
            $script:InstallDir,
            'gosched.exe'
        )
        try {
            if (-not (Test-Path -LiteralPath $installedCli -PathType Leaf)) {
                throw "Installed candidate CLI is missing: $installedCli"
            }
            $script:AccessProbeCommand = "$installedCli health"
            $output = & $installedCli health 2>&1
            $script:AccessProbeExitCode = $LASTEXITCODE
            $script:AccessProbeOutput = (($output | Out-String).Trim() `
                -replace '[\r\n]+', ' ') -replace '`', "'"
            if ($LASTEXITCODE -ne 0) {
                $problems.Add("gosched health failed: $($output -join ' ')")
            }
        } catch {
            $problems.Add("gosched health could not run: $($_.Exception.Message)")
        }
        $status = if ($problems.Count) { 'failed' } else { 'proven' }
        Write-Evidence -Status $status -Problems $problems
        if ($problems.Count) { exit 1 }
        Write-Log 'Non-elevated access probe proven.' -Level Success `
            -Source 'access'
        exit 0
    }

    $preconditions = [System.Collections.Generic.List[string]]::new()
    if (-not $isAdministrator) {
        $preconditions.Add('PowerShell is not elevated.')
    }
    if (Get-InstalledProduct) {
        $preconditions.Add('go-schedule is already registered as installed.')
    }
    if (Test-Path -LiteralPath $script:InstallDir) {
        $preconditions.Add("install directory exists: $script:InstallDir")
    }
    if ((Get-InstallPathCount -PathValue (Get-MachinePath)) -ne 0) {
        $preconditions.Add('machine PATH already contains the install directory.')
    }
    if (Get-Service -Name 'goschedd' -ErrorAction SilentlyContinue) {
        $preconditions.Add('goschedd service already exists.')
    }
    if (Get-LocalGroup -Name $script:AdminGroupName `
        -ErrorAction SilentlyContinue) {
        $preconditions.Add('goschedadmin already exists; reset the disposable host.')
    }
    if ($Scenario -eq 'upgrade') {
        if (-not $PriorMsiPath -or -not $PriorArtifactOrigin) {
            $preconditions.Add(
                'Upgrade requires PriorMsiPath and PriorArtifactOrigin.'
            )
        } else {
            $script:ResolvedPriorMsi = `
                (Resolve-Path -LiteralPath $PriorMsiPath).Path
            $script:PriorHash = (Get-FileHash `
                -LiteralPath $script:ResolvedPriorMsi `
                -Algorithm SHA256).Hash.ToLowerInvariant()
            if (-not (Test-ReleaseOrigin -Origin $PriorArtifactOrigin `
                -ExactUrl $script:ExpectedPriorUrl)) {
                $preconditions.Add(
                    'Prior origin is not the pinned v0.9.0 release asset URL.'
                )
            }
            if ($script:PriorHash -ne $script:ExpectedPriorHash) {
                $preconditions.Add(
                    'Prior MSI SHA-256 does not match published v0.9.0.'
                )
            }
        }
    }
    if ($preconditions.Count) {
        Write-Evidence -Status 'unavailable' -Problems $preconditions
        Write-Host ("FAIL: lifecycle prerequisites unavailable.`n - " +
            ($preconditions -join "`n - ")) -ForegroundColor Red
        exit 2
    }

    $installed = $false
    $actionSkipped = $false
    $problems = [System.Collections.Generic.List[string]]::new()
    try {
        if ($Scenario -eq 'upgrade') {
            if (-not $PSCmdlet.ShouldProcess(
                $script:AdminGroupName,
                'Preprovision v0.9.0 local group and membership'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            New-LocalGroup -Name $script:AdminGroupName | Out-Null
            Add-LocalGroupMember -Group $script:AdminGroupName `
                -Member $script:CurrentIdentity
            if (-not $PSCmdlet.ShouldProcess(
                $script:ResolvedPriorMsi,
                'Install published v0.9.0 baseline'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            $script:CleanupPackage = $script:ResolvedPriorMsi
            Invoke-MsiOperation -Label 'prior-install' `
                -PackagePath $script:ResolvedPriorMsi `
                -OperationArguments @('/i') -RequireGroupOrder $false |
                Out-Null
            $installed = $true
            $priorState = Get-SecurityState -Label 'v0.9.0 installed'
            $script:States.Add($priorState)
            Assert-InstalledState -State $priorState
            $baselineSid = $priorState.GroupSid
            $baselineMembers = $priorState.Members -join '|'
            if (-not $PSCmdlet.ShouldProcess(
                $script:ResolvedMsi,
                'Upgrade v0.9.0 to candidate'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            $script:CleanupPackage = $script:ResolvedMsi
            Invoke-MsiOperation -Label 'candidate-upgrade' `
                -PackagePath $script:ResolvedMsi `
                -OperationArguments @('/i') -RequireGroupOrder $true |
                Out-Null
            $candidateState = Get-SecurityState -Label 'candidate upgraded'
            $script:States.Add($candidateState)
            Assert-InstalledState -State $candidateState
            if ($candidateState.GroupSid -ne $baselineSid -or `
                ($candidateState.Members -join '|') -ne $baselineMembers) {
                throw 'Upgrade changed the administrative group SID or members.'
            }
        } else {
            if (-not $PSCmdlet.ShouldProcess(
                $script:ResolvedMsi,
                'Install candidate MSI'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            $script:CleanupPackage = $script:ResolvedMsi
            Invoke-MsiOperation -Label 'candidate-install' `
                -PackagePath $script:ResolvedMsi `
                -OperationArguments @('/i') -RequireGroupOrder $true |
                Out-Null
            $installed = $true
            $installState = Get-SecurityState -Label 'candidate installed'
            $script:States.Add($installState)
            Assert-InstalledState -State $installState
            $baselineSid = $installState.GroupSid
            $baselineMembers = $installState.Members -join '|'
            if (-not $PSCmdlet.ShouldProcess(
                $script:ResolvedMsi,
                'Repair candidate MSI'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            Invoke-MsiOperation -Label 'candidate-repair' `
                -PackagePath $script:ResolvedMsi `
                -OperationArguments @('/fa') -RequireGroupOrder $false |
                Out-Null
            $repairState = Get-SecurityState -Label 'candidate repaired'
            $script:States.Add($repairState)
            Assert-InstalledState -State $repairState
            if (-not $PSCmdlet.ShouldProcess(
                $script:ResolvedMsi,
                'Reinstall candidate MSI'
            )) {
                $actionSkipped = $true
                Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                    -Level Warn -Source 'lifecycle'
                return
            }
            Invoke-MsiOperation -Label 'candidate-reinstall' `
                -PackagePath $script:ResolvedMsi `
                -OperationArguments @(
                    '/i','REINSTALL=ALL','REINSTALLMODE=vomus'
                ) -RequireGroupOrder $false | Out-Null
            $reinstallState = Get-SecurityState -Label 'candidate reinstalled'
            $script:States.Add($reinstallState)
            Assert-InstalledState -State $reinstallState
            foreach ($state in @($repairState, $reinstallState)) {
                if ($state.GroupSid -ne $baselineSid -or `
                    ($state.Members -join '|') -ne $baselineMembers) {
                    throw "$($state.Label) changed the group SID or members."
                }
            }
        }

        if (-not $PSCmdlet.ShouldProcess(
            $script:ResolvedMsi,
            'Uninstall candidate MSI'
        )) {
            $actionSkipped = $true
            Write-Log 'Lifecycle action skipped; stopping without evidence.' `
                -Level Warn -Source 'lifecycle'
            return
        }
        Invoke-MsiOperation -Label 'candidate-uninstall' `
            -PackagePath $script:ResolvedMsi `
            -OperationArguments @('/x') -RequireGroupOrder $false |
            Out-Null
        $installed = $false
        $finalState = Get-SecurityState -Label 'candidate uninstalled'
        $script:States.Add($finalState)
        if ($finalState.ProductPresent -or `
            $finalState.InstallDirectoryPresent -or `
            $finalState.PathEntries -ne 0 -or `
            -not $finalState.GroupPresent -or `
            -not $finalState.IntendedMemberPresent -or `
            $finalState.ServicePresent -or `
            $finalState.GroupSid -ne $baselineSid -or `
            ($finalState.Members -join '|') -ne $baselineMembers) {
            throw 'Final uninstall state violates cleanup or preservation.'
        }
    } catch {
        $problems.Add($_.Exception.Message)
    } finally {
        if (-not $actionSkipped -and `
            ($installed -or [bool](Get-InstalledProduct))) {
            try {
                $registeredProduct = Get-InstalledProduct
                $cleanupTarget = if ($registeredProduct) {
                    [string]$registeredProduct.PSChildName
                } else { $script:CleanupPackage }
                if ($PSCmdlet.ShouldProcess(
                    $cleanupTarget,
                    'Cleanup uninstall after interrupted lifecycle'
                )) {
                    Invoke-MsiOperation -Label 'cleanup-uninstall' `
                        -PackagePath $cleanupTarget `
                        -OperationArguments @('/x') `
                        -RequireGroupOrder $false | Out-Null
                }
            } catch {
                $problems.Add(
                    "Cleanup uninstall failed: $($_.Exception.Message)"
                )
            }
            $script:States.Add((Get-SecurityState -Label 'final cleanup state'))
        }
    }

    $status = if ($problems.Count) { 'failed' } else { 'proven' }
    Write-Evidence -Status $status -Problems $problems
    if ($problems.Count) {
        Write-Host ("FAIL: installer lifecycle failed.`n - " +
            ($problems -join "`n - ")) -ForegroundColor Red
        exit 1
    }
    Write-Log "Lifecycle proven; evidence: $script:ResolvedEvidence" `
        -Level Success -Source 'lifecycle'
    exit 0

#_______________________________________________________________________________
## End of script
