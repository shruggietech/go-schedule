---
title: Actor Permissions and Management Audit
nav_order: 5.7
---

# Actor Permissions and Management Audit

go-schedule uses one server-owned authorization vocabulary for protected local management and the opt-in remote transport defined by the [remote access architecture](remote-access.md). Local IPC remains the default and authenticates the operating-system boundary. An explicitly configured HTTPS listener authenticates each opaque bearer credential to its current actor and applies the same operation catalog without adding a general policy language.

## Capability hierarchy

Capabilities are closed and monotonic. A stronger capability includes every weaker capability.

| Capability | Intended operations |
|---|---|
| Observe | Safe task, schedule, run, alert, log, manifest, and status reads |
| Operate | Observe plus run-now, trigger fire, alert acknowledgement, and notification tests |
| Manage | Operate plus task, group, chain, watcher, trigger, notification, and daemon-name configuration |
| Enroll | Manage plus actor, credential, listener, identity, secret rotation, secret reveal, and audit administration |

Every registered management route has a stable operation identifier, minimum capability, target kind, and audit classification in one shared catalog. Unknown operations and unknown capability values fail closed.

## Local trust boundary

Schema v17 creates exactly one built-in `local_os` actor with Enroll capability. Existing Unix-socket or Windows named-pipe ownership and permissions continue to authenticate local clients, and the API resolves those requests to the built-in actor without a new header, token, or login. The actor cannot be revoked, expired, deleted, or reduced. Local CLI and desktop workflows therefore remain compatible while their management actions gain attribution at the operating-system boundary.

Non-built-in actors may use the `desktop`, `cli`, `json`, or `mcp` kind and one capability. Actor records contain no credential material. A remote pairing transaction creates one actor and one separate digest-only credential; creating an actor directly grants no network access by itself. Actor and credential state are reloaded for every remote request, so a committed revocation denies the next request even when a transport connection is reused.

## Management audit

Mutations and privileged reads require a durable audit intent before the protected handler runs. Failure to persist the intent returns `503 audit_unavailable` and prevents the operation. The intent starts as `uncertain`, then the same event becomes `succeeded` or `failed` after the handler returns. A crash or process interruption leaves `uncertain` instead of claiming success. Authorization denials are recorded directly as `denied` without running the operation.

Each event stores only its event identifier, known actor identifier, daemon identifier, operation identifier, target classification and optional identifier, result, correlation identifier, occurrence time, and optional completion time. It never stores request or response bodies, authorization headers, raw errors, commands, environment variables, standard input, filesystem paths, cryptographic keys, credentials, or secrets. Scheduler run history and logs remain separate from management audit.

On every new audit record, the same transaction removes events older than 90 days and trims history to the newest 10,000 events. Queries default to 100 events and permit at most 1,000. Results sort by occurrence time and event identifier. Filters cover actor, operation, result, start time, and end time; newline-delimited JSON export uses the same deterministic filter and ordering rules.

## Actor administration

The local API and typed client support listing, creating, updating, and revoking actors. The CLI equivalents are:

```text
gosched actor list
gosched actor create "Release operator" --kind cli --capability manage
gosched actor update <actor-id> --capability observe
gosched actor update <actor-id> --expires 2026-12-01T00:00:00Z
gosched actor update <actor-id> --clear-expiration
gosched actor revoke <actor-id>
```

Names are trimmed, contain 1 through 80 Unicode characters, and reject control characters. Revocation is irreversible. Expiration can be set or cleared on a non-built-in actor.

## Audit inspection and export

```text
gosched audit list
gosched audit list --actor <actor-id> --operation tasks.create --result succeeded --since 2026-09-01T00:00:00Z --limit 250
gosched audit export --result denied
gosched audit export --output management-audit.ndjson
```

The export file is created with owner-only permissions. Treat it as administrative evidence even though sensitive execution content is structurally excluded.

## Upgrade and future integration

Opening an older database applies forward-only schema v18, preserves scheduler and daemon identity data, creates actor, audit, pairing, and credential storage as needed, and initializes one built-in local actor. Reopening is idempotent. The remote HTTPS adapter authenticates the presented credential to its current actor on every request, then calls this same operation catalog and authorizer. Pairing creation and credential lifecycle remain local Enroll operations with intent-first audit; bearer values, phrase verifiers, request bodies, and certificate keys never enter audit records.
