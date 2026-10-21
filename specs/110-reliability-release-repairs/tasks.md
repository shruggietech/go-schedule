# Tasks: S110 reliability and release-promotion repairs

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, model and contracts.

## Phase 1: Setup

- [x] T001 Verify synchronized baseline, literal branch/scope and installed commands in specs/110-reliability-release-repairs/spec.md.
- [x] T002 Complete specify, clarify and requirements checklists in specs/110-reliability-release-repairs/checklists/.

## Phase 2: Foundation

- [x] T003 Record demonstrated mechanisms, architecture alternatives and contracts in specs/110-reliability-release-repairs/{research,data-model,plan}.md and contracts/maintenance.md.
- [x] T004 Resolve speckit-analyze coverage and constitutional checks in specs/110-reliability-release-repairs/analysis.md before production edits.

## Phase 3: US1 exact release promotion (#262)

Independent outcome: authenticated draft or published recovery remains bound to one reviewed release ID and immutable contents.

- [x] T005 [P] [US1] Add failing HTTP fixture regressions for pagination, ambiguous/missing identity and mutation refusal in scripts/release-promotion/*_test.go (FR-001/FR-002).
- [x] T006 [US1] Implement versioned atomic selection and finite redacted HTTP operations in scripts/release-promotion/ (FR-001/FR-002).
- [x] T007 [US1] Add asset-ID download and initial/final inventory/hash checks, including redirects and content changes, in scripts/release-promotion/ (FR-003).
- [x] T008 [US1] Cover and implement exact checksum reuse, absent-draft upload and interrupted/published recovery in scripts/release-promotion/ (FR-004).
- [x] T009 [US1] Cover and implement draft publication, latest-only recovery and already-latest no-op in scripts/release-promotion/ (FR-004).
- [x] T010 [US1] Wire trusted workflow-source helper build before tagged checkout and preserve staging/candidate/comparison/checksum/tag controls in .github/workflows/promote-release-waived.yml (FR-003).
- [x] T011 [US1] Cover actual workflow mutation ordering and update scripts/automation-check.sh plus test/scripts/automation-check_test.sh and test/integration/ standing-waiver contract tests (FR-003).
- [x] T012 [US1] Document corrected operator path, resumed outcomes and refusal conditions in test/windows/README.md and docs/INSTALL-windows.md (FR-005).

## Phase 4: US2 watcher recovery (#271)

Independent outcome: replacements crossing observer recovery dispatch once, with unchanged revision isolation and lifecycle behavior.

- [x] T013 [P] [US2] Promote the failing fake-clock outage reproduction into internal/watcher/recovery_test.go and add pending/registration/repeated-failure cases (FR-006/FR-012).
- [x] T014 [US2] Retain same-revision file identity and pending work through full/root recovery in internal/watcher/manager.go (FR-006/FR-007).
- [x] T015 [US2] Suspend inactive candidate timers, retain bounded dispatch identity and suppress recovery/late-hint duplicates in internal/watcher/manager.go (FR-007/FR-008).
- [x] T016 [US2] Cover changed/disabled/deleted definitions, explicit Reload, missing file, subsequent replacements, health/revision and rejected dispatch semantics in internal/watcher/recovery_test.go (FR-008).

## Phase 5: US3 heartbeat record integrity (#335)

Independent outcome: distinct concurrent invocations persist distinct records and preserve real failure diagnostics with existing retry bounds.

- [x] T017 [P] [US3] Add failing exact-text and distinct numeric-looking session-ID regressions in test/integration/testscripts_binding_test.go for supported script twins (FR-009/FR-011/FR-012/FR-013).
- [x] T018 [US3] Encode non-null values as UTF-8 hex text-cast expressions in test/scripts/lib/sqlite.sh and test/scripts/lib/Sqlite.ps1, preserving nulls and SQL numeric casts (FR-013).
- [x] T019 [US3] Cover concurrent first-use and initialized writes, bounded lock retry/exhaustion and non-contention failure output in test/integration/testscripts_binding_test.go without weakening TestScriptsConcurrentWriters (FR-009/FR-010/FR-011).
- [x] T020 [US3] Document binding and diagnosis scope in docs/test-scripts.md and specs/110-reliability-release-repairs/verification.md, distinguishing historical cause from demonstrated correction (FR-010/FR-012).

## Phase 6: Integration and local delivery

- [x] T021 Review issue acceptance coverage and run focused native-platform race checks; record failing-before and passing-after evidence in specs/110-reliability-release-repairs/verification.md (SC-001/SC-002/SC-003).
- [x] T022 Update CHANGELOG.md with all three outcomes and dated pinned-workflow/architecture decision, plus specs/README.md lifecycle inventory (SC-004).
- [x] T023 Run foreground sh scripts/verify.sh all and resolve all eight named gates; record exact outcomes in specs/110-reliability-release-repairs/verification.md.
- [x] T024 Complete all required tasks, advance specification to Implemented, prepare local delivery/pull-request text and commit with attribution in specs/110-reliability-release-repairs/delivery.md (SC-004).

## Dependencies and execution order

T001-T004 precede all implementation. Each story follows failing tests -> correction -> focused validation. US1, US2 and US3 have independent files and may proceed concurrently after the foundation; T010/T011 stay sequential within US1. Shared changelog/lifecycle/delivery updates are coordinated only by the primary agent. Integration T021-T024 follows all three functional outcomes. Publication, PR creation/review, merging and release are outside the implementation task list.

## Parallel examples

The installed speckit-tasks team strategy permits one developer per independent story: US1 owns scripts/release-promotion and waiver workflow/policy/tests/docs; US2 owns internal/watcher recovery; US3 owns test-script twins and their integration regressions. No overlapping source writes or concurrent aggregate suites. The primary agent integrates and executes canonical verification.

## Implementation strategy

Deliver all three authorized stories in one coherent maintenance slice, with no deferred functional outcome represented as delivered. Independently validate each story and preserve separate issue links. Any newly discovered out-of-scope defect gets its own prepared planning record before publication approval.
