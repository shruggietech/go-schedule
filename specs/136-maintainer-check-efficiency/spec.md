# Feature Specification: Maintainer check efficiency

**Feature Branch**: `codex/136-maintainer-check-efficiency`

**Created**: 2026-11-08

**Status**: Implemented

**Delivery**: Local implementation on review branch codex/136-maintainer-check-efficiency; all eight foreground gates passed. Official PR and owner merge remain separate delivery steps.

**Input**: Owner kickoff of S136 under autopilot, automatic official PR publication and CI resolution, at most two Codex reviews including automatic review, stop before owner final review/merge. Outcome: #433, executable action-field correction #434, publication dependency correction #435 and pending-source exact-policy correction #436.

## User Scenarios & Testing

### User Story 1 - Evaluate repository policy efficiently (Priority: P1)

A maintainer runs the existing automation command against the checkout or an isolated fixture and gets the same policy decisions and useful failures with substantially less Windows process startup overhead.

**Why this priority**: S135's whole canonical run took 4355.058 seconds; isolated current measurements establish which checker costs this slice can actually reduce.

**Independent Test**: Existing positive and rejection fixtures retain their intended decisions and diagnostic context; the same checkout is timed before and after.

**Acceptance Scenarios**:

1. **Given** approved workflow, release and dependency configuration, **When** the maintainer runs the command, **Then** all existing policy categories are evaluated and success retains the eight-gate contract.
2. **Given** an invalid action pin, publication bypass or missing source binding, **When** the same command runs, **Then** it fails with file/rule context.
3. **Given** a commented or script-string action example, **When** action references are evaluated, **Then** only executable workflow action fields are audited; malformed or duplicate-key workflow configuration fails explicitly.

### User Story 2 - Evaluate specification state efficiently (Priority: P2)

A maintainer scans the whole specification inventory using the existing lifecycle command without launching text/path tools for each specification.

**Why this priority**: The 133-spec baseline takes 31.541 seconds on this host and is also included in automation.

**Independent Test**: Isolated lifecycle fixtures cover every existing state, missing/invalid delivery evidence, task counts, missing inventory and disposition requirements.

**Acceptance Scenarios**:

1. **Given** a consistent inventory, **When** evaluated, **Then** the command reports the specification count and succeeds.
2. **Given** inconsistent status, tasks, delivery or inventory, **When** evaluated, **Then** the command fails with the affected file and reason.

### User Story 3 - Maintain cross-platform execution (Priority: P3)

A maintainer uses the same commands and isolated roots on Windows, Linux or macOS, with hidden noninteractive console children on Windows and no new runtime dependency.

**Independent Test**: Existing shell entry points and the unchanged canonical gate manifest work on supported hosts; the synthetic corpus reuses a built checker.

**Acceptance Scenarios**:

1. **Given** resolved project tooling, **When** repeated fixtures run, **Then** policy evaluation does not compile for every fixture or launch a utility chain per assertion.
2. **Given** a child failure, missing input or invalid mode, **When** invoked, **Then** the original failure remains visible and the aggregate stops normally.

### Edge Cases

Missing directories/files; empty inventories; every lifecycle state; CRLF input; paths containing spaces; unknown/obsolete action versions; quoted/folded YAML; action-like comments and script strings; malformed/duplicate YAML keys; comparison checks later in the same run block; conditional or soft-failed source/comparison checks; unexpected children or helper failures.

## Requirements

### Functional Requirements

- **FR-001**: Preserve the existing automation and lifecycle command names, optional explicit root, success/failure status and actionable diagnostic context.
- **FR-002**: Preserve action/runtime pins, CodeQL, Dependabot, release-source preflight/order, comparison binding/order, release notes, brand, manifest and lifecycle policies. Audit actual action fields and reject malformed configuration instead of copying comment matching.
- **FR-003**: Evaluate assertions and specification rows without per-assertion/per-specification external text/path tools, and reuse one built checker for the fixture corpus.
- **FR-004**: Keep existing meaningful positive/rejection cases, add discriminating cases for changed parsing and lifecycle boundaries, and preserve failure propagation.
- **FR-005**: Keep Windows console children hidden and noninteractive, preserve the original child failure, avoid new runtime dependencies and keep already-resolved operation offline.
- **FR-006**: Document isolated before/after timings and their scope. Do not attribute all canonical duration to the checker, invent an overall time promise, alter watchdogs, skip gates or retry failures into success.
- **FR-007**: Preserve frozen release/QA evidence, application behavior, current-state policy and owner merge/publication boundaries.

### Key Entities

- **Evaluation**: One root, checker mode and ordered findings; input files are read once during that evaluation, never cached across subsequent fixture mutations.
- **Finding**: A path, optional source line and actionable policy reason; any finding yields nonzero completion.
- **Timing observation**: Source identity, host, command, elapsed duration and native result; historical aggregate observations remain distinct.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Both isolated existing checker commands complete faster on the same native Windows checkout than their retained unchanged baseline, with the affected categories and source identity stated.
- **SC-002**: Every existing meaningful fixture decision and failure category remains available; only separately recorded parser defects change classification.
- **SC-003**: Repeated fixture evaluation requires no per-fixture compilation and no per-assertion or per-specification text/path subprocesses.
- **SC-004**: Maintainers retain the eight canonical gates and the same documented entry points on Windows, Linux and macOS.

## Assumptions

This is maintainer tooling, not a product release. Existing brand and manifest helpers remain bounded trusted checkout operations. No application runtime or installed service is changed. #355/#356 and clustered scheduling stay outside this slice.

## Clarifications

### Session 2026-11-08

- Scope is the existing two policy checkers and their corpus, not a new QA framework. Full canonical verification remains required engineering evidence, separate from issue closure.
- Isolated unchanged baseline: lifecycle 31.541 seconds and automation 44.701 seconds (automation includes lifecycle). These cannot be added into independent canonical savings.
