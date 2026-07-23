---
name: go-schedule-verify
description: Run the go-schedule CI-parity verification gates correctly - the exact six commands, in the foreground, watched to completion - and interpret the two local-environment traps that make them fail for reasons unrelated to the code. Use before any commit, before the pre-push halt, when asked to "verify", "run the tests", "check CI parity", or when a gate has failed and the cause is unclear. Also use when reporting verification results, because the honesty rules about skipped gates live here.
---

# go-schedule Verification

The CI-parity procedure for this repository. It exists as a skill rather than
only as prose in `CLAUDE.md` because prose gets paraphrased, and the specific
shortcut it guards against — launching the test suite in the background and
polling — is one an agent will rationalize into under time pressure. It has
already caused a misdiagnosed hang in this project.

## The rule that matters most

**Run every gate in the foreground and watch it finish. Never background the
test suite and poll for output.**

`go test` buffers a package's output until that package completes. A backgrounded
run is therefore indistinguishable from a dead one, and treating the second as
the first is how you conclude a suite passed when it never ran.

## The gates

Run all six. They mirror `.github/workflows/ci.yml`.

```bash
gofmt -l internal cmd test
```

Must print **nothing**. Any output is a failure.

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

```bash
sh scripts/coverage-gate.sh
```

The race run excludes the cgo-only GUI entry point and the Fyne widget package —
races there live inside Fyne's own font cache, not this project's code.
`gui/viewmodel` stays race-tested, and the GUI is covered by the headless run.

For changes touching `test/scripts/`, add:

```bash
shellcheck test/scripts/*.sh test/scripts/lib/*.sh
```

```bash
pwsh -File .claude/skills/shruggie-powershell/scripts/Test-ScriptCompliance.ps1 -Path <script.ps1>
```

## Coverage: do not substitute a different measurement

`scripts/coverage-gate.sh` is the **single implementation** of the gate — CI runs
this exact script, so the local result and the CI result are one measurement
rather than two approximations of one.

**Do not use `go test -cover` instead.** It reports per-package coverage; the
gate measures *cross-package* coverage with `-coverpkg`, where a package's
statements count as covered when any test in the tree reaches them. The two
disagree, and the disagreement is not a rounding error. When the local check and
the CI check were two different measurements, a push went out that CI rejected.

The six core packages must stay at or above 80 percent.

## The two local-environment traps

Neither indicates a problem with the repository.

### golangci-lint refuses to start

> the Go language version (go1.x) used to build golangci-lint is lower than the
> targeted Go version

Your **base** Go toolchain is older than the `go` line in `go.mod`. `go version`
can still report the newer one, because `GOTOOLCHAIN=auto` upgrades transparently
inside this repo — but `go run <linter>@<ver>` builds the linter under *its*
go.mod, which the older base toolchain already satisfies, so no upgrade happens
and the linter compiles against the older version.

Either upgrade the base Go install to match `go.mod`, or force it for that one
command:

```bash
GOTOOLCHAIN=go1.25.0 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```

**Do not "fix" this by editing `.golangci.yml` or `go.mod`.** CI installs the Go
version from `go.mod` as its base toolchain and the pinned setup passes there.
Both files are pinned artifacts anyway.

### The race run needs a C toolchain

`-race` requires cgo. A machine with no `gcc` on `PATH` fails with

> cgo: C compiler "gcc" not found

before any test runs. Install a C toolchain (MSYS2/MinGW-w64 on Windows), or
rely on CI for the race gate — **and say so explicitly** rather than reporting
the suite as passing.

Check before you start:

```bash
command -v gcc || command -v cc || command -v clang || echo "NO C COMPILER - race gate cannot run locally"
```

## Reporting results honestly

This is part of the procedure, not an afterthought.

- A gate that did not run is **not** a gate that passed. Say which ones ran,
  which did not, and why.
- A skipped test is not a passing test. Print the skip reason.
- If a linter is unavailable, report it as not-run. Do not omit it from the
  report — an absent line reads as an absent problem.
- Paste the actual output for failures. "Tests fail" without the output is not a
  report.

A skipped gate reported as green silently invalidates every other claim in the
report, which is why this section exists at all.

## When a gate is red

A red local run is a **halt**, not something to push and sort out afterwards.
CI runs on every push to `main` but reports *after* the push rather than blocking
a merge, so the local run is the real gate. See
[`docs/build-autopilot.md`](../../../docs/build-autopilot.md) and constitution
principle II.
