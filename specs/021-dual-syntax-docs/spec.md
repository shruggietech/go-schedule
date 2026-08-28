# Feature Specification: Dual-Syntax Product Documentation

**Feature Branch**: `codex/021-dual-syntax-docs`

**Created**: 2026-08-28

**Status**: Draft

**Input**: User description: "Complete issue #52's documentation follow-through for first-class human and cron schedule input, reconcile the product posture across every current surface, explicitly supersede obsolete boundary-only policy, and complete parent epic #50."

## Context

S018 through S020 delivered local string conversion, one shared human-or-cron
task-input boundary, source retention, and GUI adoption. Current documentation
was changed only where each implementation slice required it, leaving the
product introduction, command overview, API contracts, dialect guide, master
specification, and historical boundary-only decision inconsistent with shipped
behavior. S021 is the documentation completion slice for issues #52 and #50.
Phase 0 research also found that wildcard steps in calendar fields are silently
reduced to an unrelated daily schedule. S021 includes the narrow named-refusal
fix required to document and close the parent contract honestly. Cron fidelity
breadth remains issue #22.

## Clarifications

### Session 2026-08-28

- Q: Should completed historical specifications be rewritten to look as though they always supported cron? → A: Preserve delivery history and add explicit supersession notes that point to S019-S021; update only the still-authoritative master contract in place.
- Q: Should the guide imply full cron compatibility or describe the actual accepted subset? → A: Name the exact five-field input contract, explain named refusals and semantic traps, and keep unsupported breadth linked to #22.
- Q: Does this slice finish only #52, or also the parent epic? → A: The pull request uses `Closes #52` and `Closes #50` because S018-S021 together satisfy the parent acceptance criteria.
- Q: What should happen after research proves calendar-field wildcard steps are silently approximated? → A: Add the smallest correctness fix and regression tests so those lossy forms are refused by name; supporting them remains #22.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Choose Either Authoring Syntax Confidently (Priority: P1)

An operator arriving through the README, CLI help, GUI guidance, or API
contract understands that human phrases are the approachable default and
supported cron expressions are equally valid recurring task input.

**Why this priority**: The current anti-cron headline and scattered legacy copy
contradict a shipped product capability at the first point of discovery.

**Independent Test**: Follow the primary task-creation documentation from each
surface and confirm it presents copy/pasteable, equivalent human and cron
examples with consistent acceptance, preview, retention, and error language.

**Acceptance Scenarios**:

1. **Given** a reader at the product introduction, **When** they evaluate scheduling input, **Then** human and supported cron syntax are both welcomed without implying that either creates a second execution engine.
2. **Given** a CLI user, **When** they inspect root and task help, **Then** the help names both accepted forms and retains a concise human-first example.
3. **Given** an API or GUI user, **When** they inspect the relevant contract, **Then** syntax detection, optional explicit identity, exact expression retention, preview behavior, and named refusal are described consistently.
4. **Given** either syntax, **When** a reader copies the documented task-creation example, **Then** both examples represent the same weekday-at-09:00 recurrence.

---

### User Story 2 - Understand Cron Fidelity Before Migrating (Priority: P2)

A cron user can determine what a schedule expression means, which five-field
features are faithfully accepted, which constructs are refused, and when to use
conversion, explanation, import, or export.

**Why this priority**: Cron-like text can carry different semantics across
dialects. Clear fidelity boundaries prevent silent scheduling mistakes.

**Independent Test**: Use the cron guide alone to identify field count and
order, names/macros/extensions, timezone handling, day-of-month/day-of-week
semantics, field-local steps, refusal behavior, and the difference between an
expression and a full crontab line/file.

**Acceptance Scenarios**:

1. **Given** a five-field cron expression, **When** a reader checks the dialect contract, **Then** supported fields, names, macros, lists/ranges/steps, and extensions are stated precisely rather than described as generic "cron".
2. **Given** restricted day-of-month and day-of-week fields, **When** a reader checks fidelity, **Then** traditional OR semantics and go-schedule's named refusal are explicit.
3. **Given** a stepped field or timezone-sensitive schedule, **When** a reader checks fidelity, **Then** field-local step meaning and timezone/DST execution semantics are explicit.
4. **Given** a crontab file, **When** a reader chooses a workflow, **Then** task input, `convert`, `explain`, `import`, and `export` are distinguished by source, result, daemon use, and mutation behavior.

---

### User Story 3 - Leave One Coherent Policy of Record (Priority: P3)

A maintainer can distinguish current requirements from completed historical
decisions and can prove that no active policy still mandates cron rejection.

**Why this priority**: Contradictory specifications and policy tests invite a
future regression even after user-facing prose is fixed.

**Independent Test**: Search current documentation, help, contracts, source
comments, specifications, and policy tests for boundary-only language; confirm
current artifacts use the dual-syntax contract and historical S008 material is
clearly superseded without falsifying its chronology.

**Acceptance Scenarios**:

1. **Given** the master specification, **When** a maintainer reads schedule-input requirements, **Then** human input remains usable without requiring cron while supported cron is also first-class.
2. **Given** the completed S008 specification and tasks, **When** a maintainer encounters the boundary-only decision, **Then** a prominent note identifies the later slices that supersede it while the historical record remains intact.
3. **Given** documentation-policy automation, **When** categorical no-cron positioning is reintroduced into current surfaces, **Then** the offline documentation gate fails with an actionable reason.
4. **Given** the completed S021 pull request, **When** GitHub applies its closing keywords, **Then** issues #52 and #50 close while #22 remains open for dialect breadth.

### Edge Cases

- Historical changelog entries and completed task text can accurately describe
  an earlier boundary-only release; they are not current product guidance and
  must not be globally rewritten.
- `internal/schedule.Parse` remains a human-language parser and its comments may
  accurately say so; multi-syntax callers use the separate central boundary.
- "Cron expression" means schedule fields only. A full crontab line may also
  contain a command, environment, shell, user, or other file-level context.
- Supported descriptors or extensions must not be described more broadly than
  the parser accepts.
- DOM/DOW behavior must not imply intersection or claim faithful conversion
  where traditional cron uses OR and the product refuses the combination.
- Examples must avoid shell quoting that works on only one supported platform,
  or clearly label platform-specific forms.
- Links and anchors must resolve in both repository Markdown and the published
  documentation site.
- Non-ASCII arrows or punctuation used in examples must remain valid UTF-8
  without BOM or mojibake.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The README and product-level CLI description MUST replace categorical anti-cron positioning with concise dual-syntax positioning that keeps human phrases approachable.
- **FR-002**: Current README, CLI, GUI, and API guidance MUST agree that recurring task input accepts human phrases or the supported cron subset.
- **FR-003**: At least one copy/pasteable recurring task-creation example in human syntax and one in cron syntax MUST express the same schedule.
- **FR-004**: CLI task help and documentation MUST describe dual-syntax `--schedule` input without adding a client-side parser or changing command behavior.
- **FR-005**: CLI cron documentation MUST distinguish `convert`, `explain`, `import`, and `export`, including daemon dependency, mutation, and input/output purpose.
- **FR-006**: GUI guidance MUST cover cron entry, automatic selection, live preview, retained edit prefill, syntax switching, and named validation refusal while keeping human examples first.
- **FR-007**: API documentation MUST define automatic syntax selection, the optional explicit syntax discriminator, no-fallback validation, normalized expression retention, and response source identity.
- **FR-008**: The cron guide MUST name the exact default field count/order and accurately document supported numbers, names, lists, ranges, steps, descriptors/macros, and extensions.
- **FR-009**: The cron guide MUST explain traditional DOM/DOW OR semantics, go-schedule's refusal of jointly restricted DOM/DOW input, and link future coverage to #22.
- **FR-010**: The cron guide MUST explain field-local step semantics and identify any forms refused because they cannot be represented faithfully.
- **FR-011**: Timezone and DST guidance MUST state that the task timezone governs compilation and execution through the single recurrence model.
- **FR-012**: Documentation MUST distinguish a five-field schedule expression from a crontab line or file and explain why import remains a separate workflow.
- **FR-013**: The authoritative master specification and contracts MUST adopt the dual-syntax posture without weakening the requirement that users can schedule tasks without knowing cron.
- **FR-014**: Completed boundary-only Spec 008 artifacts MUST retain their historical text and receive prominent supersession notes pointing to S019-S021.
- **FR-015**: Current package/source comments that falsely describe cron as interchange-only MUST be corrected; comments that accurately describe a human-only component MUST remain specific.
- **FR-016**: Offline documentation-policy tests MUST reject categorical current-surface claims that cron is prohibited or interchange-only and MUST require the core dual-syntax examples/contracts.
- **FR-017**: Existing documentation link, front-matter, fence, theme, CLI help, and behavior tests MUST remain green.
- **FR-018**: The slice MUST update the chronological Unreleased changelog with the completed documentation posture and the preserve-history/supersession decision.
- **FR-019**: Wildcard steps greater than one in day-of-month, month, or day-of-week MUST be refused by name rather than silently compiled as an unstepped daily schedule.
- **FR-020**: The pull request MUST use `Closes #52` and `Closes #50`; issue #22 MUST remain open.
- **FR-021**: Apart from FR-019's fidelity correction, the slice MUST NOT change syntax breadth, API shape, storage, migrations, engine timing, GUI behavior, IPC/security, daemon lifecycle, or command execution.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of current user-facing schedule-input surfaces state the same human-or-supported-cron contract.
- **SC-002**: The documentation contains at least one equivalent, copy/pasteable task-creation example for each syntax and explains their identical meaning.
- **SC-003**: Every dialect item named in issue #52 (field count/order, names/descriptors, timezone/DST, DOM/DOW OR, and field-local steps) is covered by an accurate, linked contract.
- **SC-004**: A repository-wide policy search finds zero unsuperseded current claims that cron is categorically prohibited or interchange-only.
- **SC-005**: All historical boundary-only artifacts selected for this slice carry an explicit supersession note while their original chronological content remains intact.
- **SC-006**: Documentation-policy fixtures prove both acceptance of the required dual-syntax posture and rejection of obsolete categorical language.
- **SC-007**: All eight repository verification gates pass, all core packages remain at or above 80% coverage, and every changed text file passes strict UTF-8/no-BOM/mojibake validation.
- **SC-008**: All tested calendar-field wildcard-step expressions produce a specific fidelity refusal and no task preview/create mutation.

## Assumptions

- S018-S020 behavior and their central input/API contracts are complete and are
  the source of truth; this slice documents them rather than changing them.
- The cron fidelity table remains the canonical feature-by-feature contract.
- Human-readable schedules remain the teaching default, but not the only valid
  authoring form.
- Historical changelog entries and completed Spec-Kit task descriptions remain
  chronological evidence, not active policy.
- Issue #22 remains the sole tracker for increasing cron dialect breadth.

## Out of Scope

- Any new cron construct, detector/fallback behavior, or timing semantics; the
  single FR-019 parser change only prevents an existing silent approximation.
- API, schema, engine, daemon, GUI interaction, security, packaging, release,
  or command-execution changes.
- Rewriting historical changelog entries or completed specifications as though
  dual-syntax behavior existed before S019-S021.
- Closing or implementing any portion of issue #22.

## Traceability

- Completes GitHub issue #52 and the remaining documentation outcome in #50.
- Documents string conversion from #51/S018, central authoring from S019, and
  GUI adoption from S020.
- Explicitly supersedes the boundary-only authoring policy from #12/S008 while
  preserving it as delivery history.
- Pull request uses `Closes #52` and `Closes #50`; #22 remains open.
