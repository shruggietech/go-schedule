# Implementation Plan: Dependency Consolidation

**Branch**: `codex/073-dependency-consolidation` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/073-dependency-consolidation/spec.md`

## Summary

Apply the eight dependency updates represented by pull requests #201 through #208 directly to the current `main` baseline, regenerate both Go module graphs and the frontend lockfile, and add only the runtime and peer updates required for a valid combined graph. Align the frontend runtime with the requested Node 26 type definitions and jsdom 30 engine range, align Vite 8.2.2 with the requested React plugin 6.1.1 peer range, protect the Node 26 hosted baseline in repository automation tests, and validate the result through clean restoration, focused compatibility checks, canonical local verification, and final-head hosted CI.

## Technical Context

**Language/Version**: Go 1.25.0; Node.js 26; TypeScript 5.6.3

**Primary Dependencies**: Wails 2.15.0, modernc.org/sqlite 1.58.0, fsnotify 1.10.1, React 19.2.8, Testing Library jest-dom 7.0.1, jsdom 30.0.1, Vite React plugin 6.1.1, Node type definitions 26.5.0, and Vite 8.2.2 as a required peer

**Storage**: Existing SQLite scheduler database with no schema change

**Testing**: Go module verification and race suites; npm clean install, unit tests, type checking, production bundle, Chromium accessibility and responsive tests; native Wails build; Windows MSI contract; canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux products and hosted runners; Node.js 26 build environment

**Project Type**: Cross-platform Go daemon, CLI, Wails desktop application, and React frontend

**Performance Goals**: Preserve the established scheduler dispatch p99 below 100 ms and existing core-package coverage thresholds

**Constraints**: Apply one coherent baseline without force or legacy dependency resolution; keep root Go version and toolchain lines unchanged; include no product feature or release operation; record pinned workflow changes in the changelog; keep Windows child processes hidden

**Scale/Scope**: Eight source pull requests, three dependency graphs, ten direct requested packages, two required companion baseline updates, three hosted operating systems, and the existing full product verification surface

## Constitution Check

*GATE: Passed before research and re-checked after design.*

- **I. Code Quality**: The implementation changes dependency declarations, integrity files, and narrowly required runtime contracts only. It does not replicate the invalid peer or runtime combinations present in isolated Dependabot branches.
- **II. Testing Standards**: Existing safety-critical tests remain unchanged and the full race, GUI, coverage, packaging, and automation gates remain mandatory. Clean restoration is added as explicit dependency evidence.
- **III. User Experience Consistency**: No user-facing capability or interaction is redesigned. Existing desktop, accessibility, responsive, and packaging contracts protect behavior.
- **IV. Performance Requirements**: No scheduling hot path changes. Existing benchmark and p99 enforcement remain part of canonical verification.
- **V. Autonomous Build-Phase Execution**: S073 traces to open issue #215, follows the complete Spec Kit sequence, uses a review branch and pull request, and has explicit push and PR authorization.
- **Engineering Constraints**: SQLite migration and persistence tests, fsnotify watcher tests, Wails native builds, frontend tests, and local access-control tests continue across supported platforms.
- **Pinned artifacts**: `.github/workflows/ci.yml` and `.github/workflows/release.yml` must move their Node baseline from 24 to 26. The required dated decision will be recorded in `CHANGELOG.md`.
- **Post-design re-check**: Passed. The design introduces no new application data, public interface, product behavior, or constitutional exception.

## Project Structure

### Documentation (this feature)

```text
specs/073-dependency-consolidation/
├── checklists/
│   ├── dependency-integrity.md
│   └── requirements.md
├── contracts/
│   └── dependency-baseline.md
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
go.mod
go.sum
desktop/
├── go.mod
├── go.sum
└── frontend/
    ├── package.json
    └── package-lock.json
.github/workflows/
├── ci.yml
└── release.yml
scripts/
├── automation-check.sh
└── verify.sh
test/scripts/
└── automation-check_test.sh
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Preserve the existing module boundaries and regenerate each graph from its native package manager. Apply requested versions to current `main` instead of cherry-picking stale Dependabot commits, because those branches predate shipped S067 through S072 work and their historical tree differences are outside dependency scope. Runtime and peer compatibility are enforced at existing manifest and workflow boundaries rather than through new abstractions.

## Complexity Tracking

No constitutional violations or complexity exceptions are required.
