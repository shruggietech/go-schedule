# Data Model: Maintainer Automation Baseline

This feature persists no runtime or user data. The model below describes the repository automation records whose consistency is enforced in versioned files.

## ActionReference

An external action invocation in a workflow.

| Field | Meaning | Validation |
| --- | --- | --- |
| workflow | Repository-relative workflow path | Must be under `.github/workflows/` |
| location | Step position or line | Must identify one `uses:` entry |
| owner_repository | Action owner and repository | Must be in the audited allowlist |
| selected_major | Floating major selector | Must equal the approved Node 24 major |
| declared_runtime | Runtime from the selected action metadata | Must not be Node 20 |
| used_contract | Inputs, outputs, permissions, artifacts | Must remain compatible after upgrade |

**Uniqueness**: Each workflow location contains one action reference. The same approved action may appear at multiple locations.

## VerificationGate

A non-mutating check in the local/CI verification contract.

| Field | Meaning | Validation |
| --- | --- | --- |
| name | Stable command mode | One of the eight required names |
| command | Canonical executable behavior | Defined once in `scripts/verify.sh` |
| prerequisites | Required host tools | Missing prerequisites fail closed |
| ci_consumers | CI jobs that invoke the mode | At least one where the gate runs in CI |
| local_consumer | Aggregate verification | Every required gate appears exactly once |
| success | Exit behavior | Zero only when the child gate passes |
| failure | Diagnostic behavior | Non-zero and identifies the gate |
| mutating | Whether repository files may change | False for every verification gate |

The required gate identities are `format`, `vet`, `lint`, `race`, `gui`, `coverage`, `docs`, and `automation`.

## VerificationContract

The ordered set of VerificationGate records that defines a local green result.

- **Manifest**: All eight required names, with no duplicates or omissions.
- **Aggregate transition**: `not_started -> running(gate) -> passed` or `not_started -> running(gate) -> failed(gate)`.
- **Failure rule**: The first failing or unavailable gate terminates aggregate execution with a non-zero result.
- **Cleanliness rule**: A passing run over a clean checkout leaves it clean.

## AutomationPolicy

The independent offline policy used to catch drift.

| Field | Meaning | Validation |
| --- | --- | --- |
| approved_actions | Audited action-major references | Exactly the four decisions in research D1 |
| required_gates | Verification manifest | Exactly the eight gates above |
| repository_root | Root to inspect | Defaults to current repository; overrideable for fixtures |
| network_required | Whether validation fetches metadata | Always false |

**Relationships**:

- AutomationPolicy admits many ActionReference records.
- VerificationContract contains eight VerificationGate records.
- AutomationPolicy compares its independent `required_gates` with the VerificationContract manifest.

## PinnedArtifactDecision

A dated changelog record for a governed repository file.

| Field | Meaning | Validation |
| --- | --- | --- |
| date | Decision date | ISO calendar date |
| path | Pinned file changed | Workflow or Makefile path |
| decision | What changed | Names runtime or verification contract |
| rationale | Why it changed | Connects to #21 or #41 and preserved constraints |
