# Installer Lifecycle Requirements Checklist: Windows Setup Lifecycle Control

**Purpose**: Review the completeness, clarity, consistency, and measurability of the installer UX and destructive-removal requirements
**Created**: 2026-09-02
**Feature**: [spec.md](../spec.md)

**Audience and depth**: Formal pull-request review checklist covering installer UX and destructive cleanup safety.

## Requirement Completeness

- [x] CHK001 Are requirements defined for fresh install, maintenance modification, repair, upgrade, uninstall, rollback, cancellation, and unattended execution? [Completeness, Spec §Context and Scope]
- [x] CHK002 Are all four shortcut-selection states and all four completion-action states explicitly required? [Coverage, Spec §SC-001, Spec §SC-002]
- [x] CHK003 Does the removal inventory distinguish software, application data, user preferences, and security state? [Completeness, Spec §FR-013]
- [x] CHK004 Are preserve, wipe, cancel, reinstall, multiple-profile, locked-file, and partial-cleanup outcomes specified? [Coverage, Spec §FR-014–FR-025]

## Requirement Clarity

- [x] CHK005 Are shortcut defaults and completion-action defaults unambiguous for a fresh install? [Clarity, Spec §FR-002, Spec §FR-008]
- [x] CHK006 Is the documentation target and its system-handler behavior stated exactly? [Clarity, Spec §FR-009]
- [x] CHK007 Is destructive consent defined separately for attended and unattended removal? [Clarity, Spec §FR-015]
- [x] CHK008 Is the meaning of “all application-related user data” bounded to declared product-owned roots and registered profiles? [Clarity, Spec §Assumptions]
- [x] CHK009 Are unsafe deletion-target classes enumerated rather than described only as “safe”? [Clarity, Spec §FR-018]

## Requirement Consistency

- [x] CHK010 Are independent completion checkboxes consistent across user scenarios, functional requirements, and measurable outcomes? [Consistency, Spec §US2, Spec §FR-007–FR-010, Spec §SC-002]
- [x] CHK011 Is preserve-by-default behavior consistent across attended uninstall, unattended uninstall, and reinstall scenarios? [Consistency, Spec §US3, Spec §FR-012–FR-015, Spec §SC-003]
- [x] CHK012 Does post-commit cleanup timing remain consistent with cancel, failure, and rollback protections? [Consistency, Spec §Edge Cases, Spec §FR-017]
- [x] CHK013 Is shared security-state preservation consistent with the declared wipe scope and out-of-scope section? [Consistency, Spec §FR-022, Spec §Out of Scope]

## Acceptance Criteria Quality

- [x] CHK014 Can shortcut presence, completion action counts, preserved bytes, erased roots, and untouched controls be measured objectively? [Measurability, Spec §SC-001–SC-007]
- [x] CHK015 Does every refusal or partial-cleanup requirement have a non-success result and retained evidence expectation? [Measurability, Spec §FR-020, Spec §SC-007]
- [x] CHK016 Are issue-level completion and downstream #94 release proof kept distinct? [Traceability, Spec §Dependencies and Traceability]

## Non-Functional and Edge-Case Coverage

- [x] CHK017 Are elevation boundaries specified for completion actions and privileged cleanup? [Security, Spec §FR-009, Spec §FR-018]
- [x] CHK018 Are unrelated files, exports, adjacent shortcuts, reparse points, redirected locations, and orphaned profiles covered explicitly? [Edge Case, Spec §Edge Cases, Spec §FR-018–FR-019]
- [x] CHK019 Are diagnostics required to remain actionable without exposing unrelated user information? [Security, Spec §US4, Spec §FR-020]
- [x] CHK020 Are S039 automated verification and downstream #94 attended candidate evidence both represented without substituting one for the other? [Dependency, Spec §FR-024–FR-026, Spec §Out of Scope]

## Notes

- All 20 requirement-quality questions passed after the clarify scan.
- Focus areas were installer UX and destructive cleanup safety; no unresolved ambiguity requires operator input before planning.
