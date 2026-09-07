# Implementation Plan: Desktop Settings, Information, and Recovery

**Branch**: `codex/065-options-info` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/065-options-info/spec.md`

## Summary

Replace the Settings and Connections placeholders with complete Wails workspaces. A new `desktop/settings` service will own versioned preferences, one-time legacy appearance migration, authoritative storage inventory, product metadata, and bounded native copy and link actions. React will apply the persisted appearance at startup, expose accessible Settings sections, and turn the existing connection manager diagnosis into an inline recovery route without modal repetition.

## Technical Context

**Language/Version**: Go 1.25, TypeScript 5.9, React 19

**Primary Dependencies**: Go standard library, existing local daemon client, Wails v2 runtime, React, Vite, Vitest, Testing Library, Playwright

**Storage**: Versioned per-user JSON desktop preference file plus read-only existing daemon runtime metadata; no database schema change

**Testing**: Go unit and race tests, Vitest and axe accessibility tests, Playwright Chromium interaction and zoom tests, native Wails build, canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop application

**Project Type**: Go desktop application with an embedded React frontend

**Performance Goals**: Load local settings without daemon dependency, render the complete inventory, and apply a preference mutation within two seconds under ordinary local filesystem conditions

**Constraints**: One-time legacy migration; system default; no new dependency; no arbitrary native action input; daemon paths remain authoritative; atomic preference replacement; offline local usability; no repeated modal recovery; supported browser zoom from 80 through 200 percent

**Scale/Scope**: Two routes, one backend service package, one small versioned preference document, a bounded storage inventory, two fixed product links, and focused Go, React, accessibility, zoom, and native-build coverage

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **Code quality**: The settings service owns filesystem, migration, inventory, and action authorization. React owns presentation and current route state. Errors are explicit and recoverable.
- **Testing**: Migration matrices, atomic-write failures, invalid values, storage authority, native-action authorization, offline behavior, retry sequencing, accessibility, and zoom are mandatory.
- **UX consistency**: Existing shell, appearance tokens, notices, buttons, status announcements, and connection manager are reused. Settings consolidates legacy Options and Info outcomes with a simpler hierarchy.
- **Performance**: Local inspection is bounded to known paths, runtime metadata is one aggregate daemon read, and no polling or per-row daemon request is introduced.
- **Security and truth**: Preferences contain no secrets. File permissions are restrictive where supported. Frontend input is a stable identifier, not arbitrary clipboard content or URLs. Daemon and external path ownership is not inferred.
- **Review workflow**: Work remains on `codex/065-options-info`; the user explicitly authorized push, PR publication, review fixes, and at most one second Codex review round.
- **Pinned artifacts**: No pinned build or policy artifact is expected to change.

## Project Structure

### Documentation (this feature)

```text
specs/065-options-info/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── desktop-settings-bridge.md
├── checklists/
│   ├── preferences-diagnostics.md
│   └── requirements.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
desktop/
├── app.go
├── main.go
├── settings/
│   ├── local.go
│   ├── model.go
│   ├── preferences.go
│   ├── service.go
│   └── service_test.go
└── frontend/
    ├── e2e/
    │   └── settings.spec.ts
    └── src/
        ├── App.tsx
        ├── App.test.tsx
        ├── Shell.tsx
        ├── Shell.test.tsx
        ├── styles.css
        └── settings/
            ├── ConnectionsPage.tsx
            ├── ConnectionsPage.test.tsx
            ├── SettingsPage.tsx
            ├── SettingsPage.test.tsx
            ├── bridge.ts
            ├── model.ts
            ├── store.ts
            └── store.test.ts
```

**Structure Decision**: Add one desktop feature package parallel to `desktop/operations` and one React feature folder parallel to existing workspaces. The service adapts the existing runtime-info client through a narrow interface and injects filesystem, clipboard, browser, and platform boundaries for deterministic tests. The shell continues to own global appearance application and the existing connection manager remains the single recovery authority.

## Complexity Tracking

No constitutional violations require justification.
