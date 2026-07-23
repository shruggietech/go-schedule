# Implementation Plan: Issue cleanup, README refresh, and documentation completion

**Branch**: `007-issue-cleanup-docs` | **Date**: 2026-07-23 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/007-issue-cleanup-docs/spec.md`

## Summary

Close the three open issues and finish the documentation set. The installer
declares its install directory on the machine `PATH` so the documented commands
resolve by name; the Windows `service status` path opens the service control
manager and the service handle with the minimum rights the query needs, so an
unprivileged user gets a real answer instead of a misleading access error; and
the repository gains structured issue forms so reports arrive carrying the
facts that decide triage. Alongside those, `README.md` is rewritten to the
house Markdown style, a user-facing CLI reference and per-platform install
guides are added, and the stale roadmap and missing community-health documents
are brought into line with what the project actually is.

No scheduling behavior changes. The only Go code that changes is the service
control layer's status path.

## Technical Context

**Language/Version**: Go 1.25.0 (`go.mod`), PowerShell 7 for the WiX sanity
check, WiX v5 for the installer source, GitHub issue-form YAML, Markdown.

**Primary Dependencies**: `golang.org/x/sys/windows` (already a direct
dependency) for the least-privilege service query; `github.com/kardianos/service`
stays in place for every other service action.

**Storage**: N/A — nothing here touches the store or its migrations.

**Testing**: `go test -race`; a portable table test for the status-string
mapping and a Windows-tagged test for the not-installed versus access-denied
distinction. Windows is already in the CI race matrix.

**Target Platform**: Windows for the two defects; Linux, macOS, and Windows for
the documentation.

**Project Type**: Cross-platform daemon + CLI + desktop GUI.

**Performance Goals**: Unchanged. The dispatch-latency budget is untouched.

**Constraints**: `build/**` and `docs/INSTALL-windows.md` are pinned artifacts
and each requires a dated `CHANGELOG.md` decision. Status output wording must
not change on any platform. Elevation requirements for every non-status service
action must not be relaxed.

**Scale/Scope**: Two small Go files plus a routing change, one installer
element, one sanity-check assertion, four issue and PR templates, and eleven
Markdown documents.

## Constitution Check

*Re-checked after design. No violations; no Complexity Tracking entries.*

- **I. Code Quality.** The new code is two small platform-split files with a
  clear seam. The Windows path is additive; the non-Windows stub keeps the
  existing behavior literally unchanged rather than reimplementing it.
- **II. Testing Standards (non-negotiable).** No safety-critical surface is
  weakened. Clock injection, timezone and DST resolution, store migrations,
  restart and catch-up recovery, goroutine termination, and IPC access control
  are all untouched. Coverage: `internal/service` is not among the six packages
  in `scripts/coverage-gate.sh`, so the gate is unaffected, and tests are added
  rather than removed.
- **III. User Experience Consistency.** Status wording is preserved verbatim
  across platforms — the point of the fix is that the *answer* becomes
  reachable, not that it reads differently. The documentation change makes the
  written commands match the commands that work.
- **IV. Performance Requirements.** Not engaged.
- **V. Autonomous Build-Phase Execution.** Full spec-kit sequence, single halt
  before push. Both pinned-artifact changes carry dated decisions.

One point deserves stating rather than assuming, because issue #7 explicitly
asks for it: the pinned list in `CLAUDE.md` names `.github/workflows/**`, not
`.github/**`. Issue and PR templates are therefore unpinned and need no
decision entry. This was checked, not inferred.

## Project Structure

### Documentation (this feature)

```text
specs/007-issue-cleanup-docs/
├── spec.md
├── plan.md              # this file
├── research.md          # the three root causes, and the options weighed
├── quickstart.md        # how to verify each fix
├── checklists/
└── tasks.md
```

### Source Code (repository root)

```text
build/windows/
├── goschedule.wxs           # + <Environment> on the Gosched component  [pinned]
└── verify_wxs.ps1           # + PATH-declaration assertion              [pinned]

internal/service/
├── service.go               # status routes through the platform path
├── status_windows.go        # new — least-privilege SCM + service query
├── status_other.go          # new — non-Windows stub, falls through
├── service_test.go          # new — statusString mapping
└── status_windows_test.go   # new — not-installed vs access-denied

.github/
├── ISSUE_TEMPLATE/{bug_report,feature_request,config}.yml   # new
└── PULL_REQUEST_TEMPLATE.md                                 # new

docs/
├── README.md                # new — index
├── cli.md                   # new — user-facing command reference
├── INSTALL-linux.md         # new
├── INSTALL-macos.md         # new
└── INSTALL-windows.md       # bare `gosched` commands          [pinned]

README.md CONTRIBUTING.md SECURITY.md CODE_OF_CONDUCT.md TODO.md CHANGELOG.md
```

**Structure Decision**: No new packages and no new module dependencies. The
platform split inside `internal/service` follows the convention already used in
`internal/platform` (`platform_unix.go` / `platform_windows.go`), so the shape
is one a reader of this repository already recognizes.

## Design

### Installer `PATH` entry

The `<Environment>` element goes on the `Gosched` component — the one whose
`KeyPath` is `gosched.exe` — rather than on the daemon or a standalone
component. Binding it to the CLI is what makes its lifetime correct for free:
MSI reference-counts by component, so the entry is written when the CLI is
installed, replaced in place on a major upgrade, and removed when the CLI is
removed. `System="yes"` matches the `perMachine` package scope,
`Permanent="no"` is what removes it on uninstall, and `Part="last"` appends
rather than replacing `PATH` wholesale.

`verify_wxs.ps1` gains one assertion in the same shape as its existing checks.
It runs in CI before `wix build`, so a future edit that drops the element fails
the release rather than shipping the regression again.

### Least-privilege service status

`svc.Status()` in the upstream library reaches a helper that opens the service
with `SERVICE_QUERY_CONFIG|SERVICE_QUERY_STATUS|SERVICE_START|SERVICE_STOP`.
The Interactive Users ACE carries the query rights but not start or stop, so
`OpenService` fails before the status read is attempted. The mask is the whole
defect.

`Control` gains a platform hook. On Windows, `queryStatus` opens the SCM with
`SC_MANAGER_CONNECT` and the service with `SERVICE_QUERY_STATUS`, then queries.
On every other platform the stub reports "not handled" and `Control` falls
through to `svc.Status()` exactly as today — the non-Windows path is unchanged
rather than reimplemented, which is what keeps FR-008 true by construction.

`statusString` is reused unmodified, so running, stopped, and not-installed all
render exactly as they do now.

### Documentation

`README.md`, the new guides, and the refreshed roadmap are authored to the
house Markdown style: single H1, a labeled front-matter block, a manual table
of contents built on GitHub's automatic anchor slugs, prose where prose belongs
and bullets only for genuinely enumerable items, hard wrap at 80 columns with
wrap-safe continuation lines, and language-tagged fences.

`docs/cli.md` is written from `internal/cli/*.go` rather than from
`specs/001-task-scheduler/contracts/cli.md`. The contract is a spec artifact
describing what the CLI must do; the reference describes what it does. Deriving
the user doc from the contract would make the two indistinguishable and leave
the reference silently wrong the first time the implementation moved.

## Complexity Tracking

No constitution violations. Table intentionally empty.
