# Feature Specification: Trigger Sets

**Feature Branch**: `codex/055-trigger-sets`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/055-trigger-sets`; focused persistence, migration, API, client, event, CLI, and headless GUI tests plus canonical format, vet, lint, race, coverage, documentation, and automation gates passed 2026-09-05.

**Input**: GitHub issue #134 under parent epic #17, adding bulk trigger-key administration after the individual trigger lifecycle delivered by S054.

## Clarifications

### Session 2026-09-05

- Q: What relationship exists between Trigger Sets and task Groups? A: Trigger Sets are an independent administration concept and never alter task Group hierarchy or eligibility semantics.
- Q: How are set members identified and ordered? A: Every member remains an ordinary trigger with its own stable identity and a permanent one-based position within one optional set; bulk output uses ascending position.
- Q: How are set-level mutations applied? A: Create, retarget, enable, disable, rotate, and delete operations are transactional and report no success until every requested member mutation commits.
- Q: Can one member of a set target a different task? A: No; individual member updates may change the member name, enabled state, key, or existence, but target changes occur only for the complete set.
- Q: What happens when the last member is deleted individually? A: The now-empty Trigger Set is deleted in the same transaction so every persisted set owns between 1 and 99 members.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create and Copy a Trigger Set (Priority: P1)

An operator creates between 1 and 99 distinct trigger keys for one task in a single operation and copies the complete invocation commands in a stable order.

**Why this priority**: Batch creation and exact bulk copy are the primary value of Trigger Sets.

**Independent Test**: Create a set at both supported count boundaries, retrieve its explicit secret output, and verify every unique member command invokes the same selected task.

**Acceptance Scenarios**:

1. **Given** an existing task, **When** an operator creates a named Trigger Set with a count from 1 through 99, **Then** the set and every uniquely keyed member are persisted atomically with the same target.
2. **Given** a newly created set, **When** the result is shown, **Then** all keys and complete commands are available in ascending member position and can be copied together.
3. **Given** an existing set, **When** an operator explicitly reveals its secrets, **Then** the same stable ordered output is available without exposing keys through ordinary list or detail views.
4. **Given** a requested count outside 1 through 99, **When** creation is attempted, **Then** nothing is persisted and an actionable validation result identifies the count constraint.

---

### User Story 2 - Administer Members Independently (Priority: P1)

An operator can identify and administer one trigger within a set without changing its siblings.

**Why this priority**: Separate keys are useful only when one compromised or retired caller can be isolated safely.

**Independent Test**: Disable, enable, rename, rotate, and delete one member, then verify every unaffected sibling retains its identity, key, target, position, and enabled state.

**Acceptance Scenarios**:

1. **Given** a multi-member set, **When** one member is disabled, enabled, renamed, rotated, or deleted, **Then** no sibling changes.
2. **Given** a set member, **When** an individual target change is attempted, **Then** the request is rejected with guidance to retarget the set.
3. **Given** a one-member set, **When** its member is deleted, **Then** the empty set is removed atomically.
4. **Given** a member deleted from the middle of a set, **When** remaining members are listed or copied, **Then** their permanent positions remain unchanged and the output remains ordered.

---

### User Story 3 - Administer a Complete Set (Priority: P2)

An operator can retarget, enable, disable, rotate, or delete every member through one explicit set-level action.

**Why this priority**: Set-level administration prevents repetitive work and inconsistent partial configuration.

**Independent Test**: Exercise each set-level action and verify all members change atomically, preserve stable identities where applicable, and retain exact provenance.

**Acceptance Scenarios**:

1. **Given** an existing set, **When** its target is changed, **Then** every member points to the new task in one transaction and no member retains the old target.
2. **Given** an existing set, **When** it is enabled or disabled, **Then** every member changes to the requested state atomically.
3. **Given** an existing set, **When** all keys are rotated after explicit confirmation, **Then** all previous keys become invalid together and the complete replacement secret output is shown once.
4. **Given** an existing set, **When** deletion is confirmed, **Then** the set and all members are removed atomically and every former key becomes unknown.

---

### User Story 4 - Use Trigger Sets from CLI and Desktop (Priority: P2)

An operator can complete the Trigger Set lifecycle through either machine-friendly commands or the desktop Triggers view.

**Why this priority**: Interface parity keeps scripted administration and approachable desktop workflows consistent.

**Independent Test**: Create and administer a set through the CLI and desktop independently, comparing visible identities, ordered output, confirmation, error, and redaction behavior.

**Acceptance Scenarios**:

1. **Given** the command line, **When** an operator uses Trigger Set lifecycle commands, **Then** human and JSON output represent the same stable identities and results.
2. **Given** the Triggers view, **When** set members are displayed, **Then** their set name and position are visible without exposing raw keys.
3. **Given** a selected set member, **When** an operator requests a set-level action, **Then** the confirmation identifies the set and affected member count.
4. **Given** a copied bulk command list, **When** it is pasted into a line-oriented script or terminal one line at a time, **Then** each nonblank line is one complete usable invocation command.

### Edge Cases

- Concurrent set creation cannot produce duplicate trigger keys or duplicate member positions.
- A set name must be nonblank but need not be globally unique; stable set identity disambiguates duplicates.
- Member positions are assigned once from 1 through the requested creation count and are never renumbered after individual deletion.
- A set containing disabled members may be enabled as one atomic action, but this does not enable its target task.
- Disabling, deleting, or retargeting set members must preserve the existing automatic-source invariant for both old and new target tasks.
- A failure to generate any key or persist any member leaves no set and no members behind.
- A set-level mutation failure leaves every set and member value unchanged.
- Ordinary set responses, events, logs, Activity entries, and errors never contain raw keys.
- Explicit bulk reveal and rotation may return raw keys and commands, but output contains no blank lines and ends with exactly one newline in human-readable CLI form.
- Existing triggers created outside a set remain valid and ungrouped after migration.
- Deleting a target task removes its Trigger Sets and members through the same atomic task-deletion operation.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST let an operator create a Trigger Set with a required nonblank display name, one existing target task, a requested member count from 1 through 99 inclusive, and an enabled state that defaults to enabled.
- **FR-002**: Trigger Set creation MUST atomically persist the set and exactly the requested number of ordinary trigger members, each with a unique generated key, stable trigger identity, permanent one-based position, shared target task, and requested enabled state.
- **FR-003**: Trigger Set names MAY repeat; every set MUST have a stable identity and creation and update timestamps.
- **FR-004**: Each trigger MUST belong to at most one Trigger Set, and existing standalone triggers MUST remain valid without set membership.
- **FR-005**: Every persisted Trigger Set MUST own between 1 and 99 members; deleting its final member individually MUST delete the empty set in the same transaction.
- **FR-006**: Deleting a nonfinal member MUST preserve every sibling identity, key, target, position, and enabled state without renumbering positions.
- **FR-007**: Individual member rename, enable, disable, rotate, reveal, fire, and delete operations MUST retain the S054 trigger contract and affect no sibling.
- **FR-008**: The system MUST reject an individual target change for a set member and identify the set-level retarget action as the supported path.
- **FR-009**: The system MUST support listing, inspecting, retargeting, enabling, disabling, rotating, revealing, and deleting complete Trigger Sets through versioned local interfaces.
- **FR-010**: Set-level retarget, enable, disable, rotate, and delete operations MUST be atomic and MUST report failure without a persisted partial result.
- **FR-011**: Set-level retargeting MUST update every member target while preserving set identity, member identities, positions, keys, and enabled states.
- **FR-012**: Set-level key rotation MUST replace every member key atomically, invalidate all previous keys before success is reported, and return replacement secrets only through the explicit rotation result.
- **FR-013**: Ordinary Trigger Set list and detail output MUST omit raw keys; create, explicit reveal, and rotate results MAY contain ordered raw keys and complete commands.
- **FR-014**: Bulk secret output MUST order members by permanent ascending position, include stable set and trigger identities in structured output, and represent human-readable commands as one `gosched trigger fire <key>` command per nonblank line with exactly one final newline.
- **FR-015**: The CLI MUST provide consistent verb-noun commands for Trigger Set create, list, show, retarget, enable, disable, rotate, reveal, and remove operations with human-readable and `--json` output where output is produced.
- **FR-016**: The desktop Triggers view MUST identify set membership and permanent position for every member while keeping standalone triggers visually distinguishable.
- **FR-017**: The desktop MUST support Trigger Set creation, bulk command copy, retarget, enable, disable, rotate, and delete through a selected member, with confirmations naming the set and affected member count for destructive or broad mutations.
- **FR-018**: Desktop create, reveal, rotate, and copy flows MUST provide visible accessible status without placing raw keys in ordinary tables or notifications.
- **FR-019**: Trigger Set and member lifecycle changes MUST propagate to connected desktop clients without application restart.
- **FR-020**: Existing task readiness, group eligibility, overlap, concurrency, dispatch, run provenance, and trigger firing behavior MUST remain authoritative for every member.
- **FR-021**: Disabling, deleting, or retargeting members MUST apply the existing automatic-source invariant atomically for each affected task, and set creation or enabling MUST NOT automatically enable a target task.
- **FR-022**: Deleting a target task MUST atomically remove every Trigger Set targeting it and all of their members.
- **FR-023**: Trigger Set storage migration MUST be forward-only, preserve all existing standalone triggers, and be verified from schema version 12.
- **FR-024**: Set and member identifiers MAY appear in logs, events, Activity, and errors, but raw trigger keys MUST appear in none of those surfaces.
- **FR-025**: Persistence, migration, ordering, transaction rollback, individual isolation, automatic-source invariants, redaction, API, CLI, event, and desktop behavior MUST have deterministic automated coverage.
- **FR-026**: Creating or mutating the maximum-size set MUST complete within one second under nominal local load, excluding user interaction time.
- **FR-027**: S055 MUST add no network listener, remote invocation, arbitrary payload, filesystem watcher, Chain target, or task Group behavior.

### Key Entities

- **Trigger Set**: A persisted named administration boundary with stable identity, one target task, timestamps, and between 1 and 99 ordinary trigger members.
- **Set Member**: An ordinary external trigger with optional set identity and a permanent one-based position; its key, enabled state, firing behavior, and run provenance follow the S054 contract.
- **Bulk Secret Result**: An explicitly requested ordered collection containing set identity plus each member's position, stable trigger identity, raw key, and complete invocation command.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can create and copy a set of 99 usable commands in one create action followed by one copy action.
- **SC-002**: Creating sets at counts 1 and 99 yields exactly the requested number of unique keys and stable positions in 100 consecutive automated trials.
- **SC-003**: Every individual lifecycle operation changes only one selected member, with zero sibling field changes across sets of 2 and 99 members.
- **SC-004**: Every set-level mutation either changes all intended members or changes none when an injected failure occurs.
- **SC-005**: Bulk human output is byte-stable across repeated reads, contains one complete command per nonblank line in ascending position, and contains exactly one final newline.
- **SC-006**: Automated redaction tests find zero raw-key occurrences across ordinary set and trigger lists, details, events, logs, Activity records, errors, and run history.
- **SC-007**: A desktop operator can complete create, membership inspection, bulk copy, retarget, enable, disable, rotate, and delete without using the CLI.
- **SC-008**: Creating, revealing, retargeting, enabling, disabling, rotating, or deleting a 99-member set completes within one second under nominal local load.
- **SC-009**: Existing standalone trigger, task, Group, Chain, schedule, execution, and release verification suites pass unchanged.

## Assumptions

- Trigger Sets organize trigger credentials only and are unrelated to hierarchical task Groups.
- A member position is a permanent creation ordinal rather than a dense current row number.
- Set-level operations use transactions, so precise partial-success reporting is unnecessary because no partial persisted success is permitted.
- Generated member display names use the set name plus a zero-padded position, while later individual renames remain allowed.
- Explicit bulk copy uses complete commands by default because they are immediately executable; structured output includes both keys and commands.
- Existing local IPC authorization and recoverable raw-key storage remain the security boundary established by S054.
- Filesystem watchers, remote invocation, arbitrary payloads, direct Chain targets, adding members after creation, and moving existing standalone triggers into or between sets are outside S055.
