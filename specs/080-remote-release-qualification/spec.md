# Feature Specification: Remote Access Release Qualification

**Feature Branch**: `codex/080-remote-release-qualification`

**Created**: 2026-09-10

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Service-bound configuration, package-shaped remote lifecycle qualification, three-platform hosted coverage, and the complete operator runbook completed for [#173](https://github.com/shruggietech/go-schedule/issues/173), with canonical eight-gate verification passed 2026-09-10 on review branch `codex/080-remote-release-qualification`; no release was published.

**Input**: Complete issue #173 as the v1.4 remote-access documentation and qualification gate without publishing a release.

## Clarifications

### Session 2026-09-10

- Q: What evidence is sufficient before a v1.4 tag exists? -> A: Cross-platform package-shaped CI and deterministic integration evidence qualify the reviewed source; public release artifacts remain a separate authorized ritual.
- Q: How should a registered headless service receive remote configuration? -> A: `service install --config` records one validated absolute configuration path in the service definition, while a flag-free install preserves current defaults.
- Q: Which deployment modes are recommended? -> A: Private-network HTTPS and SSH-tunneled HTTPS are recommended; reverse proxy and direct public HTTPS are supported only with explicit operator-owned prerequisites.

## User Scenarios & Testing

### User Story 1 - Operate a Remote Headless Daemon (Priority: P1)

As an operator, I can install a headless daemon with an explicit configuration, enable HTTPS, pair a client, perform an authorized operation, revoke it, disable HTTPS, and upgrade without accidental exposure or state loss.

**Why this priority**: This is the complete release journey promised by issue #173.

**Independent Test**: On each supported operating system, build package-shaped daemon and CLI binaries, run a configured daemon through enable, pairing, use, revocation, disable, restart, and replacement-binary cycles, and verify local IPC remains available throughout.

**Acceptance Scenarios**:

1. **Given** a default or explicitly disabled configuration, **When** the daemon starts, **Then** no remote TCP listener is opened and local IPC remains healthy.
2. **Given** an enabled TLS configuration, **When** an administrator creates a pairing and a client exchanges it, **Then** the client can verify the pinned daemon identity and perform only authorized operations.
3. **Given** a revoked credential, **When** it is used again, **Then** the remote request fails while local administration remains available.
4. **Given** a registered service installed with a configuration path, **When** the service manager launches it, **Then** it uses the validated absolute path across restart and binary upgrade.

### User Story 2 - Choose a Supported Deployment Safely (Priority: P1)

As an operator, I can distinguish recommended remote deployments from advanced examples and understand which TLS, DNS, routing, firewall, proxy, and key-lifecycle duties remain mine.

**Why this priority**: Documentation that blurs product guarantees and operator infrastructure would create unsafe public exposure.

**Independent Test**: Validate the documentation contract for explicit mode rankings, TLS requirements, service configuration, pairing, revocation, disablement, upgrade, backup, incident recovery, and public-listener prerequisites.

**Acceptance Scenarios**:

1. **Given** a private network or SSH access path, **When** an operator follows the guide, **Then** application TLS remains enabled and no insecure client override is offered.
2. **Given** a reverse proxy or public address, **When** an operator reads the guide, **Then** certificate issuance and renewal, DNS, firewalling, monitoring, backend trust, and exposure acknowledgement are explicit prerequisites.
3. **Given** a certificate, daemon identity, address, or credential incident, **When** recovery guidance is followed, **Then** trust is repaired deliberately without silent acceptance, automatic profile rewriting, or credential recovery claims.

### User Story 3 - Review Release Evidence (Priority: P1)

As a maintainer, I can audit one source-bound qualification record showing which functional, compatibility, permission, audit, artifact, documentation, and platform checks passed.

**Why this priority**: The milestone can close only from reproducible evidence, not association with prior slices.

**Independent Test**: Run the canonical verification suite and the named three-platform v1.4 qualification workflow, then map every issue #173 criterion to a test, document, or hosted result.

**Acceptance Scenarios**:

1. **Given** an incomplete source tree, **When** qualification runs, **Then** missing deployment guidance or lifecycle coverage fails with an actionable diagnostic.
2. **Given** a reviewed commit, **When** all local and hosted checks pass, **Then** the evidence names the exact commit and supported platforms without claiming that a public release was published.

### Edge Cases

- A relative service configuration path could resolve differently under a service manager and is resolved before registration.
- A configuration file can disappear after registration; the daemon then fails closed rather than starting with defaults.
- Port allocation can race in tests; the harness uses bounded startup diagnostics and never treats an unreachable endpoint as proof of disablement without confirming local health.
- Revocation, grant denial, version skew, identity change, network loss, and certificate change remain distinct evidence classes.
- Existing configuration, daemon identity, tasks, actors, credentials, and audit records survive a binary replacement.

## Requirements

### Functional Requirements

- **FR-001**: Service installation MUST accept an optional daemon configuration file, resolve it to an absolute path, validate it before registration, and persist it as the daemon's `--config` argument.
- **FR-002**: A flag-free daemon or service start MUST load an optional `config.json` from the platform data directory and MUST preserve safe built-in defaults when that file is absent.
- **FR-003**: A missing, unreadable, malformed, or invalid configuration supplied for service installation MUST fail before the service definition is changed.
- **FR-004**: Package-shaped cross-platform qualification MUST prove default-off behavior, explicit TLS enablement, pairing, identity verification, authorized use, permission denial, audit attribution, revocation, disablement, restart, and replacement-binary state preservation.
- **FR-005**: Qualification MUST cover Windows, macOS, and Linux through a named, fail-fast-disabled hosted matrix.
- **FR-006**: The operator guide MUST recommend private-network and SSH-tunneled HTTPS before reverse-proxy or direct-public modes.
- **FR-007**: Every network mode MUST retain application TLS and state its product-owned and operator-owned responsibilities.
- **FR-008**: Direct public HTTPS guidance MUST require explicit exposure acknowledgement, trusted certificate issuance and renewal, DNS, firewall restriction, monitoring, and a documented revocation response.
- **FR-009**: The guide MUST document enablement, service registration, pairing for desktop, CLI, and JSON clients, capability selection, use, rotation, revocation, disablement, upgrade, backup, and incident recovery.
- **FR-010**: Supported configurations MUST be visibly distinguished from illustrative commands and operator-owned infrastructure.
- **FR-011**: Documentation and tests MUST never expose bearer credentials, pairing phrases, certificate private keys, task inputs, or raw sensitive transport details.
- **FR-012**: S080 MUST retain local IPC behavior and MUST NOT publish a tag, release, or public artifact.

### Key Entities

- **Service Configuration Binding**: One validated absolute path persisted in a service definition as the daemon's startup argument.
- **Deployment Mode**: A supported connectivity shape with explicit application and operator responsibilities.
- **Qualification Journey**: A source-bound sequence covering initial default, enabled remote use, revoked access, disabled restart, and upgraded restart.
- **Qualification Evidence**: Local and hosted results mapped to issue #173 acceptance criteria.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of service-configuration acceptance tests prove the exact registered argument list or fail before registration.
- **SC-002**: The package-shaped lifecycle passes on Windows, macOS, and Linux without platform-specific behavior exclusions.
- **SC-003**: Every issue #173 acceptance criterion maps to at least one executable check or published documentation section.
- **SC-004**: Default and disabled starts expose zero remote listeners while retaining a healthy local connection in all platform jobs.
- **SC-005**: Secret-canary scans find zero protected values in committed evidence, logs, diagnostics, or documentation.
- **SC-006**: All eight canonical verification gates pass without exclusions added by S080.

## Assumptions

- Issues #165 through #172 are complete and their focused tests remain the detailed security and client-behavior evidence.
- GitHub-hosted Windows, macOS, and Linux runners are the supported cross-platform source qualification environments.
- A public v1.4.0 tag and immutable release artifacts require a separate explicit release request after merge.

## Dependencies

- Parent: [#18](https://github.com/shruggietech/go-schedule/issues/18).
- Completes: [#173](https://github.com/shruggietech/go-schedule/issues/173).
- Depends on: #145 and #165 through #172, all complete.

## Scope Boundaries

**In scope**: Service configuration persistence, package-shaped three-platform remote lifecycle qualification, deployment and recovery documentation, release evidence, and issue traceability.

**Out of scope**: Release tags, public artifacts, certificate automation, NAT traversal, service discovery, browser CORS, general user accounts, clustered execution, new remote routes, or remote MCP authority.
