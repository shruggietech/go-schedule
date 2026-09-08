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

The v1.2 candidate promotes the Wails control center as the sole maintained desktop application. It preserves the `gosched-gui` launch name used by earlier releases, so supported in-place upgrades replace the implementation without changing shortcuts or command habits. Existing daemon-owned tasks and run history remain in their platform data location, and the first Wails start migrates only a supported legacy appearance choice. The latest public release remains v1.1.1 until the separately authorized v1.2 release ritual completes.
