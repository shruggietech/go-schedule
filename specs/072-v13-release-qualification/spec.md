# Feature Specification: v1.3 Notifications and Local Agent Access Qualification

**Feature Branch**: `codex/072-v13-release-qualification`

**Created**: 2026-09-09

**Status**: Implemented

**Delivery**: Package-shaped Windows qualification, focused notification and MCP race suites, documentation boundary audit, workflow-contract fixtures, and canonical eight-gate verification passed on 2026-09-09. Hosted Windows, macOS, and Linux qualification runs on the S072 pull request before merge; no release tag or public artifact is part of this slice.

**Input**: GitHub issue [#190](https://github.com/shruggietech/go-schedule/issues/190), qualify webhook notifications and optional observe-only local MCP as the v1.3.0 release boundary without publishing a release.

## User Scenarios & Testing

### User Story 1 - Preserve safe installation defaults (Priority: P1)

A person can install or upgrade go-schedule on Windows, macOS, or Linux and retain the local offline scheduler without silently enabling webhook delivery or an MCP network listener.

**Why this priority**: The release must not turn optional outbound or listening behavior into an installation side effect.

**Independent Test**: Start a package-shaped daemon from a fresh data directory and again from its retained state, then inspect notification channels, delivery history, and local MCP status before any explicit setup.

**Acceptance Scenarios**:

1. **Given** a fresh supported-platform installation, **When** the daemon first becomes healthy, **Then** it has no notification channels, no notification deliveries, and no active MCP HTTP listener.
2. **Given** an existing installation that has never enabled either optional surface, **When** it is upgraded or restarted, **Then** the same empty and disabled state is preserved.
3. **Given** either clean or retained state, **When** no network feature is configured, **Then** ordinary local scheduling remains available through protected local IPC without an account or Internet connection.

---

### User Story 2 - Qualify webhook delivery without affecting task truth (Priority: P1)

A maintainer can demonstrate that webhook configuration is explicit, protected values stay write-only, delivery failures remain isolated, and task outcomes never change because a receiver succeeds or fails.

**Why this priority**: Notifications are valuable only when they cannot corrupt scheduling results or disclose credentials.

**Independent Test**: Exercise test and run-outcome deliveries against successful, failing, and unavailable receivers, restart during recoverable work, and inspect the resulting bounded, redacted evidence and original run outcome.

**Acceptance Scenarios**:

1. **Given** an explicitly configured webhook, **When** a matching run completes, **Then** one bounded payload reaches the receiver with stable delivery identity and without task output, environment values, or stored authorization appearing in read responses.
2. **Given** a receiver failure or daemon interruption, **When** delivery is retried or terminally failed, **Then** scheduling continues and the source run retains its original outcome.
3. **Given** disabled notification configuration, **When** runs complete, **Then** no new outbound delivery is created.

---

### User Story 3 - Qualify local Observe access and truthful release claims (Priority: P1)

A maintainer can connect supported MCP clients through process-launched stdio or explicitly enabled authenticated localhost HTTP and prove that both expose the same bounded Observe-only resources with no mutation authority.

**Why this priority**: The local agent surface crosses a security boundary and must remain narrower than the future remote and mutation roadmap.

**Independent Test**: Use official clients against package-shaped stdio and localhost HTTP, cover both supported protocol revisions, every resource and continuation template, hostile scheduler content, authorization rejection, disconnect, rotation, revocation, and daemon restart.

**Acceptance Scenarios**:

1. **Given** a supported MCP host, **When** it launches the packaged stdio command, **Then** it discovers exactly the documented Observe resources and templates, discovers no tools, opens no listener, and exits with the host.
2. **Given** localhost HTTP is explicitly enabled, **When** a client presents the current credential and valid request metadata, **Then** it receives the same Observe contract; invalid credentials, hosts, or origins receive no scheduler data.
3. **Given** credentials are rotated, access is revoked, or the daemon restarts, **When** an old client reconnects, **Then** its previous credential is unusable and the listener remains off after restart.
4. **Given** release-facing documentation, **When** a reader evaluates available authority, **Then** Observe, stdio, localhost HTTP, and future remote or mutation capabilities are clearly distinguished.

### Edge Cases

- The daemon process is stopped abruptly and restarted with retained storage.
- A webhook receiver redirects, times out, returns a bounded error, or becomes unavailable between attempts.
- Scheduler-controlled strings contain Markdown, XML-like instructions, ANSI escapes, credentials, or protocol-shaped content.
- An MCP host disconnects during initialization or resource pagination.
- A stale localhost credential is reused after rotation, disablement, or restart.
- A supported-platform prerequisite is unavailable; the evidence must report the gate as unavailable rather than passing or silently skipping it.

## Requirements

### Functional Requirements

- **FR-001**: Qualification MUST exercise fresh and retained daemon state on Windows, macOS, and Linux using package-shaped executables.
- **FR-002**: Fresh and retained state MUST contain zero notification channels and deliveries until a user explicitly configures them.
- **FR-003**: Fresh and restarted state MUST report localhost MCP HTTP disabled with no endpoint, credential fingerprint, client identity, or access evidence.
- **FR-004**: Qualification MUST demonstrate that ordinary local health and scheduling access remain available without notifications, MCP HTTP, an account, or Internet connectivity.
- **FR-005**: Webhook evidence MUST cover success, bounded retry or terminal failure, daemon recovery, disabled-channel behavior, redaction, and preservation of the source run outcome.
- **FR-006**: MCP evidence MUST cover both supported protocol revisions, all five resources, all four continuation templates, zero tools, bounded errors, redaction, and hostile-content isolation.
- **FR-007**: Package-shaped stdio MUST initialize through an official client on every supported platform and terminate cleanly when its host disconnects.
- **FR-008**: Localhost HTTP evidence MUST cover valid authorization plus invalid credential, Host, Origin, rotation, revocation, and restart cases without exposing protected values.
- **FR-009**: The supported-platform qualification MUST be a named pull-request check so missing, skipped, or failed platform evidence is visible before merge.
- **FR-010**: Documentation MUST distinguish shipped webhook, stdio Observe, and localhost Observe behavior from future SMTP, native desktop delivery, remote MCP, Operate, and Manage work.
- **FR-011**: The slice MUST record exact local and hosted evidence, including any unavailable prerequisite, without claiming that a release tag or public artifact was produced.
- **FR-012**: S072 MUST NOT publish a release, create a tag, enable an optional network surface by default, add mutation tools, or expand remote access.
- **FR-013**: GitHub issue [#190](https://github.com/shruggietech/go-schedule/issues/190) MUST remain traceable through the specification, task list, changelog, pull request, and verification record.

### Key Entities

- **Qualification Matrix**: The required platform and scenario set with explicit pass, fail, or unavailable evidence.
- **Package-Shaped Daemon State**: An isolated executable, configuration, data directory, database, IPC endpoint, and logs used to prove first-start and restart defaults.
- **Notification Evidence**: Redacted channel, delivery, retry, terminal state, and source-run facts that prove isolation from scheduling truth.
- **MCP Evidence**: Protocol, transport, resource, template, authority, authorization, lifecycle, and hostile-content results produced by official clients.
- **Release Boundary**: The reviewed commit and hosted checks that qualify v1.3 behavior without creating a release tag or public artifact.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All three supported operating-system jobs start the package-shaped daemon twice and observe zero notification channels, zero deliveries, and zero active MCP HTTP listeners before explicit configuration.
- **SC-002**: One hundred percent of exercised webhook success and failure cases preserve the recorded task outcome and expose zero protected authorization or task-output values through read surfaces.
- **SC-003**: Official clients discover exactly five resources, four resource templates, and zero tools through each qualified MCP transport and supported protocol revision.
- **SC-004**: One hundred percent of rejected localhost credentials, hosts, and origins return no scheduler resource data, and old credentials fail after rotation, revocation, and restart.
- **SC-005**: Every supported-platform qualification result is visible as pass, fail, or explicitly unavailable; no required scenario is represented by a silent skip.
- **SC-006**: The complete canonical eight-gate verification passes before publication, and hosted pull-request checks pass before maintainer review.

## Assumptions

- S067 through S071 are the implemented notification and MCP baseline; S072 adds integrated release evidence and fixes only defects discovered by that qualification.
- Pull-request artifacts and package-shaped test executables are sufficient for this qualification boundary. Tagging and public release promotion require a separate explicit ritual.
- Existing platform CI runners provide the supported execution environments, while the local workstation provides Windows evidence before publication.
- The v1.3.0 milestone contains only issue #190 after its seven implementation dependencies have closed.
