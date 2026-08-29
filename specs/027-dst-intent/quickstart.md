# Quickstart: Validate explicit DST scheduling intent

## Focused tests

```powershell
go test ./internal/timezone ./internal/schedule ./internal/catchup ./internal/engine ./internal/store ./internal/api/server ./internal/cli ./gui
```

Expected: policy resolution, real-transition matrices, migration, lifecycle, and client-boundary tests pass.

## Manual API sequence

1. Preview and create `every 6 hours starting at 09:00` in `America/New_York` once with `wall_clock` and once with `elapsed`.
2. Inspect next runs around 2026-03-08 and compare local readings and UTC gaps.
3. Preview a 01:30 wall-clock recurrence with overlap values `first`, `both`, and `last` around 2026-11-01.
4. Restart the daemon and confirm task detail and calendar return the same policy values and instants.
5. Attempt `elapsed` with a monthly ordinal recurrence and confirm a non-mutating `time_basis` validation error.

## Full verification

Run all eight gates through the repository's canonical verification workflow. On this Windows host, run the race gate through the established WSL host split when the native C compiler is unavailable. Record commands, results, benchmark comparison, encoding audit, and acceptance trace in `verification.md`.
