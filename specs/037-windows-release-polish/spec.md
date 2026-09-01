# Feature Specification: Windows Release Polish

**Feature Branch**: `codex/037-windows-release-polish`

**Created**: 2026-08-31

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: review branch `codex/037-windows-release-polish`; local verification completed 2026-08-31

**Input**: GitHub issues #88 and #89: revise the Windows MSI Explorer subject and open the desktop GUI at a bounded, DPI-aware windowed size.

## User Scenarios & Testing

### User Story 1 - Reachable Windowed Launch (Priority: P1)

A Windows desktop user launches go-schedule into a restored window that is large enough for normal work but leaves visible desktop margins. The full frame and primary actions remain reachable on constrained and high-DPI displays.

**Why this priority**: A fresh launch currently consumes essentially the whole work area and can exceed very small displays, so the first desktop interaction is misleading and potentially unreachable.

**Independent Test**: Exercise the startup-size contract for unknown work area, 800x600, 1024x768, 1366x768, 1920x1080, and 200 percent scaling, then confirm the top-level views fit within the effective 800x600 startup viewport.

**Acceptance Scenarios**:

1. **Given** at least 1423x889 logical units of work area, **When** the GUI launches without saved placement, **Then** it requests a restored 1280x800 window.
2. **Given** a smaller work area, **When** the GUI launches, **Then** neither requested dimension exceeds 90 percent of the available logical dimension.
3. **Given** a work area smaller than 640x480, **When** the GUI launches, **Then** the requested size stays within the reported work area rather than being enlarged to a nominal minimum.
4. **Given** a 200 percent display scale, **When** the GUI computes its size, **Then** physical work-area pixels are converted to logical units before the cap is applied.
5. **Given** an explicit user maximize or restore action after launch, **When** Windows applies it, **Then** normal window-state behavior remains available.

---

### User Story 2 - Concise Explorer Metadata (Priority: P1)

A Windows user inspecting a newly built MSI in Explorer sees `go-schedule: cross-platform task scheduler` as its Subject, with no U+2014 character, while all unrelated installer identity and behavior remain stable.

**Why this priority**: The published installer exposes copy that conflicts with the project's current prose convention, and source-only checking cannot prove the compiled artifact users receive.

**Independent Test**: Build a candidate MSI, inspect Summary Information PID 3 through the Windows Installer API, and record the artifact hash and exact Subject value.

**Acceptance Scenarios**:

1. **Given** a newly compiled candidate MSI, **When** its Summary Information PID 3 is inspected, **Then** the exact value is `go-schedule: cross-platform task scheduler`.
2. **Given** the candidate MSI in Windows Explorer, **When** its tooltip or Properties metadata is viewed, **Then** the revised Subject appears.
3. **Given** the metadata-only source change, **When** the compiled installer is inspected, **Then** Title, Author, Comments, ProductName, manufacturer, upgrade identity, install behavior, and release filename are unchanged.

### Edge Cases

- Unknown, zero, or invalid display scale falls back to the existing 1280x800 restored request.
- Work-area width and height are capped independently so unusually narrow or short displays remain reachable.
- The Windows work-area query must use the monitor associated with the launch context rather than assuming the primary monitor.
- The 90 percent cap is a hard upper bound even when it is below the old 640x480 nominal minimum.
- Compiled-MSI inspection must fail for a missing or merely similar Subject value and must report the observed value.
- Source checks and compiled-table checks do not claim that Explorer was observed; native Explorer observation remains distinct evidence.

## Requirements

### Functional Requirements

- **FR-001**: A fresh GUI launch MUST request `min(1280x800, 90 percent of the current monitor work area)` in logical units, with each dimension capped independently.
- **FR-002**: Physical work-area dimensions MUST be divided by the effective display scale before the 90 percent cap is applied.
- **FR-003**: A valid work area MUST be a hard upper bound; no nominal minimum may enlarge a requested dimension beyond that cap.
- **FR-004**: An unknown work area or invalid scale MUST retain the 1280x800 restored fallback.
- **FR-005**: Windows work-area discovery MUST select the monitor nearest the launch pointer and use that monitor's reserved work area, including taskbar reservation.
- **FR-006**: The title bar, resize borders, primary navigation, and primary action on every top-level tab MUST remain reachable at an effective 800x600 display.
- **FR-007**: The change MUST preserve normal user-controlled maximize, restore, resize, minimize, and fullscreen behavior and MUST NOT persist placement as part of this slice.
- **FR-008**: Specification 003 FR-001 and its verification record MUST document that S037 intentionally supersedes the former full-work-area startup requirement.
- **FR-009**: The Windows MSI Summary Information PID 3 Subject MUST equal `go-schedule: cross-platform task scheduler` and MUST contain no U+2014 character.
- **FR-010**: The compiled-artifact inspector MUST assert the exact PID 3 value and record both the inspected Subject and candidate SHA-256 hash.
- **FR-011**: The MSI Title, Author, Comments, ProductName, manufacturer, UpgradeCode, install behavior, release filename, and all non-Subject metadata MUST remain unchanged.
- **FR-012**: Source verification MUST prevent the approved Subject from drifting before compilation, while compiled-MSI verification remains the authoritative artifact check.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All six required sizing cases produce the exact expected logical dimensions in automated table tests.
- **SC-002**: At effective 800x600, every top-level tab keeps its primary navigation and primary action within a 720x540 launch viewport.
- **SC-003**: A Windows constrained-monitor check demonstrates that sizing comes from the monitor nearest the launch pointer and not a larger primary monitor.
- **SC-004**: A compiled candidate MSI inspection records one SHA-256 hash and the exact approved Subject, with zero U+2014 characters.
- **SC-005**: Existing installer identity, behavior, and eight canonical verification gates have zero regressions.

## Assumptions

- Fyne window dimensions are logical units and the existing canvas scale represents the effective display scale.
- The monitor nearest the launch pointer is the least surprising launch monitor before the native window has a stable placement.
- Native Explorer inspection can be recorded only on Windows with a compiled candidate MSI; source and script tests cannot substitute for it.
- Task-editor dialog sizing and saved window placement remain out of scope.

## Clarifications

### Session 2026-08-31

- Q: Which MSI Subject copy is approved? -> A: `go-schedule: cross-platform task scheduler`.
- Q: How is the launch monitor selected before the window is shown? -> A: Use the monitor nearest the launch pointer, then center the restored window normally.
- Q: Does an unavailable native candidate MSI block source implementation? -> A: No, but compiled and Explorer evidence must be reported honestly and cannot be inferred.
