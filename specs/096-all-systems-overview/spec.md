# Feature Specification: All Systems Operational Overview

**Feature Branch**: `codex/096-all-systems-overview`

**Created**: 2026-09-21

**Status**: Implemented

**Delivery**: review branch `codex/096-all-systems-overview`

**Input**: Work slice S096, complete issue #182 by adding a bounded, read-only All Systems overview for every registered daemon without introducing clustered execution or shared task ownership.

## Clarifications

### Session 2026-09-21

- Q: Which daemons belong in All Systems? -> A: This computer and every saved desktop connection profile appear exactly once, including profiles that cannot currently connect.
- Q: What remains visible after a target refresh fails? -> A: Successful summaries remain available only for the current desktop session and are marked stale with their observation time; no cross-session operational cache is introduced.
- Q: How is fan-out bounded? -> A: At most four targets refresh concurrently, each target receives a five-second deadline, and one target failure never cancels another target result.
- Q: What does drill-down do? -> A: It selects exactly one identified daemon and opens the relevant Tasks, Schedule, Activity, or Notifications destination with the source record or task context preserved when available.
- Q: How are operational counts bounded? -> A: The summary uses the latest 24 hours for recent failures and notification problems, the next 24 hours for upcoming work, and bounded representative records rather than transferring complete histories.

## User Scenarios & Testing

### User Story 1 - Understand every daemon at a glance (Priority: P1)

As an operator, I can open one All Systems destination and immediately distinguish healthy, stale, disconnected, unauthorized, incompatible, and not-yet-contacted daemons.

**Why this priority**: A trustworthy fleet overview is the foundational outcome and must be useful before cross-daemon actions are introduced.

**Independent Test**: Register a mixture of reachable and unreachable profiles, open All Systems, and confirm each saved target plus This computer appears once with explicit identity, state, and observation time.

**Acceptance Scenarios**:

1. **Given** This computer and several saved profiles, **When** All Systems loads, **Then** every target appears exactly once with daemon identity, connection state, version, platform, and freshness.
2. **Given** two targets with the same display name, **When** they appear together, **Then** endpoint or local designation and shortened daemon identity keep them distinguishable.
3. **Given** one target times out, **When** other targets finish refreshing, **Then** their current summaries appear without waiting for or being replaced by the failed target.
4. **Given** no remote profiles, **When** All Systems loads, **Then** This computer remains a useful one-system overview and the empty guidance explains how additional systems are registered.

---

### User Story 2 - Triage operational attention (Priority: P1)

As an operator, I can compare active tasks, upcoming work, recent failures, unacknowledged alerts, and notification problems without opening every daemon individually.

**Why this priority**: Connection health alone does not identify which daemon needs attention or why.

**Independent Test**: Seed deterministic daemon summaries with upcoming work, failures, alerts, and delivery problems, then sort and filter the overview to isolate each attention category.

**Acceptance Scenarios**:

1. **Given** mixed healthy and attention-needed systems, **When** I filter to systems needing attention, **Then** only targets with connection problems, failures, alerts, or notification problems remain.
2. **Given** a reachable daemon, **When** its summary is current, **Then** active-task count, nearest upcoming work, recent-failure count, unacknowledged-alert count, and notification-problem count identify that daemon and observation time.
3. **Given** large or mixed-version profile collections, **When** results arrive in any order, **Then** stable sorting and bounded representative details keep the overview usable.
4. **Given** a previously successful target now fails to refresh, **When** session data exists, **Then** its last successful summary remains visible as stale and is never presented as current.

---

### User Story 3 - Drill into one safe target (Priority: P1)

As an operator, I can open the relevant destination for one summary item while keeping the daemon identity explicit and unchanged.

**Why this priority**: Aggregate information is actionable only when it leads to the correct source without creating wrong-target risk.

**Independent Test**: Open task, schedule, failure, alert, and notification-problem drill-downs from same-named daemons and confirm the intended profile is selected before the relevant source context is displayed.

**Acceptance Scenarios**:

1. **Given** a target summary, **When** I open Tasks, **Then** that exact daemon becomes the selected target before the Tasks destination loads.
2. **Given** a representative upcoming occurrence, failure, alert, or notification problem, **When** I open it, **Then** the relevant destination receives its task or record context and never infers a target from the display name.
3. **Given** a target that cannot be selected because its credential, trust, or connection state changed, **When** I request drill-down, **Then** the overview remains in place and explains the target-specific recovery action.
4. **Given** stale summary data, **When** I request drill-down, **Then** the desktop reconnects and re-establishes current target identity before showing mutable destination controls.

---

### User Story 4 - Refresh and navigate accessibly (Priority: P2)

As a keyboard or assistive-technology user, I can refresh, sort, filter, inspect, and leave the overview at supported window sizes without losing context.

**Why this priority**: Multi-system density must not recreate the cramped or inaccessible administration surfaces previously corrected.

**Independent Test**: Exercise the overview by keyboard and screen reader at 800 by 600, 200 percent zoom, reduced motion, and a large target set while refreshes are active and partially failing.

**Acceptance Scenarios**:

1. **Given** a refresh in progress, **When** results complete independently, **Then** progress and result changes are announced without moving keyboard focus.
2. **Given** 100 saved profiles, **When** I filter or sort, **Then** controls remain reachable, rows remain identifiable, and the page introduces no document-level horizontal scrolling.
3. **Given** reduced motion or high zoom, **When** state changes, **Then** status remains understandable without animation or color-only meaning.

### Edge Cases

- A saved profile duplicates another profile's label or references the same daemon identity through a different profile; each profile remains an independently identified registration and no results are silently merged.
- The local daemon is unavailable while remote profiles remain reachable; the overview still returns remote results and identifies This computer separately.
- A profile is removed during refresh; its late result is discarded rather than restoring a deleted registration.
- The user initiates another refresh before the current refresh completes; the older generation is canceled and cannot overwrite newer results.
- A target returns malformed, unsupported, or future-version summary data; it is marked incompatible without blocking compatible targets.
- A target has no tasks, upcoming work, failures, alerts, or notification deliveries; zero values are explicit and do not resemble missing data.
- A target has more matching records than the summary limit; counts remain authoritative while representative records remain bounded.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST provide an All Systems destination that includes This computer and every saved desktop connection profile exactly once.
- **FR-002**: Every system entry MUST identify its registration, daemon identity when known, display label, local or remote location, platform, version, connection state, and observation time without exposing credentials or certificate contents.
- **FR-003**: The overview MUST distinguish healthy, stale, unavailable, timed out, unauthorized, forbidden, revoked, trust-changed, identity-changed, incompatible, and not-yet-contacted states using text and shape in addition to color.
- **FR-004**: Refresh MUST run independently per target with no more than four target operations in flight and a five-second deadline per target.
- **FR-005**: Starting a newer refresh MUST cancel the older refresh generation, and late older results MUST NOT replace newer results or reintroduce removed profiles.
- **FR-006**: A successful summary MUST report observation time, active-task count, nearest upcoming scheduled work in the next 24 hours, recent failed-run count in the previous 24 hours, unacknowledged-alert count, and failed or retrying notification-delivery count in the previous 24 hours.
- **FR-007**: Each successful summary MUST include at most one representative upcoming occurrence, recent failure, unacknowledged alert, and notification problem with safe identifiers needed for target-specific drill-down.
- **FR-008**: Operational summary queries MUST be bounded independently of total stored history and MUST NOT return task commands, arguments, environments, run output, log details, notification payloads, endpoints, authorization material, or other secret-bearing content.
- **FR-009**: The desktop MAY retain the latest successful summary for each registration during the current process lifetime and MUST mark it stale with its original observation time when a later refresh fails.
- **FR-010**: Operational summaries MUST NOT be persisted across desktop restarts and MUST NOT alter daemon state, task state, connection profiles, or delivery records.
- **FR-011**: Refresh results MUST remain useful when any subset of targets fails, and each failed target MUST expose a safe state-specific explanation and next action.
- **FR-012**: Users MUST be able to filter by text, connection state, and attention-needed status, and sort by display label, state severity, freshness, upcoming time, recent failures, alerts, and notification problems.
- **FR-013**: Same-named targets MUST remain distinguishable through immutable or registration-specific identity and location context in every row, detail, announcement, and drill-down.
- **FR-014**: Drill-down MUST resolve one profile identifier or the local target directly, establish that exact connection, and only then navigate to Tasks, Schedule, Activity, or Notifications.
- **FR-015**: Representative records MUST preserve task or record identifiers when available so the destination can expose the relevant source context without inferring identity from names.
- **FR-016**: Failed selection or reconnection MUST leave the user on All Systems, preserve current results, and provide a target-specific recovery path without opening mutable controls against stale identity.
- **FR-017**: Observe-authorized remote credentials MUST be sufficient to retrieve the read-only operational summary; no Operate or Manage authority is required.
- **FR-018**: The daemon summary contract MUST be versioned, additive, and independently retrievable through local and enabled remote transports.
- **FR-019**: The All Systems surface MUST support keyboard-only operation, visible focus, screen-reader status announcements, reduced motion, 200 percent zoom, an 800 by 600 viewport, and up to 100 saved profiles without document-level horizontal scrolling.
- **FR-020**: All labels and documentation MUST describe independent daemons, registered systems, partial results, and target ownership without implying a cluster, shared execution state, failover, synchronization, or exactly-once behavior.
- **FR-021**: S096 MUST complete the functional acceptance criteria of issue #182 while leaving cross-daemon mutation, bulk actions, portable bundles, drift reconciliation, SMTP, and native desktop notification delivery outside the slice.

### Key Entities

- **System registration**: This computer or one saved desktop connection profile, identified independently from its mutable display label.
- **Operational summary**: One bounded, read-only observation from a daemon with counts and representative records for current operational triage.
- **System observation**: The registration, current connection classification, freshness, successful summary if available, safe failure guidance, and refresh generation.
- **Drill-down intent**: A target registration plus one destination and optional task or record identity used only after the exact connection is established.

## Success Criteria

### Measurable Outcomes

- **SC-001**: For collections from one through 100 registrations, every registration appears exactly once and remains distinguishable even when all display labels match.
- **SC-002**: One target that never responds is classified within five seconds while every other target can complete independently.
- **SC-003**: No refresh executes more than four target operations concurrently, and superseded refresh results never alter the current generation.
- **SC-004**: Counts and representative records describe only the documented 24-hour windows and remain bounded regardless of stored task or history volume.
- **SC-005**: A failed refresh with prior session data retains exactly one stale summary with its original observation time; a failed refresh without prior data never invents operational values.
- **SC-006**: Every drill-down either selects the exact registration and opens its relevant source context or leaves the overview unchanged with a specific recovery action.
- **SC-007**: The overview is fully operable by keyboard with no serious accessibility findings or document-level horizontal overflow at 800 by 600 and 200 percent zoom.
- **SC-008**: Read-only overview requests produce no daemon mutations, no stored overview cache, and no secret-bearing response fields.

## Assumptions

- Existing local IPC, remote HTTPS, pinned daemon identity, saved connection profiles, native credential storage, and typed connection failures remain authoritative.
- The daemon's own current clock defines each summary's observation time and 24-hour windows; the desktop displays freshness without comparing operational records across clock domains as if perfectly synchronized.
- A profile registration remains distinct even if another profile points to the same daemon. Silent deduplication could hide authority or trust differences and is outside this slice.
- Session-only stale data is sufficient for partial-failure continuity. Durable fleet history and background monitoring remain outside S096.
- The maximum saved-profile contract remains 100, so bounded concurrency and progressive result updates are sufficient without pagination of registrations.

## Dependencies and Traceability

- Parent: #174.
- Completes: #182 when its functional acceptance criteria are delivered.
- Depends on completed #166, #170, and #172.
- Enables: #183.
- Leaves #176, #177, #184, and #185 open.
