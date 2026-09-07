# Release Candidate Contract

## External identity

- The supported desktop command name remains `gosched-gui`.
- `gosched gui` locates that executable beside the CLI.
- Windows shortcuts and installer process handling target `gosched-gui.exe`.
- The macOS bundle is `gosched-gui.app` and its executable is `gosched-gui`.
- Linux desktop integration launches `gosched-gui`.

## Required desktop distributions

| Platform | Architecture | Distribution | Required payload |
| --- | --- | --- | --- |
| Windows | amd64 | MSI | gosched-gui.exe, goschedd.exe, gosched.exe, installer-private cleanup helper, service and shortcut contracts |
| macOS | arm64 | tar.gz containing gosched-gui.app | Wails executable, goschedd, gosched, icon, bundle metadata, README, LICENSE, CHANGELOG |
| Linux | amd64 | tar.gz | Wails executable, goschedd, gosched, freedesktop entry, icon tree, README, LICENSE, CHANGELOG |

Server-only root-module archives remain Linux and macOS amd64/arm64 and contain no desktop runtime.

## Required candidate gates

1. Root format, vet, lint, race, coverage, documentation, automation, and specification lifecycle checks.
2. Desktop Go race suite.
3. Frontend dependency audit, unit suite, production build, and Chromium accessibility/responsive suite.
4. Native Wails build on Windows, macOS, and Linux.
5. Windows MSI source inspection, compiled MSI inspection, and silent install/upgrade/repair/uninstall lifecycle.
6. Package payload assertions for all desktop matrix members.
7. Current-product residue scan that rejects Fyne dependencies, source entry points, CI jobs, release steps, and present-tense support claims.

Every hosted result is valid only for the pull-request head revision. Any required missing, skipped, cancelled, timed-out, or failed result rejects the candidate.

## Upgrade and rollback boundary

- The Windows UpgradeCode, product name, install directory, service identity, CLI PATH component, shortcut feature identities, bundle identifier, and stable executable name remain unchanged.
- Upgrade replaces the Fyne binary at the same path with Wails and preserves daemon-owned data.
- Wails performs its bounded one-time appearance migration from the legacy preference location.
- Default uninstall preserves data; explicit wipe remains opt-in and limited to documented application-owned roots.
- Rolling back the code before release is a Git revert. Rolling back an installed candidate uses the documented installer lifecycle and does not reinterpret persisted scheduler data.

## Publication boundary

Passing this contract makes a reviewed commit eligible for maintainer acceptance. It does not create a tag, publish a draft or public release, sign an artifact, close the milestone, or claim availability.
