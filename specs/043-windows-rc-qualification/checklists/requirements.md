# Requirements Checklist: Windows Demo Qualification

**Purpose**: Verify that the S043 specification is complete and testable before implementation **Created**: 2026-09-03 **Feature**: [spec.md](../spec.md)

## Requirement quality

- [x] CHK001 Scope distinguishes demo qualification from formal release-candidate proof.
- [x] CHK002 Every requirement is observable and uses mandatory language.
- [x] CHK003 Artifact identity includes source, version, ProductCode, size, and SHA-256.
- [x] CHK004 Missing native prerequisites remain unavailable rather than passing.
- [x] CHK005 No push, PR, tag, release, or issue closure is implicitly authorized.
- [x] CHK006 #94's error-spam proof-before-commit ordering is preserved.

## Traceability

- [x] CHK007 #94, #98, #96, #101, #104, #105, and #106 are explicitly linked.
- [x] CHK008 #102 is explicitly excluded unless it becomes release-blocking evidence.
- [x] CHK009 Automated, attended-demo, and formal-candidate evidence are distinct.
- [x] CHK010 Success criteria cover both readiness and truthful claim boundaries.

## Testability

- [x] CHK011 Canonical and focused commands can be recorded exactly.
- [x] CHK012 The compiled MSI can be independently hashed and inspected.
- [x] CHK013 Operator observations bind to one immutable demo hash.
- [x] CHK014 Rebuild invalidation and failed-observation handling are defined.
