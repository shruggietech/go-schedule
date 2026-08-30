# Feature Specification: GUI Activity Clarity

**Feature Branch**: `codex/014-gui-activity-clarity`

**Created**: 2026-08-27

**Status**: Implemented

**Delivery**: [PR #46](https://github.com/shruggietech/go-schedule/pull/46)

**Traceability**: Closes GitHub issues
[#26](https://github.com/shruggietech/go-schedule/issues/26),
[#28](https://github.com/shruggietech/go-schedule/issues/28), and
[#30](https://github.com/shruggietech/go-schedule/issues/30).

## Overview

The GUI combines daemon log records and scheduler alerts in a tab named
"Logs." Its alert badge can grow without bound, and its "Dismiss All" action
uses a trash icon even though it deletes no persisted records. This slice makes
that mixed activity surface accurate at a glance and clear about its temporary,
non-destructive action.

### Scope in

- Use "Activity" as the user-facing name for the mixed log and alert view.
- Cap the unacknowledged-alert badge at `99+`.
- Replace destructive-looking clear language and imagery with accurate copy.
- Keep an explanation of the clear action visible without requiring hover.
- Preserve the existing alert acknowledgement and current-view cutoff behavior.

### Scope out

- Exposing the daemon log-file path or adding API fields (issue #31).
- Changing log retention, persistence, filtering, or collection.
- Deleting log records or alerts.
- Changing daemon, CLI, scheduler, or notification behavior.
- Adding tooltip infrastructure solely for this action.

## User Story 1 - Understand the mixed activity view (Priority: P1)

As a GUI user, I want the tab to identify itself as Activity so I understand
that it contains both log records and scheduler alerts.

**Why this priority**: The current label misstates the contents of a primary
navigation destination.

**Independent Test**: Open the GUI and confirm the fourth tab is named
"Activity" both with and without an alert-count badge.

### Acceptance Scenarios

1. **Given** no unacknowledged alerts, **when** the GUI is displayed, **then**
   the mixed log and alert tab is labeled `Activity`.
2. **Given** one or more unacknowledged alerts, **when** the badge is updated,
   **then** the count is appended to `Activity` rather than `Logs`.

---

## User Story 2 - Read a bounded alert badge (Priority: P1)

As a GUI user, I want a compact badge so a large alert backlog does not produce
an unwieldy tab label.

**Why this priority**: The badge shares limited navigation space with every
other tab and must remain stable at high counts.

**Independent Test**: Evaluate the displayed label at zero, one, 99, and 100
unacknowledged alerts.

### Acceptance Scenarios

1. **Given** between one and 99 unacknowledged alerts, **when** the badge is
   updated, **then** the exact count is displayed.
2. **Given** 100 or more unacknowledged alerts, **when** the badge is updated,
   **then** the label displays `Activity (99+)`.
3. **Given** no unacknowledged alerts, **when** the badge is updated, **then**
   no count is displayed.

---

## User Story 3 - Clear the current view without fearing data loss (Priority: P1)

As a GUI user, I want the clearing action to explain its effect so I can hide
current activity without believing records will be deleted.

**Why this priority**: Destructive wording and a trash icon communicate a false
data-loss risk for a routine view action.

**Independent Test**: Inspect and invoke the Activity clear control, then
confirm its visible explanation, current rows disappear, visible alerts are
acknowledged, and persisted records are not deleted.

### Acceptance Scenarios

1. **Given** current activity rows, **when** the Activity view is displayed,
   **then** the action uses non-destructive wording and a clear-content icon.
2. **Given** the action is available, **when** a user reads the view, **then**
   visible text explains that current activity is hidden, visible alerts are
   acknowledged, and records are not deleted.
3. **Given** current activity rows, **when** the user activates the action,
   **then** rows at or before that moment are hidden for the current view and
   visible alerts are acknowledged through the existing backend behavior.
4. **Given** new activity arrives after the action, **when** the view refreshes,
   **then** the new activity appears.

### Edge Cases

- A negative or otherwise invalid badge count is treated as no alerts.
- Exactly 99 alerts display the exact value; 100 begins the capped display.
- Clearing an empty or filtered view remains safe and does not delete data.
- The explanation remains available on devices or environments without hover.

## Requirements

### Functional Requirements

- **FR-001**: The mixed daemon-log and scheduler-alert tab MUST use `Activity`
  as its user-facing base label.
- **FR-002**: The tab badge MUST display exact unacknowledged-alert counts from
  1 through 99, `99+` for counts of 100 or more, and no badge for zero or
  invalid negative counts.
- **FR-003**: The current-view action MUST use non-destructive wording and MUST
  NOT use a delete or trash icon.
- **FR-004**: The Activity view MUST visibly explain that the action hides
  current activity, acknowledges visible alerts, and does not delete records.
- **FR-005**: Activating the action MUST preserve the existing behavior of
  hiding entries at or before the action time and acknowledging alerts visible
  in the current filtered view.
- **FR-006**: Activity created after the action time MUST remain visible.
- **FR-007**: The feature MUST NOT change persistence, API contracts, daemon
  behavior, log retention, or alert semantics.
- **FR-008**: Automated tests MUST cover the canonical tab label, badge
  boundaries, explanatory copy, and non-destructive control presentation.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All user-facing references for this GUI destination use
  `Activity`; no visible `Logs` tab label remains.
- **SC-002**: Automated tests pass for badge counts 0, 1, 99, 100, and a count
  above 100.
- **SC-003**: The Activity view presents a clear-content action and an always
  visible explanation that explicitly says records are not deleted.
- **SC-004**: Existing merge, sort, filter, acknowledgement, and post-clear
  activity behavior remains covered and passes the full local verification
  aggregate.

## Assumptions

- "Activity" is the canonical user-facing term; internal function and field
  names may remain unchanged when renaming them would add churn without user
  value.
- Visible explanatory copy is preferable to hover-only help because the GUI may
  run in environments without hover and Fyne has no existing tooltip pattern in
  this project.
- The existing current-view cutoff and alert acknowledgement are intentional;
  this slice clarifies them rather than redesigning them.
