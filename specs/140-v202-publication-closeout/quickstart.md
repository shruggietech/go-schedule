# S140 verification scenarios

Use the existing hidden foreground launcher with redirected noninteractive children. Keep every run under a new ignored S140 directory.

1. Check eight products, ten stable rows, 80 uniquely addressed cells, valid sources/dates and preserved opportunity links. Review the ten status changes and five retained evidence dates against research.md.
2. Create the new review with node site/scripts/review-release.mjs v2.0.2. Confirm the existing guard rejects review-required, retain that observation, then complete all ten reviewed decisions. Validate with go run ./scripts/comparison-check -tag v2.0.2 -commit <actual-40-character-HEAD> and the historical aggregate.
3. Run go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v2.0.2. This establishes source preparation only.
4. Build the production website with SITE_RELEASE_TAG=v2.0.2 and run the existing eleven headless Chromium contracts with private outputs. Restore the default production build, confirm its comparison is v2.0.1 and run the same contracts with separate outputs.
5. Run scripts/verify.sh all in the foreground: format, vet, lint, race, gui, coverage, docs and automation. The gui gate builds/tests without launching the application. Preserve exact source manifests and original evidence bytes.
6. Run lifecycle, GitHub-format, UTF-8/BOM/mojibake and diff checks on final records, then publish the official source PR. Confirm exact-head terminal applicable checks, review findings and no auto-merge before owner handoff.
