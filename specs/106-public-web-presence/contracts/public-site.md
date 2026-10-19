# Public site and release contracts

- `/`, `/docs/` and `/compare/` are the home, docs and comparison. Include static 404, sitemap, robots and canonical metadata.
- `/docs/<slug>/` serves canonical Markdown; installation pages use `/docs/install/<platform>/`. Legacy case-sensitive `/<source>.html` aliases preserve query/hash.
- `/static.json` is a build-produced browser search index.
- Appearance respects system preference, explicit light/dark choice, keyboard use and storage failure.
- `site/data/comparisons/<tag>.json` is selected by deployment's latest-public-release tag. Local preview uses the recorded published baseline.
- `go run ./scripts/comparison-check` validates structure; release mode takes `-tag` and `-commit`, requires matching review/digests and names corrective action on failure. Staging and both promotions invoke it before release mutation.
- No source preparation operation authorizes live DNS/Pages writes, deployment, push, PR, tag or release during this build phase.
