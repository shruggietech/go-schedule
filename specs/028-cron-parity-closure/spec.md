# Feature Specification: Cron Parity Closure

**Feature Branch**: `codex/028-cron-parity-closure`

**Created**: 2026-08-29

**Status**: Draft

**Input**: Close issue #22 with one substantial fidelity slice: implement the remaining cron expression and crontab-file features that map cleanly to go-schedule, explicitly ratify incompatible features, and bring documentation and marketing into agreement with the resulting contract.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Import an operational crontab faithfully (Priority: P1)

As an operator, I can import a user or system crontab without losing its schedule timezone, environment, shell behavior, run-as user, or standard-input payload.

**Why this priority**: The current importer can create tasks that look valid but execute differently from cron, which is more dangerous than a visible refusal.

**Independent Test**: Dry-run and import fixtures containing `CRON_TZ`, quoted and unquoted environment assignments, `SHELL`, shell operators, escaped and unescaped `%`, and system-crontab user fields, then inspect the resulting task requests and execute a task with standard input.

**Acceptance Scenarios**:

1. **Given** `CRON_TZ=America/New_York`, **When** later job lines are imported, **Then** their schedules use that timezone until another `CRON_TZ` assignment changes it.
2. **Given** ordinary environment assignments, including `TZ`, `PATH`, `HOME`, and `SHELL`, **When** later jobs are imported, **Then** each task receives the effective environment without retroactively changing earlier jobs.
3. **Given** a cron command containing pipelines, redirections, substitutions, quoting, or `&&`, **When** it is imported, **Then** the task invokes the effective cron shell with the original command text rather than whitespace-splitting it.
4. **Given** unescaped `%` characters, **When** the line is imported, **Then** the first separates the command from standard input, later unescaped percent signs become newlines, and escaped percent signs remain literal.
5. **Given** system-crontab mode and a username field, **When** a line is imported, **Then** the username becomes the task's run-as identity and is never mistaken for the command.
6. **Given** user-crontab mode, **When** an otherwise identical line is imported, **Then** no username column is consumed and the command remains intact.
7. **Given** a crontab file using six-field Quartz timing, **When** it is imported with the explicit Quartz dialect option, **Then** exactly six timing fields are consumed; the default Unix dialect continues to consume exactly five.

---

### User Story 2 - Use seconds-precision cron where it maps cleanly (Priority: P1)

As an operator, I can explain, preview, create, import, and export supported six-field Quartz-style schedules with seconds and `?` without switching to a different scheduler.

**Why this priority**: The scheduler already supports sub-minute recurrence, so refusing the corresponding cron input is an artificial interoperability gap.

**Independent Test**: Round-trip representative six-field expressions through parsing, compilation, persistence, preview, task input, import, and export, then compare exact run instants.

**Acceptance Scenarios**:

1. **Given** `*/30 * * * * *`, **When** it is compiled, **Then** it fires at seconds 0 and 30 of each minute and retains the original expression for editing.
2. **Given** `0 0 12 ? * MON`, **When** it is compiled, **Then** `?` acts as no specific day-of-month value and the task fires at noon on Mondays.
3. **Given** `?` outside the day-of-month or day-of-week field, **When** it is parsed, **Then** the input fails with a field-specific error and creates no task.
4. **Given** a recurrence with non-zero or multiple seconds that cron can carry, **When** it is exported, **Then** export emits a six-field expression; minute-resolution schedules continue to export as conventional five-field cron.
5. **Given** a six-field expression using unsupported Quartz year or incompatible day semantics, **When** it is submitted, **Then** it is refused by name rather than approximated.
6. **Given** a numeric day-of-week, **When** it is parsed, **Then** five-field input uses Unix cron numbering and six-field input uses Quartz numbering, so neither dialect shifts the requested weekday.

---

### User Story 3 - Know exactly where parity ends (Priority: P1)

As an operator, I can consult one complete fidelity matrix that states whether every audited cron feature is supported, deferred to a named capability, or deliberately out of scope.

**Why this priority**: Issue #22 cannot close honestly while marketing implies universal parity or audit rows remain undecided.

**Independent Test**: Compare every A1-A12 and B1-B9 row from issue #22 with the published matrix and verify each row has one decision, a rationale, and behavior matching tests.

**Acceptance Scenarios**:

1. **Given** the issue #22 audit, **When** the documentation is reviewed, **Then** every row A1-A12 and B1-B9 has an explicit supported, deferred, or out-of-scope decision.
2. **Given** `@reboot` or notification directives, **When** users consult the matrix, **Then** it points to the existing event or notification work rather than claiming schedule conversion support.
3. **Given** DOM-and-DOW OR semantics, unsupported modifier composites, run-parts directories, or anacron files, **When** they are encountered, **Then** the system refuses or excludes them explicitly and the matrix explains why.
4. **Given** README and command help, **When** cron support is described, **Then** claims are scoped to the documented faithful subset and never promise that every cron feature is supported.

---

### User Story 4 - Upgrade without breaking existing tasks (Priority: P2)

As an existing operator, I can upgrade and continue running current tasks while imported standard input and seconds schedules survive restart and use the same execution path as every other task.

**Why this priority**: Fidelity additions must not create a second evaluator or silently alter existing five-field schedules.

**Independent Test**: Migrate a pre-S028 database, compare existing schedules before and after, restart with new standard-input tasks, and run the canonical suite and scheduling benchmarks.

**Acceptance Scenarios**:

1. **Given** a pre-S028 database, **When** it is opened, **Then** existing tasks receive empty standard input and otherwise retain their stored behavior.
2. **Given** a task with imported standard input, **When** the daemon restarts and runs it, **Then** the exact stored payload is supplied to the child process.
3. **Given** any supported five- or six-field cron input, **When** it executes, **Then** it uses the existing durable recurrence evaluator rather than a cron-specific runtime evaluator.
4. **Given** existing five-field inputs and exports, **When** the upgrade is applied, **Then** their accepted syntax and minute-resolution output remain compatible.

### Edge Cases

- Environment assignments may contain optional whitespace around `=`, matching quotes, empty values, and embedded non-leading spaces; malformed assignments remain visible errors or ordinary command lines according to crontab grammar.
- Reassignments affect only following jobs. `CRON_TZ` is schedule context, while `TZ` remains child-process environment.
- `MAILTO` and `MAILFROM` are not ordinary child environment for parity purposes because their cron behavior is output delivery; they remain visible deferred warnings tied to notification work.
- An empty `SHELL` value, a shell path with unusable arguments, or an absent shell must fail visibly at validation or execution without falling back to whitespace splitting.
- Backslash handling around `%` must preserve shell-visible escaping for characters other than `%` and remove only the escape that makes `%` literal.
- A `%` inside quotes is still special to cron unless escaped, because cron performs the split before the shell.
- System-crontab lines missing either a user or command are errors. User mode never guesses that a token is a username.
- File import never guesses a timing dialect from command tokens. Unix five-field mode is the compatibility default and Quartz six-field mode is explicit.
- `?` is accepted only in six-field Quartz-style day-of-month or day-of-week positions and is normalized to an unrestricted selector without weakening existing DOM-and-DOW safety.
- Numeric weekdays follow their source dialect: Unix five-field accepts Sunday as 0 or 7, while Quartz six-field uses 1 through 7 for Sunday through Saturday; named weekdays map identically.
- Seconds value 60 and unsupported seven-field Quartz expressions are rejected.
- Windows may preview system-crontab run-as jobs but creation continues to reject unsupported run-as execution explicitly.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Crontab scanning MUST maintain ordered per-file context so assignments affect subsequent job lines and never preceding lines.
- **FR-002**: `CRON_TZ` MUST set the schedule timezone for subsequent jobs, MUST be validated as a resolvable timezone before creation, and MUST NOT be injected into the child environment solely because it controls scheduling.
- **FR-003**: Ordinary crontab environment assignments, including `TZ`, `PATH`, `HOME`, and `SHELL`, MUST be copied into each subsequent task's environment with quoted boundary whitespace handled according to crontab rules; `LOGNAME` overrides MUST be ignored visibly because cron does not permit them.
- **FR-004**: Reassigning or emptying a variable MUST update the effective environment for subsequent jobs; each imported task MUST own an independent snapshot.
- **FR-005**: `MAILTO` and `MAILFROM` MUST remain visible deferred warnings and MUST NOT be represented as working notification delivery.
- **FR-006**: Imported commands MUST execute through the effective `SHELL`, defaulting to `/bin/sh`, with the original command portion passed as one shell command string.
- **FR-007**: Import MUST implement cron percent semantics exactly: unescaped `%` becomes newline, the first separates executable command from standard input, and escaped `%` remains literal.
- **FR-008**: Tasks MUST support an optional persisted standard-input payload that is exposed through the local API and supplied verbatim to the child process on every run.
- **FR-009**: Existing tasks and callers that omit standard input MUST behave exactly as before and migration MUST require no operator action.
- **FR-010**: `gosched cron import` MUST offer an explicit system-crontab mode that consumes and validates the username column and maps it to the existing run-as task field.
- **FR-011**: User-crontab mode MUST remain the default and MUST never infer a username column.
- **FR-011a**: User-crontab import MUST accept an explicit run-as owner because the file itself contains no owner identity; it MUST NOT be combined with system layout.
- **FR-011b**: Unix run-as execution MUST establish the selected account's `LOGNAME`, `USER`, and default `HOME`, while preserving an explicit task `HOME` override.
- **FR-012**: Crontab import MUST use an explicit `unix` or `quartz` dialect option, default to Unix compatibility, and consume exactly five or six timing fields respectively without guessing from command text.
- **FR-013**: Import dry-run output MUST show each job's effective timezone, shell command, run-as identity when present, environment context, and standard-input presence without exposing a multiline payload ambiguously.
- **FR-014**: The cron parser MUST accept both conventional five-field expressions and supported six-field Quartz-style expressions whose first field is seconds.
- **FR-015**: Five-field input MUST imply second zero; six-field input MUST support seconds 0 through 59 with the existing list, range, and step forms.
- **FR-016**: In six-field input, `?` MUST be accepted only as the complete day-of-month or day-of-week field and MUST mean no specific value; it MUST be rejected elsewhere or when embedded in another token.
- **FR-017**: Five-field numeric weekdays MUST retain Unix cron numbering, while six-field numeric weekdays MUST use Quartz numbering; named weekdays MUST remain equivalent in both dialects.
- **FR-018**: Six-field expressions MUST compile directly into the existing durable recurrence model and MUST NOT add a cron runtime evaluator or new scheduling dependency.
- **FR-019**: Explain, convert, task input, preview, import, persistence, edit, execution, and export MUST agree on accepted six-field semantics and exact run instants.
- **FR-020**: Export MUST retain five-field output when seconds are exactly zero and MUST emit six-field Quartz numbering when a recurrence requires expressible sub-minute seconds.
- **FR-021**: Cron export MUST refuse task execution context it cannot serialize faithfully, including standard input, environment, or run-as identity, rather than emitting a lossy line.
- **FR-022**: Unsupported seven-field Quartz year expressions, DOM-and-DOW OR semantics, and unsupported modifier composites MUST be refused by name without approximation.
- **FR-023**: Documentation MUST contain a complete decision matrix for issue #22 rows A1-A12 and B1-B9, with a rationale and related issue where applicable.
- **FR-024**: `@reboot` MUST remain deferred to event-trigger work; `MAILTO` and `MAILFROM` delivery MUST remain deferred to notification work.
- **FR-025**: Cron DOM-and-DOW OR semantics and unsupported `L`, `W`, and `#` composites MUST be ratified outside the current single-recurrence representation until a faithful composite model exists.
- **FR-026**: Run-parts directory discovery and anacron-file parsing MUST be ratified as separate import formats, outside the crontab-file converter.
- **FR-027**: README, CLI help, and cron documentation MUST describe faithful interoperability with a documented subset and MUST remove universal-parity language.
- **FR-028**: Every implemented audit row MUST have focused parser or importer tests plus integration coverage at the API, persistence, or executor boundary it crosses.
- **FR-029**: The implementation MUST add no external service, permission, or third-party dependency.
- **FR-030**: Representative five- and six-field occurrence evaluation MUST remain within the existing p99 dispatch budget and MUST NOT regress the relevant benchmark by more than ten percent without recorded justification.

### Key Entities

- **Crontab context**: The effective schedule timezone, environment, and shell accumulated while scanning ordered assignment and job lines.
- **Imported job**: A parsed crontab line carrying timing, shell command text, standard input, environment snapshot, timezone, and optional run-as identity.
- **Standard-input payload**: Optional task data supplied verbatim to the child process independently of command arguments and environment.
- **Cron dialect**: The structural distinction between conventional five-field minute-resolution input and supported six-field Quartz-style seconds input.
- **Parity decision**: The supported, deferred, or out-of-scope disposition and rationale for one issue #22 audit row.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A mixed crontab fixture containing two timezone contexts, environment reassignment, shell operators, percent input, and two system users imports with 100 percent of supported semantics preserved in task requests.
- **SC-002**: Representative seconds expressions produce the exact expected instants through parser, API preview, persisted task detail, restart, and scheduler enumeration with zero duplicates or drift.
- **SC-003**: All pre-existing five-field cron regression fixtures retain their prior accepted behavior and canonical five-field export.
- **SC-004**: A pre-S028 database migrates with 100 percent of existing tasks assigned empty standard input and no other task-field changes.
- **SC-005**: Every invalid system line, malformed environment assignment, invalid timezone, misplaced `?`, seven-field expression, and incompatible day combination fails visibly and creates zero affected tasks.
- **SC-006**: The published matrix contains exactly one explicit disposition for all 21 audit rows, and README and CLI descriptions make zero universal cron-parity claims.
- **SC-007**: Representative cron evaluation benchmarks remain within ten percent of baseline and within the existing p99 dispatch budget.
- **SC-008**: All eight canonical verification gates, whitespace checks, UTF-8-without-BOM checks, and mojibake audits pass.

## Clarifications

### Session 2026-08-29

- `CRON_TZ` controls subsequent schedule interpretation. `TZ` is preserved as child-process environment and does not silently replace the schedule timezone.
- Imported commands use the effective cron shell. This intentionally replaces the old whitespace-splitting behavior because the old behavior changed pipelines, redirects, quoting, and conditionals.
- System-crontab parsing is opt-in. The importer does not guess whether a sixth token is a username.
- Quartz crontab-file parsing is also opt-in because a Unix command can begin with a cron-shaped token; structural guessing would corrupt valid commands.
- Seconds precision is the supported six-field boundary. Quartz's optional year field remains out of scope.
- Six-field numeric weekdays follow Quartz's Sunday-through-Saturday values 1 through 7; five-field inputs retain Unix cron's 0/7 Sunday convention.
- This slice closes issue #22 by combining implementation with explicit disposition, not by pretending `@reboot`, notifications, run-parts, anacron, or cron's DOM/DOW OR semantics are already supported.

## Assumptions

- Crontab import targets Unix cron semantics even when a file is previewed from another platform.
- The existing run-as implementation remains authoritative; S028 only maps system-crontab identity into it.
- Standard input is bounded by the existing local API request and SQLite practical limits; no new streaming protocol is required.
- Imported shell paths are operator-controlled and may fail at execution if absent on the target host, just as any imported command may.

## Out of Scope

- Boot-event triggers, notification delivery, directory discovery for run-parts, or anacron parsing.
- A composite schedule model for cron's restricted DOM-or-DOW behavior.
- Quartz year fields, calendars, misfire instructions, or unsupported `L`, `W`, and `#` composites.
- Windows run-as implementation or automatic installation of Unix shells.
- General secret storage, environment encryption, or command sandboxing.
