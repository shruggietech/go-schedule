# Verification: Authenticated Localhost MCP

**Branch**: `codex/070-localhost-mcp-http`

**Date**: 2026-09-08

**Issue**: [#163](https://github.com/shruggietech/go-schedule/issues/163)

## Spec-kit analysis

The specification, plan, research, data model, contract, checklists, and tasks were audited for coverage and consistency. One mismatch was found and resolved: the first protocol test proved HTTP discovery and a resource read but did not compare the transport contract with stdio. The completed test compares official HTTP and in-memory stdio discovery plus normalized resource payloads from the same Observe fixture. Hosted Linux verification then exposed a listener startup race in which disable could reach `http.Server.Shutdown` before the Serve goroutine registered its listener; disable now closes the manager-owned listener immediately after credential revocation and before graceful HTTP shutdown, and the existing listener-closure test is the regression proof. First-round Codex review identified missing narrow CORS responses for explicitly allowed browser origins and a stale credential policy snapshot; exact-origin preflight and success responses plus one locked current-policy authorization point now cover both findings. Second-round Codex review identified browser default-port serialization and an expected wrapped listener-closure diagnostic; canonical origin normalization and explicit `net.ErrClosed` suppression now cover both findings. No unresolved clarification, requirement, constitution, or traceability finding remains.

## Focused evidence

`go test ./...` passed every root module package, including the real-loopback security and protocol integration suite.

`go test -race ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli ./cmd/goschedd` passed. The manager suite exercises default-off state, origin normalization, occupied ports, repeated enable, credential rotation, idempotent disable, listener closure, concurrent status and rotation, exact Host and Origin matrices, malformed and duplicate authorization, stale credentials, official SDK discovery, zero tools, both supported revisions, stdio parity, request-body bounds, and cancellation propagation.

## Canonical verification

`sh scripts/verify.sh all` passed in the foreground on the completed review-branch tree:

- `format`: repository-authored Markdown and GitHub text contain no Unicode em dash or hard-wrapped prose defects.
- `vet`: passed.
- `lint`: passed with zero issues.
- `race`: passed across daemon, CLI, core, scripts, and integration packages, including the new localhost MCP packages.
- `gui`: desktop Go suites, native Wails Windows build, all 73 frontend tests, and production frontend bundle passed.
- `coverage`: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.7%, catchup 88.9%, and logbus 91.1%.
- `docs`: documentation policy, fixtures, links, front matter, fences, theme, and product policy passed across 17 pages.
- `automation`: workflow, CodeQL, Dependabot, release, brand, lifecycle, eight-gate, and fixture audits passed.

## Security and publication audit

Changed content passed `git diff --check`, the repository GitHub formatter, UTF-8 and mojibake scans, secret-boundary inspection, spec lifecycle validation, and issue #163 traceability review. No credential fixture value is committed as a production credential, no plaintext active credential is retained by the manager, and no pinned artifact changed.
