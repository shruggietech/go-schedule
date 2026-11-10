# Feature Specification: Correlated parallel task joins

**Feature Branch**: `codex/144-parallel-task-joins`

**Created**: 2026-11-10

**Status**: Implemented

**Delivery**: Implemented on review branch codex/144-parallel-task-joins in PR #473. All 32 implementation tasks and the complete corrected eight-gate local run pass. Both permitted Codex reviews are consumed, and all three findings have functional corrections; final #476 changes are qualified by focused/full tests and fresh exact-head CI rather than a third Codex review. Final owner review and merge remain pending.

**Input**: Owner kickoff for S144 implements [#470](https://github.com/shruggietech/go-schedule/issues/470), prioritized from the maintained competitor matrix's dependent-workflow gap.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run a complete parallel block (Priority: P1)

An operator defines Prepare as the entry, Export A and Export B as required parallel branch invocations, and Publish as the joined target. A successful Prepare execution starts those branches; Publish becomes eligible only after both specific branch invocations succeed.

**Why this priority**: Current independent completion edges cannot express a successor waiting for several branches of the same execution.

**Independent Test**: Complete the two branches in either order and observe no premature Publish admission, followed by exactly one join target admission after both succeed.

**Acceptance Scenarios**:

1. **Given** a configured parallel join, **When** its entry execution succeeds, **Then** its required branch invocations are created as one attributable instance and dispatched through ordinary task eligibility and overlap policy.
2. **Given** one successful branch and one outstanding branch, **When** the operator inspects the instance, **Then** the waiting branch is named and the joined target has not been admitted.
3. **Given** all required branches succeeded, **When** the last result is recorded, **Then** the joined target is admitted once with source correlation.
4. **Given** two overlapping entry executions, **When** branch results interleave, **Then** no contribution crosses instance boundaries.

### User Story 2 - Author and inspect through every interface (Priority: P2)

An authorized operator configures, edits, inspects and removes parallel joins from the desktop, CLI, local/remote JSON interfaces and controlled MCP, and round-trips them through portable definitions.

**Why this priority**: Backend-only orchestration does not deliver the approved end-to-end capability.

**Independent Test**: Author the example, inspect a waiting/blocked/completed instance and round-trip its references through each supported surface without changing task execution inputs or authority.

**Acceptance Scenarios**:

1. **Given** existing tasks, **When** an operator saves a join definition, **Then** the entry, named required branches and target are unambiguous and invalid references/cycles receive actionable diagnostics.
2. **Given** Observe versus Manage authority and a selected daemon, **When** the client attempts a mutation or read, **Then** existing authority, identity and redaction rules apply.
3. **Given** portable export and import with remapped task identities, **When** the definition is imported, **Then** its entry, branches and target retain their intended relationships.
4. **Given** a keyboard user or a stale/unavailable target, **When** operating the desktop, **Then** controls remain accessible and uncertain or unavailable data is not presented as successful mutation.

### User Story 3 - Retain truthful state through failures and restart (Priority: P3)

An operator can explain why a join is waiting or blocked and recover durable progress after daemon restart without confusing deliberate cancellation with successful work.

**Why this priority**: A partial workflow that silently stalls or releases incorrect work is worse than the current explicit chains.

**Independent Test**: Record partial branch completion, reopen valid current storage, resume pending delivery and observe the same instance and target decision. Fail, skip, cancel or remove required work and inspect its explicit blocked/resolved reason.

**Acceptance Scenarios**:

1. **Given** a required branch fails, is skipped or canceled, **When** its result is recorded, **Then** the success join is blocked for that instance and the reason identifies the branch.
2. **Given** repeated observations or concurrent last-branch completions, **When** they are recorded, **Then** target admissions are deduplicated.
3. **Given** a definition is edited after an instance starts, **When** its remaining branches finish, **Then** the instance uses its original entry/branch/target bindings.
4. **Given** a definition or required task is removed, **When** outstanding work is reconciled, **Then** no new invalid target is launched and the instance has an explicit resolution.

### Edge Cases

- Several branches may reference the same task but remain distinct invocation identities; completion of one must not satisfy another.
- Branch delivery can be delayed by worker limits or queue-one policy; queued markers are not successful branch completion. Collapsed or superseded deliveries require an explicit resolution.
- The entry can be manual, scheduled, catch-up or event driven. Unsuccessful entry executions do not launch the success parallel block.
- Task or ancestor-group state can change between definition, admission and execution. Target admission rechecks current eligibility.
- A native command may perform external side effects before a crash prevents its completion record. Existing at-least-once delivery remains explicit.
- Canceling one branch does not disable its task, cancel unrelated work or accidentally run the joined target. Unconfirmed cleanup remains owned as in S143.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Operators MUST be able to define, edit, inspect and remove a reusable parallel block with an entry task, one or more explicitly identified required branch invocations and one target task.
- **FR-002**: A successful entry execution MUST create one durable instance per applicable definition, with exact immutable entry, branch and target bindings; the block itself launches its required branch invocations.
- **FR-003**: Required branch contributions MUST be tied to their exact instance and invocation identities, including repeated references to the same task. Ambient or latest task results MUST NOT satisfy the join.
- **FR-004**: The target MUST become eligible only after every required branch invocation succeeds, with one durable target admission per instance despite repeated/concurrent completion observations.
- **FR-005**: Branch and target work MUST use existing task/group eligibility, worker, overlap and delivery rules. Waiting, collapsed, skipped, superseded or missing work MUST have an explicit outcome, never implicit success.
- **FR-006**: Operators MUST observe instance states, exact source run correlation, each outstanding/finished branch, target admission and actionable blocked/resolved reasons.
- **FR-007**: Valid current storage MUST retain definitions, partial results and target decisions across restart; interrupted pending work MUST preserve the existing at-least-once external-execution contract.
- **FR-008**: Required branch failure, deliberate cancellation, unresolved cleanup or invalidated work MUST NOT release a success join. Existing single-run cancellation and ordinary success/failure/any chains MUST retain their documented semantics.
- **FR-009**: Definition edits MUST affect future instances without retargeting in-flight work. Deletion/invalidation MUST prevent future invalid dispatch and visibly resolve affected instances while retaining permitted historical attribution.
- **FR-010**: Reference and directed-cycle validation MUST include ordinary completion edges and parallel entry/branch/target relationships, with atomic mutation and actionable errors.
- **FR-011**: Desktop, human/JSON CLI, local/remote JSON API and controlled MCP MUST deliver equivalent authoring and bounded observation using existing authority, daemon identity, mutation audit, retry and redaction boundaries.
- **FR-012**: Portable definitions MUST round-trip join configuration and correctly remap task identities without exposing execution inputs beyond existing bundle contracts or importing runtime instances.
- **FR-013**: Desktop controls MUST be keyboard accessible, prevent duplicate/stale-target mutation, retain useful observations on failure and distinguish uncertain submission from authoritative success.
- **FR-014**: Existing independent tasks, task scheduling, completion chains, notifications, per-run cancellation, native command arguments/input/run-as and output limits MUST retain their current contracts.
- **FR-015**: Implementation and public evidence MUST be unattended, acquire no operator desktop focus and contain no private workstation, account, address or diagnostic values.
- **FR-016**: Documentation MUST distinguish the new source increment from public v2.0.3 and full visual-workflow parity; historical comparison/release records and owner-deferred cluster work MUST remain unchanged.

- **FR-017**: Automatic delivery queued or waiting for worker capacity MUST recheck current task/group eligibility and definition/instance validity immediately before execution (#471).
- **FR-018**: Shared Automation mutations, including ordinary and secret-bearing results, MUST preserve uncertain outcomes and prevent resubmission until authoritative same-target reconciliation (#472).
- **FR-019**: The full Automation workspace MUST load and operate for enrolled remote targets, including its existing chain collections, while retaining Observe/Manage separation and selected-target attribution.

### Key Entities

- **Parallel join definition**: Reusable entry, required branch identities and joined target, plus editable name and lifecycle metadata.
- **Join instance**: One entry execution and its snapshotted branch/target bindings, durable decision state and attributable reason.
- **Branch invocation**: Exact work identity within one instance, referenced task, pending/running/final delivery evidence and native run correlation.
- **Joined target admission**: Unique durable eligibility decision and resulting execution or resolution for the same instance.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The Prepare / two parallel Exports / Publish example completes end to end from supported authoring surfaces, with zero premature target admissions and one target admission per successful instance.
- **SC-002**: Two overlapping instances with interleaved results never share branch contributions or target decisions.
- **SC-003**: Restart after partial completion retains all committed progress and the correct remaining prerequisites, with no duplicate target admission.
- **SC-004**: Failure, cancellation, supersession or deletion produces an attributable visible waiting/blocked/resolved state instead of false successful continuation.
- **SC-005**: Existing ordinary task/chain scheduling and access boundaries remain usable while parallel joins are configured and running.

## Assumptions

- A parallel block owns its branch invocations, so the operator need not also configure independent fan-out edges to create those same branches. Existing independently configured chains remain separate work and do not count as its contributions.
- The initial capability joins successful branch results. Existing individual success/failure/any chains remain available; an arbitrary-expression language and separate visual-canvas framework are not prerequisites.
- The implementation uses the current prelaunch storage policy: initialize and reopen the current format, refuse obsolete/malformed state with explicit reset guidance, and never automatically convert or delete it.
- The owner authorizes this source branch/official PR and in-scope review fixes. Final merge, tagging and package release remain outside this kickoff.

## Clarifications

### Session 2026-11-10

- Q: Does a join observe arbitrary existing task history? A: No. The configured parallel block launches and owns explicit branch invocations after one successful entry execution, preserving unambiguous instance correlation. Existing independent chains remain separate configured work.
- Q: Can two branches use the same task? A: Yes, each branch has its own invocation identity and must complete independently.
- Q: Do definition edits retarget existing work? A: No, instance bindings are snapshotted. Deletion/invalidation resolves outstanding decisions explicitly without silently retargeting or canceling unrelated native executions.
- Q: What delivery guarantee is claimed after a crash? A: Durable decisions and unique target admissions, with the existing at-least-once external-command boundary preserved.
- Q: Are publication and native desktop interaction authorized? A: Push and official PR publication plus in-scope CI/review fixes are authorized, with two total Codex reviews. Final merge/release remain outside scope; all testing stays unattended and never acquires operator desktop focus.
