# Quickstart: Reviewing Bundle Target Compatibility

1. Select the intended local or remote daemon in the desktop or CLI profile.
2. Run `gosched bundle validate bundle.json` or use the desktop bundle validator. Missing bundle or source-family capabilities appear as named findings.
3. For a compatible bundle with new watchers, supply target-local absolute paths during preview. The selected daemon interprets those paths using its own platform rules.
4. Review the target identity, conflicts, and target-only drift. Apply only the returned single-use plan if the destination is correct.
5. If the daemon lacks `bundles`, upgrade or select a capable target. Do not assume a failed endpoint call means the bundle was applied.
