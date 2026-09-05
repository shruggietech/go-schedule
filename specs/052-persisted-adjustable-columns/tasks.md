# Tasks: Persisted Adjustable Columns

**Input**: Design documents from `specs/052-persisted-adjustable-columns/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required by FR-015 and constitution principle II. Each behavioral test must fail before its implementation is added.

## Phase 1: Setup and Baseline

**Purpose**: Confirm the existing structured-table and preference contracts.

- [x] T001 Record the clean baseline and existing structured-table test result in specs/052-persisted-adjustable-columns/verification.md
- [x] T002 Verify existing ignore files remain sufficient for the Go/Fyne project and record the result in specs/052-persisted-adjustable-columns/verification.md

---

## Phase 2: Foundational Preference and Allocation Model

**Purpose**: Create the validated, reusable sizing model required by both views.

- [x] T003 [P] Add failing serialization, normalization, invalid-value, version, and view-isolation tests in gui/column_preferences_test.go
- [x] T004 [P] Add failing profile allocation, conservation, adjacent-transfer, minimum, and narrow-width tests in gui/structured_table_test.go
- [x] T005 Implement versioned per-view column profiles and preference persistence in gui/column_preferences.go
- [x] T006 Extend responsive allocation and adjacent adjustment around a shared profile in gui/structured_table.go

**Checkpoint**: The shared sizing model validates, falls back, allocates, and persists independently.

---

## Phase 3: User Story 1 - Make Important Values Readable (Priority: P1) MVP

**Goal**: Resize every Schedule and Activity boundary by pointer or keyboard without overflow or regression.

**Independent Test**: Exercise all boundaries headlessly and verify adjacent-only movement, minimums, alignment, accessibility, and existing selection behavior.

- [x] T007 [US1] Add failing drag, keyboard, focus, cursor, accessible-label, refresh, and alignment tests in gui/structured_table_test.go
- [x] T008 [US1] Implement focusable draggable header boundary controls and shared layout refresh in gui/structured_table.go
- [x] T009 [US1] Increase practical default When proportions in gui/schedule.go and gui/logs.go

**Checkpoint**: Both tables expose one consistent accessible resize interaction.

---

## Phase 4: User Story 2 - Keep Each View's Layout (Priority: P1)

**Goal**: Restore valid per-view layouts across recreation and reject unsafe stored values.

**Independent Test**: Save distinct layouts, rebuild both views at multiple widths, and inject malformed and obsolete preferences.

- [x] T010 [US2] Add failing Schedule and Activity integration tests for recreation, view isolation, and fallback in gui/calendar_test.go and gui/logs_test.go
- [x] T011 [US2] Wire stable preference identities into Schedule and Activity construction in gui/schedule.go and gui/logs.go

**Checkpoint**: Each view durably restores only its own safe layout.

---

## Phase 5: User Story 3 - Restore Practical Defaults (Priority: P2)

**Goal**: Reset either table independently and persist the result.

**Independent Test**: Customize both tables, reset each one, and recreate the views.

- [x] T012 [US3] Add failing reset-isolation and persistence tests in gui/calendar_test.go and gui/logs_test.go
- [x] T013 [US3] Add per-view Reset columns actions in gui/schedule.go and gui/logs.go

**Checkpoint**: Reset is immediate, durable, and isolated.

---

## Phase 6: Polish and Verification

**Purpose**: Close traceability, documentation, and quality gates.

- [x] T014 [P] Update CHANGELOG.md and specs/README.md for S052 scope and decisions
- [x] T015 Run focused GUI tests, focused race tests, full tests, and canonical eight-gate verification; record exact evidence in specs/052-persisted-adjustable-columns/verification.md
- [x] T016 Audit FR-001 through FR-016 and SC-001 through SC-007, complete all task checkboxes, and set the S052 lifecycle evidence in specs/052-persisted-adjustable-columns/spec.md

---

## Dependencies and Execution Order

- Setup precedes all implementation.
- T003 and T004 may be authored independently, then T005 precedes T006.
- User Story 1 depends on the shared model and is the independently testable MVP.
- User Story 2 depends on User Story 1's adjustable table constructor.
- User Story 3 depends on persistence wiring from User Story 2.
- Documentation can proceed separately; final verification depends on all code.

## Parallel Example: Foundation

```text
T003: Preference serialization and validation tests in gui/column_preferences_test.go
T004: Allocation and adjustment tests in gui/structured_table_test.go
```

## Implementation Strategy

1. Establish the failing preference and allocator contracts.
2. Deliver the adjustable shared table as the MVP.
3. Wire independent persistence into both views.
4. Add isolated reset actions.
5. Run the complete traceability and verification gate before committing.
