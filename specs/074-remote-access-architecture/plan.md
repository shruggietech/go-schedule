# Implementation Plan: Remote Access Architecture

**Branch**: `codex/074-remote-access-architecture` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/074-remote-access-architecture/spec.md`

## Summary

Resolve issue #165 with one reviewed, durable architecture boundary before any remote listener is implemented. Preserve the current local IPC server, define a separate opt-in HTTPS allowlist around shared daemon operations, freeze the minimum capability, credential, deployment, compatibility, failure, threat, dependency, and downstream sequencing decisions, then protect those decisions with a fixture-backed repository check in the canonical documentation gate.

## Technical Context

**Language/Version**: Markdown, POSIX shell, Go 1.27.1 publication formatter

**Primary Dependencies**: No new runtime dependency; standard shell utilities for the architecture check

**Storage**: No schema change; conceptual future entities only

**Testing**: Positive architecture contract check, adversarial shell fixtures, docs integrity gate, specification lifecycle gate, canonical `scripts/verify.sh all`

**Target Platform**: Repository documentation and verification on Linux, macOS, and Windows-hosted Git Bash

**Project Type**: Architecture and governance increment for an existing daemon, CLI, and Wails desktop application

**Performance Goals**: The offline contract check completes in seconds and adds no runtime work

**Constraints**: No network listener, runtime dependency, endpoint implementation, schema migration, issue closure, release, or public behavior claim; UTF-8 without BOM; GitHub prose remains unwrapped and contains no Unicode em dash

**Scale/Scope**: Issue #165 only, one maintained architecture page, one executable check, one fixture suite, and one complete spec-kit package

## Constitution Check

### Pre-research gate

- **I. Code Quality**: PASS. The only executable production-tree addition is a small POSIX shell policy check with explicit failures and deterministic fixture coverage.
- **II. Testing Standards**: PASS. The contract check and negative fixtures are authored before the maintained architecture page, and canonical verification remains mandatory.
- **III. User Experience Consistency**: PASS. The architecture preserves local defaults and defines target, error, time, and compatibility behavior for future remote clients.
- **IV. Performance Requirements**: PASS. S074 adds no runtime code; future request and stream bounds are explicit review requirements.
- **V. Autonomous Build-Phase Execution**: PASS. The work traces to #165, follows the complete spec-kit sequence, uses a review branch and pull request, and the user has explicitly authorized publication.
- **Engineering constraints**: PASS. Standard library is preferred; every future third-party component has a bounded owner and review rule; Linux and Windows remain required; secrets are excluded from logs and durable plaintext.

### Post-design gate

- **I. Code Quality**: PASS. The remote adapter shares daemon operations without exporting the mixed local route mux.
- **II. Testing Standards**: PASS. Threats, deployment modes, operation records, drift detection, and adversarial mutations all have specified evidence.
- **III. User Experience Consistency**: PASS. Remote access is explicit, target identity precedes action, stable errors remain JSON, and uncertain mutations are not replayed.
- **IV. Performance Requirements**: PASS. Two-stage rate limits, request caps, bounded streams, timeouts, and ephemeral limiter pruning prevent unbounded network work.
- **V. Autonomous Build-Phase Execution**: PASS. #166 is intentionally excluded because #165 requires this architecture to be reviewed first.

## Project Structure

### Documentation

```text
docs/
├── api.md
├── architecture.md
└── remote-access.md

specs/074-remote-access-architecture/
├── checklists/
│   ├── remote-security.md
│   └── requirements.md
├── contracts/
│   └── remote-boundary.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Verification

```text
scripts/
├── docs-check.sh
└── remote-architecture-check.sh

test/scripts/
└── remote-architecture-check_test.sh
```

**Structure Decision**: Keep product-facing architecture in `docs/remote-access.md`, implementation provenance in the S074 spec directory, and the regression contract beside existing shell-based repository policy checks. Integrate the new check through `docs-check.sh` so the canonical docs gate owns architecture drift.

## Implementation Strategy

1. Write the remote architecture checker and adversarial fixtures first; demonstrate that the missing maintained page fails.
2. Add the maintained architecture page and link it from the current API and architecture pages.
3. Record the architecture decision in `CHANGELOG.md` and add S074 to the lifecycle inventory.
4. Run focused contract, docs, lifecycle, and publication-format checks.
5. Run canonical verification and record evidence.
6. Publish the verified branch and pull request, process at most two review rounds, and leave final merge authority to the maintainer.

## Complexity Tracking

No constitutional violation or justified complexity exception is present.
