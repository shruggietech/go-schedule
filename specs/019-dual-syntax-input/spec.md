# Feature Specification: Dual-Syntax Task Input Foundation

**Feature Branch**: `codex/019-dual-syntax-input`

**Created**: 2026-08-28

**Status**: Implemented

**Delivery**: [PR #54](https://github.com/shruggietech/go-schedule/pull/54)

**Input**: User description: "S019: establish the first focused vertical slice
of issue #50 by accepting and retaining supported cron expressions through the
central schedule-input boundary, CLI task authoring, API preview/create/update,
stored task round-tripping, and cron import. Keep RRULE/anchor authoritative;
defer GUI adoption and the broad documentation rewrite in issue #52."

## Clarifications

### Session 2026-08-28

- Q: Must source syntax be stored in a new database field? → A: No. Retain the
  original expression in the existing durable field, derive syntax centrally,
  and expose the derived identity in API responses.
- Q: How can an API caller avoid automatic-classification ambiguity? → A:
  Preview, create, and update accept an optional `human` or `cron` hint; omitted
  hints use the shared deterministic detector.
- Q: Which parts of issue #50 form this slice's end-to-end boundary? → A: Central
  input, API, CLI task authoring, storage round-trip, and import retention are
  included; GUI adoption and issue #52 documentation follow later.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create and preview a task with cron (Priority: P1)

As a CLI or API user, I can preview and create a recurring task using a
supported cron expression anywhere the non-GUI task workflow accepts a schedule,
without translating it into English first.

**Why this priority**: This is the smallest end-to-end result that changes cron
from a conversion-only format into usable task-authoring input while retaining
the existing scheduling engine.

**Independent Test**: Submit `0 9 * * 1-5` to schedule preview and task creation,
then verify that both produce the same weekday 09:00 recurrence and that the
created task reports the original cron expression.

**Acceptance Scenarios**:

1. **Given** a supported five-field cron expression, **When** a user previews it,
   **Then** the response identifies cron as the source syntax and returns the
   same recurrence and upcoming runs that task creation will use.
2. **Given** a supported cron expression, **When** a user creates a task through
   the API or `gosched task add --schedule`, **Then** the task is created with
   RRULE/anchor as its execution definition and the submitted cron expression
   as its editable source.
3. **Given** a cron-shaped expression that is invalid or cannot be represented
   faithfully, **When** it is previewed or submitted, **Then** it is rejected
   with a specific cron/fidelity reason and is never retried as human text.
4. **Given** an API caller with an unambiguous syntax requirement, **When** the
   caller supplies an explicit `human` or `cron` syntax hint, **Then** validation
   uses only that syntax and reports a field-specific error for an invalid hint.

---

### User Story 2 - Edit and retrieve the syntax originally supplied (Priority: P2)

As a task maintainer, I can edit an existing task with either supported syntax
and later receive the same source form I supplied, while existing human phrases
continue to behave exactly as before.

**Why this priority**: First-class input is incomplete if reading or editing a
task silently replaces cron with generated prose or changes existing human task
semantics.

**Independent Test**: Create a human task and a cron task, replace each schedule
through the update API/CLI, fetch them again, and verify expression, source
syntax, RRULE, anchor, timezone behavior, and upcoming runs.

**Acceptance Scenarios**:

1. **Given** a cron-authored task, **When** it is fetched or updated without a
   schedule replacement, **Then** its cron expression and cron syntax identity
   are returned unchanged.
2. **Given** a human-authored task, **When** it follows the same create, preview,
   update, and read paths, **Then** its retained phrase and run times remain
   unchanged from the pre-feature behavior.
3. **Given** a recurring schedule replacement in either syntax, **When** the
   update succeeds, **Then** source expression and syntax change together while
   unrelated task fields and missing-date policy retain their prior behavior.
4. **Given** a one-off task, **When** it is returned by the API, **Then** it has
   no recurring source-syntax identity and its existing behavior is unchanged.

---

### User Story 3 - Keep cron source through import (Priority: P3)

As a user importing a crontab, I can later view or edit each supported imported
task using the cron expression from the source line instead of an English
replacement.

**Why this priority**: Import is the main onboarding path for cron users. Losing
their source immediately after import defeats the first-class-input goal.

**Independent Test**: Import a supported crontab line, fetch the created task,
and verify that its source is the timing expression from the line, its syntax is
cron, and its upcoming runs match import preview.

**Acceptance Scenarios**:

1. **Given** a supported crontab job line, **When** it is imported, **Then** task
   creation receives the original timing expression with an explicit cron hint.
2. **Given** import preview and creation for the same line, **When** upcoming
   runs are evaluated, **Then** both use the central input behavior and agree.
3. **Given** a declined or invalid cron line, **When** import runs, **Then** its
   existing per-line refusal behavior remains and no task is created.

### Edge Cases

- `@` shorthands and five fields with a cron-shaped minute field are classified
  as cron automatically; existing five-word human phrases remain human.
- A cron-shaped invalid expression is not reinterpreted as human input.
- An explicit syntax hint overrides automatic detection and accepts only
  `human` or `cron`; case and surrounding whitespace do not create extra values.
- Surrounding whitespace is trimmed once before a directly submitted expression
  is retained; internal spelling, capitalization, and spacing are otherwise
  preserved. Import retains the scanner's existing single-space-normalized
  timing expression rather than reconstructing prose.
- Supported cron field-local steps keep their current semantics. Constructs
  outside the existing faithful converter remain refused rather than broadened.
- Combined restricted day-of-month and day-of-week input keeps the current
  named refusal. Expanding cron coverage or introducing different semantics is
  issue #22, not this slice.
- Timezone, DST, missing-date policy, overlap, catch-up, and one-off behavior do
  not gain alternate execution paths.
- Existing schedule rows whose expression is empty remain readable and report
  no source syntax; no backfill or destructive migration is performed.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide one central recurring schedule-input
  boundary shared by task preview, create, update, CLI-backed requests, and cron
  import.
- **FR-002**: Automatic classification MUST treat trimmed `@` input and exactly
  five fields with a cron-shaped minute field as cron, and every other value as
  human input, consistent with the S018 converter.
- **FR-003**: API preview, create, and update requests MUST accept an optional
  source-syntax hint whose only non-empty values are `human` and `cron`; invalid
  values MUST produce a validation error naming the hint field.
- **FR-004**: An explicit hint MUST select only the named parser and MUST NOT
  fall back to the other syntax after failure.
- **FR-005**: Supported cron input MUST compile through the existing faithful
  cron-to-human conversion into the existing recurrence model; raw cron MUST
  NOT become an engine execution input.
- **FR-006**: A successfully parsed recurring schedule MUST retain the trimmed
  submitted expression for editing while RRULE/anchor remain the only
  authoritative timing definition.
- **FR-007**: Cron-shaped invalid or lossy input MUST be refused with the
  existing specific parser/fidelity reason and MUST NOT be retried as human.
- **FR-008**: Preview, create, and update responses MUST identify the recurring
  source syntax as `human` or `cron`; one-off and legacy expressionless
  schedules MUST report no source syntax.
- **FR-009**: Task reads and edits that do not replace a schedule MUST preserve
  the stored expression, execution timing, and inferred source identity.
- **FR-010**: Existing human schedule phrases MUST retain their expression,
  RRULE, anchor, summary, timezone/DST behavior, and validation behavior across
  every changed path.
- **FR-011**: `gosched task add` and `gosched task edit` MUST accept either
  supported syntax through `--schedule` and describe that contract in their
  command help without adding a second CLI-only parser.
- **FR-012**: Supported crontab imports MUST submit the original timing
  expression with cron identity, while preserving existing preview reporting,
  command parsing, partial-success, and refusal behavior.
- **FR-013**: Preview and persisted execution for the same input, timezone, and
  missing-date policy MUST produce identical upcoming runs across DST and month
  boundaries.
- **FR-014**: The feature MUST require no schedule-storage schema migration;
  source identity for existing rows MUST be derived from retained expression
  through the central classifier.
- **FR-015**: Relevant tests, comments, help, API contracts, and chronological
  project records that mandate human-only non-GUI input MUST be deliberately
  superseded; historical Spec-Kit artifacts remain historical records rather
  than being silently rewritten.
- **FR-016**: The feature MUST NOT change IPC access, authorization, secret
  handling, daemon lifecycle, task execution, or any security boundary.
- **FR-017**: Existing GUI requests MUST explicitly retain human-only parsing
  until the GUI adoption slice; this preservation wiring MUST NOT add GUI cron
  validation, preview, prefill, help, or save behavior.

### Key Entities

- **Schedule Input**: A submitted recurring expression, optional explicit
  syntax hint, detected source syntax, and either a compiled schedule or a named
  refusal.
- **Stored Schedule**: The existing authoritative RRULE/anchor recurrence plus
  the inert editable source expression. Its storage shape is unchanged.
- **Task Request/Response**: Preview, create, and update contracts that accept an
  optional hint and return source syntax alongside the existing schedule data.
- **Imported Job**: A crontab job whose timing expression becomes the retained
  task source while command and arguments retain their existing import meaning.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All supported cron examples in this slice complete preview,
  creation, retrieval, and update with one retained cron expression and no
  manual conversion step.
- **SC-002**: Human and cron inputs representing the same schedule produce
  identical upcoming runs in every tested DST/month-boundary window.
- **SC-003**: 100% of cron-shaped invalid or lossy fixtures fail with no task
  mutation, no human fallback, and a named reason.
- **SC-004**: 100% of pre-existing human and expressionless compatibility
  fixtures retain their prior timing and response behavior except for the new
  non-empty syntax identity on recurring rows with retained expressions.
- **SC-005**: A supported imported job can be fetched after creation with the
  scanner's timing expression and `cron` source identity.
- **SC-006**: The complete repository verification suite passes with all eight
  gates green and no core package below its coverage threshold.

## Assumptions

- S018's supported dialect, structural detection, and named fidelity boundaries
  are the source of truth for this slice.
- `Schedule.Expression` is already sufficient durable source storage, so a
  syntax column and database migration would add state that can drift without
  improving execution correctness.
- The API exposes explicit source identity so later GUI work does not need to
  duplicate detection, even though identity can be reconstructed from the
  retained expression.
- The CLI remains a thin client and relies on the daemon's central input
  boundary for authoritative validation.
- GUI task-editor adoption, broad README/docs posture changes in issue #52,
  additional cron constructs in issue #22, and closing epic #50 are out of
  scope. This slice will reference #50 rather than close it.

## Scope Boundaries

### In Scope

- Central dual-syntax recurring input parsing and source classification.
- Optional API syntax hint and response syntax identity.
- API preview/create/update/read behavior.
- CLI task add/edit help and end-to-end request behavior.
- Cron import source retention and preview/create parity.
- Human-only GUI request hints that prevent accidental scope expansion.
- Compatibility, fidelity, timezone/DST, and persistence regressions.

### Out of Scope

- GUI validation, preview, prefill, or save behavior.
- The global dual-syntax documentation rewrite tracked by #52.
- New cron dialect features or changed DOM/DOW semantics tracked by #22.
- Raw-cron execution, a second scheduler, schema migration, or source backfill.
- Release, packaging, security-governance, or unrelated issue work.

## Traceability

- Partial delivery of GitHub issue #50; the eventual pull request uses
  `Refs #50` and leaves it open.
- Builds on closed issue #51 and Spec 018's pure conversion boundary.
- Supersedes human-only authoring requirements from Spec 008 only for the
  non-GUI surfaces explicitly listed here.
