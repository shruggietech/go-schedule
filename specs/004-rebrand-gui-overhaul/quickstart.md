# Quickstart / Validation Guide

Runnable checks proving each user story. Run from repo root unless noted. `go test` runs headless
(no display/cgo needed for the daemon/CLI/view-model and Fyne test-driver GUI tests).

## Prerequisites

- Go toolchain (`go.mod` version).
- For the GUI build only: C toolchain + OpenGL (WinLibs MinGW on this machine; see memory).
- For the MSI: WiX Toolset v5 (CI-provided; local optional).

## US1 — Rename to go-schedule

```bash
go build ./...              # builds under the new module path github.com/shruggietech/go-schedule
go test ./...              # all existing tests pass post-rename
```

- Repo-wide search shows only history residue:
  ```bash
  grep -rn "go-scheduler" --exclude-dir=.git .   # only CHANGELOG history + specs/001..003 allowed
  ```
- Launch the GUI; window title reads **go-schedule**. CLI `gosched --help` branding reads
  go-schedule. (SC-001)

## US2 — Windows MSI (manual, on a clean Windows VM)

CI proof (automated): `build/windows/verify_wxs.ps1` runs in the release workflow and `wix build`
produces `go-schedule_<ver>_windows_amd64.msi`. The behaviors below need a real machine.

1. Build/download `go-schedule_<ver>_windows_amd64.msi`.
2. Double-click → UAC prompt → complete wizard. Confirm:
   - Files in `C:\Program Files\go-schedule\` (not Downloads). (FR-005)
   - Service present & running: `sc query goschedd` → `RUNNING`; `sc qc goschedd` → `START_TYPE: 2
     AUTO_START`. (FR-006)
   - Start-Menu **go-schedule** entry launches the GUI, no console window. (FR-007)
3. Reboot without logging in elsewhere; a scheduled task still fires. (SC-003)
4. Apps & features → Uninstall → binaries/service/shortcut gone (`sc query goschedd` →
   `1060 service does not exist`); `C:\ProgramData\goschedule\` retained. (FR-009/SC-004)

### FR-008 — elevation behavior (negative test)
- Run the `.msi` and **decline** the UAC prompt. Expected: the install cancels cleanly with a
  clear message; nothing is left in `C:\Program Files\go-schedule\` and `sc query goschedd`
  reports the service does not exist (no partial install). The `perMachine` scope in the wxs is
  what forces the elevation prompt.

### FR-010 — GUI reuses the running service
- With the service installed and `RUNNING`, launch **go-schedule** from the Start Menu, then check
  processes: there is exactly **one** `goschedd.exe` (the service). The GUI's health-check finds
  the running daemon and does **not** spawn a second one; the single-instance lock in
  `C:\ProgramData\goschedule\goschedd.lock` would block a second anyway. (FR-010)

## US3 — Logs view

```bash
go test ./internal/logbus/...     # ring, rotator, handler severity mapping
go test ./gui/...                 # Logs view: filter, detail dialog, Dismiss All (headless)
go test ./internal/api/...        # GET /v1/logs incl. severity filter
```

Manual end-to-end:
1. Seed a failing task (bad command) and let it run → an `error` LogRecord + run_failed alert.
2. Open Logs view: entries show newest-first with severities. (US3 #1)
3. Filter **Errors** → only error rows; **Warnings** → only warnings; clear → all. (US3 #2)
4. Click the error → detail shows full message + `attrs` (exit code, error chain). (US3 #3)
5. **Dismiss All** → list clears; on-disk file still has the records:
   ```bash
   tail -n 5 "<DataDir>/logs/goschedule.log"   # JSONL records persist (FR-016)
   ```
6. With Logs open, trigger another failure → it appears live, no manual refresh. (FR-018)
- Retention: drive sustained logging; total log bytes stay ≤ `LogMaxSizeBytes × LogMaxFiles`.
  (SC-006)

## US4 — Triggers removed

```bash
go build ./...                    # compiles with internal/trigger deleted
go test ./internal/store/...      # migration v3 regression: DB-with-triggers starts clean;
                                  # DB-without is a no-op
```

- GUI has no Triggers tab. `gosched --help` lists no trigger command; `gosched trigger ...` → unknown.
- `curl --unix-socket <ipc> http://x/v1/triggers` → `404 not_found`.
- Start the daemon against a pre-v3 DB that had triggers → starts cleanly; a warning LogRecord
  notes removed trigger rows. (FR-020/SC-007)

## US5 — Real-time updates

```bash
go test ./gui/viewmodel/...       # ApplyEvent folds task/group/log events (upsert/remove)
```

Manual two-client check:
1. Open the GUI on Tasks. From another shell: `gosched task add ...`.
2. The new task appears in the GUI within ~2s, no Refresh pressed. (SC-008)
3. Edit/delete from CLI → GUI reflects it live. Schedule view updates as runs complete.
4. Confirm **no Refresh button** exists in any view (Tasks, Schedule, Groups, Logs). (FR-023)
5. Kill & restart the daemon → GUI reconnects and re-syncs automatically. (FR-024)

## US6 — Calendar view

```bash
go test ./gui/...                 # calendar widget: occurrence placement; toggle preserves window
```

Manual:
1. Schedule tab → toggle **Calendar**. Occurrences appear on correct dates. (US6 #1)
2. Toggle back to **List** → same occurrences for the same window (window preserved). (US6 #2/SC-009)
3. With calendar open, complete a run / add a task → calendar updates live. (US6 #3)

## Full gate

```bash
gofmt -l . && go vet ./... && go test -race ./...
```
All clean; core scheduling package coverage ≥ 80%.
