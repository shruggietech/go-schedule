# Tasks: S127 Clustered-Execution Deferral

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contract](contracts/decision.md).

**Tests**: No runtime behavior changes; use existing engineering gates and meaningful documentation/browser checks. No new tests that mirror prose.

## Phase 1: Setup

- [x] T001 Create the S127 specification, requirement checklists and feature selector in specs/127-cluster-decision-deferral/ and .specify/feature.json.

## Phase 2: Foundational

- [x] T002 Resolve owner direction and record the existing issue dependency/authority model in specs/127-cluster-decision-deferral/plan.md and contracts/decision.md.
- [x] T003 Complete the read-only Spec-Kit cross-artifact gate and retain its report in specs/127-cluster-decision-deferral/analysis.md before source implementation.

## Phase 3: US1 - Explicit owner direction

**Goal**: Current source and native records communicate explicit deferral, distinct from historical hold, final no-go or delivered capability.
**Independent Test**: Trace current page and issue dispositions to the owner answer and retained unmet outcomes.

- [x] T004 [US1] Write the attributable decision and all six issue dispositions in specs/127-cluster-decision-deferral/decision.md.
- [x] T005 [US1] Update docs/distributed-scheduling.md with current deferral and accurate public v1.7.0 comparison attribution while preserving historical evidence.
- [x] T006 [US1] Reconcile native GitHub #186-#189, #20 and #146 with the decision contract; preserve exact before/after/readback receipts in .test-temp/s127/ and keep all six issues open.

## Phase 4: US2 - Deliberate reactivation

**Goal**: Reactivation has an owner boundary and preserves all staged decisions without a date or automation.
**Independent Test**: Trace a workload report and a separate owner revision through the retained four-stage sequence.

- [x] T007 [US2] Define later owner revisit, evidence intake and retained semantics-approval boundaries in docs/distributed-scheduling.md and specs/127-cluster-decision-deferral/decision.md.
- [x] T008 [US2] Record explicit scope, incomplete outcomes and excluded runtime/release operations in specs/127-cluster-decision-deferral/delivery.md.

## Phase 5: Polish and verification

- [x] T009 Add the dated owner deferral to CHANGELOG.md, update the S127 row in specs/README.md and managed plan reference in CLAUDE.md.
- [x] T010 Format public prose and run full foreground sh scripts/verify.sh all, applicable site browser/lifecycle checks and frozen evidence/source integrity checks; retain actual results in specs/127-cluster-decision-deferral/verification.md and .test-temp/s127/.
- [x] T011 Resolve all required tasks and advance the S127 lifecycle with factual delivery evidence in specs/127-cluster-decision-deferral/spec.md and specs/README.md.

## Dependencies and execution order

T001 -> T002 -> T003 -> T004/T005 -> T006 -> T007/T008 -> T009 -> T010 -> T011. T007 refines the same documents as T004/T005 and is executed sequentially to avoid conflicting edits.

## Parallel opportunities

Read-only source and native issue discovery may be batched. The decision and current page share authority and are authored sequentially; no agent delegation is needed. No implementation task has a misleading parallel marker.

## Implementation strategy

US1 supplies the smallest useful outcome: an attributable current deferral. US2 makes it reversible and actionable without promising delivery. Verification and lifecycle finalization follow both stories. Push, PR, review and merge are publication workflow evidence, not unchecked implementation tasks.
