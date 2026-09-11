# Verification: Task-First Notifications

## Scope

S085 completes issue #232 by making Notifications task-first while preserving existing destination, assignment, retry, diagnostic, and secret-handling behavior. The slice does not create a release, tag, installer, or public artifact.

## Analysis

The final specification contains 17 functional requirements and 9 measurable success criteria. Every requirement maps to implementation tasks and verification evidence, with no unresolved clarification, contradictory requirement, or untracked implementation scope.

## Test-First Evidence

- Configured coverage tests initially failed to compile because the workspace did not expose coverage records or completeness. The implementation added secret-free task and group projections through existing notification backend methods, deterministic ordering, enabled-destination counts, and partial-failure preservation.
- The task-first React regressions initially failed seven tests against the former destination-first page. The implementation added the status overview, configured coverage, task-prioritized recent results, state guidance, and initially collapsed advanced disclosures while preserving all existing operations.
- The first focused browser run rejected an ambiguous count locator. The assertion was made exact, after which the complete Notifications journey passed at the supported viewport and zoom boundary.
- The first full frontend suite found one obsolete app-level assertion for the removed advanced-first empty headings. The assertion now verifies the task-first setup and recent-result guidance.

## Focused Results

- `cd desktop && go test -race ./notifications`: passed.
- `cd desktop/frontend && npm test -- --run src/notifications/NotificationsPage.test.tsx src/notifications/store.test.ts`: 2 files and 16 tests passed after first-round review fixes.
- `cd desktop/frontend && npm test -- --run`: 24 files and 127 tests passed.
- `cd desktop/frontend && npm run build`: TypeScript and Vite production build passed.
- `cd desktop/frontend && npx playwright test e2e/notifications.spec.ts`: 2 Notifications browser journeys passed, including six overview states with axe checks.
- `cd desktop/frontend && npx playwright test`: all 27 browser journeys passed.

## Native Windows Evidence

The canonical GUI gate compiled the production Wails desktop for `windows/amd64` with Wails v2.15.0. The build completed successfully at `desktop/build/bin/gosched-gui.exe` after generating bindings, installing frontend dependencies, compiling the frontend, generating assets, and compiling the native application.

## Canonical Results

`C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` passed on `codex/085-task-first-notifications` on 2026-09-11, including the first-round review corrections. The eight gates covered formatting, vet, lint, Go race tests, desktop Go race tests, native Wails build, 127 frontend tests, frontend production build, coverage thresholds, documentation and architecture policies, and automation plus negative-fixture contracts.

## Requirement Evidence

- The initial page reports deterministic setup, paused, active, in-progress, retrying, and attention states with text labels and useful next actions.
- Configured task and group coverage is derived from authoritative direct and effective assignments without returning endpoint or authorization secrets.
- The recent list is bounded to five records and prioritizes task outcomes before transport tests, while complete redacted history remains available on demand.
- Destination setup, assignment rules, and delivery diagnostics use native keyboard-accessible disclosures that begin collapsed.
- Existing create, edit, test, enable, disable, remove, policy inheritance, direct override, filters, selected detail, stale snapshot, and write-only secret behaviors remain covered.
- The browser journey covers 800 by 600, 100 through 200 percent zoom, keyboard disclosure use, long values, and axe validation without horizontal page overflow.

## First-Round Review

- Configured coverage now uses at most eight concurrent backend reads instead of one serial round trip per task or group. A race-enabled regression proves multiple reads overlap without exceeding the bound.
- Successful assignment-rule saves refresh the workspace before returning, so configured coverage and counts immediately reflect the accepted policy mutation.
- An enabled but unassigned destination, assignments that only target disabled destinations, and an incomplete coverage projection each receive honest non-healthy status and a specific next action.

## Release Boundary

S085 changes source, tests, and planning records only. No version, release note, tag, GitHub Release, installer publication, or immutable release claim is part of this slice.
