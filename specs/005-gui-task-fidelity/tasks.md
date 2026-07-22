---

description: "Task list for 005-gui-task-fidelity"
---

# Tasks: GUI task fidelity — schedule round-trip and group assignment

**Input**: Design documents from `/specs/005-gui-task-fidelity/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md),
[data-model.md](data-model.md), [contracts/task-update.md](contracts/task-update.md)

**Tests**: REQUIRED, not optional. Constitution principle II is non-negotiable
and FR-023/FR-024 mandate a regression test per fixed defect plus a
migration-survival test. Every test task is written and observed **failing**
before its implementation task.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 (schedule fidelity), US2 (assign to group), US3 (ungroup),
  US4 (hierarchy membership), or FOUND (shared prerequisite)

## Path Conventions

Single Go module at repository root. Production code in `internal/**` and
`gui/**`; tests live beside the code they cover (`*_test.go`), per existing
project convention.

---

## Phase 1: Setup

No setup tasks. The module, toolchain, linter, and CI parity commands already
exist and are unchanged by this feature (see plan.md → Pinned artifacts).

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Persist and recover the schedule phrase, and make the group
tri-state expressible. Every user story depends on some part of this.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

### Storage and domain

- [x] T001 [FOUND] Add `Expression string \`json:"expression,omitempty"\`` to
      `domain.Schedule` in `internal/domain/domain.go`, with a doc comment
      stating it is the operator's phrase, is re-parseable, and is **inert with
      respect to execution** (FR-011a, data-model.md).
- [x] T002 [FOUND] Write the migration-survival test in
      `internal/store/store_test.go` **first, and observe it fail**: build a
      database at schema v3 (apply only migrations 1–3), insert schedule rows of
      each kind, then run the full migration and assert (a) the `expression`
      column exists and defaults to `''`, (b) every pre-existing schedule row is
      unchanged in every other column, (c) re-opening the store is a no-op and
      the version does not advance twice. This is the constitution's
      forward-only-migration safety surface (FR-002, FR-024, R4).
- [x] T003 [FOUND] Append migration **v4** to the `migrations` slice in
      `internal/store/store.go`:
      `ALTER TABLE schedules ADD COLUMN expression TEXT NOT NULL DEFAULT '';`
      Do not alter the migration runner. Makes T002 pass.
- [x] T004 [FOUND] Add `expression` to the schedule `INSERT` and `SELECT` column
      lists in `internal/store/crud.go` (~lines 168 and 180), and extend the
      round-trip assertions in `internal/store/store_test.go` to cover it.

### Phrase capture and recovery

- [x] T005 [FOUND] Set `sch.Expression = strings.TrimSpace(input)` in `finish()`
      in `internal/schedule/parse.go` — one place covers all four parse branches.
      Extend `parse_test.go` to assert the phrase is captured. `NewOneOff` in
      `recur.go` is deliberately left alone (data-model.md).
- [x] T006 [FOUND] Write `internal/schedule/render_test.go` **first, and observe
      it fail**: for every phrase in the existing `parse_test.go` table, assert
      `Parse → Render → Parse` yields an identical `RRULE` (FR-004); assert
      `Render` returns `""` for an unrecognized RRULE and for `event` kind
      (FR-003); assert **no** output ever contains `starting at` (FR-005, R3).
- [x] T007 [FOUND] Implement `schedule.Render(sch domain.Schedule, tzName string) string`
      in the new file `internal/schedule/render.go`, inverting the four parser
      shapes per the table in research.md R10. Exported doc comment states it is
      a best-effort inverse that returns `""` rather than guessing. Makes T006
      pass.

### Contract: tri-state group membership

- [x] T008 [FOUND] Write the group tri-state tests in
      `internal/api/server/update_test.go` **first, and observe them fail**:
      omitted `group_id` leaves membership unchanged; `""` clears it; a valid ID
      assigns it; an unknown ID returns 400 with code `validation` and field
      `group_id` (FR-014, FR-016, contracts/task-update.md).
- [x] T009 [FOUND] Change `TaskUpdateRequest.GroupID` to `*string` in
      `internal/api/server/update.go`, replace the `!= ""` test with the
      nil/empty/non-empty tri-state, and validate a non-empty ID via
      `store.GetGroup` returning `CodeValidation` on miss. **Single atomic task**
      — this is a compile-breaking type change, so update every in-repo call site
      in the same commit: `internal/cli/task.go`, `gui/editor.go`, and any
      server test constructing the struct. Makes T008 pass.
- [x] T010 [FOUND] Write a CLI test asserting `--group ""` ungroups while an
      omitted `--group` leaves membership unchanged, then implement it in
      `internal/cli/task.go` with `cmd.Flags().Changed("group")` (FR-015,
      contracts/task-update.md).

**Checkpoint**: the phrase persists and is recoverable, and all three group
intents are expressible end to end. User stories can begin.

---

## Phase 3: User Story 1 — See a task's real schedule when editing (P1) 🎯 MVP

**Goal**: The Edit dialog reflects the task's actual mode and timing, and an
untouched save changes nothing.

**Independent Test**: Create one recurring and one one-off task, reopen each,
confirm the dialog matches what was created, save untouched, confirm next-run
times are identical. No group work involved.

### Tests for User Story 1 ⚠️ write first, observe failing

> **Expected red state**: T013–T017 are written against the
> `*server.TaskResponse` editor signature that T021 introduces, so their first
> failure is a **compile** error, not an assertion failure. That is the correct
> red state here — do not "fix" it by reordering T021 ahead of the tests.

- [x] T011 [P] [US1] `internal/api/server/tasks_test.go`: `taskDetail` fills
      `Expression` from `Render` when the stored column is empty, and leaves a
      stored non-empty value alone (FR-003, R2).
- [x] T012 [P] [US1] `internal/api/server/update_test.go`: changing only the
      timezone and resubmitting the phrase re-anchors the recurrence in the new
      zone (FR-011). Note this requirement is satisfied *as a consequence of*
      T022 — once the editor prefills the phrase, a normal save resubmits it and
      the server re-parses in the new zone. The test pins that behavior so it
      cannot silently regress.
- [x] T013 [P] [US1] `gui/editor_test.go` (headless): prefill from a **one-off**
      `TaskResponse` selects Mode `One-off` and fills Date and Time expressed in
      the task's timezone (FR-006, FR-007).
- [x] T014 [P] [US1] `gui/editor_test.go`: prefill from a **recurring**
      `TaskResponse` selects Mode `Recurring`, fills Schedule from `Expression`,
      and splits a trailing `starting at HH:MM` into the Start-at field so
      `effectiveSchedule()` reconstructs the identical phrase without doubling
      the anchor (FR-006).
- [x] T015 [P] [US1] `gui/editor_test.go`: an editor opened on an existing task
      and not modified reports `isDirty() == false` (FR-008).
- [x] T016 [P] [US1] `gui/editor_test.go`: switching Mode on an existing task
      with the new mode's timing fields empty leaves Save disabled; not switching
      Mode still allows a blank schedule to mean "keep existing" (FR-011b, R11).
- [x] T017 [P] [US1] `gui/editor_test.go`: with an empty `Expression`, Schedule
      stays blank and the preview shows the schedule's `HumanSummary` (FR-010).
- [x] T017a [P] [US1] `internal/api/server/update_test.go`: create a task, then
      submit an update carrying **no** schedule and no `at` (the untouched-save
      shape), and assert the schedule ID and the computed next-run times are
      byte-identical before and after. This is SC-002 — the feature's central
      promise — and it must be verified automatically, not only in the manual
      walkthrough. *(analyze finding E3)*

### Implementation for User Story 1

- [x] T018 [US1] In `internal/api/server/tasks.go`, have `taskDetail` fill
      `sch.Expression = schedule.Render(sch, task.Timezone)` when it is empty —
      one place, covering create, get, and update responses. Never write it back
      to the store (R2). Makes T011 pass.
- [x] T019 [US1] Add `GetTask(ctx, id) (server.TaskResponse, error)` to the
      `Backend` interface in `gui/app.go` and to the fake in `gui/app_test.go`.
      **Do this before any other GUI task** or every `gui/` test fails to
      compile.
- [x] T020 [US1] In `gui/tasks.go`, make the Edit action fetch task detail via
      `GetTask` and pass the `*server.TaskResponse` to the editor. On failure,
      still open the dialog and surface the actionable message from research.md
      R7 rather than blocking the edit (FR-009).
- [x] T021 [US1] In `gui/editor.go`, change `showTaskEditor`, `newTaskEditor`,
      and the `taskEditor` struct to carry `*server.TaskResponse` instead of
      `*domain.Task`. Makes T013–T017 compile.
- [x] T022 [US1] Rewrite `prefill()` in `gui/editor.go`, replacing the hardcoded
      `e.mode.SetSelected(modeRecurring)` block (currently lines 349–351): select
      the mode from `Schedule.Kind`; for one-off format `*RunAt` in the task's
      timezone into the date and time entries; for recurring set Schedule from
      `Expression`, splitting any anchor clause into Start-at. Take
      `e.baseline = e.snapshot()` **after** prefill (FR-006, FR-007, FR-008).
- [x] T023 [US1] Scope the blank-keeps-existing allowance in `valid()` in
      `gui/editor.go` to "mode unchanged from the stored mode" rather than
      `existing == nil` (FR-011b, R11). Makes T016 pass.
- [x] T023a [US1] In `gui/editor.go`, when the prefilled `Expression` is empty
      but a schedule exists, seed the schedule-preview label with the schedule's
      `HumanSummary` instead of the "Type a schedule to see upcoming runs"
      placeholder, so the operator can still read what is currently set (FR-010).
      Makes T017 pass. *(analyze finding E1 — this requirement had a test but no
      implementation task)*

**Checkpoint**: US1 fully functional and independently testable. This alone is a
shippable fix for issue #4.

---

## Phase 4: User Story 2 — Assign a task to a group (P2)

**Goal**: A group can be chosen when creating or editing a task, and membership
is visible in the task list.

**Independent Test**: Create a group, assign a task to it from the editor,
confirm the task list shows it, disable the group and confirm the cascade
suppresses the task.

### Tests for User Story 2 ⚠️ write first, observe failing

- [x] T024 [P] [US2] `gui/editor_data_test.go`: the group choice helper emits
      `(none)` plus one entry per group, each showing its hierarchy path, and
      maps a chosen label back to the right group ID — including same-named
      groups at different levels (FR-012, FR-013, R8).
- [x] T025 [P] [US2] `gui/editor_test.go`: the editor prefills the group choice
      from the task's `GroupID`, an unresolvable ID prefills as `(none)` without
      error (FR-019a), and `buildForm()` round-trips the selection into the
      submitted request (FR-012).
- [x] T026 [P] [US2] `gui/tasks_test.go` (**new file**): the task list row shows
      the task's group name, and shows nothing group-like for an ungrouped task
      (FR-017).

### Implementation for User Story 2

- [x] T027 [US2] Add the group choice helpers to `gui/editor_data.go` beside the
      existing `overlapChoices`/`overlapValue` pattern: build the ordered label
      list (with `(none)`) from a `[]domain.Group` using `task.BuildForest` and
      `task.ByID`, and map label↔ID. **Single source for both the editor and the
      Groups tab** (R8). Makes T024 pass.
- [x] T028 [US2] Add a Group `widget.Select` to the "What to run" form in
      `gui/editor.go`, populated from `a.model.Snapshot().Groups` via T027 and
      prefilled from the task. Add it to `editorSnapshot` so dirty detection
      covers it, and to `taskForm` and `submitTask` so it reaches
      `TaskCreateRequest.GroupID` and the tri-state update field. Makes T025
      pass.
- [x] T029 [US2] Show the resolved group name in the task row in `gui/tasks.go`.
      Makes T026 pass.

**Checkpoint**: US1 and US2 both work independently.

---

## Phase 5: User Story 3 — Take a task back out of a group (P2)

**Goal**: A task can be returned to ungrouped from the GUI and is still visible
afterwards.

**Independent Test**: Assign a task to a group, remove it from all groups,
confirm it is listed as ungrouped and still runs.

### Tests for User Story 3 ⚠️ write first, observe failing

- [x] T030 [P] [US3] `gui/editor_test.go`: choosing `(none)` for a task that
      currently belongs to a group submits the **clear** intent (an empty-string
      pointer), not the unchanged intent (FR-014, Story 3).

### Implementation for User Story 3

- [x] T031 [US3] In `gui/editor.go`, map the `(none)` selection to the clear
      intent when the task previously had a group, and to the unchanged intent
      when it did not, so a no-op edit does not emit a pointless write. Makes
      T030 pass.

**Note**: the server and CLI halves of this story are T009 and T010 in the
foundational phase — the capability had to exist at the shared boundary before
any client could use it.

**Checkpoint**: assignment and its reverse both work from the GUI and the CLI.

---

## Phase 6: User Story 4 — See group membership in the hierarchy (P3)

**Goal**: The Groups tab shows member tasks, an always-present ungrouped area,
and a move action.

**Independent Test**: With tasks across two groups and one ungrouped, confirm
every task appears exactly once under the right parent, then move one and
confirm it relocates live.

### Tests for User Story 4 ⚠️ write first, observe failing

- [x] T032 [P] [US4] `gui/groups_test.go` (**new file**): the tree exposes member tasks under
      their group and ungrouped tasks under the ungrouped node; every task
      appears exactly once (FR-018, SC-005); group and task node identities never
      collide (Edge Cases).
- [x] T033 [P] [US4] `gui/groups_test.go`: the ungrouped node is present even
      when it has no members (FR-019).
- [x] T034 [P] [US4] `gui/groups_test.go`: the move action issues an
      `UpdateTask` carrying the expected tri-state value, including the clear
      case, against the fake backend (FR-020).
- [x] T035 [P] [US4] `gui/groups_test.go`: group-only actions (New Group,
      Enable/Disable, Delete) do not act when a task node is selected (FR-021).

### Implementation for User Story 4

- [x] T036 [US4] In `gui/groups.go`, prefix tree node IDs (`g:<id>` / `t:<id>`)
      so groups and tasks can share the tree without collision, and extend
      `childIDs` with each group's member tasks read from
      `a.model.Snapshot().Tasks`. Makes T032 pass.
- [x] T037 [US4] Add the synthetic always-present ungrouped node collecting tasks
      whose `GroupID` is empty or unresolvable (FR-019, FR-019a). Makes T033
      pass.
- [x] T038 [US4] Render task rows distinguishably from group rows (name + state)
      in the tree's update callback (FR-018).
- [x] T039 [US4] Add the **Move to group…** toolbar action, enabled only for a
      task selection, opening a form with the same choice list from T027
      (including `(none)`) and calling `UpdateTask` with the tri-state value.
      Makes T034 pass.
- [x] T040 [US4] Guard New Group, Enable/Disable, and Delete so a task selection
      does not trigger them. Makes T035 pass.

- [x] T040a [P] [US4] `gui/groups_test.go`: applying a `KindTask` event that
      changes a task's `GroupID` moves the task under the new group in the tree
      with no manual refresh, and applying a `KindTask` delete removes it from
      the tree (FR-022, SC-006). *(analyze finding E2 — FR-022 was asserted by
      inspection with no test. The Groups tab now reads `Tasks`, a dependency it
      did not previously have, so the assumption needs pinning rather than
      trusting.)*

**Note**: no new refresh *path* is needed — the existing `KindTask` fold in
`gui/viewmodel/viewmodel.go` plus `onModelChange` fanning out to every
registered refresher already delivers this. T040a exists to prove that, not to
build it.

**Checkpoint**: all four stories independently functional.

---

## Phase 7: Polish & Cross-Cutting

- [x] T041 [P] Add a `## Group` section to `docs/gui-fields.md` and note that
      Edit now prefills Mode, Schedule, and the one-off date and time.
- [x] T042 [P] Document `--group ""` as the ungroup form in
      `specs/001-task-scheduler/contracts/cli.md`.
- [x] T043 [P] Add a roadmap line under **Polish & cross-cutting** in `TODO.md`
      so the feature traces to the roadmap as the autopilot protocol requires.
- [x] T044 `CHANGELOG.md` — Unreleased **Fixed** entries for issues #3 and #4,
      plus dated **Decisions** entries for migration v4, the read-time renderer
      (R2), and the tri-state `group_id` contract change (R5).
- [x] T045 Confirm core-package coverage stays at or above 80% and no exported
      addition lacks a doc comment (principles I and II).
- [x] T046 Run the full CI-parity suite from [quickstart.md](quickstart.md) in
      the **foreground**, watched to completion — never backgrounded.
      *Status*: `gofmt` clean, `go vet` clean, `golangci-lint@v2.1.6` reports
      **0 issues**, full test suite passes. The one command not run is the
      `-race` suite: `-race` requires cgo and this machine has no C compiler
      (`cgo: C compiler "gcc" not found`). CI runs it as the merge gate. Both
      local traps are now documented in `CLAUDE.md`.
- [x] T046a Re-run the `/speckit-analyze` coverage check: every FR and every
      buildable SC maps to at least one task. Findings E1, E2, E3, F1, and C1
      from the 2026-07-22 analysis are resolved by T023a, T040a, T017a, the
      expected-red-state note in Phase 3, and the new-file markers respectively.
- [ ] T047 Run the manual GUI walkthrough in [quickstart.md](quickstart.md),
      including section A against a **pre-upgrade database**, which is the only
      check that proves the fix reaches the people who reported it.

---

## Dependencies & Execution Order

### Phase dependencies

- **Phase 2 (Foundational)** blocks everything. Within it:
  T001 → T002 → T003 → T004 (domain field before storage before persistence);
  T005 needs T001; T006 → T007; T008 → T009 → T010.
  T005–T007 are independent of T008–T010 and may proceed in parallel.
- **US1 (Phase 3)** needs T001–T007. **US2 (Phase 4)** needs T009.
  **US3 (Phase 5)** needs T009 and T028. **US4 (Phase 6)** needs T027.
- **Phase 7** needs every story intended for this release.

### Hard sequencing constraints

1. **T009 is atomic.** The `*string` change breaks compilation across the
   server, CLI, GUI, and tests; splitting it leaves the tree unbuildable.
2. **T019 precedes every other GUI task.** Without `GetTask` on the interface and
   the fake, the whole `gui/` package fails to compile.
3. **T027 precedes T028 and T039.** Both consume the single choice-list source;
   building either list independently reintroduces the drift CHK045 flagged.
4. **Every test task precedes its implementation task and must be observed
   failing first** (principle II, FR-023).

### Parallel opportunities

- T005–T007 ‖ T008–T010 (schedule package vs API/CLI)
- T011–T017 (all US1 tests, distinct concerns) once Phase 2 completes
- T024–T026, T032–T035 within their stories
- T041–T043 (three different documents)

---

## Implementation Strategy

**MVP first**: Phase 2 → Phase 3 → validate US1 independently. That closes issue
#4, the more severe defect, and is shippable on its own.

**Then incrementally**: US2 → US3 → US4, each validated independently, together
closing issue #3.

Single-operator execution, so the parallel-team strategy in the template does not
apply; the `[P]` markers indicate only that tasks touch disjoint files.

---

## Notes

- `[P]` = different files, no dependencies
- Verify every test fails before implementing against it
- No pinned artifact is modified by any task in this list
- The single autopilot halt comes after T046 and the commit, before any push
