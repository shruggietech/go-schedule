---

description: "Task list for Rebrand to go-schedule + GUI & Installer Overhaul"
---

# Tasks: Rebrand to go-schedule + GUI & Installer Overhaul

**Input**: Design documents from `specs/004-rebrand-gui-overhaul/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: INCLUDED — the project constitution (II. Testing Standards, NON-NEGOTIABLE) requires
tests alongside every behavioral change, run under `go test -race`.

**Organization**: Tasks are grouped by user story. Priority order: US1 (rename) and US2 (MSI) are
P1; US3 (logs), US4 (remove triggers), US5 (real-time) are P2; US6 (calendar) is P3.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependency on an incomplete task)
- **[Story]**: US1–US6 (setup/foundational/polish have no story label)

## Path note (multi-file hotspots)

`internal/api/server/server.go`, `gui/app.go`, `gui/viewmodel/viewmodel.go`, and
`gui/schedule.go` are each edited by **multiple** stories. Those edits are sequenced (not `[P]`)
across stories. Run stories in priority order to avoid same-file conflicts.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Branch and baseline.

- [X] T001 Create and switch to branch `004-rebrand-gui-overhaul` from `main`
- [X] T002 Capture green baseline before changes: run `gofmt -l . && go vet ./... && go test -race ./...` and record the result

**Checkpoint**: Clean working tree on the feature branch with a known-green baseline.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: None of the stories share a blocking code prerequisite beyond the rename itself, and
the rename is delivered as US1 (the MVP). All later phases assume the **new** module path
`github.com/shruggietech/go-schedule` and the renamed data dir/DB from US1.

**⚠️ CRITICAL**: Complete Phase 3 (US1) before starting US3–US6 edits, because they reference the
new import path and the renamed config paths.

**Checkpoint**: After US1 lands, the remaining stories can proceed (US4 is fully independent; US3
and US5 share broker/view-model surfaces and should be done in order).

---

## Phase 3: User Story 1 - Consistent project name "go-schedule" (Priority: P1) 🎯 MVP

**Goal**: Rename the live project from go-scheduler to go-schedule (module path, build/release
config, user-facing strings) and rename on-disk identity (data dir, DB, log dir) with a
best-effort startup move from the old paths. History preserved.

**Independent Test**: `go build ./... && go test ./...` pass under the new module path; a repo-wide
search for `go-scheduler` returns only history (`CHANGELOG.md`, `specs/001..003`); GUI title and
CLI branding read "go-schedule".

### Tests for User Story 1

- [X] T003 [P] [US1] Update/assert GUI window title and app name read "go-schedule" in `gui/app_test.go`
- [X] T004 [P] [US1] Add startup data-dir/DB rename-move test (old `goscheduler` dir/DB present → moved to `goschedule`; new present → no-op; move failure → falls back without data loss) in `internal/config/migrate_paths_test.go`

### Implementation for User Story 1

- [X] T005 [US1] Change the module path in `go.mod` from `github.com/shruggietech/go-scheduler` to `github.com/shruggietech/go-schedule`
- [X] T006 [US1] Rewrite every import path `github.com/shruggietech/go-scheduler/...` → `.../go-schedule/...` across all `.go` files (mechanical pass; verify with `go build ./...`)
- [X] T007 [P] [US1] Update ldflags `-X` buildinfo path in `.github/workflows/release.yml` and `Makefile` to the new module path
- [X] T008 [P] [US1] Update service display name to "go-schedule" in `internal/service/service.go` (keep service Name `goschedd`)
- [X] T009 [P] [US1] Update GUI window title + package doc strings in `gui/app.go` to "go-schedule"
- [X] T010 [P] [US1] Update macOS bundle strings (CFBundleName/DisplayName/Identifier → `tech.shruggie.goschedule`) in `.github/workflows/release.yml` and Windows `cmd/gosched-gui/versioninfo.json`
- [X] T011 [US1] Rename data dir → `goschedule` and DB → `goschedule.db`, add log-dir derivation, in `internal/config/config.go` (`DBPath`, new log fields per data-model) and `internal/platform/platform_windows.go` (ProgramData `goschedule`)
- [X] T012 [US1] Implement best-effort one-time startup path migration (move old `goscheduler` dir/DB → new) wired into `cmd/goschedd/main.go` (before `store.Open`); log a warning on fallback (depends on T011)
- [X] T013 [P] [US1] Update remaining live name references in `CLAUDE.md` (prose only; SPECKIT block already updated), `.golangci.yml`, `.golangci.bck.yml`, and the `README.md`/`CHANGELOG.md` *headers* (leave CHANGELOG history entries and `specs/001..003` untouched)
- [X] T014 [US1] Verify: `gofmt -l . && go vet ./... && go build ./... && go test ./...` clean; run the SC-001 grep and confirm only history residue remains

**Checkpoint**: Project builds/tests under `go-schedule`; branding and on-disk paths renamed.

---

## Phase 4: User Story 2 - Formal Windows installation via MSI (Priority: P1)

**Goal**: A WiX-built `.msi` that installs to Program Files, registers `goschedd` as an auto-start
service, adds a Start-Menu shortcut, upgrades/uninstalls cleanly; the portable zip is removed.

**Independent Test**: On a clean Windows VM, the `.msi` installs the service (running, auto-start),
creates a Start-Menu entry that opens the GUI with no console, tasks fire after reboot with no
login, and uninstall removes binaries/service/shortcut (ProgramData data retained).

### Tests for User Story 2

- [X] T015 [P] [US2] Add a WiX authoring sanity check (component file names match the three built binaries; service Name == `goschedd`) as a CI script `build/windows/verify_wxs.ps1` invoked in the release workflow
- [X] T016 [US2] Document the manual install/reboot/uninstall validation checklist in `specs/004-rebrand-gui-overhaul/quickstart.md` US2 section (already drafted — confirm steps match the final wxs)

### Implementation for User Story 2

- [X] T017 [P] [US2] Author `build/windows/goschedule.wxs`: Product (per-machine, `ALLUSERS=1`, MajorUpgrade), components for `gosched-gui.exe`/`goschedd.exe`/`gosched.exe` → `C:\Program Files\go-schedule\`, per contracts/msi-package.md
- [X] T018 [P] [US2] Add WiX `ServiceInstall` (Name `goschedd`, DisplayName "go-schedule", Start=auto, LocalSystem) + `ServiceControl` (start on install; stop+remove on uninstall) to `build/windows/goschedule.wxs`
- [X] T019 [P] [US2] Add Start-Menu shortcut to `gosched-gui.exe` and any required `license.rtf`/branding assets under `build/windows/`
- [X] T020 [US2] Update `.github/workflows/release.yml` Windows GUI job to install WiX v5 and build `go-schedule_<ver>_windows_amd64.msi` from the staged binaries (depends on T017–T019)
- [X] T021 [US2] Update `.github/workflows/release.yml` to **remove** Windows `.zip` outputs (daemon-only and desktop bundle), add the `.msi` to release files + `SHA256SUMS.txt`, and rewrite the release-body Windows guidance (MSI, no "run the exe")
- [X] T022 [US2] Rewrite `docs/INSTALL-windows.md` for MSI-only install (install/upgrade/uninstall, service auto-start, data location, SmartScreen note); remove zip/"run from folder" instructions
- [X] T023 [US2] Update the Windows install section of `README.md` to point at the `.msi`
- [X] T023a [US2] FR-008 coverage: confirm the `.msi` triggers UAC elevation (per-machine/`ALLUSERS=1`) and, when elevation is declined/unavailable, Windows Installer surfaces a clear failure (no partial install); document the expected message in the quickstart US2 checklist
- [X] T023b [US2] FR-010 coverage: in the quickstart US2 manual test, explicitly verify that launching the Start-Menu GUI against the MSI-installed running service reuses it (health-check passes → no second `goschedd` spawns; single-instance lock holds)

**Checkpoint**: Release pipeline emits a working `.msi`; docs describe a formal system install only.

---

## Phase 5: User Story 3 - Unified Logs view (Priority: P2)

**Goal**: Replace Alerts with a Logs view fed by a new daemon log pipeline (rotating JSONL file +
bounded ring + `KindLog` broker events), with severity filters, click-through detail, Dismiss All,
and live updates. Existing alerts merge into the same view.

**Independent Test**: Seed a failing task → an `error` record + alert appear in Logs; filter
Errors shows only errors; clicking shows full cause; Dismiss All clears the view while the JSONL
file retains records; new records appear live.

### Tests for User Story 3

- [X] T024 [P] [US3] Ring buffer test (bounded cap, newest-first, dedupe by id) in `internal/logbus/ring_test.go`
- [X] T025 [P] [US3] Rotating-writer test (rotates at size threshold; retains `LogMaxFiles`; bound holds) using an injected size/clock in `internal/logbus/rotate_test.go`
- [X] T026 [P] [US3] Handler test: slog level→severity mapping and tee to file+ring+broker in `internal/logbus/handler_test.go`
- [X] T027 [P] [US3] API test for `GET /v1/logs` incl. `severity` filter in `internal/api/server/logs_test.go`
- [X] T028 [P] [US3] View-model test: `ApplyEvent` folds `KindLog` (prepend, dedupe) and `Refresh` loads logs in `gui/viewmodel/viewmodel_test.go`
- [X] T029 [P] [US3] GUI Logs view test (severity filter, detail open, Dismiss All clears, live append) in `gui/logs_test.go`

### Implementation for User Story 3

- [X] T030 [P] [US3] Add `LogRecord` type (fields per data-model.md) to `internal/domain/domain.go`
- [X] T031 [P] [US3] Implement bounded ring in `internal/logbus/ring.go`
- [X] T032 [P] [US3] Implement in-house size-based rotating writer in `internal/logbus/rotate.go`
- [X] T033 [US3] Implement teeing `slog.Handler` (file via rotator + ring + broker publish) in `internal/logbus/handler.go` (depends on T031, T032)
- [X] T034 [US3] Add `KindLog` + `*domain.LogRecord` payload and `PublishLog` to `internal/events/broker.go`
- [X] T035 [US3] Add log config fields (`LogFilePath`, `LogMaxSizeBytes`, `LogMaxFiles`, `LogRingSize`) + fail-fast validation in `internal/config/config.go`
- [X] T036 [US3] Wire the logbus handler into the daemon logger (replace/extend `config.NewLogger` usage) in `cmd/goschedd/main.go` so the broker is available to the handler (depends on T033, T034, T035)
- [X] T037 [US3] Implement `GET /v1/logs` handler (severity/limit/since filters) in `internal/api/server/logs.go` and register the route in `internal/api/server/server.go`
- [X] T038 [US3] Add `ListLogs` to the API client (`internal/api/client/logs.go` + methods) and to `gui/viewmodel`'s `API` interface and the GUI `Backend` interface in `gui/app.go`
- [X] T039 [US3] Extend `gui/viewmodel/viewmodel.go`: `State.Logs`, load logs in `Refresh`, fold `KindLog` in `ApplyEvent`
- [X] T040 [US3] Implement `gui/logs.go`: unified Logs list (merge `LogRecord` + `Alert`), severity/type filter control, click→detail dialog (full message + attrs/cause), "Dismiss All" (clear view + ack shown alerts); register the live refresher
- [X] T041 [US3] Update `gui/app.go`: replace the "Alerts" tab with "Logs" (use `buildLogsTab`); update the tab badge to an error count; delete `gui/alerts.go`
- [X] T042 [P] [US3] CLI: rename `alerts` command to `logs` (keep `alerts` as a deprecated alias) reading `GET /v1/logs` with `--severity`/`--json` in `internal/cli/cli.go` (+ new `internal/cli/logs.go`); update client as needed
- [X] T043 [US3] Verify `go test -race ./internal/logbus/... ./internal/api/... ./gui/...` pass

**Checkpoint**: Logs view fully functional; records persist to disk and stream live.

---

## Phase 6: User Story 4 - Remove the Triggers feature (Priority: P2)

**Goal**: Delete Triggers across GUI, CLI, API, client, engine, store, and domain; add store
migration v3 dropping `triggers`/`dedup_ledger`. Independent of other stories.

**Independent Test**: Build succeeds with `internal/trigger` deleted; no Triggers tab; no CLI
trigger command; `/v1/triggers` → 404; a daemon started against a pre-v3 DB with triggers starts
clean (migration v3); a DB without triggers is a no-op.

### Tests for User Story 4

- [X] T044 [P] [US4] Store migration v3 test: DB seeded with `triggers`/`dedup_ledger` → dropped, starts clean; DB without them → no-op, in `internal/store/store_test.go`
- [X] T045 [P] [US4] Remove now-obsolete trigger tests: delete `internal/trigger/dedup_test.go` and `test/integration/triggers_test.go`; adjust any references in other integration tests

### Implementation for User Story 4

- [X] T046 [US4] Add migration v3 (`DROP TABLE IF EXISTS dedup_ledger; DROP TABLE IF EXISTS triggers;`, optional warn-count) to `internal/store/store.go`
- [X] T047 [P] [US4] Delete the `internal/trigger/` package (`dispatcher.go` and tests)
- [X] T048 [P] [US4] Delete `internal/store/trigger.go`; prune trigger/event references in `internal/store/crud.go` and `internal/schedule/recur.go`
- [X] T049 [P] [US4] Delete `internal/api/server/triggers.go` and remove the three `/v1/triggers` routes in `internal/api/server/server.go`
- [X] T050 [P] [US4] Delete `internal/api/client/triggers.go` and remove `Create/List/DeleteTrigger` from the client methods
- [X] T051 [P] [US4] Delete `internal/cli/trigger.go` and remove `newTriggerCmd()` (cli.go:60) in `internal/cli/cli.go`
- [X] T052 [US4] Delete `gui/triggers.go`; remove the Triggers tab and the three trigger methods from the `Backend` interface in `gui/app.go`
- [X] T053 [US4] Remove engine trigger wiring: `FireEvent`, `SetCompletionHook`, `SetStartupHook` (+ `disp` dispatcher) in `internal/engine/engine.go` and `cmd/goschedd/main.go`
- [X] T054 [US4] Remove domain trigger types (`Trigger`, `DedupLedger`, `TriggerOutcome`) and stop producing `event`-kind schedules/runs in `internal/domain/domain.go` and `internal/schedule/parse.go` (retain historical enum values for `runs.trigger` per data-model)
- [X] T055 [P] [US4] Remove triggers from advertised features in `README.md` and add a CHANGELOG note
- [X] T056 [US4] Verify `go build ./... && go test -race ./...` clean with triggers fully removed

**Checkpoint**: No trigger surface remains anywhere; migration v3 verified both ways.

---

## Phase 7: User Story 5 - Real-time GUI updates (Priority: P2)

**Goal**: Add task/group change events to the broker, publish them from API mutation handlers, fold
them into the view-model, re-sync on reconnect, and remove all manual Refresh controls.

**Independent Test**: With the GUI open, a CLI mutation appears in the relevant view within ~2s
with no Refresh pressed; no Refresh button exists anywhere; killing/restarting the daemon
re-syncs the GUI automatically.

### Tests for User Story 5

- [X] T057 [P] [US5] View-model test: `ApplyEvent` folds `KindTask`/`KindGroup` (created/updated upsert, deleted remove) in `gui/viewmodel/viewmodel_test.go`
- [X] T058 [P] [US5] API test: create/update/delete/enable/disable task and group handlers publish the corresponding events in `internal/api/server/events_test.go` (or handler tests)
- [X] T059 [P] [US5] GUI test: no toolbar Refresh control remains in tasks/groups/schedule/logs views in `gui/app_test.go`

### Implementation for User Story 5

- [X] T060 [US5] Add `KindTask` and `KindGroup` event kinds with `{verb, entity|id}` payloads and `PublishTask`/`PublishGroup` helpers in `internal/events/broker.go`
- [X] T061 [US5] Publish task events after successful writes in the task handlers (`internal/api/server/tasks.go`, `update.go`) — create/update/delete/enable/disable
- [X] T062 [US5] Publish group events after successful writes in `internal/api/server/groups.go`
- [X] T063 [US5] Fold `KindTask`/`KindGroup` into `State` (upsert/remove tasks & groups) in `gui/viewmodel/viewmodel.go`
- [X] T064 [US5] On stream reconnect, perform a full `Refresh()` before resuming event application in `gui/app.go` `streamEvents` (FR-024)
- [X] T065 [P] [US5] Remove the manual Refresh button + handler from `gui/tasks.go`
- [X] T066 [P] [US5] Remove the manual Refresh button from `gui/groups.go`
- [X] T067 [US5] Remove the Refresh button from `gui/schedule.go` (keep the range selector; rely on live refresher)
- [X] T068 [US5] Verify cross-client live updates and absence of Refresh controls; `go test -race ./gui/... ./internal/api/...`

**Checkpoint**: Every view updates live; no manual Refresh anywhere; reconnect re-syncs.

---

## Phase 8: User Story 6 - Calendar view under Schedule (Priority: P3)

**Goal**: Add a toggleable month-grid calendar under Schedule over the existing `GET /v1/calendar`
data; toggling preserves the selected window; the calendar updates live.

**Independent Test**: Toggle Schedule to Calendar → occurrences on correct dates; toggle back to
List → same occurrences/window; completing a run/adding a task updates the calendar live.

### Tests for User Story 6

- [X] T069 [P] [US6] Calendar widget test: occurrences map to the correct day cells; empty window renders empty (no error) in `gui/calendar_test.go`
- [X] T070 [P] [US6] Schedule toggle test: switching List⇄Calendar preserves the selected window and shared occurrences in `gui/schedule_test.go`

### Implementation for User Story 6

- [X] T071 [US6] Implement `gui/calendar.go`: month-grid widget rendering `[]server.Occurrence`, day cells marked when occurrences exist, day selection lists that day's occurrences
- [X] T072 [US6] Add a List/Calendar view toggle to `gui/schedule.go`, hold the selected window in tab state, share loaded occurrences, and register the calendar with the same live refresher (FR-027/FR-028)
- [X] T073 [P] [US6] Mention the calendar view in `README.md` Features/GUI section
- [X] T074 [US6] Verify `go test -race ./gui/...` pass

**Checkpoint**: Calendar view available, toggles cleanly, updates live.

---

## Phase 9: Polish & Cross-Cutting Concerns

- [X] T075 [P] Rewrite `README.md` intro/Features/layout to reflect go-schedule, MSI install, Logs (with file location), real-time GUI, calendar, and no triggers
- [X] T076 [P] Add a consolidated `CHANGELOG.md` entry for this release (rename, MSI, Logs, triggers removed, real-time, calendar)
- [X] T077 Run the full gate: `gofmt -l . && go vet ./... && go test -race ./...`; confirm core scheduling package coverage ≥ 80%
- [X] T077a Run engine benchmarks `go test -bench ./internal/engine/...` and confirm dispatch latency is within 10% of the pre-change baseline captured in T002 (constitution IV)
- [~] T078 Execute `specs/004-rebrand-gui-overhaul/quickstart.md` validations and record results.
  Automatable validations PASS (build, full test suite, SC-001 rename grep, migration v3, logbus,
  real-time publish, calendar). REMAINING (manual, needs a clean Windows VM + WiX): MSI
  install/reboot/uninstall, UAC-decline, and service-reuse checks (US2 §FR-008/FR-010). The CI
  release job builds the `.msi` and runs the wxs sanity check.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies.
- **Foundational (Phase 2)**: trivial here; the effective prerequisite for US3–US6 is that **US1
  (rename) has landed** (new import path + renamed config paths).
- **US1 (Phase 3, P1)**: start after Setup. **Blocks** the others by virtue of the module rename.
- **US2 (Phase 4, P1)**: after US1 (it references the renamed product/binaries). Otherwise
  self-contained (packaging + docs).
- **US3 (Phase 5, P2)**: after US1. Adds broker `KindLog`, logbus, `/v1/logs`, Logs view.
- **US4 (Phase 6, P2)**: after US1. Fully independent of US2/US3/US5/US6.
- **US5 (Phase 7, P2)**: after US1; best after US3 (shares `broker.go`, `viewmodel.go`,
  `app.go`, `schedule.go`) and after US4 (so removed Triggers tab isn't re-touched).
- **US6 (Phase 8, P3)**: after US1; shares `schedule.go` with US5, so do it after US5.
- **Polish (Phase 9)**: after all desired stories.

### Cross-story shared files (sequence, do not parallelize across stories)

- `internal/api/server/server.go`: US3 (add /v1/logs) + US4 (remove /v1/triggers)
- `gui/app.go`: US1 + US3 (Logs tab) + US4 (drop Triggers tab) + US5 (reconnect)
- `gui/viewmodel/viewmodel.go`: US3 (logs) + US5 (task/group)
- `internal/events/broker.go`: US3 (KindLog) + US5 (KindTask/KindGroup)
- `gui/schedule.go`: US5 (remove Refresh) + US6 (calendar toggle)
- `internal/domain/domain.go`: US3 (add LogRecord) + US4 (remove Trigger types)

### Within each story

- Tests first (must fail before implementation), then models → services → endpoints → UI.

---

## Parallel Opportunities

- **Setup**: T002 after T001.
- **US1**: T003/T004 (tests) parallel; T007–T010, T013 parallel after the import rewrite (T005–T006).
- **US2**: T017–T019 (wxs authoring/assets) parallel; then T020–T023.
- **US3**: tests T024–T029 parallel; impl T030–T032 parallel, then T033→T034→…; T042 (CLI) parallel with GUI tasks.
- **US4**: deletions T047–T051, T055 parallel (different files) after T046 migration.
- **US5**: tests T057–T059 parallel; T065/T066 parallel.
- **US6**: tests T069/T070 parallel; T073 parallel with code.

---

## Implementation Strategy

### MVP scope

US1 (rename) + US2 (MSI) are the P1 MVP: a correctly-named product that installs as a formal
Windows system service. Ship/validate before P2 work.

### Incremental delivery

1. Setup → US1 (rename) → validate build/tests + branding. (MVP part 1)
2. US2 (MSI) → validate install/reboot/uninstall on a clean VM. (MVP part 2)
3. US4 (remove Triggers) → independent, reduces surface area early.
4. US3 (Logs) → durable logging + troubleshooting view.
5. US5 (real-time) → live everywhere, drop Refresh.
6. US6 (calendar) → advertised view.
7. Polish → README/CHANGELOG, full gate, quickstart.

### Notes

- [P] = different files, no incomplete-task dependency.
- Commit per task or logical group; keep `go test -race` green at each checkpoint.
- Tests must fail before implementation (constitution II).
