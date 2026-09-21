# Tasks: Remote MCP Authorization

## Phase 1: Specification and contracts

- [x] T001 Define enablement, OAuth discovery, resource binding, credential isolation, scope, lifecycle, audit, protocol, and deployment requirements.
- [x] T002 Record configuration, access-grant, token endpoint, and MCP resource contracts.

## Phase 2: Authorization and transport

- [x] T003 Add failing tests for independent enablement, configuration validation, metadata, client authentication, resource binding, scope bounds, and pairing-phrase rejection.
- [x] T004 Implement a bounded memory-only access-grant service backed by existing MCP credential and actor lifecycle validation.
- [x] T005 Implement SDK-backed protected-resource metadata, authorization-server metadata, client-credentials token exchange, and bearer middleware.
- [x] T006 Mount remote MCP on the existing TLS listener behind shared host, origin, request, concurrency, source, and actor limits.

## Phase 3: Shared authority and protocol behavior

- [x] T007 Add actor-bound in-process API clients and compose existing Observe, Operate, and Manage MCP surfaces without a surrogate actor.
- [x] T008 Prove authority discovery, audit attribution, redaction, deduplication, uncertainty, revocation, rotation, expiry, wrong target, and failure behavior.
- [x] T009 Verify current and compatibility MCP protocol revisions with the official SDK and reject unsupported versions.

## Phase 4: Documentation and delivery

- [x] T010 Document direct, private-network, and reverse-proxy setup, provisioning, security, recovery, and client configuration.
- [x] T011 Update changelog, specification inventory, and agent context.
- [x] T012 Run focused race tests, deployment matrices, canonical verification, and GitHub publication formatting checks.
- [x] T013 Record final evidence and resolve the specification lifecycle state.
