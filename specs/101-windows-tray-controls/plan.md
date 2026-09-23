# Implementation Plan: Windows local service tray and controls

**Branch**: `codex/101-windows-tray-controls` | **Date**: 2026-09-22 | **Spec**: [spec.md](spec.md)

## Summary

Deliver issue #255 with a separate windowless Windows user-session companion, shared local service state/control logic, GUI integration, one-instance GUI activation, and complete MSI/build lifecycle. Preserve deliberate service stops. Linux tray, native desktop notifications, and SMTP remain out of scope.

## Technical Context

**Language/Version**: Repository Go version, TypeScript/React in the Wails desktop module
**Primary Dependencies**: Wails v2, kardianos/service, golang.org/x/sys/windows, Windows SCM and Shell APIs
**Storage**: Existing GUI preferences only; no new persistent state
**Testing**: Go unit/integration tests, frontend tests, Windows installer contract checks, `scripts/verify.sh all`, hosted CI
**Target Platform**: Windows 10/11 installed service and GUI; Linux build behavior must remain unchanged
**Project Type**: Multi-module Go desktop application with WiX MSI
**Performance Goals**: Status polling bounded and nonblocking, no busy loop or duplicate process
**Constraints**: No console flash, user-session UI separate from service, narrow UAC, truthful SCM plus IPC health, no elevation for reads
**Scale/Scope**: One companion per logon session and one GUI per user session

## Constitution Check

- I Code quality: shared state model, explicit errors, bounded goroutines and context cancellation, no ignored outcomes.
- II Testing: add state mapping, action, lifecycle, GUI and installer regressions alongside code; full local parity and CI.
- III UX: precise statuses, confirmation, cancellation/failure feedback, consistent tray and GUI meanings.
- IV Performance: bounded status checks and timers, no scheduler hot-path change.
- V Autopilot: spec-kit complete before implementation, review branch and PR. The user's S101 instruction explicitly authorizes the pre-publication gate and push/PR.
- Engineering constraints: no new dependency without clear justification; Windows-specific code behind build tags; Linux tests continue to pass.

**Gate**: Pass, subject to analysis of detailed artifacts.

## Project Structure

```text
specs/101-windows-tray-controls/
  spec.md
  plan.md
  research.md
  data-model.md
  contracts/
  quickstart.md
  tasks.md
internal/service/
internal/autostart/
internal/desktopcontrol/
cmd/gosched-tray/
desktop/
desktop/frontend/src/
build/windows/
test/integration/
.github/workflows/
```

**Structure Decision**: Put platform-neutral state/control orchestration in an internal package; isolate SCM, notification-area, and activation details in Windows-tagged files. Keep the Wails GUI, headless daemon, and per-session companion as separate process roles.

## Design Decisions

1. **Separate companion**: The daemon runs outside the user session, while the GUI can close. Hosting the icon in either would violate the lifecycle requirement.
2. **State projection**: Preserve raw SCM state, then combine with bounded local health. SCM Running without healthy IPC is Unknown/unreachable, not Running.
3. **Narrow elevation**: Query status as ordinary user. Launch only the service-mutation helper through Windows UAC; observe SCM/health after it exits.
4. **GUI startup**: Skip bundled-daemon autospawn whenever the installed service exists, including Stopped or transitional states.
5. **One-instance activation**: Companion and GUI use a session-scoped activation handshake. Open focuses the existing GUI, otherwise starts it.
6. **Installer**: Include signed companion binary, logon registration, upgrade/removal shutdown handling, and user documentation. A Run entry is per-logon and Windows may delay its execution; the product should not promise immediate icon appearance at sign-in.
7. **Artwork exception**: Approve named monochrome reduced-mark exports with exact canonical geometry for tray contrast; record them in the brand guide and integrity manifest rather than silently substituting full-detail marks.

## Complexity Tracking

| Added element | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| Separate tray binary | Independent session lifetime | A Wails-hosted icon vanishes when GUI closes |
| Narrow elevated helper mode | UAC without elevating UI | Elevating Wails or tray grants a broad process unnecessary rights |
| GUI activation handshake | Focus existing window | Launching the GUI executable alone creates duplicates |

## Post-Design Constitution Check

The design preserves all five principles. Windows-specific native interop and one additional process are necessary for #255's service/session boundary, not general architectural expansion. No scheduling data schema changes.
