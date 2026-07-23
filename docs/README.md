# go-schedule documentation

**Audience:** users and maintainers of go-schedule\
**Applies to:** go-schedule 0.6.0 and later\
**Project home:** [README.md](../README.md)

An index of what is here and who each document is for. If you are installing,
start with the guide for your platform; if you are contributing, start with
[CONTRIBUTING.md](../CONTRIBUTING.md) in the repository root.

## Using go-schedule

| Document | What it covers |
| --- | --- |
| [Windows install guide](INSTALL-windows.md) | The `.msi`, the service, `PATH`, upgrading, uninstalling. |
| [Linux install guide](INSTALL-linux.md) | Release archive, systemd registration, data paths. |
| [macOS install guide](INSTALL-macos.md) | Desktop bundle versus headless, launchd, and the boot-persistence caveat that catches people out. |
| [`gosched` command reference](cli.md) | Every command and flag the CLI exposes, with exit codes and elevation requirements. |
| [GUI field reference](gui-fields.md) | What every field in the desktop task editor accepts and means. |

## Verifying an installation

| Document | What it covers |
| --- | --- |
| [Maintainer test scripts](test-scripts.md) | Proving that a real, installed daemon fires on time, survives restarts, catches up after downtime, and honors overlap policies — with recorded evidence rather than a hopeful glance at a log. |

These are written for maintainers but are equally useful to anyone who wants
their install demonstrated rather than assumed. Their output makes an excellent
attachment to a bug report.

## Developing go-schedule

| Document | What it covers |
| --- | --- |
| [Contributing](../CONTRIBUTING.md) | Trunk-based workflow, the spec-kit requirement, the CI-parity gates, pinned artifacts. |
| [Build-phase autopilot](build-autopilot.md) | The protocol under which features are executed end to end with a single human halt. Maintainer-facing. |
| [Constitution](../.specify/memory/constitution.md) | The engineering principles the project is governed by. |
| [Master specification](../specs/001-task-scheduler/spec.md) | What the scheduler is and why, with the plan, data model, and contracts alongside it. |
| [Changelog](../CHANGELOG.md) | Release history, including the dated decisions behind changes to pinned artifacts. |

## Reporting problems

Open an [issue](https://github.com/shruggietech/go-schedule/issues/new/choose).
The forms ask for version, component, install method, OS, and whether you were
elevated — each of those has, at some point on this project, been the fact that
decided the diagnosis.

For security issues, please use the private route described in
[SECURITY.md](../SECURITY.md) rather than a public issue.
