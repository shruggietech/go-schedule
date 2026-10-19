## What this changes

Replace the documentation-index home and Jekyll pipeline with a branded static Next.js/Fumadocs website. Canonical Markdown stays in docs/, with static search, original heading IDs, legacy destinations and byte-preserved brand assets. The concise README surfaces installation and first use and restores linked ShruggieTech attribution.

The dedicated comparison uses a sourced 80-cell v1.5.1 baseline, product/category filters and issue-linked opportunities. Staging and both promotion paths require a tag-specific review bound to dataset, changelog and release-note digests. Successful promotions refresh the site with the latest public snapshot.

Refs #9, #278, #279, #280, #281, #282. Parent #276. These references do not automatically close issues. #278's live outcome remains pending domain cutover; authentication and the recorded empty DNS inventory are preparation only.

## Verification

The foreground canonical aggregate completed with exit 0. Selected output:

```sh
sh scripts/verify.sh all
```

```text
github-format: OK - no em dashes or hard-wrapped Markdown prose
0 issues.
engine      82.9%  OK
schedule    89.1%  OK
timezone    91.3%  OK
store       80.1%  OK
catchup     88.9%  OK
logbus      91.1%  OK
docs-check: OK, 21 pages, links, metadata, fences, static export, and product policy clean
automation-check: OK - actions, CodeQL, Dependabot, release operations, release notes, brand, lifecycle, and 8 gates
automation-check-test: OK (automation)
```

- [x] The complete foreground canonical aggregate passed locally.

Final website: seven Node tests, static export with link/fragment/asset audits and ten Chromium tests passed. Clean dependency install reports zero vulnerabilities. Windows Wails build and 158 desktop frontend tests passed. See the [S106 delivery record](https://github.com/shruggietech/go-schedule/blob/codex/106-public-web-presence/specs/106-public-web-presence/delivery.md) for evidence and public-delivery boundaries.

## Checklist

- [x] No safety-critical test surface was weakened or skipped: clock injection, timezone and DST resolution, store migrations, restart and catch-up recovery, goroutine termination, IPC access control.
- [x] Pinned workflows changed; CHANGELOG.md carries the dated 2026-10-19 S106 build-protocol decision.
- [x] Documentation updated where behavior changed.

## Publication dependency

Coordinate Pages source/domain settings, organization ownership TXT, nine Cloudflare DNS-only traffic records and HTTPS using [the cutover procedure](https://github.com/shruggietech/go-schedule/blob/codex/106-public-web-presence/docs/website.md). No live settings changed during local implementation. Public deployment, hosted CI/CodeQL and DNS/TLS outcomes are not claimed from local tests.
