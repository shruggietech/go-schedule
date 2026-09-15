# S086 Research

## Installer delay

Static inspection found no MSI-authored prerequisite download. GUI and documentation launch only after Finish, so WebView2 startup acquisition cannot explain the pre-wizard delay. ServiceControl waits might explain execution-phase waits, not pre-wizard preparation. Existing wrappers wait indefinitely and lack phase timing. Mapped-share latency, Installer initialization, endpoint protection, and service operation remain unverified hypotheses. Stage identical bytes guest-locally and retain flushed logs and elapsed timing before claiming any cause or fix.

## Preparation design

Microsoft documents XML `.wsb` configurations with mapped folders established before LogonCommand. Map inputs read-only and only a dedicated export folder read-write; use fixed guest paths and XML-escaped host mapping data. Closing Sandbox destroys guest state. [Microsoft Sandbox configuration](https://learn.microsoft.com/en-us/windows/security/application-security/application-isolation/windows-sandbox/windows-sandbox-configure-using-wsb-file).

Windows Installer supports verbose logging and `!` to flush each line. Use `/L*vx!`, retain logs, and bound diagnostic waits without killing the Installer service. [Microsoft msiexec reference](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/msiexec).

Supply the hashed offline Evergreen Standalone WebView2 installer instead of another application-startup download wait. [Microsoft WebView2 distribution](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/distribution).

## Evidence and refresh boundaries

Reuse the attended collector's unavailable placeholders and shared gate. Mechanical logs cannot attest visual behavior, normal-user access, mixed DPI, or upgrade preservation. Public v1.1.1 baseline must be independently sourced. Existing v1.4.0 commit 57555ffa413df641cb21784847ad598c113bf199 excludes three UI merges. Preserve old assets; reviewed tag replacement, hosted staging, and promotion need separate authority. Do not close #226/#228 from tooling tests.
