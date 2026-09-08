# Verification: Local Observe-Only MCP

**Branch**: `codex/069-local-mcp-observe`

**Date**: 2026-09-08

## Focused evidence

- `go test -race ./internal/mcpobserve ./internal/cli ./cmd/gosched` passed. This covers allowlisted mapping, prohibited-field canaries, UTF-8 bounds, output and message truncation, deterministic pagination, opaque cursor validation, bounded safe errors, deadlines, cancellation, SDK discovery, the `2026-07-28` protocol revision, real stdio initialization at `2025-11-25`, zero tool capability, all five resources, all four continuation templates, hidden Windows subprocess creation, protocol-only stdout, and host-disconnect shutdown.
- `go test ./...` passed across every root Go package and integration package before the canonical run.
- `go vet ./internal/mcpobserve ./internal/cli ./cmd/gosched` passed.
- `go run ./scripts/github-format` passed with no Unicode em dash or hard-wrapped Markdown prose.

## Dependency evidence

- `github.com/modelcontextprotocol/go-sdk v1.7.0` is the official SDK selected by issue #161 and is pinned in `go.mod` and `go.sum`.
- The SDK license file states its Apache-2.0 and MIT transition terms. New transitive modules `github.com/google/jsonschema-go`, `github.com/segmentio/asm`, `github.com/segmentio/encoding`, `github.com/yosida95/uritemplate/v3`, `golang.org/x/oauth2`, `golang.org/x/sync`, and `golang.org/x/time` each include a permissive MIT or BSD-style license file in the resolved module cache.
- No standard-library-only implementation provides the required official MCP interoperability, so the dependency is justified by the protocol conformance and compatibility requirements.

## Canonical verification

`C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` passed in the foreground on 2026-09-08 after one formatting-only import-group correction found by the first lint attempt.

| Gate | Evidence |
| --- | --- |
| format | `gofmt` and `github-format` clean. |
| vet | Root `go vet ./...` passed. |
| lint | golangci-lint v2.12.0 reported zero issues. |
| race | Every root package and the integration suite passed under the race detector; `test/integration` completed in 58.825 seconds. |
| gui | Desktop Go packages, native Windows Wails production build, 20 Vitest files with 73 tests, TypeScript, and Vite production bundle passed. |
| coverage | engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.4%, catchup 88.9%, logbus 91.1%. |
| docs | Product policy, fixtures, 17 pages, links, front matter, fences, theme, and product copy passed. |
| automation | Workflow, CodeQL, Dependabot, release, brand, lifecycle, eight-gate, and automation fixture checks passed. |

## Boundary audit

- `internal/mcpobserve` imports no network-listener, persistence, executor, scheduler, notification, trigger, or elevation package. Its injected interface contains four read methods only.
- Five static resources and four continuation templates are registered. No tool, prompt, subscription, TCP transport, or HTTP transport is registered.
- Dedicated MCP response types contain no command, arguments, environment, stdin, working directory, run-as identity, trigger key, notification endpoint, authorization, raw schedule, IPC path, or filesystem path field.
- A hostile `SECRET_CANARY` fixture occupies every prohibited daemon task and provenance field; encoded MCP resource assertions find zero matches.
- User-controlled display fields and output are JSON strings, byte-bounded, truncation-aware, and named in each envelope's untrusted-field metadata beside the fixed data-only trust notice.
- Git diff checks, UTF-8 and BOM checks in the format gate, publication formatting, and mojibake inspection are clean. The only secret-canary text in changed files is the deliberate test fixture and its verification description.

## Result

S069 satisfies issues #161 and #162 locally. Hosted CI and third-party Codex review remain publication evidence on the pull request and do not change the implemented specification state.
