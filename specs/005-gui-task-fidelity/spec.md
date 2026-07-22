# Feature Specification: GUI task fidelity — schedule round-trip and group assignment

**Feature Branch**: `005-gui-task-fidelity`

**Created**: 2026-07-22

**Status**: Draft

**Input**: User description: "GUI task fidelity: schedule round-trip and group
assignment. Fixes GitHub issues #3 and #4 against v0.3.0. (1) The Edit Task
dialog always shows Mode as Recurring and blanks the Schedule and one-off
Date/Time fields, so the user cannot see what a task is currently set to and a
user who 'corrects' the mode silently rewrites the schedule. (2) Groups cannot
be populated from the GUI: the task editor has no group field and the Groups tab
shows no member tasks and no assign action; additionally no interface can move a
task back out of a group, because an empty group value means 'leave unchanged'."

## Clarifications

### Session 2026-07-22

Answered under the Build-Phase Autopilot Protocol decision policy
(constitution principle V): each was resolved against the constitution, the
master specification, and existing code patterns rather than escalated. The
rationale is recorded with each answer.

- Q: Is the retained schedule phrase authoritative for execution, or a
  display-and-round-trip aid only? → A: **Aid only.** The stored recurrence
  remains the single source of truth for evaluation; nothing in the execution
  path may read the phrase. *Rationale*: making the phrase authoritative would
  create two sources of truth for timing that can silently disagree, which
  principle II's determinism requirement and the master specification's
  UTC-internal model both forbid. The phrase is written once at creation and
  read only by clients.
- Q: When an operator switches an existing task's mode (recurring ↔ one-off),
  may they save with the new mode's timing fields left blank? → A: **No.** Blank
  timing fields mean "leave the schedule alone", which is impossible once the
  mode has changed — the existing schedule is the other kind. Switching mode
  makes the new mode's timing fields required. *Rationale*: the "blank keeps the
  existing schedule" affordance exists to let an operator edit a command without
  restating the schedule; carrying it across a mode switch would either silently
  ignore the switch or produce a task whose mode and schedule disagree.
- Q: What does the group selector show when a task's stored group cannot be
  resolved? → A: **"No group"**, with no error. *Rationale*: deleting a group
  already releases its tasks to ungrouped, so a dangling reference should not
  occur; treating it as ungrouped is the defensive reading that matches what the
  operator would see everywhere else, and failing the whole dialog over a
  display-only lookup would violate principle III's actionable-error rule for no
  gain.
- Q: Is the ungrouped area in the group hierarchy always shown, or only when it
  has members? → A: **Always shown.** *Rationale*: it is the visible destination
  for removing a task from a group (Story 3). A node that appears only once it
  already has members makes the un-group operation undiscoverable in exactly the
  state where the operator first needs it — the same discoverability failure
  that produced the original report.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See a task's real schedule when editing it (Priority: P1)

An operator opens an existing task to change something about it — the command,
an argument, the timezone. The dialog shows that task as it actually is: the
correct mode (one-off or recurring), and either the schedule phrase that drives
it or the concrete date and time it will fire, expressed in the task's own
timezone. Nothing about when the task runs changes unless the operator changes
it.

**Why this priority**: This is a correctness and trust defect, not a
convenience gap. Today the dialog asserts something false about every one-off
task, and an operator who acts on what they are shown destroys their own
schedule. It also blocks every other edit: an operator cannot safely change a
task's command without first proving to themselves that the schedule survived.

**Independent Test**: Create one recurring and one one-off task, reopen each for
editing, and confirm the dialog reflects what was created. Save without touching
anything and confirm the next-run times are unchanged. Delivers a trustworthy
edit path with no group work whatsoever.

**Acceptance Scenarios**:

1. **Given** a recurring task scheduled with the phrase "every day at 09:00",
   **When** the operator opens it for editing, **Then** Mode reads "Recurring"
   and the Schedule field contains a phrase that describes that same recurrence.
2. **Given** a one-off task set to fire at a specific future date and time,
   **When** the operator opens it for editing, **Then** Mode reads "One-off" and
   the Date and Time fields show that instant expressed in the task's timezone.
3. **Given** any existing task, **When** the operator opens it for editing and
   immediately saves without changing a field, **Then** the task's upcoming run
   times are identical to what they were before.
4. **Given** any existing task, **When** the operator opens it for editing and
   immediately cancels, **Then** the dialog closes without prompting to discard
   changes, because nothing was changed.
5. **Given** a sub-daily interval task created with an explicit first-cycle
   anchor, **When** the operator opens it for editing, **Then** the anchor is
   shown in the field dedicated to it and re-saving does not duplicate or drop
   it.
6. **Given** a recurring task, **When** the operator changes only its timezone
   and saves, **Then** the recurrence is re-interpreted in the new timezone
   rather than remaining anchored to the old one.
7. **Given** a recurring task open for editing, **When** the operator switches
   Mode to "One-off" but leaves the date and time empty, **Then** saving is
   blocked until they supply them.

---

### User Story 2 - Put tasks into groups from the GUI (Priority: P2)

An operator who has organized their work into nested groups can assign a task to
a group while creating or editing it, and can see at a glance which tasks belong
where. Groups stop being empty containers and start doing what the product
advertises: cascading enable and disable over real members.

**Why this priority**: A headline capability is entirely unreachable for the
GUI's whole audience. It ranks below Story 1 only because Story 1 is active data
misrepresentation while this is an absent capability with a command-line
workaround.

**Independent Test**: Create a group, create or edit a task and assign it to that
group, and confirm the assignment is visible in the task list and reflected in
the group's membership. Disable the group and confirm the member task stops
being eligible to run. Testable with no dependence on Story 1.

**Acceptance Scenarios**:

1. **Given** at least one group exists, **When** the operator creates a task,
   **Then** they can choose a group for it, and the created task belongs to that
   group.
2. **Given** an existing task in no group, **When** the operator edits it and
   chooses a group, **Then** the task becomes a member of that group.
3. **Given** a task in group A, **When** the operator edits it and chooses group
   B, **Then** the task belongs to B and no longer to A.
4. **Given** a nested group hierarchy, **When** the operator picks a group,
   **Then** the choices show enough of each group's position in the hierarchy to
   distinguish two groups that share a name at different levels.
5. **Given** a task assigned to a group, **When** the operator views the task
   list, **Then** the task's group is shown without opening the editor.
6. **Given** a task assigned to a group, **When** the group is disabled,
   **Then** the task is suppressed by the existing cascade.

---

### User Story 3 - Take a task back out of a group (Priority: P2)

An operator who put a task in the wrong group, or who no longer wants a task
grouped at all, can return it to ungrouped from the GUI and see it land somewhere
visible rather than disappearing from the group view.

**Why this priority**: Equal in priority to Story 2 and inseparable from it in
practice — an assignment capability with no reverse is a trap. It is stated
separately because it is the only part of this feature that changes behavior
shared with the command line, and because no interface can do it today.

**Independent Test**: Assign a task to a group, then remove it from all groups
and confirm it is listed as ungrouped and still runs on its own schedule.

**Acceptance Scenarios**:

1. **Given** a task that belongs to a group, **When** the operator removes it
   from all groups, **Then** the task belongs to no group and continues to run on
   its own schedule.
2. **Given** a task that belongs to no group, **When** the operator views the
   group hierarchy, **Then** the task is visible under a clearly labeled
   ungrouped area rather than absent from the view.
3. **Given** an operator using the command line, **When** they explicitly request
   that a task be removed from its group, **Then** the task is ungrouped — and
   when they simply do not mention the group at all, membership is left
   unchanged.

---

### User Story 4 - See group membership in the group hierarchy (Priority: P3)

An operator looking at the group hierarchy sees the tasks each group contains,
not just the group names, and can move a selected task to a different group from
that same view.

**Why this priority**: Discoverability and verification. Stories 2 and 3 make
assignment possible; this makes the result legible and gives the operator a
second, more natural place to do the work. Lowest priority because the task list
and editor already answer "what group is this task in" once Story 2 lands.

**Independent Test**: With tasks spread across two groups and one ungrouped,
open the group hierarchy and confirm every task appears exactly once under the
right parent; move one task to another group from that view and confirm it
relocates.

**Acceptance Scenarios**:

1. **Given** groups containing tasks, **When** the operator views the group
   hierarchy, **Then** each group's member tasks are listed beneath it and are
   visually distinguishable from groups.
2. **Given** a task selected in the group hierarchy, **When** the operator moves
   it to another group, **Then** the change takes effect and the view updates.
3. **Given** the group hierarchy is open, **When** a task's group membership
   changes from anywhere, **Then** the view reflects it without the operator
   asking for a refresh.
4. **Given** a task is selected rather than a group, **When** the operator uses a
   group-only action, **Then** the action does not fire against the task.

---

### Edge Cases

- The operator switches an existing task from recurring to one-off, or back,
  without filling in the new mode's timing: saving is blocked until they do,
  because there is no existing schedule of the new kind to fall back on
  (FR-011b).
- A task references a group that cannot be resolved: it is presented as
  ungrouped everywhere rather than erroring (FR-019a).
- No task is currently ungrouped: the ungrouped area is still shown, empty, so
  the operator can see where an un-grouped task would land (FR-019).
- A one-off task whose fire time has already passed: reopening it shows the
  stored past date and time, and the existing rule that a one-off must be in the
  future applies only when the operator actually changes it.
- Assigning a task to a group that no longer exists (deleted between the dialog
  opening and the save): the operator gets a clear, actionable message naming the
  field, not a generic failure.
- Deleting a group that contains tasks: the existing behavior stands — the group
  is removed and its tasks become ungrouped, where they must now be visible in
  the ungrouped area.
- A group hierarchy with a group and a task sharing the same identifier: the
  hierarchy view must not confuse one for the other.
- An operator opens the editor, the task is changed or deleted elsewhere, and
  then they save: the existing conflict behavior is unchanged by this feature.
- The detail lookup that populates the editor fails: the operator can still open
  and use the dialog, with a clear indication that the current schedule could not
  be read, rather than being blocked from editing.

## Requirements *(mandatory)*

### Functional Requirements

#### Schedule fidelity

- **FR-001**: The system MUST retain, for every recurring schedule it creates,
  the human-readable phrase the schedule was created from, and MUST make that
  phrase available when the task is retrieved.
- **FR-002**: The schema change that stores the phrase MUST be forward-only and
  non-destructive: applying it MUST NOT read, rewrite, or drop any existing
  value, and no stored task timing may shift as a result.
- **FR-006**: The task editor MUST populate its mode, schedule phrase, one-off
  date, one-off time, and first-cycle anchor fields from the task's actual
  schedule when an existing task is opened.
- **FR-007**: The editor MUST express a one-off task's date and time in that
  task's timezone.
- **FR-008**: An editor opened on an existing task and not modified MUST be
  treated as unmodified: saving MUST leave the schedule untouched, and cancelling
  MUST NOT prompt about discarding changes.
- **FR-009**: When the editor cannot obtain a task's schedule, it MUST still
  open, MUST leave the schedule fields empty, and MUST tell the operator that the
  current schedule could not be read.
- **FR-011**: Changing only a task's timezone MUST re-interpret its recurrence in
  the new timezone.
- **FR-011a**: The retained phrase MUST NOT influence execution. The stored
  recurrence remains the sole input to scheduling decisions; no component on the
  execution path may read the phrase.
- **FR-011b**: When the operator switches an existing task's mode between
  recurring and one-off, the timing fields for the newly selected mode MUST be
  required before the change can be saved. The "blank leaves the schedule
  unchanged" allowance applies only while the mode is unchanged.

#### Group assignment

- **FR-012**: The task editor MUST let the operator choose the group a task
  belongs to, including an explicit "no group" choice, when creating and when
  editing a task.
- **FR-013**: Group choices MUST convey each group's position in the hierarchy
  well enough to distinguish same-named groups at different levels.
- **FR-014**: The system MUST be able to express three distinct intents when
  updating a task: leave group membership unchanged, assign it to a named group,
  and remove it from all groups. An omitted group MUST mean "unchanged".
- **FR-015**: The command line MUST be able to express all three intents, and its
  existing behavior — that not mentioning the group leaves membership unchanged —
  MUST be preserved.
- **FR-016**: Assigning a task to a group that does not exist MUST be rejected
  with an actionable message naming the group field, not an internal error.
- **FR-017**: The task list MUST show each task's group.
- **FR-018**: The group hierarchy view MUST show each group's member tasks
  beneath it, visually distinguishable from groups.
- **FR-019**: The group hierarchy view MUST show tasks that belong to no group
  under a clearly labeled ungrouped area. That area MUST be present whether or
  not it currently has members, so the destination for removing a task from a
  group is always visible.
- **FR-019a**: Where a task's stored group cannot be resolved, every view and the
  editor MUST present the task as belonging to no group, without raising an
  error.
- **FR-020**: The operator MUST be able to move a selected task to another group,
  or to no group, from the group hierarchy view.
- **FR-021**: Actions that apply only to groups MUST NOT act on a selected task.
- **FR-022**: Group membership changes MUST appear in every affected view without
  the operator requesting a refresh.

#### Verification

- **FR-023**: Every defect fixed by this feature MUST ship with a regression test
  that fails before the fix and passes after it.
- **FR-024**: The schema change MUST be covered by a test proving the
  forward-only, non-destructive property of FR-002 — that stored schedules
  survive it intact and that applying it twice is a no-op. This holds
  independently of whether any deployment currently has data worth keeping; it
  guards the mechanism, which the constitution names a safety-critical surface.

### Key Entities

- **Task**: a unit of work. Gains no new attribute here, but its group
  membership becomes settable to "none" as a distinct, expressible intent rather
  than an unreachable state.
- **Schedule**: the timing definition for a task. Gains the human-readable
  phrase it was created from, retained alongside the recurrence and the
  plain-language summary it already carries. The phrase is what the operator
  typed; the summary is what the system says back; the recurrence is what the
  engine evaluates. Only the phrase is new, and it is inert with respect to
  execution (FR-011a) — it exists so a client can put the operator's own words
  back in front of them.
- **Group**: a named container in a nested hierarchy. Unchanged, but its
  membership becomes visible and editable from the GUI.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of tasks, recurring and one-off, display their actual mode and
  timing when opened for editing.
- **SC-002**: 100% of tasks opened for editing and saved without modification
  have identical upcoming run times before and after.
- **SC-003**: An operator can create a group, put a task in it, and confirm the
  membership without leaving the GUI and without consulting documentation.
- **SC-004**: An operator can remove a task from all groups and locate it
  afterwards, entirely within the GUI.
- **SC-005**: Every task in the system appears exactly once in the group
  hierarchy view — under its group or under the ungrouped area.
- **SC-006**: Group membership changes are visible in every affected view within
  one second of being made, with no manual refresh.
- **SC-007**: Applying the schema change to a database that already contains
  tasks changes zero upcoming run times.
- **SC-008**: The two reported defects are reproducible before the change and not
  reproducible after it, demonstrated by tests.

## Assumptions

- The operator is a single local user administering their own machine; there is
  no multi-user or permission dimension to group membership.
- **There is no installed base to carry forward.** The software has no working
  deployments; the only databases that exist are the maintainers' own, and none
  of them is functional. So this feature does not reconstruct phrases for
  schedules stored before it shipped — every recurring schedule is created
  through the parser, which retains its phrase, and a database predating that
  simply shows a blank schedule field on edit (which means "keep unchanged" and
  is harmless). FR-002 and FR-024 still hold: they guard the schema-change
  *mechanism*, which the constitution names safety-critical and which matters as
  soon as there is a real user.
- The existing plain-language schedule vocabulary is the vocabulary for retained
  phrases. This feature adds no new schedule forms.
- One-off schedules need no retained phrase: their date and time are recovered
  from the stored instant.
- The existing live-update mechanism already carries task changes to the views,
  so FR-022 is satisfied by making the views read membership rather than by a new
  notification path.
- Group deletion semantics (children cascade, member tasks become ungrouped) are
  existing behavior and are not changed here.
- Renaming and reparenting groups from the GUI, drag-and-drop assignment, and any
  change to how schedules are evaluated or executed are out of scope.
- The GUI and the command line remain thin clients over the same daemon; any
  capability added for the GUI is added at the shared boundary, so the command
  line gains it too.
