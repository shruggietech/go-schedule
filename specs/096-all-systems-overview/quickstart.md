# Quickstart: S096 All Systems Operational Overview

## Build and exercise locally

1. Start the daemon through the repository's normal development workflow.
2. Start the desktop frontend and backend.
3. Open **All Systems** from the sidebar.
4. Confirm **This computer** appears even when no remote profiles exist.
5. Add or retain saved profiles that cover reachable, unavailable, unauthorized, incompatible, and duplicate-label cases.
6. Refresh and confirm every current registration appears once, one slow target does not erase successful summaries, and attention filters isolate failures, alerts, or delivery problems.
7. Open representative Tasks, Schedule, Activity, and Notifications actions and confirm the exact source profile is selected before the route changes.

## Focused development checks

```powershell
go test ./internal/store ./internal/api/server ./internal/api/client ./desktop/systems
```

```powershell
Set-Location desktop/frontend
npm test
npm run build
```

## CI-parity

Run the repository CI-parity command documented by `docs/build-autopilot.md`, then run the GitHub publication formatter before committing:

```powershell
go run ./scripts/github-format
```

## Expected boundaries

- No overview state survives a desktop restart.
- Refresh never changes the currently selected daemon.
- Drill-down selects one exact target before navigation.
- Observe authority is sufficient for the summary endpoint.
- The overview does not provide bulk actions, shared ownership, failover, reconciliation, or cluster semantics.
