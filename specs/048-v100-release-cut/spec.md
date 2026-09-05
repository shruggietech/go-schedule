# Feature Specification: v1.0.0 Release Cut

**Feature Branch**: `codex/048-v100-release-cut`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/048-v100-release-cut`; release-preparation boundary implemented and locally verified; hosted CI and review remain pull-request evidence

**Input**: User request to prepare the S048 v1.0.0 release slice under Spec Kit, publish its reviewed preparation pull request, address up to two Codex review rounds, and hand the green pull request back for the final merge ritual.

**Tracking**: GitHub issue [#122](https://github.com/shruggietech/go-schedule/issues/122)

## Scope and authorization boundary

S048 prepares the reviewed v1.0.0 source boundary and the exact post-merge qualification ritual. This pull request does not create a tag, stage a draft release, upload attended evidence, dispatch promotion, or close issues whose acceptance depends on that evidence. Repository policy treats tags and releases as separate publication authority. Those operations begin only after this pull request merges and the maintainer explicitly authorizes the v1.0.0 tag.

This is an explicit scope correction from the initial high-level S048 recommendation. Trying to qualify a candidate before merging its release preparation would tag a commit that omits S048's reviewed version boundary and release notes. The preparation PR must therefore come first.

## User Scenarios & Testing

### User Story 1 - Review the exact v1 release boundary (Priority: P1)

As the maintainer, I can review one coherent pull request that identifies the complete v1.0.0 change set, tag-specific health example, publication-aware release badge, and release notes before any tag is created.

**Why this priority**: The tag is difficult to correct after publication. The reviewed repository state must be the sole candidate source.

**Independent Test**: Compare the prepared changelog with v0.9.1 and verify that every accumulated Unreleased entry moves intact beneath a dated v1.0.0 heading, a new empty Unreleased section remains, comparison links are correct, the health example identifies v1.0.0, and the release badge follows the latest published GitHub release.

**Acceptance Scenarios**:

1. **Given** the accumulated post-v0.9.1 history, **When** S048 prepares the boundary, **Then** all entries appear once beneath `1.0.0` and none remain in the new Unreleased section.
2. **Given** the Release workflow's README preflight, **When** a future v1.0.0 tag evaluates the reviewed commit, **Then** the health example matches the tag and the release badge remains bound to GitHub's latest published release.
3. **Given** no v1.0.0 tag or release exists, **When** the preparation PR is reviewed or merged, **Then** the visible release badge continues to identify v0.9.1 and publication is described as pending rather than shipped.

---

### User Story 2 - Understand v1.0.0 quickly (Priority: P1)

As a prospective user, I can scan concise v1.0.0 highlights and follow one link to the complete tagged changelog.

**Why this priority**: v1.0.0 contains a large body of work. Release notes must surface the user-visible outcomes without reproducing the complete history.

**Independent Test**: Validate the tag-specific notes with the repository's offline release-copy policy and manually confirm four to six meaningful bullets, one Highlights heading, one versioned changelog link, and no generated inventory.

**Acceptance Scenarios**:

1. **Given** the v1.0.0 draft release, **When** a reader scans its body, **Then** four to six highlights summarize the release without an exhaustive commit or pull-request list.
2. **Given** the reader wants detail, **When** they follow the changelog link, **Then** it targets the v1.0.0 section at the v1.0.0 tag.
3. **Given** a later tag, **When** its workflow runs, **Then** dynamic tag lookup prevents reuse of the v1.0.0 notes.

---

### User Story 3 - Execute the post-merge release safely (Priority: P1)

As the release operator, I have one exact, chronological contract for staging the draft, validating its Windows MSI, collecting all 47 native observations, reconciling every remaining v1 issue, and promoting only the tested artifacts.

**Why this priority**: S047 deliberately separates local-demo confidence from formal candidate proof. S048 must preserve that distinction through release.

**Independent Test**: Review the contract against the Release and Promote Release workflows and verify that each command binds tag, commit, workflow run, candidate manifest, MSI, attended archive, issue dispositions, and promotion.

**Acceptance Scenarios**:

1. **Given** the merged S048 commit and explicit tag authorization, **When** the tag is created, **Then** it points to the synchronized reviewed `main` commit and the Release workflow leaves a draft release with the complete staged set.
2. **Given** the staged MSI and manifest, **When** formal evidence is collected, **Then** all 47 required observations bind to that exact candidate and pass the production validator.
3. **Given** one failed/unavailable observation or unmet issue criterion, **When** reconciliation occurs, **Then** the affected issue and draft release remain open and promotion is not dispatched.
4. **Given** all evidence and issue criteria pass, **When** promotion runs, **Then** it verifies the complete checksummed asset set before publishing the same draft release.

### Edge Cases

- `main` advances after the S048 merge but before tagging; the boundary must be re-reviewed rather than tagging an unreviewed commit.
- A v1.0.0 tag or release unexpectedly exists; staging must stop before mutation.
- The Release workflow partially uploads artifacts or a rerun uses another attempt; the candidate manifest and successful Windows job must identify the accepted run attempt exactly.
- README or release notes mention v1.0.0 while the changelog link or heading is stale; offline verification must fail.
- A local-demo MSI resembles the candidate; only the draft-release artifact and manifest may enter formal evidence.
- An attended attachment has a misleading extension or media type; its raster bytes and hash must satisfy the S047 validator.
- One issue shares observations with another but has unmet individual criteria; it remains open independently.
- Promotion succeeds but the latest-release pointer, checksum inventory, or milestone state is inconsistent; #122 remains open until corrected.

## Requirements

### Functional Requirements

- **FR-001**: S048 MUST prepare version `v1.0.0` from the complete reviewed change set after `v0.9.1`.
- **FR-002**: `CHANGELOG.md` MUST retain every current Unreleased entry exactly once beneath a dated `1.0.0` heading and MUST retain a new empty Unreleased section for future work.
- **FR-003**: Changelog comparison links MUST compare Unreleased against v1.0.0 and v1.0.0 against v0.9.1.
- **FR-004**: `README.md` MUST contain exactly one release badge whose image is derived from GitHub's latest published release and one `daemon ok (version 1.0.0)` health example before tag staging. The tag preflight MUST validate both forms without requiring an unpublished badge.
- **FR-005**: `.github/release-notes/v1.0.0.md` MUST contain one Highlights heading, four to six concise bullets, one tagged full-changelog link, and no generated or exhaustive change inventory.
- **FR-006**: The release preparation MUST preserve dynamic tag-specific notes, draft-only staging, exact-candidate manifests, attended evidence validation, final checksum generation, and promotion ordering.
- **FR-007**: The post-merge contract MUST stop unless local `main`, `origin/main`, and the reviewed S048 merge commit are identical and clean.
- **FR-008**: The post-merge contract MUST require absence of an existing v1.0.0 tag and public release before tag creation.
- **FR-009**: Staging MUST produce exactly the nine pre-checksum assets required by the current promotion workflow, including the Windows MSI, candidate manifest, and attended evidence archive.
- **FR-010**: Formal qualification MUST use all 47 S047 scenarios and MUST bind the evidence archive to the exact v1.0.0 MSI, manifest, tag commit, Release run ID, and run attempt.
- **FR-011**: Formal Windows evidence MUST remain attended and MUST NOT copy pass results from S043 or S047 local-demo evidence.
- **FR-012**: S048 MUST map #96, #98, #101, #104, #105, #106, #109, #111,
  #112, and #113 to their exact formal observations and reconcile each issue
  independently.
- **FR-013**: Issue #122 and milestone `v1.0.0 - Release readiness` MUST remain open through the preparation PR and until the published release and all issue dispositions are audited.
- **FR-014**: The preparation PR MUST use `Refs #122`; merging preparation alone MUST NOT close the release issue or claim v1.0.0 is published.
- **FR-015**: Creating the tag, uploading evidence, dispatching promotion, closing issues, and closing the milestone MUST require post-merge authority and successful upstream gates.
- **FR-016**: All focused release checks and all eight canonical repository gates MUST pass locally before S048 is pushed.
- **FR-017**: All hosted pull-request checks and permitted Codex review rounds MUST be satisfied before the maintainer is asked to merge.
- **FR-018**: S048 MUST NOT implement Post-v1 issues #102 or #118 through #120, nor include their behavior in the v1.0.0 release boundary.

### Key Entities

- **Release Boundary**: The reviewed merge commit, tag, date, changelog range, tag-specific README health example, publication-aware badge contract, and curated notes that define v1.0.0.
- **Staged Candidate**: The draft-release Windows MSI and candidate manifest produced by one successful Release workflow attempt from the tagged commit.
- **Formal Evidence Bundle**: The attended, hashed archive containing exactly 47 passing observations bound to the staged candidate.
- **Issue Disposition**: One v1 issue's acceptance criteria, mapped observations, evidence links, and final open/closed decision.
- **Promotion Record**: The successful workflow run, final asset inventory, checksums, public release URL, latest pointer, and milestone audit.

## Success Criteria

### Measurable Outcomes

- **SC-001**: One hundred percent of the post-v0.9.1 Unreleased entries are retained beneath v1.0.0 with no duplication or omission.
- **SC-002**: Release notes contain four to six highlights, one full-changelog link, and zero generated commit or pull-request entries.
- **SC-003**: README health output, changelog, release notes, tag contract, and expected asset names identify v1.0.0 consistently, while the README badge continues to reflect only the latest published GitHub release.
- **SC-004**: All focused checks, eight local canonical gates, hosted CI checks, and permitted Codex review comments complete without unresolved failures.
- **SC-005**: The post-merge contract accounts for all 47 observations, all nine pre-checksum assets, and all ten remaining v1.0.0 issues.
- **SC-006**: No tag, release, evidence upload, issue closure, milestone closure, or promotion occurs merely because the preparation PR is merged.

## Clarifications

### Session 2026-09-03

- Q: Does the S048 preparation PR itself create or publish v1.0.0? A: No. The reviewed preparation merge must precede separately authorized tag staging.
- Q: Can the accepted S047 local demo satisfy formal qualification? A: No. All 47 observations must be repeated against the exact draft-release MSI.
- Q: Are remaining v1 issues closed together? A: No. Shared evidence can support several issues, but each issue closes only after its own criteria pass.
- Q: Does S048 include Post-v1 usability work? A: No.

## Assumptions

- PR #121 is merged at `d4ad27540331eae943e43a3830d4f5c9a6c38afa` and `main` is clean and synchronized when S048 begins.
- v0.9.1 is the latest tag and public release; v1.0.0 does not exist.
- The target release date remains 2026-09-03 unless review crosses the local calendar boundary.
- GitHub Actions, draft releases, and the configured Codex reviewer remain available.
- The physical Windows environment needed for formal evidence will be available after tag staging; lack of required attended evidence blocks promotion.

## Out of Scope

- Product behavior changes without a newly reproduced v1 release blocker.
- Post-v1 issues #102 and #118 through #120.
- Adding supported platforms, package formats, signing, or distribution stores.
- Weakening or replacing S047's 47-scenario evidence contract.
- Creating the v1.0.0 tag or public release before the reviewed preparation merge and explicit post-merge authorization.
