# Tasks: MCP Operate Authority

## Phase 1: Specification and contracts

- [x] T001 Define Observe and Operate discovery, authority, target, retry, audit, and revocation requirements in `specs/091-mcp-operate-authority/`.
- [x] T002 Record the runtime session, tool input, result, and deduplication contracts.

## Phase 2: Runtime identity foundation

- [x] T003 Add a memory-only MCP session registry with actor creation, digest-only secrets, resolution, and revocation.
- [x] T004 Add protected local API and client contracts for creating, using, and revoking runtime MCP sessions.
- [x] T005 Prove revocation, expiry, malformed secrets, and actor attribution through unit and API tests.

## Phase 3: Operate tools

- [x] T006 Add typed run-now, enable, and disable tools with exact daemon, task, and request identity.
- [x] T007 Add bounded deduplication and explicit accepted, rejected, denied, and uncertain results.
- [x] T008 Change only task enable and disable authorization from Manage to Operate and retain Manage for definition mutations and deletion.
- [x] T009 Add conformance, hostile-content, wrong-target, readiness, retry, timeout, and authority-boundary tests.

## Phase 4: Transport integration

- [x] T010 Add explicit stdio permission and client-name flags while retaining Observe defaults.
- [x] T011 Extend localhost HTTP lifecycle state, CLI, and desktop projection with explicit permission and session revocation.
- [x] T012 Prove credential rotation, listener disable, process shutdown, and existing-connection revocation behavior.

## Phase 5: Documentation and delivery

- [x] T013 Update local MCP guidance, changelog, specification index, and agent context.
- [x] T014 Run focused tests and all canonical verification gates.
- [ ] T015 Record final evidence, run GitHub publication formatting checks, and publish the reviewed S091 pull request.
