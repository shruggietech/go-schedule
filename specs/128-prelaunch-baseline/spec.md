# Feature Specification: S128 Prelaunch Baseline Simplification

**Feature Branch**: `codex/128-prelaunch-baseline`
**Created**: 2026-10-31
**Status**: Implemented
**Delivery**: Review branch codex/128-prelaunch-baseline; source and local verification complete, official PR and owner merge pending.
**Native scope**: [#408](https://github.com/shruggietech/go-schedule/issues/408), bounded verification corrections [#409](https://github.com/shruggietech/go-schedule/issues/409) and [#410](https://github.com/shruggietech/go-schedule/issues/410), milestone 18 (Maintenance - Prelaunch baseline).

**Input**: The owner states that go-schedule has no users and migrations of any kind are irrelevant, approves the prelaunch simplification proposal, and authorizes the complete Spec-Kit/autopilot sequence, push and one official PR, in-scope CI/review corrections and at most two Codex review requests. Stop before merge.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Start and reopen a current database (Priority: P1)

A developer starts the current daemon with an empty data directory and gets the complete current database directly, without replaying historical upgrades. Restarting the same current database preserves its tasks, run history, credentials, identities and audit data.

**Why this priority**: Fresh installations currently execute 24 historical schema stages despite having no installed user base to migrate.
**Independent Test**: Initialize empty memory and file databases, exercise current persistence/constraints, close and reopen the file database, and compare identities and stored state.

**Acceptance Scenarios**:

1. **Given** empty storage, **When** the current daemon initializes, **Then** one atomic current-schema initialization supplies all supported tables, indexes and defaults without intermediate historical schemas.
2. **Given** a current database, **When** it is reopened, **Then** current data and daemon/client identity persist without conversion or schema replay.
3. **Given** an unsupported or foreign development database, **When** it is opened, **Then** startup returns an actionable explicit-reset instruction without silently converting, replacing or deleting stored state.
4. **Given** an initialization failure, **When** the transaction ends, **Then** no partial application schema or successful version marker remains.

### User Story 2 - Use current desktop preferences (Priority: P1)

A developer uses the Wails preferences directly. Old Fyne files are neither imported nor reported through settings, CLI, API bindings or desktop migration messages.

**Why this priority**: Legacy conversion and its UI metadata are unnecessary product and maintenance complexity.
**Independent Test**: Exercise default/current preference load, change and reload appearance/popups, and place an unrelated legacy file beside the current preference root to show it has no effect.

**Acceptance Scenarios**:

1. **Given** no current preferences, **When** desktop or headless preference commands run, **Then** current defaults are used regardless of obsolete Fyne files.
2. **Given** valid current preferences, **When** settings are changed and reloaded, **Then** current appearance and popup behavior persist without migration transition metadata.
3. **Given** malformed current preferences, **When** they are loaded, **Then** current error handling stays explicit without converting or silently overwriting the file.
4. **Given** the settings page or current response models, **When** preferences are inspected, **Then** they contain no obsolete migration-status UI or legacy-import path.

### User Story 3 - Maintain a truthful prelaunch baseline (Priority: P2)

A maintainer follows current documentation and engineering policy that match the owner's no-migration requirement. Verification concentrates on fresh install and current-state behavior, rather than historical upgrade matrices.

**Why this priority**: Legacy promises and mandatory migration gates would otherwise reintroduce the removed machinery.
**Independent Test**: Trace current constitution/autopilot/install/QA guidance and active test entry points against the same prelaunch contract, without rewriting historical release evidence.

**Acceptance Scenarios**:

1. **Given** current project policy, **When** new work is planned, **Then** historical migration compatibility is not a required outcome or test gate.
2. **Given** obsolete developer state, **When** reset guidance is followed intentionally, **Then** the developer selects fresh state explicitly; production startup never wipes it automatically.
3. **Given** current installer testing, **When** fresh install, repair, maintenance and uninstall are exercised, **Then** useful current behavior, access controls, ownership limits and preserve/wipe choices retain their coverage while historical migration-specific matrices are retired.
4. **Given** archived specifications, release notes, comparisons and frozen QA inputs, **When** S128 is delivered, **Then** historical artifacts remain attributable and immutable.

### Edge Cases

- An empty database file, memory database, foreign SQLite database, old schema marker, newer schema marker or partially initialized schema must have an explicit current disposition.
- Existing current-schema development state remains useful; no need to renumber it just to remove history.
- Old preference metadata may exist in a current file; it has no current migration behavior or UI authority.
- Normal installer mechanics may replace an already installed development build; this is distinct from supporting historical migration matrices.
- Retired tests may contain current behavior assertions that must be retained in fresh/current-state tests.
- Current authentication, identity, SQL cancellation, durable audit, rollback and uninstall ownership behavior remain required despite having no users.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Empty application and helper-metrics databases MUST initialize their complete current schema directly and atomically, without a historical migration registry or ordered replay.
- **FR-002**: Current-version file databases MUST reopen with unchanged current data, identity and relational constraints.
- **FR-003**: Unsupported/foreign development databases MUST receive actionable reset guidance without automatic conversion or deletion.
- **FR-004**: Current preference defaults, validation, persistence and headless operations MUST remain; legacy Fyne import behavior and its transition metadata/path/UI MUST be removed.
- **FR-005**: The current constitution, autopilot procedure and active documentation MUST reflect the owner-approved prelaunch policy rather than mandate historical migrations or version-upgrade promises.
- **FR-006**: Historical migration/upgrade test matrices MUST be retired or replaced with meaningful fresh/current-state coverage; security, ownership, cancellation, persistence, transactional correctness, repair and uninstall assertions MUST remain.
- **FR-007**: Reset instructions MUST use an explicit developer action and accurately distinguish fresh state from preserving a current database.
- **FR-008**: Frozen evidence, historical release artifacts, comparison bindings and prior outcome attribution MUST remain unchanged.
- **FR-009**: SQLite 1.60.1/libc 1.77.1 and unrelated runtime features MUST remain; cluster work remains deferred and #355/#356/#400/#404 retain their separate unresolved outcomes.
- **FR-010**: S128 specification, branch, feature selector, native planning record, source scope and PR MUST agree; the approved governance-policy change MUST have a dated decision and synchronized guidance.

### Key Entities

- **Current schema baseline**: complete current tables/indexes/defaults and a current identity marker, without an ordered conversion history.
- **Current preferences**: supported appearance/popup values and persistence version, without migration-transition state.
- **Development reset**: explicit fresh-state selection, with existing files preserved unless separately intentionally removed.
- **Owner policy decision**: no users, no migration obligation, bounded S128 source authorization and separately retained delivery/review boundaries.

## Success Criteria *(mandatory)*

- **SC-001**: Empty memory/file databases expose the same required current capabilities without executing any historical schema stage.
- **SC-002**: Current-state reopen and persistence preserve all tested task/run/identity/audit/credential outcomes; initialization failure leaves no partial successful schema.
- **SC-003**: No active preference import path, transition-status contract or migration UI remains.
- **SC-004**: No active verification entry point requires historical schema/Fyne/installer migration qualification; current security and lifecycle coverage remains effective.
- **SC-005**: Current policy and guidance describe one consistent prelaunch contract, with zero automatic data resets or rewritten historical evidence.

## Assumptions

- The latest explicit owner instruction supersedes the old constitution's mandatory compatibility/migration clauses; S128 includes their formal amendment, not a weakened interpretation.
- Retaining the current schema identifier avoids unnecessary disruption to current developer databases and does not constitute migration support.
- No new reset framework, storage backend, dependency update, cluster feature, tag or release is required.
- The existing publication authorization removes the default pre-push permission pause; merge remains owner-controlled.

## Clarifications

### Session 2026-10-31

- Q: Are historical users or migration guarantees a requirement? A: No. The owner states there are no users and migrations of any kind are irrelevant.
- Q: Can startup silently reset old state? A: No. The approved proposal calls for explicit developer control and a straightforward reset path.
- Q: Should current state still reopen and persist? A: Yes. Current-state persistence, constraints and security remain core product behavior.
- Q: Should immutable historical evidence be rewritten to match the new policy? A: No. Update active source and guidance; retain prior facts and artifacts.
