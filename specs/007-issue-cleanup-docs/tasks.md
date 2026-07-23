# Tasks: Issue cleanup, README refresh, and documentation completion

**Feature**: `007-issue-cleanup-docs` | **Date**: 2026-07-23

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md)

Tasks are grouped by user story so each group can be completed and demonstrated
on its own. `[P]` marks tasks that touch disjoint files and may run in parallel.

## Phase 1: US1 — the documented commands work after an MSI install (P1)

- **T001** Add `<Environment Id="PathEnv" Name="PATH" Value="[INSTALLFOLDER]"
  Permanent="no" Part="last" Action="set" System="yes" />` to the `Gosched`
  component in `build/windows/goschedule.wxs`, with a comment explaining why the
  element lives on the CLI's component. **Pinned artifact.**
- **T002** Extend `build/windows/verify_wxs.ps1` with an assertion that the
  `PATH` environment element is declared. **Pinned artifact.**
- **T003** Verify T002 fails when the element is absent (remove, run, restore).
- **T004** Rewrite `docs/INSTALL-windows.md` CLI examples as bare `gosched`,
  add the new-shell caveat, and keep the full-path form only as the fallback
  for an already-open shell. **Pinned artifact.**

## Phase 2: US2 — an ordinary user can query service status (P1)

- **T005** Add `internal/service/status_windows.go`: open the SCM with
  `SC_MANAGER_CONNECT`, the service with `SERVICE_QUERY_STATUS`, query, and map
  `ERROR_SERVICE_DOES_NOT_EXIST` to the existing not-installed result.
- **T006** Add `internal/service/status_other.go`: a stub reporting "not
  handled" so `Control` falls through unchanged on Linux and macOS.
- **T007** Route the `status` action in `internal/service/service.go` through
  the platform hook, reusing `statusString` unmodified.
- **T008** [P] Add `internal/service/service_test.go` — table test for
  `statusString`, including the unknown case.
- **T009** [P] Add `internal/service/status_windows_test.go` — a query against
  a service name that cannot exist must report not-installed rather than an
  access error.
- **T010** Manually confirm from a non-elevated shell that `service status`
  answers and `service start` still refuses.

## Phase 3: US3 — reports arrive triaged (P2)

- **T011** [P] `.github/ISSUE_TEMPLATE/bug_report.yml` with version, component,
  install method, OS, and elevation all required, and the Summary /
  Reproduction / Evidence / Impact body shape.
- **T012** [P] `.github/ISSUE_TEMPLATE/feature_request.yml` — problem before
  solution, plus spec traceability.
- **T013** [P] `.github/ISSUE_TEMPLATE/config.yml` — `blank_issues_enabled:
  false` and contact links.
- **T014** [P] `.github/PULL_REQUEST_TEMPLATE.md` — honest about trunk-based
  development; asks for the CI-parity result.

## Phase 4: US4 — the documentation set (P2)

- **T015** Rewrite `README.md` to house style: single H1, front-matter block,
  TOC, 80-column wrap, prose density, Mermaid architecture diagram, quick start
  with bare `gosched`, real project layout, `specs/` links confined to a
  development section.
- **T016** [P] `docs/cli.md` — every command from `internal/cli/*.go`, with
  flags, an example, and elevation requirements for the service subcommands.
- **T017** [P] `docs/INSTALL-linux.md`.
- **T018** [P] `docs/INSTALL-macos.md`, including the desktop bundle's
  boot-persistence caveat.
- **T019** [P] `docs/README.md` — index separating user guides from maintainer
  material.
- **T020** [P] Rewrite `TODO.md` to delivered state; remove event triggers.
- **T021** [P] `CONTRIBUTING.md` — trunk-based workflow, spec-kit requirement,
  six CI-parity commands, both local-environment traps, pinned artifacts.
- **T022** [P] `SECURITY.md` — supported versions, private reporting, and the
  threat model the project actually holds.
- **T023** [P] `CODE_OF_CONDUCT.md` — Contributor Covenant 2.1.

## Phase 5: release preparation and verification

- **T024** `CHANGELOG.md` `[0.6.0]` — narrative lead, Fixed / Added / Changed,
  and two dated decision entries for the pinned artifacts.
- **T025** Resolve every relative link in `README.md` and `docs/**`.
- **T026** Re-read every authored document against the house-style rules.
- **T027** Run all six CI-parity gates in the foreground, watched to
  completion.
- **T028** Work the checklist in `checklists/requirements.md`.
- **T029** Commit, then halt before push with the CI result, the two pinned
  decisions, the #5 verification gap, and the proposed v0.6.0 release.

## Dependencies

- T003 depends on T002; T004 depends on T001 being the reason it is correct.
- T007 depends on T005 and T006; T008–T010 depend on T007.
- T015 depends on T016–T019 existing as link targets.
- T024–T029 depend on everything above.
- Phases 1–4 are independent of each other and could be delivered separately.
