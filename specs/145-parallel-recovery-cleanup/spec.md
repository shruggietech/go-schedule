# Feature Specification: Correlated parallel recovery and cleanup

**Feature Branch**: `codex/145-parallel-recovery-cleanup`

**Created**: 2026-11-10

**Status**: Implemented

**Delivery**: Implemented on review branch `codex/145-parallel-recovery-cleanup`, official PR #480: #477 with #478/#479/#481 corrections. Corrected complete frozen eight-gate acceptance passed, with actual Windows/Linux seven-case recovery execution and 39 headless browser cases on unchanged frontend source; [verification](verification.md) and [delivery](delivery.md). One of at most two Codex reviews used at local acceptance; final hosted/PR review facts are recorded on the PR. Stop for owner final review and merge.

**Input**: Extend the accepted S144 parallel block with outcome-based recovery and final cleanup across all maintained interfaces.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Recover one exact parallel instance (Priority: P1)

An operator configures Prepare to start Export A and Export B, Publish after all succeed, and Recover for selected unsuccessful branch conditions. The follow-up belongs to this Prepare execution, rather than the latest results of those tasks.

**Why this priority**: A single failed export needs one correlated recovery after all siblings finish, not an unrelated task failure chain or an early cleanup race.

**Independent Test**: Complete one branch unsuccessfully while its sibling remains active, then finish the sibling. No early follow-up occurs; exactly the configured recovery becomes eligible for that instance.

**Acceptance Scenarios**:

1. **Given** all branches succeed, **When** the last branch settles, **Then** Publish is selected once and Recover is not selected.
2. **Given** a selected unsuccessful condition and an outstanding sibling, **When** the sibling settles, **Then** Recover is selected once and Publish is not selected.
3. **Given** unsuccessful conditions not selected by the recovery policy, **When** every branch settles, **Then** neither follow-up runs and the unresolved or unselected reason is visible.
4. **Given** overlapping entries and repeated references to one branch task, **When** results interleave, **Then** each exact instance makes its own attributable decision.

### User Story 2 - Finish with cleanup (Priority: P2)

The operator optionally adds Cleanup, which runs after every branch and the chosen Publish or Recover follow-up has finished or definitively cannot start. Cleanup still matters if the chosen follow-up fails, is canceled or becomes ineligible.

**Why this priority**: One final step can remove temporary work without overlapping still-running commands.

**Independent Test**: Settle branches, hold the selected follow-up active, then finish it unsuccessfully. Cleanup waits for it, becomes eligible once, and retains the actual branch/follow-up outcome separately.

**Acceptance Scenarios**:

1. **Given** a selected follow-up still running, **When** it has not terminated, **Then** cleanup remains waiting.
2. **Given** the selected follow-up has terminated or explicitly cannot start, **When** cleanup is configured and currently eligible, **Then** cleanup is admitted once regardless of that follow-up outcome.
3. **Given** no follow-up is selected, **When** all branches are settled, **Then** configured cleanup may proceed with the unselected follow-up reason retained.
4. **Given** cleanup fails or cannot start, **When** the instance is inspected, **Then** the exact cleanup result or resolution is visible without rewriting earlier outcomes.

### User Story 3 - Configure and operate the same complete flow everywhere (Priority: P2)

An authorized operator authors the recovery conditions and cleanup, previews routing order, inspects exact runs, changes future configuration and round-trips portable definitions using desktop, CLI, JSON and controlled agent interfaces.

**Why this priority**: A backend-only addition leaves the daily workflow incomplete.

**Independent Test**: Author and inspect the example using maintained local/remote surfaces and round-trip its complete configuration without runtime history or text truncation.

**Acceptance Scenarios**:

1. **Given** a complete definition, **When** a target-only edit or portable remap occurs, **Then** complete names, branch identities, selected conditions and all task references survive.
2. **Given** a stale revision or changed selected daemon, **When** an edit or late read completes, **Then** no stale or wrong-target mutation is accepted silently.
3. **Given** no new recovery/cleanup options, **When** an existing join executes, **Then** its established all-success behavior remains unchanged.

### Edge Cases

- Failure, successful cancellation, overlap skip and definitive no-start resolution remain distinct selectable causes; queued observations and uncertain/interrupted work are not terminal success or definitive no-start.
- A failed cancellation request is not proof that the native command terminated. Late actual native completion continues to own settlement.
- Restart can replay interrupted external work under the existing at-least-once boundary. No inferred absence of history permits premature cleanup.
- Concurrent duplicate branch/follow-up completions, same-task slots, stale edits, disabled/deleted targets or groups, command removal, overlap collapse/queue and capacity delays retain exact attribution and current admission checks.
- Deleting a definition or a referenced task resolves future automatic admissions without falsely terminating already-started native commands or creating new recovery/cleanup after invalidation.
- A recovery target or cleanup cannot introduce an automatic execution cycle in the combined chain/join graph.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: An operator MUST be able to add an optional recovery task, explicit unsuccessful-branch conditions and optional cleanup task to one existing parallel definition, retaining its successful target.
- **FR-002**: Route selection MUST wait until every owned branch invocation has an actual terminal result or definitive no-start resolution. Pending, claimed, queued, active and uncertain work MUST not count as settled.
- **FR-003**: All successful branch results MUST select the existing success target. Any unsuccessful branch MUST prevent that success route; recovery MUST be selected only when at least one unsuccessful cause matches the configured conditions.
- **FR-004**: Failure, cancellation, skip and definitive no-start causes MUST remain distinct. Recovery with an omitted condition list defaults to failure; an explicit empty list matches no cause; users MUST be able to include the other documented unsuccessful causes explicitly.
- **FR-005**: Configured cleanup MUST wait for all branches and the selected follow-up to settle. If no follow-up is selected, or the chosen follow-up cannot start, its explicit disposition MUST allow cleanup without fabricating a native run.
- **FR-006**: One exact instance MUST admit at most one selected follow-up and one cleanup delivery, with stable separate stage and actual native run identities, despite concurrent/duplicate observations and repeated task references.
- **FR-007**: Durable branch settlement, chosen route, recovery and cleanup decisions MUST survive current-state reopen and preserve the existing at-least-once external-command crash boundary.
- **FR-008**: Every new automatic stage MUST recheck current relationship ownership, task/group/command eligibility, authority and overlap/capacity state at actual admission boundaries. Failed admission MUST have an attributable resolution.
- **FR-009**: Edits MUST affect future instances only; definition/task deletion MUST resolve outstanding automatic decisions without reviving them through late native results. Already-started native outcomes MUST remain truthful.
- **FR-010**: Combined chain/join relationship validation and atomic configuration replacement MUST reject real automatic cycles, dangling references and invalid routing configuration without partial writes.
- **FR-011**: Complete authoring and exact observation MUST be available through human/structured CLI, local and enrolled remote JSON interfaces, controlled MCP and desktop Automation/Activity with existing authority, audit and daemon identity rules.
- **FR-012**: Portable definitions MUST retain recovery/cleanup configuration, remap every new task reference, include it in drift/final graph validation and exclude all execution history.
- **FR-013**: Desktop forms MUST explain routing order and recovery defaults, expose exact waiting/selected/unselected/resolved results, support keyboard use and preserve current-target fences, revision protection, uncertainty and duplicate-submit guards.
- **FR-014**: Writable text MUST remain complete during editing and sparse updates, while observation text, collections and nested branches remain bounded and cancelable with explicit continuation/truncation.
- **FR-015**: Existing joins without the new options, ordinary completion chains, manual execution, exact-run cancellation and independent scheduling MUST retain their established behavior.
- **FR-016**: Current-only state initialization/reopen and explicit refusal of obsolete state MUST remain non-mutating; no historical migration or silent reset is required.
- **FR-017**: Documentation and all outputs MUST distinguish accepted/running/settled work, native cancellation versus cancellation-request failure, source capability versus released packages, and known runtime evidence versus unverified native UI/network behavior.

### Key Entities

- **Parallel definition**: existing entry/branches/success target plus optional recovery target, recovery condition selection and optional final cleanup target.
- **Exact instance snapshot**: immutable configured intent for one successful entry run, owned branch invocations and separately attributable follow-up/cleanup stages.
- **Settlement and route decision**: actual terminal or no-start evidence, selected success/recovery/none route, reasons and unique stage admissions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The Prepare / two exports / Publish-or-Recover / Cleanup example produces the configured order and no premature recovery or cleanup in every defined success, failure, cancellation, skip and no-start scenario.
- **SC-002**: Two overlapping entries with repeated branch-task references produce independent routes, at most one selected follow-up and at most one cleanup admission per instance.
- **SC-003**: An operator can author, inspect and round-trip all routing configuration through every maintained interface without losing names, identities, conditions or references.
- **SC-004**: Current-state reopen and duplicate observations preserve exact prior decisions and visible runtime uncertainty without false completed work.
- **SC-005**: Existing nominal dispatch latency remains below the documented 100ms p99 target; material paired scheduling regressions receive the constitution-required recorded justification.

## Assumptions

- Builds on accepted main cd4eb94f41e98383f6505099f15de84205b98d13 and completed #470/S144. Optional new fields preserve existing behavior by default.
- The daemon/store delivery architecture and current task authority remain authoritative; no expression language, loops, retries, separate workflow service or visual-canvas framework is required.
- Cluster work under #20/#186-#189/#146 remains owner-deferred. Released v2.0.3 comparison/release bytes remain historical, and no release publication is part of S145.
- Verification uses hidden unattended native processes and synthetic headless browser views. It never acquires desktop focus or requests operator interaction.

## Clarifications

### Session 2026-11-10

Autopilot resolved the following from #477, the accepted delivery architecture and owner-directed scope, without operator questions.

- Q: Does a blocked success route imply that all siblings are finished? A: No; actual terminal or definitive no-start evidence is required for each owned invocation before recovery or cleanup.
- Q: Which unsuccessful conditions enable recovery by default? A: Failure only when recovery is configured with an omitted list; an explicit empty list matches nothing. Cancellation, skip and no-start remain separate opt-in conditions. Any configured matching unsuccessful cause may select recovery after all branches settle.
- Q: What if no follow-up is selected or the selected task cannot start? A: Retain that explicit disposition and admit configured cleanup after all branches settle, without fabricating a run.
- Q: Does cleanup run after a failed or canceled selected follow-up? A: Yes, after actual native termination or definitive no-start; a failed cancellation request or uncertain command is not termination.
- Q: Do deletion and edits change already-owned work? A: Edits are future-only; deletion resolves future admissions and preserves actual late native evidence. No extra human interaction or release operation is introduced.
