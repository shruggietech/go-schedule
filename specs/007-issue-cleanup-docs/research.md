# Research: Issue cleanup, README refresh, and documentation completion

**Feature**: `007-issue-cleanup-docs`

**Date**: 2026-07-23

Three defects, each traced to a specific line rather than described, plus the
options weighed on each.

## Why the MSI leaves the CLI unreachable

`build/windows/goschedule.wxs` declares the three binaries as `<File>` elements
inside `ComponentGroup Id="AppComponents"` and installs them under
`INSTALLFOLDER` (`C:\Program Files\go-schedule\`). There is no `<Environment>`
element anywhere in the package, so nothing ever touches `PATH`.

The reason it survived several releases is worth recording, because it is the
kind of blind spot that recurs: every machine where the product is developed or
tested already has the install directory on `PATH`, put there by hand. The
defect is invisible from inside the project and unmissable from outside it.

### Options weighed

| Option | Verdict |
| --- | --- |
| `<Environment>` on the CLI's component | **Chosen.** Component reference-counting gives correct install, upgrade, and uninstall behavior with no custom action. |
| A custom action editing `PATH` | Rejected. Hand-rolled `PATH` editing is the classic source of duplicated and truncated `PATH` values, and it has to implement rollback itself. |
| Per-user `PATH` (`System="no"`) | Rejected. The package is `perMachine` and registers a system service; a per-user entry would be written for whoever ran the installer and be invisible to everyone else on the machine. |
| Document the full-path invocation everywhere | Rejected. It is what the current Windows guide does, and it is the reason the README and `docs/test-scripts.md` disagree with reality. It also makes every copied command longer and platform-specific for no benefit. |

`Part="last"` rather than `Part="first"`: appending means an existing tool of
the same name on the machine keeps winning, which is the conservative choice
for an installer that runs elevated.

## Why `service status` reports a lie

`internal/service/service.go` handles the `status` action by calling
`svc.Status()` on the `kardianos/service` handle. On Windows that reaches
`lowPrivSvc`, which opens the service handle with:

```text
SERVICE_QUERY_CONFIG | SERVICE_QUERY_STATUS | SERVICE_START | SERVICE_STOP
```

The service ACL grants Interactive Users `CCLCSWLOCRRC` — query config, query
status, enumerate dependents, interrogate, user-defined control, read control.
It does **not** grant `RP` (start) or `WP` (stop), correctly. `OpenService`
evaluates the whole requested mask at once, so the call fails with
`Access is denied` before any status is read. The helper's name is misleading:
it is low-privilege only relative to `SERVICE_ALL_ACCESS`.

What makes this worth fixing rather than documenting is that the error is
actively wrong. It states that access was denied, from which a reader
reasonably concludes the ACL forbids the query — when the ACL permits it.

### Options weighed

| Option | Verdict |
| --- | --- |
| A Windows-only status path in this repo | **Chosen.** Smallest change, no fork, no new dependency (`golang.org/x/sys` is already direct), and it leaves every other action on the library. |
| Patch or fork `kardianos/service` | Rejected. `lowPrivSvc` is shared with paths that legitimately need start and stop rights, so a correct upstream fix is a larger change than the one this defect needs, and it puts a fork on the critical path of every future upgrade. |
| Shell out to `sc.exe query` | Rejected. Parsing localized console output to answer a question the API answers directly, and it adds a process spawn to a read. |
| Tell users to run status elevated | Rejected. It is the current behavior, and it is what the issue is about. |

The stub-and-fall-through shape on non-Windows is deliberate: the Linux and
macOS paths keep calling `svc.Status()` on the same line of `Control` they
always did, so FR-008 holds by construction rather than by test.

## Why issue templates belong here specifically

`.github/` currently holds only `ci.yml` and `release.yml`. The argument for
forms is not generic hygiene; it is that this repository has already had
triage decided by fields nobody thought to supply:

- Issue #5 reproduces only via the MSI install path. A `go install` user never
  sees it, so without an install-method field the report reads as "the CLI is
  broken", which is not true for every user.
- Issue #6 turns entirely on whether the reporter is an administrator — a
  question nobody answers unprompted.
- Issue #3's version had to be reconstructed from its title.
- Issue #4 arrived with no version, OS, or install method at all.

A project shipping a daemon, a CLI, and a GUI across three platforms with two
install paths cannot triage from prose. YAML forms are chosen over legacy
Markdown templates for one reason: `required: true` is enforced, whereas a
Markdown heading is a suggestion.

The body shape is lifted from issues #5 and #6 themselves — Summary,
Reproduction, Evidence, Root cause, Impact, Suggested fix. The *Evidence*
section is what made those two actionable, because it carried pasted output
rather than a description of output, and the form says so explicitly.

## Documentation gaps

Read as a set, the current documentation describes a project that is still
being built:

- `TODO.md` has every checkbox unticked though the roadmap is delivered, and it
  still lists event triggers under Phase 6 — a feature removed outright by
  feature 004, including its store migration. Leaving it advertises something
  that does not exist.
- `README.md` sends command-level questions into `specs/001-task-scheduler/`,
  which are spec artifacts. It also labels its project layout "(target)" and
  prefixes its feature list with build-log checkmarks.
- Only Windows has an install guide. Linux and macOS get four README lines
  each, and the macOS desktop bundle's boot-persistence caveat — the
  auto-started daemon keeps running but does not survive a reboot unless the
  service is registered — appears only in a README blockquote.
- There is no `CONTRIBUTING.md`, so the trunk-based workflow, the six CI-parity
  commands, the two local-environment traps, and the pinned-artifact rule live
  only in `CLAUDE.md`, which is addressed to an agent rather than a person.
- There is no `SECURITY.md` and no `CODE_OF_CONDUCT.md`; GitHub surfaces both
  absences on the repository page.

## Sources

- `build/windows/goschedule.wxs`, `build/windows/verify_wxs.ps1`
- `internal/service/service.go`; `kardianos/service@v1.2.4` `service_windows.go`
- `.github/workflows/{ci,release}.yml`
- Issues #3 through #7 on `shruggietech/go-schedule`
- `CLAUDE.md` (pinned-artifact list, verification procedure)
- `.specify/memory/constitution.md` v2.0.0
