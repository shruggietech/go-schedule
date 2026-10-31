# Implementation Plan: S128 Prelaunch Baseline

**Branch**: `codex/128-prelaunch-baseline` | **Date**: 2026-10-31 | **Issue**: [#408](https://github.com/shruggietech/go-schedule/issues/408) | **Spec**: [spec.md](spec.md)

## Summary

Replace active historical conversion machinery with direct current-state initialization and current-only qualification. The owner explicitly states no users/no migrations and approved this scope; preserve current developer data, explicit reset control, security and immutable historical evidence.

## Technical Context

**Language/Version**: Existing Go 1.26, TypeScript/React, POSIX shell, PowerShell and Markdown toolchains.
**Dependencies**: Existing SQLite 1.60.1/libc 1.77.1, Wails and site dependencies, unchanged.
**Storage**: Current application schema 24, current preference version 1, one direct current metrics baseline; no migration registry/converter.
**Testing**: Meaningful red/green bootstrap/preference/current qualification contracts, current integration/race/coverage, shell and PowerShell compliance, native desktop and site browser tests, complete foreground canonical verification.
**Platforms**: Windows, Linux and macOS as existing CI supports; attended desktop observation remains unperformed under the standing waiver.
**Performance**: Startup changes measured with a focused fresh-store benchmark; scheduling hot paths and dispatch budget unchanged.
**Constraints**: No silent resets, no visible consoles, no historical artifact rewrites, no tag/release or cluster/dependency work.
**Scope**: One coherent native issue #408 under milestone 18, Maintenance - Prelaunch baseline.

## Constitution Check

The old 3.1.0 requirement for forward migration and migration qualification conflicts with the latest explicit owner instruction. The owner-approved S128 includes a formal 4.0.0 amendment (a governance MAJOR because mandatory historical compatibility is removed), with dated Sync Impact Report/history and synchronized current guidance. This resolves the conflict rather than silently reinterpreting it. No product release/version is selected by that governance version.

Principles I-IV retain Go quality, current persistence/security/race/coverage, deterministic scheduling and measured performance. Principle V retains Spec-Kit analyze, foreground canonical gates and feature-branch PR integration. Explicit push/PR authorization satisfies the normal pre-push approval boundary. Issue completion stays functional.

Pre-research and post-design checks: PASS under the explicit owner-authorized amendment; no unresolved constitution conflict or new exemption.

## Research decisions

1. Capture accepted pre-S128 current DDL privately before changes. The capture has 22 tables and 33 indexes. Embed that final schema directly, with no ALTER, DROP or historical data-copy steps.
2. Retain application marker 24. Fresh initialization writes one current marker atomically; existing complete current databases (including their existing marker rows) reopen without rewriting state. Reject unsupported/foreign/incomplete state before persistent journal/identity writes. A private initializer seam supports meaningful rollback failure testing.
3. Use current preferences only: remove Fyne paths/imports, transition type and Upgrade card, including CLI/QA injection. Retain version 1, current defaults, validation, locks, atomic persistence and explicit reset. Remove the missing-popup historical fallback; existing validation then rejects incomplete current content. Unknown unrelated JSON keys remain ignored, without interpretation.
4. Test-script metrics helpers also convert old schemas. Establish one complete current metrics layout for both heartbeat and system writers, preserving shared explicit database paths. Use a new current marker 4 to distinguish this unified complete layout; reject obsolete2/3 state with an empty-path/reset instruction rather than convert it. Fresh/current readers and writers share the current contract.
5. Retire historical baseline-MSI construction, upgrade session/retained-profile observations and exact old-package identities from active qualification. Retain ordinary MSI replacement/exclusion mechanics and current install/repair/maintenance/uninstall, access, rollback, identity, hash and ownership assertions.
6. The current full release-evidence packet becomes a direct current contract revision 2 without converters. Immutable schema 1 packets remain historical and cannot be relabeled current qualification. Remove obsolete Fyne-era measurement branches where they exclusively exist for old candidates; retain current Wails measurement/identity rules.
7. Current documentation describes the S128 source baseline as unreleased. Published v1.7.0 notes, assets, comparison bindings, historical specs and 209 frozen inputs remain unchanged.
8. No new reset framework, schema-migration tooling, dependency update, cluster work or historical incident campaign. Supersede the old conditional S128 allocation on #404 with the owner's revised #408 scope; #400/#404 remain independent.
9. Installed checklist prerequisite requires plan.md, despite checklist-before-plan ordering; supported PathsOnly resolution supplied the real identity before planning. The optional agent-context hook runs after this plan exists.
10. Spec-Kit Phase 0 dispatched three read-only research agents for store, preferences and policy/qualification unknowns. Their concrete findings are retained privately and consolidated here. Implementation tasks have disjoint file groups and follow the installed tasks/implement parallel execution workflow; only focused checks run during those groups, with integrated verification after all changes settle.

11. Complete verification exposed calendar-sensitive desktop preview assertions (#409) on the eve of New York's 2026 repeated hour. Keep production DST conversion unchanged, inject the preview current-time dependency with its existing production default, and use fixed ordinary/past/equal/first/last expectations. This is a bounded verification correction, not a new scheduling feature.

12. The coverage gate exposed #410: current helper preflight can misclassify transient SQLite contention as an unsupported format. Reuse the established bounded retry/timeout behavior for read-only validation, retain the underlying diagnostic, and add a deterministic lock/refusal contract for both twins. Do not change timeouts or row-count requirements; preserve the first failure independently of #355/#356.

## Structure and ownership

- Store group: internal/store/, current metrics helpers/readers under test/scripts/, and test/integration/testscripts_test.go.
- Preference group: internal/userprefs/, preference CLI/QA wiring, desktop settings/current preference UI and corresponding fixtures/tests. Generated Wails bindings are ignored; regenerate through the normal build.
- Policy/qualification group: constitution and current docs/context, CI/automation contracts, installer/session/attended collector/releasegate source and their current synthetic tests.
- S128 artifacts: spec, plan, research, data-model, contracts, quickstart, tasks, analysis, delivery and verification plus requirements/domain checklists.
- Private evidence: .test-temp/s128, including accepted schema, first failures, focused/full logs, protected-input manifests and bounded review ledger.

## Verification sequence

Create failing current-contract tests before each behavior change, retaining original red logs. Preserve useful assertions from migration-named tests instead of deleting whole groups indiscriminately. Compare the resulting current application schema and behavior with the accepted 55-object baseline; test atomicity and non-mutating refusal. Measure fresh startup before/after. Run relevant metrics twins, preference and qualification contracts, then all Go suites, unchanged canonical eight gates and current desktop/site browser cases. Add shellcheck/PowerShell compliance for changed scripts. Use headless owned fixtures only, retaining 209 frozen inputs and all new receipts.

No new test merely mirrors prose. No test or review is an issue closure condition. Source/code acceptance and public release acceptance remain distinct.

## Complexity Tracking

The direct current metrics format 4 and release-evidence contract 2 remove ambiguous historical formats without adding migration compatibility. The constitution 4.0.0 policy removal is explicit owner scope. No other architectural exception.
