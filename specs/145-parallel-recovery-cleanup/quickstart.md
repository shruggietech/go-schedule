# S145 validation guide

## Prerequisites and boundaries

Use current source and an owned disposable schema-26 store, existing hidden native command runners and synthetic headless browser fixtures. No installed service, operator desktop focus or interactive credential flow is needed. Old stores are refused rather than migrated/deleted.

Prepare existing task IDs for Prepare, Export A, Export B, Publish, Recover and Cleanup. Branch definitions use the existing branches-file JSON format from [the CLI guide](../../docs/cli.md).

## Configure and observe

```text
gosched join create --name "Recover exports" --entry PREPARE_ID --target PUBLISH_ID --branches-file branches.json --recovery RECOVER_ID --recovery-on failure,canceled,skipped,not_started --cleanup CLEANUP_ID
gosched join show JOIN_ID --json
gosched join instances --join JOIN_ID --json
gosched join instance INSTANCE_ID --json
gosched join update JOIN_ID --recovery-on none
gosched join update JOIN_ID --recovery= --cleanup=
```

Run the entry through existing task execution. With all-success branches, only Publish is selected; with one selected unsuccessful branch, Recover waits for its sibling. Cleanup waits for the chosen follow-up, even if it fails or cannot start. Inspect exact selected route, native stage IDs, settlements and reasons. An explicit empty condition list survives readback and selects no recovery.

## Owned implementation qualification

```text
go test -race ./internal/store ./internal/engine ./internal/api/server ./internal/bundle
go test -race ./test/integration -run TestPackagedParallelRecovery
sh scripts/verify.sh all
```

Focused tests also cover API/client/CLI/MCP/desktop, duplicate/concurrent completions, failed cancellation feedback while native work remains active, no-start-only wake, immutable edited snapshots, deletion with late native results, portable v3/v4 and full-name sparse editing. The full synthetic browser suite exercises cause selection, vertical routing help and exact stage inspection; it does not establish native Wails window behavior.

Run fresh actual Windows and Linux/WSL packaged workflows and retained paired baseline/candidate engine benchmarks through verified hidden foreground launchers. WSL networking limitations remain separate from actual Linux execution. Freeze source/protected inputs for the complete aggregate, retain original failures and qualify exact final-head hosted checks after the authorized push.
