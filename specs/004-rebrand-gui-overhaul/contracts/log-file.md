# Contract: On-disk log file format & rotation

## Location

- Default: `<DataDir>/logs/goschedule.log`
  - Windows (MSI/service): `C:\ProgramData\goschedule\logs\goschedule.log`
  - Linux/macOS: `<platform data dir>/goschedule/logs/goschedule.log`
- Overridable via `config.LogFilePath`. The `logs/` directory is created (0o755) on daemon start.

## Format: JSON Lines (JSONL)

One JSON object per line, append-only, UTF-8. Fields mirror `LogRecord`:

```json
{"id":"123","time":"2026-06-20T17:04:11Z","severity":"error","source":"executor","message":"task run failed","task_id":"tsk_abc","run_id":"run_def","exit_code":1,"error":"file does not exist"}
```

- Reserved keys: `id`, `time`, `severity`, `source`, `message`, `task_id`, `run_id`. Remaining
  structured slog attributes are written as additional top-level keys (the GUI shows them as the
  cause/context detail).
- The file is the durable troubleshooting record (FR-016) and survives daemon restarts.

## Rotation & retention (FR-017)

- **Trigger**: when the active file would exceed `LogMaxSizeBytes` (default 10 MiB), it is rotated.
- **Scheme**: `goschedule.log` → `goschedule.log.1` → … → `goschedule.log.<LogMaxFiles-1>`; the
  oldest is discarded. Default `LogMaxFiles = 5` ⇒ ≈ 50 MiB ceiling.
- **Writer**: in-house size-based rotating `io.Writer` (no third-party dependency). Single daemon
  writer; rotation guarded by the writer's mutex. Writes are buffered; rotation must not block the
  engine's hot path (logging publish is non-blocking; file write happens in the handler off the
  dispatch path).
- **Crash safety**: a torn final line on crash is tolerated by readers (skip unpar.seable lines).

## Behavior notes

- "Dismiss All" in the GUI does **not** truncate or delete this file (durable record).
- Uninstalling via MSI leaves this file in place by default (data preservation, see
  [msi-package.md](msi-package.md)).
- The retention bound MUST hold under sustained logging (SC-006): total bytes ≤
  `LogMaxSizeBytes × LogMaxFiles` plus the in-progress file.
