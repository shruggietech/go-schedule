# Requirements Checklist: Issue cleanup, README refresh, and documentation completion

**Purpose**: Verify each functional requirement is satisfied and that nothing
adjacent regressed — particularly the privilege restrictions that are *correct*
and must survive a fix aimed at loosening a neighbouring one.

**Created**: 2026-07-23

**Feature**: [spec.md](../spec.md)

## Installer

- [ ] CHK001 `<Environment>` is declared on the component whose `KeyPath` is
      `gosched.exe`, not on the daemon or a standalone component (FR-001).
- [ ] CHK002 The element carries `System="yes"`, `Permanent="no"`,
      `Part="last"`, `Action="set"`, `Name="PATH"` (FR-001, FR-002).
- [ ] CHK003 `verify_wxs.ps1` fails when the element is absent — verified by
      temporarily removing it, not by reading the assertion (FR-003).
- [ ] CHK004 No custom action edits `PATH`; the entry's lifetime is entirely
      component reference-counting (FR-002).

## Service status

- [ ] CHK005 The Windows query opens the SCM with `SC_MANAGER_CONNECT` only
      and the service with `SERVICE_QUERY_STATUS` only (FR-005).
- [ ] CHK006 No start or stop right appears anywhere in the status path
      (FR-005, FR-006).
- [ ] CHK007 `start`, `stop`, `restart`, `install`, `uninstall` still route
      through the existing library path and still require elevation (FR-006).
- [ ] CHK008 Status wording — running, stopped, not-installed — is byte-for-byte
      what it was, on every platform (FR-007).
- [ ] CHK009 The non-Windows path still calls `svc.Status()` on the same line
      of `Control`; it is not reimplemented (FR-008).
- [ ] CHK010 A missing service reports not-installed, not an access error
      (FR-007).

## Issue and PR templates

- [ ] CHK011 `blank_issues_enabled: false` is set (FR-009).
- [ ] CHK012 Version, component, install method, OS, and elevation are all
      `required: true` on the bug form (FR-010).
- [ ] CHK013 The bug form asks for pasted output, and says so, rather than for
      a description of output (FR-010).
- [ ] CHK014 The bug form notes that `gosched --version` and the daemon version
      from `gosched health` can differ after a partial upgrade (FR-010).
- [ ] CHK015 The feature form asks for the problem before the solution and for
      traceability to the master specification (FR-011).
- [ ] CHK016 The PR template describes the trunk-based reality and does not
      imply a review gate that internal work does not use (FR-012).

## Documentation

- [ ] CHK017 A newcomer reaches a first running task from `README.md` without
      opening a `specs/` artifact (FR-013).
- [ ] CHK018 `docs/cli.md` covers every command and subcommand the binary
      exposes, checked against `internal/cli/*.go` (FR-014).
- [ ] CHK019 `docs/cli.md` states which service subcommands need elevation and
      which do not, matching post-fix behavior (FR-014).
- [ ] CHK020 An install guide exists for Linux, macOS, and Windows, each
      standing alone (FR-015).
- [ ] CHK021 `docs/README.md` separates user-facing guides from maintainer
      material (FR-016).
- [ ] CHK022 `TODO.md` lists no removed feature — specifically, no event
      triggers (FR-017).
- [ ] CHK023 `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md` exist and
      the first carries the six CI-parity commands and both local traps
      (FR-018).
- [ ] CHK024 Every relative link in `README.md` and `docs/**` resolves.
- [ ] CHK025 Every authored document: single H1, headings blank-line-delimited,
      80-column wrap, wrap-safe continuation lines, language-tagged fences, TOC
      anchors matching real heading slugs.

## Process

- [ ] CHK026 `CHANGELOG.md` carries a dated decision for the `build/**` change
      and one for `docs/INSTALL-windows.md` (FR-019).
- [ ] CHK027 The pinned-artifact list in `CLAUDE.md` was read, and the
      conclusion that `.github/ISSUE_TEMPLATE/**` is unpinned is recorded
      rather than assumed (FR-019).
- [ ] CHK028 All six CI-parity gates green, run in the foreground (SC-008).

## Notes

- CHK003 is a real removal-and-rerun, not a read. An assertion that cannot fail
  is worse than no assertion, because it reports confidence it has not earned.
- CHK007 is the one most likely to be skipped and the most costly to get wrong:
  the fix loosens one privilege check, and the neighbouring checks that must
  stay tight are not exercised by any automatic test.
