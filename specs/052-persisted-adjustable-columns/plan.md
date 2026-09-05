# Implementation Plan: Persisted Adjustable Columns

**Branch**: `codex/052-persisted-adjustable-columns` | **Date**: 2026-09-05 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/052-persisted-adjustable-columns/spec.md`

## Summary

Extend the responsive structured table with shared, focusable boundary controls. Persist versioned normalized proportions independently for Schedule and Activity through current-user preferences, derive constrained widths from current space, and preserve the fixed header, virtualized list, and no-horizontal-scroll model.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library and existing Fyne v2.7.4

**Storage**: Existing Fyne current-user preferences; no daemon or SQLite change

**Testing**: Go unit tests, Fyne headless widget tests, race detector, coverage gate, golangci-lint, documentation and automation checks

**Target Platform**: Windows, Linux, and macOS Fyne desktop GUI

**Project Type**: Single Go module with desktop GUI plus daemon and CLI

**Performance Goals**: Resizing remains constant-time in four-column views and does not alter scheduler dispatch performance

**Constraints**: No new dependency, daemon storage, or horizontal scrolling; preserve virtualization and stable identity

**Scale/Scope**: Two four-column views, six boundary controls, two preference keys, shared infrastructure, and focused GUI tests

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Code Quality**: PASS. Allocation, validation, serialization, and interaction remain focused GUI responsibilities with safe fallback.
- **II. Testing Standards**: PASS. Tests precede implementation and cover allocation, persistence, malformed input, isolation, reset, keyboard, drag, alignment, and regression. Canonical verification remains mandatory.
- **III. User Experience Consistency**: PASS. Both views share accessible controls, predictable key steps, reset terminology, and responsive behavior.
- **IV. Performance Requirements**: PASS. Calculation and adjustment are O(column count), with four columns per view. No benchmark is warranted.
- **V. Autonomous Build-Phase Execution**: PASS. S052 is traceable to #119, uses a `codex/` branch, and has explicit PR publication authorization.
- **Engineering constraints**: PASS. State is current-user presentation data; no dependency, daemon schema, IPC, or platform support changes.

## Phase 0 Research Decisions

1. Persist versioned JSON with ordered column identities and normalized proportions. Reject the entire value on any validation failure.
2. Transfer proportions only between adjacent columns. Clamp to logical minimums when space permits and retain proportional compression below them.
3. Put a visible, focusable boundary widget over each header gap with a resize cursor, accessible label, drag handling, and arrow-key handling.
4. Drive header and virtual row layouts from one shared profile and refresh without recreating row data or identity.
5. Put **Reset columns** beside each affected view's existing controls.

See [research.md](research.md) for alternatives and detailed rationale.

## Phase 1 Design

- Add an immutable serialized record and validated in-memory column profile.
- Extend allocation to accept profile proportions while retaining minimums and exact width conservation.
- Refresh rows during changes and persist only completed interactions.
- Give Schedule and Activity stable, independent preference identities.
- Keep Tasks on its fixed profile for backward compatibility.

### Post-design constitution re-check

All gates remain PASS. The design adds no dependency, API, migration, or scheduler work. Corrupt preferences are non-fatal and all new interaction has deterministic headless coverage.

## Project Structure

### Documentation (this feature)

```text
specs/052-persisted-adjustable-columns/
├── spec.md / plan.md / research.md / data-model.md / quickstart.md
├── contracts/column-layout-contract.md
├── checklists/requirements.md / ux.md
└── tasks.md / verification.md
```

### Source Code (repository root)

```text
gui/
├── schedule.go / calendar_test.go
├── logs.go / logs_test.go
├── structured_table.go / structured_table_test.go
└── column_preferences.go / column_preferences_test.go
```

**Structure Decision**: Extend `gui` because this is desktop-only presentation state. Keep serialization focused while shared table code owns allocation and widgets. No new package is needed.

## Complexity Tracking

No constitution violation or exceptional complexity is required.
