# Contract: v1.1.1 Publication Recovery

## Historical preservation

The existing v1.1.0 tag remains at its original commit and is never moved, deleted, recreated, or published. Before deleting the associated GitHub release draft, operations must prove that the release is still a draft and still targets v1.1.0. Draft deletion must not delete the Git tag.

## Preparation

The preparation branch adds the corrective v1.1.1 changelog boundary, updates the README example, adds four concise highlights with one tagged changelog link, and passes the full local and hosted verification suites. The existing release coordinator and milestone are revised for v1.1.1 rather than duplicated.

## Tagging

After the preparation pull request is reviewed and merged, fetch and prune origin, fast-forward `main`, require an empty tree, require local and remote main equality, revalidate the preserved v1.1.0 tag, retire only the unpublished v1.1.0 draft, and require local tag, remote tag, and GitHub release state for v1.1.1 to be absent. Create one annotated v1.1.1 tag at the reviewed S058 merge commit and never move it.

## Staging

Accept only the successful tag-push Release workflow whose head SHA matches v1.1.1. Its preflight must locate the push-triggered CI run on `main` for that exact SHA and require a successful conclusion before any package is built or uploaded. A missing, pending, failed, cancelled, or timed-out CI result cannot stage release artifacts. The release must remain draft and contain the complete expected packages plus `windows-candidate-manifest.json`.

## Qualification

Use the exact staged MSI and candidate identity. Do not substitute a local build or prior release artifact. Preserve only genuine attended evidence and never manufacture an observation that was not performed.

## Promotion

Upload the validated attended archive, dispatch `Promote Release` for v1.1.1, and require the workflow to verify the candidate, evidence, tag, asset set, and checksums before changing the existing draft to public. Promotion must not rebuild artifacts.

## Final audit

Freshly download every public asset, verify `SHA256SUMS.txt`, confirm the public release is latest, and compare the release tag, commit, README, changelog, release notes, manifest, and sampled binary versions. Close issue #140 and milestone v1.1.1 only after the audit passes.
