---
title: Local API
nav_order: 5.5
---

# Local and remote JSON API

**Audience:** client and integration authors\
**Applies to:** the current unreleased local API contract\
**Transport:** protected local IPC by default; optional authenticated TLS 1.3 HTTPS on the reviewed `/api/v1` allowlist

The CLI and desktop app use the same versioned JSON handlers hosted by `goschedd`. Local IPC paths begin at `/v1`; the opt-in remote adapter exposes only reviewed operations below `/api/v1`. Errors use `{"error":{"code":"...","field":"...","message":"..."}}`.

The [remote access guide](remote-access.md) defines listener setup, deployment modes, pairing, profiles, revocation, and trust ownership. Installing or upgrading never enables a listener.

## Direct JSON client workflow

A direct client must retain the operator-provided certificate, expected daemon installation ID, credential ID, and bearer value in protected storage. First exchange a one-time phrase without following redirects, then compare both the response `daemon_id` and `GET /api/v1/manifest` `installation_id` to the expected value before retaining or using the credential. Never disable TLS verification or accept an identity change automatically.

Read the phrase and bearer from protected input rather than a URL, source file, or command argument. This shell sketch intentionally keeps both values out of command history:

```sh
read -rs PAIRING_PHRASE
printf '{"daemon_id":"%s","pairing_id":"%s","phrase":"%s","display_name":"JSON client","kind":"json","capability":"observe"}' "$DAEMON_ID" "$PAIRING_ID" "$PAIRING_PHRASE" | curl --fail --proto =https --tlsv1.3 --max-redirs 0 --cacert daemon.pem --json @- "$ENDPOINT/api/v1/enroll"
read -rs BEARER
printf 'header = "Authorization: Bearer %s"\n' "$BEARER" | curl --fail --proto =https --tlsv1.3 --max-redirs 0 --cacert daemon.pem --config - "$ENDPOINT/api/v1/tasks?limit=100"
unset PAIRING_PHRASE BEARER
```

Collection responses use bounded limits and continuation parameters where documented by the generated OpenAPI contract in `api/openapi/remote-v1.yaml`. Classify stable error envelope codes before considering HTTP text. Treat `401 unauthorized` as missing or invalid authentication, `401 credential_revoked` as a profile that requires new enrollment material, `403` as insufficient actor authority, `409` as identity or compatibility conflict, `429` as rate limiting, and `5xx` as server failure. Use bounded timeouts for every request.

GET and other retry-safe operations may be retried deliberately. The desktop connection owner retries transient health and subscription failures with jittered delays that begin below one second and cap at thirty seconds. It stops for authentication, revocation, authorization, compatibility, certificate-trust, and daemon-identity failures. If a mutation response is lost or otherwise ambiguous, do not replay it automatically: report the outcome as uncertain, refresh the authoritative resource, determine whether the first mutation committed, then require a deliberate retry if needed. Server-Sent Events are invalidation hints; reconnecting clients refresh authoritative state rather than depending on replay.

## Daemon identity and capability manifest

The [daemon identity lifecycle](daemon-identity.md) defines stable installation identity, editable names, restore and clone behavior, and reset safety. The manifest is deliberately separate from health so existing liveness clients remain compatible.

| Method | Path | Result |
|---|---|---|
| `GET` | `/v1/manifest` | Stable identity, display name, product version, protocol versions, operating mode, sorted capabilities, and safe platform facts |
| `PATCH` | `/v1/manifest` | Rename from `{"display_name":"Workshop scheduler"}` and return the resulting manifest |
| `POST` | `/v1/manifest/reset` | Compare and replace identity from `{"confirm_installation_id":"<current-id>"}` and return the resulting manifest |

The current manifest reports local API `v1`, no remote API version, and `local_only` operating mode. It excludes hostnames, network addresses, storage paths, accounts, environment values, commands, credentials, trigger keys, scheduler records, and lifecycle timestamps. A stale reset confirmation returns `409 conflict`; invalid names and malformed requests return `400 validation_failed`; neither failure mutates state. `GET /v1/health` remains unchanged.

## Actor permissions and management audit

The [actor permissions and management audit contract](access-control.md) applies one Observe, Operate, Manage, and Enroll hierarchy to every registered management operation. Protected local IPC resolves to the built-in local actor without new authentication input.

| Method | Path | Minimum capability | Result |
|---|---|---|---|
| `GET` | `/v1/access/actors` | Enroll | List credential-independent actors |
| `POST` | `/v1/access/actors` | Enroll | Create a non-built-in actor |
| `PATCH` | `/v1/access/actors/{id}` | Enroll | Update name, capability, state, or expiration |
| `POST` | `/v1/access/actors/{id}/revoke` | Enroll | Irreversibly revoke an actor |
| `GET` | `/v1/audit` | Enroll | List filtered retained audit events |
| `GET` | `/v1/audit/export` | Enroll | Export filtered events as newline-delimited JSON |

Audit filters are `actor_id`, `operation`, `result`, `since`, `until`, and `limit`. The default limit is 100 and the maximum is 1,000. Invalid filters return `400 validation_failed`. Authorization denial returns `403 forbidden`; failure to persist mandatory audit evidence returns `503 audit_unavailable`.

## External triggers

An ordinary trigger representation contains `id`, `name`, optional `set_id`, optional `set_name`, optional `set_position`, `target_task_id`, `target_task_name`, `enabled`, `readiness`, `reason`, `created_at`, and `updated_at`. It never contains the raw key. A set member cannot be retargeted individually; use the set endpoint so every member retains the shared target invariant.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/triggers` | `{"triggers": [...]}` with redacted records |
| `POST` | `/v1/triggers` | Create from name, target task, and optional enabled state; `201` with the key |
| `GET` | `/v1/triggers/{id}` | One redacted trigger |
| `PATCH` | `/v1/triggers/{id}` | Update name or target task |
| `DELETE` | `/v1/triggers/{id}` | Delete; `204` |
| `POST` | `/v1/triggers/{id}/enable` | Enable the trigger |
| `POST` | `/v1/triggers/{id}/disable` | Disable the trigger |
| `POST` | `/v1/triggers/{id}/rotate` | Atomically replace and return the key |
| `POST` | `/v1/triggers/{id}/reveal` | Explicitly return the current key |
| `POST` | `/v1/triggers/fire` | Accept `{"key":"gst_..."}` and submit one run request; `202` |

Fire failures use stable codes: `trigger_unknown`, `trigger_disabled`, `trigger_target_missing`, `trigger_command_incomplete`, `trigger_task_inactive`, `trigger_task_disabled`, `trigger_group_blocked`, or `trigger_dispatch_unavailable`. Error responses never echo the submitted key.

## Filesystem watchers

A filesystem watcher representation contains `id`, `name`, `kind`, `path`, optional `pattern`, `recursive`, string durations `debounce` and `stability`, target identity, enabled state, runtime `health`, readiness, reason, and timestamps. `kind` is `file` for one exact path or `directory` for basename-glob selection. A directory pattern defaults to `*`; recursive selection excludes linked directories. Durations accept Go duration syntax such as `250ms`, `2s`, or `1m` and must be from 25 milliseconds through one hour.

| Method | Path | Behavior |
|---|---|---|
| `GET` | `/v1/filesystem-watchers` | List definitions joined with current runtime health |
| `POST` | `/v1/filesystem-watchers` | Create a definition and reload observation; `201` |
| `GET` | `/v1/filesystem-watchers/{id}` | Show one definition and current health |
| `PATCH` | `/v1/filesystem-watchers/{id}` | Atomically update selection, timing, target, enabled state, or name and reload observation |
| `DELETE` | `/v1/filesystem-watchers/{id}` | Delete, cancel pending candidates, and reload observation; `204` |
| `POST` | `/v1/filesystem-watchers/{id}/enable` | Enable and reload observation |
| `POST` | `/v1/filesystem-watchers/{id}/disable` | Disable, cancel pending candidates, and reload observation |

Watcher lifecycle and health events use `kind: "filesystem_watcher"` with identity, name, verb, and optional health only. They omit configured and matched paths. Filesystem-originated runs have `trigger: "filesystem_watcher"` and `source_watcher_id`; they never retain the matched path.

### Trigger Sets

Ordinary Trigger Set representations include stable set identity, name, target, member and enabled counts, ordered redacted members, and timestamps. Create, reveal, and rotate responses additionally contain ordered member keys and complete commands.

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/trigger-sets` | List redacted Trigger Sets and members |
| `POST` | `/v1/trigger-sets` | Atomically create 1 through 99 members and return ordered secrets |
| `GET` | `/v1/trigger-sets/{id}` | Show one redacted Trigger Set |
| `PATCH` | `/v1/trigger-sets/{id}` | Atomically retarget every member |
| `DELETE` | `/v1/trigger-sets/{id}` | Atomically delete the set and every member; `204` |
| `POST` | `/v1/trigger-sets/{id}/enable` | Atomically enable every member |
| `POST` | `/v1/trigger-sets/{id}/disable` | Atomically disable every member |
| `POST` | `/v1/trigger-sets/{id}/rotate` | Atomically rotate every key and return ordered replacement secrets |
| `POST` | `/v1/trigger-sets/{id}/reveal` | Explicitly return current ordered secrets |

## Runtime storage information

`GET /v1/runtime-info` returns the absolute effective paths used by the running daemon: `data_dir`, `database_path`, optional `config_path`, `log_path`, and `lock_path`. Desktop clients use this endpoint for read-only storage disclosure, including when the daemon was launched with a custom configuration path.

## Completion chains

The chain representation contains `id`, `source_task_id`, `source_task_name`, `target_task_id`, `target_task_name`, `on_outcome`, `created_at`, and `updated_at`.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/chains` | `{"chains": [...]}` |
| `POST` | `/v1/chains` | Create from source, target, and outcome; `201` |
| `GET` | `/v1/chains/{id}` | One chain |
| `PATCH` | `/v1/chains/{id}` | Update any non-empty subset of mutable fields |
| `DELETE` | `/v1/chains/{id}` | Delete; `204` |

Create example:

```text
{
  "source_task_id": "source-id",
  "target_task_id": "target-id",
  "on_outcome": "success"
}
```

`on_outcome` is `success`, `failure`, or `any`. Missing resources return `not_found`. Invalid outcomes, self-links, duplicates, and direct or indirect cycles return `validation_failed` without partial mutation.

## Correlated history and events

Completion-triggered entries from `GET /v1/runs` have `trigger` set to `completion` and include optional `source_task_id` and `source_run_id`. Externally triggered entries have `trigger` set to `external_trigger` and include `source_trigger_id`. Filesystem-triggered entries have `trigger` set to `filesystem_watcher` and include `source_watcher_id`. Raw keys and matched file paths are never stored in run history.

`GET /v1/events` emits `kind: "chain"` with a `created`, `updated`, or `deleted` verb. Create and update include the current chain; delete carries its stable ID.

Trigger lifecycle events use `kind: "trigger"` and contain only a redacted trigger or its stable deletion ID.

Trigger Set lifecycle events use `kind: "trigger_set"` and contain set identity, name, member count, and verb without member keys. One set-level mutation publishes one event after its transaction commits.

## Webhook notifications

Notification channels are reusable write-only webhook destinations. Ordinary channel responses contain `id`, `name`, `kind`, `endpoint_summary`, `has_authorization`, `enabled`, `created_at`, and `updated_at`. See [Webhook notifications](notifications.md) for the receiver payload, precedence, retry, duplicate, and security contracts.

| Method | Path | Result |
| --- | --- | --- |
| `GET` | `/v1/notification-channels` | `{"notification_channels": [...]}` with redacted metadata |
| `POST` | `/v1/notification-channels` | Create from `name`, `endpoint`, optional `authorization`, and optional `enabled`; `201` |
| `GET` | `/v1/notification-channels/{id}` | One redacted channel |
| `PATCH` | `/v1/notification-channels/{id}` | Update optional `name`, `endpoint`, or `enabled` |
| `DELETE` | `/v1/notification-channels/{id}` | Remove assignments and unfinished work while preserving safe terminal history; `204` |
| `POST` | `/v1/notification-channels/{id}/enable` | Enable new delivery creation |
| `POST` | `/v1/notification-channels/{id}/disable` | Disable new delivery creation |
| `POST` | `/v1/notification-channels/{id}/rotate` | Replace or clear the write-only `authorization` value |
| `POST` | `/v1/notification-channels/{id}/test` | Queue a transport-equivalent test delivery; `202` |
| `GET` | `/v1/tasks/{id}/notifications` | List direct task assignments |
| `PUT` | `/v1/tasks/{id}/notifications` | Atomically replace direct task assignments |
| `GET` | `/v1/tasks/{id}/notifications/effective` | Return selected source scope and complete effective assignments |
| `GET` | `/v1/groups/{id}/notifications` | List direct group assignments |
| `PUT` | `/v1/groups/{id}/notifications` | Atomically replace direct group assignments |
| `GET` | `/v1/notification-deliveries` | List redacted evidence with optional `channel`, `task`, `run`, `state`, and `limit` filters |

An assignment contains `channel_id`, `on_success`, and `on_failure`; at least one outcome must be true and one channel may occur only once per scope. The nearest non-empty task or group scope replaces all more distant assignments. Delivery responses remain distinct from run responses and include a safe parsed `event`, attempts, timestamps, last status, and bounded diagnostic without protected endpoint or authorization fields.

## Remote HTTPS API

The optional network API is a separate allowlisted adapter under `/api/v1`; it does not expose the complete local `/v1` mux. Its canonical OpenAPI 3.1 description is `api/openapi/remote-v1.yaml`. Health, manifest discovery, and one-time enrollment are public within the trusted TLS boundary. Every other listed operation requires one strict `Authorization: Bearer` header and is authorized against the credential's current actor capability before the existing local handler runs.

The enrollment request is `POST /api/v1/enroll` with JSON fields `pairing_id`, `phrase`, and `daemon_id`. Successful exchange returns the bearer token exactly once. Authentication failures use `401 unauthorized`, browser-origin requests use `403 origin_rejected`, unknown or deliberately excluded routes use `404 not_found`, and rate limits use `429 rate_limited` with `Retry-After: 1`.
