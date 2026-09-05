# Quickstart: Pure Schedule String Conversion

**Date**: 2026-08-28

## Prerequisites

- Repository root on branch `codex/018-cron-string-conversion`.
- Go toolchain selected from `go.mod`.
- No daemon is required. Stop it or point the environment at an unavailable endpoint when manually proving locality.

## Focused tests

```sh
go test ./internal/cron ./internal/cli -run 'Test.*Convert' -count=1
```

## Build the CLI

```sh
go build -o ./tmp/gosched-convert ./cmd/gosched
```

PowerShell may use `./tmp/gosched-convert.exe`; POSIX shells use `./tmp/gosched-convert`.

## Cron to human

```sh
./tmp/gosched-convert cron convert "0 9 * * 1-5"
```

Expected stdout:

```text
weekdays at 09:00
```

## Human to cron

```sh
./tmp/gosched-convert cron convert "weekdays at 09:00"
```

Expected stdout:

```text
0 9 * * 1-5
```

## Forced direction and refusal

```sh
./tmp/gosched-convert cron convert --to human "61 9 * * *"
```

Expected: empty stdout, a minute-range diagnostic on stderr, and exit code 2.

## Structured mode

```sh
./tmp/gosched-convert --json cron convert "weekdays at 09:00"
./tmp/gosched-convert --json cron convert "every 7 minutes"
```

The first command writes the five-field success object to stdout. The second writes the same shape with a refusal reason to stderr and exits 2.

## Shell quoting

POSIX:

```sh
./tmp/gosched-convert cron convert 'weekdays at 09:00'
```

PowerShell:

```powershell
.\tmp\gosched-convert.exe cron convert 'weekdays at 09:00'
```

In both cases, one quoted value reaches the command. Unquoted words are multiple arguments and are rejected by the normal command parser.

## Full verification

Run the canonical aggregate in the foreground:

```sh
sh scripts/verify.sh all
```

Record focused red/green evidence, all eight gate results, checklist completion, and issue #51 closure eligibility in `verification.md`.
