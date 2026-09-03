# Research: Windows Demo Qualification

## Decision 1: Separate demo qualification from formal candidate proof

**Decision**: Use a local `s043-demo` MSI for pre-PR attended testing. Preserve
the formal #94 requirement for the later byte-identical draft-release MSI.

**Rationale**: The user requires all demo testing before a pull request. The
Release workflow requires a tag on reviewed source and stages a draft GitHub
release. One artifact cannot satisfy both ordering constraints. Explicit evidence
classes preserve truthful claims.

**Alternatives considered**:

- Push an unreviewed tag to obtain an exact candidate: rejected because it
  publishes before the required human gate and bypasses the PR-first constitution.
- Call a local MSI the release candidate: rejected because its provenance cannot
  satisfy #94 or the release gate's workflow identity checks.
- Defer all native testing until after PR merge: rejected because it ignores the
  operator's request to validate the demo before opening the PR.

## Decision 2: Use v1.0.0 numeric MSI identity with an explicit demo binary version

**Decision**: Compile WiX with ProductVersion `1.0.0`; embed
`1.0.0-s043-demo.<short-commit>` in binaries; use an `s043-demo` filename.

**Rationale**: Windows Installer requires a numeric product version, while the UI
and CLI need an unmistakable non-release marker. This closely exercises intended
v1 upgrade behavior without presenting the file as an official release artifact.

## Decision 3: Reuse reviewed build and inspection commands

**Decision**: Mirror the Windows section of `.github/workflows/release.yml` with
pinned WiX 6.0.2 and existing source/compiled inspectors.

**Rationale**: A second packaging path would weaken comparability. Reusing the
release commands gives the demo the same contents and authoring semantics, with
provenance remaining the intentional difference.

## Decision 4: Keep destructive automation fail-closed

**Decision**: Do not bypass the disposable-runner/elevation checks in the hosted
installer contract. Local automation compiles and inspects; the operator drives
installation and uninstall interactions.

**Rationale**: The scripts intentionally protect developer machines from service,
group, profile, and ProgramData mutation. Removing that guard to gain a local
green check would violate the safety model and produce weaker evidence.

## Decision 5: One condensed demo checklist, formal matrix retained

**Decision**: The demo handoff focuses on high-value visible journeys and records
results against the demo hash. The complete 36-observation formal matrix remains
in `test/windows/README.md` for the later staged candidate.

**Rationale**: Exploratory pre-PR testing should efficiently detect regressions.
It must not silently redefine the release gate or imply that omitted formal
observations passed.

