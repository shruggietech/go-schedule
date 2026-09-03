# Specification Lifecycle

Feature specifications describe implementation maturity separately from publication. Use exactly one state:

- `Draft`: requirements are still being shaped and at least one task remains actionable.
- `Ready`: requirements and tasks are executable, but implementation has not started.
- `In Progress`: implementation has both completed and actionable required tasks.
- `Implemented`: required implementation and verification are resolved in the repository. The `Delivery` field cites objective evidence.
- `Deferred`: intentionally postponed, with a recorded disposition.
- `Superseded`: replaced by another specification or decision, with a recorded disposition.
- `Abandoned`: intentionally ended without delivery, with a recorded disposition.

Publication bookkeeping (push, pull request, review, merge, and branch cleanup) is workflow evidence, not implementation work. Do not leave it as an open feature task. A historical task that was not executed must be resolved with an explicit waiver or supersession note, never silently represented as performed.

Run `sh scripts/spec-lifecycle-check.sh .` locally. The canonical automation gate runs the same check and rejects invalid states, missing delivery evidence, contradictory task state, and inventory drift.

## Inventory

| Feature | State | Delivery |
| --- | --- | --- |
| 001-task-scheduler | Implemented | releases `v0.1.0` through `v0.3.0` |
| 002-gui-task-editor-ux | Implemented | commit `c10aac4` |
| 003-gui-editor-refinements | Implemented | commit `a21c22e` |
| 004-rebrand-gui-overhaul | Implemented | PR [#1](https://github.com/shruggietech/go-schedule/pull/1) |
| 005-gui-task-fidelity | Implemented | commit `2603986` |
| 006-maintainer-test-scripts | Implemented | commit `78b33b8` and release `v0.5.0` |
| 007-issue-cleanup-docs | Implemented | commit `eb76c8e` and release `v0.6.0` |
| 008-cron-interop | Implemented | commit `a71172a` and release `v0.7.0` |
| 009-ci-bench-p99 | Implemented | commit `95eca21` |
| 010-docs-site-pages | Implemented | commit `3e99d6a` and release `v0.8.0` |
| 011-maintainer-automation-baseline | Implemented | PR [#43](https://github.com/shruggietech/go-schedule/pull/43) |
| 012-github-security-baseline | Implemented | PR [#44](https://github.com/shruggietech/go-schedule/pull/44) |
| 013-pr-governance | Implemented | PR [#45](https://github.com/shruggietech/go-schedule/pull/45) |
| 014-gui-activity-clarity | Implemented | PR [#46](https://github.com/shruggietech/go-schedule/pull/46) |
| 015-docs-dark-theme | Implemented | PR [#47](https://github.com/shruggietech/go-schedule/pull/47) |
| 016-windows-installer-identity | Implemented | PR [#48](https://github.com/shruggietech/go-schedule/pull/48) |
| 017-gui-info-navigation | Implemented | PR [#49](https://github.com/shruggietech/go-schedule/pull/49) |
| 018-cron-string-conversion | Implemented | PR [#53](https://github.com/shruggietech/go-schedule/pull/53) |
| 019-dual-syntax-input | Implemented | PR [#54](https://github.com/shruggietech/go-schedule/pull/54) |
| 020-gui-dual-syntax | Implemented | PR [#55](https://github.com/shruggietech/go-schedule/pull/55) |
| 021-dual-syntax-docs | Implemented | PR [#56](https://github.com/shruggietech/go-schedule/pull/56) |
| 022-activity-diagnostics | Implemented | PR [#57](https://github.com/shruggietech/go-schedule/pull/57) |
| 023-ordinal-weekday-cron | Implemented | PR [#58](https://github.com/shruggietech/go-schedule/pull/58) |
| 024-last-weekday-cron | Implemented | PR [#59](https://github.com/shruggietech/go-schedule/pull/59) |
| 025-monthly-calendar-parity | Implemented | PR [#60](https://github.com/shruggietech/go-schedule/pull/60) |
| 026-cron-expression-breadth | Implemented | PR [#61](https://github.com/shruggietech/go-schedule/pull/61) |
| 027-dst-intent | Implemented | PR [#62](https://github.com/shruggietech/go-schedule/pull/62) |
| 028-cron-parity-closure | Implemented | PR [#63](https://github.com/shruggietech/go-schedule/pull/63) |
| 029-v09-closure | Implemented | PR [#65](https://github.com/shruggietech/go-schedule/pull/65) |
| 030-ipc-admin-group | Implemented | PR [#68](https://github.com/shruggietech/go-schedule/pull/68) |
| 031-reboot-startup-trigger | Implemented | PR [#69](https://github.com/shruggietech/go-schedule/pull/69) |
| 032-brand-system-integration | Implemented | PR [#71](https://github.com/shruggietech/go-schedule/pull/71) |
| 033-task-completion-chains | Implemented | PR [#78](https://github.com/shruggietech/go-schedule/pull/78) |
| 034-v090-release-cut | Implemented | review branch `codex/034-v090-release-cut`; local verification completed 2026-08-30 |
| 035-windows-msi-local-group | In Progress | implementation and verification in progress on `codex/083-fix-msi-local-group` |
| 036-ipc-access-denied-recovery | Implemented | native Windows diagnosis, deterministic recovery, and canonical eight-gate verification passed 2026-09-01 |
| 037-windows-release-polish | Implemented | review branch `codex/037-windows-release-polish`; local verification completed 2026-08-31 |
| 038-windows-installed-core-recovery | Implemented | PR [#95](https://github.com/shruggietech/go-schedule/pull/95), local canonical verification, and hosted LocalSystem run [33508591264](https://github.com/shruggietech/go-schedule/actions/runs/33508591264) |
| 039-windows-setup-lifecycle-control | Implemented | canonical eight-gate verification and local WiX 6.0.2 compiled-MSI inspection in `verification.md`; hosted silent evidence runs in pull-request CI; #94 retains attended Windows 11 acceptance for #97/#98 |
