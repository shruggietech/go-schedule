# Implementation Plan: Desktop Visual and Shell Foundations

**Branch**: `codex/083-desktop-foundations` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/083-desktop-foundations/spec.md`

## Summary

Replace accumulated page-level presentation defaults with a compact shared control and spacing system, make the application shell own the viewport, and classify status feedback into transient toasts or dismissible persistent notices. Preserve all daemon, storage, routing, and page-domain behavior while closing #229 and #230 with component, browser, accessibility, canonical, and focused native Windows evidence.

## Technical Context

**Language/Version**: TypeScript 5.6.3, React 19.2.8, CSS, Go 1.25.0 with toolchain 1.25.1 for the Wails host

**Primary Dependencies**: Existing React component catalog, Wails desktop host, Vite 8.2.2, Vitest 5.0.0, Testing Library, Playwright 1.63.0, axe-core 4.13.0

**Storage**: Existing desktop appearance preference only; no schema or persisted-data change

**Testing**: Vitest component tests, Playwright Chromium shell and accessibility tests, focused native Windows build and observation, `sh scripts/verify.sh all`

**Target Platform**: Wails desktop on Windows with cross-platform source compatibility; 800 by 600 minimum viewport and 200 percent zoom reflow

**Project Type**: Desktop application with a React frontend and Go/Wails host

**Performance Goals**: No new network or persistence work; feedback timers and palette listeners remain constant-space and event-driven

**Constraints**: Local assets only, WCAG 2.2 AA contrast and reflow, keyboard parity, reduced-motion behavior, no daemon or public interface changes, no release mutation

**Scale/Scope**: Eight routes, five shared action variants, three appearance preferences, two feedback lifetimes, one active-page scroll region

## Constitution Check

- **I. Code Quality**: Centralize reusable tokens and feedback behavior, remove message concatenation, and avoid page-specific exceptions.
- **II. Testing Standards**: Add failing component and browser tests before implementation, retain axe coverage, and run all eight canonical gates.
- **III. User Experience Consistency**: Apply one action hierarchy, spacing scale, theme resolution contract, and feedback classification across the desktop.
- **IV. Performance Requirements**: Use one media-query listener and one bounded timer with cleanup; do not add polling or work to daemon hot paths.
- **V. Autonomous Build-Phase Execution**: Complete the required spec-kit phases, review branch, pull request, hosted checks, bot review, and maintainer merge gate.

**Gate result before design**: PASS. No constitution violation or unjustified complexity is required.

**Gate result after design**: PASS. The design adds no dependency, storage entity, external interface, or background service.

## Project Structure

### Documentation

```text
specs/083-desktop-foundations/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── contracts/
│   └── desktop-foundations.md
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
│   ├── accessibility.test.tsx
│   ├── styles.css
│   ├── components/
│   │   ├── index.tsx
│   │   ├── components.test.tsx
│   │   └── Shell.tsx
│   └── tasks/
│       ├── TaskActions.tsx
│       └── TaskActions.test.tsx
└── e2e/
    └── shell.spec.ts

CHANGELOG.md
docs/decisions.md
CLAUDE.md
```

**Structure Decision**: Extend the existing component catalog and shell rather than adopt a third-party design-system dependency. The defects are caused by inconsistent shared contracts and viewport ownership, so repairing those seams is smaller, easier to audit, and immediately consumable by #231 through #233.

## Decision Log

- **Tokenized CSS foundation**: Define semantic variables at the application root and consume them in shared controls, panels, dialogs, and shell regions. This keeps palette resolution and density auditable without migrating every page to a new library.
- **Viewport-owned shell**: Set the application to the available viewport height, keep the rail and workspace chrome within that grid, and make `main` the sole vertical scroll container. At narrow widths, retain a compact fixed shell with horizontally scrollable navigation instead of hiding Exit.
- **Feedback classification**: Add a stateful toast for routine messages and a dismissible Notice contract for persistent problems. New transient messages replace the prior one, timers pause on hover or focus, and empty messages remove the visual toast while preserving the live-region contract.
- **Theme-aware identity**: Resolve `system` through `prefers-color-scheme`, expose the resolved palette to CSS and accessible tests, and use the existing neutral mark with palette-aware treatment rather than introducing unreviewed brand artwork.
- **Release boundary**: S083 changes source and tests only. The existing v1.4.0 tag and draft release are immutable inputs that this slice does not modify.

## Phase Plan

1. Complete specification, clarification record, UX requirements checklist, research, interaction contract, and executable tasks.
2. Add failing tests for control variants, dialog action grouping, theme resolution, fixed shell overflow, and transient feedback timing.
3. Implement shared primitives and semantic token states, then wire affirmative and destructive task actions.
4. Implement viewport ownership, narrow reflow, theme-aware identity, dismissible notices, and bounded toast behavior.
5. Run focused frontend tests, browser accessibility and reflow checks, a native Windows build and focused visual observation, then all canonical gates.
6. Record verification and publish one reviewable pull request closing #229 and #230 without mutating release state.
