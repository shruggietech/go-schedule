# Tasks: Actor Permissions and Management Audit

**Input**: Design documents from `specs/076-actor-permissions-audit/`

**Tests**: Required for every persisted, authorization, API, client, and CLI behavior in issue #167.

## Phase 1: Specification and Design

- [x] T001 Specify and clarify issue #167 in `spec.md`.
- [x] T002 Validate requirements and security boundaries in both checklists.
- [x] T003 Record capability, local actor, catalog, audit, redaction, retention, and administration decisions in `research.md`.
- [x] T004 Define persistence entities and lifecycle in `data-model.md`.
- [x] T005 Define authorization and audit contracts in `contracts/authorization-audit.md`.
- [x] T006 Record the implementation plan and operator verification flow.

## Phase 2: Analysis Gate

- [x] T007 Analyze all specification artifacts for ambiguity, duplication, missing issue criteria, constitutional conflicts, and scope leakage; resolve every critical or high finding.

## Phase 3: Domain and Persistence Foundation

- [x] T008 Write failing domain tests for capabilities, actor kinds and states, display names, expiration, and audit enums.
- [x] T009 Write failing schema v17 migration and built-in actor initialization tests that preserve prior scheduler and daemon identity state.
- [x] T010 Write failing actor lifecycle tests for create, list, update, revoke, built-in protection, reload, and contextual errors.
- [x] T011 Write failing audit tests for intent, completion, denial, deterministic filters, export ordering, redaction shape, and count plus age pruning.
- [x] T012 Implement domain models, schema v17, actor persistence, and audit persistence.
- [x] T013 Run focused domain and store tests and resolve every failure.

## Phase 4: Authorization and API Enforcement

- [x] T014 Write failing hierarchy, unknown-operation, expired-actor, revoked-actor, and per-request reload authorization tests.
- [x] T015 Write failing catalog tests proving every registered management route has one classification.
- [x] T016 Write failing server tests for local actor compatibility, allowed and denied operations, intent-first behavior, result completion, and audit failure blocking.
- [x] T017 Write failing actor and audit endpoint tests for validation, lifecycle protection, filtering, and newline-delimited export.
- [x] T018 Implement the shared operation catalog, authorizer, middleware, and actor plus audit handlers.
- [x] T019 Run focused authorization and server tests and resolve every failure.

## Phase 5: Typed Client and CLI

- [x] T020 Write failing client tests for actor lifecycle, audit filters, structured errors, and deterministic export.
- [x] T021 Write failing CLI tests for actor list, create, update, revoke and audit list, export in human and JSON forms.
- [x] T022 Implement typed client access operations and `gosched actor` plus `gosched audit` command families.
- [x] T023 Run focused client and CLI tests and resolve every failure.

## Phase 6: Documentation and Verification

- [x] T024 Document capability mapping, local trust, actor lifecycle, audit schema, redaction, retention, filters, export, migration, and future transport integration.
- [x] T025 Update API, CLI, architecture, changelog, and specification inventory references.
- [x] T026 Record test-first and focused evidence in `verification.md`, advance the spec to In Progress, and update issue #167 project fields to S076 and In progress.
- [x] T027 Run focused tests, catalog coverage, spec lifecycle, docs, publication format, diff, encoding, mojibake, and scope checks.

## Phase 7: Canonical Verification and Delivery

- [x] T028 Run `bash scripts/verify.sh all` and resolve every failure.
- [x] T029 Mark all tasks complete, advance the specification and inventory to Implemented, and record exact evidence.
- [x] T030 Prepare the verified commit as `feat(076): add actor permissions and audit` with the required co-author trailer.

## Authorized Publication Runbook

Push the authorized branch, open a structured pull request closing #167, move the project item to PR review, and process every CI and review result. Request at most one manual second `@Codex` review round. Do not request a third round. Stop for the maintainer's final review and merge ritual.

## Scope Guard

- Do not close #167 before merge.
- Do not add remote listeners, TLS, pairing, credentials, keyring persistence, local login, MCP mutations, or general policy configuration.
- Do not store sensitive request or execution content in audit records.
- Do not merge the pull request; stop for the maintainer's final review and merge ritual.
