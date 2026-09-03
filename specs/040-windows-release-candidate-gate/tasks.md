# Tasks: Windows Release Candidate Gate

**Input**: Design documents from `/specs/040-windows-release-candidate-gate/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/windows-release-evidence.md`, `quickstart.md`

**Tests**: Required by #94, #98, and constitution principle II. Validator and workflow mutation tests are written and observed failing before implementation.

## Phase 1: Specification and Baseline

- [x] T001 Complete specification, clarification record, requirements checklists, plan, research, data model, release-evidence contract, and quickstart under `specs/040-windows-release-candidate-gate/`.
- [x] T002 Run focused green baselines for repository automation, installer contracts, Windows PowerShell parsing, and the canonical gate manifest before adding S040 behavior.
- [x] T003 Run Spec Kit's read-only consistency analysis across `spec.md`, `plan.md`, and `tasks.md`; resolve every critical or high-severity finding before implementation.

---

## Phase 2: Foundational Evidence Model

**Purpose**: Establish exact identity, strict decoding, complete scenario coverage, and attachment safety before any release workflow consumes evidence.

- [x] T004 Add failing tests for strict schema decoding, candidate identity syntax, ordered timestamps, unique references, and every explicit outcome in `internal/releasegate/validate_test.go`.
- [x] T005 Add failing tests for safe bundle-relative attachment paths, regular-file enforcement, streaming byte counts, SHA-256 matching, ZIP traversal refusal, duplicates, and archive bounds in `internal/releasegate/bundle_test.go`.
- [x] T006 Implement the versioned evidence types, strict decoder, deterministic diagnostic collection, candidate hashing, safe path checks, and bounded ZIP reader in `internal/releasegate/model.go`, `validate.go`, and `bundle.go`.
- [x] T007 Add a complete inert passing fixture and named mutation cases under `test/fixtures/windows-release-gate/`; prove fixture data is labeled non-native and cannot be confused with a real candidate record.

**Checkpoint**: Malformed, ambiguous, unsafe, or altered evidence fails closed before scenario semantics are evaluated.

---

## Phase 3: User Story 1 - Prepare One Traceable Candidate (Priority: P1)

**Goal**: Bind one staged Windows MSI to repository, tag, commit, workflow, version, ProductCode, filename, size, and digest.

**Independent Test**: Validate a matching artifact and systematically mutate every identity field and artifact byte.

- [x] T008 [US1] Add failing candidate-identity and command exit-contract tests in `internal/releasegate/validate_test.go` and `scripts/windows-release-gate/main_test.go`.
- [x] T009 [US1] Implement local-directory and ZIP-bundle validation commands with explicit expected repository, tag, and commit flags in `scripts/windows-release-gate/main.go`.
- [x] T010 [US1] Add candidate-manifest creation to the Windows staging build, using compiled MSI inspection as the ProductCode/version source and retaining run ID/attempt in `.github/workflows/release.yml`.
- [x] T011 [US1] Add automation mutation tests that reject a public staging upload, a missing Windows candidate manifest, staging-time final checksums, and any release rebuild after attended evidence.

**Checkpoint**: One byte or identity-field change invalidates the candidate, and the tag workflow can only stage a draft.

---

## Phase 4: User Story 2 - Prove Installed Core Behavior (Priority: P1)

**Goal**: Validate complete attended access, native window, connection-error, and production-task evidence without weakening the proof-before-commit rule.

**Independent Test**: The passing fixture covers every scenario and every targeted mutation fails for the expected reason.

- [x] T012 [US2] Add failing access and environment tests for Windows 11 client, clean snapshot, medium-integrity intended user, LocalSystem service, unrelated-user denial, fresh PATH resolution, and post-uninstall PATH absence.
- [x] T013 [US2] Add failing native-window tests for standard/high-or-mixed DPI diversity, clean/retained profiles, rectangle completeness, restored state, margins, 1280-by-800 preferred sizing, independent 90-percent caps, state transitions, and relaunch.
- [x] T014 [US2] Add failing connection tests for all required categories, at least 120 seconds, timestamps, one in-frame incident, zero modal/additional top-level errors, reachable Retry, accurate denial guidance, and successful recovery.
- [x] T015 [US2] Add failing task tests for public interfaces, production process creation, manual/scheduled successful run IDs, output/marker/history evidence, nonzero exit, and process-start failure diagnostics.
- [x] T016 [US2] Implement access, environment, window, error, and task semantic validation in `internal/releasegate/validate.go`.
- [x] T017 [US2] Add failing tests and implement a narrowly scoped opt-in GUI evidence file for exact Fyne canvas size/scale, then implement `Initialize` and exact-PID `CaptureWindow` actions in `test/windows/Invoke-ReleaseCandidateAttended.ps1` with native HWND/process/token/session/DPI/monitor/work-area/state capture and safe workspace creation.
- [x] T018 [US2] Implement `RecordObservation` and `Finalize` actions with no-overwrite resume behavior, attachment hashing, shared-validator invocation, and ZIP creation only after a passing result.
- [x] T019 [US2] Run PowerShell parser and project compliance checks; prove console children are hidden, redirected, and noninteractive and that fixture automation never claims attended proof.

**Checkpoint**: Native evidence requirements are measurable and enforceable while product GUI and error behavior remain unchanged absent an allowed reproduction.

---

## Phase 5: User Story 3 - Prove the Attended Installer Lifecycle (Priority: P1)

**Goal**: Validate the attended S039 setup, maintenance, uninstall, multiple-profile, cancellation, and reinstall matrix required by #98 and #96.

**Independent Test**: Lifecycle fixture scenarios pass only with explicit selections, process integrity, before/after fingerprints, multiple genuine profiles, unaffected controls, and accurate partial cleanup.

- [x] T020 [US3] Add failing setup tests for shortcut defaults/matrix, completion matrix, unelevated finish launch, handler target, cancel, maintenance, and upgrade.
- [x] T021 [US3] Add failing removal tests for preserve, wipe, cancel, multiple profiles, locked partial cleanup, unaffected controls, security-state disposition, and both reinstall outcomes.
- [x] T022 [US3] Implement installer/removal semantic validation and required attachment rules in `internal/releasegate/validate.go`.
- [x] T023 [US3] Generate fail-closed setup/removal templates and operator prompts during collector initialization without automating consent or destructive choices.

**Checkpoint**: S039's hosted silent proof and S040's attended evidence remain distinguishable, and every incomplete destructive result blocks.

---

## Phase 6: User Story 4 - Block or Release the Exact Candidate (Priority: P1)

**Goal**: Make the validated evidence and exact staged MSI the only supported path from draft to public release.

**Independent Test**: Workflow fixtures reject every bypass and the promotion contract orders draft inspection, exact-asset validation, final checksum creation, and publication.

- [x] T024 [US4] Convert every tag-workflow asset upload to explicit draft staging and remove staging-time final checksum publication from `.github/workflows/release.yml`.
- [x] T025 [US4] Add failing automation fixtures for missing manual dispatch, unsafe tag input, absent draft-state/target checks, missing gate invocation, pre-validation publication, incomplete checksum ordering, and public-release mutation.
- [x] T026 [US4] Implement `.github/workflows/promote-release.yml` with strict tag validation, draft/target inspection, exact asset download, candidate/evidence validation, all-asset checksums, and final draft promotion.
- [x] T027 [US4] Extend `scripts/automation-check.sh` to enforce the staging and promotion contracts with all findings aggregated.
- [x] T028 [US4] Run the real repository automation check plus all positive/negative fixture cases and confirm existing release-note, brand, CodeQL, Dependabot, and gate-manifest contracts remain enforced.

**Checkpoint**: The supported release path cannot become public without complete evidence for the exact MSI, and no workflow rebuild occurs after testing.

---

## Phase 7: Documentation, Evidence, and Lifecycle

- [x] T029 Update `test/windows/README.md` and `docs/INSTALL-windows.md` with the exact candidate, attended matrix, data/privacy, two-minute observation, multi-profile, evidence upload, draft, promotion, and failure-recovery procedures.
- [x] T030 Update `CHANGELOG.md` with S040 behavior and dated architectural decisions for pinned release workflows and Windows operator documentation.
- [x] T031 Run all focused Go tests with race detection, PowerShell parser/compliance checks, release-gate fixture validation, workflow automation fixtures, WiX source checks, and locally available compiled-MSI checks; record exact outcomes in `verification.md`.
- [x] T032 Run all eight canonical gates with `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` and record exact outcomes in `verification.md`.
- [x] T033 Check every changed text file for UTF-8 without BOM and mojibake, resolve every task, and update `spec.md` plus `specs/README.md` to `Implemented` with an explicit #94/#98/#96 evidence disposition.
- [x] T034 Perform independent final reviews for architecture, evidence fail-closure, PowerShell process safety, workflow ordering, issue traceability, and accidental unrelated changes; remediate all confirmed findings.

## Dependencies & Execution Order

- Phase 1 completes the Spec Kit gate before production edits.
- Strict model and bundle safety in Phase 2 precede workflow or collector consumers.
- Candidate identity in Phase 3 precedes native and lifecycle observations.
- User Stories 2 and 3 share the evidence model but keep scenario-specific validation independently testable.
- Promotion work in Phase 6 depends on a stable validator and candidate manifest.
- Documentation, canonical verification, publication preparation, and review follow implementation. Branch publication and pull-request lifecycle are delivery bookkeeping recorded outside the implementation task list.

## Parallel Opportunities

- After Phase 2, semantic validator tests and PowerShell collector design may proceed independently when file ownership is separated.
- Workflow automation fixtures and operator documentation can proceed in parallel after candidate and promotion names stabilize.
- Final reviews may run from independent architecture, Windows process-safety, and release-security perspectives.

## Implementation Strategy

1. Observe the focused green baseline and complete the blocking analyze pass.
2. Make strict decoding, identity, path, archive, and scenario mutation tests fail first.
3. Implement the shared validator and command until the complete inert fixture passes.
4. Add the native attended collector without modifying product behavior.
5. Convert tag publication to draft staging and add exact-asset promotion.
6. Strengthen repository automation against workflow bypasses.
7. Document, run focused and canonical verification, review, publish the PR, and resolve CI/reviews.
