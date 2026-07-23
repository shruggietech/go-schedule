# Quickstart: verifying feature 007

**Feature**: `007-issue-cleanup-docs`

**Date**: 2026-07-23

How to check each acceptance scenario, and — where a check cannot be run before
release — what stands in for it and what stays open.

## CI parity

Run all six gates in the **foreground**, watched to completion. Never background
the suite: `go test` buffers a package's output until that package finishes, so
a backgrounded run cannot be told apart from a dead one.

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

```bash
sh scripts/coverage-gate.sh
```

## Installer `PATH` (SC-001, SC-002)

The WiX source is checked before the build:

```bash
pwsh build/windows/verify_wxs.ps1
```

The end-to-end check needs an installer produced by the release workflow and a
machine that has never had the repository checked out — a development machine
already has the directory on `PATH`, which is precisely what hid the defect.
After the v0.6.0 MSI is installed there, in a **new** shell:

```powershell
gosched --version
```

```powershell
([Environment]::GetEnvironmentVariable('Path','Machine') -split ';') |
  Where-Object { $_ -like '*go-schedule*' }
```

Expect exactly one entry. Repeat that second command after an uninstall and
expect none, and after a major upgrade and expect one rather than two.

This check cannot run before the release. The pre-release evidence is the
sanity check above plus a read of the generated element; issue #5 stays open
until the post-release check passes.

## Service status (SC-003, SC-004)

From a **non-elevated** shell, as a user who is not in `Administrators`, against
the installed service:

```powershell
gosched service status
```

Expect `status: running` or `status: stopped`. An `Access is denied` here means
the fix did not take.

Then confirm the restriction that must survive:

```powershell
gosched service start
```

Expect this to still fail on privileges. The Interactive Users ACE grants no
start or stop right, and relaxing that would be a regression, not a bonus.

With no service installed, `service status` must report the not-installed
wording unchanged — that distinguishes "we asked for too much access" from "the
service is absent", which is the exact confusion the issue is about.

On Linux and macOS, every service subcommand must behave as it did before.

## Issue forms (SC-005)

On the repository, click **New issue**. Expect the two forms and no blank-issue
route. Submit the bug form with version, component, install method, OS, or
elevation left empty and confirm it is rejected.

## Documentation (SC-006, SC-007)

Compare `docs/cli.md` against the command tree the binary exposes:

```bash
gosched --help
```

Every command and subcommand listed there — and every subcommand of `task`,
`group`, `runs`, and `service` — must appear in the reference.

Then resolve every relative link in `README.md` and `docs/**` against the
working tree. Several documents are new and cross-reference each other, so a
typo here is not caught by anything else.

Finally, re-read each authored document for the house-style rules no tool
checks: a single H1, hard wrap at 80 columns, wrap-safe continuation lines, and
a table of contents whose anchors match the actual heading slugs.
