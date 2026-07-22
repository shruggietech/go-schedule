# Phase 0 Research: GUI task fidelity

**Feature**: [spec.md](spec.md) · **Date**: 2026-07-22

This resolves every unknown in the plan's Technical Context and closes the five
checklist items left open by `/speckit-checklist` (CHK004, CHK005, CHK006,
CHK044, CHK045). Decisions were made under the Build-Phase Autopilot Protocol
decision policy and are binding on `tasks.md`.

## R1 — How the operator's schedule phrase is recovered

**Decision**: Persist the phrase on the schedule record (new `expression`
column, migration v4). Nothing reconstructs phrases for schedules stored before
this feature.

**Rationale**: The stored `human_summary` cannot serve as the recovery path — it
is not re-parseable. Confirmed by reading the parsers in
[`internal/schedule/parse.go`](../../internal/schedule/parse.go):

| Stored `human_summary` | Re-parses? | Why |
|---|---|---|
| `Every day at 09:00` | yes | `reInterval` accepts it case-insensitively |
| `Every weekday at 09:00` | **no** | `reDayset` requires the literal `weekdays` |
| `Every weekend day` | **no** | `reDayset` requires `weekends` |
| `The 3rd Wednesday of every month at 14:00` | **no** | `reOrdinal` requires a leading ordinal |

So the phrase has to be retained at creation time. `schedule.Parse` is the only
producer of recurring schedules (verified: the sole call sites are the create,
update, and preview handlers), so retaining it there covers every recurring
schedule the system can make.

**Superseded — recovery for pre-existing rows (2026-07-22).** This decision
originally also specified `schedule.Render`, an RRULE→phrase inverse applied at
read time to serve schedules stored before the `expression` column existed. That
was built on a wrong premise: the issues were filed against v0.3.0, and the
design inferred an installed base that must not be broken. There is none — the
software has no working deployments, and the only databases in existence are the
maintainers' own, all non-functional. The renderer (~180 lines plus its
round-trip test suite) served exclusively that phantom population, so it and the
read-time backfill were removed. A database predating the column now shows a
blank schedule field on edit, which means "keep unchanged" and is harmless.

**Alternatives considered**:

- *Re-render `human_summary` into parseable wording*: rejected — changes an
  existing user-visible string for every task and still would not recover the
  operator's own words.
- *Reconstruct from the RRULE at read time*: implemented, then removed. See the
  superseded note above.

## R4 — Migration safety *(closes CHK005 and CHK006)*

**Decision**: No new rollback or idempotency machinery. Both properties are
already provided by the existing migration runner and will be pinned by tests
rather than re-implemented.

**Rationale**: reading `migrate()` at
[`internal/store/store.go:156`](../../internal/store/store.go:156):

- **Partial-failure rollback (CHK005)**: each migration runs inside its own
  transaction; any error triggers `tx.Rollback()` before returning, and the
  version row is written in the same transaction as the schema change. A failed
  v4 therefore leaves the database at v3 exactly, with no half-applied state.
- **Re-run on an upgraded installation (CHK006)**: the runner reads
  `COALESCE(MAX(version), 0)` and skips every migration with
  `m.version <= current`. Re-opening an upgraded database is a no-op.
- **Forward-only and non-destructive (FR-002)**: `ALTER TABLE schedules ADD
  COLUMN expression TEXT NOT NULL DEFAULT ''` adds one column with a total
  default. No existing column, row, or value is read, rewritten, or dropped, so
  no stored timing can move.

**Consequence for `tasks.md`**: a test that opens a v3-shaped database, migrates,
and asserts (a) the column exists, (b) every pre-existing schedule row is
byte-identical apart from the new default, and (c) re-opening is a no-op. This
is the constitution's forward-only-migration safety surface; it is mandatory.

## R5 — Expressing "remove from all groups" *(FR-014)*

**Decision**: `TaskUpdateRequest.GroupID` becomes `*string`. `nil` = unchanged,
`""` = clear, non-empty = assign. The CLI decides which to send with
`cmd.Flags().Changed("group")`.

**Rationale**: This is the established convention in this codebase, not a new
one — `GroupUpdateRequest.Parent` at
[`internal/api/server/groups.go:22`](../../internal/api/server/groups.go:22) is
already `*string` for exactly this reason (distinguishing "don't reparent" from
"make top-level"). Reusing it keeps one idiom for tri-state optional fields.
`Changed()` preserves the CLI's current contract — omitting `--group` still
means "unchanged" (FR-015) — while making `--group ""` mean "ungroup".

**Alternatives considered**: a sentinel string such as `"-"` or `"none"`
(rejected: collides with a legal group ID and invents a second idiom); a
separate `--ungroup` boolean flag plus an `Ungroup bool` field (rejected: two
fields encoding one tri-state, and they can contradict each other).

## R6 — Validating a group assignment *(FR-016)*

**Decision**: The update handler resolves a non-empty group ID via
`store.GetGroup` before assigning, and returns HTTP 400 `CodeValidation` naming
the `group_id` field when it does not exist.

**Rationale**: Without the check, a bad ID reaches the foreign key and surfaces
as an opaque 500, violating principle III (errors state what failed and what to
do). The existing `groupErr` helper already maps store errors to validation
responses for the group endpoints; this follows it.

## R7 — What the editor shows when the detail lookup fails *(closes CHK044)*

**Decision**: The editor opens with schedule fields blank and the preview pane
showing `Could not read this task's current schedule — leave Schedule blank to
keep it unchanged.` The dialog is never blocked.

**Rationale**: Principle III requires the message to say what failed and what
the operator can do; naming the safe action ("leave it blank") is the actionable
half that a bare failure notice would omit. Blocking the dialog would make a
transient read failure prevent an unrelated edit (renaming a task), which is a
worse outcome than editing with one field unpopulated — and the blank-keeps-
existing rule means the schedule is still safe.

## R8 — One source of group choices *(closes CHK045)*

**Decision**: A single helper builds the group choice list (`(none)` plus every
group as an indented hierarchy path) from the view-model snapshot. The task
editor and the Groups-tab move action both call it; neither builds its own list.

**Rationale**: Two independently constructed lists would drift — FR-012 and
FR-020 must offer the same destinations, including `(none)`, or the two paths to
the same operation disagree. Placing it beside the existing
`overlapChoices`/`overlapValue` label↔value helpers in
[`gui/editor_data.go`](../../gui/editor_data.go) matches how this codebase
already separates presentation data from widget wiring.

## R9 — Fetching task detail in the GUI

**Decision**: Add `GetTask` to the GUI's `Backend` interface and call it from the
Tasks tab's Edit action, passing the resulting `server.TaskResponse` to the
editor. The view-model's cached `[]domain.Task` is not extended to carry
schedules.

**Rationale**: `client.GetTask` already exists
([`internal/api/client/methods.go:39`](../../internal/api/client/methods.go:39));
only the GUI's interface omits it. Caching a schedule per task in the view-model
would add a second staleness surface that the event stream would have to keep
coherent, for data needed only while one dialog is open. A single fetch at open
is simpler and always current.

## R11 — Mode switching on an existing task *(FR-011b)*

**Decision**: The editor's validity check treats "creating" as
`existing == nil || mode changed from the stored mode`. The blank-keeps-existing
allowance applies only when the selected mode still matches the stored one.

**Rationale**: Today `valid()` keys the allowance off `e.existing == nil` alone
([`gui/editor.go:464`](../../gui/editor.go:464)). Once the editor knows the
stored mode (R9), the allowance can be scoped correctly: leaving the schedule
alone is only meaningful when the schedule being kept is of the kind the
operator is looking at. This closes a hole this feature would otherwise open —
switching to One-off and saving with empty date/time would silently keep a
recurring schedule on a task the operator believes is now one-off.
