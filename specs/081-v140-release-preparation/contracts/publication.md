# Contract: Cumulative v1.4.0 Publication

## Preparation

The preparation branch converts the complete post-v1.1.1 Unreleased record into one dated v1.4.0 section, retains a new empty Unreleased section, adds exactly four concise release highlights and one tagged changelog link, synchronizes durable candidate-aware documentation, and passes focused plus canonical verification. It does not create a tag or release.

## Historical integrity

The v1.2.0 and v1.3.0 milestones remain closed implementation records. Release operations do not create tags or GitHub Releases for those versions from current main and do not imply that they were previously public. The cumulative v1.4.0 changelog comparison begins at v1.1.1.

## Tagging

After the preparation pull request is reviewed and merged, fetch and prune origin, fast-forward `main`, require an empty tree, require local and remote main equality, require successful push CI for the exact merge commit, and require local tag, remote tag, and GitHub release state for v1.4.0 to be absent. Create one annotated v1.4.0 tag at that exact commit only after explicit tag authorization and never move it.

## Staging

Accept only the successful tag-push Release workflow whose head SHA matches v1.4.0. Before artifact mutation, its preflight must require the matching README health version, one dated v1.4.0 changelog section, the tag-specific release note, successful main CI for the exact SHA, and absent or draft release state. The staged GitHub Release remains draft and contains four headless archives, two Wails desktop archives, one Windows MSI, and one candidate manifest.

## Qualification

Download the exact staged MSI and candidate manifest. On one clean Windows 11 snapshot, complete the fresh-install path and the established attended matrix. On a separate snapshot, download the public v1.1.1 MSI from its canonical GitHub release URL, require SHA-256 `f7ac8f56f28330b016eb6e505e424b19e9bfbe435591cfbc54a723c91ac8e567`, install it, create representative tasks and preferences, capture its identity and service state, then apply the exact staged v1.4.0 MSI. Prove state retention, Wails desktop usability, service and local CLI operation, and machine-validated disabled notification, localhost MCP, and remote HTTPS states. All 47 attended observations and their attachments must bind to the staged candidate. Never substitute a local build or manufacture a result.

## Promotion

Upload the validated attended archive and dispatch Promote Release for v1.4.0 only after explicit promotion authorization. The workflow must revalidate candidate provenance, the evidence bundle, immutable tag, exact asset inventory, and final checksums before making the existing draft public. Promotion never rebuilds artifacts.

## Final audit

Freshly download every public asset, verify `SHA256SUMS.txt`, confirm the release is public and latest, and compare the tag, main commit, README, changelog, release note, candidate manifest, package names, sampled binary versions, issue, project item, and milestone. Only then complete issue #226, set its project item to Done, and close the v1.4.0 milestone.

## Failure handling

A missing, pending, failed, cancelled, skipped, stale, mismatched, extra, rebuilt, or unverifiable input stops the sequence. Before tagging, preserve reviewed source. After staging, leave the release draft. Do not close planning records or represent v1.4.0 as public until the final audit passes.
