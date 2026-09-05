# Research: v1.0.0 Release Execution and Audit

## Decision 1: Preserve S049 as the immutable tag boundary

**Decision**: The authorized annotated `v1.0.0` tag remains fixed at `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`. S050 is a post-tag audit branch.

**Rationale**: S049 is the last reviewed commit that changes release tooling and documentation without changing runtime behavior. Moving the tag to include an S050 audit record would make the record part of the event it is supposed to observe and would create an endless new-boundary problem.

**Alternatives considered**:

- Merge S050 before tagging. Rejected because it violates the reviewed S049 publication contract.
- Omit S050 documentation entirely. Rejected because the user requested a Spec Kit slice and durable PR-reviewed release evidence is useful after publication.

## Decision 2: Reuse only production release gates

**Decision**: Use `verify-candidate`, the attended collector, `verify-bundle`, and `render-dispositions` exactly as reviewed in S047-S049.

**Rationale**: The promotion workflow invokes the same candidate and bundle validators. Reuse prevents a weaker S050-only acceptance path.

**Alternatives considered**:

- Populate evidence from the checked fixture. Rejected because fixtures are explicitly non-native and cannot authorize release.
- Copy the accepted S047 local-demo observations. Rejected because those results bind different bytes and provenance.

## Decision 3: Separate sensitive/raw evidence from the repository audit

**Decision**: Store the working evidence, screenshots, manifest, MSI, and issue packet in a fixed-volume workspace outside Git. Commit only safe environment facts, immutable identifiers, hashes, counts, URLs, and summarized outcomes.

**Rationale**: Screenshots and user-profile evidence may contain local machine details. The release archive is the canonical immutable evidence payload; Git should not duplicate it or leak unnecessary host data.

**Alternatives considered**:

- Commit the full ZIP and screenshots. Rejected because release assets already provide the canonical storage and Git would duplicate binary evidence.
- Commit no audit detail. Rejected because the PR would lack independently reviewable traceability.

## Decision 4: Reconcile leaves before coordinators

**Decision**: Review and update #98, #101, #104, #105, #106, #109, #111, #112, and #113 first. Then reconcile #96 from actual child/prerequisite states. Keep
#122 open until promotion and public audit succeed.

**Rationale**: This preserves issue-level acceptance criteria and prevents a coordinator checkbox from becoming evidence for its own children.

**Alternatives considered**:

- Close all issues from one bulk comment. Rejected because association with a release is not acceptance evidence.
- Close #96 when the evidence archive validates. Rejected because #96 also depends on actual child state and lifecycle criteria.

## Decision 5: Promote only through the reviewed workflow

**Decision**: Dispatch `Promote Release` with input `v1.0.0` only after the formal archive is uploaded and all readiness dispositions pass.

**Rationale**: The workflow independently re-downloads the exact staged set, validates provenance and evidence, generates checksums, verifies the final set, checks tag immutability, and only then publishes.

**Alternatives considered**:

- Run `gh release edit --draft=false`. Rejected because it bypasses the release gate and checksum creation.
- Rebuild locally and replace an asset. Rejected because any byte substitution invalidates the formal candidate.

## Decision 6: Treat attended proof honestly

**Decision**: A scenario is `pass` only from genuine observation on the exact installed candidate in the required Windows environment. Unavailable hardware, accounts, display classes, or interactive access remains a stop condition.

**Rationale**: The validator can prove schema and identity, but it cannot turn an unobserved visual, input, security, or lifecycle result into native evidence.

**Alternatives considered**:

- Infer pass from automated GUI/unit tests. Rejected because those tests cannot prove Windows rasterization, physical input, UAC/session behavior, or visible interaction states.
- Mark unavailable results as acceptable. Rejected because promotion requires all 47 observations to pass.
