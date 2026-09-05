# Quickstart: Task Execution Safety and Diagnostics

## 1. Run focused persistence and interface tests

```bash
go test ./internal/store ./internal/api/server ./internal/api/client
```

Expected: migration v10 preserves v9 data and rollback behavior; create requests
honor omitted/false/true enabled intent; exact run lookup returns the requested
run or not-found; correlation and truncation fields round-trip.

## 2. Run engine, executor, and group-policy tests

```bash
go test ./internal/engine ./internal/executor ./internal/task
```

Expected: every new failed-run alert carries the recorded run ID, output remains
bounded with an explicit truncation flag, and nearest disabled ancestors are
identified without changing existing chain-enabled behavior.

## 3. Run headless desktop tests

```bash
go test ./gui/...
```

Expected:

- fresh creation offers a cleared activation checkbox;
- explicit false/true values reach task creation atomically;
- edit state is preserved and has no creation-only activation control;
- Tasks rows separate Enabled, Lifecycle, and Effective values;
- blocked rows name the nearest disabled group with full disclosure;
- failed-run detail is selectable and handles exact data, empty/truncated output,
  legacy alerts, missing runs, and missing tasks.

## 4. Run race-focused affected packages

```bash
go test -race ./internal/store ./internal/api/server ./internal/api/client ./internal/engine ./internal/executor ./internal/task ./gui/viewmodel
```

Expected: PASS with no data races.

## 5. Run canonical verification

```bash
sh scripts/verify.sh all
```

Expected: all eight gates pass in order: format, vet, lint, race, GUI, coverage,
docs, and automation. Each core package remains at or above 80 percent.

## 6. Review issue traceability

- #102: exact task/run identity, exit/launch status, bounded combined output,
  empty/truncated states, trigger coverage, and privacy boundary.
- #118: cleared creation default, atomic inactive creation, explicit opt-in,
  validation retention, and edit-state preservation.
- #120: runnable/group-blocked/own-disabled/lifecycle states, nearest group,
  live refresh, full disclosure, and accessibility.
