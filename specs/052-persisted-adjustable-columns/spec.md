# Feature Specification: Persisted Adjustable Columns

**Feature Branch**: `codex/052-persisted-adjustable-columns`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/052-persisted-adjustable-columns`; focused,
race, full GUI, light/dark appearance, and canonical eight-gate verification
passed 2026-09-05 for issue #119

**Input**: GitHub issue [#119](https://github.com/shruggietech/go-schedule/issues/119), Add persisted adjustable columns to Schedule and Activity.

## Problem Statement

The fixed Schedule and Activity table allocations can hide most of a timestamp,
especially in the **When** column. Operators need to adapt each table to their
work without losing those choices when the application, display scale, font, or
window size changes. The interaction must remain usable without requiring a
horizontal scrollbar or sacrificing full-value disclosure and keyboard access.

## Scope

### In scope

- Independently adjustable columns in the Schedule list and Activity table.
- Pointer and keyboard resizing through discoverable header controls.
- Per-user, per-view persistence using normalized proportions.
- Safe fallback and clamping for missing, malformed, obsolete, or impossible
  stored values.
- A per-view action that restores default proportions.
- A more useful default allocation for each **When** column.
- Preservation of frozen headers, stable row identity, full-value disclosure,
  alternating rows, themes, fonts, and responsive narrow-window behavior.
- Deterministic headless tests plus canonical repository verification.

### Out of scope

- Adjustable columns in Tasks, Groups, Chains, Options, or Calendar mode.
- Reordering, hiding, sorting, or adding columns.
- Daemon-owned or machine-wide storage of desktop layout preferences.
- Horizontal table scrolling.
- Redesigning table data, time formats, or row selection behavior.

## Clarifications

### Session 2026-09-05

- Q: What width representation best survives DPI, font, and window changes? → A: Persist normalized proportions for each view and derive responsive widths from current available space.
- Q: How should keyboard and pointer users resize columns? → A: Give every boundary a focusable header control that supports dragging plus left/right keyboard adjustment and exposes an accessible label.
- Q: How should invalid or obsolete preferences recover? → A: Validate the complete versioned value, clamp usable proportions, and fall back atomically to that view's defaults when the value cannot be recovered safely.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Make Important Values Readable (Priority: P1)

An operator widens a Schedule or Activity column from the frozen header and the
adjacent column contracts, keeping the table within the available view.

**Why this priority**: The feature exists to let operators see the fields that
matter without changing the window or accepting permanent truncation.

**Independent Test**: Resize every boundary by pointer and keyboard in each
view, then verify alignment, minimum usability, and full-value disclosure at
the supported default and narrow window sizes.

**Acceptance Scenarios**:

1. **Given** a Schedule or Activity list, **when** a boundary is dragged, **then** the two adjacent columns resize and all other boundaries stay fixed.
2. **Given** a boundary has keyboard focus, **when** left or right is pressed, **then** the same adjacent pair changes by a predictable step.
3. **Given** either adjacent column reaches its protected minimum, **when** resizing continues toward it, **then** the boundary stops and no content overlaps.
4. **Given** a resized table becomes narrower, **when** available width cannot satisfy every minimum, **then** all columns scale proportionally without horizontal scrolling.

---

### User Story 2 - Keep Each View's Layout (Priority: P1)

An operator's Schedule and Activity layouts persist independently across an
application restart and adapt to a different window size, display scale, or
font.

**Why this priority**: Adjustment would become repeated busywork if it were not
durable, while raw physical widths would become invalid across environments.

**Independent Test**: Save distinct layouts for both views, recreate the UI,
and verify each layout restores as proportions under changed available widths
and preferences.

**Acceptance Scenarios**:

1. **Given** distinct Schedule and Activity adjustments, **when** the application is reopened, **then** each view restores only its own layout.
2. **Given** a stored layout, **when** the window, display scale, or font changes, **then** its relative intent is preserved within current constraints.
3. **Given** missing, malformed, non-finite, wrong-length, or obsolete stored data, **when** a view opens, **then** it uses safe defaults without an error dialog.

---

### User Story 3 - Restore Practical Defaults (Priority: P2)

An operator can restore one table to practical default proportions without
changing the other table.

**Why this priority**: A reversible customization prevents users from becoming
stuck with an accidental or no-longer-useful layout.

**Independent Test**: Customize both views, reset one, and verify only that
view returns to defaults now and after restart.

**Acceptance Scenarios**:

1. **Given** a customized view, **when** its restore action is invoked, **then** default proportions apply immediately and persist.
2. **Given** both views are customized, **when** one is restored, **then** the other view remains unchanged.
3. **Given** no customization, **when** either list first opens, **then** **When** receives enough default width to expose a practical timestamp portion at the supported default window size.

### Edge Cases

- A drag ends outside the table or at a negative coordinate.
- Available width is zero or smaller than the sum of protected minimums.
- A stored value is empty, truncated, wrong-version, wrong-column-count,
  nonnumeric, non-finite, non-positive, or extremely disproportionate.
- A future release adds, removes, or renames columns.
- The view is rebuilt after a theme or font change while preferences remain.
- Repeated keyboard adjustment reaches a minimum and must stop deterministically.
- Reset is invoked while a boundary has focus.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Schedule list and Activity MUST expose an adjustable control at every boundary between adjacent columns.
- **FR-002**: Each boundary MUST support pointer dragging and focused left/right keyboard adjustment with a consistent logical step.
- **FR-003**: Each boundary MUST expose a meaningful accessible label naming both adjacent columns and MUST be discoverable without relying on color alone.
- **FR-004**: Resizing MUST transfer width only between the two adjacent columns and MUST preserve the table's total usable width.
- **FR-005**: Resizing MUST protect defined per-column minimums when the available width can satisfy them and MUST stop cleanly at those limits.
- **FR-006**: When available width cannot satisfy all minimums, columns MUST scale proportionally without overlap or mandatory horizontal scrolling.
- **FR-007**: Completed adjustments MUST be persisted as normalized proportions in current-user desktop preferences.
- **FR-008**: Schedule and Activity MUST use independent, versioned preference identities and MUST NOT affect one another.
- **FR-009**: Restored proportions MUST be normalized and safely constrained for the current columns and available width.
- **FR-010**: Missing, malformed, non-finite, non-positive, wrong-length, obsolete, or otherwise unrecoverable preference data MUST fall back atomically to that view's defaults.
- **FR-011**: Each affected view MUST provide a clearly labeled action that restores and persists only that view's default proportions.
- **FR-012**: Default proportions MUST allocate a more practical share to **When** than the pre-S052 layout while preserving useful space for all remaining columns.
- **FR-013**: Header and body cells MUST remain aligned during and after resize, reset, data refresh, theme/font change, and application recreation.
- **FR-014**: Existing fixed-header behavior, row virtualization, stable selection identity, activation, alternating rows, keyboard list navigation, and selectable full-value disclosure MUST remain intact.
- **FR-015**: Automated tests MUST cover allocation conservation, adjacent-only resizing, bounds, persistence, per-view isolation, invalid-value fallback, schema-version fallback, reset, pointer and keyboard behavior, and rebuilt-view alignment.
- **FR-016**: The complete change MUST pass the canonical format, vet, lint, race, GUI, coverage, documentation, and automation gates.

### Key Entities

- **Column layout profile**: One view's schema version, ordered column identity,
  normalized positive proportions, and defaults.
- **Column boundary control**: The accessible interaction between two adjacent
  columns that transfers width within constraints.
- **View layout preference**: The current user's persisted Schedule or Activity
  profile, isolated by a stable preference identity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every boundary in both affected views can be adjusted by pointer and keyboard while total allocated width differs from available usable width by no more than one display unit.
- **SC-002**: One hundred percent of tested valid custom layouts restore independently after UI recreation and remain usable at three representative window widths.
- **SC-003**: One hundred percent of malformed, obsolete, non-finite, and impossible preference fixtures recover to a usable default without an error dialog or inaccessible column.
- **SC-004**: Resetting either view restores its documented defaults immediately and after recreation while leaving the other view unchanged.
- **SC-005**: At the supported default window size, the default **When** allocation displays a practical timestamp portion and is larger than its pre-S052 allocation.
- **SC-006**: Existing structured-table interaction and presentation regression tests remain green in both appearance modes.
- **SC-007**: The canonical eight-gate repository verification completes successfully with all core packages at or above 80 percent coverage.

## Assumptions

- The existing per-user desktop preference store is available before the views
  are constructed and is the proper owner for presentation-only state.
- Normalized proportions express user intent more robustly than device pixels.
- A boundary resize changes only its two neighbors, matching common table
  behavior and keeping unrelated columns stable.
- A small header action labeled **Reset columns** is sufficiently discoverable
  and does not require a new Options-page section.
- Tasks table adjustment remains separately deferrable because issue #119 names
  only Schedule and Activity.
