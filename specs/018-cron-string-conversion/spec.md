# Feature Specification: Pure Schedule String Conversion

**Feature Branch**: `codex/018-cron-string-conversion`

**Created**: 2026-08-28

**Status**: Draft

**Input**: User description: "S018: add a pure local cron-to-human and
human-to-cron string converter to the CLI, closing issue #51 and laying focused
groundwork for the broader first-class cron work in #50."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Translate cron without a daemon (Priority: P1)

As a shell user, I can provide one supported cron expression and receive only
the equivalent human-readable schedule so I can use the result in a script,
document, or go-schedule task command.

**Why this priority**: Cron-to-human conversion builds directly on the shipped
interoperability surface while removing the current daemon and report-format
requirements for the common one-string translation job.

**Independent Test**: Stop or make the daemon unavailable, convert
`0 9 * * 1-5`, and observe exactly `weekdays at 09:00` plus a trailing newline
on standard output with a successful exit.

**Acceptance Scenarios**:

1. **Given** the daemon is unavailable, **When** a user converts
   `0 9 * * 1-5`, **Then** the command succeeds locally and prints only
   `weekdays at 09:00` followed by a newline.
2. **Given** a supported shorthand or five-field expression, **When** the user
   explicitly requests human output, **Then** it is validated as cron and
   translated according to the existing fidelity rules.
3. **Given** malformed or unsupported cron, **When** automatic detection or an
   explicit human-output request selects cron input, **Then** the command names
   the problem and never retries the text as a human schedule.

---

### User Story 2 - Produce one faithful cron expression (Priority: P1)

As a shell user, I can provide one human-readable recurring schedule and receive
one canonical cron expression when cron can reproduce it exactly.

**Why this priority**: Symmetric conversion is the missing half of the existing
cron tooling and gives users a mutation-free way to translate a schedule before
the broader task-authoring work in #50.

**Independent Test**: With no daemon, convert `weekdays at 09:00` and observe
exactly `0 9 * * 1-5` plus a trailing newline; then convert the result back and
confirm equivalent runs across a month boundary and daylight-saving change.

**Acceptance Scenarios**:

1. **Given** a supported recurring human schedule, **When** the user converts
   it, **Then** the command prints one documented canonical cron expression.
2. **Given** a one-off, elapsed interval, non-default missing-date behavior, or
   another schedule cron cannot reproduce, **When** conversion is attempted,
   **Then** the command refuses by name rather than approximating run times.
3. **Given** a human phrase whose shape could be mistaken for cron, **When** the
   user supplies the explicit cron-output override, **Then** the phrase is
   validated as human input and converted or faithfully refused.

---

### User Story 3 - Automate conversion reliably (Priority: P2)

As a script author, I can select a destination syntax and request structured
output with stable fields, streams, and exit behavior.

**Why this priority**: The converter is intentionally pipeline-oriented.
Predictable machine behavior prevents scripts from scraping the richer
`explain` report or accidentally treating a refusal as a successful string.

**Independent Test**: Exercise success, malformed input, and faithful-refusal
cases in text and structured modes; assert exact output streams, fields, and
exit codes without any daemon or stored task.

**Acceptance Scenarios**:

1. **Given** a successful conversion in structured mode, **When** the command
   completes, **Then** standard output contains one object naming the input
   syntax, output syntax, normalized input, output, and an empty refusal reason.
2. **Given** malformed or unfaithful input in structured mode, **When** the
   command refuses, **Then** standard output remains empty, standard error
   contains one object with the same stable fields and a specific refusal
   reason, and the exit code is 2.
3. **Given** normal shell quoting on Windows PowerShell or a POSIX shell,
   **When** one quoted schedule string is supplied, **Then** internal whitespace
   is preserved as one argument and surrounding whitespace does not affect the
   conversion.

### Edge Cases

- An empty or whitespace-only argument is invalid even when supplied as the one
  positional value.
- Any `@`-prefixed input and any five-field input whose first field has cron's
  numeric/wildcard minute shape is treated as cron during automatic detection;
  invalid cron cannot fall through into the human parser.
- Existing five-field human phrases such as `every 15 minutes from 9am` and
  `3rd wednesday monthly at 14:00` remain human input automatically; the
  destination override remains an explicit validation escape hatch.
- Cron day-of-month plus day-of-week restrictions retain traditional OR
  semantics and are refused when no faithful human representation exists.
- Field-local steps, six-field dialects, and unsupported extensions retain
  their existing named fidelity outcomes.
- A human schedule whose recurrence depends on its creation instant or cannot
  be represented in five-field cron is refused rather than anchored to the
  conversion time.
- Equivalent spellings, extra surrounding whitespace, and case differences
  produce the same canonical output.
- Default text output never includes labels, arrows, previews, timestamps, or
  daemon status.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The CLI MUST provide a local string-conversion operation beneath
  the existing cron command group.
- **FR-002**: The operation MUST accept exactly one non-empty schedule string
  and MUST create, edit, or delete no task, schedule, configuration, or file.
- **FR-003**: With no override, the operation MUST deterministically classify an
  `@`-prefixed value or five whitespace-delimited fields beginning with a
  cron-shaped minute field as cron input; every other value MUST be classified
  as human input.
- **FR-004**: Once input is classified as cron, malformed or unsupported cron
  MUST NOT be retried as human text.
- **FR-005**: Users MUST be able to select `cron` or `human` as the destination
  syntax; this selection MUST force validation of the opposite input syntax.
- **FR-006**: Supported cron input MUST translate through the same dialect and
  fidelity rules used by the existing cron explanation and import behavior.
- **FR-007**: Supported human input MUST produce one canonical five-field cron
  expression, and any conversion that changes run times MUST be refused by
  name.
- **FR-008**: Traditional day-of-month/day-of-week OR semantics, field-local
  step behavior, supported names and shorthands, and all existing named
  unsupported constructs MUST retain their documented meaning.
- **FR-009**: Successful default output MUST be exactly the converted string and
  one trailing newline on standard output, with no other output and exit code 0.
- **FR-010**: Malformed and unfaithful conversions MUST write no standard
  output, MUST provide a specific diagnostic on standard error, and MUST exit
  with code 2.
- **FR-011**: Structured mode MUST expose stable `input_syntax`,
  `output_syntax`, `input`, `output`, and `refusal_reason` fields. Success MUST
  write the object to standard output; failure MUST write it to standard error
  while preserving FR-010's exit behavior.
- **FR-012**: Surrounding whitespace and case variation MUST not change meaning,
  while the normal quoted-argument behavior of Windows PowerShell and POSIX
  shells MUST be documented and protected by command-boundary tests.
- **FR-013**: Conversion MUST succeed with the daemon stopped or unavailable and
  MUST perform no network, IPC, storage, or task mutation.
- **FR-014**: CLI help and the CLI guide MUST include copyable examples in both
  directions and clearly distinguish conversion from explanation, import, and
  export.
- **FR-015**: Automated tests MUST cover both directions, automatic and forced
  syntax selection, exact streams and exit behavior, structured results, named
  refusals, and faithful round trips across a month boundary and a
  daylight-saving transition.
- **FR-016**: This slice MUST NOT make cron valid task-authoring input, change
  APIs or the GUI, broaden the supported cron dialect, add persisted
  source-syntax metadata, or revise the product-wide dual-syntax posture; those
  outcomes remain tracked by #50, #52, and #22.

### Key Entities

- **Conversion request**: One input string plus an optional destination syntax
  and output-format selection.
- **Syntax identity**: The deterministic classification `cron` or `human` used
  to choose validation and the opposite output syntax.
- **Conversion result**: The normalized input, both syntax identities, converted
  string on success, and empty refusal reason.
- **Conversion refusal**: The same traceable identity fields plus a specific
  reason explaining why the input is malformed or cannot be represented
  faithfully; it never carries an approximated output.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can translate one supported string in either direction
  with one command and receive exactly one converted line without starting the
  daemon.
- **SC-002**: One hundred percent of invalid and unfaithful conversion cases in
  the documented fidelity table produce no default standard output, a named
  diagnostic, and exit code 2.
- **SC-003**: Automatic detection and explicit destination selection agree on
  the output for every supported example, while every cron-shaped invalid input
  is rejected without human-parser fallback.
- **SC-004**: Structured success and refusal tests account for all five stable
  result fields and their required output stream.
- **SC-005**: The documented round-trip corpus produces equivalent upcoming
  runs across at least one month boundary and one daylight-saving transition,
  with no approximated case accepted.
- **SC-006**: All eight repository verification gates pass without reducing
  core-package coverage or changing existing task, API, GUI, storage, or
  scheduling-engine behavior.

## Assumptions

- The existing supported five-field cron dialect and named fidelity boundaries
  remain authoritative; expanding them belongs to #22.
- Destination selection names the desired output (`--to cron` means human input;
  `--to human` means cron input), matching common conversion-tool conventions.
- Structured failures are emitted as one machine-readable diagnostic on
  standard error so both the no-standard-output failure contract and structured
  automation contract remain true.
- Human-readable schedules remain the normal task-authoring syntax in this
  slice. The converter is a pure boundary utility and does not preempt #50's
  source-retention and cross-surface design.
- No new dependency, daemon endpoint, configuration, migration, or branding
  asset is needed.
