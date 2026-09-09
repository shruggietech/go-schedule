# Contract: Local Daemon Manifest and Identity Lifecycle

All routes use the existing protected local IPC transport and JSON error envelope. They do not open a TCP listener or provide remote authorization.

## Read manifest

`GET /v1/manifest`

Successful response: `200 OK`

```json
{
  "installation_id": "095a4f62-c580-4f62-9ec4-f52f672b9ad7",
  "display_name": "go-schedule daemon",
  "product_version": "dev",
  "local_api_versions": ["v1"],
  "remote_api_versions": [],
  "operating_mode": "local_only",
  "capabilities": ["activity", "agent-access", "chains", "groups", "notifications", "schedule", "tasks", "triggers", "watchers"],
  "platform": {
    "os": "windows",
    "architecture": "amd64"
  }
}
```

The response must not include hostname, network address, storage path, account name, environment value, credential, trigger key, command, scheduler record, or identity lifecycle timestamp.

## Rename daemon

`PATCH /v1/manifest`

Request:

```json
{
  "display_name": "Workshop scheduler"
}
```

Successful response: `200 OK` with the resulting manifest.

Invalid JSON, unknown fields, blank names, names longer than 80 Unicode code points, or names containing control characters return `400 Bad Request` and do not mutate identity state.

## Reset installation identity

`POST /v1/manifest/reset`

Request:

```json
{
  "confirm_installation_id": "095a4f62-c580-4f62-9ec4-f52f672b9ad7"
}
```

Successful response: `200 OK` with a manifest containing a newly generated installation identifier and the unchanged display name.

Invalid JSON or a missing confirmation returns `400 Bad Request`. A well-formed confirmation that does not exactly match the current identifier returns `409 Conflict`. Every failure leaves identity and scheduler state unchanged.

## Compatibility

`GET /v1/health` retains its existing status code and JSON fields. Existing local routes retain their paths, methods, and response contracts. Unknown future manifest fields remain safe for Go JSON clients to ignore.
