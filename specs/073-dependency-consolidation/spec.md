# Feature Specification: Dependency Consolidation

**Feature Branch**: `codex/073-dependency-consolidation`

**Created**: 2026-09-09

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Dependency graphs, companion runtime updates, focused compatibility evidence, and canonical verification completed on review branch `codex/073-dependency-consolidation`; hosted checks and pull-request review follow publication

**Input**: User description: "Consolidate every dependency update represented by Dependabot pull requests #201 through #208 into one coherent, fully green replacement for issue #215, including lockfiles, generated artifacts, and narrowly required compatibility fixes without unrelated product or dependency-policy work."

## User Scenarios & Testing

### User Story 1 - Maintain One Coherent Dependency Baseline (Priority: P1)

As a maintainer, I want the eight pending dependency updates combined into one internally consistent baseline so that I can review and merge one trustworthy change instead of reconciling overlapping automated pull requests.

**Why this priority**: The open updates overlap across two module graphs and one frontend tree. Several individual pull requests are red, so merging them independently would leave compatibility and lockfile state uncertain.

**Independent Test**: Start from the current default branch, apply the complete selected version set, install each dependency graph from its committed manifests and lockfiles, and verify that no requested update is absent or silently downgraded.

**Acceptance Scenarios**:

1. **Given** pull requests #201 through #208 are open, **When** the consolidated baseline is inspected, **Then** every represented direct dependency is present at the requested version or a documented newer compatible version.
2. **Given** a requested update introduces a peer or runtime conflict, **When** the baseline is resolved, **Then** the smallest required companion update is included and its rationale is recorded.
3. **Given** the root module, desktop module, and frontend tree share transitive constraints, **When** their manifests are regenerated, **Then** every graph installs reproducibly without requesting uncommitted changes.

---

### User Story 2 - Preserve Supported Product Behavior (Priority: P1)

As a user, I want dependency maintenance to preserve scheduling, storage, notifications, agent access, desktop interaction, and packaging behavior across supported platforms.

**Why this priority**: Dependency updates have no delivery value if they regress the product, and the affected libraries sit beneath storage, filesystem watching, desktop runtime, frontend rendering, and test infrastructure.

**Independent Test**: Run the repository's complete local verification contract and the pull request's Windows, macOS, and Linux checks from clean dependency installations.

**Acceptance Scenarios**:

1. **Given** the consolidated dependency baseline, **When** all local quality gates run, **Then** format, vet, lint, race, GUI, coverage, documentation, and automation complete successfully.
2. **Given** clean supported hosted environments, **When** the final pull request head is tested, **Then** every required build, test, security, accessibility, and packaging check is green on that exact head.
3. **Given** an update requires a compatibility adjustment, **When** the adjustment is reviewed, **Then** it is bounded to restoring existing supported behavior and introduces no unrelated product capability.

---

### User Story 3 - Retain Review and Supersession Traceability (Priority: P2)

As a reviewer, I want the replacement pull request to identify every source update, version decision, compatibility fix, and validation result so that automated pull requests can be closed only after the replacement is merged.

**Why this priority**: One consolidated pull request must remain auditable and must not erase the provenance or completion status of the work it supersedes.

**Independent Test**: Review the final specification, change set, pull request description, check results, and linked records, then account for each pull request from #201 through #208 and issue #215 without closing any source pull request before merge.

**Acceptance Scenarios**:

1. **Given** the replacement pull request is open, **When** its description is reviewed, **Then** it links #201 through #208, closes #215 on merge, lists selected versions, and explains every companion change or omission.
2. **Given** review feedback or CI changes the branch, **When** the final head is published, **Then** all evidence is rerun or refreshed against that exact head.
3. **Given** the replacement has not merged, **When** repository planning records are inspected, **Then** #201 through #208 and #215 remain open.

### Edge Cases

- A package publishes a newer version after issue #215 but before implementation. The slice may select it only when it is compatible, remains within the represented update line, and is documented.
- A requested package requires a newer runtime or peer dependency. The smallest supported baseline change is included instead of bypassing or ignoring the constraint.
- A lockfile can be generated locally but not installed reproducibly in a clean environment. The lockfile is invalid until clean installation succeeds.
- A module tidy operation changes transitive requirements in the other Go module. Both module graphs must be reconciled and committed together.
- A frontend major update passes unit tests but breaks production bundling, native packaging, browser accessibility, or responsive behavior. The update remains incomplete until all affected gates pass.
- A dependency update reveals an upstream defect without a safe compatibility fix. The exact update may be omitted only with evidence, a bounded technical rationale, and explicit pull request disclosure.
- A source Dependabot branch is updated, closed, or superseded while S073 is active. The final audit uses the latest authoritative state without broadening beyond the eight named updates.

## Requirements

### Functional Requirements

- **FR-001**: The consolidated baseline MUST account for every direct dependency update represented by pull requests #201 through #208.
- **FR-002**: Each represented dependency MUST use the requested version or a documented newer compatible version; any omission MUST include technical evidence and a reviewer-visible rationale.
- **FR-003**: The root module, desktop module, and frontend dependency tree MUST each have internally consistent committed manifests and integrity files.
- **FR-004**: Clean dependency restoration MUST complete without force flags, legacy peer resolution, ignored engine constraints, or uncommitted manifest changes.
- **FR-005**: Runtime and type-definition baselines MUST describe the same supported major version wherever a selected package makes that alignment necessary.
- **FR-006**: Companion dependency updates MUST be limited to versions required to satisfy the represented updates' declared compatibility constraints.
- **FR-007**: Compatibility code or configuration changes MUST be limited to preserving behavior supported on the default branch.
- **FR-008**: The consolidated baseline MUST preserve all existing safety-critical scheduling, storage migration, restart recovery, concurrency, and local access-control verification.
- **FR-009**: Desktop unit, production bundle, native application, browser accessibility, responsive, and Windows installer contracts MUST pass without weakening their assertions.
- **FR-010**: All repository quality gates and hosted pull request checks MUST pass on the final reviewed head.
- **FR-011**: The replacement pull request MUST link #201 through #208, use `Closes #215`, enumerate selected versions, and record compatibility decisions and verification evidence.
- **FR-012**: Pull requests #201 through #208 and issue #215 MUST remain open until the replacement pull request merges; post-merge supersession is housekeeping outside implementation tasks.
- **FR-013**: No unrelated product feature, dependency-policy change, release tag, release publication, or broad toolchain modernization may be included.
- **FR-014**: Every human, Codex, and security-bot review comment MUST receive an evidence-based disposition, and warranted fixes MUST be revalidated before publication.

### Key Entities

- **Source Update**: One of pull requests #201 through #208, identified by dependency, requested version, dependency graph, update class, and final disposition.
- **Consolidated Baseline**: The exact direct and transitive dependency versions committed across the root module, desktop module, frontend manifest, and integrity files.
- **Compatibility Adjustment**: A narrowly required runtime, peer dependency, generated artifact, source, test, build, or packaging change with a recorded cause and validation result.
- **Verification Evidence**: A quality-gate or hosted-check result tied to the exact commit that was tested.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All eight source updates are represented in one replacement pull request with zero unexplained omissions.
- **SC-002**: All three dependency graphs restore cleanly and leave zero uncommitted manifest or integrity-file changes.
- **SC-003**: All eight canonical local quality gates complete successfully on the implemented branch.
- **SC-004**: One hundred percent of required hosted checks, including Windows, macOS, Linux, security, accessibility, and packaging checks, pass on the final reviewed head.
- **SC-005**: Every review comment receives an explicit disposition and zero actionable review conversations remain unresolved.
- **SC-006**: The final change set contains zero unrelated product capabilities, release operations, or weakened tests.
- **SC-007**: Until merge, all eight source pull requests and issue #215 retain their open state; after merge, the repository has one authoritative merged dependency update and eight traceably superseded automated pull requests.

## Assumptions

- Issue #215 and pull requests #201 through #208 are the authoritative scope boundary.
- A newer compatible patch or minor version may replace the exact Dependabot proposal when it remains within the same intended update and reduces immediate staleness.
- A package's declared runtime or peer-dependency range is a compatibility requirement, not a warning to bypass.
- The repository's supported platform and quality-gate definitions remain unchanged unless a represented dependency requires a narrowly justified runtime-major alignment.
- Source pull-request closure happens only after the maintainer merges S073 and is therefore excluded from implementation task completion.
