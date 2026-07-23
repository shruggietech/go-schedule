# Calendar Correctness Checklist: missing-date policy and by-date grammar

**Purpose**: Validate that the requirements governing schedules addressing a
date that does not exist in every period are complete, unambiguous, and pinned
to real dates — and that they guarantee no existing schedule changes behavior.
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)

**Note**: This is a requirements-quality checklist. Every item asks whether the
specification says enough, clearly enough, to be implemented and verified — not
whether the code works.

## Requirement Completeness

- [x] CHK031 Are all three policy settings defined by what they produce, rather
      than by name alone? [Completeness, Spec §FR-019]
- [x] CHK032 Is the set of rule shapes the policy applies to enumerated —
      by-date monthly, yearly by-date, ordinal weekday? [Completeness,
      Spec §FR-021]
- [x] CHK033 Is "last valid" defined separately for a by-date rule and for an
      ordinal-weekday rule, since the two mean different things? [Completeness,
      Spec §FR-021]
- [x] CHK034 Are requirements defined for the new grammar forms' optional
      time-of-day clause, consistent with the existing forms? [Completeness,
      Spec §FR-015, §FR-016]
- [x] CHK035 Are requirements defined for how a task's policy is surfaced in
      every place a schedule is shown, rather than only in the editor?
      [Completeness, Spec §FR-022]
- [x] CHK036 Is the migration requirement stated as forward-only and
      non-destructive, naming what must not change? [Completeness, Spec §FR-026]

## Requirement Clarity

- [x] CHK037 Is "the period that has no matching date" defined unambiguously for
      each frequency (month for a monthly rule, year for a yearly rule)?
      [Clarity, Spec §FR-019]
- [x] CHK038 Is the roll-forward behavior specified precisely enough to settle
      where the run lands and whether the following period's own occurrence is
      affected? [Ambiguity, Spec §Edge Cases]
- [x] CHK039 Is the default stated as reproducing prior behavior exactly, rather
      than as "similar to" it? [Clarity, Spec §FR-020]
- [x] CHK040 Is the requirement that a description name its policy stated in
      terms of what it may not say, so the existing false "every month" label is
      unambiguously a defect? [Clarity, Spec §FR-023]
- [x] CHK041 Is the policy's inertness on date-less schedules stated as a
      requirement rather than left as an implication? [Clarity, Spec §FR-024]

## Requirement Consistency

- [x] CHK042 Do the grammar requirements and the policy requirements agree on
      which phrases can produce a period with no matching date? [Consistency,
      Spec §FR-015–§FR-021]
- [x] CHK043 Are the policy's interaction with daylight-saving resolution stated
      once and consistently, with this feature changing neither? [Consistency,
      Spec §FR-025]
- [x] CHK044 Do the command-line and graphical requirements describe the same
      three settings with the same meanings? [Consistency, Spec §FR-022]
- [x] CHK045 Is the independence of the policy and the schedule phrase stated
      consistently for both create and edit paths? [Consistency, Spec §FR-024a]

## Acceptance Criteria Quality

- [x] CHK046 Are the expected run times for each policy pinned to named real
      dates rather than described relationally? [Measurability, Spec §SC-004,
      §US2]
- [x] CHK047 Is "no existing task's run times change" expressed as something
      that can be measured against a pre-existing database? [Measurability,
      Spec §FR-026, §SC-006]
- [x] CHK048 Can "no description asserts a rule fires in every period when it
      does not" be checked mechanically across the descriptions the system
      produces? [Measurability, Spec §SC-005]
- [x] CHK049 Do the acceptance scenarios cover a full year including a non-leap
      February, rather than a single sampled month? [Coverage, Spec §SC-004]

## Scenario Coverage

- [x] CHK050 Are requirements defined for the 29 February case under all three
      policies? [Coverage, Spec §US2 scenario 4]
- [x] CHK051 Are requirements defined for the 30th and 31st in February, where
      fall-back and roll-forward differ by more than one day? [Coverage, Edge
      Cases]
- [x] CHK052 Are requirements defined for the fifth-weekday rule, whose existing
      behavior this feature must correct rather than merely extend? [Coverage,
      Spec §US2 scenario 5]
- [x] CHK053 Are requirements defined for a policy-resolved date that lands on a
      daylight-saving transition? [Coverage, Edge Cases]
- [x] CHK054 Are requirements defined for a database created before this feature
      — both its stored values and its computed run times? [Recovery, Spec
      §US2 scenario 7]
- [x] CHK055 Are requirements defined for setting a policy on a schedule that
      can never be affected by it? [Coverage, Spec §FR-024]

## Dependencies & Assumptions

- [x] CHK056 Is the exclusion of daylight-saving anchoring options recorded with
      the reason and the tracking issue, so the feature's boundary is legible?
      [Assumption, Spec §Assumptions]
- [x] CHK057 Is the assumption that the default policy makes the feature
      non-breaking stated explicitly? [Assumption, Spec §Assumptions]

## Ambiguities & Conflicts

- [x] CHK058 Does any requirement leave open whether the policy is stored per
      task or per schedule in a way that changes observable behavior?
      [Ambiguity, Gap]
- [x] CHK059 Is there any conflict between the requirement that descriptions
      name the policy and the requirement that a stored phrase round-trips
      unchanged? [Conflict, Spec §FR-018, §FR-023]

## Notes

- Check items off as completed: `[x]`
- An unchecked item is a spec defect to fix before `/speckit-plan`, not a code
  defect.
- Numbering continues from `fidelity.md` (CHK001–CHK030) so the two checklists
  can be cited without collision.
