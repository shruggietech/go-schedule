# S133 Source Handoff and Same-slice Release Continuation

## Authorized execution boundary

The owner authorizes S133 source autopilot, automatic push/official PR, CI monitoring and in-scope fixes, with at most two total Codex requests including automatic PR-open review. Stop ready for final owner review and merge. No merge, tag, draft release, staging dispatch or promotion is authorized.

## Prepared behavior

v2.0.1 packages will include the merged selected own Run ID and complete first-task guidance, Wails runtime/build-tool 2.17.1, YAML 3.0.6, Node types 26.6.6 and Fumadocs core/UI 16.17.9. Public latest remains v2.0.0; candidate source checks and new comparison records do not select a public release. Current state/API and /v2 identity remain. #424 supplies the source correction; #419 retains actual public delivery.

## After owner integration

1. Verify the actual merged source/tree and clean main, then terminal successful exact-main CI/CodeQL and website acceptance. Read live issue/release/tag state; retain the source PR and review ledger.
2. Re-run existing release-source-check for v2.0.1 and comparison-check bound to the exact accepted main SHA. Ask for approval naming that immutable source for the tag and draft staging. Never reuse this source-PR approval as release approval.
3. After exact approval, create/push the immutable v2.0.1 tag and monitor the existing Release workflow to successful completion. Resolve the actual draft ID and complete platform packages plus windows-candidate-manifest.json; inspect actual archive/MSI metadata and SHA256 versus GitHub/manifest digests into fresh owned paths. Preserve all v2.0.0 bytes. Source QA is distinct from those staged bytes.
4. Present the actual draft, source, workflow attempt and inspected asset inventory for separate concrete-asset promotion approval. Use existing promote-release-waived.yml only after that approval; preserve accepted bytes and add/verify matching SHA256SUMS.txt through existing identity guards.
5. Verify public/latest identity, downloads/checksums, Pages deployment and HTTPS discovery independently. Reconcile current fallback/README/first-task availability through normal reviewed source integration within S133. Close #419 only for its actual public package/discovery outcome.

Attended native installation/GUI observations remain untested under the standing waiver. No active-workstation service or GUI mutation, historical upgrade matrix, new release framework, unchanged Windows replay or cluster reactivation is part of this continuation. Original #355/#356 causes remain unresolved.
