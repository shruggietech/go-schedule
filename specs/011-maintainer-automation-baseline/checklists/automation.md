# Automation Requirements Checklist: Maintainer Automation Baseline

**Purpose**: Test whether the automation and verification requirements are
complete, precise, consistent, and ready for implementation
**Created**: 2026-08-26
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are all workflow files and all external action references included
  in the runtime-upgrade scope? [Completeness, Spec §Scope in, FR-001]
- [x] CHK002 Are the required local gates enumerated rather than summarized as
  an ambiguous "CI parity" promise? [Completeness, Spec §FR-005]
- [x] CHK003 Are both aggregate and individual verification entry points covered
  by the single-source requirement? [Completeness, Spec §FR-004, FR-007]
- [x] CHK004 Are contributor, autopilot, Makefile, CI, and changelog surfaces all
  named where consistency or updates are required? [Completeness, Spec §FR-010,
  FR-012]

## Requirement Clarity

- [x] CHK005 Is "supported runtime" bounded by the hosted automation platform
  and an explicit prohibition on Node 20? [Clarity, Spec §FR-002]
- [x] CHK006 Is "non-mutating" made observable through a clean-checkout status
  assertion and an unformatted-file failure scenario? [Clarity, Spec §US2,
  SC-003]
- [x] CHK007 Is the complete local gate set precise enough to distinguish format
  checking from format rewriting and threshold coverage from informational
  coverage? [Clarity, Spec §FR-005]
- [x] CHK008 Is failure behavior defined for both a failing child gate and an
  unavailable prerequisite? [Clarity, Spec §FR-006, FR-009]

## Requirement Consistency

- [x] CHK009 Do the local verification requirements preserve the constitution's
  mandatory lint, race, coverage, and foreground-execution rules? [Consistency,
  Constitution Principles I-II, Spec §FR-005]
- [x] CHK010 Does static release-workflow validation remain consistent with the
  prohibition on tags, releases, and external side effects? [Consistency, Spec
  §FR-013, SC-007]
- [x] CHK011 Do action compatibility requirements preserve workflow behavior
  rather than using runtime modernization as permission to redesign jobs?
  [Consistency, Spec §FR-003, Assumptions]
- [x] CHK012 Is the Makefile contract consistent with the documented distinction
  between mutating convenience commands and verification gates? [Consistency,
  Spec §FR-008, US2]

## Acceptance Criteria Quality

- [x] CHK013 Can complete action-runtime coverage be measured as a percentage of
  inventoried references? [Measurability, Spec §SC-001]
- [x] CHK014 Can complete local-gate coverage be measured against one explicit
  required-gate inventory? [Measurability, Spec §SC-003]
- [x] CHK015 Does the deliberate-negative validation require every child gate to
  propagate a named non-zero failure? [Acceptance Criteria, Spec §SC-004]
- [x] CHK016 Does the drift criterion cover both action regression and local-gate
  omission rather than only one side of the automation contract?
  [Acceptance Criteria, Spec §SC-005]

## Scenario and Edge-Case Coverage

- [x] CHK017 Are ordinary CI, release-workflow audit, clean local verification,
  unformatted input, child-gate failure, and mutating-target discoverability all
  represented by acceptance scenarios? [Coverage, Spec §US1-US2]
- [x] CHK018 Are missing C toolchains, missing POSIX shell access, older base Go
  toolchains, action contract breaks, release side effects, and future CI jobs
  addressed as boundary cases? [Edge Cases, Spec §Edge Cases]
- [x] CHK019 Is the behavior for a gate that cannot run explicit enough to prevent
  a partial suite from being reported as green? [Recovery, Spec §FR-009]

## Dependencies, Boundaries, and Traceability

- [x] CHK020 Are issues #21 and #41 explicitly traceable while #23, #38, #39,
  #40, and #42 are clearly excluded? [Traceability, Spec §Traceability, Scope
  out]
- [x] CHK021 Are the existing workflow contracts and supported shell assumptions
  documented rather than left implicit? [Assumption, Spec §Assumptions]
- [x] CHK022 Is every pinned-artifact change tied to the dated decision rule?
  [Dependency, Spec §FR-012]
- [x] CHK023 Is product behavior explicitly protected from this maintenance-only
  slice? [Boundary, Spec §Scope out, SC-007]

## Notes

- Depth: formal pre-implementation gate.
- Audience: maintainer and autopilot reviewer.
- Focus: workflow-runtime safety and CI/local verification parity.
- All 23 requirements-quality checks pass against the current specification.
