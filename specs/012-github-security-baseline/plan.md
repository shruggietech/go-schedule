# Implementation Plan: GitHub Security Baseline

**Branch**: `main` (trunk-based) | **Date**: 2026-08-26 | **Spec**:
`specs/012-github-security-baseline/spec.md`

**Input**: Feature specification from
`specs/012-github-security-baseline/spec.md`

## Summary

Restore the repository's private security-reporting promise and add a
reviewable GitHub-native detection baseline without changing product behavior.
Locally, update `SECURITY.md`, add an advanced CodeQL workflow for Go, and
extend the existing offline automation policy with contract fixtures for the
workflow's actions, triggers, permissions, and required analysis steps. After
the single autopilot halt and explicit authorization, publish the local commit,
enable each supported repository security setting, and collect non-sensitive
hosted evidence.

## Technical Context

**Language/Version**: POSIX shell; YAML; Markdown; Go version selected by
`go.mod` for the hosted analysis build

**Primary Dependencies**: GitHub Actions; `github/codeql-action@v4`;
`actions/checkout@v7`; `actions/setup-go@v7`; GitHub REST API via `gh api`

**Storage**: Repository files and GitHub repository settings; no application
storage changes

**Testing**: `test/scripts/automation-check_test.sh`, ShellCheck, offline
`scripts/automation-check.sh`, and the canonical eight-gate
`sh scripts/verify.sh all`

**Target Platform**: GitHub-hosted Ubuntu runners and POSIX-compatible local
shells; hosted setting activation on `shruggietech/go-schedule`

**Project Type**: Cross-platform Go desktop/service repository with
repository-owned security automation

**Performance Goals**: No application runtime impact; offline automation
contract remains fast enough for the existing aggregate gate; CodeQL runs at
ordinary GitHub Actions cadence rather than on product execution paths

**Constraints**: No product or dependency changes; no remote mutation before
authorization; least-privilege workflow permissions; no repository secrets;
manual cgo-free headless build; missing GitHub-plan capabilities reported
individually rather than treated as passing

**Scale/Scope**: Two GitHub issues (#38 and #39), one new workflow, one policy
script, one fixture suite, one policy document, feature artifacts, and
repository-level security settings

## Constitution Check

*GATE: Passed before research and re-checked after design.*

| Gate | Result | Evidence |
| --- | --- | --- |
| Code quality | PASS | Shell changes remain POSIX, fail closed, and receive ShellCheck plus fixture coverage. No Go implementation changes. |
| Testing standards | PASS | Behavioral policy changes are test-first and cover all five required drift classes. The full race/coverage aggregate remains mandatory. |
| UX consistency | PASS | The reporting journey has one private primary route, one fallback, named ownership, and the existing one-week expectation. Product UI is unchanged. |
| Performance | PASS | No scheduler path changes. The policy check is static and offline; CodeQL executes only in hosted automation. |
| Autonomous execution | PASS | Scope traces to #38/#39. Work proceeds through the full Spec-Kit sequence with one halt immediately before push and repository-setting mutation. |
| Pinned artifacts | PASS | The new `.github/workflows/codeql.yml` and changed automation policy receive dated Unreleased decisions in `CHANGELOG.md`. |
| Integration | PASS | Work stays on local `main`; CI-parity must pass before the local commit and pre-push halt. No pull request or release is created. |
| Safety surfaces | PASS | Clock, timezone, migrations, recovery, concurrency, and IPC are untouched; their existing verification gates still run. |

**Post-design re-check**: PASS. Advanced CodeQL setup does not weaken the
existing CI contract, the manual `CGO_ENABLED=0 go build ./...` analysis build
avoids a new native desktop dependency, and hosted mutations remain beyond the
local implementation boundary until authorization.

## Project Structure

### Documentation (this feature)

```text
specs/012-github-security-baseline/
├── checklists/
│   ├── requirements.md
│   └── security.md
├── contracts/
│   ├── codeql-workflow.md
│   └── security-settings.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

### Repository files

```text
.github/workflows/
├── ci.yml                         # unchanged
└── codeql.yml                     # new advanced CodeQL workflow
scripts/
└── automation-check.sh            # extend action + workflow contract
test/scripts/
└── automation-check_test.sh       # extend temporary fixture regressions
SECURITY.md                        # reporting fallback and triage ownership
CHANGELOG.md                       # changed entry and dated decisions
CLAUDE.md                          # active feature plan pointer
```

**Structure Decision**: Keep the feature entirely in repository policy,
automation, and Spec-Kit artifacts. There is no new application package,
database entity, API, or dependency.

## Complexity Tracking

No constitution violations require justification.
