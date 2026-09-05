# Feature Specification: Maintainer Automation Baseline

**Feature Branch**: `011-maintainer-automation-baseline` (trunk-based: committed directly to `main`)

**Created**: 2026-08-26

**Status**: Implemented

**Delivery**: [PR #43](https://github.com/shruggietech/go-schedule/pull/43)

**Input**: User description: "Combine GitHub issues #21 and #41 into one focused Spec-Kit work slice that can run end-to-end under the build-phase autopilot protocol: move all repository workflows off deprecated Node 20 actions and converge the Makefile onto one canonical, non-mutating local CI-parity verification path."

**Traceability**: Closes GitHub issues [#21](https://github.com/shruggietech/go-schedule/issues/21) and [#41](https://github.com/shruggietech/go-schedule/issues/41).

## Overview

The repository currently has two maintenance hazards that share one surface: the automation used to decide whether a change is safe.

First, every first-party action used by the CI and release workflows targets a deprecated runtime. Hosted runners temporarily translate those actions onto a supported runtime, but that compatibility bridge is scheduled to disappear. When it does, otherwise-correct workflows will fail before project commands can start.

Second, the Makefile offers commands whose names imply project verification but whose behavior differs from the gates documented for contributors and enforced by CI. A maintainer can therefore receive a local green result without running the same lint version, race-test selection, coverage threshold, or documentation check that protects `main`.

This feature treats CI and local verification as one maintainer-facing contract: the hosted workflows remain runnable on supported action runtimes, and one documented, non-mutating local entry point exercises the complete repository gate using the same underlying commands as CI.

### Scope in

- All action references in the CI and release workflows that currently target the deprecated Node 20 action runtime.
- The Makefile targets used for formatting checks, vetting, linting, race tests, GUI tests, coverage, documentation validation, and aggregate verification.
- Contributor and autopilot documentation that names the canonical local gate.
- Regression checks that detect action-runtime or verification-command drift.
- Dated changelog decisions for every pinned artifact changed by the feature.

### Scope out

- Product, scheduler, daemon, CLI, GUI, store, and IPC behavior.
- Go module version upgrades or automated dependency-update proposals (#40).
- Repository security settings, CodeQL, secret scanning, or private vulnerability reporting (#38 and #39).
- Spec-Kit lifecycle cleanup (#42).
- Changing the trunk-based integration model or requiring pull requests (#23).
- Cutting a release, pushing to `main`, or running a release workflow.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Hosted automation starts on a supported runtime (Priority: P1)

As a maintainer, I need every action used by CI and release automation to target a supported hosted-runner runtime so that validation and packaging do not fail before project commands begin.

**Why this priority**: The current compatibility shim is temporary. Its removal would disable the repository's validation and release paths independently of the quality of the project itself.

**Independent Test**: Audit every action reference in both workflow files, resolve each referenced release's declared runtime, and run the ordinary CI workflow. The story is complete when no reference targets Node 20 and the CI run contains no Node 20 deprecation annotation.

**Acceptance Scenarios**:

1. **Given** the CI workflow, **when** each action reference is resolved to its declared runtime, **then** every action uses a supported runtime and none uses Node 20.
2. **Given** the release workflow, **when** each action reference is resolved to its declared runtime without executing a release, **then** every action uses a supported runtime and none uses Node 20.
3. **Given** an ordinary validation run, **when** workflow setup completes, **then** no deprecation warning identifies an action in this repository as a Node 20 action.
4. **Given** the action-major upgrades, **when** their inputs and outputs are compared with existing workflow usage, **then** checkout, toolchain caching, benchmark upload, release publication, and artifact attachment retain their existing contracts.

---

### User Story 2 - Maintainers have one definition of green (Priority: P2)

As a maintainer or autopilot agent, I need one documented local verification entry point that does not edit the working tree and exercises the same required checks as CI, so that a local green result means the change is ready for the single pre-push review halt.

**Why this priority**: The current Makefile drift can hide failures locally, but CI remains available today. Converging the contract removes that ambiguity and prevents the two paths from drifting again.

**Independent Test**: Start from a clean checkout, invoke the documented aggregate verification entry point, and record every child check it runs. The story is complete when the command leaves the checkout unchanged, invokes the project-pinned linter, uses the canonical race selection and coverage gate, validates documentation and the automation contract, and returns non-zero when any child check fails.

**Acceptance Scenarios**:

1. **Given** a clean, valid checkout, **when** the canonical local verification entry point completes, **then** every required local CI-parity check passes and `git status --short` remains empty.
2. **Given** an unformatted Go file, **when** verification runs, **then** it reports the file and fails without rewriting it.
3. **Given** a failure in any format, vet, lint, race, GUI, coverage, documentation, or automation-contract gate, **when** aggregate verification runs, **then** it stops with a non-zero result that identifies the failing gate.
4. **Given** a maintainer invokes an individual Make target, **when** that target corresponds to a CI gate, **then** it delegates to the same command or script used by aggregate verification rather than maintaining a second definition.
5. **Given** a target that formats or otherwise mutates files, **when** a maintainer reads the Makefile help or contributor guidance, **then** that mutating behavior is explicitly distinguished from verification targets.

### Edge Cases

- A host without a C toolchain cannot run the race gate. Verification must fail clearly and preserve the documented rule that an unrun gate is never reported as passing.
- A Windows host may not expose POSIX `sh` through `PATH`. The canonical entry point must either perform an actionable prerequisite check or route the shell checks through an already-supported project mechanism; it must not silently omit coverage or documentation validation.
- The base Go installation may be older than the module's selected toolchain. The pinned linter command must retain the documented toolchain override path instead of encouraging changes to the linter configuration or module version.
- An action can publish a new major version without preserving an input or output used by this repository. Runtime support alone is insufficient; the workflow contract must also be checked.
- The release workflow cannot be exercised end-to-end without creating a tag and publishing artifacts. Validation must prove its action references and syntax without triggering a release.
- A new required CI job can be added later. Drift protection must make the corresponding local-gate decision explicit rather than allowing silent divergence.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST inventory every external action reference in both the validation and release workflows.
- **FR-002**: Every inventoried action MUST target a release whose declared execution runtime is supported by the hosted automation platform and is not Node 20.
- **FR-003**: Action upgrades MUST preserve every existing workflow capability, input, output, permission boundary, trigger, matrix, artifact name, and release file selection unless a separately recorded correction is required.
- **FR-004**: The repository MUST provide exactly one documented aggregate local entry point for the complete pre-push verification contract.
- **FR-005**: Aggregate verification MUST be non-mutating and MUST include format checking, vetting, the project-pinned linter, the canonical race-test package selection, headless GUI tests, the enforced core-package coverage thresholds, documentation integrity validation, and automation-contract drift checking.
- **FR-006**: Aggregate verification MUST return a non-zero result when any required gate fails and MUST identify the failed gate.
- **FR-007**: Individual verification targets MUST share their commands or scripts with aggregate verification so package selections, tool versions, and thresholds are not independently duplicated.
- **FR-008**: Mutating convenience targets MUST be named or documented so they cannot be mistaken for verification gates.
- **FR-009**: The canonical verification contract MUST preserve the project's rule that unavailable prerequisites are reported as unrun failures, never as passing checks.
- **FR-010**: Contributor guidance, the autopilot protocol, and the Makefile MUST name the same canonical local verification entry point and the same required gate set.
- **FR-011**: The repository MUST contain an automated drift check that fails when a workflow reintroduces a deprecated action runtime reference or when the documented aggregate verification entry point omits a required gate.
- **FR-012**: Changes to pinned workflows and the Makefile MUST be accompanied by dated changelog decisions explaining the runtime upgrades and the selected single-source verification design.
- **FR-013**: Validation of the release workflow MUST NOT create a tag, publish a release, upload release assets, or otherwise perform an external release-side effect.

### Key Entities

- **Action Reference**: A workflow's external action owner, repository, selected release, declared execution runtime, used inputs/outputs, and workflow locations.
- **Verification Gate**: A named, non-mutating check with one canonical command, prerequisites, success condition, failure behavior, and CI/local consumers.
- **Verification Contract**: The complete ordered set of required gates that defines a trustworthy local pre-push result.
- **Pinned-Artifact Decision**: A dated record naming the pinned file, the change, the reason, and any compatibility constraint.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of action references in the validation and release workflows resolve to supported, non-Node-20 execution runtimes.
- **SC-002**: An ordinary CI run completes without any Node 20 action deprecation annotation attributable to this repository.
- **SC-003**: One documented local command invokes 100% of the required verification gates and leaves a clean checkout unchanged.
- **SC-004**: A controlled child-gate failure produces a non-zero aggregate result naming that gate; no required failure is swallowed.
- **SC-005**: A mechanical drift check detects both a deliberately reintroduced deprecated action reference and a deliberately omitted local verification gate.
- **SC-006**: All changed pinned artifacts have a dated decision in the unreleased changelog before the feature is committed.
- **SC-007**: The feature completes without changes to product behavior, release publication, repository settings, or the trunk-based integration model.

## Assumptions

- GitHub-hosted Actions remains the repository automation platform.
- The existing CI job behaviors, release triggers, permissions, matrices, and artifact contracts are correct; this feature preserves rather than redesigns them.
- The required local gate set is the union already enforced or documented by the project: format check, vet, pinned lint, race-selected tests, GUI tests, core coverage thresholds, and documentation integrity.
- Maintainers may need Git Bash or an equivalent POSIX shell for existing shell gates on Windows; this feature must make that dependency actionable but does not replace the existing shell scripts solely to remove it.
- Release workflow validation is static plus contract-focused because tags and releases remain explicit operator-authorized actions.
- No new third-party runtime dependency is required to implement the feature.
