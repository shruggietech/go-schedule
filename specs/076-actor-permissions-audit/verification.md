# Verification: Actor Permissions and Management Audit

## Analysis gate

- The specification, requirements checklist, security checklist, research decisions, data model, contract, plan, tasks, and quickstart consistently implement issue #167.
- The four capability levels and five actor kinds match the approved S074 architecture.
- The audit model covers intent, completion, denial, interruption, retention, redaction, filtering, migration, and export.
- Remote transport, credentials, pairing, local login, and MCP mutation authority remain excluded for #168 and #169.
- No unresolved markers, constitutional violations, or critical or high analysis findings remain.

## Test-first evidence

- The first focused compile after persistence and middleware scaffolding failed because an unused server import remained and prior migration tests still asserted schema v16. Those failures were resolved by removing the import and updating the migration expectations to v17.
- The first focused server run after capability publication failed the existing bounded manifest contract because the capability count had not been advanced for actor authorization and management audit. The assertion now requires both new capability names and the expanded deterministic set.
- New domain, catalog, store, server, client, and CLI acceptance tests cover the S076 behavior before canonical verification.

## Focused verification

| Command | Result |
|---|---|
| `go test ./...` | PASS |
| `go test ./...` from `desktop/` | PASS |
| `go test -race ./internal/domain ./internal/authorization ./internal/store ./internal/api/server ./internal/api/client ./internal/cli ./internal/mcpobserve` | PASS |
| `go vet ./...` | PASS |
| `bash scripts/spec-lifecycle-check.sh .` | PASS, 76 specifications lifecycle-consistent |
| `bash scripts/docs-check.sh` | PASS, 20 pages plus policy and remote-architecture fixtures |
| `go run ./scripts/github-format` | PASS |
| `git diff --check` | PASS |

## Acceptance coverage

- Domain tests exercise all four capability levels in both directions and reject unknown capability, kind, state, and audit values.
- Schema v17 initializes exactly one protected local actor while existing v15 and v16 upgrade tests preserve prior state.
- Actor store tests cover normalization, creation, update, irreversible revocation, built-in protection, and deterministic listing.
- Audit tests cover uncertain intent, successful completion, denial, actor and result filters, deterministic ordering, age pruning, and the 10,000-event limit.
- Catalog tests reject duplicate routes and identifiers, validate every definition, prefer static paths over parameterized paths, and fail closed for unknown operations or invalid actors.
- A source-backed router completeness test proves every registered management route has exactly one catalog entry and every catalog entry has a registered route.
- Server tests prove local compatibility, request-by-request revocation, denied evidence, successful actor mutation evidence, protected local identity, newline-delimited export, and structural exclusion of sensitive fields.
- Typed-client and CLI tests exercise actor lifecycle, filters, listing, export, and both human and JSON-compatible data paths.

## Canonical verification

`bash scripts/verify.sh all` passed on 2026-09-09 after correcting one ineffectual test assignment and adding actor-store branch coverage to restore the 80 percent core-package floor. All eight gates passed: format, vet, lint, race, GUI, coverage, documentation, and automation. Store coverage measured 80.7 percent.

## Review remediation

- The first CI run and its retries reproduced an external Google Chrome apt index hash mismatch before Linux desktop or Playwright project commands ran. The affected ephemeral jobs now remove both the legacy `.list` and deb822 `.sources` forms of that unused hosted-runner source before resolving Ubuntu and Playwright prerequisites, preserving package-integrity checks without blind retries.
- The automatic first Codex review failed to execute without findings. The one authorized manual second review was requested with `@codex review`; no further review round will be triggered.
