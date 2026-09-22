# S098: Portable automation bundles and drift reporting

## Summary

Implements the safe v1 portable-bundle workflow for #184 across the daemon API, authenticated remote API, CLI, and desktop control center.

## Delivered

- Adds deterministic `go-schedule.bundle/v1` documents for groups, task scheduling intent, and completion chains.
- Persists durable portable identities separately from daemon-local IDs through SQLite migration v22.
- Adds export, validate, compare, target-bound preview, and single-use apply endpoints with Observe and Manage authorization.
- Applies only reviewed create or update operations in dependency order. Imported tasks are command-free disabled drafts, and updates preserve existing local execution inputs.
- Adds `gosched bundle export`, `validate`, `compare`, `preview`, and `apply` commands.
- Adds a desktop Bundles page with selected-target context, validation findings, read-only drift comparison, explicit review confirmation, and per-item outcomes.

## Safety boundaries

- Bundles never export daemon identity, credentials, trigger keys, notification endpoints or authorization, run history, audit records, commands, environments, users, directories, or stdin.
- Preview is stored by the selected daemon for 15 minutes and consumed once by apply. Changed targets, mismatched bindings, expired previews, and replay attempts are rejected before mutation.
- Omitted target objects are reported as `target_only` drift. V1 has no delete, disable, synchronization, background reconciliation, transaction, or retry behavior.
- Watchers, external triggers, trigger sets, and notification-policy references are reported as exclusions until they have an approved secret-free portable representation.

## Validation performed

- `go run ./scripts/github-format`
- `go test ./...`
- `go test ./...` from `desktop`
- `npm run build` from `desktop/frontend`

## Follow-up

#184 remains open and In progress for the excluded source families. This pull request does not close it implicitly.
