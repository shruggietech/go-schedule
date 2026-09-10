# Quickstart: Validate Remote Connection Resilience

## Prerequisites

- One local desktop build and one remote daemon profile paired through the S078 workflow.
- A disposable remote daemon and credential suitable for interruption and revocation tests.
- Go, Node.js, the Wails build prerequisites, a POSIX shell, and the repository verification toolchain.

## Focused validation

1. Run the API client transport and uncertainty tests.
2. Run the desktop connection, profile-service, task, and operations tests under the race detector.
3. Run frontend connection-store, shell, Connections, task-editor, schedule, and activity tests.
4. Run Chromium connection-recovery accessibility and responsive scenarios.
5. Build the native desktop application to regenerate and validate Wails bindings.

## Manual recovery scenario

1. Select a paired remote profile and load Tasks, Schedule, and Activity.
2. Stop or isolate the remote daemon.
3. Confirm the same target remains selected, prior complete data is marked stale, retry timing is visible, and daemon-backed mutations are disabled.
4. Restore the daemon without changing its identity or certificate.
5. Confirm the desktop reconnects automatically, refreshes authoritative workspaces, clears stale status, and resumes live updates.
6. Drop a mutation response after the daemon receives the request and confirm the client reports an uncertain outcome without replaying it.
7. Revoke the credential, change the certificate, and substitute another daemon identity in separate trials. Confirm each stops automatic retry with distinct recovery guidance.

## Canonical verification

```bash
sh scripts/verify.sh all
```

The run is complete only when format, vet, lint, race, GUI, coverage, docs, and automation gates all pass in order.
