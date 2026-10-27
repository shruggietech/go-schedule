# Tasks: S121 Windows incident capture and causal resolution

## Phase 1: Setup

- [x] T001 Freeze accepted base and original artifact hashes in .test-temp/s121/evidence-before.json.
- [x] T002 Complete installed Spec-Kit specification, clarification, checklists and research in specs/121-windows-incident-capture/.

## Phase 2: Shared incident foundation

- [x] T003 Add failing bounded artifact, exclusive-path, identity and format regressions in internal/incident/incident_test.go (FR-001, FR-002, FR-006, FR-007).
- [x] T004 Implement fixed incident metadata and exclusive bounded artifact collection in internal/incident/incident.go (FR-002, FR-003, FR-006).

## Phase 3: US1 Owned runtime failure

**Goal**: Opt-in child-only crash capture with the original result unchanged.

**Independent test**: Native controlled crash and ordinary nonzero exit have distinct incident outcomes and preserve their exact native exit.

- [x] T005 [US1] Add native failure, unavailable/oversize, concurrency and disabled-mode regressions in test/integration/testscripts_incident_test.go (FR-001, FR-003, FR-004, FR-007, FR-008).
- [x] T006 [US1] Integrate CoreCLR child-only capture and same-process pre-payload identity in test/integration/testscripts_observer_test.go and testscripts_binding_test.go (FR-003, FR-004).
- [x] T007 [US1] Run retained/current runtime and six-writer/SQLite contrasts, recording actual dispositions in specs/121-windows-incident-capture/investigation.md (FR-008, FR-009).

## Phase 4: US2 Owned daemon stall

**Goal**: Live bounded evidence independent of the blocked API/storage path.

**Independent test**: A real held-pool IPC request is visible while waiting, independently of later cancellation and cleanup.

- [x] T008 [US2] Add real IPC held-pool/handler, wrong-target, canceled/absent responder and bounded tracker regressions in internal/incident/daemon_test.go and internal/store/request_cancellation_test.go (FR-005, FR-006, FR-007, FR-008).
- [x] T009 [US2] Implement bounded private responder and optional authored-route/stage tracking in internal/incident/, internal/api/server/, internal/store/store.go and cmd/goschedd/main.go (FR-005, FR-006).
- [x] T010 [US2] Add capture/resume identity, live watchdog and result-preservation regressions in internal/clievaluate/incident_test.go (FR-001, FR-005, FR-007, FR-011).
- [x] T011 [US2] Propagate --incidents through scripts/cli-qa.ps1, cmd/gosched-qa/main.go and internal/clievaluate coordinator/worker/session, retaining strict report and original budgets (FR-005, FR-007, FR-010, FR-011).

## Phase 5: US3 Supported dispositions

**Goal**: Actionable maintainer instructions and separately attributable follow-up outcomes.

**Independent test**: Fixed incident fields and existing QA reports distinguish unavailable collection from the original command/inspection failure.

- [x] T012 [US3] Document collection, local artifact interpretation and original causal limits in docs/cli-qa.md, docs/test-scripts.md and specs/121-windows-incident-capture/investigation.md (FR-009, FR-010).
- [x] T013 [US3] Reconcile merged S120 delivery in specs/120-windows-causal-diagnostics/ and specs/README.md without altering release identity (FR-009, FR-011).
- [x] T014 [US3] Run one fresh source-bound Windows/WSL capture and message review using existing report schemas, recording any demonstrated defects separately in GitHub and specs/121-windows-incident-capture/verification.md (FR-008, FR-009, FR-011).

## Phase 6: Acceptance

- [x] T015 Run targeted native/WSL race checks and canonical foreground eight-gate aggregate, recording source manifests and original evidence hashes in specs/121-windows-incident-capture/verification.md.
- [x] T016 Finalize CHANGELOG.md, implemented lifecycle/delivery records and GitHub formatting, then recheck final lifecycle/automation metadata in specs/121-windows-incident-capture/.

## Dependencies and execution

T001/T002 precede foundation. T003 precedes T004. US1 T005 precedes T006/T007; US2 T008 precedes T009 and T010 precedes T011. US3 follows both integrations; T014 follows source changes. T015/T016 complete implementation before authorized publication. Branch, selector, feature directory and issue mapping remain S121 throughout. Publishing and review are workflow evidence, not unchecked implementation tasks.

US1 and US2 test design can be read in parallel after foundation, but overlapping shared-file edits execute sequentially. MVP is the shared incident contract plus US1; the authorized bundle delivers all three stories. No additional implementation agents are required.
