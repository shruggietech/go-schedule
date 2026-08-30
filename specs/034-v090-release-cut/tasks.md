# Tasks: v0.9.0 Release Cut

**Input**: Design documents from `specs/034-v090-release-cut/`

**Tests**: Required by constitution Principle II and the S034 release-copy
contract.

## Phase 1: Setup and Baseline

**Purpose**: Capture the exact pre-release boundary and active hosted state.

- [x] T001 Record the current tag, release, milestone, issue, branch, and
  workflow baseline in `specs/034-v090-release-cut/verification.md`.
- [x] T002 Record the current Unreleased changelog range and expected v0.9.0
  artifact inventory in `specs/034-v090-release-cut/verification.md`.
- [x] T003 Update the active Spec-Kit context in `CLAUDE.md` and add S034 to
  `specs/README.md`.

## Phase 2: Foundational Release-Copy Contract

**Purpose**: Add failing policy fixtures before changing the release workflow.

- [x] T004 [P] Add a passing tag-specific highlights fixture to
  `test/scripts/automation-check_test.sh`.
- [x] T005 [P] Add failing fixtures for generated release notes and a fixed
  v0.9.0 body path to `test/scripts/automation-check_test.sh`.
- [x] T006 [P] Add failing fixtures for missing changelog linkage, invalid
  highlight count, and exhaustive-copy headings to
  `test/scripts/automation-check_test.sh`.
- [x] T007 Run the focused automation fixtures and record the expected red
  baseline in `specs/034-v090-release-cut/verification.md`.

## Phase 3: User Story 1 - Understand the Release Quickly (Priority: P1)

**Goal**: Publish only concise highlights with one path to full detail.

**Independent Test**: The offline checker accepts the reviewed v0.9.0 notes and
rejects generated, stale, overlong, or unlinked variants.

- [x] T008 [US1] Implement the tag-specific release-note policy in
  `scripts/automation-check.sh` until the new fixtures pass.
- [x] T009 [US1] Replace generated release copy with dynamic `body_path` lookup
  in `.github/workflows/release.yml`.
- [x] T010 [US1] Write four to six v0.9.0 highlights and the full-changelog link
  in `.github/release-notes/v0.9.0.md`.
- [x] T011 [US1] Run the focused automation fixtures and the repository
  automation checker against the new release-copy contract.

## Phase 4: User Story 2 - Preserve an Honest Version Boundary (Priority: P1)

**Goal**: Cut the detailed post-v0.8.0 history under v0.9.0 without losing a
future Unreleased section.

**Independent Test**: Every pre-cut Unreleased entry remains under v0.9.0, the
new Unreleased section is empty, and comparison links identify both boundaries.

- [x] T012 [US2] Add the new Unreleased and dated v0.9.0 headings to
  `CHANGELOG.md` without altering the accumulated entry order.
- [x] T013 [US2] Update the Unreleased and v0.9.0 comparison links in
  `CHANGELOG.md`.
- [x] T014 [US2] Add the S034 pinned-workflow decision and release-note policy
  to the v0.9.0 Decisions section in `CHANGELOG.md`.
- [x] T015 [US2] Audit pre-cut and post-cut headings, entries, and links against
  the baseline in `specs/034-v090-release-cut/verification.md`.

## Phase 5: User Story 3 - Download a Complete Verified Build (Priority: P1)

**Goal**: Make the existing tag workflow ready to publish and objectively audit
the complete supported artifact set after authorization.

**Independent Test**: Static workflow inspection preserves every build and
checksum job, and the publication contract provides exact post-tag audit steps.

- [x] T016 [US3] Verify release-note lookup does not alter package jobs,
  dependency ordering, explicit asset globs, or checksum behavior in
  `.github/workflows/release.yml`.
- [x] T017 [US3] Complete the preparation, tag, publication, and completion
  contracts in `specs/034-v090-release-cut/contracts/publication.md`.
- [x] T018 [US3] Complete the runnable local and post-tag audit procedure in
  `specs/034-v090-release-cut/quickstart.md`.
- [x] T019 [US3] Record the exact expected asset names, tag guard, README audit,
  and issue #79 closure gate in
  `specs/034-v090-release-cut/verification.md`.

## Phase 6: Analysis, Documentation, and Verification

**Purpose**: Prove the preparation is complete without prematurely tagging or
closing the release issue.

- [x] T020 Run the blocking cross-artifact analysis and resolve all critical or
  high findings across `spec.md`, `plan.md`, and `tasks.md`.
- [x] T021 Run ShellCheck for changed shell files and the focused automation
  fixture suite.
- [x] T022 Run `sh scripts/verify.sh all` in the foreground through format, vet,
  lint, race, GUI, coverage, docs, and automation.
- [x] T023 Audit UTF-8 without BOM, mojibake, trailing whitespace, changelog
  retention, release-note size, issue references, and branch cleanliness;
  record results in `specs/034-v090-release-cut/verification.md`.
- [x] T024 Complete both S034 checklists, mark the release preparation
  Implemented in `specs/034-v090-release-cut/spec.md`, and record local delivery
  evidence while leaving issue #79 open for publication.

## Dependencies and Execution Order

- Phase 1 establishes the immutable baseline before any release boundary moves.
- Phase 2 must fail for the new negative fixtures before US1 implementation.
- US1 and US2 may proceed independently after Phase 2, but both are required
  before the full release preparation is coherent.
- US3 depends on the final workflow and documentation shapes from US1 and US2.
- Phase 6 requires every implementation task to be complete.

## Parallel Opportunities

- T004 through T006 touch independent fixture variants and can be drafted
  together before the checker changes.
- The v0.9.0 highlight copy and changelog boundary can be reviewed independently
  before their final cross-link audit.
- Static workflow inspection and publication-procedure documentation touch
  different files after the release-copy contract is fixed.

## Implementation Strategy

Deliver one substantial release-preparation pull request. The local
implementation ends with a verified workflow, changelog boundary, curated notes,
and publication contract. Branch publication, review, merge, tag creation,
workflow execution, artifact audit, and issue closure remain chronological
workflow evidence outside this task list, with separate authorization at the
branch and tag boundaries.
