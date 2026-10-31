---
title: Installation
nav_order: 2
has_children: true
---

# Installing go-schedule

go-schedule ships per-platform. Pick the guide for your operating system:

- **[Windows](INSTALL-windows.md)**, the `.msi`, the Windows service, `PATH`, upgrading, and uninstalling.
- **[Linux](INSTALL-linux.md)**, the release archive, systemd registration, and data paths.
- **[macOS](INSTALL-macos.md)**, the desktop bundle versus headless, launchd, and the boot-persistence caveat that catches people out.

Each guide is self-contained. Once installed, the [`gosched` command reference](cli.md) covers every command and flag.

The Wails control center is the maintained desktop application and uses the `gosched-gui` launch name. Published v2.0.0 is available from the [latest-release page](https://github.com/shruggietech/go-schedule/releases/latest). This release adopts a current-only prelaunch state policy: fresh stores initialize directly, valid current stores reopen, and historical formats require an explicit developer-controlled reset. Desktop preferences start from current defaults when absent and do not import Fyne state. See [Current state and explicit reset](#current-state-and-explicit-reset).

## Current state and explicit reset

Published v2.0.0 provides the current-only prelaunch baseline. Earlier release assets retain their original behavior. Application databases retain marker 24 and their current schema. Current data and daemon identity reopen unchanged. Unsupported, foreign, partial, or malformed state is refused without conversion or automatic deletion. Metrics helpers use their shared current format 4 and likewise refuse older layouts.

For disposable development state, stop every process using the affected path. Either select a new empty data/database path or explicitly remove the obsolete development database together with its SQLite sidecar files after confirming the exact path and any desired backup. A new application store creates a new daemon identity. Do not delete an active database or treat another application's database as disposable. There is no automatic reset command for application or helper databases. Current user preferences have the separate explicit reset command `gosched preferences reset --scope all`; it does not alter daemon data.
