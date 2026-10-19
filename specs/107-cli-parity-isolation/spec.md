# Feature Specification: CLI parity and isolated operation

**Feature Branch**: `codex/107-cli-parity-isolation`

**Created**: 2026-10-19

**Status**: Implemented

**Delivery**: [PR #305](https://github.com/shruggietech/go-schedule/pull/305) publishes review branch `codex/107-cli-parity-isolation`, implementing #283 through #288 and corrections #304, #306, #307 and #308. The corrected source passes all eight local canonical gates, the six-target build matrix, and focused native Linux/WSL race and packaged-CLI checks; see [verification.md](verification.md). Hosted CI and Codex review evidence is maintained in the PR. The feature remains unreleased and unmerged.

**Input**: Owner kickoff: run S107 under autopilot with the full six-issue bundle, #283, #284, #285, #286, #287 and #288, stopping before publication. Parent: #277. Baseline: `bae303071d1dd1af050cffb201267052ae58b2c7`.

## Clarifications

### Session 2026-10-19

- Q: Does this slice include the comprehensive unattended Windows/WSL runner and reporting product? A: No. Deliver every functional outcome in #283 through #288; #289 and #290 remain the following slice.
- Q: Does isolation require changing normal production credential handling or restricting ordinary product use? A: No. Preserve production defaults; inject explicitly scoped recording dependencies into the harness and require ownership within the QA fixture lifecycle.
- Q: Must persistent preference changes launch or live-reconfigure an open desktop? A: No. Headless commands persist the shared contract; the next desktop startup reads the saved default connection.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Operate an isolated target (Priority: P1)

A maintainer explicitly targets a disposable daemon and exercises command handlers without touching the installed daemon, personal settings, saved connections, native credentials, services or desktop.

**Why this priority**: The remaining authoring operations must be usable against owned state, and the future unattended evaluator depends on this boundary (#288).

**Independent Test**: Select a disposable target, read its identity and paths, change a fixture task and remove only the fixture's processes/state. Exercise service and launcher intents with recording dependencies.

**Acceptance Scenarios**:

1. **Given** installed and disposable daemons, **When** an explicit local target is supplied, **Then** only that target receives the command; an invalid selector fails without fallback.
2. **Given** recording service, secret and launcher dependencies with isolated persistence, **When** the actual command handlers run, **Then** their outputs and recorded intents are available without host service operations, native credential changes or a visible window.
3. **Given** a workspace with one owned daemon, **When** ownership checks or cleanup run, **Then** mismatched identity/paths are rejected and only exact owned processes and contained state are removed.

### User Story 2 - Author tasks and groups directly (Priority: P1)

An operator creates disabled tasks/groups atomically, edits execution inputs and hierarchy, clears individual inputs explicitly, selects schedule interpretation and previews a complete timing policy without creating tasks (#283, #287).

**Why this priority**: These are direct administration gaps and include a shared desktop environment-clear defect.

**Independent Test**: Create a disabled fixture task/group, set and clear every input, rename/reparent a group and preview human/cron timing with preserved state before and after.

**Acceptance Scenarios**:

1. **Given** populated inputs, **When** arguments, environment overrides, stdin, run identity or working directory are explicitly cleared, **Then** those inputs are empty while omitted fields remain unchanged. Clearing overrides retains the process's base environment.
2. **Given** a fully configured task or group, **When** it is created disabled, **Then** no intermediate enabled state exists.
3. **Given** a nested group, **When** its name/parent is edited or it is moved to root, **Then** the existing hierarchy and cycle rules apply.
4. **Given** a human or supported cron expression and all date/time policies, **When** preview runs with explicit syntax intent, **Then** understandable text or structured predicted occurrences are returned without mutation.
5. **Given** the desktop removes its final task environment row, **When** the task is saved/read again, **Then** no previous explicit overrides remain.

### User Story 3 - Maintain independent notification policies (Priority: P2)

An operator reads, modifies and reapplies different conditions for individual destinations directly on a task/group, with the same supported inheritance behavior (#284).

**Independent Test**: Apply a mixed policy from file/stdin, read and reapply it, reject a malformed replacement and clear direct assignments to resume inheritance.

**Acceptance Scenarios**:

1. **Given** two destinations with different supported conditions, **When** one structured replacement is applied, **Then** their settings remain distinct across read/modify/reapply.
2. **Given** invalid/duplicate assignments or conflicting input forms, **When** replacement is requested, **Then** no partial mutation occurs and the error identifies the correction.
3. **Given** direct assignments, **When** an empty replacement is supplied, **Then** existing inheritance semantics resume; existing shorthand remains compatible.

### User Story 4 - Inspect operations and connected systems (Priority: P2)

An operator queries calendar windows, active/exact executions, output/provenance, runtime/log metadata and summaries/search across selected or explicitly all registered systems (#285).

**Independent Test**: Inspect a fixture execution and calendar, then query multiple fixture targets including an unreachable/mismatched target and retain successful identified results.

**Acceptance Scenarios**:

1. **Given** active and recorded runs, **When** active/exact/detail queries execute, **Then** output, timestamps, truncation and provenance are available in readable and structured forms without changing existing listing formats.
2. **Given** a selected target, **When** calendar/runtime/log inspection runs, **Then** predicted events and execution records are distinct and paths belong to that target; unavailable remote logs never become local logs.
3. **Given** multiple registrations, **When** an explicit all-system summary/search runs, **Then** concurrency, time and result counts are bounded, every result retains registration/daemon identity and failures remain independent of successes.

### User Story 5 - Administer durable desktop choices headlessly (Priority: P2)

A user reads, updates or resets appearance, popup conditions/severities/source filters and the saved desktop default connection without opening the desktop or emitting popups (#286).

**Independent Test**: Use an isolated user store to migrate preferences, patch/reset each durable setting and set/reset a registered default connection, then read the same persisted choices through the desktop component.

**Acceptance Scenarios**:

1. **Given** existing or legacy choices, **When** desktop or CLI reads/updates them, **Then** both use identical validation, migration, defaults and atomic persistence without resetting unrelated state.
2. **Given** malformed or omitted preference fields, **When** a partial update runs, **Then** malformed values do not alter persistence and omitted values remain unchanged.
3. **Given** a saved connection registration, **When** its desktop-default choice is set/reset, **Then** the next desktop startup reads that choice while invocation-only CLI targeting remains independent.

### Edge Cases

- Missing/malformed/conflicting target selectors, identity mismatch and workspace paths outside ownership never cause default-target fallback or host mutation.
- Explicit empty strings/collections differ from omitted values; mutually exclusive set/clear flags fail before mutation.
- Unsupported syntax, invalid time/date policies, cycles, invalid preferences, duplicate destinations and malformed/trailing input return actionable diagnostics with no partial writes.
- Remote timeouts, denied capabilities, missing credentials, stale identity and absent local-only logs are per-target failures, never substituted local results.
- Cancellation terminates streams/fan-out and owned children; cleanup retains evidence and never terminates by executable name or removes unrelated files.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Explicit local targeting MUST select the specified daemon or fail, with incompatible remote target selection rejected and no fallback (#288).
- **FR-002**: Each command invocation MUST have independently scoped client, profile/preference storage, secrets, service, launcher, clock and I/O dependencies, retaining production defaults and native secret handling (#288).
- **FR-003**: QA fixtures MUST identify owned workspace, daemon identity, effective state paths and exact child processes; mutating fixture operations and cleanup MUST reject ownership mismatches (#288).
- **FR-004**: Recording service/launcher/secret adapters MUST exercise actual handlers and diagnostics without host service/session/desktop changes; evidence MUST distinguish packaged CLI, command harness and deterministic fixtures (#288).
- **FR-005**: Task authoring MUST set/clear stdin, run identity, arguments, working directory and environment overrides while omission leaves values unchanged; the desktop uses the same environment-clear semantics (#283, #287).
- **FR-006**: Task/group creation MUST support disabled state in the creation operation; group editing MUST support rename/reparent/root with existing hierarchy validation. Rejected combined edits leave the entire group unchanged (#283, discovered defect #304).
- **FR-007**: Task create/edit/preview MUST support explicit automatic/human/cron interpretation; preview MUST include all existing timing policies without mutation (#283).
- **FR-008**: Direct notification assignment MUST support heterogeneous structured file/stdin input, round trips, complete validation before replacement, inheritance restoration and compatible shorthand (#284).
- **FR-009**: Calendar, active/exact/detail run and runtime/log inspection MUST expose available target-specific data in human/structured forms while preserving existing formats (#285).
- **FR-010**: Selected/all-registered summaries/search MUST preserve registration and daemon identity, bounded work, successful partial results and independent failures under existing authorization (#285).
- **FR-011**: Durable desktop preference read/update/reset MUST share current migration/defaults/validation/atomic persistence with the desktop, remain headless and isolate user state when requested (#286).
- **FR-012**: Saved desktop default connection read/set/reset MUST use the existing registration store, remain distinct from invocation targeting and be consumed at next desktop startup (#286).
- **FR-013**: New commands MUST document intent, constraints, input/output shapes and corrective diagnostics; successes go to stdout and diagnostics to stderr, with structured output where applicable.
- **FR-014**: Windows console children MUST remain hidden/noninteractive with redirected I/O; owned cleanup MUST use exact process identity and contained paths, preserve evidence and handle cancellation.

### Key Entities

- **Invocation**: Selected target, scoped dependencies, input/output and cancellation for one command execution.
- **Owned workspace**: Unique run identity, daemon identity, effective paths, process ownership and retained evidence.
- **Edit intent**: Omitted, replace or explicitly clear values, with conflict validation.
- **Operational result**: Target identity, predicted/recorded evidence, data or independent failure.
- **Desktop preferences**: Versioned user choices and migrations, separate from daemon notification policy and saved connection registrations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All six included issues' functional acceptance criteria have direct delivered commands/shared behavior and issue-level traceability.
- **SC-002**: Every explicit isolated-target operation either affects only owned state or fails without affecting the installed target, real preferences, native credentials, services or desktop.
- **SC-003**: All supported task inputs and durable desktop choices can be changed/reset without opening the desktop; unrelated values remain unchanged.
- **SC-004**: Multi-system queries finish within their configured deadlines, retain each successful target result and identify every failed target separately.
- **SC-005**: Existing CLI scripting output and persistent user settings remain compatible while new modes expose the tracked missing capabilities.

## Assumptions

- Existing scheduling, notification, authorization, transport and persistence behavior remains the shared product foundation. No parallel scheduler or new product-use restrictions are introduced.
- The comprehensive Windows/WSL coordinator, scenario manifest, communication review/report schema and maintainer prompts remain #289/#290 for the next slice. This slice delivers their isolation prerequisite and the complete six-issue parity bundle.
- No real service reconfiguration, login/logout, host restart, public DNS change, release or publication is required. The normal pre-publication boundary remains in force.
