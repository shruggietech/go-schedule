# Implementation Plan: Desktop Notification Management

**Branch**: `codex/068-desktop-notifications` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/068-desktop-notifications/spec.md`

## Summary

Add one transport-neutral `desktop/notifications` service over the existing notification, task, and group client methods, then expose it through the Wails facade to a React Notifications workspace. The backend will return only redacted summaries, bounded history, safe outcomes, and authoritative direct/effective policy. The frontend will provide explicit write-only replacement controls, accessible channel and policy workflows, stable selection, request sequencing, debounced event refresh, and delivery-state explanation.

## Technical Context

**Language/Version**: Go 1.25, TypeScript 5.9, React 19

**Primary Dependencies**: Go standard library, existing local daemon client and notification API, Wails v2 runtime, React, Vite, Vitest, Testing Library, Playwright

**Storage**: Existing notification channel, assignment, and delivery tables only; no migration or new persistence

**Testing**: Go unit and race tests, Vitest and axe accessibility tests, Playwright Chromium interaction and zoom tests, native Wails build, canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop application

**Project Type**: Go desktop application with an embedded React frontend

**Performance Goals**: Load one complete workspace of current channels, tasks, groups, and 200 deliveries within three seconds on the protected local transport; filter 200 records without perceptible input delay

**Constraints**: No protected value may cross a desktop response; no broad per-task policy fan-out; all daemon calls bounded; authoritative inheritance only; stale response suppression; last-complete-snapshot preservation; no new dependency; 80 through 200 percent zoom

**Scale/Scope**: One new route, one backend feature package, seven facade methods, one React feature folder, up to 200 history records, and focused Go, frontend, accessibility, browser, native, and canonical verification

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: Notification composition and safe mapping live in one desktop service. React owns presentation and transient drafts. Existing server precedence remains authoritative.
- **Testing**: Failure-first tests cover redaction, validation, mutation refresh, inheritance, state mapping, stale responses, duplicate activation, accessibility, zoom, and native build.
- **UX consistency**: Existing shell, panels, notices, tables, dialogs, buttons, focus behavior, and live announcements are reused.
- **Performance**: Initial workspace uses four bounded aggregate reads. Policy detail is fetched only for the selected scope. History is capped at 200 and filtered locally.
- **Security and truth**: Desktop models contain no endpoint, authorization, or payload field. Existing protected transport and persistence are reused. Server responses decide precedence and mutation results.
- **Review workflow**: Work remains on `codex/068-desktop-notifications`; the user explicitly authorized push, PR publication, review fixes, and at most one manually triggered second Codex review round.
- **Pinned artifacts**: No pinned workflow, toolchain, release, or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/068-desktop-notifications/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── desktop-notifications-bridge.md
├── checklists/
│   ├── notification-trust.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
desktop/
├── app.go
├── app_test.go
├── main.go
├── notifications/
│   ├── local.go
│   ├── model.go
│   ├── service.go
│   └── service_test.go
└── frontend/
    ├── e2e/
    │   └── notifications.spec.ts
    └── src/
        ├── App.tsx
        ├── App.test.tsx
        ├── components/Shell.tsx
        ├── connection/model.ts
        ├── styles.css
        └── notifications/
            ├── bridge.ts
            ├── model.ts
            ├── store.ts
            ├── store.test.ts
            ├── NotificationsPage.tsx
            └── NotificationsPage.test.tsx
```

**Structure Decision**: Add one desktop feature package parallel to `desktop/operations` and one React feature folder parallel to its schedule and activity workspaces. The package adapts existing daemon client calls through a narrow interface, produces secret-free complete snapshots, and fetches policy details only on selection. This avoids coupling Notifications to the task editor or duplicating server policy logic.

## Complexity Tracking

No constitutional violations require justification.
