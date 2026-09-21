# Feature Specification: Cross-Daemon Search and Target-Safe Actions

**Feature Branch**: `codex/097-cross-daemon-search`

**Created**: 2026-09-21

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: review branch `codex/097-cross-daemon-search`

**Input**: Work slice S097, complete issue #183 by adding bounded search across registered daemons and explicit, identity-safe actions without introducing clustered ownership or distributed transactions.

## Clarifications

### Session 2026-09-21

- Q: When does cross-daemon search run? -> A: Search runs only after an explicit user submission, requires a non-empty query, and progressively reports independent per-target results.
- Q: Which records are searchable? -> A: Tasks, groups, failed runs, upcoming schedule occurrences, and unacknowledged alerts are searchable through bounded summaries that exclude commands, output, credentials, and other secret-bearing detail.
- Q: How are multi-result mutations handled? -> A: Users may select results across daemons only for a single compatible action, must review a confirmation grouped by exact daemon identity, and receive independent per-target outcomes without rollback or atomicity claims.
- Q: What happens when a result becomes stale before mutation? -> A: The desktop reconnects to the recorded registration, confirms the current daemon identity and authority, and revalidates each selected object before enabling submission.
- Q: How are permission and version differences handled? -> A: Every result exposes current action availability and a reason for each unsupported action; mixed selections containing unsupported or unverified objects cannot be submitted until those objects are removed or refreshed.

## User Scenarios & Testing

### User Story 1 - Find automation across registered daemons (Priority: P1)

As an operator, I can search tasks, groups, failures, upcoming schedule occurrences, and alerts across This computer and saved daemon profiles without losing machine ownership.

**Why this priority**: Cross-daemon discovery is the foundational outcome. Actions are unsafe if results cannot first identify their source and freshness unambiguously.

**Independent Test**: Register daemons containing duplicate names across every supported result type, submit one query, and confirm that progressive results identify each registration, daemon, connection state, and observation time while individual target failures remain visible.

**Acceptance Scenarios**:

1. **Given** This computer and several connected profiles, **When** I submit a query, **Then** matching tasks, groups, failures, schedule occurrences, and alerts appear with their exact source registration and daemon identity.
2. **Given** two daemons containing identically named objects, **When** both match, **Then** label, location, shortened immutable identity, connection state, and freshness keep the results distinguishable.
3. **Given** one unavailable, unauthorized, incompatible, or timed-out target, **When** other targets finish, **Then** successful results remain usable and the failed target reports its own state and recovery guidance.
4. **Given** more matches than the result limit, **When** search completes, **Then** the interface reports truncation per target and preserves deterministic ordering.

---

### User Story 2 - Open the exact source safely (Priority: P1)

As an operator, I can open a search result on its owning daemon without the desktop inferring ownership from a display name.

**Why this priority**: Opening the wrong daemon creates the same operational risk as mutating the wrong daemon and must use the same identity boundary.

**Independent Test**: Open same-named results from two profiles and confirm each request reconnects to the recorded registration, validates daemon identity, selects the correct destination, and preserves the object context.

**Acceptance Scenarios**:

1. **Given** a current result, **When** I open it, **Then** the desktop selects the recorded registration and navigates to its relevant destination with the object context preserved.
2. **Given** a stale result, **When** I open it, **Then** the desktop reconnects and confirms current daemon identity before exposing destination controls.
3. **Given** changed trust, identity, authorization, or object existence, **When** open revalidation fails, **Then** the search remains visible and explains the target-specific recovery action.

---

### User Story 3 - Perform explicit target-safe actions (Priority: P1)

As an operator, I can acknowledge alerts and enable, disable, or run tasks only after reviewing the exact daemons and objects that will be affected.

**Why this priority**: Cross-daemon convenience must not weaken authorization, identity pinning, or deliberate operator intent.

**Independent Test**: Select compatible objects on one or more daemons, review the grouped confirmation, change one daemon's identity or permissions before submission, and confirm that only revalidated targets can execute with independent outcomes.

**Acceptance Scenarios**:

1. **Given** one selected current task, **When** I choose enable, disable, or run now, **Then** confirmation names the daemon, registration, task, requested action, and current authority before submission.
2. **Given** selected alerts across multiple daemons, **When** I choose acknowledge, **Then** confirmation groups every alert by daemon and execution produces a separate result for each target.
3. **Given** a mixed selection containing unsupported, stale, or wrong-kind objects, **When** I request an action, **Then** submission is refused and each incompatible object explains why it must be removed or refreshed.
4. **Given** several confirmed target groups, **When** one target fails after another succeeds, **Then** successful changes remain applied, failed objects retain their failure result, and the interface makes no rollback or atomicity claim.
5. **Given** a stale selected result, **When** the desktop reconnects to a different daemon identity or no longer finds the object, **Then** that object is not mutated and the mismatch is explicit.

---

### User Story 4 - Search and act accessibly at scale (Priority: P2)

As a keyboard or assistive-technology user, I can submit searches, inspect large result sets, select compatible results, review confirmations, and understand partial outcomes without losing focus or source context.

**Why this priority**: A dense multi-system workspace must remain usable under the supported window, zoom, and profile limits.

**Independent Test**: Exercise the complete workflow by keyboard at 800 by 600, 200 percent zoom, reduced motion, 100 registered profiles, duplicate labels, and partial connection failures.

**Acceptance Scenarios**:

1. **Given** progressive target completions, **When** results and failures arrive, **Then** status changes are announced without moving focus or changing the active selection unexpectedly.
2. **Given** a large result set, **When** I filter by kind, target, connection state, or action availability, **Then** controls and source identity remain reachable with no document-level horizontal scrolling.
3. **Given** a confirmation or outcome summary, **When** I navigate by keyboard or assistive technology, **Then** target groups, object names, action, success, and failure states are explicit without color-only meaning.

### Edge Cases

- The query consists only of whitespace; no network work begins and the search control explains that a query is required.
- A profile is added, edited, or removed while search is running; late results from the superseded registration generation cannot restore removed or outdated registrations.
- A newer search begins before the older search completes; the older generation is canceled and cannot overwrite the current query.
- A target returns malformed, unsupported, or future-version search data; that target is marked incompatible without blocking compatible targets.
- An object is renamed, deleted, changes type, or changes state between search and confirmation; revalidation reports the current condition and prevents unintended mutation.
- The same daemon is registered through multiple profiles; every registration remains distinct and actions use the deliberately selected registration and its authority.
- A target succeeds for some objects and fails for others; outcomes remain attached to each object and no target-wide success is inferred.
- A run-now request is repeated after an uncertain transport result; the request uses the existing retry-safe action contract and does not silently create an additional run.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST provide a Cross-Daemon Search destination covering This computer and every saved connection profile.
- **FR-002**: Search MUST begin only after explicit submission of a non-empty normalized query and MUST allow filtering by result kind, target, connection state, and action availability.
- **FR-003**: Search MUST cover task names and identifiers, group names and identifiers, failed-run task context, upcoming schedule task context, and unacknowledged alert summaries without returning commands, arguments, environment values, output, credentials, endpoints containing secrets, or notification payloads.
- **FR-004**: Every result MUST identify its registration, display label, local or remote location, daemon identity when known, connection state, observation time, result kind, object identity, and safe display summary.
- **FR-005**: Same-named registrations and objects MUST remain distinguishable in result rows, detail views, selection summaries, confirmations, announcements, and outcomes.
- **FR-006**: Search MUST execute independently per target with no more than eight target operations in flight, a three-second deadline per target, and progressive target completion.
- **FR-007**: Each target MUST return at most 50 matches with deterministic ordering, and the desktop MUST report target-specific truncation instead of implying completeness.
- **FR-008**: A newer search or registration generation MUST cancel the previous generation, and late results MUST NOT overwrite current results or reintroduce removed or changed registrations.
- **FR-009**: Per-target search failures MUST preserve successful results from other targets and expose safe, state-specific recovery guidance.
- **FR-010**: Search results MUST remain read-only, process-local observations and MUST NOT be persisted as an operational cache.
- **FR-011**: Every result MUST describe whether open, acknowledge, enable, disable, and run-now are available, unavailable, or require refresh, including a concise reason based on object kind, connection state, daemon capability, version, and current authority.
- **FR-012**: Opening a result MUST resolve the recorded local target or saved profile identifier directly, re-establish the connection, confirm current daemon identity, and only then navigate with object context.
- **FR-013**: A mutating action MUST require deliberate selection of every affected object and a confirmation that groups objects by exact daemon identity and names the requested action, registration, objects, and authority before submission.
- **FR-014**: Bulk selection MAY span daemons only when one requested action is compatible with every selected object; ambiguous mixed-kind, mixed-action, unsupported, stale, or unverified selections MUST be refused before submission.
- **FR-015**: Immediately before mutation, the desktop MUST reconnect through each recorded registration, confirm its current daemon identity and authority, and revalidate object identity, kind, existence, and relevant state.
- **FR-016**: Identity, trust, authorization, capability, version, connection, or object revalidation failure MUST prevent mutation of the affected object and MUST NOT redirect the action to another registration or similarly named object.
- **FR-017**: Cross-daemon mutations MUST execute independently per target and report one outcome per object; successful mutations MUST remain applied when another target fails, and the desktop MUST NOT claim atomicity or automatic rollback.
- **FR-018**: Alert acknowledgement MUST require Operate authority and task enable, disable, and run-now MUST use the existing Operate-authorized retry-safe action contracts and audit attribution.
- **FR-019**: Observe authority MUST be sufficient for search and open workflows, while mutation availability MUST reflect current Operate authority without exposing authorization secrets.
- **FR-020**: The daemon search contract MUST be versioned, additive, bounded independently of stored history, and available through local and enabled remote transports.
- **FR-021**: Search, selection, confirmation, cancellation, target progress, and per-object outcomes MUST support keyboard-only operation, visible focus, screen-reader announcements, reduced motion, 200 percent zoom, an 800 by 600 viewport, and up to 100 saved profiles without document-level horizontal scrolling.
- **FR-022**: Product labels and documentation MUST describe independent daemons, registrations, target-specific operations, and partial outcomes without implying cluster membership, shared ownership, failover, synchronization, distributed transactions, or exactly-once execution.
- **FR-023**: S097 MUST complete the functional acceptance criteria of issue #183 while leaving cross-daemon task creation or editing, shared ownership, rollback orchestration, saved search history, and background fleet indexing outside the slice.

### Key Entities

- **Search request**: A normalized query, filters, generation, and documented per-target limit submitted explicitly by the operator.
- **Target search observation**: One registration's current connection classification, daemon identity, authority, capabilities, observation time, matches, truncation state, or safe failure guidance.
- **Search result**: A bounded, secret-free task, group, failed-run, schedule-occurrence, or alert summary tied to one immutable registration and daemon observation.
- **Action availability**: Per-result status and rationale for open, acknowledge, enable, disable, and run-now based on current kind, state, capability, version, and authority.
- **Action intent**: One requested action plus deliberately selected object references grouped by their exact target registrations and daemon identities.
- **Action outcome**: One object's success, rejection, or failure tied to its target and request identity without any cross-target rollback guarantee.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Across one through 100 registrations, every displayed match identifies exactly one source registration and remains distinguishable when all target and object display names match.
- **SC-002**: A responsive target can publish progressive results without waiting for an unresponsive target, and every target reaches a result or target-specific timeout classification within three seconds of beginning its own search.
- **SC-003**: No search executes more than eight target operations concurrently, no target returns more than 50 matches, and superseded generations never alter the current result set.
- **SC-004**: Search responses contain no commands, arguments, environment values, run output, credentials, secret-bearing endpoints, notification payloads, or unbounded history.
- **SC-005**: Every supported mutation presents the requested action, exact daemon identity, registration, authority, and affected object names before submission.
- **SC-006**: Every mutation either revalidates the same target and object immediately before execution or produces an explicit non-mutating mismatch outcome; no action is redirected by display-name similarity.
- **SC-007**: Multi-target execution produces one outcome per selected object, retains successful outcomes when other targets fail, and never presents the operation as atomic or rolled back.
- **SC-008**: The complete search, filter, selection, confirmation, cancellation, and outcome workflow is keyboard operable with no serious accessibility finding or document-level horizontal overflow at 800 by 600 and 200 percent zoom.

## Assumptions

- Existing local IPC, remote HTTPS, saved profile identity pinning, connection-state classification, and native credential storage remain authoritative.
- Existing Observe and Operate authority semantics remain unchanged. Search adds no new authority level and mutations reuse established audited action paths.
- The supported saved-profile ceiling remains 100. Explicit submission, bounded fan-out, per-target limits, and progressive results are sufficient without a durable fleet index.
- Result freshness is an observation property, not a mutation guarantee. Reconnection and revalidation are mandatory even when a displayed result appears current.
- Same-daemon duplicate registrations remain independent because credentials, trust, and intended routing may differ.

## Dependencies and Traceability

- Parent: #174.
- Completes: #183 when its functional acceptance criteria are delivered.
- Depends on completed #167, #170, #172, and #182.
- Reuses capabilities from #178, #180, #181, and #182.
- Leaves #176, #177, #184, and #185 open.
