# Tasks: v0.9 Closure and Maintenance Automation

**Input**: Design documents from `specs/029-v09-closure/` **Tests**: Required by constitution Principle II and the S029 lifecycle/automation contract.

## Phase 1: Setup and Evidence

- [x] T001 Record the pre-change spec inventory, delivery references, issue states, and repository security settings in `specs/029-v09-closure/verification.md`.
- [x] T002 Update the active feature context between the Spec-Kit markers in `CLAUDE.md`.

## Phase 2: Foundational Lifecycle Contract

- [x] T003 Define lifecycle states, transitions, evidence, and task-disposition rules in `specs/README.md`.
- [x] T004 [P] Add failing lifecycle contradiction fixtures in `test/scripts/spec-lifecycle-check_test.sh`.
- [x] T005 Implement the offline lifecycle contract in `scripts/spec-lifecycle-check.sh` until the fixtures pass.
- [x] T006 Integrate the lifecycle checker into `scripts/automation-check.sh` and its fixtures in `test/scripts/automation-check_test.sh`.
- [x] T007 Update `.specify/templates/spec-template.md` and `docs/build-autopilot.md` so future slices maintain status and exclude publication bookkeeping from implementation completion.

## Phase 3: User Story 1 - Trust the specification history (Priority: P1)

**Independent Test**: The audited repository passes the lifecycle checker and each contradiction fixture fails for its expected reason.

- [x] T008 [US1] Add delivery evidence and accurate Implemented status to every historical `specs/001-*/spec.md` through `specs/028-*/spec.md`.
- [x] T009 [US1] Reconcile stale checkbox metadata in `specs/003-gui-editor-refinements/tasks.md`, `specs/005-gui-task-fidelity/tasks.md`, `specs/006-maintainer-test-scripts/tasks.md`, and `specs/012-github-security-baseline/tasks.md` through `specs/015-docs-dark-theme/tasks.md` without claiming unperformed manual evidence.
- [x] T010 [US1] Record the legacy task-format disposition for `specs/007-issue-cleanup-docs/tasks.md` and complete the 29-spec inventory in `specs/README.md`.
- [x] T011 [US1] Run `sh test/scripts/spec-lifecycle-check_test.sh` and `sh scripts/spec-lifecycle-check.sh .`.

## Phase 4: User Story 2 - Receive dependency work as proposals (Priority: P1)

**Independent Test**: The offline automation fixtures accept the exact two-ecosystem Dependabot contract and reject missing or noisy variants.

- [x] T012 [US2] Add failing Dependabot ecosystem, cadence, grouping, label, and limit fixtures to `test/scripts/automation-check_test.sh`.
- [x] T013 [US2] Add the low-noise Go module and GitHub Actions policy to `.github/dependabot.yml`.
- [x] T014 [US2] Extend `scripts/automation-check.sh` to enforce the Dependabot contract, then run the focused automation fixtures.
- [x] T015 [US2] Confirm Dependabot security updates remain enabled through hosted repository readback and record the result in `specs/029-v09-closure/verification.md`.

## Phase 5: User Story 3 - Finish hosted security controls (Priority: P1)

**Independent Test**: GitHub readback reports the requested three controls enabled and unrelated controls unchanged.

- [x] T016 [US3] Enable push protection, non-provider patterns, and validity checks individually through the GitHub repository API.
- [x] T017 [US3] Read back all security-and-analysis fields and record requested, observed, and unchanged states in `specs/029-v09-closure/verification.md`.

## Phase 6: User Story 4 - Keep v0.9 actionable (Priority: P2)

**Independent Test**: Issue #33 remains open in Post-v1 with P3 and `needs: verification`, and the v0.9 milestone contains only S029 merge-closing issues.

- [x] T018 [US4] Move issue #33 to `Post-v1`, replace P1 with P3, retain `needs: verification`, and post an honest deferral comment.
- [x] T019 [US4] Verify issue #33 and v0.9 milestone state through GitHub readback and record it in `specs/029-v09-closure/verification.md`.

## Phase 7: Documentation and Verification

- [x] T020 Add Unreleased entries and dated process decisions to `CHANGELOG.md` for lifecycle enforcement, Dependabot, hosted security settings, and proportional #33 deferral.
- [x] T021 Run shellcheck for changed shell tests and scripts, then run `sh scripts/verify.sh all` in the foreground through all eight gates.
- [x] T022 Audit the diff, hosted state, UTF-8 without BOM, mojibake, trailing whitespace, and issue-closing language; record results in `specs/029-v09-closure/verification.md`.
- [x] T023 Complete `specs/029-v09-closure/checklists/maintenance-contract.md`, mark all required tasks resolved, and update `specs/029-v09-closure/spec.md` to Implemented with local delivery evidence.

## Dependencies and Execution Order

- Phase 2 establishes the contract before historical reconciliation or automation configuration.
- US1, US2, and US3 are independently testable after Phase 2. US4 is independent of code but precedes the final milestone audit.
- T021 through T023 require every implementation and hosted mutation to be complete.

## Parallel Opportunities

- T004 can proceed while the evidence inventory is assembled.
- Historical spec headers in T008 are mechanically independent, but one final audit must ensure inventory consistency.
- Hosted security activation and issue #33 triage touch separate GitHub resources after their contracts are fixed.

## Implementation Strategy

Deliver the complete closure slice as one pull request. No sub-slice is published independently because the value is the trustworthy v0.9 boundary: accurate history, durable maintenance channels, finished supported security controls, and a backlog containing only actionable work.
