# Implementation Plan: Cumulative v1.4.0 Release Preparation

**Branch**: `codex/081-v140-release-preparation` | **Date**: 2026-09-10 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/081-v140-release-preparation/spec.md` and release coordinator [#226](https://github.com/shruggietech/go-schedule/issues/226).

## Summary

Cut one cumulative v1.4.0 source boundary after public v1.1.1, add concise tag-specific release copy, make release preflight validate every source-owned identity surface before artifact mutation, remove stale v1.2-candidate and Fyne-era release wording, and define the exact post-merge fresh-install, upgrade, staging, qualification, promotion, and audit sequence without creating a tag or release.

## Technical Context

**Language/Version**: GitHub Actions YAML, POSIX shell, Go 1.25 release validators, PowerShell Windows qualification, and GitHub-flavored Markdown

**Primary Dependencies**: GitHub Actions, GitHub Releases, WiX 6.0.2, existing release-gate tooling, and the repository GitHub-format validator

**Storage**: Git objects, release metadata, draft assets, candidate manifests, Windows evidence archives, and project planning records

**Testing**: Release automation fixture mutations, release-copy validation, documentation and lifecycle checks, `scripts/verify.sh all`, hosted CI, candidate verification, evidence validation, and checksum audit

**Target Platform**: GitHub-hosted Linux, macOS, and Windows builders plus attended Windows 11 fresh and upgrade environments

**Project Type**: Cross-platform daemon, CLI, Wails desktop, and release operation

**Performance Goals**: Not applicable; release correctness, provenance, and fail-closed behavior dominate elapsed time

**Constraints**: No tag or release in S081, no retroactive v1.2.0 or v1.3.0 publication, exact-commit CI, draft-only staging, no artifact rebuild, no false public-release claims, no hard-wrapped GitHub prose, and no Unicode em dash

**Scale/Scope**: One cumulative version boundary, four highlights, seven distributable packages, one candidate manifest, one attended evidence archive, one checksum inventory, two Windows installation paths, and 47 established attended observations

## Constitution Check

- **I. Code Quality**: PASS. Release-preflight changes stay deterministic and fixture-tested; no product dependency or runtime branch is introduced.
- **II. Testing Standards**: PASS. Automation changes begin with failing fixture mutations, all canonical gates remain mandatory, and physical Windows assertions remain attended rather than simulated.
- **III. User Experience Consistency**: PASS. Public copy, install guides, version examples, and artifact names describe one consistent release boundary.
- **IV. Performance Requirements**: PASS. No scheduler or hot-path behavior changes.
- **V. Autonomous Build-Phase Execution**: PASS. S081 uses the complete spec-kit sequence, review branch, official pull request, hosted CI, third-party review, and final maintainer merge gate. The user's explicit instruction authorizes branch push and PR creation but not a tag or release.

**Gate result before and after design**: PASS. No principle deviation or unjustified complexity is required.

## Project Structure

```text
.github/release-notes/v1.4.0.md
.github/workflows/release.yml
.specify/feature.json
CHANGELOG.md
CLAUDE.md
README.md
docs/INSTALL-linux.md
docs/INSTALL-macos.md
docs/INSTALL-windows.md
docs/api.md
docs/architecture.md
docs/install.md
internal/releasegate/validate.go
internal/releasegate/validate_test.go
specs/081-v140-release-preparation/
├── checklists/
│   ├── release.md
│   └── requirements.md
├── contracts/publication.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
specs/README.md
test/scripts/automation-check_test.sh
test/integration/windows_release_gate_contract_test.go
test/windows/Invoke-ReleaseCandidateAttended.ps1
test/windows/README.md
```

**Structure Decision**: Extend the established release source surfaces and automation fixture rather than introduce a second release tool. The release workflow gains a metadata preflight in its existing gate. The attended collector emits generic version-two native-window evidence from Win32 measurements, while the validator retains version-one Fyne attachment compatibility for historical archives. Spec Kit artifacts hold the post-merge operational contract and evidence boundary.

## Complexity Tracking

No constitution violation requires justification.
