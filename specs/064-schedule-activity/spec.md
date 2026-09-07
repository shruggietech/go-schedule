# Feature Specification: Operational Schedule and Activity

**Feature Branch**: `codex/064-schedule-activity`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Operational Schedule and Activity workspaces, focused Go race and frontend suites, Chromium accessibility and scale coverage, native Windows Wails build, and canonical eight-gate verification passed 2026-09-07 on review branch `codex/064-schedule-activity` for [#155](https://github.com/shruggietech/go-schedule/issues/155).

**Input**: GitHub issue [#155](https://github.com/shruggietech/go-schedule/issues/155), rebuild Schedule and Activity around fast operational understanding.

## User Scenarios & Testing

### User Story 1 - Understand scheduled work (Priority: P1)

A user opens Schedule and can distinguish future predictions from completed execution records, choose a useful time range, switch between agenda and calendar views, and inspect one occurrence without losing context.

**Why this priority**: Upcoming work is the primary operational question, and confusing predictions with recorded execution creates false conclusions.

**Independent Test**: Load mixed future predictions and completed outcomes, switch ranges and views, select records, and verify every state is named in text or symbol plus text.

**Acceptance Scenarios**:

1. **Given** predicted and completed occurrences in one range, **When** Schedule loads, **Then** each row identifies whether it is a prediction or recorded run and exposes task, time, state, and outcome where applicable.
2. **Given** the agenda and calendar views, **When** a user changes view or range, **Then** the selected range is preserved and the resulting window is stated explicitly.
3. **Given** a selected occurrence, **When** live data refreshes, **Then** selection and keyboard focus remain on that identity when it still exists, otherwise the UI explains that it is no longer in the selected range.

---

### User Story 2 - Investigate activity and failures (Priority: P1)

A user opens Activity and can scan recent run records, daemon logs, and scheduler alerts without mistaking one for another, filter the feed, and inspect complete failure diagnostics.

**Why this priority**: Operators need to move from a visible failure to its persisted evidence quickly and accurately.

**Independent Test**: Load runs, logs, and alerts with mixed severities and outcomes, filter them, open each detail type, and verify failure output and provenance remain available.

**Acceptance Scenarios**:

1. **Given** running, successful, failed, skipped, caught-up, queued, and unavailable records, **When** Activity loads, **Then** each state is distinguishable without relying on color.
2. **Given** daemon logs, alerts, and run records, **When** the user filters by record type, severity, outcome, or text, **Then** only matching records remain and the result count is announced.
3. **Given** a failed run with captured output and provenance, **When** its detail opens, **Then** the user can inspect task identity, trigger, timestamps, exit status, truncation, output, and source identifiers.

---

### User Story 3 - Control alerts and the current view safely (Priority: P2)

A user can acknowledge an alert or clear the currently visible Activity view while understanding that persisted run records and daemon logs are not deleted.

**Why this priority**: Operational cleanup must be useful without implying or causing data loss.

**Independent Test**: Acknowledge one alert, clear a filtered view, verify visible alerts are acknowledged, existing rows are hidden locally, and later activity appears.

**Acceptance Scenarios**:

1. **Given** an unacknowledged alert, **When** the user acknowledges it, **Then** the refreshed record is marked acknowledged and no unrelated alert changes.
2. **Given** a filtered Activity view, **When** the user activates Clear View, **Then** visible rows at or before that moment are hidden locally, visible alerts are acknowledged, and text states that records are not deleted.
3. **Given** a cleared view, **When** new activity arrives later, **Then** it appears normally.

### Edge Cases

- The daemon returns an empty range, an empty activity collection, or an exact range boundary.
- A future prediction has no run identifier or outcome; a persisted run always remains labeled as recorded activity.
- An occurrence or activity item disappears between selection and refresh.
- Multiple records have matching timestamps, messages, or task names but different stable identities.
- A task referenced by a historical run or alert has been deleted.
- A run has no exit code because launch or setup failed, empty output, or truncated output.
- The log path is unavailable until the daemon responds and must never be guessed.
- One constituent read fails after a previous successful load.
- At least 100 schedule rows and at least 100 activity rows are loaded.
- A rapid event burst causes several refresh requests to overlap.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST provide dedicated Schedule and Activity routes in the desktop application.
- **FR-002**: Schedule MUST offer agenda and calendar views over user-selectable 1-day, 7-day, and 30-day windows.
- **FR-003**: Schedule MUST label computed future occurrences as predictions and persisted occurrences as recorded runs.
- **FR-004**: Schedule MUST display stable task identity, task name, occurrence time, record kind, and outcome where available.
- **FR-005**: Calendar view MUST expose occurrence counts per local day and an accessible list for the selected day.
- **FR-006**: Activity MUST present run records, daemon log records, and alerts as distinct record types in one recent operational workspace.
- **FR-007**: Activity MUST support text search plus record-type, severity, and outcome filters without changing backend data.
- **FR-008**: Running, success, failure, skipped, caught-up, queued, predicted, acknowledged, unacknowledged, and unavailable states MUST use explicit text or symbol plus text and MUST NOT rely on color alone.
- **FR-009**: Activity detail MUST preserve complete available run diagnostics, including trigger and source provenance, timestamps, exit status, output truncation, and retained output.
- **FR-010**: Activity MUST show the exact daemon-reported log path when available and state when it is unavailable; the desktop MUST NOT invent, probe, or open a fallback path.
- **FR-011**: Users MUST be able to acknowledge one alert through the existing daemon behavior.
- **FR-012**: Clear View MUST hide currently visible records at or before its activation time, acknowledge visible alerts, preserve persisted data, and allow later activity to appear.
- **FR-013**: Live run, log, alert, and task events MUST refresh relevant data through a bounded debounce and request-ordering guard.
- **FR-014**: Refresh MUST preserve selected range, view, filters, keyboard focus, and selected stable identity when that identity remains present.
- **FR-015**: A failed or partial refresh MUST preserve the last complete snapshot with a retryable explanation and MUST NOT present partial data as current.
- **FR-016**: While disconnected, the last complete snapshot MUST remain visible as read-only context and alert mutations MUST be disabled.
- **FR-017**: Schedule and Activity controls, filters, tables, calendar days, details, status announcements, and empty states MUST be keyboard accessible, screen-reader meaningful, and understandable without color.
- **FR-018**: Schedule and Activity MUST remain responsive and readable with at least 100 rows in each workspace, with ordinary filtering and selection completing within two seconds.
- **FR-019**: The daemon remains authoritative for occurrences, run outcomes, logs, alerts, acknowledgement, and log-path metadata; the desktop MUST NOT duplicate scheduler or retention policy.
- **FR-020**: Automated Go, React, accessibility, stable-selection, request-ordering, large-collection, and canonical repository verification MUST pass.
- **FR-021**: Issue [#155](https://github.com/shruggietech/go-schedule/issues/155) MUST remain traceable through the specification, tasks, change log, pull request, and verification record.

### Key Entities

- **Schedule Window**: The requested from and to instants plus the complete occurrence collection returned for that window.
- **Schedule Occurrence**: A stable presentation record for either a future prediction or a persisted run occurrence.
- **Activity Workspace**: One complete snapshot of recent runs, daemon logs, alerts, exact log-path metadata, and load time.
- **Activity Record**: A typed run, log, or alert summary with stable identity, timestamp, state, searchable text, and detail payload.
- **Activity Detail**: Full available diagnostic fields for one selected record without inventing missing data.
- **View Cutoff**: A local timestamp that hides records currently visible when Clear View is activated.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A keyboard-only user can choose every Schedule range and view, select an occurrence, filter Activity, inspect each record type, acknowledge an alert, and clear the current view.
- **SC-002**: Every tested operational state is conveyed by text or symbol plus text in both light and dark appearances.
- **SC-003**: Automated tests prove that overlapping refreshes cannot replace newer data or move selection when the selected identity remains available.
- **SC-004**: Search, filtering, view switching, and selection over fixtures with at least 100 schedule rows and 100 activity rows complete within two seconds.
- **SC-005**: Automated tests prove a failed constituent read preserves the previous complete snapshot and exposes no partly-current workspace.
- **SC-006**: Focused tests and `sh scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-07

- Q: Should Schedule and Activity be merged into one route? A: No. They answer different questions and remain dedicated routes, but share one transport-neutral operations service and consistent refresh behavior.
- Q: How are predictions distinguished from history? A: Future calendar calculations are explicitly labeled predictions; records with run identity are labeled recorded runs and never blended into prediction semantics.
- Q: What does Activity combine? A: It presents three visibly distinct sections or record types: persisted runs, recent daemon logs, and alerts. Filters may combine them, but type remains explicit.
- Q: What survives live refresh? A: Range, view, filters, keyboard focus, and selected identity survive when valid; an absent selected record produces an explanatory empty detail instead of silently selecting another.
- Q: What does Clear View mutate? A: It stores only a local cutoff and acknowledges alerts visible in the filtered result. It never deletes runs, logs, or alerts.
- Q: What happens on partial failure or disconnection? A: The last complete snapshot remains visible as read-only context; no partial replacement is published and mutations are disabled.

## Assumptions

- Existing calendar, runs, logs, alerts, acknowledgement, runtime metadata, and event-stream APIs already provide the required source-of-truth data.
- The Wails desktop shell and shared components from slices 060 through 063 are the target UI foundation.
- The Activity workspace is a bounded recent operational view; full daemon history remains in the daemon-reported log file.
- One hundred rows in each workspace is the practical desktop interaction baseline for this slice.
- Connections, Settings, legacy Fyne removal, and release cutover remain assigned to issues #156 and #157.
- This slice fully resolves issue #155 and does not claim completion of parent issue #147.
