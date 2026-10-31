# Tasks: SQLite convergence and release delivery

**Input**: Design documents from specs/130-sqlite-release-convergence/.

**Prerequisites**: spec.md, plan.md, research.md, data-model.md, contracts/storage-release.md and quickstart.md. Tests are required by the constitution and this plan.

## Phase 1: Setup

- [x] T001 Establish exact S130 identity and frozen 221-input baseline in .specify/feature.json and specs/130-sqlite-release-convergence/spec.md.
- [x] T002 Complete published candidate/license and storage/release research in specs/130-sqlite-release-convergence/research.md.

## Phase 2: Foundational

- [x] T003 Resolve the blocking cross-artifact analysis and preserve its result in specs/130-sqlite-release-convergence/analysis.md before implementation.
- [x] T004 Add real-driver cancellation/prepared correctness and benchmark fixtures in internal/store/sqlite_driver_test.go; retain old-pair results before changing dependencies.

## Phase 3: User Story 1 - Coherent candidate (Priority: P1)

**Goal**: Both application modules resolve one complete corrected graph without metadata writes.

**Independent Test**: Normal -mod=readonly resolution, affected desktop connection tests and unchanged hashes of all four metadata files.

- [x] T005 [US1] Select exact published SQLite1.61.1/libc1.79.0 and tidy root go.mod/go.sum without unrelated upgrades; retain expected stale-desktop read-only failure.
- [x] T006 [US1] Regenerate desktop/go.mod and desktop/go.sum, preserving /v2 namespace, required root v2.0.0 and local replacement.
- [x] T007 [US1] Resolve both graphs read-only, verify checksums and affected desktop/remote tests, recording zero metadata drift in specs/130-sqlite-release-convergence/verification.md.

## Phase 4: User Story 2 - Current persistence and cancellation (Priority: P1)

**Goal**: Current storage behavior remains and the real lock-cancellation regression is corrected by dependency selection.

**Independent Test**: Focused driver tests, identical before/after prepared benchmark and existing current-state/audit/security/race suites.

- [x] T008 [US2] Demonstrate selected-driver cancellation and prepared-query green in internal/store/sqlite_driver_test.go, preserving first red without weakening its assertions.
- [x] T009 [US2] Compare repeated identical prepared benchmark samples before/after and record the 10 percent policy disposition in specs/130-sqlite-release-convergence/verification.md.
- [x] T010 [US2] Correct misleading old migration descriptions and weak marker assertion in internal/store/store_test.go to state the current schema contract exactly.

## Phase 5: User Story 3 - Existing release continuation (Priority: P2)

**Goal**: Current candidate facts and existing release continuation agree with the adopted pair while public delivery remains honestly pending.

**Independent Test**: Existing source/comparison guards, candidate content bindings, unchanged historical records and public selection.

- [x] T011 [US3] Record the chronological S130 dependency decision in CHANGELOG.md and the prepared .github/release-notes/v2.0.0.md without changing historical releases.
- [x] T012 [US3] Refresh only candidate changelog/notes bindings in site/data/reviews/v2.0.0.json and preserve comparison dataset and historical bytes.
- [x] T013 [US3] Update the factual dependency disposition in specs/129-prelaunch-release/handoff.md and document exact existing release continuation plus two-request review cap in specs/130-sqlite-release-convergence/delivery.md; reconcile #404/#412 native planning records.

## Phase 6: Polish and Cross-Cutting Verification

- [x] T014 Run all eight foreground gates through scripts/verify.sh all, retaining complete terminal/source/input evidence in specs/130-sqlite-release-convergence/verification.md.
- [x] T015 Correct the stale current dependency guide in desktop/README.md and rerun scripts/docs-check.sh plus final source/format/lifecycle checks; preserve the earlier aggregate source boundary.
- [x] T016 Complete source/comparison, dual-module read-only, current security/refusal/audit and history/input integrity acceptance in specs/130-sqlite-release-convergence/verification.md.
- [x] T017 Finalize factual delivery/lifecycle records in specs/130-sqlite-release-convergence/spec.md, specs/130-sqlite-release-convergence/tasks.md and specs/README.md, then pass GitHub formatting, lifecycle and UTF-8/no-BOM checks.

## Dependencies and Execution Order

T001-T003 precede implementation. T004 baseline precedes T005-T006. T007 and T008 require the combined graph; T009 compares T004 and selected-pair measurements. T010 is a proportional current-test correction. T011-T013 follow selected-pair acceptance; T014 freezes complete implementation for the full aggregate. T015 corrects only the current build guide after the quiescent aggregate. T016-T017 verify the final documentation/metadata and finalize factual records without changing evaluated code/dependencies. User stories have independent acceptance evidence but converge into one source candidate.

## Parallel Opportunities

Upstream/storage/release research is read-only and was dispatched in parallel by /speckit-plan. After dependency acceptance, candidate documentation and private evidence analysis concern separate files but are performed sequentially here to avoid source drift. Verification runs with source quiescent. No extra implementation agents are needed.

## Implementation Strategy

Deliver the complete root/desktop correction as the minimum useful increment, preserve current store behavior with specific regressions and benchmarks, then bind the existing release continuation to that source. Push/PR/review/owner merge and later actual release operations are publication evidence, not unchecked implementation tasks. #404 can close on functional integration; #412 remains open for actual public delivery.
