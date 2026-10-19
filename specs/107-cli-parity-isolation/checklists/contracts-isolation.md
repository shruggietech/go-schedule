# Requirements Checklist: CLI contracts and isolation

**Purpose**: Review authoring, target identity, persistent state and evidence requirements before coding.

**Created**: 2026-10-19

**Feature**: [spec.md](../spec.md)

## Completeness and clarity

- [x] CHK001 Are omitted, replace and clear intents unambiguous for every task execution input? [Clarity, FR-005]
- [x] CHK002 Are disabled creation and hierarchy changes specified separately? [Completeness, FR-006]
- [x] CHK003 Are preview syntax and all timing policies required without mutation? [Completeness, FR-007]
- [x] CHK004 Are heterogeneous replacement and empty-set inheritance semantics specified? [Clarity, FR-008]
- [x] CHK005 Are every persistent desktop field and saved connection choice in scope? [Completeness, FR-011, FR-012]

## Security and ownership

- [x] CHK006 Are explicit selector failure and incompatible remote selection defined without fallback? [Coverage, FR-001]
- [x] CHK007 Are native secret defaults distinct from recording harness dependencies? [Consistency, FR-002, FR-004]
- [x] CHK008 Are identity, path containment and exact owned process cleanup requirements explicit? [Coverage, FR-003, FR-014]
- [x] CHK009 Is unavailable remote data prevented from becoming substituted local data? [Coverage, FR-009]

## Compatibility and failure behavior

- [x] CHK010 Are existing scripting formats and persistence migrations retained? [Consistency, SC-005, FR-011]
- [x] CHK011 Are multi-target deadlines, counts, cancellation and independent failures specified? [Measurability, FR-010, FR-014]
- [x] CHK012 Are output channels and actionable errors explicit? [Clarity, FR-013]
- [x] CHK013 Are evidence kinds distinct from the later complete QA report feature? [Consistency, FR-004, Assumptions]

## Review

All 13 requirement-quality items pass before planning. No missing requirement or unresolved conflict was found.
