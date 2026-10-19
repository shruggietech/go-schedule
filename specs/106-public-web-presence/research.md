# S106 research and decisions

## Static content and hosting

Decision: static Next.js with Fumadocs; canonical `docs/*.md` generates ignored `.md` inputs. Rationale: fragcap's static implementation demonstrates source loader, catch-all routes and search; Markdown mode preserves literal CLI placeholders. Alternatives: retaining Jekyll misses the approved migration; a custom renderer repeats solved work; copying sibling Worker hosting adds unrelated changes.

References: https://nextjs.org/docs/app/guides/static-exports, https://www.fumadocs.dev/docs/deploying/static, https://www.fumadocs.dev/docs/headless/search/orama and sibling `A:/Code/fragcap/site`. Next's official registry version is 16.4.2; Fumadocs core/UI 16.17.6 and MDX 15.4.7 were validated together by the production export after the older sibling pins failed Markdown compilation.

## Documentation continuity

Decision: source-to-route inventory, link rewriting outside code, and static legacy HTML aliases preserving query/hash. Pages cannot execute runtime redirects. Copy `docs/assets` and derive navigation from existing metadata. Render each original source H1 once, preserving its fragment ID rather than substituting a second generated title. Project-authored modified Mermaid uses TB; do not alter user-owned literal examples.

## Competitive evidence

Decision: cells identify status, explanation, official source, edition/deployment and date. Unknown replaces unsupported absence claims. The seven official sources are documented in #282. JAMS also advertises MCP, so no uniqueness claim is warranted. Keep published v1.5.1 visible until a newer public GitHub release has a matching snapshot; candidate data must not advance public claims.

## Release identity

Decision: exact tag/source ref, reviewed dataset/release-content digests and every capability's decision; bind the resolved commit at runtime using existing workflow identity checks. Embedding a candidate SHA in its own committed review is self-referential. No arbitrary review expiry; changed tag/content invalidates review.

## Authentication and publication

Decision: request owner authentication for read-only Cloudflare inventory, never persist tokens in repository artifacts, prepare record-level changes/rollback and defer mutation to pre-publication review. DNS-only apex A/AAAA and www CNAME target GitHub; retain ownership TXT and unrelated DNS. Official references and owner authority are in #278.

Final evidence refinement: the vendor-linked JAMS documentation establishes failover and ACLs; Control-M 9.0.22 documents self-hosted failover and a preview MCP server with explicit limitations; OpCon documents REST-backed CLI access and VisualCron/MFT-qualified event/dependency integration. Those cells link the exact official pages rather than treating overview omissions as absence. Product/category filters retain the go-schedule baseline while reducing grid width.
