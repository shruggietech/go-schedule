# Tasks: MCP Manage Authority

## Phase 1: Specification and contracts

- [x] T001 Define authority, family tools, secret boundaries, confirmation, retry, audit, and failure requirements in `specs/092-mcp-manage-authority/`.
- [x] T002 Record request, result, session policy, and discovery contracts.

## Phase 2: Manage executor

- [x] T003 Add failing tests for discovery isolation, family actions, exact daemon identity, confirmation, redaction, deduplication, denial, and uncertainty in `internal/mcpmanage/`.
- [x] T004 Implement six typed family tools that delegate to the existing API client and return redacted results in `internal/mcpmanage/`.
- [x] T005 Prove hostile task content and authority-shaped definition content cannot raise permissions or leak through results.

## Phase 3: Session and transport integration

- [x] T006 Permit explicit Manage runtime sessions while retaining Observe defaults and live revocation in `internal/mcpsession/`.
- [x] T007 Add Manage and optional confirmation flags to stdio and localhost HTTP in `internal/cli/` and `internal/mcphttp/`.
- [x] T008 Add Manage selection and confirmation policy to the desktop Agent Access projection and frontend.
- [x] T009 Prove Observe and Operate cannot discover or invoke Manage tools across stdio and HTTP.

## Phase 4: Documentation and delivery

- [x] T010 Update MCP guidance, changelog, specification index, and agent context.
- [x] T011 Run focused tests and all canonical verification gates.
- [x] T012 Record final evidence and run GitHub publication formatting checks before publishing S092.
