# Research: Windows MSI Local-Group Recovery

## R1: Local-group routing

**Decision**: Omit `Domain` from `GoScheduleAdminGroup`.

**Rationale**: WiX 6.0.2 selects its domain path for every non-empty group
domain and its elevated local path for an empty value.

**Alternatives considered**: `Domain=""` is noisier than omission;
`[ComputerName]` reproduces the defect; a custom action duplicates WiX.

## R2: Installing-account identity

**Decision**: Preserve `[LogonUser]` with `[%USERDOMAIN]`.

**Rationale**: It locates the existing interactive account and is separate from
group-operation classification.

**Alternatives considered**: Removing it broadens scope and risks domain-backed
accounts without root-cause evidence.

## R3: Regression contract

**Decision**: Require an empty effective group domain and reject a mutation that
adds `[ComputerName]` with a routing-specific error.

**Rationale**: The current positive tests falsely encode the bug.

**Alternatives considered**: Literal substring checks are formatting-sensitive;
positive-only fixtures could silently regress.

## R4: Compiled package

**Decision**: Query and record the `Wix4Group` row.

**Rationale**: Runtime selection depends on that row, not on whether both custom
actions exist in the package.

**Alternatives considered**: Sequence/custom-action presence cannot identify
the selected path.

## R5: Lifecycle isolation

**Decision**: Use separately reset hosts for fresh and v0.9.0 upgrade evidence.

**Rationale**: Uninstall must preserve the group; v0.9.0 needs a preprovisioned
group/member to establish its upgrade baseline.

**Alternatives considered**: Deleting the group invalidates preservation;
v0.8.0 does not exercise the released package named by issue #83.

## R6: Runtime proof

**Decision**: Combine post-install state with verbose-log chronology.

**Rationale**: Final state alone cannot prove group-before-service ordering,
while static sequence alone cannot prove successful execution.

**Alternatives considered**: Either evidence class alone repeats the original
verification gap.

## R7: Evidence environment

**Decision**: Require disposable supported Windows 11 for qualifying runtime
evidence; hosted Windows Server evidence is supplementary.

**Rationale**: Platform substitution and an elevated CLI probe do not prove the
stated platform or refreshed group-token access.

**Alternatives considered**: A new hosted job adds cost without satisfying the
acceptance boundary.
