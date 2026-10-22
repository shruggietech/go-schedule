# Feature Specification: SMTP notification channels

**Feature Branch**: `codex/111-smtp-notification-channels`

**Created**: 2026-10-21

**Status**: Implemented

**Delivery**: [PR #345](https://github.com/shruggietech/go-schedule/pull/345), open and unmerged; review branch codex/111-smtp-notification-channels; [local delivery](delivery.md) and [verification](verification.md), all eight canonical gates and finalized Windows/WSL evaluation passed 2026-10-22; authorized publication follows, maintainer retains merge authority.

**Input**: Owner-authorized S111 delivers #176 in full and completes notification coordinator #19 using the recommended end-to-end scope, including discovered remote-administration defect #341. Baseline: e5aa3240676b1bbc0d215bdb6a7b120ce1240d86.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Receive unattended task email (Priority: P1)

An operator configures an SMTP destination, assigns existing task/group notification conditions, and receives readable email even while the desktop is closed.

**Why this priority**: Email adds a broadly usable unattended destination to the existing notification system.

**Independent Test**: An owned mail receiver observes a test message and task outcome message through a disposable daemon; a failed mail destination does not affect a task or webhook destination.

**Acceptance Scenarios**:

1. **Given** an authenticated or unauthenticated mail service with a supported explicit encryption mode, **When** a channel is configured and tested, **Then** its recipients receive a readable message and delivery history identifies the outcome without revealing credentials.
2. **Given** task/group conditions and both email and webhook destinations, **When** a matching event occurs, **Then** each destination applies its independent policy and a mail failure does not change the task outcome.
3. **Given** multiple recipients and one rejected address, **When** delivery is attempted, **Then** acceptance, rejection and retry behavior is explicit per recipient and already confirmed recipients are not deliberately resent during a normal retry.

### User Story 2 - Administer mail through every supported interface (Priority: P1)

An administrator uses either the desktop or headless tools to create, inspect, edit, disable, enable, rotate credentials, test and delete SMTP channels, and inspect delivery results on the selected local or remote daemon.

**Why this priority**: A transport is incomplete if its supported administration surfaces lack parity.

**Independent Test**: Equivalent lifecycle operations work through desktop bindings and headless/local/remote interfaces with shared validation and permission behavior.

**Acceptance Scenarios**:

1. **Given** a channel, **When** its settings or credential change, **Then** subsequent deliveries use the new configuration and ordinary reads never reveal credentials.
2. **Given** a remote target, **When** an authorized caller manages mail, **Then** changes apply only to the selected daemon and existing access checks remain effective.
3. **Given** existing webhook channels and queued deliveries, **When** stored data advances, **Then** webhook definitions, policies and delivery records retain their meaning.

### User Story 3 - Move configuration and diagnose problems unattended (Priority: P2)

Maintainers transport SMTP definitions using portable configuration, provision secrets explicitly at the destination, and evaluate complete mail scenarios with predictable reports and useful plain-language diagnostics.

**Why this priority**: Headless operation and reproducible evaluation make the new capability maintainable.

**Independent Test**: A redacted bundle round-trip and isolated SMTP/TLS evaluation exercise success and intentional errors without a production account or installed-service changes.

**Acceptance Scenarios**:

1. **Given** an SMTP channel with credentials, **When** its configuration is exported, **Then** the portable definition carries no credential and identifies required destination provisioning.
2. **Given** invalid settings, certificate failure, missing authentication, rejected recipients or unavailable service, **When** an operation fails, **Then** a clear diagnostic identifies the problem and correction without including credentials or raw server text containing them.
3. **Given** the unattended evaluator, **When** mail cases run, **Then** owned fixtures provide reproducible results, useful communication reviews and existing canonical report formats.

### Edge Cases

- Blank/invalid server, port, sender or recipients; duplicate addresses; header/command injection; excessive settings; invalid timeout; unknown encryption mode; wrong-kind configuration.
- Required STARTTLS unavailable, invalid/expired/untrusted certificate, wrong hostname, invalid authentication, cancellation and a server that stalls.
- Mixed recipient outcomes, temporary versus permanent rejection, a lost connection after message submission, process restart and bounded retries.
- Unicode task/group names and content requiring HTML escaping; no task object for test and daemon-health messages.
- Disabled/deleted/edited channels with queued work; restart recovery; credential rotation/clearing; bundle dry-run, missing secret binding and drift.
- Empty/loading/error desktop states, keyboard navigation, human/JSON output separation and exact target selection.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Add optional SMTP destinations with server, port, explicit encryption mode, sender, one or more recipients, optional authentication, bounded timeout and optional simple HTML. Preserve existing channel identities and webhook behavior.
- **FR-002**: Support required STARTTLS, implicit TLS and explicitly selected unencrypted relay mode. TLS modes require certificate and hostname verification. Authentication is optional and explicitly configured; never infer downgrade after TLS failure.
- **FR-003**: Reject invalid settings, unsupported modes and header/command injection at authoring boundaries with field-specific actionable diagnostics. Bound recipient count, field sizes, message size and attempt time.
- **FR-004**: Store credentials using the existing protected write-only secret model. Reads, list output, errors, logs, delivery evidence and exports must not disclose them. Rotation and explicit clearing must work through supported administration interfaces.
- **FR-005**: Reuse task/group policy inheritance, independent destination conditions, daemon-health events and durable bounded delivery work. Mail delivery must not change task results, block scheduling or prevent webhook delivery.
- **FR-006**: Render restrained plain-text messages, with optional escaped simple HTML, from common event content. Unicode task/group names must remain readable; messages include useful task/outcome or health/test context and stable correlation.
- **FR-007**: Record recipient-specific outcomes, distinguish retryable temporary failures from permanent rejection, and retry only recipients not already confirmed as accepted. Retain bounded retry behavior and explain ambiguous post-submission transport failures; do not claim exactly-once delivery.
- **FR-008**: Provide equivalent lifecycle, credential rotation/clearing, test-delivery and history operations through desktop, CLI/local IPC and remote JSON access with existing authorization and target selection.
- **FR-009**: Portable configuration includes typed SMTP settings and secret references/provisioning requirements without exporting credentials. Preview/apply/drift must represent SMTP accurately and leave changes unapplied when required bindings are missing.
- **FR-010**: Stored-state migration is forward-only and non-destructive; old webhook channels and queued work remain operable after migration and daemon restart.
- **FR-011**: SMTP failures produce contextual, redacted messages useful to a layperson, including encryption/certificate, authentication, rejection, timeout and connectivity failures. Help text and documentation explain configuration, provider examples, retries and transport uncertainty.
- **FR-012**: Extend the existing unattended command inventory/scenarios and fixed report schema using owned local SMTP/TLS fixtures. No installed-service, user-session or production-mail changes belong in routine evaluation.
- **FR-013**: All new child processes are hidden/noninteractive on Windows; transport cancellation and daemon shutdown terminate owned work within bounded time. Preserve existing scheduler performance budgets.
- **FR-014**: The desktop form has labelled keyboard-accessible transport-specific settings, pending/error feedback and credential-presence indicators without exposing stored secrets.
- **FR-015**: Preserve portable-bundle schema compatibility through optional typed additions or an explicit supported version path; webhook-only definitions continue to import without invented mail settings.

### Key Entities

- **SMTP channel**: Reusable mail destination with transport settings, sender/recipients, protected optional authentication, enablement and shared notification policy.
- **Mail delivery**: Durable notification event with frozen destination settings and redacted recipient outcomes, bounded attempts and final state.
- **Portable SMTP definition**: Logical channel identity/settings and explicit destination-secret binding, excluding credentials.
- **Evaluation fixture**: Owned local mail service with deterministic success, rejection, TLS and timeout behavior, isolated from production resources.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All supported channel lifecycle operations are available through both desktop and headless administration with the same validation and permission semantics.
- **SC-002**: For every configured recipient, delivery history identifies accepted, retryable or permanently rejected outcome; confirmed recipients are excluded from subsequent normal retries.
- **SC-003**: No SMTP credential appears in ordinary reads, exports, task history, delivery diagnostics or logs.
- **SC-004**: Existing webhook definitions and queued events retain their meaning after stored-state migration and restart.
- **SC-005**: Intentional mail failures leave task outcomes and another configured notification destination unaffected.
- **SC-006**: Maintainers can perform documented isolated evaluation and obtain the existing predictable human/machine reports without provider credentials or user interaction.

## Assumptions

- Existing #158 notification storage/delivery, #18 remote access, #284 destination policies, #288 isolation and #289/#290 evaluation/report contracts are merged prerequisites.
- Provider OAuth enrollment, inbound mail and custom provider integrations are outside this slice. Existing SMTP username/password or relay credentials are sufficient.
- Internationalized task/message text is required; internationalized envelope addresses are accepted only where the selected transport explicitly supports them and must otherwise fail clearly.
- A live-provider exercise and release publication are separate authorized operations. This slice includes local implementation, push, official PR, CI and at most two Codex review rounds, then stops before merge.
- GitHub issue closure follows functional delivery, not test or CI activity.

## Clarifications

### Session 2026-10-21

- Q: Should the transport be delivered separately from its UI/headless administration? A: Deliver the full end-to-end #176 scope in S111, preserving parity.
- Q: Which encryption modes and certificate behavior are supported? A: Required STARTTLS, implicit TLS and explicit unencrypted relay mode; TLS verification is mandatory and there is no opportunistic downgrade.
- Q: How should partial recipients and retries behave? A: Keep durable per-recipient outcomes and omit confirmed accepted recipients from normal retries; state unavoidable uncertainty after interrupted submission.
- Q: Does implementation need a production email account? A: No, use owned disposable fixtures for routine verification and defer live-provider credentials.
- Q: Does the normal pre-publication halt apply? A: Publication is explicitly authorized in this kickoff; proceed automatically through push/PR and bounded reviews, retaining the final pre-merge owner review.
