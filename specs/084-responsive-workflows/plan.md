# Implementation Plan: Responsive Task and Administration Workflows

**Branch**: `codex/084-responsive-workflows` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/084-responsive-workflows/spec.md`

## Summary

Move task creation and editing into the existing accessible dialog system, add small reusable administration composition primitives, and apply them across Agent Access, Connections, Pairing, and Settings. Preserve every existing bridge and data contract while closing #231 and #233 with record-scoped copy feedback, deterministic component tests, responsive browser coverage, a native Windows production build, and canonical verification.

## Technical Context

**Language/Version**: TypeScript 5.6.3, React 19.2.8, CSS, Go 1.25.0 with toolchain 1.25.1 for the Wails host

**Primary Dependencies**: Existing React component catalog, Wails desktop host, Vite 8.2.2, Vitest 5.0.0, Testing Library, Playwright 1.63.0, axe-core 4.13.0

**Storage**: Existing desktop appearance preference and daemon-owned data only; no schema or persistence change

**Testing**: Vitest component and store tests, Playwright Chromium workflow and accessibility checks, native Windows production build, `sh scripts/verify.sh all`

**Target Platform**: Wails desktop on Windows with cross-platform source compatibility; 800 by 600 minimum viewport and 200 percent zoom reflow

**Project Type**: Desktop application with a React frontend and Go/Wails host

**Performance Goals**: No new network or daemon work; UI state remains bounded by one task draft, existing workspace data, and record-keyed in-flight settings actions

**Constraints**: Local assets only, WCAG 2.2 AA, keyboard parity, no secret display, no external UI dependency, no daemon contract or release-state changes

**Scale/Scope**: One task dialog, three administration destinations, one pairing form, five reusable presentation primitives, and record-scoped settings actions

## Constitution Check

- **I. Code Quality**: Replace repeated ad hoc form and definition layouts with narrowly scoped shared presentation primitives and preserve typed bridge boundaries.
- **II. Testing Standards**: Add failing regression tests before implementation, cover failures and overlapping copy operations, retain accessibility scans, and run all eight canonical gates.
- **III. User Experience Consistency**: Use one dialog, form-grid, description-list, status-label, path, disclosure, and action-group contract across the affected workflows.
- **IV. Performance Requirements**: Keep transient state local and bounded, avoid new polling or backend work, and use CSS reflow rather than JavaScript layout measurement.
- **V. Autonomous Build-Phase Execution**: Complete every spec-kit phase, review branch, pull request, hosted check, permitted bot review round, and maintainer merge gate.

**Gate result before design**: PASS. No constitution violation or unresolved clarification exists.

**Gate result after design**: PASS. The design introduces no dependency, schema, background process, network route, or secret-bearing state.

## Project Structure

### Documentation

```text
specs/084-responsive-workflows/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── contracts/
│   └── responsive-workflows.md
└── checklists/
    ├── requirements.md
    └── ux.md
```

### Source Code

```text
desktop/frontend/
├── src/
│   ├── App.tsx
│   ├── App.test.tsx
│   ├── styles.css
│   ├── components/
│   │   ├── index.tsx
│   │   └── components.test.tsx
│   ├── tasks/
│   │   ├── TasksPage.tsx
│   │   ├── TasksPage.test.tsx
│   │   ├── TaskEditor.tsx
│   │   ├── TaskEditor.test.tsx
│   │   ├── TaskActions.tsx
│   │   └── TaskActions.test.tsx
│   ├── agentaccess/
│   │   ├── AgentAccessPage.tsx
│   │   └── AgentAccessPage.test.tsx
│   ├── settings/
│   │   ├── ConnectionsPage.tsx
│   │   ├── ConnectionsPage.test.tsx
│   │   ├── SettingsPage.tsx
│   │   ├── SettingsPage.test.tsx
│   │   ├── store.ts
│   │   └── store.test.ts
│   └── remotepairing/
│       ├── PairingForm.tsx
│       └── PairingForm.test.tsx
└── e2e/
    └── shell.spec.ts

CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Extend the existing component catalog and domain pages. The work is presentation and transient interaction state over stable bridges, so another component framework, global state layer, or backend abstraction would add migration risk without satisfying an unmet requirement.

## Decision Log

- **One create and edit dialog**: TaskEditor will compose the shared Dialog directly and receive the exact invoker from TasksPage. A scrollable dialog body and persistent footer keep the workflow bounded while one component preserves create and edit parity.
- **Small semantic presentation primitives**: Add FormGrid, DescriptionList, StatusLabel, CardSection, and PathDisplay as markup and class contracts. They own composition semantics but not domain data or mutations.
- **Native disclosure**: Continue using the existing Details-based Disclosure component for advanced task fields, inactive localhost setup, and remote pairing. This preserves keyboard semantics and avoids custom tab or accordion state.
- **Basic information first**: Agent Access shows safe process and authority status before optional network configuration. Connections shows diagnosis and saved targets before a collapsed pairing surface. A repair request expands pairing because the user has already chosen that task.
- **Record-keyed settings operations**: Replace one page-wide pending flag with a set of action keys. Copy progress and completion belong to the selected storage record, concurrent different actions remain independent, and late results cannot relabel another record.
- **Complete operational values**: Use selectable monospace PathDisplay for filesystem paths and safe wrapping for endpoints, fingerprints, origins, and identifiers. Do not truncate the only visible value.
- **Release boundary**: S084 changes source, tests, specifications, and changelog only. It does not move, rebuild, retag, publish, or promote the existing v1.4.0 draft.

## Phase Plan

1. Complete specification, clarification record, UX requirements checklist, research, view-state model, interaction contract, quickstart, executable tasks, and analysis.
2. Add failing component and store regressions for modal task authoring, focus restoration, initial disclosure, responsive administration semantics, complete path presentation, and record-scoped copy state.
3. Extend shared presentation primitives and dialog composition, then move create and edit into the bounded task modal and retain task confirmation behavior.
4. Recompose Agent Access, Connections, Pairing, and Settings with basic information first, responsive grids, safe long-value wrapping, and independent copy feedback.
5. Add browser-level workflow, reflow, zoom, theme, keyboard, and accessibility coverage, then run the native Windows production build and all canonical gates.
6. Record verification and publish one pull request closing #231 and #233 while leaving #232 and release state unchanged.
