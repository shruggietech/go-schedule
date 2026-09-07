# Feature Specification: Desktop Settings, Information, and Recovery

**Feature Branch**: `codex/065-options-info`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Desktop preference migration, Settings and Connections workspaces, authoritative storage inventory, bounded native actions, focused Go race and frontend suites, Chromium accessibility and zoom coverage, native Windows Wails build, and canonical eight-gate verification passed 2026-09-07 on review branch `codex/065-options-info` for [#156](https://github.com/shruggietech/go-schedule/issues/156). Merge-time branch cleanup remains assigned to [#195](https://github.com/shruggietech/go-schedule/issues/195).

**Input**: GitHub issue [#156](https://github.com/shruggietech/go-schedule/issues/156), rebuild Options, Info, diagnostics, and desktop preferences. Repository-governance issue [#195](https://github.com/shruggietech/go-schedule/issues/195) is included only for merge-time verification of GitHub's already-enabled automatic branch deletion.

## User Scenarios & Testing

### User Story 1 - Upgrade with predictable preferences (Priority: P1)

An existing desktop user opens the Wails application and retains a valid light, dark, or system appearance choice without depending on the retired Fyne framework. A new user receives the platform-following system appearance. The Settings page explains what was migrated or retired and allows the current appearance to be changed or restored.

**Why this priority**: An unexplained preference reset is an immediate upgrade regression, while carrying obsolete framework settings forward would preserve controls that no longer affect the redesigned application.

**Independent Test**: Start with valid, invalid, malformed, absent, and previously migrated legacy preference fixtures, then verify deterministic startup, one-time migration, durable appearance changes, retirement disclosure, and restore behavior.

**Acceptance Scenarios**:

1. **Given** a valid legacy `appearance.mode` and no Wails preference file, **When** the desktop starts, **Then** that mode is migrated once into the versioned Wails preference store and applied before the initial workspace is presented.
2. **Given** no usable legacy appearance, **When** the desktop starts, **Then** system appearance is used and the application remains fully usable with a concise transition explanation.
3. **Given** an established Wails preference file, **When** a legacy file later changes, **Then** the Wails preference remains authoritative and is not overwritten.
4. **Given** a user changes appearance or restores defaults, **When** the application restarts, **Then** the saved selection or restored system default is applied predictably.

---

### User Story 2 - Understand storage and product information (Priority: P1)

A user opens Settings and can identify application-owned and externally owned storage, whether each location exists, how removal affects it, and copy an exact available path. The same workspace provides trustworthy product, version, source, and documentation information using platform-native link behavior.

**Why this priority**: Storage visibility and accurate removal claims are necessary for support, backup, and privacy decisions, while product links and version details provide the context required to act on that information.

**Independent Test**: Load present, absent, unavailable, daemon-authoritative, and external storage fixtures, copy eligible paths, activate each fixed product link, and confirm the application never invents a path or accepts an arbitrary URL.

**Acceptance Scenarios**:

1. **Given** local and daemon-authoritative storage records, **When** Settings loads, **Then** each row names its owner, scope, existence, removal behavior, and exact path when known.
2. **Given** a row with an available path, **When** Copy is activated, **Then** the exact current path is written through the native clipboard and success is announced; unavailable rows do not offer a misleading copy action.
3. **Given** product links, **When** a link is activated, **Then** the operating system opens only the application-defined HTTPS destination in the default browser.
4. **Given** the daemon is unavailable, **When** Settings loads, **Then** local settings and product information remain usable while daemon-owned paths are quietly marked unavailable with a connection recovery route.

---

### User Story 3 - Recover the local connection without interruption (Priority: P2)

A user opens Connections or follows a recovery action and can understand the current daemon diagnosis, retry safely, and continue using local Settings without repeated modal interruptions.

**Why this priority**: Connection diagnosis already exists in the shell, but a dedicated recovery workspace is needed to make the information actionable and persistent without disrupting unrelated local tasks.

**Independent Test**: Present connected, unavailable, timed out, access denied, and incompatible snapshots, retry each state, and verify status and recovery guidance remain inline and keyboard accessible.

**Acceptance Scenarios**:

1. **Given** any supported connection diagnosis, **When** Connections opens, **Then** it displays explicit status, endpoint or platform context when safe, and diagnosis-specific recovery guidance without automatically opening a modal.
2. **Given** a disconnected state, **When** Retry is activated repeatedly, **Then** only one retry is pending, the result is announced, and the current route and focus remain stable.
3. **Given** a connection status change, **When** the shell receives it, **Then** Connections and Settings update inline without stealing focus or repeating a dialog.

### Edge Cases

- The legacy preference file is missing, malformed, unreadable, contains an unknown mode, or contains valid mode data alongside retired fields.
- The new preference directory does not exist, a write fails, or a temporary file cannot be renamed.
- Two preference mutations arrive close together.
- A configured appearance value has the correct JSON type but unsupported casing or whitespace.
- The executable path or user configuration directory cannot be resolved.
- A storage path is absent, permission-restricted, relative, daemon-unavailable, or externally owned.
- The daemon reports paths that differ from platform defaults.
- Clipboard or system-browser integration fails.
- A frontend caller supplies an unknown storage identifier or product-link key.
- The connection changes while Settings or Connections is open.
- The interface is used entirely by keyboard at 80, 100, 150, and 200 percent browser zoom.

## Requirements

### Functional Requirements

- **FR-001**: The desktop MUST provide complete Settings and Connections routes in place of their placeholders.
- **FR-002**: Settings MUST consolidate the useful outcomes of the legacy Options and Info windows into appearance, preference transition, storage, and product-information sections.
- **FR-003**: The Wails preference store MUST be a versioned per-user file independent of Fyne and MUST persist only settings that affect the redesigned desktop.
- **FR-004**: A valid legacy `appearance.mode` value of `system`, `light`, or `dark` MUST migrate exactly once when no Wails preference file exists.
- **FR-005**: New installations and unusable legacy appearance values MUST default to `system`; malformed or unreadable legacy files MUST NOT prevent startup.
- **FR-006**: Once the Wails preference file exists, it MUST be authoritative and later legacy changes MUST NOT overwrite it.
- **FR-007**: Legacy font selection and scroll sensitivity MUST be retired rather than migrated because the Wails design system and browser scrolling own those behaviors; Settings and documentation MUST disclose this decision.
- **FR-008**: Users MUST be able to select `system`, `light`, or `dark` appearance, persist it immediately, and restore the desktop preference defaults.
- **FR-009**: Preference writes MUST replace the complete file atomically where the platform permits, avoid secrets, reject unsupported values, and surface a non-destructive actionable failure.
- **FR-010**: Settings MUST list known machine data, task database, configuration, logs, runtime state, desktop application data, desktop preferences, executable directory, installed documentation when present, and Windows maintenance evidence when applicable.
- **FR-011**: Every storage record MUST state ownership, scope, existence or unavailable status, normal removal behavior, and explicit-wipe behavior without claiming that externally owned data is removed.
- **FR-012**: Daemon-reported effective paths MUST remain authoritative for daemon storage; when the daemon is unavailable, the desktop MUST mark those records unavailable instead of guessing.
- **FR-013**: Copy actions MUST operate on a backend-resolved storage identifier, copy only the exact current available path through native clipboard integration, and announce success or failure.
- **FR-014**: Product information MUST include application name, version, publisher, source link, and documentation link using the repository's existing build metadata.
- **FR-015**: External-link actions MUST accept only fixed backend-defined keys and open their corresponding HTTPS destinations through the operating system default browser.
- **FR-016**: Connections MUST display connected, unavailable, timed out, access denied, incompatible, and retrying states with diagnosis-specific inline recovery guidance.
- **FR-017**: Connection recovery MUST reuse the existing connection manager, suppress duplicate pending retries, preserve route and focus, and avoid automatic or repeated modal dialogs.
- **FR-018**: Local preferences and product information MUST remain usable while the daemon is unavailable; daemon-backed storage rows and mutations MUST accurately reflect their unavailable state.
- **FR-019**: Settings and Connections controls, status, disclosures, tables or lists, copy actions, links, and recovery actions MUST be keyboard accessible, screen-reader meaningful, understandable without color, and readable at 80 through 200 percent zoom.
- **FR-020**: Automated Go, React, accessibility, migration, atomic-write failure, authorization-boundary, offline, retry, zoom, native-build, and canonical repository verification MUST pass.
- **FR-021**: Issue [#156](https://github.com/shruggietech/go-schedule/issues/156) MUST remain traceable through the specification, tasks, change log, pull request, and verification record.
- **FR-022**: Issue [#195](https://github.com/shruggietech/go-schedule/issues/195) MUST remain open until S065 is merged and its same-repository remote topic branch is confirmed automatically deleted; S065 MUST add no branch-cleanup workflow or credential.

### Key Entities

- **Desktop Preferences**: Versioned durable settings, currently containing only the appearance mode and migration disposition.
- **Preference Transition**: The recorded result of the one-time legacy inspection, including migrated, no legacy value, invalid legacy value, or unreadable legacy source.
- **Settings Workspace**: One response containing preferences, transition explanation, storage records, product information, and daemon availability.
- **Storage Record**: A stable identifier, label, ownership, scope, existence state, exact path when authoritative, and removal descriptions.
- **Product Link**: A stable backend-defined key and human label mapped to a fixed HTTPS destination.
- **Connection Snapshot**: Existing connection state, diagnosis, safe endpoint context, guidance, and retry availability.
- **Native Action Result**: Accepted, rejected, or unavailable outcome with a concise user-safe message.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Automated fixtures for valid, absent, malformed, unreadable, invalid, and already-migrated legacy preferences all produce the specified deterministic appearance and transition result.
- **SC-002**: A keyboard-only user can change and restore appearance, copy every eligible storage path, open each product link, inspect all connection states, and retry the daemon connection.
- **SC-003**: Every storage fixture states ownership and both removal behaviors, and no daemon or external path is guessed or claimed as application-owned.
- **SC-004**: Tests prove that unknown storage and link identifiers cannot copy arbitrary text or open arbitrary destinations.
- **SC-005**: Settings and Connections remain usable and readable at 80, 100, 150, and 200 percent zoom with no clipped primary action or horizontal page overflow at the supported desktop viewport.
- **SC-006**: Focused tests, native Wails build, and `sh scripts/verify.sh all` pass before publication.

## Clarifications

### Session 2026-09-07

- Q: Which legacy preferences survive the framework transition? A: Only valid system, light, or dark appearance intent migrates. Fyne font choice and scroll sensitivity are retired because they no longer govern the Wails experience.
- Q: What is the default for new or damaged preference state? A: System appearance. This intentionally differs from the legacy dark default so new installs follow the operating system, while valid existing intent is preserved.
- Q: Are Options and Info separate destinations? A: No. Settings combines appearance, transition disclosure, storage inventory, and product information. Connections remains separate because recovery is an operational task.
- Q: Which paths may the desktop infer? A: It may resolve its own preferences, application-data, executable, installed-documentation, and applicable Windows maintenance locations. Daemon-owned effective paths come only from runtime metadata and are unavailable when that source cannot be reached.
- Q: What can the frontend ask the operating system to open or copy? A: Only stable identifiers. The backend resolves the current storage path and fixed HTTPS product destination before invoking the native integration.
- Q: How is connection failure presented? A: Inline through the shell and Connections route, with one explicit retry action. No automatic repeated modal is introduced.

## Assumptions

- The existing Wails shell, appearance tokens, connection manager, runtime-info API, build metadata, and legacy storage inventory behavior provide the implementation foundation.
- The versioned Wails preference file lives beneath the current user's configuration directory and contains no credential or daemon configuration.
- Existing daemon storage paths can vary from platform defaults and remain authoritative whenever runtime metadata is available.
- Packaging and full legacy Fyne code removal remain assigned to cutover issue #157; this slice documents transition behavior but does not remove the legacy desktop implementation.
- GitHub automatic branch deletion was enabled before S065. Its real merge behavior can only be verified after the maintainer merges this pull request.
- This slice fully resolves issue #156. Issue #195 is resolved only after the post-merge branch-deletion observation, and parent issue #147 remains open for cutover.
