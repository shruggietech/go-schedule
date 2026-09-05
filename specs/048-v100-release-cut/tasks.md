# Tasks: v1.0.0 Release Cut

**Input**: Design documents from `specs/048-v100-release-cut/`

**Tests**: Required by the constitution and the release-security scope. Existing release-copy, workflow, evidence, and canonical gates must all pass.

**Organization**: Tasks are chronological and grouped by independently testable release outcomes. Tagging, release publication, hosted PR review, and branch cleanup are workflow evidence rather than implementation tasks.

## Phase 1: Setup and Baseline

- [x] T001 Create `codex/048-v100-release-cut`, activate the Spec Kit workspace, and create release coordinator issue #122
- [x] T002 Record the clean S047 merge, v0.9.1 tag/release, v1.0.0 absence, Unreleased cardinality/hash, commit range, milestone state, and expected assets in `specs/048-v100-release-cut/verification.md`
- [x] T003 Add S048 as In Progress in `specs/README.md` and update active feature context in `CLAUDE.md`

---

## Phase 2: Specification, Clarification, and Design

- [x] T004 Complete the S048 specification and clarification record in `specs/048-v100-release-cut/spec.md`
- [x] T005 Complete the 20-item requirement-quality checklist in `specs/048-v100-release-cut/checklists/requirements.md`
- [x] T006 Complete the constitution-checked implementation plan and seven research decisions in `specs/048-v100-release-cut/plan.md` and `specs/048-v100-release-cut/research.md`
- [x] T007 Define release preparation, staged candidate, asset, evidence, issue-disposition, and release-state entities in `specs/048-v100-release-cut/data-model.md`
- [x] T008 Define the fail-closed post-merge ritual in `specs/048-v100-release-cut/contracts/publication.md` and complete its 22-item checklist

**Checkpoint**: The preparation PR and separately authorized release operation have one unambiguous boundary.

---

## Phase 3: User Story 1 - Review the Exact v1 Boundary (Priority: P1)

**Goal**: Reviewers can approve the complete v1.0.0 source and version boundary before any tag exists.

**Independent Test**: All 33 pre-cut entries remain exactly once beneath the dated v1.0.0 heading; future Unreleased is empty; comparison links and health output identify v1.0.0; the README badge follows the latest published release.

- [x] T009 [US1] Cut the current `CHANGELOG.md` Unreleased content beneath `## [1.0.0] - 2026-09-03` while preserving every entry
- [x] T010 [US1] Update the Unreleased and v1.0.0 comparison links in `CHANGELOG.md`
- [x] T011 [US1] Bind the single README badge to the latest published GitHub release and change the health output version to v1.0.0
- [x] T012 [US1] Audit changelog entry retention, section boundaries, comparison links, and README version cardinality against the recorded baseline

---

## Phase 4: User Story 2 - Understand v1.0.0 Quickly (Priority: P1)

**Goal**: The GitHub release body communicates the major outcomes concisely and links to the full immutable record.

**Independent Test**: The existing offline policy accepts exactly one v1.0.0 notes file with five bullets, one Highlights heading, and one tagged changelog link while rejecting all established malformed fixtures.

- [x] T013 [US2] Add five concise v1.0.0 highlights and one tagged full-changelog link in `.github/release-notes/v1.0.0.md`
- [x] T014 [US2] Run the release-note fixture suite and repository automation checker against the v1.0.0 notes

---

## Phase 5: User Story 3 - Execute the Post-Merge Release Safely (Priority: P1)

**Goal**: The release operator has a complete, ordered, candidate-bound ritual without introducing another workflow or bypass.

**Independent Test**: Static workflow and integration-contract checks agree with the documented nine-payload, 47-observation, ten-issue promotion sequence.

- [x] T015 [US3] Inspect and preserve the existing Release/Promote workflow ordering, draft state, candidate manifest, evidence, checksum, and immutable-tag guards
- [x] T016 [US3] Map all ten remaining v1 readiness issues to exact formal observations in `specs/048-v100-release-cut/data-model.md`
- [x] T017 [US3] Document preparation, tag, staging, formal evidence, reconciliation, promotion, and final audit stop conditions in `specs/048-v100-release-cut/contracts/publication.md`
- [x] T018 [US3] Provide the runnable local and post-merge operator sequence in `specs/048-v100-release-cut/quickstart.md`

---

## Phase 6: Analysis and Verification

- [x] T019 Run `/speckit-analyze`, resolve every critical/high finding and all requirement/task coverage gaps, and record the result in `specs/048-v100-release-cut/verification.md`
- [x] T020 Run focused release-gate, workflow-contract, release-copy fixture, automation, and documentation checks
- [x] T021 Run foreground `scripts/verify.sh all` through all eight canonical gates
- [x] T022 Audit UTF-8 without BOM, mojibake, unresolved placeholders, task format, release-copy cardinality, changelog retention, and `git diff --check`
- [x] T023 Mark S048 Implemented with objective local evidence and create the verified release-preparation commit

---

## Phase 7: Pull Request Review

- [x] T024 [US1] Resolve PR #123's P1 finding by binding the badge to published GitHub release state and decoupling badge validation from the future tag
- [x] T025 Add a negative regression fixture and rerun focused, canonical, lifecycle, and integrity checks after the review fix
- [x] T026 Correct PR #123's second-round audit finding by documenting the narrow Release metadata-preflight change accurately
- [x] T027 Rerun documentation, automation, lifecycle, encoding, placeholder, task-identity, and whitespace checks after the second-round fix

---

## Dependencies and Execution Order

```text
baseline -> specification/design -> reviewed v1 boundary -> release highlights
         -> post-merge contract -> analysis -> focused verification
         -> canonical verification -> committed preparation -> review fix
```

- The changelog baseline must be recorded before the boundary moves.
- README and release notes must be present in the commit that will later be tagged, so preparation review precedes tag staging.
- Formal evidence cannot exist before the tag-triggered Release workflow stages the exact candidate and remains outside this implementation task list.
- Issue closure and promotion depend on the future formal evidence and cannot be represented as completed by the preparation PR.

## Implementation Strategy

Deliver the smallest coherent v1.0.0 preparation: one preserved changelog cut, one curated note file, two synchronized README lines, and a complete reuse-based release contract. Do not modify runtime code or the already sufficient release workflows unless a focused gate proves a release-blocking defect.
