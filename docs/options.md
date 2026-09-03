---
title: Desktop options and local storage
nav_order: 4.5
---

# Desktop options and local storage

The desktop navigation rail keeps the application views in this order: **Tasks**,
**Groups**, **Chains**, **Schedule**, **Activity**, **Options**, and **Info**.
**Exit** is a separate command at the bottom of the rail. It performs the same
orderly shutdown as the title-bar close control.

## Appearance

**Options** provides two independent appearance choices:

| Setting | Choices | Default |
| --- | --- | --- |
| Color mode | Dark, Light, Follow system | Dark |
| Interface font | Brand, System, Monospace | Brand |

Changes apply immediately across the current interface and are saved for the
current user. **Restore defaults** returns both settings to Dark and Brand.
Unsupported or damaged saved values also fall back to those defaults rather
than preventing startup.

The Brand choice uses the bundled Geist body face, Space Grotesk headings, and
Geist Mono where the interface requests fixed-width text. System delegates to
the platform font selected by the GUI framework. Monospace uses bundled Geist
Mono throughout non-symbol interface text. The choices do not load external
font files.

## Application storage

The Application storage section resolves known locations for the running
platform and labels each with:

- its ownership scope;
- whether the exact path is present, absent, or could not be inspected;
- what software-only removal does; and
- what an explicitly confirmed data wipe does.

Available paths are selectable and have a **Copy** action. The inventory covers
the machine data root, task database, machine configuration, logs, runtime
state, per-user desktop application data and preferences, the running executable
directory, installed documentation when discoverable, and Windows maintenance
evidence when applicable.

Daemon-owned rows come from the connected daemon's effective configuration, so
launching `goschedd` with a custom configuration file displays that file and the
actual derived data, database, log, and lock paths. A path configured outside
the standard application-owned roots is shown as external and preserved rather
than being described as an uninstall or wipe target. On Linux and macOS, the
view states that the Windows-only guided data wipe is unavailable.

This view is informational and read-only. It does not create, open, scan,
delete, or relocate paths. User-created exports are not discovered, and
administrator-configured locations outside the application-owned defaults are
never represented as wipe targets. A development executable is identified as
the running application location, not promised to be installer-owned.
