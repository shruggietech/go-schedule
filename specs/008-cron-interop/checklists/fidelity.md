# Conversion Fidelity Checklist: Cron interoperability

**Purpose**: Validate that the requirements governing cron conversion are
complete, unambiguous, and leave no room for a silent approximation — in either
direction, and between what a preview promises and what an import produces.
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)

**Note**: This is a requirements-quality checklist. Every item asks whether the
specification says enough, clearly enough, to be implemented and verified — not
whether the code works.

## Requirement Completeness

- [x] CHK001 Is the accepted cron dialect named precisely enough to bound the
      work — field count, field order, and which named shorthands are in?
      [Completeness, Spec §FR-001]
- [x] CHK002 Is the set of explicitly declined inputs enumerated rather than
      described as "non-standard extensions"? [Completeness, Spec §FR-002]
- [x] CHK003 Are requirements defined for what an import carries across as the
      task payload, field by field (command, arguments, working directory)?
      [Completeness, Spec §FR-006]
- [x] CHK004 Are requirements defined for crontab content that is neither a
      schedule line nor a comment — variable assignments, `MAILTO`, blank lines?
      [Completeness, Spec §FR-007, Edge Cases]
- [x] CHK005 Is the composition of the import summary specified — which totals
      it reports and which fidelity statements it must make? [Completeness,
      Spec §FR-009, §FR-010]
- [x] CHK006 Are requirements defined for the naming of imported tasks, given
      that crontabs do not name jobs? [Assumption, Spec §Assumptions]
- [x] CHK007 Are requirements defined for reading a crontab from standard input
      as well as from a file? [Completeness, Spec §FR-004]

## Requirement Clarity

- [x] CHK008 Is "never approximate" stated as a testable rule rather than an
      aspiration — i.e. is there a defined observable outcome (a named refusal)
      for every input the system will not convert? [Clarity, Spec §FR-002]
- [x] CHK009 Is the single conversion route (cron → phrase → schedule) stated as
      a requirement, so that a second route would be a spec violation rather
      than a design preference? [Clarity, Spec §FR-003a]
- [x] CHK010 Is "cron is never an authoring syntax" expressed as a checkable
      prohibition naming the surfaces it binds? [Clarity, Spec §FR-014]
- [x] CHK011 Is the divides-evenly rule for step values stated with its range
      basis, so that `*/15` on minutes and `*/5` on hours can each be judged
      without interpretation? [Clarity, Spec §FR-003b, Edge Cases]
- [x] CHK012 Is the day-of-month/day-of-week decline stated with the reason it
      exists (cron's OR versus the recurrence model's AND), so a later reader
      cannot mistake it for an unimplemented case? [Clarity, Spec
      §Clarifications, §FR-003b]

## Requirement Consistency

- [x] CHK013 Do the import, explain, and export requirements use one vocabulary
      for the same concepts (declined / unsupported / warning / error)?
      [Consistency]
- [x] CHK014 Are the preview and the real import specified as producing an
      identical report, with creation as the only difference? [Consistency,
      Spec §FR-005, §SC-002a]
- [x] CHK015 Do the export requirements and the import requirements agree on
      what "cron can carry" means, so that a round trip cannot decline on one
      side what it accepted on the other? [Consistency, Spec §FR-012, §FR-013]
- [x] CHK016 Is the treatment of a disabled task on export consistent with the
      general refusal rule rather than a separate special case? [Consistency,
      Spec §FR-012a]

## Acceptance Criteria Quality

- [x] CHK017 Is the round-trip success criterion measurable — does it name the
      comparison window and the boundaries it must cross? [Measurability,
      Spec §FR-013, §SC-003]
- [x] CHK018 Can "no line is silently dropped or silently altered" be verified
      by counting, rather than by inspection? [Measurability, Spec §SC-002]
- [x] CHK019 Is the fidelity statement the import must make specified concretely
      enough that its absence would fail a test? [Measurability, Spec §FR-008,
      §FR-009]

## Scenario Coverage

- [x] CHK020 Are requirements defined for a crontab that is entirely
      unsupported — does the run succeed with zero creations, or fail?
      [Coverage, Exception Flow]
- [x] CHK021 Are requirements defined for a partially supported crontab in a
      real (non-preview) import, including whether supported lines are still
      created? [Coverage, Edge Cases]
- [x] CHK022 Are requirements defined for a malformed expression, distinct from
      a well-formed but unsupported one? [Coverage, Spec §US3 scenario 3]
- [x] CHK023 Are requirements defined for repeating an import of the same
      crontab? [Coverage, Spec §Assumptions, Edge Cases]
- [x] CHK024 Are requirements defined for an export of an empty task set?
      [Coverage, Gap]
- [x] CHK025 Are recovery requirements defined for an import that fails partway
      through — is a partial import acceptable, and is it visible? [Recovery,
      Gap]

## Dependencies & Assumptions

- [x] CHK026 Is the assumption that the conversion targets Vixie/ISC-style cron,
      and not Quartz or systemd timers, recorded rather than implied?
      [Assumption, Spec §Assumptions]
- [x] CHK027 Is the dependency on the schedule grammar's new by-date and yearly
      forms stated, so that the ordering within the feature is not left to
      chance? [Dependency, Spec §FR-015–§FR-017]
- [x] CHK028 Is the decision not to apply crontab environment assignments
      recorded with its rationale? [Assumption, Spec §Assumptions]

## Ambiguities & Conflicts

- [x] CHK029 Does any requirement leave open whether a declined line is an error
      (non-zero exit) or a reported outcome (zero exit)? [Ambiguity, Gap]
- [x] CHK030 Is there any remaining requirement that would permit the export to
      omit a task entirely rather than refuse it visibly? [Conflict, Spec
      §FR-012]

## Notes

- Check items off as completed: `[x]`
- An unchecked item is a spec defect to fix before `/speckit-plan`, not a code
  defect.
