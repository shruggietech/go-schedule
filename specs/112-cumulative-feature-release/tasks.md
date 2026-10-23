# Tasks: S112 cumulative feature release

**Input**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md).

## Phase 1: Setup

- [x] T001 Verify clean synced main and S112 branch/selector identity in .specify/feature.json.
- [x] T002 Produce spec, clarification, checklists and design artifacts in specs/112-cumulative-feature-release/.

## Phase 2: Foundation

- [x] T003 Retain unchanged merge-head CI evidence and reproduce queued observation masking using test/windows/Invoke-ServiceCoreCI.ps1 without installing a service.
- [x] T004 Record blocking Spec-Kit analysis coverage in specs/112-cumulative-feature-release/analysis.md before source implementation.

## Phase 3: US1 - Accurate CI outcomes

**Independent test**: executable queued/terminal and failure-retention fixtures plus the real packaged helper authentication regression.

- [x] T005 [US1] Add failing executable PowerShell fixtures in test/integration/service_core_probe_test.go for queued/terminal selection and primary/cleanup evidence.
- [x] T006 [US1] Correct terminal selection and bounded failure checkpoint/cleanup in test/windows/Invoke-ServiceCoreCI.ps1.
- [x] T007 [US1] Regress setup/runtime lifetime ownership in internal/clientsecret/context_test.go, then separate bounded fixture build and runtime contexts with setup diagnostics.
- [x] T008 [US1] Rerun targeted regressions and update disposable probe documentation in test/windows/README.md.

## Phase 4: US2 - Dependency convergence

**Independent test**: module graph, production builds, policy negative fixtures and canonical gates.

- [x] T009 [US2] Reconcile #351 and #273 into root go.mod/go.sum; preserve SMTP and record supersession pending merge.
- [x] T010 [US2] Apply #336 to desktop/go.mod/go.sum and align active Wails CLI build pins/fixtures to 2.17.0.
- [x] T011 [P] [US2] Apply #337 desktop/frontend/package.json and package-lock.json; install matching managed test browsers.
- [x] T012 [P] [US2] Apply #295 site/package.json and package-lock.json; retain matching runtime/tool versions.
- [x] T013 [US2] Migrate .github/workflows/pages.yml to reviewed action majors and update scripts/automation-check.sh plus negative fixtures in test/scripts/automation-check_test.sh.

## Phase 5: US3 - Release preparation

**Independent test**: version-bound comparison/docs/release contracts and unchanged published release pointer.

- [x] T014 [US3] Prepare .github/release-notes/v1.6.0.md, CHANGELOG.md and README.md version/upgrade references while preserving pending publication status.
- [x] T015 [US3] Prepare site/data/comparisons/v1.6.0.json with exact-content review and current feature/opportunity claims, preserving site/data/published.json.
- [x] T016 [US3] Reconcile specs/README.md, S111 delivery metadata and candidate preparation instructions in specs/112-cumulative-feature-release/delivery.md.

## Phase 6: Polish and publication

- [x] T017 Run all eight canonical gates in the foreground and owned Windows/WSL QA; record exact evidence and communication review in specs/112-cumulative-feature-release/verification.md.
- [x] T018 Run go run ./scripts/github-format, encoding sanity and diff/issue coverage review in specs/112-cumulative-feature-release/.

## Dependencies and execution order

T001-T004 precede implementation. US1 corrections precede dependency changes. T009 precedes T010; npm tasks may run independently of one another after US1. US3 consumes the converged baseline. All eight gates and QA precede the first commit/push. Hosted checks and bounded reviews follow publication as workflow evidence, separate from feature maturity. Public release staging/promotion remains a later owner-authorized operation under #354.

## Parallel opportunities

US1 fixture/deadline analysis can be researched separately. T011 and T012 modify separate module lockfiles and may be evaluated independently. Research agents are read-only; source writes remain sequential.

## Implementation strategy

Deliver US1 first with before/after evidence, then integrate dependencies as one reviewable baseline, then bind release preparation to that source. Keep all user scenarios traceable and report publication as outstanding until actually authorized and delivered.

## Authorized publication workflow (not feature maturity)

WP1: Commit with the conventional subject and project trailer, push codex/112-cumulative-feature-release, open and attach the official PR, and update #354 traceability without closing publication outcomes.

WP2: Monitor exact-head CI and every bot review, resolve all warranted findings and every thread, request at most one second Codex round and record the final state in review.md. Stop before merge, tags or release publication. These obligations remain part of the owner request; they are not open feature-implementation tasks.
