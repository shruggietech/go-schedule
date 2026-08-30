# Tasks: GitHub Security Baseline

**Input**: Design documents from `specs/012-github-security-baseline/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`,
`contracts/`, `quickstart.md`

**Tests**: Required by constitution Principle II. Workflow-policy fixtures are
written and observed failing before implementation.

**Task Reconciliation**: PR #44 delivered the local baseline. Hosted activation was completed across subsequent maintainer sessions, with the final requested secret-scanning controls verified by S029 on 2026-08-30.

## Phase 1: Setup and Baseline Evidence

**Purpose**: Fix the feature boundary and preserve evidence of the pre-change
repository contract.

- [X] T001 Confirm local `main`, active Feature 012, and the pre-change remote
  security state in `specs/012-github-security-baseline/quickstart.md` using
  read-only repository and GitHub API inspection.
- [X] T002 Confirm `.github/workflows/ci.yml`, product code, dependency files,
  branch protection, and release automation are out of scope and capture the
  planned file set in `specs/012-github-security-baseline/plan.md`.

---

## Phase 2: Foundational Workflow Contract

**Purpose**: Establish failing regressions before the security workflow and
policy implementation exist.

**CRITICAL**: Required negative fixtures must fail for the expected reasons
before the checker is extended.

- [X] T003 [US2] Extend `test/scripts/automation-check_test.sh` with a valid
  advanced CodeQL fixture plus obsolete CodeQL v3, missing required trigger,
  insufficient `security-events` permission, and missing analysis-step
  fixtures; retain the existing unknown-action regression.
- [X] T004 [US2] Run the extended automation fixture against the existing
  `scripts/automation-check.sh` and record the expected pre-implementation
  failure in `specs/012-github-security-baseline/quickstart.md`.

**Checkpoint**: Tests express all five required drift classes and demonstrate
that the previous automation policy cannot yet certify the feature.

---

## Phase 3: User Story 1 - Researchers can report privately (Priority: P1)

**Goal**: Make the reporting policy truthful and ready for activation, with an
independent fallback and accountable triage.

**Independent Test**: Read `SECURITY.md` and confirm it contains the repository
private-advisory URL, `info@shruggie.tech`, the public-disclosure warning,
repository-administrator ownership, and one-week acknowledgement target.

- [X] T005 [US1] Update `SECURITY.md` to document
  `info@shruggie.tech` as the fallback when private reporting is unavailable,
  while retaining the instruction not to open a public issue.
- [X] T006 [US1] Name repository administrators as initial triage owners in
  `SECURITY.md` and preserve the measurable one-week acknowledgement target.
- [X] T007 [US1] Validate the complete reporting-policy contract locally and
  record that hosted route/API activation remains pending authorization in
  `specs/012-github-security-baseline/quickstart.md`.

**Checkpoint**: The repository policy is internally complete and no longer
depends on an undocumented fallback.

---

## Phase 4: User Story 2 - Maintainers receive automated findings (Priority: P2)

**Goal**: Add repeatable Go CodeQL analysis and offline drift prevention without
changing the existing CI or product build.

**Independent Test**: Run the temporary fixture suite and the real repository
automation check. The valid workflow passes and every required drift class
fails with an actionable diagnostic.

- [X] T008 [US2] Add `.github/workflows/codeql.yml` with push and pull-request
  coverage for `main`, a weekly schedule, `contents: read`,
  `security-events: write`, Ubuntu hosting, and no secret consumption.
- [X] T009 [US2] Configure checkout v7, setup-go v7 from `go.mod`, CodeQL init
  v4 for Go/manual mode, `CGO_ENABLED=0 go build ./...`, and CodeQL analyze v4
  in `.github/workflows/codeql.yml`.
- [X] T010 [US2] Add the two researched CodeQL v4 references to the explicit
  allowlist in `scripts/automation-check.sh`.
- [X] T011 [US2] Implement the canonical CodeQL workflow contract checks and
  actionable diagnostics in `scripts/automation-check.sh` for triggers,
  default branch, permissions, language/manual mode, headless build, and
  required init/analyze steps.
- [X] T012 [US2] Run `sh test/scripts/automation-check_test.sh automation`,
  `sh scripts/automation-check.sh .`, and ShellCheck; make all positive cases
  pass while retaining all expected negative failures.

**Checkpoint**: The local repository has a complete, offline-certified advanced
CodeQL workflow; hosted execution and secret-setting activation remain pending.

---

## Phase 5: Decisions and Autopilot Readiness

**Purpose**: Close governance, traceability, full verification, and the local
commit before the one remote-operations halt.

- [X] T013 Add an Unreleased Changed entry for #38/#39 and dated decisions for
  `.github/workflows/codeql.yml` plus the extended offline automation policy in
  `CHANGELOG.md`.
- [X] T014 Confirm `.github/workflows/ci.yml`, product code, dependencies,
  branch protection, and release automation have no diff; record the result in
  `specs/012-github-security-baseline/quickstart.md`.
- [X] T015 Run `sh scripts/verify.sh all` in the foreground through all eight
  gates and record the result plus non-mutating `git status` evidence in
  `specs/012-github-security-baseline/quickstart.md`.
- [X] T016 Audit the complete diff for requirement/task coverage, exact issue
  traceability, UTF-8 without BOM, mojibake, and prohibited remote or product
  changes.
- [X] T017 Mark completed local tasks in this file, set the accurate lifecycle
  state in `specs/012-github-security-baseline/spec.md`, and commit the verified
  slice directly to local `main` with #38/#39 traceability.
- [X] T018 Halt exactly once before `git push origin main` and all GitHub
  security-setting mutations; present the decisions, implementation, evidence,
  exact remote plan, and residual hosted validation.

---

## Phase 6: Hosted Activation (after authorization)

**Purpose**: Publish and activate the already verified local feature.

- [x] T019 Push local `main` to `origin/main` after explicit authorization.
- [x] T020 [US1] Enable private vulnerability reporting and validate the API,
  authenticated/unauthenticated route behavior, and administrator triage access
  without submitting a fabricated report.
- [x] T021 [US2] Request enablement of secret scanning, push protection,
  non-provider patterns, and validity checks individually; read back and classify
  each result without exposing alert contents.
- [x] T022 [US2] Observe the CodeQL `main` run and available code-scanning API
  evidence; record success, token-scope limitations, and aggregate alert counts
  only.
- [x] T023 Update the feature evidence and lifecycle state with hosted results,
  run the canonical aggregate for any resulting documentation commit, and
  publish it within the already authorized activation continuation.

## Dependencies and Execution Order

- Phase 1 precedes all implementation.
- Phase 2 is test-first and blocks the CodeQL policy implementation.
- User Story 1 and User Story 2 are locally independent after Phase 2, but run
  chronologically as P1 then P2 under the single-agent autopilot.
- Phase 5 depends on both local stories and ends at the mandatory halt.
- Phase 6 is forbidden until the operator authorizes remote publication and
  security-setting mutation.

## Requirement Coverage

| Requirements | Tasks |
| --- | --- |
| FR-001, FR-002 | T019-T020 |
| FR-003, FR-004 | T005-T007 |
| FR-005, FR-006 | T021 |
| FR-007-FR-010 | T008-T012, T022 |
| FR-011 | T003-T004, T010-T012 |
| FR-012, FR-016 | T002, T014, T016 |
| FR-013 | T013 |
| FR-014 | T020-T023 |
| FR-015 | T001-T018 |
