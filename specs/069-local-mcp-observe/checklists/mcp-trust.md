# MCP Trust Boundary Checklist

**Purpose**: Validate the Observe authority, redaction, untrusted-content, transport, and future-capability boundary before implementation and publication

**Created**: 2026-09-08

**Feature**: [spec.md](../spec.md)

## Authority

- [x] Observe is the only enabled permission class
- [x] No tool or mutation capability is approved
- [x] Future Operate and Manage requirements are documented without enabling them
- [x] Existing daemon IPC authorization remains authoritative
- [x] No fallback identity, transport, database access, or elevation is allowed

## Information Boundary

- [x] Dedicated allowlisted response types are required
- [x] Commands, arguments, environment, stdin, working directories, run-as identities, trigger keys, notification credentials, and raw schedules are excluded
- [x] Output and user-controlled text have explicit byte limits and truncation metadata
- [x] Untrusted fields and resource envelopes have an explicit data-only trust notice
- [x] Host-visible errors use a bounded safe vocabulary

## Transport and Lifecycle

- [x] Stdio is the only enabled MCP transport
- [x] No TCP listener or background service is introduced
- [x] Protocol stdout and diagnostic stderr are separated
- [x] Host disconnect, cancellation, deadlines, and clean shutdown are specified
- [x] Cross-platform subprocess verification requires hidden Windows child-process creation

## Verification

- [x] Secret canaries cover every prohibited category
- [x] Discovery proves zero tools and mutation capabilities
- [x] Pagination and size limits have measurable gates
- [x] IPC denial and unavailable-daemon behavior have tests
- [x] Official SDK negotiation and supported revisions have protocol tests
