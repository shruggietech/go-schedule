# Quickstart: Validate Authenticated Remote Access and Pairing

## Prerequisites

- A clean checkout with the supported Go, Node, POSIX shell, and native build prerequisites.
- A temporary trusted test certificate and a disposable data directory.
- No production credential or certificate material.

## Scenarios

1. Run `sh scripts/verify.sh all` and require all eight gates.
2. Start with default configuration and confirm only the existing local IPC endpoint exists.
3. Start with invalid or unacknowledged wildcard remote configuration and confirm failure occurs before bind.
4. Start with an exact loopback address and test certificate, negotiate TLS 1.3, and confirm plaintext plus browser-origin requests fail.
5. Create one pairing phrase over local IPC, exchange it once from the desktop pairing form with the matching daemon identity and trusted certificate, verify native protected storage, and confirm replay fails generically.
6. Use the issued credential for Observe, Operate, and Manage examples, then prove capability denial produces no side effect.
7. Rotate the credential locally and confirm the old value fails on the next request; revoke the replacement and confirm it also fails.
8. Exercise expiry, cancellation, attempt exhaustion, concurrent exchange, body limit, source limit, actor limit, and secret-canary tests.
9. Compare OpenAPI operation IDs with the runtime remote table and confirm excluded local paths remain unreachable.

## Expected Outcome

The daemon exposes no network port by default. Explicit TLS configuration enables only the reviewed `/api/v1` surface. One short-lived phrase establishes one named actor and one opaque credential, current server-owned authority gates every request, and no secret reaches durable or diagnostic output.
