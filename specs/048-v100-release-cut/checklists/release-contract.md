# Release Contract Checklist: v1.0.0

**Purpose**: Validate preparation, staging, qualification, and promotion coverage

**Created**: 2026-09-03

**Feature**: [spec.md](../spec.md)

## Preparation

- [x] CHK001 Does the contract preserve all current Unreleased entries beneath v1.0.0? [Completeness]
- [x] CHK002 Does it require an empty future Unreleased section and correct comparison links? [Consistency]
- [x] CHK003 Does it validate the tag-specific health example and publication-aware badge before tagging? [Release preflight]
- [x] CHK004 Does it constrain public notes to the established tag-specific highlights policy? [UX]
- [x] CHK005 Does the preparation PR leave #122 and evidence-dependent issues open? [Lifecycle]

## Tag and staging

- [x] CHK006 Are clean/synchronized reviewed-main checks required immediately before tag creation? [Integrity]
- [x] CHK007 Are local tag, remote tag, and GitHub release absence checked? [Idempotence]
- [x] CHK008 Is explicit post-merge tag authority required? [Governance]
- [x] CHK009 Is the tag immutable and bound to the reviewed merge commit? [Integrity]
- [x] CHK010 Must the Release run be successful, tag-triggered, and commit-identical? [Provenance]
- [x] CHK011 Are all seven build packages plus the candidate manifest inventoried? [Completeness]

## Formal evidence

- [x] CHK012 Is only the exact staged MSI eligible for qualification? [Provenance]
- [x] CHK013 Are all 47 S047 scenarios mandatory and passing? [Completeness]
- [x] CHK014 Are local-demo results explicitly non-transferable? [Evidence boundary]
- [x] CHK015 Are raster bytes, attachment hashes, environment identity, operator attestation, and run identity validated? [Integrity]
- [x] CHK016 Is archive upload named and protected against unexplained overwrite? [Safety]

## Reconciliation and promotion

- [x] CHK017 Are all ten remaining v1 issues mapped and reconciled individually? [Traceability]
- [x] CHK018 Does #96 close only after child and coordinator criteria pass? [Hierarchy]
- [x] CHK019 Does promotion require the exact nine pre-checksum payloads and add one checksum per payload? [Completeness]
- [x] CHK020 Does promotion publish the existing draft without rebuilding? [Provenance]
- [x] CHK021 Does final audit cover release/latest identity, assets, notes, README, issues, milestone, and clean main? [Closure]
- [x] CHK022 Do all failure paths keep the draft and affected issues open? [Fail closed]

## Notes

- Complete. The contract preserves every S040/S047 release gate and the constitution's separate tag/release authorization boundary.
