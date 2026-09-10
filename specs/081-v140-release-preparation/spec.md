# Feature Specification: Cumulative v1.4.0 Release Preparation

**Feature Branch**: `codex/081-v140-release-preparation`

**Created**: 2026-09-10

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Local focused verification and all eight canonical gates passed on 2026-09-10 for [#226](https://github.com/shruggietech/go-schedule/issues/226); no tag or release is authorized by this slice.

**Input**: Prepare the cumulative v1.4.0 release boundary after the completed v1.2.0, v1.3.0, and v1.4.0 delivery milestones, then publish the reviewed preparation through a pull request without creating a tag or release.

## Clarifications

### Session 2026-09-10

- Q: What is the honest next public version after v1.1.1? -> A: v1.4.0, because current main contains the completed v1.2 through v1.4 increments and no earlier release tag can represent a subset of that history now.
- Q: Should v1.2.0 and v1.3.0 be tagged retrospectively? -> A: No. Preserve their closed milestones as delivery history and state that v1.4.0 is cumulative from v1.1.1.
- Q: What may S081 publish? -> A: One review branch and pull request only. Tagging, draft staging, attended installation, promotion, and public release require later explicit authorization.
- Q: What Windows upgrade must the release ritual prove? -> A: Both a clean v1.4.0 install and an in-place upgrade from the latest public v1.1.1 MSI using the exact staged v1.4.0 candidate.

## User Scenarios & Testing

### User Story 1 - Understand the cumulative release (Priority: P1)

A person evaluating v1.4.0 can quickly see the major desktop, automation, notification, local-agent, and remote-access outcomes delivered since v1.1.1, then follow the changelog for the complete record.

**Why this priority**: A cumulative release must explain a large delivery boundary without pretending intermediate versions were separately published.

**Independent Test**: Review the tag-specific release note, versioned changelog boundary, and installation entry points together and confirm that they describe one cumulative v1.4.0 release after v1.1.1.

**Acceptance Scenarios**:

1. **Given** the v1.4.0 release note, **When** a reader scans it, **Then** four concise highlights identify the principal outcomes and one final link opens the tagged changelog.
2. **Given** the completed v1.2.0 and v1.3.0 milestones, **When** a reader inspects release history, **Then** those milestones remain visible as delivery history without nonexistent publication claims or retroactive tags.

### User Story 2 - Stage one trustworthy artifact set (Priority: P1)

A maintainer can create one immutable v1.4.0 tag from the reviewed merge commit and have automation stage a complete draft artifact set only after exact-commit CI and release metadata agree.

**Why this priority**: Every later qualification result depends on source, metadata, and artifact identity being the same immutable boundary.

**Independent Test**: Evaluate the release workflow against v1.4.0 metadata and fixture mutations for the README version, changelog section, release note, source commit, package names, draft state, and expected asset inventory.

**Acceptance Scenarios**:

1. **Given** a tag whose commit lacks matching release metadata or successful main CI, **When** staging begins, **Then** automation fails before creating or mutating release artifacts.
2. **Given** the exact reviewed v1.4.0 commit and complete metadata, **When** staging completes, **Then** the draft contains the expected headless archives, Wails desktop archives, Windows MSI, and candidate manifest without becoming public.

### User Story 3 - Qualify installation and upgrade before promotion (Priority: P1)

A maintainer can prove that the exact Windows candidate installs cleanly, upgrades the latest public v1.1.1 release without losing scheduler state, and keeps every newly introduced network or agent surface disabled by default.

**Why this priority**: v1.4.0 changes the Windows desktop implementation and adds several opt-in capabilities, so a fresh-only or source-only check is insufficient.

**Independent Test**: On clean Windows 11 snapshots, run the attended release matrix once as a fresh install and once from the public v1.1.1 MSI, binding both paths and all required observations to the exact staged v1.4.0 candidate.

**Acceptance Scenarios**:

1. **Given** a clean supported Windows host, **When** the exact candidate is installed, **Then** the service, CLI, Wails desktop, shortcuts, removal choices, and local offline defaults satisfy the established release gate.
2. **Given** an installed public v1.1.1 baseline with representative tasks and preferences, **When** the exact candidate upgrades it, **Then** tasks, run history, supported appearance intent, daemon identity, service operation, and local access remain usable while remote HTTPS, notifications, and localhost MCP remain disabled until explicitly configured.
3. **Given** incomplete, stale, substituted, or failed evidence, **When** promotion is requested, **Then** the release remains draft.

### Edge Cases

- A local tag, remote tag, draft, or public v1.4.0 release already exists and must be reconciled before tag creation.
- A v1.2.0 or v1.3.0 tag appears after preparation and would need separate provenance review rather than automatic acceptance.
- The reviewed merge commit differs from the successful pull-request head or lacks successful push CI on main.
- Release notes or documentation claim v1.4.0 is public before the latest-release pointer identifies it.
- An upgrade preserves persisted state but silently enables a listener, delivery channel, or credential, which fails qualification.
- A candidate rebuild, renamed MSI, changed byte, extra asset, or missing asset invalidates the qualification boundary.

## Requirements

### Functional Requirements

- **FR-001**: The repository MUST contain an empty Unreleased section followed by one dated v1.4.0 section containing the complete accumulated history after v1.1.1 without content loss.
- **FR-002**: The changelog MUST identify v1.4.0 as a direct cumulative comparison from v1.1.1 and MUST NOT represent v1.2.0 or v1.3.0 as published tags.
- **FR-003**: The v1.4.0 release note MUST contain exactly four concise one-line highlights and MUST end with one link to the tagged v1.4.0 changelog section.
- **FR-004**: The README health example MUST identify version 1.4.0 before tag creation, while installation documentation MUST distinguish the prepared candidate from the latest public release until promotion finishes.
- **FR-005**: Candidate-era wording across the README and platform installation guides MUST be synchronized with the cumulative v1.4.0 boundary and MUST NOT continue to describe completed Wails work as merely a v1.2 candidate.
- **FR-006**: Release preflight MUST reject a tag when the README version, dated changelog section, tag-specific release note, or successful main CI for the exact tagged commit is missing or inconsistent.
- **FR-007**: Draft staging MUST produce the current expected headless archives, Wails desktop archives, Windows MSI, and Windows candidate manifest from one immutable commit, and MUST remain draft.
- **FR-008**: The post-merge qualification contract MUST require fresh installation and in-place upgrade from the public v1.1.1 Windows MSI to the exact staged v1.4.0 MSI.
- **FR-009**: Upgrade qualification MUST cover persisted tasks, run history, supported appearance intent, daemon identity, service operation, local access, and default-off notification, localhost MCP, and remote HTTPS behavior.
- **FR-010**: Promotion MUST validate the exact candidate, complete attended evidence, immutable tag, expected asset inventory, and final checksums without rebuilding artifacts.
- **FR-011**: The release issue, project item, and v1.4.0 milestone MUST remain open through preparation and close only after public promotion and the final identity audit satisfy their individual criteria.
- **FR-012**: S081 MUST NOT create or push a release tag, stage or publish a GitHub release, install a candidate, manufacture attended evidence, or close issue #226.
- **FR-013**: All canonical verification gates and focused release-automation regression checks MUST pass before the preparation branch is published.
- **FR-014**: Native Windows window evidence MUST be capturable from the installed Wails desktop without a retired Fyne-only metrics file, while historical version-one evidence remains validatable.

### Key Entities

- **Release Boundary**: The reviewed commit, semantic version, changelog section, release note, and documentation state that define cumulative v1.4.0 source.
- **Staged Artifact Set**: The draft-only packages and candidate manifest built from one immutable tag and commit.
- **Upgrade Baseline**: The public v1.1.1 Windows installation and representative persisted state used before applying the v1.4.0 candidate.
- **Qualification Evidence**: Machine-validated and attended records bound to the exact staged candidate, including fresh and upgrade paths.
- **Publication Record**: Issue #226, project lifecycle state, milestone, release metadata, and final audit result.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of accumulated changelog content after v1.1.1 appears under the dated v1.4.0 boundary, with zero synthetic v1.2.0 or v1.3.0 publication sections.
- **SC-002**: Release notes contain exactly four one-line highlights and one final tagged changelog link.
- **SC-003**: Release preflight mutation tests reject every missing or mismatched README, changelog, release-note, exact-commit CI, and release-state input.
- **SC-004**: All nine required draft payloads before checksums are non-empty and identify the v1.4.0 tag or its exact candidate provenance as applicable.
- **SC-005**: Fresh-install and v1.1.1-upgrade evidence both bind to the same staged MSI, and all 47 established attended observations pass.
- **SC-006**: Post-upgrade checks find zero implicitly enabled notification channels, localhost MCP listeners, or remote HTTPS listeners.
- **SC-007**: All eight canonical local gates and every required hosted pull-request check pass with no S081 exclusions.
- **SC-008**: Before separate release authorization, zero v1.4.0 tags or GitHub releases are created and issue #226 remains open.
- **SC-009**: The attended collector contains zero required Fyne-only inputs and the validator accepts current generic desktop evidence plus the retained historical fixture.

## Assumptions

- Public v1.1.1 is the latest immutable upgrade baseline and the latest-release pointer currently resolves to it.
- The closed v1.2.0, v1.3.0, and v1.4.0 implementation milestones are accurate delivery records even though v1.2.0 and v1.3.0 were not public releases.
- The existing Release and Promote Release workflows remain the authoritative draft-staging and no-rebuild promotion pipeline.
- The established 47 observation identities remain applicable to the Wails candidate, while the native-window attachment evolves compatibly from Fyne-specific version one to generic desktop version two evidence.

## Dependencies

- Parent: [#146](https://github.com/shruggietech/go-schedule/issues/146).
- Completes repository preparation for [#226](https://github.com/shruggietech/go-schedule/issues/226), which remains open for post-merge staging, qualification, promotion, and audit.
- Depends on completed issue [#173](https://github.com/shruggietech/go-schedule/issues/173) and merged pull request [#225](https://github.com/shruggietech/go-schedule/pull/225).

## Scope Boundaries

**In scope**: Cumulative release identity, changelog and release copy, candidate-aware installation documentation, release-preflight regression coverage, a precise post-merge publication contract, and complete local verification.

**Out of scope**: Tag creation, hosted draft staging, package installation, attended observation collection, release promotion, public publication, issue closure, milestone closure, unrelated backlog implementation, or new product behavior.
