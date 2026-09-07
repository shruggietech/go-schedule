# Quickstart: Unix Credential Bounds

## Prerequisites

- Go 1.25.0 toolchain.
- A non-Windows environment for the build-tagged focused tests.
- Repository prerequisites for the canonical verification aggregate.

## Scenario 1 - Exact Valid Boundaries

1. Run the focused Unix executor tests with synthetic UID and GID values of `0` and `4294967295`.
2. Confirm every returned credential retains the exact input value.
3. Run the existing current-account tests and confirm named-account environment behavior remains unchanged.

## Scenario 2 - Rejected UID Classes

1. Supply negative, empty, alphabetic, whitespace-padded, fractional, hexadecimal-prefixed, and `4294967296` UID text with a valid GID.
2. Confirm each request returns an invalid-UID error.
3. Confirm the command's process attributes, credential, and environment retain their original state.

## Scenario 3 - Rejected GID and Atomicity

1. Supply a valid UID with each rejected GID class.
2. Confirm each request returns an invalid-GID error.
3. Confirm no UID-only credential is assigned and any caller-supplied command state remains unchanged.

## Scenario 4 - Conversion Audit

1. Search the Unix executor for assignments to process UID and GID fields.
2. Confirm each assignment receives an already validated unsigned 32-bit value.
3. Confirm no signed or architecture-sized intermediate is cast into either field.

## Focused Verification

```sh
go test ./internal/executor
go test -race ./internal/executor
```

## Canonical Verification

```sh
sh scripts/verify.sh all
```

The aggregate must pass format, vet, lint, race, GUI, coverage, documentation, and automation gates in the foreground.
