# Verification: v0.9.0 Release Cut

## Baseline (2026-08-30)

- Repository: `shruggietech/go-schedule`.
- Starting commit: `123590cead230feefd699eda87b29a84904980bd` on a clean, synchronized `main`.
- Working branch: `codex/034-v090-release-cut`.
- Latest tag and GitHub release: `v0.8.0`; neither a `v0.9.0` tag nor release existed.
- Milestone #1, `v0.9.0 - Polish and hardening`, was open with one open issue.
- Issue #79, `Cut v0.9.0 and verify release artifacts`, was open in milestone #1 with the release-preparation labels defined by the repository taxonomy.
- The pre-cut Unreleased range contained 434 lines and 61 top-level entries between `## [Unreleased]` and `## [0.8.0] - 2026-07-23`. Its normalized section SHA-256 was `e08ac07dcc70d865636faca72c0bfb80d0672a3464cdc37bb4ba649771679fca`.
- The release workflow used `generate_release_notes: true` plus a long inline installation body. This conflicted with the requested highlights-only public release copy.

## Expected v0.9.0 Assets

1. `go-schedule_v0.9.0_linux_amd64.tar.gz`
2. `go-schedule_v0.9.0_linux_arm64.tar.gz`
3. `go-schedule_v0.9.0_darwin_amd64.tar.gz`
4. `go-schedule_v0.9.0_darwin_arm64.tar.gz`
5. `go-schedule-desktop_v0.9.0_linux_amd64.tar.gz`
6. `go-schedule-desktop_v0.9.0_darwin_arm64.tar.gz`
7. `go-schedule_v0.9.0_windows_amd64.msi`
8. `SHA256SUMS.txt`

The seven non-checksum assets must be non-empty and represented exactly once in `SHA256SUMS.txt`.

## Test-First Evidence

After adding the release-copy fixtures but before changing the checker, the focused suite failed as intended:

```text
automation-check-test: FAIL: generated-notes should fail
```

The failure proved that the former automation contract did not reject GitHub-generated release notes.

## Publication Guards

- Tag creation requires a separate explicit authorization after the preparation pull request is reviewed and merged.
- Immediately before tagging, clean local `main`, `origin/main`, and the reviewed S034 merge commit must be identical.
- Both `git tag --list v0.9.0` and `gh release view v0.9.0` must confirm absence before tag creation.
- The README badge and sample health version must identify v0.9.0 after the release workflow's synchronization job finishes.
- Issue #79 remains open until the workflow, release body, complete asset set, checksums, and README synchronization are evidenced on the issue.
- Manual install/uninstall observation and Windows VM work are outside S034.

## Results

The first complete gate run reached lint and stopped before analyzing source:

```text
Error: can't load config: the Go language version (go1.24) used to build
golangci-lint is lower than the targeted Go version (1.25.0)
```

This was an existing verifier-pin incompatibility rather than a lint finding. The release preparation updates the canonical and documented pin from v2.1.6 to v2.12.0, the latest upstream release observed on 2026-08-30 whose declared build baseline remains Go 1.25.

The final verification results were:

- ShellCheck: PASS for `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`.
- Focused automation fixtures: PASS, including generated-note, fixed-path, missing-link, highlight-count, and exhaustive-copy rejections.
- Complete canonical gate: PASS for format, vet, lint, race, GUI, coverage, docs, and automation. The six coverage thresholds were 84.1% through 91.3%.
- Changelog retention: PASS. The 61 original top-level Unreleased entries are retained under v0.9.0; the section now has 64 top-level entries after adding the S034 release-note decision, lint-gate fix, and review hardening.
- Workflow scope: PASS. The workflow diff changes only the release body's note source; package jobs, dependency ordering, asset globs, matrices, and final checksum job are unchanged.
- Release-copy audit: PASS at nine lines, five highlight bullets, one second-level heading, and one tagged full-changelog link.
- Repository hygiene: PASS for `git diff --check`, UTF-8 without BOM across all 20 changed or added files, and zero mojibake candidates.
- Hosted guards: issue #79 remains open, and no v0.9.0 tag or release was created during implementation.

## Review Hardening

Codex review identified that an H2 count alone did not reject extra plain paragraphs or H3 sections. The checker now validates the complete physical document shape, and focused negative fixtures cover both bypasses in addition to the original H2 installation case.
