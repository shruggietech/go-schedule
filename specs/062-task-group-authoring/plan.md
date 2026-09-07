# Implementation Plan: Task and Group Authoring

**Branch**: `codex/062-task-group-authoring` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/062-task-group-authoring/spec.md`

## Summary

Replace the S061 Tasks placeholder with a complete task and group authoring workflow while keeping the daemon authoritative. Add a backward-compatible detailed task-list representation, a transport-neutral desktop task/group service with safe result models, React stores and screens for overview and editing, and one root-owned platform suggestion shared with the still-shipping Fyne editor. Preserve every task field and group lifecycle, add stale-draft protection, and prove safe platform examples on Windows, macOS, and Linux without migrating Activity or changing packaging.

## Technical Context

**Language/Version**: Go 1.25.0, TypeScript 5.9, React 19

**Primary Dependencies**: Existing standard library, internal command-line/schedule/task packages, Wails 2.14.0, Vite 7.3.1, Vitest 4.0.18, Testing Library 16.3.2, Playwright 1.57.0, existing Fyne 2.8.1 shipping surface

**Storage**: Existing SQLite schema and daemon store, no migration

**Testing**: Go unit, HTTP contract, client, service, race, and safe native command tests; Vitest component/store tests; Playwright Chromium accessibility, keyboard, responsive, and high-volume contracts; Wails cross-platform builds

**Target Platform**: Offline desktop on supported Windows, macOS, and Linux packages through protected local IPC

**Project Type**: Go daemon and API plus nested Wails desktop module with React frontend; temporary Fyne compatibility adapter until #157

**Performance Goals**: Local search, filter, selection, hierarchy expansion, and editor disclosure within 200 milliseconds for one hundred tasks and twenty group levels; one bounded task-detail list request instead of per-row requests

**Constraints**: Daemon-authoritative writes; no optimistic entity mutation; no secret or command details in ordinary events; new tasks inactive; direct execution only; stable selection and drafts across live refresh; canonical eight-gate verifier; no new dependency

**Scale/Scope**: Tasks and Groups route, five user stories, thirty functional requirements, one hundred task rows, twenty hierarchy levels, three platform suggestions, current local daemon only

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I, Correctness and Determinism**: PASS. Existing schedule, timezone, command-line, readiness, group-chain, and daemon mutation authorities are reused. UI drafts never become authoritative records.
- **Principle II, Tests Are the Product Contract**: PASS. Contract tests precede server, client, service, Fyne, React, and cross-platform behavior. Race, accessibility, high-volume, stale-response, and native example checks are explicit.
- **Principle III, Simplicity and Scope Discipline**: PASS. No new dependency, storage migration, remote transport, Activity migration, or packaging cutover is introduced. One detailed-list query avoids an N+1 client architecture.
- **Principle IV, Performance and Observability**: PASS. The plan bounds local interactions at 200 milliseconds, retrieves task detail in one IPC request, and retains safe existing event observability without sensitive payloads.
- **Principle V, Autonomous Build-Phase Execution**: PASS. S062 follows specify, clarify, checklist, plan, tasks, analyze, implement, verify, commit, authorized publication, hosted CI, and at most two review rounds.
- **Security constraints**: PASS. Protected local IPC remains unchanged, suggested commands are read-only and bounded, field/backend detail is sanitized, and confirmations avoid secret-bearing values.
- **Pre-design gate**: PASS with no exception or complexity waiver.

## Project Structure

### Documentation (this feature)

```text
specs/062-task-group-authoring/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)
```text
internal/
├── api/
│   ├── client/               # detailed task list and group update client methods
│   └── server/               # backward-compatible detail query and complete clear semantics
└── commandexample/           # canonical platform suggestion and native contract tests

gui/
├── editor.go                 # temporary shipping Fyne suggestion presentation
└── command_suggestion.go     # focused Tab insertion adapter

desktop/
├── app.go                    # stable Wails facade methods
├── taskgroup/
│   ├── model.go              # sanitized frontend contract
│   ├── local.go              # protected local client adapter
│   └── service.go            # authority, validation, stale checks, derived states
└── frontend/
    ├── src/
    │   ├── tasks/            # bridge models, store, overview, editor, groups, tests
    │   ├── components/       # selection, confirmation, field and state primitives
    │   └── App.tsx
    └── e2e/                  # keyboard, accessibility, scale, narrow-window contracts

docs/gui-fields.md            # safe platform examples and guided path
.github/workflows/ci.yml      # explicit native suggestion contract in desktop matrix
```

**Structure Decision**: Keep domain logic and the canonical safe example in the root module, add one bounded `desktop/taskgroup` application service beside S061's connection package, and organize the frontend by the Tasks feature. This preserves one-way dependency from the nested desktop module into the root module, avoids duplicating schedule and readiness logic in TypeScript, and leaves other Wails routes isolated for later slices.

## Planned Architecture

### Read and event flow

1. `GET /v1/tasks?details=true` returns the existing task detail shape as a list while the default response remains unchanged for older clients.
2. The local desktop adapter fetches detailed tasks and groups through the existing protected client.
3. The task/group service derives full paths, declared/effective state, safe summaries, valid parent choices, and stable timestamps into transport-neutral models.
4. The React store loads one workspace snapshot, preserves selection by identity, and refreshes after accepted writes or task/group events.
5. An open dirty editor ignores replacement data, marks itself stale when its source timestamp changes, and requires explicit reload or overwrite intent.

### Write flow

1. The frontend submits a complete typed draft and original update timestamp to the Wails facade.
2. The service parses the direct command line, validates environment rows and schedule mode, fetches current authority for edits, and rejects stale drafts unless overwrite is explicit.
3. The local adapter maps the draft to existing create/update/preview/group calls, including explicit clear intent for supported optional fields.
4. The facade returns a safe field-aware result and refreshed authoritative task, group, or workspace data; rejected writes never alter local entities.
5. Pending controls suppress duplicate activation. Late results remain associated with the operation identity and do not change a different selection.

### Safe example flow

1. `internal/commandexample` owns the Windows, macOS, and Linux command, parsed invocation, explanation, and recognizable-output predicate.
2. Wails obtains the mapping from the active daemon platform in the connection snapshot, not browser runtime identity.
3. The Fyne adapter uses the same root mapping and implements one-shot unmodified Tab insertion only for an empty new-task command field.
4. Cross-platform Go tests execute the current platform example with a five-second deadline and reject mismatched, mutating, interactive, or network-bearing mappings by exact allowlist.

## Post-Design Constitution Check

- All pre-design gates remain PASS after the data model and contracts.
- The detailed-list mode is a backward-compatible read extension, not a new endpoint or persistence model.
- Centralizing the platform mapping is an architecture-affecting compatibility choice and will be recorded in `CHANGELOG.md`.
- Explicit task clear flags repair a pre-existing partial-update weakness only for fields the replacement must faithfully edit; no broader API redesign is planned.
- No complexity-table entry is required.
