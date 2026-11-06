# S134 Research and Decision Register

## Observed baseline and research method

Accepted main is e09231bf0074480e5765a776f7fbf61c8b5610b1, tree 84ae5aeb7e346fe2d3f3fde7cc84f044c3e4a456. Read-only research used the installed planning command's research-agent requirement; it is not a PR review request. Evidence and exact package/license metadata are retained in .test-temp/s134/research-findings.json and associated files.

## Navigation and diagnostic retention

**Decision**: Await the expected first-task route in every platform loop before accepting its existing heading and advancing. Retain generated failure context/traces with the existing CI artifact action.

**Rationale**: [#429](https://github.com/shruggietech/go-schedule/issues/429) records navigation to the macOS installer interrupted by a still-pending first-task navigation in run 42225299776. The loop currently waits only for a heading; the initial entry already waits for its URL. The existing website job does not upload the trace/context paths named in that failure.

**Alternatives considered**: Sleeps, retries, removed platform entries and reduced assertions hide the ordering contract. A router abstraction or new test framework is unnecessary.

## Effective Pages policy

**Decision**: Set and require 600000 ms for the selected action, preserving outer jobs and the remaining S116 allowances. Document the supported ceiling and existing recovery.

**Rationale**: [#430](https://github.com/shruggietech/go-schedule/issues/430) retains the failed and successful retry receipts for the same source/artifact. The [exact action implementation](https://github.com/actions/deploy-pages/blob/368f82528645a54fb793d4d04e342629a3f51346/src/internal/deployment.js#L26) enforces MAX_TIMEOUT = 600000. The repository checker currently requires the unsupported one-hour input and must change with its regression cases. Closed #373 explicitly permits supported-ceiling classification.

**Alternatives considered**: Keeping a silently clamped input misstates the effective budget. A custom Pages API polling/deployment framework is disproportionate; existing source/artifact-bound reruns already recovered the observed deployment without replacing its failure evidence.

## Bounded dependency candidate

**Decision**: Include pflag 1.1.0, desktop plugin-react 6.2.1/Vitest 5.1.0, Fumadocs core/UI 16.17.10, Mermaid 12.3.0, Next 16.4.5 and Tailwind/PostCSS 4.3.4.

**Rationale**: #426-#428 share the maintenance source boundary. Published engines/peers remain compatible with Node 26, React 19.3.1 and the existing toolchains. The provided npm proposals add no installed package paths; the site proposal removes two nested entries. Verify the final 16.4.5 lockfile, clean installs, audits and actual candidate before drawing qualification conclusions. Read-only resolution of #426 succeeded in root and desktop without metadata changes.

**Alternatives considered**: Separate source PRs duplicate one integration boundary. Forcing an unused desktop pflag indirect entry would replicate a past module-graph assumption without a demonstrated need.

## Next patch correction and applicability

**Decision**: Supersede #428's 16.4.4 component with published 16.4.5.

**Rationale**: The [official November release](https://nextjs.org/blog/november-2026-security-release), [version release](https://github.com/vercel/next.js/releases/tag/v16.4.5) and published package identify the patch baseline. Node >=20.9.0, peers and MIT metadata are compatible. This is a same-package patch within the planned website update.

The project uses static export, unoptimized images and a force-static generated search endpoint. No announced affected request-time cache, Edge-cookie, proxy-response or image-optimizer surface is demonstrated here; that applicability assessment is an inference from the source and [static-export documentation](https://nextjs.org/docs/app/guides/static-exports). Advisory ranges are not identical: GHSA-rhh4-g77r-r9qj currently ends at 16.4.2 with no patched version listed despite appearing in the release list. Do not claim all seven affect 16.4.4 or claim independent security qualification.

## License audit

**Decision**: Record actual changed metadata and continue the offered EPL-2.0 option for elkjs; preserve the project's root license.

**Rationale**: Pflag's LICENSE and module directive bytes are unchanged. The npm proposals' other changed license fields remain MIT/MPL-2.0. Transitive elkjs 0.9.3 to 0.12.1 changes SPDX from EPL-2.0 to EPL-2.0 OR GPL-3.0-or-later. Published license bytes retain the EPL text and replace the secondary-license placeholder with GPL v3-or-later wording. Old LICENSE SHA256: 89591d4578fb1ebd91501312a3d25f021bd865a2e436641c1cf7b1bc7e3c1617. New: 637e81f4a1b6b4079535c499fca05e238cf7605c1ff5b76d60d7da7ce96700c9. This is an offered OR choice, not GPL-only metadata.

**Alternatives considered**: Calling all license bytes unchanged would be inaccurate. No existing maintained dependency-notice inventory was found; introducing a broad inventory subsystem is not warranted by this bounded update.

## Source and public delivery

**Decision**: Keep prepared v2.0.1, update source-bound release records, preserve current public v2.0.0 and carry #419 into the established post-merge controls.

**Rationale**: Navigation/diagnostic/policy fixes and compatible dependency maintenance do not change the public application contract or current-state formats. Same-day S133 competitor facts and existing qualifiers remain valid; comparison conclusions stay no-change with refreshed content hashes. The owner's explicit stop before final merge prevents tag/draft/promotion in this run.

## Fresh responsive observation

The complete S134 browser run reached the first-task route from all platforms, then observed 38 pixels of overflow during the live 390-pixel resize. Current-bundle animation-frame measurements were 47, 38, 0, 0, 0 and 0 pixels. The settled screenshot shows the intended mobile layout; long code remains inside local scroll containers. Preserve the one-pixel requirement with a condition-based Playwright assertion, zero test retries and no fixed sleep. Retain the original trace rather than replacing its failed outcome.
