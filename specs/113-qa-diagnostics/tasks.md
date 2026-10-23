# Tasks: S113 reliable unattended QA and actionable delivery diagnostics

**Input**: spec.md, plan.md, research.md, data-model.md and contracts/diagnostics.md.
**Tests**: Regression-first behavior checks are required by the constitution. Publication, bot review and merge are delivery workflow, not unchecked implementation tasks.

## Phase 1: Setup

- [x] T001 Confirm baseline, issue/milestone traceability and exact S113 selector/branch/context in .specify/feature.json, CLAUDE.md and specs/README.md.
- [x] T002 Preserve read-only S112 failure/input identities and establish focused reproduction commands in specs/113-qa-diagnostics/verification.md.

## Phase 2: Foundational contracts

- [x] T003 Add compatible report-contract and redaction regressions for optional diagnostics and incomplete observations in internal/clievaluate/model_test.go, validate_test.go and artifacts_test.go before behavior changes.

## Phase 3: US1, trust an interrupted evaluation

**Goal**: Selected budgets and bounded finalization produce truthful retained outcomes and recovery guidance.
**Independent test**: Nondefault/invalid policies, canceled workers, latest progress recovery, unchanged evidence and changed redaction.

- [x] T004 [US1] Write timing-policy, wrapper and interrupted capture regressions in cmd/gosched-qa/main_test.go, internal/clievaluate/coordinator_windows_test.go and test/integration/cli_qa_wrapper_test.go.
- [x] T005 [US1] Write exact-byte no-rewrite/redaction/accounting regression and benchmark in internal/clievaluate/artifacts_test.go.
- [x] T006 [US1] Implement one validated timing policy and capture/resume budget resolution in internal/clievaluate/timing_policy.go and cmd/gosched-qa/main.go.
- [x] T007 [US1] Consume policy, preserve timeout output and identify original recovery inputs in scripts/cli-qa.ps1.
- [x] T008 [US1] Separate execution cancellation from bounded worker shutdown and aggregate finalization in internal/clievaluate/process.go and coordinator.go.
- [x] T009 [US1] Preserve latest durable interrupted case outcomes and stage/recovery information in internal/clievaluate/checkpoint.go, freeze.go, resume.go and render.go.
- [x] T010 [US1] Avoid unchanged sanitized physical replacements while preserving confinement, quota and hashes in internal/clievaluate/artifacts.go; run focused lifecycle tests and the artifact benchmark.

## Phase 4: US2, identify runtime and observation failures

**Goal**: Useful partial evidence survives, complete assertions stay strict, and demonstrated causes receive accurate dispositions.
**Independent test**: An intermediate request fails after successful observations, command success precedes inspection failure, and both script twins reject bad candidates while preserving original outcomes.

- [x] T011 [US2] Write partial-snapshot, command/inspection separation and bounded request/phase regressions in internal/clievaluate/diagnostics_test.go, ipc_diagnostics_test.go and internal/api/client/observation_test.go.
- [x] T012 [US2] Add context-scoped safe request observations and bounded phase timing in internal/api/client/ and internal/clievaluate/session.go without tokens, bodies or arbitrary URLs.
- [x] T013 [US2] Retain partial snapshot/protocol data and failed-phase facts, require complete snapshots for state assertions, and persist blocked-case diagnostics in internal/clievaluate/scenarios.go, worker.go, model.go, schemas/, validate.go and render.go.
- [x] T014 [US2] Write both-twin candidate rejection/version/fallback regressions and interpreter outcome checks in test/integration/testscripts_probe_test.go and testscripts_binding_test.go.
- [x] T015 [US2] Implement bounded candidate/version diagnostics, strict numeric versions and explicit-selection behavior in test/scripts/lib/Sqlite.ps1 and sqlite.sh, reusing existing hidden process patterns.
- [x] T016 [US2] Preserve bounded interpreter identity, original crash/timeout/exit output and usable fixture prerequisite checks in test/integration/testscripts_binding_test.go and testscripts_test.go, retaining all six writers and record assertions.
- [x] T017 [US2] Perform bounded independent runtime and sustained owned IPC reproductions, compare relevant existing/current baselines, apply only demonstrated source corrections, and record exact unresolved causes in specs/113-qa-diagnostics/verification.md and the existing #355/#356 records.

## Phase 5: US3, distinguish delivery details from errors

**Goal**: Successful SMTP evidence is neutral, real failures remain actionable and old saved history survives.
**Independent test**: Migration23->24, success after retry, partial rejection, uncertain acknowledgment, real SMTP delivery and typed local/remote/desktop representations.

- [x] T018 [US3] Write migration23->24/history and SMTP completion regressions in internal/store/migration_v24_test.go and notifications_smtp_test.go before changing persistence.
- [x] T019 [US3] Add Diagnostic and forward migration24, copy old SMTP evidence and clear only succeeded SMTP errors in internal/domain/notification.go, internal/store/store.go, notifications.go and notifications_smtp.go.
- [x] T020 [US3] Write real SMTP, typed remote/client, CLI and desktop neutral/error regressions in test/integration/smtp_notifications_test.go, internal/remote/notifications_test.go, internal/cli/notification_smtp_test.go and desktop/notifications/service_test.go.
- [x] T021 [US3] Update optional OpenAPI evidence fields and generated bindings, CLI human labels and desktop adapters/models/details in api/openapi/remote-v1.yaml, internal/remote/openapi.gen.go, internal/cli/notification.go, desktop/notifications/, desktop/frontend/src/notifications/ and desktop/frontend/wailsjs/go/models.ts.

## Phase 6: US4, current documentation tooling

**Goal**: Incorporate the existing bounded dependency proposal without changing public-release discovery.
**Independent test**: Exact site lockfile installation, audit and production documentation/browser checks.

- [x] T022 [P] [US4] Apply PR #360's core/UI16.17.7 and MDX15.4.8 patch to site/package.json and site/package-lock.json with upstream/license review retained.
- [x] T023 [US4] Run the production website checks and preserve source/public discovery in site/ and the existing docs policy.

## Phase 7: Integration and evidence

- [x] T024 Update canonical docs/cli-qa.md, notifications.md, test-scripts.md and INSTALL-windows.md for diagnostic/reporting semantics and schema24 forward upgrade/old-state rollback; update CHANGELOG.md with dated decisions and reconcile S112/S113 metadata in specs/README.md.
- [x] T025 Refresh pending .github/release-notes/v1.6.0.md and site/data/reviews/v1.6.0.json binding where content changes; preserve site/data/published.json and published snapshots.
- [x] T026 Run relevant security, strict-schema, source-identity, migration and original-outcome regressions, then foreground sh scripts/verify.sh all; resolve actual in-scope failures with evidence in specs/113-qa-diagnostics/verification.md.
- [x] T027 Run a fresh complete owned Windows/WSL evaluation through scripts/cli-qa.ps1 with exact-bound communication reviews; finalize/validate/render, retain truthful per-issue functional dispositions and cleanup in specs/113-qa-diagnostics/verification.md and delivery.md.

## Dependencies and execution order

T001-T003 precede all implementation. US1 and US2 share lifecycle/report contracts and must agree before integration; US1's durable reporting precedes lengthy US2 reproduction. US3 is functionally independent but must land before final whole-catalog message review. US4 touches separate site lockfiles and can proceed independently after the foundation; T023 follows T022. All stories precede final integration and publication readiness.

## Parallel examples

Independent read-only research covered runtime/prerequisites, QA lifecycle and SMTP persistence. T022 can run alongside diagnostic implementation because it owns only two site files. Distinct focused unit suites may run together when resource bounds permit; never run competing full catalog captures or foreground aggregate suites. Shared internal/clievaluate files have one implementation owner.

## Implementation strategy

Deliver US1's interruption contract first, then US2 diagnostics and empirical investigations, US3's complete persistence-to-UI correction and US4's site patch. Every step keeps tests meaningful. This is the full approved bundle; MVP is only an ordering aid, not permission to omit stories. Unconfirmed #355/#356 causes remain open and explicitly incomplete even if other source outcomes are delivered. The final PR uses closing keywords only for individually fulfilled issues.

- [x] T028 Resolve PR #361 round-1 checkpoint finding with durable exact-identity freeze transitions, verified interruption recovery and corruption rejection regressions (internal/clievaluate/checkpoint.go, freeze.go and checkpoint_test.go).

- [x] T029 Address #362 from S113 Linux CI: identify listener writer results by persistent row ID, model the fallback chronology explicitly and retain original/fixture rows plus writer output on failure (test/integration/testscripts_test.go).
