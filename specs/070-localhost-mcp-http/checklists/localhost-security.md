# Localhost MCP Security Checklist

**Purpose**: Validate the network, browser, credential, authority, and lifecycle boundary before implementation and publication

**Created**: 2026-09-08

**Feature**: [spec.md](../spec.md)

## Network Boundary

- [x] Numeric IPv4 loopback is the only permitted bind address
- [x] Callers cannot provide a host, interface, wildcard, or forwarded address
- [x] Exact Host validation is required before MCP dispatch
- [x] SDK localhost protection remains enabled
- [x] Port conflicts leave lifecycle state unchanged

## Browser and Credential Boundary

- [x] Every request requires a separately provisioned bearer credential
- [x] Credentials are random, runtime-only, returned once, and never logged
- [x] Rotation, disablement, and restart immediately stale prior credentials
- [x] Origin-bearing requests require an exact explicit loopback origin
- [x] Missing, malformed, duplicate, wildcard, opaque, and non-loopback cases fail closed

## Authority and Lifecycle

- [x] HTTP reuses the exact Observe server and advertises zero tools
- [x] Local IPC, GUI, CLI, stdio, and remote settings remain independent
- [x] Disablement and daemon shutdown revoke credentials and bound graceful shutdown
- [x] Concurrent lifecycle operations and in-flight requests are specified
- [x] No endpoint, origin, credential, or enablement state is persisted

## Verification

- [x] DNS-rebinding, Origin, and authorization matrices are required
- [x] Official SDK revision and stdio parity tests are required
- [x] Port conflict, cancellation, request-size, and shutdown tests are required
- [x] Race verification and canonical repository verification are required
