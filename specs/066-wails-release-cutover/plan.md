# Implementation Plan: Wails Release Cutover

**Branch**: `codex/066-wails-cutover` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/066-wails-release-cutover/spec.md`

## Summary

Promote the completed `desktop/` Wails application into the production release path while retaining the external `gosched-gui` launch name. Package it natively with the existing daemon, CLI, branding, and Windows lifecycle contract; replace Fyne and proof-only CI with production Wails gates; remove the retired Fyne source and root dependencies; and leave one auditable parity, deviation, and exact-candidate qualification record. This completes #157 and makes parent #147 eligible for post-merge closure without tagging or publishing a release.

## Technical Context

**Language/Version**: Go 1.25.0 in the root and desktop modules; TypeScript 5.6.3 and React 19.1.0 in `desktop/frontend`

**Primary Dependencies**: Wails v2.14.0, Node.js 24, Vite 7.3.6, Vitest 5.0.0, Playwright 1.63.0, WiX 6.0.2, existing protected local IPC client

**Storage**: Existing daemon SQLite state plus versioned per-user desktop JSON preferences; no schema change

**Testing**: Root Go tests with race and coverage gates, desktop Go race tests, frontend unit/build/audit, Playwright accessibility and responsive suite, Wails native builds on Windows/macOS/Linux, Windows compiled MSI and silent lifecycle contract, documentation and automation policy checks

**Target Platform**: Windows amd64 MSI, macOS arm64 application archive, Linux amd64 portable desktop archive; existing cgo-free daemon and CLI archives for Linux and macOS amd64/arm64

**Project Type**: Multi-module desktop application with daemon, CLI, release automation, installer, documentation, and immutable historical specifications

**Performance Goals**: Preserve the existing scheduler dispatch budget; desktop startup and interaction remain bounded by the existing two-second local connection operations; no new hot scheduling path

**Constraints**: Stable `gosched-gui` external name, windowless Windows application, native Wails compilation, offline-capable assets, no release tag or public publication, UTF-8 without BOM, no current Fyne implementation residue, no skipped gate reported as passing

**Scale/Scope**: One final v1.2 cutover across three desktop package targets, one Windows upgrade lifecycle, one production Wails module, and the complete #147 workflow inventory

## Constitution Check

*GATE: Passed before research and re-checked after design.*

- **I. Code Quality**: PASS. Retirement removes an obsolete UI implementation and duplicated proof target. Release-contract changes receive focused tests, and generated or shell orchestration remains explicit and reviewable.
- **II. Testing Standards**: PASS. Contract tests precede workflow changes; root and desktop race suites, core coverage, frontend tests, native builds, accessibility, package inspection, and installer lifecycle remain mandatory.
- **III. User Experience Consistency**: PASS. The stable launcher name, daemon authority, CLI behavior, preference migration, platform identity, and documented package entry points remain consistent while the implementation changes.
- **IV. Performance Requirements**: PASS. No scheduler hot path changes. Existing dispatch benchmarks and p99 tests remain in the canonical verifier and hosted CI.
- **V. Autonomous Build-Phase Execution**: PASS. S066 uses the review branch `codex/066-wails-cutover`, runs every spec-kit stage including read-only analyze, and publishes only because the operator explicitly authorized automatic push and PR creation. No release or tag is authorized.
- **Engineering constraints**: PASS. Supported builds remain native on Windows, macOS, and Linux; persistence and IPC contracts do not change; dependency removal reduces supply-chain surface.

## Project Structure

### Documentation (this feature)

```text
specs/066-wails-release-cutover/
├── checklists/
│   ├── release-readiness.md
│   └── requirements.md
├── contracts/
│   └── release-candidate.md
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
.github/
├── ISSUE_TEMPLATE/
└── workflows/
    ├── ci.yml
    └── release.yml
build/windows/
desktop/
├── frontend/
├── settings/
├── wails.json
└── README.md
docs/
scripts/
test/
├── integration/
└── scripts/
CHANGELOG.md
README.md
go.mod
go.sum
```

The retired `gui/`, `cmd/gosched-gui/`, and `experiments/wails-foundation/` trees are removed after contract tests define the replacement boundary. The Wails module remains separate because native desktop dependencies must not contaminate cgo-free daemon and CLI builds.

**Structure Decision**: Keep the two-module architecture selected in S061. Production packaging enters `desktop/` for Wails and the root module for daemon, CLI, cleanup helper, installer contract, and core verification. Use the Wails configuration to emit the compatibility filename `gosched-gui`; do not copy a differently named binary into a legacy alias because one canonical identity is easier to inspect and support.

## Delivery Phases

1. Establish exact-candidate, package-payload, parity, and residue contracts in tests and specification artifacts.
2. Promote Wails into CI, release packaging, Windows MSI staging, and the stable launcher identity.
3. Remove Fyne source, dependencies, proof-only application, and superseded gates after the replacement path is testable.
4. Reconcile user, contributor, issue-template, build, and release guidance; preserve explicitly historical and migration references.
5. Run focused and canonical verification locally, publish the PR, then require the hosted native and packaging matrix on the same head revision.

## Post-Design Constitution Re-check

PASS. The design retains current external contracts, reduces duplicate implementation and dependency surface, adds failure-closed release checks, preserves the full quality gate, and keeps irreversible publication outside the slice. No constitutional deviation is required.

## Complexity Tracking

No constitution violations require justification.
