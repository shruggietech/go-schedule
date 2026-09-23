# Implementation Plan: Linux desktop daemon presence and controls

**Branch**: `codex/102-linux-daemon-indicator` | **Date**: 2026-09-23 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #256 with a per-session Linux StatusNotifier indicator, local service state and control in both indicator and GUI, one-window activation, safe system-service authorization, and portable desktop-bundle startup assets. Preserve a daemon that runs without any graphical session. Native outcome notifications and SMTP remain separate.

## Technical Context

**Language/Version**: Repository Go 1.26, TypeScript/React in the Wails desktop module
**Primary Dependencies**: Existing godbus and kardianos/service; `github.com/gogpu/systray` v0.3.0 (MIT) for StatusNotifierItem and dbusmenu
**Storage**: No scheduler or user-data schema change; session singleton uses runtime lock and socket
**Testing**: Go state/action/activation tests, frontend component tests, Linux build and packaging checks, `scripts/verify.sh all`, hosted CI
**Target Platform**: Linux amd64 graphical sessions with a D-Bus session bus and StatusNotifier host; GUI fallback works without a host; daemon remains headless
**Project Type**: Multi-module Go Wails desktop app and portable Linux archive
**Performance Goals**: Bounded health probe, no daemon hot-path change, low-rate status refresh
**Constraints**: No shell interpretation for actions, narrow graphical authorization, no GUI elevation, no false Running state, no duplicate session item
**Scale/Scope**: One indicator and GUI per signed-in session

## Constitution Check

- I Code quality: platform-specific behavior isolated in Linux-tagged files; errors explicit; polling and activation goroutines terminate with context.
- II Testing: state classification, action failures, singleton, and GUI no-autostart regressions are tested; full parity remains mandatory.
- III UX: status language is shared with S101; no remote/local confusion; missing host is explained through docs and GUI fallback.
- IV Performance: bounded status checks and low-rate polling outside scheduler paths.
- V Autopilot: spec-kit sequence and blocking analysis precede implementation; review branch and PR. The S102 user request expressly authorizes publication.
- Engineering constraints: new dependency is MIT and avoids CGO/system toolkit coupling for the indicator; Linux desktop build remains the only consumer.

**Gate**: Pass, subject to post-design analysis.

## Project Structure

```text
specs/102-linux-daemon-indicator/
internal/service/
internal/desktopcontrol/
cmd/gosched-indicator/
desktop/
desktop/frontend/src/
brand/platform/linux/
docs/INSTALL-linux.md
.github/workflows/release.yml
scripts/automation-check.sh
```

## Design Decisions

1. **Separate session process**: The system daemon cannot own user-session UI and the GUI may close. Package one per-session indicator launched by the desktop autostart contract.
2. **StatusNotifier library**: Use `gogpu/systray` v0.3.0 for SNI and dbusmenu rather than hand-writing two D-Bus protocols. It is pure Go, MIT, and uses the already present godbus transport. Keep it in the Linux-only command so Windows code paths remain unchanged.
3. **Service state**: Linux installed-service truth comes from a bounded systemd query; Running additionally needs fresh local IPC health. An absent unit remains distinct from Stopped, and an error remains Unknown.
4. **Mutations**: Request only fixed `systemctl` verbs and the fixed `goschedd.service` unit through `pkexec` when unprivileged. No shell, user-supplied unit, or elevated GUI. Observe the final state after helper exit.
5. **GUI startup**: Auto-spawn a bundled daemon only if the installed service is absent. This corrects a Linux-specific S101 gap where opening the GUI could reverse a deliberate system-service stop.
6. **Session uniqueness**: File locks guard both GUI and indicator; a same-user local activation socket focuses an existing GUI. No session bus is required for GUI fallback.
7. **Packaging**: Ship the indicator binary and an XDG autostart entry in the Linux desktop archive. Installing the autostart entry is explicit; the daemon never depends on it.

## Complexity Tracking

| Added element | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| Separate indicator binary | Outlives GUI, lives in user session | GUI-hosted item vanishes with GUI |
| StatusNotifier library | Menu and icon D-Bus protocols | Hand-written protocol has higher correctness risk |
| Linux session activation socket | One GUI even without D-Bus | D-Bus-only activation would break the no-bus fallback |

## Post-Design Constitution Check

The design keeps service-manager mutations outside the daemon and GUI, preserves least authority and truthful status, and leaves scheduler data and hot paths unchanged. The deliberate Linux startup correction is required by FR-005 and is explicitly noted as a deviation from the pre-S102 Linux GUI behavior.
