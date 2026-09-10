# Tasks: Cumulative v1.4.0 Release Preparation

**Input**: Design documents from `specs/081-v140-release-preparation/`

**Tests**: Release identity, copy, workflow preflight, upgrade safety, repository gates, candidate integrity, and public checksums are mandatory.

## Phase 1: Setup and specification

- [x] T001 Verify synchronized main, public v1.1.1 state, absent v1.4.0 tag and release state, completed v1.2 through v1.4 milestones, and the intended `shruggietech/go-schedule` remote.
- [x] T002 Create issue #226, reopen milestone v1.4.0, add the project item, create branch `codex/081-v140-release-preparation`, and initialize the Spec Kit artifact set.
- [x] T003 Complete specification, autonomous clarifications, requirements and release checklists, research, data model, plan, publication contract, and quickstart in `specs/081-v140-release-preparation/`.
- [x] T004 Analyze requirement coverage, constitution compliance, chronological ordering, release history, and authorization boundaries before implementation.

## Phase 2: Fail-closed release metadata (User Story 2)

**Goal**: Reject incomplete or inconsistent v1.4.0 source metadata before a release can mutate artifacts.

**Independent Test**: Mutate each required source-owned identity in the approved workflow fixture and require the automation validator to reject it.

- [x] T005 [US2] Add failing release-workflow fixture mutations for missing changelog and release-note metadata in `test/scripts/automation-check_test.sh`.
- [x] T006 [US2] Extend the pre-artifact release metadata gate for the README, changelog, tag-specific note, and exact commit in `.github/workflows/release.yml` and `scripts/automation-check.sh`.
- [x] T007 [US2] Generalize `test/windows/Invoke-ReleaseCandidateAttended.ps1` to version-two native window attachments and retain historical version-one validation in `internal/releasegate/validate.go` and `internal/releasegate/validate_test.go` after the current-desktop contract test fails.

## Phase 3: Cumulative public identity (User Story 1)

**Goal**: Prepare concise and historically honest v1.4.0 public source surfaces.

**Independent Test**: Count release highlights, resolve the tagged changelog link, compare the dated version boundary to v1.1.1, and scan all publication-facing prose for stale candidate or latest-release claims.

- [x] T008 [US1] Cut the complete accumulated Unreleased record into a dated v1.4.0 boundary and retain an empty Unreleased section in `CHANGELOG.md`.
- [x] T009 [US1] Add exactly four concise highlights and one final tagged changelog link in `.github/release-notes/v1.4.0.md`.
- [x] T010 [US1] Update the health example and durable cumulative-release language in `README.md`, `docs/install.md`, `docs/INSTALL-windows.md`, `docs/INSTALL-macos.md`, `docs/INSTALL-linux.md`, `docs/api.md`, and `docs/architecture.md`.

## Phase 4: Exact candidate qualification contract (User Story 3)

**Goal**: Make the later fresh-install, v1.1.1-upgrade, candidate, evidence, promotion, and audit sequence executable without pretending it has run.

**Independent Test**: Follow the contract from reviewed merge through final audit and confirm every transition has an objective input, stop condition, and safe failure state.

- [x] T011 [US3] Replace stale Fyne and v1.0.0 examples with candidate-bound Wails v1.4.0 fresh and v1.1.1-upgrade instructions in `test/windows/README.md`.
- [x] T012 [US3] Reconcile release issue #226, project Slice and Status fields, and parent issue #146 with the S081 preparation and post-merge authorization boundary.

## Phase 5: Verification and evidence

- [x] T013 Run focused automation fixture tests, release-copy checks, GitHub formatting, specification lifecycle, and diff integrity.
- [x] T014 Run all eight canonical verification gates in the foreground with `C:\Program Files\Git\bin\sh.exe scripts/verify.sh all`.
- [x] T015 Record exact results in `specs/081-v140-release-preparation/verification.md`, mark required tasks complete, and transition S081 to Implemented in `specs/081-v140-release-preparation/spec.md` and `specs/README.md`.

## Phase 6: Publication operations

Publication operations remain actionable under issue #226 rather than being represented as completed repository implementation tasks. After the preparation pull request is reviewed, green, and merged, they proceed chronologically: exact-commit main CI, explicit tag authorization, immutable v1.4.0 tag, draft staging, exact-candidate fresh and v1.1.1-upgrade qualification, evidence upload, explicit promotion authorization, no-rebuild promotion, final checksum and identity audit, then issue, project, and milestone closure.

## Dependencies and execution order

T001 through T004 establish release authority and an internally consistent design. T005 must fail before T006 implements the preflight. T007 removes the retired-toolkit blocker before current candidate instructions rely on the collector. T008 through T010 establish source identity after the validator contract is known. T011 and T012 make the later physical and administrative lifecycle executable. T013 and T014 validate the complete branch before T015 records implementation evidence.

## Parallel opportunities

After T006, T007 and T008 affect separate files and may be prepared together. T009 and T010 are documentation changes with overlapping release terminology and therefore remain sequential. Publication operations are deliberately excluded from parallel execution because every step depends on the immutable output of the previous step.

## Implementation strategy

The minimum useful increment is the fail-closed metadata gate plus one consistent cumulative source boundary. Documentation and qualification contracts then make that boundary operable. No product behavior, dependency, tag, draft release, installed candidate, or public artifact is introduced by S081.
