# Feature Specification: Cron interoperability and calendar-anomaly policy

**Feature Branch**: `008-cron-interop`

**Created**: 2026-07-23

**Status**: Draft

**Input**: User description: "Cron interoperability and calendar-anomaly policy.
Give users with an existing crontab a native, bidirectional cron conversion
surface exposed only through the CLI (`gosched cron explain|import|export`) —
cron remains an interchange format at the boundary and is never an authoring
syntax in the GUI or the human phrase grammar. Import reads a crontab file or
stdin, prints each line's cron expression, the human phrase it maps to, and the
next few run times, and creates the tasks; `--dry-run` creates nothing and
serves as the migration preview. `explain` translates a single expression with
no side effects. `export` emits crontab lines for tasks whose schedule cron can
carry, and an explicit commented refusal where it cannot — never a silent
approximation. Fidelity must be stated rather than implied: cron has no timezone
(import takes `--timezone`, defaulting to the task default), no catch-up, no
overlap policy and no restart recovery (imported tasks take project defaults,
reported in the import summary), and non-standard extensions (`@reboot`, Quartz
seconds fields, `L`/`W`/`#`) plus `MAILTO` and shell variable assignments must be
reported as explicit warnings rather than dropped. Prerequisite in the same
feature: the human-readable schedule grammar today has no by-date monthly form
and no yearly frequency, so ordinary cron lines like `0 9 1 * *` and
`0 0 29 2 *` have no target representation; add by-date monthly ("on the 15th of
every month") and yearly ("every year on february 29", "every 12 months") forms.
Adding those raises the calendar-anomaly question, so also add a per-task
missing-date policy — skip (default, current behavior), last-valid (fall back to
the last day that exists: Feb 29 to Feb 28, the 31st to the 30th, the 5th Friday
to the last Friday), and next-valid (roll forward into the following period) —
surfaced in the CLI task flags, the GUI editor's advanced section, and the task's
next-run output, with a forward-only additive store migration. Human summaries
and the calendar view must stop asserting 'every month' for a rule that does not
fire every month and must name the policy instead. Out of scope for this feature:
DST anchoring options (wall-clock / elapsed / utc) and the per-task skipped-hour
/ repeated-hour policy, which remain open on issue #8. Closes issue #12 and the
missing-date half of issue #8. Also in scope as a documentation fix: hyperlink
the ShruggieTech attribution in the README footer and the other project-facing
documents to https://shruggie.tech (issue #9)."

## Clarifications

### Session 2026-07-23

Answered under the Build-Phase Autopilot Protocol decision policy (constitution
principle V): each was resolved against the constitution, the master
specification, and existing code patterns rather than escalated. The rationale is
recorded with each answer.

- Q: Does a cron expression become a schedule directly, or by way of the
  human-readable phrase? → A: **By way of the phrase, always.** A cron expression
  is converted to the phrase a user would have typed, and that phrase goes
  through the existing schedule parser. An expression with no phrase is declined,
  never converted by a second path. *Rationale*: this is the single decision that
  makes the preview trustworthy — what the operator is shown in the preview is
  literally what is parsed and stored, so the preview cannot disagree with the
  result. It also keeps one implementation of phrase→recurrence rather than two
  that can drift, and it makes "cron is not an authoring syntax" structural
  rather than a matter of discipline: the conversion has no privileged path into
  the engine that a user's own input does not have.
- Q: How is a cron line that restricts both day-of-month and day-of-week handled
  (`0 0 13 * 5`, which cron fires on the 13th *or* on any Friday)? → A:
  **Declined, with the reason named.** *Rationale*: cron's OR semantics for that
  combination have no faithful equivalent in the recurrence model, which treats
  the two restrictions as an intersection. Emitting the intersection would
  silently produce a task that fires a handful of times a year where the original
  fired weekly — exactly the silent meaning change this feature exists to
  prevent. FR-002 already forbids approximation; this is an instance of it.
- Q: How is a step value that does not divide its range evenly handled (`*/7` in
  the minutes field, which cron restarts at the top of each hour rather than
  running continuously)? → A: **Accepted only when the step divides its range
  evenly; otherwise declined with the reason.** *Rationale*: `*/15` and `*/5` are
  exactly equivalent to a fixed interval anchored at the top of the hour, and
  those are the overwhelming majority of real crontab lines. `*/7` is not — it
  fires at :00, :07 … :56 and then again at :00, a 4-minute gap the interval model
  cannot reproduce. Declining is honest and cheap; representing it would require
  a by-minute-list phrase form that no user has asked for and that would enlarge
  the authoring grammar to serve an import edge case.
- Q: What does export do with a task cron cannot represent, and with a disabled
  task? → A: **Both produce a commented refusal naming the task and the reason;
  neither produces a bare crontab line.** A disabled task is declined because
  cron has no disabled state. *Rationale*: the export is a migration and
  diffing tool, so a silently omitted task is worse than a visible refusal — the
  reader must be able to account for every task. Emitting a live line for a task
  that is currently disabled would activate work the operator had deliberately
  stopped.
- Q: Does changing a task's schedule phrase reset its missing-date policy, or are
  the two independent? → A: **Independent.** Editing the phrase leaves the policy
  alone, and editing the policy leaves the phrase alone. *Rationale*: the policy
  is an answer to "what should this task do when the calendar cannot honor it",
  which is a property of the operator's intent, not of the particular phrase.
  Resetting it on an unrelated edit would silently change run times — the failure
  mode issue #4 already produced once in the task editor, and the reason the
  editor now round-trips every field it displays.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Move an existing crontab across without retyping it (Priority: P1)

An operator running scheduled work on cron today points the tool at their
crontab and sees, line by line, what each entry means in plain language and when
it would next run — before anything is created. Satisfied that the translation
is right, they run it again without the preview flag and the tasks exist.

**Why this priority**: This is the adoption path. The population that would
benefit from restart recovery, catch-up, overlap policies and a GUI is precisely
the population that already has a crontab, and today the only way across is to
retype every job by hand and hope nothing changed meaning.

**Independent Test**: Point the import at a sample crontab in preview mode and
confirm every line is echoed with its phrase and next run times and that no task
was created; run it again without preview and confirm the tasks exist with those
schedules.

**Acceptance Scenarios**:

1. **Given** a crontab containing `0 9 * * 1-5 /usr/bin/report`, **When** the
   operator previews an import of that file, **Then** the output shows the cron
   expression, the phrase "weekdays at 09:00", and the next several run times,
   and no task is created.
2. **Given** the same crontab, **When** the operator imports it without the
   preview flag, **Then** a task exists whose recurrence produces exactly the
   run times shown in the preview.
3. **Given** a crontab line the tool cannot represent, **When** the operator
   previews the import, **Then** that line is reported as unsupported with the
   reason, the remaining lines are still processed, and the summary counts both.
4. **Given** a crontab containing `MAILTO=ops@example.com` and a shell variable
   assignment, **When** the operator imports it, **Then** each is reported as an
   explicit warning naming what was not carried across, rather than dropped
   silently.
5. **Given** an import with no timezone specified, **When** the operator imports,
   **Then** the tasks take the default timezone and the summary states that cron
   carries no timezone, that a timezone was chosen, and which one.
6. **Given** an import, **When** it completes, **Then** the summary states which
   behaviors the imported tasks gained that cron does not have (catch-up,
   overlap policy, restart recovery) and what defaults they were given.

---

### User Story 2 - Schedule by calendar date, and know what happens when the date does not exist (Priority: P1)

An operator schedules a task for the 31st of every month, or for the 29th of
February, or for the fifth Friday of the month — and can state what the task
should do in a period that has no such date: skip it, use the closest earlier
date that exists, or roll into the next period. Whatever they choose, the
schedule describes itself honestly wherever it is displayed.

**Why this priority**: Ordinary cron lines address dates by number, so cron
import cannot be complete without these forms. Adding them without a stated
policy would silently create tasks that fire seven months in twelve. The defect
already exists for the fifth-weekday form, whose summary claims "every month"
for a rule that fires four times a year.

**Independent Test**: Create a by-date task under each policy and compare its
next run times across a year containing a short month and a non-leap February;
confirm the displayed description names the policy.

**Acceptance Scenarios**:

1. **Given** a schedule for the 31st of every month with the default policy,
   **When** its next runs are listed across a year, **Then** it runs only in
   months that have a 31st, and its description says so rather than claiming
   every month.
2. **Given** the same schedule set to fall back to the last valid date, **When**
   its next runs are listed, **Then** it runs in every month, on the 31st, 30th,
   or last day of February as applicable.
3. **Given** the same schedule set to roll forward, **When** its next runs are
   listed, **Then** a period without a 31st produces a run on the first day of
   the following period.
4. **Given** a yearly schedule on 29 February, **When** its next runs are listed
   across a non-leap year under each policy, **Then** the result is respectively
   no run, a run on 28 February, and a run on 1 March.
5. **Given** an ordinal-weekday schedule for the fifth Friday of the month,
   **When** it is set to fall back to the last valid occurrence, **Then** a month
   with only four Fridays runs on the fourth (last) Friday.
6. **Given** any of these tasks, **When** the operator views it in either
   interface, **Then** the policy is visible and editable, and the human
   description names it.
7. **Given** a database created before this feature, **When** it is opened,
   **Then** every existing schedule keeps its stored timing exactly and takes the
   default policy, which reproduces its previous behavior.

---

### User Story 3 - Understand one cron expression without creating anything (Priority: P2)

Someone who cannot read `0 0 1 * *` at a glance asks the tool what it means and
gets a plain-language phrase and the next few times it would fire.

**Why this priority**: Useful on its own to anyone maintaining a crontab, and it
is the smallest honest demonstration that the conversion is trustworthy — the
thing an evaluator tries first. It is a strict subset of the import path, so it
costs almost nothing once import exists.

**Independent Test**: Ask for an explanation of several expressions and confirm
the phrase and run times are correct and that nothing was created or changed.

**Acceptance Scenarios**:

1. **Given** a standard five-field cron expression, **When** the operator asks
   for an explanation, **Then** the phrase and the next several run times are
   printed and no task, schedule, or run record is created.
2. **Given** an expression using a non-standard extension, **When** the operator
   asks for an explanation, **Then** the response names the extension and states
   that it is unsupported, rather than guessing.
3. **Given** a malformed expression, **When** the operator asks for an
   explanation, **Then** the error states which field is wrong and what was
   expected, and the command reports a validation failure.

---

### User Story 4 - Get the jobs back out (Priority: P3)

An operator evaluating the tool, or comparing two machines, asks for the current
task set as crontab lines. Everything cron can carry comes out as a crontab
line; everything it cannot comes out as a comment saying so and why.

**Why this priority**: It removes the lock-in objection at low cost — the
conversion is a pure function of the stored schedule — and doubles as an
operational diff tool. It is last because nothing depends on it.

**Independent Test**: Export a task set containing both expressible and
inexpressible schedules and confirm each appears either as a valid crontab line
or as a commented refusal naming the reason.

**Acceptance Scenarios**:

1. **Given** a task whose recurrence cron can express, **When** the operator
   exports, **Then** a crontab line is emitted whose run times match the task's.
2. **Given** a one-off task, **When** the operator exports, **Then** a comment
   states that cron cannot express a single-occurrence schedule and no line is
   emitted for it.
3. **Given** a task whose behavior depends on a policy cron has no notion of,
   **When** the operator exports, **Then** the emitted comment names what is lost
   rather than approximating it.
4. **Given** a crontab that was imported, **When** it is exported again, **Then**
   the resulting lines produce the same run times as the originals over a window
   crossing both a daylight-saving transition and a month boundary.

---

### User Story 5 - Find out who maintains this (Priority: P3)

A reader who reaches the end of a project document and wants to know who is
behind the project can click the organization's name.

**Why this priority**: A one-line documentation fix with no dependencies. It is
included because it is the last unlinked proper noun in a document that links
every other one, and it costs nothing to carry alongside the docs this feature
already touches.

**Independent Test**: Open each project-facing document and confirm the
attribution is a link, and that every document points to the same destination.

**Acceptance Scenarios**:

1. **Given** the project README, **When** a reader reaches the footer, **Then**
   the organization name is a link to the organization's site.
2. **Given** the other project-facing documents that name the organization,
   **When** they are read, **Then** they point to the same destination as the
   README.

---

### Edge Cases

- A crontab line whose schedule is valid but whose command is empty, or which is
  a comment or blank line: comments and blanks are skipped silently; a schedule
  with no command is a reported error on that line, not a created task.
- A cron field combining a day-of-month restriction and a day-of-week
  restriction (cron's documented OR semantics, e.g. `0 0 13 * 5`): declined with
  the reason named. It is never reinterpreted as an intersection.
- A cron step that does not divide its range evenly (`*/7` on minutes): declined
  with the reason named, because cron restarts the sequence at the top of each
  hour and a fixed interval does not. Steps that do divide evenly (`*/5`,
  `*/15`, `*/30`) are accepted and are exactly equivalent.
- An import naming a timezone that does not exist: fail fast before creating
  anything, naming the field.
- An import of a file that is partly unsupported: nothing is created in preview
  mode; in a real import the supported lines are created and the unsupported
  ones are reported, with the summary stating both counts.
- The same crontab imported twice: the second import creates a second set of
  tasks. Deduplication is not attempted; the summary makes the count visible so
  the operator can see what happened.
- A by-date schedule for the 30th under the roll-forward policy in February:
  rolling forward lands on 1 or 2 March; the run must not be lost or duplicated
  in the following period's own occurrence.
- A missing-date policy applied to a schedule that has no date component at all
  (an interval, a weekday rule): the setting is inert and must not alter any run
  time.
- A daylight-saving transition on a date resolved by the fall-back or
  roll-forward policy: the existing transition rules still apply to the resolved
  instant, unchanged by this feature.

## Requirements *(mandatory)*

### Terminology

- **Declined** — the outcome when the system will not convert an input. Always
  named and always visible; never a silent omission and, on its own, never a
  command failure.
- **Unsupported** — one reason for a decline: the input is well-formed cron but
  outside what can be represented (`@reboot`, `L`/`W`/`#`, a non-dividing step, a
  day-of-month plus day-of-week combination).
- **Warning** — something not carried across that does not stop the line being
  converted (`MAILTO`, a variable assignment).
- **Error** — a failure of the run itself: unreadable input, an unknown
  timezone, or a malformed expression.

### Functional Requirements

#### Cron conversion

- **FR-001**: The system MUST accept standard five-field cron expressions,
  including wildcards, lists, ranges, steps, and month and day-of-week names, and
  the common named shorthands (`@hourly`, `@daily`, `@midnight`, `@weekly`,
  `@monthly`, `@yearly`, `@annually`).
- **FR-002**: The system MUST recognize, name, and report as unsupported every
  input it cannot faithfully represent, including `@reboot`, seconds-field
  (six-field) expressions, and the `L`, `W`, and `#` extensions. It MUST NOT
  approximate, guess, or silently drop such an input.
- **FR-003**: Users MUST be able to translate a single cron expression into a
  plain-language phrase plus its next several run times without creating or
  modifying anything.
- **FR-003a**: Every conversion from cron MUST produce a human-readable phrase
  first and derive the schedule from that phrase by the same route a user typing
  it would take. An expression with no phrase MUST be declined. There MUST be no
  second route from cron into a stored schedule.
- **FR-003b**: The system MUST decline, naming the reason, any expression that
  restricts both day-of-month and day-of-week, and any step value that does not
  divide its field's range evenly.
- **FR-004**: Users MUST be able to import a crontab from a file or from standard
  input, seeing for each line the cron expression, the phrase it maps to, and its
  next run times.
- **FR-005**: The import MUST offer a preview mode that produces the identical
  report but creates nothing.
- **FR-005a**: In a real import, a declined line MUST NOT prevent the supported
  lines from being created. If creating a task fails partway through, the tasks
  already created MUST remain and the summary MUST report the failure alongside
  the created count; the import MUST NOT silently roll back or silently
  continue.
- **FR-006**: The import MUST carry across each line's command, arguments, and
  working directory as the task's payload.
- **FR-007**: The import MUST report `MAILTO`, shell variable assignments, and
  any other crontab directive it does not carry across, as explicit per-line
  warnings.
- **FR-008**: The import MUST accept an explicit timezone and, when none is
  given, apply the project's default and state in the summary which timezone was
  applied and that cron itself carries none.
- **FR-009**: The import summary MUST state that imported tasks gain catch-up,
  overlap-policy, and restart-recovery behavior that cron has no notion of, and
  which defaults they received.
- **FR-010**: The import MUST report per-run totals: lines read, tasks created
  (zero in preview), lines skipped as comments or blanks, and lines reported as
  unsupported or erroneous.
- **FR-010a**: A declined or unsupported line MUST be a reported outcome, not a
  failure: an import or preview that read its input successfully MUST report
  success even when every line was declined. Failure MUST be reserved for inputs
  the run could not process at all — an unreadable file, an unknown timezone, or
  a malformed expression given to the single-expression translator — and MUST
  follow the project's existing exit-code contract for validation versus runtime
  errors.
- **FR-011**: Users MUST be able to export the task set, or a single task, as
  crontab lines.
- **FR-011a**: Exporting an empty task set MUST succeed and produce output that
  is recognizably an empty export rather than nothing at all.
- **FR-012**: The export MUST emit a line only where cron can carry the
  schedule's meaning, and MUST otherwise emit a comment naming the task and the
  reason it was declined. Every task MUST appear in the output as either a line
  or a refusal; none may be silently omitted.
- **FR-012a**: A disabled task MUST be exported as a commented refusal, never as
  a live crontab line, because cron has no disabled state.
- **FR-013**: A schedule that was imported from cron and then exported MUST
  produce run times identical to the original expression's over a window
  spanning at least one daylight-saving transition and one month boundary.
- **FR-014**: Cron MUST NOT become an authoring syntax: it MUST NOT be accepted
  where a human-readable schedule phrase is accepted, and MUST NOT appear as an
  input in the graphical interface.

#### Schedule grammar

- **FR-015**: Users MUST be able to express a monthly schedule by calendar date
  (for example, the 15th of every month), with an optional time of day.
- **FR-016**: Users MUST be able to express a yearly schedule by month and date
  (for example, 29 February each year), with an optional time of day.
- **FR-017**: Users MUST be able to express a fixed interval in months or years
  (for example, every 12 months).
- **FR-018**: Every newly accepted phrase MUST round-trip: the phrase the user
  typed is retained, and the description the system gives back describes the same
  rule.

#### Missing-date policy

- **FR-019**: Each task MUST carry a missing-date policy with three settings:
  skip the period, fall back to the last valid date in the period, or roll
  forward into the next period.
- **FR-019a**: Rolling forward MUST land on the first day of the following
  period at the rule's time of day, and MUST NOT suppress, displace, or
  duplicate that following period's own occurrence.
- **FR-020**: The default MUST be to skip, which MUST reproduce the behavior of
  schedules created before this feature, bit for bit.
- **FR-021**: The policy MUST apply to by-date monthly rules, yearly by-date
  rules, and ordinal-weekday rules (where "last valid" means the last occurrence
  of that weekday in the period).
- **FR-022**: The policy MUST be settable and viewable from the command line and
  from the graphical editor, and MUST be visible wherever a task's schedule is
  shown.
- **FR-023**: Any human-readable description of a rule that does not fire in
  every period MUST say so and MUST name the policy in effect. No description may
  assert "every month" for a rule that skips months.
- **FR-024**: The policy MUST have no effect on schedules with no date component.
- **FR-024a**: The policy and the schedule phrase MUST be independently editable:
  changing one MUST NOT reset or alter the other.
- **FR-025**: Existing daylight-saving resolution MUST continue to apply
  unchanged to the instant the policy resolves to.
- **FR-026**: Stored schedules MUST migrate forward without rewriting any stored
  timing value; no existing task's run times may change as a result of this
  feature.

#### Documentation

- **FR-027**: The documentation MUST state the conversion's fidelity explicitly,
  listing what is supported, what is declined, and what cron cannot carry in
  either direction.
- **FR-028**: Project-facing documents MUST link the maintaining organization's
  name to a single, consistent destination.

### Key Entities

- **Cron expression**: An externally supplied timing string in the crontab
  interchange format, together with the command line it governs. Input and output
  only; never stored.
- **Schedule**: The existing stored timing definition for a task. Gains a
  missing-date policy and the ability to be described by calendar date or by
  year.
- **Missing-date policy**: A per-task setting stating what a schedule does in a
  period that has no matching date. One of three values, defaulting to the
  behavior that predates this feature, and editable independently of the
  schedule phrase.
- **Import report**: The per-line and summary account of a conversion run —
  expression, phrase, next runs, warnings, and totals. Produced identically in
  preview and in a real import; the only difference is whether tasks were
  created.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator with an existing crontab can see every line of it
  translated, with next run times, in a single command that changes nothing.
- **SC-002**: Every line of a crontab composed only of standard five-field
  expressions and common shorthands is either imported or reported with a
  reason; no line is silently dropped or silently altered.
- **SC-002a**: The phrase shown for a line in the preview is the phrase the
  created task reports afterwards, for every imported line.
- **SC-003**: A crontab imported and then exported produces run times identical
  to the original expressions over a window that includes a daylight-saving
  transition and a month boundary.
- **SC-004**: For each of the three missing-date policies, a schedule addressing
  a date that does not exist in every period produces the stated run times across
  a full calendar year, verified against real dates including a non-leap
  February.
- **SC-005**: No schedule description asserts that a rule fires in every period
  when it does not.
- **SC-006**: Opening a database created before this feature changes no stored
  timing value and no task's computed run times.
- **SC-007**: A reader can reach the maintaining organization from any
  project-facing document, and every such document points to the same place.

## Assumptions

- The population targeted is people migrating from standard Unix cron
  (Vixie/ISC-style five-field crontabs). Quartz, systemd timers, and Windows
  Task Scheduler are not conversion sources for this feature.
- Imported tasks are given a generated name derived from their command when the
  crontab offers no name, since crontabs do not name jobs.
- Cron's environment-variable lines are reported rather than applied, because
  applying them would change the meaning of every subsequent line in ways the
  operator cannot see in the preview.
- Deduplication against already-imported tasks is out of scope; a repeated
  import creates a second set, and the visible totals are how the operator
  notices.
- The default missing-date policy is the current behavior, so this feature
  cannot alter any existing task.
- DST anchoring (wall-clock versus elapsed-time versus UTC) and per-task
  skipped-hour/repeated-hour resolution remain out of scope and stay open on the
  existing advanced-options issue.
- The graphical interface gains the missing-date setting but gains no cron
  surface of any kind.
