# Quickstart: Cross-Daemon Search and Target-Safe Actions

## Prerequisites

- A local daemon and desktop development environment.
- At least two saved remote profiles for the duplicate-name and partial-failure scenarios.
- One Observe-only profile and one Operate-capable profile.
- Node dependencies installed under `desktop/frontend`.

## Focused daemon checks

```powershell
go test -race ./internal/store ./internal/api/server ./internal/api/client ./internal/authorization ./internal/remote
```

Expected outcome: bounded search, validation, authorization, redaction, local client, and generated remote transport tests pass.

## Focused desktop service checks

```powershell
go test -race ./desktop/search ./desktop/connection ./desktop/systems
```

Expected outcome: fan-out limits, generation cancellation, duplicate registrations, identity revalidation, authority differences, wrong-target protection, partial failures, and per-object action outcomes pass.

## Focused frontend checks

```powershell
Set-Location desktop/frontend
npm test -- --run src/search src/App.test.tsx src/components/Shell.test.tsx
npm run build
```

Expected outcome: query validation, filtering, progressive results, compatible selection, confirmation, partial outcomes, focus, and route integration pass.

## Browser workflow

```powershell
Set-Location desktop/frontend
npm run test:e2e -- --grep "cross-daemon search"
```

Expected outcome: the workflow remains keyboard operable at 800 by 600 and 200 percent zoom with duplicate names, 100 profiles, mixed authority, and partial target failure.

## End-to-end scenarios

1. Create the same task name and group name on two daemons, produce one failed run and one unacknowledged alert, then search the shared name. Confirm every row exposes registration, daemon identity, state, and freshness.
2. Disable one remote daemon and repeat the search. Confirm successful targets remain usable while the failed target reports its own recovery guidance.
3. Select compatible tasks across two Operate-capable daemons, request run-now, and inspect the confirmation grouping. Stop one daemon after confirmation. Confirm independent accepted and unavailable outcomes with no rollback claim.
4. Search through an Observe-only profile. Confirm results and open remain available while mutation actions explain the missing Operate authority before confirmation.
5. Replace or edit a profile between search and submission. Confirm daemon identity mismatch prevents mutation.
6. Remove a task between search and submission. Confirm the exact missing object receives a rejected outcome and no similarly named task is selected.

## Repository verification

```powershell
& 'C:\Program Files\Git\bin\sh.exe' scripts/verify.sh all
```

Expected outcome: all canonical repository gates pass before publication.
