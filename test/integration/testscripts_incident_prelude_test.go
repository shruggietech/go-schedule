package integration

// This is authored into the original owned harness after its param block.
// It does not wrap invocation or replace the interpreter's native exit.
func powerShellIncidentPrelude() string {
	return `
if ($env:GOSCHEDULE_INCIDENT_IDENTITY) {
  try {
    if ([System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Windows)) {
      Add-Type -ErrorAction Stop -TypeDefinition 'using System.Runtime.InteropServices; public static class OwnedIncidentErrorMode { [DllImport("kernel32.dll")] public static extern uint GetErrorMode(); [DllImport("kernel32.dll")] public static extern uint SetErrorMode(uint mode); }'
      [void][OwnedIncidentErrorMode]::SetErrorMode([OwnedIncidentErrorMode]::GetErrorMode() -bor 3)
    }
    $module = [Diagnostics.Process]::GetCurrentProcess().Modules | Where-Object { $_.ModuleName -in @('coreclr.dll','libcoreclr.so','libcoreclr.dylib') } | Select-Object -First 1
    $identity = @{ observation='pre-payload'; pid=$PID; powershell=$PSVersionTable.PSVersion.ToString(); framework=[Runtime.InteropServices.RuntimeInformation]::FrameworkDescription; runtime_module=$null; runtime_sha256=$null }
    if ($module) { $identity.runtime_module=$module.FileName; $identity.runtime_sha256=(Get-FileHash -ErrorAction Stop -LiteralPath $module.FileName -Algorithm SHA256).Hash.ToLowerInvariant() }
    [IO.File]::WriteAllText($env:GOSCHEDULE_INCIDENT_IDENTITY, ($identity | ConvertTo-Json -Compress), [Text.UTF8Encoding]::new($false))
  } catch { }
}
`
}
