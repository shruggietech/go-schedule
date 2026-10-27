# S119 validation guide

Prerequisites: Existing authenticated GitHub access, hidden native runner, current toolchain/caches and existing private WSL automation bridge. No installed-product or account-session operation.

1. Verify remote, approved source, tag absence/existing identity and exact-main successful push CI. Run existing release-source/comparison/format guards without mutating releases.
2. After Spec-Kit analysis and exact-tag approval, use the existing tag-only staging route. Retain the eight assets and validate actual run/attempt/Windows manifest identities using existing windows-release-gate verify-candidate.
3. Present complete candidate identities for separate promotion approval. After approval, run the existing standing-waiver promotion workflow and verify actual release/latest, nine final assets and SHA256SUMS.txt.
4. Reproduce the current responsive browser version mismatch with an owned prepared v1.6.0 build; retain the original failed expectation. Correct only the demonstrated assertion, then validate exact selected-version visibility with the existing browser runner.
5. After actual publication, reconcile the default selected tag and current canonical status prose. Preserve immutable tagged notes/changelog/comparison bytes and record their original receipt separately from the reconciled current main changelog digest/date; verify current source using existing guards and source snapshot hashes.
6. Run all eight foreground gates with scripts/verify.sh all, keeping finished inputs frozen. Finalize factual source tasks/lifecycle, publish the authorized PR and complete exact-head CI plus up to two Codex rounds. Owner decides merge.
