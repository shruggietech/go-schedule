# Feature Specification: Post-release dependency refresh

**Feature Branch**: `codex/090-post-release-dependency-refresh`  
**Created**: 2026-09-17  
**Status**: In Progress
**Input**: Work slice S090 and [#243](https://github.com/shruggietech/go-schedule/issues/243), consolidating Dependabot pull requests #240, #241, and #242.

## User Scenarios and Testing

### User Story 1 - Maintain one supported dependency baseline (Priority: P1)

Maintainers receive one coherent dependency change based on current `main` instead of three overlapping automated branches.

**Independent Test**: Restore every dependency graph from its manifest and integrity file, then confirm the selected direct versions and clean graph state.

**Acceptance Scenarios**:

1. **Given** the current root, desktop, and frontend dependency graphs, **When** S090 applies the selected updates, **Then** all manifests and integrity files resolve without force, legacy-peer, or manual checksum edits.
2. **Given** the selected Go modules require Go 1.26, **When** the baseline is updated, **Then** both Go modules, automation, release tooling, and contributor guidance agree on Go 1.26.

### User Story 2 - Preserve released behavior and security boundaries (Priority: P1)

Users retain the released daemon, CLI, MCP, remote-access, desktop, accessibility, Windows packaging, and scheduling behavior after the dependency refresh.

**Independent Test**: Run focused affected-surface tests followed by all eight canonical verification gates without weakening assertions.

**Acceptance Scenarios**:

1. **Given** the upgraded MCP SDK and Go security/platform modules, **When** the race and security suites run, **Then** stdio, localhost, remote pairing, authorization, audit, and Windows lifecycle behavior remains passing.
2. **Given** React 19.3 and Vite 8.3, **When** the frontend tests, type check, production build, browser accessibility checks, and native desktop build run, **Then** the released interaction and presentation contracts remain passing.

### User Story 3 - Retain supersession and review traceability (Priority: P2)

Maintainers can audit exactly how each automated proposal was represented and why its source pull request no longer needs independent integration.

**Independent Test**: Compare #240, #241, and #242 with the S090 manifests, verification record, replacement pull request, and source-pull-request disposition comments.

**Acceptance Scenarios**:

1. **Given** the replacement pull request is published, **When** the automated source pull requests are closed, **Then** each receives a comment linking to the complete replacement and no issue is closed prematurely.
2. **Given** CI and external reviews operate on the published branch, **When** findings arrive, **Then** each finding receives an explicit response and any warranted fix is reverified on the latest head with no more than two requested Codex rounds.

### Edge Cases

- An upstream module forces a newer Go language baseline than the repository currently declares.
- Root and desktop Go graphs resolve different transitive versions because the desktop replaces the root module locally.
- A clean npm install resolves a companion package not visible in the direct-version proposal.
- A source Dependabot branch is stale relative to post-release `main` and therefore has unrelated CI failures.
- A review comment arrives after an earlier green check run; the final disposition must refer to the latest head.

## Requirements

### Functional Requirements

- **FR-001**: Consolidate #240, #241, and #242 from current `main` into one official S090 pull request tracked by #243.
- **FR-002**: Select go-sdk 1.8.0, x/crypto 0.57.0, x/sys 0.48.0, and x/time 0.16.0 plus only the transitive companions selected by native module resolution.
- **FR-003**: Advance both module declarations and every documented or automated Go baseline from 1.25 to 1.26 because the selected Go modules require Go 1.26.
- **FR-004**: Select React and React DOM 19.3.0, matching 19.3.0 type packages, Node types 26.5.1, and Vite 8.3.0 plus native lockfile resolution.
- **FR-005**: Regenerate `go.sum`, `desktop/go.sum`, and `desktop/frontend/package-lock.json` only through native dependency tooling.
- **FR-006**: Preserve all released public behavior and security boundaries; do not weaken tests or compatibility assertions to accommodate updates.
- **FR-007**: Verify root and desktop graph cleanliness, frontend clean restoration and audit, focused affected surfaces, native desktop production build, and canonical verification.
- **FR-008**: Record the dated Go 1.26 pinned-baseline decision in the changelog and update contributor and desktop requirements guidance.
- **FR-009**: Publish a structured replacement pull request, close #240, #241, and #242 as superseded with traceable comments, and close #243 only through the replacement merge.
- **FR-010**: Address every review finding and reach green latest-head CI with no more than two requested Codex review rounds.

### Key Entities

- **Source proposal**: Automated pull request, requested versions, source head, and supersession disposition.
- **Dependency baseline**: Direct selections, resolved companions, runtime requirements, manifests, and integrity files.
- **Compatibility surface**: Product, security, accessibility, platform, packaging, and release contracts affected by a dependency.
- **Review evidence**: Final commit, CI state, review finding, response, and resolution.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All three source pull requests are represented by one clean S090 dependency baseline with zero omitted direct updates.
- **SC-002**: Root and desktop `go mod tidy` and `go mod verify` complete with zero residual diffs; frontend `npm ci` and high-severity audit complete without bypass flags.
- **SC-003**: Focused MCP, remote security, Windows, desktop, frontend, browser, and packaging evidence passes on the selected baseline.
- **SC-004**: All eight canonical gates pass on the final local head and all required hosted checks pass on the final PR head.
- **SC-005**: The Go baseline is consistent at 1.26 across both modules, all workflow consumers, and maintainer guidance, with zero stale 1.25 operational references.
- **SC-006**: Every review comment is dispositioned, all source pull requests have an explicit supersession record, and no third Codex review round is requested.

## Assumptions

- Current post-v1.4.0 `main` is the authoritative source boundary.
- `GOTOOLCHAIN=auto` may provision Go 1.26 locally; CI obtains the module-declared toolchain through `setup-go`.
- Node 26 remains the supported frontend runtime and requires no runtime-major change.
- The maintainer explicitly authorized branch publication, PR creation, in-scope review fixes, and up to two Codex review rounds.
