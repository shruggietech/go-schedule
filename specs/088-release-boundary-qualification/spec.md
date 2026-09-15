# Feature Specification: Final v1.4.0 boundary and qualification

**Feature Branch**: `codex/088-release-boundary-qualification`

**Created**: 2026-09-15

**Status**: In Progress

**Delivery**: Source preparation implemented; review and post-merge candidate execution pending. Use Windows Sandbox without host changes or restart. Native checks unavailable in Sandbox may remain untested under the maintainer's explicit release waiver.

**Input**: User description: "Run S088 under spec-kit autopilot according to the recommended v1.4.0 qualification approach, automatically push and open a PR, handle up to two external review rounds, and report green CI for final maintainer review."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review a coherent corrected release boundary (Priority: P1)

The maintainer sees one cumulative v1.4.0 boundary that includes the reviewed desktop fixes and describes unpublished software honestly.

**Why this priority**: Candidate source, release notes, and the changelog must agree before final qualification.

**Independent Test**: Audit the cumulative changelog, version health example, release notes, draft status, and reviewed source identity with the existing release preflight.

**Acceptance Scenarios**:

1. **Given** reviewed desktop fixes and an unpublished draft, **When** the release boundary is reconciled, **Then** its source-owned summaries include those corrections without advertising a public release.
2. **Given** an unreviewed boundary edit, **When** final staging is requested, **Then** it waits for maintainer merge and exact-main verification rather than tagging the review branch.

### User Story 2 - Qualify exact repaired candidate bytes (Priority: P1)

The maintainer receives genuine fresh-install, public-v1.1.1 upgrade, task-create/delete, and corrected desktop results from disposable Windows 11 environments.

**Why this priority**: Passing source tests does not prove the installed native application.

**Independent Test**: Validate the candidate identity and all required observations and attachments through the existing native evidence gate.

**Acceptance Scenarios**:

1. **Given** verified reviewed source and explicit release-operation authority, **When** the draft is refreshed, **Then** all eight assets and the manifest agree with the exact staged revision, preserving obsolete provenance separately.
2. **Given** suitable clean Windows 11 environments, **When** qualification runs, **Then** actual normal-user, profile, window/display, lifecycle, and task observations satisfy their individual criteria.
3. **Given** missing environment capability, stale bytes, or a timeout, **When** detected, **Then** execution records the actual nonpassing result; unavailable Sandbox checks may remain explicitly untested under the maintainer waiver, but stale bytes and known failures are not treated as approved passes.

### User Story 3 - Complete an evidence-backed handoff (Priority: P2)

The maintainer reviews tested changes and individual review dispositions without confusing a merge-ready prerequisite PR with a qualified public release.

**Why this priority**: Open release work must not disappear through an issue-closing keyword or a misleading state transition.

**Independent Test**: Compare individual issue criteria, slice tasks, release evidence, latest-head CI, and external review records.

**Acceptance Scenarios**:

1. **Given** a complete verified local deliverable, **When** published for review, **Then** its structured PR body matches the submitted Markdown, CI runs, and every review finding receives a disposition within the two-round limit.
2. **Given** incomplete native criteria or absent public-promotion approval, **When** a checkpoint is reported, **Then** #231, #228, and #226 remain open as appropriate and the release stays a draft.

### Edge Cases

- Source boundary changes require review and merge before staging; candidate identity cannot be silently substituted afterward.
- Sandbox does not provide every required profile, display, or physical-input environment.
- Fresh and upgrade sessions cannot share installation residue.
- A previous release-operation approval does not authorize unrelated host feature changes or public promotion.
- A second review round finds defects; fixes remain authorized, but no third round is requested.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Reconcile cumulative v1.4.0 release notes, changelog, and version documentation with reviewed corrections, retaining honest unpublished status.
- **FR-002**: Use the existing Windows Sandbox approach. Inventory supported scenarios and explicitly identify unsupported checks before installation; do not provision another virtualization platform or restart the host.
- **FR-003**: Stage only exact reviewed main source after release-operation approval and successful exact-main checks; preserve old tag, draft metadata, and assets before replacement.
- **FR-004**: Verify all eight staged asset identities and the exact manifest/MSI before observations; never substitute locally rebuilt or obsolete bytes.
- **FR-005**: Execute the fresh-install and public-v1.1.1 upgrade checks supported by Sandbox. Retain the full matrix, but record unavailable checks as untested with the maintainer's release waiver instead of claiming they passed.
- **FR-006**: Complete the corrected desktop walkthrough for #228 and native task-create/delete verification for #231 against the exact final candidate.
- **FR-007**: Reuse existing qualification preparation and collection; missing capability, failed, unavailable, timed-out, and partial observations remain nonpassing. A waiver authorizes release despite specified missing tests, not fabrication of passing evidence or disregard of known failures.
- **FR-008**: Run all eight canonical local verification gates in the foreground before commit/publication; handle every external review finding and verify latest-head CI with no more than two requested review rounds.
- **FR-009**: Preserve issue-level traceability and close each issue only after its own acceptance criteria pass. Exclude public promotion until explicitly approved.
- **FR-010**: Do not enable host virtualization features, provision a new VM platform, alter security settings, or install the product on the active development host without specific authority for that operation.

### Key Entities

- **Release boundary**: Version, cumulative history, release notes, reviewed source, and publication state.
- **Capability inventory**: Available disposable environments, command/control paths, user/profile support, and display/input evidence capability.
- **Candidate identity**: Repository, tag, reviewed commit, run/attempt, asset lengths/digests, and installer product identity.
- **Observation set**: Scenario, environment, actual status, measurements, and validated attachments.
- **Review handoff**: Local verification, PR body integrity, review dispositions, current-head CI, and explicit incomplete dependencies.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All release-boundary summaries and candidate identities agree, with zero obsolete observations counted as passes.
- **SC-002**: Account for all 47 native observations as verified, failed, or explicitly untested. Declare full qualification only if all required evidence validates; distinguish maintainer-authorized release with testing gaps from full qualification.
- **SC-003**: No unnecessary installation retry occurs against stale bytes or before the source-review dependency is satisfied; waived unavailable environments do not block supported Sandbox checks.
- **SC-004**: Every review finding is dispositioned and latest-head checks are green, with zero third-round requests or unauthorized publication actions.

## Assumptions

- Work bundles #231, #228, and #226; later v1.5 features remain out of scope.
- Reviewed S087 main is `c48ee096251b33187deda61200eaa25edd0e36ab`; exact-main CI and CodeQL completed successfully.
- Existing draft assets identify S086 revision `951d864d9683ec3bdcb1e37535ecddd67738bf13` and do not contain the S087 selector repair.
- Push and official PR creation are explicitly authorized. Public promotion and material host-environment changes are not authorized by that publishing instruction.
- Source-boundary changes require a maintainer merge before final staging. Do not redefine the slice as complete merely to bypass that dependency.

## Clarifications

### Session 2026-09-15

- Q: Which source and evidence are eligible? A: Final exact reviewed main source after boundary reconciliation; the S086 draft and prior S087 walkthrough remain historical, not passing evidence for new bytes.
- Q: Should missing environments be replaced by hosted server tests or another evidence framework? A: No. Reuse existing tooling and record unavailable checks honestly; the later maintainer waiver permits release despite specified testing gaps, not substitution of fixtures for native evidence.
- Q: Is public promotion included? A: No. Push/PR authority is explicit; promotion and material development-host changes require separate specific authority.
- Q: Can the full source-edit, stage, qualify, and final-merge sequence finish before maintainer merge? A: No. Source edits require review and merge before reviewed-main staging. Record this dependency explicitly rather than silently recasting S088 as an implemented qualification slice.
- Q: What overrides the earlier environment blocker? A: The maintainer explicitly requires reuse of Windows Sandbox, prohibits restarting the host, and authorizes release of incomplete items without testing. This supersedes the requirement to block the entire slice on unavailable native environments. Record each waived test as untested, preserve available tests and CI, and do not call the release fully qualified.
