# Implementation Plan: Actionable Notification Conditions

**Branch**: `codex/095-notification-conditions` | **Date**: 2026-09-21 | **Spec**: [spec.md](spec.md)

## Summary

Extend the existing durable webhook policy and delivery pipeline with per-assignment problem conditions, durable per-task and channel evaluation state, structured process-start failure evidence, condition explanations, and opt-in daemon-health heartbeats. Evaluation remains inside the source-run transaction, while outbound work remains asynchronous and isolated from scheduler workers.

## Technical Context

**Language/Version**: Go 1.25, TypeScript 5, React 18

**Primary Dependencies**: Standard Go library, existing SQLite driver, Cobra, Wails v2, React, Vitest, Playwright

**Storage**: Existing daemon-owned SQLite database with forward-only migration v20

**Testing**: Go unit, migration, API, client, CLI, integration, race, frontend component, accessibility, and Playwright tests

**Target Platform**: Windows, macOS, and Linux daemon, CLI, API, and Wails desktop

**Project Type**: Local scheduler daemon with CLI, JSON API, MCP, and desktop application

**Performance Goals**: Constant work per matching channel assignment and bounded heartbeat scheduling; no additional scheduler worker occupancy

**Constraints**: Preserve current assignments, secrets, retry semantics, inheritance, and webhook compatibility; add no background service, message broker, or dependency

**Scale/Scope**: Existing task, group, channel, and 1,000-record terminal delivery bounds; failure thresholds at most 100; time settings at most 30 days except heartbeat at most 24 hours

## Constitution Check

- **I. Code Quality**: PASS. The evaluator is a focused deterministic function, persistence remains transactional, and the dispatcher retains one owned goroutine with cancellation.
- **II. Testing Standards**: PASS. Migration, state transition, restart, API, CLI, desktop, accessibility, integration, race, and full verification coverage are planned with injected time.
- **III. User Experience Consistency**: PASS. API values use seconds while CLI accepts Go durations; desktop labels explain each condition, inheritance, and next eligibility.
- **IV. Performance Requirements**: PASS. Evaluation is bounded by the effective assignment count and uses indexed state keys; heartbeat creation is bounded by enabled channels.
- **V. Autonomous Build-Phase Execution**: PASS. S095 is explicitly authorized, issue-traceable, spec-kit managed, review-branch based, and will publish through one PR.

Post-design check: PASS. Migration v20 is additive except for a data-preserving rebuild of the delivery table required to extend its event-kind check constraint. No new dependency or authority surface is introduced.

## Technical Decisions

1. Existing `on_failure` becomes threshold one so upgrades retain behavior.
2. One primary problem per run and channel prevents duplicate deliveries and yields an explainable state machine.
3. Policy fingerprints reset state when inherited or direct configuration changes, avoiding stale streaks without expensive descendant rewrites.
4. First problem and recovery transitions bypass quiet periods; routine success and reminders are suppressible.
5. Daemon health uses opt-in heartbeats. Absence detection belongs to the receiver because local outage notification is impossible while the daemon is stopped.
6. Heartbeats reuse the existing durable delivery and retry pipeline and never recursively react to delivery failure.

## Project Structure

```text
specs/095-notification-conditions/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/notification-conditions.md
├── checklists/requirements.md
└── tasks.md

internal/domain/
internal/store/
internal/executor/
internal/notification/
internal/api/server/
internal/api/client/
internal/cli/
internal/mcpmanage/
desktop/notifications/
desktop/frontend/src/notifications/
desktop/frontend/e2e/
docs/
```

**Structure Decision**: Extend the existing notification layers in place. No package or service boundary is added because policy evaluation belongs beside current transactional delivery creation and heartbeat delivery belongs to the existing dispatcher lifecycle.

## Complexity Tracking

No constitution violation requires justification. Rebuilding the delivery table is the smallest safe SQLite mechanism for expanding its checked event-kind vocabulary while preserving existing rows and indexes.
