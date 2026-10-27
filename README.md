<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/assets/brand/go-schedule-horizontal-color.svg">
  <img alt="go-schedule" src="docs/assets/brand/go-schedule-horizontal-black.svg" width="100%">
</picture>

<p align="center">
  <a href="https://github.com/shruggietech/go-schedule/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/shruggietech/go-schedule/actions/workflows/ci.yml/badge.svg"></a>
<a href="https://github.com/shruggietech/go-schedule/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/shruggietech/go-schedule?color=58A6FF"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-62D9B7"></a>
  <a href="https://go-schedule.io/docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-go--schedule.io-58A6FF"></a>
</p>

**Open-source workload automation with readable schedules and supported cron.**

Schedule scripts, react to files and external events, and connect dependent jobs. A background daemon owns the schedule; the CLI, Wails desktop control center and JSON API let you manage it.

**Windows, Linux and macOS.** [Documentation](https://go-schedule.io/docs/) | [Latest release](https://github.com/shruggietech/go-schedule/releases/latest) | [Changelog](CHANGELOG.md) | [Compare capabilities](https://go-schedule.io/compare/)

## What you can do

- Write readable schedules such as `weekdays at 09:00` or supported cron such as `0 9 * * 1-5`, with explicit timezone, DST, missed-date, catch-up and overlap policies.
- Start tasks from schedules, daemon startup, file changes, external triggers or another task's success/failure.
- Group related tasks and inspect upcoming work, execution history, output and diagnostics.
- Manage explicitly enrolled independent daemons, search across systems and move portable configuration.
- Integrate through webhooks, the JSON API and optional permission-controlled MCP access.

## Install

Download from the [latest public release](https://github.com/shruggietech/go-schedule/releases/latest). Verify artifacts against `SHA256SUMS.txt`; artifacts are currently unsigned.

| Platform | Package | Guide |
| --- | --- | --- |
| Windows amd64 | `go-schedule_<ver>_windows_amd64.msi` | [Windows](docs/INSTALL-windows.md) |
| Linux desktop amd64 | `go-schedule-desktop_<ver>_linux_amd64.tar.gz` | [Linux](docs/INSTALL-linux.md) |
| macOS desktop arm64 | `go-schedule-desktop_<ver>_darwin_arm64.tar.gz` | [macOS](docs/INSTALL-macos.md) |
| Linux/macOS headless | `go-schedule_<ver>_<os>_<arch>.tar.gz` | [All guides](docs/install.md) |

The Windows MSI installs the service and updates PATH; open a new terminal afterward. Linux and macOS service registration is covered in the platform guides. The macOS desktop's first-launch daemon is not itself a boot-persistent service.

## First task

After installation, check the daemon:

```sh
gosched health
```

The health example below targets the prepared v1.7.0 source, which is not yet published. Use the [latest published release](https://github.com/shruggietech/go-schedule/releases/latest) to identify currently available packages.

```text
daemon ok (version 1.7.0)
```

Create a harmless weekday task. On Windows:

```powershell
gosched task add hello --command cmd.exe --arg /c --arg "echo hello" --schedule "weekdays at 09:00"
```

On Linux or macOS:

```sh
gosched task add hello --command /bin/echo --arg hello --schedule "weekdays at 09:00"
```

The response includes the task ID, interpreted schedule and upcoming runs. Substitute that ID below to run it immediately and inspect its history:

```sh
gosched task run-now <task-id>
gosched runs --task <task-id>
```

Open the desktop with `gosched gui` if installed. Continue with the [CLI reference](docs/cli.md), [desktop fields](docs/gui-fields.md) and [cron guide](docs/cron.md) for complete options and workflow examples.

## Documentation and support

- [Settings and storage](docs/options.md), [remote access](docs/remote-access.md) and [architecture](docs/architecture.md)
- [JSON API](docs/api.md), [notifications](docs/notifications.md) and [MCP](docs/mcp.md)
- [Maintainer checks](docs/test-scripts.md) and [website maintenance](docs/website.md)
- [Report a bug or request a feature](https://github.com/shruggietech/go-schedule/issues/new/choose)

Development uses [Spec Kit](https://github.com/github/spec-kit), review branches and pull requests. Start with [CONTRIBUTING.md](CONTRIBUTING.md). The [constitution](.specify/memory/constitution.md) defines engineering requirements; [brand integration](brand/REPOSITORY.md) describes approved assets and updates.

Report security issues privately through [SECURITY.md](SECURITY.md). Participation follows the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

[Apache License 2.0](LICENSE). © 2026 [ShruggieTech](https://shruggie.tech).

<sub>A [ShruggieTech](https://shruggie.tech) project.</sub>
