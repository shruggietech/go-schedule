# Quickstart: Cumulative v1.4.0 Release Preparation

## 1. Validate source-owned release metadata

```powershell
go run ./scripts/github-format
git diff --check
bash test/scripts/automation-check_test.sh
```

Confirm an empty Unreleased section, one dated v1.4.0 section, a direct v1.1.1 comparison, exactly four release-note highlights, one tagged changelog link, README health version 1.4.0, and durable candidate-aware installation wording.

## 2. Run the canonical local gates

```powershell
& 'C:\Program Files\Git\bin\sh.exe' scripts/verify.sh all
```

All eight gates must pass in the foreground without S081 exclusions.

## 3. Review and merge

Push `codex/081-v140-release-preparation`, open one pull request that references #226 without closing it, satisfy hosted CI and every review finding, then ask the maintainer to perform the final review and squash merge. The merge does not authorize a tag or release.

## 4. Stage only after explicit authorization

Synchronize clean `main`, verify exact-commit push CI, confirm v1.4.0 is absent locally and remotely, create the immutable annotated tag at the reviewed merge commit, and wait for the Release workflow to stage the complete draft artifact set.

## 5. Qualify fresh and upgrade paths

Use separate clean Windows 11 snapshots. Qualify a fresh install on one. On the other, install the public v1.1.1 MSI, create representative scheduler state, then upgrade with the exact staged v1.4.0 MSI. Complete and validate all 47 attended observations against that candidate.

## 6. Promote and audit only after explicit authorization

Upload the candidate-bound evidence archive, dispatch Promote Release, verify the no-rebuild public artifact set and checksums, confirm v1.4.0 is latest, then reconcile issue #226, its project item, and the v1.4.0 milestone.
