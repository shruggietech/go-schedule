# Release Evidence Requirements Checklist: Windows Release Qualification

**Purpose**: Test whether S047's release-evidence requirements are complete, unambiguous, measurable, safe, and issue-traceable before planning
**Created**: 2026-09-03
**Feature**: [spec.md](../spec.md)

## Evidence Identity and Boundaries

- [x] CHK001 Is the distinction between a local demo and the formal workflow-staged candidate explicit and irreversible? [Clarity, Spec §Scope]
- [x] CHK002 Are the actions that require later publication or release authority explicitly excluded? [Completeness, Spec §Scope out]
- [x] CHK003 Is every native claim bound to an identified environment, candidate, observation, and attachment? [Traceability, Spec §FR-003]
- [x] CHK004 Does the specification prevent synthetic, headless, or local-demo evidence from satisfying the production class? [Consistency, Spec §FR-014, §FR-019]
- [x] CHK005 Are rebuild and byte-identity edge cases addressed without permitting evidence transfer? [Coverage, Spec §Edge Cases]

## Desktop Appearance and Interaction

- [x] CHK006 Are palette, DPI, System-font, lifecycle, and sharpness requirements complete for the affected text? [Completeness, Spec §FR-004]
- [x] CHK007 Are the normal-text and essential non-text contrast floors quantified? [Measurability, Spec §FR-005]
- [x] CHK008 Are all shared control families and interactive states named rather than represented by a sample control only? [Coverage, Spec §FR-005]
- [x] CHK009 Is non-color state communication required consistently across interaction and semantic-row evidence? [Accessibility, Spec §FR-005, §FR-010]
- [x] CHK010 Are native screenshots required in addition to metric assertions for visually judged outcomes? [Completeness, Spec §FR-003]

## Navigation, Options, and Input

- [x] CHK011 Are rail spacing, boundary, destination order, Options placement, and Exit placement/state specified together? [Consistency, Spec §FR-006]
- [x] CHK012 Are compact storage rows, unavailable rows, selector alternatives, Copy fidelity, and overflow behavior all covered? [Completeness, Spec §FR-006]
- [x] CHK013 Are both supported content sizes explicitly required for responsive navigation and Options evidence? [Measurability, Spec §FR-006]
- [x] CHK014 Are conventional-wheel sensitivity levels, affected surfaces, persistence, nesting, and keyboard behavior explicit? [Completeness, Spec §FR-007]
- [x] CHK015 Is optional touchpad hardware handled without weakening mandatory conventional-wheel evidence? [Edge Case, Spec §FR-008]

## Structured Tables

- [x] CHK016 Is the minimum populated-row count measurable for each applicable table view? [Measurability, Spec §FR-009, §FR-010]
- [x] CHK017 Are every Tasks header and the Enabled/Lifecycle distinction named explicitly? [Clarity, Spec §FR-009]
- [x] CHK018 Are Schedule and Activity headers, severity casing, glyph/text semantics, and supported event states specified? [Completeness, Spec §FR-010]
- [x] CHK019 Are truncation disclosure, frozen headers, no horizontal overflow, and both window sizes covered? [Coverage, Spec §FR-009, §FR-010]
- [x] CHK020 Are hover, focus, selection, parity, live-refresh identity, and view-specific actions covered without conflating the views? [Consistency, Spec §FR-009, §FR-010]

## Collector and Gate Safety

- [x] CHK021 Is exactly-one placeholder/template generation required for every canonical observation? [Measurability, Spec §FR-011]
- [x] CHK022 Are overwrite, duplication, malformed input, attachment, and non-passing-status behaviors specified fail-closed? [Coverage, Spec §FR-012, §FR-013]
- [x] CHK023 Does fixture coverage require every new semantic rule without creating a production bypass? [Security, Spec §FR-014]
- [x] CHK024 Are existing installer, uninstall, access, task, error, and window checks explicitly preserved? [Regression, Spec §FR-001]
- [x] CHK025 Is any corrective product work gated by reproducible test-first evidence? [Consistency, Spec §FR-020]

## Delivery and Issue Traceability

- [x] CHK026 Are all eight desktop issues mapped to formal observations while #98 and #96 retain lifecycle traceability? [Traceability, Spec §FR-015]
- [x] CHK027 Are local MSI identity and compiled-inspection fields exhaustive and measurable? [Completeness, Spec §FR-016–FR-018]
- [x] CHK028 Is the pre-push stop condition explicit about unfinished attended, formal, issue, and publication work? [Clarity, Spec §FR-021]
- [x] CHK029 Is Post-v1 diagnostic work explicitly excluded? [Scope, Spec §FR-022]
- [x] CHK030 Does issue closure depend on individual acceptance evidence rather than bundled association? [Governance, Spec User Story 4]

## Notes

- All 30 release-evidence requirements checks passed on 2026-09-03.
- The checklist is for requirements quality. Runtime verification is defined in the later quickstart and task artifacts.
