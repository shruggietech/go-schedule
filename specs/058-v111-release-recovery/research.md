# Research: v1.1.1 Release Recovery

## Historical tag disposition

**Decision**: Preserve the existing v1.1.0 tag at its original commit and never publish it.

**Rationale**: A Git tag is an immutable release identity even when the associated GitHub release is still a draft. Moving or recreating it would erase provenance and make previously staged checksums ambiguous.

## Corrective version

**Decision**: Advance to v1.1.1.

**Rationale**: PR #143 corrects release-blocking behavior after the v1.1.0 tag was created. The smallest honest semantic version increment is a patch, and v1.1.1 can carry the complete v1.1 feature set as its first public distribution.

## Draft disposition

**Decision**: Delete only the unpublished v1.1.0 GitHub draft during post-merge release operations, after rechecking that it is still a draft and still targets the historical tag.

**Rationale**: Keeping the tag preserves source history, while retiring the draft removes the risk that unqualified artifacts are later published accidentally.

## Planning records

**Decision**: Rename milestone #4 and revise issue #140 for v1.1.1 instead of creating new records.

**Rationale**: The desired outcome remains one public v1.1 release. Updating the interrupted release coordinator preserves its discussion and evidence while preventing duplicate tracking.

## Release-note shape

**Decision**: Publish exactly four short highlights and one final tagged changelog link.

**Rationale**: This is the shortest shape accepted by the repository validator and matches the maintainer's standing preference for brief release copy.

## Artifact path

**Decision**: Reuse the existing draft staging and no-rebuild promotion workflows after the exact-commit CI correction in PR #143.

**Rationale**: The workflow now requires successful main-branch CI for the precise tag commit before building. Replacing the pipeline would expand scope and risk without improving the release guarantee.

## Automation regression closure

**Decision**: Update the approved release-workflow fixture for the exact-commit CI contract and execute its fixture suite from the canonical automation gate.

**Rationale**: Direct workflow validation can pass while its synthetic approved baseline is stale. Running both from one gate prevents future validator changes from escaping CI without fixture reconciliation.
