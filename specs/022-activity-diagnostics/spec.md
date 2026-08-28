# Feature Specification: Activity Diagnostics Clarity

**Feature Branch**: `codex/022-activity-diagnostics`
**Created**: 2026-08-28
**Status**: Draft
**Input**: Improve Activity diagnostics so operators can find the complete daemon log and recognize daemon startup as a completed event.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Find the complete daemon log (Priority: P1)

As an operator reviewing Activity, I can tell that the visible records are a bounded recent view and see the exact configured on-disk log path so that I know where to investigate older records.

**Why this priority**: The Activity view currently looks complete even though it is bounded, which can mislead troubleshooting.

**Independent Test**: Open Activity with default, overridden, empty, and unavailable log-path metadata and verify that the scope and path guidance are truthful without disrupting existing Activity controls.

**Acceptance Scenarios**:

1. **Given** the daemon uses its default log path, **When** Activity loads, **Then** it says the view contains recent records and displays the exact default path reported by the daemon.
2. **Given** the daemon uses an overridden log path containing spaces or platform-specific separators, **When** Activity loads, **Then** it displays that value unchanged.
3. **Given** the daemon has not yet returned log-path metadata, **When** Activity renders, **Then** it states that the path is unavailable until the daemon responds and does not invent a default.
4. **Given** the daemon returns no recent records, **When** Activity loads, **Then** the reported log path remains available.
5. **Given** an operator filters, selects, refreshes, clears the view, or acknowledges alerts, **When** the new guidance is present, **Then** those existing behaviors remain unchanged.

---

### User Story 2 - Recognize completed daemon startup (Priority: P2)

As an operator reading daemon logs, I can distinguish the one-time startup completion event from ongoing listening behavior and see the endpoint, database path, and log path used by the running daemon.

**Why this priority**: Clear event wording makes startup investigations easier without adding a lifecycle or uptime feature.

**Independent Test**: Capture the daemon startup record and verify its discrete completion wording and structured diagnostic fields.

**Acceptance Scenarios**:

1. **Given** the daemon has initialized and is ready to serve, **When** it emits its startup record, **Then** the message describes a completed startup event rather than an ongoing action.
2. **Given** the startup record is emitted, **When** its structured fields are inspected, **Then** the endpoint, database path, and configured log path are present.
3. **Given** Activity renders startup records, **When** an operator views them, **Then** no uptime display or continuously updated lifecycle state is introduced.

### Edge Cases

- The configured log path is relative, inaccessible, or points to a file that has not yet been created.
- The configured path contains whitespace, Unicode characters, or Windows and Unix path separators.
- The daemon returns an empty path because metadata is not available.
- The in-memory recent-log ring is empty or unavailable while configuration metadata is still known.
- A client that ignores newly added response metadata continues consuming recent records.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Activity MUST explicitly identify its records as a bounded recent view.
- **FR-002**: Activity MUST display the exact configured daemon log path reported by the daemon.
- **FR-003**: The displayed path MUST reflect either the resolved default or configured override used by the daemon.
- **FR-004**: Log-path metadata MUST remain available when there are no recent records.
- **FR-005**: Before usable metadata arrives, Activity MUST state that the path is unavailable until the daemon responds and MUST NOT hardcode or infer a platform default.
- **FR-006**: Activity MUST preserve path whitespace, Unicode characters, and platform-specific separators.
- **FR-007**: Activity MUST NOT open, validate, probe, create, or modify the reported file or path.
- **FR-008**: The daemon MUST emit a one-time startup completion record whose wording describes a discrete completed event.
- **FR-009**: The startup completion record MUST include structured endpoint, database-path, and log-path fields.
- **FR-010**: Activity ordering, filtering, details, live updates, refresh, Clear View behavior, alerts, and badge behavior MUST remain compatible.
- **FR-011**: Existing CLI human-readable and JSON log output shapes MUST remain compatible.
- **FR-012**: User and local API documentation MUST explain the recent-view boundary and log-path metadata.
- **FR-013**: The change MUST be recorded in the changelog.
- **FR-014**: The delivery pull request MUST close GitHub issues #27 and #31.
- **FR-015**: File launching, permission checks, file-existence checks, copy controls, uptime display, logging retention or rotation changes, configuration changes, migrations, lifecycle redesign, and network exposure are out of scope.

### Key Entities

- **Recent activity response**: The bounded set of recent structured records plus daemon-reported diagnostic metadata.
- **Configured log path**: The exact path selected by daemon configuration for the rotating on-disk log.
- **Startup completion record**: The single structured event emitted after daemon initialization, containing endpoint and storage locations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Default and overridden configured paths, including paths with spaces and platform-specific separators, display exactly as reported in automated tests.
- **SC-002**: Activity presents one explicit explanation that the visible records are recent and where the full log is stored.
- **SC-003**: Empty-record and pre-response states provide truthful path guidance without an inferred fallback.
- **SC-004**: Automated tests verify discrete startup-completion wording and all three required structured fields.
- **SC-005**: Existing Activity behavior and CLI log-output tests remain green without compatibility regressions.
- **SC-006**: All eight canonical verification gates, encoding checks, and whitespace checks pass.

## Clarifications

### Session 2026-08-28

- No uptime indicator is included; the existing timestamp and discrete startup wording meet the diagnostic need.
- The UI reports the configured path only; it does not open the file or inspect filesystem state.
- Log-path metadata travels with the recent-activity response so each refresh remains internally consistent and requires no extra request.
- Until the daemon responds with usable metadata, Activity uses an explicit unavailable state instead of a hardcoded fallback.

## Assumptions

- The GUI and daemon run in the same local environment, so a daemon-reported path is meaningful to the operator.
- The configured value remains the truthful diagnostic value even if the file is temporarily inaccessible or has not yet been created.
- Existing record timestamps provide enough temporal context for a discrete startup event.

## Out of Scope

- Opening or copying the log path from Activity.
- Testing path permissions or file existence.
- Showing daemon uptime or continuously changing lifecycle status.
- Changing log retention, rotation, configuration semantics, storage, or transport security.
