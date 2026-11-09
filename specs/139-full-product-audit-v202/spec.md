# Feature Specification: Full-product audit and v2.0.2 maintenance delivery

**Feature Branch**: `codex/139-full-product-audit-v202`

**Created**: 2026-11-08

**Status**: Implemented

**Delivery**: Review branch codex/139-full-product-audit-v202; [delivery](delivery.md), [verification](verification.md); owner final review and merge pending

**Input**: Owner-approved S139, tracked by [#444](https://github.com/shruggietech/go-schedule/issues/444): audit the whole current product, all CLI surfaces, actual Windows/WSL2 interoperability and actual isolated desktop views where feasible; sanitize public evidence, file every actionable finding and deliver compatible corrections toward v2.0.2. Run unattended without additional confirmation, with no operator-desktop focus acquisition and only Windows administrator approval popups permitted. Push/open the official PR, monitor CI and resolve findings within two total Codex reviews including automatic review; stop for owner final review/merge. Tagging and promotion remain separate publication boundaries.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand the complete current product (Priority: P1)

A maintainer receives an exact-version assessment of every command surface and major product workflow on Windows and Linux, with actionable defects and explicit coverage gaps.

**Why this priority**: A fresh audit must assess actual current behavior rather than inherit historical passing claims.

**Independent Test**: Reconcile the current command inventory and product-area register with an assessment that names every executed, failed, blocked or inapplicable surface and its evidence boundary.

**Acceptance Scenarios**:

1. **Given** fixed source and candidate identities, **When** the audit completes, **Then** every shipped or documented maintainer command has an explicit disposition, including behavior, human/structured output, help and failure handling.
2. **Given** a demonstrated defect, **When** the assessment is published, **Then** it maps to a deduplicated GitHub issue with functional acceptance criteria and dependencies.
3. **Given** a correction, **When** subsequent observations are reported, **Then** original failures remain attributable and fresh results identify their corrected source and binaries.

### User Story 2 - Observe real cross-platform operation (Priority: P1)

A maintainer sees actual Linux-to-Windows and Windows-to-Linux remote control using the existing local WSL2 environment and Windows 11 host.

**Why this priority**: Separate platform runs do not establish cross-platform transport, identity or execution behavior.

**Independent Test**: Each direction establishes a real authenticated connection, changes an audit-owned object and verifies execution/history on the intended host.

**Acceptance Scenarios**:

1. **Given** disposable owned daemons on both hosts, **When** each opposite-OS client connects, **Then** enrollment, authorization, task operations, persistence and recovery outcomes are independently attributable.
2. **Given** revoked authority or an identity mismatch, **When** the client acts, **Then** the documented refusal is preserved and no unintended target mutation occurs.

### User Story 3 - Inspect the desktop without interrupting the owner (Priority: P1)

The audit attempts actual desktop interaction and screenshots only where isolated display operation can be established before launch.

**Why this priority**: The owner is working in parallel and explicitly prohibits any desktop-focus acquisition.

**Independent Test**: Native evidence identifies its real platform/build/display boundary, while unavailable isolation yields a literal blocked result without launching in the operator session.

**Acceptance Scenarios**:

1. **Given** a verified isolated display, **When** the actual application is exercised, **Then** screenshots and workflow observations use synthetic data and cannot acquire the operator desktop's focus.
2. **Given** unavailable isolation, **When** UI coverage is assessed, **Then** native coverage is explicitly blocked and headless frontend observations remain separately labeled.
3. **Given** any unavailable prerequisite or declined elevation, **When** execution continues, **Then** no human walkthrough, password entry, console prompt or focus-taking workaround is requested.

### User Story 4 - Receive coherent maintenance delivery (Priority: P2)

The maintainer receives compatible corrections, accurate release preparation and an official reviewable PR toward v2.0.2.

**Why this priority**: Findings should lead to concrete functional corrections and an honest release path.

**Independent Test**: Trace each finding to its issue and source disposition, and inspect the PR's exact candidate, verification and remaining publication boundary.

**Acceptance Scenarios**:

1. **Given** compatible defects, **When** source delivery is complete, **Then** each delivered outcome is traced independently and uncompleted outcomes remain open.
2. **Given** final source acceptance, **When** release work continues, **Then** package identity, staging and public promotion remain distinct from audit observations and owner merge.

### Edge Cases

- Missing native UI isolation, credentials, toolchain, delegated resource controls or transport reachability.
- Expected refusal versus actual product/harness defect, interrupted capture and incomplete cleanup.
- Identical display names on different daemons, stale profiles, revoked credentials and uncertain mutations.
- Sensitive data in errors, reports, filenames, screenshots, traces or source/environment metadata.
- Source or binaries changing between capture and review; old failures must not be replaced.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The audit MUST bind each observation to exact source, binary, platform and execution-boundary identities.
- **FR-002**: Every current CLI command/alias/flag and documented maintainer entrypoint MUST have explicit coverage or a specific blocked/inapplicable disposition; help alone MUST NOT substitute for functional behavior.
- **FR-003**: The assessment MUST cover architecture, dependencies, documentation, scheduling/timezones, storage/recovery, execution, IPC/access, remote control, MCP, notifications, desktop and installation.
- **FR-004**: Actual WSL2-to-Windows and Windows-to-WSL2 remote operations MUST be attempted and independently verified using supported transport and audit-owned state.
- **FR-005**: Actual desktop interaction/screenshots MUST be attempted only after display isolation is established; frontend adapters MUST remain distinct from native evidence.
- **FR-006**: Execution MUST be unattended and MUST NEVER acquire the operator desktop's focus; only Windows administrator approval popups may request operator interaction.
- **FR-007**: Missing prerequisites MUST remain explicit while independent work continues, without manual handoffs, altered expected outcomes or unapproved host reconfiguration.
- **FR-008**: Public evidence MUST exclude PII and secrets; private originals and sanitized derivatives MUST retain separate provenance and integrity.
- **FR-009**: Every actionable finding MUST map to an existing or new issue under established labels, milestones, parent/dependency and functional-outcome conventions.
- **FR-010**: Compatible corrections MUST preserve original failures, include discriminating regressions and identify fresh corrected observations without converting historical outcomes.
- **FR-011**: Source delivery MUST use canonical verification, an official PR and at most two total Codex reviews including automatic review, stopping before owner merge/tag/promotion.
- **FR-012**: Original #355/#356 causes and owner-deferred clustered execution MUST remain separate from controlled demonstrations and this audit's delivered outcomes.

### Key Entities

- **Candidate**: Immutable source and executable identities, platform and build provenance.
- **Observation**: Surface/workflow, expected/actual outcome, execution boundary and evidence references.
- **Finding**: Independently actionable functional defect, issue identity, source disposition and dependencies.
- **Evidence**: Private original or inspected sanitized derivative with explicit provenance and hash.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every inventoried CLI surface and product area has an explicit attributable disposition, with zero silent omissions.
- **SC-002**: Both cross-platform directions have actual observations or precise blocked outcomes; no simulated adapter result is claimed as cross-host execution.
- **SC-003**: Operator-focus acquisition and non-admin interaction requests remain zero.
- **SC-004**: Every actionable finding maps to a deduplicated issue, and every public artifact is inspected for private data.
- **SC-005**: Original evidence remains unchanged, corrected observations identify fresh candidates, and all remaining functional/publication outcomes are stated literally.

## Assumptions

- The owner approved this full scope and unattended execution; routine decisions are resolved and recorded without further questions.
- Existing source/main is 8b10419b916e04891609b3cef9727da16f486a23 and latest published packages are v2.0.1 at kickoff. v2.0.2 is the maintenance delivery target, not an existing release.
- Existing owned fixtures, reports and hidden launchers are reused; no new general QA framework or historical migration work is authorized.
- Optional context hooks are omitted because the plan's core context update supplies the same active feature context.
