# Quickstart: Task and Group Authoring

## Prerequisites

- A clean checkout on `codex/062-task-group-authoring`
- Go 1.25 or automatic toolchain support
- Node.js 24 and exact frontend dependencies installed with `npm ci`
- A C compiler for local race verification
- Wails 2.14.0 build prerequisites for the current operating system

## Focused Go contracts

```bash
go test -race ./internal/commandexample ./internal/api/server ./internal/api/client ./test/integration
cd desktop
go test -race ./taskgroup ./...
```

Expected: detailed list compatibility, complete field clearing, task/group authority, stale drafts, safe errors, full paths, effective state, exact command parsing, coalesced writes, and the current platform's bounded safe command pass. The guided integration creates the task through the product boundary, invokes Run now, and observes a successful Activity record with recognizable captured output.

## Frontend contracts

```bash
cd desktop/frontend
npm ci
npm test
npm run build
npm run test:e2e
```

Expected: empty and one-hundred-row overviews, stable selection, task editing, validation focus, platform suggestion consistency, one-shot Tab insertion and traversal, group hierarchy, target-aware confirmations, stale-draft preservation, accessibility, reduced motion, narrow windows, and offline assets pass.

## Manual local journey

1. Start the local daemon using the normal development configuration.
2. Start the production-intent desktop application from `desktop/`.
3. Open Tasks and create a task from the empty command field.
4. Press Tab once to insert the platform suggestion and inspect Program and Arguments.
5. Save the task inactive, confirm Run now against This computer, and follow the Activity handoff message.
6. Create nested groups, move the task, disable an ancestor, and inspect declared versus effective state.
7. Edit advanced scheduling values, preview five occurrences, cancel once, then save and verify the authoritative refresh.

Expected: no unsafe default execution, hidden shell interpretation, lost focus, optimistic state, secret-bearing announcement, stale overwrite, or ambiguous group destination occurs.

## Canonical repository gate

```bash
sh scripts/verify.sh all
```

Expected: format, vet, lint, race, GUI, coverage, documentation, and automation gates all pass. Hosted CI additionally builds the Wails desktop and executes frontend contracts on Windows, macOS, and Linux as configured.
