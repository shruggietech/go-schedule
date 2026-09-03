# Release Evidence Requirements Checklist: Windows Release Candidate Gate

**Purpose**: Review the completeness, consistency, security, and measurability of the exact-candidate evidence and fail-closed publication requirements
**Created**: 2026-09-02
**Feature**: [spec.md](../spec.md)

**Audience and depth**: Formal pull-request review checklist for release maintainers and Windows evidence operators.

## Candidate and Promotion Integrity

- [x] CHK001 Does identity bind repository, commit, version, ProductCode, workflow run and attempt, filename, size, and SHA-256? [Completeness, Spec §FR-001–FR-003]
- [x] CHK002 Is byte-identical promotion required and rebuild-after-test prohibited? [Consistency, Spec §Clarifications, Spec §FR-026]
- [x] CHK003 Do missing candidate, missing evidence, mismatches, and non-pass results leave the release unpublished? [Fail Closed, Spec §FR-024–FR-027]
- [x] CHK004 Are final checksums defined only after the evidence asset joins the final asset set? [Completeness, Spec §SC-008]

## Native Windows Coverage

- [x] CHK005 Are intended-user access, LocalSystem service identity, and unrelated-user denial all required? [Coverage, Spec §FR-008]
- [x] CHK006 Are native rectangles, logical content and work area, DPI, monitor, and state flags distinguished? [Clarity, Spec §FR-009]
- [x] CHK007 Are first launch, retained profile, both DPI classes, window controls, and subsequent launch covered? [Coverage, Spec §FR-010–FR-012]
- [x] CHK008 Is conservative sizing expressed as 1280 by 800 on sufficiently large work areas and independent 90 percent caps on smaller work areas? [Measurability, Spec §FR-011]

## Connection and Execution Coverage

- [x] CHK009 Are unavailable, access denied, timeout, stream disconnect, repeated refresh or reconnect, retry, and recovery enumerated? [Coverage, Spec §FR-013–FR-015]
- [x] CHK010 Is the two-minute interval and visible in-frame, modal, and top-level surface count measurable? [Measurability, Spec §FR-014, Spec §SC-004]
- [x] CHK011 Is the proof-before-commit order preserved for any new recurring-error correction? [Governance, Spec §FR-016]
- [x] CHK012 Are manual success, scheduled success, nonzero exit, and process-start failure proven through production execution? [Coverage, Spec §FR-017–FR-018]

## Installer Lifecycle Coverage

- [x] CHK013 Are shortcut defaults/selections, completion combinations, normal-token finish launch, cancel, repair, and upgrade represented? [Coverage, Spec §FR-019]
- [x] CHK014 Are preserve, wipe, locked cleanup, genuine multiple profiles, reinstall, controls, and security-state treatment represented? [Coverage, Spec §FR-019–FR-020]
- [x] CHK015 Does the specification prohibit silently converting partial cleanup or unavailable observation into success? [Consistency, Spec §FR-006–FR-007]

## Evidence Safety and Auditability

- [x] CHK016 Are bundle-relative attachment paths, traversal refusal, size, and digest checks required? [Security, Spec §FR-005]
- [x] CHK017 Are secrets and unnecessary personal data excluded while audit context is retained? [Privacy, Spec §FR-029]
- [x] CHK018 Are fixture evidence and hosted CI explicitly prevented from masquerading as attended Windows 11 proof? [Integrity, Spec §FR-023, Spec §Out of Scope]
- [x] CHK019 Is safe resume required without overwriting earlier evidence? [Reliability, Spec §FR-022]
- [x] CHK020 Are issue closure and real-candidate acceptance deferred until evidence actually passes? [Traceability, Spec §Clarifications, Spec §Dependencies and Traceability]

## Notes

- All 20 requirement-quality questions passed.
- The checklist deliberately treats native attended observation and automated validator correctness as complementary evidence, not substitutes.
