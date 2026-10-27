# Tasks: Historical Windows incident differential diagnosis

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [quickstart.md](quickstart.md).

## Phase 1: Setup

- [x] T001 Synchronize accepted main and select S123 consistently in .specify/feature.json and specs/123-historical-windows-diagnosis/spec.md.
- [x] T002 Complete clarify, requirements/domain checklists and design decisions in specs/123-historical-windows-diagnosis/plan.md.

## Phase 2: Foundations

- [x] T003 Freeze original evidence and source-index/binary identities in .test-temp/s123/evidence-before.json.
- [x] T004 Recover historical target source from the base archive plus indexed overlay in .test-temp/s123/historical-source and audit the full-tree digest, explicitly recording the approved reconstruction fallback if it differs.
- [x] T005 Prepare source-bound current inputs and original/current race test executables under .test-temp/s123 after T010 replay selections; record preparation independently.

## Phase 3: US1 runtime comparison

**Goal**: Attribute exact original versus current wrapper/adapter outcomes.

**Independent test**: Six concurrent writers on fresh/existing schemas for each declared payload/runtime selection.

- [x] T006 [US1] Declare the complete fixed matrix and shared deadline in specs/123-historical-windows-diagnosis/investigation.md before execution.
- [x] T007 [US1] Run historical/current binding tests with retained/current PowerShell and Bash, preserving native exits and row assertions in .test-temp/s123/campaign-01.
- [x] T008 [US1] Analyze any actual runtime/prerequisite recurrence, file separate demonstrated defects and correct affected paths with regressions; record distinct conclusions in specs/123-historical-windows-diagnosis/investigation.md.

## Phase 4: US2 IPC comparison

**Goal**: Compare original operation shapes while preserving controller/component provenance.

**Independent test**: Historical/current packaged WAL candidates at independently selected allowances with original catalog assertions.

- [x] T009 [US2] Add test-first duration/collector selection and invalid-input coverage in internal/clievaluate/incident_replay_test.go.
- [x] T010 [US2] Implement explicit replay selections and provenance in internal/clievaluate/incident_replay_test.go without changing ordinary QA defaults.
- [x] T011 [US2] Run the declared packaged comparison, preserving actual component identities and unavailable historical collection in .test-temp/s123/campaign-01.
- [x] T012 [US2] Analyze observed IPC divergence, correct separately tracked demonstrated defects with regressions and record independent results in specs/123-historical-windows-diagnosis/investigation.md.

## Phase 5: US3 disposition and guidance

**Goal**: Provide an actionable causal conclusion without inventing a historical cause.

**Independent test**: Every symptom has precise evidence, outcome, correction or unresolved disposition.

- [x] T013 [US3] Deliver the six-section hypothesis/outcome report in specs/123-historical-windows-diagnosis/investigation.md and update docs/maintainers/windows-incident-replay.md.
- [x] T014 [US3] Reconcile merged S122 source/CI/review identities in specs/122-windows-incident-replay/delivery.md and verification.md.
- [x] T015 [US3] Reconcile #355/#356 and separate child issue records with actual source outcomes; retain unresolved causal outcomes in specs/123-historical-windows-diagnosis/delivery.md.

## Phase 6: Integration and completion

- [x] T016 Run existing full source-bound CLI QA after application changes, or record explicit non-applicability for diagnostic/documentation-only work in specs/123-historical-windows-diagnosis/verification.md.
- [x] T017 Run canonical foreground scripts/verify.sh all, formatting/UTF-8 checks and frozen-source/evidence audits; preserve first failures in .test-temp/s123.
- [x] T018 Finalize specs/123-historical-windows-diagnosis/delivery.md, verification.md, task states, CHANGELOG.md and specs/README.md with measured outcomes and exact boundaries.

## Dependencies

T001-T002 precede the blocking analyze gate. T003-T005 establish immutable comparison inputs. T009-T010 must precede the compiled current IPC executable in T005 and campaign execution; T006 freezes the matrix before T007/T011. Runtime and IPC groups execute sequentially under one shared deadline. T008/T012 depend on actual findings, not invented mandatory production edits. T013-T015 follow those dispositions; T016-T018 provide integrated acceptance and factual completion.

## Parallel opportunities

Read-only historical source research and IPC compatibility research used independent Spec-Kit research agents. Optional independent documentation preparation can proceed separately from input hashing, but campaign groups and mutations stay sequential. No parallel writer groups are added beyond each original six-writer test.

## Implementation strategy

First recover and attribute exact inputs, then add only missing existing replay selections under TDD, prepare/freeze inputs and execute one bounded campaign. Correct real defects and preserve uncertainty where no matching recurrence exists. Publication/review/merge bookkeeping is not an open implementation task.

## Task dispositions

T004 used the specification-approved reconstruction fallback: target wrapper/adapter blobs are recoverable, but exact full-source recovery is unproven because the reconstructed digest differs from the original receipt. This is explicitly not a claim that the missing exact-source outcome was performed. T008/T012 found no new production defect or matching original recurrence, so no speculative application fix or separate product bug was invented. T016 is not applicable because the changes are diagnostic tests and documentation only; no unchanged full-catalog repetition is represented as performed. Canonical source verification remains required under T017.
