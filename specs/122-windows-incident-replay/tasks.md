# Tasks: S122 Windows incident replay and causal triage

## Phase 1: Setup

- [x] T001 Verify accepted main, exact S122 branch/selector and preserve original artifact hashes in .test-temp/s122/evidence-before.json (FR-001, FR-009).
- [x] T002 Complete installed Spec-Kit specify, clarify, checklist and plan artifacts in specs/122-windows-incident-replay/ (FR-001 through FR-010).

## Phase 2: Discriminating development regressions

**Goal**: Establish separately attributable first failures without claiming original causality.

- [x] T003 [US1] Add and run the existing-but-unlaunchable interpreter regression in test/integration/testscripts_incident_test.go, preserving native exit/error and collection independence (FR-004, FR-008).
- [x] T004 [US2] Add and run canceled/queued lifecycle mutation regressions in internal/mcphttp/manager_cancellation_test.go, retaining listener/credential and access semantics (FR-003, FR-008).

## Phase 3: Replay preparation

**Goal**: Reuse actual owned contracts with direct observation attribution.

**Independent test**: All original assertions remain active and each selected observation names its condition and incident.

- [x] T005 [US1] Add opt-in per-writer/contrast incident correlation in test/integration/testscripts_binding_test.go and testscripts_contrasts_test.go (FR-002, FR-004, FR-007).
- [x] T006 [US2] Add the opt-in unchanged catalog replay and owned memory/WAL plus packaged-daemon fixture cleanup in internal/clievaluate/incident_replay_test.go (FR-003, FR-004, FR-006, FR-009).
- [x] T007 [US3] Declare the finite condition matrix, shared campaign deadline, prepared source/binary identities and fixed report in specs/122-windows-incident-replay/investigation.md (FR-001, FR-005, FR-007).

## Phase 4: Bounded investigation

**Goal**: Assess original symptoms with the delivered collectors before causal conclusions.

**Independent test**: Every planned condition has an observed or explicit incomplete disposition with source identity and cleanup.

- [x] T008 [US1] Run retained/current PowerShell and Bash six-writer/selection contrasts, plus prerequisite fault cases once; retain exact outcomes and artifacts under .test-temp/s122/ (FR-002, FR-004, FR-005).
- [x] T009 [US2] Run all five original scenario shapes in memory/WAL component fixtures and packaged-daemon WAL, respecting collection count and shared deadline; retain private observations under .test-temp/s122/ (FR-003, FR-005, FR-006).
- [x] T010 [US3] Analyze actual observations, file separately demonstrated defects after duplicate search, and record causal or unresolved dispositions in specs/122-windows-incident-replay/investigation.md and GitHub #355/#356 (FR-001, FR-007, FR-008, FR-010).

## Phase 5: Supported corrections

**Goal**: Deliver narrowly demonstrated corrections without rewriting original failures.

- [x] T011 [US1] Correct demonstrated launch incident classification in test/integration/testscripts_observer_test.go and rerun the failing regression, or record evidence disproving that candidate (FR-004, FR-008).
- [x] T012 [US2] Correct demonstrated lifecycle cancellation in internal/mcphttp/manager.go with context-aware serialization and unchanged committed cleanup, or record evidence disproving that candidate (FR-003, FR-008).
- [x] T013 [US3] Run targeted Windows/native Linux race, cancellation, ownership, credential and collector checks; record exact evidence in specs/122-windows-incident-replay/verification.md (FR-004, FR-006, FR-008, FR-009).
- [x] T014 [US3] After application changes run and fully review/finalize one fresh source-bound Windows/WSL catalog using existing report schemas; retain original failed/campaign inputs unchanged (FR-007, FR-008, FR-010).

## Phase 6: Acceptance and delivery

- [x] T015 [US3] Publish the exact bounded replay recipe and fixed six-section investigation report outline in docs/maintainers/windows-incident-replay.md; link from docs/cli-qa.md and docs/test-scripts.md (FR-002, FR-003, FR-005, FR-007).
- [x] T016 Reconcile actual S121 merge/main CI/Pages in specs/121-windows-incident-capture/ and specs/README.md; update CHANGELOG.md and current S122 delivery records without release identity changes (FR-010).
- [x] T017 Run complete foreground sh scripts/verify.sh all with frozen before/after source and original evidence manifests; record results in specs/122-windows-incident-replay/verification.md (FR-001, FR-008, FR-009).
- [x] T018 Finalize implemented lifecycle/tasks and factual delivery in specs/122-windows-incident-replay/, validate final formatting/lifecycle records and commit on codex/122-windows-incident-replay (FR-007, FR-010).

## Dependencies and execution

T001/T002 precede all implementation. T003/T004 capture first failures before T011/T012. T005/T006 precede T007, then T008/T009 run sequentially under the common deadline. T010 analyzes/follows up before corrections. T013 follows all actual fixes; T014 follows application changes. T015/T016 precede frozen aggregate T017 and factual finalization T018. Each newly demonstrated defect has its own GitHub record. Original #355/#356 remain open if causal outcomes are incomplete.

The two read-only research agents required by the installed plan command completed. No additional implementation delegation is needed. US1 and US2 design could be read in parallel, but owned campaigns and overlapping edits execute sequentially. MVP is attributable runtime replay; the authorized slice delivers all three stories and the complete IPC/store cross-product. Publication/review are the owner-authorized workflow after implementation, not unchecked code tasks. Push/open one official PR, handle at most two Codex rounds and stop before merging/tagging/promotion.
