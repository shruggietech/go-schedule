# Research: Lightweight Pull-Request Workflow

## Decision 1 - Pull requests are a review venue, not an enforcement system

**Decision**: Use a review branch and pull request for maintainer and agent work, primarily to obtain third-party AI review and preserve the discussion.

**Rationale**: This is a one-developer project with no users. A PR provides the desired review surface without requiring enterprise governance controls.

## Decision 2 - Keep enforcement social and lightweight

**Decision**: Do not add branch protection, required approvals, required conversation resolution, fixed status-check contexts, administrator rules, or a repository-owned settings payload.

**Rationale**: Those controls add maintenance and failure modes without solving a present project risk. The maintainer already controls the final merge.

## Decision 3 - Keep existing verification unchanged

**Decision**: Run `sh scripts/verify.sh all` before publication and allow the existing PR-triggered CI to provide additional evidence. Do not change CI or turn its job names into governance API configuration.

**Rationale**: The aggregate already defines local green. Hosted results remain useful to reviewers without becoming a new policy subsystem.

## Decision 4 - Respond to review with judgment

**Decision**: Consider each AI review comment. Implement suggestions that improve correctness or clarity; reply with a concise rationale when a suggestion is not warranted. The maintainer decides when review is sufficient.

**Rationale**: AI feedback is advisory. Treating it as mechanically binding would replace engineering judgment with ceremony.

## Decision 5 - Publication authorization covers the review cycle

**Decision**: The authorization that pushes the branch and opens its PR also covers verified, in-scope pushes to that same branch for CI or review fixes until the PR is merged or closed. It excludes material scope expansion, another PR, tags, and releases.

**Rationale**: Requiring another halt for each review-fix commit would contradict the one-halt protocol and turn ordinary review iteration into the red tape this feature deliberately avoids.
