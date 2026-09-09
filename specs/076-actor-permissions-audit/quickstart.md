# Quickstart: Actor Permissions and Management Audit

## Focused verification

1. Run store migration, actor lifecycle, retention, and audit intent tests.
2. Run authorization catalog, server middleware, local compatibility, client, and CLI tests.
3. Confirm every registered management route is represented in the operation catalog.
4. Create and revoke a non-built-in actor, then verify a subsequent authorization attempt is denied.
5. Exercise successful, failed, denied, and interrupted audit outcomes and inspect redacted fields.
6. Export filtered audit events and verify deterministic newline-delimited JSON ordering.

## Full verification

Run `bash scripts/verify.sh all` from the repository root. Then run `go run ./scripts/github-format`, `git diff --check`, encoding checks, specification lifecycle checks, and a final scope audit.

## Scope confirmation

Confirm no TCP listener, TLS certificate, pairing phrase, bearer credential, keyring persistence, local login, or MCP mutation endpoint was added. S076 prepares the actor and audit substrate consumed by #168 and #169.
