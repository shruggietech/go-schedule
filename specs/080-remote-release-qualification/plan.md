# Implementation Plan: Remote Access Release Qualification

**Branch**: `codex/080-remote-release-qualification` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/080-remote-release-qualification/spec.md`

## Summary

Close the v1.4 implementation milestone with one reproducible source qualification. Add optional validated configuration binding to service installation, exercise package-shaped remote access across Windows, macOS, and Linux, publish operator-focused deployment and recovery guidance, and map issue #173 to exact local and hosted evidence. No tag or release is created.

## Technical Context

**Language/Version**: Go 1.25.0 with toolchain 1.25.1; GitHub Actions YAML; GitHub-flavored Markdown

**Primary Dependencies**: Go standard library; existing Cobra CLI, kardianos service adapter, remote HTTPS server, enrollment service, typed client, and OpenAPI contract

**Storage**: Existing JSON daemon configuration, SQLite daemon state, and user-scoped secret-free profiles

**Testing**: Table-driven CLI tests, package-shaped subprocess integration, existing focused security/race suites, three-platform hosted matrix, `scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux daemon/CLI; current Wails desktop clients

**Project Type**: Cross-platform daemon, CLI, desktop, release workflows, and documentation site

**Performance Goals**: Preserve existing bounded startup, request, retry, and shutdown behavior; no new hot path

**Constraints**: Remote access remains opt-in, TLS 1.3 only, local IPC unchanged, no secrets in configuration arguments or evidence, no release publication

**Scale/Scope**: One daemon per host, one configured listener, existing bounded client and request limits

## Constitution Check

- **I. Code Quality**: Keep service argument preparation isolated, deterministic, and directly tested; add no dependency.
- **II. Testing Standards**: Write failing CLI and lifecycle tests first, retain race coverage, and run all eight gates.
- **III. User Experience Consistency**: Validate configuration before service registration and return actionable field and path errors.
- **IV. Performance Requirements**: No scheduling or hot-path change; existing timing budgets remain unchanged.
- **V. Autonomous Build-Phase Execution**: Use the required spec-kit sequence, review branch, PR, hosted review, and final maintainer merge gate.

**Gate result**: PASS. No principle deviation or unjustified complexity is required.

## Project Structure

### Documentation

```text
specs/080-remote-release-qualification/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── tasks.md
├── verification.md
├── contracts/
│   └── qualification.md
└── checklists/
    ├── release.md
    └── requirements.md
```

### Source Code

```text
internal/cli/service.go
internal/cli/service_test.go
test/integration/v14_remote_release_qualification_test.go
.github/workflows/ci.yml
docs/remote-access.md
docs/cli.md
docs/INSTALL-linux.md
docs/INSTALL-macos.md
docs/INSTALL-windows.md
CHANGELOG.md
```

**Structure Decision**: Extend the existing service command, cross-platform integration package, CI matrix, and canonical remote guide. A separate release-gate executable would duplicate contracts already enforced by Go tests and the documentation checker.

## Decision Log

- **Service configuration binding**: A flag-free daemon checks the platform data directory for optional `config.json`, while `service install --config` requires an existing valid file, converts it to an absolute path, and registers `--config <absolute-path>`. The standard path survives Windows Installer service replacement, custom paths avoid service working-directory ambiguity, and an absent standard file preserves safe built-in defaults.
- **Qualification level**: Use package-shaped daemon binaries plus real IPC, TLS, enrollment, authorization, audit, revocation, restart, and binary replacement on all three hosted systems. System-service-manager mutation stays in focused command tests because GitHub service-manager behavior is already exercised by platform-specific packaging jobs and is not uniform across disposable runners.
- **Deployment posture**: Recommend private networks and SSH tunnels. Reverse proxy and direct public HTTPS remain advanced supported modes with operator-owned certificate, DNS, firewall, monitoring, and proxy duties.
- **Publication boundary**: Completing #173 qualifies reviewed source and closes the implementation milestone; producing v1.4.0 artifacts remains a separately authorized release slice.

## Phase Plan

1. Establish specification, clarification decisions, release checklist, plan, contracts, and tasks.
2. Add failing service configuration and package-shaped remote lifecycle tests.
3. Implement service configuration argument preparation and cross-platform qualification workflow.
4. Complete supported deployment, recovery, upgrade, and client guidance.
5. Run analysis, focused checks, canonical verification, and record source-bound evidence.
