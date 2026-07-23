# Specification Quality Checklist: Maintainer Test Scripts and Vendored Skills

**Purpose**: Validate specification completeness and quality before proceeding to planning
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)

## Content Quality

- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

## Requirement Completeness

- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Success criteria are technology-agnostic (no implementation details)
- [x] All acceptance scenarios are defined
- [x] Edge cases are identified
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

## Feature Readiness

- [x] All functional requirements have clear acceptance criteria
- [x] User scenarios cover primary flows
- [x] Feature meets measurable outcomes defined in Success Criteria
- [x] No implementation details leak into specification

## Notes

Validation ran in two iterations.

**Iteration 1 findings, all resolved in the spec:**

1. *Implementation detail leak.* The first draft named the concrete artifacts —
   `Test-Heartbeat.ps1`, `sqlite3`, `drift_ms`, `PRAGMA journal_mode=WAL`, `.gitignore`,
   `.claude/skills/`. These are the approved plan's answers, not the specification's
   questions. Rewritten to describe the capability: "the heartbeat script", "the external
   command-line SQLite tool", "the difference between actual and scheduled moment",
   "tolerate concurrent writers by waiting and retrying", "the agent skills directory".
2. *Untestable success criteria.* SC-002 originally read "drift p99 under 100ms", which
   bakes in both a unit and a threshold the scheduler owns, not this feature. Restated as
   producing a figure "directly comparable against the project's documented dispatch
   budget" — this feature's job is to *measure*, not to *pass*.
3. *Missing edge case.* Nothing covered a schedule that does not supply an expected firing
   moment, which is the normal case when a script is invoked by hand. Added, with the
   inference fallback stated.

**Iteration 2: all items pass.** No [NEEDS CLARIFICATION] markers were needed — the
operator's four pre-execution decisions (skills-directory-only tracking, the vendored skill
set, detect-with-opt-in-installer, and the v0.5.0 release) removed every question that
would otherwise have qualified.

**Post-clarify re-validation (2026-07-23): 16/16 → 16/16, no regressions.** The clarify
pass added a `## Clarifications` section and five requirement changes; every checklist item
that passed before still passes, and the changes strengthened three of them.

**The boundary flagged during the first pass is now closed.** It read: FR-003 requires drift
against a scheduler-supplied scheduled moment, but whether `goschedd` exports that moment
into the spawned environment was unverified. It was then verified against
`internal/executor/executor.go` — **it does not**. A spawned task gets the inherited
environment plus its own configured variables and nothing else. This inverted an assumption
the spec had stated as fact, so the assumption was rewritten, FR-003 was split into a
three-source precedence with a mandatory source label, FR-002 gained a run-finish moment so
overlap became decidable, and two edge cases were added for the unmeasurable and
near-ambiguous cases. The full rationale, including the two rejected alternatives, is in the
spec's Clarifications section.

This is the finding to carry into planning: **drift here is derived, not reported.** The
implementation must keep it labelled as such everywhere it surfaces.
