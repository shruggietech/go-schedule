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

**Options** provides three independent appearance choices:

| Setting | Choices | Default |
| --- | --- | --- |
| Color mode | Dark, Light, Follow system | Dark |
| Interface font | System, Geist (brand), Inter, Ubuntu, Monospace | System |
| Scroll sensitivity | 1x through 4x in 0.5x steps | 2x |

Changes apply immediately across the current interface and are saved for the
current user. **Restore defaults** returns the settings to Dark, System, and
2x. Unsupported or damaged saved values also fall back to those defaults rather
than preventing startup.

System delegates to the platform font selected by the GUI framework. Geist
(brand) uses the bundled Geist body face, Space Grotesk headings, and Geist
Mono where the interface requests fixed-width text. Inter and Ubuntu package
their upstream regular and bold faces; fixed-width text still uses Geist Mono.
Monospace uses bundled Geist Mono throughout non-symbol interface text. Symbols
always use the GUI framework's symbol face. The packaged Inter and Ubuntu files
retain their upstream OFL and Ubuntu Font Licence notices in `gui/assets/fonts/`;
no choice loads font files from the network or enumerates locally installed
fonts.

The font and color selectors show the active value when closed and omit it from
the menu of alternatives. Conventional vertical mouse-wheel steps are multiplied
by Scroll sensitivity in application-owned long views. Precision touchpad
deltas, keyboard navigation, scrollbar use, and drag scrolling retain the GUI
framework's behavior.

## Application storage

The Application storage section resolves known locations for the running
platform and labels each with:

- its ownership scope;
- whether the exact path is present, absent, or could not be inspected;
- what software-only removal does; and
- what an explicitly confirmed data wipe does.

Locations appear as compact aligned rows under Category, Location and removal
details, and Action headers. Long paths wrap vertically, and no horizontal
scrolling is introduced. Available paths are selectable and have a **Copy**
action. Unavailable locations mute the complete row and disable selection and
Copy. The inventory covers
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
