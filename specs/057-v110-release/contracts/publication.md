# Contract: v1.1.0 Publication

## Preparation

The preparation branch cuts the changelog boundary without losing entries, updates the README example, adds four concise highlights with one tagged changelog link, and passes the full local and hosted verification suites. The push-triggered CI workflow for the reviewed merge commit must also complete successfully before release staging can begin.

## Tagging

Immediately before creating `v1.1.0`, fetch and prune origin, fast-forward `main`, require an empty tree, require local and remote main equality, and require the local tag, remote tag, and GitHub release to be absent. Create one annotated tag at the reviewed S057 merge commit and never move it.

## Staging

Accept only the successful tag-push Release workflow whose head SHA matches the tag. Its preflight must locate the push-triggered CI run on `main` for that exact SHA and require a successful conclusion before any package is built or uploaded. A missing, pending, failed, cancelled, or timed-out CI result cannot stage release artifacts. The release must remain draft and contain the complete expected packages plus `windows-candidate-manifest.json`. Verify the manifest and MSI before any installation.

## Qualification

Use the exact staged MSI and candidate identity. Do not substitute a local build or prior release artifact. Preserve only genuine attended evidence and never manufacture an observation that was not performed.

## Promotion

Upload the validated attended archive, dispatch `Promote Release` for `v1.1.0`, and require the workflow to verify the candidate, evidence, tag, asset set, and checksums before changing the existing draft to public. Promotion must not rebuild artifacts.

## Final audit

Freshly download every public asset, verify `SHA256SUMS.txt`, confirm the public release is latest, and compare the release tag, commit, README, changelog, release notes, manifest, and sampled binary versions. Close issue #140 and milestone v1.1.0 only after the audit passes.
