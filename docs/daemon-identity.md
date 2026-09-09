---
title: Daemon identity
nav_order: 8.6
---

# Daemon identity

Every initialized daemon has one opaque installation identifier and one editable display name. The identifier distinguishes same-named daemon targets without deriving identity from a hostname, address, account, storage path, product version, or scheduler record. The default display name is `go-schedule daemon` and can be changed without changing identity.

## Lifecycle

| Operation | Installation identity | Display name | Scheduler data |
|---|---|---|---|
| Clean install without retained data | New | Generic default | New store |
| Ordinary restart | Preserved | Preserved | Preserved |
| Software upgrade | Preserved | Preserved | Preserved |
| Database backup restore | Restored | Restored | Restored |
| Database copy or clone | Copied until deliberate reset | Copied | Copied |
| Rename | Preserved | Updated | Preserved |
| Confirmed identity reset | Replaced | Preserved | Preserved |

A database backup represents one logical daemon, so restoring it restores identity. Copying a database also copies identity. Before operating both copies independently at the same time, reset one clone by confirming its exact current identifier:

```sh
gosched daemon manifest
gosched daemon reset-identity --confirm <current-installation-id>
```

Reset is intentionally compare-and-confirm. A missing, mismatched, or stale identifier is rejected without changing anything. A successful reset atomically replaces only the installation identifier. It does not delete schedules, tasks, groups, automation sources, history, notification configuration, or the display name.

## Rename and discovery

```sh
gosched daemon manifest
gosched daemon manifest --json
gosched daemon rename "Workshop scheduler"
```

Display names are trimmed, must contain 1 through 80 Unicode characters, and cannot contain control characters. Names are labels only. Clients must use the installation identifier when they need stable target identity.

The manifest reports the product version, supported local and remote API versions, operating mode, sorted feature capabilities, and OS/architecture compatibility facts. It does not report hostname, network address, storage path, account, environment, command, credential, trigger key, scheduler record, or identity timestamps.

## Current access boundary

Identity discovery and lifecycle operations use the existing protected local IPC endpoint. S075 advertises local API `v1`, an empty remote API version list, and `local_only` mode. It does not open a listener, issue credentials, or provide actor-scoped authorization. Those controls remain staged in the [remote access architecture](remote-access.md).
