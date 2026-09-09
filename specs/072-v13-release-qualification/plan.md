# Implementation Plan: v1.3 Notifications and Local Agent Access Qualification

**Branch**: `codex/072-v13-release-qualification` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/072-v13-release-qualification/spec.md`

## Summary

Qualify the completed v1.3 notification and local Observe capabilities as one release boundary for issue #190. Add a named Windows, macOS, and Linux pull-request matrix that starts a package-shaped daemon from fresh and retained state, proves notification and MCP HTTP defaults stay opt-in, and reruns the existing race-checked webhook, MCP SDK, authorization, redaction, recovery, and migration evidence. Protect the workflow contract with repository automation checks, record exact verification evidence, and do not tag or publish a release.

## Technical Context

**Language/Version**: Go 1.25.0; GitHub Actions YAML; POSIX shell for repository automation checks

**Primary Dependencies**: Go standard library, `github.com/modelcontextprotocol/go-sdk` v1.7.0, existing local IPC client and daemon packages; no new dependencies

**Storage**: Existing SQLite scheduler state in isolated temporary data directories; no schema changes

**Testing**: Go unit and integration tests under `go test -race`; package-shaped subprocess tests; repository automation contract; canonical `sh scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux hosted runners plus the local Windows workstation

**Project Type**: Cross-platform daemon, CLI, desktop client, and release automation repository

**Performance Goals**: Preserve the existing p99 dispatch target below 100 ms; the qualification harness adds no production hot-path work

**Constraints**: Optional outbound notifications and localhost MCP remain disabled by default; local scheduling remains offline; protected values remain redacted; no release tag or artifact publication; all child console processes remain hidden on Windows

**Scale/Scope**: One release issue (#190), three supported operating systems, two daemon starts per platform, the complete webhook qualification surface, and the complete Observe-only MCP transport and security surface

## Constitution Check

*GATE: Passed before research and re-checked after design.*

- **I. Code Quality**: The implementation adds a focused integration harness and CI contract without changing production architecture or dependencies. Process lifecycle and failure diagnostics are explicit.
- **II. Testing Standards**: Qualification is test-first, race-checked, and reuses the detailed notification, persistence, SDK, authorization, and hostile-content suites instead of replacing them with weaker smoke tests.
- **III. User Experience Consistency**: Fresh and retained state are inspected through the supported local API. Optional surfaces remain explicit and release-facing claims remain bounded to shipped behavior.
- **IV. Performance Requirements**: No production path changes. Canonical verification retains the existing dispatch-latency and benchmark evidence.
- **V. Autonomous Build-Phase Execution**: The slice is traceable to open issue #190, follows the complete Spec Kit sequence, uses a review branch and pull request, and relies on the user's explicit instruction as publication authorization.
- **Engineering Constraints**: The three-platform matrix exercises supported targets, storage remains forward-compatible, and security checks cover redaction and denied access.
- **Post-design re-check**: Passed. The design introduces no new runtime interface, persistence state, dependency, or constitutional deviation.

## Project Structure

### Documentation (this feature)

```text
specs/072-v13-release-qualification/
├── checklists/
│   ├── release-security.md
│   └── requirements.md
├── contracts/
│   └── qualification-gate.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
.github/workflows/ci.yml
scripts/automation-check.sh
test/integration/v13_release_qualification_test.go
docs/notifications.md
docs/mcp.md
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Keep qualification at the existing repository integration boundary. The new test owns the package-shaped fresh and retained daemon journey, the existing focused package tests retain detailed behavioral authority, and CI composes both into a visible three-platform gate. No production package or new abstraction is warranted.

## Complexity Tracking

No constitutional violations or complexity exceptions are required.
