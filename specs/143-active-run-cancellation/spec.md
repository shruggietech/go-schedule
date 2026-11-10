# Feature Specification: Active-run cancellation

**Feature Branch**: `codex/143-active-run-cancellation`

**Created**: 2026-11-09

**Status**: Implemented

**Delivery**: Implemented on review branch codex/143-active-run-cancellation and locally verified across all eight foreground gates, native owned Windows/Linux packaged cancellation and headless desktop coverage. Official PR publication, CI and in-scope fixes continue under the explicit kickoff, with at most two total Codex reviews including automatic review. Final owner review and merge remain pending; public v2.0.3 is unchanged. See [verification](verification.md) and [delivery](delivery.md).

**Input**: S143 kickoff implements [#465](https://github.com/shruggietech/go-schedule/issues/465), allowing an operator to stop an exact active execution through local/remote CLI, desktop and programmatic interfaces.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Stop one running job (Priority: P1)

An operator selects an active execution and stops its owned work without restarting the daemon or disturbing another execution.

**Why this priority**: A stuck command currently cannot be stopped through the scheduler's run controls.

**Independent Test**: Start two owned long-running executions, cancel one by its Run ID and observe that it ends while the other keeps running and the daemon accepts more work.

**Acceptance Scenarios**:

1. **Given** an authorized operator and an active execution with child processes, **When** the operator requests cancellation for its Run ID, **Then** the selected owned execution is terminated, its retained output remains available and the daemon reports its final cancellation result.
2. **Given** simultaneous executions of the same task, **When** one Run ID is canceled, **Then** only that execution is affected.
3. **Given** a completed execution or a repeated cancellation request, **When** cancellation is requested again, **Then** the actual existing result is retained and no later process or duplicate history record is affected.

### User Story 2 - Use the same control from every interface (Priority: P2)

An operator uses the command line, desktop, direct client or controlled agent with consistent target identity and permissions.

**Why this priority**: Local and enrolled-remote work must share one reliable operational contract.

**Independent Test**: Request cancellation through each interface and compare its acknowledged Run ID, eventual history result and authorization behavior.

**Acceptance Scenarios**:

1. **Given** an enrolled remote daemon and sufficient existing operating authority, **When** cancellation is requested, **Then** the request targets that daemon and exact Run ID and is audited using existing conventions.
2. **Given** an observer, wrong daemon identity, disconnected/stale desktop or revoked credential, **When** cancellation is attempted, **Then** the existing permission/connection rules refuse mutation without stopping work.
3. **Given** an accepted request, **When** the process has not finished cleanup yet, **Then** the interface distinguishes requested cancellation from confirmed termination and later presents the authoritative result.

### User Story 3 - Understand the resulting workflow state (Priority: P3)

An operator sees deliberate cancellation separately from natural success, ordinary command failure and unconfirmed cleanup.

**Why this priority**: Cancellation must not masquerade as a successful job or launch surprising failure-recovery work.

**Independent Test**: Cancel an execution with completion chains and notifications configured and inspect one terminal result with explicit cancellation semantics.

**Acceptance Scenarios**:

1. **Given** canceled work, **When** terminal history is read, **Then** its timestamps, retained output, capture-cap information and genuine available native exit information remain attributable to that Run ID.
2. **Given** an explicit cancellation, **When** completion handling runs, **Then** no success/failure/any completion-chain target is launched from that cancellation; future independent schedule/trigger requests retain their existing behavior.
3. **Given** cancellation cleanup cannot be confirmed, **When** the final result is shown, **Then** it explicitly reports the uncertainty/failure and does not claim successful termination.

### Edge Cases

- Cancellation arrives before command launch, during execution, while output drains, after natural completion or during final history persistence.
- Two callers cancel the same Run ID; a different execution of the same task is still running.
- The root exits but an owned descendant retains inherited output handles.
- A child exits naturally at the same time cancellation is requested; the native result and actual winning observation are retained.
- A caller loses the response, a terminal record is removed by existing retention/removal behavior, or a Run ID is unknown.
- Process ownership setup, termination, output drain or cleanup fails; each failure is surfaced rather than hidden by a canceled label.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Operators MUST be able to request cancellation for an exact currently executing Run ID without stopping the daemon.
- **FR-002**: Cancellation MUST target that execution's owned process tree, with a five-second termination/output cleanup allowance and explicit failure or unconfirmed-cleanup reporting on Windows, Linux and macOS. If the OS has not reaped owned work, it MUST remain visibly active with failed/unconfirmed cancellation and retain worker/ownership resources until actual exit; a timeout alone cannot create a terminal result.
- **FR-003**: Successful cancellation MUST produce one distinct canceled terminal result retaining available output, cap/truncation metadata, timestamps, correlation and genuine native exit information.
- **FR-004**: Cancellation MUST preserve unrelated executions, including another concurrent execution of the same task, and normal future scheduling/trigger behavior.
- **FR-005**: Repeated requests and completion races MUST be idempotent for the selected identity, preserve an already completed result and avoid reused-process targeting or duplicate terminal history.
- **FR-006**: Local/remote command line, desktop, direct JSON clients and controlled MCP access MUST expose the same request/final-result meaning, using existing target identity, operating authority and audit rules.
- **FR-007**: Desktop exact-record refresh MUST retain the selected run (#468). Desktop controls MUST be keyboard-accessible, identify the selected execution, distinguish pending from finished cancellation and respect existing disconnected/stale/observer mutation rules.
- **FR-008**: An accepted request MUST NOT be presented as proof of termination; interfaces MUST expose the authoritative eventual result and actionable missing/not-active/failure outcomes.
- **FR-009**: Deliberately canceled work MUST NOT dispatch completion chains for success, failure or any outcome. Existing run-finished events MUST carry the actual cancellation outcome; outbound success/failure/duration/recovery policies MUST not classify cancellation as ordinary job success/failure or change their failure/recovery streaks.
- **FR-010**: Process launch and collection MUST retain existing command arguments, environment, input, run-as and output-cap behavior, and hidden Windows console launches. Ordinary natural completion MUST NOT introduce new restrictions on commands or implicitly cancel unrelated/background work.
- **FR-011**: Valid current stores MUST retain durable reopen and cancellation history without automatic conversion/deletion. Reuse existing storage where possible; per-task runtime limits, retries, clustering and release preparation are outside S143.
- **FR-012**: Execution and public evidence MUST remain unattended, avoid operator desktop focus and exclude private workstation/account/network/diagnostic contents.

- **FR-013**: Desktop connections MUST negotiate the actually advertised supported local/remote API contract while preserving health, product-identity agreement, pinned daemon identity and access checks (#466).
- **FR-014**: Finishing or canceling one concurrent same-task execution MUST retain remaining in-flight accounting so skip and queue-one policies remain correct (#467).

### Key Entities

- **Active execution**: Stable Run ID, task identity, selected daemon, start/correlation information and ownership of its currently running work.
- **Cancellation request**: Exact Run ID, existing authorized actor and an acknowledgement that is distinct from the eventual execution result.
- **Terminal execution result**: One immutable observed completion with output/native-result metadata and an explicit natural/canceled/cleanup-failed distinction.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can stop one selected active execution from every supported control surface without restarting the daemon.
- **SC-002**: Successful cancellation ends its owned work within a five-second termination/collection allowance, preserves one final result and frees execution capacity.
- **SC-003**: In a simultaneous two-execution scenario, cancellation changes only the selected Run ID and produces no duplicate terminal history or downstream completion-chain execution.
- **SC-004**: All interfaces distinguish requested cancellation, canceled completion, already finished work and unconfirmed cleanup; denied requests leave execution unchanged.
- **SC-005**: The resulting history remains readable after restart, while current execution isolation, scheduling, permissions and output caps remain intact.

## Assumptions

- Existing active-run observation and enrolled-daemon access are reused. The selected run has begun acquiring execution resources; queued history markers are not active execution handles.
- Cancellation is an explicit operator action. It does not undo external side effects or promise to control intentionally detached/escaped work outside OS-enforced ownership. Such limits and cleanup failures must remain explicit.
- Natural completion wins when its terminal completion is already observed; a later cancellation request cannot rewrite it.
- The normal PR is the integration deliverable. Tags, package publication and automatic merge are not authorized by this kickoff.

## Clarifications

### Session 2026-11-09

- Q: Does cancellation mean a whole task is disabled? A: It targets one started Run ID; other executions and future independent triggers retain their existing scheduling semantics.
- Q: Which terminal result wins a completion race? A: Already observed natural completion is immutable; otherwise cancellation remains a request until actual termination/cleanup is observed.
- Q: Should deliberate cancellation dispatch completion chains? A: No success/failure/any chain is dispatched from a deliberate cancellation; cancellation remains distinct in run-finished notifications.
- Q: Do push authorization and review limits require another confirmation? A: The kickoff authorizes this branch and official PR plus verified in-scope fixes, with two total Codex reviews and owner final review/merge.

- Q: How are existing source defects incorporated? A: Confirmed #466 API negotiation and #467 concurrent accounting are directly required for this capability and are corrected in the same slice with separate bug traceability.
- Q: What happens if the operating system cannot finish termination within five seconds? A: Surface failed/unconfirmed cancellation promptly, retain active ownership and capacity until actual native exit, and preserve that failure in the eventual terminal result.
