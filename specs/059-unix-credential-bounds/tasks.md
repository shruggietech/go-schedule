# Tasks: Unix Credential Bounds

**Input**: Design documents from `specs/059-unix-credential-bounds/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Required by issue #145, the feature specification, and constitution. Regression tests precede implementation and must cover the complete credential boundary without host account provisioning.

**Organization**: Tasks are grouped by independently testable user story and execute chronologically under one-agent autopilot.

## Phase 1: Setup and Spec-Kit Gates

**Purpose**: Establish traceability and clear all pre-implementation gates.

- [x] T001 Register S059 as Draft in `specs/README.md` and update active plan context in `CLAUDE.md`
- [x] T002 Validate `specs/059-unix-credential-bounds/checklists/requirements.md` and `specs/059-unix-credential-bounds/checklists/security.md`
- [x] T003 Run the blocking Spec-Kit cross-artifact analysis over `specs/059-unix-credential-bounds/spec.md`, `plan.md`, and `tasks.md`, resolving every critical or high finding before implementation

---

## Phase 2: User Story 1 - Preserve the Resolved Unix Identity (Priority: P1) MVP

**Goal**: Accept exactly representable UID/GID pairs and reject all unsafe values before command mutation.

**Independent Test**: Synthetic resolved accounts cover both inclusive boundaries and every rejected UID/GID class, including a valid UID paired with an invalid GID, while caller-supplied command state remains unchanged on failure.

### Tests for User Story 1

- [x] T004 [US1] Add synthetic boundary and rejected-class regression cases for UID and GID in `internal/executor/runas_unix_test.go`
- [x] T005 [US1] Add command-state preservation cases for UID failure and valid-UID/invalid-GID failure in `internal/executor/runas_unix_test.go`
- [x] T006 [US1] Run the focused test before implementation and record the expected failure in `specs/059-unix-credential-bounds/verification.md`

### Implementation for User Story 1

- [x] T007 [US1] Implement unsigned 32-bit credential-pair validation with field and account context in `internal/executor/runas_unix.go`
- [x] T008 [US1] Refactor `applyRunAs` to assign process attributes, the complete credential pair, and environment only after validation succeeds in `internal/executor/runas_unix.go`
- [x] T009 [US1] Run focused executor tests and record passing evidence in `specs/059-unix-credential-bounds/verification.md`

---

## Phase 3: User Story 2 - Keep Valid Run-As Behavior Compatible (Priority: P2)

**Goal**: Preserve empty, named-account, numeric-account fallback, and home-directory behavior.

**Independent Test**: Existing and focused integration tests demonstrate unchanged valid behavior while the new parser accepts exact boundaries without host account dependencies.

- [x] T010 [US2] Extend compatibility assertions for empty `run_as`, current named account, credential values, and explicit or inherited home behavior in `internal/executor/runas_unix_test.go`
- [x] T011 [US2] Run the executor package test and race test on an available Unix target and record any platform prerequisite honestly in `specs/059-unix-credential-bounds/verification.md`

---

## Phase 4: User Story 3 - Produce Reviewable Security Evidence (Priority: P3)

**Goal**: Demonstrate that the shared CodeQL cause is removed without suppressions or unrelated changes.

**Independent Test**: Static search finds no signed or architecture-sized narrowing into process credential fields, local verification passes, and hosted CodeQL reports no equivalent finding.

- [x] T012 [US3] Audit process UID/GID assignments and document the boundary disposition in `specs/059-unix-credential-bounds/verification.md`
- [x] T013 [US3] Add S059 correction and decision entries under `[Unreleased]` in `CHANGELOG.md`
- [x] T014 [US3] Update issue #145 closure eligibility and remaining hosted CodeQL evidence in `specs/059-unix-credential-bounds/verification.md`

---

## Phase 5: Polish and Completion Gates

**Purpose**: Prove full repository readiness and prepare the authorized review publication.

- [x] T015 Run final Spec-Kit analysis and validate both S059 checklists against the delivered implementation
- [x] T016 Run `go run ./scripts/github-format` and resolve every repository-publication formatting defect
- [x] T017 Audit changed files for UTF-8 without BOM, mojibake, whitespace errors, and unintended scope changes
- [x] T018 Run `sh scripts/verify.sh all` in the foreground and record format, vet, lint, race, GUI, coverage, docs, and automation evidence in `specs/059-unix-credential-bounds/verification.md`
- [x] T019 Advance S059 to Implemented with objective delivery evidence in `specs/059-unix-credential-bounds/spec.md` and `specs/README.md`
- [x] T020 Commit the review-ready slice as `feat(059): reject invalid Unix credential IDs` with the required co-author trailer

---

## Dependencies and Execution Order

- Phase 1 establishes the Spec-Kit and lifecycle gates.
- User Story 1 is the security MVP and blocks all later phases.
- User Story 2 depends on the corrected boundary and proves compatibility.
- User Story 3 depends on the final implementation and produces closure evidence.
- Phase 5 depends on all user stories and publishes only after the already granted authorization is reconfirmed by a green local gate.

## Parallel Opportunities

- No active delegation is authorized. The `[P]` marker is intentionally unused because the focused tests and implementation share one build-tagged module and must preserve test-first order.
- Documentation artifacts can be reviewed independently after implementation, but their final evidence depends on the same verified commit.

## Implementation Strategy

1. Complete the Spec-Kit gate and mark S059 In Progress before behavioral implementation.
2. Write the synthetic boundary and atomicity regression tests and record the red result.
3. Implement the smallest unsigned fixed-width parser and validate the complete pair before mutation.
4. Re-run focused compatibility and security tests, then audit all credential assignments.
5. Finish publication records, canonical verification, Implemented lifecycle state, and the local commit.
6. Use the user's explicit authorization to push, open the PR, resolve CI and every review comment, request at most one second Codex review round, and stop for the maintainer's final review and merge ritual.
