# go-schedule

Cross-platform (Linux/macOS/Windows) task scheduler in Go. A system-wide daemon
(`goschedd`) hosts the scheduling engine, SQLite store, and executor; the CLI
(`gosched`) and the Go-native Fyne GUI (`gosched-gui`) are thin clients over a
local IPC API (Unix socket / Windows named pipe). The master specification is
`specs/001-task-scheduler/spec.md`; the ordered roadmap is `TODO.md`, whose
authoritative task list is `specs/001-task-scheduler/tasks.md`. Per-feature
specs live under `specs/NNN-name/`.

## Build-phase autopilot

Standing authorization: every feature traceable to the master specification and
the roadmap runs under the Build-Phase Autopilot Protocol. A verbal kickoff
("kick off the catch-up feature", "run the next feature", "autopilot this")
authorizes the full spec-kit feature sequence end to end (specify, clarify,
checklist, plan, tasks, analyze, implement, verify, commit) with no pause for
authorization between steps. Every feature MUST be spec'd through the spec-kit
framework before implementation; the master specification scopes a feature but
never substitutes for its spec.

Default to deciding, not asking: enumerate the alternatives, evaluate them
against the constitution (`.specify/memory/constitution.md`), the master
specification, and the feature scope, pick the best, proceed, and record the
rationale. Halt to the user only when no option is clearly best on an
irreversible or architecture-defining choice, the feature intent is genuinely
ambiguous, or a constitution CRITICAL conflict needs a human call.

Halt exactly once per feature: right before anything leaves the machine, with a
breakdown of notable decisions and what was built. Never push, open a pull
request, tag, run the release workflow, or modify pinned artifacts without
explicit authorization.

The full procedure is `docs/build-autopilot.md`; the governing principle is
constitution principle V. This applies to features traceable to the master
specification and the roadmap, and to any feature or task the operator
explicitly places under autopilot; unrelated requests with no such kickoff use
normal interactive mode.

This project uses the hyphenated spec-kit command form (`/speckit-specify`,
`/speckit-plan`). `.claude/` is gitignored, so a fresh clone has no command
skills until `specify integration upgrade claude` restores them.

## Integration workflow

Work integrates through pull requests. Autopilot commits onto a feature branch
named for the feature (for example `005-catchup-hardening`), never onto `main`;
the constitution forbids direct pushes to the default branch. The single halt
precedes both the branch push and the pull request, and CI is the merge gate.

## Running verification (read before verifying)

Run CI parity in the foreground and watch it finish. NEVER launch the test suite
in the background and poll for its output. `go test` buffers a package's output
until that package completes, so a background run cannot be distinguished from a
dead one.

```bash
gofmt -l internal cmd test
```

```bash
go vet ./...
```

```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```

```bash
CGO_ENABLED=1 go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')
```

```bash
go test ./gui/...
```

`gofmt` must print nothing. The race run excludes the cgo-only GUI entry point
and the Fyne widget package (races there are inside Fyne's own font cache, not
this project's code); `gui/viewmodel` stays race-tested and the GUI is covered
by the headless run. Core packages must stay at or above 80 percent coverage.

## Non-negotiables

- Safety-critical test surfaces are never weakened or skipped: clock injection
  (no direct `time.Now()` in engine code), timezone and DST resolution,
  forward-only non-destructive store migrations, restart and catch-up recovery,
  goroutine termination under the race detector, and local IPC access control.
- CI parity before any commit, run in the foreground and watched to completion.
- Pinned artifacts (`.github/workflows/**`, `build/**`, `Makefile`,
  `.golangci.yml`, the `go`/`toolchain` lines of `go.mod`, `.gitattributes`,
  `.gitignore`, `LICENSE`, `docs/INSTALL-windows.md`) change only with a dated
  decision recorded in `CHANGELOG.md`.
- Cutting a `vX.Y.Z` tag always requires explicit authorization.

## Key conventions

Internal scheduling in UTC; per-task IANA timezone with DST (next-valid /
first-occurrence); recurrence via RFC 5545 RRULE (rrule-go) behind a
human-readable layer; injected `Clock` interface; `log/slog` structured logs;
`go test -race`; dispatch latency p99 < 100ms. The GUI is built windowless
(`-H windowsgui`) and tasks spawn with no console window.

<!-- SPECKIT START -->
# go-schedule — Active Plan

Governing documents:
- Constitution: `.specify/memory/constitution.md` (v1.1.0 — code quality, testing, UX, performance, autopilot)
- Spec: `specs/001-task-scheduler/spec.md`
- Plan: `specs/001-task-scheduler/plan.md`
- Design: `specs/001-task-scheduler/research.md`, `data-model.md`, `contracts/`, `quickstart.md`

Active feature:
- Plan: `specs/004-rebrand-gui-overhaul/plan.md` (rename go-scheduler→go-schedule; Windows .msi install w/ auto-start service; Alerts→unified Logs view w/ filters + on-disk JSONL + detail; remove Triggers entirely (migration v3); real-time GUI via broker task/group/log events (drop manual Refresh); toggleable calendar view under Schedule)
- Prior: `specs/003-gui-editor-refinements/plan.md` (GUI editor refinements: maximized window, two-pane modal + Help, code-block preview, custom collapsible, cancel-confirm, app-wide pointer cursor)
- Prior: `specs/002-gui-task-editor-ux/plan.md` (GUI task-editor UX overhaul + interval anchor)
<!-- SPECKIT END -->
