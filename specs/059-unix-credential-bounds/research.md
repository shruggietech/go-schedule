# Research: Unix Credential Bounds

## R1 - Numeric Conversion Boundary

**Decision**: Parse resolved UID and GID text directly with base 10 and a 32-bit unsigned limit.

**Rationale**: The destination contract is an unsigned 32-bit value. Validating that exact representation prevents negative wrapping, architecture-dependent acceptance, and truncation.

**Alternatives considered**: `strconv.Atoi` plus range checks was rejected because it introduces an unnecessary architecture-sized intermediate and is harder for CodeQL to prove safe. Parsing at machine width and checking afterward was rejected for the same reason. A regular expression plus conversion was rejected as duplicate validation.

## R2 - Failure Atomicity

**Decision**: Construct the full credential pair before allocating or mutating command process attributes and before writing environment variables.

**Rationale**: A valid UID followed by an invalid GID must not partially configure the command. A pure construction step makes this invariant direct and preserves any state the caller supplied.

**Alternatives considered**: Assigning UID then parsing GID was rejected because rollback is easy to miss and can overwrite caller state. Clearing the credential on failure was rejected because that is also a mutation and would destroy a pre-existing value.

## R3 - Test Seam

**Decision**: Extract credential construction from a supplied resolved account record and keep production lookup unchanged.

**Rationale**: Boundary values and malformed account data cannot be provisioned portably in the host account database. Passing a synthetic record to a pure helper gives deterministic coverage without mutable package globals.

**Alternatives considered**: Replacing `lookupUser` with a mutable function variable was rejected because it creates shared state and race risk. Depending only on the current account was rejected because it cannot exercise invalid or maximum identifiers.

## R4 - Accepted Text Form

**Decision**: Accept the decimal syntax supported by the standard unsigned parser at explicit base 10, including leading zeroes, while rejecting whitespace, negative values, fractional text, and alternate-base prefixes.

**Rationale**: Operating-system account records are expected to provide decimal text. Explicit base 10 avoids accidental octal or hexadecimal interpretation and exact numeric preservation is the security property.

**Alternatives considered**: Requiring an additional digits-only regular expression was rejected because it adds a second parser without protecting the numeric boundary. Base-zero parsing was rejected because it accepts notation outside the account-record contract.

## R5 - Audit Boundary

**Decision**: Audit `internal/executor/runas_unix.go` and other assignments to `syscall.Credential.Uid` or `.Gid`; do not alter the Unix IPC administrator-group conversion without evidence of the same destination narrowing.

**Rationale**: Issue #145 concerns fixed-width child-process credentials. The IPC path produces the platform `int` required by `os.Chown`, has separate access-control tests, and is not an assignment into the affected credential fields.

**Alternatives considered**: Refactoring every `strconv.Atoi` call was rejected as unrelated scope with different valid ranges and destination contracts.

## R6 - Checklist Ordering Workaround

**Decision**: Generate the domain checklist directly from the resolved feature path before planning.

**Rationale**: The installed prerequisite rejects a missing plan even though both Spec-Kit command order and project autopilot require checklist first. Direct feature-path generation preserves the governing semantic order.

**Alternatives considered**: Planning first and skipping the checklist were both rejected because each violates the explicit autopilot procedure. Patching Spec-Kit tooling inside this security slice was rejected as unrelated process-tool scope.
