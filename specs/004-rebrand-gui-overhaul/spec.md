# Feature Specification: Rebrand to go-schedule + GUI & Installer Overhaul

**Feature Branch**: `004-rebrand-gui-overhaul`

**Created**: 2026-06-20

**Status**: Draft

**Input**: User description: "Change the name of the repository from go-scheduler to go-schedule and update all references; replace the Windows install with a formal .msi system install that sets up the background service; replace Alerts with a unified Logs view (type filters, Dismiss All, on-disk persistence, click-through detail); remove the Triggers feature entirely; make the GUI refresh in real time across the board (remove manual Refresh); add a toggleable calendar view under Schedule."

## Clarifications

### Session 2026-06-20

- Q: What should the new "Logs" view show? → A: Full daemon log stream — surface the daemon's structured log records (info/warn/error) and the existing alert events together, persisted to a log file on disk, with type/severity filters and click-through detail.
- Q: How should Windows distribution change? → A: MSI only — the `.msi` becomes the sole Windows distribution; it installs to a system location, registers the daemon as an auto-start service, adds Start Menu shortcuts, and provides clean uninstall. The portable zip and "run from folder" guidance are removed.
- Q: How deep should Triggers removal go? → A: Full removal across all layers — GUI, CLI, API, engine dispatcher, store schema/table, and domain types, including a migration that drops existing trigger data.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consistent project name "go-schedule" (Priority: P1)

A user who clones the project, reads the documentation, or runs the tooling encounters the name
"go-schedule" consistently. No stale references to the former "go-scheduler" name remain in
documentation, build/release configuration, the module/import paths, or user-facing strings.

**Why this priority**: The rename is foundational — every other change in this feature touches
files that also carry the old name, and import paths must be consistent for the code to build at
all. Doing it first avoids churn and merge conflicts in later work.

**Independent Test**: Search the entire repository for the literal "go-scheduler"; the only
permitted matches are historical entries (e.g. existing CHANGELOG/spec history) explicitly
designated as immutable. The project builds, tests pass, and the GUI window title and CLI help
text read "go-schedule".

**Acceptance Scenarios**:

1. **Given** a fresh checkout, **When** a developer builds the daemon, CLI, and GUI, **Then**
   the build succeeds with the updated module/import path and no references to the old path.
2. **Given** the running GUI, **When** the main window opens, **Then** its title and branding
   read "go-schedule".
3. **Given** the published documentation (README, install guide, changelog header), **When** a
   user reads them, **Then** the product is referred to as "go-schedule" throughout.

---

### User Story 2 - Formal Windows installation via MSI (Priority: P1)

A Windows user downloads a single `.msi` installer, runs it, and the product is installed as a
proper system application: files placed in a standard install location, the daemon registered as
an auto-starting background service, Start Menu shortcuts created, and the GUI launchable from the
Start Menu. The user never has to extract a zip or run a loose `.exe` from their Downloads folder.
Uninstalling through the standard Windows mechanism removes everything cleanly.

**Why this priority**: The current "extract and double-click an exe" flow is the user's primary
complaint and undermines the "runs reliably in the background" promise. A formal install is the
expected baseline for desktop software on Windows.

**Independent Test**: On a clean Windows machine, run the `.msi`; verify the service is installed
and running, the app appears in the Start Menu and in "Apps & features", tasks fire after a
reboot without anyone logging in, and uninstall removes binaries, service registration, and
shortcuts (with a documented choice about user data).

**Acceptance Scenarios**:

1. **Given** the `.msi`, **When** the user runs it and completes the wizard, **Then** the
   product is installed to a standard system location (not the Downloads folder) and the daemon
   is registered as a background service set to start automatically.
2. **Given** a completed install, **When** the machine is rebooted and no user logs in, **Then**
   scheduled tasks continue to fire.
3. **Given** an installed product, **When** the user launches the app from the Start Menu,
   **Then** the GUI opens and connects to the already-running service without spawning a second
   daemon and without a visible console window.
4. **Given** an installed product, **When** the user uninstalls it via the standard Windows
   mechanism, **Then** binaries, the service registration, and shortcuts are removed cleanly.
5. **Given** the release process, **When** a Windows release is produced, **Then** the `.msi` is
   the published Windows artifact and no portable-zip "run from folder" path is offered.

---

### User Story 3 - Unified Logs view for troubleshooting (Priority: P2)

A user opens a single "Logs" view (replacing "Alerts") that shows everything noteworthy the
system has reported: structured daemon log records (informational, warning, error) together with
scheduler events (run failures, overlaps, missed runs). The user can filter by type/severity to
focus on, say, only Errors or only Warnings; click any entry to open a detail view showing the
full message and the context that caused it; and clear the list with "Dismiss All". Logs are also
written to a file on disk in the installed environment so they survive restarts and can be
inspected outside the GUI.

**Why this priority**: Troubleshooting value is the explicit goal. It depends on the rename
(US1) being settled but is independent of the installer and other GUI changes.

**Independent Test**: Trigger a failing task and an informational event; open the Logs view; the
entries appear with correct severities; filtering to "Errors" hides the informational entry;
clicking an error opens detail with the full message/cause; "Dismiss All" clears the view; and
the on-disk log file contains the same records.

**Acceptance Scenarios**:

1. **Given** the GUI, **When** the user opens the Logs view, **Then** it lists log entries of
   multiple severities (at minimum Info, Warning, Error) drawn from both daemon logs and
   scheduler events, newest first.
2. **Given** a populated Logs view, **When** the user selects the "Errors" filter, **Then** only
   error-severity entries are shown; selecting "Warnings" shows only warnings; clearing the
   filter shows all.
3. **Given** a log entry, **When** the user clicks it, **Then** a detail view shows the full
   message, severity, timestamp, source, and the cause/context (e.g. which task, exit code, or
   error chain).
4. **Given** logs are present, **When** the user chooses "Dismiss All", **Then** the visible log
   list is cleared.
5. **Given** an installed environment, **When** events are logged, **Then** they are persisted to
   a log file on disk whose location is documented, and the file survives a daemon restart.
6. **Given** new log entries are produced while the Logs view is open, **When** they occur,
   **Then** they appear in the view without a manual refresh.

---

### User Story 4 - Remove the Triggers feature (Priority: P2)

The Triggers feature is removed entirely from the product. Users no longer see a Triggers area in
the GUI, no trigger commands in the CLI, and no trigger endpoints in the API. Existing trigger
data is removed via a forward migration so installs continue to start cleanly.

**Why this priority**: Explicitly unwanted by the stakeholder; removing it reduces surface area
and simplifies the remaining work. Independent of the other changes.

**Independent Test**: After the change, the GUI has no Triggers tab, the CLI rejects/omits
trigger commands, the API has no trigger routes, and a daemon started against a database that
previously held triggers starts cleanly with the trigger data/schema removed.

**Acceptance Scenarios**:

1. **Given** the GUI, **When** the user views the available areas, **Then** there is no Triggers
   tab or trigger configuration anywhere.
2. **Given** the CLI, **When** the user lists commands, **Then** no trigger-related commands are
   present.
3. **Given** an existing installation whose database contains triggers, **When** the upgraded
   daemon starts, **Then** it migrates the store to remove triggers and starts without error.
4. **Given** the documentation, **When** a user reads the feature list, **Then** triggers are no
   longer advertised.

---

### User Story 5 - Real-time GUI updates (Priority: P2)

The GUI reflects state changes as they happen, across every view (tasks, schedule, groups, logs),
without the user pressing a Refresh button. When a task is created, edited, runs, fails, or a log
is produced — by this user, another client, or the scheduler itself — the relevant views update
automatically. The manual "Refresh" controls are removed.

**Why this priority**: A live, self-updating UI is the desired experience and removes a confusing
manual step. The event-streaming foundation already exists; this story makes it comprehensive and
removes the now-redundant manual controls.

**Independent Test**: With the GUI open, perform a mutation from the CLI (e.g. add a task); the
GUI's task list updates on its own within a couple of seconds. Confirm no Refresh button remains
in any view.

**Acceptance Scenarios**:

1. **Given** the GUI is open on the Tasks view, **When** a task is created or modified by any
   client, **Then** the Tasks view updates automatically without user action.
2. **Given** the GUI is open on the Schedule view, **When** an occurrence runs or a task's
   schedule changes, **Then** the Schedule/calendar updates automatically.
3. **Given** the GUI is open, **When** any view is displayed, **Then** no manual "Refresh"
   control is present.
4. **Given** the live update connection drops, **When** it is restored, **Then** the GUI
   re-synchronizes to current state automatically.

---

### User Story 6 - Calendar view under Schedule (Priority: P3)

Under Schedule, the user can toggle between the existing agenda/timeline list and a calendar view
that lays out past and upcoming occurrences on a date grid, matching what the README advertises.
The toggle is a view option; switching does not lose the selected time window.

**Why this priority**: A nice-to-have that fulfills an advertised capability. The backend
calendar data already exists; this is primarily a presentation addition, so it is lowest risk to
defer.

**Independent Test**: Open Schedule, toggle to the calendar view; occurrences appear on the
correct dates; toggle back to the list and the same data is shown; the calendar updates live as
occurrences change.

**Acceptance Scenarios**:

1. **Given** the Schedule view, **When** the user selects the calendar view option, **Then**
   occurrences are displayed on a date-grid calendar.
2. **Given** the calendar view, **When** the user switches back to the list/agenda view, **Then**
   the same underlying occurrences are shown for the selected window.
3. **Given** the calendar view is open, **When** occurrences change (a run completes, a task is
   added), **Then** the calendar reflects the change automatically.

---

### Edge Cases

- **Rename collisions**: occurrences of "go-scheduler" that are immutable history (existing
  CHANGELOG entries, prior spec files) must be preserved, while live references are updated. The
  scope must clearly distinguish the two.
- **Per-user binary names**: the executable names (`gosched`, `goschedd`, `gosched-gui`) do not
  contain "go-scheduler"; the rename must not gratuitously rename binaries unless the stakeholder
  wants it (see Assumptions).
- **MSI without admin rights**: installing a system service requires elevation. The installer
  must behave predictably (request elevation or fail with a clear message) when run without it.
- **Upgrade over an existing zip install**: a user who previously used the portable zip and now
  runs the MSI should not end up with two competing daemons. The single-instance guard and
  install flow must account for a pre-existing auto-started daemon.
- **Log volume**: the on-disk log file must not grow without bound; a retention/rotation policy is
  required so disk usage stays bounded.
- **Empty Logs**: filtering or "Dismiss All" on an empty list must be a safe no-op.
- **Dismiss All semantics**: RESOLVED (planning) — "Dismiss All" clears the in-GUI list and
  acknowledges the shown alerts only; it does NOT delete the on-disk log file, which remains the
  durable troubleshooting record. So clearing never loses error history.
- **Triggers data referenced elsewhere**: if any task or group references a trigger, removal must
  not orphan or break those records.
- **Calendar with no occurrences**: an empty window must render an empty calendar, not an error.

## Requirements *(mandatory)*

### Functional Requirements

**Rename (US1)**

- **FR-001**: The project MUST be renamed from "go-scheduler" to "go-schedule" across all live
  references, including the module/import path, build and release configuration, documentation,
  and user-facing strings (GUI title, CLI help/branding).
- **FR-002**: The rename MUST preserve immutable historical records (existing changelog entries
  and prior spec documents) rather than rewriting history; the project MUST define which files are
  treated as immutable history.
- **FR-003**: After the rename, the daemon, CLI, and GUI MUST build and all existing tests MUST
  pass.

**Windows MSI install (US2)**

- **FR-004**: The Windows distribution MUST be a single `.msi` installer; the portable zip and any
  "run the exe from a folder" instructions MUST be removed from the release and documentation.
- **FR-005**: The installer MUST place application files in a standard system install location
  (not a user Downloads/temp folder).
- **FR-006**: The installer MUST register the daemon as a background service configured to start
  automatically, including across reboots without an interactive login.
- **FR-007**: The installer MUST create Start Menu entry points for launching the GUI, and the GUI
  MUST launch without a visible console window.
- **FR-008**: The installer MUST require/obtain the elevation needed to install the service and
  MUST fail with a clear, actionable message if elevation is unavailable.
- **FR-009**: Uninstalling via the standard Windows mechanism MUST remove binaries, the service
  registration, and shortcuts; the product MUST document what happens to user data (tasks
  database, logs) on uninstall.
- **FR-010**: The installed GUI MUST detect and reuse the already-running service rather than
  starting a second daemon instance.

**Logs (US3)**

- **FR-011**: The GUI MUST replace the "Alerts" view with a "Logs" view that aggregates both the
  daemon's structured log records and scheduler events (failures, overlaps, missed runs) into a
  single chronological list, newest first.
- **FR-012**: Each log entry MUST carry at minimum a severity (Info, Warning, Error), a timestamp,
  a source, a short message, and detailed cause/context.
- **FR-013**: The Logs view MUST provide filters to show only entries of a chosen severity/type
  (at minimum Errors and Warnings) and to clear the filter.
- **FR-014**: The Logs view MUST let the user open a selected entry to see its full message and
  cause/context for troubleshooting.
- **FR-015**: The Logs view MUST provide a "Dismiss All" action that clears the displayed logs.
- **FR-016**: In an installed environment, log records MUST be persisted to a documented log file
  on disk that survives daemon restarts.
- **FR-017**: On-disk logs MUST be bounded by a retention/rotation policy so disk usage does not
  grow without limit.
- **FR-018**: New log entries MUST appear in the open Logs view in real time without a manual
  refresh.

**Triggers removal (US4)**

- **FR-019**: The Triggers feature MUST be removed from every layer: GUI, CLI, API, scheduling
  engine, store schema, and domain model.
- **FR-020**: A forward store migration MUST remove existing trigger data and schema so upgraded
  daemons start cleanly; the migration MUST not corrupt or block startup for databases that never
  held triggers.
- **FR-021**: Documentation and advertised feature lists MUST no longer mention triggers.

**Real-time GUI (US5)**

- **FR-022**: Every GUI view (tasks, schedule, groups, logs) MUST update automatically in
  response to relevant state changes originating from any client or the scheduler itself.
- **FR-023**: Manual "Refresh" controls MUST be removed from all GUI views.
- **FR-024**: When the live update connection drops and is restored, the GUI MUST automatically
  re-synchronize to current state.

**Calendar view (US6)**

- **FR-025**: The Schedule area MUST offer a calendar view as a toggleable view option alongside
  the existing agenda/list view.
- **FR-026**: The calendar view MUST display past and upcoming occurrences positioned on a date
  grid for the selected time window.
- **FR-027**: Toggling between calendar and list views MUST preserve the selected time window and
  underlying occurrence data.
- **FR-028**: The calendar view MUST update automatically as occurrences change (consistent with
  FR-022).

### Key Entities *(include if feature involves data)*

- **Log Entry**: A unified record shown in the Logs view. Attributes: severity (Info/Warning/
  Error), timestamp, source/component, task or job correlation identifier (when applicable), short
  message, and detailed cause/context. Originates from daemon structured logs and from scheduler
  events; persisted to an on-disk log file in installed environments.
- **Occurrence**: An item shown in the Schedule/calendar — a past run (with outcome) or an
  upcoming scheduled run, with a time and the owning task. (Existing concept; reused by the
  calendar view.)
- **Trigger** *(being removed)*: The existing event-trigger entity and its persisted table are
  deleted as part of this feature.
- **Windows Install Package**: The `.msi` artifact and its declared components — installed files,
  service registration, shortcuts, and uninstall behavior.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A repository-wide search for the live product name returns "go-schedule" only; the
  sole "go-scheduler" matches are files explicitly designated as immutable history.
- **SC-002**: On a clean Windows machine, a non-expert user can go from downloading the `.msi` to
  a running, auto-starting installed product in under 2 minutes and without extracting any archive
  or running a loose executable.
- **SC-003**: After install and a reboot with no user logged in, scheduled tasks still fire.
- **SC-004**: Uninstalling removes 100% of installed binaries, the service registration, and
  shortcuts (verified by their absence afterward).
- **SC-005**: In the Logs view, a user can isolate all error entries in a single action, and open
  any entry to read its full cause — verified by a troubleshooting walkthrough on a seeded failure.
- **SC-006**: Log records produced by the daemon are present in the on-disk log file and persist
  across a daemon restart; the file size stays within the configured retention bound under
  sustained logging.
- **SC-007**: No trigger functionality is reachable from the GUI, CLI, or API, and an upgraded
  daemon starts cleanly against a database that previously contained triggers.
- **SC-008**: With the GUI open, a change made from another client is reflected in the relevant
  view within 2 seconds without any manual refresh, and no Refresh control exists in any view.
- **SC-009**: A user can toggle the Schedule area to a calendar view and see occurrences on the
  correct dates, then toggle back without losing the selected window.

## Assumptions

- **Module/import path**: The Go module path is updated from the `.../go-scheduler` form to the
  `.../go-schedule` form. The owning organization segment of the path is unchanged unless the
  stakeholder specifies otherwise.
- **Binary names unchanged**: The executables `gosched`, `goschedd`, and `gosched-gui` keep their
  names (they do not contain "go-scheduler"); only the project/repo name and paths change. The
  GitHub remote/repository rename itself is an external action performed by the maintainer; this
  feature updates in-repo references.
- **Immutable history**: Existing `CHANGELOG.md` historical entries and prior `specs/00*`
  documents are treated as immutable history and are not rewritten by the rename.
- **MSI is per-machine**: The installer performs a per-machine (all-users) install requiring
  elevation, consistent with installing a system service. A per-user install mode is out of scope.
- **MSI tooling/signing**: Producing the `.msi` is handled by the release pipeline. Code-signing
  the installer is desirable but treated as a separate concern; an unsigned MSI is acceptable for
  this feature if signing infrastructure is unavailable, with SmartScreen guidance documented.
- **Log persistence location**: On Windows the log file lives under a standard per-machine
  application-data location established by the installer; on Linux/macOS it follows the existing
  service/data-directory conventions.
- **Dismiss All scope**: RESOLVED — "Dismiss All" clears the in-GUI log list (and acknowledges the
  underlying alert events) but does NOT delete the on-disk log file, which remains the durable
  troubleshooting record.
- **Calendar data source**: The calendar view reuses the existing calendar/occurrence backend API
  already present in the codebase; no new backend query is required beyond what exists.
- **Real-time foundation**: The existing event-stream (SSE) and view-model change mechanism are
  the basis for comprehensive real-time updates; this feature extends coverage and removes manual
  refresh rather than introducing a new transport.
- **Cross-platform install scope**: This feature's installer work targets Windows (the `.msi`).
  Linux/macOS service installation continues to use existing mechanisms and is out of scope here
  except where the rename touches them.
