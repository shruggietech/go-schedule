# Contract: Windows MSI package

Built with WiX Toolset v5 from `build/windows/goschedule.wxs` on the `windows-latest` release job.

## Package properties

| Property            | Value                                                              |
|---------------------|-------------------------------------------------------------------|
| Product name        | `go-schedule`                                                     |
| Manufacturer        | ShruggieTech                                                      |
| Install scope       | per-machine (`ALLUSERS=1`) → requires elevation (UAC)            |
| Install dir         | `C:\Program Files\go-schedule\`                                    |
| Upgrade behavior    | `MajorUpgrade` (uninstall prior, install new; preserves data dir) |
| Artifact name       | `go-schedule_<version>_windows_amd64.msi`                          |

## Components / files

Installed to the install dir:

- `goschedd.exe` — daemon (hosted as a service, below)
- `gosched-gui.exe` — desktop GUI (windowless build; Start-Menu shortcut target)
- `gosched.exe` — CLI

## Service (WiX `ServiceInstall` + `ServiceControl`)

| Attribute     | Value                                                  |
|---------------|--------------------------------------------------------|
| Name          | `goschedd` (matches `internal/service` svcName)        |
| Display name  | `go-schedule`                                          |
| Description   | `Cross-platform task scheduler daemon`                 |
| Start type    | `auto` (starts on boot, no interactive login — FR-006) |
| Account       | LocalSystem (default)                                  |
| On install    | start the service                                      |
| On uninstall  | stop + remove the service                              |

The service name matches what the CLI `service` control layer expects, so
`gosched service status/start/stop` work against the MSI-installed service. The daemon binary is
unchanged (it already runs under the Windows SCM via `service.Run`).

## Shortcuts

- Start-Menu shortcut → `gosched-gui.exe` (FR-007). Launching it opens the GUI with no console
  window and reuses the already-running service (FR-010) via the existing health-check/autostart
  logic (the service answers the ping, so no second daemon spawns).

## Data directory

- `C:\ProgramData\goschedule\` holds the SQLite DB (`goschedule.db`) and `logs/`. Created by the
  daemon on first run; the MSI does not seed it.

## Uninstall behavior (FR-009)

- Removes: all installed binaries, the service registration, and shortcuts.
- **Leaves**: `C:\ProgramData\goschedule\` (DB + logs) so reinstall preserves tasks. Documented;
  a manual "remove all data" step is described in the install guide.
- Add/Remove Programs (Apps & features) lists `go-schedule` with working Uninstall.

## Out of scope (documented)

- Code-signing the MSI (separate infra; unsigned acceptable with SmartScreen guidance).
- Per-user (non-elevated) install mode.
- ARM64 MSI (amd64 only for this feature unless trivially added).

## CI / release

- The Windows GUI job builds the three binaries, then invokes WiX to produce the `.msi`.
- The portable Windows `.zip` (both the daemon-only and the desktop bundle) is **removed** from
  the Windows release outputs; `SHA256SUMS.txt` covers the `.msi`.
- Linux/macOS archives are unchanged.
