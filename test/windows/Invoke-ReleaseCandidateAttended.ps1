<#
.SYNOPSIS
    Collect and finalize attended Windows release-candidate evidence.

.DESCRIPTION
    Creates a resumable evidence workspace for one exact go-schedule MSI,
    captures native HWND and Fyne canvas measurements, imports explicit
    operator-reviewed observation fragments, and produces an evidence archive
    only after the shared release gate accepts every required scenario.

    The script never installs, removes, or publishes software. It refuses to
    overwrite workspaces, observations, evidence attachments, and final ZIPs.
    Console child processes run hidden with redirected noninteractive I/O.

.PARAMETER Action
    Operation to perform: Initialize, CaptureWindow, RecordObservation, or
    Finalize.
    Alias: a

.PARAMETER MsiPath
    Absolute path to the exact staged Windows MSI.
    Alias: m

.PARAMETER WorkspacePath
    Absolute local fixed-volume evidence workspace.
    Alias: w

.PARAMETER Repository
    GitHub owner/repository recorded during Initialize.
    Default: 'shruggietech/go-schedule'.
    Alias: r

.PARAMETER Tag
    Candidate tag in vMAJOR.MINOR.PATCH form. Required for Initialize.
    Alias: t

.PARAMETER Commit
    Forty-character lowercase tag commit. Required for Initialize.
    Alias: c

.PARAMETER RunId
    GitHub Actions staging run identifier. Required for Initialize.
    Alias: i

.PARAMETER RunAttempt
    GitHub Actions staging run attempt.
    Default: 1.
    Alias: n

.PARAMETER ProcessId
    Exact installed GUI process ID used by CaptureWindow.
    Alias: p

.PARAMETER ObservationId
    Fixed window scenario identifier used by CaptureWindow.
    Alias: o

.PARAMETER EnvironmentPath
    JSON environment record used by CaptureWindow.
    Alias: e

.PARAMETER FyneEvidencePath
    Opt-in JSON evidence emitted by the installed GUI process.
    Alias: f

.PARAMETER ScreenshotPath
    Existing screenshot beneath the workspace attachments directory.
    Alias: s

.PARAMETER ObservationPath
    Operator-reviewed fragment imported by RecordObservation.
    Alias: j

.PARAMETER OperatorRole
    Non-personal operator role written by Finalize.
    Default: 'release maintainer'.
    Alias: q

.PARAMETER Help
    Print detailed help.
    Alias: h

.EXAMPLE
    .\Invoke-ReleaseCandidateAttended.ps1 -Action Initialize `
        -MsiPath C:\Evidence\go-schedule_v1.0.0_windows_amd64.msi `
        -WorkspacePath C:\Evidence\v1.0.0 -Tag v1.0.0 `
        -Commit 0123456789abcdef0123456789abcdef01234567 -RunId 1234

.EXAMPLE
    .\Invoke-ReleaseCandidateAttended.ps1 -Action Finalize `
        -MsiPath C:\Evidence\go-schedule_v1.0.0_windows_amd64.msi `
        -WorkspacePath C:\Evidence\v1.0.0

.NOTES
    Attended observations require a clean Windows 11 client and genuine user
    interaction. Fixture or hosted-server output cannot replace that proof.
#>
[CmdletBinding(SupportsShouldProcess=$true,ConfirmImpact='Medium',
    DefaultParameterSetName='Default')]
Param(
    [Parameter(Mandatory=$true,ParameterSetName='Default')]
    [Alias('a')]
    [ValidateSet('Initialize','CaptureWindow','RecordObservation','Finalize')]
    [string]$Action,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('m')]
    [string]$MsiPath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('w')]
    [string]$WorkspacePath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('r')]
    [string]$Repository = 'shruggietech/go-schedule',

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('t')]
    [string]$Tag,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('c')]
    [string]$Commit,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('i')]
    [long]$RunId,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('n')]
    [int]$RunAttempt = 1,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('p')]
    [int]$ProcessId,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('o')]
    [ValidateSet(
        'window.clean-standard',
        'window.clean-high-or-mixed',
        'window.retained-profile',
        'window.state-transitions',
        'window.subsequent-launch'
    )]
    [string]$ObservationId,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('e')]
    [string]$EnvironmentPath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('f')]
    [string]$FyneEvidencePath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('s')]
    [string]$ScreenshotPath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('j')]
    [string]$ObservationPath,

    [Parameter(Mandatory=$false,ParameterSetName='Default')]
    [Alias('q')]
    [string]$OperatorRole = 'release maintainer',

    [Parameter(Mandatory=$true,ParameterSetName='HelpText')]
    [Alias('h')]
    [Switch]$Help
)
#_______________________________________________________________________________
## Declare Functions

    function Assert-AbsolutePath {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Name,
            [Parameter(Mandatory=$true)]
            [string]$Value
        )

        if (-not [System.IO.Path]::IsPathFullyQualified($Value)) {
            throw "$Name must be an absolute path."
        }
    }

    function Get-MsiProperty {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Path,
            [Parameter(Mandatory=$true)]
            [string]$Property
        )

        $installer = New-Object -ComObject WindowsInstaller.Installer
        $database = $installer.OpenDatabase($Path, 0)
        $query = 'SELECT `Value` FROM `Property` WHERE `Property`=?'
        $view = $database.OpenView($query)
        $record = $installer.CreateRecord(1)
        $record.StringData(1) = $Property
        $view.Execute($record) | Out-Null
        $row = $view.Fetch()
        if (-not $row) {
            throw "MSI property '$Property' was not found."
        }
        return $row.StringData(1)
    }

    function Write-JsonNoBom {
        Param(
            [Parameter(Mandatory=$true)]
            [object]$Value,
            [Parameter(Mandatory=$true)]
            [string]$Path
        )

        $json = $Value | ConvertTo-Json -Depth 20
        $stream = [System.IO.File]::Open(
            $Path,
            [System.IO.FileMode]::CreateNew,
            [System.IO.FileAccess]::Write,
            [System.IO.FileShare]::None
        )
        $completed = $false
        try {
            $writer = [System.IO.StreamWriter]::new(
                $stream,
                [System.Text.UTF8Encoding]::new($false)
            )
            try {
                $writer.Write("$json`n")
                $writer.Flush()
                $completed = $true
            } finally {
                $writer.Dispose()
            }
        } finally {
            $stream.Dispose()
            if (-not $completed -and (Test-Path -LiteralPath $Path)) {
                [System.IO.File]::Delete($Path)
            }
        }
    }

    function Save-JsonReplacingFile {
        Param(
            [Parameter(Mandatory=$true)]
            [object]$Value,
            [Parameter(Mandatory=$true)]
            [string]$Path
        )

        $temporary = "$Path.tmp"
        if (Test-Path -LiteralPath $temporary) {
            throw "Stale temporary evidence file exists: '$temporary'."
        }
        $json = $Value | ConvertTo-Json -Depth 20
        [System.IO.File]::WriteAllText(
            $temporary,
            "$json`n",
            [System.Text.UTF8Encoding]::new($false)
        )
        Move-Item -LiteralPath $temporary -Destination $Path -Force
    }

    function Invoke-HiddenProcess {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$FilePath,
            [Parameter(Mandatory=$true)]
            [string[]]$Arguments,
            [Parameter(Mandatory=$true)]
            [string]$WorkingDirectory
        )

        $info = [System.Diagnostics.ProcessStartInfo]::new()
        $info.FileName = $FilePath
        $info.WorkingDirectory = $WorkingDirectory
        $info.UseShellExecute = $false
        $info.CreateNoWindow = $true
        $info.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $info.RedirectStandardOutput = $true
        $info.RedirectStandardError = $true
        $info.RedirectStandardInput = $true
        foreach ($argument in $Arguments) {
            [void]$info.ArgumentList.Add($argument)
        }
        $process = [System.Diagnostics.Process]::new()
        $process.StartInfo = $info
        if (-not $process.Start()) {
            throw "Failed to start '$FilePath'."
        }
        $process.StandardInput.Close()
        $stdoutTask = $process.StandardOutput.ReadToEndAsync()
        $stderrTask = $process.StandardError.ReadToEndAsync()
        $process.WaitForExit()
        $stdout = $stdoutTask.GetAwaiter().GetResult()
        $stderr = $stderrTask.GetAwaiter().GetResult()
        return [pscustomobject]@{
            ExitCode = $process.ExitCode
            Stdout = $stdout
            Stderr = $stderr
        }
    }

    function Add-NativeWindowType {
        if ('GoSchedule.ReleaseEvidence.NativeWindow' -as [type]) {
            return
        }
        Add-Type -TypeDefinition @'
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Runtime.InteropServices;
using System.Security.Principal;

namespace GoSchedule.ReleaseEvidence {
    public sealed class NativeWindowResult {
        public long Hwnd { get; set; }
        public Rect Outer { get; set; }
        public Rect Client { get; set; }
        public Rect Monitor { get; set; }
        public Rect WorkArea { get; set; }
        public string MonitorId { get; set; }
        public uint Dpi { get; set; }
        public uint ShowCommand { get; set; }
        public bool Visible { get; set; }
        public bool Maximized { get; set; }
        public bool Minimized { get; set; }
        public bool Fullscreen { get; set; }
        public bool Restored { get; set; }
        public string ProcessUserSid { get; set; }
        public int IntegrityRid { get; set; }
    }

    public struct Rect {
        public int Left;
        public int Top;
        public int Right;
        public int Bottom;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal struct Point {
        public int X;
        public int Y;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal struct WindowPlacement {
        public int Length;
        public int Flags;
        public uint ShowCommand;
        public Point MinPosition;
        public Point MaxPosition;
        public Rect NormalPosition;
    }

    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Auto)]
    internal struct MonitorInfo {
        public int Size;
        public Rect Monitor;
        public Rect WorkArea;
        public uint Flags;
        [MarshalAs(UnmanagedType.ByValTStr, SizeConst = 32)]
        public string Device;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal struct SidAndAttributes {
        public IntPtr Sid;
        public uint Attributes;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal struct TokenMandatoryLabel {
        public SidAndAttributes Label;
    }

    public static class NativeWindow {
        private delegate bool EnumWindowsCallback(IntPtr hwnd, IntPtr state);

        [DllImport("user32.dll")]
        private static extern bool EnumWindows(
            EnumWindowsCallback callback,
            IntPtr state
        );

        [DllImport("user32.dll")]
        private static extern uint GetWindowThreadProcessId(
            IntPtr hwnd,
            out uint processId
        );

        [DllImport("user32.dll")]
        private static extern bool IsWindowVisible(IntPtr hwnd);

        [DllImport("user32.dll")]
        private static extern bool GetWindowRect(IntPtr hwnd, out Rect rect);

        [DllImport("user32.dll")]
        private static extern bool GetClientRect(IntPtr hwnd, out Rect rect);

        [DllImport("user32.dll")]
        private static extern bool ClientToScreen(IntPtr hwnd, ref Point point);

        [DllImport("user32.dll")]
        private static extern bool GetWindowPlacement(
            IntPtr hwnd,
            ref WindowPlacement placement
        );

        [DllImport("user32.dll")]
        private static extern bool IsZoomed(IntPtr hwnd);

        [DllImport("user32.dll")]
        private static extern bool IsIconic(IntPtr hwnd);

        [DllImport("user32.dll")]
        private static extern IntPtr MonitorFromWindow(
            IntPtr hwnd,
            uint flags
        );

        [DllImport("user32.dll", CharSet = CharSet.Auto)]
        private static extern bool GetMonitorInfo(
            IntPtr monitor,
            ref MonitorInfo info
        );

        [DllImport("user32.dll")]
        private static extern uint GetDpiForWindow(IntPtr hwnd);

        [DllImport("user32.dll")]
        private static extern IntPtr SetThreadDpiAwarenessContext(
            IntPtr context
        );

        [DllImport("kernel32.dll", SetLastError = true)]
        private static extern IntPtr OpenProcess(
            uint desiredAccess,
            bool inheritHandle,
            uint processId
        );

        [DllImport("advapi32.dll", SetLastError = true)]
        private static extern bool OpenProcessToken(
            IntPtr processHandle,
            uint desiredAccess,
            out IntPtr tokenHandle
        );

        [DllImport("advapi32.dll", SetLastError = true)]
        private static extern bool GetTokenInformation(
            IntPtr tokenHandle,
            int tokenInformationClass,
            IntPtr tokenInformation,
            int tokenInformationLength,
            out int returnLength
        );

        [DllImport("advapi32.dll")]
        private static extern IntPtr GetSidSubAuthorityCount(IntPtr sid);

        [DllImport("advapi32.dll")]
        private static extern IntPtr GetSidSubAuthority(
            IntPtr sid,
            uint subAuthority
        );

        [DllImport("kernel32.dll", SetLastError = true)]
        private static extern bool CloseHandle(IntPtr handle);

        public static NativeWindowResult Capture(uint expectedProcessId) {
            IntPtr previous = SetThreadDpiAwarenessContext(new IntPtr(-4));
            try {
                return CaptureAware(expectedProcessId);
            } finally {
                if (previous != IntPtr.Zero) {
                    SetThreadDpiAwarenessContext(previous);
                }
            }
        }

        private static NativeWindowResult CaptureAware(
            uint expectedProcessId
        ) {
            var matches = new List<IntPtr>();
            EnumWindows((hwnd, state) => {
                uint processId;
                GetWindowThreadProcessId(hwnd, out processId);
                if (processId == expectedProcessId && IsWindowVisible(hwnd)) {
                    matches.Add(hwnd);
                }
                return true;
            }, IntPtr.Zero);
            if (matches.Count != 1) {
                throw new InvalidOperationException(
                    "Expected exactly one visible top-level window for PID " +
                    expectedProcessId + "; found " + matches.Count + "."
                );
            }
            IntPtr hwnd = matches[0];
            Rect outer;
            Rect client;
            if (!GetWindowRect(hwnd, out outer) ||
                !GetClientRect(hwnd, out client)) {
                throw new InvalidOperationException("Could not read rectangles.");
            }
            var origin = new Point();
            if (!ClientToScreen(hwnd, ref origin)) {
                throw new InvalidOperationException("Could not map client area.");
            }
            client.Right += origin.X;
            client.Bottom += origin.Y;
            client.Left = origin.X;
            client.Top = origin.Y;
            var placement = new WindowPlacement();
            placement.Length = Marshal.SizeOf<WindowPlacement>();
            if (!GetWindowPlacement(hwnd, ref placement)) {
                throw new InvalidOperationException("Could not read placement.");
            }
            var info = new MonitorInfo();
            info.Size = Marshal.SizeOf<MonitorInfo>();
            IntPtr monitor = MonitorFromWindow(hwnd, 2);
            if (monitor == IntPtr.Zero || !GetMonitorInfo(monitor, ref info)) {
                throw new InvalidOperationException("Could not read monitor.");
            }
            bool maximized = IsZoomed(hwnd);
            bool minimized = IsIconic(hwnd);
            bool fullscreen = !maximized && !minimized &&
                outer.Left <= info.Monitor.Left &&
                outer.Top <= info.Monitor.Top &&
                outer.Right >= info.Monitor.Right &&
                outer.Bottom >= info.Monitor.Bottom;
            string processUserSid;
            int integrityRid;
            ReadToken(expectedProcessId, out processUserSid, out integrityRid);
            return new NativeWindowResult {
                Hwnd = hwnd.ToInt64(),
                Outer = outer,
                Client = client,
                Monitor = info.Monitor,
                WorkArea = info.WorkArea,
                MonitorId = info.Device,
                Dpi = GetDpiForWindow(hwnd),
                ShowCommand = placement.ShowCommand,
                Visible = true,
                Maximized = maximized,
                Minimized = minimized,
                Fullscreen = fullscreen,
                Restored = !maximized && !minimized && !fullscreen,
                ProcessUserSid = processUserSid,
                IntegrityRid = integrityRid
            };
        }

        private static void ReadToken(
            uint processId,
            out string userSid,
            out int integrityRid
        ) {
            const uint ProcessQueryLimitedInformation = 0x1000;
            const uint TokenQuery = 0x0008;
            IntPtr process = OpenProcess(
                ProcessQueryLimitedInformation,
                false,
                processId
            );
            if (process == IntPtr.Zero) {
                throw new Win32Exception(Marshal.GetLastWin32Error());
            }
            IntPtr token = IntPtr.Zero;
            try {
                if (!OpenProcessToken(process, TokenQuery, out token)) {
                    throw new Win32Exception(Marshal.GetLastWin32Error());
                }
                IntPtr userBuffer = ReadTokenBuffer(token, 1);
                try {
                    var user = Marshal.PtrToStructure<SidAndAttributes>(
                        userBuffer
                    );
                    userSid = new SecurityIdentifier(user.Sid).Value;
                } finally {
                    Marshal.FreeHGlobal(userBuffer);
                }
                IntPtr integrityBuffer = ReadTokenBuffer(token, 25);
                try {
                    var label = Marshal.PtrToStructure<TokenMandatoryLabel>(
                        integrityBuffer
                    );
                    byte count = Marshal.ReadByte(
                        GetSidSubAuthorityCount(label.Label.Sid)
                    );
                    integrityRid = Marshal.ReadInt32(GetSidSubAuthority(
                        label.Label.Sid,
                        (uint)(count - 1)
                    ));
                } finally {
                    Marshal.FreeHGlobal(integrityBuffer);
                }
            } finally {
                if (token != IntPtr.Zero) {
                    CloseHandle(token);
                }
                CloseHandle(process);
            }
        }

        private static IntPtr ReadTokenBuffer(IntPtr token, int kind) {
            int length;
            GetTokenInformation(token, kind, IntPtr.Zero, 0, out length);
            if (length <= 0) {
                throw new Win32Exception(Marshal.GetLastWin32Error());
            }
            IntPtr buffer = Marshal.AllocHGlobal(length);
            if (!GetTokenInformation(token, kind, buffer, length, out length)) {
                int error = Marshal.GetLastWin32Error();
                Marshal.FreeHGlobal(buffer);
                throw new Win32Exception(error);
            }
            return buffer;
        }
    }
}
'@
    }

    function Convert-Rect {
        Param(
            [Parameter(Mandatory=$true)]
            [object]$Rect
        )

        return [ordered]@{
            left = $Rect.Left
            top = $Rect.Top
            right = $Rect.Right
            bottom = $Rect.Bottom
        }
    }

    function Get-RelativeAttachmentPath {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Workspace,
            [Parameter(Mandatory=$true)]
            [string]$Path
        )

        $workspaceFull = [System.IO.Path]::GetFullPath($Workspace)
        $pathFull = [System.IO.Path]::GetFullPath($Path)
        $relative = [System.IO.Path]::GetRelativePath($workspaceFull, $pathFull)
        $relative = $relative.Replace('\', '/')
        if (-not $relative.StartsWith('attachments/')) {
            throw "Evidence attachment must be beneath '$workspaceFull'."
        }
        if ($relative.Contains('../')) {
            throw 'Evidence attachment path traversal is not allowed.'
        }
        $current = if (Test-Path -LiteralPath $pathFull) {
            Get-Item -LiteralPath $pathFull
        } else {
            Get-Item -LiteralPath (Split-Path $pathFull -Parent)
        }
        while ($current.FullName -ne $workspaceFull) {
            if (($current.Attributes -band
                [System.IO.FileAttributes]::ReparsePoint) -ne 0) {
                throw "Reparse attachment path is not allowed: '$pathFull'."
            }
            $current = $current.Parent
            if (-not $current) {
                throw "Evidence attachment escaped workspace: '$pathFull'."
            }
        }
        return $relative
    }

    function New-UnavailableObservation {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Id,
            [Parameter(Mandatory=$true)]
            [datetime]$Timestamp
        )

        return [ordered]@{
            id = $Id
            environment_id = ''
            status = 'unavailable'
            started_at = $Timestamp.ToUniversalTime().ToString('o')
            completed_at = $Timestamp.ToUniversalTime().ToString('o')
            summary = 'not yet observed'
            metrics = [ordered]@{}
            attachment_paths = @()
        }
    }

    function New-ScenarioTemplate {
        Param(
            [Parameter(Mandatory=$true)]
            [string]$Id,
            [Parameter(Mandatory=$true)]
            [datetime]$Timestamp
        )

        $metrics = [ordered]@{}
        if ($Id.StartsWith('setup.')) {
            $metrics = [ordered]@{
                install_session_id = ''
                installer_process_id = 0
                installer_session_id = 0
                installer_process_owner_role = ''
                installer_process_integrity = ''
                selected_options_sha256 = ''
                observed_targets_sha256 = ''
                before_fingerprint = ''
                after_fingerprint = ''
                owned_data_cleanup_invoked = $false
            }
        } elseif ($Id.StartsWith('remove.')) {
            $metrics = [ordered]@{
                owned_roots_count = 0
                before_content_sha256 = ''
                after_content_sha256 = ''
                controls_before_sha256 = ''
                controls_after_sha256 = ''
                security_state_disposition = ''
                reinstall_result = ''
            }
        }
        $templateId = switch ($Id) {
            'desktop.interaction-states-scaled' { 'desktop.interaction-states' }
            'desktop.navigation-options-scaled' { 'desktop.navigation-options' }
            'desktop.tasks-table-scaled' { 'desktop.tasks-table' }
            'desktop.schedule-activity-tables-scaled' { 'desktop.schedule-activity-tables' }
            default { $Id }
        }
        switch ($templateId) {
            'setup.shortcut-defaults' {
                $metrics.start_menu_default = $false
                $metrics.desktop_default = $false
                $metrics.defaults_visible = $false
                $metrics.effects_verified = $false
            }
            'setup.shortcut-matrix' {
                $metrics.combinations_verified = 0
                $metrics.targets_verified = $false
            }
            'setup.completion-matrix' {
                $metrics.combinations_verified = 0
                $metrics.independent_choices = $false
                $metrics.default_handler_verified = $false
                $metrics.launch_default = $false
                $metrics.documentation_default = $true
            }
            'setup.finish-launch-integrity' {
                $metrics.process_integrity = ''
                $metrics.launch_count = 0
            }
            'setup.cancel' {
                $metrics.state_unchanged = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'setup.maintenance' {
                $metrics.transitions_verified = $false
                $metrics.repair_verified = $false
                $metrics.completion_actions_absent = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'setup.upgrade' {
                $metrics.choices_preserved = $false
                $metrics.completion_actions_absent = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'setup.invalid-input' {
                $metrics.input_rejected = $false
                $metrics.state_unchanged = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'setup.rollback' {
                $metrics.rollback_completed = $false
                $metrics.state_unchanged = $false
            }
            'remove.preserve' {
                $metrics.mode = ''
                $metrics.software_removed = $false
                $metrics.owned_bytes_preserved = $false
                $metrics.controls_unchanged = $false
                $metrics.preserve_default_visible = $false
                $metrics.owned_inventory_reviewed = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'remove.wipe' {
                $metrics.mode = ''
                $metrics.software_removed = $false
                $metrics.owned_roots_removed = $false
                $metrics.controls_unchanged = $false
                $metrics.security_state_preserved = $false
                $metrics.wipe_explicitly_selected = $false
                $metrics.wipe_confirmed = $false
                $metrics.owned_inventory_reviewed = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'remove.cancel' {
                $metrics.software_unchanged = $false
                $metrics.data_unchanged = $false
                $metrics.owned_data_cleanup_invoked = $false
            }
            'remove.multiple-profiles' {
                $metrics.profile_count = 0
                $metrics.all_profiles_accounted = $false
            }
            'remove.locked-partial' {
                $metrics.cleanup_result = ''
                $metrics.residual_count = 0
                $metrics.truthfully_reported = $false
            }
            'remove.reinstall-after-preserve' {
                $metrics.prior_tasks_available = $false
                $metrics.prior_preferences_available = $false
                $metrics.prior_config_available = $false
                $metrics.prior_logs_available = $false
            }
            'remove.reinstall-after-wipe' {
                $metrics.prior_tasks_available = $true
                $metrics.prior_preferences_available = $true
                $metrics.prior_config_available = $true
                $metrics.prior_logs_available = $true
            }
            'desktop.appearance-standard' {
                $metrics = [ordered]@{
                    palettes = 'dark,light'
                    effective_dpi = 96
                    system_font_default = $false
                    system_font_restored = $false
                    font_persistence_verified = $false
                    info_text_sharp = $false
                    body_text_sharp = $false
                    labels_centered = $false
                    labels_unclipped = $false
                    resize_verified = $false
                    minimize_restore_verified = $false
                    reopen_verified = $false
                    fonts_exercised = 'system,geist,inter,ubuntu,monospace'
                }
            }
            'desktop.appearance-scaled' {
                $metrics = [ordered]@{
                    palettes = 'dark,light'
                    effective_dpi = 0
                    system_font_default = $false
                    system_font_restored = $false
                    font_persistence_verified = $false
                    info_text_sharp = $false
                    body_text_sharp = $false
                    labels_centered = $false
                    labels_unclipped = $false
                    resize_verified = $false
                    minimize_restore_verified = $false
                    reopen_verified = $false
                    fonts_exercised = 'system,geist,inter,ubuntu,monospace'
                }
            }
            'desktop.interaction-states' {
                $metrics = [ordered]@{
                    palettes = 'dark,light'
                    control_families = 'navigation,selector,ordinary,primary,danger,dialog,table-row'
                    states = 'rest,hover,focus,pressed,selected,disabled'
                    minimum_text_contrast = 0
                    minimum_non_text_contrast = 0
                    labels_readable = $false
                    glyphs_readable = $false
                    selection_identifiable = $false
                    focus_visible = $false
                    non_color_cues_present = $false
                }
            }
            'desktop.navigation-options' {
                $metrics = [ordered]@{
                    palettes = 'dark,light'
                    content_sizes = '1280x800,800x600'
                    destination_order = 'tasks,groups,chains,schedule,activity,options,info'
                    rail_spacing_balanced = $false
                    labels_unclipped = $false
                    boundary_full_height = $false
                    boundary_subtle = $false
                    exit_bottom_right = $false
                    exit_never_selected = $false
                    exit_semantic_glyph = $false
                    storage_rows_compact = $false
                    unavailable_rows_muted = $false
                    copy_exact = $false
                    selector_current_omitted = $false
                    horizontal_scrollbar_present = $true
                }
            }
            'desktop.scroll-input' {
                $metrics = [ordered]@{
                    sensitivities = '1x,2x,4x'
                    surfaces = 'options,info,editor-command,editor-schedule,editor-help'
                    wheel_detents_responsive = $false
                    immediate_apply = $false
                    persistence_verified = $false
                    nested_multiplier_absent = $false
                    keyboard_scroll_preserved = $false
                    touchpad_available = $false
                    touchpad_fine_deltas_preserved = $false
                    touchpad_unavailable_reason = ''
                }
            }
            'desktop.tasks-table' {
                $metrics = [ordered]@{
                    row_count = 0
                    palettes = 'dark,light'
                    content_sizes = '1280x800,800x600'
                    headers = 'task,enabled,lifecycle,time-zone,group'
                    row_states = 'odd,even,hover,focus,selected'
                    headers_frozen = $false
                    status_dimensions_distinct = $false
                    bracket_decoration_absent = $false
                    full_values_discoverable = $false
                    horizontal_scrollbar_present = $true
                    refresh_identity_stable = $false
                    removed_selection_clears = $false
                    toolbar_actions_work = $false
                    double_click_edits = $false
                }
            }
            'desktop.schedule-activity-tables' {
                $metrics = [ordered]@{
                    schedule_row_count = 0
                    activity_row_count = 0
                    palettes = 'dark,light'
                    content_sizes = '1280x800,800x600'
                    schedule_headers = 'when,task,event,outcome'
                    activity_headers = 'when,severity,source,summary'
                    schedule_states = 'scheduled,success,failure,skipped,caught-up,queued,missing,unknown'
                    severities = 'INFO,WARNING,ERROR'
                    row_states = 'odd,even,hover,focus,selected'
                    headers_frozen = $false
                    semantic_text_glyphs_match = $false
                    non_color_cues_present = $false
                    full_values_discoverable = $false
                    horizontal_scrollbar_present = $true
                    refresh_identity_stable = $false
                    removed_selection_clears = $false
                    detail_activation_accurate = $false
                    range_calendar_switching = $false
                    filter_clear_acknowledge = $false
                }
            }
        }
        return [ordered]@{
            instructions = @(
                'Perform this scenario interactively against the exact MSI.',
                'Replace every placeholder with observed values and digests.',
                'Add a screenshot beneath attachments/ and keep status unavailable until reviewed.',
                'Set status to pass only after every claimed result is visible and verified.'
            )
            environment = [ordered]@{
                id = ''
                snapshot = ''
                clean_snapshot = $false
                windows_edition = ''
                windows_build = ''
                account_role = ''
                account_sid = ''
                integrity = ''
                integrity_rid = 0
                service_identity = ''
                display_class = 'not-applicable'
                effective_dpi = 0
                profile_state = 'not-applicable'
            }
            observation = [ordered]@{
                id = $Id
                environment_id = ''
                status = 'unavailable'
                started_at = $Timestamp.ToUniversalTime().ToString('o')
                completed_at = $Timestamp.ToUniversalTime().ToString('o')
                summary = 'not yet observed'
                metrics = $metrics
                attachment_paths = @(
                    'attachments/screenshots/' + $Id + '.png'
                )
            }
        }
    }

    function Invoke-Initialize {
        if (-not $MsiPath -or -not $WorkspacePath -or -not $Tag -or
            -not $Commit -or $RunId -le 0 -or $RunAttempt -le 0) {
            throw 'Initialize requires MSI, workspace, tag, commit, and run ID.'
        }
        Assert-AbsolutePath -Name 'MsiPath' -Value $MsiPath
        Assert-AbsolutePath -Name 'WorkspacePath' -Value $WorkspacePath
        if ($Tag -notmatch '^v[0-9]+\.[0-9]+\.[0-9]+$') {
            throw 'Tag must match vMAJOR.MINOR.PATCH.'
        }
        if ($Commit -notmatch '^[0-9a-f]{40}$') {
            throw 'Commit must be 40 lowercase hexadecimal characters.'
        }
        if (Test-Path -LiteralPath $WorkspacePath) {
            throw "Refusing to overwrite workspace '$WorkspacePath'."
        }
        $resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
        $file = Get-Item -LiteralPath $resolvedMsi
        $hash = (Get-FileHash -LiteralPath $resolvedMsi -Algorithm SHA256).Hash
        $hash = $hash.ToLowerInvariant()
        $productVersion = Get-MsiProperty -Path $resolvedMsi `
            -Property 'ProductVersion'
        $productCode = Get-MsiProperty -Path $resolvedMsi `
            -Property 'ProductCode'
        if (-not $PSCmdlet.ShouldProcess(
            $WorkspacePath,
            'Create attended evidence workspace')) {
            return
        }
        $workspaceParent = Split-Path $WorkspacePath -Parent
        $workspaceStage = Join-Path $workspaceParent (
            '.goschedule-evidence-workspace-' +
                [guid]::NewGuid().ToString('N')
        )
        $now = [datetime]::UtcNow
        $scenarioIds = @(
            'access.intended-user',
            'access.unrelated-user-denied',
            'access.fresh-path-resolution',
            'access.path-removed',
            'window.clean-standard',
            'window.clean-high-or-mixed',
            'window.retained-profile',
            'window.state-transitions',
            'window.subsequent-launch',
            'error.daemon-unavailable',
            'error.access-denied',
            'error.timeout',
            'error.stream-disconnect',
            'error.repeated-refresh',
            'error.manual-retry',
            'error.recovery',
            'task.manual-success',
            'task.scheduled-success',
            'task.nonzero-exit',
            'task.start-failure',
            'setup.shortcut-defaults',
            'setup.shortcut-matrix',
            'setup.completion-matrix',
            'setup.finish-launch-integrity',
            'setup.cancel',
            'setup.maintenance',
            'setup.upgrade',
            'setup.invalid-input',
            'setup.rollback',
            'remove.preserve',
            'remove.wipe',
            'remove.cancel',
            'remove.multiple-profiles',
            'remove.locked-partial',
            'remove.reinstall-after-preserve',
            'remove.reinstall-after-wipe'
            'desktop.appearance-standard'
            'desktop.appearance-scaled'
            'desktop.interaction-states'
            'desktop.interaction-states-scaled'
            'desktop.navigation-options'
            'desktop.navigation-options-scaled'
            'desktop.scroll-input'
            'desktop.tasks-table'
            'desktop.tasks-table-scaled'
            'desktop.schedule-activity-tables'
            'desktop.schedule-activity-tables-scaled'
        )
        $observations = @(
            foreach ($scenarioId in $scenarioIds) {
                New-UnavailableObservation -Id $scenarioId -Timestamp $now
            }
        )
        $evidence = [ordered]@{
            schema_version = 1
            evidence_class = 'attended-windows'
            candidate = [ordered]@{
                repository = $Repository
                tag = $Tag
                commit = $Commit
                workflow = 'Release'
                run_id = $RunId
                run_attempt = $RunAttempt
                filename = $file.Name
                bytes = $file.Length
                sha256 = $hash
                product_version = $productVersion
                product_code = $productCode.ToUpperInvariant()
            }
            operator = [ordered]@{
                role = ''
                attested_at = $now.ToString('o')
                statement = ''
            }
            started_at = $now.ToString('o')
            completed_at = $now.ToString('o')
            environments = @()
            observations = $observations
            attachments = @()
        }
        try {
            New-Item -ItemType Directory -Path $workspaceStage | Out-Null
            New-Item -ItemType Directory -Path (
                Join-Path $workspaceStage 'attachments/windows'
            ) | Out-Null
            New-Item -ItemType Directory -Path (
                Join-Path $workspaceStage 'attachments/screenshots'
            ) | Out-Null
            New-Item -ItemType Directory -Path (
                Join-Path $workspaceStage 'attachments/tasks'
            ) | Out-Null
            New-Item -ItemType Directory -Path (
                Join-Path $workspaceStage 'fragments'
            ) | Out-Null
            Write-JsonNoBom -Value $evidence `
                -Path (Join-Path $workspaceStage 'evidence.json')
            foreach ($scenarioId in $scenarioIds) {
                if ($scenarioId.StartsWith('setup.') -or
                    $scenarioId.StartsWith('remove.') -or
                    $scenarioId.StartsWith('desktop.')) {
                    $templatePath = Join-Path $workspaceStage (
                        'fragments/' + $scenarioId + '.template.json'
                    )
                    Write-JsonNoBom -Value (
                        New-ScenarioTemplate -Id $scenarioId -Timestamp $now
                    ) -Path $templatePath
                }
            }
            [System.IO.Directory]::Move($workspaceStage, $WorkspacePath)
        } finally {
            if (Test-Path -LiteralPath $workspaceStage) {
                $resolvedStage = [System.IO.Path]::GetFullPath($workspaceStage)
                if ((Split-Path $resolvedStage -Parent) -ne
                    [System.IO.Path]::GetFullPath($workspaceParent) -or
                    -not (Split-Path $resolvedStage -Leaf).StartsWith(
                        '.goschedule-evidence-workspace-'
                    )) {
                    throw 'Refusing unsafe temporary workspace cleanup.'
                }
                [System.IO.Directory]::Delete($resolvedStage, $true)
            }
        }
        Write-Output "attended-evidence: initialized $WorkspacePath"
    }

    function Invoke-CaptureWindow {
        if (-not $WorkspacePath -or $ProcessId -le 0 -or
            -not $ObservationId -or -not $EnvironmentPath -or
            -not $FyneEvidencePath -or -not $ScreenshotPath) {
            throw ('CaptureWindow requires workspace, process ID, ' +
                'observation ID, environment JSON, Fyne evidence JSON, ' +
                'and a screenshot.')
        }
        Assert-AbsolutePath -Name 'WorkspacePath' -Value $WorkspacePath
        $workspace = (Resolve-Path -LiteralPath $WorkspacePath).Path
        $environment = Get-Content -LiteralPath $EnvironmentPath -Raw |
            ConvertFrom-Json
        $fyne = Get-Content -LiteralPath $FyneEvidencePath -Raw |
            ConvertFrom-Json
        if ($fyne.process_id -ne $ProcessId) {
            throw 'Fyne evidence process ID does not match ProcessId.'
        }
        $process = Get-Process -Id $ProcessId -ErrorAction Stop
        Add-NativeWindowType
        $native = [GoSchedule.ReleaseEvidence.NativeWindow]::Capture(
            [uint32]$ProcessId
        )
        if ($environment.account_sid -ne $native.ProcessUserSid -or
            $environment.integrity_rid -ne $native.IntegrityRid) {
            throw 'Native process token does not match the environment record.'
        }
        $snapshotPath = Join-Path $workspace (
            'attachments/windows/' + $ObservationId + '.json'
        )
        $fragmentPath = Join-Path $workspace (
            'fragments/' + $ObservationId + '.json'
        )
        if ((Test-Path -LiteralPath $snapshotPath) -or
            (Test-Path -LiteralPath $fragmentPath)) {
            throw ('Refusing to overwrite an existing window-capture ' +
                "artifact for '$ObservationId'.")
        }
        $snapshotRelative = Get-RelativeAttachmentPath `
            -Workspace $workspace -Path $snapshotPath
        $fyneRelative = Get-RelativeAttachmentPath `
            -Workspace $workspace -Path $FyneEvidencePath
        $screenshotRelative = $null
        if ($ScreenshotPath) {
            $screenshotRelative = Get-RelativeAttachmentPath `
                -Workspace $workspace -Path $ScreenshotPath
        }
        $snapshot = [ordered]@{
            schema_version = 1
            kind = 'native-window-v1'
            observation_id = $ObservationId
            captured_at = [datetime]::UtcNow.ToString('o')
            process_id = $ProcessId
            process_path = $process.Path
            process_sha256 = (
                Get-FileHash -LiteralPath $process.Path -Algorithm SHA256
            ).Hash.ToLowerInvariant()
            process_session_id = $process.SessionId
            process_user_sid = $native.ProcessUserSid
            process_integrity_rid = $native.IntegrityRid
            hwnd = ('0x{0:x16}' -f $native.Hwnd)
            outer_rect = Convert-Rect -Rect $native.Outer
            client_rect = Convert-Rect -Rect $native.Client
            monitor_rect = Convert-Rect -Rect $native.Monitor
            work_area_rect = Convert-Rect -Rect $native.WorkArea
            monitor_id = $native.MonitorId
            effective_dpi = $native.Dpi
            show_command = $native.ShowCommand
            visible = $native.Visible
            maximized = $native.Maximized
            minimized = $native.Minimized
            fullscreen = $native.Fullscreen
            restored = $native.Restored
            fyne = $fyne
        }
        if (-not $PSCmdlet.ShouldProcess(
            $snapshotPath,
            'Write native window measurement attachment')) {
            return
        }
        $attachmentPaths = @(
            $snapshotRelative
            $fyneRelative
        )
        if ($screenshotRelative) {
            $attachmentPaths += $screenshotRelative
        }
        $dpiScale = [double]$native.Dpi / 96.0
        $logicalWorkWidth = (
            $native.WorkArea.Right - $native.WorkArea.Left
        ) / $dpiScale
        $logicalWorkHeight = (
            $native.WorkArea.Bottom - $native.WorkArea.Top
        ) / $dpiScale
        $fragment = [ordered]@{
            environment = $environment
            observation = [ordered]@{
                id = $ObservationId
                environment_id = $environment.id
                status = 'unavailable'
                started_at = [datetime]::UtcNow.ToString('o')
                completed_at = [datetime]::UtcNow.ToString('o')
                summary = 'Review native measurements and visual evidence.'
                metrics = [ordered]@{
                    pid = $ProcessId
                    executable_sha256 = $snapshot.process_sha256
                    process_session_id = $snapshot.process_session_id
                    process_user_sid = $snapshot.process_user_sid
                    process_integrity_rid = $snapshot.process_integrity_rid
                    hwnd = $snapshot.hwnd
                    outer_rect = $snapshot.outer_rect
                    client_rect = $snapshot.client_rect
                    monitor_rect = $snapshot.monitor_rect
                    work_area_rect = $snapshot.work_area_rect
                    fyne_content_width = $fyne.content_width
                    fyne_content_height = $fyne.content_height
                    fyne_scale = $fyne.canvas_scale
                    logical_work_area_width = $logicalWorkWidth
                    logical_work_area_height = $logicalWorkHeight
                    effective_dpi = $native.Dpi
                    monitor_id = $native.MonitorId
                    restored = $native.Restored
                    maximized = $native.Maximized
                    minimized = $native.Minimized
                    fullscreen = $native.Fullscreen
                    margins_visible = $false
                    title_bar_reachable = $false
                    resize_borders_reachable = $false
                    taskbar_reachable = $false
                }
                attachment_paths = $attachmentPaths
            }
        }
        switch ($ObservationId) {
            'window.clean-standard' {
                $fragment.observation.metrics.launch_sequence = 0
            }
            'window.state-transitions' {
                $fragment.observation.metrics.maximize_worked = $false
                $fragment.observation.metrics.restore_worked = $false
                $fragment.observation.metrics.resize_worked = $false
                $fragment.observation.metrics.minimize_worked = $false
                $fragment.observation.metrics.close_worked = $false
            }
            'window.subsequent-launch' {
                $fragment.observation.metrics.prior_process_id = 0
                $fragment.observation.metrics.launch_sequence = 0
                $fragment.observation.metrics.fresh_process = $false
                $fragment.observation.metrics.prior_process_closed = $false
            }
        }
        $snapshotTemporary = "$snapshotPath.$([guid]::NewGuid().ToString('N')).tmp"
        $fragmentTemporary = "$fragmentPath.$([guid]::NewGuid().ToString('N')).tmp"
        $snapshotPublished = $false
        try {
            Write-JsonNoBom -Value $snapshot -Path $snapshotTemporary
            Write-JsonNoBom -Value $fragment -Path $fragmentTemporary
            [System.IO.File]::Move($snapshotTemporary, $snapshotPath)
            $snapshotPublished = $true
            [System.IO.File]::Move($fragmentTemporary, $fragmentPath)
        } catch {
            if ($snapshotPublished -and (Test-Path -LiteralPath $snapshotPath)) {
                [System.IO.File]::Delete($snapshotPath)
            }
            throw
        } finally {
            foreach ($temporary in @($snapshotTemporary, $fragmentTemporary)) {
                if (Test-Path -LiteralPath $temporary) {
                    [System.IO.File]::Delete($temporary)
                }
            }
        }
        Write-Output "attended-evidence: review fragment $fragmentPath"
    }

    function Invoke-RecordObservation {
        if (-not $WorkspacePath -or -not $ObservationPath) {
            throw 'RecordObservation requires workspace and observation path.'
        }
        $workspace = (Resolve-Path -LiteralPath $WorkspacePath).Path
        $evidencePath = Join-Path $workspace 'evidence.json'
        $evidence = Get-Content -LiteralPath $evidencePath -Raw |
            ConvertFrom-Json
        $fragment = Get-Content -LiteralPath $ObservationPath -Raw |
            ConvertFrom-Json
        if (-not $fragment.environment.id -or -not $fragment.observation.id) {
            throw 'Fragment requires environment and observation records.'
        }
        $matches = @(
            $evidence.observations |
                Where-Object id -EQ $fragment.observation.id
        )
        if ($matches.Count -ne 1) {
            throw 'Fragment observation ID is not one required placeholder.'
        }
        if ($matches[0].summary -ne 'not yet observed') {
            throw "Observation '$($fragment.observation.id)' already exists."
        }
        $environmentMatches = @(
            $evidence.environments |
                Where-Object id -EQ $fragment.environment.id
        )
        if ($environmentMatches.Count -eq 0) {
            $evidence.environments = @(
                $evidence.environments
            ) + @($fragment.environment)
        } elseif (
            ($environmentMatches[0] | ConvertTo-Json -Depth 10 -Compress) -ne
            ($fragment.environment | ConvertTo-Json -Depth 10 -Compress)
        ) {
            throw "Environment '$($fragment.environment.id)' changed."
        }
        for ($index = 0; $index -lt $evidence.observations.Count; $index++) {
            if ($evidence.observations[$index].id -eq
                $fragment.observation.id) {
                $evidence.observations[$index] = $fragment.observation
                break
            }
        }
        if (-not $PSCmdlet.ShouldProcess(
            $evidencePath,
            "Record observation '$($fragment.observation.id)'")) {
            return
        }
        Save-JsonReplacingFile -Value $evidence -Path $evidencePath
        Write-Output ('attended-evidence: recorded ' +
            $fragment.observation.id)
    }

    function Invoke-Finalize {
        if (-not $WorkspacePath -or -not $MsiPath) {
            throw 'Finalize requires workspace and MSI path.'
        }
        $workspace = (Resolve-Path -LiteralPath $WorkspacePath).Path
        $resolvedMsi = (Resolve-Path -LiteralPath $MsiPath).Path
        $evidencePath = Join-Path $workspace 'evidence.json'
        $evidence = Get-Content -LiteralPath $evidencePath -Raw |
            ConvertFrom-Json
        $attachmentPaths = @(
            $evidence.observations.attachment_paths |
                Where-Object { $_ } |
                Sort-Object -Unique
        )
        $attachments = @(
            foreach ($relative in $attachmentPaths) {
                if ($relative -notmatch '^attachments/' -or
                    $relative.Contains('..') -or $relative.Contains('\')) {
                    throw "Unsafe attachment path '$relative'."
                }
                $full = Join-Path $workspace $relative
                $file = Get-Item -LiteralPath $full
                $current = $file
                while ($current.FullName -ne $workspace) {
                    if (($current.Attributes -band
                        [System.IO.FileAttributes]::ReparsePoint) -ne 0) {
                        throw "Reparse attachment path is not allowed: '$full'."
                    }
                    $current = $current.Parent
                    if (-not $current) {
                        throw "Attachment escaped workspace: '$full'."
                    }
                }
                $mediaType = switch ($file.Extension.ToLowerInvariant()) {
                    '.json' { 'application/json' }
                    '.png' { 'image/png' }
                    '.jpg' { 'image/jpeg' }
                    '.jpeg' { 'image/jpeg' }
                    default { 'text/plain' }
                }
                $purpose = switch -Regex ($relative) {
                    '^attachments/windows/window\..+\.json$' {
                        'native window measurement'
                        break
                    }
                    '^attachments/tasks/.+\.json$' {
                        'task run evidence'
                        break
                    }
                    default { 'attended release-candidate evidence' }
                }
                [ordered]@{
                    path = $relative
                    bytes = $file.Length
                    sha256 = (
                        Get-FileHash -LiteralPath $full -Algorithm SHA256
                    ).Hash.ToLowerInvariant()
                    media_type = $mediaType
                    purpose = $purpose
                }
            }
        )
        $now = [datetime]::UtcNow.ToString('o')
        $evidence.completed_at = $now
        $evidence.operator.role = $OperatorRole
        $evidence.operator.attested_at = $now
        $evidence.operator.statement = (
            'I attest that these observations were made against the ' +
            'identified candidate in the identified environments.'
        )
        $evidence.attachments = $attachments
        if (-not $PSCmdlet.ShouldProcess(
            $evidencePath,
            'Finalize and validate attended evidence')) {
            return
        }
        Save-JsonReplacingFile -Value $evidence -Path $evidencePath
        $root = (Resolve-Path -LiteralPath (
            Join-Path $PSScriptRoot '../..'
        )).Path
        $arguments = @(
            'run', './scripts/windows-release-gate', 'validate',
            '--evidence', $evidencePath,
            '--artifact', $resolvedMsi,
            '--repository', $evidence.candidate.repository,
            '--tag', $evidence.candidate.tag,
            '--commit', $evidence.candidate.commit
        )
        $result = Invoke-HiddenProcess -FilePath 'go' `
            -Arguments $arguments -WorkingDirectory $root
        if ($result.ExitCode -ne 0) {
            [Console]::Error.Write($result.Stderr)
            throw "Release evidence validation failed ($($result.ExitCode))."
        }
        Write-Output $result.Stdout.Trim()
        $archiveName = 'go-schedule_{0}_windows-attended-evidence.zip' -f (
            $evidence.candidate.tag
        )
        $archivePath = Join-Path (Split-Path $workspace -Parent) $archiveName
        $archiveParent = Split-Path $archivePath -Parent
        $bundleStage = Join-Path $archiveParent (
            '.goschedule-evidence-bundle-' + [guid]::NewGuid().ToString('N')
        )
        New-Item -ItemType Directory -Path $bundleStage | Out-Null
        try {
            Copy-Item -LiteralPath $evidencePath -Destination (
                Join-Path $bundleStage 'evidence.json'
            )
            foreach ($relative in $attachmentPaths) {
                $source = Join-Path $workspace $relative
                $destination = Join-Path $bundleStage $relative
                New-Item -ItemType Directory -Path (
                    Split-Path $destination -Parent
                ) -Force | Out-Null
                Copy-Item -LiteralPath $source -Destination $destination
            }
            $temporaryArchive = Join-Path $archiveParent (
                '.goschedule-evidence-archive-' +
                    [guid]::NewGuid().ToString('N') + '.zip'
            )
            [System.IO.Compression.ZipFile]::CreateFromDirectory(
                $bundleStage,
                $temporaryArchive,
                [System.IO.Compression.CompressionLevel]::Optimal,
                $false
            )
            [System.IO.File]::Move($temporaryArchive, $archivePath)
        } catch {
            if ($temporaryArchive -and
                (Test-Path -LiteralPath $temporaryArchive)) {
                [System.IO.File]::Delete($temporaryArchive)
            }
            throw
        } finally {
            $resolvedStage = [System.IO.Path]::GetFullPath($bundleStage)
            if (-not $resolvedStage.StartsWith(
                [System.IO.Path]::GetFullPath($archiveParent) +
                    [System.IO.Path]::DirectorySeparatorChar
            ) -or -not (
                Split-Path $resolvedStage -Leaf
            ).StartsWith('.goschedule-evidence-bundle-')) {
                throw 'Refusing unsafe temporary bundle cleanup.'
            }
            if (Test-Path -LiteralPath $resolvedStage) {
                [System.IO.Directory]::Delete($resolvedStage, $true)
            }
        }
        Write-Output "attended-evidence: finalized $archivePath"
    }

#_______________________________________________________________________________
## Declare Variables and Arrays

    $ThisScriptPath = $MyInvocation.MyCommand.Path
    $ErrorActionPreference = 'Stop'

#_______________________________________________________________________________
## Execute Operations

    if (($Help) -or ($PSCmdlet.ParameterSetName -eq 'HelpText')) {
        Get-Help $ThisScriptPath -Detailed
        exit 0
    }
    if (-not $IsWindows) {
        throw 'Attended release evidence collection requires Windows.'
    }

    switch ($Action) {
        'Initialize' { Invoke-Initialize }
        'CaptureWindow' { Invoke-CaptureWindow }
        'RecordObservation' { Invoke-RecordObservation }
        'Finalize' { Invoke-Finalize }
        default { throw "Unsupported action '$Action'." }
    }

#_______________________________________________________________________________
## End of script
