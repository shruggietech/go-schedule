# Implementation Plan: Wails Foundation and Experience Direction

**Branch**: `codex/060-wails-foundation-direction` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/060-wails-foundation-direction/spec.md`

## Summary

Resolve issues #149 and #150 together by committing an isolated Wails v2.14.0 proof built with React, TypeScript, Vite, and locally hosted brand assets. The proof uses typed Go boundaries for the current local daemon and native dialogs, provides a fresh target-aware control-center prototype, and remains outside shipped binaries. Go, frontend, browser, accessibility, offline-asset, and cross-platform build checks make the architectural and experience decisions reproducible before production shell work begins.

## Technical Context

**Language/Version**: Go 1.25.0; TypeScript 5.6.3; Node.js 24 LTS

**Primary Dependencies**: Wails v2.14.0 stable; React 19.1.0; React DOM 19.1.0; Vite 7.3.6; Vitest 5.0.0; Testing Library; axe-core 4.13.0; Playwright 1.63.0

**Storage**: No new persistence; representative fixture data and the existing local daemon IPC client

**Testing**: Go unit and race tests; TypeScript type-check and Vitest component tests; axe accessibility checks; Playwright keyboard, responsive, offline-reference, and screenshot checks; Wails cross-platform build matrix

**Target Platform**: Windows 10/11 amd64, macOS 10.13+ release and current development runners, Linux amd64 with GTK3 and WebKitGTK 4.1

**Project Type**: Existing Go daemon/CLI/Fyne desktop plus an isolated non-shipping Wails desktop proof

**Performance Goals**: Prototype first content visible within 1 second in the browser harness; input feedback within one 180 ms normal-motion token; no user workflow depends on animation

**Constraints**: Complete offline application assets; no Electron or bundled browser; no production package or runtime changes; no remote daemon implementation; 200 percent zoom and 900 by 650 compact viewport; Windows console remains hidden for any child process

**Scale/Scope**: Two GitHub issues, six representative views, nine state classes, four typed proof boundaries, three hosted operating systems, one selected direction

## Constitution Check

### Pre-research gate

| Principle | Requirement | Status |
| --- | --- | --- |
| I. Code Quality | Typed boundaries, contextual errors, bounded goroutine lifecycle, documented exported contracts | PASS. The proof keeps daemon streaming under a cancelable application context and uses explicit adapters. |
| II. Testing Standards | Tests alongside behavior, race coverage, no ignored flakes | PASS. Go, component, browser, accessibility, and hosted platform checks are planned before acceptance. |
| III. UX Consistency | Predictable language, states, target identity, actionable errors | PASS. The experience contract makes these explicit and measurable. |
| IV. Performance | No unjustified hot-path change and measured interactive budget | PASS. Production scheduling paths are untouched; browser checks enforce the prototype response budget. |
| V. Autonomous Execution | Full Spec-Kit order, analyze gate, review branch, one publication authorization boundary | PASS. The operator explicitly authorized S060 publication and up to two Codex rounds. |

**Engineering constraints**: The stable release and every dependency are justified in [research.md](research.md). The proof does not alter stored data or the production application. Hosted evidence adds macOS to the issue-specific proof without changing the constitution's Linux and Windows minimum.

## Project Structure

### Documentation (this feature)

```text
specs/060-wails-foundation-direction/
├── checklists/
│   ├── foundation-ux.md
│   └── requirements.md
├── contracts/
│   ├── experience-contract.md
│   └── foundation-contract.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
experiments/wails-foundation/
├── build/
│   └── appicon.png
├── frontend/
│   ├── e2e/
│   ├── public/fonts/
│   ├── src/
│   ├── index.html
│   ├── package-lock.json
│   ├── package.json
│   ├── playwright.config.ts
│   ├── tsconfig.json
│   └── vite.config.ts
├── app.go
├── app_test.go
├── go.mod
├── go.sum
├── main.go
├── README.md
└── wails.json

.github/workflows/ci.yml
scripts/automation-check.sh
test/scripts/automation-check_test.sh
brand/repository-consumers.json
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: A nested experimental module keeps Wails and frontend dependencies out of the production Go module and existing release graph while allowing exact, reproducible builds. The proof imports the repository's existing protected local API through a local module replacement, but its Wails adapter, fixtures, and generated frontend bindings stay isolated. Production code graduates selectively in #151 and #152 after this decision is accepted.

## Implementation Sequence

1. Establish the nested proof manifests, asset mapping, typed boundary models, and test doubles.
2. Write failing Go tests for health/list/event/native behavior and shutdown, then implement the Wails-facing application adapter.
3. Write failing component and browser acceptance tests for hierarchy, states, keyboard flow, contrast, reduced motion, zoom, responsive reflow, and offline assets.
4. Implement one target-aware control-center prototype and its local fixture fallback.
5. Add hosted Windows, macOS, and Linux proof jobs plus the minimum action-policy update and fixtures.
6. Run analysis remediation, dependency/license audit, formatting, focused proof checks, and full canonical verification.

## Post-design Constitution Check

| Principle | Design result | Status |
| --- | --- | --- |
| I. Code Quality | The disposable boundary is small, typed, cancelable, and isolated | PASS |
| II. Testing Standards | Every proof behavior and experience criterion maps to an automated task before implementation | PASS |
| III. UX Consistency | One tokenized direction covers hierarchy, states, language, target identity, and accessibility | PASS |
| IV. Performance | No scheduler path changes; prototype budgets are browser-observable and motion is bounded | PASS |
| V. Autonomous Execution | Spec, clarification, checklist, plan, tasks, analysis, implementation, verification, commit, publication, and review remain ordered | PASS |

## Process Deviation

The installed checklist prerequisite requires `plan.md`, while the governing Spec-Kit command order and project autopilot require checklist before plan. S060 generated and validated `checklists/foundation-ux.md` directly from the resolved feature path before planning. This is the established repository workaround, preserves the required semantic order, and does not weaken checklist review.

## Complexity Tracking

No constitution violation requires justification. The nested module is deliberate isolation for a disposable proof, not a new production subsystem; placing Wails in the root module would prematurely couple the shipping Fyne application and every root verification command to a decision artifact.
