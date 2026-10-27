# Tasks: S120 Windows causal diagnostics

**Input**: specs/120-windows-causal-diagnostics/spec.md and plan.md.

## Phase 1: Setup

- [x] T001 Select feature 120 consistently in .specify/feature.json, CLAUDE.md and specs/README.md, synchronizing accepted main before evidence gathering.
- [x] T002 Complete specify, clarify, requirements checklists and plan artifacts under specs/120-windows-causal-diagnostics/.
- [x] T003 Perform blocking semantic analysis and implementation prerequisites, recording findings in specs/120-windows-causal-diagnostics/analysis.md.

## Phase 2: Foundation

- [x] T004 Preserve original S112/S115/S117 evidence identities and exact accepted source in .test-temp/s120/evidence-before.json (FR-001, FR-010).

## Phase 3: US1 Runtime and prerequisite investigation (P1)

- [x] T005 [US1] Add an owned inherited-output behavioral regression in test/integration/testscripts_observation_test.go and capture its failure before changing observation (FR-007, FR-009).
- [x] T006 [US1] Search existing issues and file any reproduced observation defect against parent #355, milestone 16 (FR-011).
- [x] T007 [US1] Bound script execution/output collection and preserve root exit, collection and cleanup separately in test/integration/testscripts_binding_test.go and test/integration/testscripts_observation_test.go, reusing internal/clifixture (FR-003, FR-007, FR-008).
- [x] T008 [US1] Add opt-in outer executable/harness provenance and exercise metadata-only versus writer and implicit versus explicit SQLite selection using existing probe fixtures (FR-002, FR-004, FR-010).
- [x] T009 [US1] Retain six-writer fresh/existing-schema assertions and record distinct interpreter/prerequisite dispositions in specs/120-windows-causal-diagnostics/investigation.md (FR-001 through FR-004).

## Phase 4: US2 IPC and store cancellation (P1)

- [x] T010 [US2] Add test-only owned pool acquisition/stats in internal/store/export_test.go and real native IPC cancellation regression in internal/store/request_cancellation_test.go; capture failure before correction (FR-005, FR-009).
- [x] T011 [US2] Search existing issues and file the reproduced cancellation defect against #356, milestone 16 (FR-011).
- [x] T012 [US2] Add explicit context-aware authorization, identity, snapshot and cleanup persistence methods in internal/store/{access,identity,crud,remote_access}.go with legacy delegates (FR-003, FR-005).
- [x] T013 [US2] Propagate request contexts in affected internal/api/server handlers, retaining durable audit intent and bounded independent terminal completion (FR-005, FR-008, FR-009).
- [x] T014 [US2] Exercise cancellation at authentication/read/intent/cleanup boundaries, verify audit semantics, existing partial snapshot recovery and original service/MCP/group request shapes; record individual dispositions (FR-005, FR-006, FR-009, FR-010).

## Phase 5: US3 Attributable delivery (P2)

- [x] T015 [US3] Publish concise source-bound investigation instructions in docs/test-scripts.md and docs/cli-qa.md and actual increment in CHANGELOG.md; update #355/#356 without unsupported closure (FR-001, FR-011).
- [x] T016 [US3] Reconcile S119 PR #380 merge in specs/119-release-execution/ and specs/README.md, retaining immutable public release identities (FR-012).

## Phase 6: Cross-cutting acceptance

- [x] T017 Run targeted native Windows/Linux experiments, authorization/audit/fixture regressions and quickstart validation; resolve attributable failures without weakening contracts (FR-004 through FR-010).
- [x] T018 Freeze source inputs, run all eight foreground scripts/verify.sh gates, audit original evidence unchanged and record exact inputs/results in specs/120-windows-causal-diagnostics/verification.md (FR-013).
- [x] T019 Finalize lifecycle/task metadata and delivery record, run GitHub formatting, commit only after the gates pass (FR-012, FR-013).

## Dependencies and execution order

T001-T004 precede story implementation. T005-T006 precede T007; T008-T009 complete US1. T010-T011 precede T012-T013; T014 completes US2. US1 and US2 have independent component boundaries but execute sequentially with existing owned resources. T015-T016 depend on attributable findings, then T017-T019 depend on both stories. No additional implementation delegation is required.

## Delivery protocol

After completed implementation tasks, push the authorized branch and publish one official PR. Inspect every applicable CI check and review thread on the exact final head, address warranted findings, and allow at most two Codex review rounds. Publication and external review monitoring are delivery operations, not unfinished implementation tasks. Stop before merging, tagging or promotion (FR-014).
