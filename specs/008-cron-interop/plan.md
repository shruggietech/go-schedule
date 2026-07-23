# Implementation Plan: Cron interoperability and calendar-anomaly policy

**Branch**: `008-cron-interop` (trunk-based — committed directly onto `main`) |
**Date**: 2026-07-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/008-cron-interop/spec.md`

## Summary

Expose the cron equivalence the project already asserts as a user-facing
conversion surface on the CLI (`gosched cron explain|import|export`), and close
the grammar and calendar-policy gaps that make that conversion possible without
silent meaning changes.

Three pieces, in dependency order:

1. **Grammar** — two new parse arms in `internal/schedule/parse.go` for by-date
   monthly and yearly rules, plus month/year interval units. Without these,
   ordinary cron lines like `0 9 1 * *` have no target representation.
2. **Missing-date policy** — a third per-task policy (`skip` / `last_valid` /
   `next_valid`) resolved inside `internal/schedule` before the existing DST
   normalization, defaulting to today's behavior, with an additive store
   migration and honest human summaries.
3. **Conversion** — a new dependency-free `internal/cron` package that maps a
   cron expression to the *phrase a user would have typed* and back, and a
   `internal/cli/cron.go` command surface over it. Everything the mapping cannot
   carry is refused by name.

Plus issue #9: hyperlink the ShruggieTech attribution to `https://shruggie.tech`
across the project-facing documents.

## Technical Context

**Language/Version**: Go 1.25.0 (`go.mod`), `GOTOOLCHAIN=auto`

**Primary Dependencies**: `github.com/teambition/rrule-go` (recurrence),
`github.com/spf13/cobra` (CLI), `modernc.org/sqlite` (store), `fyne.io/fyne/v2`
(GUI). **No new dependency is added** — the cron parser is written in-tree
(research D3).

**Storage**: embedded SQLite, forward-only numbered migrations in
`internal/store/store.go`. This feature adds migration **v5**.

**Testing**: `go test -race`; table-driven unit tests in each package;
`test/integration` for cross-package behavior; `scripts/coverage-gate.sh` for the
≥80% core-package gate.

**Target Platform**: Linux, macOS, Windows — daemon + CLI cgo-free; GUI cgo.

**Project Type**: single Go module — daemon, CLI, and GUI over a shared internal
tree.

**Performance Goals**: unchanged. p99 dispatch latency < 100 ms. The policy
resolution runs on the next-run path, so it must stay allocation-light and must
not turn an O(1) computation into an unbounded search.

**Constraints**: internal scheduling in UTC; injected `Clock` (no direct
`time.Now()` in engine code); no stored timing value may change.

**Scale/Scope**: ~9 packages touched, 1 new package, 1 migration, 1 new CLI
command group with 3 subcommands, 1 new GUI form row.

## Constitution Check

*GATE: passed before Phase 0, re-checked after Phase 1 design.*

| Principle | Assessment |
| --- | --- |
| **I. Code Quality** | New package carries doc comments stating intent; refusals are values, not panics; errors wrapped with `%w`. The cron field parser is decomposed per field rather than one large function, to stay under the linter's complexity gate. |
| **II. Testing (NON-NEGOTIABLE)** | Safety-critical surfaces touched: **timezone/DST resolution** and **forward-only migrations**. Both gain tests rather than losing any: policy resolution is tested at real DST transitions, and migration v5 gets a v4→v5 test asserting no stored value changes. Clock stays injected — every new test passes explicit instants. The existing `cronparity_test.go` is extended, never weakened. |
| **III. UX Consistency** | New CLI follows verb-noun (`cron explain`), honors `--json`, writes results to stdout and refusals to stderr where they are diagnostics, and uses the established exit codes. The new flag's values use underscores to match `--overlap`/`--catchup` (contracts/cli.md). Times remain RFC 3339 in machine output. |
| **IV. Performance** | Policy resolution is a bounded walk over at most a handful of candidate periods, not a search; the default path (`skip`) is unchanged code. `BenchmarkNextRun` is kept and re-run to confirm no regression beyond the 10% bar. |
| **V. Autonomous execution** | This feature is spec'd through spec-kit, `/speckit-analyze` runs before implementation, and there is exactly one halt before the push. |

No violations. **Complexity Tracking is therefore empty and omitted.**

One item is worth surfacing at the halt rather than hiding: adding a parameter to
`schedule.NextRun` and `schedule.UpcomingRuns` touches six call sites across the
engine, catch-up, and API packages. This is deliberate (research D1) and is
compile-checked; the alternative stored the policy where an existing code path
would silently reset it.

## Project Structure

### Documentation (this feature)

```text
specs/008-cron-interop/
├── plan.md              # This file
├── spec.md              # Feature specification (+ Clarifications session)
├── research.md          # Phase 0: decisions D1–D6 and rejected alternatives
├── data-model.md        # Phase 1: entities, migration v5, signature changes
├── quickstart.md        # Phase 1: end-to-end verification walkthrough
├── contracts/
│   └── cli.md           # Phase 1: the `gosched cron` command contract
├── checklists/
│   ├── requirements.md  # Spec quality (16/16)
│   ├── fidelity.md      # Conversion fidelity requirements quality (CHK001–030)
│   └── calendar.md      # Calendar correctness requirements quality (CHK031–059)
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (repository root)

```text
internal/
├── cron/                     # NEW — dependency-free cron interchange
│   ├── cron.go               #   Spec, Parse, field grammar
│   ├── phrase.go             #   Spec → human phrase (or a named refusal)
│   ├── export.go             #   domain.Schedule → crontab line (or refusal)
│   └── *_test.go
├── schedule/
│   ├── parse.go              # + by-date monthly, yearly, month/year units
│   ├── recur.go              # + missing-date resolution before DST normalization
│   ├── missingdate.go        # NEW — period walk and date resolution
│   └── *_test.go             # + cronparity corpus, policy tests at real dates
├── domain/domain.go          # + MissingDatePolicy enum and Task field
├── store/
│   ├── store.go              # + migration v5 (tasks.missing_date_policy)
│   ├── crud.go               # + column in task insert/scan
│   └── migration_v5_test.go  # NEW
├── engine/engine.go          # NextRun call sites carry the policy
├── catchup/catchup.go        # Evaluate carries the policy
├── api/server/
│   ├── tasks.go, update.go   # + missing_date_policy request field + validation
│   └── calendar.go           # UpcomingRuns call site
└── cli/
    ├── cron.go               # NEW — explain / import / export
    ├── cli.go                # + newCronCmd() in newRoot()
    └── task.go               # + --missing-date on add/edit; shown in show

gui/
├── editor.go                 # + Missing dates row in Advanced Settings
└── editor_data.go            # + missingDateChoices via existing policyChoice[T]

docs/
└── cron.md                   # NEW — conversion guide + fidelity table
```

**Structure Decision**: the existing single-module layout is kept. The one new
package, `internal/cron`, is justified by boundary isolation: cron is an
interchange format that must not leak into the scheduling engine, and giving it
its own package makes that a compile-time fact — `internal/schedule` does not
import `internal/cron`, only the reverse.

## Implementation approach

### Phase A — grammar (`internal/schedule/parse.go`)

Two regexes and two parse arms, registered in `Parse` before `parseInterval`:

- `reByDate` — `on the <n>(st|nd|rd|th) of (the|each|every) month`, and the
  shorter `the <n>th monthly`, with the shared optional `at <time>` tail →
  `FREQ=MONTHLY;BYMONTHDAY=<n>`.
- `reYearly` — `every year on <month> <day>` → `FREQ=YEARLY;BYMONTH=<m>;
  BYMONTHDAY=<d>`.

`unitToFreq` gains `month`/`months`/`mo` → `MONTHLY` and `year`/`years`/`y` →
`YEARLY`, which makes `every 12 months` fall out of the existing interval arm
with no new code. Summaries follow the established `"Every "` / `"The "` phrasing
and are extended by the policy clause in Phase B.

### Phase B — missing-date policy

`internal/schedule/missingdate.go` holds the resolution, kept out of `recur.go`
so the DST path stays readable:

- **Applicability** — a rule is *date-bearing* when its RRULE carries
  `BYMONTHDAY`, or `BYDAY` with an ordinal prefix (`+5FR`). Anything else is
  inert, satisfying FR-024 by construction.
- **`skip`** — return the rrule-go occurrence unchanged. This is literally the
  current code path, which is how FR-020's "bit for bit" is guaranteed rather
  than asserted.
- **`last_valid` / `next_valid`** — walk candidate periods forward from the
  anchor: for each month (or year), compute the target date; if it exists, use
  it; if not, clamp to the period's last valid day, or roll to the first day of
  the following period. Return the first result strictly after `after`. The walk
  is bounded (a hard cap on periods examined, returning "no further run" rather
  than spinning) so a pathological rule cannot stall the dispatch loop —
  principle IV.
- The resolved instant then goes through the **unchanged** `timezone.WallTime`
  normalization, so DST behavior is inherited, not re-implemented.

`HumanSummary` gains a policy clause for date-bearing rules only: "…, or the last
valid date when the month has none", "…, skipped in months that have none". This
is what discharges FR-023 and the false label issue #8 documents.

### Phase C — plumbing

Domain enum → migration v5 → store CRUD → API request validation (copying the
`overlap_policy` block verbatim) → CLI flag and `task show` line → GUI selector
through the existing generic `policyChoice[T]`. Each is a repetition of an
established path; the only judgment call is that update must leave the policy
alone when the field is empty (FR-024a).

### Phase D — conversion (`internal/cron`)

- `Parse` builds a `Spec` from five fields, resolving names, lists, ranges and
  steps, and refusing six-field input, `@reboot`, and `L`/`W`/`#` by name.
- `Phrase` maps a `Spec` to a phrase, refusing anything with no phrase — a
  non-dividing step, a DOM+DOW combination, or a field pattern outside the
  grammar. The phrase then goes through `schedule.Parse`; there is no second
  route.
- `Export` inverts a stored schedule into a line, or refuses.
- The parity corpus in `cronparity_test.go` is reused as the converter's table,
  extended with the by-date and yearly cases the new grammar makes reachable, and
  strengthened into a round trip across a DST transition and a month boundary.

### Phase E — CLI, docs, and #9

`internal/cli/cron.go` implements the contract in `contracts/cli.md`. `docs/cron.md`
carries the fidelity table. `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, and
`CODE_OF_CONDUCT.md` get the attribution link. `CHANGELOG.md` records the feature
and the dated decisions (D1 and D3 are the architecture-affecting ones).

## Verification

CI parity via the `go-schedule-verify` skill, foreground, watched to completion:
`gofmt` (silent), `go vet`, `golangci-lint`, `go test -race` on the GUI-excluded
set, `go test ./gui/...`, and `scripts/coverage-gate.sh`. Then the end-to-end
walkthrough in [quickstart.md](./quickstart.md) against built binaries.

A red gate is a halt in itself. If the machine lacks a C toolchain, the race gate
is reported as not run rather than as passing.
