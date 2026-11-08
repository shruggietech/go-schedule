# Tasks: Early-crash incident evidence

## Phase 1: Setup

- [x] T001 Preserve historical evidence, record the approved scope and complete specification/design in specs/138-early-crash-diagnostics/.

## Phase 2: Foundation

- [x] T002 Add failing-first attribution, decoding and bounded collector regressions in internal/incident/runtime_events_test.go.

## Phase 3: US1 - Early-crash evidence

Independent criterion: a failed owned Windows child receives a bounded event receipt without changing its native outcome.

- [x] T003 [US1] Implement the sanitized receipt/filter and bounded platform query in internal/incident/runtime_events*.go and test/integration/script_runtime_events*test.go.
- [x] T004 [US1] Integrate failed Windows PowerShell collection in test/integration/testscripts_observer_test.go while preserving dump status and primary results.
- [x] T005 [US1] Add the hidden pre-identity native crash probe and outcome regressions in test/integration/*incident*test.go.

## Phase 4: US2 - Honest historical attribution

Independent criterion: the investigation distinguishes the recovered event from definitive writer attribution and original causal resolution.

- [x] T006 [US2] Record recovered evidence and limits in specs/138-early-crash-diagnostics/investigation.md and docs/test-scripts.md; reconcile GitHub parent/child traceability.

## Phase 5: Polish and verification

- [x] T007 Update CHANGELOG.md and specs/README.md, run focused checks and all eight foreground canonical gates, and audit source/evidence integrity.
- [x] T008 Finalize specs/138-early-crash-diagnostics/verification.md and delivery.md with factual completion records and post-record consistency checks.

## Dependencies and implementation strategy

T001 precedes T002. Failing-first contracts precede T003, which feeds T004/T005. T006 is independently authored from retained evidence but uses the agreed receipt contract. T007/T008 follow implementation. Parser authoring and historical documentation affect different files and could run in parallel, but the primary agent executes source changes sequentially. The skill-directed research agent is read-only. Publish/review/merge operations are delivery evidence, not unchecked implementation tasks.
