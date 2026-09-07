# Implementation Plan: Unix Credential Bounds

**Branch**: `codex/059-unix-credential-bounds` | **Date**: 2026-09-07 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/059-unix-credential-bounds/spec.md`

## Summary

Replace signed machine-sized parsing followed by narrowing with one small unsigned 32-bit parser, derive a complete credential pair before mutating the command, and cover the boundary through synthetic resolved-account tests. Preserve the existing lookup and environment contracts, add no dependency or persisted state, and provide hosted CodeQL evidence for issue #145.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library (`os/user`, `strconv`, `syscall`); no new dependency

**Storage**: N/A, no schema or configuration change

**Testing**: Go unit tests on non-Windows targets, existing executor integration tests, race detector, CodeQL, and canonical repository verification

**Target Platform**: Supported Unix daemon targets, with repository-wide Windows build and test compatibility preserved by build tags

**Project Type**: Single Go module with daemon, CLI, local API, and desktop GUI

**Performance Goals**: Constant-time validation during process setup with no scheduler hot-path or steady-state resource change

**Constraints**: Validate the inclusive unsigned 32-bit range, mutate no command state on failure, keep error context, preserve valid `run_as` behavior, and add no dependency

**Scale/Scope**: One Unix executor implementation file, one focused test file, Spec-Kit artifacts, lifecycle inventory, changelog, and verification evidence

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Code Quality**: PASS. One focused helper owns the fixed-width boundary, errors retain field and account context, and the successful path remains linear.
- **II. Testing Standards**: PASS. Regression tests precede implementation and cover both fields, inclusive boundaries, rejected classes, failure atomicity, and compatibility. Race and coverage gates remain mandatory.
- **III. User Experience Consistency**: PASS. Existing `run_as` input and lookup behavior remain stable; new failures clearly distinguish invalid UID from invalid GID.
- **IV. Performance Requirements**: PASS. Two bounded numeric parses occur only during child-process preparation. No scheduling algorithm, concurrency, or performance-sensitive loop changes, so no benchmark is warranted.
- **V. Autonomous Build-Phase Execution**: PASS. S059 traces to open issue #145, uses a review branch, runs the full Spec-Kit sequence, and has explicit authorization to publish the PR and verified review fixes.
- **Engineering constraints**: PASS. The standard library is sufficient, no persistent format changes, inputs are validated at the boundary, and supported-platform build tags remain intact.

## Phase 0 Research Decisions

1. Parse UID and GID directly as base-10 unsigned 32-bit values rather than parse to `int` and cast.
2. Separate pure credential-pair construction from account lookup so synthetic user records cover impossible-to-provision boundaries without global test hooks.
3. Return the complete credential pair before touching `SysProcAttr`, `Credential`, or `Env` so every failure preserves caller state.
4. Preserve the existing `lookupUser` function and environment-writing order on success.
5. Limit the conversion audit to the Unix executor process-credential path. The IPC administrator-group path converts to the `int` contract required by `os.Chown` and does not narrow into `syscall.Credential`.
6. Use the established direct-feature-path workaround for the checklist phase because the installed prerequisite script contradicts the mandated checklist-before-plan order.

See [research.md](research.md) for alternatives and rationale.

## Phase 1 Design

- Add a pure helper that accepts the requested account label and resolved account data, validates UID and GID with explicit bit size 32, and returns one `syscall.Credential` value or a contextual error.
- Keep `applyRunAs` responsible for the empty-value no-op, account lookup, successful command mutation, and environment setup.
- Call the pure helper immediately after lookup. Assign process attributes and environment only after it returns successfully.
- Test the helper with synthetic user records for zero, maximum, leading zero, negative, empty, malformed, whitespace, alternate-base notation, and overflow cases for both fields.
- Test `applyRunAs` compatibility through the existing current-account scenarios and assert that empty `run_as` remains a no-op.
- Record focused and canonical evidence in `verification.md`; use hosted CodeQL and PR review results as publication evidence.

### Post-design constitution re-check

All gates remain PASS. The pure helper is the smallest testable design that avoids host-account dependencies. It does not introduce a new abstraction layer or dependency, and validating before mutation directly satisfies the security and test requirements.

## Project Structure

### Documentation (this feature)

```text
specs/059-unix-credential-bounds/
├── spec.md / plan.md / research.md / data-model.md / quickstart.md
├── contracts/credential-boundary-contract.md
├── checklists/requirements.md / security.md
└── tasks.md / verification.md
```

### Source Code (repository root)

```text
internal/executor/
├── runas_unix.go
└── runas_unix_test.go
```

**Structure Decision**: Keep the fix inside the existing build-tagged executor module. The parsing rule has no independent domain identity outside Unix child-process preparation, so a new package would add indirection without reuse.

## Complexity Tracking

No constitution violation or exceptional complexity is required.

## Tooling Deviation

The installed checklist prerequisite requires `plan.md`, while the governing autopilot order requires checklist before plan. S059 generated and validated `checklists/security.md` directly from the feature path before planning. This follows the repository's established workaround, preserves the required order, and does not weaken the checklist.
