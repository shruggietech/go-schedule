# Research: Windows Release Qualification

## Decision 1: Extend the one promotion gate

**Decision**: Add S044/S045 desktop observations to `internal/releasegate` and the existing evidence bundle instead of creating a second visual-acceptance document or command.

**Rationale**: Promotion currently has one fail-closed dependency and one candidate-bound archive. A second gate would introduce ordering, identity, and partial-pass ambiguity. The remaining issues share environments and evidence, so one bundle is also the most efficient attended workflow.

**Alternatives considered**:

- Leave the new requirements as prose in `test/windows/README.md`. Rejected because prose can be skipped while promotion still succeeds.
- Add a second archive and validator. Rejected because candidate identity and attachment rules would be duplicated.
- Close the issues from headless tests. Rejected because the issues explicitly require native rendering and physical-input evidence.

## Decision 2: Eleven outcome-oriented observations

**Decision**: Use eleven `desktop.` scenarios organized around appearance, interaction, navigation/Options, input, and the two table families. Every visual family has separate 96-DPI and scaled-DPI observations; physical scroll input remains a standard-DPI scenario.

**Rationale**: One scenario per checkbox would be unmanageable, while one generic `desktop.acceptance` boolean would be unauditable. Eleven scenarios align with the implementation boundaries, prevent standard-DPI evidence from standing in for scaled rendering, and allow an individual issue to point at a small, coherent evidence set.

**Alternatives considered**:

- One observation for each open issue. Rejected because #101/#106 and #104/#105/#109 overlap heavily and would duplicate screenshots.
- One observation per palette, DPI, size, and control permutation. Rejected because it would create dozens of fragments without improving issue-level meaning.

## Decision 3: Flat typed metrics with exact normalized sets

**Decision**: Continue using `map[string]any`, validate booleans/numbers through existing helpers, and add a normalized exact-set helper for palettes, sizes, states, surfaces, fonts, headers, and semantic values.

**Rationale**: This preserves schema version 1 and its strict top-level decoder, keeps hand-authored fragments approachable, and gives deterministic diagnostics. Exact sets reject duplicate or omitted coverage without requiring a new JSON schema dependency.

**Alternatives considered**:

- Add nested typed desktop records and bump the evidence schema. Rejected as a migration burden with no product consumer beyond this validator.
- Accept free-form summary prose. Rejected because promotion could not prove completeness.
- Require ordered comma-separated strings. Rejected because ordering has no semantic value and creates needless operator errors.

## Decision 4: Require screenshots for every desktop observation

**Decision**: Each `desktop.` observation references at least one raster-image attachment validated from its content signature rather than trusted metadata or filename; structured metrics remain mandatory alongside it.

**Rationale**: Metrics prove checklist completeness, but font sharpness, restrained color, clipping, row appearance, and state legibility require native visual review. Content inspection excludes mislabeled text/vector files, and hashing attachments binds the review evidence to the bundle.

**Alternatives considered**:

- Metrics only. Rejected because operator booleans would have no reviewable visual support.
- Screenshots only. Rejected because images do not prove all palettes, sizes, row counts, input devices, persistence, or interaction sequences.

## Decision 5: Treat touchpad as explicitly optional, not silently skipped

**Decision**: Conventional mouse-wheel evidence is mandatory. The scroll observation also records `touchpad_available`; a true value requires passing fine-delta behavior, while false requires a non-empty unavailability reason.

**Rationale**: Issue #111 says precision touchpad evidence is required when the hardware is available. A pass-by-omission is misleading, while making hardware availability a global release blocker exceeds the approved acceptance criteria.

**Alternatives considered**:

- Require a touchpad on every release. Rejected because the issue explicitly conditions this check on availability.
- Omit touchpad fields when absent. Rejected because reviewers cannot distinguish unavailable hardware from a forgotten check.

## Decision 6: Keep local demo and formal candidate identities separate

**Decision**: Repeat the S043 `local-demo` build pattern with an S047 marker, then retain all formal observations for the later workflow-staged MSI.

**Rationale**: The user requires a halt before push. The formal artifact cannot exist until reviewed source is merged and an authorized tag runs the Release workflow. Calling the local build exact-candidate evidence would contradict the S040 promotion contract.

**Alternatives considered**:

- Push or tag before the halt. Rejected by explicit user instruction and the constitution.
- Skip pre-push native testing. Rejected because S043 already showed that real Windows use finds defects headless tests cannot.

## Decision 7: Do not reopen closed #94 automatically

**Decision**: Extend the S040 mechanism and report that #96's child checklist is stale, but do not reopen #94 or claim its unchecked criteria passed.

**Rationale**: The issue was manually closed by a repository member. Reopening would rewrite current project-management intent without a clear authorization. The open leaf issues independently retain the native criteria S047 enforces.

## Decision 8: Preserve the mandated Spec Kit order despite a prerequisite defect

**Decision**: Run `setup-plan.ps1` only to create the untouched plan template before `/speckit-checklist`, then populate the plan after the checklist.

**Rationale**: The installed checklist prerequisite rejects a missing plan even though both the command skill and autopilot protocol place checklist before plan. This established project workaround preserves semantic execution order without modifying unrelated Spec Kit tooling.

**Alternatives considered**:

- Populate the plan before the checklist. Rejected because it violates the required sequence.
- Skip the checklist. Rejected because S047 has substantial release and native evidence risks.
- Repair Spec Kit inside S047. Rejected as unrelated scope.

## Repository Findings

- S040's validator requires 36 fixed observations and already enforces exact candidate, environment, status, attachment, task-run, lifecycle, and native window rules.
- S044 and S045 added detailed native runbook prose but did not add their desktop observations to `RequiredScenarioIDs`, so promotion can currently succeed without those later acceptance outcomes.
- `Invoke-ReleaseCandidateAttended.ps1` creates templates only for setup and removal observations. Generalizing that helper is smaller than a second collector and preserves overwrite safety.
- The checked-in passing fixture is explicitly synthetic and remains suitable for semantic regression tests if the production validator continues to reject its evidence class.
- The development host has previously built and inspected a WiX 6.0.2 local demo successfully. Destructive lifecycle automation remains reserved for a disposable elevated environment.
- Issue #94 is closed, but #96 still lists it unchecked and its own completion criteria remain open. S047 must surface that state rather than infer either completion or authority to reopen it.
