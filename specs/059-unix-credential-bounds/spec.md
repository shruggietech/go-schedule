# Feature Specification: Unix Credential Bounds

**Feature Branch**: `codex/059-unix-credential-bounds`

**Created**: 2026-09-07

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Pull request [#191](https://github.com/shruggietech/go-schedule/pull/191); local canonical verification, hosted Linux/macOS/Windows race tests, CodeQL, packaging, and two Codex review rounds completed 2026-09-07

**Input**: Resolve GitHub issue [#145](https://github.com/shruggietech/go-schedule/issues/145) by rejecting invalid or out-of-range Unix UID and GID values before task process credentials are assigned, then complete the slice under autopilot and publish it for CI and third-party review.

## Problem Statement

Unix task execution resolves an account and converts its UID and GID into the fixed-width values used to start the child process. The current conversion accepts signed machine-sized integers and narrows them afterward. Negative or oversized values can therefore wrap or truncate into a different identity. A configured `run_as` task could execute under credentials other than the resolved account.

S059 must make the identity boundary explicit, reject every value that cannot be represented exactly, and leave the command unchanged when either identifier is invalid. The fix must preserve ordinary named-user and numeric-user behavior and close the four CodeQL findings attached to issue #145 through code correction rather than dismissal.

## Scope

### In scope

- Validation of resolved Unix UID and GID text as decimal values representable from zero through 4,294,967,295 inclusive.
- Construction and assignment of process credentials only after both identifiers have passed validation.
- Clear contextual errors that identify whether UID or GID validation failed and which requested account was affected.
- Deterministic tests using synthetic resolved account data, independent of the host account database.
- Preservation of existing environment setup and valid named-user and numeric-user `run_as` behavior.
- An audit of the surrounding Unix task process-credential path for equivalent narrowing conversions governed by the same boundary.
- Focused, race, full-suite, CodeQL, and canonical repository verification evidence.

### Out of scope

- Changing account lookup order or the public `run_as` configuration contract.
- Adding supplementary groups, privilege transitions, user namespaces, capabilities, or Windows impersonation.
- Changing Unix IPC administrator-group ownership logic unless the audit proves that it performs the same unsafe fixed-width process-credential narrowing.
- Dismissing or suppressing the CodeQL alerts without correcting their shared cause.
- Any Wails control-center, notification, MCP, remote-access, or distributed-scheduling roadmap work.

## Clarifications

### Session 2026-09-07

- Q: Which values are valid at the Unix process-credential boundary? -> A: Decimal UID and GID text representing the complete unsigned 32-bit range, with malformed, signed-negative, empty, and overflowing values rejected.
- Q: What may change when either identifier is invalid? -> A: Nothing on the command; both identifiers are validated before process attributes, credentials, or environment values are mutated.
- Q: How far does the required conversion audit extend? -> A: The Unix task process-credential path governed by issue #145; unrelated integer parsing remains outside S059 unless it narrows into the same credential fields.
- Q: Does the operator's up-front publication direction satisfy the usual autopilot pre-push authorization? -> A: Yes for this S059 review branch, pull request, and verified in-scope review fixes only; merge, tag, and release remain unauthorized.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preserve the Resolved Unix Identity (Priority: P1)

An operator who configures a Unix task to run as another account can trust that the child process will receive exactly that account's resolved UID and GID or will not be configured to start.

**Why this priority**: Executing under an unintended numeric identity is a security boundary failure and the repository's only current P0 issue.

**Independent Test**: Supply synthetic accounts whose UID and GID cover valid boundaries and every rejected class, then inspect the resulting command without consulting the host account database.

**Acceptance Scenarios**:

1. **Given** resolved UID and GID text at zero or 4,294,967,295, **when** process credentials are prepared, **then** both exact numeric values are retained without wrapping or truncation.
2. **Given** either resolved identifier is negative, empty, malformed, non-decimal, or greater than 4,294,967,295, **when** process credentials are prepared, **then** the request fails before any command state is changed.
3. **Given** a valid UID and an invalid GID, **when** preparation fails on the GID, **then** no UID-only or replacement credential is left on the command.

---

### User Story 2 - Keep Valid Run-As Behavior Compatible (Priority: P2)

An operator using an ordinary named or numeric Unix account continues to receive the existing account lookup, environment, and execution behavior after the security correction.

**Why this priority**: The correction must close the unsafe boundary without breaking supported task execution.

**Independent Test**: Exercise the existing current-account paths with explicit and inherited home-directory behavior alongside the new boundary tests.

**Acceptance Scenarios**:

1. **Given** an existing valid named account, **when** `run_as` is applied, **then** the resolved UID, GID, username, and home-directory behavior remain correct.
2. **Given** an existing valid numeric account identifier, **when** account lookup falls back to numeric lookup, **then** credential preparation remains supported.
3. **Given** no `run_as` value, **when** task execution is prepared, **then** it remains a no-op for credentials and environment.

---

### User Story 3 - Produce Reviewable Security Evidence (Priority: P3)

A maintainer can verify that the shared cause of all four CodeQL alerts is fixed and that no equivalent narrowing remains in the affected execution path.

**Why this priority**: The issue is not complete until the correction is reviewable, mechanically checked, and suitable for hosted security reanalysis.

**Independent Test**: Inspect the affected conversion path, run focused and canonical checks, and confirm hosted CodeQL reports no replacement finding after publication.

**Acceptance Scenarios**:

1. **Given** the final Unix execution path, **when** it is audited, **then** no signed or architecture-sized intermediate is narrowed into a process UID or GID.
2. **Given** the published pull request, **when** hosted CodeQL completes, **then** alerts #1 through #4 close through the code change with no equivalent replacement alert.
3. **Given** the complete repository, **when** canonical verification runs, **then** all required gates pass without weakening an existing security or test gate.

### Edge Cases

- Zero is a valid Unix identifier and must not be confused with a missing value.
- The maximum unsigned 32-bit value is valid even on a host whose machine-sized signed integer is narrower.
- Leading and trailing whitespace, signs that do not represent a permitted nonnegative identifier, fractional values, hexadecimal notation, and other non-decimal text are invalid.
- UID validation can succeed before GID validation fails; the command must remain byte-for-byte equivalent in all fields governed by this operation.
- A command may already carry process attributes or credentials supplied by its caller; a validation failure must preserve that prior state.
- Account lookup failure remains distinct from resolved-account identifier validation failure.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Resolved Unix UID and GID values MUST be validated as decimal values within the inclusive range 0 through 4,294,967,295 before conversion.
- **FR-002**: Negative, empty, malformed, non-decimal, and out-of-range UID values MUST produce a contextual UID error.
- **FR-003**: Negative, empty, malformed, non-decimal, and out-of-range GID values MUST produce a contextual GID error.
- **FR-004**: The validated numeric values MUST be preserved exactly, including both inclusive boundaries.
- **FR-005**: Both identifiers MUST be validated before the command's process attributes, credential pointer, or environment is changed.
- **FR-006**: Failure after one identifier validates MUST leave all command state governed by `run_as` unchanged, including any pre-existing process attributes or credentials.
- **FR-007**: A successful request MUST assign the UID and GID together as one complete credential pair.
- **FR-008**: Existing valid named-user lookup, numeric-user fallback, empty-value behavior, and environment handling MUST remain compatible.
- **FR-009**: Tests MUST exercise valid boundaries and every rejected class independently for UID and GID using synthetic account records rather than host account configuration.
- **FR-010**: Tests MUST prove that a valid UID paired with an invalid GID cannot partially configure or replace process credentials.
- **FR-011**: The affected Unix task process-credential path MUST contain no equivalent narrowing conversion after implementation.
- **FR-012**: All four linked CodeQL alerts MUST be resolved by the merged correction without dismissal or an equivalent replacement alert.
- **FR-013**: The focused Unix tests and the repository's format, vet, lint, race, GUI, coverage, documentation, and automation gates MUST pass.
- **FR-014**: Files created or changed by S059 MUST remain UTF-8 without BOM and contain no mojibake.

### Key Entities

- **Resolved Unix account**: Account data returned by the existing lookup flow, including textual UID, textual GID, username, and home directory.
- **Validated credential pair**: Exact nonnegative UID and GID values that both fit the process credential representation.
- **Command preparation state**: Process attributes, credentials, and environment that must be updated only after validation succeeds.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One hundred percent of valid boundary cases, zero and 4,294,967,295 for both identifiers, retain their exact values.
- **SC-002**: One hundred percent of negative, empty, malformed, non-decimal, and overflowing test cases are rejected before command mutation.
- **SC-003**: Every mixed-validity case leaves pre-existing command credentials and environment unchanged.
- **SC-004**: Existing valid named-account, numeric-account, empty-value, and home-directory tests continue to pass.
- **SC-005**: The affected execution path contains zero narrowing conversions from signed or architecture-sized integers into process UID or GID fields.
- **SC-006**: Hosted security analysis reports zero open instances corresponding to CodeQL alerts #1 through #4 and no equivalent replacement finding.
- **SC-007**: All eight canonical repository gates pass, core package coverage remains at or above 80 percent, and the changed files pass UTF-8, BOM, mojibake, and whitespace checks.

## Assumptions

- The operating-system account lookup contract supplies UID and GID as base-10 text for valid accounts.
- The process credential fields remain unsigned 32-bit values on supported Unix targets.
- The standard runtime conversion library is sufficient; no dependency is needed.
- Hosted CodeQL closure occurs after publication and merge processing, so the pull request records the corresponding clean analysis as release evidence when available.
- Issue #145 is the only GitHub issue closed by S059; later roadmap issues remain independent.
