# S134 Validation Quickstart

## Prerequisites

Use the existing Go 1.26.0, Node 26, C/WebView toolchain and hidden noninteractive Windows launch paths. Run long commands in the foreground, watching completion. Keep private logs and evidence under .test-temp/s134; preserve frozen prior observations.

## Focused source contracts

1. Retain #429's original CI failure and #430's original/recovery Pages receipts.
2. Run the updated timeout-policy cases against the old checker to retain the failing-first result, then after the supported-ceiling correction.
3. Exercise the existing production first-task browser flow after its route-completion correction, preserving all three platform/content/responsive expectations.
4. Inspect failure-only artifact wiring against the actual generated context/trace paths and ensure no unsuccessful exercise is hidden.
5. Restore exact dependency inputs; run root/desktop tidy/read-only graph/checksum checks, both clean npm installs/audits and the actual license/lock delta audit.

## Complete candidate qualification

Run the unmodified canonical driver:

```bash
sh scripts/verify.sh all
```

Watch all eight gates finish. Run the existing desktop and website browser contracts against their production bundles. Use scripts/cli-qa.ps1 to build and capture fresh Windows/WSL candidate evidence, review every current message under the established rubric, then finalize, validate and render. Keep the original unsuccessful observations, explicit not-applicable cases and owned cleanup results.

Run release-source-check for v2.0.1 and comparison-check against the exact committed source when available. Run github-format, lifecycle and history/encoding checks before publishing text. Source and protected-input snapshots must show no unexpected drift.

## Handoff

Record source identity, all actual local/hosted outcomes, the two-request review ledger, unresolved threads and unverified native observations. Use handoff.md for later owner integration and exact-source/actual-asset approvals. A source PR and QA result do not publish v2.0.1 or complete #419.
