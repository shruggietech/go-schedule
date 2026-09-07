# Data Model: Unix Credential Bounds

S059 changes no persisted data. The model describes the transient security boundary used while preparing a Unix child process.

## Resolved Unix Account

Existing account lookup output:

- `requested`: the original nonempty `run_as` value used for error context.
- `username`: canonical account name used for `LOGNAME` and `USER`.
- `home_directory`: resolved home used when no explicit task home overrides it.
- `uid_text`: base-10 account identifier text.
- `gid_text`: base-10 primary-group identifier text.

## Validated Credential Pair

Transient all-or-nothing result:

- `uid`: exact unsigned 32-bit value parsed from `uid_text`.
- `gid`: exact unsigned 32-bit value parsed from `gid_text`.

The pair does not exist when either parse fails. No partially valid form is returned or assigned.

## Command Preparation State

Existing mutable command fields governed by `applyRunAs`:

- `SysProcAttr`: existing process attributes, possibly absent or caller-supplied.
- `SysProcAttr.Credential`: existing credential pointer, possibly absent or caller-supplied.
- `Env`: existing environment, possibly absent.

## State Transitions

1. Empty `run_as`: return success with no state transition.
2. Lookup failure: return an unknown-user error with no state transition.
3. UID validation failure: return an invalid-UID error with no state transition.
4. GID validation failure: return an invalid-GID error with no state transition.
5. Complete validation success: ensure process attributes exist, assign the complete credential pair, then apply the existing environment rules.

## Invariants

- UID and GID are each within 0 through 4,294,967,295 inclusive.
- No signed or architecture-sized intermediate is narrowed into the credential pair.
- Failure preserves the complete pre-call command state.
- Success never assigns only one member of the pair.
