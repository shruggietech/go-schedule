# Feature Specification: Filesystem Watchers

**Feature Branch**: `codex/056-filesystem-watchers`

**Created**: 2026-09-05

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Implemented on review branch `codex/056-filesystem-watchers`; pull-request verification pending.

**Input**: GitHub issue [#135](https://github.com/shruggietech/go-schedule/issues/135), child of event-driven run epic [#17](https://github.com/shruggietech/go-schedule/issues/17), plus the operator-approved S056 autopilot scope.

## Clarifications

### Session 2026-09-05

- Q: How is a missing configured path distinguished as a file or directory? -> A: The operator selects an explicit watcher kind.
- Q: What does a directory glob match? -> A: It matches file base names, while recursion independently controls depth.
- Q: Which filesystem outcomes request runs? -> A: Present matching regular files after create, content-write, or rename-into-place signals become stable.
- Q: What happens to events missed during downtime or degraded health? -> A: They are not replayed or claimed; observation resumes prospectively.
- Q: How are symbolic links, junctions, network paths, and observer failures handled? -> A: Link traversal is excluded, network paths are best effort, and failures produce deduplicated degraded health with bounded recovery.

## User Scenarios & Testing

### User Story 1 - Run a task when a file becomes ready (Priority: P1)

An operator creates a watcher for one file or for files within a directory, assigns an ordinary task, and expects a matching file creation, content change, or rename into place to request exactly one run after the file stops changing.

**Why this priority**: Reliable event-to-task dispatch is the feature's core value and must preserve the scheduler's existing eligibility, overlap, concurrency, diagnostics, and history behavior.

**Independent Test**: Configure a watcher against a temporary directory, produce matching files through direct writes and temporary-file renames, and observe bounded runs with watcher provenance only after the configured quiet and stability periods.

**Acceptance Scenarios**:

1. **Given** an enabled watcher targeting an eligible task, **When** a matching regular file is created and remains unchanged for the configured stability period, **Then** the task receives one run request through the normal dispatcher and the run identifies the watcher as its source.
2. **Given** several rapid writes to the same matching file, **When** the write storm ends, **Then** the watcher produces one run request after its debounce and stability conditions are met.
3. **Given** a producer that writes a temporary file and renames it to a matching final name, **When** the final file is stable, **Then** the watcher produces one run request for the final path.
4. **Given** a matching event whose task is disabled, inactive, incomplete, blocked by a group, or constrained by overlap policy, **When** dispatch is attempted, **Then** the existing scheduler behavior remains authoritative and the watcher does not bypass it.

---

### User Story 2 - Configure and administer watchers (Priority: P1)

An operator can create, inspect, update, enable, disable, and delete file or directory watchers through both supported administration interfaces, with clear path, selection, timing, target, readiness, and health information.

**Why this priority**: A daemon capability without complete administration and visibility cannot be used safely or diagnosed by ordinary operators.

**Independent Test**: Complete the full lifecycle through the command line and desktop interface, then confirm every mutation takes effect without restarting the daemon and produces one live lifecycle event.

**Acceptance Scenarios**:

1. **Given** a task and a path that may or may not currently exist, **When** the operator creates a file watcher or directory watcher, **Then** the configuration is saved and its runtime health accurately reports whether observation is active.
2. **Given** an existing watcher, **When** the operator changes its path, kind, file-name pattern, recursion, debounce, stability delay, target, or enabled state, **Then** the daemon atomically adopts the new configuration without requiring a restart.
3. **Given** an existing watcher, **When** the operator disables or deletes it, **Then** future matching events do not request task runs and pending unsettled events are discarded.
4. **Given** a watcher in any lifecycle or health state, **When** it appears in ordinary output, **Then** its configuration and actionable health are visible without exposing unrelated filesystem contents.

---

### User Story 3 - Recover from filesystem changes and daemon restarts (Priority: P2)

An operator expects watchers to recover predictably when a root is missing, replaced, temporarily inaccessible, or affected by an observer failure, and expects restart behavior to avoid claiming that unobserved historical events were processed.

**Why this priority**: Filesystems change independently of the daemon, so truthful health and bounded recovery are necessary for long-running reliability.

**Independent Test**: Start with a missing root, create it, replace it, remove permissions where supported, inject an observer error, and restart the watcher runtime while asserting state transitions, recovery, and absence of replay.

**Acceptance Scenarios**:

1. **Given** a configured root that is missing or inaccessible, **When** the watcher runtime starts, **Then** the watcher becomes degraded with an actionable reason and periodically retries without flooding Activity or logs.
2. **Given** a degraded watcher whose root becomes observable, **When** the next bounded recovery attempt succeeds, **Then** health becomes active and subsequent matching events can dispatch.
3. **Given** an active watcher whose root is removed or replaced, **When** observation is lost, **Then** health becomes degraded and the watcher re-establishes observation against the configured root rather than following the removed object.
4. **Given** a daemon restart, **When** configured watchers are restored, **Then** observation starts from the new daemon lifecycle and no run is claimed for events that occurred while the daemon was stopped.

### Edge Cases

- A file watcher has an explicit `file` kind, so a path that does not exist at creation time is not ambiguously interpreted as a directory.
- Directory patterns match file base names, not path separators; recursive selection controls depth independently from the pattern.
- Directory creation, metadata-only changes, removals, and non-regular files do not dispatch tasks; a rename into place dispatches only when the resulting matching regular file becomes stable.
- Recursive observation does not traverse symbolic-link directories or Windows junctions, preventing cycles and scope escape; a configured root that itself resolves through such a link is rejected as unsupported.
- Network and removable paths are best effort: they become active only if the platform observer accepts them and otherwise remain degraded with the platform error summarized.
- Observer queue overflow or channel failure degrades affected watchers, clears their pending events, and starts bounded recovery without replay.
- A configuration mutation racing with a filesystem event uses a generation boundary so an event from the replaced configuration cannot dispatch afterward.
- Deleting a target task cascades its watchers; disabling the final automatic source disables a timeless task using the existing activation-readiness contract.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST persist named filesystem watchers with a stable identifier, explicit `file` or `directory` kind, configured path, optional file-name pattern, recursive selection, debounce duration, stability duration, target task, enabled state, and creation and update timestamps.
- **FR-002**: The system MUST accept watcher creation for a currently missing path because the explicit kind removes path-type ambiguity, while reporting the watcher as degraded until the path becomes observable.
- **FR-003**: A file watcher MUST select only its exact configured regular file; a directory watcher MUST select regular files under its root at the configured depth whose base name matches its pattern, with `*` as the default pattern.
- **FR-004**: The system MUST treat create, content-write, and rename-into-place events that leave a selected regular file present as dispatch candidates, and MUST ignore removals, directory-only changes, metadata-only changes, and unsupported non-regular files.
- **FR-005**: Dispatch candidates for the same watcher and file MUST be coalesced until no candidate event occurs for the configured debounce duration.
- **FR-006**: After debounce, the system MUST require the selected file's size and modification time to remain unchanged for the configured stability duration before requesting a run; a change or temporary disappearance restarts settling rather than dispatching early.
- **FR-007**: Each settled candidate MUST request one run through the existing overlap-aware task dispatcher and MUST preserve task command readiness, active state, enabled state, ancestor-group eligibility, worker limits, overlap policy, diagnostics, and run history.
- **FR-008**: Filesystem-originated runs MUST carry the `filesystem_watcher` trigger classification and the stable watcher identifier as provenance without recording the matched path in append-only run history.
- **FR-009**: The daemon MUST load all enabled watchers at startup, observe only future events after registration, and MUST NOT scan or replay changes that occurred while stopped or degraded.
- **FR-010**: The daemon MUST apply watcher create, update, enable, disable, and delete mutations without restart, and MUST cancel pending candidates from an obsolete or disabled configuration before they can dispatch.
- **FR-011**: The system MUST expose runtime health as `active`, `disabled`, or `degraded`, with a bounded actionable reason and last transition time; health MUST be derived from the current daemon runtime rather than persisted as historical truth.
- **FR-012**: Missing roots, replaced roots, permission loss, unsupported links, observer errors, and overflow MUST transition affected watchers to degraded health, discard pending candidates, and retry observation at a bounded interval.
- **FR-013**: Health transitions MUST produce one structured lifecycle event and one transition log record, while repeated identical failures MUST NOT create repeated events or log messages.
- **FR-014**: Recursive directory watchers MUST register existing and newly created real subdirectories without traversing symbolic-link directories or Windows junctions.
- **FR-015**: Network and removable paths MUST be handled as best effort and MUST report the observer's failure rather than claiming portable semantics that the host does not provide.
- **FR-016**: The local API MUST support watcher create, list, show, update, enable, disable, and delete operations with the standard error envelope and current health information.
- **FR-017**: The command line MUST support the complete watcher lifecycle with human-readable and JSON output, conventional exit codes, parseable duration input, and actionable field-specific errors.
- **FR-018**: The desktop Triggers view MUST present filesystem watchers separately from opaque-key triggers, show selection and health, and provide complete lifecycle actions with accessible controls and confirmation for deletion.
- **FR-019**: Watcher configuration mutations MUST preserve task activation readiness: an enabled watcher is an automatic activation source, while retargeting, disabling, deleting, or cascading deletion MUST disable a timeless task that no longer has any enabled automatic source.
- **FR-020**: Watcher paths MUST be cleaned and made absolute before persistence, and pattern validation MUST reject malformed patterns, path separators, and values that cannot represent a file base name.
- **FR-021**: Debounce and stability durations MUST each accept values from 25 milliseconds through 1 hour, defaulting to 250 milliseconds and 500 milliseconds respectively.
- **FR-022**: A watcher lifecycle operation or health query MUST complete within one second under a nominal load of 100 configured watchers, excluding filesystem and operating-system delays outside the application.
- **FR-023**: Every watcher-owned goroutine, timer, observer handle, and recovery loop MUST terminate when the daemon stops or the watcher configuration is replaced.
- **FR-024**: Migration from schema version 13 MUST be forward-only, preserve all existing data, add watcher persistence and watcher run provenance, and be proven against a real version-13 schema fixture.
- **FR-025**: Automated tests MUST document and verify the supported behavior on Windows, Linux, and macOS without assuming identical low-level event sequences.

### Key Entities

- **Filesystem Watcher**: A durable operator-defined mapping from a file selection rule to one target task, including timing controls and enabled state.
- **Watcher Runtime Health**: Ephemeral daemon-owned status for a watcher, including state, reason, and transition time.
- **Pending Candidate**: Ephemeral generation-bound settling state for one watcher and selected file.
- **Filesystem Run Provenance**: The run classification and watcher identifier that explain why a task was requested without retaining a potentially sensitive matched path.

## Success Criteria

### Measurable Outcomes

- **SC-001**: In 100 repeated write-storm trials, each trial produces exactly one accepted run request for the selected file after both timing conditions are satisfied.
- **SC-002**: In 100 temporary-file rename trials on each supported platform test environment, every supported rename-into-place workflow produces one accepted run request and no premature request.
- **SC-003**: A missing or replaced local root reaches degraded health within 2 seconds and reaches active health within 5 seconds after it becomes observable again under nominal load.
- **SC-004**: Restart tests produce zero replayed run requests for changes made while the daemon runtime is stopped.
- **SC-005**: Create, update, state change, delete, and health-read operations complete within one second at 100 configured watchers in deterministic test conditions.
- **SC-006**: Repeated identical observer failures produce one health transition event and one transition log record until health changes again.
- **SC-007**: The complete watcher lifecycle can be performed independently through both the command line and desktop interface, including identifying and recovering a degraded watcher.
- **SC-008**: The full race-enabled suite reports no data races and shutdown tests report no retained watcher goroutines or observer handles.
- **SC-009**: All core packages retain at least 80 percent test coverage and every canonical project verification gate passes.

## Assumptions

- The existing local daemon IPC authorization boundary remains the only administration boundary; this slice introduces no network listener or remote delivery.
- A directory pattern applies to file base names at any selected depth, which keeps recursion and matching independent and portable.
- Filesystem notifications are hints rather than a durable event log, so restart and degraded periods intentionally provide at-most-once observation with no replay claim.
- File readiness is determined by stable size and modification time after a quiet period; applications requiring locks, checksums, or domain-specific completion markers should encode those rules in the invoked task.
- A fixed two-second recovery interval provides prompt recovery while keeping missing or unavailable roots from becoming a busy polling loop.
- Existing task overlap policy bounds multiple settled files that target the same task; the watcher runtime does not invent a second concurrency policy.
- Full filesystem paths are operator configuration and appear in authorized watcher administration output, but matched paths are omitted from durable run provenance and transition logs.

## Dependencies and Scope

- Depends on completed external trigger dispatch and desktop trigger administration from #132 and #133, Trigger Sets from #134, and trigger-ready task authoring from #129.
- Closes #135 and completes the final implementation child of #17 when all acceptance criteria pass.
- Remote event delivery, arbitrary payloads, direct Chain invocation, durable filesystem event replay, content inspection, and operating-system-specific shell hooks are out of scope.
