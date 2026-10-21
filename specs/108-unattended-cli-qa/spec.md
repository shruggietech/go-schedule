# Feature Specification: Unattended CLI evaluation and canonical QA reports

**Feature Branch**: `codex/108-unattended-cli-qa`

**Created**: 2026-10-21

**Status**: Implemented

**Delivery**: review branch codex/108-unattended-cli-qa, locally verified and unpublished; [delivery record](delivery.md) and [verification](verification.md).

**Input**: Owner kickoff: S108 under autopilot, bundling #289 and #290 in full; parent #277. Baseline main: `57461fc`. Run autonomously and halt before publication.

## Clarifications

### Session 2026-10-21

- Q: Does unattended include semantic review without another AI service? A: Yes. The executing agent captures, reviews the immutable corpus, then finalizes through the coordinator; no human or external model credentials are needed.
- Q: Can discovered application defects be reported while delivering the evaluator? A: Yes. Record failing cases and independently actionable findings, file/update distinct product defects, and preserve accurate outcomes rather than weakening expectations.
- Q: What happens when prerequisites or resource controls are unavailable? A: Preserve every planned case as blocked with actionable reasons and complete reports; never install, elevate or reconfigure the host.
- Q: Must resume accept changed subjects or evidence? A: No. Reuse requires identical subject/manifest identity and verified retained artifacts; otherwise start a new run.
- Q: Should Windows and WSL run concurrently? A: No. Default to sequential environments and low concurrency within each, matching the issue's resource objective.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Evaluate the CLI without disturbing the host (Priority: P1)

A maintainer gives an agent a single runbook and receives evaluation of the actual binaries on Windows and a selected existing WSL environment, including intentional bad inputs, with isolated data and quiet execution.

**Why this priority**: Broad repeatable application evaluation is the requested capability; containment enables unattended use on a running computer.

**Independent Test**: Evaluate fixture binaries against a disposable daemon and inspect the case ledger, exact owned cleanup and unchanged installed state.

**Acceptance Scenarios**:

1. **Given** native and selected WSL prerequisites, **When** the run starts, **Then** both environments execute planned cases sequentially with recorded resource budgets and no further interaction.
2. **Given** service, GUI-launch, credential or preference cases, **When** evaluation runs, **Then** scoped adapters and disposable stores exercise handlers without accessing installed services, user credentials or desktop sessions.
3. **Given** an unavailable environment, deadline or interruption, **When** execution stops, **Then** every planned case has a truthful result, evidence survives and only owned resources are cleaned.

### User Story 2 - Detect coverage drift and defects (Priority: P1)

A maintainer can see every command, alias, flag and meaningful behavior class represented in a versioned inventory; changed metadata surfaces uncovered work.

**Why this priority**: Silent omissions would invalidate comprehensive review.

**Independent Test**: Change command metadata or remove a case and observe drift or completeness rejection.

**Acceptance Scenarios**:

1. **Given** the actual command tree, **When** inventory is captured, **Then** leaves, aliases, inherited/local flags, generated help/completion and output modes map to stable cases.
2. **Given** valid and invalid input, **When** a case executes, **Then** exit, output, state, history, audit and logs determine the result; expected rejection is a passing negative case.
3. **Given** the full scenario manifest, **When** it runs, **Then** each required functional group and meaningful lifecycle/error class is exercised with its execution boundary identified.

### User Story 3 - Review communication without another model service (Priority: P1)

The executing agent reviews frozen help, descriptions, normal/error/verbose output and logs for a layperson, then submits structured reviews for finalization.

**Why this priority**: Scenario-aware semantic review cannot be replaced with a text-presence check.

**Independent Test**: Capture a corpus, submit reviews for its immutable IDs and finalize; reject missing or stale reviews.

**Acceptance Scenarios**:

1. **Given** captured communication, **When** capture completes, **Then** bounded redacted messages are frozen and associated with case/channel IDs.
2. **Given** complete review records, **When** finalization runs, **Then** rubric verdicts and reproducible findings appear in machine and human reports.
3. **Given** missing reviews, **When** reporting runs, **Then** overall status remains blocked with pending_review.

### User Story 4 - Consume consistent reports and file actionable issues (Priority: P1)

A human reads a fixed-outline report and a subsequent agent files independently closeable defects from exact source identity, reproduction and evidence.

**Why this priority**: Predictable artifacts are the other half of the bundle and must agree with execution.

**Independent Test**: Validate and render reports, then reject modified hashes, duplicate IDs, missing cases, escaping paths and unsupported pass claims.

**Acceptance Scenarios**:

1. **Given** a completed or interrupted run, **When** artifacts are written, **Then** report.json, deterministic report.md, events.jsonl and a redacted evidence manifest exist.
2. **Given** a finding, **When** another agent reads it, **Then** identity, preconditions, invocation, fixture reproduction, expected/actual behavior, impact and bounded evidence suffice for triage.
3. **Given** unchanged canonical data, **When** rendering repeats, **Then** Markdown is identical and all counts agree.

### Edge Cases

- Missing tools/distro/limit capabilities, occupied or unwritable paths, malformed/truncated input, unknown commands/flags, invalid enums/ranges/timezones/certificates, unreachable daemons and stream cancellation.
- Unicode/quoting, empty results, pagination, duplicate/cyclic relationships, DST/month boundaries, catch-up/overlap, task failure/output limits and webhook retry/timeouts.
- Reused/rotated/deleted trigger keys, stale credentials, denied capabilities, TLS/daemon identity mismatch, partial multi-system results and reconnects.
- Changed binaries/manifest on resume, incomplete events, parent death, orphan descendants, junction/symlink escapes, corpus tampering, excessive output and split-write secret redaction.
- No emitted message is explicit; missing evidence/review never becomes pass or a free-form skip.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deliver one unattended coordinator and versioned agent execution/issue-triage prompts for all #289/#290 outcomes.
- **FR-002**: Inventory the real command tree and map all leaves, aliases, meaningful local/inherited flags, output modes and help/completion to stable cases; surface drift.
- **FR-003**: Cover task/group lifecycle and edits/clears/inheritance; real benign child execution with stdin/env/cwd/quoting/Unicode/exit/output limits; human/cron/one-off/startup schedules, syntax/conversion/import/export/timezones/DST/month ends/catch-up/overlap/calendar/preview; chains/cycles; triggers/sets/reusable and rotated keys; bounded watchers; notifications/inheritance/webhooks/retry/history/popup logic; bundles/round trips/drift; actors/credentials/grants/audit/pairing/profiles/TLS/identity; MCP stdio/HTTP authority/confirmation; history/active/exact runs/alerts/log streams/summaries/cross-system failures; service/GUI handlers and durable preferences. Include #304/#306/#307/#308/#309 regression boundaries.
- **FR-004**: Use explicit packaged_cli, command_harness and deterministic_engine evidence kinds. Exercise all meaningful negative classes from #289 and correlate exit/output/state/history/audit/log expectations. Missing relevant evidence prevents pass.
- **FR-005**: Execute native Windows and a selected existing WSL distro/user sequentially with independent owned workspaces. Remain unelevated, hidden and noninteractive; preserve installed services/groups, sessions, real preferences/profiles/credentials and host configuration.
- **FR-006**: Enforce recorded process CPU/memory/process-count budgets, low priority, bounded parallelism, per-case/total deadlines, output/log/artifact caps and watchdogs. Bound Linux work independently inside WSL without changing global VM configuration.
- **FR-007**: Validate exact daemon/PID/IPC/path ownership before mutation; contain child trees and clean only owned resources. Preserve evidence and actionable cleanup failures.
- **FR-008**: Emit one result per planned environment/case using pass, fail, blocked or manifest-defined not_applicable. Preserve uncompleted and partial-environment outcomes; expected bad-input rejection can pass.
- **FR-009**: Freeze a deduplicated context-bearing message corpus for help/descriptions, normal/error/verbose output and logs. Redact secrets before persistence and explicitly represent no emitted message.
- **FR-010**: Support capture, agent review and finalize/validate/render without human interaction or separately configured model-service credentials. Rubric covers plain language, action/outcome, correction, severity, terminology, identifiers/paths/units and secret absence.
- **FR-011**: Missing reviews produce blocked/pending_review; reject stale corpus IDs, missing reviews and unrelated case references.
- **FR-012**: Publish versioned JSON Schema 2020-12 contracts with required fields, strict properties, explicit nulls/enums, UTC RFC3339 timestamps and integer counts/durations for reports, events, manifests and reviews.
- **FR-013**: Canonical reports include exact commit/dirty state/version/binary hashes, runner/manifest/budgets, environments, complete cases/counts, communication reviews, findings, artifacts and cleanup.
- **FR-014**: Findings include stable deduplication keys, severity/category, affected identity/platform/cases, preconditions, exact redacted invocation, fixture reproduction, expected/actual result, layperson impact, hashed evidence spans, functional acceptance suggestions and existing issue URL or null. Distinguish product and harness defects.
- **FR-015**: Validate schemas and semantic ID/reference integrity, path containment, artifact hashes/sizes/redaction, counts, timestamp ordering, scenario completeness and aggregate truthfulness.
- **FR-016**: Render the fixed seven sections from canonical data: identity/outcome, environment/budget, coverage, functional results, communication review, consistently shaped defect blocks, cleanup/artifact index.
- **FR-017**: Interrupted runs retain complete reports and resource status. Resume validates run identity, exact binaries/manifest and retained artifacts before reusing outcomes.
- **FR-018**: Runner never changes source or GitHub records. Separate triage prompt verifies remote, searches open/closed issues, deduplicates findings and uses project conventions.

### Key Entities

- **Manifest**: Versioned command metadata, stable scenarios, environments, applicability and expected checks.
- **Run**: Exact subject/runner identity, budgets, planned cases, events, state and cleanup ownership.
- **Case result**: Expected/observed behavior, boundary, status, evidence, reviews and findings.
- **Corpus/review**: Immutable message/context IDs and structured agent rubric verdicts.
- **Finding**: Independently actionable product/harness defect with reproducible evidence.
- **Report/artifact manifest**: Canonical outcome and bounded, hashed, redacted evidence with deterministic rendering.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every actual leaf and meaningful metadata/behavior class is represented or explicitly uncovered; command changes cannot silently retain a clean inventory.
- **SC-002**: Every planned environment/case has exactly one truthful outcome after completion/interruption, and totals/references agree.
- **SC-003**: The full routine requires no human interaction after agent kickoff, preserves installed state and obeys recorded limits.
- **SC-004**: Every relevant distinct message is reviewed or the run remains blocked; each defect supports direct evidence-based triage.
- **SC-005**: Unchanged canonical data renders identically; tampered, incomplete or contradictory data is rejected.
- **SC-006**: Both issues deliver all functional outcomes in one independently traceable source increment, remaining unreleased until publication.

## Assumptions

- Reuse S107 ownership guards, hidden launchers and scoped dependencies. Resource enforcement and the complete evaluator are S108 work.
- Existing selected WSL/toolchains are prerequisites; report missing capabilities rather than install software or change host configuration.
- The executing agent performs semantic review of the frozen corpus; the coordinator validates it without fabricating verdicts or requiring another model service.
- Product defects remain separate findings/issues; failing evaluated cases do not make a correctly delivered evaluator incomplete or become false passes.
- Existing canonical eight-gate verification remains mandatory. No push, merge, release or DNS change is authorized.
