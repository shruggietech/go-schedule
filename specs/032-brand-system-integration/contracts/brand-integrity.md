# Brand Integrity Command Contract

## Invocation

From repository root:

```text
go run ./scripts/brand-check
```

The existing automation gate invokes the same command. It is offline and uses only the Go standard library.

## Inputs

- `brand/manifest.json`: canonical standalone-kit inventory and SHA-256 values.
- `brand/repository-consumers.json`: exact canonical-to-consumer copy mappings.
- Repository filesystem rooted at the current checkout.

All manifest and mapping paths are slash-separated, relative, clean, and forbidden from escaping their declared root.

## Success output

One concise line on standard output:

```text
brand-check: OK - <artifacts> artifacts, <svgs> SVGs, <consumers> consumers
```

Exit status is `0`.

## Failure output

Each finding is printed to standard error with an actionable path, followed by a summary:

```text
brand/manifest.json: missing artifact: logos/svg/example.svg
docs/assets/brand/example.svg: differs from canonical brand/logos/svg/example.svg
brand-check: FAILED with 2 issue(s)
```

Exit status is `1` for contract failures and `2` for invalid invocation or an unreadable repository root.

## Required checks

1. Parse and validate both JSON control files.
2. Reject absolute, unclean, duplicate, or escaping paths.
3. Require every manifest artifact and match its byte count and SHA-256 digest.
4. Validate every textual artifact as UTF-8 without BOM and scan for common mojibake signatures.
5. Require accessible SVG titles and reject live text or font-family dependencies.
6. Require every mapped source and target and compare their complete bytes.
7. Require the standalone report and repository guidance files.
8. Accumulate independent findings rather than stopping after the first drift.

## Test contract

Tests create isolated fixture repositories and prove at least these negative cases:

- missing manifest artifact;
- hash or length mismatch;
- malformed or escaping manifest path;
- UTF-8 BOM and invalid UTF-8;
- SVG missing a title;
- SVG containing live text or a font-family declaration;
- missing consumer;
- consumer bytes different from canonical;
- malformed consumer map;
- complete valid repository.
