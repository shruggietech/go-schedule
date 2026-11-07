# S135 Validation Quickstart

## Prerequisites

Use the existing hidden noninteractive Windows launchers and fresh .test-temp/s135 output paths. Keep accepted source 25f42d6eb35ff0cec2753f0f2934396ba50df954 and exact-main CI/CodeQL/Pages receipts separate from the future S135 PR candidate. No tag, staging or promotion runs before its required concrete authorization.

## Release and package readback

Run release-source-check for v2.0.1 and comparison-check bound to accepted source. After staging approval, monitor the tag-triggered Release workflow and retain its actual attempt. Download the complete draft asset set by selected identity; compare sizes/digests, archive/build metadata and Windows MSI/candidate manifest without installing. After actual-asset approval and waived promotion, independently download all nine public assets, verify every checksum entry and unchanged original asset identity, and read public/latest plus Pages/HTTPS discovery.

## Current discovery source

Update only current availability notices and the live fallback after publication. Preserve historical v2.0.0 records and tagged v2.0.1 notes/comparison/review bytes. Update the existing production browser availability assertion. Run:

```bash
sh scripts/verify.sh all
```

Watch all eight gates complete in the foreground. Exercise the existing production site browser suite with a fresh S135 output configuration, then run exact-source release/comparison, lifecycle, encoding and GitHub formatting checks. Retain source-before/after and protected-input inventories; record the intended fallback transition separately.

## Handoff

Record actual release/run/asset IDs, public readback, candidate SHA/checks, review count and unresolved findings. Attended native installation/GUI remains untested. The official PR remains open for owner final review and merge; #419's pending discovery integration/deployment must remain explicit.
