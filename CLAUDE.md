# go-schedule

Cross-platform (Linux/macOS/Windows) task scheduler in Go. A system-wide daemon
(`goschedd`) hosts the scheduling engine, SQLite store, and executor; the CLI
(`gosched`) and the Go-native Fyne GUI (`gosched-gui`) are thin clients over a
local IPC API (Unix socket / Windows named pipe). The master specification is
`specs/001-task-scheduler/spec.md`, whose authoritative task list is
`specs/001-task-scheduler/tasks.md`; that scope is delivered, and the roadmap of
work still open is the GitHub issue tracker (`gh issue list`). Per-feature specs
live under `specs/NNN-name/`.

## Build-phase autopilot

Standing authorization: every feature traceable to the master specification or
to an open issue runs under the Build-Phase Autopilot Protocol. A verbal kickoff
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
explicit authorization. Authorization to publish a PR also covers verified,
in-scope review-fix pushes to that same branch until the PR closes or merges;
material scope expansion, another PR, tags, and releases remain outside it.

The full procedure is `docs/build-autopilot.md`; the governing principle is
constitution principle V. This applies to features traceable to the master
specification or an open issue, and to any feature or task the operator
explicitly places under autopilot; unrelated requests with no such kickoff use
normal interactive mode.

This project uses the hyphenated spec-kit command form (`/speckit-specify`,
`/speckit-plan`). `.claude/` is gitignored, so a fresh clone has no command
skills until `specify integration upgrade claude` restores them.

## Integration workflow

Development is pull-request based. Every maintainer, automation agent, and
outside contributor works on a review branch and integrates through a pull
request targeting `main`. Direct pushes to `main` are prohibited.

Autopilot runs local CI parity and commits on the review branch, then halts
exactly once before pushing the branch and opening its pull request. After
authorization, existing CI and third-party AI reviewers provide additional
evidence. Consider each review comment, implement warranted changes, and explain
why other suggestions do not fit. The maintainer retains final merge judgment;
verified, in-scope review fixes may be pushed to this same PR under its publication
authorization without another halt. This workflow adds no branch-protection or
approval ceremony. A PR that fully completes an issue uses `Closes #N`; partial
work uses `Refs #N`.

## Running verification (read before verifying)

Run CI parity in the foreground and watch it finish. NEVER launch the test suite
in the background and poll for its output. `go test` buffers a package's output
until that package completes, so a background run cannot be distinguished from a
dead one.

```bash
sh scripts/verify.sh all
```

This is the single definition of green. It runs `format`, `vet`, `lint`, `race`,
`gui`, `coverage`, `docs`, and `automation`, in that order, and stops on the
first failure. Use `sh scripts/verify.sh <gate>` only to diagnose an individual
gate. The format gate must print no files. The race gate excludes the cgo-only
GUI entry point and the Fyne widget package (races there are inside Fyne's own
font cache, not this project's code); `gui/viewmodel` stays race-tested and the
GUI is covered by the headless gate.

`scripts/coverage-gate.sh` is the core-package coverage gate: the six core
packages must stay at or above 80 percent. CI runs this exact script, so the
local result and the CI result are the same measurement rather than two
approximations of one — do not substitute `go test -cover`, which reports
per-package coverage and will disagree, because the gate measures cross-package
coverage with `-coverpkg` (a package's statements count as covered when *any*
test in the tree reaches them).

Two local-environment traps, neither of which indicates a problem with the repo:

- **golangci-lint refuses to start** with "the Go language version (go1.x) used
  to build golangci-lint is lower than the targeted Go version". Your *base* Go
  toolchain is older than the `go` line in `go.mod`. `go version` can still
  report the newer one, because `GOTOOLCHAIN=auto` upgrades transparently inside
  this repo — but `go run <linter>@<ver>` builds the linter under *its* go.mod,
  which the older base toolchain already satisfies, so no upgrade happens and the
  linter is compiled with the older version. The canonical driver derives
  `GOTOOLCHAIN` from `go.mod`; for a direct diagnostic invocation, either
  upgrade the base Go install or force the matching toolchain for that one
  command:

  ```bash
  GOTOOLCHAIN=go1.25.0 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.0 run ./...
  ```

  Do not "fix" this by editing `.golangci.yml` or `go.mod` — CI installs the Go
  version from `go.mod` as its base toolchain and the pinned setup passes there.
- **The race run needs a C toolchain.** `-race` requires cgo, so a machine with
  no `gcc` on `PATH` fails with `cgo: C compiler "gcc" not found` before any test
  runs. Install a C toolchain (MSYS2/MinGW-w64 on Windows) or rely on CI for the
  race gate, and say so explicitly rather than reporting the suite as passing.

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
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan
at specs/045-structured-data-tables/plan.md
<!-- SPECKIT END -->
