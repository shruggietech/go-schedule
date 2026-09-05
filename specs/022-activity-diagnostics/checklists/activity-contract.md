# Activity Diagnostics Contract Checklist

**Purpose**: Validate the quality, completeness, and consistency of the Activity diagnostics contract before implementation
**Created**: 2026-08-28
**Feature**: [spec.md](../spec.md)

## Recent-View Semantics

- [x] CHK001 - Is the bounded nature of Activity stated explicitly and independently of the current record count? [Clarity, Spec §FR-001]
- [x] CHK002 - Is the daemon, rather than the GUI, established as the source of truth for the configured log path? [Consistency, Spec §FR-002–FR-003]
- [x] CHK003 - Are empty-record and pre-response states distinguished so the UI cannot imply an invented path? [Edge Case, Spec §FR-004–FR-005]
- [x] CHK004 - Is exact path preservation required for whitespace, Unicode, and platform-specific separators? [Completeness, Spec §FR-006]
- [x] CHK005 - Are file opening, probing, creation, and modification expressly excluded? [Scope, Spec §FR-007]

## Compatibility

- [x] CHK006 - Does the contract require metadata to remain available independently of recent records? [Completeness, Spec §FR-004]
- [x] CHK007 - Is additive metadata compatible with consumers that only need the record collection? [Compatibility, Spec Edge Cases]
- [x] CHK008 - Are existing CLI human-readable and JSON output shapes protected? [Compatibility, Spec §FR-011]
- [x] CHK009 - Are Activity ordering, filtering, details, refresh, live updates, Clear View, alerts, and badge behavior all protected? [Coverage, Spec §FR-010]

## Startup Event

- [x] CHK010 - Is startup wording required to describe one discrete completed event? [Clarity, Spec §FR-008]
- [x] CHK011 - Are endpoint, database path, and log path all required as structured startup fields? [Completeness, Spec §FR-009]
- [x] CHK012 - Is continuously updated uptime or lifecycle state explicitly outside the slice? [Scope, Spec §FR-015]

## Delivery

- [x] CHK013 - Are documentation and changelog obligations explicit? [Traceability, Spec §FR-012–FR-013]
- [x] CHK014 - Are the two originating issues identified for closure by the delivery pull request? [Traceability, Spec §FR-014]
- [x] CHK015 - Are measurable tests defined for exact paths, unavailable metadata, empty records, startup fields, and compatibility? [Measurability, Spec §SC-001–SC-005]

## Notes

- All contract-quality checks pass. No critical ambiguity remains after the recorded autopilot clarifications.
