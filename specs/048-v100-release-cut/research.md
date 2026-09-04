# Research: v1.0.0 Release Cut

## Decision 1: Merge preparation before tagging

**Decision**: Complete and merge the S048 release-preparation PR before any
v1.0.0 tag is created.

**Rationale**: The Release workflow reads static README version lines,
tag-specific notes, and the changelog from the tagged commit. Tagging the
pre-S048 commit would omit the very boundary being reviewed and fail README
preflight. The constitution also reserves tag/release mutations for separate
explicit authority.

**Alternatives considered**:

- Tag the S047 merge and add evidence later. Rejected because the tagged source
  would lack reviewed v1.0.0 notes and version strings.
- Create and later move the tag. Rejected because mutable release identity is
  unsafe and prohibited by the promotion guard.

## Decision 2: v1.0.0 is the next release

**Decision**: Cut v1.0.0, not another 0.9.x release.

**Rationale**: The authoritative milestone is `v1.0.0 - Release readiness`, the
remaining ten issues are verification-only release criteria, and the merged
S038-S047 work constitutes the intended Windows and desktop v1 baseline.

## Decision 3: Keep the existing workflow unchanged

**Decision**: Reuse `.github/workflows/release.yml` and
`.github/workflows/promote-release.yml` without modification.

**Rationale**: S040 and S047 already provide draft-only staging, exact candidate
identity, 47-scenario attended validation, complete-asset enforcement, checksum
generation, and immutable promotion. A new orchestration path would duplicate
security-critical logic.

## Decision 4: Publish five concise highlights

**Decision**: Summarize scheduling capability, execution reliability, Windows
lifecycle control, desktop usability, and release assurance in five bullets.

**Rationale**: These are the broad user outcomes accumulated after v0.9.1. Five
bullets satisfy the established four-to-six policy while the tagged changelog
retains all implementation detail.

## Decision 5: Treat formal evidence as post-merge release work

**Decision**: The preparation contract specifies but does not fabricate or
pre-complete formal evidence.

**Rationale**: The candidate MSI does not exist until the reviewed tag workflow
finishes. S047 local-demo evidence is intentionally non-transferable, and native
rendering, physical input, maintenance, and multi-profile removal outcomes
require attended Windows observations against the exact staged bytes.

## Decision 6: Reconcile issues independently

**Decision**: Shared observation evidence may support multiple issues, but each
issue receives its own acceptance audit and closes independently.

**Rationale**: #96 is a coordinator and the other nine records have distinct
criteria. Closing by release association would violate the repository's GitHub
planning policy and could hide a partial acceptance failure.

## Decision 7: Keep Post-v1 findings out of v1.0.0

**Decision**: #102 and #118 through #120 remain outside the release boundary.

**Rationale**: They are explicitly assigned to the Post-v1 milestone. Pulling
them into S048 would invalidate the reviewed boundary and delay a release target
the maintainer already accepted.
