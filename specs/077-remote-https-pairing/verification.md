# Verification: Authenticated Remote Access and Pairing

## Analysis gate

- The first analysis pass found one HIGH scope mismatch: issue #169 requires a desktop enrollment journey and native credential storage, while the initial plan deferred all desktop work. The specification, plan, research, tasks, and implementation were corrected to include the bounded pairing form and protected-storage handoff; persistent profiles and target switching remain in #170.
- Enrollment requires the administrator-approved daemon identity, client display name, desktop kind, and capability. Phrase consumption and actor plus credential creation occur in one transaction.
- The runtime remote table is independent of the local mux, carries operation, capability, target, audit, retry, body-limit, and secret-exclusion classifications, and is checked against the committed OpenAPI source.
- No unresolved markers, constitutional violations, or remaining critical or high analysis findings remain.

## Focused verification

| Command | Result |
| --- | --- |
| `go test ./internal/... ./cmd/goschedd ./cmd/gosched` | PASS |
| `go test -race ./internal/enrollment ./internal/remote ./internal/store` | PASS |
| `go test ./...` from `desktop/` | PASS |
| `npm test -- --run` from `desktop/frontend/` | PASS, 24 files and 83 tests |
| `npm run build` from `desktop/frontend/` | PASS |
| `go vet ./...` | PASS |
| `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.0 run ./...` | PASS, zero issues |
| `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | PASS, zero called-symbol or imported-package vulnerabilities |
| `go tool oapi-codegen --config api/openapi/remote-v1.cfg.yaml api/openapi/remote-v1.yaml` | PASS, clean generated client boundary |
| `go run ./scripts/github-format` | PASS |
| `git diff --check` | PASS |

## Acceptance evidence

- Configuration tests prove the listener defaults disabled, requires a numeric address and certificate pair, and rejects unacknowledged wildcard or public exposure.
- A real loopback listener test proves trusted TLS 1.3 succeeds, TLS 1.2 fails, and cancellation completes bounded graceful shutdown.
- Route tests prove public health, isolated enrollment, protected bearer use, local-only route exclusion, and browser-origin rejection.
- Enrollment tests cover single use, concurrent exchange, wrong daemon, cancellation, expiry, five-attempt exhaustion, rotation, immediate old-token invalidation, and revocation.
- Schema v18 preserves the prior migration fixtures and stores only Argon2id verifier material and credential digests plus safe fingerprints.
- Native-storage tests prove daemon-scoped credential entries and pre-exchange failure when protected storage is unavailable. Desktop tests prove the phrase is cleared after the bounded handoff.
- Generated OpenAPI and runtime-table checks cover every reachable remote operation. Unsupported local secret and runtime administration paths are unreachable.
- The shared S076 authorization and intent-first audit suites continue to cover all four capability levels, denied side effects, actor reload, and protected-value exclusion.

## Canonical verification

The canonical eight-gate run passed after three gate-driven corrections: an unsupported `json` fence label was changed to `text`, implemented delivery evidence was rewritten in the lifecycle check's established review-branch form, and direct remote-storage failure tests restored the core persistence coverage floor after final fail-closed branches were added. Format, vet, lint, race, GUI, coverage, documentation, and automation are green. Final coverage measured engine 82.9 percent, schedule 89.1 percent, timezone 91.3 percent, store 80.2 percent, catchup 88.9 percent, and logbus 91.1 percent.

## First-round review remediation

- Source limiter entries now expire after ten minutes without traffic when capacity pressure requires reclamation. A deterministic clock test proves that a full limiter admits a new source after stale eviction while preserving active-source rejection.
- Desktop enrollment rejects every HTTP redirect so a 307 or 308 response cannot replay the one-time pairing phrase to another trusted certificate endpoint. The regression test trusts both test certificates and proves the redirect target receives no request.
- The OpenAPI source now assigns concrete request and success-response schemas to every JSON operation, including the one-time issued credential, and records the actual bodyless statuses for delete, enable, disable, run-now, and alert acknowledgement. The generated client now exposes typed operation payloads instead of arbitrary maps.

## Dependency and integrity review

- Direct runtime dependencies are `golang.org/x/crypto` v0.55.0, `golang.org/x/time` v0.15.0, `github.com/zalando/go-keyring` v0.2.8, and `github.com/oapi-codegen/runtime` v1.7.0. The generator is pinned as the Go tool `github.com/oapi-codegen/oapi-codegen/v2` v2.8.0.
- The selected versions retain the repository's Go 1.25 baseline. Standard library crypto owns randomness, SHA-256, constant-time comparison, certificates, and TLS.
- Govulncheck reports no reachable vulnerability. It identifies three module-only advisories in unused `x/crypto/ssh` and `x/crypto/openpgp` packages; S077 imports only `x/crypto/argon2`, and the SSH fixes require the Go 1.26-only v0.56.0 release.
- Repository publication formatting, diff whitespace, generated-contract, secret-name, em-dash, and mojibake scans are clean.
