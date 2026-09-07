# Contract: Desktop Settings Bridge

## Workspace

### `SettingsWorkspace() SettingsResult`

Loads or initializes the versioned desktop preferences, resolves all local storage records, requests authoritative daemon runtime metadata once, and returns product information. Daemon failure does not fail the complete workspace; affected records are marked unavailable and `daemonAvailable` is false. A local preference or inventory initialization failure returns unavailable without a partial workspace.

## Preferences

### `SaveAppearance(mode) SettingsResult`

Accepts only `system`, `light`, or `dark`. Atomically persists the complete current preference document and returns a refreshed workspace. An invalid value is rejected without a write. A failed write preserves the prior file and active frontend appearance.

### `RestoreDesktopPreferences() SettingsResult`

Atomically persists the current schema with system appearance and a `not_required` transition, then returns a refreshed workspace. It does not delete daemon configuration, task data, logs, or external content.

## Native actions

### `CopyStoragePath(id) SettingsResult`

Resolves the current Settings workspace, finds the exact backend-owned record identifier, and writes its available path using native clipboard integration. Unknown, unavailable, or non-copyable identifiers are rejected. Success and failure are returned as an announcement-safe result.

### `OpenProductLink(key) SettingsResult`

Resolves the key against the fixed application allowlist and opens the HTTPS destination using native browser integration. Unknown keys and non-HTTPS destinations are rejected.

## Shared guarantees

- Every daemon call has an existing bounded timeout.
- Frontend callers cannot provide arbitrary clipboard content or URLs.
- Error messages use a small safe vocabulary and do not reflect raw preference content.
- Preference writes use a same-directory temporary file and atomic rename where the platform permits.
- Established Wails preferences remain authoritative after migration.
- Settings remains locally useful when the daemon is unavailable.
- Product destinations and daemon paths are never guessed in React.

## Connection recovery

Connections uses the existing frontend connection manager contract:

- read the current `ConnectionSnapshot`;
- call its single retry operation;
- suppress duplicate pending retries;
- preserve route, focus, and last successful application data;
- present diagnosis and guidance inline without automatic modal display.
