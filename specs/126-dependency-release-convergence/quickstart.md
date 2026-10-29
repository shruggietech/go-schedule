# S126 Validation and Continuation

Use the established hidden launcher on Windows; never open a visible console or attended app on the owner's active desktop.

1. Regenerate selected Go modules using go mod tidy in root and desktop, inspect exact diffs, then exercise read-only graphs with GOFLAGS=-mod=readonly.
2. Restore both npm projects with npm ci, retaining their selected versions and coupled peers.
3. Build/capture fresh QA with scripts/cli-qa.ps1 under .test-temp/s126/qa-inputs-01 and qa-run-01, existing Ubuntu-24.04 user h8rt3rmin8r, and incidents. Review corpus.json against exact occurrences, then finalize, validate and render. See [maintainer execution contract](https://github.com/shruggietech/go-schedule/blob/main/docs/maintainers/cli-qa-run-v1.md).
4. Run full headless frontend browser contracts with one worker, inspect static site routes/search/anchors and dark/light responsive brand rendering, then run sh scripts/verify.sh all in the foreground.
5. Run release-source-check for v1.7.0 and comparison-check against the supplied source identity. These source guards do not establish tag/package existence.
6. Audit exact PR head CI and every review comment, no more than two Codex rounds. Stop for owner merge review.

After owner merge, continue actual release delivery within S126 using [the existing handoff](../125-reliability-delivery-decision/handoff.md): new exact-main acceptance, exact tag/draft approval, real staged asset inspection, separate asset-promotion approval and public-discovery reconciliation. No release identity is preselected by this source guide.
