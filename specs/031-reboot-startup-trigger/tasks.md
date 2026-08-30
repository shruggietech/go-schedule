# Tasks: Scheduler Startup Trigger

**Input**: Design documents from `/specs/031-reboot-startup-trigger/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/startup-schedule.md, quickstart.md

**Tests**: Required by the feature specification and constitution. Each behavior task begins with a failing focused test.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can proceed independently in a different file after its phase prerequisites
- **[Story]**: User story traceability from spec.md

## Phase 1: Setup and specification gates

**Purpose**: Establish the review slice and prove the design is internally complete.

- [x] T001 Record S031 as Draft in `specs/README.md`, activate it in `.specify/feature.json` and `CLAUDE.md`, and validate both requirements checklists
- [x] T002 Run the Spec-Kit cross-artifact analysis over `specs/031-reboot-startup-trigger/spec.md`, `plan.md`, and `tasks.md`, resolving all critical or high findings before implementation

---

## Phase 2: Foundational event model

**Purpose**: Define one startup identity shared by every authoring and execution surface.

- [x] T003 [P] Add failing domain tests for startup event identity and distinct run origin in `internal/domain/domain_test.go`
- [x] T004 [P] Add failing schedule tests for the canonical human phrase, no upcoming runs, and event-policy compatibility in `internal/schedule/parse_test.go` and `internal/schedule/recur_test.go`
- [x] T005 Implement the typed `scheduler_startup` schedule identity, constructor/predicate, and `startup` run origin in `internal/domain/domain.go` and `internal/schedule/parse.go`
- [x] T006 Update source-syntax derivation for retained event expressions and cover cron/human startup identity in `internal/scheduleinput/input.go` and `internal/scheduleinput/input_test.go`

**Checkpoint**: One durable startup event exists without a clock occurrence or restored trigger subsystem.

---

## Phase 3: User Story 1 - Once per daemon start (Priority: P1)

**Goal**: Dispatch every eligible startup task once per engine lifecycle and never on reload.

**Independent Test**: Two fresh engines over one persisted store produce two startup runs total; reload inside one engine produces none; disabled task/group cases produce none.

- [x] T007 [US1] Add a failing controlled lifecycle test for two starts, one reload, injected startup time, and startup run history in `internal/engine/engine_extra_test.go`
- [x] T008 [P] [US1] Add failing startup eligibility tests for disabled tasks, completed state, disabled group ancestry, and unknown event identities in `internal/engine/engine_extra_test.go`
- [x] T009 [US1] Implement one startup snapshot and dispatch boundary outside recompute in `internal/engine/engine.go`
- [x] T010 [US1] Keep catch-up and clock advancement inert for event schedules and cover mixed stored schedules in `internal/catchup/catchup_test.go` and `internal/engine/engine_extra_test.go`
- [x] T011 [US1] Add a failing startup-versus-manual overlap test for queue, skip, and concurrent policies, then preserve queued trigger provenance in `internal/engine/overlap_test.go`, `internal/engine/overlap.go`, and `internal/engine/engine.go`
- [x] T012 [US1] Add persistence/reopen coverage for the event row and startup run origin, including absence of removed trigger tables, in `internal/store/store_test.go`

**Checkpoint**: Startup execution is correct, restartable, reload-proof, persisted, and overlap-safe.

---

## Phase 4: User Story 2 - Author and inspect everywhere (Priority: P1)

**Goal**: Accept both canonical forms through shared parsing, API, CLI, and GUI without fake next times.

**Independent Test**: Preview, create, update, inspect, and reopen startup tasks through each client surface with the same event schedule and no next runs.

- [x] T013 [P] [US2] Replace cron refusal tests with failing supported explain/compile cases in `internal/cron/cron_test.go`, `internal/cron/compile_test.go`, and `internal/cron/convert_test.go`
- [x] T014 [US2] Implement cron parse, explain, compile, and symmetric conversion for `@reboot` in `internal/cron/cron.go`, `internal/cron/phrase.go`, `internal/cron/compile.go`, and `internal/cron/convert.go`
- [x] T015 [P] [US2] Add failing API preview/create/update/detail cases with empty next runs and retained source syntax in `internal/api/server/tasks_test.go` and `internal/api/server/update_test.go`
- [x] T016 [US2] Make API response metadata and summaries event-aware in `internal/scheduleinput/input.go` and `internal/api/server/tasks.go`
- [x] T017 [P] [US2] Replace CLI explain refusal coverage and add task add/edit help and output cases in `internal/cli/cron_test.go` and `internal/cli/task_test.go`
- [x] T018 [US2] Update CLI wording and event explanation rendering in `internal/cli/cron.go` and `internal/cli/task.go`
- [x] T019 [P] [US2] Add failing desktop validation, preview, submission, and prefill tests for both syntaxes in `gui/editor_test.go` and `gui/editor_prefill_test.go`
- [x] T020 [US2] Make the desktop recurring-mode Schedule field prefill event schedules and document startup syntax in `gui/editor.go`

**Checkpoint**: Every supported authoring surface exposes one consistent startup event contract.

---

## Phase 5: User Story 3 - Cron import/export fidelity (Priority: P2)

**Goal**: Import, export, and convert startup jobs with existing operational-context guarantees.

**Independent Test**: Dry-run mutates nothing, real user/system import retains context and deduplicates normally, faithful export emits `@reboot`, and context-bearing tasks remain refused.

- [x] T021 [P] [US3] Replace crontab refusal expectations with supported user/system job and context cases in `internal/cron/crontab_test.go` and `internal/cli/cron_test.go`
- [x] T022 [US3] Preserve startup source expression through crontab scan and existing import request construction in `internal/cron/crontab.go` and `internal/cli/cron.go`
- [x] T023 [P] [US3] Add failing startup export and round-trip cases plus operational-context refusals in `internal/cron/export_test.go` and `internal/cron/convert_test.go`
- [x] T024 [US3] Export startup event schedules as `@reboot` while preserving existing task-context refusal boundaries in `internal/cron/export.go`

**Checkpoint**: `@reboot` is faithful cron interoperability, not merely accepted task syntax.

---

## Phase 6: Documentation, verification, and delivery

**Purpose**: Publish the behavior contract, run all gates, and leave a review-ready local commit.

- [x] T025 [P] Update supported syntax and daemon-start caveats in `README.md`, `docs/cli.md`, `docs/gui-fields.md`, and `docs/cron.md`
- [x] T026 [P] Update installed-path validation guidance and master contract traceability in `docs/test-scripts.md`, `test/scripts/README.md`, and `specs/001-task-scheduler/contracts/cli.md`
- [x] T027 Add the feature and architecture decision to `CHANGELOG.md`
- [x] T028 Run focused package tests and controlled lifecycle validation described in `specs/031-reboot-startup-trigger/quickstart.md`
- [x] T029 Advance the spec to In Progress, run `sh scripts/verify.sh all`, and capture all eight gate results in `specs/031-reboot-startup-trigger/verification.md`
- [x] T030 Mark all tasks complete, advance S031 to Implemented with local review-branch evidence in `specs/031-reboot-startup-trigger/spec.md` and `specs/README.md`, perform UTF-8/mojibake and diff audits, and commit locally

---

## Dependencies and execution order

- Phase 1 gates all implementation.
- Phase 2 establishes the shared model and blocks the three user stories.
- User Story 1 owns lifecycle semantics; User Story 2 and User Story 3 can then proceed independently around the shared parser and event model.
- Phase 6 depends on all story checkpoints.
- Within every story, failing tests precede implementation.

## Parallel opportunities

- Domain, schedule, API, CLI, GUI, import/export, and documentation test files marked `[P]` can be authored independently after their phase prerequisite.
- Implementation remains sequential where packages share parser or engine files.
- This autopilot run uses one working agent because repository instructions do not authorize subagent delegation; `[P]` records dependency structure rather than active delegation.

## Implementation strategy

Deliver the full slice as one review unit. The minimal internal model and lifecycle behavior land first, then every authoring surface, then cron interoperability and documentation. No partial publication checkpoint is introduced because the user requested one substantial end-to-end slice.
