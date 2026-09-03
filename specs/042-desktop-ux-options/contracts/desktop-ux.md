# Contract: Desktop UX and Options

## Navigation contract

1. Ordinary destinations appear in this order: Tasks, Groups, Chains,
   Schedule, Activity, Options, Info.
2. Options is immediately above Info.
3. Exactly one ordinary destination is visibly selected.
4. The rail reserves enough width for the longest supported label and symmetric
   horizontal padding at the supported 800 by 600 launch clamp.
5. Exit is keyboard reachable, visually separate from destinations, and pinned
   to the rail's bottom-right without becoming selected content.
6. Activity badge updates do not reorder destinations or collapse the rail.

## Appearance contract

1. Missing or invalid preferences produce Dark mode with Brand font.
2. The user can select Dark, Light, or Follow system and Brand, System, or
   Monospace independently.
3. Every valid combination applies during the initiating interaction and is
   restored through the established application identifier on next launch.
4. Reset applies and persists Dark plus Brand.
5. Follow system consumes the current Fyne theme variant. Explicit Dark and
   Light ignore that variant.
6. Symbol fonts and unhandled semantic theme roles delegate to Fyne defaults.
7. Version and attribution labels remain centered, unwrapped, semantically
   unchanged, and sourced from the existing dynamic version value.

## Storage contract

1. Options distinguishes machine data from current-user desktop preferences.
2. Daemon-owned rows use absolute effective paths reported by the connected
   daemon, including a custom configuration file and its derived data paths.
3. Each available row displays a clean resolved path, ownership scope, existence state,
   software-only removal behavior, and explicit-wipe behavior.
4. Selectable text and Copy produce the exact same path string.
5. Missing paths are reported absent without being created.
6. Indeterminate or platform-inapplicable paths are omitted or marked
   unavailable; no path is fabricated.
7. Resolution inspects only the declared path itself and never walks a parent,
   profile, drive, or sibling directory.
8. External user-created and custom-configured locations are not claimed as
   application-owned or wipe targets.
9. Non-Windows platforms do not claim that the Windows-only explicit data-wipe
   workflow will remove any path.
10. A running development executable is not represented as installer-owned or
   promised to be removed by uninstall.

`GET /v1/runtime-info` returns `data_dir`, `database_path`, optional
`config_path`, `log_path`, and `lock_path` for the running daemon. The paths are
absolute and reflect the daemon's effective configuration.

## Task interaction contract

1. A single click selects a task and opens no editor.
2. A double activation resolves the rendered task ID against current model
   state and enters the existing detail/editor workflow.
3. Reordering cannot redirect the activation to the new occupant of an old
   index; removal makes the activation a no-op.
4. Detail lookup failure retains the established degraded editor behavior.
5. New, toolbar Edit, and row double activation share a one-editor ownership
   guard released on every dialog close path.

## Shutdown contract

Title-bar close, connection-card Exit, and navigation Exit call one idempotent
coordinator. Across repeated or concurrent requests, run-context cancellation
and window close each occur at most once.

## Evidence boundary

Headless evidence may prove state, widget properties, theme resources, layout
geometry, and callbacks. It must not be described as proof of native Windows
glyph sharpness or DPI rendering. Exact-candidate Windows 11 evidence under #94
must cover 100 percent and at least one scaled-DPI setting, Dark and Light modes,
Info text, navigation spacing, and Exit placement before dependent visual issues
close.
