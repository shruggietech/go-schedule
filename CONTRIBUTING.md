# Contributing to go-schedule

**Audience:** anyone proposing a change to go-schedule\
**Governed by:** [the project constitution](.specify/memory/constitution.md)\
**See also:** [docs/](docs/README.md) · [CHANGELOG.md](CHANGELOG.md)

Contributions are welcome. This document describes how the project actually
works rather than how a Go project conventionally works, because on two points
— integration and specification — it differs deliberately.

## Contents

- [Before you write code](#before-you-write-code)
- [How work is integrated](#how-work-is-integrated)
- [Verification gates](#verification-gates)
- [Two local-environment traps](#two-local-environment-traps)
- [Pinned artifacts](#pinned-artifacts)
- [What is never weakened](#what-is-never-weakened)
- [Documentation](#documentation)
- [Conventions](#conventions)

## Before you write code

**Open an issue first for anything feature-shaped.** Every feature on this
project is specified through [Spec Kit](https://github.com/github/spec-kit)
before implementation, under `specs/NNN-name/`. A substantial change that
arrives without a spec is likely to need reworking rather than merging — not
because the code is wrong, but because the reasoning that produced it is not
recorded anywhere, and this project treats that record as part of the
deliverable.

Bug fixes and documentation corrections do not need this. Use your judgement:
if you can describe the change in a sentence and it does not add a capability,
just send it.

Build requirements: Go at the version in `go.mod`, and — for the race gate — a C
toolchain, since `-race` requires cgo. The desktop GUI additionally needs
OpenGL. The daemon and CLI are cgo-free and need neither.

## How work is integrated

Development is **trunk-based**. Maintainer work is committed directly onto
`main`: no feature branches, no internal pull requests. This is a one-to-two
developer project, and a pull request with no reviewer is ceremony rather than
review. The real gate is a single human halt immediately before anything is
pushed.

Outside contributions come as pull requests, which are reviewed here. The
[pull-request template](.github/PULL_REQUEST_TEMPLATE.md) asks for the
verification output described below.

CI runs on every push to `main`, but it reports *after* the fact rather than
blocking a merge. That is why the **local** run is the real gate: it must be
green before a change lands, and a red run is a stop rather than something to
push and sort out afterwards.

## Verification gates

Run these six, in the **foreground**, watched to completion. Never launch the
test suite in the background and poll for output: `go test` buffers a package's
output until that package completes, so a backgrounded run cannot be
distinguished from a dead one.

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

`gofmt` must print nothing. The race run excludes the cgo-only GUI entry point
and the Fyne widget package — races there live inside Fyne's own font cache
rather than this project's code — while `gui/viewmodel` stays race-tested and
the GUI is covered by the headless run.

`scripts/coverage-gate.sh` is the core-package coverage gate: six packages must
stay at or above 80 percent. CI runs this exact script, so the local number and
the CI number are one measurement rather than two approximations of it. Do not
substitute `go test -cover`; it reports per-package coverage and will disagree,
because the gate measures cross-package coverage with `-coverpkg`, where a
package's statements count as covered when *any* test in the tree reaches them.

**Report results honestly.** If a gate did not run, say which one and why. A
suite reported as passing when one gate was skipped is worse than no report.

## Two local-environment traps

Neither indicates a problem with the repository.

**golangci-lint refuses to start**, saying the Go version used to build it is
lower than the targeted Go version. Your *base* Go toolchain is older than the
`go` line in `go.mod`. `go version` can still report the newer one, because
`GOTOOLCHAIN=auto` upgrades transparently inside this repository — but
`go run <linter>@<ver>` builds the linter under *its* `go.mod`, which the older
base toolchain already satisfies, so no upgrade happens and the linter compiles
against the older version. Either upgrade your base Go install to match
`go.mod`, or force it for that one command:

```bash
GOTOOLCHAIN=go1.25.0 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```

Do not "fix" this by editing `.golangci.yml` or `go.mod`. CI installs the Go
version from `go.mod` as its base toolchain, and the pinned setup passes there.

**The race run needs a C toolchain.** `-race` requires cgo, so a machine with no
`gcc` on `PATH` fails with `cgo: C compiler "gcc" not found` before any test
runs. Install one — MSYS2 or MinGW-w64 on Windows — or rely on CI for the race
gate, and say explicitly that you did rather than reporting the suite as
passing.

## Pinned artifacts

These files change only with a **dated decision entry** in
[`CHANGELOG.md`](CHANGELOG.md) explaining why:

- `.github/workflows/**`
- `build/**`
- `Makefile`
- `.golangci.yml`
- the `go` and `toolchain` lines of `go.mod`
- `.gitattributes`, `.gitignore`, `LICENSE`
- `docs/INSTALL-windows.md`

The list is process infrastructure — things whose quiet drift would be
expensive and hard to notice. Note what is *not* on it:
`.github/ISSUE_TEMPLATE/**` and `.github/PULL_REQUEST_TEMPLATE.md` are unpinned,
because only `workflows/**` is named. Check the list rather than inferring from
the directory.

Cutting a `vX.Y.Z` tag always requires explicit authorization from a maintainer.

## What is never weakened

These test surfaces are not relaxed, skipped, or made conditional, whatever the
schedule pressure:

- clock injection — no direct `time.Now()` in engine code;
- timezone and DST resolution, including next-valid and first-occurrence
  behavior across transitions;
- forward-only, non-destructive store migrations;
- restart and catch-up recovery;
- goroutine termination under the race detector;
- local IPC access control.

If a change appears to require weakening one of these, that is a design
question, not a testing question. Raise it.

## Documentation

User-facing documentation lives in `docs/` and is the **single source of
truth**. It is published as a site —
[shruggietech.github.io/go-schedule](https://shruggietech.github.io/go-schedule/)
— served directly by GitHub Pages from the `docs/` folder on `main`, so the
Markdown you edit *is* the page that ships. There is no separate build step and
no generated output to commit; a page's `title` and `nav_order` front matter is
all that places it in the site's navigation.

Two rules keep the docs from drifting, and `scripts/docs-check.sh` enforces
both:

- **Subdirectory READMEs are pointers, not copies.** A README under, say,
  `test/scripts/` signposts the relevant `docs/` page rather than restating it.
  Duplicated prose rots; a pointer cannot.
- **Links out of `docs/` are absolute.** Because the site is rooted at `docs/`,
  a relative link that climbs out of it (`../README.md`) would 404 once
  published. Reference anything outside `docs/` by its full
  `https://github.com/…` URL instead.

Run the check before pushing any documentation change (CI runs the same script):

```bash
sh scripts/docs-check.sh
```

## Conventions

Commit messages follow Conventional Commits with the feature number where one
applies:

```text
fix(006): parse anchor timestamps portably on BSD date (macOS)
feat(005): assign a task to a group from the GUI editor
chore(release): prepare v0.5.0
```

Explain *why* in the body when the reason is not obvious from the diff. The
commit message is where a deviation from the standard workflow is recorded,
since there is no pull-request description to hold it.

Code follows ordinary Go idiom, matching the surrounding file's comment density
and naming. Internal scheduling is UTC; per-task IANA timezones handle
presentation and DST. Logging is structured, via `log/slog`.

Markdown follows the house style: a single `#` heading, hard-wrapped at 80
columns, prose where prose belongs and bullets only for genuinely enumerable
things, and language-tagged code fences.
