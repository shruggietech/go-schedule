# Quickstart: Stable Daemon Identity and Capability Manifest

## Review the slice

```bash
git diff --check
go run ./scripts/github-format
go test ./internal/store ./internal/api/server ./internal/api/client ./internal/cli ./desktop/connection
```

Confirm that `GET /v1/health` is unchanged, manifest collections are sorted and duplicate-free, no prohibited host data appears in the response type, and reset uses exact current-identity acknowledgement.

## Exercise the operator workflow

With the daemon running on its protected local transport:

```bash
gosched daemon manifest
gosched daemon manifest --json
gosched daemon rename "Workshop scheduler"
gosched daemon reset-identity --confirm <current-installation-id>
```

Restart the daemon after rename and confirm the identifier and name persist. Repeat the reset with the old identifier and confirm it fails without changing state.

## Verify lifecycle behavior

Run the focused migration and store tests to verify fresh-store uniqueness, restart stability, prior-schema preservation, restore preservation, clone semantics, and reset invariants. Finish with canonical verification:

```bash
sh scripts/verify.sh all
```
