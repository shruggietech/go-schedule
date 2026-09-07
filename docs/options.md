---
title: Desktop settings and local storage
nav_order: 4.5
---

# Desktop settings and local storage

The Wails desktop navigation rail keeps **Connections** and **Settings** available beside the task and operations workspaces. **Exit** remains a separate orderly application command. Settings combines the useful outcomes from the former Fyne Options and Info windows: appearance, preference transition, storage visibility, version details, and product links.

## Appearance and upgrade behavior

The Wails desktop stores one durable appearance choice for the current user:

| Setting | Choices | Default |
| --- | --- | --- |
| Color mode | System, Light, Dark | System |

Changes are saved immediately in the versioned Wails preference file beneath the current user's configuration directory at `go-schedule/desktop/preferences.json`. **Restore desktop defaults** changes only this desktop preference and returns appearance to System. It does not alter scheduler configuration, tasks, history, logs, or external files.

On the first Wails start, a valid `appearance.mode` value from the earlier Fyne `tech.shruggie.goschedule` preference file is migrated when no Wails preference file exists. Valid System, Light, and Dark choices are preserved. A missing, malformed, unreadable, or unsupported legacy value selects System without preventing startup. After the Wails preference file exists, it is authoritative and later changes to the legacy file do not overwrite it.

The former Fyne interface-font and scroll-sensitivity preferences are intentionally retired. Wails uses the bundled responsive type system and browser-native scrolling, so those toolkit-specific controls would no longer describe real behavior. They are not copied into the new preference file.

## Application storage

The Application storage section resolves known locations for the running platform and labels each with ownership, scope, whether the exact path is present, absent, or unavailable, what normal software removal does, and what an explicitly confirmed data wipe does.

Available exact paths have a **Copy path** action backed by the operating-system clipboard. The frontend sends only the storage record identifier; the Go service resolves the current path before copying it. Unavailable locations show no copy action. The inventory covers the machine data root, task database, machine configuration, logs, runtime state, per-user Wails application data and preferences, running executable directory, installed documentation when discoverable, and Windows maintenance evidence when applicable.

Daemon-owned rows come only from the connected daemon's effective runtime metadata, so custom data, database, configuration, log, and lock paths are represented accurately. When the daemon is unavailable, these rows remain visible as unavailable rather than displaying guessed defaults. A path configured outside the standard application-owned machine root is labeled external and preserved rather than described as an uninstall or wipe target.

This inventory is informational and read-only. It does not create, open, scan, delete, or relocate paths. User-created exports are not discovered, and administrator-configured locations outside application-owned defaults are never represented as wipe targets.

## Product information and links

Settings shows the application name, build version, publisher, source repository, and documentation. Link buttons send a fixed product-link identifier to Go, which opens only the corresponding application-defined HTTPS destination in the operating system default browser. Arbitrary frontend URLs are not accepted.

## Connection recovery

Connections displays the current protected local daemon state, safe platform and version context, capabilities, last successful connection when available, and diagnosis-specific guidance for unavailable, timed out, access denied, incompatible, degraded, and recovering states. **Try again** uses the existing bounded connection manager and suppresses duplicate pending retries.

Connection changes remain inline. They do not automatically open or repeatedly display a modal. Desktop-local preferences and product information remain usable while the daemon is offline, while daemon-backed storage information and scheduler mutations accurately remain unavailable.
