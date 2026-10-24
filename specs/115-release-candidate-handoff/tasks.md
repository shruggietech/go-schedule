# Tasks: S115 release candidate handoff

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, data model and contracts.

**Delivery**: Implementation complete through PR #369 on codex/115-release-candidate-handoff, with final canonical verify-all-04 accepted; pushed-head CI and review workflow are recorded in the PR, outside implementation checkboxes.

## Phase 1: Setup

- [x] T001 Establish S115 branch, selector and complete design in `specs/115-release-candidate-handoff/`.

## Phase 2: Foundation

- [x] T002 Complete blocking cross-artifact analysis for `specs/115-release-candidate-handoff/spec.md`, `plan.md` and `tasks.md` before any implementation.

## Phase 3: US1 current source handoff (P1)

Goal: reusable pre-tag source validation and truthful release preparation. Independent criterion: candidate metadata succeeds; inconsistent fixture metadata fails with an actionable file-specific error.

- [x] T003 [US1] Incorporate authorized 18e6fd8 reconciliation into its existing `specs/114-release-readiness/` and release records, preserving archived 574d071.
- [x] T004 [US1] Demonstrate dated-heading refusal gap and add source metadata fixture matrix in `scripts/release-source-check/main_test.go`.
- [x] T005 [US1] Implement read-only validated arguments, badge/health/calendar-heading/notes-link preflight in `scripts/release-source-check/main.go`.
- [x] T006 [US1] Wire shared command into `.github/workflows/release.yml` and replace obsolete inline policy assertions/negative fixtures in `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`.
- [x] T007 [US1] Document pre-tag invocation, upgrade/rollback and immutable later staging/promotion sequence in existing release docs and `specs/115-release-candidate-handoff/handoff.md`; refresh governed comparison digests if needed.

## Phase 4: US2 native outcomes (P1)

Goal: attributable simultaneous native workload and exact current candidate observations. Independent criterion: successful complete snapshots preserve all state during exactly sixteen mutations, and writer reports retain six original outcomes.

- [x] T008 [US2] Extend `internal/clievaluate/ipc_diagnostics_test.go` with owned memory/file-backed concurrent full snapshots and sixteen exact create/delete pairs, bounded joined workers and existing cancellation assertions.
- [x] T009 [US2] Run original native Windows/Linux diagnostic six-writer and IPC paths once per declared experiment; record source/runtime identities and distinct #355/#356 dispositions in `specs/115-release-candidate-handoff/observations.md`.
- [x] T010 [US2] Build six receipt-bound binaries and perform fresh full owned capture, actual message review, finalize/validate/render using `scripts/cli-qa.ps1`; summarize exact source, binary, command, observation and cleanup outcomes in `specs/115-release-candidate-handoff/observations.md`.
- [x] T011 [US2] Correct any demonstrated in-scope defect with regression coverage in its affected source path and native issue references in `specs/115-release-candidate-handoff/observations.md`; preserve unknown historical causes.

## Phase 5: US3 identity and recovery (P2)

Goal: stale/incomplete inputs remain distinguishable and future publication preserves identities. Independent criterion: existing refusal/recovery regressions and written handoff agree without creating tags or packages.

- [x] T012 [US3] Exercise existing stale-source, binary receipt, interrupted observation and publication identity refusal/recovery regressions in `internal/clievaluate/`, `internal/clifixture/` and release tooling; record applicable outcomes in `specs/115-release-candidate-handoff/observations.md`.
- [x] T013 [US3] Reconcile exact frozen candidate versus later metadata identities, absent staged/public artifacts and issue-level functional dispositions in `specs/115-release-candidate-handoff/handoff.md`.

## Verification discovery adjustment

- [x] T016 [US2] Reuse bounded once-per-target packaged fixture preparation in `test/integration/package_build_test.go`, preserving independent package shapes and all existing runtime assertions in CLI/MCP/v1.3/v1.4 integration cases.

- [x] T017 [US2] Synchronize the private timeout fixture in `test/integration/cli_qa_wrapper_test.go` on child output readiness while retaining the original 150 ms policy, timeout output and input assertions.

- [x] T018 [US2] Correct #370 PowerShell persistence help and #371 whole-session MCP wording in `internal/cli/command_contract.go` and `internal/mcpobserve/server.go`, with actual CLI/SDK regressions and unchanged generator/authority contracts.

## Phase 6: Cross-cutting completion

- [x] T014 Run foreground canonical eight gates through `scripts/verify.sh`, formatting and encoding checks; record exact platform split/results in `specs/115-release-candidate-handoff/observations.md`.
- [x] T015 Complete functional task and lifecycle metadata in `specs/115-release-candidate-handoff/spec.md`, `tasks.md` and current context files using real source delivery references.

## Dependencies and implementation strategy

T001 -> T002 blocks all implementation. US1 T003-T007 proceeds in order; US2 T008 can follow T002 independently of source tooling, but T009/T010 require the final candidate source. T014 canonical acceptance precedes the frozen source checkpoint for T010; rerun applicable canonical gates after any behavioral correction in T011. US3 T012 can accompany focused verification and T013 follows current observations. T015 follows functional completion. Release-source validation is the smallest useful increment, but this authorized slice delivers the full bundle.

## Parallel opportunities

US1 fixture review and US2 IPC design use separate paths after T002; US3 refusal tests can be read alongside focused tests. Actual native experiments and full evaluator run sequentially to preserve the resource budget. No additional implementation agents are required.

## Publication workflow

Owner explicitly authorized one branch push and official PR, exact-head CI monitoring and up to two Codex rounds. Record actual findings and round identities in the PR. Stop before merge, tag, staging or promotion; #354 stays open for its publication outcome and unresolved #355/#356 remain open.
