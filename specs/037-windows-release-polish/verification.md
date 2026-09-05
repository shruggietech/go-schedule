# Verification: Windows Release Polish

**Date**: 2026-08-31

## Test-First Evidence

- The new Windows monitor-selection test initially failed to compile because `winPoint` and `workAreaForPoint` did not exist. This proves the constrained secondary-monitor contract preceded its adapter implementation.
- The sizing table initially failed against the old full-work-area and 640x480 clamp behavior.
- After the source verifier was strengthened but before WiX authoring changed, `build/windows/verify_wxs.ps1` failed with: `SummaryInformation Description must be "go-schedule: cross-platform task scheduler"`.

## Issue #89: Bounded Windowed Launch

- `go test ./gui -run "TestWindowSizeFor|TestWorkAreaForPoint|TestRootContentFits" -count=1`: PASS.
- `go test ./gui/... -count=1`: PASS.
- Table coverage includes unknown work area, 800x600, 1024x768, 1366x768, 1920x1080, a 3840x2160 physical work area at 200 percent scaling, a work area below 640x480, and invalid scale.
- Expected effective 800x600 startup viewport: 720x540 logical units. The complete root content minimum fits within that viewport.
- Windows monitor selection is tested with a constrained 800x560 secondary work area at positive desktop coordinates and proves the selected monitor handle supplies the size.
- Specification 003 FR-001, User Story 5, SC-001, assumptions, and verification record now document the intentional reversal. No maximize, restore, resize, minimize, fullscreen, or placement-persistence behavior was added or removed.

### Review Resolution (2026-09-01)

- The mixed-DPI review found that the selected monitor's physical work area was divided by the not-yet-positioned Fyne window's scale, which could belong to a different display.
- The Windows adapter now obtains `MDT_EFFECTIVE_DPI` with `GetDpiForMonitor` from the same `HMONITOR` used by `GetMonitorInfoW`. The initial canvas scale is retained only as a compatibility fallback when monitor DPI is unavailable.
- The regression test proves both lookups receive the same selected monitor handle and that 192 DPI produces a 2.0 scale. A separate fallback test proves a valid work area remains usable when the DPI API is unavailable.

## Issue #88: MSI Explorer Subject

- `pwsh build/windows/verify_wxs.ps1`: PASS with the exact approved source value.
- Both changed PowerShell scripts parse without errors under PowerShell 7.
- A local WiX 6.0.2 candidate build completed from the S037 working tree.
- Candidate SHA-256: `648bfa402752af1504155fa911cc60679c34e8ea20c1350f22c74cb5ff48bde5`.
- Summary PID 3 and native Shell `System.Subject` both equal `go-schedule: cross-platform task scheduler`.
- Preserved compiled identity: Title `Installation Database`, Author and Manufacturer `ShruggieTech`, standard WiX Comments, ProductName `go-schedule`, and UpgradeCode `{B6F3C2E1-7A4D-4C9E-9B2A-1F6D8E5A0C34}`.
- Full generated evidence: [candidate-msi-evidence.md](candidate-msi-evidence.md).
- The release filename expression remains `go-schedule_${VERSION}_windows_amd64.msi`; install and lifecycle authoring is unchanged.

## Canonical Eight-Gate Verification

`sh scripts/verify.sh all` ran in one foreground process and exited 0:

| Gate | Result | Evidence |
| --- | --- | --- |
| format | PASS | No files printed. |
| vet | PASS | `go vet ./...` completed without findings. |
| lint | PASS | golangci-lint v2.12.0 reported `0 issues`; cache path-relativity warnings named other worktree paths only. |
| race | PASS | All selected packages and the 56.708-second integration suite passed under the race detector. |
| gui | PASS | `gui` and `gui/viewmodel` passed. |
| coverage | PASS | engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, logbus 91.1%. |
| docs | PASS | 14 pages, links, front matter, fences, theme, and product policy clean. |
| automation | PASS | Actions, CodeQL, Dependabot, release notes, brand, lifecycle, and eight-gate contract clean. |

## Delivery Audit

- Specification lifecycle state: Implemented after all required tasks and the canonical aggregate completed.
- Scope remains exclusive to issues #88 and #89.
- Pinned changes are limited to Windows installer authoring/source verification and the release workflow's new compiled-artifact gate, with dated decisions in `CHANGELOG.md`.
- Constitution principles I through V pass with no deviation: code and tests are focused, no safety-critical surface was weakened, the restored startup is consistent, no scheduler performance path changed, and review publication follows the authorized S037 protocol.
