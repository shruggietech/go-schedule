# Feature Specification: GUI Dual-Syntax Scheduling

**Feature Branch**: `codex/020-gui-dual-syntax`

**Created**: 2026-08-28

**Status**: Draft

**Input**: User description: "Complete issue #50's remaining functional gap by allowing supported cron expressions in the existing GUI task editor, including validation, preview, edit prefill, and save/update behavior."

## Context

S019 made supported human and cron schedules first-class inputs through the
central task boundary, API, CLI, and import paths. The desktop editor was kept
explicitly human-only to prevent accidental behavior expansion. S020 adopts
that delivered boundary in the existing Schedule field so a user who starts
with cron can remain a cron user while creating or editing a task.

This is the remaining functional slice of issue #50. The broad dual-syntax
documentation rewrite remains issue #52 and additional cron dialect breadth
remains issue #22.

## Clarifications

### Session 2026-08-28

- Q: Should the editor add a syntax selector or retain one Schedule field with deterministic detection? → A: Retain one field and use the shared deterministic selection rule, avoiding duplicate UI state.
- Q: How should an imported or API-created cron task open in the editor? → A: Prefill the retained cron expression exactly and preview/save it as cron without translating it to prose.
- Q: Does this slice complete the documentation follow-through or close issue #50? → A: Update only GUI-local guidance and use `Refs #50`; issue #52 remains the final documentation follow-through before the epic closes.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create a Cron Task in the Desktop Editor (Priority: P1)

An operator can enter a supported cron expression in the existing Schedule
field, see the same explanation and upcoming runs produced by the central
preview boundary, and save the recurring task without converting it manually.

**Why this priority**: Creation is the first point where the GUI currently
rejects an otherwise valid first-class schedule input.

**Independent Test**: Enter `0 9 * * 1-5` for a new recurring task, confirm the
preview identifies the weekday schedule, save it, and inspect the create request
to prove it carries the original expression and cron source identity.

**Acceptance Scenarios**:

1. **Given** a new recurring task form, **When** the operator enters a supported cron expression, **Then** the editor accepts it and shows the central human summary and upcoming runs.
2. **Given** a valid cron expression and the other required task fields, **When** the operator saves, **Then** the create request retains the exact normalized expression and identifies it as cron.
3. **Given** a valid human schedule phrase, **When** the operator previews and saves, **Then** the existing human behavior and identity remain unchanged.
4. **Given** invalid or unsupported cron-shaped input, **When** validation runs, **Then** Save remains disabled and the editor displays the specific parser or fidelity reason without retrying the input as human text.

---

### User Story 2 - Edit a Task in Its Original Syntax (Priority: P2)

An operator opening a recurring task sees the retained human or cron expression
that created it and can save it unchanged or replace it with either supported
syntax.

**Why this priority**: Round-trip editing is the capability that lets imported
and API-created cron users remain cron users in the desktop application.

**Independent Test**: Open a task whose schedule response contains expression
`0 9 * * 1-5` and source syntax `cron`, prove the field and preview retain that
expression, save it, then replace it with a human phrase and prove the update
request changes source identity accordingly.

**Acceptance Scenarios**:

1. **Given** a recurring cron task, **When** its editor opens, **Then** the Schedule field contains the retained cron expression rather than a generated human replacement.
2. **Given** an unchanged cron task, **When** the operator saves, **Then** the update request retains cron input and source identity.
3. **Given** a cron task, **When** the operator replaces the Schedule field with a valid human phrase, **Then** preview and update select human syntax.
4. **Given** a human task, **When** the operator replaces the Schedule field with supported cron, **Then** preview and update select cron syntax.
5. **Given** an expressionless legacy recurrence, **When** the editor opens, **Then** the existing degraded blank-field guidance and schedule-preservation behavior remain unchanged.

---

### User Story 3 - Understand the Field's Dual-Syntax Contract (Priority: P3)

An operator can discover from the editor and GUI field guide that the Schedule
field accepts approachable human phrases and the supported five-field cron
subset, while unsupported cron remains an explicit refusal.

**Why this priority**: A capability hidden behind a previously human-only field
is incomplete, but the full cross-product documentation rewrite belongs to
issue #52.

**Independent Test**: Inspect the editor help and GUI field guide and verify
they show copy/pasteable human and cron examples, name the five-field boundary,
and direct broader fidelity questions to the cron guide.

**Acceptance Scenarios**:

1. **Given** the task editor, **When** the operator opens Schedule help, **Then** both human and cron examples are present with human language remaining the approachable default.
2. **Given** the GUI field guide, **When** a reader checks Schedule, **Then** it no longer claims cron is prohibited and links to the fidelity contract.

### Edge Cases

- Leading and trailing whitespace is normalized consistently with S019 while
  the meaningful source expression is preserved.
- A cron-shaped value that is invalid or outside the faithful dialect is never
  retried as a human phrase.
- Five-word human phrases such as `3rd wednesday monthly at 14:00` remain human.
- Switching between human and cron text in one editing session updates preview
  and save identity from the current text, not from stale prefill metadata.
- An empty recurring field disables Save for new tasks; on an existing degraded
  legacy task it continues to mean "keep the current schedule".
- One-off task mode sends no recurring expression or syntax identity and keeps
  its existing date/time validation.
- Preview transport failures remain visible without changing the locally
  determined validity of a supported expression.
- Timezone, DST, month-boundary, and missing-date behavior use the central
  preview and execution model and gain no GUI-specific evaluator.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The existing GUI Schedule field MUST accept every human phrase and supported cron expression accepted by the S019 central input boundary.
- **FR-002**: The GUI MUST use the shared deterministic syntax selection and no-fallback rules; it MUST NOT maintain an independent cron detector or timing evaluator.
- **FR-003**: A valid recurring input MUST enable save eligibility when all other required fields are valid, regardless of whether its source syntax is human or cron.
- **FR-004**: Invalid or unsupported cron-shaped input MUST disable Save and show the specific central parser or fidelity reason without human fallback.
- **FR-005**: Recurring preview requests MUST send the current normalized Schedule expression with the selected `human` or `cron` source identity.
- **FR-006**: Create and schedule-replacing update requests MUST send the current normalized expression with the same source identity used for preview.
- **FR-007**: Editing a recurring task MUST prefill the retained `Schedule.Expression`; cron input MUST NOT be replaced by `HumanSummary`.
- **FR-008**: The syntax used for preview and save MUST follow the current field contents after an edit rather than stale response metadata.
- **FR-009**: Existing human schedule creation, preview, editing, help, timezone, start-at, missing-date, overlap, catch-up, and error behavior MUST remain compatible except where guidance now includes cron.
- **FR-010**: Existing one-off creation/editing and expressionless legacy recurrence preservation MUST remain unchanged and MUST send no recurring syntax hint.
- **FR-011**: The GUI MUST display actionable validation or preview errors without exposing internal conversion steps or creating a second schedule summary.
- **FR-012**: Editor help and `docs/gui-fields.md` MUST describe both accepted syntax forms, include copy/pasteable examples, and link dialect/fidelity detail to `docs/cron.md`.
- **FR-013**: The slice MUST add deterministic tests for automatic classification at the GUI boundary, preview/create/update identity parity, syntax switching, retained cron prefill, invalid/refused cron, five-word human regression, one-off isolation, and legacy preservation.
- **FR-014**: The slice MUST NOT change the API, storage schema, scheduling engine, cron dialect, IPC/security boundary, daemon lifecycle, or command execution path.
- **FR-015**: The chronological changelog MUST record GUI dual-syntax adoption, the single-field decision, and that broad documentation remains in issue #52.
- **FR-016**: The pull request MUST use `Refs #50`; issue #50 remains open until issue #52 completes the product-wide documentation contract.

### Key Entities

- **Editor Schedule Input**: The current recurring field text, its normalized expression, selected source syntax, and either a compiled previewable schedule or a named validation refusal.
- **Task Detail Schedule**: The retained expression and response source identity used to initialize an edit without translating the operator's syntax.
- **Editor Submission**: A create or update request whose recurring expression and source identity match the most recent successful local validation and preview.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of the supported cron fixtures selected for this slice can be entered, previewed, and submitted from a new GUI task without manual translation.
- **SC-002**: Human and cron inputs representing the same recurrence produce identical previewed upcoming runs in every tested timezone, DST, and month-boundary scenario.
- **SC-003**: 100% of invalid or faithfully unsupported cron fixtures keep Save disabled, show a named reason, and produce no create or update request.
- **SC-004**: A retained cron task can be opened and saved with the identical normalized expression and cron identity, then switched to human input in the same editor session.
- **SC-005**: 100% of existing GUI human, one-off, policy, start-at, and expressionless-legacy regression fixtures retain their prior behavior.
- **SC-006**: The complete repository verification suite passes all eight gates with no core package below its coverage threshold.

## Assumptions

- S019's central schedule-input boundary and API source identity are the source
  of truth for selection, validation, and retained identity.
- The current single Schedule field is preferable to a syntax toggle because
  input classification is deterministic and duplicate state could disagree.
- Preview remains a daemon request after local validation; no GUI-only run-time
  calculation is introduced.
- Existing editor help structure and inline preview area can carry the added
  guidance and named refusals without a layout redesign.
- Issue #52 owns the complete README, CLI, API, and cron-dialect documentation
  rewrite; S020 changes GUI-local help and field documentation necessary to
  make the new behavior discoverable, plus the single `docs/cron.md` sentence
  that would otherwise falsely claim the GUI still rejects cron.

## Out of Scope

- New cron syntax, parser extensions, DOM/DOW policy changes, or issue #22.
- A syntax selector, cron builder, separate cron field, or visual editor
  redesign.
- API, database, recurrence-engine, IPC, authorization, service, packaging, or
  execution changes.
- The product-wide documentation rewrite and closure of issues #50/#52.

## Traceability

- Functional continuation of GitHub issue #50 and S019.
- Pull request uses `Refs #50` and leaves it open for issue #52.
- Fidelity breadth continues under issue #22.
