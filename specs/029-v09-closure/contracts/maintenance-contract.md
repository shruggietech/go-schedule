# Maintenance Contract: S029

## Lifecycle checker

The checker accepts an optional repository root and exits 0 only when:

- every `specs/[0-9][0-9][0-9]-*/spec.md` has one allowed `**Status**` value;
- every implemented spec has a non-empty `**Delivery**` value;
- a Draft or Ready spec has at least one actionable unchecked task;
- an Implemented spec has no actionable unchecked task; and
- the inventory covers exactly the discovered feature directories.

Failures go to standard error, name the offending file, and exit nonzero.

## Dependabot policy

- Go modules: `/`, weekly, grouped minor/patch routine updates, five-proposal cap.
- GitHub Actions: `/`, monthly, grouped minor/patch routine updates, five-proposal cap.
- Both: `dependencies` label, target `main` implicitly through the default branch.
- Major and security updates are not placed in routine groups.

## Hosted security policy

S029 requests `enabled` for:

- `secret_scanning_push_protection`
- `secret_scanning_non_provider_patterns`
- `secret_scanning_validity_checks`

The operation is successful only after readback confirms each control. No other security-and-analysis field is mutated.

## Issue #33 policy

The issue stays open, keeps `needs: verification`, moves to `Post-v1`, and changes from P1 to P3. The triage comment must state that no current reproduction was performed.
