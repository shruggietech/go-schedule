# S106 local validation

Install the exact website lockfile dependencies in `site/` with `npm ci`, then run `npm test` and `npm run build`. Serve `site/out` locally for production browser checks.

Inspect `/`, `/docs/`, `/compare/`, platform installation and `/cli.html#global-flags`. Search a CLI term, change appearance, reload and use a narrow viewport. Search needs no runtime endpoint.

Run `go test ./scripts/comparison-check` and structural `go run ./scripts/comparison-check`. Future releases need matching snapshot/review content before staging. Run `sh scripts/verify.sh all` in the foreground and retain gate outcomes. Inspect six-issue delivery status and pending domain actions before owner review; publish nothing.
