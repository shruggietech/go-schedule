# Tasks: Remote Access Architecture

**Input**: Design documents from `specs/074-remote-access-architecture/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/remote-boundary.md`, and both checklists

**Tests**: Required because S074 adds an executable architecture-regression contract.

## Phase 1: Specification and Research

- [x] T001 [US1] Create and clarify the issue #165 specification in `specs/074-remote-access-architecture/spec.md`.
- [x] T002 [US1] Validate specification quality in `specs/074-remote-access-architecture/checklists/requirements.md`.
- [x] T003 [US1] Validate security, deployment, compatibility, and ownership coverage in `specs/074-remote-access-architecture/checklists/remote-security.md`.
- [x] T004 [US3] Record primary-source research and selected library ownership in `specs/074-remote-access-architecture/research.md`.
- [x] T005 [US3] Define future entity boundaries without adding schema in `specs/074-remote-access-architecture/data-model.md`.
- [x] T006 [US1] Define the ordered trust, operation, capability, failure, deployment, non-goal, and sequencing contract in `specs/074-remote-access-architecture/contracts/remote-boundary.md`.
- [x] T007 [US1] Document review and verification commands in `specs/074-remote-access-architecture/quickstart.md`.

## Phase 2: Analysis Gate

- [x] T008 Audit the specification, plan, research, data model, contract, quickstart, checklists, and task coverage through the spec-kit analyze gate; resolve every critical or high finding before implementation.

## Phase 3: Executable Architecture Contract

- [x] T009 [US1] Write `test/scripts/remote-architecture-check_test.sh` first with a passing fixture plus negative mutations for transport, local separation, every deployment mode, control ownership, threat verification, non-goals, and issue order.
- [x] T010 [US1] Demonstrate the focused fixture suite fails because `scripts/remote-architecture-check.sh` and the maintained architecture page do not yet exist.
- [x] T011 [US1] Implement `scripts/remote-architecture-check.sh` with deterministic aggregated diagnostics and UTF-8-safe POSIX behavior.
- [x] T012 [US1] Integrate the focused checker and its fixture suite into `scripts/docs-check.sh`.

## Phase 4: Maintained Architecture

- [x] T013 [US1] Author `docs/remote-access.md` with current-state warning, ordered trust boundary, allowlist record, capability and credential lifecycles, and failure behavior.
- [x] T014 [US2] Add all five deployment modes, explicit product and operator ownership, supported posture, and fail-closed configuration rules to `docs/remote-access.md`.
- [x] T015 [US3] Add API versioning, OpenAPI generation, SSE, retry, dependency maintenance, threat-to-test, non-goal, and issue-sequencing sections to `docs/remote-access.md`.
- [x] T016 [US1] Link the maintained boundary from `docs/architecture.md` and distinguish current local behavior from future remote behavior in `docs/api.md`.
- [x] T017 [US3] Record the dated S074 architecture decision in `CHANGELOG.md`.

## Phase 5: Focused Verification and Lifecycle

- [x] T018 Run the remote architecture fixture suite, documentation gate, specification lifecycle check, and GitHub publication formatter; resolve every failure.
- [x] T019 Confirm no runtime source, dependency manifest, listener configuration, or endpoint behavior changed and that issues #165 through #173 remain open.
- [x] T020 Advance the specification to `In Progress`, update the project item to S074 and In progress, and record focused evidence in `specs/074-remote-access-architecture/verification.md`.

## Phase 6: Canonical Verification and Delivery

- [x] T021 Run `sh scripts/verify.sh all` in the foreground with the repository's Node 26 baseline and resolve every failure.
- [x] T022 Mark every required task complete, advance the specification and inventory to `Implemented`, and record exact local delivery evidence.
- [x] T023 Run UTF-8, BOM, mojibake, Unicode em dash, Markdown wrapping, diff, and repository-status audits.
- [x] T024 Prepare the complete S074 slice for commit as `feat(074): define remote access architecture` with the required co-author trailer.

## Authorized Publication Runbook

Publication and review are workflow evidence rather than implementation work. After the verified commit, push the authorized branch, open a structured pull request that closes #165 without closing downstream issues, move the project item to PR review, and process every CI and review result. At most one manual second `@codex review` round may be requested. Do not request a third round. Once final-head CI is green and all review conversations are resolved, ask the maintainer to perform the final review and merge ritual.

## Dependencies and Execution Order

- Phase 1 completes the specify, clarify, checklist, and plan artifacts before task execution.
- T008 is the blocking analysis gate.
- T009 and T010 establish red evidence before T011 through T015 make the contract pass.
- T013 through T015 can be reviewed independently, but T016 depends on the maintained page existing.
- T018 through T024 require all implementation tasks.
- The publication runbook depends on a clean commit and canonical verification but does not affect the specification's implementation lifecycle state.

## Scope Guard

- Do not implement daemon identity, actor persistence, HTTPS listeners, pairing, remote clients, reconnection, or release qualification.
- Do not add `oapi-codegen`, `zalando/go-keyring`, or any other runtime or tool dependency in S074.
- Do not close #165 before merge and do not close, retitle, or reorder #166 through #173.
- Do not create a release, tag, package, or public remote-access claim.
