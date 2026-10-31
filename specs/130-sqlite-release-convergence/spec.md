# Feature Specification: SQLite convergence and release delivery

**Feature Branch**: `codex/130-sqlite-release-convergence`

**Created**: 2026-10-31

**Status**: Implemented

**Delivery**: Implemented on review branch codex/130-sqlite-release-convergence; all eight foreground canonical gates passed in 1464.688 seconds on 2026-10-31, with focused red/green, paired benchmarks and final documentation/integrity evidence in [verification.md](verification.md). Source publication/review is authorized; owner merge and actual release delivery remain separate.

**Input**: Owner-authorized S130 follows the post-S129 recommendation: adopt the now-published corrected storage dependency coherently, preserve current behavior, and continue the existing v2.0.0 release outcome under #412. Traceability: #404 and proposal #400; source baseline a58651724bfe4a9582d52464273e9825943771c7.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build one coherent storage candidate (Priority: P1)

As the maintainer, I can build the daemon, CLI, desktop and remote-control clients from one selected corrected storage dependency set without any build rewriting dependency metadata.

**Why this priority**: The held dependency proposal updates only part of the product and fails desktop/remote builds. A single coherent candidate resolves #404 and makes the upstream correction usable.

**Independent Test**: Resolve and build both application modules with dependency writes disabled; compare their selected storage/runtime pair and unchanged metadata.

**Acceptance Scenarios**:

1. **Given** the selected published correction and its required runtime dependency, **When** either application module is built normally, **Then** it resolves the same compatible storage pair without rewriting its manifests/checksums.
2. **Given** a root-only dependency update with stale desktop metadata, **When** the candidate is assessed, **Then** the inconsistency is detected before publication instead of being repaired implicitly during the build.

### User Story 2 - Preserve current tasks and bounded cancellation (Priority: P1)

As an operator of a current-format development installation, I retain tasks and audit history across restart, and a cancelled operation waiting for a storage lock returns control without waiting through the full busy timeout.

**Why this priority**: Dependency convergence must preserve the scheduler's current persistence and concurrency contracts while adopting the corrected driver's prepared-statement behavior.

**Independent Test**: Exercise current store creation, reopen, refusal, prepared operations, concurrent lock cancellation and durable audit boundaries on the selected candidate.

**Acceptance Scenarios**:

1. **Given** valid current state, **When** tasks and audit records are written and the store is reopened, **Then** committed records remain intact and prepared queries return the intended rows repeatedly.
2. **Given** a lock-blocked operation, **When** its supplied context is cancelled, **Then** the caller receives the context error promptly, without a partial committed change.
3. **Given** obsolete or malformed persisted state, **When** it is opened, **Then** existing refusal/reset guidance remains and opening neither converts nor deletes it.

### User Story 3 - Continue the real release outcome honestly (Priority: P2)

As the owner, I receive one reviewed source candidate and a concrete continuation of the existing v2.0.0 release workflow, with source readiness distinguished from actual package publication.

**Why this priority**: S129 already supplies release machinery and candidate identity. This dependency implementation should advance that outcome without another release framework or premature public claims.

**Independent Test**: Read the release continuation and current public discovery: the candidate remains v2.0.0, actual public latest stays v1.7.0 until approved promotion, and #412 retains actual publication work.

**Acceptance Scenarios**:

1. **Given** owner-merged corrected source, **When** release staging is considered, **Then** the existing exact-main, immutable tag/draft and concrete-asset promotion decisions apply to that source in order.
2. **Given** source PR readiness only, **When** issue and public status are reported, **Then** #404 is eligible for closure on functional integration, #400 is superseded only after replacement integration, and #412 stays open for actual public delivery/discovery.

### Edge Cases

- A published upstream correction is a candidate, not proof of adopted product behavior; a failed candidate retains its first failure and does not silently fall back or weaken checks.
- A cancelled write or commit must not turn into a successful partial mutation; lock-wait tests distinguish caller cancellation from an unrelated SQL error.
- The local root-module replacement and major-version namespace must remain intact in the desktop module.
- No users means no migrations, historical upgrade matrix or automatic persisted-state conversion/deletion. Current correctness, atomicity, access security, repair and removal still apply.
- Original #355/#356 causes remain unconfirmed and clustered execution remains explicitly owner-deferred.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Both application modules MUST describe one compatible selected published storage/runtime dependency pair and complete required transitive metadata.
- **FR-002**: Normal desktop and remote builds MUST consume that graph without requesting or performing manifest/checksum rewrites, and current build guidance MUST use the module-pinned tool.
- **FR-003**: Root and desktop module identity, local replacement and public interface versions MUST remain consistent with the prepared major release.
- **FR-004**: Current-state direct initialization, atomic persistence, durable reopen and audit records MUST retain their existing behavior.
- **FR-005**: Unsupported or malformed state MUST retain explicit refusal/reset guidance without automatic conversion or deletion.
- **FR-006**: Repeated prepared operations MUST preserve result correctness and MUST NOT introduce a storage performance regression above the constitution's 10 percent threshold without measured, recorded justification.
- **FR-007**: Lock-wait cancellation MUST preserve context-error attribution and transaction atomicity, with no new goroutine lifecycle leak.
- **FR-008**: Existing access controls, secret handling and local/remote security behavior MUST remain unchanged; their tests remain required.
- **FR-009**: The source handoff MUST continue existing #412/S129 v2.0.0 delivery through existing workflows, preserving separate exact tag/draft and concrete-asset promotion approval and accepted asset bytes.
- **FR-010**: Public latest selection and immutable historical release/comparison/evidence records MUST remain unchanged during source preparation.
- **FR-011**: Issue records MUST distinguish #404's functional correction, #400's proposal disposition and #412's still-incomplete actual release outcome.
- **FR-012**: The source PR MUST receive at most two Codex review requests total, including automatic PR-open review; every genuine finding is addressed within the authorized PR and reviewed heads remain accurately attributed.

### Key Entities *(include if feature involves data)*

- **Dependency candidate**: Published storage correction, required runtime pair and complete root/desktop dependency metadata.
- **Current persisted state**: Existing application/helper formats with current initialization, valid reopen and explicit unsupported-state refusal.
- **Release candidate**: Existing v2.0.0 identity tied to actual accepted source, later staged assets and separate public promotion.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every supported application build uses the same selected storage pair and causes zero dependency metadata changes.
- **SC-002**: Current persisted tasks and audit records survive reopen without loss; rejected state remains byte-for-byte unchanged.
- **SC-003**: Cancelled lock waits return the supplied context error before the configured busy timeout, without a partial committed write.
- **SC-004**: Repeated prepared workloads preserve results and stay within the existing 10 percent performance-regression policy or carry measured justification.
- **SC-005**: Source readiness, actual public release and attended native behavior are reported separately; public release selection remains unchanged until promotion.

## Assumptions

- The selected candidate is published SQLite v1.61.1 and its declared libc v1.79.0 pair, subject to fresh metadata/license/behavior verification. The root and desktop modules already use Go 1.26 and /v2 identities.
- #404 is the existing defect record; #412 is the actual-release coordinator. No duplicate preparation issue or release framework is needed.
- All 221 retained inputs are frozen. Existing S129 release/comparison bindings are preserved; candidate dependency notes may be updated factually without rewriting historical records.
- The owner grants automatic push/official PR and in-scope corrections, with final review/merge as the stop. No tag, draft, promotion, merge or recurring monitor is authorized in this kickoff.

## Clarifications

### Session 2026-10-31

- Q: Does this kickoff authorize actual release staging or public promotion? A: No. It authorizes source implementation/push/PR and stops for final owner review/merge; continue the existing #412 release plan after its distinct approvals.
- Q: Should this introduce migrations or revive clustered scheduling? A: No. The current-only state policy and explicit cluster deferral remain.
- Q: Should the root-only 1.61.0 proposal be merged? A: No. Select the published corrected pair and converge both graphs, superseding #400 only after integration.

All functional, data, interaction, quality, integration, edge-case, constraint, terminology and completion categories are clear. Three routine decisions were resolved from owner direction and the constitution; zero owner questions remain.
