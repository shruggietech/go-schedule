# Feature Specification: Production Wails Shell and Local Connection

**Feature Branch**: `codex/061-wails-shell-connection`

**Created**: 2026-09-07

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Focused Go race, frontend, Chromium, native Windows, dependency, brand, lifecycle, and canonical eight-gate verification passed 2026-09-07 on review branch `codex/061-wails-shell-connection`; hosted three-platform evidence runs on the pull request

**Input**: GitHub issues #151 and #152, following the approved S060 foundation and Calm Operations experience direction.

## Clarifications

### Session 2026-09-07

- Q: Does this slice replace the shipping Fyne application? → A: No. It creates the production Wails foundation, but packaging and Fyne retirement remain in #157.
- Q: How much scheduler functionality belongs in the shell slice? → A: Only connection identity, health, capabilities, lifecycle, events, recovery, and representative shell states; feature workflows remain in #153 through #156.
- Q: How should transient local connection loss recover? → A: Retry automatically with a bounded backoff, retain a manual retry action, and keep access, compatibility, timeout, and unavailable failures distinct.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Open the local control center safely (Priority: P1)

An operator opens the new desktop foundation and immediately sees a usable application frame, the active target named This computer, its platform and version context, and a clear connection state without registration, login, or a network listener.

**Why this priority**: Every later desktop workflow depends on a trustworthy shell and local connection, and the existing offline-first security boundary must remain intact.

**Independent Test**: Start the desktop foundation with connected and unavailable fake local services, then with the real protected local service, and verify that target identity and honest state appear without exposing transport details or credentials.

**Acceptance Scenarios**:

1. **Given** the local daemon is available, **When** the operator opens the application, **Then** the shell identifies This computer, reports connected health and supported capabilities, and makes the navigation frame usable.
2. **Given** the local daemon is not running, **When** the operator opens the application, **Then** the shell remains usable, reports an unavailable state with a recovery action, and does not show invented daemon data.
3. **Given** the application is installed on a supported platform, **When** it starts with no configuration changes, **Then** it uses only the protected local operating-system transport and creates no network listener.

---

### User Story 2 - Understand and recover connection failures (Priority: P1)

An operator can tell whether a local connection is unavailable, denied, incompatible, timed out, degraded, or recovering, and can retry without restarting the desktop application.

**Why this priority**: A vague disconnected state hides materially different security, compatibility, and service failures and makes recovery harder.

**Independent Test**: Drive deterministic connection fakes through every connection state, automatic retry, manual retry, event interruption, recovery, and shutdown while checking state transitions, user guidance, and cancellation.

**Acceptance Scenarios**:

1. **Given** local access is denied, **When** the initial request fails, **Then** the application names the access problem and offers platform-appropriate guidance without exposing sensitive endpoint details.
2. **Given** the daemon version or capability contract is incompatible, **When** negotiation completes, **Then** the application reports incompatibility distinctly and does not treat the connection as healthy.
3. **Given** a connected event stream is interrupted, **When** ordinary requests still work, **Then** the application becomes degraded, preserves last-known nonsensitive identity, and attempts bounded recovery.
4. **Given** the daemon becomes available after a failure, **When** automatic or manual retry succeeds, **Then** the application returns to connected state and resumes one coherent event lifecycle.

---

### User Story 3 - Use a consistent accessible application frame (Priority: P2)

An operator uses primary navigation, appearance preferences, notices, forms, tables, disclosures, dialogs, and status feedback consistently across ordinary and compact windows with keyboard or assistive technology.

**Why this priority**: Shared primitives must be complete before feature teams build on them, otherwise later screens will create incompatible interaction and accessibility behavior.

**Independent Test**: Exercise the component catalog and shell at supported window sizes, appearance modes, zoom levels, keyboard paths, reduced-motion settings, and representative content extremes with automated accessibility checks.

**Acceptance Scenarios**:

1. **Given** any supported appearance or connection state, **When** the operator moves through the shell by keyboard, **Then** focus remains visible, navigation order is logical, dialogs return focus to their exact invoker, and no keyboard trap occurs.
2. **Given** a compact window or 200 percent zoom, **When** the shell contains long translated-style labels and status text, **Then** target identity, page identity, primary actions, and recovery guidance remain reachable without horizontal page overflow.
3. **Given** reduced motion or system appearance preferences, **When** the shell opens or changes state, **Then** the interface respects those preferences and never relies on motion or color alone to communicate status.
4. **Given** a later feature screen uses the production component catalog, **When** it needs a table, form, help text, validation, empty state, loading state, error, confirmation, disclosure, notification, or overlay, **Then** a documented reusable primitive and contract are available.

### Edge Cases

- The service can disappear between health negotiation and a later request; state must move from connected to unavailable or degraded without retaining a false healthy claim.
- Manual retry can overlap an automatic retry; only one connection attempt and one event stream may own the active lifecycle.
- An old daemon can answer health but omit required capabilities; it must be incompatible rather than connected.
- Shutdown can begin while a request, retry delay, event read, or native action is active; all owned work must terminate within the shutdown budget.
- An event can arrive while a snapshot refresh is in flight; state updates must not regress to an older connection generation.
- Local configuration or endpoint errors can include sensitive paths or operating-system details; frontend messages must use bounded safe categories and actions.
- Long target names, versions, localized-style labels, and status messages must not obscure the active target or page.
- Browser-only component tests have no native bridge; deterministic adapters must represent native outcomes without claiming a real daemon connection.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product MUST contain one production-intent desktop foundation that is separate from the disposable S060 proof and excluded from current release artifacts until #157.
- **FR-002**: The shell MUST provide primary navigation, page framing, persistent active-target identity, appearance selection, command surfaces, overlays, notifications, and orderly application exit.
- **FR-003**: Shared interface primitives MUST cover buttons, links, status indicators, notices, tables, forms, field help, validation, empty and loading states, errors, confirmations, disclosures, dialogs, and transient notifications.
- **FR-004**: Shared primitives MUST define accessible names, focus behavior, keyboard behavior, disabled and busy behavior, destructive emphasis, error association, and non-color-only state communication.
- **FR-005**: The shell MUST support light, dark, and system appearance, visible focus, reduced motion, 200 percent zoom, and useful layouts from 900 by 650 through 1440 by 900 logical pixels.
- **FR-006**: Feature-facing frontend code MUST consume one connection-facing application contract and MUST NOT receive local endpoint names, credentials, transport objects, or raw backend errors.
- **FR-007**: The connection contract MUST represent stable target identity, display name, execution platform, daemon version, capabilities, permissions, connection generation, transition revision, health, safe guidance, and last successful contact.
- **FR-008**: The default target MUST be This computer over the existing protected local transport with no account, login, registration, remote URL, or network listener.
- **FR-009**: Connection states MUST distinguish connecting, connected, degraded, recovering, unavailable, access denied, incompatible, and timed out.
- **FR-010**: Connection transitions MUST be deterministic and reject stale responses or events from an earlier connection generation or transition revision.
- **FR-011**: Initial and manual connection attempts MUST time out within two seconds, while transient retry MUST use a bounded progression of 250 milliseconds, one second, and five seconds with no more than one active retry loop.
- **FR-012**: Manual retry MUST cancel any pending retry delay, start one fresh generation, and remain available from every recoverable failure state.
- **FR-013**: Ordinary requests and live events MUST share the connection lifecycle, cancellation root, safe error mapping, and shutdown behavior while allowing an event-only failure to report degraded service.
- **FR-014**: The shell MUST publish sanitized connection-state and daemon-domain notifications through one documented frontend event channel with stable event identity and generation.
- **FR-015**: Access denial, incompatible daemon, timeout, unavailable service, degraded events, and internal failure MUST each map to bounded safe language and a relevant recovery action without exposing sensitive values.
- **FR-016**: Closing the application MUST cancel requests, retry waits, event streams, and native actions, prevent new work, and finish owned goroutines within two seconds.
- **FR-017**: Deterministic fakes MUST support every connection state, delayed and stale results, event interruption, recovery, and shutdown without a native daemon or wall-clock sleeps.
- **FR-018**: The production foundation MUST use only local repository assets and MUST make no telemetry, analytics, CDN, font, image, script, or style request.
- **FR-019**: S061 MUST preserve the shipping Fyne application, current installer inputs, CLI, daemon, storage schema, scheduler behavior, and release identity unchanged.
- **FR-020**: The component and connection contracts MUST be documented so #153 through #156 can add feature workflows without bypassing the shared system or depending on transport details.
- **FR-021**: Automated verification MUST cover component behavior and accessibility, connection state and concurrency, real protected-local-client integration, offline assets, production builds on Windows, macOS, and Linux, and the canonical repository gates.

### Key Entities

- **Desktop target**: Stable identity for an execution host, including display name, platform, daemon metadata, capabilities, permissions, and current connection state.
- **Connection snapshot**: Immutable generation-and-revision-stamped view of target health, safe guidance, last successful contact, and supported capabilities.
- **Connection event**: Sanitized generation-stamped notification about connection or daemon state.
- **Connection generation**: Monotonic attempt identity used to reject stale results and prevent overlapping lifecycle ownership.
- **Shell route**: Stable destination in the production navigation frame, including availability and accessible page identity.
- **Appearance preference**: System, light, or dark selection applied through shared tokens.
- **Interface primitive**: Reusable presentation and interaction contract for later workflow screens.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: On all three supported desktop platforms, a clean production build opens and closes the shell while retaining visible target and page identity at every supported window size.
- **SC-002**: Connected local-service identity and capabilities appear within two seconds under nominal local conditions; every unavailable attempt produces actionable state within the same bound.
- **SC-003**: All eight connection states, every allowed transition, stale-generation rejection, event-only degradation, manual recovery, and shutdown cancellation pass deterministic automated tests with zero race findings.
- **SC-004**: One hundred consecutive connect, fail, recover, and shutdown cycles leave zero owned retry or event goroutines after each cycle.
- **SC-005**: Automated component checks report zero serious or critical accessibility violations across all appearance and connection states at ordinary, compact, and 200 percent zoom layouts.
- **SC-006**: Every required shared primitive has a documented contract and at least one behavior or accessibility test before any migrated feature screen uses it.
- **SC-007**: Browser and packaged build inspection find zero runtime requests for remote assets, telemetry, or analytics.
- **SC-008**: Windows, macOS, and Linux hosted builds pass the production foundation, while every existing canonical repository gate remains green and current shipping artifacts remain unchanged.

## Assumptions

- S060 already approved Wails v2.14.0, React 19.1.0, TypeScript 5.6.3, Vite 7.3.6, Node 24, the local asset set, state vocabulary, and Calm Operations direction.
- The daemon health response is the initial negotiation surface; S061 may derive a bounded local capability manifest from the current version and available client methods without changing the public daemon API.
- This computer is the only selectable target in v1.2 until remote access work explicitly adds enrolled targets.
- Existing Fyne preferences are not migrated here; their disposition belongs to #156.

## Dependencies

- Parent coordinator: #147.
- Resolves #151 and #152 together.
- Depends on completed #149 and #150 and the S060 architecture and experience contracts.
- Blocks #153, #154, #155, #156, and the Wails portion of #192.

## Out of Scope

- Migrating Tasks, Groups, Chains, Triggers, Watchers, Schedule, Activity, Options, Info, or preference data.
- Adding remote targets, remote credentials, network listeners, authentication, or pairing.
- Replacing current Fyne release inputs, changing installers, deleting Fyne, or publishing a release.
- Changing daemon APIs, scheduler behavior, persisted data, CLI behavior, or public configuration.
- Claiming attended native screen-reader or installer qualification reserved for #157.
