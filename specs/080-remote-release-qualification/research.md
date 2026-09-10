# Research: Remote Access Release Qualification

## Registered Service Configuration

The current service command registers only the daemon executable, so the optional daemon `--config` flag is lost when the service manager starts it. Remote enablement cannot be made boot-persistent through documented service commands without adding an installation-time configuration binding. The safe minimal design validates the file with the existing loader, resolves the path before registration, and stores no secrets in the command line because certificates and keys remain file references.

## Evidence Composition

Prior slices already provide detailed deterministic tests for TLS, pairing, rate limits, permissions, audit, desktop and CLI clients, version and identity checks, recovery, and secret exclusion. S080 should not copy those suites. One package-shaped integration journey composes the shipped daemon boundary, while a named cross-platform CI matrix reruns it and the relevant focused packages on each supported operating system.

## Service-Manager Boundary

Linux systemd, macOS launchd, and Windows Service Control Manager have materially different privilege and lifecycle behavior on hosted runners. The portable invariant owned by go-schedule is the registered executable and argument list. Existing platform and installer jobs cover native service behavior. S080 tests exact argument preparation in process and exercises the daemon itself as a subprocess so failures remain deterministic and diagnostic.

## Documentation Shape

The existing `docs/remote-access.md` is the architecture and client contract of record. Expanding it with an operator runbook prevents link and terminology drift across a second remote guide. Platform install pages will point to the shared procedure and document the platform-specific service command.
