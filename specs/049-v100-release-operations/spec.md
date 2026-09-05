# Feature Specification: v1.0.0 Release Operations

**Feature Branch**: `codex/049-v100-release-operations`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/049-v100-release-operations`; deterministic packet generation, release-runbook integration, focused/race/policy tests, and canonical eight-gate verification completed; hosted CI and review remain pull-request evidence

**Input**: Execute the evidence-backed portion of issue [#122](https://github.com/shruggietech/go-schedule/issues/122) by making exact-candidate qualification, per-issue disposition, promotion readiness, and final audit deterministic and fail-closed. The implementation pull request must not create a tag, publish a release, or treat prior local-demo evidence as formal evidence.

## Problem Statement

S048 prepared the reviewed v1.0.0 source boundary and specified the post-merge release ritual. S047 already enforces 47 exact-candidate Windows observations. The remaining high-risk manual step is translating one validated evidence archive into ten complete, consistent GitHub issue dispositions before release promotion. A missed observation, mismatched candidate identity, incomplete attachment inventory, or copied local-demo result could otherwise be hidden by an apparently complete narrative comment.

S049 supplies a deterministic disposition packet and an executable operations runbook. Actual tag creation and release publication remain later remote operations requiring separate, explicit maintainer authorization after this reviewed slice merges.

## Clarifications

### Session 2026-09-03

- Q: Which reviewed commit receives `v1.0.0` after the required S049 pull request advances `main` beyond S048? -> A: The S049 merge commit becomes the final reviewed release boundary; S049 may change only release tooling and documentation, not packaged runtime behavior.
- Q: Does a generated disposition record itself authorize or perform a GitHub issue closure? -> A: No. It is immutable review input, and each issue remains subject to individual maintainer judgment and an explicit remote operation.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate Evidence-Backed Issue Dispositions (Priority: P1)

As a release steward, I can transform one completely valid formal Windows candidate evidence archive into a separate review record for every remaining v1.0.0 readiness issue, so issue closure is based on the observations relevant to that issue rather than association with a release slice.

**Why this priority**: The release contract requires ten independent issue dispositions before promotion. Manually copying 47 observations and their candidate identity into overlapping issue comments is the most error-prone remaining part of that contract.

**Independent Test**: Provide a fully valid exact-candidate archive, its matching candidate manifest, and the matching MSI. The operation produces the complete expected issue-file set, and every file contains only its canonical observations plus the shared immutable candidate and validation evidence.

**Acceptance Scenarios**:

1. **Given** a valid 47-observation formal archive, matching manifest, and matching MSI, **when** the steward generates the disposition packet, **then** exactly ten issue records are produced for #96, #98, #101, #104, #105, #106, #109, #111, #112, and #113.
2. **Given** a generated child-issue record, **when** a reviewer reads it, **then** the record names the candidate tag, commit, workflow run, MSI hash, validator result, mapped observation results, environments, and supporting attachment paths needed to judge that issue independently.
3. **Given** the same validated inputs twice, **when** packets are generated in separate empty destinations, **then** corresponding files are byte-for-byte identical.

---

### User Story 2 - Fail Closed Before Disposition (Priority: P1)

As a release reviewer, I need report generation to reject incomplete, synthetic, mismatched, or corrupted evidence before it writes any issue record, so a polished report cannot launder invalid evidence into a closure claim.

**Why this priority**: A disposition packet is trustworthy only if it consumes the same production validation path that promotion uses and leaves no partial success output after failure.

**Independent Test**: Mutate each candidate identity, required observation, status, attachment, archive member, and evidence class in turn. Each invalid input is rejected with an actionable diagnostic and the destination remains absent or empty.

**Acceptance Scenarios**:

1. **Given** evidence from a local demo or automated fixture, **when** packet generation is attempted, **then** it is rejected and no issue file is created.
2. **Given** a formal archive whose manifest, MSI, repository, tag, commit, or attachment inventory does not match, **when** generation is attempted, **then** all independently discoverable validation failures are reported and no partial packet remains.
3. **Given** an existing output path, **when** generation is attempted, **then** the operation refuses to overwrite unexplained evidence or prior disposition records.

---

### User Story 3 - Execute and Audit the Release Ritual (Priority: P2)

As a maintainer, I can follow one chronological runbook from the reviewed merge boundary through tag staging, formal qualification, issue reconciliation, promotion, and final audit without confusing pull-request authority with tag or release authority.

**Why this priority**: The release is intentionally staged before publication. The runbook must keep the candidate immutable and make every stop condition clear while the new disposition packet is used.

**Independent Test**: A reviewer can walk the runbook against an inert example and identify the precise prerequisite, expected artifact count, issue state, and stop condition at every phase without making any remote mutation.

**Acceptance Scenarios**:

1. **Given** an unmerged S049 pull request, **when** the runbook is followed, **then** it explicitly forbids tag creation, release staging, issue closure, promotion, and milestone closure.
2. **Given** a merged S049 boundary and separate tag authorization, **when** the tag is staged, **then** the runbook binds the draft release, candidate manifest, MSI, evidence archive, and disposition packet to that exact commit.
3. **Given** ten individually satisfied readiness issues and the exact nine pre-checksum payloads, **when** promotion is authorized, **then** the runbook requires checksum generation and a final ten-payload public-release audit before #122 or the milestone can close.

### Edge Cases

- The evidence contains all canonical observation identifiers but one has a non-passing status, invalid metric, wrong display class, or unsupported attachment bytes.
- The evidence and candidate manifest agree with one another but the MSI bytes or basename differ.
- An observation maps to several issues; each issue record must include the observation without modifying or duplicating the source evidence.
- An attachment supports several observations or an observation has no required attachment; reports must preserve the exact declared relationship.
- The destination exists, is non-empty, is a file, or contains a symbolic link.
- A valid archive is generated from a different current working directory or with platform-native path separators.
- GitHub is unavailable after local packet generation; the packet remains local and no issue is represented as reconciled.
- One readiness issue fails its own acceptance criteria even though its mapped observations pass; that issue, #96, #122, promotion, and the milestone remain open.
- Promotion or checksum verification fails after issue comments are posted; the release remains a draft and the failure is recorded without moving the tag.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The disposition operation MUST consume the formal evidence archive, the independently staged candidate manifest, and the exact candidate MSI as three explicit inputs.
- **FR-002**: The operation MUST apply the complete production candidate, evidence, attachment, archive-content, and expected-identity validation rules before generating any disposition output.
- **FR-003**: The operation MUST accept only the formal attended evidence class and all 47 canonical observations in passing state.
- **FR-004**: The operation MUST generate exactly one readiness-issue record for each of #98, #101, #104, #105, #106, #109, #111, #112, and #113, plus one lifecycle-coordinator record for #96.
- **FR-005**: Each child record MUST include exactly the canonical observations mapped to that issue, preserving each observation identifier, status, summary, environment reference, timing, metrics, and attachment paths.
- **FR-006**: Every record MUST include the repository, tag, commit, workflow run link, run attempt, MSI filename, byte count, SHA-256, product version, ProductCode, and an explicit production-validator pass statement.
- **FR-007**: Every record MUST identify the complete referenced environments, including Windows version/build, account role, integrity, display class, effective DPI, and profile state.
- **FR-008**: The #96 coordinator record MUST identify its actual child and prerequisite issues #97, #98, #94, #89, and #90, preserve #98's independent closure boundary, and summarize the complete 36-observation pre-desktop lifecycle gate without claiming that generated text closed an issue.
- **FR-009**: Report ordering, field ordering, formatting, and filenames MUST be deterministic for identical inputs.
- **FR-010**: Generated text MUST be suitable for direct use as GitHub Markdown without embedding local workspace paths, secrets, access tokens, operator names, or mutable latest-release links. Evidence-recorded installed target paths remain part of their structured observation metrics.
- **FR-011**: Packet generation MUST be atomic and overwrite-protected: invalid inputs or output conflicts MUST leave no partial destination.
- **FR-012**: Diagnostics MUST state what failed and identify the affected candidate field, observation, attachment, archive member, or destination.
- **FR-013**: The operation MUST perform no network request and MUST NOT create tags, change releases, comment on or close issues, or change milestone state.
- **FR-014**: The operations runbook MUST preserve S048's tag-absence, draft-staging, nine-pre-checksum-payload, formal-evidence, promotion, and ten-final-payload gates in chronological order while replacing the obsolete S048-only commit check with the reviewed S049 merge commit.
- **FR-015**: The runbook MUST provide exact commands for generating the packet, reviewing each issue record, applying records only after individual acceptance review, and independently revalidating the evidence before promotion.
- **FR-016**: The runbook MUST distinguish the current pull-request authority from the separate authority required for tag staging and release publication.
- **FR-017**: The existing release workflows and the established 47-scenario evidence schema MUST remain backward-compatible; this slice MUST NOT weaken or replace their checks.
- **FR-018**: The change MUST introduce no new third-party runtime dependency and MUST pass all canonical repository gates.

### Canonical Issue Mapping

| Issue | Required formal observations |
| --- | --- |
| #98 | All nine `setup.*` observations and all seven `remove.*` observations |
| #101 | `desktop.appearance-standard`, `desktop.appearance-scaled` |
| #104 | `desktop.navigation-options`, `desktop.navigation-options-scaled`, `desktop.interaction-states`, `desktop.interaction-states-scaled` |
| #105 | `desktop.navigation-options`, `desktop.navigation-options-scaled`, `desktop.interaction-states`, `desktop.interaction-states-scaled` |
| #106 | `desktop.appearance-standard`, `desktop.appearance-scaled`, `desktop.navigation-options`, `desktop.navigation-options-scaled`, `desktop.scroll-input` |
| #109 | `desktop.interaction-states`, `desktop.interaction-states-scaled`, `desktop.tasks-table`, `desktop.tasks-table-scaled`, `desktop.schedule-activity-tables`, `desktop.schedule-activity-tables-scaled` |
| #111 | `desktop.scroll-input` |
| #112 | `desktop.tasks-table`, `desktop.tasks-table-scaled`, `desktop.interaction-states`, `desktop.interaction-states-scaled` |
| #113 | `desktop.schedule-activity-tables`, `desktop.schedule-activity-tables-scaled`, `desktop.interaction-states`, `desktop.interaction-states-scaled` |
| #96 | All 36 pre-desktop `access.*`, `window.*`, `error.*`, `task.*`, `setup.*`, and `remove.*` observations; coordinator references #97, #98, #94, #89, and #90 |

### Key Entities

- **Disposition Packet**: An immutable local directory containing one deterministic Markdown record per readiness issue and an index binding the files to their candidate identity.
- **Issue Disposition**: One issue number, its canonical observation set, the candidate identity, relevant environments and attachments, and a validation statement. It is review evidence, not a remote issue-state mutation.
- **Candidate Identity**: The repository, tag, commit, workflow run and attempt, MSI identity, product version, ProductCode, size, and digest shared by the staged manifest and attended evidence.
- **Formal Observation**: One canonical scenario result recorded against one identified Windows environment with metrics and supporting attachments.
- **Release Boundary**: The reviewed main commit that may receive the one immutable v1.0.0 tag only after separate authorization.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One valid formal candidate produces exactly ten Markdown records covering all 47 canonical observations and all ten remaining readiness issues with no unmapped or unexpected observation.
- **SC-002**: One mutation test for every candidate identity field, evidence class, required observation, attachment-integrity rule, archive-content rule, and destination-conflict class is rejected before any packet file is visible.
- **SC-003**: Two generations from identical inputs into separate empty destinations produce byte-for-byte identical files and file inventories.
- **SC-004**: A release steward can generate and inventory the complete packet in one command and under one minute on the supported development environment, excluding attended observation collection time.
- **SC-005**: Every generated child record provides enough traceability for a reviewer to locate its exact workflow run, candidate digest, environment, observation metrics, and attachment path without consulting another issue record.
- **SC-006**: The implementation pull request creates zero tags, releases, issue-state changes, milestone-state changes, or candidate artifacts.
- **SC-007**: The final runbook has one unambiguous stop condition for each of merge synchronization, tag staging, draft validation, formal qualification, issue reconciliation, promotion, and final audit.
- **SC-008**: All eight canonical local gates pass, including race, GUI, documentation, and release-automation policy checks.

## Scope Boundaries

### In Scope

- Deterministic, local, fail-closed generation of issue-specific formal evidence records from the existing v1.0.0 candidate inputs.
- Exact observation-to-issue mappings for the nine child issues and coordinator #96.
- Tests covering report completeness, determinism, safety, candidate binding, and command behavior.
- A chronological S049 operations runbook that reuses S048's release contract.
- Focused documentation for the new operation and its post-merge use.

### Out of Scope

- Creating or pushing `v1.0.0`, staging or publishing a GitHub release, closing issues, or closing the milestone during implementation of this pull request.
- Automating attended GUI judgment or importing S043/S047 local-demo results.
- Modifying the desktop application, daemon, task behavior, installer behavior, evidence schema, or the 47 canonical scenario definitions.
- Automatically deciding whether an issue's acceptance criteria are satisfied; generated records support, but never replace, individual reviewer judgment.
- Implementing Post-v1 issues #102, #103, #108, #110, #118, #119, or #120.

## Dependencies

- Depends on merged S047 pull request [#121](https://github.com/shruggietech/go-schedule/pull/121) and merged S048 pull request [#123](https://github.com/shruggietech/go-schedule/pull/123).
- Implements the issue-reconciliation support required by [#122](https://github.com/shruggietech/go-schedule/issues/122) without closing it.
- Reuses `specs/048-v100-release-cut/contracts/publication.md` as the authoritative tag, staging, promotion, and final-audit contract.

## Assumptions

- The reviewed S049 merge commit becomes the intended v1.0.0 source boundary. S049 is restricted to release tooling and documentation, so the packaged application behavior remains the S048-approved release candidate behavior.
- The formal evidence archive is collected only after a separately authorized tag-triggered draft exists and is never committed to the repository.
- GitHub issue comments and state changes are applied by the release steward after reviewing each generated file; credentials and network behavior remain outside the generator.
- The existing production validator is authoritative for evidence validity and is reused rather than reimplemented.
