# UX Requirements Checklist: Desktop UX and Options

**Purpose**: Review the completeness, clarity, consistency, and measurability of S042's desktop interaction and visual requirements
**Created**: 2026-09-03
**Feature**: [spec.md](../spec.md)

**Note**: This checklist evaluates the requirements themselves for author and pull-request review. It does not substitute for implementation testing.

## Requirement Completeness

- [x] CHK001 Are all ordinary navigation destinations and their exact order specified? [Completeness, Spec §FR-001]
- [x] CHK002 Is Exit explicitly distinguished from selectable view destinations? [Completeness, Spec §FR-002]
- [x] CHK003 Are appearance modes, font choices, and their initial defaults enumerated? [Completeness, Spec §FR-005–FR-006]
- [x] CHK004 Are immediate application, restart restoration, invalid-value fallback, and reset behavior all covered? [Completeness, Spec §FR-007–FR-009]
- [x] CHK005 Are required storage categories, metadata, copy behavior, and discovery boundaries specified? [Completeness, Spec §FR-012–FR-015]
- [x] CHK006 Are single-click, double-click, toolbar, keyboard, stale-row, and duplicate-editor behaviors all addressed? [Completeness, Spec §FR-016–FR-019]
- [x] CHK007 Are the semantic text and dynamic-version requirements for the affected Info labels retained? [Completeness, Spec §FR-011]

## Requirement Clarity

- [x] CHK008 Is the meaning of Follow system distinct from forced Light and forced Dark? [Clarity, Spec clarification session and §FR-005]
- [x] CHK009 Are Brand, System, and Monospace sufficiently bounded to reject arbitrary font-file loading? [Clarity, Spec §FR-006 and Scope Boundaries]
- [x] CHK010 Is the default behavior for absent, empty, malformed, and future preference values unambiguous? [Clarity, Spec §FR-008]
- [x] CHK011 Is "bottom-right" tied to the navigation rail rather than the entire application window? [Clarity, Spec §FR-002]
- [x] CHK012 Is "breathing room" translated into content-derived minimum width, balanced spacing, and unclipped-label outcomes? [Clarity, Spec §FR-004 and §SC-004]
- [x] CHK013 Is the stable identity used for task editing distinguished from a mutable list index? [Clarity, Spec §FR-017]
- [x] CHK014 Is unavailable path presentation distinguished from a fabricated or guessed location? [Clarity, Spec clarification session and §FR-015]

## Requirement Consistency

- [x] CHK015 Do the navigation order and Exit placement requirements agree across all user stories and functional requirements? [Consistency, Spec §US3 and §FR-001–FR-003]
- [x] CHK016 Do appearance defaults preserve the existing dark-first brand behavior while still enabling system and light variants? [Consistency, Spec §US1, §FR-005–FR-010]
- [x] CHK017 Do storage ownership descriptions remain consistent with preserve-or-wipe boundaries without expanding deletion scope? [Consistency, Spec §US2, §FR-013–FR-015]
- [x] CHK018 Does the Info text requirement forbid using a font preference to conceal the underlying rendering defect? [Consistency, Spec §US1 and §FR-011]
- [x] CHK019 Do double-click requirements preserve the current detail-fetch and degraded fallback contract? [Consistency, Spec §US4 and §FR-016]

## Acceptance Criteria Quality

- [x] CHK020 Can every valid appearance combination and invalid fallback be measured deterministically? [Measurability, Spec §SC-001–SC-002]
- [x] CHK021 Can storage metadata completeness and exact copied values be objectively compared? [Measurability, Spec §SC-003]
- [x] CHK022 Are navigation order, clipping, small-window behavior, and bottom anchoring tied to explicit viewport classes? [Measurability, Spec §SC-004]
- [x] CHK023 Are editor activation counts and task identity outcomes quantified for each interaction class? [Measurability, Spec §SC-005]
- [x] CHK024 Is one-shot shutdown defined with an objective at-most-once outcome? [Measurability, Spec §SC-006]
- [x] CHK025 Does native visual acceptance name both required DPI classes and the visual comparisons to make? [Measurability, Spec §SC-008]

## Scenario and Edge-Case Coverage

- [x] CHK026 Are operating-system appearance changes while Follow system is active addressed? [Coverage, Spec Edge Cases]
- [x] CHK027 Are development, installed, absent, inaccessible, and platform-specific path states addressed without filesystem mutation? [Coverage, Spec §US2 and Edge Cases]
- [x] CHK028 Are live task removal and reordering during activation addressed? [Coverage, Spec §US4 and Edge Cases]
- [x] CHK029 Are rapid overlapping Exit/title-bar close requests addressed? [Coverage, Spec §US3 and Edge Cases]
- [x] CHK030 Are long or badge-bearing navigation labels accounted for? [Coverage, Spec Edge Cases]

## Accessibility and Evidence

- [x] CHK031 Are focus visibility, keyboard reachability, and predictable navigation requirements defined for every new control? [Coverage, Spec §FR-002 and §FR-010]
- [x] CHK032 Are contrast and unclipped-control requirements applied to every supported appearance/font combination? [Coverage, Spec §FR-010]
- [x] CHK033 Is headless evidence separated from exact-candidate native Windows visual evidence? [Clarity, Spec §FR-020–FR-021]
- [x] CHK034 Is issue closure withheld until native criteria are met rather than inferred from automated tests? [Dependency, Spec Delivery and clarification session]

## Scope and Dependencies

- [x] CHK035 Are #101, #103, #104, #105, and #106 explicitly included and individually traceable? [Traceability, Spec Scope Boundaries]
- [x] CHK036 Are #94, #98, #96, task diagnostics, arbitrary fonts, configuration editing, and path mutation explicitly excluded? [Scope, Spec Scope Boundaries]
- [x] CHK037 Is the dependency on the existing exact-candidate attended gate explicit? [Dependency, Spec §FR-021]
- [x] CHK038 Is the absence of scheduler, daemon protocol, installer deletion, and promotion changes stated? [Assumption, Spec Assumptions]

## Notes

- Standard-depth reviewer checklist covering interaction, appearance, accessibility, lifecycle safety, and evidence boundaries.
- Validation passed 38 of 38 requirement-quality checks.
