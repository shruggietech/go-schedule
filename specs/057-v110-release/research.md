# Research: v1.1.0 Release

## Version selection

**Decision**: Release v1.1.0.

**Rationale**: S051 through S056 add backward-compatible task authoring, diagnostics, layout persistence, external triggers, Trigger Sets, and filesystem watchers. Semantic Versioning therefore calls for a minor increment from v1.0.0.

## Release-note shape

**Decision**: Publish exactly four short highlights and one final tagged changelog link.

**Rationale**: This is the shortest shape accepted by the repository's existing four-to-six-highlight validator and directly implements the maintainer's request for brief, clean notes.

## Artifact path

**Decision**: Reuse the existing draft staging and no-rebuild promotion workflows.

**Rationale**: Those workflows already bind artifacts to one immutable tag, generate a Windows candidate manifest, require qualification evidence, verify asset cardinality, and create a complete checksum inventory. Replacing them during a release cut would add risk without user value.

## Human interaction

**Decision**: Automate source preparation, verification, review handling, tagging, staging, downloads, candidate inspection, promotion, and final audit; request only the smallest attended Windows confirmation that cannot be established headlessly.

**Rationale**: The maintainer explicitly requested minimal input, while the project still treats native installation behavior as evidence that cannot be fabricated from CI.
