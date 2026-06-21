# Implementation Plan: Rebrand to go-schedule + GUI & Installer Overhaul

**Branch**: `004-rebrand-gui-overhaul` | **Date**: 2026-06-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/004-rebrand-gui-overhaul/spec.md`

## Summary

Six changes, sequenced so the foundational rename lands first and the riskier backend work
precedes the GUI work that depends on it:

1. **Rename** `go-scheduler` → `go-schedule` across the live tree: the Go module path
   (`github.com/shruggietech/go-scheduler` → `.../go-schedule`), build/release config, docs, and
   user-facing strings. History (CHANGELOG entries, prior `specs/00{1,2,3}`) is left intact.
2. **Windows MSI**: replace the portable zip with a WiX-built `.msi` that installs to
   `Program Files`, registers `goschedd` as an auto-start Windows service, adds Start-Menu
   shortcuts, and uninstalls cleanly. The zip/"run from folder" path is removed.
3. **Logs view**: replace the Alerts tab with a unified **Logs** view fed by a new daemon log
   pipeline — a custom `slog.Handler` that tees structured records to (a) a rotating on-disk
   file, (b) a bounded in-memory ring served by a new `GET /v1/logs`, and (c) the event broker
   for live streaming. Existing alerts are folded into the same stream. The view adds
   severity/type filters, click-through detail, and Dismiss All.
4. **Remove Triggers** entirely: delete the GUI tab, CLI commands, API routes, client methods,
   engine dispatcher wiring, store layer, and domain types; add store **migration v3** to drop
   the `triggers` and `dedup_ledger` tables.
5. **Real-time everywhere**: extend the broker with task/group change events, publish them from
   the API mutation handlers, fold them into the view-model, and **remove every manual Refresh
   control**.
6. **Calendar view**: add a toggleable month-grid calendar under Schedule over the existing
   `GET /v1/calendar` data (backend already exists — this is GUI-only).

The heavy lifting is the **log pipeline** (new, dependency-free) and the **MSI** (new build
tooling). Triggers removal is broad but mechanical. The calendar and real-time work reuse
existing infrastructure (the calendar API and the SSE broker already exist).

## Technical Context

**Language/Version**: Go (project toolchain, `go.mod`); GUI via Fyne v2.7.4.

**Primary Dependencies**:
- Existing: Fyne v2, `kardianos/service` (service control), `modernc.org/sqlite` (cgo-free),
  `rrule-go`, stdlib `log/slog`, `net/http`.
- **New build-time only (not a Go module dependency)**: **WiX Toolset v5** to compile the `.msi`,
  run on the Windows CI runner. No new *runtime* Go dependency is introduced.
- Log rotation is implemented **in-house** (a small size-based rotating `io.Writer`) rather than
  adding `lumberjack`, to honor the constitution's "stdlib where it suffices" constraint. The
  daemon is the single log writer, so a simple rotator is sufficient and testable.

**Storage**: SQLite (`modernc.org/sqlite`), single-writer. New **migration v3** drops `triggers`
and `dedup_ledger`. The on-disk log file is a rotating **JSONL** file under the data dir
(`logs/goschedule.log` + rotated siblings); alerts remain in the existing `alerts` table. Recent
log records also live in a bounded in-memory ring (no new hot table).

**Testing**: `go test -race ./...`. New unit tests: the log handler/rotator (pure, no wall-clock),
the ring buffer, the new broker event kinds, migration v3 (drops trigger tables; no-op on DBs
that never had them), view-model folding of task/group/log events, headless Fyne tests for the
Logs view (filter, detail, Dismiss All) and the calendar view (date placement, toggle preserves
window). MSI is validated by a CI build producing the artifact plus a manual install/uninstall
checklist in quickstart.

**Target Platform**: Linux, macOS, Windows for the daemon/CLI/GUI. The MSI is **Windows-only**;
Linux/macOS service install is unchanged except for rename touches.

**Project Type**: Multi-binary single Go module — daemon (`goschedd`), CLI (`gosched`), desktop
GUI (`gosched-gui`) over a local IPC API.

**Performance Goals**: Preserve dispatch latency p99 < 100ms. The log pipeline MUST stay off the
hot path: broker publish is already non-blocking (drops on full buffer); file writes are
buffered; the ring is O(1) append under a dedicated lock. No engine allocation regressions.

**Constraints**: Daemon and CLI remain **cgo-free** (only the GUI uses cgo/OpenGL). No new runtime
Go dependency. Existing 001–003 behavior preserved except the explicitly removed Triggers feature
and the removed manual Refresh. Constitution UX rules (structured logs, actionable errors) apply
directly to the new log records.

**Scale/Scope**: Rename touches ~85 files (import paths/strings). Triggers removal deletes ~6
files and prunes references in ~8 more. New code: one log-handler package, ring buffer, rotator,
`GET /v1/logs` + client method, broker event kinds + publish sites, GUI Logs view rewrite,
calendar widget, refresh-control removal across ~5 GUI files, and the WiX packaging + release
workflow changes.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality** — PASS. New pieces are small and single-purpose (handler, rotator, ring,
  calendar widget), each documented. Triggers removal *reduces* surface area. Errors wrapped with
  `%w`; no `panic` in library code. The rename is mechanical and verified by a clean build.
- **II. Testing Standards (NON-NEGOTIABLE)** — PASS (planned). Every behavioral change ships
  tests; migration v3 gets a regression test (DB with triggers → starts clean; DB without →
  no-op). Log rotation/ring tested as pure functions with an injected clock — no real
  `time.Sleep`. `-race` across the suite. Core scheduling packages keep ≥80% coverage (Triggers
  removal deletes both code and its tests symmetrically).
- **III. UX Consistency** — PASS, advanced. The Logs view unifies structured daemon logs and
  scheduler events under one consistent severity model (Info/Warning/Error) with actionable
  detail; real-time updates make every view behave uniformly; the MSI gives a standard,
  predictable install/uninstall. CLI/API error envelopes unchanged.
- **IV. Performance** — PASS. Log pipeline is non-blocking and bounded (ring cap + file
  rotation); no new work on the dispatch hot path. New broker event kinds are cheap fan-out.
  No benchmark regression expected; the engine bench set still guards dispatch latency. **Gate**:
  the engine benchmarks (`go test -bench ./internal/engine/...`) MUST be run and confirmed within
  10% of baseline (tasks T077a) — this is the automated proof for the p99<100ms claim.

**Known manual gate (acknowledged)**: the Windows MSI success criteria (SC-002/SC-003/SC-004 —
install, reboot-without-login, clean uninstall) are verified by the documented quickstart manual
checklist on a clean Windows VM (T078), not by a fully automated CI gate. CI builds the `.msi`
artifact and runs the wxs sanity check (T015/T020), but the install/reboot/uninstall behaviors
require a real machine. This is the one place coverage is manual rather than CI-enforced; called
out here so it is a conscious choice, not a gap.

No constitution violations → Complexity Tracking is empty.

## Project Structure

### Documentation (this feature)

```text
specs/004-rebrand-gui-overhaul/
├── plan.md                  # This file
├── spec.md
├── research.md              # Phase 0: decisions (rename, MSI/WiX, log pipeline, removal, real-time, calendar)
├── data-model.md            # Phase 1: LogRecord, event kinds, migration v3, removed entities
├── quickstart.md            # Phase 1: validation scenarios (install, logs, triggers gone, real-time, calendar)
├── contracts/
│   ├── api-logs.md          # GET /v1/logs + /v1/events extension (log/task/group kinds); removed /v1/triggers
│   ├── log-file.md          # On-disk JSONL log format + rotation/retention policy
│   └── msi-package.md       # MSI component/feature/service/shortcut contract + uninstall behavior
└── checklists/requirements.md
```

### Source Code (repository root)

```text
go.mod                       # module path → github.com/shruggietech/go-schedule
cmd/
├── goschedd/main.go         # wire log pipeline; drop trigger dispatcher wiring
├── gosched/main.go          # rename strings
└── gosched-gui/             # versioninfo.json (rename), Info.plist/bundle id rename

internal/
├── logbus/                  # NEW: slog.Handler teeing to file + ring + broker
│   ├── handler.go           # custom slog.Handler (level→severity, fields→detail)
│   ├── ring.go              # bounded in-memory ring of recent LogRecords
│   ├── rotate.go            # size-based rotating file writer (in-house)
│   └── *_test.go
├── events/broker.go         # add KindLog, KindTask, KindGroup + Publish helpers
├── domain/domain.go         # add LogRecord; REMOVE Trigger, DedupLedger, TriggerOutcome,
│                            #   ScheduleEvent kind usage, Schedule.TriggerID, RunTrigger event
├── store/
│   ├── store.go             # migration v3: DROP triggers, dedup_ledger
│   ├── trigger.go           # DELETE
│   └── crud.go              # prune trigger/event-schedule references
├── api/server/
│   ├── server.go            # remove /v1/triggers routes; add GET /v1/logs
│   ├── logs.go              # NEW: handleListLogs (severity filter)
│   ├── triggers.go          # DELETE
│   └── events.go            # stream new event kinds
├── api/client/
│   ├── methods.go / logs.go # NEW ListLogs; remove trigger client
│   └── triggers.go          # DELETE
├── cli/
│   ├── cli.go               # drop newTriggerCmd(); rename strings; (alerts cmd → logs, see research)
│   └── trigger.go           # DELETE
├── engine/engine.go         # remove FireEvent + completion/startup trigger hooks
├── trigger/                 # DELETE entire package (dispatcher.go, dedup_test.go)
├── service/service.go       # rename display name → "go-schedule"
├── config/config.go         # data dir/DB/log-file paths (rename handling — see research)
└── platform/platform_windows.go  # ProgramData dir name (rename handling — see research)

gui/
├── app.go                   # tabs: drop Triggers; Alerts→Logs; remove refreshAll buttons wiring
├── logs.go                  # NEW: Logs view (filters, detail dialog, Dismiss All); replaces alerts.go
├── alerts.go                # DELETE (superseded by logs.go)
├── triggers.go              # DELETE
├── schedule.go              # add list⇄calendar view toggle; remove Refresh button
├── calendar.go              # NEW: month-grid calendar widget over Occurrences
├── tasks.go / groups.go     # remove Refresh buttons; rely on live events
└── viewmodel/viewmodel.go   # State.Logs; ApplyEvent folds log/task/group events; API.ListLogs

build/windows/               # NEW: WiX packaging
├── goschedule.wxs           # product, components, ServiceInstall/Control, shortcuts
└── (license.rtf, assets)
.github/workflows/release.yml # build MSI on windows runner; drop zip; rename artifacts
docs/INSTALL-windows.md      # rewrite for MSI-only
README.md / CHANGELOG.md     # rename (history preserved), drop triggers, mention calendar/logs
```

**Structure Decision**: Single Go module, multi-binary, unchanged. The log pipeline is isolated in
a new `internal/logbus` package so the handler/ring/rotator are unit-testable without the daemon.
Triggers removal is a clean delete of the `internal/trigger` package plus symmetric pruning of its
call sites and a forward store migration. MSI packaging is a new `build/windows/` tree consumed
only by the Windows release job; it adds no runtime code. The calendar is a self-contained Fyne
widget fed by the existing calendar API.

## Complexity Tracking

> No constitution violations — section intentionally empty.
