# Quickstart: S098 Portable Bundles

1. Create two daemon workspaces with matching and deliberately different task or group state.
2. Export the source bundle with `gosched bundle export --json > bundle.json`.
3. Confirm `bundle.json` contains schema `go-schedule.bundle/v1` and no source daemon identifier, trigger key, notification endpoint, credentials, command, working directory, or environment.
4. Validate and compare with an explicitly selected target: `gosched bundle compare --file bundle.json --target <profile> --json`.
5. Inspect the returned target daemon identity, bundle digest, plan ID, changes, conflicts, incompatibilities, and target-only drift. Confirm no target records changed.
6. Apply only after reviewing the plan: `gosched bundle apply --file bundle.json --target <profile> --plan <plan-id> --json`.
7. Confirm every item reports an outcome. Re-run compare to inspect remaining drift. Do not expect omitted definitions to be deleted.

The desktop Bundles workflow provides the same sequence and keeps the selected target visible in every screen state.
