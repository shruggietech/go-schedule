# Portability & Operator-Safety Requirements Checklist: Feature 006

**Purpose**: Validate that the requirements governing cross-platform parity, the
prerequisite contract, concurrency, bounded execution, measurement honesty, and repository
hygiene are complete, unambiguous, and objectively verifiable — *before* implementation.
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)
**Depth**: Release gate (reviewed at the single pre-push halt)
**Audience**: Reviewer

These items test the requirements, not the code. An unchecked item means the spec needs a
sentence, not that a script needs a fix.

## Twin Parity (PowerShell ↔ POSIX shell)

- [x] CHK001 - Is "equivalent behavior" between the twins defined precisely enough to be
      falsifiable, or does it rely on the reader's judgment? [Measurability, Spec §FR-015]
- [x] CHK002 - Is the naming correspondence between the two option styles stated as a rule
      rather than left to per-option inspection? [Clarity, Spec §FR-015]
- [x] CHK003 - Are requirements defined for what happens when a twin is run on the platform
      it does not primarily target? [Coverage, Spec §Assumptions]
- [x] CHK004 - Do the requirements specify whether the two twins must produce
      *byte-identical* stored records or merely equivalent ones? [Ambiguity, Spec §SC-004]
- [x] CHK005 - Is there a stated requirement that both twins share one implementation of the
      prerequisite-resolution logic, or could they legitimately diverge? [Gap, Spec §FR-016]
- [x] CHK006 - Are requirements defined for the timestamp format and precision each twin
      records, so that records from different twins are comparable? [Gap, Spec §FR-002]

## Platform-Conditional Probes

- [x] CHK007 - Does the spec require the PowerShell twin to avoid assuming platform-exclusive
      host-inspection commands are available? [Completeness, Spec §Edge Cases]
- [x] CHK008 - Is "a probe that cannot run" distinguished in the requirements from "a probe
      that ran and found nothing"? [Ambiguity, Spec §FR-009]
- [x] CHK009 - Are the requirements clear on whether a *partially* successful probe records
      partial data or is treated as absent? [Clarity, Spec §FR-009]
- [x] CHK010 - Is the required fallback ordering among alternative host-inspection commands
      specified, or left to implementation choice? [Gap, Spec §FR-007]
- [x] CHK011 - Are requirements defined for how an absent datum is represented in stored
      records so that queries can distinguish it from a legitimate zero? [Gap, Spec §FR-009]
- [x] CHK012 - Does the spec state which platforms are in scope for the host-inspection
      script, and which are explicitly excluded? [Coverage, Spec §Assumptions]

## Exit-Code Contract

- [x] CHK013 - Are all three exit conditions defined mutually exclusively, with no scenario
      that could reasonably map to two of them? [Consistency, Spec §FR-019]
- [x] CHK014 - Is it specified which exit code a *deliberately induced* failure produces, and
      whether it can collide with the prerequisite-failure code? [Conflict, Spec §FR-006 and
      §FR-019]
- [x] CHK015 - Are the requirements explicit that a missing prerequisite is a usage-class
      failure rather than a runtime failure? [Clarity, Spec §FR-017]
- [x] CHK016 - Is the interaction between the exit-code contract and the scheduler's own
      success/failure interpretation documented? [Dependency, Spec §US1 scenario 4]
- [x] CHK017 - Does the spec require that a recorded outcome and the process exit code always
      agree? [Gap, Spec §FR-006]

## Prerequisite Detection & Opt-In Installer

- [x] CHK018 - Is the resolution order stated as strictly ordered precedence rather than an
      unordered set of places to look? [Clarity, Spec §FR-016]
- [x] CHK019 - Are requirements defined for what happens when a located tool exists but is
      unusable — wrong architecture, not executable, too old? [Gap, Spec §FR-016]
- [x] CHK020 - Is a minimum acceptable version of the external tool specified, or is any
      version presumed acceptable? [Gap, Spec §FR-016]
- [x] CHK021 - Does the spec state that checksum verification precedes installation, not
      merely that it occurs? [Clarity, Spec §FR-018]
- [x] CHK022 - Are requirements defined for the disposition of a downloaded artifact that
      fails verification? [Gap, Spec §FR-018]
- [x] CHK023 - Is the no-network-without-explicit-opt-in rule stated as an absolute
      constraint covering every script, including the reader? [Completeness, Spec §FR-018]
- [x] CHK024 - Are the requirements clear on where the pinned checksums live and what
      obligation exists to update them? [Gap, Spec §FR-018]
- [x] CHK025 - Is the behavior on an unsupported platform-architecture combination
      distinguishable in the requirements from the tool-simply-not-installed case?
      [Ambiguity, Spec §FR-018a and §FR-017]

## Concurrency & Bounded Execution

- [x] CHK026 - Is "tolerate concurrent writers" quantified with a bound — how long to wait,
      how many times to retry — or left unmeasurable? [Measurability, Spec §FR-021]
- [x] CHK027 - Are requirements defined for what happens when contention *exceeds* that
      bound? [Coverage, Gap, Spec §FR-021]
- [x] CHK028 - Does the spec require that the bounded-loop guarantee hold when *both* bounds
      are omitted by the caller? [Edge Case, Spec §FR-004]
- [x] CHK029 - Are the requirements explicit about which bound wins when a maximum count and
      a maximum duration would end the loop at different moments? [Ambiguity, Spec §FR-004]
- [x] CHK030 - Is there a requirement covering interaction between the deliberate-slowness
      knob and the loop bounds — can a single slow run overrun the duration bound?
      [Conflict, Spec §FR-004 and §FR-005]
- [x] CHK031 - Are requirements defined for interrupted execution — a run terminated
      mid-flight by the scheduler or the operator? [Gap, Exception Flow]

## Measurement Honesty

- [x] CHK032 - Is the requirement that a drift figure never appears without its source stated
      as covering *every* surface it appears on, including documentation examples?
      [Completeness, Spec §FR-003]
- [x] CHK033 - Are the three expected-moment sources defined so that exactly one applies to
      any given record? [Consistency, Spec §FR-003 and §FR-003a]
- [x] CHK034 - Is the unreliability threshold for boundary-derived drift quantified rather
      than described qualitatively? [Measurability, Spec §Edge Cases]
- [x] CHK035 - Does the spec require the reader to disclose how many records a query
      *excluded*, not just what it included? [Clarity, Spec §Edge Cases]
- [x] CHK036 - Are requirements defined for records whose expected-moment source differs
      within a single query result set? [Gap, Spec §FR-013]
- [x] CHK037 - Is the relationship between this feature's derived drift figure and the
      project's own dispatch budget stated as comparison rather than compliance?
      [Ambiguity, Spec §SC-002]

## Repository Hygiene

- [x] CHK038 - Is the tracking rule for the agent configuration directory expressed as
      default-exclude-with-narrow-inclusion, rather than a list of things to exclude?
      [Clarity, Spec §FR-026]
- [x] CHK039 - Are requirements defined for what happens when a credential-bearing file
      appears *inside* the newly tracked subtree? [Edge Case, Gap, Spec §FR-026]
- [x] CHK040 - Does the spec state the obligation to keep vendored copies in step with their
      upstream sources, or is drift silently accepted? [Gap, Spec §Assumptions]
- [x] CHK041 - Are the requirements explicit that the produced databases and any locally
      installed tool are excluded from tracking under *all* the paths they can occupy?
      [Completeness, Spec §FR-027]
- [x] CHK042 - Is the change to repository configuration identified as touching a pinned
      artifact, with the resulting recording obligation stated? [Traceability, Gap]

## Dependencies & Assumptions

- [x] CHK043 - Is the assumption that the scheduler injects no context validated against the
      code rather than asserted? [Assumption, Spec §Assumptions]
- [x] CHK044 - Are the requirements clear that no product behavior is changed, and is that
      claim stated in a way a reviewer can confirm? [Measurability, Spec §Assumptions]
- [x] CHK045 - Is the claim that no continuous-integration change is needed supported by a
      stated reason? [Assumption, Spec §Assumptions]
- [x] CHK046 - Are the governing authoring conventions for the POSIX twins identified, given
      that no separate standard exists for them? [Dependency, Spec §FR-025]

## Notes

Items CHK014, CHK025, CHK030, and CHK037 mark places where two requirements could be read as
disagreeing. These are the highest-value items in the list: a conflict resolved at spec time
costs a sentence, and the same conflict resolved after implementation costs a rewrite of both
twins plus their tests.

CHK042 exists because the repository-configuration change carries a process obligation the
functional requirements do not mention, and process obligations are exactly what gets
forgotten between the plan and the commit.

## Resolution Record (2026-07-23)

All 46 items pass. Two passes were required.

**Pass 1 — 39/46.** The seven failures were real gaps, not wording quibbles.

**Pass 2 — 46/46**, after amending the spec. The four items flagged in the Notes above as
potential conflicts all turned out to be genuine, and each is now resolved explicitly rather
than by a reading that happens to reconcile them:

- CHK014 — the induced-failure knob could have been handed the code reserved for unmet
  prerequisites, making a deliberate failure indistinguishable from a missing tool. FR-006
  now forbids the reserved codes.
- CHK025 — an unsupported architecture and an uninstalled tool produced the same exit and
  the same advice, sending a maintainer to install something that was never going to work.
  FR-018a now separates them.
- CHK030 — the duration bound and the deliberate-slowness knob could each be read as
  overriding the other. FR-004a picks one, and states why interrupting a run mid-write is
  the worse option.
- CHK037 — SC-002 could have been read as this feature asserting compliance with the
  scheduler's dispatch budget. It measures and compares; it does not certify.

The seven pass-1 gaps produced FR-018b, FR-018c, FR-018d, FR-021c, FR-021d, FR-021e, and the
vendoring-drift sentence in Assumptions. The most consequential is **FR-021c**: writing the
beat once at the end of the run, rather than twice, halves contention and makes an
interrupted run present as a missed firing — which is the honest signal, since a maintainer
cannot act differently on a run that vanished mid-flight than on one that never started.

