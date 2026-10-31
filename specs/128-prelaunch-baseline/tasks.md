# Tasks: S128 Prelaunch Baseline

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [current contract](contracts/current-baseline.md).
**Organization**: Three current-state user stories, with independent file groups and integrated verification after they settle.

## Phase 1: Setup

- [x] T001 Create S128 specification, requirements/domain checklists and active feature selector in specs/128-prelaunch-baseline/ and .specify/feature.json.
- [x] T002 Capture accepted current DDL/evidence inventory and create native issue #408/milestone 18, retaining private receipts under .test-temp/s128/.

## Phase 2: Foundation

- [x] T003 Complete the blocking read-only Spec-Kit cross-artifact analysis and retain its result in specs/128-prelaunch-baseline/analysis.md.

## Phase 3: US1 - Current database baseline

**Goal**: Direct complete initialization, current reopen and explicit refusal/reset without old schema conversion.
**Independent Test**: Fresh/current/unsupported/atomicity contracts and accepted current schema comparison.

- [x] T004 [P] [US1] Add failing bootstrap/reopen/refusal/atomicity contracts and a fresh-start benchmark in internal/store/, retaining baseline output in .test-temp/s128/.
- [x] T005 [US1] Replace the migration registry with direct current_schema.sql initialization and current marker 24 validation in internal/store/store.go.
- [x] T006 [US1] Retire seeded historical conversions while preserving current assertions in internal/store/ tests, including notification/identity/index/portable behavior.
- [x] T007 [US1] Replace metrics upgrades and legacy reader fallbacks with one current format 4 in test/scripts/lib/sqlite.sh, Sqlite.ps1, Test-ReadTestDB twins and test/integration/testscripts_test.go.

## Phase 4: US2 - Current preferences

**Goal**: Current defaults/persistence/controls without legacy import state or UI.
**Independent Test**: Shared preferences, CLI/evaluator, settings and browser current behavior.

- [x] T008 [P] [US2] Add failing no-import/no-transition/incomplete-current contracts in internal/userprefs/, internal/cli/ and desktop/settings/frontend tests.
- [x] T009 [US2] Remove legacy preference paths/imports/transition models and QA injection in internal/userprefs/, internal/cli/, internal/clievaluate/ and desktop/settings/.
- [x] T010 [US2] Remove migration UI/current response fields and revise meaningful settings/browser fixtures in desktop/frontend/, preserving popup/default/offline/native-action behavior.

## Phase 5: US3 - Prelaunch policy and qualification

**Goal**: Active policy/qualification matches no migration obligation while retaining current security/lifecycle outcomes.
**Independent Test**: Current fresh/session/packet/silent-installer and guidance contracts, without old baseline artifacts.

- [x] T011 [P] [US3] Add or revise failing current-only qualification contracts in internal/releasegate/, scripts/windows-qualification-session/ and test/integration/windows_installer_contract_test.go.
- [x] T012 [US3] Amend .specify/memory/constitution.md to 4.0.0 and synchronize docs/build-autopilot.md, CLAUDE.md, CONTRIBUTING.md and .github/PULL_REQUEST_TEMPLATE.md.
- [x] T013 [US3] Remove historical baseline-MSI/upgrade matrices from .github/workflows/ci.yml, test/windows/ installer harnesses and scripts/automation-check.sh contracts.
- [x] T014 [US3] Retire historical roles/observations and implement current evidence revision 2 in internal/releasegate/, scripts/windows-qualification-session/, current collector/guest tools and their synthetic tests.
- [x] T015 [US3] Update active docs/ and test/windows/README.md with current prelaunch/reset guidance; update #404's superseded S128 allocation and #408 scope, keeping historical evidence intact.

## Phase 6: Integrated verification and delivery

- [x] T016 Compare current schema, measure fresh startup, run focused race/integration/metrics twin tests, current desktop/site browser checks and shell/PowerShell compliance; retain first failures in .test-temp/s128/.
- [x] T017 Run complete foreground sh scripts/verify.sh all and protected-input/source integrity checks; record actual results in specs/128-prelaunch-baseline/verification.md.
- [x] T018 Finalize CHANGELOG.md, S128 delivery/spec/task records, specs/README.md lifecycle inventory and objective local evidence.

## Dependencies and parallel execution

T001 -> T002 -> T003 precede all implementation. US1 (T004-T007), US2 (T008-T010) and US3 (T011-T015) may run as independent parallel file groups after T003; each group runs its tests before its implementation and retains red/green receipts. Current code dependencies require all three groups complete before T016-T018.

Store group owns internal/store, metrics helper/readers and testscripts_test.go. Preference group owns shared preferences, preference CLI/evaluator wiring, desktop settings and matching frontend fixtures. Policy group owns current docs/governance/CI/installer/session/releasegate and corresponding tests; it does not edit the other groups' files.

## Implementation strategy

Deliver US1's direct current database first, US2's current preference surface and US3's policy/current qualification coherently in the same PR. The installed tasks/implement workflow's parallel execution pattern is used only for these disjoint groups. No historical conversion behavior is relabeled current or replaced by weakened security tests.

Push, PR, review, merge and branch cleanup are publication evidence, not unchecked implementation tasks. Stop before owner merge; at most two Codex requests including PR-open automation.

## Verification discovery (2026-10-31)

- [x] T019 Track #409, preserve the original canonical desktop failure, and make desktop/taskgroup preview expectations deterministic through an injected clock and explicit overlap choices. The focused race suite passes; T017 restarts after this correction.

- [x] T020 Track #410, retain the original coverage failure, and correct current helper read-only contention handling with deterministic twin lock/refusal coverage and unchanged existing budgets/counts. Complete before restarting T017.
