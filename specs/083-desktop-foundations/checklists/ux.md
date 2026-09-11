# UX Requirements Checklist: Desktop Visual and Shell Foundations

**Purpose**: Review the completeness, clarity, consistency, and measurability of the cross-cutting desktop UX requirements before implementation
**Created**: 2026-09-11
**Feature**: [spec.md](../spec.md)

## Shared Visual Requirements

- [x] CHK001 Are control dimensions, spacing, contrast, and interaction states quantified for every shared action variant? [Completeness, Spec §FR-001 through FR-004]
- [x] CHK002 Are affirmative, destructive, primary, secondary, subtle, disabled, and pending meanings consistently distinguished without relying on color alone? [Consistency, Spec §FR-002 and FR-007]
- [x] CHK003 Are card and dialog spacing requirements tied to one shared scale and objective minimum action separation? [Clarity, Spec §FR-005 and FR-006]
- [x] CHK004 Is follow-system behavior defined for both interface surfaces and visible application imagery, including a live system-theme change? [Coverage, Spec §FR-008]

## Shell and Reflow Requirements

- [x] CHK005 Is the ownership of viewport scrolling explicit and consistent for every route? [Clarity, Spec §FR-009 and FR-010]
- [x] CHK006 Are minimum viewport, zoom, horizontal overflow, and retained-functionality criteria measurable? [Acceptance Criteria, Spec §FR-011 and SC-005]
- [x] CHK007 Are narrow-view alternatives bounded without hiding the Exit action or ordinary destinations? [Coverage, Spec §FR-010 and Edge Cases]

## Feedback Requirements

- [x] CHK008 Is routine feedback distinguished from correctable and ongoing errors by lifetime, placement, and dismissal behavior? [Consistency, Spec §FR-012 through FR-016]
- [x] CHK009 Are replacement, timeout, hover, focus, reduced-motion, and empty-message scenarios specified? [Coverage, Spec §FR-012 through FR-016 and Edge Cases]
- [x] CHK010 Are safe margins and prohibited overlaps defined for shell controls, focus, navigation, and window boundaries? [Completeness, Spec §FR-014]
- [x] CHK011 Are live-region priority and explicit dismissal requirements documented for each feedback class? [Accessibility, Spec §FR-016]

## Scope and Evidence

- [x] CHK012 Are the deferred page-specific issues named so shared-foundation work cannot silently absorb their redesign scope? [Boundary, Spec §FR-018 and Scope Boundaries]
- [x] CHK013 Are automated and native evidence requirements specific enough to cover themes, pointer, keyboard, scrolling, feedback, and reflow? [Measurability, Spec §FR-017 and SC-007 through SC-009]

## Notes

- Standard-depth checklist for author and pull-request reviewers, focused on shared visual, accessibility, shell, and feedback requirements.
