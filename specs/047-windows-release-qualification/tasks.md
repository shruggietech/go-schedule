# Tasks: Windows Release Qualification

**Input**: Design documents from `specs/047-windows-release-qualification/`

**Tests**: Required by the constitution and the release-security scope. Write failing validator and collector-contract tests before implementation.

**Organization**: User stories remain independently testable. Tasks execute chronologically under one-agent autopilot even where `[P]` records safe file independence.

## Phase 1: Setup and Baseline

- [x] T001 Confirm the clean `main` baseline, create `codex/047-windows-release-qualification`, activate `.specify/feature.json`, and record S047 as Draft in `specs/README.md`
- [x] T002 Inventory issues #96, #98, #101, #104, #105, #106, #109, #111, #112, and #113 plus S040-S046 evidence boundaries in `specs/047-windows-release-qualification/spec.md`
- [x] T003 Record the local-demo/formal-candidate, scenario grouping, metric, touchpad, issue-state, and Spec Kit prerequisite decisions in `specs/047-windows-release-qualification/research.md`

---

## Phase 2: Foundational Contract and Red Tests

- [x] T004 Add the S047 scenario identities and issue mapping to `specs/047-windows-release-qualification/contracts/desktop-release-evidence.md`
- [x] T005 [P] Add failing required-scenario, missing-observation, class-boundary, image-attachment, environment, exact-set, DPI, contrast, input, row-count, and outcome mutation tests in `internal/releasegate/validate_test.go`
- [x] T006 [P] Add failing S047 collector scenario/template/parser contract assertions in `test/integration/windows_release_gate_contract_test.go`
- [x] T007 Confirm the red tests fail only because desktop scenarios and templates are not implemented, and record the red evidence in `specs/047-windows-release-qualification/verification.md`

**Checkpoint**: The existing 36-scenario fixture is proven insufficient and all new rules have a failing test before implementation.

---

## Phase 3: User Story 1 - Enforce Complete Desktop Evidence (Priority: P1)

**Goal**: Promotion cannot accept a bundle missing or misrepresenting any S047 desktop acceptance family.

**Independent Test**: A complete 47-scenario synthetic fixture passes only the fixture entry point; each missing or corrupted desktop metric fails with the scenario and metric named.

- [x] T008 [US1] Append all fixed `desktop.` identities without weakening the existing list in `internal/releasegate/validate.go`
- [x] T009 [US1] Implement trimmed, duplicate-rejecting, order-independent exact-set metric validation in `internal/releasegate/validate.go`
- [x] T010 [US1] Require routine installed-user context and native image evidence for every desktop observation in `internal/releasegate/validate.go`
- [x] T011 [US1] Implement standard/scaled appearance metric validation for #101/#106 in `internal/releasegate/validate.go`
- [x] T012 [US1] Implement interaction-state metric and contrast-floor validation for #109 and its consumers in `internal/releasegate/validate.go`
- [x] T013 [US1] Implement navigation/Options and scroll-input validation for #104-#106/#111 in `internal/releasegate/validate.go`
- [x] T014 [US1] Implement Tasks and Schedule/Activity table validation for #112/#113 in `internal/releasegate/validate.go`
- [x] T015 [US1] Update the passing 47-scenario synthetic evidence in `internal/releasegate/validate_test.go` and `test/fixtures/windows-release-gate/passing/evidence.json`
- [x] T016 [US1] Run focused validator and CLI tests, including race coverage, and record the green result in `specs/047-windows-release-qualification/verification.md`

**Checkpoint**: The cross-platform gate machine-enforces every desktop evidence family while rejecting its synthetic fixture through the production path.

---

## Phase 4: User Story 2 - Record Native Results Safely (Priority: P1)

**Goal**: The attended collector creates one complete, reviewable template for every new observation and preserves its overwrite protections.

**Independent Test**: Initialize against an inert MSI, inspect all 47 unique placeholders and 27 setup/remove/desktop templates, and prove malformed or incomplete imports remain fail-closed.

- [x] T017 [US2] Generalize scenario-template creation without changing lifecycle fields in `test/windows/Invoke-ReleaseCandidateAttended.ps1`
- [x] T018 [US2] Add complete desktop metric templates and screenshot paths in `test/windows/Invoke-ReleaseCandidateAttended.ps1`
- [x] T019 [US2] Append the eleven desktop observations and emit their templates during Initialize in `test/windows/Invoke-ReleaseCandidateAttended.ps1`
- [x] T020 [US2] Pass the PowerShell parser and Windows integration contract tests and record results in `specs/047-windows-release-qualification/verification.md`

**Checkpoint**: Operator workspaces cannot silently omit or overwrite the new desktop evidence.

---

## Phase 5: User Story 4 - Close Work by Evidence, Not Association (Priority: P2)

**Goal**: Reviewers can map each open issue to formal observations and keep incomplete issues open independently.

**Independent Test**: Audit the documented matrix against every issue's current acceptance criteria and the validator's canonical scenario list.

- [x] T021 [P] [US4] Replace overlapping S042/S044/S045 prose with the canonical S047 scenario matrix and issue links in `test/windows/README.md`
- [x] T022 [P] [US4] Document exact-candidate collection, local-demo limits, and individual issue closure in `specs/047-windows-release-qualification/quickstart.md`
- [x] T023 [US4] Record #94's closed/current #96 checklist inconsistency without reopening or rewriting history in `specs/047-windows-release-qualification/verification.md`
- [x] T024 [US4] Update `[Unreleased]` and add a dated promotion-gate decision in `CHANGELOG.md`

**Checkpoint**: The enforced matrix, operator runbook, and GitHub acceptance surface use the same scenario vocabulary.

---

## Phase 6: Automated Qualification and Local Source Identity

- [x] T025 Run `/speckit-analyze`, resolve every finding, and record 100 percent requirement/task coverage with zero critical/high findings in `specs/047-windows-release-qualification/verification.md`
- [x] T026 Advance S047 to In Progress in `specs/047-windows-release-qualification/spec.md` and `specs/README.md` before verification
- [x] T027 [P] Run the focused release-gate, PowerShell parser, Windows contract, WiX source, and full Go suites from `specs/047-windows-release-qualification/quickstart.md`
- [x] T028 Run the foreground canonical `scripts/verify.sh all` aggregate and record all eight results and coverage in `specs/047-windows-release-qualification/verification.md`
- [x] T029 Run strict UTF-8 without BOM, mojibake, unresolved-placeholder, task-format, and `git diff --check` audits and record results in `specs/047-windows-release-qualification/verification.md`
- [x] T030 Mark every completed task, preserve attended tasks as actionable, and create the local source commit `feat(047): enforce Windows release qualification` with the required co-author trailer

---

## Phase 7: User Story 3 - Produce and Inspect the Pre-Push Demo (Priority: P1)

**Goal**: Deliver one non-release S047 MSI tied to the committed source for the last native pre-push walkthrough.

**Independent Test**: Artifact, compiled inspection, and handoff agree on source commit, embedded version, ProductVersion, ProductCode, byte size, SHA-256, and `local-demo` class.

- [x] T031 [US3] Build the four Windows amd64 binaries from the committed source with `1.0.0-s047-demo.<short-commit>` embedded and stage only intended MSI inputs under untracked `dist/`
- [x] T032 [US3] Compile `dist/go-schedule_s047-demo_<short-commit>_windows_amd64.msi` through the pinned WiX 6.0.2 path
- [x] T033 [US3] Run compiled inspection with artifact class `local-demo`, verify embedded versions and hash, and record identity in `specs/047-windows-release-qualification/verification.md`
- [x] T034 [US3] Re-run documentation, lifecycle, diff, encoding, and artifact-identity checks after recording the demo and create an evidence-only local commit

**Checkpoint**: The exact linked local demo is ready for attended testing and no repository or release publication has occurred.

---

## Phase 8: Attended Pre-Push Qualification (Operator Required)

- [x] T035 [US3] Have the maintainer walk through the exact linked demo while the assistant verifies artifact/environment identity and records the observed appearance, interaction, navigation/Options, scroll, table, installer, and uninstall results in `specs/047-windows-release-qualification/checklists/attended-demo.md`
- [x] T036 [US3] Record all attended pass/fail/unavailable results and environment details in `specs/047-windows-release-qualification/verification.md` before any correction
- [x] T037 [US3] Resolve every reproducible S047 failure test-first, rebuild to a new demo identity, and repeat affected observations, or record that no correction was required
- [x] T038 [US4] Reconcile local-demo issue dispositions without claiming formal exact-candidate closure, mark S047 Implemented when its pre-push gate is complete, rerun all eight gates, and create the final local commit

---

## Phase 9: Pull-Request Review Corrections

- [x] T039 Add failing tests proving every visual family requires distinct standard/scaled evidence and declared image metadata cannot disguise non-raster bytes
- [x] T040 Add four scaled-DPI scenario identities and bind standard/scaled desktop observations to exactly 96 or greater than 96 effective DPI
- [x] T041 Validate screenshot content from supported raster signatures independently of declared media type and filename extension
- [x] T042 Expand the checked fixture and attended collector to all 47 observations and 27 generated lifecycle/desktop templates
- [x] T043 Align the S047 contract, plan, data model, runbook, quickstart, changelog, and verification record with the corrected gate
- [x] T044 Run the focused and canonical verification gates after review corrections

---

## Dependencies and Execution Order

```text
Phase 1 -> Phase 2 -> US1 gate -> US2 collector -> US4 traceability
        -> automated qualification -> committed source -> US3 demo
        -> attended pre-push qualification -> pull-request review corrections
```

- US1 supplies the canonical scenario contract consumed by US2 and US4.
- The source commit must precede demo compilation so its embedded version is immutable and auditable.
- Phase 8 requires physical Windows interaction and maintainer judgment. It is the only intended external-input boundary before publication.
- Merge, tag staging, formal 47-scenario qualification, issue closure, and release promotion are publication/release workflow, not unchecked implementation tasks.

## Parallel Opportunities

- T005 and T006 touch independent Go and integration-test files after T004.
- T021 and T022 can proceed independently once the scenario vocabulary is final.
- T027's focused commands are conceptually independent but run sequentially on this host to avoid noisy process fan-out.

## Implementation Strategy

The local MVP is US1: the release gate refuses incomplete desktop evidence. US2 makes that contract operable, US4 makes it auditable, and US3 supplies the pre-push native-risk check. Stop after T034 until the maintainer completes T035; do not weaken or mark Phase 8 complete from headless evidence.
