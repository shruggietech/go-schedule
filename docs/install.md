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

Beginning with v1.4.0, the Wails control center is the sole maintained desktop application. It preserves the `gosched-gui` launch name used by earlier releases, so supported in-place upgrades replace the implementation without changing shortcuts or command habits. Existing daemon-owned tasks and run history remain in their platform data location, and the first Wails start migrates only a supported legacy appearance choice. The [latest-release page](https://github.com/shruggietech/go-schedule/releases/latest) is authoritative while a candidate is staged and qualified.
