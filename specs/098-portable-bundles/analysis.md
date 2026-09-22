# S098 Spec and Implementation Analysis

## Result

The implementation matches the approved S098 safety boundary.

- A bundle contains deterministic, versioned group, task-scheduling, and completion-chain intent keyed by durable portable identities.
- Execution commands, environments, run-as identities, directories, stdin, daemon identity, trigger keys, watcher paths, notification credentials, run history, and audit history are absent from exported content.
- Preview stores one short-lived server-held plan. Apply consumes that plan once and rejects a changed target or a mismatched identity, digest, or fingerprint before mutation.
- Target-only state is reported as drift. Omission cannot delete or disable target records.
- The API, authenticated remote surface, CLI, and desktop page share the same daemon contract and authorization rules.

## Deliberate v1 boundary

Watchers, external triggers, trigger sets, and notification-policy references are reported as exclusions in v1 because their existing values contain paths, keys, or credential-bound references that cannot be safely translated without a separately approved portable representation. This avoids a lossy implicit approximation or accidental secret export.

## Verification evidence

- `go run ./scripts/github-format`
- `go test ./...`
- `go test ./...` in `desktop`
- `npm run build` in `desktop/frontend`
