# Verification: Remote Connection Profiles and Target-Safe Clients

## Specification Analysis

`/speckit-analyze` completed against `spec.md`, `plan.md`, `research.md`, `data-model.md`, `contracts/profiles.md`, `tasks.md`, and the constitution. All 22 functional requirements, nine success criteria, four independently testable user stories, issue #170 acceptance criteria, and issue #171 acceptance criteria map to implementation or verification evidence. No unresolved ambiguity, duplication, constitution conflict, or coverage gap remains. Automatic reconnect, event resume, certificate-change recovery, and stale offline state remain isolated to #172.

## Focused Evidence

- `go test -race ./internal/api/client ./internal/clientprofile ./internal/clientsecret ./internal/cli ./internal/remoteenroll` passed.
- `go test -race ./connection ./connections ./remotepairing` passed from `desktop/`.
- `npm test -- --run` passed all 85 frontend tests.
- `npm run build` passed TypeScript validation and the production Vite build.
- `npm run test:e2e` passed all 21 Chromium accessibility, keyboard, scale, responsive, and connection-recovery tests.
- Target-isolation tests prove in-flight requests stay on their original immutable client, stale connection generations are rejected, target selection precedes request-router changes, remote failure requires manual retry, and active removal switches to local before native credential deletion.
- Secret-exclusion tests prove profile JSON and frontend workspace projections contain no bearer values or pairing phrases; CLI list and show omit certificate bodies.

## Canonical Verification

`scripts/verify.sh all` passed the canonical format, vet, lint, race, GUI, coverage, docs, and automation gates on 2026-09-09. Core coverage remained above the required threshold: engine 82.9 percent, schedule 89.1 percent, timezone 91.3 percent, store 80.2 percent, catchup 88.9 percent, and logbus 91.1 percent. The GUI gate included desktop race tests, generated Wails bindings, a native Windows production build, frontend tests, and the production frontend build.

## Review Evidence

Pending pull-request review.
