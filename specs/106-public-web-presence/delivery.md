# S106 delivery record

## Local handoff boundary

Base: 288128690d45b18cc5b0bade50bb138269f517ea, clean and equal to fetched origin/main before kickoff; remote refreshed and unchanged before the local commit. Review branch: codex/106-public-web-presence. At the initial local handoff, no push, pull request, tag, release, deployment or live domain mutation had been performed. The subsequent authorized publication stage is recorded below.

## Issue-level outcomes

| Issue | Local result | Public delivery |
| --- | --- | --- |
| [#9](https://github.com/shruggietech/go-schedule/issues/9) | Final subordinate README attribution links ShruggieTech to https://shruggie.tech | Awaiting review/publication |
| [#281](https://github.com/shruggietech/go-schedule/issues/281) | Concise purpose, installation, minimal platform-specific first task and focused references; release metadata retained | Awaiting review/publication |
| [#279](https://github.com/shruggietech/go-schedule/issues/279) | Static Next/Fumadocs shell, 21 canonical docs, static search, original heading IDs, legacy page/home fragment destinations, source links, assets and Actions deployment | Awaiting Pages source change and reviewed deployment |
| [#280](https://github.com/shruggietech/go-schedule/issues/280) | Branded responsive home, approved local fonts/art, workflow example, project resources and persistent system-aware appearance | Awaiting reviewed deployment |
| [#282](https://github.com/shruggietech/go-schedule/issues/282) | Eight-product, 80-cell sourced public v1.5.1 baseline; product/category filters; opportunity links; exact-content release reviews in all three publication paths; promotion-to-site refresh | Awaiting reviewed deployment; future releases need their own snapshot/review |
| [#278](https://github.com/shruggietech/go-schedule/issues/278) | Owner authentication, live DNS/Pages inventory and concrete cutover/recovery procedure prepared | Functional outcome remains incomplete until live ownership verification, domain binding, DNS and HTTPS cutover |

Parent [#276](https://github.com/shruggietech/go-schedule/issues/276) and all child issues remain open. Local implementation readiness does not assert public delivery or authorize issue closure.

## Authenticated inventory, 2026-10-19

Owner completed Cloudflare browser sign-in/MFA. The authenticated dashboard confirms ShruggieTech owns the active Cloudflare Registrar registration for go-schedule.io, with a full DNS zone on the Free plan. The actual records table shows 0 of 200 records and No DNS records. No A, AAAA, CNAME, TXT, MX or CAA records were present. Assigned nameservers: norah.ns.cloudflare.com and rex.ns.cloudflare.com. Account/zone identifiers and credentials are omitted. No settings changed.

Observed GitHub Pages before-state: build_type legacy, main:/docs, cname null, HTTPS enforcement true, public URL https://shruggietech.github.io/go-schedule/. The latest public application release remains v1.5.1, sourced from 86e8644f3919f1ca7b19f4fb8a78317c2dd6d4a6.

After publication approval: refresh the inventory, obtain the exact organization Pages TXT challenge, create/retain that TXT and finish verification, bind this repository's custom domain, switch Pages to Actions, deploy the reviewed export, apply the nine DNS-only traffic records and finish certificate provisioning/HTTPS enforcement. The owner may need to renew Cloudflare sign-in or supply organization-owner/repository-admin authority if the retained session lacks it. Values and provider references are in [website maintenance](../../docs/website.md).

Rollback starts from this recorded empty traffic-DNS baseline and Pages source/domain state. Restore prior publication content from the recorded pre-cutover revision through a reviewed revert or saved deployment, coordinate restoration of Pages settings and remove only records introduced by this cutover. Retain the organization ownership TXT and avoid traffic pointing at an unbound site.

## Verification, 2026-10-19

- The final foreground sh scripts/verify.sh all aggregate completed with exit 0: format, vet, lint, race, GUI, coverage, docs and automation, including all regression fixtures. Earlier documentation failures were corrected, then the complete aggregate was repeated. No required gate or safety-critical surface was skipped.
- Core statement coverage: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1%. Production desktop: 28 test files and 158 tests passed; Windows Wails build passed.
- Final website: seven Node tests, production static build and ten headless Chromium tests passed. Export audit covers all 21 docs, internal paths/fragments and unchanged canonical asset bytes. Browser checks cover search, appearance and storage failure, keyboard entry, diagrams, legacy query/hash destinations, original headings, 390/768/1440px layouts and comparison filters.
- Comparison Go regressions passed for changed reviewed content, missing/mismatched reviews, tag/source identity, missing evidence and invalid statuses. Publication fixtures reject missing, bypassed, late, conditional and soft-failure guards in all three release entry points. A parsed-workflow regression binds successful promotion names to site refresh despite token event suppression.
- Exact-lockfile npm install passed with zero vulnerabilities. Subsequent postinstall confirmed preparation creates generated docs before MDX compilation. Dependency and font license choices are recorded in site/DEPENDENCIES.md and existing brand license files.
- GitHub formatting and whitespace checks passed; all 35 authored site text files were checked for UTF-8, BOMs, replacement characters, common mojibake and forbidden em dashes. Generated brand assets are copied byte-for-byte.
- Visual review inspected production dark/light, mobile and documentation screenshots. The in-app browser could not reach loopback; the local Chromium suite supplied the inspected screenshots. Next emits an upstream Node module.register deprecation and two CSS optimizer warnings for Fumadocs search highlights; compilation, search and browser tests pass.

Hosted GitHub CI/CodeQL, real Pages deployment, public redirects, DNS propagation and public IPv4/IPv6 TLS remain unexecuted at the requested pre-publication boundary. No application release is being cut.

## Initial review handoff

Inspect the branch diff and proposed [PR body](pr-body.md). These publication commands are prepared only:

```sh
git push -u origin codex/106-public-web-presence
gh pr create --base main --head codex/106-public-web-presence --title "feat(106): rebuild public website and release-linked comparison" --body-file specs/106-public-web-presence/pr-body.md
```

Stop for owner review. The follow-on publication stage coordinates GitHub review/CI and domain cutover without treating #278 as already complete.

## Authorized PR publication, 2026-10-19

The owner approved pushing S106, opening its PR, monitoring CI and resolving findings through at most two Codex review rounds. The boundary remains before merge or public DNS changes. Branch commit 4d10db8489f98278a62c590e4f575c4f3652df0b was pushed and [PR #291](https://github.com/shruggietech/go-schedule/pull/291) opened under milestone Public website and project discovery. PR opening automatically started Codex round 1. No duplicate manual round was requested.

The initial hosted production desktop audit failed on all three platforms because the inherited frontend lockfile resolved source-map-js 1.2.1, affected by [GHSA-68fv-2mgg-jv7q](https://github.com/advisories/GHSA-68fv-2mgg-jv7q). Defect [#292](https://github.com/shruggietech/go-schedule/issues/292) tracks the finding. The compatible lockfile-only update resolves 1.2.2; an exact install and the unchanged high-severity audit report zero vulnerabilities. The website lockfile already resolves 1.2.2. Review results and final hosted check evidence are maintained on the PR so this historical record does not imply success for an untested later revision.

Codex round 1 completed on 4d10db8 with two publication-order findings, tracked together in [#293](https://github.com/shruggietech/go-schedule/issues/293). The five unchanged legacy Jekyll source files are retained while Pages still uses main:/docs. README discovery links use the live GitHub Pages address and repository comparison data. The cutover procedure makes successful replacement delivery precede source cleanup and domain-link changes, and a transition regression covers that boundary. No live settings were changed to address these findings.

After these fixes, the full foreground sh scripts/verify.sh all aggregate completed again with exit 0, including all automation fixtures. The final website export passes eight Node regressions and all ten Chromium checks. The unchanged audit gate passes after an exact desktop install. No application release or public deployment was performed.
