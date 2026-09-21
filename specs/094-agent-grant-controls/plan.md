# Implementation Plan: Agent Grant Controls

**Branch**: `codex/094-agent-grant-controls` | **Date**: 2026-09-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/094-agent-grant-controls/spec.md`

## Summary

Extend the existing actor, pairing, credential, and audit contracts into one secret-free Agent Access workspace. Persist the intended grant deadline with a pairing, project MCP actors into transport-aware grants, copy enrollment material only through the native boundary, expose bounded narrowing, expiry, and revocation operations, and show the newest shared audit evidence without coupling permission grants to listener controls.

## Technical Context

**Language/Version**: Go 1.26.0 and TypeScript 5.9 with React 19

**Primary Dependencies**: Existing SQLite store, versioned Go API client and server, Wails 2 native bridge, React Testing Library, Vitest, and Playwright

**Storage**: Existing SQLite actor, pairing, client credential, and audit tables with one forward-only pairing grant-expiry column

**Testing**: Go unit, integration, race, API lifecycle, TypeScript unit, accessibility, Playwright desktop, MCP conformance, security redaction, and canonical repository gates

**Target Platform**: Supported Windows, macOS, and Linux daemon and desktop packages

**Project Type**: Local daemon, CLI, and Wails desktop application

**Performance Goals**: One bounded concurrent workspace load, no unbounded audit query, at most 25 recent events per selected actor, and no polling faster than the existing five-second active-listener refresh

**Constraints**: Enrollment authority remains local, network transports remain default off, no protected secret may cross into React state, existing connections must observe actor changes on their next request, and grant controls cannot widen or extend access

**Scale/Scope**: One daemon at a time, MCP actors only, three authority levels, three transport labels, five duration choices, and a bounded recent-action list

## Constitution Check

- The specification, clarification record, requirements checklists, research, data model, contracts, quickstart, tasks, analysis, and implementation phases remain in the required order.
- Tests are written or updated before the corresponding behavior and cover domain, storage, service, bridge, interface, accessibility, and end-to-end boundaries.
- Existing actor authorization and audit paths remain authoritative. The desktop adds projection and bounded administration rather than a parallel policy engine.
- Pairing phrases and durable credentials remain outside React and are not written to repository artifacts, logs, or audit records.
- Optional listeners remain independently controlled and default off. Creating a remote grant does not start a listener.
- The final gate includes full canonical verification and the S094 issue and epic acceptance matrices.

All gates pass before research. The post-design check also passes because one nullable deadline field and one desktop projection service are the smallest coherent extension of existing boundaries.

## Project Structure

### Documentation

```text
specs/094-agent-grant-controls/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── checklists/
└── contracts/
```

### Source Code

```text
internal/
├── domain/              # pairing grant deadline and actor lifecycle
├── store/               # forward migration and pairing persistence
├── enrollment/          # expiry propagation into exchanged actors
└── api/                 # existing typed actor, pairing, credential, and audit contracts

desktop/
├── agentaccess/         # secret-free grant projection and lifecycle service
├── app.go               # Wails facade methods
└── frontend/src/agentaccess/ # models, store, dialogs, grant list, and audit view

docs/                    # administrator workflow and security boundary updates
```

**Structure Decision**: Extend the established enrollment and Agent Access packages. No new policy subsystem, credential store, transport, or frontend route is introduced.

## Complexity Tracking

No constitution deviations or complexity exceptions are required.
