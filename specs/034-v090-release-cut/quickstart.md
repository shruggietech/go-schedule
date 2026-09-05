# Quickstart: Validate v0.9.0 Release Preparation

## Before Publication

Confirm the branch starts from current `origin/main`, then run the complete repository gate:

```bash
sh scripts/verify.sh all
```

Inspect the release-copy contract:

```bash
sh test/scripts/automation-check_test.sh automation
sh scripts/automation-check.sh .
```

The prepared release notes must contain four to six highlights, one link to the tagged v0.9.0 changelog, and no generated-notes inventory.

## Before Tagging

After the preparation pull request merges and tag authorization is given:

```bash
git fetch --prune origin
git switch main
git pull --ff-only origin main
git status --short --branch
git tag --list v0.9.0
gh release view v0.9.0
```

The worktree must be clean, local and remote `main` must identify the reviewed merge commit, and both tag and release lookups must show that v0.9.0 is absent.

## After Tagging

Watch the single tag-triggered workflow to completion. Then download every release asset into a temporary directory, verify `SHA256SUMS.txt`, compare the asset names with [publication.md](contracts/publication.md), inspect the public body, and confirm the README synchronization commit reached `main`.

The checksum audit must run from the directory containing the downloaded assets:

```bash
sha256sum --check SHA256SUMS.txt
```

Also confirm that the release body matches `.github/release-notes/v0.9.0.md`, is neither draft nor prerelease, and links to the v0.9.0 changelog rather than the moving Unreleased section.

Issue #79 closes only after all of that evidence is recorded.
