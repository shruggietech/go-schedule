# Implementation Plan: GUI task fidelity — schedule round-trip and group assignment

**Branch**: `005-gui-task-fidelity` | **Date**: 2026-07-22 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/005-gui-task-fidelity/spec.md`

## Summary

Two shipped defects ([#4](https://github.com/shruggietech/go-schedule/issues/4),
[#3](https://github.com/shruggietech/go-schedule/issues/3)) share one cause: the
GUI cannot see data the daemon already holds.

**Schedule fidelity.** The task editor hardcodes Mode to *Recurring* and leaves
the timing fields blank because the cached task row carries no schedule and the
stored `human_summary` is not re-parseable. Fix: persist the operator's phrase on
the schedule (`expression`, migration v4) and have the editor fetch task detail
and prefill from it.

**Group assignment.** Groups cannot be populated from the GUI at all, and no
client can un-group a task because an empty group value means "unchanged". Fix:
make `TaskUpdateRequest.GroupID` a `*string` tri-state, add a group selector to
the editor, show member tasks and an always-present ungrouped node in the Groups
tab with a move action, and show the group in the task list.

Design decisions and their rationale: [research.md](research.md).

## Technical Context

**Language/Version**: Go 1.25.0 (`go.mod`)

**Primary Dependencies**: `fyne.io/fyne/v2` v2.7.4 and `fyne.io/x/fyne` (GUI),
`github.com/teambition/rrule-go` v1.8.2 (recurrence), `modernc.org/sqlite`
v1.52.0 (pure-Go driver), `github.com/spf13/cobra` (CLI). No new dependency.

**Storage**: embedded SQLite via `internal/store`; UTC RFC 3339 timestamps;
ordered forward-only migrations (currently at v3, this feature adds v4)

**Testing**: `go test -race`; headless Fyne test driver for `gui/`; existing
table-driven parser tests in `internal/schedule`

**Target Platform**: Linux, macOS, Windows; GUI built windowless (`-H windowsgui`)

**Project Type**: desktop application + system daemon + CLI over a local IPC API
(Unix socket / Windows named pipe)

**Performance Goals**: unchanged — dispatch latency p99 < 100ms. This feature
adds nothing to the dispatch path.

**Constraints**: persisted schedules migrate forward non-destructively and no
stored timing may move (FR-002/SC-007); the retained phrase is inert with
respect to execution (FR-011a); GUI widget construction stays cgo-free so it
remains headlessly testable

**Scale/Scope**: single local operator; tens to low hundreds of tasks; ~10 source
files changed across `internal/{domain,store,schedule,api,cli}` and `gui/`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

Constitution v1.1.0.

| Principle | Gate | Initial | Post-design |
|---|---|---|---|
| I. Code Quality | `gofmt`/`go vet`/`golangci-lint` clean; doc comments on exported items; wrapped errors; documented goroutine lifecycles | PASS — the new `Expression` field carries a doc comment stating its contract; no new goroutines | PASS |
| II. Testing (NON-NEGOTIABLE) | tests alongside; regression test failing before each fix; injected clock; `-race`; ≥80% coverage on core packages | PASS — every FR maps to a named test in [quickstart.md](quickstart.md); migration-survival test is mandatory | PASS |
| III. UX Consistency | consistent CLI verb-noun and flags; actionable errors naming the field; explicit timezone handling | PASS — `--group ""` extends an existing flag rather than adding one; unknown group returns a field-named validation error (R6); lookup-failure message names the safe action (R7) | PASS |
| IV. Performance | measured dispatch budget; no regression >10%; no leaks | PASS — nothing added to the dispatch path | PASS |
| V. Autonomous Build-Phase Execution | spec'd through spec-kit; `analyze` not skipped; one halt before push | PASS — full sequence run; decisions recorded here and in research.md | PASS |

### Safety-critical surfaces touched

CLAUDE.md names surfaces that may never be weakened. This feature touches one:

- **Forward-only non-destructive store migrations** — migration v4. Covered by
  R4 and a mandatory test asserting a v3 database upgrades with every schedule
  row intact and re-opens as a no-op.

Not touched: clock injection, timezone/DST resolution (FR-011 exercises existing
resolution, adds none), restart/catch-up recovery, goroutine termination, local
IPC access control.

### Pinned artifacts

None modified. `.github/workflows/**`, `build/**`, `Makefile`, `.golangci.yml`,
the `go`/`toolchain` lines of `go.mod`, `.gitattributes`, `.gitignore`,
`LICENSE`, and `docs/INSTALL-windows.md` are all untouched, so no pinned-artifact
decision entry is required. `CHANGELOG.md` still gets a dated decision entry for
the migration and the contract change.

**Result: PASS. No violations, no Complexity Tracking entries.**

## Project Structure

### Documentation (this feature)

```text
specs/005-gui-task-fidelity/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0 — R1–R11 decisions
├── data-model.md        # Phase 1 — Schedule.Expression, migration v4
├── quickstart.md        # Phase 1 — validation guide
├── contracts/
│   └── task-update.md   # Phase 1 — tri-state group_id, expression on task detail
├── checklists/
│   ├── requirements.md  # spec quality (16/16)
│   ├── migration.md     # persisted-state quality (27/30 → gaps closed in research.md)
│   └── ux.md            # UX requirements quality (25/27 → gaps closed in research.md)
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (repository root)

```text
internal/
├── domain/domain.go              # + Schedule.Expression
├── store/
│   ├── store.go                  # + migration v4 (ALTER TABLE schedules ADD COLUMN expression)
│   └── crud.go                   # + expression in schedule INSERT/SELECT
├── schedule/parse.go             # Parse sets Expression in finish()
├── api/server/update.go          # GroupID *string tri-state + group existence validation
└── cli/task.go                   # --group tri-state via Flags().Changed

gui/
├── app.go                        # Backend interface + GetTask
├── tasks.go                      # Edit fetches detail; list row shows group
├── editor.go                     # takes *server.TaskResponse; prefill; group select; mode-switch validity
├── editor_data.go                # group choice helpers (single source, R8)
└── groups.go                     # tree shows member tasks + ungrouped node + Move to group…
```

**Structure Decision**: The existing layout is unchanged — this is a defect fix
inside an established architecture, not a new component, and it adds no new
source file. Group choice-list helpers go in `gui/editor_data.go`, which already
exists to keep presentation data out of widget wiring.

## Phase 2 handoff

`/speckit-tasks` derives `tasks.md` from this plan. Ordering constraints it must
respect:

1. `domain` → `store` → `schedule` before any API change (the `Expression` field
   must exist and persist before anything reads it).
2. Reserved (an earlier revision derived phrases for pre-existing rows; removed).
3. The `GroupID *string` change is a compile-breaking edit: server, CLI, GUI, and
   server tests must land in one task, not spread across several.
4. GUI work depends on `GetTask` being on the `Backend` interface and the fake
   backend being updated first, or every `gui/` test fails to compile.
5. Under principle II each fix's regression test is written first and observed
   failing before its fix.

## Complexity Tracking

No constitution violations. Section intentionally empty.
