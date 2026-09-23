# Desktop notification requirements checklist

**Purpose**: Test requirements quality for desktop delivery, muting, attribution, activation, and release scope.
**Created**: 2026-09-23
**Audience**: S103 author and PR reviewers.

## Completeness and clarity

- [x] CHK001 Are opt-in, mute, and persistence requirements explicit and mutually consistent? [Spec §FR-001, §FR-002]
- [x] CHK002 Are daemon, severity, and condition filter dimensions named? [Spec §FR-002]
- [x] CHK003 Is the source identity required independently of human-readable task names? [Spec §FR-003]
- [x] CHK004 Are startup and reconnection backlog expectations bounded? [Spec §FR-004]
- [x] CHK005 Is popup activation defined for available and unavailable source records? [Spec §FR-005]

## Platform and failure coverage

- [x] CHK006 Are supported desktop platforms and headless exclusions documented? [Spec §FR-006, Assumptions]
- [x] CHK007 Are permission denial and unsupported native facilities included without scheduler impact? [Spec §FR-006]
- [x] CHK008 Is closed-application delivery explicitly excluded? [Spec §FR-007]
- [x] CHK009 Is duplicate suppression defined across multiple daemon identities? [Spec §FR-003, §SC-001]
- [x] CHK010 Are concurrent preference changes and identity changes included as edge cases? [Spec §Edge Cases]

## Release and issue integrity

- [x] CHK011 Is public publication distinguished from release preparation? [Spec §FR-008]
- [x] CHK012 Are SMTP and clustered execution excluded from shipped claims? [Spec §FR-008, §FR-009]
- [x] CHK013 Are #176, #19, and #146 kept open when their outcomes remain incomplete? [Spec §FR-009]
