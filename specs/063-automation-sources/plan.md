# Implementation Plan: Connected Automation Sources

**Branch**: `codex/063-automation-sources` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/063-automation-sources/spec.md`

## Summary

Replace the remaining legacy chain, trigger, Trigger Set, and filesystem-watcher administration surfaces with one Wails Automation Sources workspace. A transport-neutral Go service will compose existing daemon APIs into secret-free presentation models, enforce stale-write and transient-secret boundaries, and expose focused lifecycle actions. React will provide searchable type sections, type-specific editors, accessible confirmations, and ephemeral secret dialogs while preserving last-known context during connection loss.

## Technical Context

**Language/Version**: Go 1.25, TypeScript 5.9, React 19

**Primary Dependencies**: Existing local daemon client, Wails v2, React, Vite, Vitest, Testing Library

**Storage**: Existing SQLite-backed daemon storage; no schema change

**Testing**: Go unit and race tests, Vitest and axe accessibility tests, canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop application

**Project Type**: Go desktop application with an embedded React frontend

**Performance Goals**: Load, search, filter, and operate on at least 100 mixed automation sources with ordinary interactions completing within two seconds

**Constraints**: Raw trigger keys only cross the Wails boundary for explicit create, reveal, or rotate actions; watcher health remains daemon-owned; no new persistence or API semantics; offline state is read-only

**Scale/Scope**: One route, four source types, roughly 20 lifecycle actions, one backend service package, and focused Go and React coverage

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: The desktop service owns composition and validation, React owns presentation, and daemon rules are not duplicated. Public transport models use explicit JSON fields and bounded error vocabulary.
- **Testing**: Backend contract, stale-write, secret-boundary, lifecycle, request-ordering, duplicate-submission, large-collection, and accessibility tests are planned. Canonical format, vet, lint, race, GUI, coverage, documentation, and automation gates remain mandatory.
- **UX consistency**: Existing shell, buttons, dialogs, status badges, notices, state panels, and visual tokens are reused. All four source types use a common source-to-target vocabulary and keyboard behavior.
- **Performance**: Workspace loading uses bounded daemon reads, presents one complete snapshot, and filters locally. No polling or per-row backend calls are introduced.
- **Security**: Ordinary models exclude keys, fire resolves and consumes its key within Go, ephemeral secret results are explicit, and errors never reflect secret material.
- **Review workflow**: Work remains on `codex/063-automation-sources`, canonical verification precedes commit, and the user explicitly authorized push, PR publication, review fixes, and at most one second Codex review round.
- **Pinned artifacts**: No pinned build or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/063-automation-sources/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── desktop-automation-bridge.md
├── checklists/
│   ├── automation-safety.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
desktop/
├── app.go
├── main.go
├── automation/
│   ├── local.go
│   ├── model.go
│   ├── service.go
│   └── service_test.go
└── frontend/src/
    ├── App.tsx
    ├── App.test.tsx
    ├── connection/model.ts
    ├── components/Shell.tsx
    ├── styles.css
    └── automation/
        ├── AutomationPage.tsx
        ├── AutomationPage.test.tsx
        ├── bridge.ts
        ├── model.ts
        ├── store.ts
        └── store.test.ts
```

**Structure Decision**: Add one desktop feature package parallel to `desktop/taskgroup` and one React feature folder parallel to `desktop/frontend/src/tasks`. Keep daemon persistence, validation, and readiness code unchanged and adapt its existing typed clients through a narrow backend interface.

## Complexity Tracking

No constitutional violations require justification.
