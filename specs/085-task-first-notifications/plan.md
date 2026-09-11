# Implementation Plan: Task-First Notifications

**Branch**: `codex/085-task-first-notifications` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/085-task-first-notifications/spec.md`

## Summary

Extend the secret-free desktop notification snapshot with configured task and group coverage, then reorganize Notifications into a task-first overview and bounded recent results followed by three collapsed advanced disclosures. Preserve all existing channel, policy, delivery, polling, redaction, and bridge behavior while adding deterministic Go, React, Playwright, native Windows, and canonical verification for #232.

## Technical Context

**Language/Version**: Go 1.25.0 with toolchain 1.25.1, TypeScript 5.6.3, React 19.2.8, CSS

**Primary Dependencies**: Existing desktop notification service, React component catalog, Wails desktop host, Vite 8.2.2, Vitest 5.0.0, Testing Library, Playwright 1.63.0, axe-core 4.13.0

**Storage**: Existing daemon notification channels, assignments, and delivery records only; no schema or persistence change

**Testing**: Go unit and race tests, Vitest component and store tests, Playwright Chromium workflow and accessibility checks, native Windows production build, `sh scripts/verify.sh all`

**Target Platform**: Wails desktop on Windows with cross-platform source compatibility; 800 by 600 minimum viewport and 200 percent zoom reflow

**Project Type**: Desktop application with a React frontend and Go/Wails host

**Performance Goals**: Overview computation remains within the existing three-second desktop operation boundary; the first surface renders at most five recent results while diagnostics retain the bounded 200-record snapshot

**Constraints**: Local assets only, WCAG 2.2 AA, keyboard parity, secret-free snapshots, no dispatch, persistence, daemon API, authorization, or release-state changes

**Scale/Scope**: One desktop service projection, one notification page, three advanced disclosures, one bounded recent-results surface, existing 200-record detailed history

## Constitution Check

- **I. Code Quality**: Keep domain projection in the existing typed notification service and split page composition into small view helpers instead of extending one monolithic return block.
- **II. Testing Standards**: Add failing service and component regressions before implementation, retain delivery polling tests, and run focused, race, browser, native, and all eight canonical gates.
- **III. User Experience Consistency**: Reuse shared CardSection, Disclosure, StatusLabel, Notice, Field, Dialog, and button contracts from S083 and S084.
- **IV. Performance Requirements**: Bound the initial result count, reuse one workspace load, avoid new polling, and allow coverage projection to degrade explicitly without hiding the core snapshot.
- **V. Autonomous Build-Phase Execution**: Complete the full spec-kit sequence, review branch, authorized publication, hosted checks, permitted review rounds, and maintainer merge gate.

**Gate result before design**: PASS. No constitution conflict or unresolved clarification exists.

**Gate result after design**: PASS. The design adds a secret-free read projection and presentation state without persistence, dispatch, authorization, network-route, dependency, or release changes.

## Project Structure

### Documentation

```text
specs/085-task-first-notifications/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── contracts/
│   └── task-first-notifications.md
└── checklists/
    ├── requirements.md
    └── ux.md
```

### Source Code

```text
desktop/
├── notifications/
│   ├── model.go
│   ├── service.go
│   └── service_test.go
└── frontend/
    ├── src/
    │   ├── notifications/
    │   │   ├── model.ts
    │   │   ├── NotificationsPage.tsx
    │   │   └── NotificationsPage.test.tsx
    │   └── styles.css
    └── e2e/
        └── notifications.spec.ts

CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Extend the existing notification service and page. No new global state, route, dependency, or backend endpoint is warranted because all authoritative capability already exists behind the desktop backend interface.

## Decision Log

- **Task-first overview**: The first region reports active destinations, configured coverage, recent task outcomes, and one state-dependent next action. Webhook vocabulary does not appear in the page introduction.
- **Derived secret-free coverage**: The desktop service resolves direct group and effective task assignments through existing backend methods and returns only names, scope identity, outcome kinds, and destination counts. Per-scope projection errors set an explicit incomplete flag while preserving the usable core snapshot.
- **Bounded recent results**: Show the five newest results initially, prioritizing task context and state guidance. The full 200-record set, filters, tests, correlation data, and failures remain in diagnostics.
- **Three native disclosures**: Destinations, Assignment rules, and Delivery diagnostics use the existing native Disclosure component and start collapsed. This is simpler and more consistent than nested routing for three related tasks.
- **Plain-language guidance map**: One deterministic mapping owns labels, tones, meanings, and next actions for every delivery and overview state. State remains understandable without color.
- **Preserved specialist capability**: Existing channel and policy forms move intact beneath disclosure headings, with concise definitions added through section prose and Field help. Secret values remain absent and write-only.
- **Release boundary**: S085 changes source, tests, specifications, and changelog only. It does not rebuild, retag, publish, or promote v1.4.0.

## Phase Plan

1. Complete specification, autonomous clarifications, UX checklist, research, view-state model, interaction contract, quickstart, tasks, and cross-artifact analysis.
2. Add failing Go tests for secret-free configured coverage and graceful partial projection, then implement the desktop projection.
3. Add failing React tests for the initial overview, state guidance, bounded recent results, collapsed advanced sections, contextual help, and preserved operations.
4. Recompose Notifications into small overview, recent-result, destination, assignment, and diagnostic sections, then add responsive styles.
5. Update browser workflows for first-time comprehension, disclosures, full administration, keyboard use, 800 by 600, 200 percent zoom, states, and accessibility.
6. Run focused checks, native build, canonical verification, record evidence, and publish one pull request closing #232 while leaving release state unchanged.
