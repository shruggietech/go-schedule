# Quickstart: Validate the Maintainer Automation Baseline

This guide is executed after implementation. It performs no push, tag, release,
or repository-setting change.

## Prerequisites

- Go at the version selected by `go.mod`
- Git
- POSIX shell (Git Bash on Windows)
- C toolchain for the race detector
- A clean checkout on local `main`

Make is optional because the canonical direct script remains available.

## Pre-change baseline (T002)

Captured on 2026-08-26 before pinned automation files are edited:

- `.github/workflows/ci.yml` triggers on pushes and pull requests targeting
  `main`, grants `contents: read`, and runs lint, three-platform race tests,
  coverage, six-target cross-builds, GUI build/tests, benchmarks, and docs.
- The CI benchmark artifact is named `engine-benchmarks`; the race matrix is
  `ubuntu-latest`, `macos-latest`, and `windows-latest`.
- `.github/workflows/release.yml` triggers only on `v*.*.*` tags, grants
  `contents: write`, and preserves the existing release names and file globs.
- Existing action majors are `actions/checkout@v4`, `actions/setup-go@v5`,
  `actions/upload-artifact@v4`, and `softprops/action-gh-release@v2`; all four
  declare the Node 20 action runtime.
- Existing Make verification commands duplicate lint, race-selection, and
  coverage behavior instead of delegating to one command source.

Implementation validation compares the final workflows and Makefile with this
baseline. Evidence is appended under the corresponding quickstart steps.

## 1. Inspect the declared contract

```sh
sh scripts/verify.sh list
```

Expected, in order:

```text
format
vet
lint
race
gui
coverage
docs
automation
```

## 2. Run automation contract regressions

```sh
sh test/scripts/automation-check_test.sh
```

Expected: the approved fixture passes; obsolete, unknown, missing, duplicate,
and extra cases fail inside the harness with their expected diagnostics; the
harness exits zero.

## 3. Check the real workflows offline

```sh
sh scripts/automation-check.sh .
```

Expected: all `uses:` references match the approved Node 24 major allowlist and
the verification manifest is exact.

Evidence (2026-08-26): `automation-check: OK - approved actions and 8-gate
manifest`. The release workflow diff changes only action major selectors. The
CI workflow retains its triggers, permissions, jobs, matrix, build and benchmark
commands, and `engine-benchmarks` artifact; its verification steps now delegate
to the driver and its lint job adds the automation-contract gate.

## 4. Run one named gate

```sh
sh scripts/verify.sh docs
```

Expected: only the documentation gate runs and reports success.

## 5. Prove failure propagation without touching tracked files

Invoke the contract-test fixture that substitutes a controlled failing child
gate:

```sh
sh test/scripts/automation-check_test.sh failure-propagation
```

Expected: the nested aggregate run stops at the injected gate and returns
non-zero; the harness recognizes that expected failure and exits zero.

## 6. Run the complete local contract

```sh
sh scripts/verify.sh all
```

Run in the foreground and watch it complete. Expected: all eight headings appear
once in contract order and every gate passes.

If Make is installed, verify the wrapper reaches the same driver:

```sh
make verify
```

## 7. Prove the checkout stayed clean

```sh
git status --short
```

Expected: no output beyond the already-reviewed feature changes present before
the verification run. The verification run itself introduces no additional
change.

## Local implementation evidence (2026-08-26)

- `sh scripts/verify.sh list` printed the exact eight-gate manifest in order.
- `sh test/scripts/automation-check_test.sh` passed every positive and negative
  fixture; its failure-propagation mode stopped at the injected `[format]` gate.
- `sh scripts/verify.sh docs` passed independently with 11 pages clean.
- `shellcheck` passed for the three new shell files.
- `sh scripts/verify.sh all` passed all eight gates in WSL Ubuntu 24.04 using
  the existing GCC 13.3 toolchain. WSL's DNS was unavailable, so the run used
  the existing Windows module-download cache as a local Go proxy and the
  contract-supported `GOTOOLCHAIN=local` override (Go 1.26.4). The default lint
  mode was also run separately on Windows and passed with the driver-derived Go
  1.25.0 toolchain.
- Coverage passed at engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%,
  catchup 87.5%, and logbus 91.1%.
- `make verify` was not run because Make is not installed; the wrapper is
  optional and its recipe was inspected as a one-line delegate.
- The post-verification `git status --short` matched the reviewed feature file
  set. Verification introduced no tracked or untracked repository change.

## 8. Hosted evidence after authorized push

After the autopilot pre-push halt and explicit operator authorization, inspect
the ordinary CI run. Expected:

- every job passes;
- no Node 20 action deprecation annotation appears;
- matrix, artifacts, and job names remain intact.

Do not tag or execute the release workflow as part of this feature. Its action
references and contract are validated statically before the halt.
