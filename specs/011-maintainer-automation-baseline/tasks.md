# Tasks: Maintainer Automation Baseline

**Input**: Design documents from
`specs/011-maintainer-automation-baseline/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`,
`contracts/`, `quickstart.md`

**Tests**: Required by constitution Principle II. Contract regressions are
written and observed failing before their implementations.

**Organization**: Tasks are grouped by user story. User Story 1 modernizes the
hosted action runtime; User Story 2 establishes the canonical local definition
of green. Shared drift protection is foundational because it validates both.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and has no
  dependency on incomplete work.
- **[Story]**: Maps the task to US1 or US2 from `spec.md`.
- Every task names the exact file or files it changes or validates.

## Phase 1: Setup and Baseline Evidence

**Purpose**: Establish the exact pre-change contract and feature context before
touching pinned artifacts.

- [X] T001 Confirm local `main`, a clean pre-feature baseline, and active feature
  `specs/011-maintainer-automation-baseline` in `.specify/feature.json`; record
  any pre-existing working-tree changes without modifying them.
- [X] T002 Record the current action-reference inventory, workflow triggers,
  permissions, matrices, artifact names, release globs, and Make verification
  commands in the validation notes of
  `specs/011-maintainer-automation-baseline/quickstart.md` so preservation can be
  checked after the pinned-file edits.

---

## Phase 2: Foundational Automation Contract

**Purpose**: Build the independent policy test seam that both user stories need.

**CRITICAL**: The regression harness must fail for the expected missing scripts
before either implementation is added.

- [X] T003 Create `test/scripts/automation-check_test.sh` with temporary fixture
  repositories covering the approved action set, obsolete
  `actions/checkout@v4`, an unknown action, complete/missing/duplicate/extra gate
  manifests, unknown verify mode, and controlled aggregate-child failure; run it
  and record the expected FAIL because `scripts/automation-check.sh` and
  `scripts/verify.sh` do not exist.
- [X] T004 Implement the offline approved-action allowlist, independent eight-gate
  manifest comparison, optional fixture root, and actionable failures in
  `scripts/automation-check.sh`; run the action/manifest fixture subset from
  `test/scripts/automation-check_test.sh` and make it pass without network or
  workspace mutation.

**Checkpoint**: Workflow action choices and the future driver manifest now have
an independent, fixture-testable policy.

---

## Phase 3: User Story 1 - Hosted automation starts on a supported runtime
(Priority: P1)

**Goal**: Replace every Node 20 action major in CI and release automation while
preserving the existing workflow contract.

**Independent Test**: Run the action-policy fixture and the real-repository
action audit, then diff workflow triggers, permissions, matrices, inputs,
outputs, artifact names, and release globs against the T002 baseline. All action
references are approved and no non-runtime contract changed.

### Tests for User Story 1

- [X] T005 [US1] Run the action-policy fixture subset in
  `test/scripts/automation-check_test.sh` against the unmodified real workflows
  and record the expected FAIL naming the four Node 20-era majors before editing
  either workflow.

### Implementation for User Story 1

- [X] T006 [US1] Update all `actions/checkout`, `actions/setup-go`, and
  `actions/upload-artifact` references to their researched Node 24 majors in
  `.github/workflows/ci.yml` without changing triggers, permissions, job
  boundaries, matrices, environment, inputs, or artifact names.
- [X] T007 [P] [US1] Update all `actions/checkout`, `actions/setup-go`, and
  `softprops/action-gh-release` references to their researched Node 24 majors in
  `.github/workflows/release.yml` without changing tag triggers, permissions,
  inputs, output consumption, release names, or file globs.
- [X] T008 [US1] Run `scripts/automation-check.sh` against the repository and
  compare `.github/workflows/ci.yml` and `.github/workflows/release.yml` with the
  T002 baseline; record the static pass and preserved contract in
  `specs/011-maintainer-automation-baseline/quickstart.md` without running a tag
  or release.

**Checkpoint**: User Story 1 is complete and independently demonstrable: all
workflow action references are approved Node 24 majors and behavior is preserved.

---

## Phase 4: User Story 2 - Maintainers have one definition of green
(Priority: P2)

**Goal**: Provide one non-mutating local driver used by Make, CI, contributor
guidance, and autopilot for the complete eight-gate contract.

**Independent Test**: Run the driver manifest, named docs gate, controlled child
failure fixture, aggregate command, optional Make wrapper, and post-run
`git status`. The manifest is exact, failures propagate, all gates pass on the
valid repository, and verification adds no working-tree change.

### Tests for User Story 2

- [X] T009 [US2] Run the verify-mode, manifest, and failure-propagation cases in
  `test/scripts/automation-check_test.sh` and record the expected FAIL before
  implementing `scripts/verify.sh`.

### Implementation for User Story 2

- [X] T010 [US2] Implement `list`, `all`, `format`, `vet`, `lint`, `race`, `gui`,
  `coverage`, `docs`, and `automation` modes, foreground fail-closed dispatch,
  stable gate headings, usage errors, and `go.mod`-derived linter toolchain
  selection in `scripts/verify.sh` according to
  `contracts/verify-command.md`.
- [X] T011 [US2] Make `verify`, `fmt-check`, `vet`, `lint`, `test-race`,
  `test-gui`, `cover`, and `docs-check` thin delegates to
  `scripts/verify.sh` in `Makefile`; retain `fmt` as an explicitly documented
  mutating convenience target and remove the drifted command copies.
- [X] T012 [US2] Replace inline format, vet, lint, race, GUI-test, coverage, and
  docs-check commands with the corresponding `scripts/verify.sh` named modes in
  `.github/workflows/ci.yml`, and add the automation-contract mode to an ordinary
  gating job without changing job boundaries or platform coverage.
- [X] T013 [P] [US2] Replace the multi-command verification instructions with the
  canonical aggregate invocation, retain the per-gate explanation and Windows/C
  toolchain prerequisites, and distinguish mutating Make targets in
  `CONTRIBUTING.md`.
- [X] T014 [P] [US2] Update the build-phase and verification sections plus the
  Spec-Kit active-plan summary to name the canonical aggregate invocation and
  eight-gate contract in `CLAUDE.md`.
- [X] T015 [P] [US2] Update the CI-parity section and pre-push evidence breakdown
  to name the canonical aggregate invocation, eight required gates, and honest
  prerequisite failure behavior in `docs/build-autopilot.md`.
- [X] T016 [US2] Run every User Story 2 scenario from
  `specs/011-maintainer-automation-baseline/quickstart.md`; record manifest,
  named-gate, failure-propagation, aggregate, optional Make-wrapper, and
  cleanliness evidence in that file.

**Checkpoint**: User Story 2 is complete and independently demonstrable: one
driver defines green, all consumers delegate, and a successful run is
non-mutating.

---

## Phase 5: Polish, Decisions, and Autopilot Readiness

**Purpose**: Close cross-story traceability, pinned-artifact governance, and the
blocking verification gates.

- [X] T017 Add an Unreleased Changed entry closing #21 and #41 plus dated pinned
  decisions for `.github/workflows/ci.yml`,
  `.github/workflows/release.yml`, and `Makefile` in `CHANGELOG.md`; explain the
  Node 24 majors, single-source driver, offline allowlist, and preserved workflow
  contract.
- [X] T018 Run `sh test/scripts/automation-check_test.sh` and
  `sh scripts/automation-check.sh .` in the foreground; confirm every positive
  fixture passes, every negative fixture is caught, and no network or repository
  mutation occurs.
- [X] T019 Run `sh scripts/verify.sh all` in the foreground through all eight
  gates, then run `git status --short` and confirm verification introduced no
  additional changes; if the C toolchain or shell prerequisite is unavailable,
  halt and report the gate as unrun rather than green.
- [X] T020 Audit the complete diff for issue traceability, workflow-contract
  preservation, task completion, UTF-8 without BOM, and mojibake across
  `specs/011-maintainer-automation-baseline/`, `.github/workflows/`, `scripts/`,
  `test/scripts/`, `Makefile`, `CONTRIBUTING.md`, `CLAUDE.md`,
  `docs/build-autopilot.md`, and `CHANGELOG.md`.
- [X] T021 Mark completed tasks in
  `specs/011-maintainer-automation-baseline/tasks.md`, set the accurate feature
  lifecycle state in `specs/011-maintainer-automation-baseline/spec.md`, commit
  locally as `feat(011): establish the maintainer automation baseline` with the
  required co-author trailer, and halt before `git push origin main`.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: No dependencies.
- **Phase 2**: Depends on the baseline inventory from Phase 1 and blocks both
  user stories.
- **User Story 1 (Phase 3)**: Depends on the action-policy foundation. It does
  not require the verification-driver implementation.
- **User Story 2 (Phase 4)**: Depends on the manifest-policy foundation. Its CI
  edit follows the User Story 1 edit because both touch `ci.yml`; the stories are
  functionally independent but file sequencing avoids a conflict.
- **Phase 5**: Depends on both selected user stories.

### User Story Dependencies

- **US1**: Independently complete after T005-T008.
- **US2**: Independently complete after T009-T016. T012 is sequenced after US1
  only because both edit `.github/workflows/ci.yml`, not because US2 behavior
  depends on the action-major upgrade.

### Parallel Opportunities

- T007 can run in parallel with T006 because the workflows are different files.
- T013, T014, and T015 can run in parallel after the driver and Make contracts
  stabilize because they edit separate guidance files.
- Fixture authoring in T003 may prepare all negative cases together, but their
  implementations remain ordered so red-before-green evidence is preserved.

## Parallel Example: User Story 1

```text
Task: T006 update .github/workflows/ci.yml action majors
Task: T007 update .github/workflows/release.yml action majors
```

## Parallel Example: User Story 2

```text
Task: T013 update CONTRIBUTING.md
Task: T014 update CLAUDE.md
Task: T015 update docs/build-autopilot.md
```

## Implementation Strategy

### MVP First: User Story 1

1. Complete baseline and policy foundation (T001-T004).
2. Observe the old-action regression fail (T005).
3. Upgrade both workflows (T006-T007).
4. Prove the real action inventory and workflow contract (T008).
5. Stop and validate the independently useful CI-survival increment.

### Incremental Delivery

1. Foundation makes action and gate drift testable.
2. US1 removes the immediate hosted-runtime deadline.
3. US2 makes future local/CI verification changes converge on one source.
4. Final phase records pinned decisions, runs the complete gate, commits locally,
   and performs the protocol's only halt.

## Format Validation

- 21 tasks use the required checkbox and sequential `T###` identifiers.
- Every user-story task carries `[US1]` or `[US2]`.
- Every task names an exact file or an exact validation command plus evidence
  destination.
- `[P]` appears only where files and dependencies are independent.
