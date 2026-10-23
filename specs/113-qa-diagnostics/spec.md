# Feature Specification: S113 reliable unattended QA and actionable delivery diagnostics

**Feature Branch**: `codex/113-qa-diagnostics`

**Created**: 2026-10-23

**Status**: In Progress

**Delivery**: Pending implementation. Owner authorization covers the full bundle, push, official PR publication, CI and at most two Codex review rounds. Stop before merging; tags and release publication remain separate.

**Input**: Kick off S113 under autopilot with #355, #356, #357, #358 and PR #360. Use the installed Spec-Kit sequence, implement the full bundle, publish the PR, monitor CI and handle review feedback. Parent release coordinator: #354. Baseline: main 11c2a9a46d4610fe8a67903cea246050224bd606.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Trust an interrupted unattended evaluation (Priority: P1)

A maintainer can select a supported evaluation budget and receive a truthful, recoverable report when the evaluation finishes or its deadline stops execution.

**Why this priority**: #357 currently imposes an unrelated outer deadline and can leave stopped work labelled running.

**Independent Test**: Select a nondefault bounded budget, interrupt an owned evaluation and inspect its retained outcomes, remaining cases, cleanup and recovery instructions.

**Acceptance Scenarios**:

1. **Given** a supported selected budget, **When** evaluation starts or resumes, **Then** the entry point honors that budget and a documented bounded cleanup/reporting allowance.
2. **Given** executed and unfinished cases, **When** a deadline stops work, **Then** completed outcomes survive, unfinished work remains interrupted/blocked, and the report distinguishes active workers from unresolved cleanup.
3. **Given** retained interrupted evidence, **When** recovery is requested, **Then** original source, binaries and budget identities are preserved; a fresh attempt remains distinguishable.

### User Story 2 - Identify runtime and observation failures (Priority: P1)

A maintainer can distinguish application execution failures, interpreter crashes, prerequisite refusal, follow-up inspection failures and cleanup failures without losing useful observations.

**Why this priority**: #355 and #356 have distinct intermittent evidence whose cause cannot be inferred from passing repetitions.

**Independent Test**: Inject failures at intermediate observations and subprocess boundaries, inspect bounded diagnostics, then exercise sustained owned operations and the existing concurrent writer cases.

**Acceptance Scenarios**:

1. **Given** several successful observations followed by a failed request, **When** inspection ends, **Then** completed partial data and failed phase/request remain visible while complete-state assertions fail.
2. **Given** a command succeeded but follow-up inspection failed, **When** the report is read, **Then** the original command result and separate observation failure are clear.
3. **Given** an interpreter crash or unusable prerequisite, **When** a payload invocation fails, **Then** the original exit/output and bounded executable/version context distinguish that failure from application or writer failure.
4. **Given** a demonstrated project cause, **When** corrected, **Then** owned operations remain responsive within existing limits. An upstream-only cause has an explicit supported disposition and authoritative guidance; an unconfirmed cause remains open.

### User Story 3 - Read delivery results without contradictory error labels (Priority: P2)

An operator can understand successful mail-server acceptance, partial rejection, retries and genuine failures consistently in delivery history, the command line and the desktop.

**Why this priority**: #358 labels successful acceptance summaries as errors despite successful delivery execution.

**Independent Test**: Inspect complete acceptance, partial rejection, temporary failure, later success and historical delivery records across supported interfaces.

**Acceptance Scenarios**:

1. **Given** all recipients accepted, **When** delivery evidence is shown, **Then** neutral acceptance diagnostics remain available without a success summary labelled as an error.
2. **Given** rejected or retrying recipients, **When** evidence is shown, **Then** actionable failure information and individual recipient states remain distinguishable.
3. **Given** older saved history or an existing client, **When** the new version reads or serves delivery results, **Then** compatibility handling is explicit and useful history is preserved.

### User Story 4 - Build current documentation on Windows paths (Priority: P3)

Maintainers retain a coherent website dependency baseline that includes the existing small documentation-tooling patch proposal.

**Why this priority**: PR #360 is a bounded patch update, including upstream corrections for Windows cross-drive imports.

**Independent Test**: Build the production documentation with the proposed dependency and matching lockfile, preserving existing references, static export and release discovery.

**Acceptance Scenarios**:

1. **Given** PR #360's proposed updates, **When** incorporated, **Then** all three patches and their lockfile changes are accounted for without unrelated dependency churn.
2. **Given** unpublished v1.6.0 source, **When** public discovery is generated, **Then** it continues to distinguish source preparation from the actual published release.

### Edge Cases

- Budget is invalid, at the supported maximum, omitted on resume, or differs from the original attempt.
- A worker exits while shutdown/reporting runs; its last durable checkpoint is incomplete, corrupt or unavailable.
- A later observation fails after useful data was collected; a command used most of the case budget; protocol cleanup also fails.
- An executable exists but cannot run, emits malformed version text, crashes, or has an unsupported version. Explicit paths retain their existing precedence and refusal semantics.
- Concurrent payload writers include failures: exactly one distinct record remains required per successful invocation, with six writers retained.
- SMTP acceptance follows a temporary failure, only some recipients accept, submission acknowledgment is uncertain, or historical success text was saved in the older error field.
- Diagnostic text contains secrets, excessively long output, or machine-specific details that must stay bounded and redacted.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Honor the selected validated evaluation budget for capture/resume with a bounded explicit cleanup/reporting allowance (#357).
- **FR-002**: Preserve terminal interruption status, executed outcomes, blocked remainder, worker/cleanup facts and actionable recovery inputs without rewriting original attempt identity (#357).
- **FR-003**: Retain partial before/after observations and bounded phase/request timings; incomplete observations cannot satisfy complete-state assertions (#356).
- **FR-004**: Keep command execution, independent inspection, protocol checks and cleanup failures distinguishable in machine and human reports (#356).
- **FR-005**: Investigate and correct demonstrated owned-operation stalls under existing limits, preserving original failure evidence and separately identifying diagnostic experiments (#356).
- **FR-006**: Distinguish interpreter/runtime crashes, prerequisite refusal and application writer failures with original outcomes and bounded executable/version diagnostics (#355).
- **FR-007**: Correct demonstrated project defects or document a supported upstream disposition with authoritative references; unconfirmed intermittent causes remain explicitly incomplete (#355/#356).
- **FR-008**: Preserve six concurrent writers, distinct literal identities, exact successful record counts, existing deadlines, resource limits, ownership checks and authentication behavior (#355/#356).
- **FR-009**: Present neutral delivery diagnostics separately from actual errors across persisted history and shared command-line/desktop interfaces, with explicit compatibility behavior (#358).
- **FR-010**: Preserve recipient-aware acceptance, rejection, retry, uncertainty and credential meanings; acceptance never claims final inbox receipt (#358).
- **FR-011**: Incorporate PR #360's three documentation-tooling patches and matching lockfile while preserving production website and release-discovery behavior.
- **FR-012**: Keep evaluation hidden, noninteractive, isolated and bounded; preserve secret redaction and owned cleanup without changing installed services, account sessions, public DNS or workstation runtimes.
- **FR-013**: Keep issue-level delivery accurate, reconcile S112 merge metadata, and leave #354 open for actual public release. Publishing this PR does not authorize merging or releasing.

### Key Entities

- **Evaluation attempt**: exact source/binary/budget identity, progress, terminal outcome and recovery state.
- **Observation**: completed partial data, completeness, failed request and bounded phase timing.
- **Process outcome**: executable identity, version/prerequisite context, original exit/output and cleanup facts.
- **Delivery evidence**: lifecycle state, neutral diagnostic, actual error, recipient outcomes and compatibility representation.
- **Dependency proposal**: existing PR identity, selected patches, lockfile and replacement disposition.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every supported selected budget is honored without a shorter unrelated wrapper cap, and every retained terminal report identifies remaining work and cleanup accurately.
- **SC-002**: Every failed observation retains completed partial data and identifies its phase/request; no incomplete observation is reported as complete.
- **SC-003**: Runtime, prerequisite, command, observation and cleanup failures have distinguishable actionable diagnostics without loss of original outcomes.
- **SC-004**: Fully accepted SMTP delivery has a neutral success explanation and no success text presented as an error; partial failure/retry/history meanings remain intact.
- **SC-005**: Every one of the four issues and the dependency proposal has an explicit functional disposition, with unconfirmed causes retained as open rather than inferred fixed.

## Assumptions

- Existing supported budget bounds and the 15-second case limit remain unchanged; reporting/cleanup gets an explicit bounded reserve rather than weakening case assertions.
- Existing schema/history and report consumers require an explicit compatible evolution; destructive migration or silent removal is excluded.
- Real native service/installer operations remain confined to disposable hosted CI. Local evaluation uses owned Windows/WSL fixtures and loopback services.
- Conditional clustered scheduling and actual v1.6.0 tagging/promotion remain outside this source slice.

## Clarifications

### Session 2026-10-23

- Q: Does capture get a fresh full budget for each environment? A: Preserve the existing selected total-attempt deadline; derive the outer wrapper bound from it plus an explicit bounded reporting/cleanup reserve, without multiplying case or attempt limits.
- Q: May incomplete observations count as successful state checks? A: No. Retain their partial data and phase/request diagnostics, while complete-state assertions remain failed.
- Q: Can a passing repetition close an intermittent failure? A: No. Deliver the demonstrated source correction or supported upstream disposition; unresolved causes remain separately open with exact evidence.
- Q: May new diagnostic fields invalidate retained v1 reports or saved delivery history? A: Use compatible additive optional report fields and explicit delivery-history compatibility handling; preserve original evidence and forward-only storage.
- Q: Is another pre-push authorization required? A: No. This kickoff explicitly authorizes push, the official PR and verified in-scope fixes through at most two Codex rounds. Stop before merge or release publication.
