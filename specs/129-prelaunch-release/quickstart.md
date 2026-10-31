# S129 Validation Guide

Use a clean review branch and the verified hidden noninteractive launcher on Windows. Never install or restart the product on the active workstation. Preserve the 215 retained inputs and original failures in private receipts.

## Source preparation

1. Run github-format, confirm UTF-8 without BOM, staged diff whitespace and matching branch/spec/selector.
2. Check both Go module graphs read-only; root module and desktop imports/dependency must use /v2 and current dependencies/checksums remain valid.
3. Exercise a disposable tag-injected CLI build and version/health metadata through owned headless paths. This is local build evidence, not a staged package or public release.
4. Run the existing guards:

```sh
go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v2.0.0
go run ./scripts/comparison-check
go run ./scripts/comparison-check -tag v2.0.0 -commit <actual-reviewed-source-sha>
```

The last source guard binds local reviewed content and the supplied SHA; it does not prove a Git tag exists. Before a commit exists, identify the baseline plus dirty source digest explicitly.

## Complete verification

Run `sh scripts/verify.sh all` in the foreground to terminal completion. Preserve all eight existing gates and deadlines. Run the existing desktop/site component and headless browser suites after clean npm restores, redirecting only new legacy-named site screenshots into S129-owned captures. Run ShellCheck for all test shell scripts if touching the module fixture and the project PowerShell compliance checker. Check source/retained-input drift, immutable historical release/comparison bytes, lifecycle and staged whitespace before commit.

## Review and later publication

After authorized push/official PR creation, verify exact-head checks/workflow conclusions and every review thread/comment/body reaction. Count automatic review as request one; at most one manual second request. Stop for owner merge after all received findings are handled and CI is green. Follow [release contract](contracts/release.md) and handoff.md for later approvals and actual asset/public acceptance. Do not dispatch release workflows during source preparation.
