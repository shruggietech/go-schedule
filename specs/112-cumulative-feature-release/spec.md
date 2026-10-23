# Feature Specification: S112 cumulative feature release

**Feature Branch**: `codex/112-cumulative-feature-release`

**Created**: 2026-10-22

**Status**: Implemented

**Delivery**: Merged through [PR #359](https://github.com/shruggietech/go-schedule/pull/359) at squash commit 11c2a9a46d4610fe8a67903cea246050224bd606. #352 and #353 closed on merge. Merge-head [CI](https://github.com/shruggietech/go-schedule/actions/runs/39902543218), CodeQL and public website deployment passed. The v1.6.0 source candidate remains pending public tagging/promotion under #354; published discovery remains v1.5.1. S113 follows up #355-#358 without claiming the intermittent CLR/IPC causes resolved.

**Input**: Owner-authorized S112 bundles #352, #353 and cumulative release coordinator #354, including dependency proposals #351, #273, #336, #337, #295, #297, #299 and #301. Push and PR publication are authorized, including at most two Codex review rounds; stop before merge or release publication.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand and correct CI failures (Priority: P1)

Maintainers can distinguish completed manual work from queued observations and receive useful failure evidence from disposable service probes. Credential fixture preparation has its own lifetime, leaving the full runtime budget for authentication assertions.

**Why this priority**: The unchanged merge source failed two CI jobs and currently lacks enough service evidence to establish the timeout cause.

**Independent Test**: Regress queued-before-terminal observations and setup/runtime deadline ownership, then exercise the unchanged authentication and disposable service contracts.

**Acceptance Scenarios**:

1. **Given** multiple queued and completed observations for a fixture task, **When** the probe seeks a manual outcome, **Then** it selects the matching terminal run and reports missing or unfinished work accurately.
2. **Given** a probe failure, **When** cleanup executes, **Then** bounded diagnostic observations and cleanup results survive without replacing the primary failure or exposing secrets.
3. **Given** successful but slow fixture preparation, **When** authentication checks start, **Then** they receive an independent runtime budget and retain all peer, parent, image and request-order checks.

### User Story 2 - Maintain a coherent dependency baseline (Priority: P2)

Maintainers can build the CLI, daemon, desktop and public website with reviewed compatible dependency updates and current publishing actions.

**Why this priority**: Shared module and lockfile dependencies must converge before release candidates identify a reproducible source.

**Independent Test**: Review each upstream change, dependency graph and policy fixture; run the production builds and canonical gates.

**Acceptance Scenarios**:

1. **Given** overlapping dependency proposals, **When** they are integrated, **Then** current dependencies and lockfiles preserve SMTP additions and desktop alignment, with every proposal explicitly accounted for.
2. **Given** Pages action upgrades, **When** publishing runs, **Then** domain, environment, permissions, artifact and release-comparison controls remain effective.

### User Story 3 - Prepare the accumulated release (Priority: P3)

Operators can understand the accumulated features, upgrade requirements and release contents from concise, mutually consistent release material, while public discovery continues to identify the actual published release.

**Why this priority**: v1.5.1 is the latest published package despite substantial merged functionality.

**Independent Test**: Validate the pending v1.6.0 comparison, release metadata, candidate build inputs, upgrade guidance and existing release integrity checks without creating a tag or public release.

**Acceptance Scenarios**:

1. **Given** pending source capabilities, **When** a user reads public discovery, **Then** source-only SMTP and CLI capabilities are distinguished from the published v1.5.1 package.
2. **Given** approved release preparation, **When** the owner later authorizes tagging, **Then** documented commands bind candidate assets, checksums and comparison to the exact reviewed source and existing draft promotion route.

### Edge Cases

- An unfinished queued run precedes a completed matching run; a skipped run is terminal but cannot satisfy successful execution assertions.
- Primary probe failure occurs before service startup, after fixture creation, or during cleanup; unavailable diagnostics must not erase the original exception.
- Fixture build cancellation, build failure and helper runtime rejection have distinct errors and owned cleanup.
- A refreshed dependency PR overlaps another update; an upstream major action changes runner requirements or output contracts.
- Candidate release is not published: existing public release links and comparison selection remain v1.5.1.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Correct the demonstrated #352 probe selection defect without changing ordinary scheduler overlap semantics, retaining manual, recurring, nonzero exit, start failure and LocalSystem/group contracts.
- **FR-002**: Persist bounded fixture-owned failure evidence and cleanup results before removal; retain the primary exception and exclude credential/environment dumps.
- **FR-003**: Separate #353 fixture construction and runtime lifetimes with explicit setup diagnostics, owned cancellation and process joining; retain zero request bytes before authentication and bounded channel reads.
- **FR-004**: Reconcile all eight dependency proposals against current source, preserving direct dependencies and aligning desktop/tooling versions; explicitly document incompatible or superseded proposals.
- **FR-005**: Update Pages actions only after upstream migration review and matching policy, fixture and dated decision changes, preserving deployment controls.
- **FR-006**: Prepare v1.6.0 release material, storage-v23 upgrade guidance and a release-bound comparison/opportunity snapshot; preserve published historical records.
- **FR-007**: Bind candidate preparation to the reviewed source, existing asset/checksum rules and draft/waived promotion flow; never create tags, publish releases or imply absent attended results passed.
- **FR-008**: Reconcile S111 merge metadata and S112 issue-level traceability, keeping #354 open until its publication outcome is delivered.
- **FR-009**: Preserve quiet isolated QA and hidden Windows process launches; installed services, account sessions and public DNS remain untouched.
- **FR-010**: Run canonical CI parity, relevant regressions and isolated unattended QA; file separately discovered defects, resolve every PR review thread and request no third Codex round.

### Key Entities

- **Probe observation**: fixture task, trigger, terminal outcome, timestamp and bounded diagnostic context.
- **Dependency proposal**: upstream PR, reviewed versions, affected module/lockfiles and final disposition.
- **Release candidate**: proposed version, immutable source identity, package assets/checksums, comparison and publication state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Both failure paths produce distinguishable outcomes and retained diagnostics without weakening any existing execution/authentication assertion.
- **SC-002**: Every one of the eight dependency proposals has an explicit source integration or supersession disposition.
- **SC-003**: Release preparation identifies one v1.6.0 source baseline, while all current-public discovery still identifies the actual latest published package.
- **SC-004**: Every included issue has an accurate functional status; no source readiness is reported as public publication.

## Assumptions

- The owner's acceptance of the proposed S112 approach selects v1.6.0 for preparation; tagging and publication remain outside this PR authorization.
- #354 spans source preparation and later publication, so this PR references it without automatic closure.
- Real LocalSystem and installer mutation runs only on disposable hosted CI. Isolated Windows/WSL QA may run locally with bounded resources.
- Conditional distributed-scheduling issues #20 and #186-#189 remain governed by their existing real-use prerequisite.

## Clarifications

### Session 2026-10-22

- Q: Does PR authorization include release tagging? A: No. Prepare source and reproducible candidate instructions, then stop for PR review/merge; exact-source release generation follows separately authorized tagging.
- Q: Can dependency convergence mask CI timing failures? A: Diagnose and regress #352/#353 on the unchanged source first, then apply reviewed dependency changes.
- Q: What if an upstream major action cannot preserve controls? A: Record the concrete incompatibility and retain the existing action until a compatible migration is available; never silently drop a proposal or weaken controls.
