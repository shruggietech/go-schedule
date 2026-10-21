# Feature Specification: CLI contract and diagnostic repairs

**Feature Branch**: `codex/109-cli-contract-diagnostics`

**Created**: 2026-10-21

**Status**: Implemented

**Delivery**: [PR #332](https://github.com/shruggietech/go-schedule/pull/332), open and unmerged; [delivery](delivery.md), [verification](verification.md) and [publication/review record](review.md).

**Input**: Owner-authorized S109 autopilot, bundling #311 through #322 in full under parent #277 and milestone 15. Baseline main is c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd, containing merged S108 PR #310.

## Clarifications

### Session 2026-10-21

- Q: Should partial cron input errors be successful if some lines are usable? A: No. Invalid lines produce exit 2 after the full report; task-creation or I/O failures remain operational exit 1.
- Q: Should watcher changes wait for runtime reconciliation? A: No. Responses distinguish saved configuration from last observed health and explain pending observation without introducing blocking waits.
- Q: May diagnostic improvements expose credential data? A: No. Object identities and field names are useful context; tokens, authorization, credential material and unsafe raw input are not diagnostic context.

The installed clarification workflow found all scope, entity, lifecycle, interaction, observability, dependency, edge-case, compatibility, terminology and completion categories clear after these autonomous decisions. Performance and scale remain unchanged by this bounded correction. Zero owner questions were required; all 10 specification-quality items remain passing.

## User Scenarios & Testing

### User Story 1 - Reliable scripted invocations (Priority: P1)

Operators can distinguish an invalid invocation from an operational failure, select the intended daemon independently of a webhook destination, and trust validation to refuse unusable inputs before side effects.

**Why this priority**: An accepted invalid invocation or an incorrect success status can cause an unintended action or conceal a broken automation step.

**Independent Test**: Exercise invalid flags, typed values, commands, extra operands, webhook targeting, bundle bindings and mixed cron input against disposable targets; inspect status, reports and unchanged state.

**Acceptance Scenarios**:

1. **Given** an unknown flag, invalid typed value, unknown command, malformed local document or other usage/validation failure, **When** invoked, **Then** exit 2 identifies the invalid input while operational, transport and OS failures retain exit 1 (#311).
2. **Given** a command with no or a bounded number of operands, **When** surplus operands are supplied, **Then** the invocation is refused before connecting, invoking a host adapter or changing state (#312).
3. **Given** local IPC, a remote profile or an explicit remote daemon, **When** creating or editing a webhook channel, **Then** the webhook destination and daemon selection are independently expressible and cannot shadow each other (#314).
4. **Given** supplied watcher path bindings, **When** previewing a bundle, **Then** each identity resolves to exactly one bundled watcher; unknown or duplicate bindings fail before a usable plan is produced (#315).
5. **Given** cron input with valid, unsupported and invalid lines, **When** importing or previewing, **Then** any invalid line yields an unsuccessful status and all per-line and aggregate results remain available (#316).

### User Story 2 - Useful human and machine output (Priority: P2)

Operators receive meaningful help and labeled operational data. Scripts receive one structured response on stdout and diagnostics on stderr.

**Why this priority**: Output must explain the action or result without requiring users to infer identities or parse diagnostic prose as data.

**Independent Test**: Review help, normal human output, structured output and warning channels for each affected command using existing scoped fixtures.

**Acceptance Scenarios**:

1. **Given** alerts requested as JSON, **When** the deprecated command runs, **Then** stdout contains one valid JSON response and the migration notice remains on the diagnostic channel (#313).
2. **Given** credential, pairing and service help, **When** viewing parent or child help, **Then** descriptions explain the operation, relevant identities, resulting state and applicable service configuration inputs (#317).
3. **Given** actor, audit, pairing, profile, channel or watcher lists, **When** rendered for a human, **Then** every displayed column has an explicit label and identity/time/count fields are distinguishable; JSON retains its contract (#320).
4. **Given** compatibility-mode IPC, **When** the daemon emits its warning, **Then** the message explains the effective access configuration and how to find or adjust the relevant setting, preserving access behavior and diagnostic metadata (#322).

### User Story 3 - Accurate diagnostic and state context (Priority: P2)

Operators can identify failed objects and invalid fields, understand a rotated credential's association, and distinguish saved watcher configuration from last observed runtime health.

**Why this priority**: Incomplete context makes successful changes appear broken and prevents useful corrective action after refusals.

**Independent Test**: Trigger missing-object and invalid-input cases, rotate a disposable credential and toggle an owned watcher while observing before/after state.

**Acceptance Scenarios**:

1. **Given** lookup, schedule, audit, certificate-target or document failures, **When** reported, **Then** diagnostics identify the relevant object kind/identity or field/file/input source and provide accepted values or a useful next step, preserving parser detail (#318).
2. **Given** a credential rotation, **When** reporting it, **Then** available actor/daemon context is accurate and unavailable context is omitted or explicit, without fabricated zero-value records; one-time token delivery remains intact (#319).
3. **Given** a watcher definition changed before runtime reconciliation, **When** replying immediately, **Then** saved configuration, last observation and pending reconciliation are distinguishable, with useful timing or transition guidance rather than a fabricated current state (#321).

### Edge Cases

- Root and namespace commands with unknown operands, custom operand validators, valid variadic operands, help/completion commands and parser failures before execution hooks.
- File read or output write failures versus malformed file content; missing objects versus server validation errors.
- Duplicate binding identities with identical or conflicting paths, absent watchers and malformed bindings containing path delimiters.
- Both destination and remote endpoint options supplied, explicit remote credential metadata and profile selection, and all existing valid local uses.
- Mixed cron results with failed task creation, invalid lines and unsupported lines; the complete report survives a nonzero result.
- Empty lists, redacted tokens/authorization, legacy response decoding, unavailable identities, repeated watcher enable/disable transitions and degraded observations.

## Requirements

### Functional Requirements

- **FR-001**: Usage/parser/local and server validation refusals return exit 2; genuine runtime/transport/OS failures remain exit 1 (#311).
- **FR-002**: Every command enforces its documented operand bounds before connection or side effects (#312).
- **FR-003**: Alerts diagnostics never contaminate structured stdout (#313).
- **FR-004**: Webhook destination and local/profile/explicit remote daemon targeting are independent (#314).
- **FR-005**: Unknown and duplicate watcher binding identities are refused before plan creation (#315).
- **FR-006**: Any invalid cron input line produces an unsuccessful result while retaining complete line and aggregate outcomes (#316).
- **FR-007**: Credential, pairing and service help meaningfully describes the command, identities, effects and relevant inputs (#317).
- **FR-008**: Errors name relevant objects, identities, fields or input sources and useful corrections while retaining diagnostic causes (#318).
- **FR-009**: Credential rotation preserves token delivery and reports accurate or explicitly unavailable association context (#319).
- **FR-010**: All affected human operational lists label each column while preserving machine output (#320).
- **FR-011**: Watcher responses distinguish saved state, observed health and pending reconciliation (#321).
- **FR-012**: IPC compatibility warnings explain effective configuration and a next step without changing access behavior or structured metadata (#322).
- **FR-013**: Preserve valid command behavior, credential privacy, persisted schemas and existing desktop/API consumers. Update public examples and inventory when the command contract changes.

### Key Entities

- **Invocation outcome**: success, usage/validation refusal or operational failure, with separated result and diagnostic channels.
- **Target selection**: local IPC, saved remote profile or explicit remote daemon, independent of a webhook destination.
- **State observation**: saved watcher configuration and runtime observation with sufficient reconciliation context.
- **Credential result**: rotated credential, one-time token and accurate available association context.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All 12 linked issue outcomes are delivered without partial or implicit closure.
- **SC-002**: Invalid invocations never trigger a host operation or application mutation and distinguish refusal from operational failure.
- **SC-003**: Each affected structured response is readable as one document without stripping warning prose.
- **SC-004**: Every affected help/error/list/state/warning message supplies the relevant operation or identity and enough meaning to take the intended next step.

## Assumptions

- Existing S108 owned fixtures and schema-defined reports are reused; no new QA framework, dependency or live production target is needed.
- Shared response changes remain compatible with existing consumers. Persisted data, scheduling semantics and IPC authorization are not redesigned.
- Existing authenticated GitHub access is sufficient for planning. No publication, production authentication or installed-service operation is authorized by this kickoff.
- Clarifications are resolved autonomously from the issues, constitution and demonstrated contracts; unresolved material decisions are escalated only when required.
