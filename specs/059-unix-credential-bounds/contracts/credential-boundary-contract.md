# Contract: Unix Process Credential Boundary

## Input

The existing `run_as` flow supplies:

- the requested account label;
- a resolved account record containing UID, GID, username, and home directory;
- the command being prepared;
- whether the caller supplied an explicit home directory.

No public configuration, CLI, or API shape changes.

## Validation

- UID and GID are parsed independently as base-10 unsigned 32-bit values.
- Accepted numeric range is 0 through 4,294,967,295 inclusive.
- Empty, negative, whitespace-padded, fractional, non-decimal, and overflowing values fail.
- A UID failure returns context containing `invalid uid` and the requested account.
- A GID failure returns context containing `invalid gid` and the requested account.

## Mutation

- Empty `run_as`, lookup failure, UID failure, or GID failure does not mutate the command.
- After both values validate, one credential pair is assigned together.
- Successful environment behavior remains unchanged: `LOGNAME` and `USER` are replaced, and `HOME` is replaced only when the task did not supply one explicitly.

## Compatibility

- Named-account lookup remains first.
- Numeric-account lookup remains the fallback.
- Existing valid accounts preserve exact credential and environment behavior.
- Windows behavior remains governed by the separate Windows implementation.

## Evidence

- Synthetic unit cases cover the complete boundary without host account provisioning.
- Existing current-account tests cover production integration.
- Repository CI and hosted CodeQL must pass without suppression.
