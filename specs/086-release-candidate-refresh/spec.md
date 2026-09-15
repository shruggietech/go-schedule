# Feature Specification: Corrected v1.4.0 candidate refresh

**Feature Branch**: `codex/086-release-candidate-refresh`

**Created**: 2026-09-15

**Status**: Implemented

**Delivery**: Implemented on review branch `codex/086-release-candidate-refresh`; focused Windows fixtures, source review, PowerShell compliance, native desktop build, and canonical eight-gate results are recorded in [verification.md](verification.md). Refs #226 and #228; native candidate qualification, tag replacement, hosted staging, and promotion remain incomplete.

**Input**: User description: "Kick off S086", following the recommendation to investigate extreme interactive installer delays, automate mechanical Windows qualification work, prepare the corrected exact-commit candidate, and minimize remaining human steps.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Launch a prepared qualification session (Priority: P1)

The maintainer launches a prepared disposable Windows session without transcribing commands or manually assembling provenance. Executable inputs are verified and copied to guest-local storage before installation.

**Why this priority**: The previous workflow imposed repetitive manual setup and unobservable waits.

**Independent Test**: Generate a session package from explicit candidate metadata; inspect its launch configuration; prove missing or changed inputs refuse execution without installing on the development host.

**Acceptance Scenarios**:

1. **Given** the candidate, public v1.1.1 baseline, and helper tools, **When** preparation runs, **Then** it produces separate fresh and upgrade launch packages with immutable input hashes and dedicated output destinations.
2. **Given** a host-mapped input folder, **When** guest setup starts, **Then** all executable inputs are verified and staged to guest-local disk before use.
3. **Given** missing, altered, or overlapping inputs, **When** validation runs, **Then** it names the failure and performs no installation.

### User Story 2 - Explain installation waits (Priority: P1)

The maintainer receives phase timings, periodic progress, and persistent verbose installer logs instead of guessing whether setup is frozen.

**Why this priority**: Reported fifteen-minute preparation delays must be investigated before further attended qualification.

**Independent Test**: Successful, failing, and timed-out child-process fixtures produce distinct non-destructive timing records and retain logs.

**Acceptance Scenarios**:

1. **Given** a running installer, **When** it continues past a progress interval, **Then** elapsed time and the diagnostic log location are reported without launching another installer.
2. **Given** failure or timeout, **When** the operation ends, **Then** its result cannot be reported as passing and subsequent conflicting installer actions stop.
3. **Given** the reported latency, **When** investigation is documented, **Then** confirmed, rejected, and unverified hypotheses remain distinct.

### User Story 3 - Preserve genuine attended verification (Priority: P1)

The maintainer reviews a bounded checklist for the corrected desktop. Mechanical facts are collected automatically; visual defaults, normal-user interactions, and native behavior remain explicitly human-reviewed.

**Why this priority**: Reducing effort must not fabricate evidence or weaken existing release gates.

**Independent Test**: Map the session checklist to the shared release gate and issues #229 through #233; prove unattended output cannot satisfy missing attended observations.

**Acceptance Scenarios**:

1. **Given** successful mechanical setup, **When** the walkthrough begins, **Then** the exact remaining human steps are displayed with issue-level UI traceability.
2. **Given** missing, failed, skipped, stale, or unattended evidence, **When** finalization is attempted, **Then** qualification remains rejected.
3. **Given** an upgrade session, **When** the candidate is applied, **Then** its baseline is the independently identified public v1.1.1 artifact, not a same-authoring CI fixture.

### User Story 4 - Refresh without corrupting release history (Priority: P2)

The maintainer receives an exact-commit refresh runbook identifying the abandoned draft/tag, preserving old diagnostics, and separating implementation from authorized release operations.

**Why this priority**: Existing v1.4.0 candidate commit 57555ffa413df641cb21784847ad598c113bf199 excludes three UI remediation merges.

**Independent Test**: Verify local preparation cannot change any tag, draft asset, public release, issue closure, or milestone closure.

**Acceptance Scenarios**:

1. **Given** the stale draft, **When** S086 prepares the refresh, **Then** old and proposed commit identities are recorded without moving the tag or changing hosted assets.
2. **Given** separately authorized staging, **When** qualification occurs, **Then** evidence binds to the exact staged bytes and provenance, not a nominally equivalent rebuild.

### Edge Cases

- Paths contain spaces, XML metacharacters, or wildcard characters; input and output overlap.
- A destination already contains evidence, inputs change after package generation, or a mapped share becomes unavailable.
- PowerShell or WebView2 is missing and network bootstrap is unavailable.
- An installer hangs, another installation is active, or guest reset destroys unexported evidence.
- Elevated visual checks cannot prove normal-user access; fresh and upgrade scenarios inherit residual state.
- The draft becomes public or its tag identity changes before an authorized refresh.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Preparation MUST create uniquely identified qualification packages without installing or removing development-host software.
- **FR-002**: Candidate, baseline, and helper inputs MUST be bound to explicit provenance, byte lengths, and SHA-256 digests; absent, altered, overlapping, or previously occupied destinations MUST be rejected.
- **FR-003**: Guest execution MUST verify and copy executable inputs to local disk before invocation and MUST never execute an MSI directly from a mapped host share.
- **FR-004**: Fresh and upgrade packages MUST use separate disposable environments; upgrade MUST identify the public v1.1.1 baseline independently.
- **FR-005**: Console child processes MUST run hidden and noninteractively with redirected output; only intentional product or installer UI may be visible.
- **FR-006**: Installer operations MUST be serialized, retain verbose logs and phase timings, and report progress at intervals no greater than 30 seconds.
- **FR-007**: Failure and timeout results MUST remain nonpassing, preserve diagnostics, and stop conflicting subsequent operations.
- **FR-008**: Desktop prerequisites MUST be prepared before the focused UI walkthrough; unavailable prerequisites MUST be reported without an indefinite invisible download wait.
- **FR-009**: Automated facts MUST NOT claim visual, DPI, normal-user, or attended interaction evidence.
- **FR-010**: Remaining human steps MUST map to the shared gate and the theme, navigation, task workflow, feedback, Notifications, and administration findings in #229 through #233.
- **FR-011**: The existing exact-byte and native release gate MUST remain unchanged in rigor; required observations MUST NOT be dropped or manufactured.
- **FR-012**: Evidence MUST export to a dedicated host output directory without overwriting previous results.
- **FR-013**: Investigation MUST document supported and unsupported hypotheses and MUST NOT claim latency resolved before native timing evidence exists.
- **FR-014**: Refresh guidance MUST distinguish local preparation, review publication, tag replacement, staging, attended installation, and public promotion authorization.
- **FR-015**: Issues #226 and #228 and the v1.4.0 milestone MUST remain open until their individual acceptance criteria and publication gates pass.
- **FR-016**: Preparation and process control MUST have non-destructive regression coverage for escaping, tampering, overlap, hidden launches, failure, and timeout.

### Key Entities

- **Session manifest**: Scenario, candidate/baseline identity, input hashes, helper inputs, provenance, and output destination.
- **Phase result**: Start/end timestamps, elapsed duration, exit code when available, success/failure/timeout state, and log reference.
- **Attended checklist**: Human-only observations linked to release-gate requirements and UI findings.
- **Refresh record**: Old draft/tag, reviewed replacement commit, authorization boundaries, and staging/qualification/promotion evidence.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Launching a prepared scenario requires no command, commit, or workflow-identifier transcription.
- **SC-002**: Every installation wait emits progress within 30 seconds and retains persistent timing/log references.
- **SC-003**: Every invalid-input, tampering, overlap, failure, and timeout fixture refuses a pass without development-host installation.
- **SC-004**: Every required release observation and all five UI remediation areas map to automated facts or explicit human steps, with zero fabricated observations.
- **SC-005**: Fresh and upgrade packages identify independent clean baselines and preserve results outside ephemeral guest storage.
- **SC-006**: Preparation changes zero tags, hosted assets, public releases, issue closures, or milestone closures.

## Assumptions

- Windows Sandbox is enabled. The S086 kickoff authorizes implementation, not tag replacement, hosted release mutation, or public promotion.
- The existing release matrix remains authoritative. A brief visual walkthrough alone does not satisfy its complete attended gate.
- Native latency is reported but not yet reproduced. Local-demo builds and hosted silent tests cannot substitute for attended Windows 11 evidence.
- Scope is release tooling, diagnostics, preparation, and proportional documented fixes, not new UI redesign or v1.5.0 capabilities.

## Clarifications

### Session 2026-09-15

- Q: May mechanical automation replace attended observations? A: No. Missing human observations remain unavailable.
- Q: Does kickoff authorize draft/tag replacement? A: No. Stop at separate publication boundaries.
- Q: Should an installer deadline kill the Windows Installer service? A: No. Record timeout and prohibit subsequent operations until guest reset.
- Q: Which prerequisites avoid another hidden network wait? A: Hashed offline WebView2 standalone installer and portable PowerShell 7 ZIP, staged locally.
