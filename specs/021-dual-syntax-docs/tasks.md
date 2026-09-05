# Tasks: Dual-Syntax Product Documentation

**Input**: Design documents from `specs/021-dual-syntax-docs/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required. Behavioral and policy tests are recorded red before their implementation changes.

## Phase 1: Setup and Baseline

- [X] T001 Record focused CLI/GUI and documentation-check baselines in `verification.md`
- [X] T002 Research the cron dialect, current surface inventory, historical policy, and primary references

---

## Phase 2: User Story 1 - Truthful Scheduling Semantics (Priority: P1)

**Goal**: Lossy calendar wildcard steps are named refusals, while supported human and cron inputs retain their existing behavior.

### Tests

- [X] T003 [US1] Add failing parser tests for day-of-month, month, and day-of-week wildcard-step refusals and `*/1` compatibility
- [X] T004 [US1] Add a failing shared-boundary regression proving a refused form does not fall back to human parsing
- [X] T005 [US1] Run focused tests and record the expected pre-implementation failures

### Implementation

- [X] T006 [US1] Reject wildcard steps greater than one in calendar fields with a named fidelity reason
- [X] T007 [US1] Run cron and schedule-input focused suites green and record evidence

---

## Phase 3: User Story 2 - Consistent Current Product Surfaces (Priority: P1)

**Goal**: README, guides, CLI help, and authoritative contracts present one human-first dual-syntax posture and equivalent examples.

- [X] T008 [US2] Add or update focused CLI help assertions for both authoring paths
- [X] T009 [US2] Align `README.md`, `docs/README.md`, `docs/cli.md`, `docs/cron.md`, and `docs/gui-fields.md`
- [X] T010 [US2] Align root CLI help and remove obsolete human-only test wording
- [X] T011 [US2] Update S001 spec, data model, quickstart, and CLI/local API contracts
- [X] T012 [US2] Run focused CLI and documentation checks green

---

## Phase 4: User Story 3 - Durable Documentation Policy (Priority: P2)

**Goal**: History remains intact, current copy cannot regress silently, and issue closure semantics are explicit.

### Tests

- [X] T013 [US3] Add a fixture harness proving aligned copy passes and targeted obsolete claims fail
- [X] T014 [US3] Run the fixture harness and record the expected pre-implementation failure

### Implementation

- [X] T015 [US3] Add S008 supersession notices to its spec, tasks, quickstart, and fidelity checklist
- [X] T016 [US3] Add the bounded current-surface policy helper and integrate it into `docs-check.sh`
- [X] T017 [US3] Update the S010 docs-check contract and run policy fixtures and docs checks green

---

## Phase 5: Polish and Cross-Cutting Verification

- [X] T018 Add chronological Unreleased changelog entries for documentation policy and the narrow named-refusal correction
- [X] T019 Validate both checklists and run Spec-Kit analysis to zero critical findings
- [X] T020 Run `git diff --check`, strict UTF-8-without-BOM, and mojibake audits across changed files
- [X] T021 Run `sh scripts/verify.sh all` in the foreground and record all eight gate results
- [X] T022 Mark every completed task `[X]`, rerun analysis, and commit locally with the required co-author trailer

## Dependencies and Execution Order

- Baseline and research precede implementation.
- US1 correctness tests and fix precede claims about the supported subset.
- US2 establishes the current copy consumed by the US3 policy inventory.
- US3 fixtures precede the policy implementation.
- Analysis, audits, and all canonical gates follow every story.

## Implementation Strategy

1. Protect timing semantics with red tests and the smallest refusal fix.
2. Align current product and authoritative contract surfaces.
3. Preserve historical evidence with supersession notices.
4. Lock the posture with fixture-backed bounded policy checks.
5. Analyze, verify, commit locally, and halt before push.
