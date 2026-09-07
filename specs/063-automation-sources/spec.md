# Feature Specification: Connected Automation Sources

**Feature Branch**: `codex/063-automation-sources`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Automation-source administration, focused Go race and frontend suites, Chromium accessibility and scale coverage, native Windows Wails build, and canonical eight-gate verification passed 2026-09-07 on review branch `codex/063-automation-sources` for [#154](https://github.com/shruggietech/go-schedule/issues/154).

**Input**: GitHub issue [#154](https://github.com/shruggietech/go-schedule/issues/154), rebuild completion chains, external triggers, Trigger Sets, and filesystem watchers as connected automation sources.

## User Scenarios & Testing

### User Story 1 - Understand every automation source (Priority: P1)

A user opens one Automation Sources workspace and can understand how completion chains, external triggers, Trigger Sets, and filesystem watchers connect an initiating source to a target task, whether each source is ready, and what action will restore an unavailable source.

**Why this priority**: Reliable administration starts with an honest, shared mental model across all four source types.

**Independent Test**: Load mixed healthy, disabled, degraded, and disconnected sources and verify each row exposes source, target, readiness or health, and an actionable explanation without exposing a secret.

**Acceptance Scenarios**:

1. **Given** a collection containing every source type, **When** the workspace loads, **Then** all sources appear under clearly labeled sections using the same source-to-target language.
2. **Given** a disabled source, deleted target, or degraded watcher, **When** it is displayed, **Then** its state and recovery action are stated in text and are not conveyed by color alone.
3. **Given** at least 100 sources, **When** the user searches or filters the workspace, **Then** the matching set updates without losing keyboard focus and ordinary interaction completes within two seconds.

---

### User Story 2 - Manage completion chains and watchers (Priority: P2)

A user can create, inspect, edit, enable or disable where supported, retarget, and delete completion chains and filesystem watchers while preserving accurate relationship and watcher-health feedback.

**Why this priority**: These sources connect internal completion or filesystem events to task execution and require safe correction when paths or targets become invalid.

**Independent Test**: Create and edit one chain and watcher, change their targets, toggle the watcher, simulate degraded health, and delete both.

**Acceptance Scenarios**:

1. **Given** runnable tasks, **When** a completion chain is created or edited, **Then** source task, target task, and completion outcome are validated and the refreshed relationship is shown.
2. **Given** a filesystem watcher, **When** its path, selection rule, timing, target, or enabled state changes, **Then** the refreshed row reports the daemon-owned health state and reason.
3. **Given** an entity changed after the editor opened, **When** the user saves, **Then** the save is rejected as stale and offers reload or explicit overwrite.

---

### User Story 3 - Manage secret-bearing trigger sources (Priority: P3)

A user can create, inspect, edit, enable or disable, rotate, explicitly reveal and copy, fire, retarget, and delete external triggers and Trigger Sets without secrets leaking into ordinary workspace state.

**Why this priority**: Full lifecycle parity is required, and secret-bearing actions need a deliberately narrower security boundary.

**Independent Test**: Exercise each trigger and Trigger Set action, assert keys appear only in a user-invoked secret dialog, and verify workspace models, events, screenshots, and diagnostics contain no raw key.

**Acceptance Scenarios**:

1. **Given** a new trigger or Trigger Set, **When** creation succeeds, **Then** generated keys appear once in an explicit secret dialog and are absent after it closes.
2. **Given** an existing trigger or Trigger Set, **When** reveal or rotate is requested, **Then** keys are returned only to the explicit secret dialog and can be copied with confirmation.
3. **Given** an external trigger, **When** the user fires it, **Then** the key is used transiently by the backend and no secret enters the durable frontend workspace.
4. **Given** a Trigger Set, **When** the user retargets or toggles it, **Then** the operation remains atomic and ordered member identity is preserved.

### Edge Cases

- A source whose target task was deleted remains visible as unavailable and can be retargeted or deleted.
- A chain whose source task was deleted remains visible with a missing-source explanation.
- A watcher path may be long, temporarily unavailable, permission-denied, or platform-specific; the UI preserves the full value and shows the daemon health reason.
- Empty, whitespace-only, duplicate, cyclic, and same-source-and-target relationships are rejected using field-level guidance where the daemon supports it.
- A reconnect refresh cannot overwrite a newer user-requested refresh, duplicate a mutation, or discard an open draft.
- When the daemon disconnects after a successful load, the last loaded collection stays visible as read-only context and mutation controls are disabled.
- Partial list failure does not present a misleading partly-current workspace; the prior snapshot remains visible with a retryable error.
- Secret material is cleared when its dialog closes, the route changes, or a newer secret operation replaces it.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST provide one Automation Sources route containing completion chains, external triggers, Trigger Sets, and filesystem watchers.
- **FR-002**: Every source summary MUST identify its type, initiating source or selection rule, target task, enabled state where supported, readiness or health, and a plain-language reason.
- **FR-003**: Users MUST be able to search all source types and filter by source type and attention state without changing backend data.
- **FR-004**: Users MUST be able to create, inspect, update, retarget, and delete completion chains, external triggers, Trigger Sets, and filesystem watchers according to their existing daemon capabilities.
- **FR-005**: Users MUST be able to enable or disable external triggers, Trigger Sets, and filesystem watchers, rotate and reveal trigger secrets, copy revealed secrets, fire external triggers, and preserve ordered Trigger Set members.
- **FR-006**: Ordinary workspace responses, events, logging, diagnostics, and screenshots MUST NOT contain raw trigger keys.
- **FR-007**: Create, reveal, and rotate MUST return secrets only through an explicit user-requested result that the frontend discards when its dialog closes.
- **FR-008**: Firing a trigger from its identifier MUST use secret material transiently within the backend boundary and MUST return a secret-free result.
- **FR-009**: The workspace MUST show daemon-owned filesystem watcher health and reason without inferring healthy status from configuration alone.
- **FR-010**: Missing source or target relationships MUST remain visible with retarget or delete recovery actions where those actions are supported.
- **FR-011**: Editors MUST validate required names, targets, chain outcomes, watcher paths, watcher kinds, timing values, and Trigger Set members before submitting.
- **FR-012**: Updates MUST detect a changed or deleted entity and require reload or explicit overwrite before replacing newer data.
- **FR-013**: Mutations MUST be single-flight per entity, suppress duplicate submission, refresh the workspace on success, and preserve the current workspace on failure.
- **FR-014**: While disconnected, the last successful workspace MUST remain visible as read-only context and all mutation actions MUST be unavailable.
- **FR-015**: All status, selection, validation, secret, and confirmation experiences MUST be keyboard accessible, screen-reader meaningful, focus-managed, and understandable without color.
- **FR-016**: The implementation MUST retain current daemon and API semantics instead of duplicating scheduling, readiness, watcher-health, or secret policy in React.
- **FR-017**: Automated Go, React, accessibility, secret-boundary, stale-write, large-collection, and canonical repository verification MUST pass.
- **FR-018**: Issue [#154](https://github.com/shruggietech/go-schedule/issues/154) MUST remain traceable through the specification, tasks, change log, pull request, and verification record.

### Key Entities

- **Automation Workspace**: A secret-free snapshot of task choices and the four source collections, with a load timestamp.
- **Completion Chain**: A relationship from a source task completion outcome to a target task.
- **External Trigger**: A named, enabled or disabled external invocation source targeting one task; its raw key is excluded from ordinary state.
- **Trigger Set**: A named, ordered collection of trigger members sharing one target and atomic lifecycle operations.
- **Filesystem Watcher**: A path and selection rule targeting one task, plus daemon-owned health and recovery detail.
- **Secret Result**: Ephemeral keys produced only by create, reveal, or rotate and never embedded in an Automation Workspace.
- **Automation Draft**: Editable fields plus the original update timestamp used for stale-write detection.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A keyboard-only user can create, edit, enable or disable where supported, retarget, invoke where supported, and delete each source type through the Wails UI.
- **SC-002**: Automated serialization and UI tests find zero raw trigger keys in ordinary workspace models, connection events, diagnostics, and non-secret screens.
- **SC-003**: Healthy, disabled, missing-target, and degraded-watcher states are identified correctly in all tested fixtures and include a textual explanation.
- **SC-004**: Search and filtering over at least 100 mixed sources completes within two seconds in automated UI tests and retains focus.
- **SC-005**: Duplicate user activation produces at most one backend mutation per entity while it is pending.
- **SC-006**: All focused tests and `sh scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-07

- Q: Should the four legacy areas remain separate destinations? A: No. They share one Automation Sources route with type sections and a shared source-to-target model, while editors remain type-specific.
- Q: Where may raw trigger keys live? A: Only in a transient backend response to create, reveal, or rotate and in the resulting explicit secret dialog until it closes.
- Q: How should offline and partial-load behavior work? A: Preserve the last complete snapshot as read-only context, disable mutations, and never replace it with a partial collection.
- Q: How are concurrent edits handled? A: Compare the original update timestamp with the current entity, reject stale saves by default, and permit only an explicit overwrite retry.
- Q: Which component determines watcher health and relationship readiness? A: The daemon remains authoritative; the desktop layer translates its facts into safe presentation models and plain-language recovery guidance.

## Assumptions

- Existing local daemon APIs from slices 033, 054, 055, and 056 remain the source of truth and need no persistence migration.
- The Wails desktop shell and shared components from slices 060 through 062 are the target UI foundation.
- One hundred mixed sources is the practical interactive desktop collection baseline for this slice.
- Schedule, Activity, Options, broader settings, Fyne removal, and release cutover remain assigned to issues #155, #156, and #157.
- This slice fully resolves issue #154 and does not claim completion of parent issue #147.
