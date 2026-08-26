# Implementation Plan: Maintainer Automation Baseline

**Branch**: `011-maintainer-automation-baseline` (trunk-based on `main`) |
**Date**: 2026-08-26 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from
`specs/011-maintainer-automation-baseline/spec.md`

## Summary

Close issues #21 and #41 as one maintainer-automation slice. Upgrade every
Node 20-based action reference in `.github/workflows/ci.yml` and
`.github/workflows/release.yml` to its verified Node 24 major. Add one POSIX
shell verification driver whose named gates are consumed by the Makefile, CI,
contributor guidance, and autopilot protocol. Add an offline automation-contract
check that allowlists the audited action majors and compares the driver's gate
manifest with the required local gate set. Preserve all workflow triggers,
permissions, matrices, inputs, outputs, artifact names, and release selection.

## Technical Context

**Language/Version**: POSIX `sh`, GNU Make syntax, GitHub Actions YAML; Go
1.25.0 commands are orchestrated but product Go code is unchanged

**Primary Dependencies**: Go toolchain, Git, a POSIX shell (Git Bash on
Windows), Make as an optional wrapper, GitHub-hosted Actions

**Storage**: Repository files only; no runtime data or schema changes

**Testing**: POSIX shell contract tests with temporary fixtures, workflow syntax
inspection, controlled negative checks, and the complete repository CI-parity
suite

**Target Platform**: Linux, macOS, and Windows maintainer environments; GitHub
hosted Linux/macOS/Windows runners

**Project Type**: Cross-platform Go application with repository automation and
maintainer tooling

**Performance Goals**: The driver adds negligible orchestration overhead and
runs each required gate exactly once in aggregate mode

**Constraints**: No network is required for drift checks; no release execution;
no product behavior change; no hidden/skipped gate; workflows and Makefile are
pinned artifacts; one foreground process owns the aggregate run

**Scale/Scope**: Two workflows, four external action families, eight local
verification gates, one Makefile, three maintainer guidance documents, one
changelog entry, and shell contract fixtures

## Constitution Check

*GATE: Passed before research and passed again after design.*

- **I. Code Quality**: The new shell scripts use `set -eu`, explicit failures,
  small gate functions, stable diagnostics, and ShellCheck-compatible POSIX
  syntax. Product Go code is untouched. Existing Go formatting, vet, and lint
  gates remain mandatory.
- **II. Testing Standards**: The automation contract receives fixture-based
  regression tests. The complete race, GUI, and coverage gates remain part of
  aggregate verification and are run before commit. No safety-critical test is
  relaxed or excluded.
- **III. User Experience Consistency**: The maintainer-facing command has named
  gates, consistent stdout progress, actionable stderr failures, conventional
  exit codes, and explicit prerequisite errors.
- **IV. Performance Requirements**: Scheduler hot paths and benchmarks are not
  modified. The aggregate driver is sequential by design so buffered foreground
  test output remains attributable and every required gate is watched.
- **V. Autonomous Build-Phase Execution**: Scope traces to open issues #21 and
  #41. The artifacts make implementation task-complete for autopilot, retain the
  analyze gate, and require the single pre-push halt. No push, tag, or release is
  part of this plan.
- **Pinned-artifact rule**: `.github/workflows/ci.yml`,
  `.github/workflows/release.yml`, and `Makefile` change only with dated
  `CHANGELOG.md` decisions. This plan also updates the autopilot guidance that
  consumes their verification contract.

No constitution violation or complexity exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/011-maintainer-automation-baseline/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── automation-policy.md
│   └── verify-command.md
├── checklists/
│   ├── automation.md
│   └── requirements.md
└── tasks.md
```

### Repository files affected by implementation

```text
.github/workflows/
├── ci.yml                         # supported action majors + shared gate calls
└── release.yml                    # supported action majors only

scripts/
├── verify.sh                      # canonical named/aggregate gate driver
├── automation-check.sh            # offline action + gate-manifest policy check
├── coverage-gate.sh               # reused unchanged
└── docs-check.sh                   # reused unchanged

test/scripts/
└── automation-check_test.sh       # temporary-fixture contract regressions

Makefile                           # convenience wrappers delegate to verify.sh
CONTRIBUTING.md                    # canonical command and prerequisites
CLAUDE.md                          # autopilot verification command/context
docs/build-autopilot.md            # protocol verification command/breakdown
CHANGELOG.md                       # feature line + dated pinned decisions
```

**Structure Decision**: Keep automation in the repository's existing POSIX
shell and Make surfaces. `scripts/verify.sh` owns the command definitions and
gate manifest. CI selects individual modes for job/matrix parallelism, while
maintainers and autopilot invoke aggregate mode in one foreground process.
`scripts/automation-check.sh` is deliberately separate so it can validate the
driver's manifest and workflow allowlist without self-certifying the content it
checks.

## Implementation Phases

### Phase 1: Contract foundation

Define the eight-gate manifest and audited action allowlist in executable shell
contracts. Add temporary-fixture tests proving an old action major, an unknown
action, and a missing gate fail with actionable output while the current
contract passes without network access.

### Phase 2: Canonical local verification

Implement the named and aggregate modes in `scripts/verify.sh`. Make
`Makefile` verification targets thin wrappers, preserve `fmt` as an explicitly
mutating convenience command, and point contributor/autopilot documentation at
the aggregate command. Aggregate mode runs in the foreground, stops on first
failure, and leaves a valid checkout unchanged.

### Phase 3: Hosted workflow modernization

Update the four audited action families to Node 24 majors and route CI gate
steps through the matching driver modes. Preserve the release workflow's
behavior exactly. Add the offline automation-contract gate to ordinary CI so a
deprecated or unaudited future action reference cannot land silently.

### Phase 4: Evidence and pinned decisions

Run controlled negative contract tests, the quickstart, and all CI-parity gates.
Record the feature and dated decisions for the two workflows and Makefile in
`CHANGELOG.md`. Re-run aggregate verification after documentation changes, then
commit locally and halt before push.

## Complexity Tracking

No constitutional violation requires justification.
