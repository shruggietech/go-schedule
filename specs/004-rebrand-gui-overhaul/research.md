# Research: Rebrand + GUI & Installer Overhaul

Phase 0 decisions. Each resolves an unknown surfaced by the spec/plan. Format: Decision /
Rationale / Alternatives considered.

## 1. Module & project rename mechanics

**Decision**: Rename the Go module path from `github.com/shruggietech/go-scheduler` to
`github.com/shruggietech/go-schedule` (org segment `shruggietech` unchanged). Update `go.mod`'s
`module` line and every import via a single mechanical pass, then update non-import references
(release scripts, ldflags `-X` paths, docs, GUI window title, service display name, macOS bundle
strings). Keep binary names `gosched`/`goschedd`/`gosched-gui` as-is (they contain no
`go-scheduler`).

**Rationale**: The module path appears in ~85 files and in build ldflags
(`-X .../internal/buildinfo.Version`); a consistent find-and-replace is the lowest-risk approach
and a clean `go build ./... && go test ./...` is the objective gate. Binary names are already
abbreviated and user-muscle-memory; renaming them would churn docs and the autostart "daemon next
to me" lookup for no user benefit.

**History preservation**: `CHANGELOG.md` existing entries and `specs/001`/`002`/`003` are
immutable history — left verbatim. Only forward-facing prose in README/CHANGELOG header is
updated. SC-001's repo-wide search treats those files as the allowed residue.

**On-disk identifiers (`goscheduler` without hyphen)**: the data dir is
`C:\ProgramData\goscheduler` (Windows) / platform dirs elsewhere, and the DB is `goscheduler.db`.

- **Decision**: rename the data dir to `goschedule` and the DB to `goschedule.db`, with a
  **one-time best-effort migration on daemon startup**: if the new dir/DB is absent but the old
  one exists, move (rename) it; if the move fails, fall back to the old path and log a warning
  (never block startup or lose data). On Windows the MSI is a fresh install so most users get the
  new layout directly; the startup migration covers prior zip users.
- **Rationale**: a full rebrand should not leave `goscheduler` paths behind, but the constitution
  requires forward migration of persisted state — the best-effort move satisfies both without a
  destructive cutover.
- **Alternatives**: (a) keep `goscheduler.db` forever — rejected as an inconsistent half-rename;
  (b) hard cutover with no fallback — rejected, risks orphaning existing users' tasks.

## 2. Windows MSI installer

**Decision**: Build the installer with **WiX Toolset v5** from a checked-in
`build/windows/goschedule.wxs`, compiled on the existing `windows-latest` release job. The MSI:

- Installs `gosched-gui.exe`, `goschedd.exe`, `gosched.exe` to
  `C:\Program Files\go-schedule\`.
- Declares a **per-machine** install (`Package InstallScope` / `ALLUSERS`) → requires elevation;
  the bootstrapper/MSI prompts via UAC.
- Registers `goschedd` as a Windows service via WiX `ServiceInstall` (Start = auto) +
  `ServiceControl` (start on install, stop+remove on uninstall). This replaces the runtime
  `kardianos/service install` step for MSI installs.
- Creates a Start-Menu shortcut to `gosched-gui.exe` (windowless build, no console).
- Sets up `C:\ProgramData\goschedule\` as the data dir (created by the daemon on first run; the
  MSI does not need to seed it) and the log dir under it.
- Uninstall removes binaries, the service, and shortcuts. **User data** (DB, logs) under
  `ProgramData` is **left in place** by default (documented), so an uninstall/reinstall preserves
  tasks; a "remove all data" note is documented for manual cleanup.

**Service identity**: the WiX `ServiceInstall` uses the same service name `goschedd` the
`kardianos/service` control layer expects, so CLI `service status/start/stop` still works against
the MSI-installed service. `goschedd` already runs correctly under the Windows SCM (it calls
`service.Run`), so no daemon code change is required for the service to host under WiX.

**Rationale**: WiX is the de-facto standard for native MSIs, supports declarative service
install/uninstall and shortcuts, and runs headless in CI. It produces a real Windows Installer
package (Add/Remove Programs entry, repair, clean uninstall) — exactly the "formal system
install" the spec demands.

**Alternatives considered**:
- `go-msi` — simpler templating but less maintained and weaker service/uninstall control.
- MSIX — modern but heavier signing/packaging requirements and worse fit for a background
  service; rejected.
- Inno Setup / NSIS — produce `.exe` installers, not `.msi`; the spec explicitly asks for `.msi`.
- Keeping `kardianos` install at runtime only — rejected: that's the current "run the exe" model
  the spec is replacing.

**Signing**: code-signing the MSI is desirable but treated as separate infra. An unsigned MSI is
acceptable for this feature; SmartScreen guidance is documented (consistent with today's unsigned
binaries).

## 3. Unified log pipeline (Logs view backend)

**Decision**: Introduce `internal/logbus` with a custom `slog.Handler` that the daemon installs as
its logger. The handler **tees** every record three ways:

1. **Rotating JSONL file** at `<dataDir>/logs/goschedule.log` — the durable troubleshooting record
   (FR-016). Rotation is size-based (default 10 MiB × 5 files ≈ 50 MiB cap, configurable) via an
   in-house rotating writer (FR-017).
2. **Bounded in-memory ring** of the most recent N records (default 1000) — served by
   `GET /v1/logs` so the GUI can load current logs and apply filters without parsing the file.
3. **Event broker** as a new `KindLog` event — so open GUI clients see new records live (FR-018).

`LogRecord` carries: severity (mapped from slog level: Debug/Info→info, Warn→warning,
Error→error), timestamp, source (logger component/`slog` group), message, optional task/run
correlation id, and the structured attributes rendered as the "cause/context" detail (FR-012,
FR-014).

**Alerts unification**: existing scheduler alerts (overlap/failure/missed) continue to be
persisted in the `alerts` table and published as `KindAlert`. The GUI Logs view merges
`KindAlert` and `KindLog` into one chronological list, tagging each with its source. We do **not**
migrate alerts into the log file; they remain queryable via the existing `/v1/alerts`. This keeps
the durable alert history intact while presenting a unified view.

**Rationale**: A `slog.Handler` is the idiomatic seam — the daemon already builds its logger in
`config.NewLogger`, so swapping in the teeing handler is localized. The ring keeps the GUI fast
and filterable without a hot SQLite write per log line (which would risk the performance budget at
high log volume). The file gives external troubleshooting and survives restarts. Non-blocking
publish keeps the engine off any logging stall.

**Alternatives considered**:
- **Persist every log to a SQLite `logs` table** — gives full history queryability but adds a DB
  write per record; rejected on performance grounds (could pressure the single-writer store under
  load). The JSONL file + ring is bounded and off the hot path.
- **`lumberjack` for rotation** — adds a runtime dependency; rejected per the constitution since a
  single-writer size rotator is small and testable in-house.
- **File-only (no ring/API)** — would force the GUI to tail/parse the file; rejected as more
  fragile and harder to filter than serving structured records over the API.

**"Dismiss All" semantics**: clears the GUI's in-memory view and acknowledges the currently shown
alerts; it does **not** delete the on-disk log file (the durable record). This matches the spec
assumption. (Flagged for stakeholder confirmation in the spec; default chosen here.)

## 4. Triggers removal (all layers) + migration

**Decision**: Delete the entire feature and add **store migration v3**:

```sql
DROP TABLE IF EXISTS dedup_ledger;
DROP TABLE IF EXISTS triggers;
```

Code deletions: `internal/trigger/` (package), `internal/store/trigger.go`,
`internal/api/server/triggers.go`, `internal/api/client/triggers.go`, `internal/cli/trigger.go`,
`gui/triggers.go`, `test/integration/triggers_test.go`. Reference pruning: remove
`newTriggerCmd()` (cli.go:60), the three `/v1/triggers` routes (server.go), the `Create/List/
DeleteTrigger` methods from the GUI `Backend` interface (app.go) and API client, the engine's
`FireEvent` + `SetCompletionHook`/`SetStartupHook` trigger wiring (engine.go, goschedd main.go),
and the domain types `Trigger`, `DedupLedger`, `TriggerOutcome`.

**Schedule "event" kind**: `ScheduleEvent` / `Schedule.TriggerID` / `RunTrigger == event` exist to
support event-driven schedules. Since event scheduling only existed to serve triggers:

- **Decision**: remove the `event` schedule kind from new validation/creation paths and treat any
  pre-existing `event`-kind schedule defensively — migration v3 leaves the `schedules` columns in
  place (no destructive column drop in SQLite) but the engine no longer dispatches `event`
  schedules. In practice no such rows are expected; the migration logs a warning if any are found.
- **Rationale**: dropping columns in SQLite requires a table rebuild; leaving unused nullable
  columns is harmless and lower-risk than a rebuild. The behavior (no event dispatch) is what
  "triggers removed" means.

**RunTrigger enum**: keep historical `runs.trigger` values as-is (append-only history); just stop
producing new `event` runs. No migration needed for `runs`.

**Rationale**: A clean delete plus a forward migration satisfies FR-019/FR-020. Symmetrically
deleting tests preserves coverage ratios. `DROP TABLE IF EXISTS` is a no-op on databases that
never had triggers (FR-020's clean-start requirement).

**Alternatives considered**: feature-flag/hide only — rejected by the clarified scope (full
removal). Rebuilding `schedules` to drop columns — rejected as unnecessary risk.

## 5. Real-time updates across all views

**Decision**: Extend the broker with `KindTask` and `KindGroup` events (carrying the changed
entity or its id + a change verb: created/updated/deleted). Publish them from the API mutation
handlers (`handleCreateTask`, `handleUpdateTask`, `handleDeleteTask`, enable/disable, and the
group equivalents) after a successful store write. The view-model's `ApplyEvent` folds these into
`State` (upsert/remove in the tasks/groups slices) in addition to today's run/alert/log folding.
The GUI then relies entirely on the stream + initial load; **all manual Refresh buttons are
removed** (`schedule.go`, `alerts.go`→`logs.go`, plus any in tasks/groups). On stream
reconnect (existing 2s backoff loop in `app.go`), the GUI performs a full `Refresh()` to
re-sync (FR-024).

**Rationale**: The SSE broker, `StreamEvents`, and `Model.ApplyEvent` already exist and already
drive live alert/run updates; today only run/alert events flow, so cross-client task/group
changes need a manual refresh. Publishing entity-change events closes that gap with minimal new
machinery and makes every view live (FR-022/FR-023). Publishing from the API layer (the single
write path) guarantees coverage regardless of which client mutated.

**Alternatives considered**:
- Client-side polling timer — rejected: reintroduces latency/refresh semantics the spec removes.
- Publishing from the store layer — rejected: the store is intentionally UI-agnostic; the API
  handler is the right seam and already has the broker.

**Reconnect/initial-load race**: on reconnect, do the full `Refresh()` *before* resuming event
application, so a missed event during the gap can't leave stale state. The existing
dedupe-by-id folding makes a redelivered event idempotent.

## 6. Calendar view under Schedule

**Decision**: Add a self-contained Fyne `calendar.go` widget that renders a **month grid** of
`server.Occurrence` items returned by the existing `GET /v1/calendar` (`GetCalendar(from,to)`).
The Schedule tab gains a view toggle (segmented control: "List" / "Calendar") that swaps the
content while preserving the selected time window; both views share the same loaded occurrences
and the same live-refresh registration. The calendar marks days with occurrences and, on day
selection, lists that day's occurrences (reusing the agenda row rendering).

**Rationale**: The backend calendar API and `Occurrence` model already exist (the current
Schedule tab even calls `GetCalendar`), so this is purely presentation. Fyne has no built-in
month-calendar in core `widget`, but `fyne.io/x/fyne` (already a dependency since 002) provides a
date-picker calendar that can be adapted, or a simple grid can be built with `container.NewGridWithColumns(7)`
of tappable day cells. Prefer the lightweight in-house grid to keep full control over
occurrence markers and live updates.

**Alternatives considered**:
- A third-party full calendar widget — none mature for Fyne; rejected.
- Reusing `xwidget.Calendar` directly — it's a date *picker*, not an occurrence display; we adapt
  the grid concept rather than its picker semantics.

**Window preservation on toggle**: the selected range (1/7/30 days, or a month for the calendar)
is held in the Schedule tab state, not in the child view, so toggling never reloads or loses it
(FR-027). The calendar registers the same `load` refresher so it updates live (FR-028).

## 7. CLI surface for logs (consistency follow-on)

**Decision**: Rename/extend the CLI `alerts` command to a `logs` command (keeping `alerts` as a
deprecated alias for one release) that can show recent log records with a `--severity`
filter and `--json`, mirroring the GUI Logs view and reading `GET /v1/logs`. This keeps the CLI
and GUI consistent (constitution III) now that "Logs" is the primary concept.

**Rationale**: The spec's Logs work is GUI-centric, but the constitution requires interface
consistency across CLI/GUI/API. Surfacing logs in the CLI is a small addition over the new
`/v1/logs` endpoint and avoids the CLI lagging behind the GUI vocabulary.

**Alternatives considered**: leave the CLI `alerts` command untouched — acceptable but
inconsistent; chosen approach is low-cost. (This is a soft requirement; if scope must shrink, the
CLI change can drop to "rename string only".)
