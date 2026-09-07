# Implementation Plan: Operational Schedule and Activity

**Branch**: `codex/064-schedule-activity` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/064-schedule-activity/spec.md`

## Summary

Replace the Schedule and Activity placeholders with operational Wails workspaces that adapt the existing calendar, runs, logs, alerts, and acknowledgement APIs. A transport-neutral Go service will produce complete presentation snapshots and explicit typed records. React will provide agenda and calendar schedule views, a filterable activity table, stable selection, detailed diagnostics, safe alert acknowledgement, and non-destructive Clear View behavior.

## Technical Context

**Language/Version**: Go 1.25, TypeScript 5.9, React 19

**Primary Dependencies**: Existing local daemon client, Wails v2, React, Vite, Vitest, Testing Library, Playwright

**Storage**: Existing SQLite run and alert records plus the daemon in-memory log ring; no schema change

**Testing**: Go unit and race tests, Vitest and axe accessibility tests, Playwright Chromium interaction tests, canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop application

**Project Type**: Go desktop application with an embedded React frontend

**Performance Goals**: Filter, switch views, and select records in at least 100-row Schedule and Activity fixtures within two seconds

**Constraints**: Predictions and recorded runs remain distinct; exact log path only; complete snapshots only; local Clear View cutoff; no daemon API, persistence, retention, or scheduling changes

**Scale/Scope**: Two routes, one backend service package, five daemon reads, one acknowledgement mutation, and focused Go, React, accessibility, scale, and native-build coverage

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: The operations service owns adaptation and complete-snapshot composition, React owns presentation state, and daemon policy is not duplicated.
- **Testing**: Backend mapping, complete-snapshot failure, acknowledgement, request ordering, stable selection, filters, Clear View, accessibility, and 100-row scale tests are mandatory.
- **UX consistency**: Existing shell, controls, dialogs, notices, tables, focus behavior, and visual tokens are reused. State always has a non-color label.
- **Performance**: Each workspace uses bounded aggregate reads and local filtering. Event bursts are debounced and no per-row backend call is introduced.
- **Security and truth**: The exact daemon log path is displayed without probing. Output and attributes are rendered as text. Acknowledgement uses the existing daemon mutation.
- **Review workflow**: Work remains on `codex/064-schedule-activity`; the user explicitly authorized push, PR publication, review fixes, and at most one second Codex review round.
- **Pinned artifacts**: No pinned build or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/064-schedule-activity/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── desktop-operations-bridge.md
├── checklists/
│   ├── operational-ux.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
desktop/
├── app.go
├── main.go
├── operations/
│   ├── local.go
│   ├── model.go
│   ├── service.go
│   └── service_test.go
└── frontend/src/
    ├── App.tsx
    ├── App.test.tsx
    ├── styles.css
    └── operations/
        ├── ActivityPage.tsx
        ├── ActivityPage.test.tsx
        ├── SchedulePage.tsx
        ├── SchedulePage.test.tsx
        ├── bridge.ts
        ├── model.ts
        ├── store.ts
        └── store.test.ts
```

**Structure Decision**: Add one desktop feature package parallel to `desktop/automation` and one React feature folder parallel to `desktop/frontend/src/automation`. Keep the two routes separate while sharing contracts and refresh state. Adapt existing daemon clients through a narrow interface instead of expanding the public API.

## Complexity Tracking

No constitutional violations require justification.
