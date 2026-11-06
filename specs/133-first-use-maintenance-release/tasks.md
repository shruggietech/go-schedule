# Tasks: S133 First-use Maintenance Release

**Input**: [plan.md](plan.md), [spec.md](spec.md), research, data model, release contract and quickstart.

**Tests**: Required failing-first automation regression and all existing canonical gates.

## Phase 1: Setup

- [x] T001 Establish S133 specification, checklists and owner boundaries in specs/133-first-use-maintenance-release/.
- [x] T002 Dispatch research, resolve candidate classification and update CLAUDE.md plan context.

## Phase 2: Foundation

- [x] T003 Freeze retained inputs and historical release records in .test-temp/s133/protected-inputs.json.
- [x] T004 Complete the blocking read-only analyze against all S133 artifacts.

## Phase 3: US1 - Coherent build baseline

**Goal**: Align every active build entry point and the four dependency proposals.

**Independent test**: Existing positive/mutation automation cases plus read-only module/checksum acceptance.

- [x] T005 [US1] Add failing-first Wails runtime/CLI mismatch cases to test/scripts/automation-check_test.sh and preserve the pre-fix result.
- [x] T006 [US1] Integrate exact PR #420-#423 module and npm manifest/lockfile changes.
- [x] T007 [US1] Align Makefile, scripts/verify.sh, desktop/README.md and CI/release Wails pins; derive existing scripts/automation-check.sh expectations from desktop/go.mod.
- [x] T008 [US1] Pass existing mutation fixtures, module read-only/checksum checks and locked frontend/site dependency restores/audits.

## Phase 4: US2 - Maintenance source and release continuation

**Goal**: Prepare accurate v2.0.1 source while public selection remains v2.0.0.

**Independent test**: Existing release-source/comparison guards and first-use Run ID/browser contracts.

- [x] T009 [US2] Prepare README.md candidate health example, CHANGELOG.md 2.0.1 and .github/release-notes/v2.0.1.md with current-public attribution.
- [x] T010 [US2] Review official comparison sources and generate site/data/comparisons/v2.0.1.json and matching content-bound review.
- [x] T011 [US2] Record exact later owner integration, staging, inspection, promotion and discovery continuation in handoff.md; keep #419 incomplete.

## Phase 5: Verification and factual closeout

- [x] T012 Run foreground sh scripts/verify.sh all against frozen source and preserve terminal gate evidence.
- [x] T013 Exercise existing desktop/site browser contracts and scripts/cli-qa.ps1 current-candidate build/capture/review/finalize/validate/render with owned paths.
- [x] T014 Audit source/evidence/history integrity, formatting/lifecycle and release guards; finalize verification.md, delivery.md and implemented status.

## Dependencies and execution

T001-T004 precede implementation. T005 precedes T006-T008; T009-T011 follow the selected dependency baseline. T012-T014 use quiescent final source. US1 is independently useful; US2 depends on its corrected candidate. Documentation and lockfile work could run in parallel on distinct files, but shared candidate/manifest finalization is sequential.

## Implementation strategy

Deliver the source pin correction and bounded dependencies first, then coherent candidate metadata and handoff. Publication, PR review/merge and later tag/release workflow evidence are recorded outside implementation task checkboxes. The requested PR stop does not claim #419's actual public outcome.
