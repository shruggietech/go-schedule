# Tasks: Unattended CLI evaluation and canonical QA reports

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md), [contracts](contracts/), [research.md](research.md).

**Scope**: Full #289/#290 bundle. Constitution-mandated behavioral/adversarial tests precede each implementation. Task completion describes implementation, not GitHub issue closure gates.

## Phase 1: Setup

- [x] T001 Establish shared maintainer package and thin command structure in internal/clievaluate and cmd/gosched-qa (FR-001).
- [x] T002 Add strict decoding, schema and report-integrity failing tests in internal/clievaluate/model_test.go (FR-012/013/015).

## Phase 2: Shared foundations

- [x] T003 Implement concrete shared model, strict decoders and embedded versioned schemas in internal/clievaluate/model.go, decode.go, schema.go and schemas/v1 (FR-012/013).
- [x] T004 Add adversarial bounded redaction, artifact path/hash/link and reference/count tests in internal/clievaluate/artifacts_test.go and validate_test.go (FR-007/009/015).
- [x] T005 Implement confined redacting artifact storage and canonical/semantic validators in internal/clievaluate/artifacts.go and validate.go (FR-007/009/015).

## Phase 3: US2 - Complete inventory and behavioral coverage

**Independent outcome**: Actual command metadata and meaningful behavior classes map to explicit stable cases; changed metadata cannot silently retain coverage.

- [x] T006 [US2] Add actual-tree, alias/generated-command/flag-shadowing and uncovered-case tests in internal/cli/inventory_test.go and internal/clievaluate/inventory_test.go (FR-002).
- [x] T007 [US2] Export initialized Cobra metadata and implement checked-in inventory/manifest drift validation in internal/cli/inventory.go, internal/clievaluate/inventory.go and manifest.go (FR-002).
- [x] T008 [US2] Define the complete command behavior, output, negative and S107 regression matrix in internal/clievaluate/scenarios*.go and the versioned manifest (FR-003/004).

## Phase 4: US1 - Quiet bounded execution

**Independent outcome**: Both selected environments execute the planned cases under verified independent limits against only owned state, with complete partial/interrupted ledgers.

- [x] T009 [US1] Add resource-bound, hidden capture, exact cleanup and excessive-output tests in internal/clifixture/process*_test.go and internal/clievaluate/environment_test.go (FR-005/006/007).
- [x] T010 [US1] Extend optional captured output and Windows CPU/memory/process limits without changing existing callers in internal/clifixture/process*.go (FR-005/006/007).
- [x] T011 [US1] Implement native/selected-WSL preflight, transient owned cgroup workers, applied-limit verification and watchdog/liveness cleanup in internal/clievaluate/environment*.go and coordinator.go (FR-005/006/007).
- [x] T012 [US1] Add complete scenario assertion/negative, interruption and partial-environment tests in internal/clievaluate/coordinator_windows_test.go, policies_test.go and process_test.go (FR-003/004/008).
- [x] T013 [US1] Implement owned packaged-daemon/CLI lifecycle and independently correlated state/history/audit/log evidence in internal/clievaluate/worker.go, session.go and scenarios*.go (FR-003/004/007/008).
- [x] T014 [US1] Implement all scoped service/GUI/profile/credential/preference/MCP handlers and failure cases in internal/clievaluate/transports.go, protocol.go, http_protocol.go, secret_boundary.go and postconditions.go (FR-003/004/005).
- [x] T015 [US1] Implement deterministic schedule/DST/month-boundary/catch-up/overlap/watcher/notification fixtures in internal/clievaluate/policies.go, recovery.go and notification_transport.go (FR-003/004).
- [x] T016 [US1] Add changed-subject, corrupt/truncated event, retained-evidence and inactive-ownership resume tests in internal/clievaluate/resume_test.go (FR-017).
- [x] T017 [US1] Implement complete checkpoint/event/interruption/resume identity handling in internal/clievaluate/coordinator.go, worker.go, events.go and resume.go (FR-008/017).

## Phase 5: US3 - Communication capture and agent review

**Independent outcome**: A frozen contextual corpus and exact complete agent reviews are required for finalization, with explicit no-message cases and pending_review outcomes.

- [x] T018 [US3] Add split-write/dynamic-secret, message identity/context, stale/missing/extra review and false-pass tests in internal/clievaluate/artifacts_test.go, freeze_test.go and validate_test.go (FR-009/010/011).
- [x] T019 [US3] Implement bounded corpus freeze, contextual IDs and review validation in internal/clievaluate/corpus.go and review.go (FR-009/010/011).

## Phase 6: US4 - Canonical reports and maintainer handoff

**Independent outcome**: Validated canonical JSON and deterministic seven-section Markdown expose exact evidence and independently actionable findings, with reusable agent prompts.

- [x] T020 [US4] Add deterministic rendering, finding/reference/span validation and fixed-outline tests in internal/clievaluate/render_test.go and validate_test.go (FR-014/015/016).
- [x] T021 [US4] Implement truthful summaries, stable finding deduplication, seven-section rendering and canonical report output in internal/clievaluate/render.go and validate.go (FR-013/014/015/016).
- [x] T022 [US4] Implement inventory/capture/resume/finalize/validate/render and worker/benign-helper entry points in cmd/gosched-qa/main.go and internal/clievaluate/coordinator.go (FR-001/010/017).
- [x] T023 [US4] Document full unattended execution, semantic review, budgets, blocked/resume behavior and issue triage in docs/cli-qa.md and docs/maintainers/cli-qa-run-v1.md and cli-qa-triage-v1.md (FR-001/010/018).

## Phase 7: Integration and polish

- [x] T024 Exercise actual native Windows and selected WSL capture-review-finalize, retain precise evidence and independently track discovered defects in specs/108-unattended-cli-qa/verification.md (SC-001/002/003/004/005/006).
- [x] T025 Update docs/test-scripts.md, docs/cli.md, CHANGELOG.md, specs/README.md and feature lifecycle/traceability without advertising a released capability (FR-001/018, SC-006).
- [x] T026 Run canonical foreground sh scripts/verify.sh all, affected native Linux race checks, supported compilation and encoding/GitHub-format/lifecycle checks; record exact outcomes in specs/108-unattended-cli-qa/verification.md (SC-006).
- [x] T027 Resolve all in-scope findings, verify complete functional issue mapping and commit locally on codex/108-unattended-cli-qa before the requested pre-publication halt (SC-006).

## Dependencies and implementation strategy

Setup -> shared foundations -> inventory -> bounded execution -> corpus/reviews -> reports/command/docs -> integration/verification/local commit. US2 inventory precedes US1 enumeration. US3/US4 share the model and consume actual capture events. Complete the entire bundle; no partial MVP substitutes for the authorized scope.

Parallel opportunities: independent research already ran through the installed plan command. During implementation, distinct contract tests or document drafts can proceed independently after their shared types stabilize; changes to shared model/capture files remain sequential. This task plan does not require additional agents.
